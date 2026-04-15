package handler

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/rs/zerolog"
)

// 缩略图尺寸常量
const (
	ThumbSmall  = 300 // 卡片预览
	ThumbMedium = 600 // 列表展示
	ThumbLarge  = 1200 // 弹窗查看
)

// ImageSizes 所有支持的缩略图尺寸
var ImageSizes = []int{ThumbSmall, ThumbMedium, ThumbLarge}

// 尺寸后缀映射
var sizeSuffixes = map[int]string{
	ThumbSmall:  "300w",
	ThumbMedium: "600w",
	ThumbLarge:  "1200w",
}

// WebPQuality WebP 质量设置
var WebPQuality = map[int]float64{
	ThumbSmall:  75, // 小图用较低质量
	ThumbMedium: 80, // 中图用中等质量
	ThumbLarge:  85, // 大图用较高质量
}

// 判断是否支持压缩
func SupportCompress(filename string) bool {
	_, err := imaging.FormatFromFilename(filename)
	return err == nil
}

// GetSizeSuffix 获取尺寸对应的后缀
func GetSizeSuffix(width int) string {
	if suffix, ok := sizeSuffixes[width]; ok {
		return suffix
	}
	return fmt.Sprintf("%dw", width)
}

// CompressImage 接收图片路径并进行压缩（兼容旧接口）
// 返回错误信息
func CompressImage(f FileHandler, originImagePath string, thumbImagePath string, quality int) error {
	f.base.log.Debug().Msgf("start compressing image, originImagePath=%s, quality=%d", originImagePath, quality)

	originImage, err := imaging.Open(originImagePath, imaging.AutoOrientation(true))
	if err != nil {
		f.base.log.Error().Msgf("failed to open image: %v", err)
		return err
	}

	thumbImageFile, err := os.Create(thumbImagePath)
	if err != nil {
		f.base.log.Error().Msgf("failed to create thumb image: %v", err)
		return err
	}
	defer thumbImageFile.Close()

	compressedImage := imaging.Resize(originImage, 600, 0, imaging.Lanczos)
	err = imaging.Encode(thumbImageFile, compressedImage, imaging.JPEG)
	if err != nil {
		f.base.log.Error().Msgf("failed to save image: %v", err)
		os.Remove(thumbImagePath)
		return err
	}

	f.base.log.Debug().Msgf("compressing image finished, originImagePath=%s, quality=%d", originImagePath, quality)
	return nil
}

// saveWebp 保存 WebP 图片
func saveWebp(img image.Image, path string, quality float64) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// webp.Encode 接受 io.Writer 和 image.Image，质量参数范围 0-100
	return webp.Encode(file, img, &webp.Options{Quality: float32(quality)})
}

// GenerateAllThumbs 生成多种尺寸的 WebP 缩略图
// originPath: 原始图片路径（临时文件）
// finalPath: 最终 WebP 文件路径（无扩展名）
// 返回生成的文件路径列表和错误
func GenerateAllThumbs(originPath string, finalPath string, log zerolog.Logger) ([]string, error) {
	log.Debug().Msgf("开始生成多尺寸缩略图: %s -> %s", originPath, finalPath)

	// 打开原图
	img, err := imaging.Open(originPath, imaging.AutoOrientation(true))
	if err != nil {
		log.Error().Msgf("打开图片失败: %v", err)
		return nil, err
	}

	originalWidth := img.Bounds().Dx()
	generatedFiles := make([]string, 0, len(ImageSizes)+1)

	// 按尺寸生成缩略图
	for _, size := range ImageSizes {
		// 如果原图小于目标尺寸，跳过（使用原图）
		if originalWidth < size {
			log.Debug().Msgf("原图宽度 %d < %d，跳过该尺寸", originalWidth, size)
			continue
		}

		suffix := GetSizeSuffix(size)
		thumbPath := fmt.Sprintf("%s_%s.webp", finalPath, suffix)
		quality := WebPQuality[size]

		// 缩放图片
		thumb := imaging.Resize(img, size, 0, imaging.Lanczos)

		// 编码为 WebP
		if err := saveWebp(thumb, thumbPath, quality); err != nil {
			log.Error().Msgf("生成缩略图 %s 失败: %v", thumbPath, err)
			continue
		}

		generatedFiles = append(generatedFiles, thumbPath)
		log.Debug().Msgf("生成缩略图成功: %s (质量: %.0f)", thumbPath, quality)
	}

	// 生成主图 WebP（始终生成，质量最高）
	mainWebP := finalPath + ".webp"
	if err := saveWebp(img, mainWebP, 85); err != nil {
		log.Error().Msgf("生成主图 WebP 失败: %v", err)
		return generatedFiles, err
	}
	generatedFiles = append(generatedFiles, mainWebP)
	log.Debug().Msgf("生成主图 WebP: %s", mainWebP)

	return generatedFiles, nil
}

// GenerateThumbnailSimple 生成单个缩略图（简单版本，用于兼容性）
func GenerateThumbnailSimple(originPath string, thumbPath string, width int) error {
	img, err := imaging.Open(originPath, imaging.AutoOrientation(true))
	if err != nil {
		return err
	}

	thumb := imaging.Resize(img, width, 0, imaging.Lanczos)

	// 打开文件进行写入
	file, err := os.Create(thumbPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 根据扩展名判断格式
	ext := strings.ToLower(filepath.Ext(thumbPath))
	switch ext {
	case ".webp":
		return webp.Encode(file, thumb, &webp.Options{Quality: 80})
	case ".jpg", ".jpeg":
		return jpeg.Encode(file, thumb, &jpeg.Options{Quality: 80})
	case ".png":
		return png.Encode(file, thumb)
	default:
		// 默认使用 WebP
		return webp.Encode(file, thumb, &webp.Options{Quality: 80})
	}
}

// GetBestThumbnailUrl 获取最佳缩略图 URL
// 根据原始 URL 返回最合适的缩略图路径
func GetBestThumbnailUrl(originalUrl string, targetWidth int) string {
	if originalUrl == "" {
		return ""
	}

	// 如果是上传文件，返回对应尺寸的缩略图
	if strings.HasPrefix(originalUrl, "/upload/") {
		ext := filepath.Ext(originalUrl)
		base := strings.TrimSuffix(originalUrl, ext)

		// 根据目标尺寸选择后缀
		var suffix string
		if targetWidth <= ThumbSmall {
			suffix = "_300w"
		} else if targetWidth <= ThumbMedium {
			suffix = "_600w"
		} else {
			suffix = "_1200w"
		}

		// 返回 WebP 缩略图
		webpUrl := base + suffix + ".webp"
		return webpUrl
	}

	return originalUrl
}
