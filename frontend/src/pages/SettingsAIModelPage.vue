<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 flex flex-col pb-24 relative z-20 transition-all">
    <!-- Header -->
    <div class="bg-white dark:bg-gray-800 px-4 py-3 md:px-6 md:py-6 rounded-b-[2rem] shadow-sm flex justify-between items-center z-10 sticky top-0 transition-all shrink-0">
      <div class="flex items-center gap-3">
        <button @click="onBack" class="p-2 -ml-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full transition-colors text-gray-600 dark:text-gray-400">
          <el-icon :size="24"><ArrowLeft /></el-icon>
        </button>
        <h1 class="text-xl md:text-2xl font-bold text-gray-800 dark:text-gray-100 flex items-center gap-2">
          <el-icon class="text-indigo-600 dark:text-indigo-400"><Cpu /></el-icon> AI 模型设置
        </h1>
      </div>
    </div>

    <!-- Content -->
    <div class="px-4 md:px-8 mt-6 flex-1">
      <div class="max-w-3xl mx-auto space-y-6 md:space-y-8">
        
        <!-- Provider Selection -->
        <div class="space-y-3">
          <label class="text-sm font-bold text-gray-500 dark:text-gray-400 uppercase tracking-wider flex items-center gap-2">
             <el-icon><Connection /></el-icon> 选择服务商
          </label>
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-3">
            <button
                v-for="(provider, key) in PROVIDERS"
                :key="key"
                @click="handleProviderChange(key)"
                class="px-2 py-4 rounded-xl text-[10px] md:text-xs font-bold border-2 transition-all flex flex-col items-center justify-center gap-2"
                :class="activeProvider === key ? 'bg-indigo-600 text-white border-indigo-600 shadow-lg transform scale-105' : 'bg-white dark:bg-gray-800 text-gray-600 dark:text-gray-300 border-gray-200 dark:border-gray-700 hover:border-indigo-300 dark:hover:border-indigo-500 hover:bg-indigo-50 dark:hover:bg-indigo-900/30'"
            >
                <span class="text-xl md:text-2xl">{{ getProviderEmoji(key) }}</span>
                <span class="text-center truncate w-full px-1">{{ provider.name }}</span>
            </button>
          </div>
        </div>

        <!-- Chat Model Config -->
        <div class="bg-white dark:bg-gray-800 p-5 md:p-6 rounded-3xl border border-gray-100 dark:border-gray-700 shadow-sm space-y-6">
          <div class="flex items-center justify-between border-b border-gray-50 dark:border-gray-700 pb-4">
              <div>
                  <h3 class="text-lg font-bold text-gray-800 dark:text-gray-100 flex items-center gap-2">
                      <el-icon :size="20" class="text-indigo-500 dark:text-indigo-400"><ChatDotRound /></el-icon> 对话模型 (Chat)
                  </h3>
                  <p class="text-[10px] md:text-xs text-gray-400 dark:text-gray-500 mt-1">用于文字生成、故事创作和汉字讲解。</p>
              </div>
              <span v-if="testStatus === 'SUCCESS'" class="bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 px-3 py-1 rounded-full text-[10px] font-bold flex items-center gap-1">
                <el-icon :size="14"><Select /></el-icon> 已连接
              </span>
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-6">
              <div class="col-span-1 md:col-span-2">
                <label class="block text-[10px] md:text-xs font-bold text-gray-500 dark:text-gray-400 mb-2 uppercase">API Host (代理地址)</label>
                <input
                  type="text"
                  :placeholder="activeProvider === 'GOOGLE' ? '默认无需填写' : 'https://api.example.com/v1'"
                  class="w-full p-3 md:p-4 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/50 text-gray-900 dark:text-gray-100 focus:bg-white dark:focus:bg-gray-900 focus:border-indigo-500 outline-none text-sm font-mono transition-all"
                  v-model="config.apiBaseUrl"
                />
              </div>
              <div class="col-span-1 md:col-span-2">
                <label class="block text-[10px] md:text-xs font-bold text-gray-500 dark:text-gray-400 mb-2 uppercase">API Key (密钥)</label>
                <div class="relative group">
                    <input
                      :type="showKey ? 'text' : 'password'"
                      placeholder="sk-..."
                      class="w-full p-3 md:p-4 pr-12 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/50 text-gray-900 dark:text-gray-100 focus:bg-white dark:focus:bg-gray-900 focus:border-indigo-500 outline-none text-sm font-mono transition-all"
                      v-model="config.apiKey"
                    />
                    <button @click="showKey = !showKey" class="absolute right-4 top-3 md:top-4 text-gray-400 hover:text-indigo-600 transition-colors">
                        <el-icon :size="20"><component :is="showKey ? Hide : View" /></el-icon>
                    </button>
                </div>
              </div>
              <div class="col-span-1 md:col-span-2">
                <label class="block text-[10px] md:text-xs font-bold text-gray-500 dark:text-gray-400 mb-2 uppercase">Model Name (模型)</label>
                <input
                  type="text"
                  placeholder="如: gemini-2.5-flash"
                  class="w-full p-3 md:p-4 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/50 text-gray-900 dark:text-gray-100 focus:bg-white dark:focus:bg-gray-900 focus:border-indigo-500 outline-none text-sm font-mono transition-all"
                  v-model="config.model"
                />
              </div>
          </div>

          <button
              @click="handleTestConnection"
              :disabled="testStatus === 'TESTING' || !config.apiKey"
              class="w-full py-3 md:py-4 rounded-xl text-sm font-bold flex items-center justify-center gap-2 transition-all border-2 disabled:opacity-50"
              :class="testStatus === 'SUCCESS' ? 'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-900/50 text-green-600 dark:text-green-400' : testStatus === 'FAIL' ? 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-900/50 text-red-600 dark:text-red-400' : 'bg-white dark:bg-gray-800 border-indigo-100 dark:border-indigo-900/50 text-indigo-600 dark:text-indigo-400 hover:bg-indigo-50 dark:hover:bg-indigo-900/30'"
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
        </div>

        <!-- Vision Model Config -->
        <div class="bg-white dark:bg-gray-800 p-5 md:p-6 rounded-3xl border border-gray-100 dark:border-gray-700 shadow-sm space-y-6">
            <div class="flex items-center justify-between border-b border-gray-50 dark:border-gray-700 pb-4">
                <div>
                    <h3 class="text-lg font-bold text-gray-800 dark:text-gray-100 flex items-center gap-2">
                        <el-icon :size="20" class="text-purple-500 dark:text-purple-400"><Picture /></el-icon> 视觉模型 (Vision)
                    </h3>
                    <p class="text-[10px] md:text-xs text-gray-400 dark:text-gray-500 mt-1">用于拍照识字和看图写话。</p>
                </div>
                <span v-if="visionTestStatus === 'SUCCESS'" class="bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 px-3 py-1 rounded-full text-[10px] font-bold flex items-center gap-1">
                    <el-icon :size="14"><Select /></el-icon> 已连接
                </span>
            </div>
            <div class="grid grid-cols-1 gap-4">
                <div>
                    <label class="block text-[10px] md:text-xs font-bold text-gray-500 dark:text-gray-400 mb-2 uppercase">Vision API Host</label>
                    <input type="text" placeholder="同上 (默认)" class="w-full p-3 md:p-4 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/50 text-gray-900 dark:text-gray-100 focus:bg-white dark:focus:bg-gray-900 focus:border-purple-500 outline-none text-sm font-mono" v-model="config.visionApiBaseUrl" />
                </div>
                <div>
                    <label class="block text-[10px] md:text-xs font-bold text-gray-500 dark:text-gray-400 mb-2 uppercase">Vision API Key</label>
                    <div class="relative group">
                        <input :type="showVisionKey ? 'text' : 'password'" placeholder="同上 (默认)" class="w-full p-3 md:p-4 pr-12 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/50 text-gray-900 dark:text-gray-100 focus:bg-white dark:focus:bg-gray-900 focus:border-purple-500 outline-none text-sm font-mono" v-model="config.visionApiKey" />
                         <button @click="showVisionKey = !showVisionKey" class="absolute right-4 top-3 md:top-4 text-gray-400 hover:text-purple-600 transition-colors">
                            <el-icon :size="20"><component :is="showVisionKey ? Hide : View" /></el-icon>
                        </button>
                    </div>
                </div>
                <div>
                    <label class="block text-[10px] md:text-xs font-bold text-gray-500 dark:text-gray-400 mb-2 uppercase">Vision Model Name</label>
                    <input type="text" placeholder="默认: gemini-2.5-flash" class="w-full p-3 md:p-4 rounded-xl border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/50 text-gray-900 dark:text-gray-100 focus:bg-white dark:focus:bg-gray-900 focus:border-purple-500 outline-none text-sm font-mono" v-model="config.visionModel" />
                </div>
            </div>
            <button
                @click="handleTestVisionConnection"
                :disabled="visionTestStatus === 'TESTING' || (!config.apiKey && !config.visionApiKey)"
                class="w-full py-3 md:py-4 rounded-xl text-sm font-bold flex items-center justify-center gap-2 transition-all border-2 disabled:opacity-50"
                :class="visionTestStatus === 'SUCCESS' ? 'bg-green-50 dark:bg-green-900/20 border-green-200 dark:border-green-900/50 text-green-600 dark:text-green-400' : visionTestStatus === 'FAIL' ? 'bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-900/50 text-red-600 dark:text-red-400' : 'bg-white dark:bg-gray-800 border-purple-100 dark:border-purple-900/50 text-purple-600 dark:text-purple-400 hover:bg-purple-50 dark:hover:bg-purple-900/30'"
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
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAIStore, PROVIDERS } from '@/stores/ai'
import { testConnection, testVisionConnection } from '@/services/aiService'
import { ArrowLeft, Cpu, Connection, ChatDotRound, Select, View, Hide, Loading, Link, CloseBold, Picture } from '@element-plus/icons-vue'

const router = useRouter()
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

const onBack = () => {
    router.back()
}

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
