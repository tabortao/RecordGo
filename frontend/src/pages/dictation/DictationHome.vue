<template>
  <SettingsShell title="听写大师" subtitle="输入内容或从词库选择，一键开始听写" :icon="Notebook" tone="violet" container-class="max-w-4xl" back-to="/homework">
    <template #headerRight>
      <button
        type="button"
        class="inline-flex h-10 w-10 items-center justify-center rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 text-gray-600 dark:text-gray-300 hover:bg-white hover:text-gray-900 dark:hover:text-gray-50 transition-colors shadow-sm"
        @click="router.push('/dictation/settings')"
      >
        <el-icon :size="20"><Setting /></el-icon>
      </button>
    </template>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
      <!-- 左侧：主要功能区 -->
      <div class="lg:col-span-2 space-y-6">
        <SettingsCard title="听写内容" description="支持空格/换行分词，也可粘贴音频链接。">
          <template #right>
            <div class="flex items-center gap-2">
              <el-button link type="primary" class="!font-extrabold" @click="showSelector = true">从词库选择</el-button>
              <el-button link type="danger" class="!font-extrabold" @click="content = ''">清空</el-button>
            </div>
          </template>
          <div class="relative">
            <el-input
              v-model="content"
              type="textarea"
              :rows="12"
              placeholder="在此输入听写内容…
例如：
apple 苹果
banana 香蕉
或者一段中文课文..."
              resize="none"
              class="dictation-input"
            />
            <div class="absolute bottom-3 right-3 text-xs text-gray-400 dark:text-gray-500 font-mono">
              {{ content.length }} 字符
            </div>
          </div>
        </SettingsCard>
        
        <SettingsCard :decor="false" class="!bg-transparent !border-none !shadow-none !p-0">
          <el-button type="primary" size="large" class="w-full !h-14 !text-lg !rounded-2xl !font-extrabold shadow-lg shadow-violet-500/20 hover:shadow-violet-500/30 transition-shadow" :disabled="!content.trim()" @click="startDictation">
            开始听写
          </el-button>
        </SettingsCard>
      </div>

      <!-- 右侧：快捷入口 -->
      <div class="lg:col-span-1 space-y-4">
        <div class="text-xs font-bold uppercase tracking-wider text-gray-400 dark:text-gray-500 px-1">常用功能</div>
        <div class="grid grid-cols-1 gap-3">
          <button class="group relative overflow-hidden rounded-[24px] border border-white/50 dark:border-gray-800/60 bg-white/60 dark:bg-gray-900/40 p-4 shadow-sm hover:shadow-md transition-all duration-300 hover:bg-white dark:hover:bg-gray-900 text-left" type="button" @click="router.push('/dictation/banks')">
            <div class="absolute right-0 top-0 -mt-4 -mr-4 h-24 w-24 rounded-full bg-blue-50 dark:bg-blue-900/10 transition-transform group-hover:scale-110"></div>
            <div class="relative flex items-center gap-4">
              <div class="h-12 w-12 shrink-0 rounded-2xl bg-blue-100/50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 flex items-center justify-center">
                <el-icon :size="24"><Collection /></el-icon>
              </div>
              <div>
                <div class="text-base font-black text-gray-900 dark:text-gray-50">词库管理</div>
                <div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">整理与筛选常用词汇</div>
              </div>
            </div>
          </button>

          <button class="group relative overflow-hidden rounded-[24px] border border-white/50 dark:border-gray-800/60 bg-white/60 dark:bg-gray-900/40 p-4 shadow-sm hover:shadow-md transition-all duration-300 hover:bg-white dark:hover:bg-gray-900 text-left" type="button" @click="router.push('/dictation/mistakes')">
            <div class="absolute right-0 top-0 -mt-4 -mr-4 h-24 w-24 rounded-full bg-rose-50 dark:bg-rose-900/10 transition-transform group-hover:scale-110"></div>
            <div class="relative flex items-center gap-4">
              <div class="h-12 w-12 shrink-0 rounded-2xl bg-rose-100/50 dark:bg-rose-900/30 text-rose-600 dark:text-rose-400 flex items-center justify-center">
                <el-icon :size="24"><Notebook /></el-icon>
              </div>
              <div>
                <div class="text-base font-black text-gray-900 dark:text-gray-50">难点收藏</div>
                <div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">积累错误，重点复盘</div>
              </div>
            </div>
          </button>

          <button class="group relative overflow-hidden rounded-[24px] border border-white/50 dark:border-gray-800/60 bg-white/60 dark:bg-gray-900/40 p-4 shadow-sm hover:shadow-md transition-all duration-300 hover:bg-white dark:hover:bg-gray-900 text-left" type="button" @click="router.push('/dictation/history')">
            <div class="absolute right-0 top-0 -mt-4 -mr-4 h-24 w-24 rounded-full bg-emerald-50 dark:bg-emerald-900/10 transition-transform group-hover:scale-110"></div>
            <div class="relative flex items-center gap-4">
              <div class="h-12 w-12 shrink-0 rounded-2xl bg-emerald-100/50 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400 flex items-center justify-center">
                <el-icon :size="24"><DataLine /></el-icon>
              </div>
              <div>
                <div class="text-base font-black text-gray-900 dark:text-gray-50">听写记录</div>
                <div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">查看历史成绩与趋势</div>
              </div>
            </div>
          </button>
        </div>
        
        <div v-if="stats" class="mt-6 rounded-[24px] bg-violet-500/5 dark:bg-violet-500/10 p-5 border border-violet-500/10">
          <div class="text-xs font-bold uppercase tracking-wider text-violet-600 dark:text-violet-400 mb-3">学习统计</div>
          <div class="flex justify-between items-end">
            <div>
              <div class="text-3xl font-black text-gray-900 dark:text-white tabular-nums">{{ stats.total_dictations || 0 }}</div>
              <div class="text-xs font-medium text-gray-500 dark:text-gray-400 mt-1">总听写次数</div>
            </div>
            <div class="text-right">
              <div class="text-xl font-black text-gray-900 dark:text-white tabular-nums">{{ stats.total_duration_minutes || 0 }}<span class="text-xs font-medium text-gray-500 ml-1">分钟</span></div>
              <div class="text-xs font-medium text-gray-500 dark:text-gray-400 mt-1">总时长</div>
            </div>
          </div>
        </div>
      </div>
    </div>

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
:deep(.dictation-input .el-textarea__inner) {
  padding: 16px;
  font-size: 16px;
  line-height: 1.6;
  border-radius: 20px;
  background-color: rgb(249 250 251); /* gray-50 */
  border: 1px solid rgb(243 244 246); /* gray-100 */
  box-shadow: none;
  transition: all 0.2s;
}

:deep(.dictation-input .el-textarea__inner:focus) {
  background-color: white;
  border-color: rgb(139 92 246); /* violet-500 */
  box-shadow: 0 0 0 4px rgb(139 92 246 / 0.1);
}

.dark :deep(.dictation-input .el-textarea__inner) {
  background-color: rgb(17 24 39 / 0.5); /* gray-900/50 */
  border-color: rgb(31 41 55); /* gray-800 */
  color: rgb(243 244 246); /* gray-100 */
}

.dark :deep(.dictation-input .el-textarea__inner:focus) {
  background-color: rgb(3 7 18); /* gray-950 */
  border-color: rgb(139 92 246); /* violet-500 */
}
</style>
