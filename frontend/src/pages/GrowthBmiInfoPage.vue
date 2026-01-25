<template>
  <SettingsShell title="BMI 指标说明" subtitle="公式、参考范围与建议" :icon="TrendCharts" tone="sky" container-class="max-w-3xl" back-to="/growth">
    <SettingsCard>
      <div class="flex items-start gap-4">
        <div class="relative shrink-0">
          <div class="absolute -inset-3 rounded-[1.5rem] bg-sky-200/55 dark:bg-sky-500/18 blur-2xl" />
          <div class="relative h-14 w-14 rounded-3xl border border-white/60 dark:border-gray-800/60 bg-white/70 dark:bg-gray-950/30 ring-1 ring-black/5 dark:ring-white/10 grid place-items-center">
            <div class="text-[13px] font-black tracking-[0.22em] text-sky-700 dark:text-sky-300">BMI</div>
          </div>
        </div>

        <div class="min-w-0 flex-1">
          <div class="text-base font-extrabold tracking-tight text-gray-900 dark:text-gray-50">身体质量指数</div>
          <div class="mt-1 text-sm leading-relaxed text-gray-600 dark:text-gray-300">
            BMI = 体重(kg) ÷ 身高(m)²；儿童建议结合年龄与性别的参考范围一起解读。
          </div>
          <div class="mt-3 flex flex-wrap gap-2">
            <span class="rounded-full border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 px-3 py-1 text-xs font-semibold text-gray-700 dark:text-gray-300">先看范围</span>
            <span class="rounded-full border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 px-3 py-1 text-xs font-semibold text-gray-700 dark:text-gray-300">再看趋势</span>
            <span class="rounded-full border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 px-3 py-1 text-xs font-semibold text-gray-700 dark:text-gray-300">少看单次波动</span>
          </div>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard title="参考图" description="点击图片可放大查看">
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/40 overflow-hidden">
          <el-image
            src="https://pub-54e58c1b9bee43ae8c6ebb5ab6e1ed3d.r2.dev/recordgo/Project-Images/BMI.webp"
            :preview-src-list="['https://pub-54e58c1b9bee43ae8c6ebb5ab6e1ed3d.r2.dev/recordgo/Project-Images/BMI.webp']"
            :preview-teleported="true"
            fit="cover"
            class="w-full h-full cursor-zoom-in"
            alt="BMI 参考图"
          />
        </div>

        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/40 px-4 py-4">
          <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">解读建议</div>
          <div class="mt-3 space-y-2 text-sm leading-relaxed text-gray-600 dark:text-gray-300">
            <div class="flex gap-2">
              <div class="mt-1 h-1.5 w-1.5 rounded-full bg-sky-500" />
              <div>先对照同年龄与性别的参考范围，再判断偏低或偏高。</div>
            </div>
            <div class="flex gap-2">
              <div class="mt-1 h-1.5 w-1.5 rounded-full bg-emerald-500" />
              <div>偏低：均衡饮食与优质蛋白摄入；必要时结合体检评估。</div>
            </div>
            <div class="flex gap-2">
              <div class="mt-1 h-1.5 w-1.5 rounded-full bg-amber-500" />
              <div>偏高：减少高糖高脂摄入；增加有氧运动与日常活动量。</div>
            </div>
          </div>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard title="参考范围" description="以下为示例范围，实际以同龄同区间参考为准">
      <template #right>
        <a
          href="https://bmi.golong.uk/article/children-bmi-report.html"
          target="_blank"
          class="inline-flex items-center gap-1 rounded-full border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 px-3 py-1 text-xs font-semibold text-sky-700 dark:text-sky-300 hover:bg-white/90 dark:hover:bg-gray-900/80 transition"
        >
          <span>来源说明</span>
          <el-icon :size="12"><TopRight /></el-icon>
        </a>
      </template>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
        <div class="rounded-3xl border border-blue-100/70 dark:border-blue-900/40 bg-blue-50/70 dark:bg-blue-900/15 px-4 py-4">
          <div class="text-xs font-extrabold text-blue-700 dark:text-blue-300">6 岁</div>
          <div class="mt-1 text-sm font-bold text-gray-900 dark:text-gray-50">约 15 - 18</div>
          <div class="mt-2 text-xs text-gray-500 dark:text-gray-400">以同年龄分位区间为准</div>
        </div>
        <div class="rounded-3xl border border-emerald-100/70 dark:border-emerald-900/40 bg-emerald-50/70 dark:bg-emerald-900/15 px-4 py-4">
          <div class="text-xs font-extrabold text-emerald-700 dark:text-emerald-300">12 岁</div>
          <div class="mt-1 text-sm font-bold text-gray-900 dark:text-gray-50">男约 17 - 21 / 女约 16 - 20</div>
          <div class="mt-2 text-xs text-gray-500 dark:text-gray-400">男女参考范围不同</div>
        </div>
        <div class="rounded-3xl border border-amber-100/70 dark:border-amber-900/40 bg-amber-50/70 dark:bg-amber-900/15 px-4 py-4">
          <div class="text-xs font-extrabold text-amber-700 dark:text-amber-300">18 岁</div>
          <div class="mt-1 text-sm font-bold text-gray-900 dark:text-gray-50">男约 19 - 24 / 女约 18 - 23</div>
          <div class="mt-2 text-xs text-gray-500 dark:text-gray-400">逐渐接近成人参考</div>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard title="小贴士" description="更有效的方式是长期可执行的习惯">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-3 text-sm text-gray-700 dark:text-gray-200">
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/60 dark:bg-gray-900/35 px-4 py-4">
          每天累计至少 1 小时户外运动，优先有氧活动
        </div>
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/60 dark:bg-gray-900/35 px-4 py-4">
          三餐均衡，注意优质蛋白与蔬果摄入
        </div>
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/60 dark:bg-gray-900/35 px-4 py-4">
          控制久坐与电子产品使用时长，保持规律作息
        </div>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { TrendCharts, TopRight } from '@element-plus/icons-vue'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'
</script>

<style scoped>
</style>
