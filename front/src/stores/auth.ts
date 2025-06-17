import * as authApi from '@/api/auth'
import * as userApi from '@/api/user'
import type { LoginRequest, RegisterRequest, User } from '@/types/auth'
import { message } from 'ant-design-vue'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<User | null>(null)
  const loading = ref(false)

  // 计算属性
  const isLoggedIn = computed(() => !!token.value && !!user.value)

  // 初始化用户信息
  const initUser = async () => {
    if (!token.value) return
    
    try {
      const response = await userApi.getProfile()
      if (response.code === 0 && response.data) {
        user.value = response.data
        // 持久化用户信息
        localStorage.setItem('user', JSON.stringify(response.data))
      }
    } catch (error) {
      console.error('获取用户信息失败:', error)
      // 清除无效token
      logout()
    }
  }

  // 登录
  const login = async (loginData: LoginRequest) => {
    loading.value = true
    try {
      const response = await authApi.login(loginData)
      if (response.code === 0 && response.data) {
        token.value = response.data.token.access_token
        user.value = response.data.user
        // 持久化到localStorage
        localStorage.setItem('token', response.data.token.access_token)
        localStorage.setItem('user', JSON.stringify(response.data.user))
        message.success('登录成功')
        return true
      } else {
        message.error(response.message || '登录失败')
        return false
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '登录失败')
      return false
    } finally {
      loading.value = false
    }
  }

  // 注册
  const register = async (registerData: RegisterRequest) => {
    loading.value = true
    try {
      const response = await authApi.register(registerData)
      if (
        response.code === 0 &&
        response.data &&
        response.data.token &&
        response.data.user
      ) {
        token.value = response.data.token.access_token
        user.value = response.data.user
        // 持久化到localStorage
        localStorage.setItem('token', response.data.token.access_token)
        localStorage.setItem('user', JSON.stringify(response.data.user))
        message.success('注册成功')
        return true
      } else {
        message.error(response.message || '注册失败')
        return false
      }
    } catch (error: any) {
      message.error(error.response?.data?.message || '注册失败')
      return false
    } finally {
      loading.value = false
    }
  }

  // 退出登录
  const logout = () => {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    message.success('已退出登录')
  }

  // 更新用户信息
  const updateUserInfo = (newUser: User) => {
    user.value = newUser
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  // 从localStorage恢复用户信息
  const restoreFromStorage = () => {
    const storedUser = localStorage.getItem('user')
    if (storedUser && token.value) {
      try {
        user.value = JSON.parse(storedUser)
      } catch (error) {
        console.error('解析用户信息失败:', error)
        logout()
      }
    }
  }

  return {
    // 状态
    token,
    user,
    loading,
    // 计算属性
    isLoggedIn,
    // 方法
    initUser,
    login,
    register,
    logout,
    updateUserInfo,
    restoreFromStorage,
  }
}) 