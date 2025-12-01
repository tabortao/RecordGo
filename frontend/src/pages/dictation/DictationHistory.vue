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

onMounted(async () => {
  try {
    const res = await http.get('/dictation/history')
    historyList.value = (res as any).data || (res as any) || []
  } catch (e) {
    console.error(e)
  }
})

const chartOption = computed(() => {
  // Show last 10 records for better visibility
  const data = [...historyList.value].reverse().slice(-10) 
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
    <div class="flex items-center gap-3 mb-6">
      <el-button circle :icon="ArrowLeft" @click="router.back()" />
      <h1 class="text-xl font-bold">听写记录</h1>
    </div>

    <div v-if="historyList.length > 0" class="bg-white rounded-xl p-4 shadow-sm mb-4 h-72">
      <v-chart class="w-full h-full" :option="chartOption" autoresize />
    </div>
    <el-empty v-else description="暂无听写记录" />

    <div class="space-y-3">
      <div v-for="item in historyList" :key="item.id" class="bg-white rounded-xl p-4 shadow-sm border-l-4" :class="item.mistake_count > 0 ? 'border-red-400' : 'border-green-400'">
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
