<template>
  <!-- 中文注释：系统设置页，左侧分类列表，右侧对应设置内容；样式采用 tailwind/shadcn 风格 -->
  <div class="p-4">
    <div class="text-xl font-semibold mb-3 flex items-center gap-2">
      <el-icon :size="18" style="color:#3b82f6"><Setting /></el-icon>
      <span>设置</span>
    </div>
    <!-- 中文注释：按钮图标大小统一为 18，与“编辑个人信息”一致；点击弹出对应设置对话框 -->
    <div class="rounded-lg border bg-white">
      <div class="px-2 py-2 space-y-1">
        <button v-for="item in items" :key="item.key" class="w-full flex items-center gap-3 px-3 py-2 rounded-md hover:bg-gray-50 transition text-left" @click="openDialog(item.key)">
          <el-icon :size="18" :style="{ color: item.fg }"><component :is="item.icon" /></el-icon>
          <span class="font-medium">{{ item.label }}</span>
        </button>
      </div>
    </div>

    <!-- 番茄钟设置对话框 -->
    <el-dialog v-model="showTomato" title="番茄钟设置" width="520px">
      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <div>
            <div class="font-medium">固定番茄钟页面</div>
            <div class="text-xs text-gray-500 mt-1">勾选后，番茄钟将保持全屏显示，不会缩小为悬浮球</div>
          </div>
          <el-switch v-model="fixed" />
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end"><el-button type="primary" @click="showTomato = false">关闭</el-button></div>
      </template>
    </el-dialog>

    <!-- 金币设置对话框 -->
    <el-dialog v-model="showCoins" title="金币设置" width="520px">
      <div class="space-y-4">
        <div class="text-sm">当前总金币：<span class="font-semibold text-amber-600">{{ coins }}</span></div>
        <div class="space-y-2">
          <label class="text-sm text-gray-600">设置新金币数量</label>
          <el-input-number v-model="newCoins" :min="0" :max="999999" />
        </div>
        <div class="space-y-2">
          <label class="text-sm text-gray-600">修改原因</label>
          <el-input v-model="reason" type="textarea" placeholder="请输入修改原因（例如：期初校准/异常修正）" />
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="showCoins = false">取消</el-button>
          <el-button type="primary" @click="saveCoins">保存</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 其他设置占位对话框：任务、朗读、学科、关于 -->
    <el-dialog v-model="showTasks" title="任务设置" width="520px"><div class="text-gray-500 text-sm">后续将提供任务默认分类、默认计划时长等设置。</div></el-dialog>
    <el-dialog v-model="showReading" title="朗读设置" width="520px"><div class="text-gray-500 text-sm">后续将提供朗读速度、发音人等设置。</div></el-dialog>
    <el-dialog v-model="showSubjects" title="任务分类设置" width="520px"><div class="text-gray-500 text-sm">后续将提供常用分类与色彩标签设置。</div></el-dialog>
    <el-dialog v-model="showAbout" title="关于" width="520px"><div class="text-gray-500 text-sm">RecordGo 任务与积分助手。前端 Vue3 + Vite，后端 Gin + GORM。</div></el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watchEffect } from 'vue'
import { Setting, Timer, List, Microphone, Notebook, Coin, InfoFilled } from '@element-plus/icons-vue'
import { useAppState } from '@/stores/appState'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'

// 中文注释：类型约束，修复模板中“string 不能分配到联合类型”的报错
type SettingsKey = 'tomato' | 'tasks' | 'reading' | 'subjects' | 'coins' | 'about'
const items: Array<{ key: SettingsKey; label: string; icon: any; fg: string }> = [
  { key: 'tomato', label: '番茄钟设置', icon: Timer, fg: '#ef4444' },
  { key: 'tasks', label: '任务设置', icon: List, fg: '#10b981' },
  { key: 'reading', label: '朗读设置', icon: Microphone, fg: '#7c3aed' },
  { key: 'subjects', label: '学科设置', icon: Notebook, fg: '#2563eb' },
  { key: 'coins', label: '金币设置', icon: Coin, fg: '#f59e0b' },
  { key: 'about', label: '关于', icon: InfoFilled, fg: '#0ea5e9' }
]
const active = ref<SettingsKey>('tomato')

// 中文注释：支持通过路由参数 tab 初始化当前设置分组（来自“我的”页入口）
const route = useRoute()
watchEffect(() => {
  const tab = String(route.query.tab || '')
  const allowed: SettingsKey[] = ['tomato', 'tasks', 'reading', 'subjects', 'coins', 'about']
  if (allowed.includes(tab as SettingsKey)) {
    active.value = tab as SettingsKey
  }
})

// 中文注释：固定番茄钟设置与全局状态联动
const store = useAppState()
const fixed = computed({
  get: () => store.tomato.fixedTomatoPage,
  set: (v: boolean) => store.updateTomato({ fixedTomatoPage: v })
})

// 中文注释：金币设置对话框字段与保存逻辑
const coins = computed(() => store.coins)
const newCoins = ref<number | null>(null)
const reason = ref('')

// 中文注释：对话框显示状态
const showTomato = ref(false)
const showCoins = ref(false)
const showTasks = ref(false)
const showReading = ref(false)
const showSubjects = ref(false)
const showAbout = ref(false)

function openDialog(k: SettingsKey) {
  active.value = k
  if (k === 'tomato') showTomato.value = true
  else if (k === 'coins') { newCoins.value = coins.value; reason.value = ''; showCoins.value = true }
  else if (k === 'tasks') showTasks.value = true
  else if (k === 'reading') showReading.value = true
  else if (k === 'subjects') showSubjects.value = true
  else if (k === 'about') showAbout.value = true
}

function saveCoins() {
  const v = Number(newCoins.value ?? coins.value)
  if (isNaN(v) || v < 0) {
    ElMessage.error('请输入有效的金币数量（≥ 0）')
    return
  }
  // 中文注释：保存到全局状态，任务页“总金币”和心愿页“可用金币”都读取该值
  store.setCoins(v)
  ElMessage.success('金币已更新')
  showCoins.value = false
}
</script>

<style scoped>
</style>