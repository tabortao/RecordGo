<template>
  <div class="min-h-screen bg-[#F5F7FA] dark:bg-gray-900 flex flex-col relative">
    <!-- Header -->
    <div class="px-4 py-3 bg-white/80 dark:bg-gray-900/80 backdrop-blur-md sticky top-0 z-30 flex items-center justify-between border-b border-gray-100 dark:border-gray-800 transition-colors">
      <div class="flex items-center gap-3">
        <div 
          class="w-8 h-8 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center cursor-pointer hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors text-gray-600 dark:text-gray-300"
          @click="router.back()"
        >
          <el-icon><ArrowLeft /></el-icon>
        </div>
        <h2 class="font-bold text-gray-800 dark:text-gray-100 text-lg">{{ title }}记录</h2>
      </div>
      <button 
        class="bg-gradient-to-r from-blue-500 to-indigo-600 text-white px-4 py-1.5 rounded-full text-sm font-bold shadow-lg shadow-blue-200 dark:shadow-blue-900/20 hover:scale-105 active:scale-95 transition-all flex items-center gap-1"
        @click="openDialog"
      >
        <el-icon><Plus /></el-icon>
        <span>添加记录</span>
      </button>
    </div>

    <div class="flex-1 overflow-y-auto p-4 sm:p-6 pb-24">
      <div class="max-w-3xl mx-auto space-y-6">
        <!-- Chart Card -->
        <div class="bg-white dark:bg-gray-800 rounded-3xl p-5 shadow-sm border border-gray-100 dark:border-gray-700">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-sm font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider flex items-center gap-2">
              <el-icon><TrendCharts /></el-icon> 数据趋势
            </h3>
            <div class="flex bg-gray-100 dark:bg-gray-700/50 p-1 rounded-xl">
              <button 
                v-for="tab in [{k: 'week', l: '周'}, {k: 'month', l: '月'}, {k: 'year', l: '年'}, {k: 'all', l: '全部'}]" 
                :key="tab.k"
                class="px-3 py-1 rounded-lg text-xs font-bold transition-all"
                :class="chartTab === tab.k ? 'bg-white dark:bg-gray-600 text-indigo-600 dark:text-indigo-300 shadow-sm' : 'text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300'"
                @click="chartTab = tab.k as any"
              >
                {{ tab.l }}
              </button>
            </div>
          </div>
          <div class="h-56 w-full">
            <VChart :option="chartOption" autoresize class="w-full h-full" />
          </div>
        </div>

        <!-- Vision Guide -->
        <div v-if="type === 'vision'" class="bg-blue-50/50 dark:bg-blue-900/10 rounded-2xl p-4 border border-blue-100 dark:border-blue-800/30 flex gap-3">
          <el-icon class="text-blue-500 text-lg flex-shrink-0 mt-0.5"><InfoFilled /></el-icon>
          <div class="text-sm text-blue-700 dark:text-blue-300 leading-relaxed">
            <p class="font-bold mb-1">视力记录说明</p>
            <p class="opacity-90">采用 5 分记录法（LogMAR）评估。5.0（对应小数 1.0）为正常视力标准。建议定期检查，关注视力变化趋势。</p>
          </div>
        </div>

        <!-- Records List -->
        <div class="space-y-4">
          <div class="flex items-center justify-between px-1">
            <h3 class="text-sm font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider">历史数据</h3>
            <span class="text-xs font-medium text-gray-400 bg-gray-100 dark:bg-gray-700/50 px-2 py-0.5 rounded-full">共 {{ records.length }} 条</span>
          </div>

          <div v-if="records.length === 0" class="flex flex-col items-center justify-center py-16 text-gray-400 bg-white dark:bg-gray-800 rounded-3xl border border-dashed border-gray-200 dark:border-gray-700">
            <div class="w-16 h-16 bg-gray-50 dark:bg-gray-700/50 rounded-full flex items-center justify-center mb-3">
              <el-icon :size="24" class="text-gray-300"><Document /></el-icon>
            </div>
            <p class="text-sm">暂无记录，点击右上角添加</p>
          </div>

          <div v-else class="space-y-3">
            <div
              v-for="item in pagedRecords"
              :key="item.id"
              class="group bg-white dark:bg-gray-800 rounded-2xl p-4 border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-md hover:-translate-y-0.5 transition-all duration-300 relative overflow-hidden"
              @click="toggleDelete(item.id)"
            >
              <!-- Delete Overlay -->
              <div 
                class="absolute inset-y-0 right-0 w-16 bg-gradient-to-l from-white via-white to-transparent dark:from-gray-800 dark:via-gray-800 flex items-center justify-end pr-4 transition-opacity duration-200"
                :class="activeDeleteId === item.id ? 'opacity-100 z-20' : 'opacity-0 -z-10 group-hover:opacity-100 group-hover:z-20'"
              >
                <button
                  class="h-8 w-8 rounded-full bg-red-50 text-red-500 hover:bg-red-500 hover:text-white flex items-center justify-center transition-colors shadow-sm"
                  @click.stop="onDelete(item.id)"
                >
                  <el-icon><Delete /></el-icon>
                </button>
              </div>

              <div class="flex items-center justify-between relative z-10">
                <div class="flex items-center gap-4">
                  <div class="w-12 h-12 rounded-xl bg-indigo-50 dark:bg-indigo-900/20 flex flex-col items-center justify-center text-indigo-600 dark:text-indigo-400 font-bold border border-indigo-100 dark:border-indigo-800/30">
                    <span class="text-xs opacity-60">{{ dayjs(item.record_date).format('MM') }}月</span>
                    <span class="text-sm leading-none">{{ dayjs(item.record_date).format('DD') }}</span>
                  </div>
                  
                  <div>
                    <div class="flex items-baseline gap-1.5">
                      <span class="text-xl font-black text-gray-900 dark:text-white tracking-tight">
                        {{ formatRecordValue(item) }}
                      </span>
                      <span class="text-xs font-bold text-gray-400" v-if="unit">{{ unit }}</span>
                    </div>
                    
                    <div class="flex items-center gap-2 mt-0.5">
                      <span class="text-xs text-gray-400">{{ dayjs(item.record_date).format('YYYY年') }}</span>
                      <div v-if="type === 'vision'" class="text-[10px] bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-300 px-1.5 py-0.5 rounded font-medium">
                        左 {{ item.left_value }} / 右 {{ item.right_value }}
                      </div>
                    </div>
                  </div>
                </div>

                <div v-if="item.remark" class="hidden sm:block text-xs text-gray-400 bg-gray-50 dark:bg-gray-700/50 px-2 py-1 rounded-lg max-w-[120px] truncate">
                  {{ item.remark }}
                </div>
              </div>
              
              <!-- Mobile Remark -->
              <div v-if="item.remark" class="sm:hidden mt-3 pt-2 border-t border-gray-50 dark:border-gray-700/50 text-xs text-gray-500 dark:text-gray-400 flex items-start gap-1">
                <el-icon class="mt-0.5 text-gray-300"><ChatLineSquare /></el-icon>
                <span>{{ item.remark }}</span>
              </div>
            </div>
          </div>

          <div v-if="records.length > pageSize" class="flex justify-center pt-4">
            <el-pagination
              :page-size="pageSize"
              :current-page="currentPage"
              layout="prev, pager, next"
              :total="records.length"
              background
              @current-change="onPageChange"
              class="custom-pagination"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Add Dialog -->
    <el-dialog 
      v-model="showDialog" 
      :title="`添加${title}记录`" 
      width="90%" 
      class="max-w-md rounded-2xl overflow-hidden"
      :show-close="false"
      align-center
    >
      <template #header>
        <div class="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
          <div class="w-1 bg-indigo-500 h-4 rounded-full"></div>
          记录新数据
        </div>
      </template>

      <div class="space-y-5 py-2">
        <div class="space-y-1.5">
          <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">记录日期</label>
          <el-date-picker 
            v-model="formDate" 
            type="date" 
            value-format="YYYY-MM-DD" 
            placeholder="选择日期" 
            class="!w-full custom-input" 
            :clearable="false"
          />
        </div>

        <div v-if="type !== 'vision'" class="space-y-1.5">
          <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">{{ title }}数值 ({{ unit }})</label>
          <div class="relative">
            <el-input-number 
              v-model="formValue" 
              :min="0" 
              :precision="2" 
              :step="0.1" 
              class="!w-full custom-input-number" 
              controls-position="right" 
              placeholder="0.00" 
            />
            <div class="absolute right-10 top-1/2 -translate-y-1/2 text-gray-400 text-sm font-bold pointer-events-none">{{ unit }}</div>
          </div>
        </div>

        <div v-else class="grid grid-cols-2 gap-4">
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">左眼视力</label>
            <el-input-number 
              v-model="formLeft" 
              :min="0" 
              :max="5.3"
              :precision="1" 
              :step="0.1" 
              class="!w-full custom-input-number" 
              controls-position="right" 
              placeholder="5.0" 
            />
          </div>
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">右眼视力</label>
            <el-input-number 
              v-model="formRight" 
              :min="0" 
              :max="5.3"
              :precision="1" 
              :step="0.1" 
              class="!w-full custom-input-number" 
              controls-position="right" 
              placeholder="5.0" 
            />
          </div>
        </div>

        <div class="space-y-1.5">
          <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">备注说明</label>
          <el-input 
            v-model="formRemark" 
            type="textarea" 
            :rows="3" 
            placeholder="例如：最近吃得比较多..." 
            class="custom-input"
            resize="none"
          />
        </div>
      </div>

      <template #footer>
        <div class="flex gap-3 pt-2">
          <button 
            class="flex-1 py-2.5 rounded-xl bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 font-bold text-sm hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
            @click="showDialog = false"
          >
            取消
          </button>
          <button 
            class="flex-1 py-2.5 rounded-xl bg-gradient-to-r from-blue-500 to-indigo-600 text-white font-bold text-sm shadow-lg shadow-blue-200 dark:shadow-blue-900/20 hover:scale-[1.02] active:scale-95 transition-all"
            @click="submit"
          >
            保存记录
          </button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { ArrowLeft, Delete, Plus, TrendCharts, InfoFilled, Document, ChatLineSquare } from '@element-plus/icons-vue'
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
.custom-pagination :deep(.el-pagination.is-background .el-pager li:not(.is-disabled).is-active) {
  background-color: #6366f1;
}
.dark .custom-pagination :deep(.el-pagination.is-background .el-pager li:not(.is-disabled).is-active) {
  background-color: #4f46e5;
}

