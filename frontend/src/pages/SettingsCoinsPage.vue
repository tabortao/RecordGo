<template>
  <div class="p-4 space-y-4">
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer" style="color:#64748b" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" style="color:#f59e0b"><Coin /></el-icon>
      <h2 class="font-semibold">金币设置</h2>
    </div>

    <el-card shadow="never">
      <div class="space-y-3">
        <div class="flex items-center gap-2">
          <span class="text-gray-600">当前可用金币：</span>
          <span class="font-semibold text-yellow-600">{{ store.coins }}</span>
        </div>

        <!-- 修改总金币：仅允许设置新总额，取消“增减值调整” -->
        <div class="space-y-2">
          <label class="text-sm text-gray-600">设置为新总金币</label>
          <el-input-number v-model="newTotal" :min="0" :max="999999" />
        </div>

        <div class="space-y-2">
          <label class="text-sm text-gray-600">修改原因（仅当变更时必填）</label>
          <el-input v-model="reason" type="textarea" :rows="3" placeholder="变更时请填写原因，例如：期初校准/异常修正" />
        </div>

        <!-- 取消底部保存/返回按钮；点击左侧返回图标时自动保存 -->
      </div>
    </el-card>

    <div class="text-xs text-gray-500">提示：金币总数更新后，任务页与心愿页会读取同一全局状态并即时刷新。</div>
  </div>
</template>

<script setup lang="ts">
import { useAppState } from '@/stores/appState'
import { ElMessage } from 'element-plus'
import { Coin, ArrowLeft } from '@element-plus/icons-vue'
import router from '@/router'

const store = useAppState()
const newTotal = ref<number | null>(null)
const reason = ref<string>('')

function goBack() {
  // 中文注释：无任何修改时允许直接返回，不强制填写原因
  const current = Number(store.coins)
  const isNoChange = (newTotal.value === null) || (Number(newTotal.value) === current)
  if (isNoChange) { router.back(); return }

  const v = Number(newTotal.value)
  if (isNaN(v) || v < 0) { ElMessage.error('请输入有效的金币数量（≥ 0）'); return }

  const r = reason.value.trim()
  if (!r) { ElMessage.error('请填写修改原因'); return }

  store.setCoins(v)
  // 中文注释：如需记录操作日志，可在此调用后端 API（预留）
  ElMessage.success('金币总数已更新')
  router.back()
}
</script>

<style scoped>
</style>