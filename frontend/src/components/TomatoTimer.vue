<template>
  <!-- 中文注释：番茄钟组件（增强版），支持倒计时/正计时、预设与自定义、开始/暂停/继续/完成 -->
  <!-- 中文注释：容器改为充满视口高度，底部区域更贴近页面底部 -->
  <div class="relative p-4 flex flex-col justify-between">

    <!-- 中文注释：移除组件内部标题，避免与弹窗标题重复显示 -->

    <!-- 中文注释：时间居中，右侧提供模式切换图标；备注在时间下方显示为灰色小字 -->
    <!-- 中文注释：夜间主题 - 时间与图标颜色调整为浅色（#B8CEE8），使对比清晰 -->
    <!-- 中文注释：中部区域使用 flex-1 居中，让时间显示位于顶部与底部之间的正中 -->
    <div class="flex-1 flex flex-col items-center justify-center gap-3 pt-16 pb-32">
      <!-- 中文注释：表盘样式 + 数字时间：浅橙色进度圈在倒计时模式下逐渐减少 -->
      <div class="relative w-64 h-64">
        <svg width="256" height="256" viewBox="0 0 256 256">
          <!-- 背景轨道 -->
          <circle cx="128" cy="128" r="100" stroke="#4a4a48" stroke-width="16" fill="none" />
          <!-- 钟表刻度：60 个刻度，5 的倍数加粗加长 -->
          <g stroke="#B8CEE8" opacity="0.7">
            <template v-for="i in 60" :key="i">
              <line
                x1="128"
                y1="28"
                x2="128"
                :y2="((i-1) % 5 === 0) ? 44 : 36"
                :stroke-width="((i-1) % 5 === 0) ? 2.5 : 1.5"
                :transform="'rotate(' + ((i-1) * 6) + ' 128 128)'"
              />
            </template>
          </g>
          <!-- 进度弧线：倒计时剩余比例，起点在上方（-90°旋转） -->
          <circle
            cx="128" cy="128" r="100"
            stroke="#F4A261" stroke-width="16" fill="none" stroke-linecap="butt"
            :style="{ transition: 'stroke-dashoffset .3s linear' }"
            :stroke-dasharray="circumference"
            :stroke-dashoffset="dashOffset"
            transform="rotate(-90 128 128)"
          />
          <!-- 进度环分钟分割线（白色）：每 1 分钟一条，覆盖在进度环之上，使橙色环呈间断 -->
          <g stroke="#FFFFFF" opacity="0.85">
            <template v-for="i in 60" :key="'sep-'+i">
              <line
                x1="128"
                y1="20"
                x2="128"
                y2="36"
                stroke-width="3"
                :transform="'rotate(' + ((i-1) * 6) + ' 128 128)'"
              />
            </template>
          </g>
          <!-- 数字刻度（0、5、10...55） -->
          <g fill="#B8CEE8" font-size="12" opacity="0.95" font-family="monospace" font-weight="700">
            <template v-for="lbl in dialLabels" :key="lbl.m">
              <text :x="lbl.x" :y="lbl.y" text-anchor="middle" dominant-baseline="middle">{{ lbl.m }}</text>
            </template>
          </g>
        </svg>
        <!-- 数字时间置于表盘中心 -->
        <!-- 中文注释：中心时间区域改为上下结构；上方笑脸、下方红色番茄图标 -->
        <div class="absolute inset-0 flex items-center justify-center">
          <div class="flex flex-col items-center justify-center">
            <div class="text-xl select-none" :style="{ color: '#B8CEE8' }">😊</div>
            <div class="text-6xl font-mono" :style="{ color: '#B8CEE8' }">{{ mm }}:{{ ss }}</div>
            <img src="@/assets/tomato.png" alt="番茄" class="w-6 h-6 mt-2 select-none" />
          </div>
        </div>
      </div>
      <!-- 模式切换图标置于时间下方，避免与中心重叠 -->
      <div>
        <el-icon
          class="cursor-pointer"
          :title="mode==='countdown' ? '切换为正计时' : '切换为倒计时'"
          :style="{ color: '#B8CEE8' }"
          @click="toggleMode"
        >
          <RefreshRight />
        </el-icon>
      </div>
      <div class="text-center text-xs" :style="{ color: '#9aa8b8' }" v-if="props.taskRemark">备注：{{ props.taskRemark }}</div>
      <!-- 中文注释：移除预计时长的小字提示，界面更简洁 -->
    </div>

    <!-- 中文注释：底部区域 - 包含时间预设/自定义与控制按钮，更贴近页面底部 -->
    <div class="fixed bottom-4 left-0 right-0 space-y-3" data-bottom="tomato-controls">
      <!-- 中文注释：倒计时模式下显示预设与自定义；正计时不显示这些 -->
      <div class="flex items-center justify-center gap-3" v-if="mode==='countdown'">
      <div class="flex items-center gap-2">
        <el-tag class="cursor-pointer" @click="setDuration(10)" :style="nightTagStyle">10分钟</el-tag>
        <el-tag class="cursor-pointer" @click="setDuration(20)" :style="nightTagStyle">20分钟</el-tag>
      </div>
      <!-- 中文注释：步进改为5分钟，点击 + / - 按5分钟增减 -->
      <el-input-number v-model="customMinutes" :min="1" :max="240" :step="5" />
      <el-button size="small" @click="applyCustom" :style="nightBtnStyle">应用</el-button>
    </div>

    <!-- 控制按钮 -->
    <div class="flex justify-center gap-2 mt-2">
      <el-button @click="start" :disabled="running" :style="nightBtnStyle">开始</el-button>
      <el-button @click="togglePauseResume" :disabled="!started" :style="nightBtnStyle">{{ running ? '暂停' : '继续' }}</el-button>
      <el-button @click="reset" :style="nightBtnStyle">重置</el-button>
      <el-button @click="complete" :style="nightBtnStyle">完成</el-button>
    </div>

    </div>
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
// 中文注释：夜间主题按钮与 Tag 样式（与深色背景协调）
const nightBtnStyle = { backgroundColor: '#3a3a38', color: '#B8CEE8', borderColor: '#4a4a48' }
const nightTagStyle = { backgroundColor: '#3a3a38', color: '#B8CEE8', borderColor: '#4a4a48' }
const running = ref(store.tomato.running)
const remaining = ref(store.tomato.remainingSeconds || (mode.value === 'countdown' ? workM.value * 60 : 0))
const started = ref(false)
const customMinutes = ref(workM.value)
let timer: any = null

