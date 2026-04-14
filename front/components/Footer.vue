<template>
  <div class="footer px-4 pb-8 pt-6">
    <div class="rounded-3xl border border-black/5 bg-white/70 px-5 py-4 text-xs text-slate-500 shadow-sm backdrop-blur dark:border-white/10 dark:bg-neutral-900/70 dark:text-slate-400">
      <div class="flex flex-col gap-3 items-center text-center">
        <div class="space-y-1">
          <div class="font-medium text-slate-700 dark:text-slate-200">
            {{ sysConfig.title || "Moments" }}
          </div>
          <div v-if="sysConfig.footerSignature || sysConfig.siteDescription">
            {{ sysConfig.footerSignature || sysConfig.siteDescription }}
          </div>
        </div>
        <div class="flex items-center gap-3 text-[11px] uppercase tracking-[0.24em] text-slate-400 dark:text-slate-500">
          <NuxtLink :to="rssLink" target="_blank">
            RSS
          </NuxtLink>
          <span v-if="sysConfig.beiAnNo">·</span>
          <a
            v-if="sysConfig.beiAnNo"
            href="https://beian.miit.gov.cn/"
            target="_blank"
          >
            {{ sysConfig.beiAnNo }}
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { SysConfigVO } from "~/types";

const sysConfig = useState<SysConfigVO>("sysConfig");
const rssLink = computed(() => sysConfig.value.rss || "/rss");
</script>

<style scoped></style>
