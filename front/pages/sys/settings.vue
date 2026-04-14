<template>
  <Header :user="currentUser"/>
  <div class="space-y-4  flex flex-col p-4 my-4 dark:bg-neutral-800">
    <div class="flex flex-col items-end text-xs text-gray-400">
      <div v-if="version" class="w-32">版本号: {{ version }}</div>
      <div v-if="commitId" class="w-32">commitId: {{ commitId }}</div>
    </div>
    <div class="rounded-2xl border border-emerald-200/70 bg-emerald-50/70 p-4 text-sm text-emerald-900 dark:border-emerald-500/20 dark:bg-emerald-500/10 dark:text-emerald-100">
      这一页现在支持直接配置首页品牌文案。适合做轻量二开，不需要额外改数据库结构。
    </div>
    <div class="space-y-1">
      <div class="text-sm font-semibold text-gray-900 dark:text-white">品牌基础</div>
      <div class="text-xs text-gray-500 dark:text-gray-400">站点标题、描述、首页主视觉文案都会影响访客第一眼感受。</div>
    </div>
    <UFormGroup label="管理员账号" name="adminUserName" :ui="{label:{base:'font-bold'}}">
      <UInput v-model="state.adminUserName"/>
    </UFormGroup>
    <UFormGroup label="网站标题" name="title" :ui="{label:{base:'font-bold'}}">
      <UInput v-model="state.title"/>
    </UFormGroup>
    <UFormGroup label="站点描述" name="siteDescription" :ui="{label:{base:'font-bold'}}">
      <UTextarea v-model="state.siteDescription" :rows="2" placeholder="一句话介绍你的站点定位、写作主题或社区气质"/>
    </UFormGroup>
    <UFormGroup label="首页主标题" name="heroTitle" :ui="{label:{base:'font-bold'}}">
      <UInput v-model="state.heroTitle" placeholder="例如：记录那些值得停留的瞬间"/>
    </UFormGroup>
    <UFormGroup label="首页副标题" name="heroSubtitle" :ui="{label:{base:'font-bold'}}">
      <UTextarea v-model="state.heroSubtitle" :rows="2" placeholder="补充说明站点内容方向、更新节奏或品牌态度"/>
    </UFormGroup>
    <UFormGroup label="顶部公告" name="announcement" :ui="{label:{base:'font-bold'}}">
      <UTextarea v-model="state.announcement" :rows="2" placeholder="可填写近期活动、欢迎语或站点提示"/>
    </UFormGroup>
    <UFormGroup label="Favicon" name="favicon"
                :ui="{label:{base:'font-bold'}}">
      <UInput type="file" size="sm" icon="i-heroicons-folder" accept="image/*" @change="uploadFavicon"/>
      <div class="text-gray-500 text-sm my-2">或者输入在线地址</div>
      <UInput v-model="state.favicon" class="mb-2"/>
      <UAvatar :src="state.favicon"/>
    </UFormGroup>
    <UFormGroup label="页脚签名" name="footerSignature" :ui="{label:{base:'font-bold'}}">
      <UInput v-model="state.footerSignature" placeholder="例如：用慢一点的节奏，写下真实的生活"/>
    </UFormGroup>
    <div class="space-y-1 pt-2">
      <div class="text-sm font-semibold text-gray-900 dark:text-white">站点行为</div>
      <div class="text-xs text-gray-500 dark:text-gray-400">这里控制首页加载、评论、注册和展示细节。</div>
    </div>
    <UFormGroup label="首页是否自动加载下一页" name="enableAutoLoadNextPage" :ui="{label:{base:'font-bold'}}">
      <UToggle v-model="state.enableAutoLoadNextPage"/>
    </UFormGroup>
    <UFormGroup label="是否启用评论" name="enableComment" :ui="{label:{base:'font-bold'}}">
      <UToggle v-model="state.enableComment"/>
    </UFormGroup>
    <UFormGroup label="是否开启注册用户" name="enableRegister" :ui="{label:{base:'font-bold'}}">
      <UToggle v-model="state.enableRegister"/>
    </UFormGroup>
    <UFormGroup label="备案号" name="beiAnNo" :ui="{label:{base:'font-bold'}}">
      <UInput v-model="state.beiAnNo" placeholder="没有可以不填写"/>
    </UFormGroup>
    <UFormGroup label="自定义CSS" name="css" :ui="{label:{base:'font-bold'}}">
      <UTextarea v-model="state.css" :rows="5"/>
    </UFormGroup>
    <UFormGroup label="自定义JS（出于安全考虑已停用前台注入）" name="js" :ui="{label:{base:'font-bold'}}">
      <UTextarea v-model="state.js" :rows="5" disabled placeholder="该字段保留历史兼容，但不会再自动注入到前台页面"/>
    </UFormGroup>
    <UFormGroup label="自定义RSS" name="rss" :ui="{label:{base:'font-bold'}}">
      <UTextarea v-model="state.rss" :rows="1"  placeholder="留空使用默认配置"/>
    </UFormGroup>
    <UFormGroup label="评论最大字数" name="maxCommentLength" :ui="{label:{base:'font-bold'}}">
      <UInput v-model.number="state.maxCommentLength"/>
    </UFormGroup>
    <UFormGroup label="发言最大高度(单位px,填0时则不限制高度)" name="memoMaxHeight" :ui="{label:{base:'font-bold'}}">
      <UInput v-model.number="state.memoMaxHeight"/>
    </UFormGroup>
    <UFormGroup label="评论排序方式(按日期)" name="commentOrder" :ui="{label:{base:'font-bold'}}">
      <USelectMenu v-model="state.commentOrder"
                   :options="[{label:'倒序,越晚发布越靠前',value:'desc'},{label:'正序,越早发布越靠前',value:'asc'}]"
                   value-attribute="value" option-attribute="label"></USelectMenu>
    </UFormGroup>
    <UFormGroup label="日期格式" name="timeFormat" :ui="{label:{base:'font-bold'}}">
      <USelectMenu v-model="state.timeFormat"
                   :options="[{label:'几分钟前',value:'timeAgo'},{label:$dayjs().format('YYYY-MM-DD HH:mm'),value:'time'}]"
                   value-attribute="value" option-attribute="label"></USelectMenu>
    </UFormGroup>
      <UFormGroup label="是否启用Google Recaptcha" name="enableGoogleRecaptcha" :ui="{label:{base:'font-bold'}}">
        <UToggle v-model="state.enableGoogleRecaptcha"/>
      </UFormGroup>
      <template v-if="state.enableGoogleRecaptcha">
        <UFormGroup label="SiteKey" name="googleSiteKey" :ui="{label:{base:'font-bold'}}">
          <UInput v-model="state.googleSiteKey"/>
        </UFormGroup>
        <UFormGroup label="SecretKey" name="googleSecretKey" :ui="{label:{base:'font-bold'}}">
          <UInput v-model="state.googleSecretKey"/>
        </UFormGroup>
      </template>
      <UFormGroup label="是否启用S3存储" name="s3" :ui="{label:{base:'font-bold'}}">
        <UToggle v-model="state.enableS3"/>
      </UFormGroup>
      <template v-if="state.enableS3">
        <UFormGroup label="Bucket 域名（资源访问地址）" name="domain" :ui="{label:{base:'font-bold'}}">
          <UInput v-model="state.s3.domain" placeholder="https://moments-test-bucket.oss-cn-hangzhou.aliyuncs.com" />
        </UFormGroup>
        <UFormGroup label="Endpoint 地址" name="endpoint" :ui="{label:{base:'font-bold'}}">
          <UInput v-model="state.s3.endpoint" placeholder="https://oss-cn-hangzhou.aliyuncs.com" />
        </UFormGroup>
        <UFormGroup label="Bucket 名称" name="bucket" :ui="{label:{base:'font-bold'}}">
          <UInput v-model="state.s3.bucket" placeholder="moments-test-bucket" />
        </UFormGroup>
        <UFormGroup label="Bucket 地区" name="region" :ui="{label:{base:'font-bold'}}">
          <UInput v-model="state.s3.region" placeholder="oss-cn-hangzhou" />
        </UFormGroup>
        <UFormGroup label="AccessKey" name="accessKey" :ui="{label:{base:'font-bold'}}">
          <UInput v-model="state.s3.accessKey"/>
        </UFormGroup>
        <UFormGroup label="SecretKey" name="secretKey" :ui="{label:{base:'font-bold'}}">
          <UInput v-model="state.s3.secretKey"/>
        </UFormGroup>
        <UFormGroup label="图片后缀（在访问缩略图时会追加在图片地址后）" name="thumbnailSuffix" :ui="{label:{base:'font-bold'}}">
          <UInput v-model="state.s3.thumbnailSuffix"/>
        </UFormGroup>
      </template>
      <UFormGroup label="是否启用邮件通知" name="enableEmail" :ui="{label:{base:'font-bold'}}">
      <UToggle v-model="state.enableEmail"/>
      </UFormGroup>
      <template v-if="state.enableEmail">
      <UFormGroup label="smtp服务器" name="smtpHost" :ui="{label:{base:'font-bold'}}">
        <UInput v-model="state.smtpHost" placeholder="smtp.qq.com"/>
      </UFormGroup>
      <UFormGroup label="smtp端口" name="smtpPort" :ui="{label:{base:'font-bold'}}">
        <UInput v-model="state.smtpPort" placeholder="465"/>
      </UFormGroup>
      <UFormGroup label="smtp用户名" name="smtpUsername" :ui="{label:{base:'font-bold'}}">
        <UInput v-model="state.smtpUsername" placeholder="******@qq.com"/>
      </UFormGroup>
      <UFormGroup label="smtp密码/授权码" name="smtpPassword" :ui="{label:{base:'font-bold'}}">
        <UInput v-model="state.smtpPassword" type="password"/>
      </UFormGroup>
      </template>
    <div class="space-y-1 pt-4">
      <div class="text-sm font-semibold text-gray-900 dark:text-white">用户管理</div>
      <div class="text-xs text-gray-500 dark:text-gray-400">查看所有用户，支持编辑昵称/邮箱/签名、重置密码，以及删除没有内容的普通用户。</div>
    </div>
    <div class="rounded-2xl border border-black/5 bg-white/70 p-4 shadow-sm dark:border-white/10 dark:bg-neutral-900/70">
      <div class="mb-4 flex items-center justify-between gap-3">
        <div class="text-sm text-gray-500 dark:text-gray-400">共 {{ users.length }} 个用户</div>
        <UButton size="sm" color="gray" variant="soft" @click="loadUsers" :loading="userLoading">刷新列表</UButton>
      </div>
      <div v-if="users.length" class="space-y-3">
        <div
          v-for="user in users"
          :key="user.id"
          class="rounded-2xl border border-black/5 bg-white/90 p-4 dark:border-white/10 dark:bg-neutral-800/80"
        >
          <div class="flex items-start justify-between gap-4">
            <div class="min-w-0 space-y-2">
              <div class="flex items-center gap-2">
                <span class="truncate text-base font-semibold text-gray-900 dark:text-white">
                  {{ user.nickname || user.username }}
                </span>
                <UBadge v-if="user.id === 1" size="xs" color="emerald" variant="soft">管理员</UBadge>
              </div>
              <div class="text-sm text-gray-500 dark:text-gray-400">@{{ user.username }}</div>
              <div class="text-sm text-gray-600 dark:text-gray-300">{{ user.email || "未填写邮箱" }}</div>
              <div class="text-sm text-gray-600 dark:text-gray-300">{{ user.slogan || "未填写签名" }}</div>
              <div class="text-xs text-gray-400 dark:text-gray-500">
                创建于 {{ formatDate(user.createdAt) }}
              </div>
            </div>
            <div class="flex shrink-0 items-center gap-2">
              <UButton size="sm" color="gray" variant="soft" @click="openUserEditor(user)">编辑</UButton>
              <UButton
                size="sm"
                color="red"
                variant="soft"
                :disabled="user.id === 1 || userDeletePending"
                @click="confirmDeleteUser(user)"
              >
                删除
              </UButton>
            </div>
          </div>
        </div>
      </div>
      <div v-else class="rounded-xl border border-dashed border-black/10 px-4 py-8 text-center text-sm text-gray-500 dark:border-white/10 dark:text-gray-400">
        当前还没有可展示的用户数据
      </div>
    </div>
    <UButton class="justify-center" color="red" @click="showCleanFileModal = true">清理已上传的文件</UButton>
    <UButton class="justify-center" @click="save">保存</UButton>
  </div>

  <UModal
    v-model="showCleanFileModal"
    :ui="{
      container:
        'flex justify-center items-center backdrop-blur',
    }"
  >
    <div class="p-4 bg-white dark:bg-neutral-800 rounded-lg shadow-md">
      <p class="text-lg font-bold mb-2">谨慎操作</p>
      <p class="text-gray-600 mb-4">确认要清理未使用的文件（图片、视频）吗？清理后，文件将被移动到 {uploadDir}/removed 目录下，请在检查后手动删除文件以释放空间。</p>
      <div class="flex justify-end gap-2 mt-4">
        <UButton color="white" @click="showCleanFileModal = false">取消</UButton>
        <UButton @click="cleanFile">确认清理</UButton>
      </div>
        </div>
  </UModal>

  <UModal
    v-model="showUserEditor"
    :ui="{
      container: 'flex justify-center items-center backdrop-blur',
    }"
  >
    <div class="w-full max-w-lg rounded-2xl bg-white p-5 shadow-md dark:bg-neutral-800">
      <div class="mb-4">
        <div class="text-lg font-bold text-gray-900 dark:text-white">编辑用户</div>
        <div class="mt-1 text-sm text-gray-500 dark:text-gray-400">@{{ editingUser.username }}</div>
      </div>
      <div class="space-y-4">
        <UFormGroup label="昵称" name="nickname" :ui="{label:{base:'font-bold'}}">
          <UInput v-model="editingUser.nickname"/>
        </UFormGroup>
        <UFormGroup label="邮箱" name="email" :ui="{label:{base:'font-bold'}}">
          <UInput v-model="editingUser.email"/>
        </UFormGroup>
        <UFormGroup label="签名" name="slogan" :ui="{label:{base:'font-bold'}}">
          <UTextarea v-model="editingUser.slogan" :rows="3"/>
        </UFormGroup>
        <UFormGroup label="重置密码（留空则不修改）" name="password" :ui="{label:{base:'font-bold'}}">
          <UInput v-model="editingUser.password" type="password" placeholder="输入新密码"/>
        </UFormGroup>
      </div>
      <div class="mt-5 flex justify-end gap-2">
        <UButton color="gray" variant="soft" @click="showUserEditor = false">取消</UButton>
        <UButton :loading="userSavePending" @click="saveUser">保存用户</UButton>
      </div>
    </div>
  </UModal>

  <UModal
    v-model="showDeleteUserModal"
    :ui="{
      container: 'flex justify-center items-center backdrop-blur',
    }"
  >
    <div class="w-full max-w-md rounded-2xl bg-white p-5 shadow-md dark:bg-neutral-800">
      <div class="text-lg font-bold text-gray-900 dark:text-white">确认删除用户</div>
      <div class="mt-2 text-sm text-gray-600 dark:text-gray-300">
        即将删除用户 <span class="font-semibold">{{ deletingUser?.username }}</span>。如果该用户还有动态或评论，系统会阻止删除。
      </div>
      <div class="mt-5 flex justify-end gap-2">
        <UButton color="gray" variant="soft" @click="showDeleteUserModal = false">取消</UButton>
        <UButton color="red" :loading="userDeletePending" @click="deleteUser">确认删除</UButton>
      </div>
    </div>
  </UModal>
