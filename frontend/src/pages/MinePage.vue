<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 pb-24 transition-colors duration-300" style="padding-bottom: calc(env(safe-area-inset-bottom) + 96px)">
    
    <!-- 顶部 Hero 区域 -->
    <div class="relative bg-gradient-to-br from-blue-600 to-cyan-500 dark:from-blue-900 dark:to-cyan-900 pt-5 pb-7 px-6 rounded-b-[2rem] shadow-xl overflow-hidden mb-6">
       <!-- 装饰背景 -->
       <div class="absolute top-0 left-0 w-full h-full overflow-hidden opacity-10 pointer-events-none">
          <div class="absolute -top-10 -right-10 w-40 h-40 bg-white rounded-full blur-3xl"></div>
          <div class="absolute top-20 -left-10 w-32 h-32 bg-yellow-300 rounded-full blur-3xl"></div>
       </div>

       <!-- 用户信息卡片 -->
       <div class="relative z-10 flex flex-row items-center mt-2 gap-5 animate-fade-in-up">
          <div class="relative group cursor-pointer" @click="isParent && router.push('/settings/profile')">
             <div class="absolute -inset-1 bg-white/30 rounded-full blur group-hover:blur-md transition duration-500"></div>
             <el-avatar :size="72" :src="avatarSrc" class="relative border-[3px] border-white/90 dark:border-gray-800 shadow-xl transition-transform duration-300 group-hover:scale-105" />
          </div>
          
          <div class="text-left flex-1 min-w-0">
             <div class="flex items-center gap-2 min-w-0">
               <h2 class="text-2xl font-bold text-white tracking-tight drop-shadow-md truncate">{{ displayName }}</h2>
               <div class="shrink-0">
                 <VipBadge :user="auth.user" />
               </div>
             </div>
          </div>
       </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="px-4 space-y-4 -mt-6 relative z-20">
       
       <!-- 账号管理组 -->
       <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700/50 overflow-hidden animate-slide-up" style="animation-delay: 0.1s;">
          <div class="px-4 py-3 border-b border-gray-50 dark:border-gray-700/30 flex items-center gap-2">
             <div class="w-1 h-4 bg-blue-500 rounded-full shadow-sm"></div>
             <span class="font-bold text-gray-800 dark:text-gray-100 text-sm">账号管理</span>
          </div>
          <div class="p-1">
             <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-1">
                <!-- 编辑个人信息 -->
                <button v-if="isParent" 
                  class="group w-full flex items-center justify-between px-3 py-3 rounded-xl hover:bg-blue-50/50 dark:hover:bg-gray-700/50 active:scale-[0.99] transition-all duration-200" 
                  @click="router.push('/settings/profile')">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-full bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center text-blue-600 dark:text-blue-400">
                      <el-icon :size="16"><Edit /></el-icon>
                    </div>
                    <span class="text-gray-700 dark:text-gray-200 font-medium text-sm">编辑个人信息</span>
                  </div>
                  <el-icon class="text-gray-300 dark:text-gray-600 group-hover:text-blue-400 transition-colors"><ArrowRight /></el-icon>
                </button>

                <!-- 子账号管理 -->
                <button v-if="isParent || manageChildren" 
                  class="group w-full flex items-center justify-between px-3 py-3 rounded-xl hover:bg-green-50/50 dark:hover:bg-gray-700/50 active:scale-[0.99] transition-all duration-200" 
                  @click="onChildManage">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-full bg-green-100 dark:bg-green-900/30 flex items-center justify-center text-green-600 dark:text-green-400">
                      <el-icon :size="16"><User /></el-icon>
                    </div>
                    <span class="text-gray-700 dark:text-gray-200 font-medium text-sm">子账号管理</span>
                  </div>
                  <el-icon class="text-gray-300 dark:text-gray-600 group-hover:text-green-400 transition-colors"><ArrowRight /></el-icon>
                </button>

                <!-- 退出登录 -->
                <button 
                  class="group w-full flex items-center justify-between px-3 py-3 rounded-xl hover:bg-red-50/50 dark:hover:bg-gray-700/50 active:scale-[0.99] transition-all duration-200" 
                  @click="onLogout">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-full bg-red-100 dark:bg-red-900/30 flex items-center justify-center text-red-600 dark:text-red-400">
                      <el-icon :size="16"><SwitchButton /></el-icon>
                    </div>
                    <span class="text-gray-700 dark:text-gray-200 font-medium text-sm">退出登录</span>
                  </div>
                  <el-icon class="text-gray-300 dark:text-gray-600 group-hover:text-red-400 transition-colors"><ArrowRight /></el-icon>
                </button>
             </div>
          </div>
       </div>

       <!-- 数据管理 -->
       <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700/50 overflow-hidden animate-slide-up" style="animation-delay: 0.2s;">
          <div class="px-4 py-3 border-b border-gray-50 dark:border-gray-700/30 flex items-center gap-2">
             <div class="w-1 h-4 bg-orange-500 rounded-full shadow-sm"></div>
             <span class="font-bold text-gray-800 dark:text-gray-100 text-sm">数据管理</span>
          </div>
         <div class="p-1">
           <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-3 gap-1">
             <button
               class="group w-full flex items-center justify-between px-3 py-3 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700/50 active:scale-[0.99] transition-all duration-200"
               @click="router.push('/data/records')"
             >
               <div class="flex items-center gap-3">
                 <div class="w-8 h-8 rounded-full bg-orange-100 dark:bg-orange-900/30 flex items-center justify-center text-orange-600 dark:text-orange-400">
                   <el-icon :size="16"><List /></el-icon>
                 </div>
                 <span class="text-gray-700 dark:text-gray-200 font-medium text-sm">操作记录</span>
               </div>
               <el-icon class="text-gray-300 dark:text-gray-600 group-hover:text-gray-500 transition-colors"><ArrowRight /></el-icon>
             </button>

             <button
               class="group w-full flex items-center justify-between px-3 py-3 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700/50 active:scale-[0.99] transition-all duration-200"
               @click="showClearDialog = true"
             >
               <div class="flex items-center gap-3">
                 <div class="w-8 h-8 rounded-full bg-red-100 dark:bg-red-900/30 flex items-center justify-center text-red-600 dark:text-red-400">
                   <el-icon :size="16"><SwitchButton /></el-icon>
                 </div>
                 <span class="text-gray-700 dark:text-gray-200 font-medium text-sm">清除数据</span>
               </div>
               <el-icon class="text-gray-300 dark:text-gray-600 group-hover:text-gray-500 transition-colors"><ArrowRight /></el-icon>
             </button>

             <el-tooltip content="开发中" placement="top">
               <button
                 class="group w-full flex items-center justify-between px-3 py-3 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-all duration-200 opacity-60 cursor-not-allowed"
                 type="button"
                 disabled
               >
                 <div class="flex items-center gap-3">
                   <div class="w-8 h-8 rounded-full bg-gray-100 dark:bg-gray-700/60 flex items-center justify-center text-gray-400 dark:text-gray-500">
                     <el-icon :size="16"><DocumentIcon /></el-icon>
                   </div>
                   <span class="text-gray-600 dark:text-gray-300 font-medium text-sm">数据导出</span>
                 </div>
                 <el-icon class="text-gray-300 dark:text-gray-700"><ArrowRight /></el-icon>
               </button>
             </el-tooltip>
           </div>
         </div>
       </div>

       <!-- 系统设置 -->
       <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700/50 overflow-hidden animate-slide-up" style="animation-delay: 0.3s;">
          <div class="px-4 py-3 border-b border-gray-50 dark:border-gray-700/30 flex items-center gap-2">
             <div class="w-1 h-4 bg-purple-500 rounded-full shadow-sm"></div>
             <span class="font-bold text-gray-800 dark:text-gray-100 text-sm">系统设置</span>
          </div>
          <div class="p-1">
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-1">
              <button
                v-for="i in settingItems"
                :key="i.key"
                class="group w-full flex items-center justify-between px-3 py-3 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700/50 active:scale-[0.99] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
                :disabled="isDisabled(i.key)"
                @click="onOpenSetting(i.key)"
              >
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-full flex items-center justify-center transition-colors" 
                       :style="{ backgroundColor: isDisabled(i.key) ? '#f3f4f6' : `${i.fg}15`, color: isDisabled(i.key) ? '#9ca3af' : i.fg }">
                    <el-icon :size="16"><component :is="i.icon" /></el-icon>
                  </div>
                  <span class="text-gray-700 dark:text-gray-200 font-medium text-sm">{{ i.label }}</span>
                </div>
                <el-icon class="text-gray-300 dark:text-gray-600 group-hover:text-gray-500 transition-colors"><ArrowRight /></el-icon>
              </button>
            </div>
          </div>
       </div>

    </div>

    <!-- 清除数据弹窗 -->
    <el-dialog v-model="showClearDialog" title="清除数据" width="90%" class="rounded-2xl overflow-hidden max-w-sm">
      <div class="space-y-4 py-2">
        <div class="text-sm text-gray-600 dark:text-gray-400 bg-red-50 dark:bg-red-900/20 p-3 rounded-lg border border-red-100 dark:border-red-900/30 flex gap-2">
           <el-icon class="text-red-500 mt-0.5"><InfoFilled /></el-icon>
           <span>此操作不可逆，将清除该用户的所有数据。请输入登录密码以确认。</span>
        </div>
        <el-input v-model="clearPassword" type="password" show-password placeholder="请输入登录密码" class="w-full" size="large" />
      </div>
      <template #footer>
        <div class="flex justify-end gap-2 pt-2">
          <el-button @click="showClearDialog = false" class="!rounded-lg">取消</el-button>
          <el-button type="danger" @click="confirmClear" class="!rounded-lg shadow-red-200">确认清除</el-button>
        </div>
      </template>
    </el-dialog>

  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, watch } from 'vue'
