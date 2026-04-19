package handler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

// videoExts 需要转码的视频后缀（HEVC 通常出现在这些容器格式中）
var videoExts = map[string]bool{
	".mp4":  true,
	".mov":  true,
	".mkv":  true,
	".avi":  true,
	".webm": true,
	".m4v":  true,
}

// IsVideo 判断文件是否为视频
func IsVideo(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return videoExts[ext]
}

// videoInfo 存储视频流的关键信息
type videoInfo struct {
	Codec   string
	Profile string
	PixFmt  string
}

// probeVideoInfo 使用 ffprobe 检测视频的编码格式、profile 和像素格式
func probeVideoInfo(filePath string) (videoInfo, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=codec_name,profile,pix_fmt",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	)
	output, err := cmd.Output()
	if err != nil {
		return videoInfo{}, fmt.Errorf("ffprobe 执行失败: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	info := videoInfo{}
	if len(lines) >= 1 {
		info.Codec = strings.TrimSpace(lines[0])
	}
	if len(lines) >= 2 {
		info.Profile = strings.TrimSpace(lines[1])
	}
	if len(lines) >= 3 {
		info.PixFmt = strings.TrimSpace(lines[2])
	}
	return info, nil
}

// needsTranscode 判断视频是否需要转码
// 满足以下任一条件就需要转码：
// 1. 非 H.264 编码
// 2. H.264 但 profile 不是 Baseline（High、Main 等在部分设备不支持）
// 3. 10bit 色深（yuv420p10le 等）
func needsTranscode(info videoInfo) bool {
	codec := strings.ToLower(info.Codec)
	// 非 H.264，必须转码
	if codec != "h264" && codec != "libx264" && codec != "avc" {
		return true
	}
	// H.264 但 profile 不是 Baseline/Constrained Baseline，需要转码
	profile := strings.ToLower(info.Profile)
	if !strings.Contains(profile, "baseline") {
		return true
	}
	// 10bit 像素格式，需要转码
	if strings.Contains(strings.ToLower(info.PixFmt), "10") {
		return true
	}
	return false
}

// transcodeVideo 使用 ffmpeg 将视频转码为 H.264 + AAC
// 转码后的文件名为原文件名加 _h264 后缀（保留原后缀）
// 例如: abc123.mp4 → abc123_h264.mp4
// 转码成功后用转码文件替换原文件
func transcodeVideo(filePath string, log zerolog.Logger) error {
	ext := filepath.Ext(filePath)             // .mp4
	baseNoExt := strings.TrimSuffix(filePath, ext) // /path/to/abc123
	transcodedPath := baseNoExt + "_h264" + ext  // /path/to/abc123_h264.mp4

	log.Info().Msgf("开始视频转码: %s → %s", filePath, transcodedPath)

	cmd := exec.Command("ffmpeg",
		"-i", filePath,
		"-c:v", "libx264",
		"-preset", "fast",            // 编码速度和压缩率的平衡
		"-crf", "23",                // 质量参数，18-28 之间，越小质量越好
		"-profile:v", "baseline",    // 兼容性最高，支持所有设备
		"-level", "4.0",             // 4.0 支持 4K(3840x2160)@30fps
		"-pix_fmt", "yuv420p",       // 强制8bit色深
		"-colorspace:v", "bt709",    // 强制转换色彩空间为 BT.709 SDR，解决 HLG/HDR 兼容性问题
		"-color_primaries:v", "bt709",
		"-color_transfer:v", "bt709",
		"-c:a", "aac",
		"-b:a", "128k",
		"-movflags", "+faststart", // Web 播放优化：将 moov atom 移到文件开头
		"-y",                    // 覆盖输出文件
		transcodedPath,
	)
	cmd.Stderr = nil // 避免大量日志输出

	if err := cmd.Run(); err != nil {
		// 转码失败，清理临时文件
		os.Remove(transcodedPath)
		return fmt.Errorf("ffmpeg 转码失败: %w", err)
	}

	// 检查转码后文件大小，如果比原文件还大（罕见情况），也替换
	// 获取转码后文件信息
	transcodedInfo, err := os.Stat(transcodedPath)
	if err != nil {
		os.Remove(transcodedPath)
		return fmt.Errorf("无法读取转码文件: %w", err)
	}

	// 用转码后的文件替换原文件
	originalPath := baseNoExt + ext
	backupPath := baseNoExt + "_original" + ext

	// 1. 备份原文件
	if err := os.Rename(originalPath, backupPath); err != nil {
		os.Remove(transcodedPath)
		return fmt.Errorf("备份原视频失败: %w", err)
	}

	// 2. 将转码文件重命名为原文件名
	if err := os.Rename(transcodedPath, originalPath); err != nil {
		// 恢复原文件
		os.Rename(backupPath, originalPath)
		os.Remove(transcodedPath)
		return fmt.Errorf("替换转码文件失败: %w", err)
	}

	// 3. 删除备份
	if err := os.Remove(backupPath); err != nil {
		log.Warn().Msgf("删除原视频备份失败: %v (不影响使用)", err)
	}

	log.Info().Msgf("视频转码完成: %s (%s)", originalPath, formatSize(transcodedInfo.Size()))
	return nil
}

// TranscodeVideoIfNeeded 检测视频编码，如果需要则转码为 H.264 Baseline
// 这是给 Upload handler 调用的入口函数
func TranscodeVideoIfNeeded(filePath string, log zerolog.Logger) {
	if !IsVideo(filePath) {
		return
	}

	// 检测编码格式、profile、像素格式
	info, err := probeVideoInfo(filePath)
	if err != nil {
		log.Warn().Msgf("检测视频编码失败，跳过转码: %s, err: %v", filePath, err)
		return
	}

	log.Info().Msgf("视频信息检测: %s → codec=%s profile=%s pix_fmt=%s",
		filepath.Base(filePath), info.Codec, info.Profile, info.PixFmt)

	if !needsTranscode(info) {
		log.Debug().Msgf("视频已符合兼容要求（H.264 Baseline 8bit），无需转码: %s", filepath.Base(filePath))
		return
	}

	if err := transcodeVideo(filePath, log); err != nil {
		log.Error().Msgf("视频转码失败: %s, err: %v", filePath, err)
	}
}

// formatSize 格式化文件大小
func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
