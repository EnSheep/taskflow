<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()

const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function submit() {
  error.value = ''
  if (!username.value.trim()) {
    error.value = '请输入用户名'
    return
  }
  if (!password.value) {
    error.value = '请输入密码'
    return
  }
  loading.value = true
  try {
    await auth.login(username.value.trim(), password.value)
    router.push('/')
  } catch (e) {
    error.value = e instanceof Error ? e.message : '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <div class="bg-shapes">
      <div class="shape shape-1" />
      <div class="shape shape-2" />
      <div class="shape shape-3" />
    </div>

    <div class="login-card glass">
      <div class="card-header">
        <div class="logo-wrap">
          <span class="logo-icon">◆</span>
        </div>
        <h1 class="title">TaskFlow</h1>
        <p class="subtitle">登录以管理你的任务</p>
      </div>

      <form class="form" @submit.prevent="submit">
        <div class="field">
          <label for="username">用户名</label>
          <input
            id="username"
            v-model="username"
            type="text"
            placeholder="输入用户名"
            autocomplete="username"
          />
        </div>
        <div class="field">
          <label for="password">密码</label>
          <input
            id="password"
            v-model="password"
            type="password"
            placeholder="输入密码"
            autocomplete="current-password"
          />
        </div>
        <p v-if="error" class="error">{{ error }}</p>
        <button type="submit" class="btn-primary" :disabled="loading">
          {{ loading ? '登录中…' : '登 录' }}
        </button>
        <p class="link-row">
          还没有账号？
          <router-link to="/register">立即注册</router-link>
        </p>
      </form>
    </div>
  </div>
</template>

<style scoped>
.login-page {
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

.login-card {
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
