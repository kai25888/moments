<template>
  <NuxtLayout>
    <NuxtPage/>
    <Toaster richColors position="top-center"/>
  </NuxtLayout>
</template>

<script lang="ts" setup>
import { Toaster } from 'vue-sonner'
import '~/assets/simple-markdown.scss'

// 🔥 动态导入 Fancybox CSS - 减少首屏体积
const fancyboxLoaded = ref(false)
onMounted(async () => {
  if (!fancyboxLoaded.value) {
    await import('@fancyapps/ui/dist/fancybox/fancybox.css')
    fancyboxLoaded.value = true
  }
})
</script>

<style>
div.grecaptcha-badge{
  display: none !important;
}
</style>
