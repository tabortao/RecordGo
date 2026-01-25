<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-300">
    <div class="fixed top-0 left-0 right-0 stats-topbar z-40">
      <div class="px-4 py-2 flex items-center gap-2">
        <el-icon :size="18" class="cursor-pointer" @click="goBack"><ArrowLeft /></el-icon>
        <el-icon :size="18" class="text-emerald-600 dark:text-emerald-300"><DataAnalysis /></el-icon>
        <div class="font-semibold text-gray-900 dark:text-gray-100">任务统计</div>
      </div>
    </div>
    <div class="h-12"></div>

    <div class="p-4 space-y-4">
      <el-card shadow="never" class="stats-card overflow-hidden">
        <el-tabs v-model="activeTab" @tab-change="onTabChange" class="stats-tabs">
          <el-tab-pane label="本周" name="week" />
          <el-tab-pane label="本月" name="month" />
          <el-tab-pane label="本年" name="year" />
          <el-tab-pane label="全部" name="all" />
        </el-tabs>
      </el-card>

      <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
        <el-card shadow="never" class="stats-card">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2 text-sky-600 dark:text-sky-400"><el-icon :size="18"><List /></el-icon><span class="text-xs">总任务数</span></div>
            <div class="font-bold text-sky-700 dark:text-sky-400">{{ anim.totalTasks }}</div>
          </div>
        </el-card>
        <el-card shadow="never" class="stats-card">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2 text-emerald-600 dark:text-emerald-400"><el-icon :size="18"><Clock /></el-icon><span class="text-xs">总时长（分）</span></div>
            <div class="font-bold text-emerald-700 dark:text-emerald-400">{{ anim.totalMinutes }}</div>
          </div>
        </el-card>
        <el-card shadow="never" class="stats-card">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2 text-emerald-600 dark:text-emerald-400"><el-icon :size="18"><CircleCheck /></el-icon><span class="text-xs">已完成</span></div>
            <div class="font-bold text-emerald-700 dark:text-emerald-400">{{ metrics.completed }}</div>
          </div>
        </el-card>
        <el-card shadow="never" class="stats-card">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2 text-indigo-600 dark:text-indigo-400"><el-icon :size="18"><Clock /></el-icon><span class="text-xs">日均时长（分）</span></div>
            <div class="font-bold text-indigo-700 dark:text-indigo-400">{{ metrics.avgDailyMinutes }}</div>
          </div>
        </el-card>
        <el-card shadow="never" class="stats-card">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2 text-violet-600 dark:text-violet-400"><el-icon :size="18"><DataAnalysis /></el-icon><span class="text-xs">平均完成率</span></div>
            <div class="font-bold text-violet-700 dark:text-violet-400">{{ metrics.rate }}%</div>
          </div>
        </el-card>
        <el-card shadow="never" class="stats-card">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2 text-amber-600 dark:text-amber-400"><el-icon :size="18"><Coin /></el-icon><span class="text-xs">总金币数</span></div>
            <div class="font-bold text-amber-700 dark:text-amber-400">{{ anim.totalCoins }}</div>
          </div>
        </el-card>
      </div>

      <el-card shadow="never" class="stats-card">
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

      <el-card shadow="never" class="stats-card">
        <template #header>
          <div class="flex items-center gap-2">
            <el-icon :size="18" class="text-pink-600"><PieChart /></el-icon>
            <span class="font-semibold text-pink-700">分类时间占比</span>
          </div>
        </template>
        <div class="flex flex-col items-center gap-4">
          <div class="w-40 h-40 rounded-full shadow-md" :style="{ background: donutGradient }"></div>
          <div class="flex flex-wrap justify-center gap-3">
            <el-tooltip v-for="c in categoryShare" :key="c.name" :content="c.name + '：' + c.minutes + '分钟（' + c.percent + '%）'" placement="top" :hide-after="0">
              <div class="flex items-center gap-2">
                <span class="inline-block w-3 h-3 rounded" :style="{ backgroundColor: c.color }"></span>
                <span class="text-sm">{{ c.name }}</span>
                <span class="text-xs text-gray-500">{{ c.minutes }}分</span>
              </div>
            </el-tooltip>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import { listTasks, listTaskOccurrences, type TaskItem } from '@/services/tasks'
