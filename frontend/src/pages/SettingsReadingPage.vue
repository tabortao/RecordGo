<template>
  <SettingsShell title="朗读设置" subtitle="系统语音或自定义接口" :icon="Microphone" tone="violet" container-class="max-w-4xl">
    <SettingsCard>
      <div class="space-y-4">
        <div class="flex items-center justify-between gap-4">
          <div class="min-w-0">
            <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">启用朗读</div>
            <div class="mt-1 text-xs text-gray-500 dark:text-gray-400">关闭后，任务卡片的小喇叭将不再播放语音</div>
          </div>
          <el-switch v-model="enabled" @change="saveSettings" />
        </div>

        <el-tabs v-if="enabled" v-model="activeTab" class="settings-tabs" @tab-change="onTabChange">
        <!-- Tab 1: System -->
        <el-tab-pane label="系统默认" name="system">
          <div class="space-y-6 pt-4">
             <div class="space-y-2">
                <label class="text-sm text-gray-600">朗读语音（中文/英文）</label>
                <el-select v-model="voiceURI" placeholder="选择朗读语音" filterable style="width: 100%" @change="saveSettings">
                  <el-option v-for="v in voices" :key="v.voiceURI" :label="`${v.name}（${v.lang}）`" :value="v.voiceURI" />
                </el-select>
                <div class="text-xs text-gray-500">提示：语音列表由浏览器提供，若为空请稍后重试或检查系统内置语音</div>
             </div>
             
             <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div class="space-y-2">
                  <label class="text-sm text-gray-600">语速（0.5 ~ 2.0）</label>
                  <el-input-number v-model="rate" :min="0.5" :max="2" :step="0.1" @change="saveSettings" />
                </div>
                <div class="space-y-2">
                  <label class="text-sm text-gray-600">音调（0 ~ 2）</label>
                  <el-input-number v-model="pitch" :min="0" :max="2" :step="0.1" @change="saveSettings" />
                </div>
             </div>

             <div class="pt-2">
                <el-button type="primary" @click="testSystemSpeak">测试系统朗读</el-button>
                <span class="text-xs text-gray-500 ml-2">示例：“语文，认2个汉字。”</span>
             </div>
          </div>
        </el-tab-pane>

        <!-- Tab 2: Custom -->
        <el-tab-pane label="自定义接口" name="custom">
           <div class="space-y-4 pt-4">
              <div class="flex justify-between items-center">
                 <div class="text-sm text-gray-600">选择要使用的配置</div>
                 <el-button size="small" type="primary" :icon="Plus" @click="openEditDialog()">新增配置</el-button>
              </div>

              <!-- Profile List -->
              <div v-if="customProfiles.length === 0" class="text-center py-10 text-gray-400 text-sm border-2 border-dashed rounded-lg">
                 暂无配置，请点击右上角新增
              </div>
              <div v-else class="space-y-3">
                 <div v-for="p in customProfiles" :key="p.id" 
                      class="border rounded-lg p-4 flex items-center justify-between cursor-pointer hover:border-primary transition-all duration-200"
                      :class="{'border-primary bg-blue-50 dark:bg-blue-900/20 shadow-sm': activeCustomId === p.id}"
                      @click="selectCustomProfile(p.id)">
                    <div class="flex items-center gap-3 overflow-hidden">
                       <el-radio :model-value="activeCustomId || ''" :label="p.id" class="!mr-0" @change="selectCustomProfile(p.id)">{{ '' }}</el-radio>
                       <div class="min-w-0">
                          <div class="font-medium truncate">{{ p.name }}</div>
                          <div class="text-xs text-gray-500 truncate">{{ p.apiUrl }}</div>
                          <div class="text-[10px] text-gray-400 mt-0.5 flex gap-2">
                            <span v-if="p.voiceId">Voice: {{ p.voiceId }}</span>
                            <span>Speed: {{ p.speed }}x</span>
                          </div>
                       </div>
                    </div>
                    <div class="flex gap-2 shrink-0">
                       <el-button circle size="small" :icon="Edit" @click.stop="openEditDialog(p)" />
                       <el-button circle size="small" type="danger" :icon="Delete" @click.stop="deleteProfile(p.id)" />
                    </div>
                 </div>
              </div>
              
              <div class="text-xs text-gray-500 pt-2">
                 提示：选中“自定义接口”标签页且选择了配置时，将优先使用自定义接口进行朗读。
              </div>
           </div>
        </el-tab-pane>
      </el-tabs>
      </div>
    </SettingsCard>

    <!-- Edit Dialog -->
    <el-dialog v-model="showDialog" :title="editingId ? '编辑配置' : '新增配置'" width="90%" class="max-w-lg" append-to-body destroy-on-close>
       <el-form label-position="top" class="space-y-2">
          <el-form-item label="配置名称" required>
             <el-input v-model="form.name" placeholder="例如：OpenAI TTS" />
          </el-form-item>
          <el-form-item label="接口地址 (API URL)" required>
             <el-input v-model="form.apiUrl" placeholder="https://api.openai.com/v1/audio/speech" />
             <div class="text-xs text-gray-500 mt-1 leading-tight">支持 OpenAI 格式 (/v1/audio/speech) 或 Edge-TTS 兼容接口</div>
          </el-form-item>
          <el-form-item label="API Key (可选)">
             <el-input v-model="form.apiKey" type="password" show-password placeholder="sk-..." />
          </el-form-item>
          
          <div class="grid grid-cols-2 gap-4">
             <el-form-item label="音色 ID (Voice)">
                <el-input v-model="form.voiceId" placeholder="例如: alloy / zh-CN-XiaoxiaoNeural" />
             </el-form-item>
             <el-form-item label="请求方式">
                <el-select v-model="form.method" style="width: 100%">
                   <el-option label="GET" value="GET" />
                   <el-option label="POST" value="POST" />
                </el-select>
             </el-form-item>
          </div>
          
          <div class="grid grid-cols-2 gap-4">
             <el-form-item label="语速 (Speed)">
                <el-input-number v-model="form.speed" :step="0.1" :min="0.25" :max="4.0" style="width: 100%" />
             </el-form-item>
             <el-form-item label="语调 (Pitch)">
                <el-input-number v-model="form.pitch" :step="0.1" :min="0" :max="2.0" style="width: 100%" />
             </el-form-item>
          </div>
       </el-form>
       
       <!-- Test Section -->
       <div class="bg-gray-50 dark:bg-gray-800 p-4 rounded-lg mb-2 mt-2">
          <div class="flex justify-between items-center mb-3">
             <span class="text-sm font-medium text-gray-700 dark:text-gray-300">连接测试</span>
             <el-button type="primary" size="small" :loading="isTesting" @click="testConnection">测试连接</el-button>
          </div>
          <div v-if="testResult" class="text-xs p-2 rounded border transition-all" 
               :class="testResult.success ? 'bg-green-50 border-green-200 text-green-700' : 'bg-red-50 border-red-200 text-red-700 break-all'">
             <div class="font-bold mb-1">{{ testResult.success ? '测试成功' : '测试失败' }}</div>
             {{ testResult.msg }}
          </div>
       </div>

       <template #footer>
          <div class="flex justify-end gap-2">
             <el-button @click="showDialog = false">取消</el-button>
             <el-button type="primary" @click="saveProfile">保存配置</el-button>
          </div>
       </template>
    </el-dialog>
  </SettingsShell>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Microphone, Plus, Edit, Delete } from '@element-plus/icons-vue'
