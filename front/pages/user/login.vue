<template>
  <div class="auth-shell">
    <div class="auth-card">
      <div class="auth-header">
        <h1 class="auth-title">Moments</h1>
        <p class="auth-subtitle">登录你的账号</p>
      </div>

      <form class="auth-form" @submit.prevent="doLogin">
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
            autocomplete="current-password"
          />
        </label>

        <button type="submit" class="auth-submit" :disabled="pending">
          {{ pending ? "登录中..." : "登录" }}
        </button>
      </form>

      <p v-if="sysConfig.enableRegister" class="auth-switch">
        还没有账号?
        <NuxtLink to="/user/reg" class="auth-link">立即注册</NuxtLink>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { LoginResp, SysConfigVO } from "~/types";
import { useGlobalState } from "~/store";
import { toast } from "vue-sonner";

definePageMeta({
  layout: false,
});

const global = useGlobalState();
const sysConfig = useState<SysConfigVO>("sysConfig", () => ({ enableRegister: true } as SysConfigVO));
const state = reactive({
  username: "",
  password: "",
});
const pending = ref(false);

const sysConfigVO = await useMyFetch<SysConfigVO>("/sysConfig/get");
sysConfig.value = sysConfigVO;

useHead({
  title: `登录 - ${sysConfigVO.title || "Moments"}`,
  bodyAttrs: {
    class: "auth-page-body",
  },
});

const doLogin = async () => {
  pending.value = true;
  try {
    global.value.userinfo = await useMyFetch<LoginResp>("/user/login", state);
    toast.success("登录成功,跳转到首页...");
    window.location.href = "/";
  } catch (error) {
    toast.error(error instanceof Error ? error.message : "登录失败");
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
  min-height: 630px;
  padding: 60px 54px 48px;
  border-radius: 26px;
  background: #ffffff;
  box-shadow: 0 20px 60px rgba(47, 36, 103, 0.16);
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
  margin-top: 66px;
}

.auth-field {
  display: block;
}

.auth-field + .auth-field {
  margin-top: 34px;
}

.auth-field-password {
  margin-top: 32px;
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

.auth-submit {
  width: 100%;
  height: 78px;
  margin-top: 42px;
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
  margin: 42px 0 0;
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
    padding: 44px 28px 36px;
    border-radius: 24px;
  }

  .auth-title {
    font-size: 42px;
  }

  .auth-form {
    margin-top: 44px;
  }

  .auth-input,
  .auth-submit {
    height: 64px;
  }
}
</style>