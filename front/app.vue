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
</script>

<style>
div.grecaptcha-badge{
  display: none !important;
}
</style>
