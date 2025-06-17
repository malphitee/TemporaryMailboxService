<template>
  <div class="change-password">
    <a-layout>
      <a-layout-header class="header">
        <div class="header-content">
          <a-button type="text" @click="$router.back()">
            <ArrowLeftOutlined />
            返回
          </a-button>
          <h1>修改密码</h1>
          <div></div>
        </div>
      </a-layout-header>
      
      <a-layout-content class="content">
        <a-card>
          <template #title>
            <LockOutlined />
            安全设置
          </template>
          
          <a-form
            :model="formData"
            :rules="rules"
            @finish="handleSubmit"
            layout="vertical"
            size="large"
          >
            <a-form-item label="当前密码" name="current_password">
              <a-input-password
                v-model:value="formData.current_password"
                placeholder="请输入当前密码"
              >
                <template #prefix>
                  <LockOutlined />
                </template>
              </a-input-password>
            </a-form-item>
            
            <a-form-item label="新密码" name="new_password">
              <a-input-password
                v-model:value="formData.new_password"
                placeholder="请输入新密码（至少6位）"
              >
                <template #prefix>
                  <SafetyOutlined />
                </template>
              </a-input-password>
            </a-form-item>
            
            <a-form-item label="确认新密码" name="confirm_password">
              <a-input-password
                v-model:value="confirmPassword"
                placeholder="请再次输入新密码"
              >
                <template #prefix>
                  <SafetyOutlined />
                </template>
              </a-input-password>
            </a-form-item>
            
            <a-form-item>
              <a-space>
                <a-button
                  type="primary"
                  html-type="submit"
                  :loading="userStore.loading"
                >
                  修改密码
                </a-button>
                <a-button @click="$router.back()">
                  取消
                </a-button>
              </a-space>
            </a-form-item>
          </a-form>
          
          <a-divider />
          
          <a-alert
            message="安全提示"
            description="为了您的账户安全，建议定期更换密码，并使用包含字母、数字和特殊字符的强密码。"
            type="info"
            show-icon
          />
        </a-card>
      </a-layout-content>
    </a-layout>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from '@/stores/user'
import type { ChangePasswordRequest } from '@/types/auth'
import {
    ArrowLeftOutlined,
    LockOutlined,
    SafetyOutlined
} from '@ant-design/icons-vue'
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const userStore = useUserStore()

const formData = reactive<ChangePasswordRequest>({
  current_password: '',
  new_password: '',
})

const confirmPassword = ref('')

const rules = {
  current_password: [
    { required: true, message: '请输入当前密码' },
  ],
  new_password: [
    { required: true, message: '请输入新密码' },
    { min: 6, message: '新密码至少6位' },
  ],
  confirm_password: [
    { required: true, message: '请确认新密码' },
    {
      validator: (_: any, value: string) => {
        if (value !== formData.new_password) {
          return Promise.reject('两次输入的密码不一致')
        }
        return Promise.resolve()
      },
    },
  ],
}

const handleSubmit = async () => {
  const success = await userStore.changePassword(formData)
  if (success) {
    // 清空表单
    formData.current_password = ''
    formData.new_password = ''
    confirmPassword.value = ''
    router.back()
  }
}
</script>

<style scoped>
.change-password {
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