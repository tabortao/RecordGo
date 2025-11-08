<template>
  <div class="p-4 space-y-4">
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer" style="color:#64748b" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" style="color:#7c3aed"><Microphone /></el-icon>
      <h2 class="font-semibold">朗读设置</h2>
    </div>
    <el-card shadow="never">
      <div class="text-gray-600">配置项待完善：TTS 引擎、语速音量、语言等。</div>
    </el-card>
    <!-- 取消底部返回按钮；统一使用标题左侧返回图标 -->
  </div>
</template>

<script setup lang="ts">
import { Microphone, ArrowLeft } from '@element-plus/icons-vue'
import router from '@/router'
function goBack() { router.back() }
</script>

<style scoped>
</style>