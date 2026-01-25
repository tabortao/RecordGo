<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import http from '@/services/http'
import dayjs from 'dayjs'
import { DataLine } from '@element-plus/icons-vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, TitleComponent, LegendComponent } from 'echarts/components'
import VChart from 'vue-echarts'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

use([CanvasRenderer, LineChart, BarChart, GridComponent, TooltipComponent, TitleComponent, LegendComponent])

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
  <SettingsShell title="听写记录" subtitle="练习趋势与错题变化" :icon="DataLine" tone="violet" container-class="max-w-5xl">
    <SettingsCard>
      <div class="grid grid-cols-2 gap-3">
        <div class="rounded-3xl border border-indigo-100/70 dark:border-indigo-900/40 bg-indigo-50/70 dark:bg-indigo-900/20 px-4 py-3">
          <div class="text-[11px] font-extrabold uppercase tracking-wider text-indigo-700/80 dark:text-indigo-300/80">总练习次数</div>
          <div class="mt-1 text-2xl font-extrabold tracking-tight text-indigo-800 dark:text-indigo-200">{{ totalStats.count }}</div>
        </div>
        <div class="rounded-3xl border border-violet-100/70 dark:border-violet-900/40 bg-violet-50/70 dark:bg-violet-900/20 px-4 py-3">
          <div class="text-[11px] font-extrabold uppercase tracking-wider text-violet-700/80 dark:text-violet-300/80">总时长</div>
          <div class="mt-1 text-2xl font-extrabold tracking-tight text-violet-800 dark:text-violet-200">{{ totalStats.durationText }}</div>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard title="筛选日期">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <el-date-picker v-model="dateRange[0]" type="date" placeholder="开始日期" class="w-full" />
        <el-date-picker v-model="dateRange[1]" type="date" placeholder="结束日期" class="w-full" />
      </div>
    </SettingsCard>

    <SettingsCard v-if="filteredList.length > 0" title="趋势图" description="展示最近 10 条记录的时长与错题数。">
      <div class="h-72">
        <v-chart class="w-full h-full" :option="chartOption" autoresize />
      </div>
    </SettingsCard>
    <SettingsCard v-else>
      <el-empty description="暂无听写记录" />
    </SettingsCard>

    <SettingsCard title="记录列表">
      <div class="space-y-3">
        <div
          v-for="item in filteredList"
          :key="item.id"
          class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/55 backdrop-blur px-4 py-4 shadow-sm"
        >
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0">
              <div class="text-xs text-gray-500 dark:text-gray-400">{{ dayjs(item.created_at).format('YYYY-MM-DD HH:mm:ss') }}</div>
              <div class="mt-2 text-sm text-gray-800 dark:text-gray-100 break-all line-clamp-3">{{ item.content_snapshot }}</div>
            </div>
            <div class="shrink-0 text-right">
              <div class="text-sm font-extrabold text-blue-700 dark:text-blue-300">{{ formatDuration(item.duration_seconds) }}</div>
              <div class="mt-2 text-xs font-extrabold" :class="item.mistake_count > 0 ? 'text-red-600 dark:text-red-300' : 'text-emerald-600 dark:text-emerald-300'">
                <span v-if="item.mistake_count > 0">错题 {{ item.mistake_count }}</span>
                <span v-else>全对</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<style scoped>
/* Chart styles */
</style>
