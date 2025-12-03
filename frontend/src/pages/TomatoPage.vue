<template>
  <!-- 中文注释：独立的番茄钟页面（夜间主题），顶部返回图标；深色背景响应式居中展示定时器 -->
  <div ref="rootRef" class="h-screen p-4 overflow-hidden flex flex-col" :style="{ backgroundColor: '#30302E' }">
    <div ref="topRef" class="pb-3">
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer" :style="{ color: '#B8CEE8' }" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" :style="{ color: '#B8CEE8' }"><Clock /></el-icon>
      <h2 class="font-semibold" :style="{ color: '#B8CEE8' }">番茄钟</h2>
      <!-- 右上角系统时间显示（字号与标题一致） -->
      <h2 class="ml-auto font-mono font-semibold" :style="{ color: '#B8CEE8', fontSize: '1.3em' }">{{ systemTime }}</h2>
      <el-dropdown @command="onTopMenu">
        <span class="el-dropdown-link">
          <svg width="18" height="18" viewBox="0 0 24 24" :style="{ color: '#B8CEE8' }">
            <rect x="3" y="5" width="18" height="2" rx="1" fill="currentColor" />
            <rect x="3" y="11" width="18" height="2" rx="1" fill="currentColor" />
            <rect x="3" y="17" width="18" height="2" rx="1" fill="currentColor" />
          </svg>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="adjust">修改番茄钟倒计时</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
    <!-- 任务标题：靠近页面标题，避免紧邻下方为时钟 -->
    <div class="mt-1 text-center font-bold" :style="{ color: '#B8CEE8', fontSize: '1.2em' }" v-if="taskName">任务：{{ taskName }}</div>
    <!-- 备注：显示在任务描述下方，颜色与番茄钟页面统一 -->
    <div class="text-center text-xs" :style="{ color: '#9aa8b8' }" v-if="taskRemark">备注：{{ taskRemark }}</div>
    </div>
    <!-- 中文注释：中部容器高度按 calc(100vh - 顶部高度 - 底部高度) 计算，确保垂直居中且无滚动 -->
    <div class="flex items-center justify-center" :style="{ height: midHeight.value }">
      <TomatoTimer ref="timerRef" :work-minutes="workMinutes.value" :break-minutes="5" :task-name="taskName.value" :task-remark="taskRemark.value" :task-id="taskId" @complete="onTomatoComplete" />
    </div>
  </div>
  <el-dialog v-model="showAdjust" title="修改番茄钟倒计时" width="360px" :teleported="false" :modal="false" class="dark-dialog">
    <div class="flex items-center justify-between">
      <span>分钟</span>
      <el-input-number v-model="adjustMinutes" :min="1" :max="240" :step="5" />
    </div>
    <template #footer>
      <div class="flex justify-end gap-2">
        <el-button @click="showAdjust.value=false" class="dark-btn">取消</el-button>
        <el-button type="primary" @click="applyAdjust" class="dark-btn-primary">应用</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
// 中文注释：番茄钟独立页面，从路由参数中读取任务ID并加载任务信息
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { ArrowLeft, Clock } from '@element-plus/icons-vue'
import router from '@/router'
import TomatoTimer from '@/components/TomatoTimer.vue'
import { useRoute } from 'vue-router'
import { useAppState } from '@/stores/appState'
import { getTask, completeTomato, updateTaskStatus } from '@/services/tasks'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'

const route = useRoute()
const store = useAppState()
const timerRef = ref<InstanceType<typeof TomatoTimer> | null>(null)
const taskId = Number(route.params.id)
function goBack() {
  if (store.tomato.running) {
    if (store.tomato.fixedTomatoPage) {
      // 中文注释：固定番茄钟页面下，返回视为未完成，立即停止并重置到预计时长；不显示悬浮球
      timerRef.value?.stop?.()
      const resetSec = (store.tomato.durationMinutes || 20) * 60
      store.updateTomato({ running: false, remainingSeconds: resetSec, currentTaskId: null, showFloating: false, runningMode: null })
    } else {
      store.updateTomato({ showFloating: true })
    }
  }
  router.back()
}

