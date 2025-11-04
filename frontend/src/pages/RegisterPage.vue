<template>
  <div class="min-h-screen bg-white flex items-center justify-center">
    <!-- 中文注释：简洁美观的注册卡片，字段与数据库一致 -->
    <el-card class="w-[360px] shadow">
      <template #header>
        <div class="font-semibold text-center">任务积分助手 · 注册</div>
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
        <div class="flex justify-between mt-3 text-sm text-gray-500">
          <button class="hover:text-green-600" type="button" @click="toLogin">去登录</button>
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
/* 中文注释：页面基础样式占位 */
</style>
