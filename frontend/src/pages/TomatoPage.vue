<template>
  <!-- 中文注释：独立的番茄钟页面（夜间主题），顶部返回图标；深色背景响应式居中展示定时器 -->
  <div class="h-screen p-4 space-y-4" :style="{ backgroundColor: '#30302E' }">
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer" :style="{ color: '#B8CEE8' }" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" :style="{ color: '#B8CEE8' }"><Clock /></el-icon>
      <h2 class="font-semibold" :style="{ color: '#B8CEE8' }">番茄钟</h2>
      <!-- 右上角系统时间显示（字号与标题一致） -->
      <h2 class="ml-auto font-mono font-semibold" :style="{ color: '#B8CEE8' }">{{ systemTime }}</h2>
    </div>
    <!-- 中文注释：将任务标题、备注、预计时长融入到上方定时器组件，取消下方信息卡片；居中显示 -->
    <div class="max-w-3xl mx-auto">
      <TomatoTimer :work-minutes="workMinutes" :break-minutes="5" :task-name="taskName" :task-remark="taskRemark" @complete="onTomatoComplete" />
    </div>
  </div>
</template>

<script setup lang="ts">
// 中文注释：番茄钟独立页面，从路由参数中读取任务ID并加载任务信息
import { ref, onMounted, onUnmounted } from 'vue'
import { ArrowLeft, Clock } from '@element-plus/icons-vue'
import router from '@/router'
import TomatoTimer from '@/components/TomatoTimer.vue'
import { useRoute } from 'vue-router'
import { getTask, completeTomato, updateTaskStatus } from '@/services/tasks'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'

const route = useRoute()
const taskId = Number(route.params.id)
function goBack() { router.back() }

const workMinutes = ref<number>(20)
const taskName = ref<string>('')
const taskRemark = ref<string>('')
const systemTime = ref<string>('')
let clockTimer: any = null

onMounted(async () => {
  // 初始化系统时间显示并每秒更新
  systemTime.value = dayjs().format('HH:mm:ss')
  clockTimer = setInterval(() => {
    systemTime.value = dayjs().format('HH:mm:ss')
  }, 1000)
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

onUnmounted(() => {
  if (clockTimer) clearInterval(clockTimer)
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