<template>
  <div class="p-6 max-w-md mx-auto">
    <!-- 中文注释：登录卡片，占位实现，前端本地存储标记登录状态 -->
    <el-card>
      <template #header>
        <div class="font-semibold">登录</div>
      </template>
      <el-form label-width="80px" @submit.prevent>
        <el-form-item label="用户名">
          <el-input v-model="username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <div class="flex gap-2">
          <el-button type="primary" @click="doLogin">登录</el-button>
          <el-button @click="clearCache">清空缓存</el-button>
          <el-button type="success" @click="toRegister">去注册</el-button>
        </div>
        <div class="mt-3 text-gray-500 text-sm">提示：当前为本地占位登录，无后端校验。</div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import router from '@/router'

// 中文注释：简单的本地登录占位实现（后续可接入后端 JWT）
const username = ref('')
const password = ref('')

function doLogin() {
  if (!username.value || !password.value) {
    ElMessage.error('请输入用户名和密码')
    return
  }
  // 中文注释：本地存储登录标记与用户名，后续可替换为后端返回的 token
  localStorage.setItem('auth_logged', 'true')
  localStorage.setItem('auth_username', username.value)
  ElMessage.success('登录成功（本地占位）')
  router.push('/tasks')
}

function clearCache() {
  // 中文注释：清空前端缓存，用于调试
  localStorage.clear()
  ElMessage.success('已清空缓存')
}

function toRegister() {
  router.push('/register')
}
</script>

<style scoped>
/* 中文注释：基础内边距与宽度样式 */
</style>

