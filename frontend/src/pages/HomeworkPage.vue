<template>
  <div class="p-4 space-y-4 pb-24">
    <div class="flex items-center gap-2">
      <el-icon :size="24" class="text-indigo-600"><Reading /></el-icon>
      <h1 class="text-xl font-bold text-gray-800 dark:text-white">作业家</h1>
    </div>

    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div 
        v-for="subject in subjects" 
        :key="subject.id"
        class="bg-white dark:bg-gray-800 rounded-xl p-4 shadow-sm border border-gray-100 dark:border-gray-700 flex flex-col items-center justify-center gap-3 cursor-pointer hover:shadow-md transition active:scale-95"
        @click="onSubjectClick(subject)"
      >
        <div class="w-16 h-16 rounded-full flex items-center justify-center" :class="subject.bgClass">
          <span class="text-3xl">{{ subject.icon }}</span>
        </div>
        <span class="font-semibold text-gray-700 dark:text-gray-200">{{ subject.name }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Reading } from '@element-plus/icons-vue'

const router = useRouter()

const subjects = [
  { id: 'chinese', name: '语文', icon: '📖', path: '/homework/chinese', bgClass: 'bg-red-100 text-red-600' },
  { id: 'math', name: '数学', icon: '🔢', path: '/homework/placeholder?title=数学', bgClass: 'bg-blue-100 text-blue-600' },
  { id: 'english', name: '英语', icon: '🔤', path: '/homework/placeholder?title=英语', bgClass: 'bg-yellow-100 text-yellow-600' },
  { id: 'science', name: '科学', icon: '🔬', path: '/homework/placeholder?title=科学', bgClass: 'bg-green-100 text-green-600' },
  { id: 'dictation', name: '听写大师', icon: '🎧', path: '/dictation', bgClass: 'bg-purple-100 text-purple-600' },
  { id: 'growth', name: '小成长', icon: '🌱', path: '/little-growth', bgClass: 'bg-green-100 text-green-600' },
]

function onSubjectClick(subject: { id: string; path: string }) {
  router.push(subject.path)
}
</script>
