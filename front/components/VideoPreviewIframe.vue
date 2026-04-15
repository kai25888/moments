<template>
  <button
    v-if="shouldDeferLoad"
    type="button"
    class="video-placeholder"
    :class="{ 'video-placeholder--has-poster': !!displayPoster }"
    @click="activated = true"
  >
    <span
      v-if="displayPoster"
      class="video-placeholder-poster"
      :style="{ backgroundImage: `url(${displayPoster})` }"
      aria-hidden="true"
    />
    <span class="video-placeholder-scrim" aria-hidden="true" />
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
import { bilibiliBvidFromPlayerUrl } from "~/utils/videoPoster";

const props = withDefaults(
  defineProps<{
    url: string;
    lazy?: boolean;
    /** 外部传入的封面（如 YouTube 缩略图、同条动态配图） */
    poster?: string;
  }>(),
  {
    lazy: false,
    poster: "",
  },
);

const activated = ref(!props.lazy);
const bilibiliFetchedPoster = ref("");

watch(
  () => props.url,
  () => {
    activated.value = !props.lazy;
  },
);

const shouldDeferLoad = computed(() => props.lazy && !activated.value && !!props.url);

const displayPoster = computed(
  () => props.poster?.trim() || bilibiliFetchedPoster.value || "",
);

async function tryFetchBilibiliPoster(playerUrl: string) {
  bilibiliFetchedPoster.value = "";
  const bvid = bilibiliBvidFromPlayerUrl(playerUrl);
  if (!bvid) {
    return;
  }
  try {
    const res = await $fetch<{ pic: string | null }>("/api/bilibili-view", {
      query: { bvid },
      timeout: 10_000,
    });
    if (res?.pic) {
      bilibiliFetchedPoster.value = res.pic;
    }
  } catch {
    // 网络失败时静默忽略
  }
}

watch(
  () => props.url,
  (u) => {
    if (!u || !u.includes("player.bilibili.com")) {
      bilibiliFetchedPoster.value = "";
      return;
    }
    if (props.poster?.trim()) {
      return;
    }
    void tryFetchBilibiliPoster(u);
  },
  { immediate: true },
);

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
  position: relative;
  overflow: hidden;
  isolation: isolate;
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

.video-placeholder--has-poster {
  align-items: flex-end;
  justify-content: flex-start;
  padding-top: 56px;
}

.video-placeholder-poster {
  position: absolute;
  inset: 0;
  z-index: 0;
  background-size: cover;
  background-position: center;
  transform: scale(1.02);
}

.video-placeholder-scrim {
  position: absolute;
  inset: 0;
  z-index: 1;
  background: linear-gradient(
    to top,
    rgba(15, 23, 42, 0.72) 0%,
    rgba(15, 23, 42, 0.2) 45%,
    rgba(15, 23, 42, 0.12) 100%
  );
  pointer-events: none;
}

.video-placeholder:hover {
  transform: translateY(-1px);
  border-color: rgba(159, 200, 74, 0.5);
  background: rgba(159, 200, 74, 0.08);
}

.video-placeholder-icon {
  position: relative;
  z-index: 2;
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

.video-placeholder-text {
  position: relative;
  z-index: 2;
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

.video-placeholder--has-poster .video-placeholder-title {
  color: rgb(248 250 252);
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.35);
}

.video-placeholder--has-poster .video-placeholder-subtitle {
  color: rgb(226 232 240);
  opacity: 0.92;
}

.video-placeholder--has-poster .video-placeholder-icon {
  background: rgba(255, 255, 255, 0.92);
  color: rgb(63 98 18);
}

.dark .video-placeholder:not(.video-placeholder--has-poster) {
  border-color: rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.04);
  color: rgb(226 232 240);
}

.dark .video-placeholder:hover {
  border-color: rgba(159, 200, 74, 0.45);
  background: rgba(159, 200, 74, 0.12);
}

.dark .video-placeholder:not(.video-placeholder--has-poster) .video-placeholder-icon {
  background: rgba(159, 200, 74, 0.2);
  color: rgb(217 249 157);
}

.dark .video-placeholder:not(.video-placeholder--has-poster) .video-placeholder-subtitle {
  color: rgb(148 163 184);
}
</style>
