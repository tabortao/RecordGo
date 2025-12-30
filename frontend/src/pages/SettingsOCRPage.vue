<template>
  <div class="p-4 space-y-4">
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" class="text-blue-500"><Document /></el-icon>
      <h2 class="font-semibold text-gray-900 dark:text-gray-100">OCR 服务配置</h2>
    </div>

    <!-- 说明卡片 -->
    <el-card shadow="never" class="!bg-blue-50 dark:!bg-blue-900/20 border-none">
      <div class="flex gap-3">
        <el-icon class="text-blue-500 mt-1"><InfoFilled /></el-icon>
        <div class="text-sm text-blue-600 dark:text-blue-300">
          配置 OCR 服务后，在“AI智能创建任务”时将优先使用 OCR 识别图片文字，提高识别准确率。
        </div>
      </div>
    </el-card>

    <!-- 模型配置 -->
    <el-card shadow="never" class="dark:bg-gray-800 dark:border-gray-700">
      <div class="font-medium mb-4 text-gray-900 dark:text-gray-100">PP-OCRv5 模型参数配置</div>
      
      <div class="space-y-4 pt-2">
        <div>
          <div class="mb-1 text-sm text-gray-600 dark:text-gray-400">API URL</div>
          <el-input v-model="configs['PP-OCRv5'].apiUrl" placeholder="请输入 API 地址" />
        </div>
        <div>
          <div class="mb-1 text-sm text-gray-600 dark:text-gray-400">Token</div>
          <el-input v-model="configs['PP-OCRv5'].token" placeholder="请输入 Token" type="password" show-password />
        </div>
      </div>
      
      <div class="mt-6 flex justify-end gap-3">
        <el-button :loading="testing" @click="testCurrentConnection">测试连接</el-button>
        <el-button type="primary" @click="saveSettings">保存配置</el-button>
      </div>
    </el-card>

  </div>
</template>

<script setup lang="ts">
import { ref, watchEffect, reactive } from 'vue'
import { ArrowLeft, Document, InfoFilled } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import { useAIStore } from '@/stores/ai'
import { ElMessage } from 'element-plus'
import type { OCRConfig } from '@/stores/ai'
import { testConnection } from '@/services/ocr'

const router = useRouter()
const store = useAIStore()

function goBack() {
  router.back()
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
