<template>
  <div class="h-full flex flex-col bg-gray-50 dark:bg-gray-900">
    <div class="bg-white dark:bg-gray-800 shadow-sm px-4 py-3 flex items-center justify-between sticky top-0 z-10">
      <div class="flex items-center gap-2">
        <el-icon :size="20" class="text-gray-600 dark:text-gray-300 cursor-pointer" @click="router.back()"><ArrowLeft /></el-icon>
        <span class="text-lg font-bold text-gray-800 dark:text-white">错题本</span>
      </div>
      <el-button type="primary" link @click="addManual">手动添加</el-button>
    </div>

    <div class="flex-1 overflow-y-auto p-3 space-y-3">
      <el-empty v-if="list.length === 0" description="暂无错题" />
      
      <div 
        v-for="item in list" 
        :key="item.id"
        class="bg-white dark:bg-gray-800 rounded-xl p-3 shadow-sm border border-gray-100 dark:border-gray-700 flex items-center justify-between"
      >
        <div>
          <div class="text-lg font-medium text-gray-900 dark:text-gray-100">{{ item.word }}</div>
          <div class="text-xs text-gray-500 mt-1" v-if="item.context">来源: {{ item.context }}</div>
          <div class="text-xs text-gray-400 mt-1">{{ new Date(item.created_at).toLocaleDateString() }}</div>
        </div>
        <el-button circle type="danger" :icon="Delete" size="small" @click="removeItem(item.id)" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Delete } from '@element-plus/icons-vue'
import { dictationApi, type Mistake } from '@/services/dictation'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const list = ref<Mistake[]>([])

onMounted(loadData)

async function loadData() {
  try {
    const res = await dictationApi.listMistakes()
    list.value = (res as any) || []
  } catch (e) {
    ElMessage.error('加载失败')
  }
}

async function removeItem(id: number) {
  try {
    await dictationApi.deleteMistake(id)
    list.value = list.value.filter(i => i.id !== id)
    ElMessage.success('已移除')
  } catch (e) {
    ElMessage.error('删除失败')
  }
}

async function addManual() {
  try {
    const { value } = await ElMessageBox.prompt('请输入错词', '手动添加', {
      confirmButtonText: '添加',
      cancelButtonText: '取消',
    })
    if (value) {
      await dictationApi.addMistake({ word: value, context: '手动添加' })
      loadData()
    }
  } catch {}
}
</script>
