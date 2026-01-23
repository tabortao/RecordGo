<template>
  <div class="min-h-screen bg-[#F5F7FA] dark:bg-gray-900">
    <div class="bg-white/80 dark:bg-gray-900/80 backdrop-blur border-b border-gray-200 dark:border-gray-800 sticky top-0 z-10">
      <div class="px-4 py-3 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <el-icon :size="18" class="cursor-pointer text-gray-600 dark:text-gray-300" @click="router.back()"><ArrowLeft /></el-icon>
          <h2 class="font-semibold text-gray-800 dark:text-gray-100">{{ title }}</h2>
        </div>
        <el-button type="primary" size="small" @click="openDialog">添加记录</el-button>
      </div>
    </div>

    <div class="p-4 space-y-4">
      <div class="rounded-xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-3">
        <el-tabs v-model="chartTab" type="card" class="mb-2">
          <el-tab-pane label="本周" name="week" />
          <el-tab-pane label="本月" name="month" />
          <el-tab-pane label="本年" name="year" />
          <el-tab-pane label="全部" name="all" />
        </el-tabs>
        <VChart :option="chartOption" autoresize style="height: 200px;" />
      </div>

      <div class="flex items-center justify-between">
        <div class="text-sm font-semibold text-gray-700 dark:text-gray-200">记录的数据</div>
      </div>
      <div v-if="records.length === 0" class="text-sm text-gray-500 dark:text-gray-400">暂无记录</div>
      <div v-else class="space-y-2">
        <el-card
          v-for="item in pagedRecords"
          :key="item.id"
          shadow="hover"
          class="rounded-lg group relative"
          @click="toggleDelete(item.id)"
        >
          <button
            class="absolute right-3 top-3 h-7 w-7 rounded-full bg-red-500 text-white flex items-center justify-center opacity-0 transition group-hover:opacity-100"
            :class="activeDeleteId === item.id ? 'opacity-100' : ''"
            @click.stop="onDelete(item.id)"
          >
            <el-icon><Delete /></el-icon>
          </button>
          <div class="flex items-center justify-between">
            <div>
              <div class="text-sm text-gray-500 dark:text-gray-400">{{ formatDate(item.record_date) }}</div>
              <div class="text-base font-semibold text-gray-800 dark:text-gray-100">
                {{ formatRecordValue(item) }} {{ unit }}
              </div>
              <div v-if="type === 'vision'" class="text-xs text-gray-500 dark:text-gray-400">
                左眼 {{ item.left_value }} 右眼 {{ item.right_value }}
              </div>
              <div v-if="item.remark" class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                {{ item.remark }}
              </div>
            </div>
          </div>
        </el-card>
      </div>
      <div v-if="records.length > pageSize" class="flex justify-center pt-2">
        <el-pagination
          :page-size="pageSize"
          :current-page="currentPage"
          layout="prev, pager, next"
          :total="records.length"
          background
          @current-change="onPageChange"
        />
      </div>

      <div v-if="type === 'vision'" class="rounded-xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-3 text-sm text-gray-600 dark:text-gray-300">
        <div class="font-semibold text-gray-700 dark:text-gray-200 mb-1">记录说明</div>
        <div>采用5分记录法（LogMAR视力记录法的一种衍生形式）对国际标准视力表进行评估时，通常将视力结果为5.0（对应小数记录法的1.0）定义为正常视力。该数值表示被检查者能够清晰分辨视角为1分角（1 arcminute）的最小视标细节，代表最佳矫正视力水平。</div>
      </div>
    </div>

    <el-dialog v-model="showDialog" title="添加记录" width="360px">
      <div class="space-y-3">
        <div class="flex items-center gap-3">
          <div class="w-14 text-sm text-gray-600 dark:text-gray-300">日期</div>
          <el-date-picker v-model="formDate" type="date" value-format="YYYY-MM-DD" placeholder="选择日期" class="flex-1" />
        </div>
        <div v-if="type !== 'vision'" class="flex items-center gap-3">
          <div class="w-14 text-sm text-gray-600 dark:text-gray-300">{{ unitName }}</div>
          <el-input-number v-model="formValue" :min="0" :precision="2" :step="0.1" class="flex-1" controls-position="right" placeholder="请输入数值" />
        </div>
        <div v-else class="space-y-2">
          <div class="flex items-center gap-3">
            <div class="w-14 text-sm text-gray-600 dark:text-gray-300">左眼</div>
            <el-input-number v-model="formLeft" :min="0" :precision="2" :step="0.1" class="flex-1" controls-position="right" placeholder="请输入左眼数值" />
          </div>
          <div class="flex items-center gap-3">
            <div class="w-14 text-sm text-gray-600 dark:text-gray-300">右眼</div>
            <el-input-number v-model="formRight" :min="0" :precision="2" :step="0.1" class="flex-1" controls-position="right" placeholder="请输入右眼数值" />
          </div>
        </div>
        <div class="flex items-start gap-3">
          <div class="w-14 text-sm text-gray-600 dark:text-gray-300 pt-2">备注</div>
          <el-input v-model="formRemark" type="textarea" :rows="2" class="flex-1" placeholder="可填写本次记录说明" />
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="showDialog = false">取消</el-button>
          <el-button type="primary" @click="submit">保存</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { ArrowLeft, Delete } from '@element-plus/icons-vue'
