<template>
  <SettingsShell title="打赏支持" :icon="Present" tone="amber" container-class="max-w-3xl">
    <SettingsCard>
      <div class="flex flex-col items-center justify-center px-4 py-8 text-center">
        <div class="relative mb-6">
          <div class="absolute -inset-4 rounded-full bg-amber-200/50 dark:bg-amber-500/15 blur-3xl animate-pulse" />
          <div class="relative h-24 w-24 rounded-[32px] bg-gradient-to-br from-amber-100 to-orange-50 dark:from-amber-900/30 dark:to-orange-900/10 grid place-items-center shadow-xl shadow-amber-500/20 ring-1 ring-white/50 dark:ring-white/10">
            <el-icon :size="48" class="text-amber-500"><Present /></el-icon>
          </div>
        </div>
        
        <h2 class="text-2xl font-black tracking-tight text-gray-900 dark:text-white mb-3">支持作者持续更新</h2>
        <p class="text-sm leading-relaxed text-gray-600 dark:text-gray-300 max-w-lg mx-auto">
          如果觉得《作业家》对你有帮助，欢迎打赏支持。你的每一份心意，都是我坚持开发与优化的动力！
        </p>
      </div>
    </SettingsCard>

    <SettingsCard title="打赏方式" description="长按保存图片或直接扫码">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
        <!-- WeChat Pay -->
        <div class="group relative overflow-hidden rounded-[32px] bg-[#2DA06F] p-1">
          <div class="absolute inset-0 bg-[url('https://api.iconify.design/ri:wechat-fill.svg?color=%23ffffff&opacity=0.1')] bg-repeat opacity-20"></div>
          <div class="relative h-full rounded-[28px] bg-white dark:bg-gray-900 p-6 flex flex-col items-center text-center transition-transform duration-300 group-hover:-translate-y-1">
            <div class="mb-4 h-12 w-12 rounded-2xl bg-[#2DA06F]/10 flex items-center justify-center text-[#2DA06F]">
              <el-icon :size="28"><ChatDotRound /></el-icon>
            </div>
            <div class="text-lg font-black text-gray-900 dark:text-white mb-1">微信支付</div>
            <div class="text-xs font-medium text-gray-500 dark:text-gray-400 mb-6">WeChat Pay</div>
            
            <div class="relative w-48 h-48 rounded-2xl overflow-hidden bg-gray-100 dark:bg-gray-800 p-2 shadow-inner group-hover:shadow-lg transition-shadow">
              <img src="@/assets/images/wechatpay.jpg" alt="微信收款码" class="w-full h-full object-cover rounded-xl" />
            </div>
            <div class="mt-4 text-xs text-gray-400">长按保存图片扫码</div>
          </div>
        </div>

        <!-- Alipay -->
        <div class="group relative overflow-hidden rounded-[32px] bg-[#1678FF] p-1">
          <div class="absolute inset-0 bg-[url('https://api.iconify.design/ri:alipay-fill.svg?color=%23ffffff&opacity=0.1')] bg-repeat opacity-20"></div>
          <div class="relative h-full rounded-[28px] bg-white dark:bg-gray-900 p-6 flex flex-col items-center text-center transition-transform duration-300 group-hover:-translate-y-1">
            <div class="mb-4 h-12 w-12 rounded-2xl bg-[#1678FF]/10 flex items-center justify-center text-[#1678FF]">
              <el-icon :size="28"><Money /></el-icon>
            </div>
            <div class="text-lg font-black text-gray-900 dark:text-white mb-1">支付宝</div>
            <div class="text-xs font-medium text-gray-500 dark:text-gray-400 mb-6">Alipay</div>
            
            <div class="relative w-48 h-48 rounded-2xl overflow-hidden bg-gray-100 dark:bg-gray-800 p-2 shadow-inner group-hover:shadow-lg transition-shadow">
              <img src="@/assets/images/alipay.jpg" alt="支付宝收款码" class="w-full h-full object-cover rounded-xl" />
            </div>
            <div class="mt-4 text-xs text-gray-400">长按保存图片扫码</div>
          </div>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard title="感谢名单" description="感谢每一位支持的朋友">
      <div class="rounded-[24px] border border-gray-100 dark:border-gray-800 bg-gray-50/50 dark:bg-gray-900/30 p-8 text-center">
        <div class="inline-flex h-16 w-16 rounded-full bg-white dark:bg-gray-800 shadow-sm items-center justify-center mb-4">
          <el-icon :size="24" class="text-amber-400"><TrophyBase /></el-icon>
        </div>
        <p class="text-sm font-bold text-gray-900 dark:text-white">虚位以待</p>
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">期待你的名字出现在这里~</p>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { Present, ChatDotRound, Money, TrophyBase } from '@element-plus/icons-vue'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'
</script>

<style scoped>
/* 增加收款码图片占位样式 */
img {
  font-size: 12px;
  color: #9ca3af;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f3f4f6;
}
.dark img {
  background-color: #1f2937;
  color: #4b5563;
}
</style>
