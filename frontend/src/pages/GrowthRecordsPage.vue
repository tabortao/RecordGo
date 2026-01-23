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

    <div class="p-4 space-y-3">
      <div v-if="records.length === 0" class="text-sm text-gray-500 dark:text-gray-400">暂无记录</div>
      <div v-else class="space-y-2">
        <el-card v-for="item in records" :key="item.id" shadow="hover" class="rounded-lg">
          <div class="flex items-center justify-between">
            <div>
              <div class="text-sm text-gray-500 dark:text-gray-400">{{ formatDate(item.record_date) }}</div>
              <div class="text-base font-semibold text-gray-800 dark:text-gray-100">
                {{ item.value }} {{ unit }}
              </div>
            </div>
            <el-button type="danger" size="small" @click="onDelete(item.id)">删除</el-button>
          </div>
        </el-card>
      </div>
    </div>

    <el-dialog v-model="showDialog" title="添加记录" width="360px">
      <div class="space-y-3">
        <el-date-picker v-model="formDate" type="date" value-format="YYYY-MM-DD" placeholder="选择日期" class="w-full" />
        <el-input-number v-model="formValue" :min="0" :precision="2" :step="0.1" class="w-full" controls-position="right" placeholder="请输入数值" />
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
import { computed, onMounted, ref } from 'vue'
import { ArrowLeft } from '@element-plus/icons-vue'
import { useRoute } from 'vue-router'
import router from '@/router'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listGrowthRecords, createGrowthRecord, deleteGrowthRecord, type GrowthMetricRecord, type GrowthMetricType } from '@/services/growth'

const route = useRoute()
const type = computed(() => route.params.type as GrowthMetricType)

const typeConfig: Record<GrowthMetricType, { title: string; unit: string }> = {
  height: { title: '身高', unit: 'cm' },
  weight: { title: '体重', unit: 'kg' },
  vision: { title: '裸眼视力', unit: '' },
  foot: { title: '脚长', unit: 'cm' }
}

const title = computed(() => typeConfig[type.value]?.title || '记录')
const unit = computed(() => typeConfig[type.value]?.unit || '')

const records = ref<GrowthMetricRecord[]>([])
const showDialog = ref(false)
const formDate = ref<string>(dayjs().format('YYYY-MM-DD'))
const formValue = ref<number | null>(null)

function formatDate(d: string) {
  return dayjs(d).format('YYYY-MM-DD')
}

function openDialog() {
  formDate.value = dayjs().format('YYYY-MM-DD')
  formValue.value = null
  showDialog.value = true
}

async function load() {
  const resp = await listGrowthRecords(type.value)
  records.value = resp.items || []
}

async function submit() {
  if (!formDate.value) {
    ElMessage.error('请选择日期')
    return
  }
  if (!formValue.value || formValue.value <= 0) {
    ElMessage.error('请输入有效数值')
    return
  }
  await createGrowthRecord({ metric_type: type.value, record_date: formDate.value, value: Number(formValue.value) })
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

onMounted(load)
</script>

<style scoped>
</style>
