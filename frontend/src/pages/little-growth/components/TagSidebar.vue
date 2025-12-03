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
        <!-- My Favorites (placed above 'All') -->
        <div 
          class="px-3 py-2 rounded-lg cursor-pointer transition-colors flex justify-between items-center"
          :class="showFavorites ? 'bg-yellow-50 dark:bg-yellow-900/30 text-yellow-600 dark:text-yellow-300 font-medium' : 'text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800'"
          @click="$emit('select-favorites')"
        >
          <span>我的收藏</span>
        </div>

        <!-- All Records -->
        <div 
          class="px-3 py-2 rounded-lg cursor-pointer transition-colors flex justify-between items-center"
          :class="!activeTagId && !showFavorites ? 'bg-purple-50 dark:bg-purple-900/30 text-purple-600 dark:text-purple-300 font-medium' : 'text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800'"
          @click="$emit('select', null)"
        >
          <span>全部记录</span>
          <span class="text-xs opacity-60">{{ totalRecords }}</span>
        </div>

        <div v-for="p in parents" :key="p.id" class="space-y-1 mt-2">
          <div 
            class="px-3 py-2 rounded-lg cursor-pointer transition-colors flex justify-between items-center group"
            :class="activeTagId === p.id ? 'font-medium bg-blue-50 dark:bg-blue-900/30' : 'hover:bg-gray-50 dark:hover:bg-gray-800'"
            :style="activeTagId === p.id ? {} : { backgroundColor: getBgColor(p.color) }"
          >
            <div class="flex items-center gap-2" :class="activeTagId === p.id ? 'text-blue-600 dark:text-blue-400' : 'text-gray-700 dark:text-gray-300'" @click="$emit('select', p.id)">
              <span>{{ p.name }}</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="text-xs bg-gray-100 dark:bg-gray-800 px-2 py-0.5 rounded-full text-gray-500 dark:text-gray-400 group-hover:bg-white dark:group-hover:bg-gray-700 transition-colors">{{ p.count }}</span>
              <el-icon v-if="childrenMap[p.id] && childrenMap[p.id].length > 0" class="text-gray-400 hover:text-gray-600 cursor-pointer" @click.stop="toggleExpand(p.id)">
                <CaretRight :class="isExpanded(p.id) ? 'rotate-90 transition-transform' : 'transition-transform'" />
              </el-icon>
            </div>
          </div>
          <div v-if="isExpanded(p.id) && childrenMap[p.id] && childrenMap[p.id].length > 0" class="pl-5 space-y-1">
            <div 
              v-for="c in childrenMap[p.id]" :key="c.id"
              class="px-3 py-2 rounded-lg cursor-pointer transition-colors flex justify-between items-center"
              :class="activeTagId === c.id ? 'font-medium bg-blue-50 dark:bg-blue-900/30' : 'hover:bg-gray-50 dark:hover:bg-gray-800'"
              :style="activeTagId === c.id ? {} : { backgroundColor: getBgColor(c.color) }"
              @click="$emit('select', c.id)"
            >
              <div class="flex items-center gap-2" :class="activeTagId === c.id ? 'text-blue-600 dark:text-blue-400' : 'text-gray-700 dark:text-gray-300'">
                <span>{{ c.name }}</span>
              </div>
              <span class="text-xs bg-gray-100 dark:bg-gray-800 px-2 py-0.5 rounded-full text-gray-500 dark:text-gray-400">{{ c.count }}</span>
            </div>
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
import { Plus, CaretRight } from '@element-plus/icons-vue'
import type { Tag } from '@/stores/littleGrowth'
import { useAuth } from '@/stores/auth'
import { presignView } from '@/services/storage'

const props = defineProps<{
  tags: Tag[]
  activeTagId: string | null
  totalRecords: number
  showFavorites?: boolean
}>()

defineEmits(['select', 'select-favorites'])

const auth = useAuth()
const user = computed(() => auth.user as any)
const nickname = computed(() => user.value?.nickname || '未登录')
const userInitial = computed(() => nickname.value?.[0]?.toUpperCase() || 'M')

const avatarSrc = ref('')

// Helper to convert hex to rgba for 50% opacity
const getBgColor = (color?: string) => {
    if (!color) return '#F3E8FF'
    // Convert to rgba with 0.5 opacity
    let c: any;
    if (/^#([A-Fa-f0-9]{3}){1,2}$/.test(color)) {
        c = color.substring(1).split('');
        if (c.length == 3) {
            c = [c[0], c[0], c[1], c[1], c[2], c[2]];
        }
        c = '0x' + c.join('');
        return 'rgba(' + [(c >> 16) & 255, (c >> 8) & 255, c & 255].join(',') + ',0.5)';
    }
    return color
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

const parents = computed(() => props.tags.filter(t => !t.parentId))

const childrenMap = computed<Record<string, Tag[]>>(() => {
  const m: Record<string, Tag[]> = {}
  for (const t of props.tags) {
    if (t.parentId) {
      const pid = String(t.parentId)
      ;(m[pid] = m[pid] || []).push(t)
    }
  }
  return m
})

const expanded = ref<Record<string, boolean>>({})
const isExpanded = (id: string) => !!expanded.value[String(id)]
const toggleExpand = (id: string) => { expanded.value[String(id)] = !expanded.value[String(id)] }
</script>
