<template>
  <div class="min-h-screen bg-[#F5F7FA] flex flex-col relative overflow-hidden">
    <!-- Header -->
    <div class="bg-white/80 backdrop-blur-md sticky top-0 z-20 px-4 py-3 shadow-sm flex items-center justify-between">
      <div class="flex items-center gap-3">
        <div 
          class="w-8 h-8 rounded-full bg-gray-100 flex items-center justify-center cursor-pointer hover:bg-gray-200 transition-colors text-gray-600"
          @click="router.push('/homework')"
        >
          <el-icon><ArrowLeft /></el-icon>
        </div>
        <div 
          class="w-8 h-8 rounded-full bg-gray-100 flex items-center justify-center cursor-pointer hover:bg-gray-200 transition-colors text-gray-600"
          @click="showSidebar = true"
        >
          <el-icon><Menu /></el-icon>
        </div>
        <h1 class="font-bold text-lg text-gray-800 tracking-tight">小成长</h1>
      </div>

      <div class="flex items-center gap-2">
        <div v-if="store.activeFilterTagId" class="flex items-center gap-2 bg-purple-100 px-3 py-1 rounded-full mr-2">
          <span class="text-xs text-purple-600 font-medium">正在筛选: {{ activeTagName }}</span>
          <el-icon class="text-purple-400 cursor-pointer hover:text-purple-600" @click="store.activeFilterTagId = null"><Close /></el-icon>
        </div>

        <!-- Search -->
        <el-input 
          v-model="searchQuery" 
          placeholder="搜索..." 
          class="!w-32 sm:!w-40 transition-all focus:!w-48" 
          size="small"
          clearable
          :prefix-icon="Search"
        />

        <!-- Calendar -->
        <div class="relative">
          <el-date-picker
            ref="datePickerRef"
            v-model="selectedDate"
            type="date"
            class="!w-0 !h-0 !border-0 !p-0 !overflow-hidden opacity-0 absolute top-0 left-0 -z-10"
            :popper-options="{ placement: 'bottom-end' }"
          >
            <template #cell="{ text, date, isCurrent }">
              <div class="w-full h-full flex flex-col items-center justify-center relative">
                <span :class="isCurrent ? 'font-bold text-purple-600' : ''">{{ text }}</span>
                <span v-if="hasRecord(date)" class="w-1.5 h-1.5 bg-purple-500 rounded-full absolute bottom-1"></span>
              </div>
            </template>
          </el-date-picker>
          
          <div 
            class="w-8 h-8 rounded-full flex items-center justify-center cursor-pointer transition-colors"
            :class="selectedDate ? 'bg-purple-100 text-purple-600' : 'bg-gray-100 text-gray-600 hover:bg-purple-50 hover:text-purple-600'"
            @click="openCalendar"
          >
             <el-icon v-if="selectedDate" @click.stop="selectedDate = null"><Close /></el-icon>
             <el-icon v-else><Calendar /></el-icon>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content (Timeline) -->
    <div class="flex-1 overflow-y-auto p-4 pb-24" ref="scrollContainer">
      <div class="max-w-2xl mx-auto w-full">
        <template v-if="filteredList.length > 0">
          <TimelineCard 
            v-for="record in filteredList" 
            :key="record.id" 
            :record="record" 
            :allTags="store.flattenedTags"
            :searchQuery="searchQuery"
            @edit="handleEdit"
            @delete="handleDelete"
            @filter-tag="handleTagSelect"
          />
        </template>
        <div v-else class="flex flex-col items-center justify-center py-20 text-gray-400">
          <el-icon :size="48" class="mb-4 opacity-50"><Files /></el-icon>
          <p>暂无记录，快去添加第一条美好回忆吧~</p>
        </div>
      </div>
    </div>

    <!-- FAB (Add Button) -->
    <div 
      class="fixed right-6 bottom-8 z-20 w-14 h-14 bg-gradient-to-tr from-purple-400 to-blue-400 rounded-full shadow-lg shadow-purple-200 flex items-center justify-center text-white cursor-pointer hover:scale-110 transition-transform active:scale-95"
      @click="router.push('/little-growth/create')"
    >
      <el-icon :size="24"><Plus /></el-icon>
    </div>

    <!-- Sidebar Drawer -->
    <el-drawer
      v-model="showSidebar"
      direction="ltr"
      :with-header="false"
      :size="drawerSize"
      class="little-growth-drawer"
    >
      <TagSidebar 
        :tags="store.tags" 
        :active-tag-id="store.activeFilterTagId"
        :total-records="store.records.length"
        @select="handleTagSelect"
      />
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Menu, Plus, Close, Files, Search, Calendar } from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { useLittleGrowthStore } from '@/stores/littleGrowth'
import { useWindowSize } from '@vueuse/core'
import TimelineCard from './components/TimelineCard.vue'
import TagSidebar from './components/TagSidebar.vue'
import dayjs from 'dayjs'

const router = useRouter()
const store = useLittleGrowthStore()
const showSidebar = ref(false)

const searchQuery = ref('')
const selectedDate = ref<Date | null>(null)
const datePickerRef = ref()

const { width } = useWindowSize()
const drawerSize = computed(() => width.value < 768 ? '50%' : '25%')

onMounted(() => {
  store.fetchRecords()
})

const activeTagName = computed(() => {
  const tag = store.flattenedTags.find(t => t.id === store.activeFilterTagId)
  return tag ? tag.name : ''
})

const filteredList = computed(() => {
  let list = store.records

  // 1. Tag Filter
  if (store.activeFilterTagId) {
    list = list.filter(r => {
        let rTags: string[] = []
        if (Array.isArray(r.tags)) {
            rTags = r.tags
        } else if (typeof r.tags === 'string') {
            try { rTags = JSON.parse(r.tags) } catch {}
        }
        return rTags.map(String).includes(String(store.activeFilterTagId))
    })
  }

  // 2. Search Filter
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(r => r.content.toLowerCase().includes(q))
  }

  // 3. Date Filter
  if (selectedDate.value) {
    const d = dayjs(selectedDate.value).format('YYYY-MM-DD')
    list = list.filter(r => r.date === d)
  }
  
  return list
})

const openCalendar = () => {
    if (selectedDate.value) {
        selectedDate.value = null
    } else {
        datePickerRef.value?.focus()
    }
}

const hasRecord = (date: Date) => {
  const d = dayjs(date).format('YYYY-MM-DD')
  return store.records.some(r => r.date === d)
}

const handleTagSelect = (tagId: string | null) => {
  store.activeFilterTagId = tagId
  showSidebar.value = false
}

const handleEdit = (id: string) => {
  router.push(`/little-growth/edit/${id}`)
}

const handleDelete = async (id: string) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除这条记录吗？删除后无法恢复。',
      '删除确认',
      {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning',
        center: true,
        draggable: true,
        lockScroll: false,
      }
    )
    await store.deleteRecord(id)
    await store.fetchRecords() // Refresh list from server
    ElMessage.success('删除成功')
  } catch {
    // cancelled
  }
}
</script>

<style>
.little-growth-drawer .el-drawer__body {
  padding: 0;
}
</style>
