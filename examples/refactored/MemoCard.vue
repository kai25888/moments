<template>
  <article 
    class="memo-card"
    :class="{ 
      'memo-card--pinned': memo.pinned,
      'memo-card--private': memo.showType === MemoShowType.Private 
    }"
  >
    <!-- 头部：用户信息 -->
    <header class="memo-card__header">
      <UserAvatar 
        :src="memo.user.avatarUrl" 
        :alt="memo.user.nickname"
        size="md"
      />
      <div class="memo-card__meta">
        <NuxtLink 
          :to="`/user/${memo.user.id}`"
          class="memo-card__username"
        >
          {{ memo.user.nickname }}
        </NuxtLink>
        <time class="memo-card__time" :datetime="memo.createdAt">
          {{ formattedTime }}
        </time>
      </div>
      <div class="memo-card__badges">
        <UBadge v-if="memo.pinned" color="primary" size="xs">
          <UIcon name="i-carbon-pin" class="mr-1" />
          置顶
        </UBadge>
        <UBadge v-if="memo.showType === MemoShowType.Private" color="red" size="xs">
          <UIcon name="i-carbon-locked" class="mr-1" />
          私密
        </UBadge>
      </div>
    </header>

    <!-- 内容区 -->
    <main class="memo-card__content">
      <div 
        ref="contentRef"
        class="memo-card__text"
        :class="{ 'memo-card__text--collapsed': isCollapsed }"
        v-html="renderedContent"
      />
      
      <UButton
        v-if="showExpandButton"
        variant="ghost"
        size="xs"
        @click="toggleExpand"
      >
        {{ isCollapsed ? '展开全文' : '收起' }}
      </UButton>

      <!-- 标签 -->
      <div v-if="tags.length > 0" class="memo-card__tags">
        <NuxtLink
          v-for="tag in tags"
          :key="tag"
          :to="`/tags/${memo.user.username}/${tag}`"
        >
          <UBadge size="xs" variant="soft">
            #{{ tag }}
          </UBadge>
        </NuxtLink>
      </div>
    </main>

    <!-- 媒体内容 -->
    <section class="memo-card__media">
      <ExternalLinkPreview
        v-if="externalLink"
        v-bind="externalLink"
      />
      
      <ImageGallery
        v-if="hasImages"
        :images="imageConfigs"
        :memo-id="memo.id"
      />

      <MusicPlayer
        v-if="media.music"
        v-bind="media.music"
      />

      <DoubanBookCard
        v-if="media.doubanBook"
        :book="media.doubanBook"
      />

      <DoubanMovieCard
        v-if="media.doubanMovie"
        :movie="media.doubanMovie"
      />

      <VideoPlayer
        v-if="media.video"
        :video="media.video"
      />
    </section>

    <!-- 位置信息 -->
    <div v-if="location" class="memo-card__location">
      <UIcon name="i-carbon-location" />
      <span>{{ location }}</span>
    </div>

    <!-- 操作栏 -->
    <footer class="memo-card__footer">
      <div class="memo-card__actions">
        <UButton
          variant="ghost"
          size="xs"
          :color="isLiked ? 'red' : 'gray'"
          @click="handleLike"
        >
          <UIcon name="i-carbon-favorite" :class="{ 'text-red-500': isLiked }" />
          <span v-if="memo.favCount > 0">{{ memo.favCount }}</span>
        </UButton>

        <UButton
          v-if="config.enableComment"
          variant="ghost"
          size="xs"
          @click="toggleComment"
        >
          <UIcon name="i-octicon-comment" />
          <span v-if="memo.commentCount > 0">{{ memo.commentCount }}</span>
        </UButton>

        <UDropdown :items="actionItems">
          <UButton variant="ghost" size="xs">
            <UIcon name="i-carbon-overflow-menu-vertical" />
          </UButton>
        </UDropdown>
      </div>
    </footer>

    <!-- 评论区 -->
    <section v-if="showComments" class="memo-card__comments">
      <CommentList
        :memo-id="memo.id"
        :comments="memo.comments"
        :memo-user-id="memo.user.id"
      />
    </section>
  </article>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useElementSize } from '@vueuse/core'
import type { MemoVO, SysConfigVO, ExtDTO } from '~/types'
import { useMemoStore } from '~/stores/memo'
import { useUserStore } from '~/stores/user'
import { useToast } from '~/composables/useToast'
import { parseTags, formatRelativeTime, renderMarkdown } from '~/utils'

// 常量定义
const MAX_CONTENT_HEIGHT = 300
const MemoShowType = {
  Private: 0,
  Public: 1
} as const

// Props 定义
interface Props {
  memo: MemoVO
}

const props = defineProps<Props>()

// Emits 定义
const emit = defineEmits<{
  (e: 'update', memo: MemoVO): void
  (e: 'delete', id: number): void
}>()

// 依赖注入
const memoStore = useMemoStore()
const userStore = useUserStore()
const toast = useToast()
const route = useRoute()

// Refs
const contentRef = ref<HTMLElement | null>(null)
const isCollapsed = ref(false)
const isLiked = ref(false)
const showComments = ref(false)

// 计算属性
const isDetailPage = computed(() => route.path.startsWith('/memo/'))
const isAuthor = computed(() => userStore.currentUser?.id === props.memo.user.id)
const isAdmin = computed(() => userStore.currentUser?.id === 1)
const canEdit = computed(() => isAuthor.value || isAdmin.value)

