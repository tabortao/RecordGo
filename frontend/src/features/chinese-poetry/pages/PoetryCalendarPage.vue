<template>
  <div class="min-h-screen bg-[#FDF6F8] dark:bg-gray-900 pb-10 font-sans">
    <!-- Header -->
    <div class="sticky top-0 z-20 bg-white/80 dark:bg-gray-800/80 backdrop-blur-md border-b border-pink-100 dark:border-gray-700 flex items-center justify-between p-4 shadow-sm">
      <div class="flex items-center gap-3">
        <el-icon :size="22" class="cursor-pointer text-gray-600 dark:text-gray-300 hover:text-pink-500 transition" @click="router.back()"><ArrowLeft /></el-icon>
        <h1 class="text-xl font-bold text-gray-800 dark:text-white">古诗词</h1>
      </div>
      <div class="flex items-center gap-3">
         <div 
            class="flex items-center gap-1 bg-pink-100 dark:bg-pink-900/30 text-pink-600 dark:text-pink-300 px-3 py-1.5 rounded-full cursor-pointer hover:bg-pink-200 transition"
            @click="router.push('/homework/chinese/poetry/library')"
         >
            <el-icon><Collection /></el-icon>
            <span class="text-sm font-medium">诗库</span>
         </div>
         <div 
            class="flex items-center gap-1 bg-purple-100 dark:bg-purple-900/30 text-purple-600 dark:text-purple-300 px-3 py-1.5 rounded-full cursor-pointer hover:bg-purple-200 transition"
            @click="router.push('/homework/chinese/poetry/stats')"
         >
            <el-icon><TrendCharts /></el-icon>
         </div>
         <div 
            class="flex items-center gap-1 bg-orange-100 dark:bg-orange-900/30 text-orange-600 dark:text-orange-300 px-3 py-1.5 rounded-full cursor-pointer hover:bg-orange-200 transition"
            @click="router.push('/homework/chinese/poetry/challenge')"
         >
            <el-icon><Trophy /></el-icon>
         </div>
      </div>
    </div>

    <!-- Week Calendar -->
    <div class="p-4">
      <div class="bg-white dark:bg-gray-800 rounded-3xl shadow-sm border border-pink-50 dark:border-gray-700 overflow-hidden p-4">
        <!-- Week Header -->
        <div class="flex justify-between items-center mb-4 px-2">
            <h2 class="text-lg font-bold text-gray-800 dark:text-white">{{ currentYearMonth }}</h2>
            <div class="flex gap-2">
                <el-button circle size="small" @click="prevWeek"><el-icon><ArrowLeft /></el-icon></el-button>
                <el-button circle size="small" @click="nextWeek"><el-icon><ArrowRight /></el-icon></el-button>
            </div>
        </div>
        <!-- Days -->
        <div class="grid grid-cols-7 gap-1">
            <div v-for="day in weekDays" :key="day.dateStr" 
                 class="flex flex-col items-center justify-center py-2 rounded-xl cursor-pointer transition-all"
                 :class="isSelectedDate(day.dateStr) ? 'bg-pink-500 shadow-md transform scale-105' : 'hover:bg-gray-50 dark:hover:bg-gray-700'"
                 @click="selectDate(day.dateStr)"
            >
                <span class="text-xs mb-1 opacity-80" :class="isSelectedDate(day.dateStr) ? 'text-white' : 'text-gray-400'">{{ day.dayName }}</span>
                <span class="text-lg font-bold" :class="isSelectedDate(day.dateStr) ? 'text-white' : 'text-gray-800 dark:text-white'">{{ day.dayNum }}</span>
                <div class="h-1.5 flex gap-0.5 mt-1">
                     <div v-if="hasTasks(day.dateStr)" class="w-1.5 h-1.5 rounded-full" :class="isSelectedDate(day.dateStr) ? 'bg-white' : 'bg-blue-400'"></div>
                     <div v-if="isCompleted(day.dateStr)" class="w-1.5 h-1.5 rounded-full" :class="isSelectedDate(day.dateStr) ? 'bg-yellow-300' : 'bg-green-400'"></div>
                </div>
            </div>
        </div>
      </div>
    </div>

    <!-- Daily Tasks -->
    <div class="px-4 space-y-4">
       <div class="flex items-center justify-between">
          <h2 class="text-lg font-bold text-gray-800 dark:text-white flex items-center gap-2">
             <el-icon class="text-pink-500"><Calendar /></el-icon>
             {{ formatDate(currentDate) }} 学习计划
          </h2>
          <div class="text-sm text-gray-500">
             完成度: <span class="font-bold text-green-500">{{ completionRate }}%</span>
          </div>
       </div>

       <div v-if="dailyPoems.length > 0" class="space-y-4">
          <div 
             v-for="item in dailyPoems" 
             :key="item.poem.id"
             class="rounded-2xl p-5 shadow-sm border border-gray-100 dark:border-gray-700 hover:shadow-md transition relative overflow-hidden cursor-pointer"
             :class="item.isMastered ? 'bg-green-50 dark:bg-green-900/20' : 'bg-white dark:bg-gray-800'"
             @click="router.push(`/homework/chinese/poetry/${item.poem.id}`)"
          >
             <!-- Status Badge (Top Right) -->
             <div class="absolute top-0 right-0 px-3 py-1 rounded-bl-xl text-xs font-bold shadow-sm"
                :class="item.isMastered ? 'bg-green-100 text-green-600 dark:bg-green-900/30 dark:text-green-400' : 'bg-blue-50 text-blue-500 dark:bg-blue-900/30 dark:text-blue-400'"
             >
                {{ item.isMastered ? '已背' : '学习中' }}
             </div>

             <div class="flex justify-between items-start pr-12">
                <div>
                   <h3 class="text-lg font-bold text-gray-800 dark:text-gray-100 mb-1">{{ item.poem.title_cns }}</h3>
                   <div class="text-sm text-gray-500 dark:text-gray-400">{{ item.poem.dynasty_cns }} · {{ item.poem.author_cns }}</div>
                </div>
                <div class="text-center">
                   <div class="text-xl font-bold text-pink-500">{{ item.timesStudied }}</div>
                   <div class="text-xs text-gray-400">学习次数</div>
                </div>
             </div>

             <div class="mt-4 flex items-center justify-between">
                <div class="flex gap-2">
                   <span class="px-2 py-0.5 rounded-md text-xs bg-gray-100 dark:bg-gray-700 text-gray-500">{{ item.type === 'review' ? '复习' : '新学' }}</span>
                </div>
                <el-button size="small" round :type="item.isTodayCompleted ? 'success' : 'primary'" class="!px-4" @click.stop="router.push(`/homework/chinese/poetry/${item.poem.id}`)">
                   {{ item.isTodayCompleted ? '今日已学' : '去学习' }}
                </el-button>
             </div>
          </div>
       </div>

       <div v-else class="text-center py-12 text-gray-400 bg-white dark:bg-gray-800 rounded-3xl border border-dashed border-gray-200 dark:border-gray-700">
          <el-icon :size="48" class="mb-2 opacity-20"><CoffeeCup /></el-icon>
          <p>今日暂无学习任务，去<span class="text-pink-500 cursor-pointer font-bold" @click="router.push('/homework/chinese/poetry/library')">诗库</span>看看吧</p>
       </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, ArrowRight, Collection, Calendar, CoffeeCup, TrendCharts, Trophy } from '@element-plus/icons-vue'
