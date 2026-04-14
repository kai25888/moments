<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-500 via-purple-500 to-purple-600">
    <div class="w-full max-w-md mx-4">
      <!-- 注册卡片 -->
      <div class="bg-white rounded-2xl shadow-2xl p-8 md:p-10">
        <!-- 标题 -->
        <div class="text-center mb-8">
          <h1 class="text-3xl font-bold text-gray-800 mb-2">Moments</h1>
          <p class="text-gray-500">创建你的账号</p>
        </div>

        <!-- 注册表单 -->
        <form class="space-y-5" @submit.prevent="doReg">
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
            <p class="mt-1.5 text-xs text-gray-400">至少 8 个字符，包含大小写字母和数字</p>
          </div>

          <!-- 确认密码 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">确认密码</label>
            <input
              v-model="state.repeatPassword"
              type="password"
              placeholder="请再次输入密码"
              class="w-full px-4 py-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all"
              :disabled="pending"
            />
          </div>

          <!-- 注册按钮 -->
          <button
            type="submit"
            :disabled="pending || !state.username || !state.password || !state.repeatPassword"
            class="w-full py-3.5 px-4 bg-gradient-to-r from-indigo-500 to-purple-600 text-white font-medium rounded-xl hover:from-indigo-600 hover:to-purple-700 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all transform hover:scale-[1.02] active:scale-[0.98]"
          >
            <span v-if="pending" class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              注册中...
            </span>
            <span v-else>注册</span>
          </button>
        </form>

        <!-- 登录链接 -->
        <div class="mt-6 text-center">
          <span class="text-gray-500">已有账号？</span>
          <NuxtLink
            to="/user/login"
            class="ml-1 text-purple-600 hover:text-purple-700 font-medium transition-colors"
          >
            立即登录
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { UserVO } from "~/types";
import { toast } from "vue-sonner";

definePageMeta({
  layout: false
});

const state = reactive({
  username: "",
  password: "",
  repeatPassword: ""
});

const pending = ref(false);
const currentUser = useState<UserVO>('userinfo');

const doReg = async () => {
  // 表单验证
  if (state.username.length < 3) {
    toast.warning("用户名最少3个字符");
    return;
  }

  if (state.password.length < 8) {
    toast.warning("密码至少8个字符");
    return;
  }

  // 密码强度检查
  const hasUpperCase = /[A-Z]/.test(state.password);
  const hasLowerCase = /[a-z]/.test(state.password);
  const hasNumbers = /\d/.test(state.password);

  if (!hasUpperCase || !hasLowerCase || !hasNumbers) {
    toast.warning("密码需包含大小写字母和数字");
    return;
  }

  if (state.password !== state.repeatPassword) {
    toast.warning("两次输入的密码不一致");
    return;
  }

  let success = false;
  pending.value = true;

  try {
    await useMyFetch('/user/reg', {
      username: state.username,
      password: state.password,
      repeatPassword: state.repeatPassword
    });
    toast.success("注册成功，快去登录吧！");
    success = true;
  } catch (error: any) {
    toast.error(error?.message || "注册失败，请稍后重试");
  } finally {
    pending.value = false;
  }

  if (success) {
    await navigateTo('/user/login');
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
