<template>
  <div class="min-h-screen bg-gradient-to-b from-sky-50 to-indigo-100">
    <div class="fixed top-0 left-0 right-0 bg-white/80 backdrop-blur z-40 border-b">
      <div class="px-4 py-2 flex items-center gap-2">
        <el-icon :size="18" class="cursor-pointer" @click="goBack"><ArrowLeft /></el-icon>
        <el-icon :size="18" class="text-indigo-600"><DataAnalysis /></el-icon>
        <div class="font-semibold text-indigo-700">任务统计</div>
      </div>
    </div>
    <div class="h-12"></div>

    <div class="p-4 space-y-4">
      <el-card shadow="never" class="rounded-xl overflow-hidden">
        <el-tabs v-model="activeTab" @tab-change="onTabChange">
          <el-tab-pane label="本周" name="week" />
          <el-tab-pane label="本月" name="month" />
          <el-tab-pane label="本年" name="year" />
          <el-tab-pane label="全部" name="all" />
        </el-tabs>
      </el-card>

      <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
        <el-card shadow="never" class="rounded-xl bg-sky-50">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2 text-sky-600"><el-icon :size="18"><List /></el-icon><span class="text-xs">总任务数</span></div>
            <div class="font-bold text-sky-700">{{ metrics.totalTasks }}</div>
          </div>
        </el-card>
        <el-card shadow="never" class="rounded-xl bg-amber-50">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2 text-amber-600"><el-icon :size="18"><Clock /></el-icon><span class="text-xs">总时长（分）</span></div>
            <div class="font-bold text-amber-700">{{ metrics.totalMinutes }}</div>
          </div>
        </el-card>
        <el-card shadow="never" class="rounded-xl bg-emerald-50">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2 text-emerald-600"><el-icon :size="18"><CircleCheck /></el-icon><span class="text-xs">已完成</span></div>
            <div class="font-bold text-emerald-700">{{ metrics.completed }}</div>
          </div>
        </el-card>
        <el-card shadow="never" class="rounded-xl bg-indigo-50">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2 text-indigo-600"><el-icon :size="18"><Clock /></el-icon><span class="text-xs">日均时长（分）</span></div>
            <div class="font-bold text-indigo-700">{{ metrics.avgDailyMinutes }}</div>
          </div>
        </el-card>
        <el-card shadow="never" class="rounded-xl bg-violet-50">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2 text-violet-600"><el-icon :size="18"><DataAnalysis /></el-icon><span class="text-xs">平均完成率</span></div>
            <div class="font-bold text-violet-700">{{ metrics.rate }}%</div>
          </div>
        </el-card>
        <el-card shadow="never" class="rounded-xl bg-pink-50">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2 text-pink-600"><el-icon :size="18"><Coin /></el-icon><span class="text-xs">总金币数</span></div>
            <div class="font-bold text-pink-700">{{ metrics.totalCoins }}</div>
          </div>
        </el-card>
      </div>

      <el-card shadow="never" class="rounded-xl">
        <template #header>
          <div class="flex items-center gap-2">
            <el-icon :size="18" class="text-emerald-600"><Histogram /></el-icon>
            <span class="font-semibold text-emerald-700">{{ seriesTitle }}</span>
          </div>
        </template>
        <div class="overflow-x-auto">
          <div class="h-40 flex items-end gap-2 justify-center">
            <div
              v-for="d in timeSeries"
              :key="d.key"
              class="w-6 relative"
              :title="d.key + ' ' + d.total + '分'"
            >
              <div class="h-[120px] flex flex-col-reverse">
                <div
                  v-for="seg in d.segments"
                  :key="seg.category"
                  :style="{ height: segHeight(seg.minutes) + 'px', backgroundColor: seg.color }"
                  :title="seg.category + '：' + seg.minutes + '分'"
                  class="rounded-t"
                ></div>
              </div>
              <div class="text-[10px] text-center mt-1">{{ labelOf(d.key) }}</div>
            </div>
          </div>
        </div>
      </el-card>

      <el-card shadow="never" class="rounded-xl">
        <template #header>
          <div class="flex items-center gap-2">
            <el-icon :size="18" class="text-pink-600"><PieChart /></el-icon>
            <span class="font-semibold text-pink-700">分类时间占比</span>
          </div>
        </template>
        <div class="flex flex-col items-center gap-4">
          <div class="w-40 h-40 rounded-full shadow-md" :style="{ background: donutGradient }"></div>
          <div class="flex flex-wrap justify-center gap-3">
            <div class="flex items-center gap-2" v-for="c in categoryShare" :key="c.name">
              <span class="inline-block w-3 h-3 rounded" :style="{ backgroundColor: c.color }"></span>
              <span class="text-sm">{{ c.name }}</span>
              <span class="text-xs text-gray-500">{{ c.percent }}%</span>
            </div>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import dayjs from 'dayjs'
