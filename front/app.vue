<template>
  <NuxtLayout>
    <NuxtPage/>
    <Toaster richColors position="top-center"/>
  </NuxtLayout>
</template>

<script lang="ts" setup>
import { Toaster } from 'vue-sonner'
import '~/assets/simple-markdown.scss'

const fancyboxLoaded = ref(false)
onMounted(async () => {
  await import('@fancyapps/ui/dist/fancybox/fancybox.css')
  fancyboxLoaded.value = true
})

// 📱 原生桥接（Capacitor App）
const { isNative } = useNativeBridge()
</script>

<style>
div.grecaptcha-badge{
  display: none !important;
}

/* 📱 Capacitor 安全区域适配 */
@supports (padding-top: env(safe-area-inset-top)) {
  :root {
    --safe-area-top: env(safe-area-inset-top);
    --safe-area-bottom: env(safe-area-inset-bottom);
    --safe-area-left: env(safe-area-inset-left);
    --safe-area-right: env(safe-area-inset-right);
  }
}
</style>
