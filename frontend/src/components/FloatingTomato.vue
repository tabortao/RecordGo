<template>
  <!-- 中文注释：悬浮番茄钟球，固定在右下角，显示 mm:ss，点击返回到独立番茄钟页面 -->
  <div class="fixed bottom-20 right-4 z-50">
    <button
      class="w-14 h-14 rounded-full shadow-lg flex items-center justify-center font-mono text-sm"
      :style="{ backgroundColor: '#F4A261', color: '#1F2937' }"
      @click="openTomatoPage"
    >
      {{ mm }}:{{ ss }}
    </button>
  </div>
  
  <!-- 可选：轻微提示文字 -->
  <!-- <div class="fixed bottom-8 right-4 text-xs" style="color:#94a3b8">点击回到番茄钟</div> -->
</template>

<script setup lang="ts">
// 中文注释：悬浮球组件在页面间继续计时（非固定模式），并在倒计时结束时写入完成记录、自动停止与隐藏
import { computed, onMounted, onUnmounted } from 'vue'
import router from '@/router'
import { useAppState } from '@/stores/appState'
import { completeTomato, updateTaskStatus } from '@/services/tasks'
import { ElMessage } from 'element-plus'

const store = useAppState()
let t: any = null
let finishedOnce = false

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
      await updateTaskStatus(id, 2)
      ElMessage.success('番茄钟完成，数据已记录')
    }
  } catch (e: any) {
    ElMessage.error(`番茄上报失败：${e?.message || e}`)
  }
}

function tick() {
  const mode = store.tomato.mode
  const cur = store.tomato.remainingSeconds || 0
  if (mode === 'countdown') {
    const next = Math.max(0, cur - 1)
    store.updateTomato({ remainingSeconds: next })
    if (next === 0 && !finishedOnce) {
      finishedOnce = true
      // 写入完成记录并关闭悬浮球
      finalizeCountdown().finally(() => {
        store.updateTomato({ running: false, showFloating: false })
        if (t) { clearInterval(t); t = null }
      })
    }
  } else {
    store.updateTomato({ remainingSeconds: cur + 1 })
  }
}

function openTomatoPage() {
  const id = store.tomato.currentTaskId
  if (id && !isNaN(id)) {
    router.push({ name: 'TaskTomato', params: { id } })
  } else {
    router.push({ name: 'Tasks' })
  }
}

onMounted(() => {
  if (!t) t = setInterval(tick, 1000)
})
onUnmounted(() => { if (t) { clearInterval(t); t = null } })
</script>

<style scoped>
</style>