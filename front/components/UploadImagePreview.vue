<template>
  <div
    v-if="
      images.length > 0 &&
      ($route.path.startsWith('/new') || $route.path.startsWith('/edit'))
    "
    ref="el"
    class="grid gap-2"
    :style="gridStyle"
  >
    <div
      v-for="(img, i) in images"
      :key="img.id"
      class="relative"
      :class="
        images.length === 1
          ? 'full-cover-image-single'
          : 'full-cover-image-mult'
      "
    >
      <img
        :src="img.url"
        class="cursor-move rounded lazy-img"
        loading="lazy"
        decoding="async"
        @load="onImgLoad"
      />
      <div
        class="absolute top-0 right-0 px-1 bg-white dark:bg-gray-900 m-2 rounded hover:text-red-500 cursor-pointer"
        @click="removeImage(i)"
      >
        <UIcon name="i-carbon-trash-can" />
      </div>
    </div>
  </div>

  <template v-else-if="imageConfigs && imageConfigs.length">
    <MyFancyBox :style="gridStyle">
      <div
        v-for="(imageConfig, z) in imageConfigs"
        :key="z"
        :ref="(el) => setItemRef(el, z)"
        :href="imageConfig.url"
        :class="
          imageConfigs.length === 1
            ? 'full-cover-image-single'
            : 'full-cover-image-mult'
        "
      >
        <img
          class="cursor-zoom-in rounded lazy-img"
          :src="imageConfig.visibleSrc"
          :onerror="imageConfig.visibleSrc ? `javascript:this.src='${imageConfig.url}';this.onerror=null` : undefined"
          decoding="async"
          @load="onImgLoad"
        />
      </div>
    </MyFancyBox>
  </template>
</template>

<script setup lang="ts">
import { useSortable } from "@vueuse/integrations/useSortable";
import { useIntersectionObserver } from "@vueuse/core";

interface ImgConfig {
  id: number;
  url: string;
  thumbUrl: string;
}

interface ImgConfigWithSrc extends ImgConfig {
  visibleSrc: string | undefined;
}

const route = useRoute();
const el = ref(null);
const props = defineProps<{ imgs?: string; imgConfigs?: ImgConfig[] }>();
const emit = defineEmits(["removeImage", "dragImage"]);

const images = ref<ImgConfig[]>([]);
const imageConfigs = ref<ImgConfigWithSrc[]>([]);

// 每个图片容器的 DOM ref，用于 IntersectionObserver
const itemRefs = ref<(Element | null)[]>([]);
const observers: ReturnType<typeof useIntersectionObserver>[] = [];

const setItemRef = (el: unknown, index: number) => {
  if (el instanceof Element) {
    itemRefs.value[index] = el;
  }
};

const stopAllObservers = () => {
  observers.forEach((obs) => obs.stop());
  observers.length = 0;
};

const setupObservers = () => {
  stopAllObservers();
  nextTick(() => {
    imageConfigs.value.forEach((config, index) => {
      const target = itemRefs.value[index];
      if (!target) return;
      const obs = useIntersectionObserver(
        target as HTMLElement,
        ([entry]) => {
          if (entry.isIntersecting && !config.visibleSrc) {
            config.visibleSrc = config.thumbUrl;
            obs.stop();
          }
        },
        { rootMargin: "200px" }
      );
      observers.push(obs);
    });
  });
};

watchEffect(() => {
  images.value = (props.imgs || "")
    .split(",")
    .filter(Boolean)
    .map((img) => ({ id: Math.random(), url: img, thumbUrl: img }));
});

watchEffect(() => {
  itemRefs.value = [];
  imageConfigs.value = (props.imgConfigs || []).map((imgConfig) => ({
    ...imgConfig,
    id: Math.random(),
    visibleSrc: undefined,
  }));
  setupObservers();
});

watchEffect(() => {
  emit(
    "dragImage",
    images.value.map((img) => img.url)
  );
});

const removeImage = async (index: number) => {
  emit("removeImage", index);
};

const onImgLoad = (e: Event) => {
  (e.target as HTMLImageElement).classList.add("loaded");
};

onMounted(() => {
  if (route.path.startsWith("/new") || route.path.startsWith("/edit")) {
    setTimeout(() => {
      useSortable(el, images);
    }, 500);
  }
});

onBeforeUnmount(() => {
  stopAllObservers();
});

const gridStyle = computed(() => {
  const count = imageConfigs.value.length || images.value.length;
  let style = "max-width:100%; display:grid; gap: 0.5rem; align-items: start;";
  switch (count) {
    case 1:
      style += "grid-template-columns: 1fr; max-width:60%;";
      break;
    case 2:
      style += "grid-template-columns: 1fr 1fr; aspect-ratio: 2 / 1;";
      break;
    case 3:
      style += "grid-template-columns: 1fr 1fr 1fr; aspect-ratio: 3 / 1;";
      break;
    case 4:
      style += "grid-template-columns: 1fr 1fr; aspect-ratio: 1;";
      break;
    default:
      style += "grid-template-columns: 1fr 1fr 1fr;";
  }
  return style;
});
</script>

<style scoped>
.full-cover-image-mult {
  width: 100%;
  max-height: 300px;
  aspect-ratio: 1 / 1;
  contain-intrinsic-size: 0 300px;

  > img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    object-position: center;
  }
}

.full-cover-image-single {
  width: fit-content;
  contain-intrinsic-size: 0 300px;

  > img {
    max-height: 300px;
    object-fit: cover;
    object-position: center;
  }
}

.lazy-img {
  opacity: 0;
  transition: opacity 0.3s ease;
  background: #e5e7eb;
}

.lazy-img.loaded {
  opacity: 1;
}
</style>
