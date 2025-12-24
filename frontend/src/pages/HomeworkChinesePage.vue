<template>
  <div class="p-4 space-y-4 pb-24">
    <div class="flex items-center gap-2 mb-6">
      <el-icon :size="20" class="cursor-pointer text-gray-600" @click="router.back()"><ArrowLeft /></el-icon>
      <h1 class="text-xl font-bold text-gray-800 dark:text-white">语文</h1>
    </div>

    <div class="grid grid-cols-2 sm:grid-cols-3 gap-4">
      <div 
        class="bg-white dark:bg-gray-800 rounded-xl p-4 shadow-sm border border-gray-100 dark:border-gray-700 flex flex-col items-center justify-center gap-3 cursor-pointer hover:shadow-md transition active:scale-95"
        @click="router.push('/homework/chinese/pinyin')"
      >
        <div class="w-14 h-14 rounded-full bg-orange-100 flex items-center justify-center text-orange-600">
          <span class="text-2xl font-bold">aoe</span>
        </div>
        <span class="font-semibold text-gray-700 dark:text-gray-200">拼音</span>
      </div>

      <div 
        class="bg-white dark:bg-gray-800 rounded-xl p-4 shadow-sm border border-gray-100 dark:border-gray-700 flex flex-col items-center justify-center gap-3 cursor-pointer hover:shadow-md transition active:scale-95"
        @click="router.push('/homework/chinese/hanzi-hero')"
      >
        <div class="w-14 h-14 rounded-full flex items-center justify-center">
          <img src="@/assets/button/hanzihero.png" alt="汉字英雄" class="w-full h-full object-contain" />
        </div>
        <span class="font-semibold text-gray-700 dark:text-gray-200">汉字英雄</span>
      </div>

      <!-- 占位符 -->
      <div class="bg-gray-50 dark:bg-gray-800/50 rounded-xl p-4 border border-dashed border-gray-300 dark:border-gray-700 flex flex-col items-center justify-center gap-2 opacity-60">
        <div class="w-12 h-12 rounded-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
          <el-icon :size="20" class="text-gray-400"><More /></el-icon>
        </div>
        <span class="text-sm text-gray-500">更多功能开发中</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ArrowLeft, More } from '@element-plus/icons-vue'

const router = useRouter()
</script>
