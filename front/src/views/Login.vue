<template>
  <div class="login-container">
    <div class="login-content">
      <!-- 左侧品牌区域 -->
      <div class="brand-section">
        <div class="brand-content">
          <h1 class="brand-title">临时邮箱系统</h1>
          <p class="brand-description">安全、便捷的临时邮箱服务</p>
        </div>
      </div>
      
      <!-- 右侧登录表单 -->
      <div class="login-form">
        <div class="login-header">
          <h2>欢迎回来</h2>
          <p>请登录您的账号</p>
        </div>
        
        <a-form
          :model="formData"
          :rules="rules"
          @finish="handleLogin"
          layout="vertical"
          size="large"
        >
          <a-form-item label="邮箱" name="email">
            <a-input
              v-model:value="formData.email"
              placeholder="请输入邮箱地址"
              type="email"
            >
              <template #prefix>
                <UserOutlined />
              </template>
            </a-input>
          </a-form-item>
          
          <a-form-item label="密码" name="password">
            <a-input-password
              v-model:value="formData.password"
              placeholder="请输入密码"
            >
              <template #prefix>
                <LockOutlined />
              </template>
            </a-input-password>
          </a-form-item>
          
          <a-form-item>
            <a-button
              type="primary"
              html-type="submit"
              :loading="authStore.loading"
              block
            >
              登录
            </a-button>
          </a-form-item>
          
          <div class="login-footer">
            <span>还没有账号？</span>
            <router-link to="/register">立即注册</router-link>
          </div>
        </a-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import type { LoginRequest } from '@/types/auth'
import { LockOutlined, UserOutlined } from '@ant-design/icons-vue'
import { reactive } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()

const formData = reactive<LoginRequest>({
  email: '',
  password: '',
})

const rules = {
  email: [
    { required: true, message: '请输入邮箱地址' },
    { type: 'email', message: '请输入有效的邮箱地址' },
  ],
  password: [
    { required: true, message: '请输入密码' },
    { min: 6, message: '密码至少6位' },
  ],
}

const handleLogin = async () => {
  const success = await authStore.login(formData)
  if (success) {
    router.push('/dashboard')
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-content {
  width: 95%;
  max-width: 1800px;
  min-height: 520px;
  background: white;
  border-radius: 18px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  display: flex;
  overflow: hidden;
}

.brand-section {
  width: 300px;
  flex-shrink: 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 36px 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.brand-content {
  text-align: center;
}

.brand-title {
  font-size: 2.2rem;
  font-weight: 600;
  margin-bottom: 16px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.brand-description {
  font-size: 1.05rem;
  opacity: 0.9;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.login-form {
  flex: 1;
  padding: 24px 80px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-width: 400px;
}

.login-form > form,
.login-form > .ant-form {
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  align-items: stretch;
}

.login-form .ant-form-item {
  width: 100%;
}

.login-form .ant-input,
.login-form .ant-input-password,
.login-form .ant-btn {
  width: 100% !important;
  box-sizing: border-box;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.login-header h2 {
  font-size: 1.875rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 8px 0;
}

.login-header p {
  font-size: 1rem;
  color: #6b7280;
  margin: 0;
}

.login-footer {
  text-align: center;
  margin-top: 24px;
  color: #6b7280;
}

.login-footer a {
  color: #1890ff;
  text-decoration: none;
  margin-left: 4px;
}

.login-footer a:hover {
  text-decoration: underline;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-content {
    flex-direction: column;
    min-height: auto;
    max-width: 480px;
  }

  .brand-section {
    width: auto;
    flex-shrink: 1;
    padding: 32px;
  }

  .login-form {
    flex: none;
    padding: 32px;
    max-width: none;
    min-width: auto;
  }
}

@media (max-width: 480px) {
  .login-container {
    padding: 0;
  }

  .login-content {
    border-radius: 0;
    min-height: 100vh;
    max-width: none;
  }

  .brand-section {
    padding: 24px;
  }

  .login-form {
    padding: 24px;
  }
}
</style> 