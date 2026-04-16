<script setup>
// 🔥 动态导入 Fancybox - 减少首屏 ~100KB
import('@fancyapps/ui/dist/index.esm.js').then(({ Fancybox }) => {
  window.Fancybox = Fancybox
})

const props = defineProps({
  options: Object,
});
const container = ref(null);

const randomId = randomHexStr();

onMounted(async () => {
  // 等待 Fancybox 加载完成
  await waitForFancybox()
  
  Array.from(container.value.children).map((el) => {
    el.setAttribute('data-fancybox', `gallery-${randomId}`);
  });
  window.Fancybox.bind(`[data-fancybox="gallery-${randomId}"]`, {
    Thumbs: {
      type: 'modern',
    },
    ...(props.options || {}),
  });
});

nextTick(async () => {
  await waitForFancybox()
  window.Fancybox.unbind(container.value);
  window.Fancybox.close();

  window.Fancybox.bind(`[data-fancybox="gallery-${randomId}"]`, {
    Thumbs: {
      type: 'modern',
    },
    ...(props.options || {}),
  });
});

function randomHexStr(len = 16, chars = '0123456789abcdefghijklmnopqrstuvwxyz') {
  let str = '';
  let length = chars.length;
  while (len > 0) {
    str += chars[Math.floor(Math.random() * length)];
    len--;
  }
  return str;
}

// 等待 Fancybox 加载
function waitForFancybox(timeout = 5000) {
  return new Promise((resolve, reject) => {
    if (window.Fancybox) {
      resolve(true)
      return
    }
    
    const start = Date.now()
    const check = () => {
      if (window.Fancybox) {
        resolve(true)
      } else if (Date.now() - start > timeout) {
        reject(new Error('Fancybox 加载超时'))
      } else {
        setTimeout(check, 50)
      }
    }
    check()
  })
}

onUnmounted(() => {
  if (window.Fancybox) {
    window.Fancybox.destroy();
  }
});
</script>

<template>
  <div ref="container">
    <slot></slot>
  </div>
</template>

<style></style>
