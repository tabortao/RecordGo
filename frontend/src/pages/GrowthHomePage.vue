<template>
  <div class="min-h-screen bg-[#F5F7FA] dark:bg-gray-900">
    <div class="bg-white/80 dark:bg-gray-900/80 backdrop-blur border-b border-gray-200 dark:border-gray-800 sticky top-0 z-10">
      <div class="px-4 py-3 flex items-center gap-2">
        <el-icon :size="18" class="cursor-pointer text-gray-600 dark:text-gray-300" @click="router.back()"><ArrowLeft /></el-icon>
        <h2 class="font-semibold text-gray-800 dark:text-gray-100">生长</h2>
      </div>
    </div>

    <div class="p-4 space-y-4">
      <div class="rounded-xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4">
        <div class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ displayName }}</div>
        <div class="mt-2 grid grid-cols-2 gap-3 text-sm text-gray-600 dark:text-gray-300">
          <div class="flex items-center gap-2">
            <span class="text-gray-500 dark:text-gray-400">年龄</span>
            <span>{{ ageText }}</span>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-gray-500 dark:text-gray-400">性别</span>
            <span>{{ genderText }}</span>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div
          v-for="card in metricCards"
          :key="card.type"
          class="rounded-xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4 cursor-pointer hover:shadow-md transition"
          @click="openRecords(card.type)"
        >
          <div class="flex items-center justify-between">
            <div class="font-semibold text-gray-800 dark:text-gray-100">{{ card.name }}</div>
            <el-icon :size="18" class="text-blue-500"><TrendCharts /></el-icon>
          </div>
          <div class="mt-3 text-sm text-gray-500 dark:text-gray-400">
            <span v-if="latestByType[card.type]">
              {{ latestByType[card.type]?.value }} {{ card.unit }}
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
import { ArrowLeft, TrendCharts } from '@element-plus/icons-vue'
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
  if (g === 'other') return '其他'
  return '未知'
})

const metricCards: Array<{ type: GrowthMetricType; name: string; unit: string }> = [
  { type: 'height', name: '身高', unit: 'cm' },
  { type: 'weight', name: '体重', unit: 'kg' },
  { type: 'vision', name: '裸眼视力', unit: '' },
  { type: 'foot', name: '脚长', unit: 'cm' }
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
