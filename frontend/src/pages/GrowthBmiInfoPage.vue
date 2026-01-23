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
            <div class="text-sm text-gray-500 dark:text-gray-400">BMI = 体重(kg) ÷ 身高(m)²，儿童需结合年龄与性别判断</div>
          </div>
        </div>
        <div class="mt-4 grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="rounded-xl bg-gradient-to-br from-blue-50 to-emerald-50 dark:from-blue-500/10 dark:to-emerald-500/10 overflow-hidden flex flex-col items-center">
            <el-image
              src="https://pub-54e58c1b9bee43ae8c6ebb5ab6e1ed3d.r2.dev/recordgo/Project-Images/BMI.webp"
              :preview-src-list="['https://pub-54e58c1b9bee43ae8c6ebb5ab6e1ed3d.r2.dev/recordgo/Project-Images/BMI.webp']"
              :preview-teleported="true"
              fit="cover"
              class="w-full h-full cursor-zoom-in"
              alt="BMI 参考图"
            />
          </div>
          <div class="rounded-xl bg-white/80 dark:bg-gray-900/40 border border-gray-100 dark:border-gray-700 p-4">
            <div class="text-sm font-semibold text-gray-800 dark:text-gray-100">解读建议</div>
            <ul class="mt-2 text-sm text-gray-600 dark:text-gray-300 space-y-2">
              <li>先对照同年龄与性别的参考范围，再判断偏低或偏高</li>
              <li>偏低：加强均衡饮食与优质蛋白摄入，必要时做体检评估</li>
              <li>偏高：减少高糖高脂摄入，增加有氧运动与日常活动量</li>
            </ul>
          </div>
        </div>
      </div>

      <div class="rounded-2xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4 shadow-sm space-y-3">
        <div class="flex items-center justify-between">
          <div class="text-base font-semibold text-gray-900 dark:text-gray-100">参考范围</div>
          <a
            href="https://bmi.golong.uk/article/children-bmi-report.html"
            target="_blank"
            class="text-xs text-blue-500 hover:text-blue-600 dark:text-blue-400 flex items-center gap-1"
          >
            <span>来源说明</span>
            <el-icon :size="12"><TopRight /></el-icon>
          </a>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
          <div class="flex items-start gap-3 rounded-xl p-3 bg-blue-50 dark:bg-blue-500/10">
            <div class="h-3 w-3 rounded-full bg-blue-500 mt-1"></div>
            <div>
              <div class="text-sm font-semibold text-blue-700 dark:text-blue-300">6 岁</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">正常范围约 15 - 18</div>
            </div>
          </div>
          <div class="flex items-start gap-3 rounded-xl p-3 bg-emerald-50 dark:bg-emerald-500/10">
            <div class="h-3 w-3 rounded-full bg-emerald-500 mt-1"></div>
            <div>
              <div class="text-sm font-semibold text-emerald-700 dark:text-emerald-300">12 岁</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">男孩约 17 - 21，女孩约 16 - 20</div>
            </div>
          </div>
          <div class="flex items-start gap-3 rounded-xl p-3 bg-amber-50 dark:bg-amber-500/10">
            <div class="h-3 w-3 rounded-full bg-amber-500 mt-1"></div>
            <div>
              <div class="text-sm font-semibold text-amber-700 dark:text-amber-300">18 岁</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">男孩约 19 - 24，女孩约 18 - 23</div>
            </div>
          </div>
        </div>
        <div class="text-xs text-gray-500 dark:text-gray-400">以上为参考范围，实际解读以同龄同区间参考为准，可点击“来源说明”进行查看。</div>
      </div>

      <div class="rounded-2xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4 shadow-sm">
        <div class="text-base font-semibold text-gray-900 dark:text-gray-100">小贴士</div>
        <div class="mt-2 grid grid-cols-1 md:grid-cols-3 gap-3 text-sm text-gray-600 dark:text-gray-300">
          <div class="rounded-xl bg-gray-50 dark:bg-gray-700/40 p-3">每天累计至少 1 小时户外运动，优先有氧活动</div>
          <div class="rounded-xl bg-gray-50 dark:bg-gray-700/40 p-3">三餐均衡，注意优质蛋白与蔬果摄入</div>
          <div class="rounded-xl bg-gray-50 dark:bg-gray-700/40 p-3">控制久坐与电子产品使用时长，保持规律作息</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ArrowLeft, TopRight } from '@element-plus/icons-vue'
import router from '@/router'
</script>

<style scoped>
</style>