import { listTasks, type TaskItem } from '@/services/tasks'
import router from '@/router'
import { ArrowLeft, DataAnalysis, Histogram, PieChart, List, Clock, CircleCheck, Coin } from '@element-plus/icons-vue'
import { useTaskCategories } from '@/stores/categories'
const cats = useTaskCategories()

function goBack() { router.back() }

const activeTab = ref<'week'|'month'|'year'|'all'>('week')
const allTasks = ref<TaskItem[]>([])

async function fetchAll() {
  const items: TaskItem[] = []
  let page = 1
  while (true) {
    const res = await listTasks({ page, page_size: 100 }) as any
    const arr: TaskItem[] = Array.isArray(res.items) ? res.items : []
    items.push(...arr)
    const total = Number(res.total || items.length)
    const size = Number(res.page_size || 100)
    if (page * size >= total || arr.length === 0) break
    page++
  }
  allTasks.value = items
}

function periodRange(tab: 'week'|'month'|'year'|'all') {
  const now = dayjs()
  if (tab === 'all') return { start: dayjs('1970-01-01'), end: dayjs('2999-12-31') }
  if (tab === 'year') return { start: now.startOf('year'), end: now.endOf('year') }
  if (tab === 'month') return { start: now.startOf('month'), end: now.endOf('month') }
  const weekday = now.day()
  const monday = now.subtract((weekday === 0 ? 6 : weekday - 1), 'day').startOf('day')
  const sunday = monday.add(6, 'day').endOf('day')
  return { start: monday, end: sunday }
}

const filtered = computed(() => {
  const { start, end } = periodRange(activeTab.value)
  return allTasks.value.filter(t => {
    const d = t.start_date ? dayjs(t.start_date) : null
    if (!d) return false
    return d.isAfter(start.subtract(1, 'millisecond')) && d.isBefore(end.add(1, 'millisecond'))
  })
})

const metrics = computed(() => {
  const tasks = filtered.value
  const totalTasks = tasks.length
  const completed = tasks.filter(t => t.status === 2).length
  const rate = totalTasks === 0 ? 0 : Math.round((completed / totalTasks) * 100)
  const totalMinutes = tasks.reduce((s, t) => s + (t.actual_minutes || 0), 0)
  const days = (() => {
    const { start, end } = periodRange(activeTab.value)
    return end.diff(start, 'day') + 1
  })()
  const avgDailyMinutes = days > 0 ? Math.round(totalMinutes / days) : 0
  const totalCoins = tasks.filter(t => t.status === 2).reduce((s, t) => s + (t.score || 0), 0)
  return { totalTasks, completed, rate, totalMinutes, avgDailyMinutes, totalCoins }
})