</template>

<script setup lang="ts">
import type {AdminUserVO, SysConfigVO, UserVO} from "~/types";
import {toast} from "vue-sonner";
import {useUpload} from "~/utils";

const currentUser = useState<UserVO>('userinfo')
const version = ref('')
const commitId = ref('')
const state = reactive({
  enableGoogleRecaptcha: false,
  googleSiteKey:"",
  googleSecretKey:"",
  enableAutoLoadNextPage: true,
  enableComment: true,
  enableRegister: true,
  maxCommentLength: 120,
  memoMaxHeight: 300,
  commentOrder: 'desc',
  timeFormat: 'timeAgo',
  adminUserName: "admin",
  title: "极简朋友圈",
  favicon: "/favicon.ico",
  siteDescription: "把日常、思考和收藏都安静地放进时间线里。",
  announcement: "欢迎来到这个被认真整理过的小站，愿你也能在这里找到想停留的片刻。",
  heroTitle: "记录那些值得停留的瞬间",
  heroSubtitle: "这里保留日常、灵感与生活切片，用更轻的界面承载更完整的表达。",
  footerSignature: "慢慢写，慢慢看，慢慢生活。",
  beiAnNo: "",
  css: "",
  js: "",
  rss: "",
  enableS3: false,
  s3: {
    domain: "",
    bucket: "",
    region: "",
    accessKey: "",
    secretKey: "",
    endpoint: "",
    thumbnailSuffix: ""
  },
  enableEmail: false,
  smtpHost: "",
  smtpPort: "",
  smtpUsername: "",
  smtpPassword: "",
})

