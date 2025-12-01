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

    <!-- Filter / Search (Advanced) -->
    <div class="p-3 bg-white dark:bg-gray-800 border-b border-gray-100 dark:border-gray-700 flex flex-col gap-2">
       <div class="flex gap-2 overflow-x-auto pb-1">
          <el-select v-model="filterStage" placeholder="学段" style="width: 100px" clearable @change="handleStageChange">
             <el-option label="全部" value="" />
             <el-option label="小学" value="小学" />
             <el-option label="初中" value="初中" />
             <el-option label="高中" value="高中" />
          </el-select>
          <el-select v-model="filterVersion" placeholder="版本" style="width: 120px" clearable filterable allow-create>
             <el-option label="全部" value="" />
             <el-option v-for="v in versionOptions" :key="v" :label="v" :value="v" />
          </el-select>
          <el-select v-model="filterGrade" placeholder="年级" style="width: 120px" clearable>
             <el-option label="全部" value="" />
             <el-option v-for="g in gradeOptions" :key="g" :label="g" :value="g" />
          </el-select>
       </div>
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
const selectedIds = ref<Set<number>>(new Set())

// Filter states
const filterStage = ref('')
const filterVersion = ref('')
const filterGrade = ref('')

const versionOptions = ['人教版', '部编版', '北师大版', '苏教版', '沪教版', '鲁教版', '浙教版', '外研版', '译林版']

const gradeOptions = computed(() => {
  const s = filterStage.value
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
  filterGrade.value = ''
}

const filteredList = computed(() => {
  let all = [...list.value]
  // Append system presets
  const presets = systemPresets.map((p, idx) => ({
    ...p,
    id: -1000 - idx,
    user_id: 0,
    created_at: new Date().toISOString(),
    is_preset: true
  } as unknown as WordBank))
  
  all = [...presets, ...all]
  
  return all.filter(i => {
      if (filterStage.value && i.education_stage !== filterStage.value) return false
      if (filterVersion.value && i.version !== filterVersion.value) return false
      if (filterGrade.value && i.grade !== filterGrade.value) return false
      return true
  })
})

onMounted(loadData)

async function loadData() {
  try {
    // Load Wordbanks
    const res = await dictationApi.listWordBanks()
    list.value = (res as any) || []

    // Load Settings for default filters
    try {
        const sRes = await dictationApi.getSettings()
        const settings = (sRes as any).data || (sRes as any)
        if (settings) {
            if (settings.default_education_stage) filterStage.value = settings.default_education_stage
            if (settings.default_version) filterVersion.value = settings.default_version
            if (settings.default_grade) filterGrade.value = settings.default_grade
        }
    } catch {}
  } catch (e) {
    ElMessage.error('加载失败')
  }
}

function extractPreview(content: string) {
  try {
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
    // Reload but keep filters
    const res = await dictationApi.listWordBanks()
    list.value = (res as any) || []
  } catch (e) {
    ElMessage.error('删除失败')
  }
}

function confirmSelection() {
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
