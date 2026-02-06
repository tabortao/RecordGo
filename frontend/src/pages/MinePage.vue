<template>
  <div class="relative min-h-screen overflow-x-hidden bg-gray-50 text-gray-900 dark:bg-gray-950 dark:text-gray-50 transition-colors duration-300" style="padding-bottom: calc(env(safe-area-inset-bottom) + 96px)">
    
    <!-- Background Decoration -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <div class="absolute -top-24 -right-24 h-72 w-72 rounded-full bg-sky-300/60 dark:bg-sky-500/25 blur-3xl opacity-35 dark:opacity-25" />
      <div class="absolute -bottom-40 -left-28 h-80 w-80 rounded-full bg-indigo-200/60 dark:bg-indigo-500/20 blur-3xl opacity-30 dark:opacity-20" />
      <div class="absolute inset-0 bg-[radial-gradient(1200px_circle_at_20%_-20%,rgba(255,255,255,.65),transparent_55%),radial-gradient(900px_circle_at_80%_0%,rgba(255,255,255,.45),transparent_55%)] dark:bg-[radial-gradient(1200px_circle_at_20%_-20%,rgba(255,255,255,.07),transparent_55%),radial-gradient(900px_circle_at_80%_0%,rgba(255,255,255,.06),transparent_55%)]" />
    </div>

    <!-- Header -->
    <div class="sticky top-0 z-30 px-4 pt-4">
      <div class="rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/75 dark:bg-gray-900/70 backdrop-blur-xl shadow-sm">
        <div class="flex items-center justify-between gap-3 px-4 py-4">
          <div class="flex items-center gap-3 min-w-0" @click="isParent && router.push('/settings/profile')">
            <div class="relative shrink-0 cursor-pointer group">
              <div class="absolute -inset-1 rounded-full bg-gradient-to-tr from-sky-400 to-indigo-400 opacity-0 group-hover:opacity-50 blur transition-opacity duration-500" />
              <div class="relative w-12 h-12 rounded-2xl ring-2 ring-white dark:ring-gray-800 overflow-hidden bg-white/60 dark:bg-gray-950/25 grid place-items-center">
                <el-avatar :size="48" :src="avatarSrc" class="w-full h-full" />
              </div>
            </div>
            <div class="min-w-0">
              <div class="text-lg font-black tracking-tight text-gray-900 dark:text-gray-50 truncate">{{ displayName }}</div>
            </div>
          </div>
          <div class="shrink-0">
            <VipBadge :user="auth.user" />
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="relative z-10 px-4 pt-5 space-y-5 max-w-5xl mx-auto">
      
      <!-- Stats Summary -->
      <div class="grid grid-cols-2 gap-3">
        <div class="rounded-[24px] bg-gradient-to-br from-amber-50 to-orange-50 dark:from-amber-900/20 dark:to-orange-900/10 border border-amber-100/60 dark:border-amber-900/30 p-4 relative overflow-hidden group">
          <div class="absolute top-0 right-0 -mt-2 -mr-2 w-16 h-16 bg-amber-200/20 dark:bg-amber-500/10 rounded-full blur-xl group-hover:scale-150 transition-transform duration-500"></div>
          <div class="relative z-10">
            <div class="flex items-center gap-2 mb-1 text-amber-600 dark:text-amber-400">
              <el-icon><Coin /></el-icon>
              <span class="text-xs font-bold uppercase tracking-wider">我的金币</span>
            </div>
            <div class="text-2xl font-black text-gray-900 dark:text-gray-50 tabular-nums">{{ store.coins }}</div>
          </div>
        </div>
        
        <div class="rounded-[24px] bg-gradient-to-br from-sky-50 to-blue-50 dark:from-sky-900/20 dark:to-blue-900/10 border border-sky-100/60 dark:border-sky-900/30 p-4 relative overflow-hidden group">
          <div class="absolute top-0 right-0 -mt-2 -mr-2 w-16 h-16 bg-sky-200/20 dark:bg-sky-500/10 rounded-full blur-xl group-hover:scale-150 transition-transform duration-500"></div>
          <div class="relative z-10">
            <div class="flex items-center gap-2 mb-1 text-sky-600 dark:text-sky-400">
              <el-icon><List /></el-icon>
              <span class="text-xs font-bold uppercase tracking-wider">今日任务</span>
            </div>
            <div class="text-2xl font-black text-gray-900 dark:text-gray-50 tabular-nums">{{ todayTasksDisplay }}</div>
          </div>
        </div>
      </div>

      <!-- Account Management -->
      <SettingsCard title="账号管理">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
          <button v-if="isParent" 
            class="group w-full flex items-center justify-between p-3 rounded-2xl hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors active:scale-[0.99]" 
            @click="router.push('/settings/profile')">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-xl bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 grid place-items-center">
                <el-icon :size="20"><Edit /></el-icon>
              </div>
              <div class="text-left">
                <div class="text-sm font-bold text-gray-900 dark:text-gray-100">编辑个人信息</div>
                <div class="text-xs text-gray-500 dark:text-gray-400">修改头像、昵称</div>
              </div>
            </div>
            <el-icon class="text-gray-300 dark:text-gray-600 group-hover:text-blue-500 transition-colors"><ArrowRight /></el-icon>
          </button>

          <button v-if="isParent || manageChildren" 
            class="group w-full flex items-center justify-between p-3 rounded-2xl hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors active:scale-[0.99]" 
            @click="onChildManage">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-xl bg-emerald-50 dark:bg-emerald-900/20 text-emerald-600 dark:text-emerald-400 grid place-items-center">
                <el-icon :size="20"><User /></el-icon>
              </div>
              <div class="text-left">
                <div class="text-sm font-bold text-gray-900 dark:text-gray-100">子账号管理</div>
                <div class="text-xs text-gray-500 dark:text-gray-400">添加或管理家庭成员</div>
              </div>
            </div>
            <el-icon class="text-gray-300 dark:text-gray-600 group-hover:text-emerald-500 transition-colors"><ArrowRight /></el-icon>
          </button>

          <button 
            class="group w-full flex items-center justify-between p-3 rounded-2xl hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors active:scale-[0.99]" 
            @click="onLogout">
            <div class="flex items-center gap-4">
              <div class="w-10 h-10 rounded-xl bg-rose-50 dark:bg-rose-900/20 text-rose-600 dark:text-rose-400 grid place-items-center">
                <el-icon :size="20"><SwitchButton /></el-icon>
              </div>
              <div class="text-left">
                <div class="text-sm font-bold text-gray-900 dark:text-gray-100">退出登录</div>
                <div class="text-xs text-gray-500 dark:text-gray-400">切换账号或注销</div>
              </div>
            </div>
            <el-icon class="text-gray-300 dark:text-gray-600 group-hover:text-rose-500 transition-colors"><ArrowRight /></el-icon>
          </button>
        </div>
      </SettingsCard>

      <!-- Data Management -->
      <SettingsCard title="数据管理">
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
          <button
            class="group flex items-center gap-3 p-3 rounded-2xl bg-gray-50/50 dark:bg-gray-900/20 hover:bg-white dark:hover:bg-gray-800 border border-transparent hover:border-gray-100 dark:hover:border-gray-700 shadow-sm hover:shadow-md transition-all duration-300"
            @click="router.push('/data/records')"
          >
            <div class="w-10 h-10 rounded-xl bg-orange-100 dark:bg-orange-900/30 text-orange-600 dark:text-orange-400 grid place-items-center shrink-0">
              <el-icon :size="20"><List /></el-icon>
            </div>
            <div class="text-left min-w-0">
              <div class="text-sm font-bold text-gray-900 dark:text-gray-100 truncate">操作记录</div>
              <div class="text-xs text-gray-500 dark:text-gray-400 truncate">查看收支明细</div>
            </div>
          </button>

          <button
            class="group flex items-center gap-3 p-3 rounded-2xl bg-gray-50/50 dark:bg-gray-900/20 hover:bg-white dark:hover:bg-gray-800 border border-transparent hover:border-gray-100 dark:hover:border-gray-700 shadow-sm hover:shadow-md transition-all duration-300"
            @click="showClearDialog = true"
          >
            <div class="w-10 h-10 rounded-xl bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 grid place-items-center shrink-0">
              <el-icon :size="20"><Delete /></el-icon>
            </div>
            <div class="text-left min-w-0">
              <div class="text-sm font-bold text-gray-900 dark:text-gray-100 truncate">清除数据</div>
              <div class="text-xs text-gray-500 dark:text-gray-400 truncate">重置所有记录</div>
            </div>
          </button>
        </div>
      </SettingsCard>

      <!-- System Settings -->
      <SettingsCard title="系统设置">
        <div class="grid grid-cols-3 sm:grid-cols-4 gap-3">
          <button
            v-for="item in settingItems"
            :key="item.key"
            class="group flex flex-col items-center gap-2 p-3 rounded-2xl hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-all active:scale-[0.95] disabled:opacity-40 disabled:cursor-not-allowed"
            :disabled="isDisabled(item.key)"
            @click="onOpenSetting(item.key)"
          >
            <div 
              class="w-12 h-12 rounded-2xl grid place-items-center transition-transform duration-300 group-hover:scale-110 shadow-sm"
              :style="{ backgroundColor: isDisabled(item.key) ? '#f3f4f6' : `${item.color}15`, color: isDisabled(item.key) ? '#9ca3af' : item.color }"
            >
              <el-icon :size="24"><component :is="item.icon" /></el-icon>
            </div>
            <span class="text-xs font-medium text-gray-600 dark:text-gray-300 text-center leading-tight">{{ item.label }}</span>
          </button>
        </div>
      </SettingsCard>

    </div>

    <!-- Clear Data Dialog -->
    <el-dialog v-model="showClearDialog" title="清除数据" width="90%" class="rounded-3xl overflow-hidden max-w-sm">
      <div class="space-y-4 py-2">
        <div class="text-sm text-gray-600 dark:text-gray-400 bg-red-50 dark:bg-red-900/20 p-4 rounded-2xl border border-red-100 dark:border-red-900/30 flex gap-3 items-start">
           <el-icon class="text-red-500 mt-0.5 shrink-0" :size="16"><InfoFilled /></el-icon>
           <span class="leading-relaxed">此操作不可逆，将清除该用户的所有数据。请输入登录密码以确认。</span>
        </div>
        <el-input v-model="clearPassword" type="password" show-password placeholder="请输入登录密码" class="w-full" size="large" />
      </div>
      <template #footer>
        <div class="flex justify-end gap-3 pt-2">
          <el-button @click="showClearDialog = false" class="!rounded-xl !h-10 !px-6">取消</el-button>
          <el-button type="danger" @click="confirmClear" class="!rounded-xl !h-10 !px-6 shadow-lg shadow-red-200 dark:shadow-none">确认清除</el-button>
        </div>
      </template>
    </el-dialog>

  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  User, Edit, SwitchButton, Setting, Timer, List, Microphone, Coin, 
  InfoFilled, Cpu, Document as DocumentIcon, ArrowRight, Delete, Monitor 
} from '@element-plus/icons-vue'
import { useAuth } from '@/stores/auth'
import { useAppState } from '@/stores/appState'
import { usePermissions } from '@/composables/permissions'
import { useTodayStats } from '@/composables/useTodayStats'
import http, { getStaticBase } from '@/services/http'
import { presignView } from '@/services/storage'
import VipBadge from '@/components/VipBadge.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'
import defaultAvatar from '@/assets/avatars/default.png'