import { isAbortError } from '@/services/http'
import router from '@/router'
import { ArrowLeft, DataAnalysis, Histogram, PieChart, List, Clock, CircleCheck, Coin } from '@element-plus/icons-vue'
import { useTaskCategories } from '@/stores/categories'
const cats = useTaskCategories()

function goBack() { router.back() }

const activeTab = ref<'week'|'month'|'year'|'all'>('week')
const allTasks = ref<TaskItem[]>([])

async function fetchAll() {
  try {
    const res = await listTasks({ page_size: 1000 }) as any
    const arr: TaskItem[] = Array.isArray(res.items) ? res.items : []
    allTasks.value = arr
  } catch (e: any) {
    if (isAbortError(e)) return
    const msg = e?.response?.data?.message || e?.message || e
    console.error('统计页加载任务失败', { message: e?.message, status: e?.response?.status, payload: e?.response?.data })
    try { ElMessage.error(`统计数据加载失败：${msg}`) } catch {}
  }
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

const periodDays = computed(() => {
  const { start, end } = periodRange(activeTab.value)
  const keys: string[] = []
  let p = start.startOf('day')
  while (p.isBefore(end.add(1, 'millisecond'))) {
    keys.push(p.format('YYYY-MM-DD'))
    p = p.add(1, 'day')
  }
  return keys
})

const occMap = ref<Record<string, Record<number, { status: number; minutes: number }>>>({})
async function fetchOccurrences() {
  try {
    const { start, end } = periodRange(activeTab.value)
    const resp = await listTaskOccurrences({ start: start.format('YYYY-MM-DD'), end: end.format('YYYY-MM-DD') })
    const map: Record<string, Record<number, { status: number; minutes: number }>> = {}
    for (const it of (resp.items || [])) {
      const day = dayjs(it.date).format('YYYY-MM-DD')
      if (!map[day]) map[day] = {}
      map[day][Number(it.task_id)] = { status: Number(it.status || 0), minutes: Number(it.minutes || 0) }
    }
    occMap.value = map
  } catch (e: any) {
    if (isAbortError(e)) return
    const msg = e?.response?.data?.message || e?.message || e
    console.error('统计页加载任务发生失败', { message: e?.message, status: e?.response?.status, payload: e?.response?.data })
    try { ElMessage.error(`统计数据加载失败：${msg}`) } catch {}
  }
}

const filtered = computed(() => {
  const { start, end } = periodRange(activeTab.value)
  return allTasks.value.filter((t: TaskItem) => {
    const s = t.start_date ? dayjs(t.start_date) : null
    if (!s) return false
    return s.isBefore(end.add(1, 'millisecond')) && s.isAfter(start.subtract(1, 'millisecond'))
  })
})

const metrics = computed(() => {
  const tasks = filtered.value
  let totalTasks = 0
  let completed = 0
  let totalMinutes = 0
  let totalCoins = 0
  for (const day of periodDays.value) {
    const occ = occMap.value[day] || {}
    for (const t of tasks) {
      const isRepeat = !!((t as any).end_date || (t as any).series_id != null) && /daily|weekdays|weekly|monthly/i.test(String((t as any).repeat || (t as any).repeat_type || 'none'))
      if (isRepeat) {
        const entry = occ[t.id]
        if (entry) {
          totalTasks++
          if (entry.status === 2) { completed++; totalCoins += (t.score || 0) }
          totalMinutes += Number(entry.minutes || 0)
        }
        continue
      }
      const key = day
      if (dayjs(t.start_date).format('YYYY-MM-DD') === key) {
        totalTasks++
        if (t.status === 2) { completed++; totalCoins += (t.score || 0); totalMinutes += Number(t.actual_minutes || 0) }
      }
    }
  }
  const { start, end } = periodRange(activeTab.value)
  const days = end.diff(start, 'day') + 1
  const avgDailyMinutes = days > 0 ? Math.round(totalMinutes / days) : 0
  const rate = totalTasks === 0 ? 0 : Math.round((completed / totalTasks) * 100)
  return { totalTasks, completed, rate, totalMinutes, avgDailyMinutes, totalCoins }
})

// 数字动画（平滑过渡）
const anim = ref({ totalTasks: 0, totalMinutes: 0, totalCoins: 0 })
function animateNumber(key: 'totalTasks'|'totalMinutes'|'totalCoins', to: number) {
  const from = Number((anim.value as any)[key] || 0)
  const duration = 400
  const start = performance.now()
  const step = (now: number) => {
    const t = Math.min(1, (now - start) / duration)
    const val = Math.round(from + (to - from) * t)
    ;(anim.value as any)[key] = val
    if (t < 1) requestAnimationFrame(step)
  }
  requestAnimationFrame(step)
}
watch(metrics, (m) => {
  animateNumber('totalTasks', Number(m.totalTasks || 0))
  animateNumber('totalMinutes', Number(m.totalMinutes || 0))
  animateNumber('totalCoins', Number(m.totalCoins || 0))
}, { immediate: true })

// 时间序列（根据周期动态：周/月=按日，年=按月，全部=按年）
const timeSeries = computed(() => {
  const tab = activeTab.value
  if (tab === 'year') {
    const now = dayjs()
    const year = now.year()
    const byMonth: Record<string, { [c: string]: number }> = {}
    for (let m = 1; m <= 12; m++) byMonth[`${year}-${String(m).padStart(2, '0')}`] = {}
    for (const day of periodDays.value) {
      const occ = occMap.value[day] || {}
      const monthKey = dayjs(day).format('YYYY-MM')
      for (const t of filtered.value) {
        const isRepeat = !!((t as any).end_date || (t as any).series_id != null) && /daily|weekdays|weekly|monthly/i.test(String((t as any).repeat || (t as any).repeat_type || 'none'))
        const c = (t.category || '未分类')
        if (isRepeat) {
          const minutes = Number((occ[t.id]?.minutes) || 0)
          if (minutes > 0) byMonth[monthKey][c] = (byMonth[monthKey][c] || 0) + minutes
        } else if (dayjs(t.start_date).format('YYYY-MM-DD') === day && t.status === 2) {
          const minutes = Number(t.actual_minutes || 0)
          if (minutes > 0) byMonth[monthKey][c] = (byMonth[monthKey][c] || 0) + minutes
        }
      }
    }
    return Object.entries(byMonth).map(([key, m]) => {
      const segments = Object.entries(m).map(([category, minutes]) => ({ category, minutes: Number(minutes), color: cats.colorOf(category) }))
      const total = segments.reduce((s, x) => s + x.minutes, 0)
      return { key, segments, total }
    })
  }
  if (tab === 'all') {
    const years = Array.from(new Set<number>(filtered.value.map((t: TaskItem) => dayjs(t.start_date).year()))).sort((a, b) => a - b)
    const byYear: Record<string, { [c: string]: number }> = {}
    for (const y of years) byYear[String(y)] = {}
    for (const day of periodDays.value) {
      const occ = occMap.value[day] || {}
      const yearKey = dayjs(day).format('YYYY')
      for (const t of filtered.value) {
        const isRepeat = !!((t as any).end_date || (t as any).series_id != null) && /daily|weekdays|weekly|monthly/i.test(String((t as any).repeat || (t as any).repeat_type || 'none'))
        const c = (t.category || '未分类')
        if (isRepeat) {
          const minutes = Number((occ[t.id]?.minutes) || 0)
          if (minutes > 0) byYear[yearKey][c] = (byYear[yearKey][c] || 0) + minutes
        } else if (dayjs(t.start_date).format('YYYY-MM-DD') === day && t.status === 2) {
          const minutes = Number(t.actual_minutes || 0)
          if (minutes > 0) byYear[yearKey][c] = (byYear[yearKey][c] || 0) + minutes
        }
      }
    }
    return Object.entries(byYear).map(([key, m]) => {
      const segments = Object.entries(m).map(([category, minutes]) => ({ category, minutes: Number(minutes), color: cats.colorOf(category) }))
      const total = segments.reduce((s, x) => s + x.minutes, 0)
      return { key, segments, total }
    })
  }
  // week / month → 按日
  const byDay: Record<string, { [c: string]: number }> = {}
  for (const day of periodDays.value) byDay[day] = {}
  for (const day of periodDays.value) {
    const occ = occMap.value[day] || {}
    for (const t of filtered.value) {
      const c = (t.category || '未分类')
      const isRepeat = !!((t as any).end_date || (t as any).series_id != null) && /daily|weekdays|weekly|monthly/i.test(String((t as any).repeat || (t as any).repeat_type || 'none'))
      if (isRepeat) {
        const minutes = Number((occ[t.id]?.minutes) || 0)
        if (minutes > 0) byDay[day][c] = (byDay[day][c] || 0) + minutes
      } else if (dayjs(t.start_date).format('YYYY-MM-DD') === day && t.status === 2) {
        const minutes = Number(t.actual_minutes || 0)
        if (minutes > 0) byDay[day][c] = (byDay[day][c] || 0) + minutes
      }
    }
  }
  return Object.entries(byDay).map(([key, m]) => {
    const segments = Object.entries(m).map(([category, minutes]) => ({ category, minutes: Number(minutes), color: cats.colorOf(category) }))
    const total = segments.reduce((s, x) => s + x.minutes, 0)
    return { key, segments, total }
  })
})

const maxSeriesTotal = computed(() => Math.max(1, ...timeSeries.value.map((r: { total: number }) => r.total)))
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
  for (const d of timeSeries.value) {
    for (const seg of d.segments) {
      sumByCat[seg.category] = (sumByCat[seg.category] || 0) + Number(seg.minutes || 0)
    }
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
  const parts = categoryShare.value.map((c: { color: string; percent: number }) => {
    const start = acc
    const end = acc + c.percent
    acc = end
    return `${c.color} ${start}% ${end}%`
  })
  return `conic-gradient(${parts.join(', ')})`
})

