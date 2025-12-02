<template>
  <div class="bg-white rounded-2xl p-5 shadow-[0_4px_20px_-4px_rgba(0,0,0,0.05)] hover:shadow-[0_8px_30px_-4px_rgba(0,0,0,0.1)] transition-shadow duration-300 mb-6">
    <!-- Header -->
    <div class="flex justify-between items-start mb-4">
      <div class="flex flex-col">
        <div class="bg-purple-50 text-purple-600 px-3 py-1 rounded-lg text-sm font-bold inline-block">
          {{ formatDateBadge(record.date) }}
        </div>
      </div>
      
      <el-dropdown trigger="click" @command="handleCommand">
        <div class="p-2 hover:bg-gray-50 rounded-full cursor-pointer transition-colors">
          <el-icon :size="20" class="text-gray-400"><MoreFilled /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="edit" :icon="Edit">编辑</el-dropdown-item>
            <el-dropdown-item command="delete" :icon="Delete" class="text-red-500">删除</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>

    <!-- Content -->
    <div class="mb-4">
      <p class="text-gray-700 whitespace-pre-wrap leading-relaxed text-base">{{ record.content }}</p>
    </div>

    <!-- Gallery -->
    <div v-if="record.images && record.images.length > 0" class="mb-4">
      <div class="grid gap-2" :class="gridClass">
        <div 
          v-for="(img, index) in record.images" 
          :key="index"
          class="relative overflow-hidden rounded-xl bg-gray-100 group"
          :class="imgClass()"
        >
          <el-image 
            :src="img" 
            :preview-src-list="record.images"
            :initial-index="index"
            fit="cover" 
            class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
            loading="lazy"
            hide-on-click-modal
          />
        </div>
      </div>
    </div>

    <!-- Footer: Tags -->
    <div v-if="displayTags.length > 0" class="flex flex-wrap gap-2">
      <span 
        v-for="tag in displayTags" 
        :key="tag.id"
        class="px-3 py-1 rounded-full bg-blue-50 text-blue-600 text-xs font-medium hover:bg-blue-100 transition-colors cursor-pointer"
      >
        #{{ tag.name }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { MoreFilled, Edit, Delete } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { type GrowthRecord, type Tag } from '@/stores/littleGrowth'

const props = defineProps<{
  record: GrowthRecord
  allTags: Tag[]
}>()

const emit = defineEmits(['edit', 'delete'])

// Helpers
const formatDateBadge = (dateStr: string) => {
  return dayjs(dateStr).format('YYYY MMM DD').toUpperCase()
}

const displayTags = computed(() => {
  return props.record.tags.map(tid => props.allTags.find(t => t.id === tid)).filter(Boolean) as Tag[]
})

// Grid Logic
const count = computed(() => props.record.images?.length || 0)

const gridClass = computed(() => {
  const n = count.value
  if (n === 1) return 'grid-cols-1'
  if (n === 2) return 'grid-cols-2'
  if (n === 4) return 'grid-cols-2'
  return 'grid-cols-3'
})

const imgClass = () => {
  const n = count.value
  if (n === 1) return 'aspect-video' // 16/9
  return 'aspect-square' // 1:1 for grid
}

const handleCommand = (cmd: string) => {
  if (cmd === 'edit') emit('edit', props.record.id)
  if (cmd === 'delete') emit('delete', props.record.id)
}
</script>