const workMinutes = ref<number>(20)
const taskName = ref<string>('')
const taskRemark = ref<string>('')
const systemTime = ref<string>('')
// 顶部菜单：调整倒计时对话框
const showAdjust = ref(false)
const adjustMinutes = ref<number>(20)
function onTopMenu(cmd: string) {
  if (cmd === 'adjust') {
    adjustMinutes.value = workMinutes.value || 20
    showAdjust.value = true
  }
}
function applyAdjust() {
  const wasRunning = store.tomato.running
  timerRef.value?.stop?.()
  workMinutes.value = Math.max(1, Math.min(240, adjustMinutes.value))
  if (wasRunning) timerRef.value?.start?.()
  showAdjust.value = false
}
let clockTimer: any = null
// 顶部与底部高度测量，用于计算中部容器高度
const rootRef = ref<HTMLElement | null>(null)
const topRef = ref<HTMLElement | null>(null)
const topH = ref(0)
const bottomH = ref(0)
const padTop = ref(0)
const padBottom = ref(0)
const midHeight = computed(() => `calc(100vh - ${topH.value}px - ${bottomH.value}px - ${padTop.value}px - ${padBottom.value}px)`)
function measureHeights() {
  topH.value = topRef.value?.offsetHeight || 0
  const bottomEl = document.querySelector('[data-bottom="tomato-controls"]') as HTMLElement | null
  bottomH.value = bottomEl?.offsetHeight || 120
  if (rootRef.value) {
    const cs = getComputedStyle(rootRef.value)
    padTop.value = parseFloat(cs.paddingTop || '0') || 0
    padBottom.value = parseFloat(cs.paddingBottom || '0') || 0
  }
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
  // 全局禁用滚动条，确保页面无滚动
  document.documentElement.style.overflowY = 'hidden'
  document.body.style.overflowY = 'hidden'
  // 中文注释：进入番茄钟页面时隐藏悬浮球，避免同时显示
  store.updateTomato({ showFloating: false })
  if (!isNaN(taskId)) {
    try {
      const t = await getTask(taskId)
      taskName.value = t.name
      taskRemark.value = t.remark || t.description
      workMinutes.value = t.plan_minutes || 20
      // 中文注释：若当前存在同一任务正在计时，则不重置剩余时间（保持连续性）；
      // 仅在未运行或切换到不同任务时，重置为预计时间。
      const expectedSec = (workMinutes.value || 20) * 60
      const isCountup = store.tomato.mode === 'countup'
      const sameRunning = store.tomato.running && store.tomato.currentTaskId === taskId
      if (store.tomato.fixedTomatoPage) {
        // 中文注释：固定页面模式下，每次进入根据设置模式重置：正计时为 0，倒计时为预计时长；保持未运行（需用户点击“开始”）
        store.updateTomato({ remainingSeconds: isCountup ? 0 : expectedSec, durationMinutes: workMinutes.value, currentTaskId: taskId, running: false, runningMode: null })
      } else {
        if (!sameRunning) {
          // 中文注释：非固定模式且未连续计时：进入页面时根据设置模式重置剩余/已用秒数
          store.updateTomato({ remainingSeconds: isCountup ? 0 : expectedSec, durationMinutes: workMinutes.value, currentTaskId: taskId, runningMode: null })
        } else {
          // 保持当前剩余时间，仅同步任务ID与时长
          store.updateTomato({ durationMinutes: workMinutes.value, currentTaskId: taskId })
        }
      }
    } catch (e: any) {
      ElMessage.error(e.message || '加载任务失败')
    }
  }
})

onUnmounted(() => {
  if (clockTimer) clearInterval(clockTimer)
  // 移除监听
  window.removeEventListener('resize', measureHeights)
  // 恢复页面滚动
  document.documentElement.style.overflowY = ''
  document.body.style.overflowY = ''
})

async function onTomatoComplete(seconds?: number) {
  try {
    if (!taskId) return
    const usedSec = Math.max(1, seconds || workMinutes.value * 60)
    const reportMinutes = Math.max(1, Math.round(usedSec / 60))
    await completeTomato(taskId, reportMinutes)
    // 中文注释：番茄钟完成允许在只读权限下标记“已完成”，后端识别 allow_by_tomato 放行
    await updateTaskStatus(taskId, 2, { allowByTomato: true })
    ElMessage.success('番茄钟完成，数据已记录')
    router.back()
  } catch (e: any) {
    ElMessage.error(`番茄上报失败：${e.message || e}`)
  }
}
</script>

<style scoped>
/* 中文注释：页面使用网格进行信息分区，移动端单列显示 */
:global(.el-overlay) { background-color: rgba(0,0,0,0.5); }
.dark-dialog { background-color: rgb(48,48,46); color: #B8CEE8; }
.dark-dialog :deep(.el-dialog) { background-color: rgb(48,48,46) !important; color: #B8CEE8 !important; border: 1px solid #4a4a48; }
.dark-dialog :deep(.el-dialog__header) { background-color: rgb(48,48,46); color: #B8CEE8; border-bottom: 1px solid #4a4a48; }
.dark-dialog :deep(.el-dialog__body) { background-color: rgb(48,48,46); color: #B8CEE8; }
.dark-dialog :deep(.el-dialog__footer) { background-color: rgb(48,48,46); border-top: 1px solid #4a4a48; }
/* 覆盖遮罩层与弹窗背景（通过 modal-class 传入）*/
:global(.dark-overlay) { background-color: rgba(0,0,0,0.6); }
:global(.dark-overlay) .el-overlay-dialog .el-dialog { background-color: rgb(48,48,46) !important; color: #B8CEE8 !important; border: 1px solid #4a4a48; }
.dark-dialog :deep(.el-input__wrapper),
.dark-dialog :deep(.el-input-number__decrease),
.dark-dialog :deep(.el-input-number__increase) { background-color: #3a3a38 !important; color: #B8CEE8 !important; box-shadow: 0 0 0 1px #4a4a48 inset !important; }
.dark-dialog :deep(.el-input__inner) { color: #B8CEE8 !important; }
.dark-btn { background-color: #3a3a38 !important; color: #B8CEE8 !important; border-color: #4a4a48 !important; }
.dark-btn-primary { background-color: #3a3a38 !important; color: #B8CEE8 !important; border-color: #4a4a48 !important; }
:global(.el-dropdown__popper .el-popper__content) { background-color: rgb(48,48,46); color: #B8CEE8; border: 1px solid #4a4a48; }
:global(.el-dropdown-menu__item) { color: #B8CEE8; }
:global(.el-dropdown-menu__item:hover) { background-color: rgb(48,48,46) !important; color: #B8CEE8 !important; }
</style>
