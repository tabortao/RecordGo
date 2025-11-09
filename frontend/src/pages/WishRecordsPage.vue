<template>
  <!-- 中文注释：独立的心愿兑换记录页面，顶部返回图标；响应式卡片列表 -->
  <div class="p-4 space-y-4">
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer text-gray-600" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" class="text-emerald-600"><List /></el-icon>
      <h2 class="font-semibold">兑换记录</h2>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <el-card v-for="item in records" :key="item.id" class="rounded-lg" shadow="hover">
        <div class="flex justify-between items-start">
          <div>
            <div class="font-medium">{{ item.wish_name }}</div>
            <div class="text-xs text-gray-500 mt-1">{{ item.created_at }}</div>
          </div>
          <el-tag :type="item.status === 'success' ? 'success' : 'info'">{{ item.status }}</el-tag>
        </div>
        <div class="mt-3 text-sm text-gray-700">数量：{{ item.amount }} {{ item.unit }}</div>
        <div class="mt-1 text-sm text-gray-700">消耗金币：{{ item.coins_used }}</div>
      </el-card>
    </div>

    <!-- 分页 -->
    <div class="flex justify-end mt-2">
      <el-pagination
        layout="prev, pager, next"
        :current-page="page"
        :page-size="pageSize"
        :total="total"
        @current-change="onPageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
// 中文注释：加载心愿兑换记录，支持分页，点击返回图标返回上一页
import { ref, onMounted } from 'vue'
import router from '@/router'
import { ArrowLeft, List } from '@element-plus/icons-vue'
import { listWishRecords, type WishRecord } from '@/services/wishes'

const userId = 1 // 中文注释：示例用户ID
function goBack() { router.back() }

const records = ref<WishRecord[]>([])
const page = ref(1)
const pageSize = ref(9)
const total = ref(0)

async function load() {
  const resp = await listWishRecords(userId, page.value, pageSize.value)
  records.value = resp.items
  total.value = resp.total
}

function onPageChange(p: number) { page.value = p; load() }

onMounted(load)
</script>

<style scoped>
/* 中文注释：布局主要使用 Tailwind 与 Element Plus 组件 */
</style>