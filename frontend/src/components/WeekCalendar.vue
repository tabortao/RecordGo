<template>
  <!-- 中文注释：周视图组件，支持左右切周、返回本周、点击日期切换 -->
  <div class="space-y-2">
    <!-- 顶部周信息与控制 -->
    <div class="flex items-center justify-center gap-3">
      <el-button text :icon="ArrowLeft" @click="prevWeek" />
      <div class="text-sm font-semibold">{{ weekTitle }}</div>
      <el-button text :icon="ArrowRight" @click="nextWeek" />
      <el-tag type="success" class="ml-2 rounded-full px-3 py-1 cursor-pointer" @click="goThisWeek">本周</el-tag>
    </div>
    <!-- 周一到周日日期格 -->
    <div class="grid grid-cols-7 gap-2">
      <div
        v-for="d in days"
        :key="d.format('YYYY-MM-DD')"
        class="p-2 rounded-lg text-center cursor-pointer border"
        :class="{
          'bg-green-50 border-green-300': d.isSame(selected, 'date'),
          'hover:bg-gray-50': !d.isSame(selected, 'date')
        }"
        @click="onPick(d)"
      >
        <div class="text-xs text-gray-500">{{ weekdayLabel(d) }}</div>
        <div class="text-lg font-semibold">{{ d.date() }}</div>
        <div class="h-5 flex items-center justify-center">
          <span v-if="countMap[d.format('YYYY-MM-DD')]" class="inline-flex items-center gap-1 text-xs text-purple-600">
            <span class="w-2 h-2 rounded-full bg-purple-400"></span>
            {{ countMap[d.format('YYYY-MM-DD')] }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 中文注释：使用 dayjs 计算当前周的起止，周一为一周开始；对外发出日期与周切换事件
import dayjs, { Dayjs } from 'dayjs'
import isoWeek from 'dayjs/plugin/isoWeek'
import weekOfYear from 'dayjs/plugin/weekOfYear'
import { ArrowLeft, ArrowRight } from '@element-plus/icons-vue'
dayjs.extend(isoWeek)
dayjs.extend(weekOfYear)

const props = defineProps<{ selectedDate: string; taskCountMap?: Record<string, number> }>()
const emit = defineEmits<{ (e: 'update:selectedDate', v: string): void }>()

const selected = computed(() => dayjs(props.selectedDate))
const startOfWeek = ref<Dayjs>(selected.value.isoWeekday(1))
const countMap = computed(() => props.taskCountMap || {})

const days = computed(() => {
  const arr: Dayjs[] = []
  for (let i = 0; i < 7; i++) {
    arr.push(startOfWeek.value.add(i, 'day'))
  }
  return arr
})

const weekTitle = computed(() => {
  const d = selected.value
  return `${d.year()} 年 ${d.month() + 1} 月第 ${d.week()} 周`
})

function weekdayLabel(d: Dayjs) {
  const map = ['一', '二', '三', '四', '五', '六', '日']
  return `周${map[d.isoWeekday() - 1]}`
}

function onPick(d: Dayjs) {
  emit('update:selectedDate', d.format('YYYY-MM-DD'))
}
function prevWeek() {
  startOfWeek.value = startOfWeek.value.add(-7, 'day')
}
function nextWeek() {
  startOfWeek.value = startOfWeek.value.add(7, 'day')
}
function goThisWeek() {
  const now = dayjs()
  startOfWeek.value = now.isoWeekday(1)
  emit('update:selectedDate', now.format('YYYY-MM-DD'))
}
</script>

<style scoped>
/* 中文注释：基本样式以边框与背景色区分选中态与悬停态 */
</style>
