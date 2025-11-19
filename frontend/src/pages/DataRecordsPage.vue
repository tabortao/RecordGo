<template>
  <div class="p-4 space-y-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <el-icon :size="18" class="cursor-pointer text-gray-600" @click="goBack"><ArrowLeft /></el-icon>
        <h2 class="font-semibold">操作记录</h2>
      </div>
      <el-date-picker v-model="selectedDate" type="date" placeholder="选择日期" :editable="false" />
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <el-card shadow="hover" class="rounded-lg">
        <div class="flex items-center gap-2 mb-2">
          <el-icon :size="18" style="color:#10b981"><CircleCheck /></el-icon>
          <div class="font-semibold">任务完成记录</div>
        </div>
        <div v-if="taskRecords.length === 0" class="text-sm text-gray-500">无记录</div>
        <div v-else class="space-y-2">
          <div v-for="r in taskRecords" :key="r.key" class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <el-tag type="success" size="small">完成</el-tag>
              <div class="font-medium">{{ r.name }}</div>
              <div class="text-xs text-gray-500">{{ r.date }}</div>
            </div>
            <div class="flex items-center gap-1">
              <el-icon :size="18" style="color:#f59e0b"><Coin /></el-icon>
              <div class="font-semibold">+{{ r.coins }}</div>
            </div>
          </div>
        </div>
      </el-card>

      <el-card shadow="hover" class="rounded-lg">
        <div class="flex items-center gap-2 mb-2">
          <el-icon :size="18" style="color:#0ea5e9"><Document /></el-icon>
          <div class="font-semibold">心愿兑换记录</div>
        </div>
        <div v-if="wishRecordsForDay.length === 0" class="text-sm text-gray-500">无记录</div>
        <div v-else class="space-y-2">
          <div v-for="w in wishRecordsForDay" :key="w.id" class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <el-tag type="warning" size="small">兑换</el-tag>
              <div class="font-medium">{{ w.wish_name }}</div>
              <div class="text-xs text-gray-500">{{ formatDate(w.created_at) }}</div>
            </div>
            <div class="flex items-center gap-1">
              <el-icon :size="18" style="color:#ef4444"><Coin /></el-icon>
              <div class="font-semibold">-{{ w.coins_used }}</div>
            </div>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { ArrowLeft, CircleCheck, Coin, Document } from '@element-plus/icons-vue'
import router from '@/router'
import dayjs from 'dayjs'
import { listTaskOccurrences, listTasks } from '@/services/tasks'
import { listWishRecords } from '@/services/wishes'
import { useAuth } from '@/stores/auth'

function goBack() { router.back() }
const auth = useAuth()
const selectedDate = ref<string>(dayjs().format('YYYY-MM-DD'))
const tasksMap = ref<Record<number, { name: string; score: number }>>({})
const taskRecords = ref<{ key: string; name: string; coins: number; date: string }[]>([])
const wishRecords = ref<any[]>([])

function formatDate(ts: string) { return dayjs(ts).format('YYYY-MM-DD HH:mm') }

async function loadTasksMap() {
  try {
    const resp = await listTasks({ page_size: 500 })
    const map: Record<number, { name: string; score: number }> = {}
    for (const t of resp.items || []) map[t.id] = { name: t.name, score: t.score }
    tasksMap.value = map
  } catch { tasksMap.value = {} }
}

async function loadTaskRecords() {
  try {
    const resp = await listTaskOccurrences({ date: selectedDate.value })
    const out: { key: string; name: string; coins: number; date: string }[] = []
    for (const it of resp.items || []) {
      if (Number(it.status) !== 2) continue
      const t = tasksMap.value[it.task_id]
      if (!t) continue
      out.push({ key: `${it.task_id}:${it.date}`, name: t.name, coins: t.score, date: it.date })
    }
    taskRecords.value = out
  } catch { taskRecords.value = [] }
}

async function loadWishRecords() {
  try {
    const uid = auth.user?.id || 0
    const resp = await listWishRecords(uid, 1, 100)
    wishRecords.value = resp.items || []
  } catch { wishRecords.value = [] }
}

const wishRecordsForDay = computed(() => {
  const d = selectedDate.value
  return (wishRecords.value || []).filter((x: any) => dayjs(x.created_at).format('YYYY-MM-DD') === d)
})

onMounted(async () => { await loadTasksMap(); await loadTaskRecords(); await loadWishRecords() })
watch(selectedDate, async () => { await loadTaskRecords() })
</script>

<style scoped>
</style>