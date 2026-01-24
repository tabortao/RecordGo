<template>
  <div class="min-h-screen bg-[#F5F7FA] dark:bg-gray-900 flex flex-col">
    <!-- Header -->
    <div class="px-4 py-3 bg-white/80 dark:bg-gray-900/80 backdrop-blur-md sticky top-0 z-30 flex items-center justify-between border-b border-gray-100 dark:border-gray-800 transition-colors">
      <div class="flex items-center gap-3">
        <div 
          class="w-8 h-8 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center cursor-pointer hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors text-gray-600 dark:text-gray-300"
          @click="router.back()"
        >
          <el-icon><ArrowLeft /></el-icon>
        </div>
        <h2 class="font-bold text-gray-800 dark:text-gray-100 text-lg">成绩分析</h2>
      </div>
      <button 
        class="bg-gradient-to-r from-blue-500 to-indigo-600 text-white px-4 py-1.5 rounded-full text-sm font-bold shadow-lg shadow-blue-200 dark:shadow-blue-900/20 hover:scale-105 active:scale-95 transition-all flex items-center gap-1"
        @click="openCreate"
      >
        <el-icon><Plus /></el-icon>
        <span>新增成绩</span>
      </button>
    </div>

    <div class="flex-1 overflow-y-auto p-4 sm:p-6 pb-24">
      <div class="max-w-6xl mx-auto space-y-6">
        <!-- Overview Card -->
        <div class="bg-white dark:bg-gray-800 rounded-3xl p-6 shadow-sm border border-gray-100 dark:border-gray-700 relative overflow-hidden">
          <div class="absolute top-0 right-0 w-64 h-64 bg-gradient-to-bl from-indigo-50 to-transparent dark:from-indigo-900/10 rounded-bl-[100px] -z-0"></div>
          
          <div class="relative z-10">
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-4">学情概览</h1>
            
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
              <div class="rounded-2xl p-5 bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 shadow-sm relative overflow-hidden group">
                <div class="absolute -right-4 -top-4 w-24 h-24 bg-emerald-50 dark:bg-emerald-900/10 rounded-full blur-xl group-hover:scale-150 transition-transform"></div>
                <div class="text-sm font-medium text-gray-500 dark:text-gray-400 flex items-center gap-1 relative z-10">
                  <el-icon class="text-emerald-500"><DataLine /></el-icon> 平均得分率
                </div>
                <div class="mt-2 text-3xl font-black tracking-tight text-gray-900 dark:text-white relative z-10">{{ avgRateText }}</div>
              </div>
              
              <div class="rounded-2xl p-5 text-white bg-gradient-to-br from-orange-400 to-rose-500 shadow-lg shadow-orange-200 dark:shadow-none relative overflow-hidden group">
                <div class="absolute -right-4 -top-4 w-24 h-24 bg-white/10 rounded-full blur-xl group-hover:scale-150 transition-transform"></div>
                <div class="text-sm font-medium opacity-90 flex items-center gap-1">
                  <el-icon><Top /></el-icon> 最高分记录
                </div>
                <div class="mt-2 text-3xl font-black tracking-tight">{{ maxScoreText }}</div>
              </div>
              
              <div class="rounded-2xl p-5 text-white bg-gradient-to-br from-indigo-400 to-purple-600 shadow-lg shadow-indigo-200 dark:shadow-none relative overflow-hidden group">
                <div class="absolute -right-4 -top-4 w-24 h-24 bg-white/10 rounded-full blur-xl group-hover:scale-150 transition-transform"></div>
                <div class="text-sm font-medium opacity-90 flex items-center gap-1">
                  <el-icon><Bottom /></el-icon> 最低分记录
                </div>
                <div class="mt-2 text-3xl font-black tracking-tight">{{ minScoreText }}</div>
              </div>
            </div>

            <!-- Filters -->
            <div class="flex flex-wrap items-center gap-3 bg-gray-50 dark:bg-gray-700/50 p-2 rounded-2xl">
              <el-select v-model="subjectFilter" size="default" class="!w-32 custom-select" placeholder="全部学科">
                <el-option label="全部学科" value="all" />
                <el-option v-for="s in subjects" :key="s" :label="s" :value="s" />
              </el-select>
              <div class="h-6 w-[1px] bg-gray-200 dark:bg-gray-600 mx-1"></div>
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                unlink-panels
                value-format="YYYY-MM-DD"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                class="!w-64 custom-date-picker"
                size="default"
                :prefix-icon="Calendar"
              />
              <button 
                v-if="subjectFilter !== 'all' || dateRange"
                class="ml-auto text-xs font-bold text-gray-500 hover:text-indigo-600 dark:text-gray-400 dark:hover:text-indigo-300 px-3 py-1.5 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
                @click="resetFilters"
              >
                重置筛选
              </button>
            </div>
          </div>
        </div>

        <!-- Charts Grid -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div class="bg-white dark:bg-gray-800 rounded-3xl p-6 shadow-sm border border-gray-100 dark:border-gray-700">
            <div class="flex items-center gap-2 mb-6">
              <div class="w-10 h-10 rounded-xl bg-indigo-50 dark:bg-indigo-900/20 flex items-center justify-center text-indigo-600 dark:text-indigo-400">
                <el-icon :size="20"><TrendCharts /></el-icon>
              </div>
              <h3 class="text-lg font-bold text-gray-800 dark:text-gray-100">各科成绩走势</h3>
            </div>
            <div class="h-64 w-full">
              <VChart v-if="multiSubjectOption" :option="multiSubjectOption" autoresize class="w-full h-full" />
              <div v-else class="w-full h-full flex flex-col items-center justify-center text-gray-400 bg-gray-50 dark:bg-gray-800/50 rounded-2xl border border-dashed border-gray-200 dark:border-gray-700">
                <el-icon :size="32" class="mb-2 opacity-50"><DataLine /></el-icon>
                <span class="text-xs">暂无趋势数据</span>
              </div>
            </div>
          </div>

          <div class="bg-white dark:bg-gray-800 rounded-3xl p-6 shadow-sm border border-gray-100 dark:border-gray-700">
            <div class="flex items-center gap-2 mb-6">
              <div class="w-10 h-10 rounded-xl bg-purple-50 dark:bg-purple-900/20 flex items-center justify-center text-purple-600 dark:text-purple-400">
                <el-icon :size="20"><DataAnalysis /></el-icon>
              </div>
              <h3 class="text-lg font-bold text-gray-800 dark:text-gray-100">总分变化趋势</h3>
            </div>
            <div class="h-64 w-full">
              <VChart v-if="totalScoreOption" :option="totalScoreOption" autoresize class="w-full h-full" />
              <div v-else class="w-full h-full flex flex-col items-center justify-center text-gray-400 bg-gray-50 dark:bg-gray-800/50 rounded-2xl border border-dashed border-gray-200 dark:border-gray-700">
                <el-icon :size="32" class="mb-2 opacity-50"><DataLine /></el-icon>
                <span class="text-xs">暂无趋势数据</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Subject Trends List -->
        <div v-if="subjectTrendOptions.length > 0" class="space-y-4">
          <h3 class="text-sm font-bold text-gray-500 dark:text-gray-400 ml-1 uppercase tracking-wider">单科详细分析</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div
              v-for="item in subjectTrendOptions"
              :key="item.subject"
              class="bg-white dark:bg-gray-800 rounded-3xl p-5 shadow-sm border border-gray-100 dark:border-gray-700 hover:shadow-md transition-shadow"
            >
              <div class="flex items-center gap-2 mb-4">
                <span class="px-2.5 py-1 rounded-lg text-xs font-bold" :class="subjectBadgeClass(item.subject)">{{ item.subject }}</span>
                <span class="text-xs text-gray-400">得分率趋势</span>
              </div>
              <div class="h-40 w-full">
                <VChart :option="item.option" autoresize class="w-full h-full" />
              </div>
            </div>
          </div>
        </div>

        <!-- Records List -->
        <div class="space-y-4">
          <h3 class="text-sm font-bold text-gray-500 dark:text-gray-400 ml-1 uppercase tracking-wider flex items-center justify-between">
            <span>历史成绩记录</span>
            <span class="text-xs font-normal normal-case bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded-full">共 {{ filteredRecords.length }} 条</span>
          </h3>
          
          <div v-if="filteredRecords.length === 0" class="flex flex-col items-center justify-center py-16 text-gray-400 bg-white dark:bg-gray-800 rounded-3xl border border-dashed border-gray-200 dark:border-gray-700">
            <div class="w-16 h-16 bg-gray-50 dark:bg-gray-700/50 rounded-full flex items-center justify-center mb-3">
              <el-icon :size="24" class="text-gray-300"><Document /></el-icon>
            </div>
            <p class="text-sm">暂无符合条件的成绩记录</p>
          </div>
          
          <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div
              v-for="item in filteredRecords"
              :key="item.id"
              class="group bg-white dark:bg-gray-800 rounded-2xl p-4 border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-lg hover:-translate-y-0.5 transition-all duration-300 cursor-pointer relative overflow-hidden"
              @click="openDetail(item.id)"
            >
              <div class="absolute top-0 right-0 w-20 h-20 bg-gradient-to-bl from-gray-50 to-transparent dark:from-gray-700/20 rounded-bl-full -z-0 transition-transform group-hover:scale-110"></div>
              
              <div class="relative z-10">
                <div class="flex items-start justify-between mb-3">
                  <div class="flex-1 min-w-0 mr-4">
                    <h4 class="text-base font-bold text-gray-900 dark:text-white truncate mb-1 group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">{{ item.exam_name }}</h4>
                    <div class="flex items-center gap-2 text-xs text-gray-500 dark:text-gray-400">
                      <el-icon><Calendar /></el-icon>
                      <span>{{ formatDate(item.exam_date) }}</span>
                    </div>
                  </div>
                  <div class="flex flex-col items-end">
                     <span class="text-2xl font-black text-gray-900 dark:text-white leading-none">
                       {{ item.score }}<span class="text-xs font-medium text-gray-400 ml-0.5">/{{ item.full_score }}</span>
                     </span>
                     <span class="px-2 py-0.5 rounded-md text-[10px] font-bold mt-1" :class="subjectBadgeClass(item.subject)">
                       {{ item.subject }}
                     </span>
                  </div>
                </div>
                
                <div class="grid grid-cols-3 gap-2 pt-3 border-t border-gray-50 dark:border-gray-700/50">
                  <div class="text-center">
                    <div class="text-[10px] text-gray-400 mb-0.5">班级排名</div>
                    <div class="text-sm font-bold text-gray-700 dark:text-gray-200">{{ item.class_rank || '-' }}</div>
                  </div>
                  <div class="text-center border-l border-gray-50 dark:border-gray-700/50">
                    <div class="text-[10px] text-gray-400 mb-0.5">年级排名</div>
                    <div class="text-sm font-bold text-gray-700 dark:text-gray-200">{{ item.grade_rank || '-' }}</div>
                  </div>
                  <div class="text-center border-l border-gray-50 dark:border-gray-700/50">
                    <div class="text-[10px] text-gray-400 mb-0.5">班级均分</div>
                    <div class="text-sm font-bold text-gray-700 dark:text-gray-200">{{ item.class_avg || '-' }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ArrowLeft, DataAnalysis, TrendCharts, Plus, DataLine, Top, Bottom, Calendar, Document } from '@element-plus/icons-vue'