import { useAppState, type CustomTTSProfile } from '@/stores/appState'
import { listVoices, speak } from '@/utils/speech'
import { ElMessage, ElMessageBox } from 'element-plus'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

const store = useAppState()
const enabled = ref<boolean>(store.speech.enabled)
const activeTab = ref(store.speech.engine === 'custom' ? 'custom' : 'system')

// System Settings
const voiceURI = ref<string | null>(store.speech.voiceURI || null)
const rate = ref<number>(store.speech.rate)
const pitch = ref<number>(store.speech.pitch)
const voices = ref<SpeechSynthesisVoice[]>([])

// Custom Settings
const customProfiles = ref<CustomTTSProfile[]>(store.speech.customProfiles || [])
const activeCustomId = ref<string | null>(store.speech.activeCustomId || null)

// Dialog State
const showDialog = ref(false)
const editingId = ref<string | null>(null)
const form = reactive<CustomTTSProfile>({
    id: '', name: '', apiUrl: '', apiKey: '', voiceId: '', speed: 1.0, pitch: 1.0, method: 'GET'
})

// Test State
const isTesting = ref(false)
const testResult = ref<{success: boolean, msg: string} | null>(null)

// Init Voices
function refreshVoices() {
  try {
    voices.value = listVoices()
    if (!voiceURI.value) {
      const zh = voices.value.find(v => /^zh(?:-|_)?/i.test(v.lang))
      if (zh) voiceURI.value = zh.voiceURI
    }
  } catch {}
}
onMounted(() => {
    refreshVoices()
    try { window.speechSynthesis?.addEventListener?.('voiceschanged', refreshVoices) } catch {}
})

