<template>
  <div ref="containerRef" class="lazy-image-container" :class="imgClass">
    <!-- 加载中骨架屏 -->
    <div 
      v-if="loading" 
      class="lazy-skeleton"
      :style="skeletonStyle"
    />
    
    <!-- 🔥 使用 NuxtImg 组件进行图片优化 -->
    <NuxtImg
      v-if="useNuxtImage && !error"
      v-show="!loading && loaded"
      ref="nuxtImgRef"
      :src="currentSrc"
      :alt="alt"
      :style="imgStyle"
      :width="imgWidth"
      :height="imgHeight"
      :loading="lazy ? 'lazy' : 'eager'"
      :format="imageFormat"
      :quality="imageQuality"
      :placeholder="placeholder"
      @load="onImgLoad"
      @error="onImgError"
    />
    
    <!-- 回退原生 img 标签 -->
    <img
      v-else-if="!error"
      v-show="!loading && loaded"
      ref="imgRef"
      :src="currentSrc"
      :alt="alt"
      :style="imgStyle"
      :loading="lazy ? 'lazy' : 'eager'"
      :width="imgWidth"
      :height="imgHeight"
      @load="onImgLoad"
      @error="onImgError"
    />
    
    <!-- 错误占位符 -->
    <div 
      v-if="error" 
      class="lazy-error"
    >
      <UIcon name="i-carbon-image" class="w-8 h-8 text-gray-400" />
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  src: string
  alt?: string
  imgClass?: string
  imgStyle?: string
  lazy?: boolean
  placeholder?: string
  thumbSrc?: string       // 缩略图 URL
  aspectRatio?: string    // 可选：维持宽高比，如 "1/1" 或 "16/9"
  width?: number          // 图片宽度
  height?: number         // 图片高度
  quality?: number        // 图片质量 1-100
  format?: 'avif' | 'webp' | 'original'  // 输出格式
}

const props = withDefaults(defineProps<Props>(), {
  alt: '',
  imgClass: '',
  imgStyle: '',
  lazy: true,
  placeholder: '',
  quality: 80,
  format: 'webp',
})

const emit = defineEmits(['load', 'error'])

const containerRef = ref<HTMLDivElement | null>(null)
const imgRef = ref<HTMLImageElement | null>(null)
const nuxtImgRef = ref<any>(null)
const loading = ref(true)
const loaded = ref(false)
const error = ref(false)
const currentSrc = ref('')

// 检测是否可以使用 NuxtImg
const useNuxtImage = ref(true)
const imgWidth = computed(() => props.width)
const imgHeight = computed(() => props.height)

// 图片格式
const imageFormat = computed(() => {
  if (props.format === 'original') return undefined
  return props.format
})
const imageQuality = computed(() => props.quality)

// 骨架屏样式
const skeletonStyle = computed(() => {
  if (props.aspectRatio) {
    return { aspectRatio: props.aspectRatio }
  }
  return {}
})

// 优先级：缩略图 -> 占位图 -> 原图
const getSrcToLoad = () => {
  return props.thumbSrc || props.src
}

// Intersection Observer
let observer: IntersectionObserver | null = null

const initObserver = () => {
  if (!props.lazy || !containerRef.value) {
    // 非懒加载模式：直接加载
    currentSrc.value = getSrcToLoad()
    return
  }

  observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          // 进入视口时加载图片
          currentSrc.value = getSrcToLoad()
          observer?.disconnect()
        }
      })
    },
    {
      rootMargin: '200px', // 提前 200px 开始加载
      threshold: 0.01
    }
  )

  observer.observe(containerRef.value)
}

const onImgLoad = () => {
  loading.value = false
  loaded.value = true
  error.value = false
  emit('load')
}

const onImgError = () => {
  // 缩略图加载失败，尝试原图
  if (currentSrc.value !== props.src) {
    console.warn('缩略图加载失败，尝试原图:', props.src)
    currentSrc.value = props.src
    return
  }
  
  loading.value = false
  error.value = true
  emit('error')
}

onMounted(() => {
  initObserver()
})

onUnmounted(() => {
  observer?.disconnect()
})

// 监听 src 变化
watch(() => props.src, (newSrc) => {
  loading.value = true
  loaded.value = false
  error.value = false
  currentSrc.value = props.lazy ? '' : getSrcToLoad()
  
  if (!props.lazy) {
    currentSrc.value = getSrcToLoad()
  }
})
</script>

<style scoped>
.lazy-image-container {
  position: relative;
  overflow: hidden;
  background-color: #f0f0f0;
}

.dark .lazy-image-container {
  background-color: #2a2a2a;
}

.lazy-skeleton {
  width: 100%;
  height: 100%;
  min-height: 100px;
  background: linear-gradient(
    90deg,
    #f0f0f0 25%,
    #e0e0e0 50%,
    #f0f0f0 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

.dark .lazy-skeleton {
  background: linear-gradient(
    90deg,
    #2a2a2a 25%,
    #3a3a3a 50%,
    #2a2a2a 75%
  );
  background-size: 200% 100%;
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

.lazy-error {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  min-height: 80px;
  background-color: #f5f5f5;
}

.dark .lazy-error {
  background-color: #1a1a1a;
}
</style>
