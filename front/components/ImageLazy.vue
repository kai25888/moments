<template>
  <div
    class="image-lazy-wrapper"
    :style="wrapperStyle"
    ref="wrapperRef"
  >
    <!-- 骨架屏占位，图片加载完成前显示 -->
    <div
      v-if="!loaded"
      class="image-lazy-skeleton"
      :style="skeletonStyle"
    />
    <!-- 实际图片 -->
    <img
      v-if="currentSrc"
      :src="currentSrc"
      :alt="alt"
      v-bind="$attrs"
      class="image-lazy-img"
      :class="{ 'image-lazy-img--loaded': loaded }"
      @load="onLoad"
      @error="onError"
    />
  </div>
</template>

<script setup lang="ts">
interface Props {
  src: string
  thumbSrc?: string
  alt?: string
  aspectRatio?: string
}

const props = withDefaults(defineProps<Props>(), {
  alt: '',
  aspectRatio: undefined,
})

const wrapperRef = ref<HTMLElement | null>(null)
const currentSrc = ref<string | null>(null)
const loaded = ref(false)
let observer: IntersectionObserver | null = null

const wrapperStyle = computed(() => {
  const style: Record<string, string> = {
    position: 'relative',
    overflow: 'hidden',
  }
  if (props.aspectRatio) {
    style['aspect-ratio'] = props.aspectRatio
  }
  return style
})

const skeletonStyle = computed(() => {
  const style: Record<string, string> = {
    position: 'absolute',
    inset: '0',
    width: '100%',
    height: '100%',
  }
  return style
})

const getSrcToLoad = () => props.thumbSrc || props.src

const onLoad = () => {
  loaded.value = true
}

const onError = () => {
  // 缩略图加载失败时降级到原图
  if (currentSrc.value !== props.src) {
    currentSrc.value = props.src
  }
}

const setupObserver = () => {
  if (observer) {
    observer.disconnect()
    observer = null
  }
  if (!wrapperRef.value) return

  observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          currentSrc.value = getSrcToLoad()
          observer?.disconnect()
          observer = null
        }
      })
    },
    {
      rootMargin: '200px',
      threshold: 0.01,
    }
  )
  observer.observe(wrapperRef.value)
}

onMounted(() => {
  setupObserver()
})

onBeforeUnmount(() => {
  observer?.disconnect()
  observer = null
})

watch(
  () => props.src,
  () => {
    loaded.value = false
    currentSrc.value = null
    nextTick(setupObserver)
  }
)
</script>

<style scoped>
.image-lazy-wrapper {
  display: block;
}

.image-lazy-skeleton {
  background: linear-gradient(90deg, #e5e7eb 25%, #f3f4f6 50%, #e5e7eb 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: inherit;
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

.dark .image-lazy-skeleton {
  background: linear-gradient(90deg, #374151 25%, #4b5563 50%, #374151 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

.image-lazy-img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.image-lazy-img--loaded {
  opacity: 1;
}
</style>
