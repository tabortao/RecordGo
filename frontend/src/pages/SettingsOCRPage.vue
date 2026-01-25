<template>
  <SettingsShell title="OCR 服务" subtitle="识别图片文字，提高创建任务准确率" :icon="Document" tone="blue">
    <SettingsCard>
      <div class="flex gap-3">
        <div class="flex h-9 w-9 items-center justify-center rounded-2xl border border-blue-100/70 dark:border-blue-900/40 bg-blue-50/80 dark:bg-blue-900/25 text-blue-700 dark:text-blue-300">
          <el-icon :size="18"><InfoFilled /></el-icon>
        </div>
        <div class="min-w-0">
          <div class="text-sm font-bold text-gray-900 dark:text-gray-50">说明</div>
          <div class="mt-1 text-sm text-gray-600 dark:text-gray-300">
            配置 OCR 服务后，在“AI智能创建任务”时会优先使用 OCR 识别图片文字。
          </div>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard title="PP-OCRv5 模型参数配置" description="填写 API 地址与 Token 后可进行连接测试">
      <div class="space-y-4">
        <div class="space-y-2">
          <div class="text-xs font-semibold text-gray-500 dark:text-gray-400">API URL</div>
          <el-input v-model="configs['PP-OCRv5'].apiUrl" placeholder="请输入 API 地址" />
        </div>
        <div class="space-y-2">
          <div class="text-xs font-semibold text-gray-500 dark:text-gray-400">Token</div>
          <el-input v-model="configs['PP-OCRv5'].token" placeholder="请输入 Token" type="password" show-password />
          <button
            type="button"
            class="inline-flex items-center gap-1.5 text-xs font-semibold text-blue-600 dark:text-blue-300 hover:underline"
            @click="openLink('https://aistudio.baidu.com/paddleocr')"
          >
            <el-icon><Link /></el-icon>
            点击这里获取密钥
          </button>
        </div>

        <div class="flex justify-end gap-2 pt-2">
          <el-button class="!rounded-xl" :loading="testing" @click="testCurrentConnection">测试连接</el-button>
          <el-button type="primary" class="!rounded-xl" @click="saveSettings">保存配置</el-button>
        </div>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { ref, watchEffect, reactive } from 'vue'
import { Document, InfoFilled } from '@element-plus/icons-vue'
import { useAIStore } from '@/stores/ai'
import { ElMessage } from 'element-plus'
import type { OCRConfig } from '@/stores/ai'
import { testConnection } from '@/services/ocr'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

const store = useAIStore()

function openLink(url: string) {
  window.open(url, '_blank')
}

const configs = reactive<Record<string, OCRConfig>>({
  'PP-OCRv5': { apiUrl: '', token: '' },
  'PaddleOCR-VL': { apiUrl: '', token: '' },
  'PP-StructureV3': { apiUrl: '', token: '' }
})

// 初始化数据
watchEffect(() => {
  if (store.ocr && store.ocr.configs) {
    Object.assign(configs, JSON.parse(JSON.stringify(store.ocr.configs)))
  }
})

function saveSettings() {
  store.updateSettings({
    ocr: {
      activeModel: 'PP-OCRv5',
      configs: {
        'PP-OCRv5': { ...configs['PP-OCRv5'] },
        'PaddleOCR-VL': { ...configs['PaddleOCR-VL'] },
        'PP-StructureV3': { ...configs['PP-StructureV3'] }
      }
    }
  })
  ElMessage.success('OCR 配置已保存')
}

const testing = ref(false)
async function testCurrentConnection() {
  const model = 'PP-OCRv5'
  const config = configs[model]
  
  if (!config.apiUrl || !config.token) {
    ElMessage.warning(`请先填写 ${model} 的 API URL 和 Token`)
    return
  }
  
  testing.value = true
  try {
    const success = await testConnection(model, config)
    if (success) {
      ElMessage.success(`${model} 连接测试成功`)
    } else {
      ElMessage.error(`${model} 连接测试失败`)
    }
  } catch (e: any) {
    ElMessage.error(`${model} 连接测试失败: ${e.message || '未知错误'}`)
  } finally {
    testing.value = false
  }
}
</script>

<style scoped>
</style>
