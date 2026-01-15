<template>
  <div class="min-h-screen relative overflow-hidden bg-gradient-to-br from-slate-950 via-slate-950 to-indigo-950">
    <div class="absolute -top-32 -left-32 h-[420px] w-[420px] rounded-full bg-indigo-600/30 blur-3xl" />
    <div class="absolute -bottom-40 -right-40 h-[520px] w-[520px] rounded-full bg-sky-500/20 blur-3xl" />
    <div class="absolute inset-0 bg-[radial-gradient(circle_at_30%_20%,rgba(255,255,255,0.08),transparent_40%),radial-gradient(circle_at_80%_80%,rgba(255,255,255,0.05),transparent_45%)]" />

    <div class="relative mx-auto flex min-h-screen w-full max-w-5xl items-center px-4 py-10">
      <div class="grid w-full grid-cols-1 gap-6 md:grid-cols-2">
        <div class="hidden md:flex flex-col justify-center rounded-2xl border border-white/10 bg-white/5 p-8 backdrop-blur">
          <div class="text-3xl font-semibold tracking-tight text-white">作业家</div>
          <div class="mt-2 text-sm text-white/70">用更轻的操作，养成更稳的习惯</div>
          <div class="mt-6 space-y-3 text-sm text-white/80">
            <div class="rounded-lg border border-white/10 bg-white/5 p-3">任务、心愿、积分一体化管理</div>
            <div class="rounded-lg border border-white/10 bg-white/5 p-3">子账号令牌登录，家庭协作更方便</div>
            <div class="rounded-lg border border-white/10 bg-white/5 p-3">登录后支持在“我的-设置”里修改密码</div>
          </div>
        </div>

        <div class="rounded-2xl border border-white/10 bg-white/95 p-6 shadow-2xl shadow-black/30 backdrop-blur md:p-8">
          <div class="flex items-start justify-between gap-3">
            <div>
              <div class="text-xl font-semibold text-slate-900">作业家 · 登录</div>
              <div class="mt-1 text-xs text-slate-500">欢迎回来</div>
            </div>
          </div>

          <div class="mt-5 flex justify-center">
            <el-radio-group v-model="mode" size="small">
              <el-radio-button label="account">账号密码</el-radio-button>
              <el-radio-button label="token">令牌登录</el-radio-button>
            </el-radio-group>
          </div>

          <el-form class="mt-5" label-position="top" @submit.prevent>
            <template v-if="mode==='account'">
              <el-form-item label="用户名">
                <el-input v-model="username" placeholder="请输入用户名" />
              </el-form-item>
              <el-form-item label="密码">
                <el-input v-model="password" type="password" show-password placeholder="请输入密码" />
              </el-form-item>
              <div class="flex items-center justify-between">
                <el-checkbox v-model="rememberMe">记住我</el-checkbox>
                <button class="text-sm text-slate-500 hover:text-indigo-600" type="button" @click="toRegister">去注册</button>
              </div>
              <el-button type="primary" class="mt-4 w-full" :loading="loading" @click="doLogin">登录</el-button>
            </template>

            <template v-else>
              <el-form-item label="子账号令牌">
                <el-input v-model="token" placeholder="请输入子账号令牌" />
              </el-form-item>
              <div class="flex items-center justify-between">
                <el-checkbox v-model="rememberMe">记住我</el-checkbox>
                <button class="text-sm text-slate-500 hover:text-indigo-600" type="button" @click="toRegister">去注册</button>
              </div>
              <el-button type="primary" class="mt-4 w-full" :loading="loading" @click="doTokenLogin">令牌登录</el-button>
            </template>

            <div class="mt-4 flex justify-center text-xs text-slate-400">
              <button class="hover:text-slate-600" type="button" @click="clearCache">清空缓存</button>
            </div>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { apiLogin, apiTokenLogin } from '@/services/auth'
import { useAuth } from '@/stores/auth'
import { useAppState } from '@/stores/appState'

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
})

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
