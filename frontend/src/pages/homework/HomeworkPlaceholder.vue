<template>
  <SettingsShell :title="title" subtitle="功能开发中" :icon="Tools" :tone="tone" container-class="max-w-3xl" back-to="/homework">
    <SettingsCard>
      <div class="flex flex-col items-center justify-center py-10 text-center">
        <div class="relative">
          <div class="absolute -inset-6 rounded-[28px] bg-gray-200/55 dark:bg-gray-700/25 blur-2xl" />
          <div class="relative h-20 w-20 rounded-[28px] border border-white/60 dark:border-gray-800/60 bg-white/75 dark:bg-gray-950/30 flex items-center justify-center">
            <el-icon :size="34" class="text-gray-500 dark:text-gray-300"><Tools /></el-icon>
          </div>
        </div>
        <div class="mt-6 text-lg font-extrabold tracking-tight text-gray-900 dark:text-gray-50">正在加速开发</div>
        <div class="mt-1 text-sm text-gray-600 dark:text-gray-300">该学科模块即将上线，敬请期待。</div>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Tools } from '@element-plus/icons-vue'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

const route = useRoute()
const title = computed(() => (route.query.title as string) || '作业家')

const tone = computed(() => {
  const t = String(title.value || '')
  if (t.includes('数学')) return 'blue'
  if (t.includes('英语')) return 'amber'
  if (t.includes('科学')) return 'emerald'
  if (t.includes('语文')) return 'red'
  return 'sky'
})
</script>
