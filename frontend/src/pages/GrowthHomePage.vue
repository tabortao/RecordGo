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
            <div class="flex items-center gap-3">
              <el-avatar :size="44" :src="avatarSrc" />
              <div class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ displayName }}</div>
            </div>
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
          <div class="w-56 rounded-xl border border-gray-100 dark:border-gray-700 bg-white/70 dark:bg-gray-900/40 p-3">
            <div class="flex items-center justify-between">
              <div class="text-xs text-gray-500 dark:text-gray-400">BMI</div>
              <el-icon class="cursor-pointer text-gray-400 hover:text-blue-500" @click="openBmiInfo"><QuestionFilled /></el-icon>
            </div>
            <div class="mt-3 flex items-center justify-center">
              <svg viewBox="0 0 200 200" class="h-36 w-36">
                <path
                  d="M 20 140 A 80 80 0 0 1 180 140"
                  stroke="currentColor"
                  stroke-width="14"
                  fill="none"
                  class="text-gray-200 dark:text-gray-700"
                />
                <path
                  v-for="seg in bmiArcSegments"
                  :key="seg.color"
                  :d="seg.path"
                  :stroke="seg.color"
                  stroke-width="14"
                  fill="none"
                  stroke-linecap="butt"
                />
                <line
                  v-if="bmiPointer"
                  x1="100"
                  y1="140"
                  :x2="bmiPointer.x"
                  :y2="bmiPointer.y"
                  stroke="#ef4444"
                  stroke-width="4"
                  stroke-linecap="round"
                />
                <circle v-if="bmiPointer" :cx="bmiPointer.x" :cy="bmiPointer.y" r="5" fill="#ef4444" />
                <text x="100" y="112" text-anchor="middle" fill="currentColor" class="text-[12px] text-gray-500 dark:text-gray-400">BMI</text>
                <text x="100" y="132" text-anchor="middle" fill="currentColor" class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ displayBmiText }}</text>
              </svg>
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
import { computed, onMounted, ref, watch } from 'vue'
import { ArrowLeft, TrendCharts, Histogram, Medal, View, Position, QuestionFilled } from '@element-plus/icons-vue'
import { useAuth } from '@/stores/auth'
import { listGrowthRecords, type GrowthMetricRecord, type GrowthMetricType } from '@/services/growth'
import dayjs from 'dayjs'
import router from '@/router'
import defaultAvatar from '@/assets/avatars/default.png'
import { getStaticBase } from '@/services/http'
import { presignView } from '@/services/storage'

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

const avatarSrc = ref<string>(defaultAvatar)
async function updateAvatar() {
  const p = auth.user?.avatar_path
  if (!p) { avatarSrc.value = defaultAvatar; return }
  const s = String(p)
  if (/storage\/images\/avatars\/default\.png$/i.test(s) || /(^|\/)default\.png$/i.test(s)) { avatarSrc.value = defaultAvatar; return }
  if (/^https?:\/\//i.test(s)) { avatarSrc.value = s; return }
  const base = getStaticBase()
  if (/uploads\//i.test(s)) { avatarSrc.value = `${base}/api/${s.replace(/^\/+/, '')}`; return }
  try { avatarSrc.value = await presignView(s) } catch { avatarSrc.value = defaultAvatar }
}

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

const displayBmiValue = computed(() => (bmiValue.value ?? 16.64))
const displayBmiText = computed(() => (bmiValue.value === null ? '16.64' : bmiValue.value.toFixed(1)))

const bmiSegments = [
  { start: 0, end: 18.5, color: '#3b82f6' },
  { start: 18.5, end: 25, color: '#22c55e' },
  { start: 25, end: 30, color: '#f59e0b' },
  { start: 30, end: 40, color: '#fca5a5' },
  { start: 40, end: 50, color: '#ef4444' }
]

function valueToAngle(value: number) {
  const v = Math.max(0, Math.min(50, value))
  return 180 - (v / 50) * 180
}

function polarToCartesian(cx: number, cy: number, r: number, angle: number) {
  const rad = (angle - 90) * (Math.PI / 180)
  return { x: cx + r * Math.cos(rad), y: cy + r * Math.sin(rad) }
}

function describeArc(cx: number, cy: number, r: number, startAngle: number, endAngle: number) {
  const start = polarToCartesian(cx, cy, r, endAngle)
  const end = polarToCartesian(cx, cy, r, startAngle)
  const largeArcFlag = endAngle - startAngle <= 180 ? '0' : '1'
  return `M ${start.x} ${start.y} A ${r} ${r} 0 ${largeArcFlag} 0 ${end.x} ${end.y}`
}

const bmiArcSegments = computed(() => {
  return bmiSegments.map((seg) => {
    const startAngle = valueToAngle(seg.start)
    const endAngle = valueToAngle(seg.end)
    return { color: seg.color, path: describeArc(100, 140, 80, startAngle, endAngle) }
  })
})

const bmiPointer = computed(() => {
  const angle = valueToAngle(displayBmiValue.value)
  return polarToCartesian(100, 140, 62, angle)
})

function openRecords(type: GrowthMetricType) {
  router.push(`/growth/records/${type}`)
}

function openBmiInfo() {
  router.push('/growth/bmi-info')
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

onMounted(() => {
  updateAvatar()
  load()
})
watch(() => auth.user?.avatar_path, () => { updateAvatar() })
</script>

<style scoped>
</style>
