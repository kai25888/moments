package handler

import (
	"fmt"
	"image"
	"os"
	"path"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/rs/zerolog"
)

// ThumbSmall / ThumbMedium / ThumbLarge 三档缩略图宽度（像素）
const (
	ThumbSmall  = 300
	ThumbMedium = 600
	ThumbLarge  = 1200
)

// WebPQuality 各档缩略图的 WebP 质量
var WebPQuality = map[int]float32{
	ThumbSmall:  75,
	ThumbMedium: 80,
	ThumbLarge:  85,
}

// SupportCompress 判断是否支持压缩
func SupportCompress(filename string) bool {
	_, err := imaging.FormatFromFilename(filename)
	return err == nil
}

// encodeWebP 将 image.Image 编码为 WebP 文件
func encodeWebP(img image.Image, destPath string, quality float32) error {
	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return webp.Encode(f, img, &webp.Options{Lossless: false, Quality: quality})
}

// GenerateAllThumbs 对图片路径 srcPath 生成三档 WebP 缩略图。
// baseName 是不含后缀的目标文件名前缀（含目录），例如 "/data/upload/abc123"。
// 生成的文件分别为 baseName_300w.webp、baseName_600w.webp、baseName_1200w.webp。
// 返回生成成功的文件路径列表。
func GenerateAllThumbs(srcPath string, baseName string, log zerolog.Logger) ([]string, error) {
	orig, err := imaging.Open(srcPath, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("open image: %w", err)
	}

	origWidth := orig.Bounds().Dx()
	var generated []string

	for _, width := range []int{ThumbSmall, ThumbMedium, ThumbLarge} {
		destPath := fmt.Sprintf("%s_%dw.webp", baseName, width)

		// 如果文件已存在，跳过
		if _, statErr := os.Stat(destPath); statErr == nil {
			generated = append(generated, destPath)
			continue
		}

		// 如果原图比目标宽度小，直接使用原图尺寸编码
		var resized image.Image
		if origWidth <= width {
			resized = orig
		} else {
			resized = imaging.Resize(orig, width, 0, imaging.Lanczos)
		}

		q := WebPQuality[width]
		if err := encodeWebP(resized, destPath, q); err != nil {
			log.Error().Msgf("生成 WebP 缩略图失败 %s: %v", destPath, err)
			continue
		}
		log.Debug().Msgf("生成 WebP 缩略图成功: %s", destPath)
		generated = append(generated, destPath)
	}

	return generated, nil
}

// CompressImage 保留旧接口兼容性：生成单个 JPEG 缩略图
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

	compressedImage := imaging.Resize(originImage, 400, 0, imaging.Lanczos)
	err = imaging.Encode(thumbImageFile, compressedImage, imaging.JPEG, imaging.JPEGQuality(quality))
	if err != nil {
		f.base.log.Error().Msgf("failed to save image: %v", err)
		os.Remove(thumbImagePath)
		return err
	}

	f.base.log.Debug().Msgf("compressing image finished, originImagePath=%s, quality=%d", originImagePath, quality)
	return nil
}

// GetBestThumbnailUrl 根据文件是否存在返回最佳缩略图 URL，降级为原图 URL。
func GetBestThumbnailUrl(uploadDir string, sha256 string, fallbackUrl string) string {
	for _, width := range []int{ThumbMedium, ThumbSmall, ThumbLarge} {
		filename := fmt.Sprintf("%s_%dw.webp", sha256, width)
		if _, err := os.Stat(path.Join(uploadDir, filename)); err == nil {
			return "/upload/" + filename
		}
	}
	return fallbackUrl
}