import router from '@/router'
import dayjs from 'dayjs'
import { useAuth } from '@/stores/auth'
import { listScores, type ScoreRecord } from '@/services/scores'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([LineChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const auth = useAuth()
const records = ref<ScoreRecord[]>([])

const subjectFilter = ref('all')
const dateRange = ref<[string, string] | null>(null)

const displayName = computed(() => {
  const u = auth.user
  if (!u) return '未登录'
  const n = (u.nickname || '').trim()
  return n ? n : u.username
})

const ageText = computed(() => {
  const b = auth.user?.child_birthday || ''
  if (!b) return '未设置'
  const d = dayjs(b)
  if (!d.isValid()) return '未设置'
  const years = dayjs().diff(d, 'year')
  return years >= 0 ? `${years}岁` : '未设置'
})

const genderText = computed(() => {
  const g = auth.user?.child_gender || ''
  if (!g) return '未设置'
  if (g === 'male') return '男'
  if (g === 'female') return '女'
  return '未设置'
})

const subjects = computed(() => Array.from(new Set(records.value.map((r) => r.subject))).filter(Boolean))

const filteredRecords = computed(() => {
  let list = [...records.value]
  if (subjectFilter.value !== 'all') {
    list = list.filter((r) => r.subject === subjectFilter.value)
  }
  if (dateRange.value) {
    const [start, end] = dateRange.value
    list = list.filter((r) => r.exam_date >= start && r.exam_date <= end)
  }
  return list.sort((a, b) => (a.exam_date < b.exam_date ? 1 : -1))
})

const chartRecords = computed(() => {
  let list = [...records.value]
  if (subjectFilter.value !== 'all') {
    list = list.filter((r) => r.subject === subjectFilter.value)
  }
  if (dateRange.value) {
    const [start, end] = dateRange.value
    list = list.filter((r) => r.exam_date >= start && r.exam_date <= end)
  }
  return list.sort((a, b) => (a.exam_date > b.exam_date ? 1 : -1))
})

const avgRateText = computed(() => {
  if (filteredRecords.value.length === 0) return '--'
  const total = filteredRecords.value.reduce((sum, r) => sum + (r.full_score > 0 ? r.score / r.full_score : 0), 0)
  return `${Math.round((total / filteredRecords.value.length) * 100)}%`
})

const maxScoreText = computed(() => {
  if (filteredRecords.value.length === 0) return '--'
  return Math.max(...filteredRecords.value.map((r) => r.score)).toString()
})

const minScoreText = computed(() => {
  if (filteredRecords.value.length === 0) return '--'
  return Math.min(...filteredRecords.value.map((r) => r.score)).toString()
})

const multiSubjectOption = computed(() => {
  const list = chartRecords.value
  if (list.length === 0) return null
  const subjectsList = Array.from(new Set(list.map((r) => r.subject))).filter(Boolean)
  const dates = Array.from(new Set(list.map((r) => r.exam_date))).sort()
  const series = subjectsList.map((sub) => {
    const map = new Map(list.filter((r) => r.subject === sub).map((r) => [r.exam_date, r.score]))
    return {
      name: sub,
      type: 'line',
      smooth: true,
      data: dates.map((d) => map.get(d) ?? null)
    }
  })
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: subjectsList },
    grid: { left: 24, right: 16, top: 32, bottom: 24 },
    xAxis: { type: 'category', data: dates },
    yAxis: { type: 'value' },
    series
  }
})

