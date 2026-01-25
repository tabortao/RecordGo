<template>
  <SettingsShell title="听写大师" subtitle="输入内容或从词库选择，一键开始听写" :icon="Notebook" tone="violet" container-class="max-w-4xl" back-to="/homework">
    <template #headerRight>
      <button
        type="button"
        class="inline-flex h-10 w-10 items-center justify-center rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 text-gray-600 dark:text-gray-300 hover:bg-white hover:text-gray-900 dark:hover:text-gray-50 transition-colors"
        @click="router.push('/dictation/settings')"
      >
        <el-icon :size="18"><Setting /></el-icon>
      </button>
    </template>

    <SettingsCard title="听写内容" description="支持空格/换行分词，也可粘贴音频链接。">
      <template #right>
        <div class="flex items-center gap-2">
          <el-button link type="primary" class="!font-extrabold" @click="showSelector = true">从词库选择</el-button>
          <el-button link type="danger" class="!font-extrabold" @click="content = ''">清空</el-button>
        </div>
      </template>
      <el-input
        v-model="content"
        type="textarea"
        :rows="10"
        placeholder="在此输入听写内容…"
        resize="none"
      />
    </SettingsCard>

    <SettingsCard title="快捷入口">
      <div class="grid grid-cols-3 gap-3">
        <button class="group rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/70 dark:bg-gray-900/55 backdrop-blur-xl px-3 py-4 shadow-sm hover:shadow-md transition active:scale-[0.99]" type="button" @click="router.push('/dictation/banks')">
          <div class="flex flex-col items-center gap-2">
            <div class="h-11 w-11 rounded-2xl border border-blue-100/70 dark:border-blue-900/40 bg-blue-50/80 dark:bg-blue-900/25 flex items-center justify-center text-blue-700 dark:text-blue-300">
              <el-icon :size="20"><Collection /></el-icon>
            </div>
            <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">词库管理</div>
            <div class="text-[11px] text-gray-500 dark:text-gray-400">整理与筛选</div>
          </div>
        </button>
        <button class="group rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/70 dark:bg-gray-900/55 backdrop-blur-xl px-3 py-4 shadow-sm hover:shadow-md transition active:scale-[0.99]" type="button" @click="router.push('/dictation/mistakes')">
          <div class="flex flex-col items-center gap-2">
            <div class="h-11 w-11 rounded-2xl border border-rose-100/70 dark:border-rose-900/40 bg-rose-50/80 dark:bg-rose-900/25 flex items-center justify-center text-rose-700 dark:text-rose-300">
              <el-icon :size="20"><Notebook /></el-icon>
            </div>
            <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">难点收藏</div>
            <div class="text-[11px] text-gray-500 dark:text-gray-400">积累与复盘</div>
          </div>
        </button>
        <button class="group rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/70 dark:bg-gray-900/55 backdrop-blur-xl px-3 py-4 shadow-sm hover:shadow-md transition active:scale-[0.99]" type="button" @click="router.push('/dictation/history')">
          <div class="flex flex-col items-center gap-2">
            <div class="h-11 w-11 rounded-2xl border border-emerald-100/70 dark:border-emerald-900/40 bg-emerald-50/80 dark:bg-emerald-900/25 flex items-center justify-center text-emerald-700 dark:text-emerald-300">
              <el-icon :size="20"><DataLine /></el-icon>
            </div>
            <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">听写记录</div>
            <div class="text-[11px] text-gray-500 dark:text-gray-400">统计与趋势</div>
          </div>
        </button>
      </div>
    </SettingsCard>

    <SettingsCard>
      <el-button type="primary" size="large" class="w-full !h-12 !text-base !rounded-2xl !font-extrabold" :disabled="!content.trim()" @click="startDictation">
        开始听写
      </el-button>
    </SettingsCard>

    <el-drawer v-model="showSelector" title="选择词库内容" direction="btt" size="80%">
      <div class="h-full flex flex-col">
        <WordBankList selector-mode @select="onSelectContent" />
      </div>
    </el-drawer>
  </SettingsShell>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Setting, Collection, Notebook, DataLine } from '@element-plus/icons-vue'
import WordBankList from './WordBankList.vue'
import { dictationApi } from '@/services/dictation'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

const router = useRouter()
const content = ref('')
const showSelector = ref(false)
const stats = ref<any>(null)

onMounted(async () => {
  try {
    const res = await dictationApi.getStats()
    stats.value = (res as any).data || (res as any)
  } catch {}
})

function onSelectContent(text: string) {
  if (content.value.trim()) {
    content.value += '\n' + text
  } else {
    content.value = text
  }
  showSelector.value = false
}

function startDictation() {
  // Pass content via state or query? State is better for large text.
  // However, state is lost on refresh. LocalStorage is safer.
  localStorage.setItem('dictation_current_content', content.value)
  router.push('/dictation/player')
}
</script>

<style scoped>
</style>
