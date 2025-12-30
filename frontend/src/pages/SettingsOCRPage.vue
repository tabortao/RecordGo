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
          <br/>支持 PP-OCRv5、PaddleOCR-VL、PP-StructureV3 三种模型。
        </div>
      </div>
    </el-card>

    <!-- 当前使用的模型 -->
    <el-card shadow="never" class="dark:bg-gray-800 dark:border-gray-700">
      <div class="flex items-center justify-between mb-4">
        <span class="font-medium text-gray-900 dark:text-gray-100">当前使用的 OCR 模型</span>
      </div>
      <el-radio-group v-model="activeModel" class="flex flex-col gap-2 items-start">
        <el-radio label="PP-OCRv5" size="large">PP-OCRv5 (通用OCR)</el-radio>
        <el-radio label="PaddleOCR-VL" size="large">PaddleOCR-VL (视觉语言模型)</el-radio>
        <el-radio label="PP-StructureV3" size="large">PP-StructureV3 (文档结构化)</el-radio>
      </el-radio-group>
    </el-card>

    <!-- 模型配置 -->
    <el-card shadow="never" class="dark:bg-gray-800 dark:border-gray-700">
      <div class="font-medium mb-4 text-gray-900 dark:text-gray-100">模型参数配置</div>
      
      <el-tabs v-model="activeTab" type="card" class="demo-tabs">
        <el-tab-pane label="PP-OCRv5" name="PP-OCRv5">
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
        </el-tab-pane>
        
        <el-tab-pane label="PaddleOCR-VL" name="PaddleOCR-VL">
          <div class="space-y-4 pt-2">
            <div>
              <div class="mb-1 text-sm text-gray-600 dark:text-gray-400">API URL</div>
              <el-input v-model="configs['PaddleOCR-VL'].apiUrl" placeholder="请输入 API 地址" />
            </div>
            <div>
              <div class="mb-1 text-sm text-gray-600 dark:text-gray-400">Token</div>
              <el-input v-model="configs['PaddleOCR-VL'].token" placeholder="请输入 Token" type="password" show-password />
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane label="PP-StructureV3" name="PP-StructureV3">
          <div class="space-y-4 pt-2">
            <div>
              <div class="mb-1 text-sm text-gray-600 dark:text-gray-400">API URL</div>
              <el-input v-model="configs['PP-StructureV3'].apiUrl" placeholder="请输入 API 地址" />
            </div>
            <div>
              <div class="mb-1 text-sm text-gray-600 dark:text-gray-400">Token</div>
              <el-input v-model="configs['PP-StructureV3'].token" placeholder="请输入 Token" type="password" show-password />
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
      
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

const activeModel = ref<'PP-OCRv5' | 'PaddleOCR-VL' | 'PP-StructureV3'>('PP-OCRv5')
const activeTab = ref('PP-OCRv5')
const configs = reactive<Record<string, OCRConfig>>({
  'PP-OCRv5': { apiUrl: '', token: '' },
  'PaddleOCR-VL': { apiUrl: '', token: '' },
  'PP-StructureV3': { apiUrl: '', token: '' }
})

// 初始化数据
watchEffect(() => {
  if (store.ocr) {
    activeModel.value = store.ocr.activeModel
    // 确保 activeTab 跟随 activeModel，或者保持独立也可以，这里让它默认显示 activeModel 的 tab
    activeTab.value = store.ocr.activeModel
    if (store.ocr.configs) {
      Object.assign(configs, JSON.parse(JSON.stringify(store.ocr.configs)))
    }
  }
})

function saveSettings() {
  store.updateSettings({
    ocr: {
      activeModel: activeModel.value,
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
  // 测试当前 Tab 对应的配置，而不是 activeModel，因为用户可能正在配置另一个模型
  // 但需求通常是测试“当前正在编辑”的模型。
  // 考虑到 Tabs 和 Radio 是分离的（Radio 选择生效模型，Tabs 编辑配置），
  // 这里逻辑是：测试当前 Tab 页签下的配置。
  
  const model = activeTab.value as 'PP-OCRv5' | 'PaddleOCR-VL' | 'PP-StructureV3'
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
:deep(.el-tabs__item) {
  color: #4b5563; /* text-gray-600 */
}
:deep(.el-tabs__item.is-active) {
  color: #3b82f6; /* text-blue-500 */
}
.dark :deep(.el-tabs__item) {
  color: #9ca3af; /* text-gray-400 */
}
.dark :deep(.el-tabs__item.is-active) {
  color: #60a5fa; /* text-blue-400 */
}
</style>
