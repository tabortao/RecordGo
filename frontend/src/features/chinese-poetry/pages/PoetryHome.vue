<template>
  <div class="min-h-screen bg-[#FDF6F8] dark:bg-gray-900 pb-20 font-sans">
    <!-- Header with Search -->
    <div class="sticky top-0 z-20 bg-white/80 dark:bg-gray-800/80 backdrop-blur-md border-b border-pink-100 dark:border-gray-700 shadow-sm">
      <div class="flex items-center gap-3 p-4">
        <el-icon :size="22" class="cursor-pointer text-gray-600 dark:text-gray-300 hover:text-pink-500 transition" @click="router.back()"><ArrowLeft /></el-icon>
        <div class="flex-1 bg-pink-50 dark:bg-gray-700 rounded-full flex items-center px-4 py-2 transition focus-within:ring-2 ring-pink-200">
          <el-icon class="text-pink-300 mr-2"><Search /></el-icon>
          <input 
            v-model="searchQuery"
            class="bg-transparent border-none outline-none text-sm w-full text-gray-700 dark:text-gray-200 placeholder-pink-300"
            placeholder="搜索诗名、作者、诗句..."
          />
        </div>
        <div class="flex gap-2">
           <el-icon :size="24" class="cursor-pointer text-pink-500 hover:scale-110 transition" @click="router.push('/homework/chinese/poetry/create')"><Plus /></el-icon>
           <el-icon :size="24" class="cursor-pointer text-purple-500 hover:scale-110 transition" @click="router.push('/homework/chinese/poetry/challenge')"><Trophy /></el-icon>
           <el-icon :size="24" class="cursor-pointer text-blue-500 hover:scale-110 transition" @click="router.push('/homework/chinese/poetry/stats')"><DataAnalysis /></el-icon>
        </div>
      </div>
      
      <!-- Tabs -->
      <div class="px-2 pb-3 overflow-x-auto scrollbar-hide">
        <div class="flex gap-3 px-2">
          <button 
            v-for="tab in tabs" 
            :key="tab.value"
            @click="currentTab = tab.value"
            class="px-4 py-1.5 rounded-full text-xs font-bold transition-all whitespace-nowrap border-2 shadow-sm"
            :class="currentTab === tab.value 
              ? 'bg-pink-400 border-pink-400 text-white shadow-pink-200' 
              : 'bg-white border-pink-100 text-gray-500 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-300 hover:border-pink-300'"
          >
            {{ tab.label }}
          </button>
        </div>
      </div>
    </div>

    <!-- Main Content Area -->
    <div v-if="!searchQuery && currentTab === 'all'" class="p-4 space-y-6">
      
      <!-- Daily Recommendation Card -->
      <div v-if="dailyRecommendPoem" class="relative overflow-hidden bg-gradient-to-br from-pink-100 to-purple-100 dark:from-pink-900/30 dark:to-purple-900/30 rounded-2xl p-5 shadow-sm border border-pink-200 dark:border-pink-800 cursor-pointer group" @click="goToDetail(dailyRecommendPoem)">
          <div class="absolute -right-4 -top-4 w-24 h-24 bg-white/30 rounded-full blur-xl group-hover:scale-150 transition duration-700"></div>
          <div class="relative z-10">
              <div class="flex justify-between items-start mb-3">
                  <span class="bg-white/60 dark:bg-black/20 text-pink-600 dark:text-pink-300 text-[10px] font-bold px-2 py-1 rounded-lg backdrop-blur">每日一诗</span>
                  <el-icon class="text-pink-400"><StarFilled /></el-icon>
              </div>
              <h2 class="text-2xl font-bold text-gray-800 dark:text-white mb-1 font-serif">{{ dailyRecommendPoem.title_cns }}</h2>
              <p class="text-sm text-gray-600 dark:text-gray-300 mb-4">{{ dailyRecommendPoem.author_cns }} · {{ dailyRecommendPoem.dynasty_cns }}</p>
              <div class="text-base text-gray-700 dark:text-gray-200 font-serif leading-relaxed line-clamp-2 opacity-80">
                  {{ dailyRecommendPoem.paragraphs_cns.slice(0, 2).join('，') }}...
              </div>
          </div>
      </div>

      <!-- Dashboard Stats -->
      <div class="grid grid-cols-2 gap-4">
        <div class="bg-[#E0F7FA] dark:bg-cyan-900/30 p-4 rounded-2xl border border-[#B2EBF2] dark:border-cyan-800 relative overflow-hidden">
            <div class="absolute -right-2 -bottom-2 text-[#B2EBF2] dark:text-cyan-800/50"><el-icon size="60"><Timer /></el-icon></div>
            <div class="relative z-10">
                <div class="text-xs text-cyan-600 dark:text-cyan-300 font-bold mb-1">今日需学</div>
                <div class="text-3xl font-black text-cyan-700 dark:text-cyan-200">{{ todayStats.studiedCount + dailyReviewList.length }} <span class="text-xs font-normal">首</span></div>
            </div>
        </div>
        <div class="bg-[#F3E5F5] dark:bg-purple-900/30 p-4 rounded-2xl border border-[#E1BEE7] dark:border-purple-800 relative overflow-hidden">
            <div class="absolute -right-2 -bottom-2 text-[#E1BEE7] dark:text-purple-800/50"><el-icon size="60"><Calendar /></el-icon></div>
            <div class="relative z-10">
                <div class="text-xs text-purple-600 dark:text-purple-300 font-bold mb-1">累计打卡</div>
                <div class="text-3xl font-black text-purple-700 dark:text-purple-200">{{ totalDays }} <span class="text-xs font-normal">天</span></div>
            </div>
        </div>
      </div>

      <!-- Ebbinghaus Plan & Calendar -->
      <div class="bg-white dark:bg-gray-800 rounded-2xl p-4 shadow-sm border border-gray-100 dark:border-gray-700">
         <div class="flex justify-between items-center mb-4">
             <div class="flex items-center gap-2">
                <div class="w-1 h-4 bg-pink-400 rounded-full"></div>
                <span class="font-bold text-gray-700 dark:text-gray-200">学习日历</span>
             </div>
             <span class="text-xs text-gray-400 bg-gray-50 dark:bg-gray-700 px-2 py-1 rounded-lg">{{ currentMonth }}</span>
         </div>
         <div class="flex justify-between mb-4">
             <div v-for="day in weekDays" :key="day.date" 
                  class="flex flex-col items-center gap-2 cursor-pointer group"
                  :class="{'opacity-40': !day.isCurrentMonth}"
                  @click="showDayDetail(day.date)"
             >
                 <span class="text-[10px] text-gray-400 font-medium">{{ day.weekDay }}</span>
                 <div 
                   class="w-9 h-9 rounded-xl flex items-center justify-center text-xs font-bold transition-all duration-300 relative"
                   :class="getDayStatusClass(day)"
                 >
                    {{ day.dayNum }}
                    <!-- Dot for tasks -->
                    <div v-if="day.hasTask" class="absolute -bottom-1 w-1 h-1 bg-pink-500 rounded-full"></div>
                 </div>
             </div>
         </div>
      </div>

      <!-- Today's Tasks Lists -->
      <div v-if="dailyReviewList.length > 0" class="animate-fade-in-up">
        <div class="flex items-center gap-2 mb-3 px-1">
            <span class="text-lg font-bold text-gray-800 dark:text-gray-100">今日复习</span>
            <span class="bg-orange-100 text-orange-600 text-[10px] px-2 py-0.5 rounded-full font-bold">{{ dailyReviewList.length }}</span>
        </div>
        <div class="space-y-3">
           <div 
            v-for="poem in dailyReviewList" 
            :key="poem.id"
            class="bg-white dark:bg-gray-800 rounded-xl p-4 shadow-sm border-l-4 border-orange-400 flex justify-between items-center cursor-pointer hover:shadow-md transition"
            @click="goToDetail(poem)"
          >
            <div>
               <div class="font-bold text-gray-800 dark:text-gray-200 mb-1">{{ poem.title_cns }}</div>
               <div class="text-xs text-gray-500 bg-gray-100 dark:bg-gray-700 inline-block px-1.5 py-0.5 rounded">{{ poem.author_cns }}</div>
            </div>
            <el-button type="warning" size="small" round color="#FB8C00" class="!text-white font-bold shadow-orange-200 shadow-md">复习</el-button>
          </div>
        </div>
      </div>
      
      <div v-if="dailyNewList.length > 0" class="animate-fade-in-up" :class="{'mt-6': dailyReviewList.length > 0}">
        <div class="flex items-center gap-2 mb-3 px-1">
            <span class="text-lg font-bold text-gray-800 dark:text-gray-100">今日新学</span>
            <span class="bg-blue-100 text-blue-600 text-[10px] px-2 py-0.5 rounded-full font-bold">{{ dailyNewList.length }}</span>
        </div>
        <div class="space-y-3">
           <div 
            v-for="poem in dailyNewList" 
            :key="poem.id"
            class="bg-white dark:bg-gray-800 rounded-xl p-4 shadow-sm border-l-4 border-blue-400 flex justify-between items-center cursor-pointer hover:shadow-md transition"
            @click="goToDetail(poem)"
          >
            <div>
               <div class="font-bold text-gray-800 dark:text-gray-200 mb-1">{{ poem.title_cns }}</div>
               <div class="text-xs text-gray-500 bg-gray-100 dark:bg-gray-700 inline-block px-1.5 py-0.5 rounded">{{ poem.author_cns }}</div>
            </div>
            <el-button type="primary" size="small" round color="#42A5F5" class="!text-white font-bold shadow-blue-200 shadow-md">学习</el-button>
          </div>
        </div>
      </div>
    </div>

    <!-- Filtered Poem List -->
    <div class="p-4 space-y-3" v-else>
      <div 
        v-for="poem in filteredPoems" 
        :key="poem.id"
        class="bg-white dark:bg-gray-800 rounded-xl p-4 shadow-sm border border-gray-100 dark:border-gray-700 cursor-pointer active:scale-98 transition relative overflow-hidden group hover:border-pink-200"
        @click="goToDetail(poem)"
      >
        <div v-if="isMastered(poem.id)" class="absolute top-0 right-0 bg-green-100 dark:bg-green-900/50 text-green-600 dark:text-green-400 text-[10px] px-2 py-1 rounded-bl-lg font-bold z-10">
            已掌握
        </div>
        
        <div class="flex justify-between items-start">
          <div class="flex-1">
            <h3 class="text-lg font-bold text-gray-800 dark:text-gray-100 mb-2 flex items-center gap-2">
                {{ poem.title_cns }}
                <el-icon v-if="poem.audio_url || poem.isCustom" class="text-pink-400 text-sm"><Headset /></el-icon>
            </h3>
            <div class="text-sm text-gray-500 dark:text-gray-400 mb-3 flex gap-2">
              <span class="bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded text-xs text-gray-600">{{ poem.dynasty_cns }}</span>
              <span class="font-medium">{{ poem.author_cns }}</span>
            </div>
            <div class="text-gray-600 dark:text-gray-300 text-sm line-clamp-2 leading-relaxed font-serif bg-gray-50 dark:bg-gray-900/50 p-2 rounded-lg">
               {{ poem.paragraphs_cns.join('，') }}
            </div>
          </div>
        </div>
      </div>
      
      <div v-if="filteredPoems.length === 0" class="flex flex-col items-center justify-center py-20 text-gray-400">
        <el-icon :size="60" class="mb-4 text-gray-200"><CollectionTag /></el-icon>
        <p>没有找到相关古诗</p>
        <el-button v-if="currentTab !== 'all'" link type="primary" @click="currentTab = 'all'">查看全部</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Headset, Search, Plus, Trophy, StarFilled, Timer, Calendar, CollectionTag, DataAnalysis } from '@element-plus/icons-vue'
