<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 pb-20 relative overflow-hidden print:bg-white print:pb-0 print:min-h-auto">
    <!-- Header -->
    <div class="bg-white dark:bg-gray-800 shadow-sm px-4 py-3 flex items-center justify-between sticky top-0 z-20 relative print:shadow-none print:px-0 print:static">
      <div class="flex items-center gap-3 z-10 print:hidden">
        <el-icon :size="20" class="cursor-pointer" @click="router.back()"><ArrowLeft /></el-icon>
        <span class="font-bold text-base text-[#333333]">{{ config?.current_grade }} {{ config?.current_semester }}</span>
      </div>
      
      <!-- Print Title (Only visible in print) -->
      <div class="hidden print:block text-center w-full mb-4">
        <h1 class="font-bold text-2xl text-black">
          {{ config?.current_grade }} {{ config?.current_semester }} {{ authStore.user?.nickname || authStore.user?.username }}的课表
        </h1>
      </div>

      <h1 class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 font-bold text-base text-[#333333] print:hidden">
        {{ authStore.user?.nickname || authStore.user?.username }}的课表
      </h1>

      <div class="z-10 flex items-center gap-4 print:hidden">
        <el-icon :size="20" class="cursor-pointer text-gray-600 dark:text-gray-300 hover:text-primary" @click="handlePrint"><Printer /></el-icon>
        <el-icon :size="20" class="cursor-pointer text-gray-600 dark:text-gray-300 hover:text-primary" @click="router.push('/timetable/settings')"><Setting /></el-icon>
      </div>
    </div>

    <!-- Timetable -->
    <div class="p-2 overflow-x-auto relative z-10 print:p-0 print:overflow-visible">
      <div class="min-w-[320px] print:w-full">
        <!-- Week Header -->
        <div class="grid gap-1 mb-1 print:gap-0 print:mb-0 print:border-b print:border-black" :style="{ gridTemplateColumns: gridColumns }">
          <div class="w-16 print:border-r print:border-black"></div> <!-- Period Column Placeholder -->
          <div v-for="day in days" :key="day.value" class="text-center py-1 font-semibold text-sm bg-blue-50 dark:bg-gray-700 rounded text-gray-700 dark:text-gray-200 print:bg-transparent print:text-black print:border-r print:border-black print:rounded-none print:py-2 print:text-base">
            {{ day.label }}
          </div>
        </div>

        <!-- Periods -->
        <div v-for="period in periods" :key="period" class="grid gap-1 mb-1 print:gap-0 print:mb-0 print:border-b print:border-black" :style="{ gridTemplateColumns: gridColumns }">
          <!-- Period Number & Time -->
          <div class="flex flex-col items-center justify-center text-xs text-gray-500 font-medium w-16 bg-gray-100 dark:bg-gray-800 rounded py-1 print:bg-transparent print:text-black print:border-r print:border-black print:rounded-none">
            <span class="text-sm font-bold print:text-base">{{ period }}</span>
            <div v-if="getPeriodTime(period)" class="flex flex-col items-center text-[10px] scale-90 text-gray-400 leading-tight mt-0.5 print:text-black print:scale-100 print:text-xs">
                <span>{{ getPeriodTime(period)?.start }}</span>
                <span>{{ getPeriodTime(period)?.end }}</span>
            </div>
          </div>
          <!-- Course Cells -->
          <div 
            v-for="day in days" 
            :key="`${day.value}-${period}`"
            class="h-16 rounded p-1 flex items-center justify-center text-center text-xs font-bold shadow-sm transition-colors relative overflow-hidden bg-white dark:bg-gray-800 print:shadow-none print:rounded-none print:border-r print:border-black print:h-20 print:text-base"
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
  
  return {
    backgroundColor: item.course.color,
    color: '#000',
    // In print mode, we want exact colors
    printColorAdjust: 'exact',
    WebkitPrintColorAdjust: 'exact'
  } as any
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