onMounted(async () => {
  await fetchAll()
  await fetchOccurrences()
})

watch(activeTab, async () => {
  await fetchOccurrences()
})

function onTabChange() {
  // 依赖计算属性自动刷新
}
</script>

<style scoped>
:deep(.stats-card.el-card) {
  border-radius: 24px;
  border: 1px solid rgb(255 255 255 / 0.6);
  background: rgb(255 255 255 / 0.72);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  box-shadow: 0 16px 44px rgb(0 0 0 / 0.06);
}

.dark :deep(.stats-card.el-card) {
  border: 1px solid rgb(148 163 184 / 0.14);
  background: rgb(17 24 39 / 0.62);
  box-shadow: 0 20px 62px rgb(0 0 0 / 0.44);
}

.stats-topbar {
  background: rgb(255 255 255 / 0.72);
  border-bottom: 1px solid rgb(255 255 255 / 0.55);
  box-shadow: 0 12px 34px rgb(0 0 0 / 0.06);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
}

.dark .stats-topbar {
  background: rgb(17 24 39 / 0.72);
  border-bottom: 1px solid rgb(148 163 184 / 0.14);
  box-shadow: 0 14px 42px rgb(0 0 0 / 0.42);
}

:deep(.stats-tabs .el-tabs__header) {
  margin: 0;
}

