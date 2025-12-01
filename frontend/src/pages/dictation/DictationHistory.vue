<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import http from '@/services/http'
import dayjs from 'dayjs'
import { ArrowLeft } from '@element-plus/icons-vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, TitleComponent, LegendComponent } from 'echarts/components'
import VChart from 'vue-echarts'

use([CanvasRenderer, LineChart, BarChart, GridComponent, TooltipComponent, TitleComponent, LegendComponent])

const router = useRouter()
const historyList = ref<any[]>([])
// Default to today
const today = new Date()
const dateRange = ref<[Date, Date]>([today, today])

const allHistoryList = ref<any[]>([])

const filteredList = computed(() => {
  if (!dateRange.value) return allHistoryList.value
  const [start, end] = dateRange.value
  // Start of day for start date, End of day for end date
  const startDate = dayjs(start).startOf('day')
  const endDate = dayjs(end).endOf('day')
  
  return allHistoryList.value.filter(item => {
    const d = dayjs(item.created_at)
    return d.isAfter(startDate) && d.isBefore(endDate)
  })
})

const totalStats = computed(() => {
  // Calculate stats based on ALL history, not filtered
  const list = allHistoryList.value
  const count = list.length
  const duration = list.reduce((acc, cur) => acc + (cur.duration_seconds || 0), 0)
  
  // Format duration: show hours and minutes if > 60 min, or just minutes
  let durationText = ''
  if (duration > 3600) {
      const h = Math.floor(duration / 3600)
      const m = Math.floor((duration % 3600) / 60)
      durationText = `${h}小时${m}分`
  } else {
      durationText = `${Math.ceil(duration / 60)}分`
  }

  return {
    count,
    durationText
  }
})

onMounted(async () => {
  try {
    const res = await http.get('/dictation/history')
    allHistoryList.value = (res as any).data || (res as any) || []
    // Sort by date descending
    allHistoryList.value.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
  } catch (e) {
    console.error(e)
  }
})

const chartOption = computed(() => {
  // Show last 10 records for better visibility
  const data = [...filteredList.value].reverse().slice(-10) 
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['时长(秒)', '错题数'] },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'category',
      boundaryGap: true,
      data: data.map(item => dayjs(item.created_at).format('MM-DD HH:mm')),
      axisLabel: { rotate: 45, interval: 0, fontSize: 10 }
    },
    yAxis: [
      { type: 'value', name: '时长', position: 'left', splitLine: { show: false } },
      { type: 'value', name: '错题', position: 'right', minInterval: 1 }
    ],
    series: [
      {
        name: '时长(秒)',
        type: 'bar',
        data: data.map(item => item.duration_seconds),
        yAxisIndex: 0,
        itemStyle: { color: '#3b82f6' }
      },
      {
        name: '错题数',
        type: 'line',
        data: data.map(item => item.mistake_count),
        yAxisIndex: 1,
        itemStyle: { color: '#ef4444' }
      }
    ]
  }
})

function formatDuration(sec: number) {
  const m = Math.floor(sec / 60)
  const s = sec % 60
  return `${m}分${s}秒`
}
</script>

<template>
  <div class="min-h-screen bg-gray-50 p-4">
    <div class="flex items-center gap-3 mb-4">
      <el-button circle :icon="ArrowLeft" @click="router.back()" />
      <h1 class="text-xl font-bold">听写记录</h1>
    </div>

    <!-- Stats Summary -->
    <div class="grid grid-cols-2 gap-3 mb-4">
       <div class="bg-white rounded-xl p-3 shadow-sm flex flex-col items-center justify-center">
          <div class="text-2xl font-bold text-indigo-600">{{ totalStats.count }}</div>
          <div class="text-xs text-gray-500">总练习次数</div>
       </div>
       <div class="bg-white rounded-xl p-3 shadow-sm flex flex-col items-center justify-center">
          <div class="text-2xl font-bold text-purple-600">{{ totalStats.durationText }}</div>
          <div class="text-xs text-gray-500">总时长</div>
       </div>
    </div>

    <!-- Date Filter -->
    <div class="mb-4">
       <el-date-picker
        v-model="dateRange"
        type="daterange"
        range-separator="至"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        style="width: 100%"
      />
    </div>

    <div v-if="filteredList.length > 0" class="bg-white rounded-xl p-4 shadow-sm mb-4 h-72">
      <v-chart class="w-full h-full" :option="chartOption" autoresize />
    </div>
    <el-empty v-else description="暂无听写记录" />

    <div class="space-y-3">
      <div v-for="item in filteredList" :key="item.id" class="bg-white rounded-xl p-4 shadow-sm border-l-4" :class="item.mistake_count > 0 ? 'border-red-400' : 'border-green-400'">
        <div class="flex justify-between items-start mb-2">
          <div class="text-sm text-gray-500">{{ dayjs(item.created_at).format('YYYY-MM-DD HH:mm:ss') }}</div>
          <div class="font-bold text-blue-600">{{ formatDuration(item.duration_seconds) }}</div>
        </div>
        <div class="text-gray-800 mb-2 text-sm break-all line-clamp-3">{{ item.content_snapshot }}</div>
        <div class="flex items-center gap-2 text-xs">
          <span v-if="item.mistake_count > 0" class="text-red-500 font-bold">错题: {{ item.mistake_count }}</span>
          <span v-else class="text-green-500 font-bold">全对 💯</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Chart styles */
</style>
