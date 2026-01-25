<template>
  <SettingsShell title="拼音王国" :icon="EditPen" tone="sky" container-class="max-w-5xl" back-to="/homework/chinese">
    <SettingsCard>
      <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
        <button
          type="button"
          class="group rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/70 dark:bg-gray-900/55 backdrop-blur px-4 py-4 shadow-sm hover:shadow-md transition active:scale-[0.99] text-left"
          @click="router.push('/homework/chinese/pinyin/wheel')"
        >
          <div class="flex items-center justify-between">
            <div class="h-11 w-11 rounded-2xl border border-emerald-100/80 dark:border-emerald-900/40 bg-emerald-50/80 dark:bg-emerald-900/25 flex items-center justify-center text-2xl">🎡</div>
            <div class="text-[11px] font-extrabold uppercase tracking-wider text-emerald-700/70 dark:text-emerald-300/70">轻松练习</div>
          </div>
          <div class="mt-3 text-base font-extrabold tracking-tight text-gray-900 dark:text-gray-50">拼音大转盘</div>
        </button>

        <button
          type="button"
          class="group rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/70 dark:bg-gray-900/55 backdrop-blur px-4 py-4 shadow-sm hover:shadow-md transition active:scale-[0.99] text-left"
          @click="router.push('/homework/chinese/pinyin/adventure')"
        >
          <div class="flex items-center justify-between">
            <div class="h-11 w-11 rounded-2xl border border-sky-100/80 dark:border-sky-900/40 bg-sky-50/80 dark:bg-sky-900/25 flex items-center justify-center text-2xl">🏰</div>
            <div class="text-[11px] font-extrabold uppercase tracking-wider text-sky-700/70 dark:text-sky-300/70">闯关模式</div>
          </div>
          <div class="mt-3 text-base font-extrabold tracking-tight text-gray-900 dark:text-gray-50">智能拼音大闯关</div>
        </button>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { EditPen } from '@element-plus/icons-vue'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

const router = useRouter()
</script>
