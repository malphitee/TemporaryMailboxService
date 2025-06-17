<template>
  <div class="register-container">
    <div class="register-content">
      <!-- 左侧品牌区 -->
      <div class="brand-section">
        <div class="brand-content">
          <h1 class="brand-title">临时邮箱系统</h1>
          <p class="brand-description">安全、便捷的临时邮箱服务</p>
        </div>
      </div>
      <!-- 右侧注册表单区 -->
      <div class="register-form">
        <div class="register-header">
          <h2>注册</h2>
          <p>创建您的临时邮箱账号</p>
        </div>
        <a-form
          ref="formRef"
          :model="formData"
          :rules="rules"
          @finish="handleSubmit"
          layout="vertical"
          size="large"
          :validateTrigger="'submit'"
        >
          <a-form-item label="用户名" name="username">
            <a-input
              v-model:value="formData.username"
              placeholder="请输入用户名"
            >
              <template #prefix>
                <UserOutlined />
              </template>
            </a-input>
          </a-form-item>
          <a-form-item label="邮箱" name="email">
            <a-input
              v-model:value="formData.email"
              placeholder="请输入邮箱地址"
              type="email"
            >
              <template #prefix>
                <MailOutlined />
              </template>
            </a-input>
          </a-form-item>
          <a-form-item label="昵称" name="nickname">
            <a-input
              v-model:value="formData.nickname"
              placeholder="请输入昵称"
            >
              <template #prefix>
                <ContactsOutlined />
              </template>
            </a-input>
          </a-form-item>
          <a-form-item label="密码" name="password">
            <a-input-password
              v-model:value="formData.password"
              placeholder="请输入密码（至少6位）"
            >
              <template #prefix>
                <LockOutlined />
              </template>
            </a-input-password>
          </a-form-item>
          <a-form-item label="确认密码" name="confirmPassword">
            <a-input-password
              v-model:value="formData.confirmPassword"
              placeholder="请再次输入密码"
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
              注册
            </a-button>
          </a-form-item>
          <div class="register-footer">
            <span>已有账号？</span>
            <router-link to="/login">立即登录</router-link>
          </div>
        </a-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import {
  ContactsOutlined,
  LockOutlined,
  MailOutlined,
  UserOutlined
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()

const formData = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  nickname: ''
})

const formRef = ref()

const rules = {
  username: [
    { required: true, message: '请输入用户名' },
    { min: 3, message: '用户名至少3位' },
    { max: 20, message: '用户名最多20位' },
  ],
  email: [
    { required: true, message: '请输入邮箱地址' },
    { type: 'email', message: '请输入有效的邮箱地址' },
  ],
  nickname: [
    { required: true, message: '请输入昵称' },
    { min: 1, message: '昵称至少1位' },
    { max: 20, message: '昵称最多20位' },
  ],
  password: [
    { required: true, message: '请输入密码' },
    { min: 6, message: '密码至少6位' },
  ],
  confirmPassword: [
    { required: true, message: '请确认密码' },
    {
      validator: (_: any, value: string) => {
        if (value !== formData.password) {
          return Promise.reject('两次输入的密码不一致')
        }
        return Promise.resolve()
      }
    },
  ],
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    
    const registerData = {
      username: formData.username,
      email: formData.email,
      password: formData.password,
      nickname: formData.nickname || ''
    }
    
    await authStore.register(registerData)
    
    message.success('注册成功！')
    await router.push('/dashboard')
  } catch (error) {
    console.error('注册失败:', error)
  }
}
</script>

<style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.register-content {
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

.register-form {
  flex: 1;
  padding: 24px 80px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-width: 400px;
}

.register-form > form,
.register-form > .ant-form {
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  align-items: stretch;
}

.register-form .ant-form-item {
  width: 100%;
}

.register-form .ant-input,
.register-form .ant-input-password,
.register-form .ant-btn {
  width: 100% !important;
  box-sizing: border-box;
}

.register-header {
  text-align: center;
  margin-bottom: 32px;
}

.register-header h2 {
  margin: 0 0 8px 0;
  font-size: 1.875rem;
  color: #1f2937;
  font-weight: 600;
}

.register-header p {
  margin: 0;
  color: #6b7280;
  font-size: 1rem;
}

.register-footer {
  text-align: center;
  margin-top: 16px;
}

.register-footer a {
  color: #1890ff;
  text-decoration: none;
}

.register-footer a:hover {
  text-decoration: underline;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .register-content {
    flex-direction: column;
    min-height: auto;
    max-width: 480px;
  }

  .brand-section {
    flex: none;
    padding: 32px;
  }

  .register-form {
    flex: none;
    padding: 32px;
    max-width: none;
  }
}

@media (max-width: 480px) {
  .register-container {
    padding: 0;
  }

  .register-content {
    border-radius: 0;
    min-height: 100vh;
    max-width: none;
  }

  .brand-section {
    padding: 24px;
  }

  .register-form {
    padding: 24px;
  }
}
</style> 