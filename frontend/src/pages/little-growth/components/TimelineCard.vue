<template>
  <div 
    class="bg-white rounded-2xl p-6 shadow-[0_4px_20px_-4px_rgba(0,0,0,0.05)] hover:shadow-[0_8px_30px_-4px_rgba(0,0,0,0.1)] transition-shadow duration-300 mb-4 dark:bg-gray-800 group"
  >
    <!-- Header -->
    <div class="flex justify-between items-start mb-2">
      <div class="flex flex-col">
        <div class="text-sm font-bold inline-block flex items-center gap-2" :class="textColorClass">
          <span>{{ formatDateBadge(record.date) }} {{ formatTime(record.date) }}</span>
          <el-icon v-if="record.is_pinned" class="text-purple-500"><Top /></el-icon>
        </div>
      </div>
      
      <div class="flex items-center gap-1">
        <!-- Favorite Icon -->
        <div class="p-2 hover:bg-black/5 rounded-full cursor-pointer transition-colors" @click.stop="$emit('toggle-favorite', record.id)">
          <el-icon :size="20" :class="record.is_favorite ? 'text-yellow-400' : 'text-gray-400'"><StarFilled v-if="record.is_favorite" /><Star v-else /></el-icon>
        </div>

        <!-- Comment Icon -->
        <div class="p-2 hover:bg-black/5 rounded-full cursor-pointer transition-colors" @click.stop="toggleCommentBox">
          <el-icon :size="20" class="text-gray-400"><ChatDotSquare /></el-icon>
        </div>

        <el-dropdown trigger="click" @command="handleCommand">
          <div class="p-2 hover:bg-black/5 rounded-full cursor-pointer transition-colors">
            <el-icon :size="20" class="text-gray-400"><MoreFilled /></el-icon>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="pin" :icon="Top">{{ record.is_pinned ? '取消置顶' : '置顶' }}</el-dropdown-item>
              <el-dropdown-item command="copy" :icon="CopyDocument">复制</el-dropdown-item>
              <el-dropdown-item command="edit" :icon="Edit">编辑</el-dropdown-item>
              <el-dropdown-item command="delete" :icon="Delete" class="text-red-500">删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- Content -->
    <div class="mb-2">
      <p v-if="!searchQuery" class="text-gray-700 dark:text-gray-200 whitespace-pre-wrap leading-relaxed text-lg">{{ record.content }}</p>
      <p v-else class="text-gray-700 dark:text-gray-200 whitespace-pre-wrap leading-relaxed text-lg" v-html="highlightedContent"></p>
    </div>

    <!-- Gallery -->
    <div v-if="record.images && record.images.length > 0" class="mb-2">
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
          />
        </div>
      </div>
    </div>

    <!-- Footer: Tags -->
    <div v-if="displayTags.length > 0" class="flex flex-wrap gap-2 mb-2">
      <span 
        v-for="tag in displayTags" 
        :key="tag.id"
        class="px-3 py-1 rounded-full text-xs font-medium transition-colors cursor-pointer"
        :style="getTagStyle(tag.color)"
        @click.stop="$emit('filter-tag', tag.id)"
      >
        {{ tag.name }}
      </span>
    </div>

    <!-- Comments Section -->
    <div v-if="(record.comments && record.comments.length > 0) || showInput" class="mt-2 bg-gray-50 dark:bg-gray-900/50 rounded-lg p-3 relative">
        <!-- Triangle indicator (optional, omitted for simplicity) -->
        
        <!-- Comments List -->
        <div v-if="record.comments && record.comments.length > 0" class="space-y-1 mb-2">
             <div v-for="c in record.comments" :key="c.id" class="text-sm leading-6">
                <span class="font-semibold text-blue-600 cursor-pointer">{{ c.user?.nickname || '用户' }}</span>
                <span class="mx-1 text-gray-400">:</span>
                <span class="text-gray-700 dark:text-gray-300">{{ c.content }}</span>
             </div>
        </div>

        <!-- Input Box -->
        <div v-if="showInput" class="flex gap-2 animate-fade-in">
            <el-input 
                v-model="commentText" 
                size="default" 
                placeholder="评论..." 
                @keyup.enter="submitComment" 
                ref="inputRef"
            />
            <el-button type="primary" :disabled="!commentText.trim()" @click="submitComment">发送</el-button>
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, nextTick } from 'vue'
import { MoreFilled, Edit, Delete, Top, CopyDocument, ChatDotSquare, Star, StarFilled } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'
import { type GrowthRecord, type Tag } from '@/stores/littleGrowth'
import { useClipboard } from '@vueuse/core'
import { ElMessage } from 'element-plus'

dayjs.locale('zh-cn')

const props = defineProps<{
  record: GrowthRecord
  allTags: Tag[]
  searchQuery?: string
}>()

const emit = defineEmits(['edit', 'delete', 'filter-tag', 'pin', 'toggle-favorite', 'add-comment'])

const { copy } = useClipboard()

const handleCommand = (cmd: string) => {
    if (cmd === 'edit') emit('edit', props.record.id)
    if (cmd === 'delete') emit('delete', props.record.id)
    if (cmd === 'pin') emit('pin', props.record.id)
    if (cmd === 'copy') {
        copy(props.record.content)
        ElMessage.success('已复制到剪贴板')
    }
}

// Comment Logic
const showInput = ref(false)
const commentText = ref('')
const inputRef = ref()

function toggleCommentBox() {
    showInput.value = !showInput.value
    if (showInput.value) {
        nextTick(() => {
            inputRef.value?.focus()
        })
    }
}

function submitComment() {
    if (!commentText.value.trim()) return
    emit('add-comment', props.record.id, commentText.value)
    commentText.value = ''
    showInput.value = false
}

// Helpers
const formatDateBadge = (dateStr: string) => {
  return dayjs(dateStr).format('YYYY/MM/DD ddd')
}

const formatTime = (dateStr?: string) => {
  if (!dateStr) return ''
  // Check if dateStr includes time
  if (dateStr.includes(' ') || dateStr.includes('T')) {
      return dayjs(dateStr).format('HH:mm')
  }
  return ''
}

const textColorClass = computed(() => {
  return 'text-gray-600 dark:text-gray-300'
})


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

 const getTagStyle = (color?: string) => {
    if (!color) return { backgroundColor: 'rgba(224, 231, 255, 0.5)', color: '#4F46E5' } // default with 50% opacity
    const rgba = hexToRgba(color, 0.5) // 50% opacity
    return {
        backgroundColor: rgba,
        color: '#1F2937' // dark gray
    }
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
</script>
