<template>
  <div class="min-h-screen bg-[#F5F7FA] dark:bg-gray-900">
    <div class="bg-white/80 dark:bg-gray-900/80 backdrop-blur border-b border-gray-200 dark:border-gray-800 sticky top-0 z-10">
      <div class="px-4 py-3 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <el-icon :size="18" class="cursor-pointer text-gray-600 dark:text-gray-300" @click="router.back()"><ArrowLeft /></el-icon>
          <h2 class="font-semibold text-gray-800 dark:text-gray-100">成绩</h2>
        </div>
        <el-button type="primary" size="small" @click="openCreate">新增成绩</el-button>
      </div>
    </div>

    <div class="p-4 space-y-4">
      <div class="rounded-2xl bg-gradient-to-br from-indigo-50 via-white to-rose-50 dark:from-gray-800 dark:via-gray-800 dark:to-gray-900 border border-indigo-100 dark:border-gray-700 p-4 shadow-sm">
        <div class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ displayName }}</div>
        <div class="mt-2 flex flex-wrap items-center gap-4 text-sm text-gray-600 dark:text-gray-300">
          <div class="flex items-center gap-1">
            <span class="text-gray-500 dark:text-gray-400">年龄</span>
            <span class="text-gray-800 dark:text-gray-100">{{ ageText }}</span>
          </div>
          <div class="flex items-center gap-1">
            <span class="text-gray-500 dark:text-gray-400">性别</span>
            <span class="text-gray-800 dark:text-gray-100">{{ genderText }}</span>
          </div>
        </div>
      </div>

      <div class="rounded-2xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4 shadow-sm">
        <div class="flex flex-wrap items-center gap-3">
          <el-select v-model="subjectFilter" size="small" class="min-w-[120px]" placeholder="全部学科">
            <el-option label="全部学科" value="all" />
            <el-option v-for="s in subjects" :key="s" :label="s" :value="s" />
          </el-select>
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            unlink-panels
            value-format="YYYY-MM-DD"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            class="min-w-[240px]"
            size="small"
          />
          <el-button size="small" @click="resetFilters">重置</el-button>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-3 md:grid-cols-3">
        <div class="rounded-2xl p-4 text-white bg-gradient-to-br from-emerald-400 to-cyan-500 shadow-sm">
          <div class="text-sm opacity-90">平均得分率</div>
          <div class="mt-2 text-2xl font-semibold">{{ avgRateText }}</div>
        </div>
        <div class="rounded-2xl p-4 text-white bg-gradient-to-br from-orange-400 to-pink-500 shadow-sm">
          <div class="text-sm opacity-90">最高分</div>
          <div class="mt-2 text-2xl font-semibold">{{ maxScoreText }}</div>
        </div>
        <div class="rounded-2xl p-4 text-white bg-gradient-to-br from-purple-400 to-indigo-600 shadow-sm">
          <div class="text-sm opacity-90">最低分</div>
          <div class="mt-2 text-2xl font-semibold">{{ minScoreText }}</div>
        </div>
      </div>

      <div class="rounded-2xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4 shadow-sm">
        <div class="flex items-center gap-2 text-gray-800 dark:text-gray-100 font-semibold">
          <el-icon><TrendCharts /></el-icon>
          <span>各科成绩变化趋势</span>
        </div>
        <div class="mt-3 h-64">
          <VChart v-if="multiSubjectOption" :option="multiSubjectOption" autoresize />
          <div v-else class="text-sm text-gray-400 dark:text-gray-500">暂无数据</div>
        </div>
      </div>

      <div class="rounded-2xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4 shadow-sm">
        <div class="flex items-center gap-2 text-gray-800 dark:text-gray-100 font-semibold">
          <el-icon><DataAnalysis /></el-icon>
          <span>总分变化趋势</span>
        </div>
        <div class="mt-3 h-64">
          <VChart v-if="totalScoreOption" :option="totalScoreOption" autoresize />
          <div v-else class="text-sm text-gray-400 dark:text-gray-500">暂无数据</div>
        </div>
      </div>

      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div
          v-for="item in subjectTrendOptions"
          :key="item.subject"
          class="rounded-2xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4 shadow-sm"
        >
          <div class="text-gray-800 dark:text-gray-100 font-semibold">{{ item.subject }} 得分率趋势</div>
          <div class="mt-3 h-56">
            <VChart :option="item.option" autoresize />
          </div>
        </div>
      </div>

      <div v-if="filteredRecords.length === 0" class="text-sm text-gray-500 dark:text-gray-400">暂无成绩记录</div>
      <div v-else class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div
          v-for="item in filteredRecords"
          :key="item.id"
          class="rounded-2xl border border-gray-100 dark:border-gray-700 p-4 shadow-sm transition hover:shadow-md cursor-pointer"
          :class="subjectCardClass(item.subject)"
          @click="openDetail(item.id)"
        >
          <div class="flex items-center justify-between">
            <div class="text-lg font-semibold text-gray-900 dark:text-white">{{ item.exam_name }}</div>
            <div class="text-xs text-gray-500 dark:text-gray-300">{{ formatDate(item.exam_date) }}</div>
          </div>
          <div class="mt-2 flex items-center gap-2 text-sm text-gray-700 dark:text-gray-200">
            <span class="px-2 py-0.5 rounded-full text-xs" :class="subjectBadgeClass(item.subject)">{{ item.subject }}</span>
            <span>{{ item.score }} / {{ item.full_score }}</span>
          </div>
          <div class="mt-2 text-xs text-gray-600 dark:text-gray-300 flex flex-wrap gap-3">
            <span>班级排名：{{ item.class_rank ?? '-' }}</span>
            <span>年级排名：{{ item.grade_rank ?? '-' }}</span>
            <span>班级均分：{{ item.class_avg ?? '-' }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ArrowLeft, DataAnalysis, TrendCharts } from '@element-plus/icons-vue'
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
</style>
