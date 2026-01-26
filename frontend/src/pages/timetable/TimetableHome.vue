<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 pb-20 relative overflow-hidden print:bg-white print:pb-0 print:min-h-auto">
    <div class="pointer-events-none absolute inset-0 overflow-hidden print:hidden">
      <div class="absolute -top-24 -right-24 h-72 w-72 rounded-full bg-sky-300/35 dark:opacity-0 blur-3xl" />
      <div class="absolute -bottom-40 -left-28 h-80 w-80 rounded-full bg-amber-200/35 dark:opacity-0 blur-3xl" />
      <div class="absolute inset-0 bg-[radial-gradient(1200px_circle_at_20%_-20%,rgba(255,255,255,.65),transparent_55%),radial-gradient(900px_circle_at_80%_0%,rgba(255,255,255,.45),transparent_55%)] dark:bg-transparent" />
    </div>

    <!-- Header -->
    <div class="sticky top-0 z-20 px-4 pt-4 print:px-0 print:pt-0 print:static">
      <div class="rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/75 dark:bg-gray-900/70 backdrop-blur-xl shadow-sm px-3 py-3 flex items-center justify-between relative print:shadow-none print:rounded-none print:border-black print:bg-transparent print:px-0 print:py-0">
        <div class="flex items-center gap-3 z-10 print:hidden">
          <button
            type="button"
            class="inline-flex h-10 w-10 items-center justify-center rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 text-gray-600 dark:text-gray-300 hover:bg-white hover:text-gray-900 dark:hover:text-gray-50 transition-colors"
            @click="router.back()"
          >
            <el-icon :size="20"><ArrowLeft /></el-icon>
          </button>
          <span class="font-extrabold tracking-tight text-gray-900 dark:text-gray-50 text-sm md:text-base">{{ config?.current_grade }} {{ config?.current_semester }}</span>
        </div>
      
      <!-- Print Title (Only visible in print) -->
      <div class="hidden print:block text-center w-full mb-4">
        <h1 class="font-bold text-2xl text-black">
          {{ config?.current_grade }} {{ config?.current_semester }} {{ authStore.user?.nickname || authStore.user?.username }}的课表
        </h1>
      </div>

        <h1 class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 font-extrabold tracking-tight text-gray-900 dark:text-gray-50 text-sm md:text-base print:hidden">
          {{ authStore.user?.nickname || authStore.user?.username }}的课表
        </h1>

        <div class="z-10 flex items-center gap-2 print:hidden">
          <button
            type="button"
            class="inline-flex h-10 w-10 items-center justify-center rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 text-gray-600 dark:text-gray-300 hover:bg-white hover:text-gray-900 dark:hover:text-gray-50 transition-colors"
            @click="handlePrint"
          >
            <el-icon :size="18"><Printer /></el-icon>
          </button>
          <button
            type="button"
            class="inline-flex h-10 w-10 items-center justify-center rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 text-gray-600 dark:text-gray-300 hover:bg-white hover:text-gray-900 dark:hover:text-gray-50 transition-colors"
            @click="router.push('/timetable/settings')"
          >
            <el-icon :size="18"><Setting /></el-icon>
          </button>
        </div>
      </div>
    </div>

    <!-- Timetable -->
    <div class="p-3 overflow-x-auto relative z-10 print:p-0 print:overflow-visible">
      <div class="min-w-[320px] max-w-6xl mx-auto print:w-full">
        <!-- Week Header -->
        <div class="grid gap-1 mb-1 print:gap-0 print:mb-0 print:border-b print:border-black" :style="{ gridTemplateColumns: gridColumns }">
          <div class="w-16 print:border-r print:border-black"></div> <!-- Period Column Placeholder -->
          <div v-for="day in days" :key="day.value" class="text-center py-2 font-extrabold text-sm md:text-base bg-white/70 dark:bg-gray-900/55 border border-white/50 dark:border-gray-800/60 rounded-2xl text-gray-800 dark:text-gray-100 print:bg-transparent print:text-black print:border-r print:border-black print:rounded-none print:py-2 print:text-base">
            {{ day.label }}
          </div>
        </div>

        <!-- Periods -->
        <div v-for="period in periods" :key="period" class="grid gap-1 mb-1 print:gap-0 print:mb-0 print:border-b print:border-black" :style="{ gridTemplateColumns: gridColumns }">
          <!-- Period Number & Time -->
          <div class="flex flex-col items-center justify-center text-xs text-gray-600 dark:text-gray-300 font-semibold w-16 bg-white/70 dark:bg-gray-900/55 border border-white/50 dark:border-gray-800/60 rounded-2xl py-2 print:bg-transparent print:text-black print:border-r print:border-black print:rounded-none">
            <span class="text-sm md:text-base font-extrabold print:text-base">{{ period }}</span>
            <div v-if="getPeriodTime(period)" class="flex flex-col items-center text-[10px] text-gray-500 dark:text-gray-400 leading-tight mt-1 print:text-black print:text-xs">
                <span>{{ getPeriodTime(period)?.start }}</span>
                <span>{{ getPeriodTime(period)?.end }}</span>
            </div>
          </div>
          <!-- Course Cells -->
          <div 
            v-for="day in days" 
            :key="`${day.value}-${period}`"
            class="h-16 md:h-20 rounded-2xl p-2 flex items-center justify-center text-center text-sm md:text-base font-extrabold shadow-sm transition-colors relative overflow-hidden bg-white/70 dark:bg-gray-900/55 border border-white/50 dark:border-gray-800/60 print:shadow-none print:rounded-none print:border-r print:border-black print:h-20 print:text-base"
            :style="getCellStyle(day.value, period)"
          >
            <span class="z-10 break-words w-full">{{ getCourseName(day.value, period) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, type CSSProperties } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Setting, Printer } from '@element-plus/icons-vue'
