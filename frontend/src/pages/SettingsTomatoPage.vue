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

    <!-- 底部操作按钮：取消与确定 -->
    <div class="flex justify-end gap-2 pt-4">
      <el-button @click="cancel">取消</el-button>
      <el-button type="primary" @click="confirm">确定</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAppState } from '@/stores/appState'
import { Timer, ArrowLeft } from '@element-plus/icons-vue'
import router from '@/router'

// 读取并绑定番茄钟设置到本页面控件
const store = useAppState()
const mode = ref(store.tomato.mode)
const duration = ref(store.tomato.durationMinutes)
const fixed = ref(store.tomato.fixedTomatoPage)

function goBack() {
  // 中文注释：返回不再自动保存，直接关闭页面
  router.back()
}

function cancel() {
  // 中文注释：取消不保存设置，关闭页面
  router.back()
}

function confirm() {
  // 中文注释：确定按钮保存设置并关闭页面
  store.updateTomato({
    mode: mode.value as 'countdown' | 'countup',
    durationMinutes: duration.value,
    fixedTomatoPage: fixed.value
  })
  router.back()
}
</script>

<style scoped>
</style>