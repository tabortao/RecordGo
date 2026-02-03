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
      <!-- 中文注释：移除预计时长的小字提示，界面更简洁 -->
    </div>

    <!-- 中文注释：底部区域 - 包含时间预设/自定义与控制按钮，更贴近页面底部 -->
    <div class="fixed bottom-12 left-0 right-0 flex items-center justify-center" data-bottom="tomato-controls">
      <button
        class="w-20 h-20 rounded-full shadow-lg font-bold select-none text-[16px]"
        :style="{ backgroundColor: '#3a3a38', color: '#B8CEE8', boxShadow: '0 0 0 1px #4a4a48 inset' }"
        @click="onCircleTap"
        @pointerdown="onCircleDown"
        @pointerup="onCircleUp"
        @pointerleave="onCircleUp"
      >
        {{ running ? '暂停' : (started ? '继续' : '开始') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
// 中文注释：接收任务名与备注 + 工作/休息分钟数（合并定义）
// 中文注释：番茄钟逻辑，工作与休息两个阶段，倒计时结束后切换阶段并在工作结束时触发 complete 事件
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useAppState } from '@/stores/appState'
import { speak } from '@/utils/speech'
import JingleBellUrl from '@/assets/audio/JingleBell.mp3'
import { RefreshRight } from '@element-plus/icons-vue'

// 中文注释：新增 taskId，便于悬浮球返回到独立番茄钟页面
const props = defineProps<{ workMinutes?: number; breakMinutes?: number; taskName?: string; taskRemark?: string; taskId?: number }>()
const emit = defineEmits<{ (e: 'complete', seconds: number): void }>()
const store = useAppState()

function playBell(): Promise<void> {
  return new Promise((resolve) => {
    try {
      const audio = new Audio(JingleBellUrl)
      let done = false
      const finish = () => { if (!done) { done = true; resolve() } }
      audio.onended = finish
      audio.onerror = finish
      const p = audio.play()
      if (p && typeof (p as any).catch === 'function') (p as any).catch(finish)
    } catch {
      resolve()
    }
  })
}

async function playEndCue() {
  await playBell().catch(() => {})
  const text = (store.tomato.countdownEndText || '').trim()
  if (!store.tomato.countdownEndSpeakEnabled) return
  if (store.speech.enabled && text) {
    speak(text, { voiceURI: store.speech.voiceURI || undefined, rate: store.speech.rate, pitch: store.speech.pitch }).catch(() => {})
  }
}

const workM = computed(() => props.workMinutes ?? store.tomato.durationMinutes ?? 20)
const breakM = computed(() => props.breakMinutes ?? 5)

type Phase = 'work' | 'break'
const phase = ref<Phase>('work')
const mode = ref<'countdown' | 'countup'>(store.tomato.runningMode ?? store.tomato.mode)
const running = ref(store.tomato.running)
const remaining = ref(store.tomato.running ? (store.tomato.remainingSeconds || (mode.value === 'countdown' ? workM.value * 60 : 0)) : (mode.value === 'countdown' ? workM.value * 60 : 0))
const started = ref(false)
// 自定义分钟已迁移到页面级调整
let pressTimer: any = null
function onCircleTap() { if (!started.value) start(); else togglePauseResume() }
function onCircleDown() { if (pressTimer) clearTimeout(pressTimer); pressTimer = setTimeout(() => { complete() }, 2000) }
function onCircleUp() { if (pressTimer) { clearTimeout(pressTimer); pressTimer = null } }
let timer: any = null
// 中文注释：墙钟时间戳，支持锁屏/后台后继续准确计时
const startAtMs = ref<number | null>(store.tomato.startAtMs ?? null)
const endAtMs = ref<number | null>(store.tomato.endAtMs ?? null)
// 中文注释：Wake Lock API 哨兵，用于保持设备常亮（移动端）；iOS Chrome 回退为隐藏视频循环播放
let wakeLock: any = null
let keepAwakeVideo: HTMLVideoElement | null = null
function ensureVideoFallback() {
  try {
    if (keepAwakeVideo) { keepAwakeVideo.play().catch(() => {}); return }
    const video = document.createElement('video')
    video.setAttribute('playsinline', '')
    video.muted = true
    video.loop = true
    video.style.position = 'fixed'
    video.style.opacity = '0'
    video.style.width = '1px'
    video.style.height = '1px'
    video.style.pointerEvents = 'none'
    video.style.zIndex = '-1'
    const addSource = (type: string, dataURI: string) => {
      const s = document.createElement('source')
      s.src = dataURI
      s.type = `video/${type}`
      video.appendChild(s)
    }
    // 中文注释：内置极小的视频数据，来源于公开示例（用于保持唤醒）
    addSource('webm', 'data:video/webm;base64,GkXfo0AgQoaBAUL3gQFC8oEEQvOBCEKCQAR3ZWJtQoeBAkKFgQIYU4BnQI0VSalmQCgq17FAAw9CQE2AQAZ3aGFtbXlXQUAGd2hhbW15RIlACECPQAAAAAAAFlSua0AxrkAu14EBY8WBAZyBACK1nEADdW5khkAFVl9WUDglhohAA1ZQOIOBAeBABrCBCLqBCB9DtnVAIueBAKNAHIEAAIAwAQCdASoIAAgAAUAmJaQAA3AA/vz0AAA=')
    addSource('mp4', 'data:video/mp4;base64,AAAAHGZ0eXBpc29tAAACAGlzb21pc28ybXA0MQAAAAhmcmVlAAAAG21kYXQAAAGzABAHAAABthADAowdbb9/AAAC6W1vb3YAAABsbXZoZAAAAAB8JbCAfCWwgAAAA+gAAAAAAAEAAAEAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAIAAAIVdHJhawAAAFx0a2hkAAAAD3wlsIB8JbCAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAQAAAAAAAAAAAAAAAAAAQAAAAAAIAAAACAAAAAABsW1kaWEAAAAgbWRoZAAAAAB8JbCAfCWwgAAAA+gAAAAAVcQAAAAAAC1oZGxyAAAAAAAAAAB2aWRlAAAAAAAAAAAAAAAAVmlkZW9IYW5kbGVyAAAAAVxtaW5mAAAAFHZtaGQAAAABAAAAAAAAAAAAAAAkZGluZgAAABxkcmVmAAAAAAAAAAEAAAAMdXJsIAAAAAEAAAEcc3RibAAAALhzdHNkAAAAAAAAAAEAAACobXA0dgAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAIAAgASAAAAEgAAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABj//wAAAFJlc2RzAAAAAANEAAEABDwgEQAAAAADDUAAAAAABS0AAAGwAQAAAbWJEwAAAQAAAAEgAMSNiB9FAEQBFGMAAAGyTGF2YzUyLjg3LjQGAQIAAAAYc3R0cwAAAAAAAAABAAAAAQAAAAAAAAAcc3RzYwAAAAAAAAABAAAAAQAAAAEAAAABAAAAFHN0c3oAAAAAAAAAEwAAAAEAAAAUc3RjbwAAAAAAAAABAAAALAAAAGB1ZHRhAAAAWG1ldGEAAAAAAAAAIWhkbHIAAAAAAAAAAG1kaXJhcHBsAAAAAAAAAAAAAAAAK2lsc3QAAAAjqXRvbwAAABtkYXRhAAAAAQAAAABMYXZmNTIuNzguMw==')
    document.body.appendChild(video)
    keepAwakeVideo = video
    video.play().catch(() => {})
  } catch {}
}
async function requestWakeLock() {
  try {
    const nav: any = navigator as any
    if (nav.wakeLock && typeof nav.wakeLock.request === 'function') {
      wakeLock = await nav.wakeLock.request('screen')
      // 中文注释：若被系统释放，则在可见且仍在计时时尝试重新申请
      wakeLock.addEventListener?.('release', () => {
        wakeLock = null
        if (document.visibilityState === 'visible' && running.value) {
          requestWakeLock().catch(() => {})
        }
      })
    } else {
      // 中文注释：iOS Chrome 等不支持 Wake Lock，使用隐藏视频回退
      ensureVideoFallback()
    }
  } catch (e: any) {
    // 中文注释：Wake Lock 失败时启用视频回退方案
    ensureVideoFallback()
    ElMessage.warning('保持常亮失败，已启用兼容方案')
  }
}
async function releaseWakeLock() {
  try { if (wakeLock) { await wakeLock.release(); wakeLock = null } } catch {}
  try { if (keepAwakeVideo) { keepAwakeVideo.pause(); keepAwakeVideo.remove(); keepAwakeVideo = null } } catch {}
}

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
  store.updateTomato({ runningMode: mode.value, remainingSeconds: remaining.value })
  if (wasRunning) start()
}