const totalScoreOption = computed(() => {
  const list = chartRecords.value
  if (list.length === 0) return null
  const groups = new Map<string, { label: string; score: number; full: number; date: string }>()
  list.forEach((r) => {
    const key = `${r.exam_name}|${r.exam_date}`
    const label = `${r.exam_name} ${r.exam_date}`
    const item = groups.get(key) || { label, score: 0, full: 0, date: r.exam_date }
    item.score += r.score
    item.full += r.full_score
    groups.set(key, item)
  })
  const items = Array.from(groups.values()).sort((a, b) => (a.date > b.date ? 1 : -1))
  return {
    tooltip: { trigger: 'axis' },
    grid: { left: 24, right: 16, top: 32, bottom: 48 },
    xAxis: { type: 'category', data: items.map((i) => i.label) },
    yAxis: { type: 'value', max: 100 },
    series: [
      {
        name: '总分得分率',
        type: 'line',
        smooth: true,
        data: items.map((i) => (i.full > 0 ? Math.round((i.score / i.full) * 100) : 0))
      }
    ]
  }
})

const subjectTrendOptions = computed(() => {
  const list = chartRecords.value
  const subjectsList = Array.from(new Set(list.map((r) => r.subject))).filter(Boolean)
  return subjectsList.map((sub) => {
    const items = list.filter((r) => r.subject === sub).sort((a, b) => (a.exam_date > b.exam_date ? 1 : -1))
    const dates = items.map((r) => r.exam_date)
    const data = items.map((r) => (r.full_score > 0 ? Math.round((r.score / r.full_score) * 100) : 0))
    return {
      subject: sub,
      option: {
        tooltip: { trigger: 'axis' },
        grid: { left: 24, right: 16, top: 24, bottom: 24 },
        xAxis: { type: 'category', data: dates },
        yAxis: { type: 'value', max: 100 },
        series: [
          { type: 'line', smooth: true, data }
        ]
      }
    }
  })
})

