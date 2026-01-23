<template>
  <div class="min-h-screen bg-[#F5F7FA] dark:bg-gray-900">
    <div class="bg-white/80 dark:bg-gray-900/80 backdrop-blur border-b border-gray-200 dark:border-gray-800 sticky top-0 z-10">
      <div class="px-4 py-3 flex items-center gap-2">
        <el-icon :size="18" class="cursor-pointer text-gray-600 dark:text-gray-300" @click="router.back()"><ArrowLeft /></el-icon>
        <h2 class="font-semibold text-gray-800 dark:text-gray-100">BMI 指标说明</h2>
      </div>
    </div>

    <div class="p-4 space-y-4">
      <div class="rounded-2xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4 shadow-sm">
        <div class="flex items-center gap-3">
          <div class="h-12 w-12 rounded-xl bg-gradient-to-br from-blue-500 to-emerald-400 flex items-center justify-center text-white text-lg font-semibold">BMI</div>
          <div>
            <div class="text-base font-semibold text-gray-900 dark:text-gray-100">身体质量指数</div>
            <div class="text-sm text-gray-500 dark:text-gray-400">BMI = 体重(kg) ÷ 身高(m)²</div>
          </div>
        </div>
        <div class="mt-4 grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="rounded-xl bg-gradient-to-br from-blue-50 to-emerald-50 dark:from-blue-500/10 dark:to-emerald-500/10 p-4 flex flex-col items-center">
            <svg viewBox="0 0 200 200" class="h-36 w-36">
              <path d="M 20 140 A 80 80 0 0 1 180 140" stroke="currentColor" stroke-width="14" fill="none" class="text-gray-200 dark:text-gray-700" />
              <path v-for="seg in bmiArcSegments" :key="seg.color" :d="seg.path" :stroke="seg.color" stroke-width="14" fill="none" stroke-linecap="butt" />
              <line x1="100" y1="140" :x2="bmiPointer.x" :y2="bmiPointer.y" stroke="#ef4444" stroke-width="4" stroke-linecap="round" />
              <circle :cx="bmiPointer.x" :cy="bmiPointer.y" r="5" fill="#ef4444" />
              <text x="100" y="112" text-anchor="middle" fill="currentColor" class="text-[12px] text-gray-500 dark:text-gray-400">BMI</text>
              <text x="100" y="132" text-anchor="middle" fill="currentColor" class="text-lg font-semibold text-gray-900 dark:text-gray-100">16.6</text>
            </svg>
            <div class="text-xs text-gray-500 dark:text-gray-400">半圆色带展示不同 BMI 区间</div>
          </div>
          <div class="rounded-xl bg-white/80 dark:bg-gray-900/40 border border-gray-100 dark:border-gray-700 p-4">
            <div class="text-sm font-semibold text-gray-800 dark:text-gray-100">解读建议</div>
            <ul class="mt-2 text-sm text-gray-600 dark:text-gray-300 space-y-2">
              <li>数值偏低：注意营养均衡与优质蛋白摄入</li>
              <li>数值正常：保持运动与规律作息</li>
              <li>数值偏高：控制油脂糖分，增加有氧运动</li>
            </ul>
          </div>
        </div>
      </div>

      <div class="rounded-2xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4 shadow-sm space-y-3">
        <div class="text-base font-semibold text-gray-900 dark:text-gray-100">区间说明</div>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <div class="flex items-center gap-3 rounded-xl p-3 bg-blue-50 dark:bg-blue-500/10">
            <div class="h-3 w-3 rounded-full bg-blue-500"></div>
            <div>
              <div class="text-sm font-semibold text-blue-700 dark:text-blue-300">0 - 18.5</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">偏瘦，需关注饮食与体重增长</div>
            </div>
          </div>
          <div class="flex items-center gap-3 rounded-xl p-3 bg-emerald-50 dark:bg-emerald-500/10">
            <div class="h-3 w-3 rounded-full bg-emerald-500"></div>
            <div>
              <div class="text-sm font-semibold text-emerald-700 dark:text-emerald-300">18.5 - 25</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">正常区间，保持良好生活习惯</div>
            </div>
          </div>
          <div class="flex items-center gap-3 rounded-xl p-3 bg-amber-50 dark:bg-amber-500/10">
            <div class="h-3 w-3 rounded-full bg-amber-500"></div>
            <div>
              <div class="text-sm font-semibold text-amber-700 dark:text-amber-300">25 - 30</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">偏胖，建议提升运动频率</div>
            </div>
          </div>
          <div class="flex items-center gap-3 rounded-xl p-3 bg-rose-50 dark:bg-rose-500/10">
            <div class="h-3 w-3 rounded-full bg-rose-300"></div>
            <div>
              <div class="text-sm font-semibold text-rose-600 dark:text-rose-300">30 - 40</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">肥胖风险，注意健康管理</div>
            </div>
          </div>
          <div class="flex items-center gap-3 rounded-xl p-3 bg-red-50 dark:bg-red-500/10">
            <div class="h-3 w-3 rounded-full bg-red-500"></div>
            <div>
              <div class="text-sm font-semibold text-red-700 dark:text-red-300">40 - 50</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">高度肥胖，建议专业指导</div>
            </div>
          </div>
        </div>
      </div>

      <div class="rounded-2xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4 shadow-sm">
        <div class="text-base font-semibold text-gray-900 dark:text-gray-100">小贴士</div>
        <div class="mt-2 grid grid-cols-1 md:grid-cols-3 gap-3 text-sm text-gray-600 dark:text-gray-300">
          <div class="rounded-xl bg-gray-50 dark:bg-gray-700/40 p-3">每周至少三次运动，每次 30 分钟以上</div>
          <div class="rounded-xl bg-gray-50 dark:bg-gray-700/40 p-3">保持充足睡眠，规律作息有助于生长</div>
          <div class="rounded-xl bg-gray-50 dark:bg-gray-700/40 p-3">记录身高体重变化，观察长期趋势</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ArrowLeft } from '@element-plus/icons-vue'
import router from '@/router'

const bmiSegments = [
  { start: 0, end: 18.5, color: '#3b82f6' },
  { start: 18.5, end: 25, color: '#22c55e' },
  { start: 25, end: 30, color: '#f59e0b' },
  { start: 30, end: 40, color: '#fca5a5' },
  { start: 40, end: 50, color: '#ef4444' }
]

function valueToAngle(value: number) {
  const v = Math.max(0, Math.min(50, value))
  return 180 - (v / 50) * 180
}

function polarToCartesian(cx: number, cy: number, r: number, angle: number) {
  const rad = (angle - 90) * (Math.PI / 180)
  return { x: cx + r * Math.cos(rad), y: cy + r * Math.sin(rad) }
}

function describeArc(cx: number, cy: number, r: number, startAngle: number, endAngle: number) {
  const start = polarToCartesian(cx, cy, r, endAngle)
  const end = polarToCartesian(cx, cy, r, startAngle)
  const largeArcFlag = endAngle - startAngle <= 180 ? '0' : '1'
  return `M ${start.x} ${start.y} A ${r} ${r} 0 ${largeArcFlag} 0 ${end.x} ${end.y}`
}

const bmiArcSegments = computed(() => {
  return bmiSegments.map((seg) => {
    const startAngle = valueToAngle(seg.start)
    const endAngle = valueToAngle(seg.end)
    return { color: seg.color, path: describeArc(100, 140, 80, startAngle, endAngle) }
  })
})

const bmiPointer = computed(() => {
  const angle = valueToAngle(16.64)
  return polarToCartesian(100, 140, 62, angle)
})
</script>

<style scoped>
</style>
