<template>
  <div class="p-4 space-y-4">
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer" style="color:#64748b" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" style="color:#ef4444"><Timer /></el-icon>
      <h2 class="font-semibold text-base">番茄钟设置</h2>
    </div>

    <!-- 模式设置（倒计时/正计时） -->
    <div class="space-y-2">
      <label class="text-base text-gray-700">计时模式</label>
      <el-radio-group v-model="mode">
        <el-radio label="countdown">倒计时</el-radio>
        <el-radio label="countup">正计时</el-radio>
      </el-radio-group>
    </div>

    <!-- 预设时长（分钟） -->
    <div class="space-y-2">
      <label class="text-base text-gray-700">预设时长（分钟）</label>
      <el-input-number v-model="duration" :min="1" :max="180" />
    </div>

    <!-- 固定全屏番茄钟页面 -->
    <div class="space-y-2">
      <el-switch v-model="fixed" active-text="固定番茄钟页面" inactive-text="不固定" />
      <div class="text-base text-gray-600">开启后，进入番茄钟为固定全屏，适合专注模式。</div>
    </div>

    <!-- 取消底部保存/取消按钮；用户点击左侧返回图标时自动保存并返回 -->
  </div>
</template>

<script setup lang="ts">
import { useAppState } from '@/stores/appState'
import { ElMessage } from 'element-plus'
import { Timer, ArrowLeft } from '@element-plus/icons-vue'
import router from '@/router'

// 读取并绑定番茄钟设置到本页面控件
const store = useAppState()
const mode = ref(store.tomato.mode)
const duration = ref(store.tomato.durationMinutes)
const fixed = ref(store.tomato.fixedTomatoPage)

function goBack() {
  // 中文注释：返回时自动保存最新设置
  store.updateTomato({
    mode: mode.value as 'countdown' | 'countup',
    durationMinutes: duration.value,
    fixedTomatoPage: fixed.value
  })
  ElMessage.success('番茄钟设置已保存')
  router.back()
}
</script>

<style scoped>
</style>