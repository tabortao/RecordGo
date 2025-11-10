<template>
  <!-- 中文注释：自定义小喇叭图标（SVG），与页面图标风格统一 -->
  <svg
    class="speaker-icon"
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    :width="size"
    :height="size"
    fill="none"
    stroke="currentColor"
    stroke-width="1.8"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-label="朗读"
  >
    <!-- 喇叭主体 -->
    <path d="M4 10 L9 7 V17 L4 14 Z" />
    <!-- 手持把手/圆角 -->
    <path d="M9 8 L14 8" />
    <!-- 声波 -->
    <path d="M15 9 C17 10 17 14 15 15" />
    <path d="M17 7 C20 9 20 15 17 17" />
  </svg>
</template>

<script setup lang="ts">
// 中文注释：支持通过 prop 调整图标尺寸
defineProps<{ size?: number }>()
</script>

<style scoped>
.speaker-icon {
  display: inline-block;
  vertical-align: middle;
}
</style>