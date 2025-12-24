<template>
  <!-- 全屏覆盖容器，固定定位确保不被父容器裁剪，z-index 设为 100 覆盖底部导航 -->
  <!-- 添加 !pb-0 覆盖 App.vue 中 router-view 传递下来的 pb-20 -->
  <div class="fixed inset-0 bg-gray-50 dark:bg-gray-900 !pb-0" style="z-index: 100;">
    <!-- 悬浮返回按钮 -->
    <div class="absolute top-4 left-4 z-50">
      <div 
        class="bg-black/20 text-white backdrop-blur rounded-full p-2 cursor-pointer hover:bg-black/40 transition-all flex items-center justify-center"
        @click="router.back()"
      >
        <el-icon :size="20"><ArrowLeft /></el-icon>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="absolute inset-0 flex flex-col items-center justify-center bg-gray-50 dark:bg-gray-900 z-40">
      <el-icon class="is-loading text-blue-500 mb-2" :size="32"><Loading /></el-icon>
      <span class="text-gray-500 text-sm">正在进入汉字英雄...</span>
    </div>

    <!-- 错误状态 -->
    <div v-if="error" class="absolute inset-0 flex flex-col items-center justify-center bg-gray-50 dark:bg-gray-900 z-40">
      <el-icon class="text-red-500 mb-2" :size="48"><Warning /></el-icon>
      <span class="text-gray-600 font-medium mb-4">页面加载失败</span>
      <el-button type="primary" @click="reload">重新加载</el-button>
    </div>

    <!-- WebView (Iframe) -->
    <iframe 
      ref="iframeRef"
      :src="url" 
      class="w-full h-full border-none block"
      allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
      allowfullscreen
      @load="onLoad"
      @error="onError"
    ></iframe>
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
