<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { registerApi } from '../api/auth'

const router = useRouter()

const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const email = ref('')
const error = ref('')
const loading = ref(false)
const success = ref('')

function validateEmail(v: string) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v)
}

async function submit() {
  error.value = ''
  success.value = ''
  if (!username.value.trim()) {
    error.value = '请输入用户名'
    return
  }
  if (username.value.trim().length < 3) {
    error.value = '用户名至少 3 个字符'
    return
  }
  if (!password.value) {
    error.value = '请输入密码'
    return
  }
  if (password.value.length < 6) {
    error.value = '密码至少 6 个字符'
    return
  }
  if (password.value !== confirmPassword.value) {
    error.value = '两次输入的密码不一致'
    return
  }
  if (!email.value.trim()) {
    error.value = '请输入邮箱'
    return
  }
  if (!validateEmail(email.value.trim())) {
    error.value = '请输入有效的邮箱地址'
    return
  }
  loading.value = true
  try {
    await registerApi({
      username: username.value.trim(),
      password: password.value,
      email: email.value.trim(),
    })
    success.value = '注册成功，请登录'
    setTimeout(() => router.push('/login'), 1500)
  } catch (e) {
    error.value = e instanceof Error ? e.message : '注册失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="register-page">
    <div class="bg-shapes">
      <div class="shape shape-1" />
      <div class="shape shape-2" />
      <div class="shape shape-3" />
    </div>

    <div class="register-card glass">
      <div class="card-header">
        <div class="logo-wrap">
          <span class="logo-icon">◆</span>
        </div>
        <h1 class="title">注册账号</h1>
        <p class="subtitle">创建 TaskFlow 账号以开始管理任务</p>
      </div>

      <form class="form" @submit.prevent="submit">
        <div class="field">
          <label for="reg-username">用户名</label>
          <input
            id="reg-username"
            v-model="username"
            type="text"
            placeholder="至少 3 个字符"
            autocomplete="username"
          />
        </div>
        <div class="field">
          <label for="reg-email">邮箱</label>
          <input
            id="reg-email"
            v-model="email"
            type="email"
            placeholder="your@email.com"
            autocomplete="email"
          />
        </div>
        <div class="field">
          <label for="reg-password">密码</label>
          <input
            id="reg-password"
            v-model="password"
            type="password"
            placeholder="至少 6 个字符"
            autocomplete="new-password"
          />
        </div>
        <div class="field">
          <label for="reg-confirm">确认密码</label>
          <input
            id="reg-confirm"
            v-model="confirmPassword"
            type="password"
            placeholder="再次输入密码"
            autocomplete="new-password"
          />
        </div>
        <p v-if="error" class="error">{{ error }}</p>
        <p v-if="success" class="success">{{ success }}</p>
        <button type="submit" class="btn-primary" :disabled="loading">
          {{ loading ? '注册中…' : '注 册' }}
        </button>
        <p class="link-row">
          已有账号？
          <router-link to="/login">去登录</router-link>
        </p>
      </form>
    </div>
  </div>
</template>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
  position: relative;
  overflow: hidden;
}

.bg-shapes {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.shape {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
}

.shape-1 {
  width: 400px;
  height: 400px;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-dark) 100%);
  top: -100px;
  right: -100px;
  animation: float 8s ease-in-out infinite;
}

.shape-2 {
  width: 300px;
  height: 300px;
  background: linear-gradient(135deg, var(--surface) 0%, var(--accent) 50%);
  bottom: -50px;
  left: -50px;
  animation: float 10s ease-in-out infinite reverse;
}

.shape-3 {
  width: 200px;
  height: 200px;
  background: var(--accent);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  opacity: 0.15;
  animation: pulse 4s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(20px, -20px); }
}

@keyframes pulse {
  0%, 100% { opacity: 0.15; transform: translate(-50%, -50%) scale(1); }
  50% { opacity: 0.25; transform: translate(-50%, -50%) scale(1.1); }
}

.register-card {
  position: relative;
  width: 100%;
  max-width: 400px;
  padding: 2.5rem;
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
}

.glass {
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(20px);
}

.card-header {
  text-align: center;
  margin-bottom: 2rem;
}

.logo-wrap {
  width: 56px;
  height: 56px;
  margin: 0 auto 1rem;
  border-radius: 16px;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-dark) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 24px rgba(251, 191, 36, 0.3);
}

.logo-icon {
  font-size: 1.5rem;
  color: var(--bg);
}

.title {
  font-family: var(--font-display);
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--text);
  margin: 0 0 0.25rem;
  letter-spacing: -0.02em;
}

.subtitle {
  color: var(--text-muted);
  font-size: 0.95rem;
  margin: 0;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.field label {
  font-size: 0.85rem;
  font-weight: 500;
  color: var(--text-muted);
}

.field input {
  width: 100%;
  padding: 0.875rem 1rem;
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.05);
  color: var(--text);
  font-size: 1rem;
  transition: border-color 0.2s, box-shadow 0.2s;
  box-sizing: border-box;
}

.field input::placeholder {
  color: var(--text-muted);
  opacity: 0.7;
}

.field input:focus {
  outline: none;
  border-color: var(--accent);
  box-shadow: 0 0 0 3px rgba(251, 191, 36, 0.15);
}

.error {
  color: #f87171;
  font-size: 0.875rem;
  margin: -0.25rem 0 0;
}

.success {
  color: #4ade80;
  font-size: 0.875rem;
  margin: -0.25rem 0 0;
}

.btn-primary {
  margin-top: 0.5rem;
  padding: 0.875rem 1.5rem;
  border-radius: 12px;
  border: none;
  background: linear-gradient(135deg, var(--accent) 0%, var(--accent-dark) 100%);
  color: var(--bg);
  font-family: var(--font-display);
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.15s, box-shadow 0.2s;
  box-shadow: 0 4px 14px rgba(251, 191, 36, 0.35);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(251, 191, 36, 0.4);
}

.btn-primary:active:not(:disabled) {
  transform: translateY(0);
}

.btn-primary:disabled {
  opacity: 0.8;
  cursor: not-allowed;
}

.link-row {
  margin: 0;
  font-size: 0.9rem;
  color: var(--text-muted);
  text-align: center;
}

.link-row a {
  color: var(--accent);
  text-decoration: none;
  font-weight: 500;
}

.link-row a:hover {
  text-decoration: underline;
}
</style>