.custom-input :deep(.el-input__wrapper) {
  background-color: #f9fafb;
  border-radius: 0.75rem;
  box-shadow: none !important;
  border: 1px solid #e5e7eb;
  padding: 8px 12px;
  transition: all 0.2s;
}
.dark .custom-input :deep(.el-input__wrapper) {
  background-color: #374151;
  border-color: #4b5563;
}
.custom-input :deep(.el-input__wrapper:hover),
.custom-input :deep(.el-input__wrapper.is-focus) {
  border-color: #6366f1;
  background-color: #fff;
}
.dark .custom-input :deep(.el-input__wrapper:hover),
.dark .custom-input :deep(.el-input__wrapper.is-focus) {
  border-color: #6366f1;
  background-color: #1f2937;
}

.custom-input :deep(.el-textarea__inner) {
  background-color: #f9fafb;
  border-radius: 0.75rem;
  box-shadow: none !important;
  border: 1px solid #e5e7eb;
  padding: 12px;
}
.dark .custom-input :deep(.el-textarea__inner) {
  background-color: #374151;
  border-color: #4b5563;
  color: #fff;
}
.custom-input :deep(.el-textarea__inner:hover),
.custom-input :deep(.el-textarea__inner:focus) {
  border-color: #6366f1;
  background-color: #fff;
}
.dark .custom-input :deep(.el-textarea__inner:hover),
.dark .custom-input :deep(.el-textarea__inner:focus) {
  border-color: #6366f1;
  background-color: #1f2937;
}

.custom-input-number :deep(.el-input__wrapper) {
  background-color: #f9fafb;
  border-radius: 0.75rem;
  box-shadow: none !important;
  border: 1px solid #e5e7eb;
}
.dark .custom-input-number :deep(.el-input__wrapper) {
  background-color: #374151;
  border-color: #4b5563;
}
.custom-input-number :deep(.el-input__wrapper:hover),
.custom-input-number :deep(.el-input__wrapper.is-focus) {
  border-color: #6366f1;
}
</style>
