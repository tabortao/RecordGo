<template>
  <SettingsShell title="关于" :icon="InfoFilled" tone="sky" container-class="max-w-2xl">
    <SettingsCard>
      <div class="flex flex-col items-center justify-center px-4 py-8">
        <div class="relative group">
          <div class="absolute -inset-6 rounded-full bg-sky-200/40 dark:bg-sky-500/10 blur-3xl transition-all duration-500 group-hover:bg-sky-300/50 dark:group-hover:bg-sky-500/20" />
          <img :src="iconSrc" alt="App Icon" class="relative w-24 h-24 rounded-[28px] shadow-xl ring-1 ring-white/60 dark:ring-gray-800/70 transition-transform duration-500 group-hover:scale-105 group-hover:rotate-3" />
        </div>
        <h3 class="mt-6 text-3xl font-black tracking-tight text-gray-900 dark:text-white">{{ appName }}</h3>
        <p class="mt-2 text-sm font-medium text-sky-600 dark:text-sky-400">记录成长，遇见更好的自己</p>
        
        <div class="mt-8 grid grid-cols-2 gap-4 w-full max-w-sm">
          <div class="rounded-[20px] border border-gray-100 dark:border-gray-800 bg-gray-50/50 dark:bg-gray-900/30 p-4 text-center transition-colors hover:bg-white dark:hover:bg-gray-800">
            <div class="text-[10px] font-black uppercase tracking-widest text-gray-400 dark:text-gray-500 mb-1.5">VERSION</div>
            <div class="text-base font-black text-gray-900 dark:text-gray-50 tabular-nums">{{ versionDisplay }}</div>
          </div>
          <div class="rounded-[20px] border border-gray-100 dark:border-gray-800 bg-gray-50/50 dark:bg-gray-900/30 p-4 text-center transition-colors hover:bg-white dark:hover:bg-gray-800">
            <div class="text-[10px] font-black uppercase tracking-widest text-gray-400 dark:text-gray-500 mb-1.5">BUILD DATE</div>
            <div class="text-base font-black text-gray-900 dark:text-gray-50 tabular-nums">{{ dateDisplay }}</div>
          </div>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard title="快捷入口" description="联系我们或获取帮助">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <button
          class="group w-full flex items-center justify-between p-4 rounded-[24px] border border-gray-100 dark:border-gray-800 bg-white/60 dark:bg-gray-900/40 hover:bg-white dark:hover:bg-gray-900 hover:shadow-lg hover:shadow-rose-100/50 dark:hover:shadow-none transition-all duration-300 active:scale-[0.99] disabled:opacity-50 disabled:cursor-not-allowed"
          :disabled="!xhsUrl"
          @click="openXHS"
        >
          <div class="flex items-center gap-4">
            <div class="h-12 w-12 rounded-2xl border border-rose-100/70 dark:border-rose-900/40 bg-rose-50 dark:bg-rose-900/20 flex items-center justify-center text-rose-500 dark:text-rose-400 shadow-sm group-hover:scale-110 transition-transform duration-300">
              <el-icon :size="22"><Link /></el-icon>
            </div>
            <div class="text-left">
              <div class="text-base font-black text-gray-900 dark:text-gray-50">关注小红书</div>
              <div class="text-xs font-medium text-gray-500 dark:text-gray-400 mt-0.5">获取最新动态与技巧</div>
            </div>
          </div>
          <div class="h-8 w-8 rounded-full flex items-center justify-center text-gray-300 dark:text-gray-600 group-hover:bg-rose-50 dark:group-hover:bg-rose-900/30 group-hover:text-rose-500 transition-all">
            <el-icon :size="16"><ArrowRight /></el-icon>
          </div>
        </button>

        <button
          class="group w-full flex items-center justify-between p-4 rounded-[24px] border border-gray-100 dark:border-gray-800 bg-white/60 dark:bg-gray-900/40 hover:bg-white dark:hover:bg-gray-900 hover:shadow-lg hover:shadow-sky-100/50 dark:hover:shadow-none transition-all duration-300 active:scale-[0.99]"
          @click="openMail"
        >
          <div class="flex items-center gap-4">
            <div class="h-12 w-12 rounded-2xl border border-sky-100/70 dark:border-sky-900/40 bg-sky-50 dark:bg-sky-900/20 flex items-center justify-center text-sky-500 dark:text-sky-400 shadow-sm group-hover:scale-110 transition-transform duration-300">
              <el-icon :size="22"><Message /></el-icon>
            </div>
            <div class="text-left">
              <div class="text-base font-black text-gray-900 dark:text-gray-50">邮件联系</div>
              <div class="text-xs font-medium text-gray-500 dark:text-gray-400 mt-0.5">反馈建议与问题</div>
            </div>
          </div>
          <div class="h-8 w-8 rounded-full flex items-center justify-center text-gray-300 dark:text-gray-600 group-hover:bg-sky-50 dark:group-hover:bg-sky-900/30 group-hover:text-sky-500 transition-all">
            <el-icon :size="16"><ArrowRight /></el-icon>
          </div>
        </button>

        <button
          class="group w-full h-full flex items-center justify-between p-4 rounded-[24px] border border-gray-100 dark:border-gray-800 bg-white/60 dark:bg-gray-900/40 hover:bg-white dark:hover:bg-gray-900 hover:shadow-lg hover:shadow-indigo-100/50 dark:hover:shadow-none transition-all duration-300 active:scale-[0.99]"
          @click="router.push('/settings/help')"
        >
          <div class="flex items-center gap-4">
            <div class="h-12 w-12 rounded-2xl border border-indigo-100/70 dark:border-indigo-900/40 bg-indigo-50 dark:bg-indigo-900/20 flex items-center justify-center text-indigo-500 dark:text-indigo-400 shadow-sm group-hover:scale-110 transition-transform duration-300">
              <el-icon :size="22"><QuestionFilled /></el-icon>
            </div>
            <div class="text-left">
              <div class="text-base font-black text-gray-900 dark:text-gray-50">使用帮助</div>
              <div class="text-xs font-medium text-gray-500 dark:text-gray-400 mt-0.5">常见问题与说明</div>
            </div>
          </div>
          <div class="h-8 w-8 rounded-full flex items-center justify-center text-gray-300 dark:text-gray-600 group-hover:bg-indigo-50 dark:group-hover:bg-indigo-900/30 group-hover:text-indigo-500 transition-all">
            <el-icon :size="16"><ArrowRight /></el-icon>
          </div>
        </button>

        <button
          class="group w-full h-full flex items-center justify-between p-4 rounded-[24px] border border-gray-100 dark:border-gray-800 bg-white/60 dark:bg-gray-900/40 hover:bg-white dark:hover:bg-gray-900 hover:shadow-lg hover:shadow-amber-100/50 dark:hover:shadow-none transition-all duration-300 active:scale-[0.99]"
          @click="router.push('/settings/support')"
        >
          <div class="flex items-center gap-4">
            <div class="h-12 w-12 rounded-2xl border border-amber-100/70 dark:border-amber-900/40 bg-amber-50 dark:bg-amber-900/20 flex items-center justify-center text-amber-500 dark:text-amber-400 shadow-sm group-hover:scale-110 transition-transform duration-300">
              <el-icon :size="22"><Present /></el-icon>
            </div>
            <div class="text-left">
              <div class="text-base font-black text-gray-900 dark:text-gray-50">打赏支持</div>
              <div class="text-xs font-medium text-gray-500 dark:text-gray-400 mt-0.5">支持作者持续更新</div>
            </div>
          </div>
          <div class="h-8 w-8 rounded-full flex items-center justify-center text-gray-300 dark:text-gray-600 group-hover:bg-amber-50 dark:group-hover:bg-amber-900/30 group-hover:text-amber-500 transition-all">
            <el-icon :size="16"><ArrowRight /></el-icon>
          </div>
        </button>
      </div>
    </SettingsCard>
    
    <div class="mt-8 text-center">
      <p class="text-[10px] font-bold text-gray-300 dark:text-gray-600 uppercase tracking-widest">Designed with ❤️ for kids</p>
    </div>
  </SettingsShell>
