<template>
  <div class="h-full flex flex-col bg-white border-r border-gray-100">
    <!-- Header: User Info -->
    <div class="p-6 border-b border-gray-50">
      <div class="flex items-center gap-4">
        <div class="w-12 h-12 rounded-full bg-gradient-to-br from-purple-200 to-blue-200 flex items-center justify-center text-white font-bold text-lg">
          {{ userInitial }}
        </div>
        <div>
          <h3 class="font-bold text-gray-800 text-lg">我的记录</h3>
          <p class="text-xs text-gray-400">{{ totalRecords }} 条美好回忆</p>
        </div>
      </div>
    </div>

    <!-- Tag Tree -->
    <div class="flex-1 overflow-y-auto p-4">
      <div class="space-y-1">
        <div 
          class="px-3 py-2 rounded-lg cursor-pointer transition-colors flex justify-between items-center"
          :class="!activeTagId ? 'bg-purple-50 text-purple-600 font-medium' : 'text-gray-600 hover:bg-gray-50'"
          @click="$emit('select', null)"
        >
          <span>全部记录</span>
          <span class="text-xs opacity-60">{{ totalRecords }}</span>
        </div>

        <div v-for="tag in tags" :key="tag.id" class="space-y-1 mt-2">
          <!-- Level 1 -->
          <div 
            class="px-3 py-2 rounded-lg cursor-pointer transition-colors flex justify-between items-center group"
            :class="activeTagId === tag.id ? 'bg-purple-50 text-purple-600 font-medium' : 'text-gray-700 hover:bg-gray-50'"
            @click="$emit('select', tag.id)"
          >
            <div class="flex items-center gap-2">
              <el-icon v-if="tag.children && tag.children.length" class="text-gray-400 text-xs"><ArrowDown /></el-icon>
              <span>{{ tag.name }}</span>
            </div>
            <span class="text-xs bg-gray-100 px-2 py-0.5 rounded-full text-gray-500 group-hover:bg-white transition-colors">{{ tag.count }}</span>
          </div>

          <!-- Level 2 -->
          <div v-if="tag.children && tag.children.length" class="ml-4 pl-4 border-l border-gray-100 space-y-1">
            <div 
              v-for="child in tag.children" 
              :key="child.id"
              class="px-3 py-1.5 rounded-lg cursor-pointer transition-colors flex justify-between items-center text-sm"
              :class="activeTagId === child.id ? 'bg-purple-50 text-purple-600 font-medium' : 'text-gray-500 hover:bg-gray-50'"
              @click="$emit('select', child.id)"
            >
              <span>{{ child.name }}</span>
              <span class="text-xs opacity-60">{{ child.count }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Footer: Manage -->
    <div class="p-4 border-t border-gray-50">
      <button class="w-full py-2 px-4 rounded-xl border border-dashed border-gray-300 text-gray-500 hover:bg-gray-50 hover:text-purple-600 transition-colors text-sm flex items-center justify-center gap-2">
        <el-icon><Plus /></el-icon>
        <span>管理标签</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ArrowDown, Plus } from '@element-plus/icons-vue'
import type { Tag } from '@/stores/littleGrowth'

const props = defineProps<{
  tags: Tag[]
  activeTagId: string | null
  totalRecords: number
}>()

defineEmits(['select'])

const userInitial = computed(() => 'M') // Mock
</script>
