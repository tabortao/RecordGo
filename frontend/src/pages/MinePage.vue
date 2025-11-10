<template>
  <!-- 中文注释：我的页面骨架，展示昵称、ID、头像占位 -->
  <div class="p-4 space-y-4">
    <el-card shadow="never">
      <div class="flex items-center gap-3">
        <!-- 中文注释：头像优先显示用户自定义头像，未设置则使用默认头像；修复相对路径未加 API 前缀导致无法显示的问题 -->
        <el-avatar :size="48" :src="avatarSrc" />
        <div>
          <!-- 中文注释：昵称优先显示真实昵称，未设置则显示用户名 -->
          <div class="font-semibold">{{ displayName }}</div>
          <!-- 中文注释：用户ID格式化为6位数，不足左侧补0 -->
          <div class="text-gray-500 text-sm">用户 ID：{{ displayId }}</div>
        </div>
      </div>
    </el-card>

    <!-- 中文注释：账号管理 UI（shadcn 风格 + tailwind），包含编辑个人信息/子账号管理/退出登录 -->
    <div class="rounded-lg border bg-white">
      <div class="px-4 py-3 flex items-center gap-2">
        <el-icon :size="18" style="color:#3b82f6"><User /></el-icon>
        <span class="font-semibold">账号管理</span>
      </div>
      <!-- 响应式网格：移动端单列，桌面端多列 -->
      <div class="px-2 py-2 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
        <!-- 编辑个人信息：改为跳转到独立页面 -->
        <button class="w-full flex items-center gap-3 px-3 py-2 rounded-md hover:bg-gray-50 transition" @click="router.push('/settings/profile')">
          <el-icon :size="18" style="color:#60a5fa"><Edit /></el-icon>
          <span class="text-gray-800">编辑个人信息</span>
        </button>
        <!-- 子账号管理（占位） -->
        <button class="w-full flex items-center gap-3 px-3 py-2 rounded-md hover:bg-gray-50 transition" @click="onChildManage">
          <el-icon :size="18" style="color:#22c55e"><User /></el-icon>
          <span class="text-gray-800">子账号管理</span>
        </button>
        <!-- 退出登录 -->
        <button class="w-full flex items-center gap-3 px-3 py-2 rounded-md hover:bg-gray-50 transition" @click="onLogout">
          <el-icon :size="18" style="color:#ef4444"><SwitchButton /></el-icon>
          <span class="text-gray-800">退出登录</span>
        </button>
      </div>
    </div>

    <!-- 中文注释：按最新需求，取消账号管理下方的“系统设置”卡片 -->

    <!-- 中文注释：编辑个人信息对话框已移除，改为独立页面 /settings/profile -->

    <!-- 中文注释：设置模块（与账号管理同级），展示各设置按钮 -->
    <div class="rounded-lg border bg-white">
      <div class="px-4 py-3 flex items-center gap-2">
        <el-icon :size="18" style="color:#0ea5e9"><Setting /></el-icon>
        <span class="font-semibold">设置</span>
      </div>
      <!-- 响应式网格：移动端单列，桌面端多列 -->
      <div class="px-2 py-2 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
        <button v-for="i in settingItems" :key="i.key" class="w-full flex items-center gap-3 px-3 py-2 rounded-md hover:bg-gray-50 transition text-left" @click="goSettingsTab(i.key)">
          <el-icon :size="18" :style="{ color: i.fg }"><component :is="i.icon" /></el-icon>
          <span class="text-gray-800">{{ i.label }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAppState } from '@/stores/appState'
import defaultAvatar from '@/assets/avatars/default.png'
import router from '@/router'
import { useAuth } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import { User, Edit, SwitchButton, Setting, Timer, List, Microphone, Notebook, Coin, InfoFilled } from '@element-plus/icons-vue'

// 中文注释：应用状态（用于退出登录时重置）
const store = useAppState()

// 中文注释：退出登录，清除认证信息并跳转到登录页
const auth = useAuth()

// 中文注释：展示名称逻辑——优先真实昵称，缺省则回退到用户名
const displayName = computed(() => {
  const u = auth.user
  if (!u) return '未登录'
  const n = (u.nickname || '').trim()
  return n ? n : u.username
})

// 中文注释：展示ID为6位数字字符串（不足位数左侧补0）
const displayId = computed(() => {
  const u = auth.user
  if (!u) return '------'
  return String(u.id).padStart(6, '0')
})

// 中文注释：头像地址优先使用用户头像路径，未设置则用默认头像
// 中文注释：头像地址解析——支持完整 URL；后端返回的相对路径（如 uploads/images/avatars/...）统一加上 /api 前缀通过 Vite 代理
function resolveAvatarUrl(p?: string | null) {
  if (!p) return defaultAvatar
  const s = String(p)
  // 中文注释：仅当为完整 URL 或包含 uploads 路径时才走后端；否则回退到内置默认头像
  if (/^https?:\/\//i.test(s)) return s
  if (!/uploads\//i.test(s)) return defaultAvatar
  return `/api/${s}`.replace(/\/+/g, '/').replace(/\/$/, '')
}

const avatarSrc = computed(() => resolveAvatarUrl(auth.user?.avatar_path))

// 中文注释：编辑个人信息改为独立页面，移除弹窗相关状态与函数

function onChildManage() {
  ElMessage.info('子账号管理功能将在后续版本提供')
}

function onLogout() {
  try {
    auth.logout()
    store.reset()
    router.replace('/login')
  } catch (_) {
    location.reload()
  }
}

// （移除未使用的登录/注册跳转函数，清理诊断警告）

// 中文注释：移除系统设置按钮对应的旧跳转函数

// 中文注释：设置模块按钮（图标统一 18，与“编辑个人信息”一致）
type SettingsKey = 'tomato' | 'tasks' | 'reading' | 'subjects' | 'coins' | 'about'
const settingItems: Array<{ key: SettingsKey; label: string; icon: any; fg: string }> = [
  { key: 'tomato', label: '番茄钟设置', icon: Timer, fg: '#ef4444' },
  { key: 'tasks', label: '任务设置', icon: List, fg: '#10b981' },
  { key: 'reading', label: '朗读设置', icon: Microphone, fg: '#7c3aed' },
  { key: 'subjects', label: '任务分类设置', icon: Notebook, fg: '#2563eb' },
  { key: 'coins', label: '金币设置', icon: Coin, fg: '#f59e0b' },
  { key: 'about', label: '关于', icon: InfoFilled, fg: '#0ea5e9' }
]

function goSettingsTab(k: SettingsKey) {
  // 中文注释：跳转到独立的设置页面路径，而不是通用 /settings
  const map: Record<SettingsKey, string> = {
    tomato: '/settings/tomato',
    tasks: '/settings/tasks',
    reading: '/settings/reading',
    subjects: '/settings/subjects',
    coins: '/settings/coins',
    about: '/settings/about'
  }
  router.push(map[k])
}

// （移除我的页内的设置入口列表，保留“系统设置”按钮跳转到独立页面）
</script>

<style scoped>
</style>
