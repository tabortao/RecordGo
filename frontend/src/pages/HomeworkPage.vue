<template>
  <div class="fixed inset-0 flex flex-col overflow-hidden bg-gray-50 dark:bg-gray-950 transition-colors duration-300">
    <div class="shrink-0 px-4 pt-4">
      <div class="rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/75 dark:bg-gray-900/65 backdrop-blur-xl shadow-sm">
        <div class="flex items-center justify-between gap-3 px-3 py-3">
          <div class="flex items-center gap-3 min-w-0">
            <div class="h-10 w-10 rounded-2xl border border-indigo-200/70 dark:border-indigo-900/40 bg-indigo-50/80 dark:bg-indigo-900/25 flex items-center justify-center text-indigo-700 dark:text-indigo-300">
              <el-icon :size="18"><Reading /></el-icon>
            </div>
            <div class="min-w-0">
              <div class="text-[17px] font-extrabold tracking-tight text-gray-900 dark:text-gray-50 truncate">作业家</div>
              <div class="mt-0.5 text-xs text-gray-500 dark:text-gray-400 truncate">语文 · 数学 · 英语 · 科学 · 课表</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto px-4 pt-5 space-y-4" style="padding-bottom: calc(env(safe-area-inset-bottom) + 96px)">
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
  { id: 'timetable', name: '课表', icon: '📅', path: '/timetable', bgClass: 'bg-orange-100 text-orange-600' },
  { id: 'dictation', name: '听写大师', icon: '🎧', path: '/dictation', bgClass: 'bg-purple-100 text-purple-600' },
  { id: 'growth', name: '小成长', icon: '🌱', path: '/little-growth', bgClass: 'bg-green-100 text-green-600' },
]

function onSubjectClick(subject: { id: string; path: string }) {
  router.push(subject.path)
}
</script>
