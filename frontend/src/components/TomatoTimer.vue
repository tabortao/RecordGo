<template>
  <!-- 中文注释：番茄钟组件（增强版），支持倒计时/正计时、预设与自定义、开始/暂停/继续/完成 -->
  <div class="p-4 space-y-4 relative">

    <!-- 中文注释：移除组件内部标题，避免与弹窗标题重复显示 -->

    <!-- 中文注释：时间居中，右侧提供模式切换图标；备注在时间下方显示为灰色小字 -->
    <div class="flex items-center justify-center gap-3">
      <div class="text-6xl font-mono">{{ mm }}:{{ ss }}</div>
      <el-icon
        class="cursor-pointer"
        :title="mode==='countdown' ? '切换为正计时' : '切换为倒计时'"
        :style="{ color: mode==='countdown' ? '#ef4444' : '#10b981' }"
        @click="toggleMode"
      >
        <RefreshRight />
      </el-icon>
    </div>
    <div class="text-center text-xs text-gray-500" v-if="props.taskRemark">{{ props.taskRemark }}</div>

    <!-- 中文注释：倒计时模式下显示预设与自定义；正计时不显示这些 -->
    <div class="flex items-center justify-center gap-3" v-if="mode==='countdown'">
      <div class="flex items-center gap-2">
        <el-tag class="cursor-pointer" @click="setDuration(10)">10分钟</el-tag>
        <el-tag class="cursor-pointer" @click="setDuration(20)">20分钟</el-tag>
      </div>
      <el-input-number v-model="customMinutes" :min="1" :max="240" />
      <el-button size="small" @click="applyCustom">应用</el-button>
    </div>

    <!-- 控制按钮 -->
    <div class="flex justify-center gap-2 mt-2">
      <el-button type="primary" @click="start" :disabled="running">开始</el-button>
      <el-button :type="running ? 'warning' : 'success'" @click="togglePauseResume" :disabled="!started">{{ running ? '暂停' : '继续' }}</el-button>
      <el-button type="warning" @click="reset">重置</el-button>
      <el-button type="danger" @click="complete">完成</el-button>
    </div>

    <div class="text-center text-gray-500 text-xs">完成后将上报实际耗时并自动标记任务完成</div>
  </div>
</template>

<script setup lang="ts">
// 中文注释：接收任务名与备注 + 工作/休息分钟数（合并定义）
// 中文注释：番茄钟逻辑，工作与休息两个阶段，倒计时结束后切换阶段并在工作结束时触发 complete 事件
import { computed, onUnmounted, ref, watch } from 'vue'
import { useAppState } from '@/stores/appState'
import { RefreshRight } from '@element-plus/icons-vue'

const props = defineProps<{ workMinutes?: number; breakMinutes?: number; taskName?: string; taskRemark?: string }>()
const emit = defineEmits<{ (e: 'complete', seconds: number): void }>()
const store = useAppState()

const workM = computed(() => props.workMinutes ?? store.tomato.durationMinutes ?? 20)
const breakM = computed(() => props.breakMinutes ?? 5)

type Phase = 'work' | 'break'
const phase = ref<Phase>('work')
const mode = ref<'countdown' | 'countup'>(store.tomato.mode)
const running = ref(store.tomato.running)
const remaining = ref(store.tomato.remainingSeconds || (mode.value === 'countdown' ? workM.value * 60 : 0))
const started = ref(false)
const customMinutes = ref(workM.value)
let timer: any = null

const phaseText = computed(() => (phase.value === 'work' ? '工作中' : '休息中'))
const headerText = computed(() => (mode.value === 'countdown' ? '倒计时' : '正计时'))
const mm = computed(() => String(Math.floor(remaining.value / 60)).padStart(2, '0'))
const ss = computed(() => String(remaining.value % 60).padStart(2, '0'))

function toggleMode() {
  // 中文注释：切换模式，正计时从 00:00 开始，倒计时从计划时长开始
  mode.value = mode.value === 'countdown' ? 'countup' : 'countdown'
  const wasRunning = running.value
  stopInternal()
  remaining.value = mode.value === 'countdown' ? workM.value * 60 : 0
  store.updateTomato({ mode: mode.value, remainingSeconds: remaining.value })
  if (wasRunning) start()
}

function tick() {
  if (mode.value === 'countdown') {
    if (remaining.value > 0) {
      remaining.value -= 1
      store.updateTomato({ remainingSeconds: remaining.value })
      return
    }
    stopInternal()
    if (phase.value === 'work') {
      emit('complete')
      phase.value = 'break'
      remaining.value = breakM.value * 60
    } else {
      phase.value = 'work'
      remaining.value = workM.value * 60
    }
  } else {
    // 正计时：向上累计，到达目标分钟触发完成
    remaining.value += 1
    store.updateTomato({ remainingSeconds: remaining.value })
    if (remaining.value >= workM.value * 60) {
      stopInternal()
      emit('complete')
    }
  }
}

function start() {
  if (running.value) return
  running.value = true
  started.value = true
  // 正计时从 00:00 开始
  if (mode.value === 'countup') remaining.value = 0
  store.updateTomato({ running: true, mode: mode.value, durationMinutes: workM.value, remainingSeconds: remaining.value })
  if (!timer) timer = setInterval(tick, 1000)
}
function pause() {
  running.value = false
  if (timer) {
    clearInterval(timer)
    timer = null
  }
  store.updateTomato({ running: false })
}
function togglePauseResume() {
  if (!started.value) return
  if (running.value) {
    pause()
  } else {
    running.value = true
    if (!timer) timer = setInterval(tick, 1000)
    store.updateTomato({ running: true })
  }
}
function reset() {
  pause()
  phase.value = 'work'
  remaining.value = workM.value * 60
  store.updateTomato({ remainingSeconds: remaining.value })
}
function complete() {
  // 中文注释：手动完成，触发上报（正计时秒数 / 倒计时换算秒数）
  const usedSeconds = mode.value === 'countdown' ? Math.max(0, workM.value * 60 - remaining.value) : remaining.value
  stopInternal()
  emit('complete', usedSeconds)
}
function stopInternal() {
  running.value = false
  if (timer) {
    clearInterval(timer)
    timer = null
  }
  store.updateTomato({ running: false })
}
function setDuration(min: number) {
  customMinutes.value = min
  applyCustom()
}
function applyCustom() {
  // 中文注释：更新工作时长与剩余时间
  const m = customMinutes.value
  const wasRunning = running.value
  stopInternal()
  remaining.value = mode.value === 'countdown' ? m * 60 : 0
  store.updateTomato({ durationMinutes: m, remainingSeconds: remaining.value })
  if (wasRunning) start()
}
// 中文注释：关闭按钮已移除，使用对话框右上角默认关闭控件

watch(phase, () => {
  // 中文注释：切换阶段时重置运行状态，避免自动继续
  pause()
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
/* 中文注释：使用基础文本样式与按钮布局，视觉轻量 */
</style>