import { useRoute } from 'vue-router'
import router from '@/router'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listGrowthRecords, createGrowthRecord, deleteGrowthRecord, type GrowthMetricRecord, type GrowthMetricType } from '@/services/growth'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import VChart from 'vue-echarts'

use([CanvasRenderer, BarChart, GridComponent, TooltipComponent])

const route = useRoute()
const type = computed(() => route.params.type as GrowthMetricType)

const typeConfig: Record<GrowthMetricType, { title: string; unit: string; unitName: string }> = {
  height: { title: '身高', unit: 'cm', unitName: '厘米' },
  weight: { title: '体重', unit: 'kg', unitName: '公斤' },
  vision: { title: '裸眼视力', unit: '', unitName: '视力' },
  foot: { title: '脚长', unit: 'cm', unitName: '厘米' }
}

const title = computed(() => typeConfig[type.value]?.title || '记录')
const unit = computed(() => typeConfig[type.value]?.unit || '')
const unitName = computed(() => typeConfig[type.value]?.unitName || '')

const records = ref<GrowthMetricRecord[]>([])
const pageSize = 10
const currentPage = ref(1)
const showDialog = ref(false)
const formDate = ref<string>(dayjs().format('YYYY-MM-DD'))
const formValue = ref<number | null>(null)
const formLeft = ref<number | null>(null)
const formRight = ref<number | null>(null)
const formRemark = ref('')
const chartTab = ref<'week' | 'month' | 'year' | 'all'>('year')
const activeDeleteId = ref<number | null>(null)

function formatDate(d: string) {
  return dayjs(d).format('YYYY-MM-DD')
}

function openDialog() {
  formDate.value = dayjs().format('YYYY-MM-DD')
  formValue.value = null
  formLeft.value = null
  formRight.value = null
  formRemark.value = ''
  showDialog.value = true
}

async function load() {
  const resp = await listGrowthRecords(type.value)
  records.value = resp.items || []
}

function recordNumericValue(item: GrowthMetricRecord) {
  if (type.value === 'vision') {
    const l = Number(item.left_value || 0)
    const r = Number(item.right_value || 0)
    if (!l || !r) return 0
    return Number(((l + r) / 2).toFixed(2))
  }
  return Number(item.value || 0)
}

function formatRecordValue(item: GrowthMetricRecord) {
  if (type.value === 'vision') {
    return `左${item.left_value} 右${item.right_value}`
  }
  return item.value
}

async function submit() {
  if (!formDate.value) {
    ElMessage.error('请选择日期')
    return
  }
  if (type.value === 'vision') {
    if (!formLeft.value || formLeft.value <= 0 || !formRight.value || formRight.value <= 0) {
      ElMessage.error('请输入有效视力数值')
      return
    }
    await createGrowthRecord({ metric_type: type.value, record_date: formDate.value, left_value: Number(formLeft.value), right_value: Number(formRight.value), remark: formRemark.value.trim() })
  } else {
    if (!formValue.value || formValue.value <= 0) {
      ElMessage.error('请输入有效数值')
      return
    }
    await createGrowthRecord({ metric_type: type.value, record_date: formDate.value, value: Number(formValue.value), remark: formRemark.value.trim() })
  }
  showDialog.value = false
  await load()
}

