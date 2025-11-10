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

    <!-- 分类列表（可编辑名称与颜色） -->
    <el-card shadow="never">
      <div class="space-y-2">
        <div v-for="c in categories" :key="c.name" class="flex items-center justify-between px-2 py-2 rounded border border-gray-200">
          <div class="flex items-center gap-3">
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
</template>

<script setup lang="ts">
// 中文注释：任务分类设置页面，联动 Pinia Store，实现增删改与颜色选择，并本地持久化
import { Notebook, ArrowLeft, Plus, Edit, Delete } from '@element-plus/icons-vue'
import router from '@/router'
import { computed, reactive, ref, watchEffect } from 'vue'
import { useTaskCategories } from '@/stores/categories'
function goBack() { router.back() }

const cats = useTaskCategories()
const categories = computed(() => cats.list())

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
</script>

<style scoped>
</style>