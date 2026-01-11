<template>
  <div class="min-h-screen bg-[#FDF6F8] dark:bg-gray-900 pb-10 font-sans">
    <!-- Header -->
    <div class="sticky top-0 z-20 bg-white/80 dark:bg-gray-800/80 backdrop-blur-md border-b border-pink-100 dark:border-gray-700 flex items-center justify-between p-4 shadow-sm">
      <div class="flex items-center gap-3">
        <el-icon :size="22" class="cursor-pointer text-gray-600 dark:text-gray-300 hover:text-pink-500 transition" @click="router.back()"><ArrowLeft /></el-icon>
        <h1 class="text-xl font-bold text-gray-800 dark:text-white">语文</h1>
      </div>
    </div>

    <!-- Grid -->
    <div class="p-4 grid grid-cols-2 md:grid-cols-4 gap-4">
       <!-- Pinyin -->
       <div 
          class="bg-white dark:bg-gray-800 rounded-xl p-4 shadow-sm border border-gray-100 dark:border-gray-700 flex flex-col items-center justify-center gap-3 cursor-pointer hover:shadow-md transition active:scale-95"
          @click="router.push('/homework/chinese/pinyin')"
       >
          <div class="w-16 h-16 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center">
             <span class="text-3xl font-bold">a</span>
          </div>
          <span class="font-semibold text-gray-700 dark:text-gray-200">拼音</span>
       </div>

       <!-- Hanzi Hero -->
       <div 
          class="bg-white dark:bg-gray-800 rounded-xl p-4 shadow-sm border border-gray-100 dark:border-gray-700 flex flex-col items-center justify-center gap-3 cursor-pointer hover:shadow-md transition active:scale-95"
          @click="router.push('/homework/chinese/hanzi-hero')"
       >
          <div class="w-16 h-16 rounded-full bg-orange-100 text-orange-600 flex items-center justify-center overflow-hidden">
             <img src="@/assets/button/hanzihero.png" alt="Hanzi" class="w-full h-full object-cover opacity-90" />
          </div>
          <span class="font-semibold text-gray-700 dark:text-gray-200">汉字英雄</span>
       </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'

const router = useRouter()
</script>
