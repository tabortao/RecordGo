<template>
  <div class="min-h-screen bg-[#F5F7FA] dark:bg-gray-900">
    <div class="bg-white/80 dark:bg-gray-900/80 backdrop-blur border-b border-gray-200 dark:border-gray-800 sticky top-0 z-10">
      <div class="px-4 py-3 flex items-center gap-2">
        <el-icon :size="18" class="cursor-pointer text-gray-600 dark:text-gray-300" @click="router.back()"><ArrowLeft /></el-icon>
        <h2 class="font-semibold text-gray-800 dark:text-gray-100">生长</h2>
      </div>
    </div>

    <div class="p-4 space-y-4">
      <div class="rounded-2xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4 shadow-sm">
        <div class="flex items-center justify-between gap-4">
          <div class="flex-1">
            <div class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ displayName }}</div>
            <div class="mt-2 flex items-center gap-4 text-sm text-gray-600 dark:text-gray-300">
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
          <div class="relative w-24 h-24">
            <div class="absolute inset-0 rounded-full" :style="bmiRingStyle"></div>
            <div class="absolute inset-2 rounded-full bg-white dark:bg-gray-800 flex flex-col items-center justify-center">
              <div class="text-xs text-gray-500 dark:text-gray-400">BMI</div>
              <div class="text-base font-semibold" :class="bmiColorClass">{{ bmiText }}</div>
            </div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div
          v-for="card in metricCards"
          :key="card.type"
          class="rounded-2xl border border-transparent p-4 cursor-pointer transition hover:shadow-md"
          :class="card.bgClass"
          @click="openRecords(card.type)"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2 font-semibold" :class="card.textClass">
              <el-icon :size="18" :class="card.iconClass"><component :is="card.icon" /></el-icon>
              <span>{{ card.name }}</span>
            </div>
            <el-icon :size="18" class="text-blue-500"><TrendCharts /></el-icon>
          </div>
          <div class="mt-3 text-sm text-gray-600 dark:text-gray-300">
            <span v-if="latestByType[card.type]">
              {{ formatValue(card.type, latestByType[card.type]!) }} {{ card.unit }}
            </span>
            <span v-else>暂无记录</span>
          </div>
          <div class="mt-1 text-xs text-gray-400">
            <span v-if="latestByType[card.type]">{{ formatDate(latestByType[card.type]?.record_date) }}</span>
            <span v-else>点击添加记录</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ArrowLeft, TrendCharts, Histogram, Medal, View, Position } from '@element-plus/icons-vue'
import { useAuth } from '@/stores/auth'
import { listGrowthRecords, type GrowthMetricRecord, type GrowthMetricType } from '@/services/growth'
import dayjs from 'dayjs'
import router from '@/router'

const auth = useAuth()

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

const metricCards: Array<{ type: GrowthMetricType; name: string; unit: string; icon: any; textClass: string; iconClass: string; bgClass: string }> = [
  { type: 'height', name: '身高', unit: 'cm', icon: Histogram, textClass: 'text-emerald-700 dark:text-emerald-300', iconClass: 'text-emerald-500', bgClass: 'bg-emerald-50/70 dark:bg-emerald-500/10' },
  { type: 'weight', name: '体重', unit: 'kg', icon: Medal, textClass: 'text-orange-600 dark:text-orange-300', iconClass: 'text-orange-500', bgClass: 'bg-orange-50/70 dark:bg-orange-500/10' },
  { type: 'vision', name: '裸眼视力', unit: '', icon: View, textClass: 'text-blue-600 dark:text-blue-300', iconClass: 'text-blue-500', bgClass: 'bg-blue-50/70 dark:bg-blue-500/10' },
  { type: 'foot', name: '脚长', unit: 'cm', icon: Position, textClass: 'text-purple-700 dark:text-purple-300', iconClass: 'text-purple-500', bgClass: 'bg-purple-50/70 dark:bg-purple-500/10' }
]

const latestByType = ref<Record<GrowthMetricType, GrowthMetricRecord | null>>({
  height: null,
  weight: null,
  vision: null,
  foot: null
})

function formatDate(d?: string) {
  if (!d) return ''
  return dayjs(d).format('YYYY-MM-DD')
}

function formatValue(type: GrowthMetricType, record: GrowthMetricRecord) {
  if (type === 'vision') {
    if (record.left_value && record.right_value) {
      return `左${record.left_value} 右${record.right_value}`
    }
    return ''
  }
  return record.value ?? ''
}

const bmiValue = computed(() => {
  const h = latestByType.value.height?.value
  const w = latestByType.value.weight?.value
  if (!h || !w) return null
  const m = h / 100
  if (!m) return null
  const bmi = w / (m * m)
  if (!Number.isFinite(bmi)) return null
  return Number(bmi.toFixed(1))
})

const bmiText = computed(() => {
  if (bmiValue.value === null) return '未计算'
  return bmiValue.value.toFixed(1)
})

const bmiPercent = computed(() => {
  if (bmiValue.value === null) return 0
  const min = 10
  const max = 35
  const clamped = Math.min(max, Math.max(min, bmiValue.value))
  return (clamped - min) / (max - min)
})

const bmiColorClass = computed(() => {
  if (bmiValue.value === null) return 'text-gray-400'
  if (bmiValue.value < 18.5) return 'text-sky-600 dark:text-sky-300'
  if (bmiValue.value < 24) return 'text-emerald-600 dark:text-emerald-300'
  if (bmiValue.value < 28) return 'text-orange-500 dark:text-orange-300'
  return 'text-rose-500 dark:text-rose-300'
})

const bmiRingStyle = computed(() => {
  const deg = Math.round(bmiPercent.value * 360)
  return {
    background: `conic-gradient(#10b981 ${deg}deg, #e5e7eb 0deg)`
  }
})

function openRecords(type: GrowthMetricType) {
  router.push(`/growth/records/${type}`)
}

async function load() {
  const resp = await listGrowthRecords()
  const list = resp.items || []
  const map: Record<GrowthMetricType, GrowthMetricRecord | null> = {
    height: null,
    weight: null,
    vision: null,
    foot: null
  }
  for (const item of list) {
    const t = item.metric_type as GrowthMetricType
    const current = map[t]
    if (!current) {
      map[t] = item
      continue
    }
    const a = dayjs(item.record_date)
    const b = dayjs(current.record_date)
    if (a.isAfter(b)) map[t] = item
  }
  latestByType.value = map
}

onMounted(load)
</script>

<style scoped>
</style>
