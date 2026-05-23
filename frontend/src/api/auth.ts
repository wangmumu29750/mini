import { request } from '@/api/http'
import type { CurrentUser, LoginRequest, LoginResponse, RegisterRequest } from '@/types/auth'

export function login(payload: LoginRequest) {
  return request<LoginResponse>({
    method: 'POST',
    url: '/auth/login',
    data: payload,
  })
}

export function register(payload: RegisterRequest) {
  return request<LoginResponse>({
    method: 'POST',
    url: '/auth/register',
    data: payload,
  })
}

export function fetchCurrentUser() {
  return request<CurrentUser>({
    method: 'GET',
    url: '/auth/me',
  })
}

export function logout() {
  return request<null>({
    method: 'POST',
    url: '/auth/logout',
  })
}

