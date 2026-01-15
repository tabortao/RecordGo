<template>
  <div class="min-h-screen relative overflow-hidden bg-gradient-to-br from-slate-950 via-slate-950 to-indigo-950">
    <div class="absolute -top-32 -left-32 h-[420px] w-[420px] rounded-full bg-indigo-600/30 blur-3xl" />
    <div class="absolute -bottom-40 -right-40 h-[520px] w-[520px] rounded-full bg-sky-500/20 blur-3xl" />
    <div class="absolute inset-0 bg-[radial-gradient(circle_at_30%_20%,rgba(255,255,255,0.08),transparent_40%),radial-gradient(circle_at_80%_80%,rgba(255,255,255,0.05),transparent_45%)]" />

    <div class="relative mx-auto flex min-h-screen w-full max-w-5xl items-center px-4 py-10">
      <div class="grid w-full grid-cols-1 gap-6 md:grid-cols-2">
        <div class="hidden md:flex flex-col justify-center rounded-2xl border border-white/10 bg-white/5 p-8 backdrop-blur">
          <div class="text-3xl font-semibold tracking-tight text-white">作业家</div>
          <div class="mt-2 text-sm text-white/70">从今天开始，把计划变成结果</div>
          <div class="mt-6 space-y-3 text-sm text-white/80">
            <div class="rounded-lg border border-white/10 bg-white/5 p-3">注册邮箱用于账号找回与通知</div>
            <div class="rounded-lg border border-white/10 bg-white/5 p-3">密码需包含数字 + 大写 + 小写</div>
            <div class="rounded-lg border border-white/10 bg-white/5 p-3">首次登录即可在设置中完善手机号</div>
          </div>
        </div>

        <div class="rounded-2xl border border-white/10 bg-white/95 p-6 shadow-2xl shadow-black/30 backdrop-blur md:p-8">
          <div>
            <div class="text-xl font-semibold text-slate-900">作业家 · 注册</div>
            <div class="mt-1 text-xs text-slate-500">创建新账户，同步作业与积分</div>
          </div>

          <el-form ref="formRef" class="mt-5" :model="form" :rules="rules" label-position="top" @submit.prevent>
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
            <div class="-mt-2 mb-2 text-xs text-slate-500">密码需包含数字、大写字母与小写字母</div>
            <el-form-item label="确认密码" prop="confirm">
              <el-input v-model="form.confirm" type="password" show-password placeholder="请再次输入密码" />
            </el-form-item>
            <el-button type="primary" class="mt-2 w-full" :loading="loading" @click="doRegister">注册</el-button>
            <div class="mt-4 flex justify-center text-sm text-slate-500">
              <button class="hover:text-indigo-600" type="button" @click="toLogin">去登录</button>
            </div>
          </el-form>
        </div>
      </div>
    </div>
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
