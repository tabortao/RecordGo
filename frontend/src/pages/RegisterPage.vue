<template>
  <div class="p-6 max-w-md mx-auto">
    <!-- 中文注释：注册卡片，占位实现，仅在前端展示与存储 -->
    <el-card>
      <template #header>
        <div class="font-semibold">注册</div>
      </template>
      <el-form label-width="80px" @submit.prevent>
        <el-form-item label="用户名">
          <el-input v-model="username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="确认密码">
          <el-input v-model="confirm" type="password" placeholder="请再次输入密码" />
        </el-form-item>
        <div class="flex gap-2">
          <el-button type="primary" @click="doRegister">注册</el-button>
          <el-button @click="toLogin">去登录</el-button>
        </div>
        <div class="mt-3 text-gray-500 text-sm">提示：当前为本地占位注册，无后端保存。</div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import router from '@/router'

// 中文注释：注册占位逻辑（后续可接入后端注册接口）
const username = ref('')
const password = ref('')
const confirm = ref('')

function doRegister() {
  if (!username.value || !password.value || !confirm.value) {
    ElMessage.error('请完整填写注册信息')
    return
  }
  if (password.value !== confirm.value) {
    ElMessage.error('两次输入的密码不一致')
    return
  }
  // 中文注释：模拟注册成功，将用户名写入本地存储并跳转登录
  localStorage.setItem('register_username', username.value)
  ElMessage.success('注册成功（本地占位），请登录')
  router.push('/login')
}

function toLogin() {
  router.push('/login')
}
</script>

<style scoped>
/* 中文注释：页面基础样式占位 */
</style>

