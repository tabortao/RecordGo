<template>
  <transition name="fade">
    <div v-if="store.pageLoading" class="fixed inset-0 z-50 pointer-events-none">
      <div class="fixed top-0 left-0 right-0 h-1 bg-transparent">
        <div class="h-1" :style="{ width: store.pageProgress + '%', backgroundColor: barColor }"></div>
      </div>
      <div class="absolute inset-0 flex items-center justify-center">
        <div class="px-4 py-3 rounded-xl shadow-lg bg-white/80 dark:bg-gray-800/80 backdrop-blur">
          <div class="flex items-center gap-2">
            <el-icon :size="18" class="text-indigo-600 dark:text-indigo-400"><Loading /></el-icon>
            <span class="text-sm text-gray-700 dark:text-gray-100">页面加载中…</span>
            <el-button size="small" class="pointer-events-auto" @click="refresh">刷新</el-button>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { Loading } from '@element-plus/icons-vue'
import { useAppState } from '@/stores/appState'
import router from '@/router'

const store = useAppState()
const barColor = computed(() => document.documentElement.classList.contains('dark') ? '#60a5fa' : '#2563eb')

let timeoutId: any = null
function ensureTimeout() {
  clearTimeout(timeoutId)
  timeoutId = setTimeout(() => {
    if (store.pageLoading) {
      store.setPageProgress(Math.min(95, store.pageProgress + 10))
    }
  }, 8000)
}
function refresh() { try { router.go(0) } catch { location.reload() } }
onMounted(ensureTimeout)
onUnmounted(() => { clearTimeout(timeoutId) })
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity .2s }
.fade-enter-from, .fade-leave-to { opacity: 0 }
</style>

