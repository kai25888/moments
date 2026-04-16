<template>
  <div
    v-if="$route.path !== '/new' && $route.path.indexOf('/edit/') < 0"
    class="header relative mb-14"
  >
    <!-- 顶部导航栏 -->
    <div
      :class="{ 'bg-white/90 dark:bg-neutral-800/90 backdrop-blur-md z-20 shadow-sm': y > 50 }"
      class="flex fixed justify-between items-center px-4 w-full md:w-[567px] top-0 transition-all duration-200"
    >
      <!-- 左侧：返回按钮 -->
      <NuxtLink class="flex items-center" title="返回主页">
        <UIcon
          @click="navigateTo('/')"
          name="i-carbon-chevron-left"
          class="w-5 h-5 cursor-pointer text-gray-600 dark:text-gray-300"
        />
      </NuxtLink>

      <!-- 右侧：功能按钮 -->
      <div class="flex items-center gap-3">
        <!-- 刷新按钮 -->
        <UIcon
          name="i-carbon-renew"
          class="w-5 h-5 cursor-pointer text-gray-600 dark:text-gray-300 hover:text-[#9fc84a] transition-colors"
          title="刷新"
          @click="() => location.reload()"
        />

        <!-- 日夜切换 -->
        <button
          class="w-8 h-8 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 flex items-center justify-center transition-colors"
          title="切换主题"
          @click="toggleColorMode"
        >
          <UIcon
            :name="colorMode.value === 'dark' ? 'i-carbon-sun' : 'i-carbon-moon'"
            class="w-5 h-5 text-gray-600 dark:text-gray-300"
          />
        </button>

        <!-- 用户中心 -->
        <NuxtLink to="/user/settings" title="用户中心">
          <div class="w-8 h-8 rounded-full bg-[#9fc84a]/20 flex items-center justify-center hover:bg-[#9fc84a]/30 transition-all">
            <UIcon name="i-carbon-user" class="w-5 h-5 text-[#9fc84a]" />
          </div>
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
    <!-- 用户信息 -->
    <div class="absolute right-4 bottom-[-40px]">
      <div class="flex flex-col items-end">
        <div class="flex flex-row items-center gap-3">
          <span class="text-lg font-bold text-white drop-shadow-md">
            {{ props.user.nickname }}
          </span>
          <img
            v-if="props.user.avatarUrl"
            :src="props.user.avatarUrl"
            class="avatar w-[70px] h-[70px] rounded-xl ring-2 ring-white/50 shadow-lg"
          />
          <div
            v-else
            class="avatar w-[70px] h-[70px] rounded-xl bg-gray-300 dark:bg-gray-600 flex items-center justify-center"
          >
            <UIcon name="i-carbon-person" class="w-10 h-10 text-gray-400" />
          </div>
        </div>
      </div>
    </div>
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
</script>

<style scoped></style>
