export type UserRole = 'PASSENGER' | 'ADMIN'

export interface CurrentUser {
  id: number
  username: string
  role: UserRole
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  accessToken: string
  user: CurrentUser
}

export interface RegisterRequest {
  username: string
  password: string
  realName: string
  idCardNo: string
  phone: string
  bankCardNo: string
}

