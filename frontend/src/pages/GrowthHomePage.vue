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
        <h2 class="font-bold text-gray-800 dark:text-gray-100 text-lg">生长记录</h2>
      </div>
      <!-- Right Placeholder for potential future actions -->
      <div class="w-8"></div>
    </div>

    <div class="flex-1 overflow-y-auto p-4 sm:p-6 pb-24">
      <div class="max-w-4xl mx-auto space-y-6">
        <!-- Profile & BMI Card -->
        <div class="bg-white dark:bg-gray-800 rounded-3xl p-4 sm:p-6 shadow-sm border border-gray-100 dark:border-gray-700 relative overflow-hidden">
          <!-- Background Decoration -->
          <div class="absolute top-0 right-0 w-64 h-64 bg-gradient-to-bl from-blue-50 to-transparent dark:from-blue-900/10 rounded-bl-[100px] -z-0"></div>
          
          <div class="relative z-10 flex flex-row items-stretch gap-4 sm:gap-8">
            <!-- Left: Profile Info -->
            <div class="flex-1 min-w-0 flex items-center gap-4 sm:gap-5">
               <div class="relative">
                 <el-avatar :size="72" :src="avatarSrc" class="border-4 border-white dark:border-gray-700 shadow-md scale-[0.9] sm:scale-100 origin-left" />
               </div>
               <div>
                 <h1 class="text-xl sm:text-2xl font-bold text-gray-900 dark:text-white mb-1.5 truncate">{{ displayName }}</h1>
                 <div class="flex flex-wrap gap-2">
                   <span class="px-2.5 py-0.5 rounded-lg bg-gray-100 dark:bg-gray-700 text-xs font-bold text-gray-600 dark:text-gray-300">{{ ageText }}</span>
                   <span class="px-2.5 py-0.5 rounded-lg bg-gray-100 dark:bg-gray-700 text-xs font-bold text-gray-600 dark:text-gray-300">{{ genderText }}</span>
                 </div>
               </div>
            </div>

            <!-- Right: BMI Status -->
            <div class="w-40 sm:w-64 shrink-0 bg-white/60 dark:bg-gray-700/30 rounded-2xl p-3 sm:p-4 border border-gray-100 dark:border-gray-600/30 backdrop-blur-sm">
               <div class="flex items-center justify-between mb-2 sm:mb-3">
                 <span class="text-xs font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider">BMI 指数</span>
                 <el-tooltip content="点击查看 BMI 详情与建议" placement="top">
                   <div class="cursor-pointer text-gray-400 hover:text-blue-500 transition-colors" @click="openBmiInfo">
                     <el-icon><QuestionFilled /></el-icon>
                   </div>
                 </el-tooltip>
               </div>
               
               <div class="flex items-center gap-4">
                 <div class="relative">
                   <div class="w-12 h-12 sm:w-14 sm:h-14 rounded-full bg-blue-50 dark:bg-blue-900/20 flex items-center justify-center overflow-hidden">
                     <img :src="genderAvatarSrc" class="w-9 h-9 sm:w-10 sm:h-10 object-contain" alt="Gender Avatar" />
                   </div>
                 </div>
                 <div>
                   <div v-if="bmiValue !== null" class="text-2xl sm:text-3xl font-black text-gray-800 dark:text-white leading-none mb-1">{{ bmiValue.toFixed(1) }}</div>
                   <div v-else class="text-[11px] font-bold text-gray-400 dark:text-gray-500 leading-tight mb-1">录入身高与体重后显示</div>
                   <div
                     v-if="bmiStatusText"
                     class="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-bold"
                     :class="{
                       'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300': bmiStatusText === '偏瘦',
                       'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-300': bmiStatusText === '正常',
                       'bg-orange-100 text-orange-700 dark:bg-orange-900/40 dark:text-orange-300': bmiStatusText === '偏胖',
                       'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300': bmiStatusText === '肥胖' || bmiStatusText === '高度肥胖'
                     }"
                   >
                     {{ bmiStatusText }}
                   </div>
                 </div>
               </div>
            </div>
          </div>
        </div>

        <!-- Metrics Grid -->
        <div>
           <h3 class="text-sm font-bold text-gray-500 dark:text-gray-400 mb-4 ml-1 uppercase tracking-wider">各项指标记录</h3>
           <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
             <div
               v-for="card in metricCards"
               :key="card.type"
               class="group relative bg-white dark:bg-gray-800 rounded-3xl p-5 border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-xl hover:-translate-y-1 transition-all duration-300 cursor-pointer overflow-hidden"
               @click="openRecords(card.type)"
             >
               <!-- Card Bg Decoration -->
               <div class="absolute top-0 right-0 w-24 h-24 rounded-bl-full opacity-20 transition-transform group-hover:scale-110" :class="card.bgClass"></div>
               
               <div class="relative z-10 flex flex-col h-full">
                 <div class="flex items-center justify-between mb-4">
                   <div class="w-10 h-10 rounded-2xl flex items-center justify-center transition-colors" :class="card.bgClass">
                      <el-icon :size="20" :class="card.iconClass"><component :is="card.icon" /></el-icon>
                   </div>
                   <div class="w-8 h-8 rounded-full bg-gray-50 dark:bg-gray-700 flex items-center justify-center text-gray-300 group-hover:text-blue-500 group-hover:bg-blue-50 dark:group-hover:bg-blue-900/20 transition-all">
                     <el-icon><ArrowRight /></el-icon>
                   </div>
                 </div>
                 
                 <div class="mb-1">
                    <span class="text-sm font-medium text-gray-500 dark:text-gray-400">{{ card.name }}</span>
                 </div>
                 
                 <div class="flex items-baseline gap-1 mt-auto">
                    <template v-if="latestByType[card.type]">
                      <span class="text-2xl font-black text-gray-900 dark:text-white tracking-tight">
                        {{ formatValue(card.type, latestByType[card.type]!) }}
                      </span>
                      <span class="text-sm font-bold text-gray-400">{{ card.unit }}</span>
                    </template>
                    <span v-else class="text-lg font-bold text-gray-300 dark:text-gray-600">暂无记录</span>
                 </div>
                 
                 <div class="mt-2 flex items-center gap-1.5 text-xs font-medium text-gray-400">
                   <el-icon><Clock /></el-icon>
                   <span v-if="latestByType[card.type]">更新于 {{ formatDate(latestByType[card.type]?.record_date) }}</span>
                   <span v-else>点击添加第一条记录</span>
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
import { computed, onMounted, ref, watch } from 'vue'
import { ArrowLeft, Histogram, Medal, View, Position, QuestionFilled, ArrowRight, Clock } from '@element-plus/icons-vue'
import { useAuth } from '@/stores/auth'
import { listGrowthRecords, type GrowthMetricRecord, type GrowthMetricType } from '@/services/growth'
import dayjs from 'dayjs'
import router from '@/router'
import defaultAvatar from '@/assets/avatars/default.png'
import boyAvatar from '@/assets/avatars/boy.svg'
import girlAvatar from '@/assets/avatars/girl.svg'
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

const genderAvatarSrc = computed(() => {
  const g = auth.user?.child_gender || ''
  if (g === 'male') return boyAvatar
  return girlAvatar
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

const bmiStatusText = computed(() => {
  const v = bmiValue.value
  if (v === null) return ''
  if (v < 18.5) return '偏瘦'
  if (v < 25) return '正常'
  if (v < 30) return '偏胖'
  if (v < 40) return '肥胖'
  return '高度肥胖'
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
