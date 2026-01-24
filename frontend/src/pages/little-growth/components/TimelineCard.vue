<template>
  <div 
    class="relative bg-white rounded-3xl p-5 shadow-sm border border-gray-100 transition-all duration-300 hover:shadow-lg hover:-translate-y-0.5 dark:bg-gray-800 dark:border-gray-700 group overflow-hidden"
  >
    <!-- Decoration -->
    <div class="absolute top-0 right-0 w-24 h-24 bg-gradient-to-br from-purple-50 to-blue-50 dark:from-purple-900/10 dark:to-blue-900/10 rounded-bl-full -z-0 opacity-50"></div>

    <!-- Header -->
    <div class="flex justify-between items-start mb-3 relative z-10">
      <div class="flex items-center gap-3">
         <div class="relative">
           <el-avatar :size="44" :src="avatarUrl" class="flex-shrink-0 border-2 border-white shadow-sm dark:border-gray-700" />
           <div v-if="record.is_pinned" class="absolute -top-1 -right-1 bg-amber-400 text-white rounded-full p-0.5 shadow-sm">
             <el-icon :size="10"><Top /></el-icon>
           </div>
         </div>
         <div class="flex flex-col">
            <span class="text-base font-bold text-gray-800 dark:text-gray-100 leading-tight">{{ nickname }}</span>
            <div class="text-xs text-gray-400 dark:text-gray-500 inline-flex items-center gap-1 mt-0.5 font-medium">
               <span>{{ formatDateBadge(record.date) }}</span>
               <span class="w-0.5 h-0.5 rounded-full bg-gray-300 dark:bg-gray-600"></span>
               <span>{{ formatTime(record.date) }}</span>
            </div>
         </div>
      </div>
      
      <div class="flex items-center gap-1">
        <!-- Favorite Icon -->
        <div 
          class="inline-flex items-center justify-center w-8 h-8 rounded-full cursor-pointer transition-all active:scale-90" 
          :class="record.is_favorite ? 'bg-yellow-50 text-yellow-400 dark:bg-yellow-900/20' : 'text-gray-300 hover:bg-gray-100 hover:text-gray-500 dark:hover:bg-gray-700'"
          @click.stop="$emit('toggle-favorite', record.id)"
        >
          <el-icon :size="18">
            <StarFilled v-if="record.is_favorite" />
            <Star v-else />
          </el-icon>
        </div>
        
        <!-- More Actions -->
        <el-dropdown trigger="click" @command="handleCommand">
          <div class="inline-flex items-center justify-center w-8 h-8 rounded-full cursor-pointer hover:bg-gray-100 text-gray-400 transition-colors dark:hover:bg-gray-700">
            <el-icon :size="18"><MoreFilled /></el-icon>
          </div>
          <template #dropdown>
            <el-dropdown-menu class="rounded-xl overflow-hidden">
              <el-dropdown-item command="pin" :icon="Top">{{ record.is_pinned ? '取消置顶' : '置顶' }}</el-dropdown-item>
              <el-dropdown-item command="copy" :icon="CopyDocument">复制内容</el-dropdown-item>
              <el-dropdown-item command="edit" :icon="Edit">编辑记录</el-dropdown-item>
              <el-dropdown-item command="delete" :icon="Delete" class="text-red-500">删除记录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- Content -->
    <div class="mb-3 relative z-10 pl-[56px]">
      <p v-if="!searchQuery" class="text-gray-700 dark:text-gray-200 whitespace-pre-wrap leading-relaxed text-[15px]">{{ record.content }}</p>
      <p v-else class="text-gray-700 dark:text-gray-200 whitespace-pre-wrap leading-relaxed text-[15px]" v-html="highlightedContent"></p>
    </div>

    <!-- Gallery -->
    <div v-if="record.images && record.images.length > 0" class="mb-3 pl-[56px] relative z-10">
      <div class="grid gap-2" :class="gridClass">
        <div 
          v-for="(img, index) in record.images" 
          :key="index"
          class="relative overflow-hidden rounded-2xl group/img bg-gray-100 dark:bg-gray-900"
          :class="imgClass()"
        >
          <el-image 
            :src="img" 
            :preview-src-list="record.images"
            :initial-index="index"
            fit="cover" 
            class="w-full h-full object-cover transition-transform duration-700 group-hover/img:scale-110"
            loading="lazy"
            hide-on-click-modal
            preview-teleported
          />
        </div>
      </div>
    </div>

    <!-- Footer: Tags & Comments Trigger -->
    <div class="flex items-center justify-between pl-[56px] mt-2 relative z-10">
        <div class="flex flex-wrap gap-2">
          <span 
            v-for="tag in displayTags" 
            :key="tag.id"
            class="px-2.5 py-1 rounded-lg text-xs font-bold transition-all cursor-pointer hover:opacity-80 active:scale-95"
            :style="getTagStyle(tag.color)"
            @click.stop="$emit('filter-tag', tag.id)"
          >
            # {{ tag.name }}
          </span>
        </div>
        
        <div 
            class="flex items-center gap-1.5 text-gray-400 hover:text-purple-600 cursor-pointer transition-colors text-xs font-medium px-2 py-1 rounded-full hover:bg-purple-50 dark:hover:bg-purple-900/20"
            @click.stop="toggleCommentBox"
        >
            <el-icon :size="14"><ChatDotSquare /></el-icon>
            <span>{{ record.comments?.length ? record.comments.length : '评论' }}</span>
        </div>
    </div>

    <!-- Comments Section -->
    <div v-if="(record.comments && record.comments.length > 0) || showInput" class="mt-4 ml-[56px] bg-gray-50/80 dark:bg-gray-900/50 rounded-2xl p-4 relative z-10 border border-gray-100 dark:border-gray-700/50">
        <!-- Triangle indicator (optional, omitted for simplicity) -->
        
        <!-- Comments List -->
        <div v-if="record.comments && record.comments.length > 0" class="space-y-1 mb-2">
             <div v-for="c in record.comments" :key="c.id" class="text-sm leading-6">
                <span class="font-semibold text-blue-600 cursor-pointer">{{ displayCommentName(c) }}</span>
                <span class="mx-1 text-gray-400">:</span>
                <span class="text-gray-700 dark:text-gray-300">{{ c.content }}</span>
             </div>
        </div>

        <!-- Input Box -->
        <div v-if="showInput" class="flex gap-2 animate-fade-in items-end">
            <div class="flex-1 relative">
                <el-input 
                    v-model="commentText" 
                    type="textarea"
                    :rows="1"
                    autosize
                    placeholder="评论..." 
                    @keyup.enter.ctrl="submitComment" 
                    ref="inputRef"
                    class="w-full"
                />
                <el-popover placement="top" :width="250" trigger="click">
                    <template #reference>
                        <div class="absolute right-2 bottom-2 cursor-pointer text-gray-400 hover:text-yellow-500 transition-colors z-10">
                            <el-icon :size="20"><img src="https://api.iconify.design/twemoji:grinning-face.svg" class="w-5 h-5 opacity-60 hover:opacity-100" /></el-icon>
                        </div>
                    </template>
                    <div class="flex flex-wrap gap-2 max-h-40 overflow-y-auto">
                        <span 
                            v-for="emoji in commonEmojis" 
                            :key="emoji" 
                            class="cursor-pointer text-xl hover:bg-gray-100 p-1 rounded transition-colors select-none"
                            @click="appendEmoji(emoji)"
                        >
                            {{ emoji }}
                        </span>
                    </div>
                </el-popover>
            </div>
            <el-button type="primary" :disabled="!commentText.trim()" @click="submitComment">发送</el-button>
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, nextTick, watchEffect } from 'vue'
import { MoreFilled, Edit, Delete, Top, CopyDocument, ChatDotSquare, Star, StarFilled } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'
import { type GrowthRecord, type Tag } from '@/stores/littleGrowth'
import { useClipboard } from '@vueuse/core'
import { ElMessage } from 'element-plus'
import { useAuth } from '@/stores/auth'
import defaultAvatar from '@/assets/avatars/default.png'
import { getStaticBase } from '@/services/http'
import { presignView } from '@/services/storage'

