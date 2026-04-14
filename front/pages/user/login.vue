<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-500 via-purple-500 to-purple-600">
    <div class="w-full max-w-md mx-4">
      <!-- 登录卡片 -->
      <div class="bg-white rounded-2xl shadow-2xl p-8 md:p-10">
        <!-- 标题 -->
        <div class="text-center mb-8">
          <h1 class="text-3xl font-bold text-gray-800 mb-2">Moments</h1>
          <p class="text-gray-500">登录你的账号</p>
        </div>

        <!-- 登录表单 -->
        <form class="space-y-6" @submit.prevent="doLogin">
          <!-- 用户名 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">用户名</label>
            <input
              v-model="state.username"
              type="text"
              placeholder="请输入用户名"
              class="w-full px-4 py-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all"
              :disabled="pending"
            />
          </div>

          <!-- 密码 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">密码</label>
            <input
              v-model="state.password"
              type="password"
              placeholder="请输入密码"
              class="w-full px-4 py-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all"
              :disabled="pending"
            />
          </div>

          <!-- 登录按钮 -->
          <button
            type="submit"
            :disabled="pending || !state.username || !state.password"
            class="w-full py-3.5 px-4 bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-medium rounded-xl hover:from-indigo-600 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all transform hover:scale-[1.02] active:scale-[0.98]"
          >
            <span v-if="pending" class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              登录中...
            </span>
            <span v-else>登录</span>
          </button>
        </form>

        <!-- 注册链接 -->
        <div class="mt-6 text-center">
          <span class="text-gray-500">还没有账号？</span>
          <NuxtLink
            to="/user/reg"
            class="ml-1 text-purple-600 hover:text-purple-700 font-medium transition-colors"
          >
            立即注册
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { LoginResp, SysConfigVO, UserVO } from "~/types";
import { useGlobalState } from "~/store";
import { toast } from "vue-sonner";

definePageMeta({
  layout: false
});

const sysConfig = useState<SysConfigVO>('sysConfig');
const currentUser = useState<UserVO>('userinfo');
const global = useGlobalState();

const state = reactive({
  username: "",
  password: ""
});

const pending = ref(false);

const doLogin = async () => {
  if (!state.username || !state.password) {
    toast.warning("请输入用户名和密码");
    return;
  }

  pending.value = true;
  let success = false;

  try {
    global.value.userinfo = await useMyFetch<LoginResp>('/user/login', state);
    toast.success("登录成功，正在跳转...");
    success = true;
  } catch (error: any) {
    toast.error(error?.message || "登录失败，请检查用户名和密码");
  } finally {
    pending.value = false;
  }

  if (success) {
    navigateTo('/');
  }
};
</script>

<style scoped>
/* 自定义动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.bg-white {
  animation: fadeIn 0.5s ease-out;
}
</style>
