<template>
  <div class="min-h-screen bg-[#FDF6F8] dark:bg-gray-900 pb-10 font-sans">
    <!-- Header -->
    <div class="sticky top-0 z-20 bg-white/80 dark:bg-gray-800/80 backdrop-blur-md border-b border-pink-100 dark:border-gray-700 flex items-center justify-between p-4 shadow-sm">
      <div class="flex items-center gap-3">
        <el-icon :size="22" class="cursor-pointer text-gray-600 dark:text-gray-300 hover:text-pink-500 transition" @click="router.back()"><ArrowLeft /></el-icon>
        <h1 class="text-xl font-bold text-gray-800 dark:text-white">古诗词库</h1>
      </div>
      <el-button type="primary" class="!bg-pink-500 !border-pink-500" @click="router.push('/homework/chinese/poetry/edit/new')">
        <el-icon class="mr-1"><Plus /></el-icon>
        添加诗词
      </el-button>
    </div>

    <div class="max-w-4xl mx-auto p-4 space-y-6">
      <!-- Search -->
      <div class="relative">
        <el-input
          v-model="searchQuery"
          placeholder="搜索诗词、作者..."
          class="!h-12 !text-base custom-input"
          clearable
        >
          <template #prefix>
            <el-icon class="text-gray-400"><Search /></el-icon>
          </template>
        </el-input>
      </div>

      <!-- Tabs -->
      <div class="bg-white dark:bg-gray-800 rounded-2xl p-1 shadow-sm border border-gray-100 dark:border-gray-700 flex">
        <button 
          v-for="tab in tabs" 
          :key="tab.value"
          @click="activeTab = tab.value"
          class="flex-1 py-2.5 text-sm font-medium rounded-xl transition-all relative z-10"
          :class="activeTab === tab.value ? 'bg-pink-50 text-pink-600 dark:bg-pink-900/30 dark:text-pink-300 shadow-sm' : 'text-gray-500 hover:text-gray-700 dark:text-gray-400'"
        >
          {{ tab.label }}
        </button>
      </div>

      <!-- List -->
      <div class="space-y-4">
        <div 
          v-for="poem in filteredPoems" 
          :key="poem.id"
          class="rounded-2xl p-6 shadow-sm border border-gray-100 dark:border-gray-700 hover:shadow-md transition group relative cursor-pointer"
          :class="isMastered(poem.id) ? 'bg-green-50 dark:bg-green-900/20' : 'bg-white dark:bg-gray-800'"
          @click="router.push(`/homework/chinese/poetry/${poem.id}`)"
        >
          <div class="flex justify-between items-start mb-2">
            <div class="flex items-center gap-2">
              <h3 class="text-lg font-bold text-gray-800 dark:text-gray-100 group-hover:text-pink-600 transition">{{ poem.title_cns }}</h3>
              <span v-if="getPoemTag(poem)" class="bg-gray-100 dark:bg-gray-700 text-gray-500 text-xs px-2 py-0.5 rounded-full">{{ getPoemTag(poem) }}</span>
            </div>
            <div class="flex items-center gap-2">
               <span class="px-2 py-1 rounded-md text-xs font-medium" :class="getStatusClass(poem.id)">
                  {{ getStatusText(poem.id) }}
               </span>
            </div>
          </div>
          
          <div class="text-sm text-gray-500 dark:text-gray-400 mb-3">
            {{ poem.dynasty_cns }} · {{ poem.author_cns }}
          </div>

          <p class="text-gray-600 dark:text-gray-300 text-base leading-relaxed line-clamp-2 font-serif">
            {{ poem.paragraphs_cns.join('，') }}
          </p>

          <!-- Action Buttons (Visible on Hover or Mobile) -->
          <div class="absolute right-4 bottom-4 flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity" @click.stop>
            <el-button v-if="poem.isCustom" circle size="small" @click.stop="router.push(`/homework/chinese/poetry/edit/${poem.id}`)">
              <el-icon><Edit /></el-icon>
            </el-button>
            <el-button v-if="poem.isCustom" circle size="small" type="danger" plain @click.stop="handleDelete(poem)">
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </div>

        <div v-if="filteredPoems.length === 0" class="text-center py-12 text-gray-400">
           <el-icon :size="48" class="mb-2 opacity-20"><Reading /></el-icon>
           <p>暂无相关诗词</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Plus, Search, Edit, Delete, Reading } from '@element-plus/icons-vue'
import { usePoetryStore } from '../stores/poetryStore'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Poem } from '../types'

const router = useRouter()
const poetryStore = usePoetryStore()

const searchQuery = ref('')
const activeTab = ref('all')

const tabs = [
  { label: '全部', value: 'all' },
  { label: '小学', value: 'primary' },
  { label: '初中', value: 'middle' },
  { label: '高中', value: 'high' },
  { label: '其他', value: 'other' }
]

onMounted(() => {
  poetryStore.init()
})

const filteredPoems = computed(() => {
  let poems = poetryStore.poems
  
  // Filter by Tab
  if (activeTab.value !== 'all') {
    if (activeTab.value === 'other') {
       poems = poems.filter(p => !p.tags || (!p.tags.includes('primary') && !p.tags.includes('middle') && !p.tags.includes('high')))
    } else {
       poems = poems.filter(p => p.tags && p.tags.includes(activeTab.value))
    }
  }

  // Filter by Search
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    poems = poems.filter(p => 
      p.title_cns.includes(q) || 
      p.author_cns.includes(q) ||
      p.paragraphs_cns.some(line => line.includes(q))
    )
  }

  return poems
})

const getPoemTag = (poem: Poem) => {
    if (!poem.tags) return ''
    if (poem.tags.includes('primary')) return '小学'
    if (poem.tags.includes('middle')) return '初中'
    if (poem.tags.includes('high')) return '高中'
    return ''
}

const getStatusText = (id: number) => {
    const record = poetryStore.getRecord(id)
    if (!record || record.timesStudied === 0) return '未开始'
    if (record.isMastered) return '已掌握'
    return '学习中'
}

const getStatusClass = (id: number) => {
    const record = poetryStore.getRecord(id)
    if (!record || record.timesStudied === 0) return 'bg-gray-100 text-gray-400 dark:bg-gray-700 dark:text-gray-500'
    if (record.isMastered) return 'bg-green-100 text-green-600 dark:bg-green-900/30 dark:text-green-400'
    return 'bg-blue-100 text-blue-600 dark:bg-blue-900/30 dark:text-blue-400'
}

const isMastered = (id: number) => {
    const record = poetryStore.getRecord(id)
    return record && record.isMastered
}

const handleDelete = (poem: Poem) => {
    ElMessageBox.confirm(`确定要删除《${poem.title_cns}》吗？`, '提示', {
        type: 'warning',
        confirmButtonText: '删除',
        cancelButtonText: '取消'
    }).then(() => {
        poetryStore.deleteCustomPoem(poem.id)
        ElMessage.success('删除成功')
    })
}
</script>

<style scoped>
.custom-input :deep(.el-input__wrapper) {
  border-radius: 999px;
  background-color: white;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}
.dark .custom-input :deep(.el-input__wrapper) {
  background-color: #1f2937;
  box-shadow: none;
  border: 1px solid #374151;
}
</style>