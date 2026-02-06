<template>
  <SettingsShell title="使用帮助" :icon="QuestionFilled" tone="sky" container-class="max-w-4xl">
    <!-- Header Hero -->
    <SettingsCard>
      <div class="flex flex-col sm:flex-row items-center gap-6 text-center sm:text-left">
        <div class="relative shrink-0">
          <div class="absolute -inset-4 rounded-full bg-sky-200/50 dark:bg-sky-500/15 blur-2xl animate-pulse" />
          <div class="relative h-20 w-20 rounded-[24px] border border-white/60 dark:border-gray-800/60 bg-white/70 dark:bg-gray-950/30 ring-1 ring-black/5 dark:ring-white/10 grid place-items-center overflow-hidden shadow-lg shadow-sky-100 dark:shadow-none">
            <img :src="appIconSrc" alt="App Icon" class="h-14 w-14 rounded-xl" />
          </div>
        </div>
        <div class="min-w-0 flex-1 space-y-3">
          <div>
            <div class="text-2xl font-black tracking-tight text-gray-900 dark:text-gray-50">作业家</div>
            <div class="text-sm font-bold text-sky-600 dark:text-sky-400">任务积分 + 学习激励助手</div>
          </div>
          <div class="text-sm leading-relaxed text-gray-600 dark:text-gray-300 max-w-xl">
            通过任务管理、番茄专注、积分奖励与心愿兑换，把“完成一件事”变成看得见的进步；配合生长记录，持续追踪孩子的每一点变化。
          </div>
          <div class="flex flex-wrap justify-center sm:justify-start gap-2">
            <span v-for="tag in ['任务管理', '积分激励', '番茄专注', '心愿兑换', '成长记录']" :key="tag" class="rounded-full border border-sky-100 dark:border-sky-900/30 bg-sky-50 dark:bg-sky-900/10 px-3 py-1 text-xs font-bold text-sky-700 dark:text-sky-300">
              {{ tag }}
            </span>
          </div>
        </div>
      </div>
    </SettingsCard>

    <!-- Quick Actions -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <button
        v-for="(action, idx) in quickActions"
        :key="idx"
        class="group relative overflow-hidden rounded-[24px] border border-white/60 dark:border-gray-800/60 bg-white/40 dark:bg-gray-900/40 p-5 text-left transition-all duration-300 hover:bg-white dark:hover:bg-gray-900 hover:shadow-lg hover:shadow-gray-200/50 dark:hover:shadow-black/20 hover:-translate-y-1"
        @click="router.push(action.path)"
      >
        <div class="absolute right-0 top-0 -mt-4 -mr-4 h-24 w-24 rounded-full transition-transform group-hover:scale-110 opacity-50" :class="action.bgClass"></div>
        <div class="relative flex flex-col h-full justify-between gap-4">
          <div class="h-12 w-12 rounded-2xl flex items-center justify-center text-white shadow-md" :class="action.iconBgClass">
            <el-icon :size="24"><component :is="action.icon" /></el-icon>
          </div>
          <div>
            <div class="text-lg font-black text-gray-900 dark:text-gray-50">{{ action.title }}</div>
            <div class="text-xs font-medium text-gray-500 dark:text-gray-400 mt-1">{{ action.desc }}</div>
          </div>
        </div>
      </button>
    </div>

    <!-- Steps -->
    <SettingsCard title="快速上手" description="简单三步，开启正向循环">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 relative">
        <!-- Connecting Line (Desktop) -->
        <div class="hidden md:block absolute top-8 left-[16%] right-[16%] h-0.5 bg-gradient-to-r from-transparent via-gray-200 dark:via-gray-700 to-transparent z-0"></div>
        
        <div v-for="(step, idx) in steps" :key="idx" class="relative z-10 flex flex-col items-center text-center group">
          <div class="h-16 w-16 rounded-[24px] bg-white dark:bg-gray-800 border-4 border-gray-50 dark:border-gray-900 shadow-sm flex items-center justify-center mb-4 transition-transform group-hover:scale-110 duration-300" :class="step.colorClass">
            <el-icon :size="28"><component :is="step.icon" /></el-icon>
          </div>
          <div class="text-xs font-black tracking-widest text-gray-400 dark:text-gray-500 mb-1">STEP 0{{ idx + 1 }}</div>
          <div class="text-base font-black text-gray-900 dark:text-gray-50 mb-2">{{ step.title }}</div>
          <div class="text-xs text-gray-500 dark:text-gray-400 leading-relaxed px-4">{{ step.desc }}</div>
        </div>
      </div>
    </SettingsCard>

    <!-- Tips -->
    <SettingsCard title="使用建议" description="更轻松地坚持下去">
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div v-for="(tip, idx) in tips" :key="idx" class="rounded-[24px] border border-gray-100 dark:border-gray-800 bg-gradient-to-br from-white to-gray-50 dark:from-gray-900 dark:to-gray-800/50 p-5 transition-all hover:shadow-md">
          <div class="flex items-start gap-4">
            <div class="h-10 w-10 rounded-xl flex items-center justify-center shrink-0" :class="tip.bgClass">
              <el-icon :size="20" :class="tip.textClass"><component :is="tip.icon" /></el-icon>
            </div>
            <div>
              <div class="text-sm font-black text-gray-900 dark:text-gray-50 mb-1">{{ tip.title }}</div>
              <div class="text-xs leading-relaxed text-gray-500 dark:text-gray-400">{{ tip.desc }}</div>
            </div>
          </div>
        </div>
      </div>
    </SettingsCard>

    <!-- Features Grid -->
    <SettingsCard title="核心功能一览">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div v-for="(feat, idx) in features" :key="idx" class="group rounded-[24px] border border-gray-100 dark:border-gray-800 bg-white/60 dark:bg-gray-900/40 p-4 hover:bg-white dark:hover:bg-gray-900 transition-colors">
          <div class="flex items-center gap-4">
            <div class="h-12 w-12 rounded-2xl flex items-center justify-center transition-transform group-hover:scale-110 duration-300" :class="feat.bgClass">
              <el-icon :size="24" :class="feat.textClass"><component :is="feat.icon" /></el-icon>
            </div>
            <div>
              <div class="text-base font-black text-gray-900 dark:text-gray-50">{{ feat.title }}</div>
              <div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ feat.desc }}</div>
            </div>
          </div>
        </div>
      </div>
    </SettingsCard>

    <!-- Downloads -->
    <SettingsCard title="下载地址" description="获取最新版本安装包">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <a v-for="(dl, idx) in downloads" :key="idx" :href="dl.url" target="_blank" class="block group">
          <div class="rounded-[24px] border border-gray-100 dark:border-gray-800 bg-white/60 dark:bg-gray-900/40 p-5 transition-all duration-300 hover:bg-white dark:hover:bg-gray-900 hover:shadow-lg hover:shadow-blue-500/10 hover:-translate-y-0.5 relative overflow-hidden">
            <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/40 to-transparent -translate-x-full group-hover:animate-shimmer pointer-events-none"></div>
            <div class="flex items-center justify-between relative z-10">
              <div class="flex items-center gap-4">
                <div class="h-12 w-12 rounded-2xl bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 grid place-items-center shadow-sm">
                  <el-icon :size="24"><Download /></el-icon>
                </div>
                <div>
                  <div class="text-base font-black text-gray-900 dark:text-gray-50">{{ dl.title }}</div>
                  <div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ dl.desc }} <span v-if="dl.code" class="font-bold bg-gray-100 dark:bg-gray-800 px-1.5 py-0.5 rounded text-gray-800 dark:text-gray-200">{{ dl.code }}</span></div>
                </div>
              </div>
              <div class="h-8 w-8 rounded-full bg-gray-50 dark:bg-gray-800 flex items-center justify-center text-gray-400 group-hover:bg-blue-500 group-hover:text-white transition-all">
                <el-icon><ArrowRight /></el-icon>
              </div>
            </div>
          </div>
        </a>
      </div>
    </SettingsCard>

    <!-- FAQ -->
    <SettingsCard title="常见问题" description="遇到问题先看这里">
      <div class="space-y-3">
        <details v-for="(faq, idx) in faqs" :key="idx" class="group rounded-[20px] border border-gray-100 dark:border-gray-800 bg-white/60 dark:bg-gray-900/40 overflow-hidden transition-all duration-300 open:bg-white open:shadow-md dark:open:bg-gray-900">
          <summary class="cursor-pointer list-none flex items-center justify-between gap-4 p-4 select-none">
            <div class="text-sm font-black text-gray-900 dark:text-gray-50 flex items-center gap-2">
              <div class="h-6 w-6 rounded-full bg-sky-100 dark:bg-sky-900/30 text-sky-600 dark:text-sky-400 text-xs font-bold flex items-center justify-center">Q</div>
              {{ faq.q }}
            </div>
            <div class="h-6 w-6 rounded-full bg-gray-50 dark:bg-gray-800 flex items-center justify-center transition-transform duration-300 group-open:rotate-180">
              <el-icon :size="12" class="text-gray-400"><ArrowDown /></el-icon>
            </div>
          </summary>
          <div class="px-4 pb-4 pl-[3.25rem]">
            <div class="text-sm leading-relaxed text-gray-600 dark:text-gray-300 bg-gray-50 dark:bg-gray-800/50 p-3 rounded-xl">
              {{ faq.a }}
            </div>
          </div>
        </details>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { QuestionFilled, Download, ArrowRight, ArrowDown, Edit, Clock, Present, List, Coin, TrendCharts, DataAnalysis, CirclePlusFilled, Microphone, CircleCheck, InfoFilled } from '@element-plus/icons-vue'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'
