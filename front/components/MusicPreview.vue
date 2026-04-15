<template>
  <button
    v-if="shouldDeferLoad"
    type="button"
    class="music-placeholder"
    @click="activated = true"
  >
    <div class="music-placeholder-icon">
      <UIcon name="i-carbon-play-filled-alt" class="h-5 w-5" />
    </div>
    <div>
      <div class="music-placeholder-title">点击加载音乐卡片</div>
      <div class="music-placeholder-subtitle">按需初始化播放器，减少列表首屏负担</div>
    </div>
  </button>

  <meting-js v-else-if="id && server && type && api" :server="server" :type="type" :id="id" :api="api"/>
</template>

<script setup lang="ts">
import type {MetingJSDTO} from "@/types";

const props = withDefaults(defineProps<{
  id?: MetingJSDTO["id"];
  api?: MetingJSDTO["api"];
  server?: MetingJSDTO["server"];
  type?: MetingJSDTO["type"];
  lazy?: boolean;
}>(), {
  lazy: false,
})

const activated = ref(!props.lazy)

watch(
  () => [props.id, props.server, props.type, props.api],
  () => {
    activated.value = !props.lazy
  },
)

const shouldDeferLoad = computed(() => props.lazy && !activated.value && !!props.id && !!props.server && !!props.type && !!props.api)
</script>

<style scoped>
.music-placeholder {
  width: 100%;
  border: 1px solid rgba(148, 163, 184, 0.22);
  border-radius: 14px;
  padding: 14px 16px;
  background: rgba(15, 23, 42, 0.03);
  color: rgb(71 85 105);
  display: flex;
  align-items: center;
  gap: 12px;
  text-align: left;
  transition: border-color 0.2s ease, background 0.2s ease;
}

.music-placeholder:hover {
  border-color: rgba(159, 200, 74, 0.5);
  background: rgba(159, 200, 74, 0.08);
}

.music-placeholder-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border-radius: 999px;
  background: rgba(159, 200, 74, 0.18);
  color: rgb(63 98 18);
  flex-shrink: 0;
}

.music-placeholder-title {
  font-size: 14px;
  font-weight: 600;
}

.music-placeholder-subtitle {
  margin-top: 4px;
  font-size: 12px;
  color: rgb(100 116 139);
}

.dark .music-placeholder {
  border-color: rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.04);
  color: rgb(226 232 240);
}

.dark .music-placeholder:hover {
  border-color: rgba(159, 200, 74, 0.45);
  background: rgba(159, 200, 74, 0.12);
}

.dark .music-placeholder-icon {
  background: rgba(159, 200, 74, 0.2);
  color: rgb(217 249 157);
}

.dark .music-placeholder-subtitle {
  color: rgb(148 163 184);
}
</style>