import axios, { AxiosError } from 'axios'

import type { ApiErrorPayload, ApiResponse } from '@/types/api'
import { useAuthStore } from '@/stores/auth'

export const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

http.interceptors.request.use((config) => {
  const authStore = useAuthStore()

  if (authStore.token) {
    config.headers.Authorization = `Bearer ${authStore.token}`
  }

  return config
})

http.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiResponse<null>>) => {
    const status = error.response?.status
    const payload = error.response?.data
    const normalized: ApiErrorPayload = {
      code: payload?.code || 'NETWORK_ERROR',
      message: payload?.message || '网络异常，请稍后重试',
      traceId: payload?.traceId,
      status,
    }

    if (status === 401) {
      const authStore = useAuthStore()
      authStore.clearSession()
    }

    return Promise.reject(normalized)
  },
)

export async function request<T>(config: Parameters<typeof http.request<ApiResponse<T>>>[0]) {
  const response = await http.request<ApiResponse<T>>(config)
  return response.data.data
}

