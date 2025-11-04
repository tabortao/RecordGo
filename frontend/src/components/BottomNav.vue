<template>
  <!-- 中文注释：底部导航，当前页面高亮 -->
  <nav class="h-14 border-t bg-white flex items-center justify-around fixed bottom-0 left-0 right-0 z-50">
    <router-link
      v-for="item in items"
      :key="item.to"
      :to="item.to"
      class="flex flex-col items-center justify-center w-24 text-gray-600"
      :class="{ 'text-primary font-semibold': $route.path === item.to }"
    >
      <el-icon :size="20"><component :is="item.icon" /></el-icon>
      <span class="text-xs mt-1">{{ item.label }}</span>
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