:deep(.stats-tabs .el-tabs__nav-wrap::after) {
  height: 0;
}

:deep(.stats-tabs .el-tabs__active-bar) {
  display: none;
}

:deep(.stats-tabs .el-tabs__nav) {
  background: rgb(255 255 255 / 0.7);
  border: 1px solid rgb(255 255 255 / 0.55);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-radius: 9999px;
  padding: 4px;
}

.dark :deep(.stats-tabs .el-tabs__nav) {
  background: rgb(17 24 39 / 0.6);
  border: 1px solid rgb(148 163 184 / 0.14);
}

:deep(.stats-tabs .el-tabs__item) {
  height: 34px;
  line-height: 34px;
  font-size: 13px;
  font-weight: 800;
  border-radius: 9999px;
  padding: 0 14px;
  color: rgb(55 65 81);
  transition: background-color 180ms ease, color 180ms ease, box-shadow 220ms ease, transform 220ms ease;
}

.dark :deep(.stats-tabs .el-tabs__item) {
  color: rgb(229 231 235);
}

:deep(.stats-tabs .el-tabs__item:hover) {
  background: rgb(255 255 255 / 0.55);
}

.dark :deep(.stats-tabs .el-tabs__item:hover) {
  background: rgb(31 41 55 / 0.45);
}

:deep(.stats-tabs .el-tabs__item.is-active) {
  color: rgb(4 120 87);
  background: rgb(16 185 129 / 0.14);
  border: 1px solid rgb(16 185 129 / 0.16);
  box-shadow: 0 16px 34px rgb(16 185 129 / 0.12);
  transform: translateY(-0.5px);
}

.dark :deep(.stats-tabs .el-tabs__item.is-active) {
  color: rgb(167 243 208);
  background: rgb(16 185 129 / 0.16);
  border: 1px solid rgb(16 185 129 / 0.14);
  box-shadow: 0 18px 40px rgb(0 0 0 / 0.5);
}
</style>