import { usePoetryStore } from '@/features/chinese-poetry/stores/poetryStore'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'

dayjs.locale('zh-cn')

const router = useRouter()
const poetryStore = usePoetryStore()
const currentDate = ref(new Date())
const weekStart = ref(dayjs().startOf('week'))

onMounted(() => {
   poetryStore.init()
})

const currentYearMonth = computed(() => {
    return weekStart.value.format('YYYY年MM月')
})

const weekDays = computed(() => {
    const days = []
    for (let i = 0; i < 7; i++) {
        const d = weekStart.value.add(i, 'day')
        days.push({
            dateStr: d.format('YYYY-MM-DD'),
            dayNum: d.date(),
            dayName: d.format('dd')
        })
    }
    return days
})

const prevWeek = () => {
    weekStart.value = weekStart.value.subtract(1, 'week')
}

const nextWeek = () => {
    weekStart.value = weekStart.value.add(1, 'week')
}

const selectDate = (dateStr: string) => {
    currentDate.value = new Date(dateStr)
}

const isSelectedDate = (dateStr: string) => {
   return dayjs(dateStr).isSame(currentDate.value, 'day')
}

const formatDate = (date: Date) => {
   return dayjs(date).format('MM月DD日')
}

const hasTasks = (dateStr: string) => {
   return poetryStore.dailyStats[dateStr] || dayjs(dateStr).isSame(new Date(), 'day')
}

const isCompleted = (dateStr: string) => {
   const stats = poetryStore.dailyStats[dateStr]
   return stats && stats.studiedCount >= 3
}

const dailyPoems = computed(() => {
   const dateStr = dayjs(currentDate.value).format('YYYY-MM-DD')
   const stats = poetryStore.dailyStats[dateStr]
   
   if (stats) {
      return stats.details.map(id => {
         const poem = poetryStore.getPoemById(id)
         const record = poetryStore.getRecord(id)
         if (!poem) return null
         return {
            poem,
            type: 'review',
            timesStudied: record ? record.timesStudied : 0,
            isMastered: record ? record.isMastered : false,
            isTodayCompleted: true
         }
      }).filter(Boolean) as any[]
   } else if (dayjs(currentDate.value).isSame(new Date(), 'day')) {
       const now = Date.now()
       const reviews = Object.values(poetryStore.studyRecords)
          .filter(r => r.nextReviewTime <= now && !r.isMastered)
          .map(r => {
             const poem = poetryStore.getPoemById(r.poemId)
             if (!poem) return null
             return {
                poem,
                type: 'review',
                timesStudied: r.timesStudied,
                isMastered: r.isMastered,
                isTodayCompleted: false
             }
          }).filter(Boolean) as any[]
        
       if (reviews.length === 0) {
           const newPoems = poetryStore.poems.slice(0, 3).map(p => ({
               poem: p,
               type: 'new',
               timesStudied: 0,
               isMastered: false,
               isTodayCompleted: false
           }))
           return newPoems
       }
       return reviews
   }
   
   return []
})

const completionRate = computed(() => {
    if (dailyPoems.value.length === 0) return 0
    const completed = dailyPoems.value.filter(p => p.isTodayCompleted).length
    return Math.round((completed / dailyPoems.value.length) * 100)
})
</script>
