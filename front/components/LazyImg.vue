<template>
  <img
    v-if="isVisible || !lazy"
    :src="finalSrc"
    :alt="alt"
    :class="imgClass"
    :style="imgStyle"
    loading="lazy"
    @load="onLoad"
    @error="onError"
  />
  <img
    v-else
    :src="placeholder"
    :alt="alt"
    :class="imgClass"
    :style="imgStyle"
  />
</template>

<script setup lang="ts">
interface Props {
  src: string
  alt?: string
  imgClass?: string
  imgStyle?: string
  lazy?: boolean
  placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
  alt: '',
  imgClass: '',
  imgStyle: '',
  lazy: true,
  placeholder: 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1 1"%3E%3Crect fill="%23f0f0f0"/%3E%3C/svg%3E',
})

const emit = defineEmits(['load', 'error'])

const isVisible = ref(false)
const finalSrc = ref(props.placeholder)

// Intersection Observer for lazy loading
let observer: IntersectionObserver | null = null

onMounted(() => {
  if (!props.lazy) {
    isVisible.value = true
    finalSrc.value = props.src
    return
  }

  observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          isVisible.value = true
          finalSrc.value = props.src
          observer?.disconnect()
        }
      })
    },
    { rootMargin: '100px' }
  )

  const imgEl = document.querySelector(`img[src="${props.placeholder}"]`)
  if (imgEl) {
    observer.observe(imgEl)
  }
})

onUnmounted(() => {
  observer?.disconnect()
})

watch(() => props.src, (newSrc) => {
  if (!props.lazy) {
    finalSrc.value = newSrc
  }
})

const onLoad = () => emit('load')
const onError = () => emit('error')
</script>
