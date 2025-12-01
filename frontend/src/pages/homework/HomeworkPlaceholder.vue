<template>
  <div class="h-full flex flex-col bg-gray-50 dark:bg-gray-900">
    <div class="bg-white dark:bg-gray-800 shadow-sm px-4 py-3 flex items-center gap-2 sticky top-0 z-10">
      <el-icon :size="20" class="text-gray-600 dark:text-gray-300 cursor-pointer" @click="router.back()"><ArrowLeft /></el-icon>
      <span class="text-lg font-bold text-gray-800 dark:text-white">{{ title }}</span>
    </div>

    <div class="flex-1 flex flex-col items-center justify-center p-8 text-center text-gray-500 dark:text-gray-400">
      <el-icon :size="64" class="mb-4 text-gray-300 dark:text-gray-600"><Tools /></el-icon>
      <h2 class="text-xl font-medium mb-2">功能开发中</h2>
      <p>请稍后，更多功能正在开发中……</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Tools } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const title = computed(() => (route.query.title as string) || '作业家')
</script>
