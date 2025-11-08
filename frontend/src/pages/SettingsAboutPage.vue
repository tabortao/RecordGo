<template>
  <div class="p-4 space-y-4">
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer" style="color:#64748b" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" style="color:#0ea5e9"><InfoFilled /></el-icon>
      <h2 class="font-semibold">关于</h2>
    </div>
    <el-card shadow="never">
      <div class="space-y-2">
        <div class="text-gray-700">RecordGo - 任务与积分管理工具</div>
        <div class="text-xs text-gray-500">版本：0.1.1 | 前端 Vue3 + Vite + Element Plus；后端 Gin + SQLite</div>
      </div>
    </el-card>
    <!-- 取消底部返回按钮；统一使用标题左侧返回图标 -->
  </div>
</template>

<script setup lang="ts">
import { InfoFilled, ArrowLeft } from '@element-plus/icons-vue'
import router from '@/router'
function goBack() { router.back() }
</script>

<style scoped>
</style>