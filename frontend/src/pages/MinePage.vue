<template>
  <div class="fixed top-0 left-0 right-0 bg-white/80 dark:bg-gray-800/80 backdrop-blur z-40 border-b border-gray-200 dark:border-gray-700">
    <div class="px-4 py-2 font-semibold">我的</div>
  </div>
  <div class="h-10"></div>
  <div class="p-3 space-y-3 pb-24" style="padding-bottom: calc(env(safe-area-inset-bottom) + 96px)">
    <el-card shadow="never">
      <div class="flex items-center gap-3">
        <!-- 中文注释：头像优先显示用户自定义头像，未设置则使用默认头像；修复相对路径未加 API 前缀导致无法显示的问题 -->
        <el-avatar :size="48" :src="avatarSrc" />
        <div>
          <!-- 中文注释：昵称优先显示真实昵称，未设置则显示用户名 -->
          <div class="font-semibold">{{ displayName }}</div>
          <!-- 中文注释：VIP标识（替换原用户ID显示） -->
          <div class="mt-1">
            <VipBadge :user="auth.user" />
          </div>
        </div>
      </div>
    </el-card>

    <!-- 中文注释：账号管理 UI（shadcn 风格 + tailwind），包含编辑个人信息/子账号管理/退出登录 -->
    <div class="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
      <div class="px-3 py-2 flex items-center gap-2">
        <el-icon :size="18" style="color:#3b82f6"><User /></el-icon>
        <span class="font-semibold dark:text-gray-100">账号管理</span>
      </div>
      <!-- 响应式网格：移动端单列，桌面端多列 -->
      <div class="px-2 py-2 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-1">
        <!-- 编辑个人信息：子账号隐藏，仅主账号显示 -->
        <button v-if="isParent" class="w-full flex items-center gap-3 px-3 py-2 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 transition" @click="router.push('/settings/profile')">
          <el-icon :size="18" style="color:#60a5fa"><Edit /></el-icon>
          <span class="text-gray-800 dark:text-gray-100">编辑个人信息</span>
        </button>
        <!-- 子账号管理：无权限时不显示（仅父账号或具备 manage_children 权限显示） -->
        <button v-if="isParent || manageChildren" class="w-full flex items-center gap-3 px-3 py-2 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 transition" @click="onChildManage" title="仅父账号或具备权限的账号可管理子账号">
          <el-icon :size="18" style="color:#22c55e"><User /></el-icon>
          <span class="text-gray-800 dark:text-gray-100">子账号管理</span>
        </button>
        <!-- 退出登录 -->
        <button class="w-full flex items-center gap-3 px-3 py-2 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 transition" @click="onLogout">
          <el-icon :size="18" style="color:#ef4444"><SwitchButton /></el-icon>
          <span class="text-gray-800 dark:text-gray-100">退出登录</span>
        </button>
      </div>
    </div>

    <!-- 中文注释：按最新需求，取消账号管理下方的“系统设置”卡片 -->

    <!-- 中文注释：编辑个人信息对话框已移除，改为独立页面 /settings/profile -->

    <!-- 中文注释：数据管理模块置于设置模块之前（更靠上显示） -->
    <div class="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
      <div class="px-3 py-2 flex items-center gap-2">
        <el-icon :size="18" style="color:#0ea5e9"><Setting /></el-icon>
        <span class="font-semibold dark:text-gray-100">数据管理</span>
      </div>
      <div class="px-2 py-2 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-1">
        <button class="w-full flex items-center gap-3 px-3 py-2 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700 transition" @click="router.push('/data/records')">
          <el-icon :size="18" style="color:#60a5fa"><List /></el-icon>
          <span class="text-gray-800 dark:text-gray-100">操作记录</span>
        </button>
        <button class="w-full flex items-center gap-3 px-3 py-2 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700 transition" @click="showClearDialog = true">
          <el-icon :size="18" style="color:#ef4444"><SwitchButton /></el-icon>
          <span class="text-gray-800 dark:text-gray-100">清除数据</span>
        </button>
      </div>
    </div>

    <!-- 中文注释：设置模块（顺序后移，位于数据管理下方），展示各设置按钮 -->
    <div class="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
      <div class="px-3 py-2 flex items-center gap-2">
        <el-icon :size="18" style="color:#0ea5e9"><Setting /></el-icon>
        <span class="font-semibold dark:text-gray-100">设置</span>
      </div>
      <!-- 响应式网格：移动端单列，桌面端多列 -->
      <div class="px-2 py-2 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-1">
        <button
          v-for="i in settingItems"
          :key="i.key"
          class="w-full flex items-center gap-3 px-3 py-2 rounded-md hover:bg-gray-100 dark:hover:bg-gray-700 transition text-left disabled:opacity-60 disabled:cursor-not-allowed"
          :disabled="isDisabled(i.key)"
          @click="onOpenSetting(i.key)"
        >
          <el-icon :size="18" :style="{ color: i.fg }"><component :is="i.icon" /></el-icon>
          <span class="text-gray-800 dark:text-gray-100">{{ i.label }}</span>
        </button>
      </div>
    </div>

    <el-dialog v-model="showClearDialog" title="清除数据" width="400px">
      <div class="space-y-3">
        <div>请输入登录密码以确认清除该用户的所有数据</div>
        <el-input v-model="clearPassword" type="password" show-password placeholder="登录密码" />
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="showClearDialog = false">取消</el-button>
          <el-button type="danger" @click="confirmClear">确认清除</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 中文注释：按需求取消“最近金币变动”卡片展示 -->
  </div>
  <div style="height: calc(env(safe-area-inset-bottom) + 20px)"></div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAppState } from '@/stores/appState'
