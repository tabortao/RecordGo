<template>
  <div class="p-4 space-y-4 pb-16">
    <!-- 标题栏 -->
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer" style="color:#64748b" @click="router.back()"><ArrowLeft /></el-icon>
      <el-icon :size="18" style="color:#f59e0b"><Present /></el-icon>
      <h2 class="font-semibold">打赏支持</h2>
    </div>

    <el-card shadow="never" body-class="p-4">
      <div class="text-center space-y-4">
        <p class="text-gray-600 dark:text-gray-300 leading-relaxed">
          如果觉得这个项目对您有帮助，欢迎打赏支持！<br>
          您的支持是我持续更新的最大动力 ❤️
        </p>
        
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6 mt-6">
          <div class="flex flex-col items-center">
            <div class="w-48 h-48 bg-gray-100 rounded-lg overflow-hidden mb-2">
              <el-image 
                :src="wechatpay" 
                :preview-src-list="[wechatpay]"
                fit="cover"
                class="w-full h-full"
              >
                <template #error>
                  <div class="flex items-center justify-center w-full h-full text-gray-400">
                    <el-icon :size="24"><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
            </div>
            <div class="flex items-center gap-2 text-green-600 font-medium">
              <el-icon><ChatDotSquare /></el-icon>
              <span>微信支付</span>
            </div>
          </div>

          <div class="flex flex-col items-center">
            <div class="w-48 h-48 bg-gray-100 rounded-lg overflow-hidden mb-2">
              <el-image 
                :src="alipay" 
                :preview-src-list="[alipay]"
                fit="cover"
                class="w-full h-full"
              >
                <template #error>
                  <div class="flex items-center justify-center w-full h-full text-gray-400">
                    <el-icon :size="24"><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
            </div>
            <div class="flex items-center gap-2 text-blue-500 font-medium">
              <el-icon><Wallet /></el-icon>
              <span>支付宝</span>
            </div>
          </div>
        </div>

        <div class="text-xs text-gray-400 mt-6">
          打赏请备注您的昵称，感谢您的慷慨解囊~
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ArrowLeft, Present, Picture, ChatDotSquare, Wallet } from '@element-plus/icons-vue'

import wechatpay from '@/assets/images/wechatpay.jpg'
import alipay from '@/assets/images/alipay.jpg'

const router = useRouter()
</script>

<style scoped>
</style>