const showCleanFileModal = ref<boolean>(false);
const users = ref<Array<AdminUserVO>>([])
const userLoading = ref(false)
const showUserEditor = ref(false)
const showDeleteUserModal = ref(false)
const userSavePending = ref(false)
const userDeletePending = ref(false)
const deletingUser = ref<AdminUserVO | null>(null)
const editingUser = reactive({
  id: 0,
  username: "",
  nickname: "",
  email: "",
  slogan: "",
  password: "",
})

const reload = async () => {
  const res = await useMyFetch<SysConfigVO>('/sysConfig/getFull')
  if (res) {
    Object.assign(state, res)
    version.value = res.version
    commitId.value = res.commitId
  }
}

const loadUsers = async () => {
  userLoading.value = true
  try {
    users.value = await useMyFetch<Array<AdminUserVO>>('/user/list')
  } catch (error) {
    toast.error(error instanceof Error ? error.message : '读取用户列表失败')
  } finally {
    userLoading.value = false
  }
}

const save = async () => {
  await useMyFetch('/sysConfig/save', state)
  toast.success("保存成功")
  location.reload()
}

const uploadFavicon = async (files: FileList) => {
  for (let i = 0; i < files.length; i++) {
    if (files[i].type.indexOf("image") < 0){
      toast.error("只能上传图片");
      return
    }
  }
  const result = await useUpload(files)
  if (result.length) {
    toast.success("上传成功")
    state.favicon = result[0]
  }
}