import VipBadge from '@/components/VipBadge.vue'
import defaultAvatar from '@/assets/avatars/default.png'
import router from '@/router'
import { useAuth } from '@/stores/auth'
import { usePermissions } from '@/composables/permissions'
import { getStaticBase } from '@/services/http'
import { presignView } from '@/services/storage'
import { ElMessage } from 'element-plus'
import { User, Edit, SwitchButton, Setting, Timer, List, Microphone, Coin, InfoFilled } from '@element-plus/icons-vue'
import http from '@/services/http'

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

// 中文注释：头像展示地址，支持 uploads 相对路径与 S3 对象键
const avatarSrc = ref<string>(defaultAvatar)
async function updateAvatar() {
  const p = auth.user?.avatar_path
  if (!p) { avatarSrc.value = defaultAvatar; return }
  const s = String(p)
  if (/storage\/images\/avatars\/default\.png$/i.test(s) || /(^|\/)default\.png$/i.test(s)) { avatarSrc.value = defaultAvatar; return }
  if (/^https?:\/\//i.test(s)) { avatarSrc.value = s; return }
  const base = getStaticBase()
  if (/uploads\//i.test(s)) { avatarSrc.value = `${base}/api/${s.replace(/^\/+/, '')}`; return }
  try { avatarSrc.value = await presignView(s) } catch { avatarSrc.value = defaultAvatar }
}
onMounted(() => {
  updateAvatar()
  auth.refreshUser()
})
watch(() => auth.user?.avatar_path, () => { updateAvatar() })

// 中文注释：编辑个人信息改为独立页面，移除弹窗相关状态与函数

// 中文注释：权限（模板与交互使用）；父账号或具备 manage_children 权限时启用子账号管理
const { isParent, manageChildren, canSettingTomato, canSettingTasks, canSettingCoins, canSettingReading } = usePermissions()

function onChildManage() {
  if (!isParent.value && !manageChildren.value) {
    ElMessage.warning('当前权限不允许管理子账号')
    return
  }
  router.push('/mine/subaccounts')
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
type SettingsKey = 'tomato' | 'tasks' | 'reading' | 'coins' | 'appearance' | 'about' | 'admin'

const settingItems = computed(() => {
  const items: Array<{ key: SettingsKey; label: string; icon: any; fg: string }> = [
    { key: 'tomato', label: '番茄钟设置', icon: Timer, fg: '#ef4444' },
    { key: 'tasks', label: '任务设置', icon: List, fg: '#10b981' },
    { key: 'reading', label: '朗读设置', icon: Microphone, fg: '#7c3aed' },
    { key: 'coins', label: '金币设置', icon: Coin, fg: '#f59e0b' },
    { key: 'appearance', label: '主题外观', icon: Setting, fg: '#2563eb' },
    { key: 'about', label: '关于', icon: InfoFilled, fg: '#0ea5e9' }
  ]
  // 管理员权限控制
  // if (isAdmin.value) {
  //   items.push({ key: 'admin', label: '用户管理', icon: DataAnalysis, fg: '#db2777' })
  // }
  return items
})


function goSettingsTab(k: SettingsKey) {
  if (k === 'admin') {
    router.push('/admin')
    return
  }
  // 中文注释：跳转到独立的设置页面路径，而不是通用 /settings
  const map: Record<string, string> = {
    tomato: '/settings/tomato',
    tasks: '/settings/tasks',
    reading: '/settings/reading',
    coins: '/settings/coins',
    appearance: '/settings/appearance',
    about: '/settings/about'
  }
  if (map[k]) router.push(map[k])
}

// 中文注释：根据权限判断设置按钮是否禁用
function isDisabled(k: SettingsKey): boolean {
  if (k === 'admin') return false
  if (k === 'about') return false
  if (k === 'appearance') return false
  if (k === 'tomato') return !canSettingTomato.value
  if (k === 'tasks') return !canSettingTasks.value
  if (k === 'reading') return !canSettingReading.value
  if (k === 'coins') return !canSettingCoins.value
  return true
}

// 中文注释：拦截点击：无权限不跳转并提示
function onOpenSetting(k: SettingsKey) {
  if (isDisabled(k)) {
    if (k === 'tasks' || k === 'coins') {
      ElMessage.info('仅查看或无权限，无法打开该设置')
    } else {
      ElMessage.info('无权限，无法打开该设置')
    }
    return
  }
  goSettingsTab(k)
}

const showClearDialog = ref(false)
const clearPassword = ref('')
async function confirmClear() {
  try {
    if (!clearPassword.value) { ElMessage.warning('请输入密码'); return }
    const resp: any = await http.post('/data/clear', { password: clearPassword.value })
    try { store.setCoins(Number(resp?.user_coins ?? 0)) } catch {}
    ElMessage.success('已清除该用户数据')
    showClearDialog.value = false
    clearPassword.value = ''
  } catch (e: any) {
    ElMessage.error(e?.message || '清除失败')
  }
}

// （移除我的页内的设置入口列表，保留“系统设置”按钮跳转到独立页面）

// 中文注释：按需求取消“最近金币变动”卡片相关逻辑
</script>

<style scoped>
</style>
