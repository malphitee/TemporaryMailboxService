<template>
  <div class="edit-profile">
    <a-layout>
      <a-layout-header class="header">
        <div class="header-content">
          <a-button type="text" @click="$router.back()">
            <ArrowLeftOutlined />
            返回
          </a-button>
          <h1>编辑资料</h1>
          <div></div>
        </div>
      </a-layout-header>
      
      <a-layout-content class="content">
        <a-card>
          <template #title>
            <EditOutlined />
            修改个人信息
          </template>
          
          <a-form
            :model="formData"
            :rules="rules"
            @finish="handleSubmit"
            layout="vertical"
            size="large"
            ref="formRef"
          >
            <a-form-item label="昵称" name="nickname">
              <a-input
                v-model:value="formData.nickname"
                placeholder="请输入昵称"
              >
                <template #prefix>
                  <UserOutlined />
                </template>
              </a-input>
            </a-form-item>
            
            <a-form-item label="头像URL" name="avatar">
              <a-input
                v-model:value="formData.avatar"
                placeholder="请输入头像URL（可选）"
              >
                <template #prefix>
                  <PictureOutlined />
                </template>
              </a-input>
            </a-form-item>
            
            <a-form-item label="时区" name="timezone">
              <a-select
                v-model:value="formData.timezone"
                placeholder="请选择时区"
              >
                <a-select-option value="Asia/Shanghai">Asia/Shanghai (北京时间)</a-select-option>
                <a-select-option value="UTC">UTC (协调世界时)</a-select-option>
                <a-select-option value="America/New_York">America/New_York (纽约时间)</a-select-option>
                <a-select-option value="Europe/London">Europe/London (伦敦时间)</a-select-option>
                <a-select-option value="Asia/Tokyo">Asia/Tokyo (东京时间)</a-select-option>
              </a-select>
            </a-form-item>
            
            <a-form-item label="语言" name="language">
              <a-select
                v-model:value="formData.language"
                placeholder="请选择语言"
              >
                <a-select-option value="zh-CN">简体中文</a-select-option>
                <a-select-option value="en-US">English</a-select-option>
                <a-select-option value="ja-JP">日本語</a-select-option>
              </a-select>
            </a-form-item>
            
            <a-form-item>
              <a-space>
                <a-button
                  type="primary"
                  html-type="submit"
                  :loading="loading"
                >
                  保存修改
                </a-button>
                <a-button @click="$router.back()">
                  取消
                </a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-card>
      </a-layout-content>
    </a-layout>
  </div>
</template>

<script setup lang="ts">
import { updateProfile } from '@/api/user'
import { useAuthStore } from '@/stores/auth'
import {
    ArrowLeftOutlined,
    EditOutlined,
    PictureOutlined,
    UserOutlined
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)
const formRef = ref()

const formData = reactive({
  nickname: '',
  avatar: '',
  timezone: '',
  language: ''
})

const rules = {
  nickname: [
    { required: true, message: '请输入昵称' },
    { min: 1, message: '昵称至少1位' },
    { max: 20, message: '昵称最多20位' },
  ],
  timezone: [
    { required: true, message: '请选择时区' },
  ],
  language: [
    { required: true, message: '请选择语言' },
  ],
}

const handleSubmit = async () => {
  try {
    await formRef.value?.validate()
    
    const updateData = {
      nickname: formData.nickname,
      avatar: formData.avatar,
      timezone: formData.timezone,
      language: formData.language
    }
    
    const response = await updateProfile(updateData)
    if (response.code === 0 && response.data) {
      authStore.updateUserInfo(response.data)
      message.success('个人资料更新成功！')
      await router.push('/profile')
    } else {
      message.error(response.message || '更新失败')
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '更新失败')
  }
}

onMounted(() => {
  if (authStore.user) {
    // 使用nickname字段填充表单
    formData.nickname = authStore.user.nickname || ''
    formData.avatar = authStore.user.avatar || ''
    formData.timezone = authStore.user.timezone || ''
    formData.language = authStore.user.language || ''
  }
})
</script>

<style scoped>
.edit-profile {
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