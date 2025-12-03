<template>
  <div class="p-4 space-y-4">
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer" style="color:#64748b" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" style="color:#7c3aed"><Setting /></el-icon>
      <h2 class="font-semibold">主题外观</h2>
    </div>

    <!-- 模式选择卡片（对齐番茄钟设置风格） -->
    <el-card shadow="never">
      <div class="font-medium">主题模式</div>
      <div class="text-xs text-gray-500 mt-1">选择系统/深色/浅色三种外观</div>
      <div class="mt-3">
        <el-radio-group v-model="appearance">
          <el-radio label="system">系统</el-radio>
          <el-radio label="dark">深色</el-radio>
          <el-radio label="light">浅色</el-radio>
        </el-radio-group>
      </div>
    </el-card>

    <!-- 说明卡片 -->
    <el-card shadow="never">
      <div class="text-xs text-gray-500 dark:text-gray-400">系统模式下跟随操作系统深浅色；选择深色或浅色将固定该模式。</div>
    </el-card>

    <!-- 底部操作按钮：取消与确定 -->
    <el-card shadow="never">
      <div class="flex justify-end gap-2">
        <el-button @click="cancel">取消</el-button>
        <el-button type="primary" @click="confirm">确定</el-button>
      </div>
    </el-card>
  </div>
  
</template>

<script setup lang="ts">
import { ref } from 'vue'
import router from '@/router'
import { useAppState } from '@/stores/appState'
import { Setting, ArrowLeft } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

function goBack() { router.back() }

const store = useAppState()
const appearance = ref<'system'|'dark'|'light'>(store.themeMode || 'system')

function cancel() { router.back() }

function confirm() {
  store.setThemeMode(appearance.value)
  ElMessage.success('主题外观已更新')
  router.back()
}
</script>

<style scoped>
</style>
