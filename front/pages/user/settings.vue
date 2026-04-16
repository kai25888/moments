<template>
  <Header :user="currentUser"/>

  <div class="settings-container p-4 max-w-2xl mx-auto">
    <!-- 用户资料卡片 -->
    <UCard class="mb-4">
      <template #header>
        <div class="flex items-center gap-2">
          <UIcon name="i-carbon-user" class="w-5 h-5 text-[#9fc84a]" />
          <span class="font-bold">个人资料</span>
        </div>
      </template>

      <!-- 头像 -->
      <UFormGroup label="头像" :ui="{ label: { base: 'font-medium text-gray-700 dark:text-gray-200' } }">
        <div class="flex items-start gap-6">
          <div class="avatar-preview w-24 h-24 rounded-2xl overflow-hidden ring-2 ring-[#9fc84a]/30">
            <img v-if="state.avatarUrl" :src="state.avatarUrl" class="w-full h-full object-cover" alt="头像" />
            <div v-else class="w-full h-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
              <UIcon name="i-carbon-person" class="w-12 h-12 text-gray-400" />
            </div>
          </div>
          <div class="flex-1">
            <UInput type="file" size="sm" accept="image/*" @change="uploadAvatarUrl" class="mb-3" />
            <p class="text-xs text-gray-500 mb-2">或者输入在线地址</p>
            <UInput v-model="state.avatarUrl" placeholder="https://..." />
          </div>
        </div>
      </UFormGroup>

      <!-- 封面图片 -->
      <UFormGroup label="封面图片" class="mt-6" :ui="{ label: { base: 'font-medium text-gray-700 dark:text-gray-200' } }">
        <div class="cover-preview rounded-xl overflow-hidden h-32 bg-gray-100 dark:bg-gray-700 mb-3">
          <img v-if="state.coverUrl" :src="state.coverUrl" class="w-full h-full object-cover" alt="封面" />
          <div v-else class="w-full h-full flex items-center justify-center">
            <UIcon name="i-carbon-image" class="w-8 h-8 text-gray-400" />
          </div>
        </div>
        <UInput type="file" size="sm" accept="image/*" @change="uploadCoverUrl" class="mb-2" />
        <p class="text-xs text-gray-500">或者输入在线地址</p>
        <UInput v-model="state.coverUrl" placeholder="https://..." class="mt-2" />
      </UFormGroup>
    </UCard>

    <!-- 账号信息卡片 -->
    <UCard class="mb-4">
      <template #header>
        <div class="flex items-center gap-2">
          <UIcon name="i-carbon-id" class="w-5 h-5 text-[#9fc84a]" />
          <span class="font-bold">账号信息</span>
        </div>
      </template>

      <UFormGroup label="登录名" :ui="{ label: { base: 'font-medium text-gray-700 dark:text-gray-200' } }">
        <UInput v-model="state.username" disabled />
        <template #help>
          <span class="text-xs text-gray-400">登录名不可修改</span>
        </template>
      </UFormGroup>

      <UFormGroup label="昵称" class="mt-4" :ui="{ label: { base: 'font-medium text-gray-700 dark:text-gray-200' } }">
        <UInput v-model="state.nickname" placeholder="给自己起个昵称" />
      </UFormGroup>

      <UFormGroup label="心情状态" class="mt-4" :ui="{ label: { base: 'font-medium text-gray-700 dark:text-gray-200' } }">
        <UTextarea v-model="state.slogan" placeholder="签名/心情/状态..." :rows="2" />
      </UFormGroup>

      <UFormGroup label="邮箱" class="mt-4" :ui="{ label: { base: 'font-medium text-gray-700 dark:text-gray-200' } }">
        <UInput v-model="state.email" type="email" placeholder="your@email.com" />
        <template #help>
          <span class="text-xs text-gray-400">若管理员启用了邮件通知，将在收到评论时发送邮件通知</span>
        </template>
      </UFormGroup>
    </UCard>

    <!-- 安全设置卡片 -->
    <UCard class="mb-6">
      <template #header>
        <div class="flex items-center gap-2">
          <UIcon name="i-carbon-password" class="w-5 h-5 text-[#9fc84a]" />
          <span class="font-bold">安全设置</span>
        </div>
      </template>

      <UFormGroup label="修改密码" :ui="{ label: { base: 'font-medium text-gray-700 dark:text-gray-200' } }">
        <UInput v-model="state.password" type="password" placeholder="留空则不修改密码" />
        <template #help>
          <span class="text-xs text-gray-400">建议使用强密码，包含字母、数字和特殊字符</span>
        </template>
      </UFormGroup>
    </UCard>

    <!-- 保存按钮 -->
    <div class="flex gap-3">
      <UButton size="lg" class="flex-1 justify-center" @click="save">
        <UIcon name="i-carbon-save" class="w-4 h-4 mr-2" />
        保存修改
      </UButton>
      <UButton size="lg" variant="outline" color="red" @click="logout">
        <UIcon name="i-carbon-logout" class="w-4 h-4 mr-2" />
        退出登录
      </UButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { UserVO } from "~/types";
import { toast } from "vue-sonner";
import { useUpload } from "~/utils";

definePageMeta({ middleware: 'auth' })
import { useGlobalState } from "~/store";

const global = useGlobalState()
const currentUser = useState<UserVO>('userinfo')

const state = reactive({
  password: "",
  username: "",
  nickname: "",
  slogan: "",
  avatarUrl: "",
  coverUrl: "",
  email: "",
})

const logout = async () => {
  global.value.userinfo = {}
  await navigateTo('/')
}

const save = async () => {
  await useMyFetch('/user/saveProfile', state)
  toast.success("保存成功")
  location.reload()
}

const uploadAvatarUrl = async (files: FileList) => {
  for (let i = 0; i < files.length; i++) {
    if (files[i].type.indexOf("image") < 0) {
      toast.error("只能上传图片");
      return
    }
  }
  const result = await useUpload(files)
  if (result.length) {
    toast.success("上传成功")
    state.avatarUrl = result[0]
  }
}

const uploadCoverUrl = async (files: FileList) => {
  for (let i = 0; i < files.length; i++) {
    if (files[i].type.indexOf("image") < 0) {
      toast.error("只能上传图片");
      return
    }
  }
  const result = await useUpload(files)
  if (result.length) {
    toast.success("上传成功")
    state.coverUrl = result[0]
  }
}

onMounted(async () => {
  Object.assign(state, currentUser.value)
})
</script>

<style scoped>
.settings-container {
  padding-bottom: 2rem;
}
</style>
