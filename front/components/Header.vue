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
      <!-- 左侧：返回 + 标题 -->
      <NuxtLink class="flex items-center" title="返回主页">
        <UIcon
          @click="navigateTo('/')"
          name="i-carbon-chevron-left"
          class="w-5 h-5 cursor-pointer mr-3 text-gray-600 dark:text-gray-300"
        />
        <span class="text-sm text-gray-600 dark:text-gray-300" v-if="$route.path === '/user/calendar'">日历检索</span>
        <span class="text-sm text-gray-600 dark:text-gray-300" v-else-if="$route.path === '/sys/settings'">系统设置</span>
        <span class="text-sm text-gray-600 dark:text-gray-300" v-else-if="$route.path === '/user/settings'">用户中心</span>
        <span class="text-sm text-gray-600 dark:text-gray-300" v-else-if="$route.path.indexOf('/tags/') >= 0">
          {{ route.params.tag || "话题专栏" }}
        </span>
        <span class="text-sm text-gray-600 dark:text-gray-300" v-else-if="$route.path === '/friend'">友情链接</span>
        <span class="text-sm text-gray-600 dark:text-gray-300" v-else>
          <span v-if="!global.userinfo.token && $route.path === '/user/login'">登录</span>
          <span v-else-if="!global.userinfo.token && $route.path === '/user/reg'">注册</span>
          <span v-else>{{ props.user.nickname }}</span>
        </span>
      </NuxtLink>

      <!-- 右侧：用户头像 -->
      <div class="flex items-center gap-2">
        <!-- 登录状态：点击头像进入用户中心 -->
        <NuxtLink v-if="global.userinfo.token" to="/user/settings" title="用户中心">
          <img
            v-if="props.user.avatarUrl"
            :src="props.user.avatarUrl"
            class="w-8 h-8 rounded-full ring-2 ring-white/50 shadow-md hover:ring-[#9fc84a] transition-all"
          />
          <div
            v-else
            class="w-8 h-8 rounded-full bg-gray-200 dark:bg-gray-600 flex items-center justify-center hover:bg-[#9fc84a]/20 transition-all"
          >
            <UIcon name="i-carbon-person" class="w-5 h-5 text-gray-400" />
          </div>
        </NuxtLink>

        <!-- 未登录：显示登录按钮 -->
        <NuxtLink
          v-else
          to="/user/login"
          class="text-sm text-[#9fc84a] hover:text-[#7ba428] transition-colors"
        >
          登录
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

const props = defineProps<{ user: UserVO }>();
const { y } = useWindowScroll();
</script>

<style scoped></style>
