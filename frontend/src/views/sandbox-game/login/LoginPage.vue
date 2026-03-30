<template>
  <div class="login-page">
    <form class="login-card" @submit.prevent="handleSubmit">
      <div class="brand-block">
        <p class="eyebrow">Sandbox Game</p>
        <h1>沙盘经营系统</h1>
      </div>

      <div class="card-header">
        <h2>账号登录</h2>
      </div>

      <label class="field">
        <span>用户名</span>
        <input v-model.trim="username" type="text" autocomplete="username" placeholder="请输入用户名" />
      </label>

      <label class="field">
        <span>密码</span>
        <input
          v-model="password"
          type="password"
          autocomplete="current-password"
          placeholder="请输入密码"
        />
      </label>

      <p v-if="errorMessage" class="error-message">{{ errorMessage }}</p>

      <button class="submit-button" type="submit" :disabled="submitting || !canSubmit">
        {{ submitting ? '登录中...' : '登录' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const username = ref('')
const password = ref('')
const submitting = ref(false)
const errorMessage = ref('')

const canSubmit = computed(() => username.value.trim() !== '' && password.value.trim() !== '')

async function handleSubmit() {
  if (!canSubmit.value) {
    errorMessage.value = '请输入用户名和密码'
    return
  }

  submitting.value = true
  errorMessage.value = ''
  try {
    const defaultRoute = await authStore.login(username.value, password.value)
    await router.replace(defaultRoute)
  } catch (error) {
    errorMessage.value = error instanceof Error && error.message ? error.message : '登录失败，请稍后重试'
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.login-card {
  width: min(520px, 100%);
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid #d7dfeb;
  border-radius: 24px;
  box-shadow: var(--shadow);
  padding: 34px 32px 30px;
}

.brand-block {
  margin-bottom: 28px;
}

.eyebrow {
  margin: 0;
  color: var(--accent);
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.brand-block h1 {
  margin: 10px 0 0;
  font-size: 38px;
  line-height: 1.1;
}

.card-header {
  margin-bottom: 10px;
}

.card-header h2 {
  margin: 0;
  font-size: 28px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 18px;
}

.field span {
  font-size: 14px;
  color: var(--text);
  font-weight: 600;
}

.field input {
  width: 100%;
  height: 46px;
  padding: 0 14px;
  border-radius: 14px;
  border: 1px solid var(--line-strong);
  background: #ffffff;
  outline: none;
}

.field input:focus {
  border-color: #8ab0ff;
  box-shadow: 0 0 0 3px rgba(31, 95, 211, 0.12);
}

.error-message {
  margin: 16px 0 0;
  padding: 10px 12px;
  border-radius: 12px;
  border: 1px solid #efc4c4;
  background: #fff5f5;
  color: var(--danger);
  font-size: 14px;
}

.submit-button {
  width: 100%;
  height: 48px;
  margin-top: 20px;
  border: none;
  border-radius: 14px;
  background: linear-gradient(135deg, #1f5fd3 0%, #2d78ef 100%);
  color: #ffffff;
  font-size: 16px;
  font-weight: 700;
}

.submit-button:disabled {
  opacity: 0.6;
}

@media (max-width: 640px) {
  .login-page {
    padding: 14px;
  }

  .login-card {
    padding: 24px 20px 22px;
    border-radius: 18px;
  }

  .brand-block h1 {
    font-size: 30px;
  }
}
</style>