const router = useRouter()
const auth = useAuth()
const store = useAppState()
const { isParent, manageChildren, canSettingTomato, canSettingTasks, canSettingCoins, canSettingReading } = usePermissions()

// Stats
const { total: totalTasks, completed: completedTasks, refresh: refreshStats } = useTodayStats()
const todayTasksDisplay = computed(() => `${completedTasks.value}/${totalTasks.value}`)

// User Info
const displayName = computed(() => {
  const u = auth.user
  if (!u) return '未登录'
  return (u.nickname || '').trim() || u.username
})

const avatarSrc = ref<string>(defaultAvatar)
async function updateAvatar() {
  const p = auth.user?.avatar_path
  if (!p) { avatarSrc.value = defaultAvatar; return }
  const s = String(p)
  if (s.includes('default.png')) { avatarSrc.value = defaultAvatar; return }
  if (/^https?:\/\//i.test(s)) { avatarSrc.value = s; return }
  const base = getStaticBase()
  if (s.includes('uploads/')) { avatarSrc.value = `${base}/api/${s.replace(/^\/+/, '')}`; return }
  try { avatarSrc.value = await presignView(s) } catch { avatarSrc.value = defaultAvatar }
}

onMounted(() => {
  updateAvatar()
  auth.refreshUser()
  refreshStats()
})
watch(() => auth.user?.avatar_path, () => updateAvatar())

// Actions
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

// Settings Items
type SettingsKey = 'tomato' | 'tasks' | 'ai' | 'ocr' | 'reading' | 'coins' | 'appearance' | 'about'
const settingItems = computed(() => [
  { key: 'tomato', label: '番茄钟', icon: Timer, color: '#ef4444' },
  { key: 'tasks', label: '任务配置', icon: List, color: '#10b981' },
  { key: 'ai', label: 'AI模型', icon: Cpu, color: '#6366f1' },
  { key: 'ocr', label: 'OCR识别', icon: DocumentIcon, color: '#3b82f6' },
  { key: 'reading', label: '语音朗读', icon: Microphone, color: '#7c3aed' },
  { key: 'coins', label: '金币规则', icon: Coin, color: '#f59e0b' },
  { key: 'appearance', label: '主题外观', icon: Monitor, color: '#2563eb' },
  { key: 'about', label: '关于软件', icon: InfoFilled, color: '#0ea5e9' }
])

function isDisabled(k: string): boolean {
  if (['about', 'ai', 'ocr', 'appearance'].includes(k)) return false
  if (k === 'tomato') return !canSettingTomato.value
  if (k === 'tasks') return !canSettingTasks.value
  if (k === 'reading') return !canSettingReading.value
  if (k === 'coins') return !canSettingCoins.value
  return true
}

function onOpenSetting(k: string) {
  if (isDisabled(k)) {
    ElMessage.info('无权限，无法打开该设置')
    return
  }
  const map: Record<string, string> = {
    tomato: '/settings/tomato',
    tasks: '/settings/tasks',
    ai: '/settings/ai',
    ocr: '/settings/ocr',
    reading: '/settings/reading',
    coins: '/settings/coins',
    appearance: '/settings/appearance',
    about: '/settings/about'
  }
  if (map[k]) router.push(map[k])
}

// Clear Data
const showClearDialog = ref(false)
const clearPassword = ref('')
async function confirmClear() {
  if (!clearPassword.value) { ElMessage.warning('请输入密码'); return }
  try {
    const resp: any = await http.post('/data/clear', { password: clearPassword.value })
    try { store.setCoins(Number(resp?.user_coins ?? 0)) } catch {}
    ElMessage.success('已清除该用户数据')
    showClearDialog.value = false
    clearPassword.value = ''
  } catch (e: any) {
    ElMessage.error(e?.message || '清除失败')
  }
}
</script>

<style scoped>
/* Custom Scrollbar for hidden overflow if needed, though usually native is fine */
</style>
