<template>
  <div class="p-4 space-y-4">
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer" style="color:#64748b" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" style="color:#ef4444"><Timer /></el-icon>
      <h2 class="font-semibold">番茄钟设置</h2>
    </div>

    <!-- 计时模式卡片 -->
    <el-card shadow="never">
      <div class="font-medium">计时模式</div>
      <div class="text-xs text-gray-500 mt-1">选择倒计时或正计时显示方式</div>
      <div class="mt-3">
        <el-radio-group v-model="mode">
          <el-radio label="countdown">倒计时</el-radio>
          <el-radio label="countup">正计时</el-radio>
        </el-radio-group>
      </div>
    </el-card>

    <!-- 预设时长卡片 -->
    <el-card shadow="never">
      <div class="flex items-center justify-between">
        <div>
          <div class="font-medium">预设时长（分钟）</div>
          <div class="text-xs text-gray-500 mt-1">进入番茄钟页面时默认使用的计划时长</div>
        </div>
        <el-input-number v-model="duration" :min="1" :max="180" />
      </div>
    </el-card>

    <!-- 固定全屏卡片 -->
    <el-card shadow="never">
      <div class="flex items-center justify-between">
        <div>
          <div class="font-medium">固定番茄钟页面</div>
          <div class="text-xs text-gray-500 mt-1">开启后：进入番茄钟为固定全屏，适合专注模式</div>
        </div>
        <el-switch v-model="fixed" />
      </div>
    </el-card>

    <!-- 保持常亮（Wake Lock）卡片 -->
    <el-card shadow="never">
      <div class="flex items-center justify-between">
        <div>
          <div class="font-medium">保持常亮</div>
          <div class="text-xs text-gray-500 mt-1">移动端开始计时后保持屏幕常亮，默认开启</div>
        </div>
        <el-switch v-model="keepAwake" />
      </div>
    </el-card>

    <!-- 底部操作按钮：取消与确定（与任务设置保持一致的页内按钮风格） -->
    <el-card shadow="never">
      <div class="flex justify-end gap-2">
        <el-button @click="cancel">取消</el-button>
        <el-button type="primary" @click="confirm">确定</el-button>
      </div>
    </el-card>
    <div class="mt-2">
      <el-button v-if="isAdmin" type="primary" @click="router.push('/admin')">用户管理</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAppState } from '@/stores/appState'
import { Timer, ArrowLeft } from '@element-plus/icons-vue'
import router from '@/router'
import { useAuth } from '@/stores/auth'
import { computed } from 'vue'

// 读取并绑定番茄钟设置到本页面控件
const store = useAppState()
const mode = ref(store.tomato.mode)
const duration = ref(store.tomato.durationMinutes)
const fixed = ref(store.tomato.fixedTomatoPage)
const keepAwake = ref(store.tomato.keepAwakeEnabled)
const auth = useAuth()
const isAdmin = computed(() => Number(auth.user?.id || 0) === 1)

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
    runningMode: null,
    durationMinutes: duration.value,
    fixedTomatoPage: fixed.value,
    keepAwakeEnabled: keepAwake.value
  })
  router.back()
}
</script>

<style scoped>
</style>
