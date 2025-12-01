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
        </el-form>
      </el-card>

      <div class="pt-4">
        <el-button type="primary" size="large" class="w-full" @click="save">保存设置</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { dictationApi, type DictationSettings } from '@/services/dictation'
import { ElMessage } from 'element-plus'

const router = useRouter()
const form = ref<DictationSettings>({
  user_id: 0,
  split_rule: 'newline',
  play_mode: 'dictate',
  order_mode: 'sequential',
  repeat_count: 2,
  interval_seconds: 5,
  voice_type: 'zh-CN',
  speed: 1.0
})

onMounted(async () => {
  try {
    const res = await dictationApi.getSettings()
    if (res.data) {
      form.value = { ...form.value, ...res.data }
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
