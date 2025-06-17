// 用户认证相关类型定义

export interface User {
  id: number
  username: string
  email: string
  nickname: string
  avatar: string
  timezone: string
  language: string
  is_active: boolean
  last_login_at?: string
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  nickname: string
}

export interface LoginResponse {
  token: {
    access_token: string
    refresh_token: string
    token_type: string
    expires_in: number
  }
  user: User
}

export interface UpdateProfileRequest {
  nickname: string
  avatar?: string
  timezone?: string
  language?: string
}

export interface ChangePasswordRequest {
  current_password: string
  new_password: string
}

export interface ApiResponse<T = any> {
  code?: number
  success?: boolean
  data?: T
  message?: string
  error?: string
} 