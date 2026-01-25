<template>
  <AuthShell
    brand-title="作业家"
    brand-subtitle="从今天开始，把计划变成结果"
    :tips="[
      '注册邮箱用于账号找回与通知',
      '密码需包含数字 + 大写 + 小写',
      '首次登录即可在设置中完善手机号'
    ]"
  >
    <div>
      <div class="text-xl font-extrabold text-slate-900 tracking-tight">注册</div>
      <div class="mt-1 text-xs text-slate-500">创建新账户，同步任务与积分</div>
    </div>

    <el-form ref="formRef" class="mt-6" :model="form" :rules="rules" label-position="top" @submit.prevent>
      <el-form-item label="用户名" prop="username">
        <el-input v-model="form.username" placeholder="请输入用户名（不少于4个字符）" />
      </el-form-item>
      <el-form-item label="邮箱" prop="email">
        <el-input v-model="form.email" placeholder="请输入邮箱" />
      </el-form-item>
      <el-form-item label="昵称（可选）" prop="nickname">
        <el-input v-model="form.nickname" placeholder="默认使用用户名" />
      </el-form-item>
      <el-form-item label="密码" prop="password">
        <el-input v-model="form.password" type="password" show-password placeholder="请输入密码" />
      </el-form-item>
      <div class="-mt-2 mb-2 text-[11px] text-slate-400">密码需包含数字、大写字母与小写字母</div>
      <el-form-item label="确认密码" prop="confirm">
        <el-input v-model="form.confirm" type="password" show-password placeholder="请再次输入密码" />
      </el-form-item>
      <el-button type="primary" class="mt-4 w-full !rounded-2xl !h-11 !font-extrabold" :loading="loading" @click="doRegister">注册</el-button>
      <div class="mt-5 flex justify-center text-sm text-slate-500">
        <button class="font-semibold hover:text-indigo-700" type="button" @click="toLogin">去登录</button>
      </div>
    </el-form>
  </AuthShell>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import type { FormInstance, FormRules } from 'element-plus/es/components/form'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { apiRegister } from '@/services/auth'
import AuthShell from '@/components/auth/AuthShell.vue'

const formRef = ref<FormInstance>()
const loading = ref(false)
const form = reactive({
  username: '',
  email: '',
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
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { validator: (_rule: any, value: string, cb: any) => {
      const v = (value || '').trim()
      if (!v || !v.includes('@') || !v.includes('.')) return cb(new Error('邮箱格式不正确'))
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
    await apiRegister(form.username.trim(), form.password, form.email.trim(), form.nickname.trim() || undefined)
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
