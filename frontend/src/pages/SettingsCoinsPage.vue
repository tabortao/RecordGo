<template>
  <div class="p-4 space-y-4">
    <div class="flex items-center gap-2">
      <!-- 中文注释：返回仅关闭页面，不再自动保存 -->
      <el-icon :size="18" class="cursor-pointer" style="color:#64748b" @click="onCancel"><ArrowLeft /></el-icon>
      <el-icon :size="18" style="color:#f59e0b"><Coin /></el-icon>
      <h2 class="font-semibold">金币设置</h2>
    </div>

    <!-- 中文注释：改为与“编辑个人信息”页一致的左对齐栈叠布局 -->
    <div class="space-y-4">
      <!-- 当前可用金币 -->
      <div class="space-y-2">
        <label class="text-sm text-gray-600">当前可用金币</label>
        <div class="flex items-baseline gap-2">
          <span class="text-gray-900">可用：</span>
          <span class="font-bold text-amber-600 text-2xl">{{ store.coins }}</span>
        </div>
      </div>

      <!-- 设置为新总金币 -->
      <div class="space-y-2">
        <label class="text-sm text-gray-600">设置为新总金币</label>
        <el-input-number
          v-model="newTotal"
          :min="0"
          :max="999999"
          :size="isMobile ? 'default' : 'large'"
          controls-position="right"
          class="w-full sm:w-64"
        />
      </div>

      <!-- 修改原因 -->
      <div class="space-y-2">
        <label class="text-sm text-gray-600">修改原因（仅当变更时必填）</label>
        <el-input
          v-model="reason"
          type="textarea"
          :rows="isMobile ? 3 : 4"
          placeholder="例如：期初校准 / 异常修正"
        />
      </div>

      <!-- 底部右侧按钮，与个人信息页一致 -->
      <div class="flex justify-end mt-6 gap-2">
        <el-button @click="onCancel">取消</el-button>
        <el-button type="primary" @click="onSave">确定</el-button>
      </div>
    </div>

    <div class="text-xs text-gray-500">提示：金币总数更新后，任务页与心愿页会读取同一全局状态并即时刷新。</div>
    
  </div>
</template>

<script setup lang="ts">
import { useAppState } from '@/stores/appState'
import { ElMessage } from 'element-plus'
import { Coin, ArrowLeft } from '@element-plus/icons-vue'
import router from '@/router'
import { useMediaQuery } from '@vueuse/core'
import { setCoins } from '@/services/coins'
 

const store = useAppState()
// 中文注释：根据屏幕宽度判断移动端，用于控制控件大小与布局细节
const isMobile = useMediaQuery('(max-width: 768px)')
// 中文注释：初始化为当前总金币，避免桌面端显示空值
const newTotal = ref<number | null>(Number(store.coins))
const reason = ref<string>('')
 

// 中文注释：取消（仅关闭页面，不保存）
function onCancel() { router.back() }

// 中文注释：确定（校验后保存并关闭）
async function onSave() {
  const current = Number(store.coins)
  const isNoChange = (newTotal.value === null) || (Number(newTotal.value) === current)
  if (isNoChange) { router.back(); return }

  const v = Number(newTotal.value)
  if (isNaN(v) || v < 0) { ElMessage.error('请输入有效的金币数量（≥ 0）'); return }

  const r = reason.value.trim()
  if (!r) { ElMessage.error('请填写修改原因'); return }

  try {
    const resp = await setCoins(v, r)
    store.setCoins(Number(resp.coins || v))
    ElMessage.success('金币总数已更新')
  } catch (e: any) {
    ElMessage.error(e?.message || '保存失败')
    return
  }
  router.back()
}
</script>

<style scoped>
</style>
