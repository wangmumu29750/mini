import { defineStore } from 'pinia'

import * as authApi from '@/api/auth'
import type { CurrentUser, LoginRequest, RegisterRequest, UserRole } from '@/types/auth'

const TOKEN_KEY = 'mini12306_access_token'
const USER_KEY = 'mini12306_current_user'

interface AuthState {
  token: string
  user: CurrentUser | null
}

function readStoredUser() {
  const raw = window.localStorage.getItem(USER_KEY)

  if (!raw) {
    return null
  }

  try {
    return JSON.parse(raw) as CurrentUser
  } catch {
    window.localStorage.removeItem(USER_KEY)
    return null
  }
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: window.localStorage.getItem(TOKEN_KEY) || '',
    user: readStoredUser(),
  }),
  getters: {
    isAuthenticated: (state) => Boolean(state.token),
    role: (state): UserRole | null => state.user?.role || null,
  },
  actions: {
    setSession(token: string, user: CurrentUser) {
      this.token = token
      this.user = user
      window.localStorage.setItem(TOKEN_KEY, token)
      window.localStorage.setItem(USER_KEY, JSON.stringify(user))
    },
    clearSession() {
      this.token = ''
      this.user = null
      window.localStorage.removeItem(TOKEN_KEY)
      window.localStorage.removeItem(USER_KEY)
    },
    async login(payload: LoginRequest) {
      const result = await authApi.login(payload)
      this.setSession(result.accessToken, result.user)
      return result.user
    },
    async register(payload: RegisterRequest) {
      const result = await authApi.register(payload)
      this.setSession(result.accessToken, result.user)
      return result.user
    },
    async refreshCurrentUser() {
      if (!this.token) {
        return null
      }

      const user = await authApi.fetchCurrentUser()
      this.user = user
      window.localStorage.setItem(USER_KEY, JSON.stringify(user))
      return user
    },
    async logout() {
      try {
        await authApi.logout()
      } finally {
        this.clearSession()
      }
    },
  },
})

