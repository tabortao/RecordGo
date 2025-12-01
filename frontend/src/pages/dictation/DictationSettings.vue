<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 flex flex-col">
    <div class="bg-white dark:bg-gray-800 shadow-sm px-4 py-3 flex items-center gap-2 sticky top-0 z-10">
      <el-icon :size="20" class="text-gray-600 dark:text-gray-300 cursor-pointer" @click="router.back()"><ArrowLeft /></el-icon>
      <span class="text-lg font-bold text-gray-800 dark:text-white">听写设置</span>
    </div>

    <div class="p-4 space-y-4">
      <el-card shadow="never" class="rounded-xl">
        <template #header><div class="font-medium">播放设置</div></template>
        <el-form label-position="left" label-width="100px">
          <el-form-item label="分词规则">
            <el-radio-group v-model="form.split_rule">
              <el-radio label="newline">按换行</el-radio>
              <el-radio label="space">按空格</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="播放模式">
            <el-radio-group v-model="form.play_mode">
              <el-radio label="dictate">听写模式</el-radio>
              <el-radio label="read">朗读模式</el-radio>
            </el-radio-group>
            <div class="text-xs text-gray-500">听写模式会有间隔；朗读模式连续播放。</div>
          </el-form-item>
          <el-form-item label="播放顺序">
            <el-radio-group v-model="form.order_mode">
              <el-radio label="sequential">顺序播放</el-radio>
              <el-radio label="random">乱序播放</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-form>
      </el-card>

      <el-card shadow="never" class="rounded-xl mb-4">
        <template #header><div class="font-medium">默认内容设置</div></template>
        <el-form label-position="left" label-width="100px">
          <el-form-item label="学段">
             <el-select v-model="form.default_education_stage" placeholder="请选择学段" style="width: 100%" @change="handleStageChange">
                <el-option label="小学" value="小学" />
                <el-option label="初中" value="初中" />
                <el-option label="高中" value="高中" />
             </el-select>
          </el-form-item>
          <el-form-item label="教材版本">
             <el-select v-model="form.default_version" placeholder="请选择版本" style="width: 100%" allow-create filterable default-first-option>
                <el-option v-for="v in versionOptions" :key="v" :label="v" :value="v" />
             </el-select>
          </el-form-item>
          <el-form-item label="年级">
             <el-select v-model="form.default_grade" placeholder="请选择年级" style="width: 100%">
                <el-option v-for="g in gradeOptions" :key="g" :label="g" :value="g" />
             </el-select>
          </el-form-item>
        </el-form>
      </el-card>

      <el-card shadow="never" class="rounded-xl">
        <template #header><div class="font-medium">参数调节</div></template>
        <el-form label-position="left" label-width="100px">
          <el-form-item label="重复次数">
            <el-input-number v-model="form.repeat_count" :min="1" :max="10" />
            <span class="ml-2 text-gray-500">次</span>
          </el-form-item>
          <el-form-item label="间隔时间">
            <el-input-number v-model="form.interval_seconds" :min="1" :max="60" />
            <span class="ml-2 text-gray-500">秒</span>
          </el-form-item>
          <el-form-item label="语速">
            <el-slider v-model="form.speed" :min="0.5" :max="2.0" :step="0.1" show-stops />
            <div class="text-right text-gray-500">{{ form.speed }}x</div>
          </el-form-item>
          <el-form-item label="中文音色">
             <el-select v-model="form.zh_voice_type" placeholder="中文默认" style="width: 100%">
                <el-option label="默认（自动匹配）" value="" />
                <el-option 
                  v-for="v in zhVoices" 
                  :key="v.voiceURI" 
                  :label="`${v.name} (${v.lang})`" 
                  :value="v.voiceURI" 
                />
             </el-select>
          </el-form-item>
          <el-form-item label="英文音色">
             <el-select v-model="form.en_voice_type" placeholder="英文默认" style="width: 100%">
                <el-option label="默认（自动匹配）" value="" />
                <el-option 
                  v-for="v in enVoices" 
                  :key="v.voiceURI" 
                  :label="`${v.name} (${v.lang})`" 
                  :value="v.voiceURI" 
                />
             </el-select>
          </el-form-item>
        </el-form>
      </el-card>

      <div class="pt-4">
        <el-button type="primary" size="large" class="w-full" @click="save">保存设置</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { dictationApi, type DictationSettings } from '@/services/dictation'
import { ElMessage } from 'element-plus'

const router = useRouter()
const availableVoices = ref<SpeechSynthesisVoice[]>([])
const zhVoices = ref<SpeechSynthesisVoice[]>([])
const enVoices = ref<SpeechSynthesisVoice[]>([])

const form = ref<DictationSettings>({
  user_id: 0,
  split_rule: 'newline',
  play_mode: 'dictate',
  order_mode: 'sequential',
  repeat_count: 2,
  interval_seconds: 5,
  voice_type: '',
  zh_voice_type: '',
  en_voice_type: '',
  default_education_stage: '',
  default_version: '',
  default_grade: '',
  speed: 1.0
})

const versionOptions = ['人教版', '部编版', '北师大版', '苏教版', '沪教版', '鲁教版', '浙教版', '外研版', '译林版']

const gradeOptions = computed(() => {
  const s = form.value.default_education_stage
  if (s === '小学') {
    return ['一年级上', '一年级下', '二年级上', '二年级下', '三年级上', '三年级下', '四年级上', '四年级下', '五年级上', '五年级下', '六年级上', '六年级下']
  } else if (s === '初中') {
    return ['初一上', '初一下', '初二上', '初二下', '初三上', '初三下']
  } else if (s === '高中') {
    return ['高一上', '高一下', '高二上', '高二下', '高三上', '高三下']
  }
  return []
})

function handleStageChange() {
  form.value.default_grade = ''
}


function loadVoices() {
  const voices = window.speechSynthesis.getVoices()
  availableVoices.value = voices
  zhVoices.value = voices.filter(v => v.lang.includes('zh') || v.lang.includes('CN'))
  enVoices.value = voices.filter(v => v.lang.includes('en'))
}

onMounted(async () => {
  loadVoices()
  if (window.speechSynthesis.onvoiceschanged !== undefined) {
    window.speechSynthesis.onvoiceschanged = loadVoices
  }

  try {
    const res = await dictationApi.getSettings()
    // Fix: handle response data correctly (it might be unwrapped)
    const data = (res as any).data || (res as any)
    if (data && typeof data === 'object') {
      // Only update fields that exist in response
      form.value = { ...form.value, ...data }
    }
  } catch (e) {
    // If error (e.g. first time), use defaults
  }
})

async function save() {
  try {
    await dictationApi.updateSettings(form.value)
    ElMessage.success('设置已保存')
    router.back()
  } catch (e) {
    ElMessage.error('保存失败')
  }
}
</script>
