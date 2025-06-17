import type {
    ApiResponse,
    LoginRequest,
    LoginResponse,
    RegisterRequest
} from '@/types/auth'
import request from '@/utils/request'

// 用户登录
export const login = (data: LoginRequest): Promise<ApiResponse<LoginResponse>> => {
  return request.post('/api/auth/login', data)
}

// 用户注册
export const register = (data: RegisterRequest): Promise<ApiResponse<LoginResponse>> => {
  return request.post('/api/auth/register', data)
}

// 退出登录
export const logout = (): Promise<ApiResponse> => {
  return request.post('/api/auth/logout')
} 