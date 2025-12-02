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

      <div v-if="store.activeFilterTagId" class="flex items-center gap-2 bg-purple-100 px-3 py-1 rounded-full">
        <span class="text-xs text-purple-600 font-medium">正在筛选: {{ activeTagName }}</span>
        <el-icon class="text-purple-400 cursor-pointer hover:text-purple-600" @click="store.activeFilterTagId = null"><Close /></el-icon>
      </div>
    </div>

    <!-- Main Content (Timeline) -->
    <div class="flex-1 overflow-y-auto p-4 pb-24" ref="scrollContainer">
      <div class="max-w-2xl mx-auto w-full">
        <template v-if="store.filteredRecords.length > 0">
          <TimelineCard 
            v-for="record in store.filteredRecords" 
            :key="record.id" 
            :record="record" 
            :allTags="store.flattenedTags"
            @edit="handleEdit"
            @delete="handleDelete"
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
import { ArrowLeft, Menu, Plus, Close, Files } from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { useLittleGrowthStore } from '@/stores/littleGrowth'
import { useWindowSize } from '@vueuse/core'
import TimelineCard from './components/TimelineCard.vue'
import TagSidebar from './components/TagSidebar.vue'

const router = useRouter()
const store = useLittleGrowthStore()
const showSidebar = ref(false)

const { width } = useWindowSize()
const drawerSize = computed(() => width.value < 768 ? '50%' : '25%')

onMounted(() => {
  store.fetchRecords()
})

const activeTagName = computed(() => {
  const tag = store.flattenedTags.find(t => t.id === store.activeFilterTagId)
  return tag ? tag.name : ''
})

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
        roundButton: true,
        customClass: 'little-growth-delete-dialog'
      }
    )
    await store.deleteRecord(id)
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
