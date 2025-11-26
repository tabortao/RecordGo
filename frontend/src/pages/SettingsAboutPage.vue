<template>
  <!-- 中文注释：关于页面，展示图标、名称、版本号，以及外部跳转按钮（小红书、邮件联系） -->
  <div class="p-4 space-y-4 pb-16">
    <!-- 标题栏：返回 + 关于图标 + 文本 -->
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer" style="color:#64748b" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" style="color:#0ea5e9"><InfoFilled /></el-icon>
      <h2 class="font-semibold">关于</h2>
    </div>
    
    <!-- 信息卡片：图标、名称、版本分别占一行显示 -->
    <el-card shadow="never" body-class="p-4">
      <div class="flex flex-col items-center justify-center p-4">
        <img :src="iconSrc" alt="App Icon" class="w-12 h-12 mb-2 rounded-lg shadow-md" />
        <h3 class="text-2xl font-bold text-gray-800 dark:text-white mb-1">{{ appName }}</h3>
        <p class="text-sm text-gray-600 dark:text-gray-300">版本 {{ versionDisplay }}</p>
        <p class="text-sm text-gray-600 dark:text-gray-300">{{ dateDisplay }}</p>
      </div>
    </el-card>

    <!-- 操作卡片：关注小红书、邮件联系（移动端类列表样式，左侧图标 + 右侧箭头；响应式排列） -->
    <el-card shadow="never" body-class="p-4">
      <!-- 中文注释：移动端单列、平板/桌面两列并排 -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <!-- 关注小红书 -->
        <button class="w-full flex items-center justify-between px-3 py-3 rounded-lg border dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700 transition"
                :disabled="!xhsUrl"
                @click="openXHS">
          <div class="flex items-center gap-3">
            <el-icon :size="18" style="color:#ef4444"><Link /></el-icon>
            <span class="text-gray-800 dark:text-gray-200">关注小红书</span>
          </div>
          <el-icon :size="18" class="text-gray-500 dark:text-gray-400"><ArrowRight /></el-icon>
        </button>

        <!-- 邮件联系 -->
        <button class="w-full flex items-center justify-between px-3 py-3 rounded-lg border dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-700 transition"
                @click="openMail">
          <div class="flex items-center gap-3">
            <el-icon :size="18" style="color:#0ea5e9"><Message /></el-icon>
            <span class="text-gray-800 dark:text-gray-200">邮件联系</span>
          </div>
          <el-icon :size="18" class="text-gray-500 dark:text-gray-400"><ArrowRight /></el-icon>
        </button>
      </div>
    </el-card>
  </div>
  
</template>

<script setup lang="ts">
// 中文注释：导入路由与图标；导入应用信息 JSON（名称、版本、外链）
import { InfoFilled, ArrowLeft, Link, Message, ArrowRight } from '@element-plus/icons-vue'
import router from '@/router'
import appInfo from '@/config/app-info.json'

// 中文注释：软件图标路径（Vite 处理静态资源导入）
import icon192 from '@/assets/favicon/android-chrome-192x192.png'

function goBack() { router.back() }

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
  } catch (e) {
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
