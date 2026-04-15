<template>
  <button
    v-if="shouldDeferLoad"
    type="button"
    class="video-placeholder"
    @click="activated = true"
  >
    <div class="video-placeholder-icon">
      <UIcon name="i-carbon-play-filled-alt" class="h-6 w-6" />
    </div>
    <div class="video-placeholder-text">
      <div class="video-placeholder-title">点击加载视频</div>
      <div class="video-placeholder-subtitle">嵌入播放器将按需加载</div>
    </div>
  </button>

  <iframe
    v-else-if="videoUrl"
    :src="videoUrl"
    class="w-full h-[250px] rounded"
    title="Video player iframe"
    frameborder="0"
    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
    referrerpolicy="strict-origin-when-cross-origin"
    allowfullscreen
  >
  </iframe>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    url: string;
    lazy?: boolean;
  }>(),
  {
    lazy: false,
  },
);

const activated = ref(!props.lazy);

watch(
  () => props.url,
  () => {
    activated.value = !props.lazy;
  },
);

const shouldDeferLoad = computed(() => props.lazy && !activated.value && !!props.url);

const videoUrl = computed(() => {
  if (shouldDeferLoad.value) {
    return "";
  }

  try {
    const url = new URL(props.url);
    url.searchParams.set("autoplay", "0");
    return url.toString();
  } catch {
    return "";
  }
});
</script>

<style scoped>
.video-placeholder {
  width: 100%;
  min-height: 180px;
  border: 1px solid rgba(148, 163, 184, 0.22);
  border-radius: 16px;
  padding: 20px;
  background: rgba(15, 23, 42, 0.03);
  color: rgb(71 85 105);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 14px;
  text-align: left;
  transition: border-color 0.2s ease, background 0.2s ease, transform 0.2s ease;
}

.video-placeholder:hover {
  transform: translateY(-1px);
  border-color: rgba(159, 200, 74, 0.5);
  background: rgba(159, 200, 74, 0.08);
}

.video-placeholder-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 999px;
  background: rgba(159, 200, 74, 0.18);
  color: rgb(63 98 18);
  flex-shrink: 0;
}

.video-placeholder-title {
  font-size: 14px;
  font-weight: 600;
}

.video-placeholder-subtitle {
  margin-top: 4px;
  font-size: 12px;
  color: rgb(100 116 139);
}

.dark .video-placeholder {
  border-color: rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.04);
  color: rgb(226 232 240);
}

.dark .video-placeholder:hover {
  border-color: rgba(159, 200, 74, 0.45);
  background: rgba(159, 200, 74, 0.12);
}

.dark .video-placeholder-icon {
  background: rgba(159, 200, 74, 0.2);
  color: rgb(217 249 157);
}

.dark .video-placeholder-subtitle {
  color: rgb(148 163 184);
}
</style>
