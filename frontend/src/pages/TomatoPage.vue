<template>
  <!-- 中文注释：独立的番茄钟页面（夜间主题），顶部返回图标；深色背景响应式居中展示定时器 -->
  <div class="h-screen p-4 space-y-4 overflow-hidden flex flex-col" :style="{ backgroundColor: '#30302E' }">
    <div ref="topRef">
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer" :style="{ color: '#B8CEE8' }" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" :style="{ color: '#B8CEE8' }"><Clock /></el-icon>
      <h2 class="font-semibold" :style="{ color: '#B8CEE8' }">番茄钟</h2>
      <!-- 右上角系统时间显示（字号与标题一致） -->
      <h2 class="ml-auto font-mono font-semibold" :style="{ color: '#B8CEE8', fontSize: '1.3em' }">{{ systemTime }}</h2>
    </div>
    <!-- 任务标题：靠近页面标题，避免紧邻下方为时钟 -->
    <div class="mt-1 mb-3 text-center font-bold" :style="{ color: '#B8CEE8', fontSize: '1.2em' }" v-if="taskName">任务：{{ taskName }}</div>
    </div>
    <!-- 中文注释：中部容器高度按 calc(100vh - 顶部高度 - 底部高度) 计算，确保垂直居中且无滚动 -->
    <div class="flex items-center justify-center" :style="{ height: midHeight }">
      <TomatoTimer :work-minutes="workMinutes" :break-minutes="5" :task-name="taskName" :task-remark="taskRemark" @complete="onTomatoComplete" />
    </div>
  </div>
</template>

<script setup lang="ts">
// 中文注释：番茄钟独立页面，从路由参数中读取任务ID并加载任务信息
import { ref, onMounted, onUnmounted, computed } from 'vue'
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
// 顶部与底部高度测量，用于计算中部容器高度
const topRef = ref<HTMLElement | null>(null)
const topH = ref(0)
const bottomH = ref(0)
const midHeight = computed(() => `calc(100vh - ${topH.value}px - ${bottomH.value}px)`)
function measureHeights() {
  topH.value = topRef.value?.offsetHeight || 0
  const bottomEl = document.querySelector('[data-bottom="tomato-controls"]') as HTMLElement | null
  bottomH.value = bottomEl?.offsetHeight || 120
}

onMounted(async () => {
  // 初始化系统时间显示并每秒更新（不显示秒）
  systemTime.value = dayjs().format('HH:mm')
  clockTimer = setInterval(() => {
    systemTime.value = dayjs().format('HH:mm')
  }, 1000)
  // 计算顶部与底部固定区高度
  measureHeights()
  window.addEventListener('resize', measureHeights)
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
  // 移除监听
  window.removeEventListener('resize', measureHeights)
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