// Unified Save
function saveSettings() {
    store.updateSpeech({ 
        enabled: enabled.value, 
        voiceURI: voiceURI.value || null, 
        rate: rate.value, 
        pitch: pitch.value,
        engine: activeTab.value === 'custom' ? 'custom' : 'system',
        customProfiles: customProfiles.value,
        activeCustomId: activeCustomId.value
    })
}

function onTabChange() {
    saveSettings()
}

// System Test
async function testSystemSpeak() {
  if (!enabled.value) {
    ElMessage.info('请先开启朗读开关')
    return
  }
  const ok = await speak('语文，认2个汉字。', { voiceURI: voiceURI.value || undefined, rate: rate.value, pitch: pitch.value })
  if (!ok) ElMessage.warning('浏览器暂不支持朗读或语音不可用')
}

// Custom Profiles Logic
function selectCustomProfile(id: string) {
    activeCustomId.value = id
    saveSettings()
}

function openEditDialog(p?: CustomTTSProfile) {
    testResult.value = null
    if (p) {
        editingId.value = p.id
        Object.assign(form, p)
    } else {
        editingId.value = null
        // Defaults
        Object.assign(form, {
            id: '', 
            name: '', 
            apiUrl: '', 
            apiKey: '', 
            voiceId: '', 
            speed: 1.0, 
            pitch: 1.0, 
            method: 'GET'
        })
    }
    showDialog.value = true
}

function deleteProfile(id: string) {
    ElMessageBox.confirm('确定要删除这个配置吗？', '提示', { type: 'warning' }).then(() => {
        customProfiles.value = customProfiles.value.filter(p => p.id !== id)
        if (activeCustomId.value === id) activeCustomId.value = null
        saveSettings()
    }).catch(() => {})
}

function saveProfile() {
    if (!form.name || !form.apiUrl) {
        ElMessage.warning('配置名称和接口地址不能为空')
        return
    }
    
    // URL Cleanup
    let cleanUrl = form.apiUrl.trim()
    if (!cleanUrl.startsWith('http')) cleanUrl = 'https://' + cleanUrl
    
    const newProfile: CustomTTSProfile = {
        ...form,
        apiUrl: cleanUrl,
        id: editingId.value || Date.now().toString()
    }
    
    if (editingId.value) {
        const idx = customProfiles.value.findIndex(p => p.id === editingId.value)
        if (idx > -1) customProfiles.value[idx] = newProfile
    } else {
        customProfiles.value.push(newProfile)
        // Auto select if first one
        if (customProfiles.value.length === 1) activeCustomId.value = newProfile.id
    }
    
    saveSettings()
    showDialog.value = false
    ElMessage.success('配置已保存')
}

