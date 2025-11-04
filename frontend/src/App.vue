<template>
  <!-- 中文注释：顶层布局，包含页面内容与底部导航 -->
  <div class="min-h-screen flex flex-col bg-gray-50">
    <!-- 中文注释：为页面内容增加底部留白，避免被固定底部导航遮挡 -->
    <router-view class="flex-1 pb-16" />
    <BottomNav />
  </div>
  
  <!-- 中文注释：全局加载中占位（可根据需要改造） -->
  <el-backtop />
</template>

<script setup lang="ts">
// 中文注释：注册底部导航组件
import BottomNav from '@/components/BottomNav.vue'
</script>

<style scoped>
/* 中文注释：可根据品牌色进行优化 */
</style>