dayjs.locale('zh-cn')

const commonEmojis = [
    '😀', '😂', '🤣', '😍', '🥰', '😘', '🤪', '🤩', '🥳', '😎',
    '😭', '😤', '😡', '🤯', '😱', '👍', '👎', '👏', '🤝', '🙏',
    '❤️', '💔', '🎉', '🌹', '✨', '🔥', '💯', '💪', '👀', '🐶'
]

const props = defineProps<{
  record: GrowthRecord
  allTags: Tag[]
  searchQuery?: string
}>()

const nickname = computed(() => {
    return props.record.user?.nickname || props.record.user?.username || '未知用户'
})

const avatarUrl = ref(defaultAvatar)

watchEffect(async () => {
    const p = props.record.user?.avatar_path
    if (!p) {
        avatarUrl.value = defaultAvatar
        return
    }
    const s = String(p)
    if (/storage\/images\/avatars\/default\.png$/i.test(s) || /(^|\/)default\.png$/i.test(s)) {
         avatarUrl.value = defaultAvatar
         return
    }
    if (/^https?:\/\//i.test(s)) {
        avatarUrl.value = s
        return
    }
    const base = getStaticBase()
    if (/uploads\//i.test(s)) {
        avatarUrl.value = `${base}/api/${s.replace(/^\/+/, '')}`
        return
    }
    try {
        avatarUrl.value = await presignView(s)
    } catch {
        avatarUrl.value = defaultAvatar
    }
})

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

function appendEmoji(emoji: string) {
    commentText.value += emoji
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
const auth = useAuth()
const myNickname = computed(() => {
  const u = auth.user
  return (u?.nickname || '').trim() || (u?.username || '用户')
})

function displayCommentName(c: import('@/stores/littleGrowth').GrowthComment) {
  const n1 = (c.user?.nickname || '').trim()
  if (n1) return n1
  const uid = String(c.user_id || '')
  const me = String(auth.user?.id || '')
  if (uid && me && uid === me) return myNickname.value
  const n2 = (c.user?.username || '').trim()
  return n2 || '用户'
}
</script>