</template>

<script setup lang="ts">
// 中文注释：导入路由与图标；导入应用信息 JSON（名称、版本、外链）
import { InfoFilled, Link, Message, ArrowRight, QuestionFilled, Present } from '@element-plus/icons-vue'
import router from '@/router'
import appInfo from '@/config/app-info.json'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

// 中文注释：软件图标路径（Vite 处理静态资源导入）
import icon192 from '@/assets/favicon/android-chrome-192x192.png'

// 中文注释：从 JSON 配置读取软件名称与版本；为防止字段缺失，增加回退值
const appName = (appInfo as any)?.name || '任务积分助手'
const versionDisplay = (appInfo as any)?.version || '未配置版本'
const dateDisplay = (appInfo as any)?.date || '未配置日期'
const xhsUrl: string | null = ((appInfo as any)?.xhs_url || '').trim() || null
const contactEmail: string = ((appInfo as any)?.contact_email || 'terlun@foxmail.com').trim()

const iconSrc = (icon192 as any) as string

// 中文注释：打开小红书主页（新标签）；若未配置链接则提示
function openXHS() {
  if (!xhsUrl) return
  try {
    window.open(xhsUrl, '_blank')
  } catch {
    location.href = xhsUrl
  }
}

// 中文注释：触发邮件客户端，自动填充收件人
function openMail() {
  const mailto = `mailto:${contactEmail}`
  try {
    window.location.href = mailto
  } catch {
    window.open(mailto)
  }
}
</script>

<style scoped>
/* 中文注释：页面样式采用 Element Plus + Tailwind 响应式类，保证不同屏幕良好显示 */
</style>
