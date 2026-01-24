<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 pb-20 transition-colors duration-300">
    <!-- 顶部导航栏 -->
    <div class="sticky top-0 z-30 bg-white/80 dark:bg-gray-900/80 backdrop-blur-md border-b border-gray-200 dark:border-gray-800 px-4 h-14 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <button 
          @click="goBack"
          class="p-2 -ml-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors text-gray-600 dark:text-gray-300">
          <el-icon :size="20"><ArrowLeft /></el-icon>
        </button>
        <span class="font-bold text-lg text-gray-800 dark:text-gray-100">操作记录</span>
      </div>
      <div class="flex items-center gap-2">
        <el-popover placement="bottom-end" trigger="click" width="auto" :show-arrow="false" popper-class="!p-2 !rounded-xl !shadow-xl !border-gray-100 dark:!border-gray-800 dark:!bg-gray-900">
          <template #reference>
            <button class="flex items-center gap-1.5 px-3 py-1.5 bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 rounded-lg text-sm text-gray-700 dark:text-gray-300 transition-colors">
              <el-icon><Calendar /></el-icon>
              <span>{{ dayjs(selectedDate).format('MM-DD') }}</span>
            </button>
          </template>
          <div class="space-y-2 p-1">
            <div class="text-xs font-medium text-gray-500 dark:text-gray-400 px-1">选择日期</div>
            <el-date-picker 
              v-model="selectedDate" 
              type="date" 
              placeholder="选择日期" 
              :editable="false" 
              value-format="YYYY-MM-DD" 
              :disabled-date="(d) => d > new Date()"
              class="!w-full custom-date-picker"
            />
          </div>
        </el-popover>
      </div>
    </div>

    <div class="max-w-4xl mx-auto p-4 md:p-6 space-y-6 animate-fade-in-up">
      <!-- 任务完成记录 -->
      <div class="bg-white dark:bg-gray-900 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-800 overflow-hidden">
        <div class="px-5 py-4 border-b border-gray-50 dark:border-gray-800/50 flex items-center justify-between bg-gradient-to-r from-green-50/50 to-transparent dark:from-green-900/10">
          <div class="flex items-center gap-2.5">
            <div class="w-10 h-10 rounded-full bg-green-100 dark:bg-green-900/30 flex items-center justify-center text-green-600 dark:text-green-400">
               <el-icon :size="20"><CircleCheck /></el-icon>
            </div>
            <div>
              <div class="font-bold text-gray-800 dark:text-gray-100">任务完成</div>
              <div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5" v-if="taskRecords.length > 0">
                 共获得 <span class="text-orange-500 font-bold">+{{ taskTotalCoins }}</span> 金币
              </div>
            </div>
          </div>
          <el-tag v-if="taskRecords.length > 0" type="success" effect="light" round size="small">{{ taskRecords.length }} 项</el-tag>
        </div>
        
        <div v-if="taskRecords.length === 0" class="flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500">
           <el-icon :size="48" class="mb-2 opacity-20"><CircleCheck /></el-icon>
           <span class="text-sm">暂无完成记录</span>
        </div>
        
        <div v-else class="divide-y divide-gray-50 dark:divide-gray-800/50">
          <div v-for="r in taskRecords" :key="r.key" class="p-4 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-800/30 transition-colors">
            <div class="flex items-center gap-3">
              <div class="w-1.5 h-1.5 rounded-full bg-green-400"></div>
              <div>
                <div class="font-medium text-gray-800 dark:text-gray-200">{{ r.name }}</div>
                <div class="text-xs text-gray-400 mt-0.5">{{ formatTime(r.date) }}</div>
              </div>
            </div>
            <div class="flex items-center gap-1 bg-orange-50 dark:bg-orange-900/10 px-2 py-1 rounded-lg border border-orange-100 dark:border-orange-900/20">
              <el-icon :size="14" class="text-orange-500"><Coin /></el-icon>
              <span class="font-bold text-orange-500 text-sm">+{{ r.coins }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 心愿兑换记录 -->
      <div class="bg-white dark:bg-gray-900 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-800 overflow-hidden">
        <div class="px-5 py-4 border-b border-gray-50 dark:border-gray-800/50 flex items-center justify-between bg-gradient-to-r from-blue-50/50 to-transparent dark:from-blue-900/10">
          <div class="flex items-center gap-2.5">
            <div class="w-10 h-10 rounded-full bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center text-blue-600 dark:text-blue-400">
               <el-icon :size="20"><Present /></el-icon>
            </div>
            <div>
              <div class="font-bold text-gray-800 dark:text-gray-100">心愿兑换</div>
              <div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5" v-if="wishRecordsForDay.length > 0">
                 共消耗 <span class="text-gray-500 font-bold">-{{ wishTotalCoins }}</span> 金币
              </div>
            </div>
          </div>
          <el-tag v-if="wishRecordsForDay.length > 0" type="primary" effect="light" round size="small">{{ wishRecordsForDay.length }} 项</el-tag>
        </div>
        
        <div v-if="wishRecordsForDay.length === 0" class="flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500">
           <el-icon :size="48" class="mb-2 opacity-20"><Present /></el-icon>
           <span class="text-sm">暂无兑换记录</span>
        </div>
        
        <div v-else class="divide-y divide-gray-50 dark:divide-gray-800/50">
          <div v-for="w in wishRecordsForDay" :key="w.id" class="p-4 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-800/30 transition-colors">
            <div class="flex items-center gap-3">
              <div class="w-1.5 h-1.5 rounded-full bg-blue-400"></div>
              <div>
                <div class="font-medium text-gray-800 dark:text-gray-200">{{ w.wish_name }}</div>
                <div class="text-xs text-gray-400 mt-0.5">{{ formatTime(w.created_at) }}</div>
              </div>
            </div>
            <div class="flex items-center gap-1 bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded-lg border border-gray-200 dark:border-gray-700">
              <el-icon :size="14" class="text-gray-500"><Coin /></el-icon>
              <span class="font-bold text-gray-500 text-sm">-{{ w.coins_used }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { ArrowLeft, CircleCheck, Coin, Present, Calendar } from '@element-plus/icons-vue'
import router from '@/router'
import dayjs from 'dayjs'
import { listTaskOccurrences, listTasks, type TaskItem } from '@/services/tasks'
import { listWishRecords } from '@/services/wishes'
import { useAuth } from '@/stores/auth'

function goBack() { router.back() }
const auth = useAuth()
const selectedDate = ref<string>(dayjs().format('YYYY-MM-DD'))
const tasksMap = ref<Record<number, { name: string; score: number }>>({})
const taskRecords = ref<{ key: string; name: string; coins: number; date: string }[]>([])
const wishRecords = ref<any[]>([])
const allTasks = ref<TaskItem[]>([])

function formatDate(ts: string) { return dayjs(ts).format('YYYY-MM-DD HH:mm') }
function formatTime(ts: string) { return dayjs(ts).format('HH:mm') }

const taskTotalCoins = computed(() => taskRecords.value.reduce((sum, r) => sum + r.coins, 0))
const wishTotalCoins = computed(() => wishRecordsForDay.value.reduce((sum, w) => sum + w.coins_used, 0))

async function loadTasksMap() {
  try {
    const resp = await listTasks({ page_size: 500 })
    const map: Record<number, { name: string; score: number }> = {}
    const items = resp.items || []
    for (const t of items) map[t.id] = { name: t.name, score: t.score }
    tasksMap.value = map
    allTasks.value = items as TaskItem[]
  } catch { tasksMap.value = {} }
}

async function loadTaskRecords() {
  try {
    const d = selectedDate.value
    const resp = await listTaskOccurrences({ start: d, end: d })
    const out: { key: string; name: string; coins: number; date: string }[] = []
    const added = new Set<string>()
    for (const it of resp.items || []) {
      if (Number(it.status) !== 2) continue
      const t = tasksMap.value[it.task_id]
      if (!t) continue
      const key = `${it.task_id}:${dayjs(it.date).format('YYYY-MM-DD')}`
      added.add(key)
      out.push({ key, name: t.name, coins: t.score, date: it.date })
    }
    const dKey = dayjs(d).format('YYYY-MM-DD')
    for (const t of allTasks.value || []) {
      const rep = String((t as any).repeat || 'none').toLowerCase()
      const isNonRepeat = /none|无|^$/i.test(rep) || !(t as any).end_date
      if (!isNonRepeat) continue
      if (Number(t.status) !== 2) continue
      const sKey = dayjs(t.start_date).format('YYYY-MM-DD')
      if (sKey !== dKey) continue
      const key = `${t.id}:${dKey}`
      if (added.has(key)) continue
      out.push({ key, name: t.name, coins: t.score, date: dayjs(t.start_date).format('YYYY-MM-DD HH:mm') })
    }
    taskRecords.value = out
  } catch { taskRecords.value = [] }
}

async function loadWishRecords() {
  try {
    const uid = auth.user?.id || 0
    const d = selectedDate.value
    const resp = await listWishRecords(uid, 1, 200, { start: d, end: d })
    wishRecords.value = resp.items || []
  } catch { wishRecords.value = [] }
}

const wishRecordsForDay = computed(() => {
  const d = selectedDate.value
  return (wishRecords.value || []).filter((x: any) => dayjs(x.created_at).format('YYYY-MM-DD') === d)
})

onMounted(async () => { await loadTasksMap(); await loadTaskRecords(); await loadWishRecords() })
watch(selectedDate, async () => { await loadTaskRecords(); await loadWishRecords() })
</script>

<style scoped>
@keyframes fade-in-up {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in-up {
  animation: fade-in-up 0.5s ease-out forwards;
}

:deep(.custom-date-picker .el-input__wrapper) {
  box-shadow: none !important;
  background-color: transparent !important;
}
</style>
