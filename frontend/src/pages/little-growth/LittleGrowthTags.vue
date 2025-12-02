<template>
  <div class="min-h-screen bg-gray-50 p-4">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center gap-2">
        <el-button circle :icon="ArrowLeft" @click="router.back()" />
        <h1 class="text-xl font-bold text-gray-800">标签管理</h1>
      </div>
      <el-button type="primary" @click="openDialog()">新建标签</el-button>
    </div>

    <!-- Tag List -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 max-w-2xl mx-auto">
      <div 
        v-for="tag in store.tags" 
        :key="tag.id"
        class="bg-white p-4 rounded-xl shadow-sm flex items-center justify-between group"
      >
        <div class="flex items-center gap-3">
          <div class="w-4 h-4 rounded-full" :style="{ backgroundColor: tag.color || '#A78BFA' }"></div>
          <span class="font-medium text-gray-700">{{ tag.name }}</span>
        </div>
        <div class="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
          <el-button link type="primary" :icon="Edit" @click="openDialog(tag)"></el-button>
          <el-button link type="danger" :icon="Delete" @click="handleDelete(tag)"></el-button>
        </div>
      </div>
    </div>

    <!-- Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑标签' : '新建标签'"
      width="90%"
      class="max-w-md rounded-2xl"
      center
    >
      <div class="flex flex-col gap-4">
        <el-input v-model="form.name" placeholder="请输入标签名称" />
        
        <div class="flex flex-col gap-2">
          <span class="text-sm text-gray-500">选择颜色</span>
          <div class="flex flex-wrap gap-2">
             <div 
               v-for="c in colors" 
               :key="c"
               class="w-8 h-8 rounded-full cursor-pointer border-2 transition-all"
               :class="form.color === c ? 'border-gray-800 scale-110' : 'border-transparent hover:scale-105'"
               :style="{ backgroundColor: c }"
               @click="form.color = c"
             ></div>
          </div>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer flex justify-center gap-4">
          <el-button @click="dialogVisible = false" class="!rounded-full px-6">取消</el-button>
          <el-button type="primary" @click="handleSubmit" class="!rounded-full px-6">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Edit, Delete } from '@element-plus/icons-vue'
import { useLittleGrowthStore, type Tag } from '@/stores/littleGrowth'
import { ElMessageBox, ElMessage } from 'element-plus'

const router = useRouter()
const store = useLittleGrowthStore()

const dialogVisible = ref(false)
const isEdit = ref(false)
const form = ref({ id: '', name: '', color: '' })

const colors = [
  "#FFB6C1", "#FF69B4", "#FFD700", "#FFA07A", "#90EE90", 
  "#20B2AA", "#87CEFA", "#9370DB", "#FF6347", "#40E0D0",
  "#EE82EE", "#F0E68C", "#E6E6FA", "#DDA0DD", "#B0C4DE"
]

const openDialog = (tag?: Tag) => {
  if (tag) {
    isEdit.value = true
    form.value = { id: tag.id, name: tag.name, color: tag.color || colors[Math.floor(Math.random() * colors.length)] }
  } else {
    isEdit.value = false
    form.value = { id: '', name: '', color: colors[Math.floor(Math.random() * colors.length)] }
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!form.value.name) return
  try {
    if (isEdit.value) {
      await store.updateTag(form.value.id, form.value.name, form.value.color)
    } else {
      await store.createTag(form.value.name, form.value.color)
    }
    dialogVisible.value = false
    ElMessage.success(isEdit.value ? '更新成功' : '创建成功')
  } catch (e) {
    // error handled by interceptor
  }
}

const handleDelete = async (tag: Tag) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除标签“${tag.name}”吗？删除后，相关记录将变为无标签状态，但记录不会被删除。`,
      '删除标签',
      {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning',
        center: true
      }
    )
    await store.deleteTag(tag.id)
    ElMessage.success('删除成功')
  } catch {}
}
</script>