import type { Poem } from '../types'
import { usePoetryStore } from '../stores/poetryStore'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'
import isBetween from 'dayjs/plugin/isBetween'
import isoWeek from 'dayjs/plugin/isoWeek'

dayjs.extend(isBetween)
dayjs.extend(isoWeek)
dayjs.locale('zh-cn')

const router = useRouter()
const store = usePoetryStore()

const searchQuery = ref('')
const currentTab = ref('all')

const tabs = [
  { label: '首页', value: 'all' },
  { label: '小学', value: 'primary' },
  { label: '初中', value: 'middle' },
  { label: '高中', value: 'high' },
  { label: '唐诗', value: 'tang' },
  { label: '宋词', value: 'song' },
  { label: '收藏', value: 'collection' }
]

onMounted(() => {
  store.init()
})

const dailyRecommendPoem = computed(() => store.dailyRecommendPoem)
const dailyReviewList = computed(() => store.dailyReviewList)
const dailyNewList = computed(() => store.dailyNewList)
const todayStats = computed(() => store.dailyStats[dayjs().format('YYYY-MM-DD')] || { studiedCount: 0, reviewCount: 0 })
const totalDays = computed(() => Object.keys(store.dailyStats).length)

// Calendar Logic
const today = dayjs()
const currentMonth = today.format('YYYY年M月')
const weekDays = computed(() => {
    // Start from 3 days ago to 3 days later, centered on today
    // Or just Monday to Sunday of current week?
    // Let's do Monday to Sunday for consistency with typical calendars
    const startOfWeek = today.startOf('isoWeek')
    
    const days = []
    for (let i = 0; i < 7; i++) {
        const date = startOfWeek.add(i, 'day')
        const dateStr = date.format('YYYY-MM-DD')
        const stats = store.dailyStats[dateStr]
        const hasTask = (stats && (stats.studiedCount > 0 || stats.reviewCount > 0)) || false
        
        days.push({
            date: dateStr,
            dayNum: date.date(),
            weekDay: date.format('dd'), // 周一, 周二...
            isCurrentMonth: date.month() === today.month(),
            isToday: date.isSame(today, 'day'),
            hasTask
        })
    }
    return days
})

