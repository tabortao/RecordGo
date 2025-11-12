<template>
  <!-- 中文注释：现代化登录页面，Tailwind 渐变背景 + 居中卡片 -->
  <div class="min-h-screen bg-gradient-to-br from-sky-50 to-indigo-100 flex items-center justify-center p-4">
    <el-card class="w-full max-w-[400px] rounded-xl shadow-lg border" body-class="p-5">
      <template #header>
        <div class="text-center">
          <div class="text-xl font-semibold">任务积分助手 · 登录</div>
        </div>
      </template>
      <el-form label-position="top" @submit.prevent>
        <div class="flex justify-center mb-2">
          <el-radio-group v-model="mode" size="small">
            <el-radio-button label="account">账号密码</el-radio-button>
            <el-radio-button label="token">令牌登录</el-radio-button>
          </el-radio-group>
        </div>

        <template v-if="mode==='account'">
          <el-form-item label="用户名">
            <el-input v-model="username" placeholder="请输入用户名" />
          </el-form-item>
          <el-form-item label="密码">
            <!-- 中文注释：启用 Element Plus 内置的小眼睛图标进行显示/隐藏切换 -->
            <el-input v-model="password" type="password" show-password placeholder="请输入密码" />
          </el-form-item>
          <div class="flex items-center justify-between mb-3">
            <el-checkbox v-model="rememberMe">记住我</el-checkbox>
            <button class="text-sm text-gray-500 hover:text-indigo-600" type="button" @click="toRegister">去注册</button>
          </div>
          <el-button type="primary" class="w-full" @click="doLogin">登录</el-button>
        </template>

        <template v-else>
          <el-form-item label="子账号令牌">
            <el-input v-model="token" placeholder="请输入子账号令牌" />
          </el-form-item>
          <div class="flex items-center justify-between mb-3">
            <el-checkbox v-model="rememberMe">记住我</el-checkbox>
            <button class="text-sm text-gray-500 hover:text-indigo-600" type="button" @click="toRegister">去注册</button>
          </div>
          <el-button type="primary" class="w-full" @click="doTokenLogin">令牌登录</el-button>
        </template>

        <div class="flex justify-center mt-3 text-xs text-gray-500">
          <button class="hover:text-gray-600" type="button" @click="clearCache">清空缓存</button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import router from '@/router'
import { apiLogin, apiTokenLogin } from '@/services/auth'
import { useAuth } from '@/stores/auth'

// 中文注释：登录模式与表单字段：账号密码 或 子账号令牌
const mode = ref<'account'|'token'>('account')
const username = ref('')
const password = ref('')
const token = ref('')
const auth = useAuth()
// 中文注释：记住我（默认未勾选），勾选则写入 localStorage；未勾选仅本次会话
const rememberMe = ref(false)

onMounted(() => {
  // 中文注释：若从注册页带回用户名，自动回填
  const url = new URL(location.href)
  const u = url.searchParams.get('u')
  if (u) username.value = u
})

async function doLogin() {
  if (!username.value || !password.value) {
    ElMessage.error('请输入用户名和密码')
    return
  }
  try {
    const resp = await apiLogin(username.value, password.value)
    auth.setLogin(resp.token, resp.user, rememberMe.value)
    ElMessage.success('登录成功')
    const redirect = (router.currentRoute.value.query.redirect as string) || '/tasks'
    router.replace(redirect)
  } catch (e: any) {
    ElMessage.error(e?.message || '登录失败')
  }
}

async function doTokenLogin() {
  if (!token.value) { ElMessage.error('请输入令牌'); return }
  try {
    const resp = await apiTokenLogin(token.value)
    auth.setLogin(resp.token, resp.user, rememberMe.value)
    ElMessage.success('登录成功')
    const redirect = (router.currentRoute.value.query.redirect as string) || '/tasks'
    router.replace(redirect)
  } catch (e: any) {
    ElMessage.error(e?.message || '令牌登录失败')
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
