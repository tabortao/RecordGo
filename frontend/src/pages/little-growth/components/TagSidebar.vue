<template>
  <div class="h-full flex flex-col bg-white dark:bg-gray-900 border-r border-gray-100 dark:border-gray-800">
    <!-- Header: User Info -->
    <div class="p-6 border-b border-gray-50 dark:border-gray-800">
      <div class="flex items-center gap-4">
        <div class="w-12 h-12 rounded-full bg-gray-200 dark:bg-gray-700 overflow-hidden flex items-center justify-center text-gray-500 dark:text-gray-400 font-bold text-lg shadow-sm border border-white dark:border-gray-600">
          <img v-if="avatarSrc" :src="avatarSrc" class="w-full h-full object-cover" />
          <span v-else>{{ userInitial }}</span>
        </div>
        <div>
          <h3 class="font-bold text-gray-800 dark:text-white text-lg">{{ nickname }}</h3>
          <p class="text-xs text-gray-400">{{ totalRecords }} 条美好回忆</p>
        </div>
      </div>
    </div>

    <!-- Tag Tree -->
    <div class="flex-1 overflow-y-auto p-4">
      <div class="space-y-1">
        <div 
          class="px-3 py-2 rounded-lg cursor-pointer transition-colors flex justify-between items-center"
          :class="!activeTagId ? 'bg-purple-50 dark:bg-purple-900/30 text-purple-600 dark:text-purple-300 font-medium' : 'text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800'"
          @click="$emit('select', null)"
        >
          <span>全部记录</span>
          <span class="text-xs opacity-60">{{ totalRecords }}</span>
        </div>

        <div v-for="tag in tags" :key="tag.id" class="space-y-1 mt-2">
          <!-- Level 1 -->
          <div 
            class="px-3 py-2 rounded-lg cursor-pointer transition-colors flex justify-between items-center group"
            :class="activeTagId === tag.id ? 'font-medium' : 'hover:bg-gray-50 dark:hover:bg-gray-800'"
            :style="activeTagId === tag.id ? { backgroundColor: getBgColor(tag.color), color: getTextColor(tag.color) } : {}"
            @click="$emit('select', tag.id)"
          >
            <div class="flex items-center gap-2" :class="activeTagId !== tag.id ? 'text-gray-700 dark:text-gray-300' : ''">
              <el-icon v-if="tag.children && tag.children.length" class="text-gray-400 text-xs"><ArrowDown /></el-icon>
              <span>{{ tag.name }}</span>
            </div>
            <span class="text-xs bg-gray-100 dark:bg-gray-800 px-2 py-0.5 rounded-full text-gray-500 dark:text-gray-400 group-hover:bg-white dark:group-hover:bg-gray-700 transition-colors">{{ tag.count }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Footer: Manage -->
    <div class="p-4 border-t border-gray-50 dark:border-gray-800">
      <button 
        class="w-full py-2 px-4 rounded-xl border border-dashed border-gray-300 dark:border-gray-700 text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800 hover:text-purple-600 dark:hover:text-purple-400 transition-colors text-sm flex items-center justify-center gap-2"
        @click="$router.push('/little-growth/tags')"
      >
        <el-icon><Plus /></el-icon>
        <span>管理标签</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watchEffect } from 'vue'
import { ArrowDown, Plus } from '@element-plus/icons-vue'
import type { Tag } from '@/stores/littleGrowth'
import { useAuth } from '@/stores/auth'
import { presignView } from '@/services/storage'

defineProps<{
  tags: Tag[]
  activeTagId: string | null
  totalRecords: number
}>()

defineEmits(['select'])

const auth = useAuth()
const user = computed(() => auth.user as any)
const nickname = computed(() => user.value?.nickname || '未登录')
const userInitial = computed(() => nickname.value?.[0]?.toUpperCase() || 'M')

const avatarSrc = ref('')

const getBgColor = (color?: string) => {
  // Just use the color directly as background (user asked for consistent with settings)
  // But for text readability, maybe we need to adjust?
  // User said "底纹与标签设置的颜色一致"
  return color || '#F3E8FF' // default purple-50 like
}

const getTextColor = (color?: string) => {
    // If color is dark, text should be white. If light, black/gray.
    // Simple heuristic or just default to a dark gray for pastel colors
    return '#374151' 
}

watchEffect(async () => {
    const p = user.value?.avatar_path
    if (!p) { 
        avatarSrc.value = '' 
        return 
    }
    if (p.endsWith('default.png') || p.includes('/default.png')) {
        avatarSrc.value = ''
        return
    }
    if (p.startsWith('http')) {
        avatarSrc.value = p
        return
    }
    if (p.includes('uploads/')) {
        avatarSrc.value = `${import.meta.env.VITE_API_BASE}/api/${p.replace(/^\/+/, '')}`
        return
    }
    try {
        avatarSrc.value = await presignView(p)
    } catch {
        avatarSrc.value = ''
    }
})
</script>
