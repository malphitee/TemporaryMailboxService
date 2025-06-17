<template>
  <div class="dashboard">
    <a-layout>
      <a-layout-header class="header">
        <div class="header-content">
          <h1>临时邮箱系统</h1>
          <div class="user-info">
            <a-avatar :size="32" style="margin-right: 8px">
              {{ displayName?.charAt(0) || 'U' }}
            </a-avatar>
            <a-dropdown>
              <span class="user-name">{{ displayName }}</span>
              <template #overlay>
                <a-menu>
                  <a-menu-item key="profile" @click="$router.push('/profile')">
                    <UserOutlined />
                    个人资料
                  </a-menu-item>
                  <a-menu-divider />
                  <a-menu-item key="logout" @click="handleLogout">
                    <LogoutOutlined />
                    退出登录
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
        </div>
      </a-layout-header>
      
      <a-layout-content class="content">
        <div class="welcome-section">
          <a-card>
            <h2>欢迎回来，{{ displayName }}！</h2>
            <p>您的临时邮箱系统已准备就绪。</p>
          </a-card>
        </div>
        
        <a-row :gutter="[16, 16]" class="stats-section">
          <a-col :xs="24" :sm="12" :md="6">
            <a-card class="stat-card">
              <a-statistic
                title="域名数量"
                :value="0"
                prefix="🌐"
              />
            </a-card>
          </a-col>
          <a-col :xs="24" :sm="12" :md="6">
            <a-card class="stat-card">
              <a-statistic
                title="临时邮箱"
                :value="0"
                prefix="📧"
              />
            </a-card>
          </a-col>
          <a-col :xs="24" :sm="12" :md="6">
            <a-card class="stat-card">
              <a-statistic
                title="收到邮件"
                :value="0"
                prefix="📬"
              />
            </a-card>
          </a-col>
          <a-col :xs="24" :sm="12" :md="6">
            <a-card class="stat-card">
              <a-statistic
                title="在线天数"
                :value="1"
                prefix="📅"
              />
            </a-card>
          </a-col>
        </a-row>
        
        <a-row :gutter="[16, 16]" class="action-section">
          <a-col :xs="24" :md="12">
            <a-card title="快速操作" class="action-card">
              <a-space direction="vertical" style="width: 100%">
                <a-button type="primary" block disabled>
                  <PlusOutlined />
                  添加域名
                </a-button>
                <a-button block disabled>
                  <MailOutlined />
                  创建临时邮箱
                </a-button>
                <a-button block @click="$router.push('/profile')">
                  <SettingOutlined />
                  账户设置
                </a-button>
              </a-space>
            </a-card>
          </a-col>
          <a-col :xs="24" :md="12">
            <a-card title="系统状态" class="status-card">
              <a-space direction="vertical" style="width: 100%">
                <div class="status-item">
                  <span>邮件服务</span>
                  <a-badge status="processing" text="运行中" />
                </div>
                <div class="status-item">
                  <span>DNS服务</span>
                  <a-badge status="success" text="正常" />
                </div>
                <div class="status-item">
                  <span>数据库</span>
                  <a-badge status="success" text="正常" />
                </div>
              </a-space>
            </a-card>
          </a-col>
        </a-row>
      </a-layout-content>
    </a-layout>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import {
    LogoutOutlined,
    MailOutlined,
    PlusOutlined,
    SettingOutlined,
    UserOutlined
} from '@ant-design/icons-vue'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const displayName = computed(() => {
  const user = authStore.user
  if (!user) return '未知用户'
  
  return user.nickname || user.username || '未知用户'
})
</script>

<style scoped>
.dashboard {
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
  max-width: 1200px;
  margin: 0 auto;
}

.header-content h1 {
  margin: 0;
  color: #1890ff;
  font-size: 20px;
}

.user-info {
  display: flex;
  align-items: center;
}

.user-name {
  cursor: pointer;
  color: #666;
}

.user-name:hover {
  color: #1890ff;
}

.content {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.welcome-section {
  margin-bottom: 24px;
}

.welcome-section h2 {
  margin: 0 0 8px 0;
  color: #1f2937;
}

.welcome-section p {
  margin: 0;
  color: #6b7280;
}

.stats-section {
  margin-bottom: 24px;
}

.stat-card {
  text-align: center;
}

.action-section .action-card,
.action-section .status-card {
  height: 100%;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style> 