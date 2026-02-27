import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { loginApi } from '../api/auth'

const TOKEN_KEY = 'taskflow_token'
const USERNAME_KEY = 'taskflow_username'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem(TOKEN_KEY))
  const username = ref<string>(localStorage.getItem(USERNAME_KEY) ?? '')

  const isLoggedIn = computed(() => !!token.value)

  function setAuth(t: string, name: string) {
    token.value = t
    username.value = name
    localStorage.setItem(TOKEN_KEY, t)
    localStorage.setItem(USERNAME_KEY, name)
  }

  async function login(name: string, password: string): Promise<void> {
    const res = await loginApi({ username: name, password })
    setAuth(res.token, res.user.username)
  }

  function logout() {
    token.value = null
    username.value = ''
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USERNAME_KEY)
  }

  return { token, username, isLoggedIn, setAuth, login, logout }
})
