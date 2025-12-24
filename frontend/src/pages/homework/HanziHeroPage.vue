<template>
  <!-- 页面容器，填满父容器（RouterView），保留底部导航 -->
  <div class="flex flex-col h-full bg-gray-50 dark:bg-gray-900">
    <!-- 顶部导航栏 -->
    <div class="bg-white dark:bg-gray-800 shadow-sm px-4 py-3 flex items-center gap-2 shrink-0">
      <el-icon :size="20" class="text-gray-600 dark:text-gray-300 cursor-pointer" @click="router.back()"><ArrowLeft /></el-icon>
      <span class="text-lg font-bold text-gray-800 dark:text-white">汉字英雄</span>
      <div class="flex-1"></div>
      <el-tag v-if="loading" type="info" size="small" class="animate-pulse">加载中...</el-tag>
      <el-tag v-else-if="error" type="danger" size="small">加载失败</el-tag>
      <el-tag v-else type="success" size="small">运行中</el-tag>
    </div>

    <!-- WebView (Iframe) 容器 -->
    <div class="flex-1 relative w-full overflow-hidden">
      <!-- 加载状态 -->
      <div v-if="loading" class="absolute inset-0 flex flex-col items-center justify-center bg-gray-50 dark:bg-gray-900 z-20">
        <el-icon class="is-loading text-blue-500 mb-2" :size="32"><Loading /></el-icon>
        <span class="text-gray-500 text-sm">正在进入汉字英雄...</span>
      </div>

      <!-- 错误状态 -->
      <div v-if="error" class="absolute inset-0 flex flex-col items-center justify-center bg-gray-50 dark:bg-gray-900 z-20">
        <el-icon class="text-red-500 mb-2" :size="48"><Warning /></el-icon>
        <span class="text-gray-600 font-medium mb-4">页面加载失败</span>
        <el-button type="primary" @click="reload">重新加载</el-button>
      </div>

      <iframe 
        ref="iframeRef"
        :src="url" 
        class="w-full h-full border-none"
        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
        allowfullscreen
        @load="onLoad"
        @error="onError"
      ></iframe>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Loading, Warning } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const url = 'https://hanzihero.sdgarden.top'
const loading = ref(true)
const error = ref(false)
const iframeRef = ref<HTMLIFrameElement | null>(null)

function onLoad() {
  loading.value = false
  error.value = false
}

function onError() {
  loading.value = false
  error.value = true
  ElMessage.error('汉字英雄加载失败，请检查网络连接')
}

function reload() {
  loading.value = true
  error.value = false
  if (iframeRef.value) {
    const currentSrc = iframeRef.value.src
    iframeRef.value.src = ''
    setTimeout(() => {
      if (iframeRef.value) iframeRef.value.src = currentSrc
    }, 100)
  }
}
</script>

<style scoped>
/* 确保全屏容器覆盖底部导航栏 */
.fixed.inset-0 {
  bottom: 0 !important;
  left: 0 !important;
  right: 0 !important;
  top: 0 !important;
}
</style>
