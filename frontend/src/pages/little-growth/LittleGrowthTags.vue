<template>
  <div class="min-h-screen bg-white flex flex-col">
    <!-- Header -->
    <div class="px-4 py-3 border-b border-gray-100 flex items-center justify-between sticky top-0 bg-white z-20">
      <div class="flex items-center gap-2 cursor-pointer" @click="router.back()">
        <el-icon class="text-gray-600 text-lg"><ArrowLeft /></el-icon>
        <span class="text-gray-600">返回</span>
      </div>
      <h1 class="font-bold text-lg">管理标签</h1>
      <div class="w-10"></div> <!-- Spacer -->
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto p-4">
      <div class="max-w-md mx-auto">
        <div v-if="tags.length === 0" class="text-center text-gray-400 py-10">
          暂无标签
        </div>
        <div v-else class="space-y-3">
          <div 
            v-for="tag in tags" 
            :key="tag.id"
            class="flex items-center justify-between p-3 bg-gray-50 rounded-xl group"
          >
            <div class="flex items-center gap-3">
              <span class="text-gray-700 font-medium"># {{ tag.name }}</span>
              <span class="text-xs bg-white px-2 py-0.5 rounded-full text-gray-400">{{ tag.count }} 记录</span>
            </div>
            <!-- Delete button (placeholder for now as backend doesn't support explicit tag deletion separate from records yet) -->
            <!-- But we can edit? For MVP, just list them. -->
            <div class="opacity-0 group-hover:opacity-100 transition-opacity">
              <!-- Future: Edit/Delete -->
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { useLittleGrowthStore } from '@/stores/littleGrowth'

const router = useRouter()
const store = useLittleGrowthStore()

const tags = computed(() => store.flattenedTags)
</script>
