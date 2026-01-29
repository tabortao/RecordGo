<template>
  <SettingsShell title="登录统计" subtitle="按天统计去重用户数（仅管理员可见）" :icon="TrendCharts" tone="emerald" container-class="max-w-6xl">
    <SettingsCard>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="flex items-center gap-2">
          <el-button class="!rounded-xl" :icon="ArrowLeft" @click="goBack">返回</el-button>
          <el-button class="!rounded-xl" :loading="loading" @click="load">刷新</el-button>
        </div>
        <el-radio-group v-model="rangeDays" size="small">
          <el-radio-button :label="7">近7天</el-radio-button>
          <el-radio-button :label="30">近30天</el-radio-button>
          <el-radio-button :label="365">近1年</el-radio-button>
        </el-radio-group>
      </div>

      <div class="mt-4 grid grid-cols-2 sm:grid-cols-4 gap-3">
        <div class="rounded-2xl border border-emerald-100/70 dark:border-gray-800 bg-emerald-50/70 dark:bg-gray-900/50 px-4 py-3">
          <div class="text-[11px] font-extrabold uppercase tracking-wider text-emerald-700/80 dark:text-emerald-300/80">今日登录</div>
          <div class="mt-1 text-2xl font-extrabold tracking-tight text-emerald-800 dark:text-emerald-200">{{ todayCount }}</div>
        </div>
        <div class="rounded-2xl border border-sky-100/70 dark:border-gray-800 bg-sky-50/70 dark:bg-gray-900/50 px-4 py-3">
          <div class="text-[11px] font-extrabold uppercase tracking-wider text-sky-700/80 dark:text-sky-300/80">区间累计</div>
          <div class="mt-1 text-2xl font-extrabold tracking-tight text-sky-800 dark:text-sky-200">{{ sumCount }}</div>
        </div>
        <div class="rounded-2xl border border-amber-100/70 dark:border-gray-800 bg-amber-50/70 dark:bg-gray-900/50 px-4 py-3">
          <div class="text-[11px] font-extrabold uppercase tracking-wider text-amber-700/80 dark:text-amber-300/80">日均</div>
          <div class="mt-1 text-2xl font-extrabold tracking-tight text-amber-800 dark:text-amber-200">{{ avgCount }}</div>
        </div>
        <div class="rounded-2xl border border-fuchsia-100/70 dark:border-gray-800 bg-fuchsia-50/70 dark:bg-gray-900/50 px-4 py-3">
          <div class="text-[11px] font-extrabold uppercase tracking-wider text-fuchsia-700/80 dark:text-fuchsia-300/80">峰值</div>
          <div class="mt-1 text-2xl font-extrabold tracking-tight text-fuchsia-800 dark:text-fuchsia-200">{{ maxCount }}</div>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard title="趋势图" :description="`统计范围：${stats.start_day || '--'} ~ ${stats.end_day || '--'}`">
      <div class="h-[360px] w-full">
        <VChart v-if="chartOption" :option="chartOption" autoresize class="w-full h-full" />
        <div v-else class="h-full flex items-center justify-center text-sm text-gray-400">暂无数据</div>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowLeft, TrendCharts } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import router from '@/router'
import http from '@/services/http'
import { useAuth } from '@/stores/auth'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'

use([LineChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer])

const auth = useAuth()
if (!auth.user || Number(auth.user.id) !== 1) {
  router.replace('/tasks')
}

const rangeDays = ref<7 | 30 | 365>(7)
const loading = ref(false)

const stats = reactive<{ range_days: number; start_day: string; end_day: string; items: Array<{ day: string; count: number }> }>({
  range_days: 7,
  start_day: '',
  end_day: '',
  items: []
})

const todayCount = computed(() => (stats.items.length ? Number(stats.items[stats.items.length - 1]?.count || 0) : 0))
const sumCount = computed(() => stats.items.reduce((sum, it) => sum + Number(it.count || 0), 0))
const avgCount = computed(() => {
  if (!stats.items.length) return 0
  return Math.round(sumCount.value / stats.items.length)
})
const maxCount = computed(() => (stats.items.length ? Math.max(...stats.items.map((it) => Number(it.count || 0))) : 0))

const chartOption = computed(() => {
  if (!stats.items.length) return null
  const x = stats.items.map((it) => {
    if (rangeDays.value === 365) return dayjs(it.day).format('MM-DD')
    return dayjs(it.day).format('MM/DD')
  })
  const y = stats.items.map((it) => Number(it.count || 0))
  return {
    tooltip: { trigger: 'axis' },
    grid: { left: 24, right: 16, top: 20, bottom: 28 },
    xAxis: {
      type: 'category',
      data: x,
      axisTick: { show: false },
      axisLabel: { color: '#94a3b8', interval: rangeDays.value === 365 ? 6 : 'auto' },
      axisLine: { lineStyle: { color: 'rgba(148,163,184,0.25)' } }
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      axisLabel: { color: '#94a3b8' },
      splitLine: { lineStyle: { color: 'rgba(148,163,184,0.14)' } }
    },
    series: [
      {
        name: '登录用户数',
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        lineStyle: { width: 3, color: '#10b981' },
        itemStyle: { color: '#10b981' },
        areaStyle: { opacity: 0.12, color: '#10b981' },
        data: y
      }
    ]
  }
})

function goBack() {
  router.push('/admin')
}

async function load() {
  if (loading.value) return
  loading.value = true
  try {
    const data = (await http.get('/admin/login-stats', { params: { range: rangeDays.value } })) as any
    stats.range_days = Number(data?.range_days || rangeDays.value)
    stats.start_day = String(data?.start_day || '')
    stats.end_day = String(data?.end_day || '')
    stats.items = Array.isArray(data?.items) ? data.items.map((it: any) => ({ day: String(it.day || ''), count: Number(it.count || 0) })) : []
  } catch (e: any) {
    ElMessage.error(e?.message || '加载失败')
  } finally {
    loading.value = false
  }
}

watch(rangeDays, async () => {
  await load()
})

onMounted(async () => {
  await load()
})
</script>

<style scoped></style>