// 中文注释：移除未使用的文本计算，避免无用警告
const mm = computed(() => String(Math.floor(remaining.value / 60)).padStart(2, '0'))
const ss = computed(() => String(remaining.value % 60).padStart(2, '0'))
// 中文注释：表盘进度按 60 分钟满圈映射。倒计时显示剩余与 60 的比例，正计时显示剩余（计划-已用）与 60 的比例。
const dialRatio = computed(() => {
  if (mode.value === 'countdown') {
    // 倒计时：显示剩余相对 60 分钟的比例（逐步减少）
    return Math.min(1, Math.max(0, remaining.value / 3600))
  }
  // 正计时：显示已用时间相对 60 分钟的比例（从 0 逐步增加）
  return Math.min(1, Math.max(0, remaining.value / 3600))
})
const circumference = 2 * Math.PI * 100
const dashOffset = computed(() => circumference * (1 - dialRatio.value))
// 中文注释：数字刻度位置（0、5、10...55），放到表盘外侧并加粗显示
const dialLabels = computed(() => {
  const minutes = [0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55]
  const r = 120 // 外侧半径：略大于进度环外缘（100 + 6 + margin）
  return minutes.map((m) => {
    const rad = (m * 6 - 90) * Math.PI / 180
    return {
      m,
      x: 128 + r * Math.cos(rad),
      y: 128 + r * Math.sin(rad),
    }
  })
})

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
      // 中文注释：工作阶段倒计时结束，按计划时长（秒）上报
      emit('complete', workM.value * 60)
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
      // 中文注释：正计时完成时按累计秒数上报
      emit('complete', remaining.value)
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
/* 中文注释：夜间主题 - 统一数字输入控件的底纹与文字颜色，与按钮风格一致 */
:deep(.el-input__wrapper) {
  background-color: #3a3a38 !important; /* rgb(58, 58, 56) */
  color: #B8CEE8 !important;
  box-shadow: 0 0 0 1px #4a4a48 inset !important; /* Element Plus 使用 box-shadow 模拟边框 */
}
:deep(.el-input__inner) {
  color: #B8CEE8 !important;
}
:deep(.el-input-number__decrease),
:deep(.el-input-number__increase) {
  color: #B8CEE8 !important;
  background-color: #3a3a38 !important;
  border-color: #4a4a48 !important;
}
:deep(.el-input-number__decrease:hover),
:deep(.el-input-number__increase:hover) {
  background-color: #3a3a38 !important;
}
</style>

