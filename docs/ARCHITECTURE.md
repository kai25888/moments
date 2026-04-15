# Moments-Fork 架构优化方案

## 场景定位

| 维度 | 现状 |
|------|------|
| 用户数 | 2-3人 |
| 主要操作 | 上传图片/视频 |
| 部署环境 | NAS 单机 |
| 核心痛点 | 图片加载慢 |

> **设计原则**: 保持 SQLite + 轻量增强，不引入过度复杂度

---

## 一、当前架构问题

```
┌─────────────────────────────────────────────────────┐
│                    当前架构                           │
├─────────────────────────────────────────────────────┤
│  用户 → API → SQLite + 本地文件 (NAS)               │
│                    ↓                                 │
│         无缓存 / 无缩略图策略 / 无CDN                │
└─────────────────────────────────────────────────────┘

问题清单：
1. 图片直接加载原图，无缩略图
2. 无浏览器缓存策略
3. 无服务端缓存（Nginx）
4. NAS 网络延迟放大
```

---

## 二、目标架构

```
┌─────────────────────────────────────────────────────────────┐
│                      NAS 单机部署                             │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────┐    ┌────────────┐    ┌──────────────────────┐ │
│  │   Nginx  │───→│  Go API    │───→│   SQLite (DB)       │ │
│  │ + 缓存   │    │  (Echo)    │    │   + 本地文件存储     │ │
│  │ + Gzip   │    │            │    │                      │ │
│  └──────────┘    └────────────┘    └──────────────────────┘ │
│        ↓                                                    │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              多尺寸缩略图 (生成时处理)                  │   │
│  │   {hash}_300w.jpg  {hash}_600w.jpg  {hash}_orig.webp  │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## 三、后端优化方案

### 3.1 图片处理优化

#### 3.1.1 多尺寸缩略图生成

```go
// backend/handler/compress.go

const (
    ThumbSmall  = 300  // 卡片预览
    ThumbMedium = 600  // 列表展示
    ThumbLarge  = 1200 // 弹窗查看
)

// GenerateAllThumbs 生成多种尺寸缩略图
func GenerateAllThumbs(originPath string) error {
    img, err := imaging.Open(originPath, imaging.AutoOrientation(true))
    if err != nil {
        return err
    }

    // 转为 WebP 存储（体积减少 30-50%）
    baseName := strings.TrimSuffix(originPath, filepath.Ext(originPath))

    // 小图
    small := imaging.Resize(img, ThumbSmall, 0, imaging.Lanczos)
    imaging.Encode(baseName+"_300w.webp", small, imaging.WebP, imaging.WebPQuality(75))

    // 中图
    medium := imaging.Resize(img, ThumbMedium, 0, imaging.Lanczos)
    imaging.Encode(baseName+"_600w.webp", medium, imaging.WebP, imaging.WebPQuality(80))

    // 大图（弹窗用）
    if img.Bounds().Dx() > ThumbLarge {
        large := imaging.Resize(img, ThumbLarge, 0, imaging.Lanczos)
        imaging.Encode(baseName+"_1200w.webp", large, imaging.WebP, imaging.WebPQuality(85))
    }

    // 删除原图节省空间
    os.Remove(originPath)

    return nil
}
```

#### 3.1.2 上传接口改造

```go
// backend/handler/file.go - Upload 函数修改

func (f FileHandler) Upload(c echo.Context) error {
    // ... 文件读取逻辑 ...

    for _, file := range files {
        // 计算 hash
        sha256, _ := fs_util.Sha256(reader)
        ext := filepath.Ext(file.Filename)

        // 统一转为 WebP
        filename := fmt.Sprintf("%s%s", sha256, ".webp")
        filePath := path.Join(f.base.cfg.UploadDir, filename)

        // 保存原图后转码
        if !fs_util.Exists(filePath) {
            // 临时保存原图
            tmpPath := filePath + ".tmp"
            io.Copy(dst, reader)

            // 异步生成缩略图
            go func() {
                if err := GenerateAllThumbs(tmpPath); err != nil {
                    f.base.log.Error().Msgf("缩略图生成失败: %v", err)
                }
            }()
        }

        // 返回 600w 缩略图
        result = append(result, "/upload/"+strings.TrimSuffix(filename, ".webp")+"_600w.webp")
    }

    return SuccessResp(c, result)
}
```

### 3.2 Nginx 缓存配置

```nginx
# /etc/nginx/conf.d/moments.conf

