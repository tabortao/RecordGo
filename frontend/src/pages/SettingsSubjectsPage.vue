<template>
  <div class="p-4 space-y-4">
    <!-- 顶部：返回 + 标题 -->
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer" style="color:#64748b" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" style="color:#2563eb"><Notebook /></el-icon>
      <h2 class="font-semibold">任务分类设置</h2>
    </div>

    <!-- 新增分类 -->
    <el-card shadow="never">
      <div class="flex items-center gap-3">
        <el-input v-model="newName" placeholder="输入分类名称，如：语文" style="max-width: 300px" />
        <el-color-picker v-model="newColor" />
        <el-button type="primary" @click="addCategory"><el-icon class="mr-1"><Plus /></el-icon>添加分类</el-button>
      </div>
      <div class="text-xs text-gray-500 mt-2">中文注释：支持添加、修改、删除分类，并设置颜色；分类与任务创建/编辑页面保持一致。</div>
    </el-card>

    <!-- 分类列表（可编辑名称与颜色；支持拖拽排序） -->
    <el-card shadow="never">
      <div class="space-y-2">
        <div
          v-for="(c, idx) in sortable"
          :key="c.name"
          class="flex items-center justify-between px-2 py-2 rounded border border-gray-200"
          draggable="true"
          @dragstart="onDragStart(idx)"
          @dragover="onDragOver"
          @drop="onDrop(idx)"
        >
          <div class="flex items-center gap-3">
            <!-- 中文注释：拖拽把手 -->
            <el-icon :size="16" class="cursor-grab text-gray-500"><Sort /></el-icon>
            <span class="inline-block w-4 h-4 rounded" :style="{ backgroundColor: c.color }"></span>
            <el-input v-model="editable[c.name].name" style="width: 160px" />
            <el-color-picker v-model="editable[c.name].color" />
          </div>
          <div class="flex items-center gap-2">
            <el-button size="small" @click="saveCategory(c.name)"><el-icon class="mr-1"><Edit /></el-icon>保存</el-button>
            <el-button size="small" type="danger" @click="removeCategory(c.name)"><el-icon class="mr-1"><Delete /></el-icon>删除</el-button>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 取消底部返回按钮；统一使用标题左侧返回图标 -->
  </div>
  <div class="mt-2">
    <el-button v-if="isAdmin" type="primary" @click="router.push('/admin')">用户管理</el-button>
  </div>
</template>

<script setup lang="ts">
// 中文注释：任务分类设置页面，联动 Pinia Store，实现增删改与颜色选择，并本地持久化
import { Notebook, ArrowLeft, Plus, Edit, Delete, Sort } from '@element-plus/icons-vue'
// 中文注释：显式引入 Color Picker 样式，避免在部分环境下样式资源加载中断
import 'element-plus/es/components/color-picker/style/css'
import router from '@/router'
import { useAuth } from '@/stores/auth'
import { computed, reactive, ref, watchEffect } from 'vue'
import { ElMessage } from 'element-plus'
import { useTaskCategories } from '@/stores/categories'
function goBack() { router.back() }
const auth = useAuth()
const isAdmin = computed(() => Number(auth.user?.id || 0) === 1)

const cats = useTaskCategories()
const categories = computed(() => cats.list())
// 中文注释：用于拖拽排序的本地可排序数组（显示顺序与 store 同步）
const sortable = ref<{ name: string; color: string }[]>([])
watchEffect(() => {
  sortable.value = categories.value.map(c => ({ name: c.name, color: c.color }))
})

// 中文注释：新增分类输入
const newName = ref('')
const newColor = ref('#64748b')
function addCategory() {
  const name = newName.value.trim()
  if (!name) return
  cats.add(name, newColor.value || '#64748b')
  newName.value = ''
  newColor.value = '#64748b'
}

// 中文注释：编辑缓冲区（按原名称为 key），避免直接改动 store 中数据
const editable = reactive<Record<string, { name: string; color: string }>>({})
watchEffect(() => {
  const next: Record<string, { name: string; color: string }> = {}
  for (const c of categories.value) next[c.name] = { name: c.name, color: c.color }
  Object.assign(editable, next)
})
function saveCategory(oldName: string) {
  const patch = editable[oldName]
  if (!patch) return
  cats.update(oldName, { name: (patch.name || '').trim(), color: patch.color })
}
function removeCategory(name: string) { cats.remove(name) }

// ===== 拖拽排序交互 =====
const draggingIndex = ref<number | null>(null)
function onDragStart(idx: number) { draggingIndex.value = idx }
function onDragOver(e: DragEvent) { e.preventDefault() }
function onDrop(targetIdx: number) {
  const from = draggingIndex.value
  draggingIndex.value = null
  if (from === null || from === targetIdx) return
  const arr = [...sortable.value]
  const [moved] = arr.splice(from, 1)
  arr.splice(targetIdx, 0, moved)
  sortable.value = arr
  // 中文注释：应用到 store 顺序（1..N），并提示
  for (let i = 0; i < arr.length; i++) {
    cats.setOrder(arr[i].name, i + 1)
  }
  ElMessage.success('已更新分类顺序')
}
</script>

<style scoped>
</style>
