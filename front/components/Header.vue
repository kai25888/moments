<template>
  <div
    v-if="$route.path !== '/new' && $route.path.indexOf('/edit/') < 0"
    class="header relative mb-14"
  >
    <!-- 顶部导航栏 -->
    <div
      :class="{ 'bg-white/90 dark:bg-neutral-800/90 backdrop-blur-md z-20 shadow-sm': y > 50 }"
      class="flex fixed justify-between items-center px-4 w-full md:w-[567px] top-0 transition-all duration-200"
      :style="{ paddingTop: 'var(--safe-area-top, 0px)' }"
    >
      <!-- 左侧：返回按钮 -->
      <button class="w-8 h-8 rounded-full hover:bg-white/20 dark:hover:bg-gray-700 flex items-center justify-center transition-colors" title="返回主页" @click="navigateTo('/')">
        <UIcon name="i-carbon-chevron-left" class="w-5 h-5 cursor-pointer text-white dark:text-gray-300" />
      </button>

      <!-- 右侧：功能按钮 -->
      <div class="flex items-center gap-3">
        <!-- 刷新按钮 -->
        <button class="w-8 h-8 rounded-full hover:bg-white/20 dark:hover:bg-gray-700 flex items-center justify-center transition-colors" title="刷新" @click="handleRefresh">
          <UIcon name="i-carbon-renew" class="w-5 h-5 text-white dark:text-gray-300" />
        </button>

        <!-- 日夜切换 -->
        <button
          class="w-8 h-8 rounded-full hover:bg-white/20 dark:hover:bg-gray-700 flex items-center justify-center transition-colors"
          title="切换主题"
          @click="toggleColorMode"
        >
          <UIcon
            :name="colorMode.value === 'dark' ? 'i-carbon-sun' : 'i-carbon-moon'"
            class="w-5 h-5 text-white dark:text-gray-300"
          />
        </button>

        <!-- 用户中心 -->
        <NuxtLink to="/user/settings" title="用户中心" class="w-8 h-8 rounded-full hover:bg-white/20 dark:hover:bg-gray-700 flex items-center justify-center transition-colors">
          <UIcon name="i-carbon-user" class="w-5 h-5 text-white dark:text-gray-300" />
        </NuxtLink>
      </div>
    </div>

    <!-- 封面图片 -->
    <img
      v-if="props.user.coverUrl"
      class="header-img w-full"
      :src="props.user.coverUrl"
      alt=""
    />
    <div v-else class="header-img w-full h-48 bg-gradient-to-br from-[#9fc84a] to-[#7ba428]" />
  </div>
</template>
<script setup lang="ts">
import type { UserVO } from "~/types";
import { useGlobalState } from "~/store";

const global = useGlobalState();
const route = useRoute();
const colorMode = useColorMode();

const props = defineProps<{ user: UserVO }>();
const { y } = useWindowScroll();

const toggleColorMode = () => {
  colorMode.preference = colorMode.value === 'dark' ? 'light' : 'dark';
};

const handleRefresh = () => {
  window.location.reload();
};
</script>

<style scoped></style>