const cleanFile = async () => {
  const res = await useMyFetch<{num: number}>('/file/clean', undefined)
  if (res) {
    toast.success(`成功清理 ${res.num} 个未使用的文件`)
    showCleanFileModal.value = false
  }
}

const openUserEditor = (user: AdminUserVO) => {
  editingUser.id = user.id
  editingUser.username = user.username
  editingUser.nickname = user.nickname || ""
  editingUser.email = user.email || ""
  editingUser.slogan = user.slogan || ""
  editingUser.password = ""
  showUserEditor.value = true
}

const saveUser = async () => {
  userSavePending.value = true
  try {
    await useMyFetch('/user/adminSave', {
      id: editingUser.id,
      nickname: editingUser.nickname,
      email: editingUser.email,
      slogan: editingUser.slogan,
      password: editingUser.password,
    })
    toast.success('用户信息已更新')
    showUserEditor.value = false
    await loadUsers()
  } catch (error) {
    toast.error(error instanceof Error ? error.message : '保存用户失败')
  } finally {
    userSavePending.value = false
  }
}

const confirmDeleteUser = (user: AdminUserVO) => {
  deletingUser.value = user
  showDeleteUserModal.value = true
}

const deleteUser = async () => {
  if (!deletingUser.value) {
    return
  }
  userDeletePending.value = true
  try {
    await useMyFetch('/user/delete?id=' + deletingUser.value.id)
    toast.success('用户已删除')
    showDeleteUserModal.value = false
    deletingUser.value = null
    await loadUsers()
  } catch (error) {
    toast.error(error instanceof Error ? error.message : '删除用户失败')
  } finally {
    userDeletePending.value = false
  }
}

const formatDate = (value?: string) => {
  if (!value) {
    return '未知时间'
  }
  return value.replace('T', ' ').slice(0, 16)
}

onMounted(async () => {
  await reload()
  await loadUsers()
})

</script>

<style scoped>

</style>
