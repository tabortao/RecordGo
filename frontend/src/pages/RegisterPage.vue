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
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent>
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" placeholder="请输入用户名（不少于4个字符）" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="昵称（可选）" prop="nickname">
          <el-input v-model="form.nickname" placeholder="默认使用用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" show-password placeholder="请输入密码" />
        </el-form-item>
        <div class="-mt-2 mb-2 text-xs text-gray-500">密码需包含数字、大写字母与小写字母</div>
        <el-form-item label="确认密码" prop="confirm">
          <el-input v-model="form.confirm" type="password" show-password placeholder="请再次输入密码" />
        </el-form-item>
        <el-button type="primary" class="w-full" :loading="loading" @click="doRegister">注册</el-button>
        <div class="flex justify-center mt-3 text-sm text-gray-500">
          <button class="hover:text-indigo-600" type="button" @click="toLogin">去登录</button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { apiRegister } from '@/services/auth'

const formRef = ref<FormInstance>()
const loading = ref(false)
const form = reactive({
  username: '',
  phone: '',
  nickname: '',
  password: '',
  confirm: ''
})

function isStrongPassword(pw: string) {
  const hasUpper = /[A-Z]/.test(pw)
  const hasLower = /[a-z]/.test(pw)
  const hasDigit = /\d/.test(pw)
  return hasUpper && hasLower && hasDigit
}

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { validator: (_rule: any, value: string, cb: any) => {
      if ((value || '').trim().length < 4) return cb(new Error('用户名不少于4个字符'))
      cb()
    }, trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { validator: (_rule: any, value: string, cb: any) => {
      const v = (value || '').trim()
      if (!/^1\d{10}$/.test(v)) return cb(new Error('手机号格式不正确'))
      cb()
    }, trigger: 'blur' }
  ],
  nickname: [],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { validator: (_rule: any, value: string, cb: any) => {
      if (!isStrongPassword(value || '')) return cb(new Error('密码需包含数字、大写字母与小写字母'))
      cb()
    }, trigger: 'blur' }
  ],
  confirm: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: (_rule: any, value: string, cb: any) => {
      if ((value || '') !== (form.password || '')) return cb(new Error('两次输入的密码不一致'))
      cb()
    }, trigger: 'blur' }
  ]
}

async function doRegister() {
  if (loading.value) return
  const ok = await formRef.value?.validate().catch(() => false)
  if (!ok) {
    return
  }
  try {
    loading.value = true
    await apiRegister(form.username.trim(), form.password, form.phone.trim(), form.nickname.trim() || undefined)
    ElMessage.success('注册成功，请登录')
    router.push({ path: '/login', query: { u: form.username.trim() } })
  } catch (e: any) {
    ElMessage.error(e?.message || '注册失败')
  } finally {
    loading.value = false
  }
}

function toLogin() {
  router.push('/login')
}
</script>

<style scoped>
/* 中文注释：页面样式由 Tailwind 类控制，这里无需额外样式 */
</style>
