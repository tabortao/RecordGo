<template>
  <AuthShell brand-title="RecordGo" brand-subtitle="把计划变成结果，把努力变成积分">
    <div class="flex items-start justify-between gap-3">
      <div>
        <div class="text-xl font-extrabold text-slate-900 tracking-tight">登录</div>
        <div class="mt-1 text-xs text-slate-500">欢迎回来</div>
      </div>
    </div>

    <div class="mt-5 flex justify-center">
      <el-radio-group v-model="mode" size="small">
        <el-radio-button label="account">账号密码</el-radio-button>
        <el-radio-button label="token">令牌登录</el-radio-button>
      </el-radio-group>
    </div>

    <el-form class="mt-6" label-position="top" @submit.prevent>
      <template v-if="mode==='account'">
        <el-form-item label="用户名">
          <el-input v-model="username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="password" type="password" show-password placeholder="请输入密码" />
        </el-form-item>
        <div class="flex items-center justify-between gap-2">
          <div class="flex flex-col gap-1">
            <el-checkbox v-model="rememberMe">记住账号密码</el-checkbox>
            <div class="text-[11px] text-slate-400">仅保存在本机浏览器，本功能请自行评估安全风险</div>
          </div>
          <button class="text-sm font-semibold text-slate-500 hover:text-indigo-700" type="button" @click="toRegister">去注册</button>
        </div>
        <el-button type="primary" class="mt-5 w-full !rounded-2xl !h-11 !font-extrabold" :loading="loading" @click="doLogin">登录</el-button>
      </template>

      <template v-else>
        <el-form-item label="子账号令牌">
          <el-input v-model="token" placeholder="请输入子账号令牌" />
        </el-form-item>
        <div class="flex items-center justify-between gap-2">
          <div class="flex flex-col gap-1">
            <el-checkbox v-model="rememberMe">记住令牌</el-checkbox>
            <div class="text-[11px] text-slate-400">令牌将保存在本机浏览器，用于下次快速登录</div>
          </div>
          <button class="text-sm font-semibold text-slate-500 hover:text-indigo-700" type="button" @click="toRegister">去注册</button>
        </div>
        <el-button type="primary" class="mt-5 w-full !rounded-2xl !h-11 !font-extrabold" :loading="loading" @click="doTokenLogin">令牌登录</el-button>
      </template>

      <div class="mt-5 flex justify-center text-xs text-slate-400">
        <button class="hover:text-slate-600" type="button" @click="clearCache">清空缓存</button>
      </div>
    </el-form>
  </AuthShell>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { apiLogin, apiTokenLogin } from '@/services/auth'
import { useAuth } from '@/stores/auth'
import { useAppState } from '@/stores/appState'
import AuthShell from '@/components/auth/AuthShell.vue'

// 中文注释：登录模式与表单字段：账号密码 或 子账号令牌
const mode = ref<'account'|'token'>('account')
const username = ref('')
const password = ref('')
const token = ref('')
const auth = useAuth()
const appState = useAppState()
// 中文注释：记住我（默认未勾选），勾选则写入 localStorage；未勾选仅本次会话
const rememberMe = ref(false)
const loading = ref(false)

onMounted(() => {
  // 中文注释：若从注册页带回用户名，自动回填
  const url = new URL(location.href)
  const u = url.searchParams.get('u')
  if (u) username.value = u

  const remembered = readRememberedLogin()
  if (remembered?.mode) {
    mode.value = remembered.mode
    rememberMe.value = true
    if (remembered.mode === 'account') {
      if (!u) username.value = remembered.username || username.value
      password.value = remembered.password || ''
    } else {
      token.value = remembered.token || ''
    }
  }
})

type RememberedLogin =
  | { mode: 'account'; username: string; password: string }
  | { mode: 'token'; token: string }

const REMEMBER_KEY = 'login_remember_v1'
function safeB64Encode(s: string) {
  try { return btoa(unescape(encodeURIComponent(s))) } catch { return '' }
}
function safeB64Decode(s: string) {
  try { return decodeURIComponent(escape(atob(s))) } catch { return '' }
}
function readRememberedLogin(): RememberedLogin | null {
  try {
    const raw = localStorage.getItem(REMEMBER_KEY)
    if (!raw) return null
    const obj = JSON.parse(raw) as any
    if (!obj || typeof obj !== 'object') return null
    if (obj.mode === 'account') {
      return { mode: 'account', username: safeB64Decode(String(obj.u || '')), password: safeB64Decode(String(obj.p || '')) }
    }
    if (obj.mode === 'token') {
      return { mode: 'token', token: safeB64Decode(String(obj.t || '')) }
    }
    return null
  } catch {
    return null
  }
}
function writeRememberedLogin(v: RememberedLogin) {
  try {
    const payload =
      v.mode === 'account'
        ? { mode: 'account', u: safeB64Encode(v.username || ''), p: safeB64Encode(v.password || '') }
        : { mode: 'token', t: safeB64Encode(v.token || '') }
    localStorage.setItem(REMEMBER_KEY, JSON.stringify(payload))
  } catch {}
}
function clearRememberedLogin() {
  try { localStorage.removeItem(REMEMBER_KEY) } catch {}
}

async function doLogin() {
  if (loading.value) return
  if (!username.value || !password.value) {
    ElMessage.error('请输入用户名和密码')
    return
  }
  try {
    loading.value = true
    const resp = await apiLogin(username.value, password.value)
    auth.setLogin(resp.token, resp.user, rememberMe.value)
    if (rememberMe.value) writeRememberedLogin({ mode: 'account', username: username.value, password: password.value })
    else clearRememberedLogin()
    // 中文注释：登录后立即同步全局金币，避免进入页面前显示为 0 的闪烁或不一致
    try { appState.setCoins(Number(resp.user?.coins ?? 0)) } catch {}
    ElMessage.success('登录成功')
    if ((resp.user as any)?.must_change_password) {
      ElMessage.warning('请先修改密码以保障账号安全')
      router.replace('/settings/profile')
      return
    }
    const redirect = (router.currentRoute.value.query.redirect as string) || '/tasks'
    router.replace(redirect)
  } catch (e: any) {
    ElMessage.error(e?.message || '登录失败')
  } finally {
    loading.value = false
  }
}

async function doTokenLogin() {
  if (loading.value) return
  if (!token.value) { ElMessage.error('请输入令牌'); return }
  try {
    loading.value = true
    const resp = await apiTokenLogin(token.value)
    auth.setLogin(resp.token, resp.user, rememberMe.value)
    if (rememberMe.value) writeRememberedLogin({ mode: 'token', token: token.value })
    else clearRememberedLogin()
    // 中文注释：令牌登录同样同步全局金币
    try { appState.setCoins(Number(resp.user?.coins ?? 0)) } catch {}
    ElMessage.success('登录成功')
    const redirect = (router.currentRoute.value.query.redirect as string) || '/tasks'
    router.replace(redirect)
  } catch (e: any) {
    ElMessage.error(e?.message || '令牌登录失败')
  } finally {
    loading.value = false
  }
}

function clearCache() {
  localStorage.clear()
  ElMessage.success('已清空缓存')
}

function toRegister() {
  router.push('/register')
}
</script>

<style scoped>
/* 中文注释：页面样式由 Tailwind 类控制，这里无需额外样式 */
</style>
