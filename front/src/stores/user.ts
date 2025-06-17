import * as userApi from '@/api/user'
import type { ChangePasswordRequest, UpdateProfileRequest, User } from '@/types/auth'
import { message } from 'ant-design-vue'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useAuthStore } from './auth'

export const useUserStore = defineStore('user', () => {
  // 状态
  const user = ref<User | null>(null)
  const loading = ref(false)

  // 加载用户资料
  const loadUserProfile = async () => {
    loading.value = true
    try {
      const response = await userApi.getProfile()
      if (response.success && response.data) {
        user.value = response.data
        // 同步更新authStore中的用户信息
        const authStore = useAuthStore()
        authStore.updateUserInfo(response.data)
      }
    } catch (error: any) {
      console.error('加载用户资料失败:', error)
      message.error('加载用户资料失败')
    } finally {
      loading.value = false
    }
  }

  // 更新用户资料
  const updateProfile = async (data: UpdateProfileRequest) => {
    loading.value = true
    try {
      const response = await userApi.updateProfile(data)
      if (response.success && response.data) {
        user.value = response.data
        // 同步更新authStore中的用户信息
        const authStore = useAuthStore()
        authStore.updateUserInfo(response.data)
        message.success('资料更新成功')
        return true
      } else {
        message.error(response.message || '更新失败')
        return false
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '更新失败')
      return false
    } finally {
      loading.value = false
    }
  }

  // 修改密码
  const changePassword = async (data: ChangePasswordRequest) => {
    loading.value = true
    try {
      const response = await userApi.changePassword(data)
      if (response.success) {
        message.success('密码修改成功')
        return true
      } else {
        message.error(response.message || '密码修改失败')
        return false
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '密码修改失败')
      return false
    } finally {
      loading.value = false
    }
  }

  // 重置状态
  const reset = () => {
    user.value = null
    loading.value = false
  }

  return {
    // 状态
    user,
    loading,
    // 方法
    loadUserProfile,
    updateProfile,
    changePassword,
    reset,
  }
}) 