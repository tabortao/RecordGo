<template>
  <SettingsShell title="听写设置" subtitle="播放规则、默认内容与语音参数" :icon="Setting" tone="violet" container-class="max-w-5xl">
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 items-start">
      <div class="space-y-6">
        <SettingsCard title="播放设置">
          <el-form label-position="top" class="dictation-form">
            <div class="space-y-6">
              <el-form-item label="分词规则">
                <div class="grid grid-cols-2 gap-3 w-full">
                  <div
                    class="cursor-pointer rounded-2xl border px-4 py-3 transition-all duration-200 flex items-center gap-3"
                    :class="form.split_rule === 'newline'
                      ? 'border-violet-500 bg-violet-50 dark:bg-violet-900/20 text-violet-700 dark:text-violet-300 ring-1 ring-violet-500'
                      : 'border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900/40 hover:bg-gray-50 dark:hover:bg-gray-800'"
                    @click="form.split_rule = 'newline'"
                  >
                    <div class="h-5 w-5 rounded-full border flex items-center justify-center transition-colors"
                      :class="form.split_rule === 'newline' ? 'border-violet-500 bg-violet-500' : 'border-gray-300 dark:border-gray-600'"
                    >
                      <div class="h-2 w-2 rounded-full bg-white" v-if="form.split_rule === 'newline'" />
                    </div>
                    <span class="text-sm font-bold">按换行</span>
                  </div>
                  <div
                    class="cursor-pointer rounded-2xl border px-4 py-3 transition-all duration-200 flex items-center gap-3"
                    :class="form.split_rule === 'space'
                      ? 'border-violet-500 bg-violet-50 dark:bg-violet-900/20 text-violet-700 dark:text-violet-300 ring-1 ring-violet-500'
                      : 'border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900/40 hover:bg-gray-50 dark:hover:bg-gray-800'"
                    @click="form.split_rule = 'space'"
                  >
                    <div class="h-5 w-5 rounded-full border flex items-center justify-center transition-colors"
                      :class="form.split_rule === 'space' ? 'border-violet-500 bg-violet-500' : 'border-gray-300 dark:border-gray-600'"
                    >
                      <div class="h-2 w-2 rounded-full bg-white" v-if="form.split_rule === 'space'" />
                    </div>
                    <span class="text-sm font-bold">按空格</span>
                  </div>
                </div>
              </el-form-item>

              <el-form-item label="播放模式">
                <div class="grid grid-cols-2 gap-3 w-full">
                  <div
                    class="cursor-pointer rounded-2xl border px-4 py-3 transition-all duration-200 flex flex-col gap-1"
                    :class="form.play_mode === 'dictate'
                      ? 'border-violet-500 bg-violet-50 dark:bg-violet-900/20 text-violet-700 dark:text-violet-300 ring-1 ring-violet-500'
                      : 'border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900/40 hover:bg-gray-50 dark:hover:bg-gray-800'"
                    @click="form.play_mode = 'dictate'"
                  >
                    <div class="flex items-center gap-2">
                      <div class="h-4 w-4 rounded-full border flex items-center justify-center transition-colors"
                        :class="form.play_mode === 'dictate' ? 'border-violet-500 bg-violet-500' : 'border-gray-300 dark:border-gray-600'"
                      >
                        <div class="h-1.5 w-1.5 rounded-full bg-white" v-if="form.play_mode === 'dictate'" />
                      </div>
                      <span class="text-sm font-bold">听写模式</span>
                    </div>
                    <div class="text-[10px] opacity-70 ml-6">隐藏文本内容</div>
                  </div>
                  <div
                    class="cursor-pointer rounded-2xl border px-4 py-3 transition-all duration-200 flex flex-col gap-1"
                    :class="form.play_mode === 'read'
                      ? 'border-violet-500 bg-violet-50 dark:bg-violet-900/20 text-violet-700 dark:text-violet-300 ring-1 ring-violet-500'
                      : 'border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900/40 hover:bg-gray-50 dark:hover:bg-gray-800'"
                    @click="form.play_mode = 'read'"
                  >
                    <div class="flex items-center gap-2">
                      <div class="h-4 w-4 rounded-full border flex items-center justify-center transition-colors"
                        :class="form.play_mode === 'read' ? 'border-violet-500 bg-violet-500' : 'border-gray-300 dark:border-gray-600'"
                      >
                        <div class="h-1.5 w-1.5 rounded-full bg-white" v-if="form.play_mode === 'read'" />
                      </div>
                      <span class="text-sm font-bold">朗读模式</span>
                    </div>
                    <div class="text-[10px] opacity-70 ml-6">显示文本内容</div>
                  </div>
                </div>
              </el-form-item>

              <el-form-item label="播放顺序">
                <div class="flex gap-4">
                  <el-radio-group v-model="form.order_mode" class="w-full">
                    <el-radio-button label="sequential">顺序播放</el-radio-button>
                    <el-radio-button label="random">乱序播放</el-radio-button>
                  </el-radio-group>
                </div>
              </el-form-item>
            </div>
          </el-form>
        </SettingsCard>

        <SettingsCard title="默认内容设置" description="选择常用学段/版本/年级，词库默认按此筛选。">
          <el-form label-position="top" class="dictation-form">
            <div class="grid grid-cols-1 gap-4">
              <el-form-item label="学段">
                <el-select v-model="form.default_education_stage" placeholder="请选择学段" class="w-full" size="large" @change="handleStageChange">
                  <el-option label="小学" value="小学" />
                  <el-option label="初中" value="初中" />
                  <el-option label="高中" value="高中" />
                </el-select>
              </el-form-item>
              <div class="grid grid-cols-2 gap-4">
                <el-form-item label="教材版本">
                  <el-select v-model="form.default_version" placeholder="请选择版本" class="w-full" size="large" allow-create filterable default-first-option>
                    <el-option v-for="v in versionOptions" :key="v" :label="v" :value="v" />
                  </el-select>
                </el-form-item>
                <el-form-item label="年级">
                  <el-select v-model="form.default_grade" placeholder="请选择年级" class="w-full" size="large">
                    <el-option v-for="g in gradeOptions" :key="g" :label="g" :value="g" />
                  </el-select>
                </el-form-item>
              </div>
            </div>
          </el-form>
        </SettingsCard>
      </div>

      <div class="space-y-6">
        <SettingsCard title="参数调节">
          <el-form label-position="top" class="dictation-form">
            <div class="grid grid-cols-2 gap-4">
              <el-form-item label="重复次数">
                <el-input-number v-model="form.repeat_count" :min="1" :max="10" class="w-full" size="large" controls-position="right" />
              </el-form-item>
              <el-form-item label="间隔时间（秒）">
                <el-input-number v-model="form.interval_seconds" :min="0" :max="60" class="w-full" size="large" controls-position="right" />
              </el-form-item>
            </div>
            
            <el-form-item label="语速" class="!mt-4">
              <div class="bg-gray-50 dark:bg-gray-900/50 rounded-2xl px-5 py-4 w-full">
                <div class="flex items-center justify-between mb-2">
                  <span class="text-xs font-bold text-gray-500">0.5x</span>
                  <span class="text-lg font-black text-violet-600 dark:text-violet-400 tabular-nums">{{ form.speed }}x</span>
                  <span class="text-xs font-bold text-gray-500">2.0x</span>
                </div>
                <el-slider v-model="form.speed" :min="0.5" :max="2.0" :step="0.1" show-stops />
              </div>
            </el-form-item>

            <div class="grid grid-cols-1 gap-4 mt-4">
              <el-form-item label="中文音色">
                <el-select v-model="form.zh_voice_type" placeholder="中文默认" class="w-full" size="large">
                  <el-option label="默认（自动匹配）" value="" />
                  <el-option v-for="v in zhVoices" :key="v.voiceURI" :label="`${v.name} (${v.lang})`" :value="v.voiceURI" />
                </el-select>
              </el-form-item>
              <el-form-item label="英文音色">
                <el-select v-model="form.en_voice_type" placeholder="英文默认" class="w-full" size="large">
                  <el-option label="默认（自动匹配）" value="" />
                  <el-option v-for="v in enVoices" :key="v.voiceURI" :label="`${v.name} (${v.lang})`" :value="v.voiceURI" />
                </el-select>
              </el-form-item>
            </div>
          </el-form>
        </SettingsCard>

        <div class="sticky bottom-6 rounded-[24px] border border-white/60 dark:border-gray-800/60 bg-white/80 dark:bg-gray-900/80 backdrop-blur-xl shadow-lg shadow-gray-200/20 dark:shadow-black/20 p-4 flex justify-end gap-3 z-10">
          <el-button class="!rounded-xl !h-10 !px-6 !font-bold !text-gray-600 dark:!text-gray-300" @click="router.back()">取消</el-button>
          <el-button type="primary" class="!rounded-xl !h-10 !px-8 !font-bold shadow-md shadow-violet-500/20" @click="save">保存设置</el-button>
        </div>
      </div>
    </div>
  </SettingsShell>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Setting } from '@element-plus/icons-vue'