function formatDate(d: string) {
  return dayjs(d).format('YYYY-MM-DD')
}

function subjectCardClass(subject: string) {
  const map: Record<string, string> = {
    语文: 'bg-amber-50/70 dark:bg-amber-500/10',
    数学: 'bg-indigo-50/70 dark:bg-indigo-500/10',
    英语: 'bg-emerald-50/70 dark:bg-emerald-500/10',
    物理: 'bg-purple-50/70 dark:bg-purple-500/10',
    化学: 'bg-rose-50/70 dark:bg-rose-500/10',
    生物: 'bg-fuchsia-50/70 dark:bg-fuchsia-500/10',
    历史: 'bg-yellow-50/70 dark:bg-yellow-500/10',
    地理: 'bg-lime-50/70 dark:bg-lime-500/10',
    政治: 'bg-teal-50/70 dark:bg-teal-500/10',
    信息技术: 'bg-blue-50/70 dark:bg-blue-500/10',
    科学: 'bg-cyan-50/70 dark:bg-cyan-500/10'
  }
  return map[subject] || 'bg-gray-50/70 dark:bg-gray-700/30'
}

function subjectBadgeClass(subject: string) {
  const map: Record<string, string> = {
    语文: 'bg-amber-100 text-amber-700 dark:bg-amber-500/30 dark:text-amber-200',
    数学: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-500/30 dark:text-indigo-200',
    英语: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/30 dark:text-emerald-200',
    科学: 'bg-cyan-100 text-cyan-700 dark:bg-cyan-500/30 dark:text-cyan-200',
    物理: 'bg-purple-100 text-purple-700 dark:bg-purple-500/30 dark:text-purple-200',
    化学: 'bg-rose-100 text-rose-700 dark:bg-rose-500/30 dark:text-rose-200',
    生物: 'bg-fuchsia-100 text-fuchsia-700 dark:bg-fuchsia-500/30 dark:text-fuchsia-200',
    政治: 'bg-teal-100 text-teal-700 dark:bg-teal-500/30 dark:text-teal-200',
    历史: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-500/30 dark:text-yellow-200',
    地理: 'bg-lime-100 text-lime-700 dark:bg-lime-500/30 dark:text-lime-200',
    信息技术: 'bg-blue-100 text-blue-700 dark:bg-blue-500/30 dark:text-blue-200'
  }
  return map[subject] || 'bg-gray-100 text-gray-700 dark:bg-gray-600/40 dark:text-gray-200'
}

async function load() {
  const resp = await listScores()
  records.value = resp.items || []
}

function resetFilters() {
  subjectFilter.value = 'all'
  dateRange.value = null
}

function openDetail(id: number) {
  router.push(`/grades/${id}`)
}

function openCreate() {
  router.push('/grades/create')
}

onMounted(load)
</script>

<style scoped>
.custom-select :deep(.el-input__wrapper) {
  background-color: transparent;
  box-shadow: none !important;
  padding: 0 8px;
}
.custom-select :deep(.el-input__inner) {
  font-weight: 600;
  color: #4b5563;
}
.dark .custom-select :deep(.el-input__inner) {
  color: #e5e7eb;
}

.custom-date-picker :deep(.el-input__wrapper) {
  background-color: transparent;
  box-shadow: none !important;
  padding: 0 8px;
}
.custom-date-picker :deep(.el-range-input) {
  background-color: transparent;
  font-size: 13px;
  color: #4b5563;
}
.dark .custom-date-picker :deep(.el-range-input) {
  color: #e5e7eb;
}
.custom-date-picker :deep(.el-range-separator) {
  color: #9ca3af;
}
</style>
