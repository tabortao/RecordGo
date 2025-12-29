<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 pb-20 relative overflow-hidden">
    <!-- Background Emojis -->
    <div v-if="emojiPositions.length > 0" class="fixed inset-0 pointer-events-none z-0 overflow-hidden opacity-20">
        <div v-for="(item, index) in emojiPositions" :key="index" 
             class="absolute text-6xl filter blur-[1px] select-none transform hover:scale-110 transition-transform duration-1000"
             :style="item.style">
            {{ item.emoji }}
        </div>
        <!-- Mask Overlay -->
        <div class="absolute inset-0 bg-white/40 dark:bg-gray-900/40 backdrop-blur-[1px]"></div>
    </div>

    <!-- Header -->
    <div class="bg-white dark:bg-gray-800 shadow-sm px-4 py-3 flex items-center justify-between sticky top-0 z-20 relative">
      <div class="flex items-center gap-3">
        <el-icon :size="20" class="cursor-pointer" @click="router.back()"><ArrowLeft /></el-icon>
        <span class="font-bold text-base text-[#333333] mr-2">一年级 上学期</span>
        <h1 class="font-normal text-lg" :style="{ color: 'var(--el-color-primary)' }">
          {{ authStore.user?.nickname || authStore.user?.username }}的课表
        </h1>
      </div>
      <el-icon :size="20" class="cursor-pointer text-gray-600 dark:text-gray-300" @click="router.push('/timetable/settings')"><Setting /></el-icon>
    </div>

    <!-- Timetable -->
    <div class="p-2 overflow-x-auto relative z-10">
      <div class="min-w-[320px]">
        <!-- Week Header -->
        <div class="grid gap-1 mb-1" :style="{ gridTemplateColumns: gridColumns }">
          <div class="w-16"></div> <!-- Period Column Placeholder -->
          <div v-for="day in days" :key="day.value" class="text-center py-1 font-semibold text-sm bg-blue-50 dark:bg-gray-700 rounded text-gray-700 dark:text-gray-200">
            {{ day.label }}
          </div>
        </div>

        <!-- Periods -->
        <div v-for="period in periods" :key="period" class="grid gap-1 mb-1" :style="{ gridTemplateColumns: gridColumns }">
          <!-- Period Number & Time -->
          <div class="flex flex-col items-center justify-center text-xs text-gray-500 font-medium w-16 bg-gray-100 dark:bg-gray-800 rounded py-1">
            <span class="text-sm font-bold">{{ period }}</span>
            <div v-if="getPeriodTime(period)" class="flex flex-col items-center text-[10px] scale-90 text-gray-400 leading-tight mt-0.5">
                <span>{{ getPeriodTime(period)?.start }}</span>
                <span>{{ getPeriodTime(period)?.end }}</span>
            </div>
          </div>
          <!-- Course Cells -->
          <div 
            v-for="day in days" 
            :key="`${day.value}-${period}`"
            class="h-16 rounded p-1 flex items-center justify-center text-center text-xs font-bold shadow-sm transition-colors relative overflow-hidden bg-white/50 dark:bg-gray-800/50 backdrop-blur-sm"
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
import { computed, onMounted, ref, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Setting } from '@element-plus/icons-vue'
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
const backgroundEmojis = ref<string[]>([])
const emojiPositions = ref<{ emoji: string, style: any }[]>([])

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

function getCellStyle(day: number, period: number) {
  const item = getCourse(day, period)
  if (!item || !item.course) return {} // Empty cell, use class bg-white/50
  
  return {
    backgroundColor: hexToRgba(item.course.color, 0.5),
    color: '#000' // Ensure text is dark for readability on light backgrounds
  }
}

function hexToRgba(hex: string, alpha: number) {
    let c: any;
    if(/^#([A-Fa-f0-9]{3}){1,2}$/.test(hex)){
        c= hex.substring(1).split('');
        if(c.length== 3){
            c= [c[0], c[0], c[1], c[1], c[2], c[2]];
        }
        c= '0x'+c.join('');
        return 'rgba('+[(c>>16)&255, (c>>8)&255, c&255].join(',')+','+alpha+')';
    }
    return hex; // Fallback
}

function generateEmojiPositions() {
    if (backgroundEmojis.value.length === 0) {
        emojiPositions.value = []
        return
    }

    const positions: { emoji: string, style: any, rect: {x: number, y: number, w: number, h: number} }[] = []
    const containerW = window.innerWidth
    const containerH = window.innerHeight
    const emojiSize = 60 // 估算大小

    // 尝试放置每个 emoji
    // 为了随机且不重叠，我们尝试多次放置
    const maxAttempts = 50

    // 随机选择 emoji 填充屏幕，数量可以根据屏幕大小决定，或者固定数量
    const count = 10 // 固定放置 10 个 emoji

    for (let i = 0; i < count; i++) {
        const emoji = backgroundEmojis.value[Math.floor(Math.random() * backgroundEmojis.value.length)]
        
        for (let attempt = 0; attempt < maxAttempts; attempt++) {
            const x = Math.random() * (containerW - emojiSize)
            const y = Math.random() * (containerH - emojiSize)
            
            // 检查重叠
            const overlap = positions.some(p => {
                return !(x + emojiSize < p.rect.x || 
                         x > p.rect.x + p.rect.w || 
                         y + emojiSize < p.rect.y || 
                         y > p.rect.y + p.rect.h)
            })

            if (!overlap) {
                positions.push({
                    emoji,
                    style: {
                        left: `${x}px`,
                        top: `${y}px`,
                        transform: `rotate(${Math.random() * 360}deg)`
                    },
                    rect: { x, y, w: emojiSize, h: emojiSize }
                })
                break
            }
        }
    }
    emojiPositions.value = positions
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
    
    // 解析背景 Emoji
    if (config.value.background_emojis) {
        backgroundEmojis.value = config.value.background_emojis.split(/[,，]/).map(s => s.trim()).filter(s => s)
        nextTick(() => {
            generateEmojiPositions()
        })
    }
    
    await store.fetchTimetable(config.value.current_grade, config.value.current_semester)
  }
})
</script>