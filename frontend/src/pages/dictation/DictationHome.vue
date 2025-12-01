<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 flex flex-col">
    <!-- Header -->
    <div class="bg-white dark:bg-gray-800 shadow-sm px-4 py-3 flex items-center justify-between sticky top-0 z-10">
      <div class="flex items-center gap-2">
        <el-icon :size="20" class="text-gray-600 dark:text-gray-300 cursor-pointer" @click="router.push('/homework')"><ArrowLeft /></el-icon>
        <span class="text-lg font-bold text-gray-800 dark:text-white">听写大师</span>
      </div>
      <div class="flex items-center gap-3">
        <el-icon :size="20" class="text-gray-600 dark:text-gray-300 cursor-pointer" @click="router.push('/dictation/settings')"><Setting /></el-icon>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 p-4 space-y-4 overflow-y-auto">
      
      <!-- Input Area -->
      <el-card shadow="never" class="rounded-xl">
        <template #header>
          <div class="flex items-center justify-between">
            <span class="font-medium text-gray-700 dark:text-gray-200">听写内容</span>
            <div class="flex gap-2">
              <el-button link type="primary" @click="showSelector = true">从词库选择</el-button>
              <el-button link type="danger" @click="content = ''">清空</el-button>
            </div>
          </div>
        </template>
        <el-input
          v-model="content"
          type="textarea"
          :rows="8"
          placeholder="在此输入听写内容，或从词库中选择。内容将按设置规则（空格/换行）进行分词播放。"
          resize="none"
        />
      </el-card>

      <!-- Quick Actions -->
      <div class="grid grid-cols-2 gap-3">
        <el-card shadow="hover" class="cursor-pointer rounded-xl border-none bg-blue-50 dark:bg-blue-900/20" :body-style="{ padding: '16px' }" @click="router.push('/dictation/banks')">
          <div class="flex flex-col items-center gap-2 text-blue-600 dark:text-blue-400">
            <el-icon :size="24"><Collection /></el-icon>
            <span class="font-medium">词库管理</span>
          </div>
        </el-card>
        <el-card shadow="hover" class="cursor-pointer rounded-xl border-none bg-red-50 dark:bg-red-900/20" :body-style="{ padding: '16px' }" @click="router.push('/dictation/mistakes')">
          <div class="flex flex-col items-center gap-2 text-red-600 dark:text-red-400">
            <el-icon :size="24"><Notebook /></el-icon>
            <span class="font-medium">错题本</span>
          </div>
        </el-card>
      </div>

      <!-- Start Button -->
      <div class="pt-4">
        <el-button type="primary" size="large" class="w-full h-12 text-lg rounded-xl shadow-lg shadow-blue-500/30" :disabled="!content.trim()" @click="startDictation">
          开始听写
        </el-button>
      </div>

    </div>

    <!-- Word Bank Selector Drawer -->
    <el-drawer v-model="showSelector" title="选择词库内容" direction="btt" size="80%">
      <div class="h-full flex flex-col">
        <!-- Reuse WordBankList component here or implement simple list -->
        <WordBankList selector-mode @select="onSelectContent" />
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Setting, Collection, Notebook } from '@element-plus/icons-vue'
import WordBankList from './WordBankList.vue'

const router = useRouter()
const content = ref('')
const showSelector = ref(false)

function onSelectContent(text: string) {
  if (content.value.trim()) {
    content.value += '\n' + text
  } else {
    content.value = text
  }
  showSelector.value = false
}

function startDictation() {
  // Pass content via state or query? State is better for large text.
  // However, state is lost on refresh. LocalStorage is safer.
  localStorage.setItem('dictation_current_content', content.value)
  router.push('/dictation/player')
}
</script>

<style scoped>
:deep(.el-card__header) {
  padding: 12px 16px;
}
</style>
