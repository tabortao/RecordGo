<template>
  <SettingsShell title="金币设置" subtitle="校准金币总数（需填写原因）" :icon="Coin" tone="amber">
    <SettingsCard>
      <div class="flex items-center justify-between gap-4">
        <div>
          <div class="text-xs text-gray-500 dark:text-gray-400">当前可用金币</div>
          <div class="mt-1 flex items-baseline gap-2">
            <span class="text-sm text-gray-700 dark:text-gray-300">可用</span>
            <span class="text-3xl font-extrabold tracking-tight text-amber-700 dark:text-amber-300">{{ store.coins }}</span>
          </div>
        </div>
        <div class="rounded-2xl border border-amber-100/80 dark:border-amber-900/40 bg-amber-50/70 dark:bg-amber-900/20 px-3 py-2 text-xs text-amber-800 dark:text-amber-200">
          仅当变更时需要原因
        </div>
      </div>
    </SettingsCard>

    <SettingsCard title="设置为新总金币" description="将金币总数直接设置为指定值">
      <el-input-number
        v-model="newTotal"
        :min="0"
        :max="999999"
        :size="isMobile ? 'default' : 'large'"
        controls-position="right"
        class="w-full sm:w-64"
      />
    </SettingsCard>

    <SettingsCard title="修改原因" description="仅当金币发生变更时必填">
      <el-input
        v-model="reason"
        type="textarea"
        :rows="isMobile ? 3 : 4"
        placeholder="例如：期初校准 / 异常修正"
      />
    </SettingsCard>

    <SettingsCard>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div class="text-xs text-gray-500 dark:text-gray-400">提示：金币更新后，任务页与心愿页会读取同一全局状态并即时刷新。</div>
        <div class="flex justify-end gap-2">
          <el-button class="!rounded-xl" @click="onCancel">取消</el-button>
          <el-button type="primary" class="!rounded-xl" @click="onSave">确定</el-button>
        </div>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { useAppState } from '@/stores/appState'
import { ElMessage } from 'element-plus'
import { Coin } from '@element-plus/icons-vue'
import router from '@/router'
import { useMediaQuery } from '@vueuse/core'
import { setCoins } from '@/services/coins'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'
 

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