server {
    listen 11307;
    server_name _;

    # 根目录
    root /path/to/moments-fork/backend;
    index index.html;

    # API 代理
    location /api/ {
        proxy_pass http://127.0.0.1:37892;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # 上传文件缓存策略
    location /upload/ {
        # 启用本地缓存
        proxy_cache moments_cache;
        proxy_cache_valid 200 7d;
        proxy_cache_valid 404 1m;
        proxy_cache_use_stale error timeout updating;
        add_header X-Cache-Status $upstream_cache_status;

        # 开启 gzip（针对 JSON）
        gzip on;
        gzip_types application/json text/plain;

        # 缓存控制
        expires 7d;
        add_header Cache-Control "public, max-age=604800";

        proxy_pass http://127.0.0.1:37892;
    }

    # 缩略图：更长期缓存
    location ~* /upload/.*_(300w|600w|1200w) {
        expires 30d;
        add_header Cache-Control "public, max-age=2592000, immutable";
        proxy_pass http://127.0.0.1:37892;
    }

    # 静态资源
    location / {
        try_files $uri $uri/ /index.html;
        expires 7d;
        add_header Cache-Control "public";
    }
}

# 缓存区配置
proxy_cache_path /tmp/nginx_cache levels=1:2 keys_zone=moments_cache:10m max_size=100m inactive=7d use_temp_path=off;
```

### 3.3 响应头优化

```go
// backend/middleware/cache.go

package middleware

import (
    "github.com/labstack/echo/v4"
)

func StaticCache(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        path := c.Path()

        // 缩略图长期缓存
        if match, _ := regexp.MatchString(`.*_(300w|600w|1200w)\.webp$`, path); match {
            c.Response().Header().Set("Cache-Control", "public, max-age=2592000, immutable")
            c.Response().Header().Set("Expires", "Thu, 31 Dec 2037 23:59:59 GMT")
        }

        return next(c)
    }
}
```

### 3.4 数据库索引（已添加）

```go
// backend/migrate.go - 已执行的优化

func migrateAddIndexes(tx *gorm.DB) {
    tx.Exec(`CREATE INDEX IF NOT EXISTS idx_memo_userId ON Memo(userId)`)
    tx.Exec(`CREATE INDEX IF NOT EXISTS idx_memo_createdAt ON Memo(createdAt DESC)`)
    tx.Exec(`CREATE INDEX IF NOT EXISTS idx_memo_pinned ON Memo(pinned DESC, createdAt DESC)`)
    tx.Exec(`CREATE INDEX IF NOT EXISTS idx_comment_memoId ON Comment(memoId)`)
}
```

---

## 四、实施计划

| 阶段 | 内容 | 优先级 | 工作量 |
|------|------|--------|--------|
| 1 | 图片转 WebP + 多尺寸缩略图 | 🔴 高 | 1-2天 |
| 2 | Nginx 缓存配置 | 🔴 高 | 0.5天 |
| 3 | 响应头优化 | 🟡 中 | 0.5天 |
| 4 | 现有图片批量转码 | 🟡 中 | 1天 |

---

## 五、预期效果

| 指标 | 优化前 | 优化后 |
|------|--------|--------|
| 图片加载 | 原图 2-5MB | 缩略图 50-200KB |
| 首屏时间 | 3-8s | < 1s |
| 缓存命中率 | 0% | > 80% |
| 存储空间 | 100% | ~60% (WebP) |

---

## 六、代码修改清单

| 文件 | 修改内容 |
|------|----------|
| `compress.go` | 新增 `GenerateAllThumbs()` 函数 |
| `file.go` | Upload 返回缩略图 URL |
| `router.go` | 添加 StaticCache 中间件 |
| `docker-compose.yml` | 添加 Nginx 容器 |

---

## 七、注意事项

1. **现有图片**: 需要写脚本批量转码
2. **兼容性**: WebP 浏览器支持率 > 95%，无需降级
3. **回滚方案**: 保留原图路径，渐进式切换
