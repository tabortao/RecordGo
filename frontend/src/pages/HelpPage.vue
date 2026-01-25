<template>
  <SettingsShell title="使用帮助" :icon="QuestionFilled" tone="sky" container-class="max-w-3xl">
    <SettingsCard>
      <div class="flex items-start gap-4">
        <div class="relative shrink-0">
          <div class="absolute -inset-3 rounded-[1.5rem] bg-sky-200/50 dark:bg-sky-500/15 blur-2xl" />
          <div class="relative h-14 w-14 rounded-3xl border border-white/60 dark:border-gray-800/60 bg-white/70 dark:bg-gray-950/30 ring-1 ring-black/5 dark:ring-white/10 grid place-items-center overflow-hidden">
            <img :src="appIconSrc" alt="App Icon" class="h-10 w-10 rounded-2xl shadow-sm ring-1 ring-white/60 dark:ring-gray-800/70" />
          </div>
        </div>
        <div class="min-w-0 flex-1">
          <div class="text-base font-extrabold tracking-tight text-gray-900 dark:text-gray-50">作业家：任务积分 + 学习激励</div>
          <div class="mt-1 text-sm leading-relaxed text-gray-600 dark:text-gray-300">
            通过任务、番茄专注、积分与心愿兑换，把“完成一件事”变成看得见的进步；也可以用生长记录，持续追踪成长变化。
          </div>
          <div class="mt-3 flex flex-wrap gap-2">
            <span class="rounded-full border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 px-3 py-1 text-xs font-semibold text-gray-700 dark:text-gray-300">任务</span>
            <span class="rounded-full border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 px-3 py-1 text-xs font-semibold text-gray-700 dark:text-gray-300">积分</span>
            <span class="rounded-full border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 px-3 py-1 text-xs font-semibold text-gray-700 dark:text-gray-300">番茄</span>
            <span class="rounded-full border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 px-3 py-1 text-xs font-semibold text-gray-700 dark:text-gray-300">心愿</span>
            <span class="rounded-full border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 px-3 py-1 text-xs font-semibold text-gray-700 dark:text-gray-300">成长记录</span>
          </div>
          <div class="mt-4 grid grid-cols-1 sm:grid-cols-3 gap-2">
            <button
              class="group flex items-center justify-between gap-3 rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/55 px-3 py-2.5 hover:bg-white dark:hover:bg-gray-900/75 transition active:scale-[0.99]"
              @click="router.push('/tasks/create')"
            >
              <div class="flex items-center gap-2">
                <div class="h-9 w-9 rounded-2xl bg-emerald-50 dark:bg-emerald-900/25 text-emerald-700 dark:text-emerald-300 grid place-items-center">
                  <el-icon :size="18"><CirclePlusFilled /></el-icon>
                </div>
                <div class="text-left">
                  <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">创建任务</div>
                  <div class="text-[11px] text-gray-500 dark:text-gray-400">从一个小目标开始</div>
                </div>
              </div>
              <el-icon :size="18" class="text-gray-300 dark:text-gray-600 group-hover:text-gray-500 transition-colors"><ArrowRight /></el-icon>
            </button>

            <button
              class="group flex items-center justify-between gap-3 rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/55 px-3 py-2.5 hover:bg-white dark:hover:bg-gray-900/75 transition active:scale-[0.99]"
              @click="router.push('/tasks/stats')"
            >
              <div class="flex items-center gap-2">
                <div class="h-9 w-9 rounded-2xl bg-indigo-50 dark:bg-indigo-900/25 text-indigo-700 dark:text-indigo-300 grid place-items-center">
                  <el-icon :size="18"><DataAnalysis /></el-icon>
                </div>
                <div class="text-left">
                  <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">查看统计</div>
                  <div class="text-[11px] text-gray-500 dark:text-gray-400">复盘更有动力</div>
                </div>
              </div>
              <el-icon :size="18" class="text-gray-300 dark:text-gray-600 group-hover:text-gray-500 transition-colors"><ArrowRight /></el-icon>
            </button>

            <button
              class="group flex items-center justify-between gap-3 rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/55 px-3 py-2.5 hover:bg-white dark:hover:bg-gray-900/75 transition active:scale-[0.99]"
              @click="router.push('/settings/reading')"
            >
              <div class="flex items-center gap-2">
                <div class="h-9 w-9 rounded-2xl bg-sky-50 dark:bg-sky-900/25 text-sky-700 dark:text-sky-300 grid place-items-center">
                  <el-icon :size="18"><Microphone /></el-icon>
                </div>
                <div class="text-left">
                  <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">朗读设置</div>
                  <div class="text-[11px] text-gray-500 dark:text-gray-400">辅助专注与复述</div>
                </div>
              </div>
              <el-icon :size="18" class="text-gray-300 dark:text-gray-600 group-hover:text-gray-500 transition-colors"><ArrowRight /></el-icon>
            </button>
          </div>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard title="快速上手">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/65 dark:bg-gray-900/40 px-4 py-4">
          <div class="flex items-center justify-between">
            <div class="text-[11px] font-extrabold tracking-[0.22em] text-gray-500 dark:text-gray-400">STEP 01</div>
            <div class="h-9 w-9 rounded-2xl bg-sky-50 dark:bg-sky-900/25 text-sky-700 dark:text-sky-300 grid place-items-center">
              <el-icon :size="18"><Edit /></el-icon>
            </div>
          </div>
          <div class="mt-2 text-sm font-black text-gray-900 dark:text-gray-50">创建任务</div>
          <div class="mt-1 text-xs leading-relaxed text-gray-600 dark:text-gray-300">把作业拆成可完成的小任务，并设置奖励金币。</div>
        </div>
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/65 dark:bg-gray-900/40 px-4 py-4">
          <div class="flex items-center justify-between">
            <div class="text-[11px] font-extrabold tracking-[0.22em] text-gray-500 dark:text-gray-400">STEP 02</div>
            <div class="h-9 w-9 rounded-2xl bg-amber-50 dark:bg-amber-900/25 text-amber-700 dark:text-amber-300 grid place-items-center">
              <el-icon :size="18"><Clock /></el-icon>
            </div>
          </div>
          <div class="mt-2 text-sm font-black text-gray-900 dark:text-gray-50">番茄专注</div>
          <div class="mt-1 text-xs leading-relaxed text-gray-600 dark:text-gray-300">用 20-25 分钟的专注块，完成一小段目标。</div>
        </div>
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/65 dark:bg-gray-900/40 px-4 py-4">
          <div class="flex items-center justify-between">
            <div class="text-[11px] font-extrabold tracking-[0.22em] text-gray-500 dark:text-gray-400">STEP 03</div>
            <div class="h-9 w-9 rounded-2xl bg-emerald-50 dark:bg-emerald-900/25 text-emerald-700 dark:text-emerald-300 grid place-items-center">
              <el-icon :size="18"><Present /></el-icon>
            </div>
          </div>
          <div class="mt-2 text-sm font-black text-gray-900 dark:text-gray-50">兑换心愿</div>
          <div class="mt-1 text-xs leading-relaxed text-gray-600 dark:text-gray-300">用金币兑换奖励，形成稳定的正向反馈。</div>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard title="使用建议" description="更轻松地坚持下去">
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/65 dark:bg-gray-900/40 px-4 py-4">
          <div class="flex items-center gap-3">
            <div class="h-10 w-10 rounded-2xl bg-emerald-50 dark:bg-emerald-900/25 text-emerald-700 dark:text-emerald-300 grid place-items-center">
              <el-icon :size="18"><CircleCheck /></el-icon>
            </div>
            <div class="min-w-0">
              <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">从小任务开始</div>
              <div class="mt-1 text-xs leading-relaxed text-gray-600 dark:text-gray-300">先把“能完成”放在第一位。</div>
            </div>
          </div>
        </div>
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/65 dark:bg-gray-900/40 px-4 py-4">
          <div class="flex items-center gap-3">
            <div class="h-10 w-10 rounded-2xl bg-amber-50 dark:bg-amber-900/25 text-amber-700 dark:text-amber-300 grid place-items-center">
              <el-icon :size="18"><Clock /></el-icon>
            </div>
            <div class="min-w-0">
              <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">用番茄做节奏</div>
              <div class="mt-1 text-xs leading-relaxed text-gray-600 dark:text-gray-300">20-25 分钟专注，效率更稳定。</div>
            </div>
          </div>
        </div>
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/65 dark:bg-gray-900/40 px-4 py-4">
          <div class="flex items-center gap-3">
            <div class="h-10 w-10 rounded-2xl bg-indigo-50 dark:bg-indigo-900/25 text-indigo-700 dark:text-indigo-300 grid place-items-center">
              <el-icon :size="18"><Present /></el-icon>
            </div>
            <div class="min-w-0">
              <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">奖励要具体</div>
              <div class="mt-1 text-xs leading-relaxed text-gray-600 dark:text-gray-300">心愿越明确，越有动力去完成。</div>
            </div>
          </div>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard title="核心功能">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/60 dark:bg-gray-900/35 px-4 py-4">
          <div class="flex items-center gap-3">
            <div class="h-10 w-10 rounded-2xl bg-sky-50 dark:bg-sky-900/25 text-sky-700 dark:text-sky-300 grid place-items-center">
              <el-icon :size="18"><List /></el-icon>
            </div>
            <div class="min-w-0">
              <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">任务管理</div>
              <div class="text-xs text-gray-600 dark:text-gray-300">重复、完成、统计、回收站</div>
            </div>
          </div>
        </div>
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/60 dark:bg-gray-900/35 px-4 py-4">
          <div class="flex items-center gap-3">
            <div class="h-10 w-10 rounded-2xl bg-amber-50 dark:bg-amber-900/25 text-amber-700 dark:text-amber-300 grid place-items-center">
              <el-icon :size="18"><Coin /></el-icon>
            </div>
            <div class="min-w-0">
              <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">积分激励</div>
              <div class="text-xs text-gray-600 dark:text-gray-300">赚金币、兑换心愿、历史变动</div>
            </div>
          </div>
        </div>
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/60 dark:bg-gray-900/35 px-4 py-4">
          <div class="flex items-center gap-3">
            <div class="h-10 w-10 rounded-2xl bg-emerald-50 dark:bg-emerald-900/25 text-emerald-700 dark:text-emerald-300 grid place-items-center">
              <el-icon :size="18"><TrendCharts /></el-icon>
            </div>
            <div class="min-w-0">
              <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">成长记录</div>
              <div class="text-xs text-gray-600 dark:text-gray-300">身高、体重、视力、脚长</div>
            </div>
          </div>
        </div>
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/60 dark:bg-gray-900/35 px-4 py-4">
          <div class="flex items-center gap-3">
            <div class="h-10 w-10 rounded-2xl bg-indigo-50 dark:bg-indigo-900/25 text-indigo-700 dark:text-indigo-300 grid place-items-center">
              <el-icon :size="18"><DataAnalysis /></el-icon>
            </div>
            <div class="min-w-0">
              <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">数据统计</div>
              <div class="text-xs text-gray-600 dark:text-gray-300">趋势与汇总，帮助复盘</div>
            </div>
          </div>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard title="下载地址" description="选择合适的下载方式">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <a href="https://pan.baidu.com/s/16zL_7-ZuYBxcZ3aVXXqWCw?pwd=izyj" target="_blank" class="block">
          <div class="group rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/55 dark:bg-gray-900/35 hover:bg-white/80 dark:hover:bg-gray-900/55 transition px-4 py-4">
            <div class="flex items-start justify-between gap-4">
              <div class="flex items-start gap-3">
                <div class="h-11 w-11 rounded-3xl bg-blue-50 dark:bg-blue-900/25 text-blue-700 dark:text-blue-300 grid place-items-center">
                  <el-icon :size="18"><Download /></el-icon>
                </div>
                <div class="min-w-0">
                  <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">百度网盘下载</div>
                  <div class="mt-1 text-xs text-gray-600 dark:text-gray-300">提取码：<span class="font-extrabold tabular-nums">izyj</span></div>
                </div>
              </div>
              <el-icon :size="18" class="text-gray-300 dark:text-gray-600 group-hover:text-gray-500 transition-colors"><ArrowRight /></el-icon>
            </div>
          </div>
        </a>

        <a href="https://www.123865.com/s/FQLuVv-c7ztv" target="_blank" class="block">
          <div class="group rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/55 dark:bg-gray-900/35 hover:bg-white/80 dark:hover:bg-gray-900/55 transition px-4 py-4">
            <div class="flex items-start justify-between gap-4">
              <div class="flex items-start gap-3">
                <div class="h-11 w-11 rounded-3xl bg-emerald-50 dark:bg-emerald-900/25 text-emerald-700 dark:text-emerald-300 grid place-items-center">
                  <el-icon :size="18"><Download /></el-icon>
                </div>
                <div class="min-w-0">
                  <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">123 盘下载</div>
                  <div class="mt-1 text-xs text-gray-600 dark:text-gray-300">高速下载推荐</div>
                </div>
              </div>
              <el-icon :size="18" class="text-gray-300 dark:text-gray-600 group-hover:text-gray-500 transition-colors"><ArrowRight /></el-icon>
            </div>
          </div>
        </a>
      </div>
    </SettingsCard>

    <SettingsCard title="常见问题" description="遇到问题先看这里">
      <div class="space-y-2">
        <details class="group rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/55 dark:bg-gray-900/35 px-4 py-3">
          <summary class="cursor-pointer list-none flex items-center justify-between gap-4">
            <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">金币是怎么计算的？</div>
            <el-icon :size="16" class="text-gray-300 dark:text-gray-600 group-open:text-gray-500 transition-colors"><ArrowRight /></el-icon>
          </summary>
          <div class="mt-2 text-sm leading-relaxed text-gray-600 dark:text-gray-300">
            完成任务会增加金币；兑换心愿会扣除金币。你也可以在关于/设置页面查看相关入口与提示。
          </div>
        </details>

        <details class="group rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/55 dark:bg-gray-900/35 px-4 py-3">
          <summary class="cursor-pointer list-none flex items-center justify-between gap-4">
            <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">为什么看不到底部导航？</div>
            <el-icon :size="16" class="text-gray-300 dark:text-gray-600 group-open:text-gray-500 transition-colors"><ArrowRight /></el-icon>
          </summary>
          <div class="mt-2 text-sm leading-relaxed text-gray-600 dark:text-gray-300">
            部分页面是“沉浸式页面”（例如设置、成长记录等），会隐藏底部导航以获得更大空间。
          </div>
        </details>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { QuestionFilled, Download, ArrowRight, Edit, Clock, Present, List, Coin, TrendCharts, DataAnalysis, CirclePlusFilled, Microphone, CircleCheck } from '@element-plus/icons-vue'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'
import appIcon from '@/assets/favicon/android-chrome-192x192.png'

const router = useRouter()
const appIconSrc = (appIcon as any) as string
</script>

<style scoped>
</style>