const formattedTime = computed(() => {
  return formatRelativeTime(props.memo.createdAt)
})

const renderedContent = computed(() => {
  if (!props.memo.content) return ''
  try {
    return renderMarkdown(props.memo.content)
  } catch (error) {
    console.error('Failed to render markdown:', error)
    return props.memo.content
  }
})

const showExpandButton = computed(() => {
  if (isDetailPage.value) return false
  // 实际高度检测在 onMounted 中完成
  return false // 简化示例
})

const tags = computed(() => {
  return parseTags(props.memo.tags)
})

const location = computed(() => {
  return props.memo.location?.replace(/\s+/g, ' · ') || ''
})

const externalLink = computed(() => {
  if (props.memo.externalUrl && props.memo.externalTitle) {
    return {
      url: props.memo.externalUrl,
      title: props.memo.externalTitle,
      favicon: props.memo.externalFavicon
    }
  }
  return null
})

const hasImages = computed(() => {
  return props.memo.imgConfigs && props.memo.imgConfigs.length > 0
})

const imageConfigs = computed(() => {
  return props.memo.imgConfigs || []
})

const media = computed<ExtDTO>(() => {
  try {
    return JSON.parse(props.memo.ext || '{}') as ExtDTO
  } catch {
    return {} as ExtDTO
  }
})

const config = computed<SysConfigVO>(() => {
  return useState<SysConfigVO>('sysConfig').value
})

const actionItems = computed(() => {
  const items = []
  
  if (isAdmin.value) {
    items.push({
      label: props.memo.pinned ? '取消置顶' : '置顶',
      icon: 'i-carbon-pin',
      click: handlePin
    })
  }
  
  if (isAuthor.value) {
    items.push({
      label: '编辑',
      icon: 'i-carbon-edit',
      click: handleEdit
    })
  }
  
  if (canEdit.value) {
    items.push({
      label: '删除',
      icon: 'i-carbon-trash-can',
      click: handleDelete
    })
  }
  
  return [items]
})

// 方法
function toggleExpand() {
  isCollapsed.value = !isCollapsed.value
}

async function handleLike() {
  if (isLiked.value) {
    toast.warning('您已经点赞过了')
    return
  }
  
  try {
    await memoStore.like(props.memo.id)
    isLiked.value = true
    toast.success('点赞成功')
    emit('update', { ...props.memo, favCount: props.memo.favCount + 1 })
  } catch (error) {
    toast.error('点赞失败')
  }
}

function toggleComment() {
  showComments.value = !showComments.value
}

async function handlePin() {
  try {
    await memoStore.setPinned(props.memo.id, !props.memo.pinned)
    toast.success(props.memo.pinned ? '已取消置顶' : '已置顶')
    emit('update', { ...props.memo, pinned: !props.memo.pinned })
  } catch (error) {
    toast.error('操作失败')
  }
}

async function handleEdit() {
  await navigateTo(`/edit/${props.memo.id}`)
}

async function handleDelete() {
  if (!confirm('确定要删除这条动态吗？')) return
  
  try {
    await memoStore.delete(props.memo.id)
    toast.success('删除成功')
    emit('delete', props.memo.id)
  } catch (error) {
    toast.error('删除失败')
  }
}

// 生命周期
onMounted(() => {
  // 检查是否已点赞
  const likedMemos = JSON.parse(localStorage.getItem('likedMemos') || '[]') as number[]
  isLiked.value = likedMemos.includes(props.memo.id)
  
  // 检测内容高度决定是否显示展开按钮
  if (contentRef.value && !isDetailPage.value) {
    const { height } = useElementSize(contentRef)
    if (height.value > MAX_CONTENT_HEIGHT) {
      isCollapsed.value = true
    }
  }
})
</script>

<style scoped>
.memo-card {
  @apply bg-white dark:bg-neutral-800 rounded-lg p-4 shadow-sm;
  @apply border border-gray-100 dark:border-neutral-700;
}

.memo-card--pinned {
  @apply bg-slate-50 dark:bg-neutral-700;
}

.memo-card--private {
  @apply border-red-200 dark:border-red-800;
}

.memo-card__header {
  @apply flex items-center gap-3 mb-3;
}

.memo-card__meta {
  @apply flex flex-col flex-1;
}

.memo-card__username {
  @apply font-medium text-gray-900 dark:text-white hover:text-primary-500;
}

.memo-card__time {
  @apply text-xs text-gray-500 dark:text-gray-400;
}

.memo-card__badges {
  @apply flex gap-2;
}

.memo-card__content {
  @apply mb-3;
}

.memo-card__text {
  @apply prose dark:prose-invert max-w-none;
}

.memo-card__text--collapsed {
  @apply max-h-[300px] overflow-hidden;
}

.memo-card__tags {
  @apply flex flex-wrap gap-2 mt-2;
}

.memo-card__media {
  @apply space-y-2 mb-3;
}

.memo-card__location {
  @apply flex items-center gap-1 text-xs text-gray-500 mb-3;
}

.memo-card__footer {
  @apply flex items-center justify-between;
}

.memo-card__actions {
  @apply flex items-center gap-1;
}

.memo-card__comments {
  @apply mt-4 pt-4 border-t border-gray-100 dark:border-neutral-700;
}
</style>
