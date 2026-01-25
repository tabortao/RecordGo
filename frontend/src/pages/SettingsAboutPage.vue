<template>
  <SettingsShell title="关于" subtitle="版本信息与联系入口" :icon="InfoFilled" tone="sky" container-class="max-w-2xl">
    <SettingsCard>
      <div class="flex flex-col items-center justify-center px-2 py-6">
        <div class="relative">
          <div class="absolute -inset-3 rounded-[1.5rem] bg-sky-200/50 dark:bg-sky-500/15 blur-2xl" />
          <img :src="iconSrc" alt="App Icon" class="relative w-14 h-14 rounded-2xl shadow-md ring-1 ring-white/60 dark:ring-gray-800/70" />
        </div>
        <h3 class="mt-4 text-2xl font-extrabold tracking-tight text-gray-900 dark:text-white">{{ appName }}</h3>
        <div class="mt-2 flex flex-wrap items-center justify-center gap-2">
          <span class="rounded-full border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 px-3 py-1 text-xs font-semibold text-gray-700 dark:text-gray-300">
            版本 {{ versionDisplay }}
          </span>
          <span class="rounded-full border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 px-3 py-1 text-xs font-semibold text-gray-700 dark:text-gray-300">
            {{ dateDisplay }}
          </span>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard title="快捷入口" description="打开外部链接或进入帮助页面">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <button
          class="group w-full flex items-center justify-between px-3 py-3 rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/50 dark:bg-gray-900/40 hover:bg-white/80 dark:hover:bg-gray-900/70 transition disabled:opacity-50 disabled:cursor-not-allowed"
          :disabled="!xhsUrl"
          @click="openXHS"
        >
          <div class="flex items-center gap-3">
            <div class="h-10 w-10 rounded-2xl border border-rose-100/70 dark:border-rose-900/40 bg-rose-50/80 dark:bg-rose-900/25 flex items-center justify-center text-rose-600 dark:text-rose-300">
              <el-icon :size="18"><Link /></el-icon>
            </div>
            <div class="text-left">
              <div class="text-sm font-bold text-gray-900 dark:text-gray-50">关注小红书</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">动态与使用技巧</div>
            </div>
          </div>
          <el-icon :size="18" class="text-gray-300 dark:text-gray-600 group-hover:text-gray-500 transition-colors"><ArrowRight /></el-icon>
        </button>

        <button
          class="group w-full flex items-center justify-between px-3 py-3 rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/50 dark:bg-gray-900/40 hover:bg-white/80 dark:hover:bg-gray-900/70 transition"
          @click="openMail"
        >
          <div class="flex items-center gap-3">
            <div class="h-10 w-10 rounded-2xl border border-sky-100/70 dark:border-sky-900/40 bg-sky-50/80 dark:bg-sky-900/25 flex items-center justify-center text-sky-700 dark:text-sky-300">
              <el-icon :size="18"><Message /></el-icon>
            </div>
            <div class="text-left">
              <div class="text-sm font-bold text-gray-900 dark:text-gray-50">邮件联系</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">反馈建议与问题</div>
            </div>
          </div>
          <el-icon :size="18" class="text-gray-300 dark:text-gray-600 group-hover:text-gray-500 transition-colors"><ArrowRight /></el-icon>
        </button>

        <button
          class="group w-full flex items-center justify-between px-3 py-3 rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/50 dark:bg-gray-900/40 hover:bg-white/80 dark:hover:bg-gray-900/70 transition"
          @click="router.push('/settings/help')"
        >
          <div class="flex items-center gap-3">
            <div class="h-10 w-10 rounded-2xl border border-indigo-100/70 dark:border-indigo-900/40 bg-indigo-50/80 dark:bg-indigo-900/25 flex items-center justify-center text-indigo-700 dark:text-indigo-300">
              <el-icon :size="18"><QuestionFilled /></el-icon>
            </div>
            <div class="text-left">
              <div class="text-sm font-bold text-gray-900 dark:text-gray-50">使用帮助</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">常见问题与说明</div>
            </div>
          </div>
          <el-icon :size="18" class="text-gray-300 dark:text-gray-600 group-hover:text-gray-500 transition-colors"><ArrowRight /></el-icon>
        </button>

        <button
          class="group w-full flex items-center justify-between px-3 py-3 rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/50 dark:bg-gray-900/40 hover:bg-white/80 dark:hover:bg-gray-900/70 transition"
          @click="router.push('/settings/support')"
        >
          <div class="flex items-center gap-3">
            <div class="h-10 w-10 rounded-2xl border border-amber-100/70 dark:border-amber-900/40 bg-amber-50/80 dark:bg-amber-900/25 flex items-center justify-center text-amber-700 dark:text-amber-300">
              <el-icon :size="18"><Present /></el-icon>
            </div>
            <div class="text-left">
              <div class="text-sm font-bold text-gray-900 dark:text-gray-50">打赏支持</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">支持作者持续更新</div>
            </div>
          </div>
          <el-icon :size="18" class="text-gray-300 dark:text-gray-600 group-hover:text-gray-500 transition-colors"><ArrowRight /></el-icon>
        </button>
      </div>
    </SettingsCard>
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