import { dictationApi, type DictationSettings } from '@/services/dictation'
import { ElMessage } from 'element-plus'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

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

<style scoped>
:deep(.dictation-form .el-form-item__label) {
  font-size: 12px;
  font-weight: 800;
  color: rgb(107 114 128);
  margin-bottom: 8px;
}

.dark :deep(.dictation-form .el-form-item__label) {
  color: rgb(156 163 175);
}

:deep(.dictation-form .el-input__wrapper),
:deep(.dictation-form .el-select__wrapper),
:deep(.dictation-form .el-input-number__decrease),
:deep(.dictation-form .el-input-number__increase) {
  border-radius: 16px;
  border: 1px solid rgb(229 231 235);
  background: rgb(249 250 251);
  box-shadow: none !important;
  transition: all 0.2s;
}

.dark :deep(.dictation-form .el-input__wrapper),
.dark :deep(.dictation-form .el-select__wrapper),
.dark :deep(.dictation-form .el-input-number__decrease),
.dark :deep(.dictation-form .el-input-number__increase) {
  border: 1px solid rgb(51 65 85);
  background: rgb(17 24 39 / 0.5);
  color: rgb(243 244 246);
}

:deep(.dictation-form .el-input__wrapper.is-focus),
:deep(.dictation-form .el-select__wrapper.is-focused) {
  background: white;
  border-color: rgb(139 92 246);
  box-shadow: 0 0 0 4px rgb(139 92 246 / 0.1) !important;
}

.dark :deep(.dictation-form .el-input__wrapper.is-focus),
.dark :deep(.dictation-form .el-select__wrapper.is-focused) {
  background: rgb(3 7 18);
  border-color: rgb(139 92 246);
  box-shadow: 0 0 0 4px rgb(139 92 246 / 0.1) !important;
}

:deep(.dictation-form .el-input__inner) {
  font-weight: 600;
  color: rgb(17 24 39);
}

.dark :deep(.dictation-form .el-input__inner) {
  color: rgb(243 244 246);
}
</style>