async function onDelete(id: number) {
  try {
    await ElMessageBox.confirm('确认删除该记录吗？', '提示', { type: 'warning' })
  } catch {
    return
  }
  await deleteGrowthRecord(id)
  await load()
}

function toggleDelete(id: number) {
  activeDeleteId.value = activeDeleteId.value === id ? null : id
}

function onPageChange(page: number) {
  currentPage.value = page
}

const pagedRecords = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return records.value.slice(start, start + pageSize)
})

watch(records, (list) => {
  const maxPage = Math.max(1, Math.ceil(list.length / pageSize))
  if (currentPage.value > maxPage) currentPage.value = 1
})

onMounted(load)

const weekOption = computed(() => {
  const start = dayjs().startOf('week')
  const labels = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  const values = Array(7).fill(0)
  const map = new Map<string, number>()
  for (const item of [...records.value].reverse()) {
    const d = dayjs(item.record_date).format('YYYY-MM-DD')
    map.set(d, recordNumericValue(item))
  }
  for (let i = 0; i < 7; i++) {
    const d = start.add(i, 'day').format('YYYY-MM-DD')
    values[i] = map.get(d) || 0
  }
  return {
    grid: { left: 8, right: 8, top: 10, bottom: 20, containLabel: true },
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: labels, axisLabel: { fontSize: 10 } },
    yAxis: { type: 'value', splitLine: { show: false } },
    series: [{ type: 'bar', data: values, itemStyle: { color: '#22c55e' } }]
  }
})

const monthOption = computed(() => {
  const now = dayjs()
  const days = now.daysInMonth()
  const labels = Array.from({ length: days }, (_, i) => String(i + 1))
  const values = Array(days).fill(0)
  const map = new Map<string, number>()
  for (const item of [...records.value].reverse()) {
    const d = dayjs(item.record_date)
    if (d.isSame(now, 'month')) {
      map.set(d.format('YYYY-MM-DD'), recordNumericValue(item))
    }
  }
  for (let i = 0; i < days; i++) {
    const d = now.date(i + 1).format('YYYY-MM-DD')
    values[i] = map.get(d) || 0
  }
  return {
    grid: { left: 8, right: 8, top: 10, bottom: 20, containLabel: true },
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: labels, axisLabel: { fontSize: 9 } },
    yAxis: { type: 'value', splitLine: { show: false } },
    series: [{ type: 'bar', data: values, itemStyle: { color: '#3b82f6' } }]
  }
})

const yearOption = computed(() => {
  const now = dayjs()
  const labels = Array.from({ length: 12 }, (_, i) => `${i + 1}月`)
  const values = Array(12).fill(0)
  const map = new Map<number, number>()
  for (const item of [...records.value].reverse()) {
    const d = dayjs(item.record_date)
    if (d.isSame(now, 'year')) {
      map.set(d.month(), recordNumericValue(item))
    }
  }
  for (let i = 0; i < 12; i++) {
    values[i] = map.get(i) || 0
  }
  return {
    grid: { left: 8, right: 8, top: 10, bottom: 20, containLabel: true },
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: labels, axisLabel: { fontSize: 9 } },
    yAxis: { type: 'value', splitLine: { show: false } },
    series: [{ type: 'bar', data: values, itemStyle: { color: '#8b5cf6' } }]
  }
})

const allOption = computed(() => {
  const sorted = [...records.value].sort((a, b) => dayjs(a.record_date).valueOf() - dayjs(b.record_date).valueOf())
  const map = new Map<number, number>()
  for (const item of sorted) {
    const y = dayjs(item.record_date).year()
    map.set(y, recordNumericValue(item))
  }
  const years = Array.from(map.keys()).sort((a, b) => a - b)
  const labels = years.map((y) => String(y))
  const values = years.map((y) => map.get(y) || 0)
  return {
    grid: { left: 8, right: 8, top: 10, bottom: 20, containLabel: true },
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: labels, axisLabel: { fontSize: 9 } },
    yAxis: { type: 'value', splitLine: { show: false } },
    series: [{ type: 'bar', data: values, itemStyle: { color: '#10b981' } }]
  }
})

const chartOption = computed(() => {
  if (chartTab.value === 'month') return monthOption.value
  if (chartTab.value === 'year') return yearOption.value
  if (chartTab.value === 'all') return allOption.value
  return weekOption.value
})
</script>

<style scoped>
</style>
