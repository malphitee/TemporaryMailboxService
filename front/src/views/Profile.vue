<template>
  <div class="profile">
    <a-layout>
      <a-layout-header class="header">
        <div class="header-content">
          <a-button type="text" @click="$router.back()">
            <ArrowLeftOutlined />
            返回
          </a-button>
          <h1>个人资料</h1>
          <div></div>
        </div>
      </a-layout-header>
      
      <a-layout-content class="content">
        <a-card>
          <template #title>
            <UserOutlined />
            基本信息
          </template>
          <template #extra>
            <a-button type="primary" @click="$router.push('/profile/edit')">
              <EditOutlined />
              编辑资料
            </a-button>
          </template>
          
          <a-descriptions :column="1" bordered>
            <a-descriptions-item label="用户名">
              {{ authStore.user?.username || '未设置' }}
            </a-descriptions-item>
            <a-descriptions-item label="邮箱">
              {{ authStore.user?.email || '未设置' }}
            </a-descriptions-item>
            <a-descriptions-item label="昵称">
              {{ displayName }}
            </a-descriptions-item>
            <a-descriptions-item label="头像">
              <a-avatar :size="48" :src="authStore.user?.avatar">
                {{ displayName?.charAt(0) || 'U' }}
              </a-avatar>
            </a-descriptions-item>
            <a-descriptions-item label="时区">
              {{ authStore.user?.timezone || '未设置' }}
            </a-descriptions-item>
            <a-descriptions-item label="语言">
              {{ authStore.user?.language || '未设置' }}
            </a-descriptions-item>
            <a-descriptions-item label="注册时间">
              {{ formatDate(authStore.user?.created_at) }}
            </a-descriptions-item>
            <a-descriptions-item label="最后更新">
              {{ formatDate(authStore.user?.updated_at) }}
            </a-descriptions-item>
          </a-descriptions>
        </a-card>
        
        <a-card style="margin-top: 24px;">
          <template #title>
            <SettingOutlined />
            账户操作
          </template>
          
          <a-space direction="vertical" style="width: 100%">
            <a-button type="default" block @click="$router.push('/profile/password')">
              <LockOutlined />
              修改密码
            </a-button>
            <a-button type="danger" block @click="handleLogout">
              <LogoutOutlined />
              退出登录
            </a-button>
          </a-space>
        </a-card>
      </a-layout-content>
    </a-layout>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import {
    ArrowLeftOutlined,
    EditOutlined,
    LockOutlined,
    LogoutOutlined,
    SettingOutlined,
    UserOutlined
} from '@ant-design/icons-vue'
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()

const formatDate = (dateString?: string) => {
  if (!dateString) return '未知'
  return new Date(dateString).toLocaleString('zh-CN')
}

const displayName = computed(() => {
  const user = authStore.user
  if (!user) return '未知用户'
  
  // 使用nickname字段显示昵称
  return user.nickname || user.username || ''
})

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

onMounted(async () => {
  // 初始化用户信息
  if (authStore.token && !authStore.user) {
    await authStore.initUser()
  }
})
</script>

<style scoped>
.profile {
  min-height: 100vh;
  background-color: #f0f2f5;
}

.header {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 0;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  max-width: 800px;
  margin: 0 auto;
}

.header-content h1 {
  margin: 0;
  color: #1f2937;
  font-size: 18px;
}

.content {
  padding: 24px;
  max-width: 800px;
  margin: 0 auto;
  width: 100%;
}
</style> 