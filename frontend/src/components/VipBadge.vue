<template>
  <button
    type="button"
    class="inline-flex items-center gap-1.5 px-3 py-0.5 rounded-full text-xs font-bold shadow-sm transition-all duration-300 hover:scale-105 select-none group relative overflow-hidden"
    :class="[badgeClasses, clickable ? 'cursor-pointer' : 'cursor-default']"
    @click="onClick"
  >
    
    <!-- Shine effect -->
    <div class="absolute inset-0 bg-white/20 translate-x-[-100%] group-hover:translate-x-[100%] transition-transform duration-700 skew-x-12"></div>

    <!-- Icon -->
    <el-icon :size="14" class="relative z-10 transition-transform group-hover:rotate-12">
      <component :is="iconComponent" />
    </el-icon>

    <!-- Text -->
    <span class="relative z-10 tracking-wide font-sans">
      {{ labelText }}
    </span>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Trophy, Medal, Lock } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { AuthUser } from '@/stores/auth'
import { ElMessage, ElMessageBox } from 'element-plus'

const props = defineProps<{
  user: AuthUser | null | undefined
}>()

const vipType = computed(() => {
  if (!props.user) return 'none'
  if (props.user.is_lifetime_vip) return 'lifetime'
  
  if (props.user.is_vip) {
    // Client-side expiration check for immediate feedback
    if (props.user.vip_expire_time) {
       const expire = new Date(props.user.vip_expire_time).getTime()
       if (expire < Date.now()) return 'none'
    }
    return 'regular'
  }
  return 'none'
})

const badgeClasses = computed(() => {
  if (vipType.value === 'lifetime') {
    // Gold/Amber theme for Lifetime
    return 'bg-gradient-to-r from-amber-200 via-yellow-400 to-amber-500 text-amber-900 border border-amber-300/50 shadow-lg shadow-amber-500/30'
  }
  if (vipType.value === 'regular') {
    return 'bg-gradient-to-r from-indigo-400 via-purple-500 to-pink-500 text-white border border-indigo-300/30 shadow-lg shadow-indigo-500/30'
  }
  return 'bg-white/70 dark:bg-gray-900/60 text-gray-600 dark:text-gray-300 border border-gray-200/70 dark:border-gray-700/60 shadow-sm'
})

const iconComponent = computed(() => {
  if (vipType.value === 'lifetime') return Trophy
  if (vipType.value === 'regular') return Medal
  return Lock
})

const labelText = computed(() => {
  if (vipType.value === 'lifetime') return '终身VIP'
  if (vipType.value === 'regular' && props.user?.vip_expire_time) {
    return `VIP(${dayjs(props.user.vip_expire_time).format('YYYY-MM-DD')})`
  }
  if (vipType.value === 'regular') return 'VIP'
  return '普通用户'
})

const clickable = computed(() => vipType.value === 'none')

async function copyText(text: string) {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制微信号')
  } catch {
    ElMessage.info(`微信号：${text}`)
  }
}

function onClick() {
  if (vipType.value !== 'none') return
  ElMessageBox.confirm(
    '添加微信 “tabor2024” 获取 VIP。\n备注：任务家',
    '获取 VIP',
    { confirmButtonText: '复制微信号', cancelButtonText: '关闭', type: 'info' }
  ).then(() => copyText('tabor2024')).catch(() => {})
}
</script>
