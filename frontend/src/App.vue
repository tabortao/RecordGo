<template>
  <el-config-provider :locale="zhCn">
    <!-- 中文注释：顶层布局，包含页面内容与底部导航 -->
    <div class="min-h-screen flex flex-col bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
      <!-- 中文注释：为页面内容增加底部留白，避免被固定底部导航遮挡 -->
      <router-view class="flex-1 pb-20" />
      <!-- 中文注释：登录/注册页面隐藏底部导航 -->
      <BottomNav v-if="!noNav" />
    </div>
    <!-- 中文注释：非固定模式下，返回后显示悬浮番茄钟继续计时 -->
    <FloatingTomato v-if="shouldShowFloating" />
    
    <!-- 全局页面加载动画与进度条 -->
    <PageLoader />
    <el-backtop />
  </el-config-provider>
</template>

<script setup lang="ts">
// 中文注释：注册底部导航组件与根据路由控制显示
import BottomNav from '@/components/BottomNav.vue'
import FloatingTomato from '@/components/FloatingTomato.vue'
import PageLoader from '@/components/PageLoader.vue'
import { useRoute } from 'vue-router'
import { computed } from 'vue'
import { useAppState } from '@/stores/appState'
import { ElConfigProvider } from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

const route = useRoute()
const noNav = computed(() => route.meta?.noNav === true)
const store = useAppState()
const shouldShowFloating = computed(() => store.tomato.showFloating && store.tomato.running)
</script>

<style scoped>
/* 中文注释：可根据品牌色进行优化 */
</style>
