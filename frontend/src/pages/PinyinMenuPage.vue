<template>
  <div class="p-4 space-y-4 pb-24">
    <div class="flex items-center gap-2 mb-6">
      <el-icon :size="20" class="cursor-pointer text-gray-600" @click="router.back()"><ArrowLeft /></el-icon>
      <h1 class="text-xl font-bold text-gray-800 dark:text-white">拼音王国</h1>
    </div>

    <div class="grid grid-cols-2 sm:grid-cols-3 gap-4">
      <div 
        class="bg-green-50 dark:bg-green-900/20 rounded-xl p-4 shadow-sm border border-green-100 dark:border-green-800 flex flex-col items-center justify-center gap-3 cursor-pointer hover:shadow-md transition active:scale-95 relative overflow-hidden group"
        @click="router.push('/homework/chinese/pinyin/wheel')"
      >
        <div class="absolute -right-4 -top-4 w-16 h-16 bg-green-200 rounded-full opacity-20 group-hover:scale-150 transition duration-500"></div>
        <div class="w-16 h-16 rounded-full bg-white shadow-sm flex items-center justify-center text-3xl">
          🎡
        </div>
        <span class="font-bold text-green-700 dark:text-green-400">拼音大转盘</span>
        <span class="text-xs text-green-600/70">互动游戏</span>
      </div>

      <!-- 占位符 -->
      <div class="bg-gray-50 dark:bg-gray-800/50 rounded-xl p-4 border border-dashed border-gray-300 dark:border-gray-700 flex flex-col items-center justify-center gap-2 opacity-60">
        <span class="text-3xl grayscale">🔤</span>
        <span class="text-sm text-gray-500">声母表 (待开发)</span>
      </div>
       <div class="bg-gray-50 dark:bg-gray-800/50 rounded-xl p-4 border border-dashed border-gray-300 dark:border-gray-700 flex flex-col items-center justify-center gap-2 opacity-60">
        <span class="text-3xl grayscale">🔊</span>
        <span class="text-sm text-gray-500">韵母表 (待开发)</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'

const router = useRouter()
</script>