// 时间序列（根据周期动态：周/月=按日，年=按月，全部=按年）
const timeSeries = computed(() => {
  const tab = activeTab.value
  if (tab === 'year') {
    const now = dayjs()
    const year = now.year()
    const byMonth: Record<string, { [c: string]: number }> = {}
    for (let m = 1; m <= 12; m++) byMonth[`${year}-${String(m).padStart(2, '0')}`] = {}
    for (const t of filtered.value) {
      const key = dayjs(t.start_date).format('YYYY-MM')
      const c = (t.category || '未分类')
      const minutes = Number(t.actual_minutes || 0)
      byMonth[key][c] = (byMonth[key][c] || 0) + minutes
    }
    return Object.entries(byMonth).map(([key, m]) => {
      const segments = Object.entries(m).map(([category, minutes]) => ({ category, minutes: Number(minutes), color: cats.colorOf(category) }))
      const total = segments.reduce((s, x) => s + x.minutes, 0)
      return { key, segments, total }
    })
  }
  if (tab === 'all') {
    const years = Array.from(new Set(filtered.value.map(t => dayjs(t.start_date).year()))).sort((a, b) => a - b)
    const byYear: Record<string, { [c: string]: number }> = {}
    for (const y of years) byYear[String(y)] = {}
    for (const t of filtered.value) {
      const key = String(dayjs(t.start_date).year())
      const c = (t.category || '未分类')
      const minutes = Number(t.actual_minutes || 0)
      byYear[key][c] = (byYear[key][c] || 0) + minutes
    }
    return Object.entries(byYear).map(([key, m]) => {
      const segments = Object.entries(m).map(([category, minutes]) => ({ category, minutes: Number(minutes), color: cats.colorOf(category) }))
      const total = segments.reduce((s, x) => s + x.minutes, 0)
      return { key, segments, total }
    })
  }
  // week / month → 按日
  const { start, end } = periodRange(tab)
  const byDay: Record<string, { [c: string]: number }> = {}
  let p = start.startOf('day')
  while (p.isBefore(end.add(1, 'millisecond'))) {
    byDay[p.format('YYYY-MM-DD')] = {}
    p = p.add(1, 'day')
  }
  for (const t of filtered.value) {
    const key = dayjs(t.start_date).format('YYYY-MM-DD')
    const c = (t.category || '未分类')
    const minutes = Number(t.actual_minutes || 0)
    byDay[key][c] = (byDay[key][c] || 0) + minutes
  }
  return Object.entries(byDay).map(([key, m]) => {
    const segments = Object.entries(m).map(([category, minutes]) => ({ category, minutes: Number(minutes), color: cats.colorOf(category) }))
    const total = segments.reduce((s, x) => s + x.minutes, 0)
    return { key, segments, total }
  })
})

const maxSeriesTotal = computed(() => Math.max(1, ...timeSeries.value.map(r => r.total)))
function segHeight(minutes: number) {
  const hMax = 120
  return Math.round((minutes / maxSeriesTotal.value) * hMax)
}

const seriesTitle = computed(() => {
  if (activeTab.value === 'year') return '每月时间分布（分类维度）'
  if (activeTab.value === 'all') return '每年时间分布（分类维度）'
  return '每日时间分布（分类维度）'
})

function labelOf(key: string) {
  if (activeTab.value === 'year') {
    const d = dayjs(key + '-01')
    return d.format('MM月')
  }
  if (activeTab.value === 'all') {
    return `${key}年`
  }
  const d = dayjs(key)
  if (activeTab.value === 'month') return d.format('MM-DD')
  // week
  return ['一','二','三','四','五','六','日'][d.day() === 0 ? 6 : d.day()-1]
}

const categoryShare = computed(() => {
  const sumByCat: Record<string, number> = {}
  for (const t of filtered.value) {
    const c = (t.category || '未分类')
    sumByCat[c] = (sumByCat[c] || 0) + Number(t.actual_minutes || 0)
  }
  const total = Object.values(sumByCat).reduce((s, x) => s + x, 0) || 1
  const arr = Object.entries(sumByCat).map(([name, minutes]) => {
    const percent = Math.round((minutes / total) * 100)
    return { name, minutes, percent, color: cats.colorOf(name) }
  })
  return arr.sort((a, b) => b.minutes - a.minutes)
})

const donutGradient = computed(() => {
  let acc = 0
  const parts = categoryShare.value.map(c => {
    const start = acc
    const end = acc + c.percent
    acc = end
    return `${c.color} ${start}% ${end}%`
  })
  return `conic-gradient(${parts.join(', ')})`
})

onMounted(async () => {
  await fetchAll()
})

function onTabChange() {
  // 依赖计算属性自动刷新
}
</script>

<style scoped>
</style>