import appIcon from '@/assets/favicon/android-chrome-192x192.png'

const router = useRouter()
const appIconSrc = (appIcon as any) as string

const quickActions = [
  { 
    title: '创建任务', desc: '从一个小目标开始', icon: CirclePlusFilled, path: '/tasks/create',
    bgClass: 'bg-emerald-100 dark:bg-emerald-900/30', iconBgClass: 'bg-emerald-500'
  },
  { 
    title: '查看统计', desc: '复盘更有动力', icon: DataAnalysis, path: '/tasks/stats',
    bgClass: 'bg-indigo-100 dark:bg-indigo-900/30', iconBgClass: 'bg-indigo-500'
  },
  { 
    title: '朗读设置', desc: '辅助专注与复述', icon: Microphone, path: '/settings/reading',
    bgClass: 'bg-sky-100 dark:bg-sky-900/30', iconBgClass: 'bg-sky-500'
  }
]

const steps = [
  { title: '创建任务', desc: '把作业拆成可完成的小任务，并设置奖励金币。', icon: Edit, colorClass: 'text-sky-500' },
  { title: '番茄专注', desc: '用 20-25 分钟的专注块，高效完成目标。', icon: Clock, colorClass: 'text-amber-500' },
  { title: '兑换心愿', desc: '用赚取的金币兑换奖励，形成正向循环。', icon: Present, colorClass: 'text-emerald-500' }
]

