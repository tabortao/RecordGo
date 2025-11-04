<template>
  <div class="min-h-screen bg-white flex items-center justify-center">
    <!-- 中文注释：简洁美观的登录卡片，匹配后端数据库字段 -->
    <el-card class="w-[360px] shadow">
      <template #header>
        <div class="font-semibold text-center">任务积分助手 · 登录</div>
      </template>
      <el-form label-position="top" @submit.prevent>
        <el-form-item label="用户名">
          <el-input v-model="username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-button type="primary" class="w-full" @click="doLogin">登录</el-button>
        <div class="flex justify-between mt-3 text-sm text-gray-500">
          <button class="hover:text-green-600" type="button" @click="toRegister">去注册</button>
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
import { apiLogin } from '@/services/auth'
import { useAuth } from '@/stores/auth'

// 中文注释：登录表单字段（与后端一致）
const username = ref('')
const password = ref('')
const auth = useAuth()

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
    auth.setLogin(resp.token, resp.user)
    ElMessage.success('登录成功')
    const redirect = (router.currentRoute.value.query.redirect as string) || '/tasks'
    router.replace(redirect)
  } catch (e: any) {
    ElMessage.error(e?.message || '登录失败')
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
/* 中文注释：基础内边距与宽度样式 */
</style>
