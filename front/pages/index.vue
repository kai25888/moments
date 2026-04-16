<template>
  <Header v-bind:user="currentUser"/>
  <section class="px-4 pb-4">
    <div class="overflow-hidden rounded-[28px] border border-black/5 bg-[linear-gradient(135deg,rgba(247,241,231,0.96),rgba(255,255,255,0.92))] p-5 shadow-[0_18px_60px_rgba(15,23,42,0.08)] dark:border-white/10 dark:bg-[linear-gradient(135deg,rgba(25,25,25,0.96),rgba(16,16,16,0.92))]">
      <div class="flex flex-col gap-5">
        <div v-if="heroAnnouncement" class="inline-flex w-fit items-center gap-2 rounded-full border border-emerald-200 bg-emerald-50 px-3 py-1 text-[11px] font-medium text-emerald-700 dark:border-emerald-500/20 dark:bg-emerald-500/10 dark:text-emerald-200">
          <UIcon name="i-carbon-notification" class="h-3.5 w-3.5" />
          <span>{{ heroAnnouncement }}</span>
        </div>

        <div class="flex items-start gap-4">
          <!-- 用户头像 -->
          <img
            v-if="currentUser.avatarUrl"
            :src="currentUser.avatarUrl"
            class="w-14 h-14 rounded-2xl ring-2 ring-white/50 shadow-md shrink-0"
          />
          <div
            v-else
            class="w-14 h-14 rounded-2xl bg-gray-200 dark:bg-gray-600 flex items-center justify-center shrink-0"
          >
            <UIcon name="i-carbon-person" class="w-8 h-8 text-gray-400" />
          </div>

          <div class="space-y-1 flex-1">
            <div class="text-[11px] uppercase tracking-[0.35em] text-slate-400 dark:text-slate-500">
              {{ sysConfig.title || currentUser.nickname || "Moments" }}
            </div>
            <div class="text-2xl font-semibold leading-tight text-slate-900 dark:text-white">
              {{ heroTitle }}
            </div>
          </div>
          <div class="text-sm leading-7 text-slate-600 dark:text-slate-300">
            {{ heroSubtitle }}
          </div>
          <div v-if="heroDescription" class="max-w-[30rem] text-sm leading-7 text-slate-500 dark:text-slate-400">
            {{ heroDescription }}
          </div>
        </div>

        <div class="flex items-center gap-3">
          <div class="flex-1 rounded-2xl border border-black/5 bg-white/80 px-4 py-3 dark:border-white/10 dark:bg-white/5">
            <div class="text-[11px] uppercase tracking-[0.24em] text-slate-400 dark:text-slate-500">动态</div>
            <div class="mt-2 text-2xl font-semibold text-slate-900 dark:text-white">{{ total || 0 }}</div>
            <div class="mt-1 text-xs text-slate-500 dark:text-slate-400">
              {{ total > 0 ? "已发布内容" : "还没有发布内容" }}
            </div>
          </div>
          <NuxtLink
            v-if="global.userinfo.token"
            to="/new"
            class="flex-1 inline-flex items-center justify-center rounded-full bg-[#9fc84a] px-4 py-2 text-sm font-medium text-[#1f2a10] transition hover:brightness-95"
          >
            发表新动态
          </NuxtLink>
          <NuxtLink
            v-else-if="sysConfig.enableRegister"
            to="/user/reg"
            class="flex-1 inline-flex items-center justify-center rounded-full bg-[#9fc84a] px-4 py-2 text-sm font-medium text-[#1f2a10] transition hover:brightness-95"
          >
            注册并加入
          </NuxtLink>
          <NuxtLink
            v-if="!global.userinfo.token"
            to="/user/login"
            class="flex-1 inline-flex items-center justify-center rounded-full border border-transparent px-4 py-2 text-sm font-medium text-slate-500 transition hover:text-slate-800 dark:text-slate-400 dark:hover:text-slate-100"
          >
            管理登录
          </NuxtLink>
        </div>
      </div>
    </div>
  </section>

  <div class="mx-4 mb-3 flex items-center justify-between text-xs text-slate-400 dark:text-slate-500">
    <div>时间流</div>
    <div>{{ hasNext ? "持续更新中" : "已加载全部内容" }}</div>
  </div>

  <div v-if="loading" class="flex flex-col divide-y divide-[#C0BEBF]/20">
    <div v-for="i in 5" :key="i" class="bg-white dark:bg-neutral-800 p-4">
      <div class="flex gap-4">
        <USkeleton class="w-10 h-10 rounded-full shrink-0" />
        <div class="flex-1 space-y-2">
          <USkeleton class="h-4 w-24" />
          <USkeleton class="h-20 w-full" />
        </div>
      </div>
    </div>
  </div>

  <div v-else class="flex flex-col divide-y divide-[#C0BEBF]/20 ">
    <Memo v-bind:memo="m" v-for="m in memos" :key="m.id" />
  </div>
  <div ref="loadMoreEle" class="px-4 py-5 text-center text-xs text-gray-500" v-if="hasNext">
    <button class="w-full rounded-full border border-black/5 bg-white/80 px-4 py-3 transition hover:border-[#9fc84a] hover:text-slate-700 dark:border-white/10 dark:bg-white/5 dark:hover:text-white" @click="loadMore">
      点击加载更多
    </button>
  </div>
  <div class="px-4 py-5 text-center text-xs text-gray-500" v-else-if="!loading">已经到底啦</div>
