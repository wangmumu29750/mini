import { request } from '@/api/http'
import type { CurrentUser, CreatePassengerProfileRequest, LoginRequest, LoginResponse, RegisterRequest } from '@/types/auth'
import type { PassengerSummary } from '@/types/domain'

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

export function createPassenger(payload: CreatePassengerProfileRequest) {
  return request<PassengerSummary>({
    method: 'POST',
    url: '/auth/passengers',
    data: payload,
  })
}
