package middleware

import (
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// CacheConfig 缓存中间件配置
type CacheConfig struct {
	// ThumbnailPattern 缩略图路径模式
	ThumbnailPattern string
}

// DefaultCacheConfig 默认缓存配置
var DefaultCacheConfig = CacheConfig{
	ThumbnailPattern: `.*_(300w|600w|1200w)\.webp$`,
}

// StaticCache 静态资源缓存中间件
// 为缩略图设置长期缓存，为其他资源设置适当缓存
func StaticCache() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()

			// 缩略图：长期缓存
			if isThumbnail(path) {
				c.Response().Header().Set("Cache-Control", "public, max-age=2592000, immutable")
				c.Response().Header().Set("Expires", time.Now().AddDate(1, 0, 0).Format(time.RFC1123))
				c.Response().Header().Set("Vary", "Accept-Encoding")
			}

			return next(c)
		}
	}
}

// isThumbnail 判断路径是否为缩略图
func isThumbnail(path string) bool {
	// 检查缩略图路径模式
	patterns := []string{
		`.*_(300w|600w|1200w)\.webp$`,
		`.*_thumb\.(jpg|jpeg|png|webp)$`,
	}

	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, path)
		if matched {
			return true
		}
	}

	// 检查是否有 thumbnail 相关参数
	return strings.Contains(path, "thumb")
}

// CorsHeaders CORS 头中间件
func CorsHeaders() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

			if c.Request().Method == "OPTIONS" {
				return c.NoContent(204)
			}

			return next(c)
		}
	}
}
