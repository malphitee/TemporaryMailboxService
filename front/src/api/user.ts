import type {
    ApiResponse,
    ChangePasswordRequest,
    UpdateProfileRequest,
    User
} from '@/types/auth'
import request from '@/utils/request'

// 获取用户资料
export const getProfile = (): Promise<ApiResponse<User>> => {
  return request.get('/api/user/profile')
}

// 更新用户资料
export const updateProfile = (data: UpdateProfileRequest): Promise<ApiResponse<User>> => {
  return request.put('/api/user/profile', data)
}

// 修改密码
export const changePassword = (data: ChangePasswordRequest): Promise<ApiResponse> => {
  return request.post('/api/user/change-password', data)
} 