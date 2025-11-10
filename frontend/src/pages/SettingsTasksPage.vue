<template>
  <div class="p-4 space-y-4">
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer" style="color:#64748b" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" style="color:#10b981"><List /></el-icon>
      <h2 class="font-semibold">任务设置</h2>
    </div>
    <!-- 中文注释：任务备注入口开关（默认开启）；关闭后任务卡片不显示备注图标，菜单不显示备注按钮 -->
    <el-card shadow="never">
      <div class="flex items-center justify-between">
        <div>
          <div class="font-medium">显示任务备注入口</div>
          <div class="text-xs text-gray-500 mt-1">关闭后：任务清单卡片不显示备注图标，操作菜单不显示“备注”按钮</div>
        </div>
        <el-switch v-model="taskNotesEnabled" />
      </div>
    </el-card>
    <!-- 中文注释：任务自动排序开关（默认开启）；开启后分类内已完成排下方，全部完成的分类排到未完分类之后 -->
    <el-card shadow="never">
      <div class="flex items-center justify-between">
        <div>
          <div class="font-medium">任务自动排序</div>
          <div class="text-xs text-gray-500 mt-1">开启后：分类内已完成任务排在下方；当某分类全部完成时，该分类整体排在未完分类之后</div>
        </div>
        <el-switch v-model="taskAutoSortEnabled" />
      </div>
    </el-card>
    <el-card shadow="never">
      <div class="text-gray-600">其他配置项待完善：任务默认提醒、重复规则等。</div>
    </el-card>
    <!-- 取消底部返回按钮；统一使用标题左侧返回图标 -->
  </div>
</template>

<script setup lang="ts">
import { List, ArrowLeft } from '@element-plus/icons-vue'
import router from '@/router'
import { useAppState } from '@/stores/appState'
import { computed } from 'vue'
function goBack() { router.back() }

// 中文注释：联动全局状态中的任务备注开关，使用计算属性实现双向绑定
const store = useAppState()
const taskNotesEnabled = computed({
  get: () => store.taskNotesEnabled,
  set: (v: boolean) => store.setTaskNotesEnabled(v)
})
const taskAutoSortEnabled = computed({
  get: () => store.taskAutoSortEnabled,
  set: (v: boolean) => store.setTaskAutoSortEnabled(v)
})
</script>

<style scoped>
</style>