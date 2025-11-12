<template>
  <!-- 中文注释：悬浮番茄钟球，固定在右下角，显示 mm:ss，点击返回到独立番茄钟页面 -->
  <div
    class="fixed z-50"
    :style="{ left: posLeft + 'px', top: posTop + 'px' }"
  >
    <button
      ref="ballRef"
      class="w-14 h-14 rounded-full shadow-lg flex items-center justify-center font-mono text-sm select-none"
      :style="{ backgroundColor: '#F4A261', color: '#1F2937' }"
      @click="openTomatoPage"
      @pointerdown="onPointerDown"
    >
      {{ mm }}:{{ ss }}
    </button>
  </div>
  
  <!-- 可选：轻微提示文字 -->
  <!-- <div class="fixed bottom-8 right-4 text-xs" style="color:#94a3b8">点击回到番茄钟</div> -->
</template>

<script setup lang="ts">
// 中文注释：悬浮球组件在页面间继续计时（非固定模式），并在倒计时结束时写入完成记录、自动停止与隐藏
import { computed, onMounted, onUnmounted, ref } from 'vue'
import router from '@/router'
import { useAppState } from '@/stores/appState'
import { completeTomato, updateTaskStatus } from '@/services/tasks'
import { ElMessage } from 'element-plus'

const store = useAppState()
let t: any = null
let finishedOnce = false
const ballRef = ref<HTMLElement | null>(null)
const posLeft = ref(0)
const posTop = ref(0)
const BALL_SIZE = 56 // w-14 h-14
const MARGIN = 8
let dragging = false
let dragOffsetX = 0
let dragOffsetY = 0

const mm = computed(() => {
  const sec = store.tomato.remainingSeconds || 0
  const m = Math.floor(sec / 60)
  return String(m).padStart(2, '0')
})
const ss = computed(() => String((store.tomato.remainingSeconds || 0) % 60).padStart(2, '0'))

async function finalizeCountdown() {
  const id = store.tomato.currentTaskId
  const minutes = store.tomato.durationMinutes || 20
  try {
    if (id && !isNaN(id)) {
      await completeTomato(id, minutes)
      // 中文注释：悬浮球计时完成也视为番茄完成，允许只读权限下标记“已完成”
      await updateTaskStatus(id, 2, { allowByTomato: true })
      ElMessage.success('番茄钟完成，数据已记录')
    }
  } catch (e: any) {
    ElMessage.error(`番茄上报失败：${e?.message || e}`)
  }
}

function tick() {
  // 中文注释：基于墙钟时间戳计算，避免锁屏/后台导致中断
  const now = Date.now()
  const mode = store.tomato.runningMode ?? store.tomato.mode
  if (mode === 'countdown') {
    const endAt = store.tomato.endAtMs ?? null
    if (endAt == null) return
    const next = Math.max(0, Math.round((endAt - now) / 1000))
    store.updateTomato({ remainingSeconds: next })
    if (next <= 0 && !finishedOnce) {
      finishedOnce = true
      finalizeCountdown().finally(() => {
        store.updateTomato({ running: false, showFloating: false, startAtMs: null, endAtMs: null, runningMode: null })
        if (t) { clearInterval(t); t = null }
      })
    }
  } else {
    const startAt = store.tomato.startAtMs ?? null
    if (startAt == null) return
    const next = Math.max(0, Math.round((now - startAt) / 1000))
    store.updateTomato({ remainingSeconds: next })
  }
}

function openTomatoPage() {
  const id = store.tomato.currentTaskId
  if (id && !isNaN(id)) {
    // 中文注释：进入番茄钟页面时隐藏悬浮球，避免同时显示
    store.updateTomato({ showFloating: false })
    router.push({ name: 'TaskTomato', params: { id } })
  } else {
    router.push({ name: 'Tasks' })
  }
}

onMounted(() => {
  if (!t) t = setInterval(tick, 1000)
  const onVis = () => { if (document.visibilityState === 'visible' && store.tomato.running) tick() }
  document.addEventListener('visibilitychange', onVis)
  ;(window as any).__floating_vis__ = onVis
  // 初始位置：底部居中
  const w = window.innerWidth
  const h = window.innerHeight
  const saved = localStorage.getItem('floatingTomatoPos')
  if (saved) {
    try {
      const p = JSON.parse(saved) as { left: number; top: number }
      posLeft.value = Math.min(Math.max(MARGIN, p.left), w - BALL_SIZE - MARGIN)
      posTop.value = Math.min(Math.max(MARGIN, p.top), h - BALL_SIZE - MARGIN)
    } catch {
      posLeft.value = Math.floor(w / 2 - BALL_SIZE / 2)
      posTop.value = h - BALL_SIZE - 64
    }
  } else {
    posLeft.value = Math.floor(w / 2 - BALL_SIZE / 2)
    posTop.value = h - BALL_SIZE - 64 // 约等于 bottom-16
  }
})
onUnmounted(() => {
  if (t) { clearInterval(t); t = null }
  const onVis = (window as any).__floating_vis__
  if (onVis) document.removeEventListener('visibilitychange', onVis)
})

function onPointerDown(e: PointerEvent) {
  dragging = true
  const rect = ballRef.value?.getBoundingClientRect()
  dragOffsetX = e.clientX - (rect?.left || 0)
  dragOffsetY = e.clientY - (rect?.top || 0)
  window.addEventListener('pointermove', onPointerMove)
  window.addEventListener('pointerup', onPointerUp)
}
function onPointerMove(e: PointerEvent) {
  if (!dragging) return
  const w = window.innerWidth
  const h = window.innerHeight
  const nextLeft = e.clientX - dragOffsetX
  const nextTop = e.clientY - dragOffsetY
  posLeft.value = Math.min(Math.max(MARGIN, nextLeft), w - BALL_SIZE - MARGIN)
  posTop.value = Math.min(Math.max(MARGIN, nextTop), h - BALL_SIZE - MARGIN)
}
function onPointerUp() {
  dragging = false
  window.removeEventListener('pointermove', onPointerMove)
  window.removeEventListener('pointerup', onPointerUp)
  localStorage.setItem('floatingTomatoPos', JSON.stringify({ left: posLeft.value, top: posTop.value }))
}
</script>

<style scoped>
</style>