function tick() {
  const now = Date.now()
  if (mode.value === 'countdown') {
    // 基于 endAtMs 计算剩余，避免页面不可见时中断
    const target = endAtMs.value ?? (startAtMs.value ? startAtMs.value + workM.value * 60 * 1000 : null)
    if (target == null) return
    remaining.value = Math.max(0, Math.round((target - now) / 1000))
    store.updateTomato({ remainingSeconds: remaining.value })
    if (remaining.value <= 0) {
      stopInternal()
      if (phase.value === 'work') {
        playEndCue().catch(() => {})
        emit('complete', workM.value * 60)
        phase.value = 'break'
        remaining.value = breakM.value * 60
      } else {
        phase.value = 'work'
        remaining.value = workM.value * 60
      }
    }
  } else {
    // 正计时：基于 startAtMs 计算累计秒数
    if (!startAtMs.value) return
    remaining.value = Math.max(0, Math.round((now - startAtMs.value) / 1000))
    store.updateTomato({ remainingSeconds: remaining.value })
    if (remaining.value >= workM.value * 60) {
      stopInternal()
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
  // 中文注释：记录墙钟时间戳，确保锁屏后继续计时
  const now = Date.now()
  if (mode.value === 'countdown') {
    startAtMs.value = now
    endAtMs.value = now + remaining.value * 1000
  } else {
    startAtMs.value = now
    endAtMs.value = null
  }
  store.updateTomato({ running: true, runningMode: mode.value, durationMinutes: workM.value, remainingSeconds: remaining.value, currentTaskId: props.taskId ?? null, startAtMs: startAtMs.value, endAtMs: endAtMs.value })
  if (!timer) timer = setInterval(tick, 1000)
  // 中文注释：开始后尝试申请常亮（移动端）
  if (store.tomato.keepAwakeEnabled) {
    requestWakeLock().catch(() => {})
  }
}
function pause() {
  running.value = false
  if (timer) {
    clearInterval(timer)
    timer = null
  }
  // 中文注释：暂停时保留当前 remainingSeconds，并清空时间戳
  startAtMs.value = null
  endAtMs.value = null
  store.updateTomato({ running: false, startAtMs: null, endAtMs: null, runningMode: null })
  // 中文注释：暂停时释放常亮
  releaseWakeLock().catch(() => {})
}
function togglePauseResume() {
  if (!started.value) return
  if (running.value) {
    pause()
  } else {
    running.value = true
    const now = Date.now()
    if (mode.value === 'countdown') {
      endAtMs.value = now + remaining.value * 1000
      startAtMs.value = now
    } else {
      // 正计时恢复：保持已用秒数的连续性
      startAtMs.value = now - remaining.value * 1000
      endAtMs.value = null
    }
    store.updateTomato({ running: true, startAtMs: startAtMs.value, endAtMs: endAtMs.value })
    if (!timer) timer = setInterval(tick, 1000)
  }
}
// 中文注释：主按钮动作（开始/继续）：未开始执行 start，已开始但暂停执行恢复
// 主按钮逻辑由大圆按钮事件统一处理
function reset() {
  pause()
  phase.value = 'work'
  remaining.value = workM.value * 60
  startAtMs.value = null
  endAtMs.value = null
  store.updateTomato({ remainingSeconds: remaining.value, startAtMs: null, endAtMs: null })
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
  startAtMs.value = null
  endAtMs.value = null
  store.updateTomato({ running: false, startAtMs: null, endAtMs: null, runningMode: null })
  // 中文注释：停止时释放常亮
  releaseWakeLock().catch(() => {})
}
// 时长调整由页面菜单处理
// 中文注释：关闭按钮已移除，使用对话框右上角默认关闭控件

watch(phase, () => {
  // 中文注释：切换阶段时重置运行状态，避免自动继续
  pause()
})

watch(workM, (m) => {
  if (running.value) return
  const sec = (mode.value === 'countdown' ? m * 60 : 0)
  remaining.value = sec
  store.updateTomato({ durationMinutes: m, remainingSeconds: sec })
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
// 中文注释：如果进入页面时番茄钟仍在运行（来自悬浮球），自动继续计时；同时监听可见性变化以刷新一次
onMounted(() => {
  if (store.tomato.running) {
    running.value = true
    started.value = true
    startAtMs.value = store.tomato.startAtMs ?? startAtMs.value
    endAtMs.value = store.tomato.endAtMs ?? endAtMs.value
    if (!timer) timer = setInterval(tick, 1000)
  }
  const onVis = () => {
    if (document.visibilityState === 'visible' && running.value) tick()
    // 中文注释：页面重新可见且在计时，则尝试重新申请常亮
    if (document.visibilityState === 'visible' && running.value && store.tomato.keepAwakeEnabled) {
      requestWakeLock().catch(() => {})
    }
  }
  document.addEventListener('visibilitychange', onVis)
  ;(window as any).__tomato_vis__ = onVis
})
onUnmounted(() => {
  const onVis = (window as any).__tomato_vis__
  if (onVis) document.removeEventListener('visibilitychange', onVis)
  // 中文注释：组件卸载时释放常亮
  releaseWakeLock().catch(() => {})
})
// 中文注释：监听“保持常亮”开关变化；运行中开启则申请，关闭则释放
watch(() => store.tomato.keepAwakeEnabled, (enabled) => {
  if (!running.value) return
  if (enabled) requestWakeLock().catch(() => {})
  else releaseWakeLock().catch(() => {})
})
// 中文注释：向父组件暴露停止/开始/暂停方法，便于页面“返回”时控制行为
defineExpose({ stop: stopInternal, start, pause, reset })
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