// Custom TTS Test Logic (Adapted from Reference)
async function testConnection() {
    if (!form.apiUrl) return
    isTesting.value = true
    testResult.value = null
    const testText = "你好，我是汉字小英雄。" // Using text from reference
    
    let rawUrl = form.apiUrl.trim()
    if (!rawUrl.startsWith('http')) rawUrl = 'https://' + rawUrl
    
    // Candidates generation logic
    const candidates: { url: string, method: 'GET' | 'POST' }[] = []
    
    // 1. Exact match
    if (rawUrl.endsWith('/v1/audio/speech') || rawUrl.endsWith('/speech')) {
        candidates.push({ url: rawUrl, method: 'POST' })
    }
    candidates.push({ url: rawUrl, method: 'GET' })
    
    // 2. Intelligent guessing
    let baseUrl = rawUrl.replace(/\/v1\/audio\/speech\/?$/, '').replace(/\/api\/tts\/?$/, '').replace(/\/tts\/?$/, '').replace(/\/+$/, '')
    candidates.push({ url: `${baseUrl}/v1/audio/speech`, method: 'POST' }) // OpenAI standard
    candidates.push({ url: `${baseUrl}/api/tts`, method: 'GET' }) // Common
    candidates.push({ url: `${baseUrl}/tts`, method: 'GET' })
    
    // Deduplicate
    const uniqueCandidates = candidates.filter((v, i, a) => a.findIndex(t => t.url === v.url && t.method === v.method) === i)
    
    let successfulUrl = ''
    let successfulMethod: 'GET' | 'POST' = 'GET'
    let errorDetails = ''
    let found = false
    
    for (const candidate of uniqueCandidates) {
        try {
            const headers: HeadersInit = {}
            if (form.apiKey) headers['Authorization'] = `Bearer ${form.apiKey}`
            
            let res: Response
            if (candidate.method === 'POST') {
                headers['Content-Type'] = 'application/json'
                const body = JSON.stringify({
                    model: 'tts-1',
                    input: testText,
                    voice: form.voiceId || 'zh-CN-XiaoxiaoNeural',
                    speed: form.speed || 1.0,
                    pitch: form.pitch || 1.0, // Some proxies support this
                    response_format: 'mp3'
                })
                res = await fetch(candidate.url, { method: 'POST', headers, body })
            } else {
                const urlObj = new URL(candidate.url)
                urlObj.searchParams.set('text', testText)
                if (form.voiceId) urlObj.searchParams.set('voice', form.voiceId)
                
                // Edge-TTS style rate/pitch
                const ratePercent = Math.round(((form.speed || 1.0) - 1) * 100)
                urlObj.searchParams.set('rate', (ratePercent >= 0 ? '+' : '') + ratePercent + '%')
                
                const pitchDiff = Math.round(((form.pitch || 1.0) - 1) * 20)
                urlObj.searchParams.set('pitch', (pitchDiff >= 0 ? '+' : '') + pitchDiff + 'Hz')
                
                res = await fetch(urlObj.toString(), { headers })
            }
            
            if (res.ok) {
                const type = res.headers.get('content-type')
                // Loose check
                if (type && (type.includes('audio') || type.includes('mpeg') || type.includes('wav') || type.includes('octet-stream'))) {
                    const blob = await res.blob()
                    if (blob.size > 0) {
                        const audio = new Audio(URL.createObjectURL(blob))
                        await audio.play()
                        successfulUrl = candidate.url
                        successfulMethod = candidate.method
                        found = true
                        break
                    }
                }
            } else {
                const txt = await res.text().catch(() => '')
                if (res.status !== 404 && res.status !== 405) {
                    errorDetails = `[${candidate.method} ${res.status}] ${txt.slice(0, 100)}`
                }
            }
        } catch (e: any) {
            if (!errorDetails || candidate.url === rawUrl) errorDetails = e.message || "Network Error"
        }
    }
    
    isTesting.value = false
    if (found && successfulUrl) {
        testResult.value = { success: true, msg: "测试成功！声音正常播放。" }
        // Auto update if better path found
        if (successfulUrl !== rawUrl || successfulMethod !== form.method) {
             ElMessageBox.confirm(
                `测试成功！\n检测到更有效的接口地址：\n${successfulUrl}\n方式: ${successfulMethod}\n是否更新配置？`, 
                '优化建议'
             ).then(() => {
                 form.apiUrl = successfulUrl
                 form.method = successfulMethod
             }).catch(() => {})
        }
    } else {
        testResult.value = { success: false, msg: `测试失败: ${errorDetails || '无法连接或返回格式错误'}` }
    }
}
</script>

<style scoped>
.settings-tabs :deep(.el-tabs__header) {
  margin: 0;
}
.settings-tabs :deep(.el-tabs__nav-wrap) {
  padding: 0 2px;
}
.settings-tabs :deep(.el-tabs__nav) {
  gap: 6px;
}
.settings-tabs :deep(.el-tabs__item) {
  height: 34px;
  line-height: 34px;
  border-radius: 9999px;
  padding: 0 14px;
  font-weight: 800;
  font-size: 12px;
  color: rgb(107 114 128);
  background: rgba(255, 255, 255, 0.55);
  border: 1px solid rgba(229, 231, 235, 0.9);
}
.dark .settings-tabs :deep(.el-tabs__item) {
  color: rgb(156 163 175);
  background: rgba(17, 24, 39, 0.4);
  border: 1px solid rgba(31, 41, 55, 0.9);
}
.settings-tabs :deep(.el-tabs__item.is-active) {
  color: rgb(88 28 135);
  border-color: rgba(221, 214, 254, 0.9);
  background: rgba(237, 233, 254, 0.75);
}
.dark .settings-tabs :deep(.el-tabs__item.is-active) {
  color: rgb(233 213 255);
  border-color: rgba(88, 28, 135, 0.55);
  background: rgba(88, 28, 135, 0.18);
}
.settings-tabs :deep(.el-tabs__active-bar) {
  display: none;
}
.settings-tabs :deep(.el-tabs__item:hover) {
  color: rgb(17 24 39);
}
.dark .settings-tabs :deep(.el-tabs__item:hover) {
  color: rgb(243 244 246);
}
</style>
