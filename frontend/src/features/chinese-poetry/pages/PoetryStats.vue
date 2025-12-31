<template>
  <div class="min-h-screen bg-[#FDF6F8] dark:bg-gray-900 pb-20 font-sans">
    <!-- Header -->
    <div class="sticky top-0 z-20 bg-white/80 dark:bg-gray-800/80 backdrop-blur-md border-b border-pink-100 dark:border-gray-700 shadow-sm">
      <div class="flex items-center gap-3 p-4">
        <el-icon :size="22" class="cursor-pointer text-gray-600 dark:text-gray-300 hover:text-pink-500 transition" @click="router.back()"><ArrowLeft /></el-icon>
        <h1 class="text-xl font-bold text-gray-800 dark:text-white">学习统计</h1>
      </div>
    </div>

    <div class="p-6 space-y-6 max-w-4xl mx-auto">
      <!-- Summary Cards -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div class="bg-white dark:bg-gray-800 p-4 rounded-3xl shadow-sm border border-pink-50 dark:border-gray-700 text-center">
           <div class="text-3xl font-bold text-pink-500 mb-1">{{ totalStudied }}</div>
           <div class="text-xs text-gray-500 dark:text-gray-400">累计学习(首)</div>
        </div>
        <div class="bg-white dark:bg-gray-800 p-4 rounded-3xl shadow-sm border border-pink-50 dark:border-gray-700 text-center">
           <div class="text-3xl font-bold text-green-500 mb-1">{{ masteredCount }}</div>
           <div class="text-xs text-gray-500 dark:text-gray-400">已掌握</div>
        </div>
        <div class="bg-white dark:bg-gray-800 p-4 rounded-3xl shadow-sm border border-pink-50 dark:border-gray-700 text-center">
           <div class="text-3xl font-bold text-blue-500 mb-1">{{ totalDictations }}</div>
           <div class="text-xs text-gray-500 dark:text-gray-400">默写次数</div>
        </div>
        <div class="bg-white dark:bg-gray-800 p-4 rounded-3xl shadow-sm border border-pink-50 dark:border-gray-700 text-center">
           <div class="text-3xl font-bold text-purple-500 mb-1">{{ averageScore }}</div>
           <div class="text-xs text-gray-500 dark:text-gray-400">平均分</div>
        </div>
      </div>

      <!-- Charts Row 1 -->
      <div class="grid md:grid-cols-2 gap-6">
        <!-- Study Trend -->
        <div class="bg-white dark:bg-gray-800 p-6 rounded-3xl shadow-lg shadow-pink-100/50 dark:shadow-none border border-pink-50 dark:border-gray-700">
           <h3 class="font-bold text-gray-800 dark:text-white mb-4 flex items-center gap-2">
             <el-icon class="text-pink-500"><TrendCharts /></el-icon>
             学习趋势 (近7天)
           </h3>
           <div class="h-64">
             <v-chart :option="trendOption" autoresize />
           </div>
        </div>

        <!-- Mastery Distribution -->
        <div class="bg-white dark:bg-gray-800 p-6 rounded-3xl shadow-lg shadow-pink-100/50 dark:shadow-none border border-pink-50 dark:border-gray-700">
           <h3 class="font-bold text-gray-800 dark:text-white mb-4 flex items-center gap-2">
             <el-icon class="text-green-500"><PieChart /></el-icon>
             掌握程度
           </h3>
           <div class="h-64">
             <v-chart :option="masteryOption" autoresize />
           </div>
        </div>
      </div>

       <!-- Memory Stages -->
       <div class="bg-white dark:bg-gray-800 p-6 rounded-3xl shadow-lg shadow-pink-100/50 dark:shadow-none border border-pink-50 dark:border-gray-700">
          <h3 class="font-bold text-gray-800 dark:text-white mb-4 flex items-center gap-2">
            <el-icon class="text-purple-500"><Histogram /></el-icon>
            记忆阶段分布 (艾宾浩斯)
          </h3>
          <div class="h-64">
            <v-chart :option="stageOption" autoresize />
          </div>
       </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, TrendCharts, PieChart, Histogram } from '@element-plus/icons-vue'
import { usePoetryStore } from '../stores/poetryStore'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, PieChart as EPieChart, BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import dayjs from 'dayjs'

use([CanvasRenderer, LineChart, EPieChart, BarChart, GridComponent, TooltipComponent, LegendComponent])

const router = useRouter()
const store = usePoetryStore()

onMounted(() => {
  store.init()
})

// Stats logic
const totalStudied = computed(() => Object.keys(store.studyRecords).length)
const masteredCount = computed(() => Object.values(store.studyRecords).filter(r => r.isMastered).length)
const totalDictations = computed(() => store.dictationRecords.length)
const averageScore = computed(() => {
  if (store.dictationRecords.length === 0) return 0
  const sum = store.dictationRecords.reduce((acc, cur) => acc + cur.accuracy, 0)
  return Math.round(sum / store.dictationRecords.length)
})

// Trend Chart
const trendOption = computed(() => {
  const days = []
  const data = []
  for (let i = 6; i >= 0; i--) {
    const d = dayjs().subtract(i, 'day').format('YYYY-MM-DD')
    days.push(dayjs(d).format('MM-DD'))
    data.push(store.dailyStats[d]?.studiedCount || 0)
  }

  return {
    tooltip: { trigger: 'axis' },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: days },
    yAxis: { type: 'value', minInterval: 1 },
    series: [{
      data: data,
      type: 'line',
      smooth: true,
      itemStyle: { color: '#ec4899' },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [{ offset: 0, color: '#fbcfe8' }, { offset: 1, color: '#fff' }]
        }
      }
    }]
  }
})

// Mastery Chart
const masteryOption = computed(() => {
  const mastered = masteredCount.value
  const learning = totalStudied.value - mastered
  const notStarted = Math.max(0, store.poems.length - totalStudied.value)
  
  return {
    tooltip: { trigger: 'item' },
    legend: { bottom: '0%' },
    series: [
      {
        name: '状态',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
        label: { show: false },
        data: [
          { value: mastered, name: '已掌握', itemStyle: { color: '#10b981' } },
          { value: learning, name: '学习中', itemStyle: { color: '#3b82f6' } },
          { value: notStarted, name: '未开始', itemStyle: { color: '#e5e7eb' } }
        ]
      }
    ]
  }
})

// Stage Chart
const stageOption = computed(() => {
  const counts = [0, 0, 0, 0, 0, 0]
  Object.values(store.studyRecords).forEach(r => {
    if (r.stage < 6) counts[r.stage]++
  })

  return {
    tooltip: { trigger: 'axis' },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { 
        type: 'category', 
        data: ['阶段1', '阶段2', '阶段3', '阶段4', '阶段5', '阶段6'] 
    },
    yAxis: { type: 'value', minInterval: 1 },
    series: [{
        data: counts,
        type: 'bar',
        itemStyle: { borderRadius: [5, 5, 0, 0], color: '#8b5cf6' },
        barWidth: '40%'
    }]
  }
})

</script>
