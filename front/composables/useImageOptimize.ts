/**
 * 🔥 图片优化工具
 * 提供图片 URL 转换、WebP/AVIF 支持、占位符等功能
 */

// 图片服务配置（根据实际情况修改）
const IMAGE_CONFIG = {
    // 是否启用图片优化
    enabled: true,
    // 默认质量
    quality: 80,
    // 支持的格式（按优先级）
    formats: ['avif', 'webp'] as const,
}

/**
 * 获取优化后的图片 URL
 * @param url 原始图片 URL
 * @param options 优化选项
 */
export function useOptimizedImage(
    url: string, 
    options: {
        width?: number
        height?: number
        quality?: number
        format?: 'avif' | 'webp' | 'original'
    } = {}
) {
    const { width, height, quality = IMAGE_CONFIG.quality, format = 'webp' } = options
    
    // 如果不启用优化或 URL 为空，返回原 URL
    if (!IMAGE_CONFIG.enabled || !url) {
        return {
            src: url,
            width,
            height,
            loading: 'lazy' as const,
            placeholder: '',
        }
    }
    
    // 构建优化 URL（如果有图片处理服务可以在这里配置）
    // 目前返回原 URL + loading="lazy"，实际优化需要配合图片 CDN 或 @nuxt/image 组件
    const optimizedSrc = buildOptimizedUrl(url, { width, height, quality, format })
    
    return {
        src: optimizedSrc,
        width,
        height,
        loading: 'lazy' as const,
        // 可选的 Base64 占位符（低带宽场景）
        placeholder: '',
    }
}

/**
 * 构建优化后的图片 URL
 * 如果后端有图片处理服务，可以在这里配置规则
 */
function buildOptimizedUrl(
    url: string, 
    options: { width?: number; height?: number; quality?: number; format?: string }
): string {
    // 如果有图片 CDN/处理服务，可以在这里添加处理规则
    // 例如：https://cdn.example.com?url=xxx&w=xxx&q=xxx&f=webp
    
    // 目前直接返回原 URL，@nuxt/image 组件会在使用时自动处理
    return url
}

/**
 * 生成图片占位符（可选功能）
 * 可用于实现模糊加载效果
 */
export function useImagePlaceholder() {
    // 生成一个简单的 SVG 占位符
    const generatePlaceholderSvg = (color = '#f0f0f0') => {
        return `data:image/svg+xml;base64,${btoa(`
            <svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
                <rect fill="${color}" width="100" height="100"/>
            </svg>
        `)}`
    }
    
    return {
        generatePlaceholderSvg,
    }
}
