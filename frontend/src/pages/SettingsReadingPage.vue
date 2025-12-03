<template>
  <div class="p-4 space-y-4">
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer" style="color:#64748b" @click="onBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" style="color:#7c3aed"><Microphone /></el-icon>
      <h2 class="font-semibold">朗读设置</h2>
    </div>

    <el-card shadow="never" class="space-y-4">
      <!-- 中文注释：开启/关闭朗读 -->
      <div class="flex items-center justify-between">
        <div>
          <div class="font-medium">启用朗读</div>
          <div class="text-xs text-gray-500 mt-1">关闭后，任务卡片的小喇叭将不再播放语音</div>
        </div>
        <el-switch v-model="enabled" />
      </div>

      <!-- 中文注释：语音列表（中文优先） -->
      <div class="space-y-2">
        <label class="text-sm text-gray-600">朗读语音（中文/英文）</label>
        <el-select v-model="voiceURI" placeholder="选择朗读语音" filterable style="width: 100%">
          <el-option v-for="v in voices" :key="v.voiceURI" :label="`${v.name}（${v.lang}）`" :value="v.voiceURI" />
        </el-select>
        <div class="text-xs text-gray-500">提示：语音列表由浏览器提供，若为空请稍后重试或检查系统内置语音</div>
      </div>

      <!-- 中文注释：语速与音调 -->
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div class="space-y-2">
          <label class="text-sm text-gray-600">语速（0.5 ~ 2.0）</label>
          <el-input-number v-model="rate" :min="0.5" :max="2" :step="0.1" />
        </div>
        <div class="space-y-2">
          <label class="text-sm text-gray-600">音调（0 ~ 2）</label>
          <el-input-number v-model="pitch" :min="0" :max="2" :step="0.1" />
        </div>
      </div>

      <!-- 中文注释：测试朗读按钮，依据当前设置播放示例文本 -->
      <div>
        <el-button type="primary" @click="testSpeak">测试朗读</el-button>
        <span class="text-xs text-gray-500 ml-2">示例：“语文，认2个汉字，备注：洪恩识字认识2个新汉字。”</span>
      </div>
    </el-card>
    <!-- 中文注释：统一使用左上角返回图标，返回时自动保存设置 -->
  </div>
  
  </template>

<script setup lang="ts">
import { Microphone, ArrowLeft } from '@element-plus/icons-vue'
import router from '@/router'
import { useAppState } from '@/stores/appState'
import { listVoices, speak } from '@/utils/speech'
import { ElMessage } from 'element-plus'

// 中文注释：读取与编辑全局朗读设置
const store = useAppState()
const enabled = ref<boolean>(store.speech.enabled)
const voiceURI = ref<string | null>(store.speech.voiceURI || null)
const rate = ref<number>(store.speech.rate)
const pitch = ref<number>(store.speech.pitch)

// 中文注释：语音列表；监听 voiceschanged 以异步刷新
const voices = ref<SpeechSynthesisVoice[]>([])
function refreshVoices() {
  try {
    voices.value = listVoices()
    // 中文注释：若当前未选择语音，自动选择中文语音
    if (!voiceURI.value) {
      const zh = voices.value.find(v => /^zh(?:-|_)?/i.test(v.lang))
      if (zh) voiceURI.value = zh.voiceURI
    }
  } catch {}
}
refreshVoices()
try { window.speechSynthesis?.addEventListener?.('voiceschanged', refreshVoices) } catch {}

function onBack() {
  // 中文注释：返回时保存到全局状态并持久化
  store.updateSpeech({ enabled: enabled.value, voiceURI: voiceURI.value || null, rate: rate.value, pitch: pitch.value })
  ElMessage.success('朗读设置已保存')
  router.back()
}

function testSpeak() {
  if (!enabled.value) {
    ElMessage.info('请先开启朗读开关')
    return
  }
  const ok = speak('语文，认2个汉字，备注：洪恩识字认识2个新汉字。', { voiceURI: voiceURI.value || undefined, rate: rate.value, pitch: pitch.value })
  if (!ok) ElMessage.warning('浏览器暂不支持朗读或语音不可用')
}
</script>

<style scoped>
</style>
