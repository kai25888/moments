<template>
  <div class="min-h-screen bg-[#f4efe7] px-0 py-0 text-slate-900 transition-colors dark:bg-[#111111] dark:text-slate-100">
    <div class="pointer-events-none fixed inset-0 bg-[radial-gradient(circle_at_top,_rgba(214,190,153,0.28),_transparent_38%),radial-gradient(circle_at_bottom,_rgba(159,200,74,0.12),_transparent_30%)] dark:bg-[radial-gradient(circle_at_top,_rgba(159,200,74,0.16),_transparent_35%),radial-gradient(circle_at_bottom,_rgba(255,255,255,0.05),_transparent_30%)]"></div>
    <div
      class="relative w-full md:w-[567px] mx-auto min-h-screen border-x border-black/5 bg-white/90 shadow-[0_25px_80px_rgba(15,23,42,0.08)] backdrop-blur dark:border-white/10 dark:bg-neutral-900/90"
    >
      <slot />
      <Footer />
    </div>
  </div>

  <div
    title="到顶部"
    v-if="y > 200"
    @click="y = 0"
    class="hidden sm:block bottom-[20%] sm:right-[20%] md:right-[10%] lg:right-[15%] xl:right-[20%] 2xl:right-[28%] fixed flex items-center justify-center"
  >
    <UIcon
      name="i-lets-icons-expand-top-stop"
      class="w-10 h-10 text-gray-500 cursor-pointer"
    ></UIcon>
  </div>

  <div class="sm:hidden relative">
    <div class="right-0 bottom-10 fixed flex items-center justify-end">
      <div class="flex flex-col items-center gap-2">
        <div
          v-if="y > 300"
          @click="y = 0"
          class="dark:bg-gray-900/85 mr-4 rounded-full bg-slate-50 w-10 h-10 flex items-center justify-center shadow-xl"
        >
          <UIcon
            name="i-lets-icons-expand-top-stop"
            class="w-6 h-6 text-[#9fc84a] cursor-pointer"
          ></UIcon>
        </div>
        <NuxtLink
          to="/new"
          v-if="global.userinfo.token && $route.path === '/'"
          class="dark:bg-gray-900/85 mr-4 rounded-full bg-slate-50 w-10 h-10 flex items-center justify-center shadow-xl"
        >
          <UIcon name="i-carbon-camera" class="w-6 h-6 text-[#9fc84a]"></UIcon>
        </NuxtLink>
        <div
          class="dark:bg-gray-900/85 mr-4 rounded-full bg-slate-50 w-10 h-10 flex items-center justify-center shadow-xl"
          @click="open = true"
        >
          <UIcon
            name="i-icon-park-solid-more-four"
            class="w-6 h-6 text-[#9fc84a] cursor-pointer"
          ></UIcon>
        </div>
        <NuxtLink
          to="/user/login"
          v-if="!global.userinfo.token && $route.path === '/'"
          class="dark:bg-gray-900/85 mr-4 rounded-full bg-slate-50 w-10 h-10 flex items-center justify-center shadow-xl"
        >
          <UIcon name="i-carbon-login" class="w-6 h-6 text-[#9fc84a]"></UIcon>
        </NuxtLink>
      </div>
    </div>

    <MobileNav :open="open" />
  </div>
</template>

<script lang="ts" setup>
import type { SysConfigVO, UserVO } from "~/types";
import { useGlobalState } from "~/store";

const global = useGlobalState();
const open = useState<boolean>("sidebarOpen", () => false);
const currentUser = useState<UserVO>("userinfo");
const sysConfig = useState<SysConfigVO>("sysConfig");
const currentProfile = await useMyFetch<UserVO>("/user/profile");
const sysConfigVO = await useMyFetch<SysConfigVO>("/sysConfig/get");
sysConfig.value = sysConfigVO;
if (currentProfile) {
  currentUser.value = currentProfile;
}
const { y } = useWindowScroll();
useHead({
  title: sysConfigVO.title,
  meta: [
    {
      name: "description",
      content:
        sysConfigVO.siteDescription ||
        sysConfigVO.heroSubtitle ||
        `${sysConfigVO.title} 的轻量时间流`,
    },
  ],
  link: [
    {
      rel: "shortcut icon",
      type: "image/png",
      href: sysConfigVO.favicon || "/favicon.png",
    },
    {
      rel: "apple-touch-icon-precomposed",
      href: sysConfigVO.favicon || "/favicon.png",
    },
    {
      rel: "alternate",
      type: "application/rss+xml",
      title: "我的 RSS 订阅",
      href: sysConfigVO.rss || `/rss`,
    },
  ],
  bodyAttrs: {
    class: "bg-[#f4efe7] dark:bg-[#111111]",
  },
  style: [
    {
      innerHTML: sysConfigVO.css || "",
    },
  ],
});

if (sysConfigVO.enableGoogleRecaptcha && sysConfigVO.googleSiteKey) {
  useHead({
    script: [
      {
        type: "text/javascript",
        src: `https://recaptcha.net/recaptcha/api.js?render=${sysConfigVO.googleSiteKey}`,
      },
    ],
  });
}
</script>
