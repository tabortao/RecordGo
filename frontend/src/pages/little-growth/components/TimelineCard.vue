<template>
  <div 
    class="rounded-2xl p-5 shadow-[0_4px_20px_-4px_rgba(0,0,0,0.05)] hover:shadow-[0_8px_30px_-4px_rgba(0,0,0,0.1)] transition-shadow duration-300 mb-6 dark:bg-gray-800"
    :style="cardBgStyle"
  >
    <!-- Header -->
    <div class="flex justify-between items-start mb-4">
      <div class="flex flex-col">
        <div class="text-sm font-bold inline-block" :class="textColorClass">
          {{ formatDateBadge(record.date) }} {{ formatTime(record.created_at) }}
        </div>
      </div>
      
      <el-dropdown trigger="click" @command="handleCommand">
        <div class="p-2 hover:bg-black/5 rounded-full cursor-pointer transition-colors">
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
      <p v-if="!searchQuery" class="text-gray-700 dark:text-gray-200 whitespace-pre-wrap leading-relaxed text-base">{{ record.content }}</p>
      <p v-else class="text-gray-700 dark:text-gray-200 whitespace-pre-wrap leading-relaxed text-base" v-html="highlightedContent"></p>
    </div>

    <!-- Gallery -->
    <div v-if="record.images && record.images.length > 0" class="mb-4">
      <div class="grid gap-2" :class="gridClass">
        <div 
          v-for="(img, index) in record.images" 
          :key="index"
          class="relative overflow-hidden rounded-xl group"
          :class="imgClass()"
        >
          <el-image 
            :src="img" 
            :preview-src-list="record.images"
            :initial-index="index"
            fit="contain" 
            class="w-full h-full object-contain transition-transform duration-500 group-hover:scale-105"
            loading="lazy"
            hide-on-click-modal
            preview-teleported
            :style="imageBgStyle"
          />
        </div>
      </div>
    </div>

    <!-- Footer: Tags -->
    <div v-if="displayTags.length > 0" class="flex flex-wrap gap-2">
      <span 
        v-for="tag in displayTags" 
        :key="tag.id"
        class="px-3 py-1 rounded-full text-xs font-medium transition-colors cursor-pointer brightness-110 hover:brightness-125"
        :style="{ backgroundColor: tag.color || '#E0E7FF', color: '#4F46E5' }"
        @click.stop="$emit('filter-tag', tag.id)"
      >
        {{ tag.name }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { MoreFilled, Edit, Delete } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'
import { type GrowthRecord, type Tag } from '@/stores/littleGrowth'

dayjs.locale('zh-cn')

const props = defineProps<{
  record: GrowthRecord
  allTags: Tag[]
  searchQuery?: string
}>()

const emit = defineEmits(['edit', 'delete', 'filter-tag'])

// Helpers
const formatDateBadge = (dateStr: string) => {
  return dayjs(dateStr).format('YYYY/MM/DD ddd')
}

const formatTime = (createdAt?: string) => {
  if (!createdAt) return ''
  return dayjs(createdAt).format('HH:mm')
}

const primaryTag = computed(() => {
  if (!props.record.tags || props.record.tags.length === 0) return null
  const tid = props.record.tags[0]
  return props.allTags.find(t => t.id === tid)
})

const cardBgStyle = computed(() => {
   const color = primaryTag.value?.color
   if (!color) return { backgroundColor: '#ffffff' } // Default white
   
   const rgba = hexToRgba(color, 0.1) // 10% opacity
   return {
     backgroundColor: rgba
   }
 })

 const imageBgStyle = computed(() => {
   const color = primaryTag.value?.color
   if (!color) return { backgroundColor: '#ffffff' } 
   // Match card background
   const rgba = hexToRgba(color, 0.1)
   return {
     backgroundColor: rgba
   }
 })

const textColorClass = computed(() => {
  return 'text-gray-600 dark:text-gray-300'
})

// Since we need to set background with opacity using the hex color, we need a helper or inline style hack
// Simplest: use rgba if we can convert, or just use opacity on a background layer.
// But we are in a single div.
// Let's try to just use the hex color with a very low opacity if possible, but hex doesn't support opacity easily without conversion.
// Let's use a utility to convert hex to rgba or use `mix-blend-mode`? No.
// Let's just apply the color to the background and set CSS variable for opacity? No.
// Let's assume the color is the background and we add an overlay?
// Or convert hex to rgba in JS.

function hexToRgba(hex: string, alpha: number) {
    let c: any;
    if(/^#([A-Fa-f0-9]{3}){1,2}$/.test(hex)){
        c= hex.substring(1).split('');
        if(c.length== 3){
            c= [c[0], c[0], c[1], c[1], c[2], c[2]];
        }
        c= '0x'+c.join('');
        return 'rgba('+[(c>>16)&255, (c>>8)&255, c&255].join(',')+','+alpha+')';
    }
    return hex; // fallback
}


const highlightedContent = computed(() => {
  if (!props.searchQuery) return props.record.content
  const safeContent = props.record.content.replace(/</g, '&lt;').replace(/>/g, '&gt;')
  const reg = new RegExp(`(${props.searchQuery.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi')
  return safeContent.replace(reg, '<mark class="bg-yellow-200 text-gray-900 rounded-sm px-0.5">$1</mark>')
})

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
