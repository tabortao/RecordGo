<template>
  <SettingsShell title="主题外观" subtitle="系统/深色/浅色" :icon="Setting" tone="violet">
    <SettingsCard title="主题模式" description="选择系统/深色/浅色三种外观">
      <el-radio-group v-model="appearance" class="!flex !flex-wrap gap-x-6 gap-y-2">
        <el-radio label="system">系统</el-radio>
        <el-radio label="dark">深色</el-radio>
        <el-radio label="light">浅色</el-radio>
      </el-radio-group>
      <div class="mt-3 text-xs text-gray-500 dark:text-gray-400">系统模式会跟随操作系统深浅色；深色或浅色将固定该模式。</div>
    </SettingsCard>

    <SettingsCard>
      <div class="flex justify-end gap-2">
        <el-button class="!rounded-xl" @click="cancel">取消</el-button>
        <el-button type="primary" class="!rounded-xl" @click="confirm">确定</el-button>
      </div>
    </SettingsCard>
  </SettingsShell>
  
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAppState } from '@/stores/appState'
import { Setting } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import router from '@/router'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

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
