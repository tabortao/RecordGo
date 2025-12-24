<template>
  <div class="p-4 space-y-4 pb-16">
    <!-- 标题栏 -->
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer" style="color:#64748b" @click="router.back()"><ArrowLeft /></el-icon>
      <el-icon :size="18" style="color:#0ea5e9"><QuestionFilled /></el-icon>
      <h2 class="font-semibold">使用帮助</h2>
    </div>

    <!-- 项目简介 -->
    <el-card shadow="never" body-class="p-4">
      <template #header>
        <div class="font-semibold">项目简介</div>
      </template>
      <p class="text-gray-600 dark:text-gray-300 text-sm leading-relaxed">
        作业家是一款面向学生的任务与积分、成长记录、学习帮手的小应用。通过简洁的界面与激励机制，帮助学生培养良好的学习习惯与时间管理能力。通过完成任务获得积分，积分可用于兑换心愿，实现正向反馈循环。
      </p>
    </el-card>

    <!-- 项目特点 -->
    <el-card shadow="never" body-class="p-4">
      <template #header>
        <div class="font-semibold">项目特点</div>
      </template>
      <ul class="list-disc list-inside text-gray-600 dark:text-gray-300 text-sm space-y-1">
        <li>任务管理：清晰的任务列表，支持多种重复模式。</li>
        <li>番茄专注：内置番茄钟，提升专注力。</li>
        <li>积分激励：完成任务赚取积分，兑换心愿奖励。</li>
        <li>数据统计：多维度数据统计，记录成长点滴。</li>
        <li>多端支持：支持 Windows、Web 等多种平台。</li>
      </ul>
    </el-card>

    <!-- 下载地址 -->
    <el-card shadow="never" body-class="p-4">
      <template #header>
        <div class="font-semibold">下载地址</div>
      </template>
      <div class="space-y-3">
        <a href="https://pan.baidu.com/s/16zL_7-ZuYBxcZ3aVXXqWCw?pwd=izyj" target="_blank" class="block">
          <div class="flex items-center justify-between p-3 rounded-lg bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 transition">
            <div class="flex items-center gap-3">
              <el-icon :size="20" style="color:#3b82f6"><Download /></el-icon>
              <div>
                <div class="font-medium text-gray-800 dark:text-gray-100">百度网盘下载</div>
                <div class="text-xs text-gray-500">提取码: izyj</div>
              </div>
            </div>
            <el-icon><ArrowRight /></el-icon>
          </div>
        </a>
        
        <a href="https://www.123865.com/s/FQLuVv-c7ztv" target="_blank" class="block">
          <div class="flex items-center justify-between p-3 rounded-lg bg-gray-50 dark:bg-gray-700 hover:bg-gray-100 dark:hover:bg-gray-600 transition">
            <div class="flex items-center gap-3">
              <el-icon :size="20" style="color:#10b981"><Download /></el-icon>
              <div>
                <div class="font-medium text-gray-800 dark:text-gray-100">123盘下载</div>
                <div class="text-xs text-gray-500">高速下载推荐</div>
              </div>
            </div>
            <el-icon><ArrowRight /></el-icon>
          </div>
        </a>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ArrowLeft, QuestionFilled, ArrowRight, Download } from '@element-plus/icons-vue'

const router = useRouter()
</script>

<style scoped>
</style>
