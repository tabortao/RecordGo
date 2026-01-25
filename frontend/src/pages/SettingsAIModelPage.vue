<template>
  <SettingsShell title="AI 模型设置" subtitle="服务商、对话与视觉模型" :icon="Cpu" tone="indigo" container-class="max-w-4xl">
    <div class="space-y-6 md:space-y-8">
      <SettingsCard>
        <div class="space-y-3">
          <div class="flex items-center gap-2 text-xs font-extrabold uppercase tracking-wider text-gray-500 dark:text-gray-400">
            <el-icon><Connection /></el-icon>
            选择服务商
          </div>
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-3">
            <button
              v-for="(provider, key) in PROVIDERS"
              :key="key"
              type="button"
              @click="handleProviderChange(key)"
              class="px-2 py-4 rounded-2xl text-[10px] md:text-xs font-extrabold border transition-all flex flex-col items-center justify-center gap-2"
              :class="activeProvider === key ? 'bg-indigo-600 text-white border-indigo-600 shadow-lg shadow-indigo-200/40 dark:shadow-indigo-950/30 scale-[1.02]' : 'bg-white/70 dark:bg-gray-900/50 text-gray-700 dark:text-gray-200 border-gray-100 dark:border-gray-800 hover:border-indigo-200 dark:hover:border-indigo-800 hover:bg-white dark:hover:bg-gray-900/70'"
            >
              <span class="text-xl md:text-2xl">{{ getProviderEmoji(key) }}</span>
              <span class="text-center truncate w-full px-1">{{ provider.name }}</span>
            </button>
          </div>
        </div>
      </SettingsCard>

      <SettingsCard>
        <template #header>
          <div class="flex items-start justify-between gap-4">
            <div class="min-w-0">
              <h3 class="text-sm font-extrabold text-gray-900 dark:text-gray-50 flex items-center gap-2">
                <el-icon :size="18" class="text-indigo-600 dark:text-indigo-300"><ChatDotRound /></el-icon> 对话模型 (Chat)
              </h3>
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">用于文字生成、故事创作和汉字讲解。</p>
            </div>
            <span v-if="testStatus === 'SUCCESS'" class="shrink-0 rounded-full border border-green-200/70 dark:border-green-900/40 bg-green-50/70 dark:bg-green-900/20 px-3 py-1 text-[10px] font-extrabold text-green-700 dark:text-green-300 flex items-center gap-1">
              <el-icon :size="14"><Select /></el-icon> 已连接
            </span>
          </div>
        </template>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-6">
          <div class="col-span-1 md:col-span-2">
            <label class="block text-[10px] font-extrabold text-gray-500 dark:text-gray-400 mb-2 uppercase">API Host (代理地址)</label>
            <input
              type="text"
              :placeholder="activeProvider === 'GOOGLE' ? '默认无需填写' : 'https://api.example.com/v1'"
              class="w-full p-3 md:p-4 rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-950/40 text-gray-900 dark:text-gray-100 focus:bg-white dark:focus:bg-gray-950/60 focus:border-indigo-500 outline-none text-sm font-mono transition-all"
              v-model="config.apiBaseUrl"
            />
          </div>
          <div class="col-span-1 md:col-span-2">
            <label class="block text-[10px] font-extrabold text-gray-500 dark:text-gray-400 mb-2 uppercase">API Key (密钥)</label>
            <div class="relative group">
              <input
                :type="showKey ? 'text' : 'password'"
                placeholder="sk-..."
                class="w-full p-3 md:p-4 pr-12 rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-950/40 text-gray-900 dark:text-gray-100 focus:bg-white dark:focus:bg-gray-950/60 focus:border-indigo-500 outline-none text-sm font-mono transition-all"
                v-model="config.apiKey"
              />
              <button type="button" @click="showKey = !showKey" class="absolute right-4 top-3 md:top-4 text-gray-400 hover:text-indigo-600 transition-colors">
                <el-icon :size="20"><component :is="showKey ? Hide : View" /></el-icon>
              </button>
            </div>
          </div>
          <div class="col-span-1 md:col-span-2">
            <label class="block text-[10px] font-extrabold text-gray-500 dark:text-gray-400 mb-2 uppercase">Model Name (模型)</label>
            <input
              type="text"
              placeholder="如: gemini-2.5-flash"
              class="w-full p-3 md:p-4 rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-950/40 text-gray-900 dark:text-gray-100 focus:bg-white dark:focus:bg-gray-950/60 focus:border-indigo-500 outline-none text-sm font-mono transition-all"
              v-model="config.model"
            />
          </div>
        </div>

        <button
          type="button"
          @click="handleTestConnection"
          :disabled="testStatus === 'TESTING' || !config.apiKey"
          class="mt-5 w-full py-3 md:py-4 rounded-2xl text-sm font-extrabold flex items-center justify-center gap-2 transition-all border disabled:opacity-50"
          :class="testStatus === 'SUCCESS' ? 'bg-green-50/80 dark:bg-green-900/20 border-green-200/70 dark:border-green-900/40 text-green-700 dark:text-green-300' : testStatus === 'FAIL' ? 'bg-red-50/80 dark:bg-red-900/20 border-red-200/70 dark:border-red-900/40 text-red-700 dark:text-red-300' : 'bg-white/70 dark:bg-gray-900/50 border-indigo-100/80 dark:border-indigo-900/40 text-indigo-700 dark:text-indigo-300 hover:bg-white dark:hover:bg-gray-900/70'"
        >
          <template v-if="testStatus === 'TESTING'">
            <el-icon class="is-loading" :size="18"><Loading /></el-icon> 连接测试中...
          </template>
          <template v-else-if="testStatus === 'SUCCESS'">
            <el-icon :size="18"><Link /></el-icon> 连接成功
          </template>
          <template v-else-if="testStatus === 'FAIL'">
            <el-icon :size="18"><CloseBold /></el-icon> 连接失败
          </template>
          <template v-else>
            <el-icon :size="18"><Link /></el-icon> 测试对话连接
          </template>
        </button>
      </SettingsCard>

      <SettingsCard>
        <template #header>
          <div class="flex items-start justify-between gap-4">
            <div class="min-w-0">
              <h3 class="text-sm font-extrabold text-gray-900 dark:text-gray-50 flex items-center gap-2">
                <el-icon :size="18" class="text-violet-600 dark:text-violet-300"><Picture /></el-icon> 视觉模型 (Vision)
              </h3>
              <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">用于拍照识字和看图写话。</p>
            </div>
            <span v-if="visionTestStatus === 'SUCCESS'" class="shrink-0 rounded-full border border-green-200/70 dark:border-green-900/40 bg-green-50/70 dark:bg-green-900/20 px-3 py-1 text-[10px] font-extrabold text-green-700 dark:text-green-300 flex items-center gap-1">
              <el-icon :size="14"><Select /></el-icon> 已连接
            </span>
          </div>
        </template>

        <div class="grid grid-cols-1 gap-4">
          <div>
            <label class="block text-[10px] font-extrabold text-gray-500 dark:text-gray-400 mb-2 uppercase">Vision API Host</label>
            <input
              type="text"
              placeholder="同上 (默认)"
              class="w-full p-3 md:p-4 rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-950/40 text-gray-900 dark:text-gray-100 focus:bg-white dark:focus:bg-gray-950/60 focus:border-violet-500 outline-none text-sm font-mono transition-all"
              v-model="config.visionApiBaseUrl"
            />
          </div>
          <div>
            <label class="block text-[10px] font-extrabold text-gray-500 dark:text-gray-400 mb-2 uppercase">Vision API Key</label>
            <div class="relative group">
              <input
                :type="showVisionKey ? 'text' : 'password'"
                placeholder="同上 (默认)"
                class="w-full p-3 md:p-4 pr-12 rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-950/40 text-gray-900 dark:text-gray-100 focus:bg-white dark:focus:bg-gray-950/60 focus:border-violet-500 outline-none text-sm font-mono transition-all"
                v-model="config.visionApiKey"
              />
              <button type="button" @click="showVisionKey = !showVisionKey" class="absolute right-4 top-3 md:top-4 text-gray-400 hover:text-violet-600 transition-colors">
                <el-icon :size="20"><component :is="showVisionKey ? Hide : View" /></el-icon>
              </button>
            </div>
          </div>
          <div>
            <label class="block text-[10px] font-extrabold text-gray-500 dark:text-gray-400 mb-2 uppercase">Vision Model Name</label>
            <input
              type="text"
              placeholder="默认: gemini-2.5-flash"
              class="w-full p-3 md:p-4 rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-950/40 text-gray-900 dark:text-gray-100 focus:bg-white dark:focus:bg-gray-950/60 focus:border-violet-500 outline-none text-sm font-mono transition-all"
              v-model="config.visionModel"
            />
          </div>
        </div>

        <button
          type="button"
          @click="handleTestVisionConnection"
          :disabled="visionTestStatus === 'TESTING' || (!config.apiKey && !config.visionApiKey)"
          class="mt-5 w-full py-3 md:py-4 rounded-2xl text-sm font-extrabold flex items-center justify-center gap-2 transition-all border disabled:opacity-50"
          :class="visionTestStatus === 'SUCCESS' ? 'bg-green-50/80 dark:bg-green-900/20 border-green-200/70 dark:border-green-900/40 text-green-700 dark:text-green-300' : visionTestStatus === 'FAIL' ? 'bg-red-50/80 dark:bg-red-900/20 border-red-200/70 dark:border-red-900/40 text-red-700 dark:text-red-300' : 'bg-white/70 dark:bg-gray-900/50 border-violet-100/80 dark:border-violet-900/40 text-violet-700 dark:text-violet-300 hover:bg-white dark:hover:bg-gray-900/70'"
        >
          <template v-if="visionTestStatus === 'TESTING'">
            <el-icon class="is-loading" :size="18"><Loading /></el-icon> 测试中...
          </template>
          <template v-else-if="visionTestStatus === 'SUCCESS'">
            <el-icon :size="18"><Link /></el-icon> 视觉连接成功
          </template>
          <template v-else-if="visionTestStatus === 'FAIL'">
            <el-icon :size="18"><CloseBold /></el-icon> 视觉连接失败
          </template>
          <template v-else>
            <el-icon :size="18"><Picture /></el-icon> 测试视觉连接
          </template>
        </button>
      </SettingsCard>
    </div>
  </SettingsShell>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { useAIStore, PROVIDERS } from '@/stores/ai'
