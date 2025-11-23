<template>
  <!-- 中文注释：底部导航，当前页面高亮 -->
  <nav class="h-20 border-t border-gray-200 dark:border-gray-700 bg-white/90 dark:bg-gray-800/90 backdrop-blur flex items-center justify-around fixed bottom-0 left-0 right-0 z-50">
    <router-link
      v-for="item in items"
      :key="item.to"
      :to="item.to"
      class="flex flex-col items-center justify-center w-24 pt-3"
      :class="$route.path === item.to ? 'text-emerald-600 dark:text-emerald-400 font-semibold' : 'text-gray-600 dark:text-gray-300 opacity-80'"
    >
      <el-icon :size="20"><component :is="item.icon" /></el-icon>
      <span class="text-xs mt-1 mb-6">{{ item.label }}</span>
    </router-link>
  </nav>
</template>

<script setup lang="ts">
import { List, Star, User } from '@element-plus/icons-vue'

// 中文注释：底部导航三页：任务/心愿/我的
const items = [
  { to: '/tasks', label: '任务', icon: List },
  { to: '/wishes', label: '心愿', icon: Star },
  { to: '/mine', label: '我的', icon: User }
]
</script>

<style scoped>
</style>
