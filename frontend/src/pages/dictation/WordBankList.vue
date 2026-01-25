<template>
  <div v-if="selectorMode" class="h-full flex flex-col">
    <div class="p-3">
      <div class="rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/80 dark:bg-gray-900/65 backdrop-blur-xl shadow-sm px-4 py-4">
        <div class="flex gap-2 overflow-x-auto pb-1">
          <el-select v-model="filterStage" placeholder="学段" style="width: 110px" clearable @change="handleStageChange">
            <el-option label="全部" value="" />
            <el-option label="小学" value="小学" />
            <el-option label="初中" value="初中" />
            <el-option label="高中" value="高中" />
          </el-select>
          <el-select v-model="filterVersion" placeholder="版本" style="width: 140px" clearable filterable allow-create>
            <el-option label="全部" value="" />
            <el-option v-for="v in versionOptions" :key="v" :label="v" :value="v" />
          </el-select>
          <el-select v-model="filterGrade" placeholder="年级" style="width: 140px" clearable>
            <el-option label="全部" value="" />
            <el-option v-for="g in gradeOptions" :key="g" :label="g" :value="g" />
          </el-select>
        </div>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto px-3 pb-3 space-y-3">
      <el-empty v-if="filteredList.length === 0" description="暂无词库" />

      <div
        v-for="wb in filteredList"
        :key="wb.id"
        class="group rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/55 backdrop-blur px-4 py-4 shadow-sm hover:shadow-md transition flex items-center gap-3 cursor-pointer"
        @click="toggleSelection(wb)"
      >
        <div class="flex-shrink-0">
          <el-checkbox :model-value="selectedIds.has(wb.id)" @change="() => toggleSelection(wb)" @click.stop />
        </div>

        <div class="flex-1 min-w-0">
          <div class="flex items-center justify-between gap-3">
            <div class="min-w-0 flex items-center gap-2">
              <span class="truncate text-base font-extrabold tracking-tight text-gray-900 dark:text-gray-50">{{ wb.title }}</span>
              <el-tag v-if="(wb as any).is_preset" size="small" type="warning" effect="plain" class="!rounded-full">系统</el-tag>
            </div>
            <el-tag size="small" type="info" class="!rounded-full">{{ wb.subject }}</el-tag>
          </div>
          <div class="mt-1 text-xs text-gray-500 dark:text-gray-400 truncate">
            {{ wb.education_stage }} · {{ wb.grade }} · {{ wb.version }}
          </div>
          <div class="mt-1 text-sm text-gray-600 dark:text-gray-300 truncate">{{ extractPreview(wb.content) }}</div>
        </div>
      </div>
    </div>

    <div class="p-4 border-t border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-950/30 backdrop-blur">
      <el-button type="primary" class="w-full !rounded-2xl !h-12 !font-extrabold" size="large" :disabled="selectedIds.size === 0" @click="confirmSelection">
        确认选择 ({{ selectedIds.size }})
      </el-button>
    </div>
  </div>

  <SettingsShell v-else title="词库管理" subtitle="系统预设 + 自建词库" :icon="Collection" tone="blue" container-class="max-w-5xl">
    <template #headerRight>
      <el-button type="primary" class="!rounded-2xl !font-extrabold" @click="router.push('/dictation/banks/create')">新建</el-button>
    </template>

    <SettingsCard title="筛选">
      <div class="flex gap-2 flex-wrap">
        <el-select v-model="filterStage" placeholder="学段" style="width: 110px" clearable @change="handleStageChange">
          <el-option label="全部" value="" />
          <el-option label="小学" value="小学" />
          <el-option label="初中" value="初中" />
          <el-option label="高中" value="高中" />
        </el-select>
        <el-select v-model="filterVersion" placeholder="版本" style="width: 140px" clearable filterable allow-create>
          <el-option label="全部" value="" />
          <el-option v-for="v in versionOptions" :key="v" :label="v" :value="v" />
        </el-select>
        <el-select v-model="filterGrade" placeholder="年级" style="width: 140px" clearable>
          <el-option label="全部" value="" />
          <el-option v-for="g in gradeOptions" :key="g" :label="g" :value="g" />
        </el-select>
      </div>
    </SettingsCard>

    <SettingsCard title="词库列表">
      <div class="space-y-3">
        <el-empty v-if="filteredList.length === 0" description="暂无词库" />

        <div
          v-for="wb in filteredList"
          :key="wb.id"
          class="group rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/55 backdrop-blur px-4 py-4 shadow-sm hover:shadow-md transition flex items-center gap-3"
        >
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between gap-3">
              <div class="min-w-0 flex items-center gap-2">
                <span class="truncate text-base font-extrabold tracking-tight text-gray-900 dark:text-gray-50">{{ wb.title }}</span>
                <el-tag v-if="(wb as any).is_preset" size="small" type="warning" effect="plain" class="!rounded-full">系统</el-tag>
              </div>
              <el-tag size="small" type="info" class="!rounded-full">{{ wb.subject }}</el-tag>
            </div>
            <div class="mt-1 text-xs text-gray-500 dark:text-gray-400 truncate">
              {{ wb.education_stage }} · {{ wb.grade }} · {{ wb.version }}
            </div>
            <div class="mt-1 text-sm text-gray-600 dark:text-gray-300 truncate">{{ extractPreview(wb.content) }}</div>
          </div>

          <div v-if="!(wb as any).is_preset" class="flex-shrink-0 flex items-center gap-2">
            <el-button circle size="small" class="!rounded-2xl" :icon="Edit" @click.stop="router.push(`/dictation/banks/${wb.id}`)" />
            <el-popconfirm title="确认删除?" @confirm="deleteBank(wb.id)">
              <template #reference>
                <el-button circle size="small" class="!rounded-2xl" type="danger" :icon="Delete" @click.stop />
              </template>
            </el-popconfirm>
          </div>
          <div v-else class="flex-shrink-0">
            <el-button circle size="small" class="!rounded-2xl" :icon="Edit" disabled />
          </div>
        </div>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Edit, Delete, Collection } from '@element-plus/icons-vue'
import { dictationApi, type WordBank } from '@/services/dictation'
import { ElMessage } from 'element-plus'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

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
