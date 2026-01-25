<template>
  <SettingsShell title="语文" :icon="Reading" tone="red" container-class="max-w-5xl" back-to="/homework">
    <SettingsCard>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
        <button
          type="button"
          class="group rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/70 dark:bg-gray-900/55 backdrop-blur px-3 py-4 shadow-sm hover:shadow-md transition active:scale-[0.99]"
          @click="router.push('/homework/chinese/pinyin')"
        >
          <div class="flex flex-col items-center gap-2">
            <div class="relative">
              <div class="absolute -inset-3 rounded-2xl bg-sky-200/40 dark:bg-sky-500/12 blur-xl opacity-0 group-hover:opacity-100 transition-opacity" />
              <div class="relative h-14 w-14 rounded-2xl border border-sky-100/80 dark:border-sky-900/40 bg-sky-50/80 dark:bg-sky-900/25 flex items-center justify-center text-sky-700 dark:text-sky-300">
                <span class="text-3xl font-extrabold leading-none">a</span>
              </div>
            </div>
            <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">拼音</div>
          </div>
        </button>

        <button
          type="button"
          class="group rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/70 dark:bg-gray-900/55 backdrop-blur px-3 py-4 shadow-sm hover:shadow-md transition active:scale-[0.99]"
          @click="router.push('/homework/chinese/hanzi-hero')"
        >
          <div class="flex flex-col items-center gap-2">
            <div class="relative">
              <div class="absolute -inset-3 rounded-2xl bg-amber-200/40 dark:bg-amber-500/12 blur-xl opacity-0 group-hover:opacity-100 transition-opacity" />
              <div class="relative h-14 w-14 rounded-2xl border border-amber-100/80 dark:border-amber-900/40 bg-amber-50/80 dark:bg-amber-900/20 overflow-hidden">
                <img src="@/assets/button/hanzihero.png" alt="Hanzi" class="h-full w-full object-cover opacity-90" />
              </div>
            </div>
            <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">汉字英雄</div>
          </div>
        </button>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Reading } from '@element-plus/icons-vue'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

const router = useRouter()
</script>