import { testConnection, testVisionConnection } from '@/services/aiService'
import { Cpu, Connection, ChatDotRound, Select, View, Hide, Loading, Link, CloseBold, Picture } from '@element-plus/icons-vue'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'
const store = useAIStore()

const activeProvider = ref(store.activeProvider)
const config = reactive({
    apiBaseUrl: store.apiBaseUrl,
    apiKey: store.apiKey,
    model: store.model,
    visionApiBaseUrl: store.visionApiBaseUrl,
    visionApiKey: store.visionApiKey,
    visionModel: store.visionModel
})

const showKey = ref(false)
const showVisionKey = ref(false)
const testStatus = ref<'IDLE' | 'TESTING' | 'SUCCESS' | 'FAIL'>('IDLE')
const visionTestStatus = ref<'IDLE' | 'TESTING' | 'SUCCESS' | 'FAIL'>('IDLE')

const getProviderEmoji = (key: string) => {
    const map: Record<string, string> = {
        GOOGLE: '🌍',
        ZHIPU: '🧠',
        DEEPSEEK: '🐳',
        SILICON: '🚀',
        CUSTOM: '⚙️'
    }
    return map[key] || '❓'
}

const handleProviderChange = (key: string) => {
    store.switchProvider(key)
    // Update local reactive state from store
    activeProvider.value = store.activeProvider
    config.apiBaseUrl = store.apiBaseUrl
    config.apiKey = store.apiKey
    config.model = store.model
    config.visionApiBaseUrl = store.visionApiBaseUrl
    config.visionApiKey = store.visionApiKey
    config.visionModel = store.visionModel
    
    testStatus.value = 'IDLE'
    visionTestStatus.value = 'IDLE'
}

const handleTestConnection = async () => {
    testStatus.value = 'TESTING'
    const success = await testConnection(config)
    testStatus.value = success ? 'SUCCESS' : 'FAIL'
    setTimeout(() => { testStatus.value = 'IDLE' }, 3000)
}

const handleTestVisionConnection = async () => {
    visionTestStatus.value = 'TESTING'
    const success = await testVisionConnection(config)
    visionTestStatus.value = success ? 'SUCCESS' : 'FAIL'
    setTimeout(() => { visionTestStatus.value = 'IDLE' }, 3000)
}

// Watch for changes and save to store (debounced implicitly by user typing speed vs immediate save)
watch(config, (newVal) => {
    store.updateSettings(newVal)
}, { deep: true })

</script>

<style scoped>
/* Add specific styles if needed, mostly using Tailwind */
</style>
