<template>
  <div class="auth-shell">
    <div class="auth-card auth-card-register">
      <div class="auth-header">
        <h1 class="auth-title">Moments</h1>
        <p class="auth-subtitle">创建你的账号</p>
      </div>

      <form class="auth-form" @submit.prevent="doReg">
        <label class="auth-field">
          <span class="auth-label">用户名</span>
          <input
            v-model="state.username"
            type="text"
            class="auth-input"
            placeholder="请输入用户名"
            autocomplete="username"
          />
        </label>

        <label class="auth-field auth-field-password">
          <span class="auth-label">密码</span>
          <input
            v-model="state.password"
            type="password"
            class="auth-input"
            placeholder="请输入密码"
            autocomplete="new-password"
          />
          <span class="auth-hint">至少 8 个字符，包含大小写字母和数字</span>
        </label>

        <label class="auth-field auth-field-confirm">
          <span class="auth-label">确认密码</span>
          <input
            v-model="state.repeatPassword"
            type="password"
            class="auth-input"
            placeholder="请再次输入密码"
            autocomplete="new-password"
          />
        </label>

        <button type="submit" class="auth-submit" :disabled="pending">
          {{ pending ? "注册中..." : "注册" }}
        </button>
      </form>

      <p class="auth-switch">
        已有账号?
        <NuxtLink to="/user/login" class="auth-link">立即登录</NuxtLink>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { SysConfigVO } from "~/types";
import { toast } from "vue-sonner";

definePageMeta({
  layout: false,
});

const state = reactive({
  username: "",
  password: "",
  repeatPassword: "",
});
const pending = ref(false);
const sysConfig = useState<SysConfigVO>("sysConfig");

const sysConfigVO = await useMyFetch<SysConfigVO>("/sysConfig/get");
sysConfig.value = sysConfigVO;

if (!sysConfigVO.enableRegister) {
  await navigateTo("/");
}

useHead({
  title: `注册 - ${sysConfigVO.title || "Moments"}`,
  bodyAttrs: {
    class: "auth-page-body",
  },
});

const doReg = async () => {
  if (state.username.length < 3) {
    toast.warning("用户名最少3个字符");
    return;
  }

  pending.value = true;
  try {
    await useMyFetch("/user/reg", state);
    toast.success("注册成功,快去登录吧!");
    await navigateTo("/user/login");
  } catch (error) {
    toast.error(error instanceof Error ? error.message : "注册失败");
  } finally {
    pending.value = false;
  }
};
</script>

<style scoped>
.auth-shell {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 24px;
  background: linear-gradient(135deg, #6f84ff 0%, #7e56c2 100%);
}

.auth-card {
  width: 100%;
  max-width: 560px;
  padding: 56px 54px 44px;
  border-radius: 26px;
  background: #ffffff;
  box-shadow: 0 20px 60px rgba(47, 36, 103, 0.16);
}

.auth-card-register {
  min-height: 720px;
}

.auth-header {
  text-align: center;
}

.auth-title {
  margin: 0;
  color: #2b2b2f;
  font-size: 56px;
  line-height: 1;
  font-weight: 800;
  letter-spacing: -0.04em;
}

.auth-subtitle {
  margin: 18px 0 0;
  color: #737373;
  font-size: 18px;
  line-height: 1.4;
  font-weight: 500;
}

.auth-form {
  margin-top: 54px;
}

.auth-field {
  display: block;
}

.auth-field + .auth-field {
  margin-top: 24px;
}

.auth-label {
  display: block;
  margin-bottom: 18px;
  color: #262626;
  font-size: 18px;
  line-height: 1.3;
  font-weight: 700;
}

.auth-input {
  width: 100%;
  height: 74px;
  padding: 0 26px;
  border: 1.5px solid #e4e4e7;
  border-radius: 16px;
  outline: none;
  background: #ffffff;
  color: #202020;
  font-size: 18px;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.auth-input::placeholder {
  color: #c3c5cc;
}

.auth-input:focus {
  border-color: #7a75f7;
  box-shadow: 0 0 0 4px rgba(122, 117, 247, 0.12);
}

.auth-hint {
  display: block;
  margin-top: 14px;
  color: #8a8a8a;
  font-size: 15px;
  line-height: 1.5;
  font-weight: 500;
}

.auth-submit {
  width: 100%;
  height: 78px;
  margin-top: 34px;
  border: 0;
  border-radius: 16px;
  background: linear-gradient(90deg, #6e86ff 0%, #7d4fc1 100%);
  color: #ffffff;
  font-size: 22px;
  line-height: 1;
  font-weight: 800;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease, opacity 0.2s ease;
}

.auth-submit:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 16px 28px rgba(95, 92, 214, 0.22);
}

.auth-submit:disabled {
  opacity: 0.8;
  cursor: wait;
}

.auth-switch {
  margin: 34px 0 0;
  text-align: center;
  color: #747474;
  font-size: 18px;
  line-height: 1.5;
  font-weight: 500;
}

.auth-link {
  margin-left: 8px;
  color: #5f67f3;
  font-weight: 700;
  text-decoration: none;
}

.auth-link:hover {
  text-decoration: underline;
}

@media (max-width: 640px) {
  .auth-shell {
    padding: 24px 16px;
  }

  .auth-card {
    min-height: auto;
    padding: 42px 28px 34px;
    border-radius: 24px;
  }

  .auth-title {
    font-size: 42px;
  }

  .auth-form {
    margin-top: 40px;
  }

  .auth-input,
  .auth-submit {
    height: 64px;
  }
}
</style>