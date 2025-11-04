<template>
  <!-- 中文注释：我的页面骨架，展示昵称、ID、头像占位 -->
  <div class="p-4 space-y-4">
    <el-card shadow="never">
      <div class="flex items-center gap-3">
        <!-- 中文注释：与任务页保持一致，使用默认头像图片 -->
        <el-avatar :size="48" :src="defaultAvatar" />
        <div>
          <div class="font-semibold">用户昵称（占位）</div>
          <div class="text-gray-500 text-sm">用户 ID：—</div>
        </div>
      </div>
    </el-card>

    <!-- 中文注释：账号管理区域，提供退出登录按钮 -->
    <el-card shadow="never">
      <template #header>
        <span>账号管理</span>
      </template>
      <el-form label-width="160px">
        <el-form-item label="退出登录">
          <el-button type="danger" @click="onLogout">退出登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <span>系统设置</span>
      </template>
      <el-form label-width="160px">
        <el-form-item label="固定番茄钟页面">
          <el-switch v-model="fixed" @change="onFixedChange" />
        </el-form-item>
        <el-form-item label="快速入口">
          <!-- 中文注释：快速进入登录与注册占位页，便于体验流程 -->
          <el-button @click="toLogin">去登录</el-button>
          <el-button type="success" class="ml-2" @click="toRegister">去注册</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { useAppState } from '@/stores/appState'
import { computed } from 'vue'
import defaultAvatar from '@/assets/avatars/default.png'
import router from '@/router'
import { useAuth } from '@/stores/auth'

// 中文注释：固定番茄设置持久化到 appState（后续联动后端 UserSettings）
const store = useAppState()
const fixed = computed({
  get: () => store.tomato.fixedTomatoPage,
  set: (v: boolean) => store.updateTomato({ fixedTomatoPage: v })
})

function onFixedChange() {
  // 中文注释：此处可调用后端保存设置
}

// 中文注释：退出登录，清除认证信息并跳转到登录页
const auth = useAuth()
function onLogout() {
  try {
    auth.logout()
    store.reset()
    router.replace('/login')
  } catch (_) {
    location.reload()
  }
}

// 中文注释：导航到登录与注册占位页面
function toLogin() {
  router.push('/login')
}

function toRegister() {
  router.push('/register')
}
</script>

<style scoped>
</style>
