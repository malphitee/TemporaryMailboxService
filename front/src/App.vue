<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'

const authStore = useAuthStore()

// 应用启动时恢复用户状态
onMounted(() => {
  if (authStore.token && !authStore.user) {
    authStore.restoreFromStorage()
    authStore.initUser()
  }
})
</script>

<template>
  <RouterView />
</template>

<style>
* {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
}

#app {
  min-height: 100vh;
}
</style>
