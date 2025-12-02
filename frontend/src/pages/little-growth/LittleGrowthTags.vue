<template>
  <div class="min-h-screen bg-white flex flex-col">
    <!-- Header -->
    <div class="px-4 py-3 border-b border-gray-100 flex items-center gap-3 sticky top-0 bg-white z-20">
      <div 
        class="w-8 h-8 rounded-full bg-gray-50 flex items-center justify-center cursor-pointer hover:bg-gray-100 transition-colors text-gray-600"
        @click="router.back()"
      >
        <el-icon><ArrowLeft /></el-icon>
      </div>
      <h1 class="font-bold text-lg text-gray-800">管理标签</h1>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto p-4">
      <div class="max-w-2xl mx-auto">
        <div v-if="store.tags.length === 0" class="text-center py-20 text-gray-400">
          <el-icon :size="48" class="mb-4 opacity-50"><CollectionTag /></el-icon>
          <p>还没有标签，点击下方按钮创建</p>
        </div>

        <div class="space-y-2">
          <div 
            v-for="tag in store.tags" 
            :key="tag.id"
            class="flex items-center justify-between p-4 bg-gray-50 rounded-xl border border-gray-100 group hover:border-purple-200 transition-colors"
          >
            <div class="flex items-center gap-3">
              <span class="font-medium text-gray-700 text-lg"># {{ tag.name }}</span>
              <span class="text-xs bg-white px-2 py-0.5 rounded-full text-gray-400 border border-gray-100">{{ tag.count || 0 }} 记录</span>
            </div>
            <div class="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
              <button 
                class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-full transition-colors"
                @click="openEdit(tag)"
              >
                <el-icon><Edit /></el-icon>
              </button>
              <button 
                class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-full transition-colors"
                @click="handleDelete(tag)"
              >
                <el-icon><Delete /></el-icon>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- FAB Add -->
    <div 
      class="fixed right-6 bottom-8 z-20 w-14 h-14 bg-purple-600 rounded-full shadow-lg shadow-purple-200 flex items-center justify-center text-white cursor-pointer hover:bg-purple-700 transition-colors active:scale-95"
      @click="openCreate"
    >
      <el-icon :size="24"><Plus /></el-icon>
    </div>

    <!-- Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '修改标签' : '新建标签'"
      width="90%"
      max-width="400px"
      class="rounded-2xl"
      center
    >
      <div class="py-4">
        <el-input 
          v-model="form.name" 
          placeholder="输入标签名称" 
          size="large"
          @keyup.enter="submit"
        />
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submit" :loading="submitting">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Plus, Edit, Delete, CollectionTag } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useLittleGrowthStore, type Tag } from '@/stores/littleGrowth'

const router = useRouter()
const store = useLittleGrowthStore()

const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const currentId = ref<string>('')
const form = ref({ name: '' })

onMounted(() => {
  store.fetchTags()
})

const openCreate = () => {
  isEdit.value = false
  form.value.name = ''
  dialogVisible.value = true
}

const openEdit = (tag: Tag) => {
  isEdit.value = true
  currentId.value = tag.id
  form.value.name = tag.name
  dialogVisible.value = true
}

const submit = async () => {
  if (!form.value.name.trim()) return
  
  submitting.value = true
  try {
    if (isEdit.value) {
      await store.updateTag(currentId.value, form.value.name)
      ElMessage.success('修改成功')
    } else {
      await store.createTag(form.value.name)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
  } catch (e) {
    ElMessage.error('操作失败')
  } finally {
    submitting.value = false
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