import { useAppState } from '@/stores/appState'
import VipBadge from '@/components/VipBadge.vue'
import defaultAvatar from '@/assets/avatars/default.png'
import router from '@/router'
import { useAuth } from '@/stores/auth'
import { usePermissions } from '@/composables/permissions'
import { getStaticBase } from '@/services/http'
import { presignView } from '@/services/storage'
import { ElMessage } from 'element-plus'
import { User, Edit, SwitchButton, Setting, Timer, List, Microphone, Coin, InfoFilled, Cpu, Document as DocumentIcon, ArrowRight } from '@element-plus/icons-vue'
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
type SettingsKey = 'tomato' | 'tasks' | 'ai' | 'ocr' | 'reading' | 'coins' | 'appearance' | 'about' | 'admin'

const settingItems = computed(() => {
  const items: Array<{ key: SettingsKey; label: string; icon: any; fg: string }> = [
    { key: 'tomato', label: '番茄钟设置', icon: Timer, fg: '#ef4444' },
    { key: 'tasks', label: '任务设置', icon: List, fg: '#10b981' },
    { key: 'ai', label: 'AI模型设置', icon: Cpu, fg: '#6366f1' },
    { key: 'ocr', label: 'OCR服务', icon: DocumentIcon, fg: '#3b82f6' },
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
    ai: '/settings/ai',
    ocr: '/settings/ocr',
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
  if (k === 'ai') return false
  if (k === 'ocr') return false
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
@keyframes fade-in-up {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slide-up {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in-up {
  animation: fade-in-up 0.6s ease-out forwards;
}

.animate-slide-up {
  opacity: 0; /* Initially hidden */
  animation: slide-up 0.5s ease-out forwards;
}
</style>