const tips = [
  { 
    title: '从小任务开始', desc: '先把“能完成”放在第一位，建立信心。', icon: CircleCheck,
    bgClass: 'bg-emerald-100 dark:bg-emerald-900/30', textClass: 'text-emerald-600 dark:text-emerald-400'
  },
  { 
    title: '用番茄做节奏', desc: '20-25 分钟专注 + 5 分钟休息，效率更稳定。', icon: Clock,
    bgClass: 'bg-amber-100 dark:bg-amber-900/30', textClass: 'text-amber-600 dark:text-amber-400'
  },
  { 
    title: '奖励要具体', desc: '心愿越明确、越渴望，越有动力去赚金币。', icon: Present,
    bgClass: 'bg-indigo-100 dark:bg-indigo-900/30', textClass: 'text-indigo-600 dark:text-indigo-400'
  }
]

const features = [
  { title: '任务管理', desc: '支持重复任务、完成打卡、回收站', icon: List, bgClass: 'bg-sky-100 dark:bg-sky-900/20', textClass: 'text-sky-600 dark:text-sky-400' },
  { title: '积分激励', desc: '做任务赚金币，兑换心愿奖励', icon: Coin, bgClass: 'bg-amber-100 dark:bg-amber-900/20', textClass: 'text-amber-600 dark:text-amber-400' },
  { title: '成长记录', desc: '记录身高、体重、视力变化趋势', icon: TrendCharts, bgClass: 'bg-emerald-100 dark:bg-emerald-900/20', textClass: 'text-emerald-600 dark:text-emerald-400' },
  { title: '数据统计', desc: '多维度图表，直观展示学习成果', icon: DataAnalysis, bgClass: 'bg-indigo-100 dark:bg-indigo-900/20', textClass: 'text-indigo-600 dark:text-indigo-400' }
]

const downloads = [
  { title: '百度网盘下载', desc: '提取码：', code: 'izyj', url: 'https://pan.baidu.com/s/16zL_7-ZuYBxcZ3aVXXqWCw?pwd=izyj' },
  { title: '123 盘下载', desc: '高速下载推荐', code: '', url: 'https://www.123865.com/s/FQLuVv-c7ztv' }
]

const faqs = [
  { q: '金币是怎么计算的？', a: '完成任务会增加预设的金币数量；兑换心愿会扣除相应金币。系统会自动记录每一笔收支明细，可在“我的-积分记录”中查看。' },
  { q: '为什么看不到底部导航？', a: '为了提供更好的沉浸式体验，部分功能页面（如设置、成长记录、专注计时）会自动隐藏底部导航栏，点击左上角返回即可。' },
  { q: '数据会丢失吗？', a: '目前数据存储在本地数据库中。建议定期备份数据文件，后续版本将支持云同步功能。' }
]
</script>

<style scoped>
@keyframes shimmer {
  100% { transform: translateX(100%); }
}
.animate-shimmer {
  animation: shimmer 1.5s infinite;
}
</style>
