<template>
  <div class="h-full flex flex-col bg-gray-50 dark:bg-gray-900">
    <!-- Header (only if not selector mode, as selector is inside drawer) -->
    <div v-if="!selectorMode" class="bg-white dark:bg-gray-800 shadow-sm px-4 py-3 flex items-center justify-between sticky top-0 z-10">
      <div class="flex items-center gap-2">
        <el-icon :size="20" class="text-gray-600 dark:text-gray-300 cursor-pointer" @click="router.back()"><ArrowLeft /></el-icon>
        <span class="text-lg font-bold text-gray-800 dark:text-white">词库管理</span>
      </div>
      <el-button type="primary" size="small" @click="router.push('/dictation/banks/create')">新建</el-button>
    </div>

    <!-- Filter / Search (Simplified) -->
    <div class="p-3 bg-white dark:bg-gray-800 border-b border-gray-100 dark:border-gray-700 flex gap-2 overflow-x-auto">
       <el-tag 
         v-for="sub in ['全部', '语文', '英语']" 
         :key="sub" 
         :effect="filterSubject === sub ? 'dark' : 'plain'"
         class="cursor-pointer"
         @click="filterSubject = sub"
       >{{ sub }}</el-tag>
    </div>

    <!-- List -->
    <div class="flex-1 overflow-y-auto p-3 space-y-3">
      <el-empty v-if="filteredList.length === 0" description="暂无词库" />
      
      <div 
        v-for="wb in filteredList" 
        :key="wb.id"
        class="bg-white dark:bg-gray-800 rounded-xl p-3 shadow-sm border border-gray-100 dark:border-gray-700 flex items-center gap-3"
        @click="toggleSelection(wb)"
      >
        <!-- Checkbox for selector mode -->
        <div v-if="selectorMode" class="flex-shrink-0">
          <el-checkbox :model-value="selectedIds.has(wb.id)" @change="() => toggleSelection(wb)" @click.stop />
        </div>

        <div class="flex-1 min-w-0">
          <div class="flex items-center justify-between mb-1">
            <div class="flex items-center gap-2 truncate">
              <span class="font-medium text-gray-900 dark:text-gray-100">{{ wb.title }}</span>
              <el-tag v-if="(wb as any).is_preset" size="small" type="warning" effect="plain">系统</el-tag>
            </div>
            <el-tag size="small" type="info">{{ wb.subject }}</el-tag>
          </div>
          <div class="text-xs text-gray-500 dark:text-gray-400 truncate">
            {{ wb.education_stage }} · {{ wb.grade }} · {{ wb.version }}
          </div>
          <div class="text-xs text-gray-400 mt-1 truncate">{{ extractPreview(wb.content) }}</div>
        </div>

        <!-- Actions for management mode -->
        <div v-if="!selectorMode && !(wb as any).is_preset" class="flex-shrink-0 flex items-center gap-2">
          <el-button circle size="small" :icon="Edit" @click.stop="router.push(`/dictation/banks/${wb.id}`)" />
          <el-popconfirm title="确认删除?" @confirm="deleteBank(wb.id)">
            <template #reference>
              <el-button circle size="small" type="danger" :icon="Delete" @click.stop />
            </template>
          </el-popconfirm>
        </div>
        <!-- View only for presets in management mode -->
        <div v-if="!selectorMode && (wb as any).is_preset" class="flex-shrink-0">
           <el-button circle size="small" :icon="Edit" disabled title="系统预设不可编辑" />
        </div>
      </div>
    </div>

    <!-- Selector Footer -->
    <div v-if="selectorMode" class="p-4 bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700">
      <el-button type="primary" class="w-full" size="large" :disabled="selectedIds.size === 0" @click="confirmSelection">
        确认选择 ({{ selectedIds.size }})
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Edit, Delete } from '@element-plus/icons-vue'
import { dictationApi, type WordBank } from '@/services/dictation'
import { ElMessage } from 'element-plus'

import systemPresets from '@/assets/presets/wordbanks.json'

const props = defineProps<{
  selectorMode?: boolean
}>()

const emit = defineEmits(['select'])

const router = useRouter()
const list = ref<WordBank[]>([])
const filterSubject = ref('全部')
const selectedIds = ref<Set<number>>(new Set())

const filteredList = computed(() => {
  let all = [...list.value]
  // Append system presets with a special ID logic or marker
  const presets = systemPresets.map((p, idx) => ({
    ...p,
    id: -1000 - idx, // negative ID for presets
    user_id: 0,
    created_at: new Date().toISOString(),
    is_preset: true
  } as unknown as WordBank))
  
  all = [...presets, ...all]
  
  if (filterSubject.value === '全部') return all
  return all.filter(i => i.subject === filterSubject.value)
})

onMounted(loadData)

async function loadData() {
  try {
    const res = await dictationApi.listWordBanks()
    list.value = (res as any) || []
  } catch (e) {
    ElMessage.error('加载失败')
  }
}

function extractPreview(content: string) {
  try {
    // Try parsing as JSON array if it looks like one
    if (content.startsWith('[')) {
      const arr = JSON.parse(content)
      if (Array.isArray(arr)) return arr.slice(0, 5).join(' ')
    }
  } catch {}
  return content.slice(0, 30) + '...'
}

function toggleSelection(wb: WordBank) {
  if (!props.selectorMode) {
    // In management mode, click goes to edit
    router.push(`/dictation/banks/${wb.id}`)
    return
  }
  if (selectedIds.value.has(wb.id)) {
    selectedIds.value.delete(wb.id)
  } else {
    selectedIds.value.add(wb.id)
  }
}

async function deleteBank(id: number) {
  try {
    await dictationApi.deleteWordBank(id)
    ElMessage.success('已删除')
    loadData()
  } catch (e) {
    ElMessage.error('删除失败')
  }
}

function confirmSelection() {
  // Aggregate content
  const selectedContent = filteredList.value
    .filter(wb => selectedIds.value.has(wb.id))
    .map(wb => {
      try {
        if (wb.content.startsWith('[')) {
          const arr = JSON.parse(wb.content)
          return arr.join(' ')
        }
      } catch {}
      return wb.content
    })
    .join('\n')
  
  emit('select', selectedContent)
}
</script>