const getDayStatusClass = (day: any) => {
    if (day.isToday) {
        return 'bg-pink-500 text-white shadow-lg shadow-pink-200'
    }
    if (day.hasTask) {
        return 'bg-pink-100 text-pink-600'
    }
    return 'bg-gray-50 text-gray-500 hover:bg-gray-100'
}

const showDayDetail = (date: string) => {
    // Show a modal or something? For now, just log or navigate
    // In PRD: "点击日期可查看每日学习情况"
    // Maybe filter the list below?
    // For simplicity, we just keep current view as "Today's view"
    if (date !== today.format('YYYY-MM-DD')) {
        // Maybe implement historical view later
    }
}

// Filtering
const filteredPoems = computed(() => {
    let poems = store.poems
    
    // 1. Tab Filter
    if (currentTab.value === 'collection') {
        poems = poems.filter((p: Poem) => store.collections.includes(p.id))
    } else if (currentTab.value === 'primary') {
        poems = poems.filter((p: Poem) => p.tags?.includes('primary') || (p.author_cns && ['李白','杜甫','白居易'].includes(p.author_cns))) // Mock logic
    } else if (currentTab.value === 'middle') {
        poems = poems.filter((p: Poem) => p.tags?.includes('middle') || (p.author_cns && ['苏轼','王维'].includes(p.author_cns)))
    } else if (currentTab.value === 'high') {
        poems = poems.filter((p: Poem) => p.tags?.includes('high') || (p.author_cns && ['辛弃疾','李清照'].includes(p.author_cns)))
    } else if (currentTab.value === 'tang') {
        poems = poems.filter((p: Poem) => p.dynasty_cns === '唐代')
    } else if (currentTab.value === 'song') {
        poems = poems.filter((p: Poem) => p.dynasty_cns === '宋代')
    }
    
    // 2. Search Filter
    if (searchQuery.value) {
        const q = searchQuery.value.toLowerCase()
        poems = poems.filter((p: Poem) => 
            p.title_cns.includes(q) || 
            p.author_cns.includes(q) || 
            p.paragraphs_cns.some((l: string) => l.includes(q))
        )
    }
    
    return poems
})

const isMastered = (id: number) => {
    const record = store.studyRecords[id]
    return record && record.nextReviewDate > Date.now() + 86400000 * 30 // Rough guess for mastery
}

const goToDetail = (poem: Poem) => {
    router.push(`/homework/chinese/poetry/${poem.id}`)
}
</script>

<style scoped>
.scrollbar-hide::-webkit-scrollbar {
    display: none;
}
.scrollbar-hide {
    -ms-overflow-style: none;
    scrollbar-width: none;
}
</style>