import { useTimetableStore } from '@/stores/timetable'
import { useAuth } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import type { PeriodSetting } from '@/api/timetable'

const router = useRouter()
const store = useTimetableStore()
const authStore = useAuth()
const { config, timetable } = storeToRefs(store)

const periods = Array.from({ length: 10 }, (_, i) => i + 1) // 假设每天10节课
const periodSettings = ref<PeriodSetting[]>([])
const courseColors = ref<Record<string, string>>({})

const days = computed(() => {
  const allDays = [
    { label: '周一', value: 1 },
    { label: '周二', value: 2 },
    { label: '周三', value: 3 },
    { label: '周四', value: 4 },
    { label: '周五', value: 5 },
  ]
  if (config.value?.show_saturday) allDays.push({ label: '周六', value: 6 })
  if (config.value?.show_sunday) allDays.push({ label: '周日', value: 7 })
  return allDays
})

const gridColumns = computed(() => {
  return `4rem repeat(${days.value.length}, 1fr)`
})

function getPeriodTime(period: number) {
    const s = periodSettings.value.find(p => p.period === period)
    if (s && s.start_time && s.end_time) {
        return { start: s.start_time, end: s.end_time }
    }
    return null
}

function getCourse(day: number, period: number) {
  return timetable.value.find(t => t.day_of_week === day && t.period === period)
}

function getCourseName(day: number, period: number) {
  const item = getCourse(day, period)
  return item?.course?.name || ''
}

function getCellStyle(day: number, period: number): CSSProperties {
  const item = getCourse(day, period)
  if (!item || !item.course) return {}
  
  const color = courseColors.value[item.course.name] || item.course.color

  return {
    backgroundColor: color,
    color: isLightColor(color) ? '#000' : '#fff',
    // In print mode, we want exact colors
    printColorAdjust: 'exact',
    WebkitPrintColorAdjust: 'exact'
  } as any
}

// 简单的判断颜色深浅
function isLightColor(color: string) {
    if (!color) return true
    const hex = color.replace('#', '')
    const r = parseInt(hex.substring(0, 2), 16)
    const g = parseInt(hex.substring(2, 4), 16)
    const b = parseInt(hex.substring(4, 6), 16)
    const brightness = (r * 299 + g * 587 + b * 114) / 1000
    return brightness > 155 // 阈值可调
}

function handlePrint() {
  window.print()
}

onMounted(async () => {
  await store.fetchConfig()
  if (config.value) {
    // 解析时间设置
    if (config.value.period_settings_json) {
        try {
            periodSettings.value = JSON.parse(config.value.period_settings_json)
        } catch (e) {
            console.error(e)
        }
    }

    if (config.value.course_colors_json) {
        try {
            courseColors.value = JSON.parse(config.value.course_colors_json)
        } catch (e) {
            console.error(e)
        }
    }
    
    await store.fetchTimetable(config.value.current_grade, config.value.current_semester)
  }
})
</script>

<style scoped>
@media print {
  @page {
    size: A4 landscape;
    margin: 10mm;
  }
  
  body {
    -webkit-print-color-adjust: exact !important;
    print-color-adjust: exact !important;
    background-color: white !important;
  }

  /* 强制单页打印 */
  .print\:min-h-auto {
    min-height: auto !important;
    height: auto !important;
  }
  
  /* 调整字号以适应 A4 */
  .print\:text-base {
    font-size: 10pt !important;
  }
  
  /* 隐藏滚动条 */
  ::-webkit-scrollbar {
    display: none;
  }
}
</style>
