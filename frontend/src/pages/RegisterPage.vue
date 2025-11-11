<template>
  <!-- 中文注释：现代化注册页面，Tailwind 渐变背景 + 居中卡片 -->
  <div class="min-h-screen bg-gradient-to-br from-sky-50 to-indigo-100 flex items-center justify-center p-4">
    <el-card class="w-full max-w-[420px] rounded-xl shadow-lg border" body-class="p-5">
      <template #header>
        <div class="text-center">
          <div class="text-xl font-semibold">任务积分助手 · 注册</div>
          <div class="text-xs text-gray-500 mt-1">创建新账户以同步任务与积分</div>
        </div>
      </template>
      <el-form label-position="top" @submit.prevent>
        <el-form-item label="用户名">
          <el-input v-model="username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="昵称（可选）">
          <el-input v-model="nickname" placeholder="默认使用用户名" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="确认密码">
          <el-input v-model="confirm" type="password" placeholder="请再次输入密码" />
        </el-form-item>
        <el-button type="primary" class="w-full" @click="doRegister">注册</el-button>
        <div class="flex justify-center mt-3 text-sm text-gray-500">
          <button class="hover:text-indigo-600" type="button" @click="toLogin">去登录</button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { apiRegister } from '@/services/auth'

// 中文注释：注册表单字段（与后端一致）
const username = ref('')
const nickname = ref('')
const password = ref('')
const confirm = ref('')

async function doRegister() {
  if (!username.value || !password.value || !confirm.value) {
    ElMessage.error('请完整填写注册信息')
    return
  }
  if (password.value !== confirm.value) {
    ElMessage.error('两次输入的密码不一致')
    return
  }
  try {
    await apiRegister(username.value, password.value, nickname.value || undefined)
    ElMessage.success('注册成功，请登录')
    router.push({ path: '/login', query: { u: username.value } })
  } catch (e: any) {
    ElMessage.error(e?.message || '注册失败')
  }
}

function toLogin() {
  router.push('/login')
}
</script>

<style scoped>
/* 中文注释：页面样式由 Tailwind 类控制，这里无需额外样式 */
</style>
