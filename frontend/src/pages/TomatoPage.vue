<template>
  <!-- 中文注释：独立的番茄钟页面，顶部返回图标；页面居中展示定时器，响应式布局 -->
  <div class="p-4 space-y-4">
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer text-gray-600" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" class="text-red-500"><Clock /></el-icon>
      <h2 class="font-semibold">番茄钟</h2>
    </div>
    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <el-card class="rounded-lg md:col-span-2 xl:col-span-2 flex items-center justify-center">
        <!-- 中文注释：定时器组件居中，传入任务信息与默认时间 -->
        <TomatoTimer :work-minutes="workMinutes" :break-minutes="5" :task-name="taskName" :task-remark="taskRemark" @complete="onTomatoComplete" />
      </el-card>
      <el-card class="rounded-lg">
        <div class="text-sm text-gray-600">
          任务：{{ taskName || '未指定' }}
        </div>
        <div class="text-sm text-gray-600 mt-1">备注：{{ taskRemark || '无' }}</div>
        <div class="text-sm text-gray-600 mt-1">预设工作时长：{{ workMinutes }} 分钟</div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
// 中文注释：番茄钟独立页面，从路由参数中读取任务ID并加载任务信息
import { ref, onMounted } from 'vue'
import { ArrowLeft, Clock } from '@element-plus/icons-vue'
import router from '@/router'
import TomatoTimer from '@/components/TomatoTimer.vue'
import { useRoute } from 'vue-router'
import { getTask, completeTomato, updateTaskStatus } from '@/services/tasks'
import { ElMessage } from 'element-plus'

const route = useRoute()
const taskId = Number(route.params.id)
function goBack() { router.back() }

const workMinutes = ref<number>(20)
const taskName = ref<string>('')
const taskRemark = ref<string>('')

onMounted(async () => {
  if (!isNaN(taskId)) {
    try {
      const t = await getTask(taskId)
      taskName.value = t.name
      taskRemark.value = t.remark || t.description
      workMinutes.value = t.plan_minutes || 20
    } catch (e: any) {
      ElMessage.error(e.message || '加载任务失败')
    }
  }
})

async function onTomatoComplete(seconds?: number) {
  try {
    if (!taskId) return
    const usedSec = Math.max(1, seconds || workMinutes.value * 60)
    const reportMinutes = Math.max(1, Math.round(usedSec / 60))
    await completeTomato(taskId, reportMinutes)
    await updateTaskStatus(taskId, 2)
    ElMessage.success('番茄钟完成，数据已记录')
    router.back()
  } catch (e: any) {
    ElMessage.error(`番茄上报失败：${e.message || e}`)
  }
}
</script>

<style scoped>
/* 中文注释：页面使用网格进行信息分区，移动端单列显示 */
</style>