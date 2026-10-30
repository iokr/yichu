import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import * as authApi from '@/api/auth'
import { session } from '@/api/client'
import type { UserInfo } from '@/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<UserInfo | null>(session.user<UserInfo>())
  const token = ref(session.token())
  const loading = ref(false)

  const isLoggedIn = computed(() => !!token.value)

  function sync() {
    token.value = session.token()
    user.value = session.user<UserInfo>()
  }

  async function login(username: string, password: string) {
    loading.value = true
    try {
      const data = await authApi.login({ username, password })
      session.save(data)
      sync()
      return data
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    try {
      if (token.value) await authApi.logout()
    } catch {
      // 退出失败不阻塞本地清理
    } finally {
      session.clear()
      sync()
    }
  }

  function clear() {
    session.clear()
    sync()
  }

  window.addEventListener('yichu:session', sync)

  return { user, token, loading, isLoggedIn, login, logout, clear, sync }
})