</template>

<script setup lang="ts">
import type {MemoVO, SysConfigVO, UserVO} from "~/types";
import Memo from "~/components/Memo.vue";
import {memoChangedEvent, memoReloadEvent} from "~/event";
import {useElementVisibility} from '@vueuse/core'
import { useGlobalState } from "~/store";

const currentUser = useState<UserVO>('userinfo')
const sysConfig = useState<SysConfigVO>('sysConfig')
const global = useGlobalState()

const loadMoreEle = ref(null)
const targetIsVisible = useElementVisibility(loadMoreEle)
watch(targetIsVisible, async (visible) => {
  if (visible && sysConfig.value.enableAutoLoadNextPage) {
    await loadMore()
  }
})
const hasNext = ref(false)
const total = ref(0)
const loading = ref(true)
const state = reactive({
  page: 1,
  size: 5,
})

const memos = ref<Array<MemoVO>>([])
const heroTitle = computed(() => sysConfig.value.heroTitle || currentUser.value.nickname || sysConfig.value.title || '记录那些值得停留的瞬间')
const heroSubtitle = computed(() => sysConfig.value.heroSubtitle || currentUser.value.slogan || '这里保留日常、思考和热爱的碎片，用更舒展的方式展示生活流。')
const heroDescription = computed(() => sysConfig.value.siteDescription || '')
const heroAnnouncement = computed(() => sysConfig.value.announcement || '')

onMounted(async () => {
  await reload()
  loading.value = false
})

const reload = async () => {
  state.page = 1
  const res = await useMyFetch<{
    list: Array<MemoVO>,
    total: number,
    hasNext: boolean
  }>('/memo/list', state)
  memos.value = res.list
  total.value = res.total
  hasNext.value = res.hasNext
}

const loadMore = async () => {
  state.page = state.page + 1
  const res = await useMyFetch<{
    list: Array<MemoVO>,
    total: number,
    hasNext: boolean
  }>('/memo/list', state)
  memos.value = [...memos.value, ...res.list]
  total.value = res.total
  hasNext.value = res.hasNext
}

memoReloadEvent.on(async () => {
  await reload()
})

memoChangedEvent.on(async (id: number) => {
  const res = await useMyFetch<MemoVO>('/memo/get?latest=1&id=' + id)
  const index = memos.value.findIndex(r => r.id === id)
  if (index >= 0) {
    memos.value[index] = res
  }
})
</script>

<style scoped>

</style>