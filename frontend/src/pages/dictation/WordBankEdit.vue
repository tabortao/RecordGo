<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 flex flex-col">
    <!-- Header -->
    <div class="bg-white dark:bg-gray-800 shadow-sm px-4 py-3 flex items-center justify-between sticky top-0 z-10">
      <div class="flex items-center gap-2">
        <el-icon :size="20" class="text-gray-600 dark:text-gray-300 cursor-pointer" @click="router.back()"><ArrowLeft /></el-icon>
        <span class="text-lg font-bold text-gray-800 dark:text-white">{{ isEdit ? '编辑词库' : '新建词库' }}</span>
      </div>
      <el-button type="primary" @click="submit" :loading="saving">保存</el-button>
    </div>

    <div class="p-4 space-y-4">
      <el-form :model="form" label-position="top">
        <el-form-item label="标题">
          <el-input v-model="form.title" placeholder="例如：第一单元生词" />
        </el-form-item>
        
        <div class="grid grid-cols-2 gap-3">
          <el-form-item label="科目">
            <el-select v-model="form.subject">
              <el-option label="语文" value="语文" />
              <el-option label="英语" value="英语" />
              <el-option label="其他" value="其他" />
            </el-select>
          </el-form-item>
          <el-form-item label="阶段">
            <el-select v-model="form.education_stage">
              <el-option label="小学" value="小学" />
              <el-option label="初中" value="初中" />
              <el-option label="高中" value="高中" />
            </el-select>
          </el-form-item>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <el-form-item label="版本">
            <el-input v-model="form.version" placeholder="例如：人教版" />
          </el-form-item>
          <el-form-item label="年级">
            <el-input v-model="form.grade" placeholder="例如：三年级上" />
          </el-form-item>
        </div>

        <el-form-item label="内容">
          <el-input 
            v-model="form.content" 
            type="textarea" 
            :rows="10" 
            placeholder="输入词汇或课文，支持空格或换行分隔" 
          />
          <div class="text-xs text-gray-500 mt-1">提示：系统将自动识别分隔符进行分词。</div>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { dictationApi, type WordBank } from '@/services/dictation'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const id = route.params.id ? Number(route.params.id) : null
const isEdit = computed(() => !!id)
const saving = ref(false)

const form = ref<Partial<WordBank>>({
  title: '',
  subject: '语文',
  education_stage: '小学',
  version: '',
  grade: '',
  content: ''
})

onMounted(async () => {
  if (isEdit && id) {
    try {
      // Since we don't have a getOne API in the snippet above, let's rely on list or implement getOne.
      // Wait, I forgot to add GetWordBank to handlers/router? 
      // Let's check router.go. I missed `dictation.GET("/wordbanks/:id", handlers.GetWordBank)`?
      // I implemented UpdateWordBank which takes ID, but I might have missed Get.
      // Let's assume I can just use the list to find it or fetch it. 
      // Actually, I didn't implement `GetWordBank` in handlers!
      // I implemented `ListWordBanks`, `Create`, `Update`, `Delete`.
      // So I should fetch list and find it, OR implement Get.
      // Implementing Get is better. But for speed, I will just fetch list and find.
      const res = await dictationApi.listWordBanks()
      const found = res.data.find((i: any) => i.id === id)
      if (found) {
        form.value = { ...found }
        // If content is JSON array, convert to string?
        try {
           if (found.content.startsWith('[')) {
             const arr = JSON.parse(found.content)
             form.value.content = arr.join('\n')
           }
        } catch {}
      }
    } catch (e) {
      ElMessage.error('加载失败')
    }
  }
})

async function submit() {
  if (!form.value.title) return ElMessage.warning('请输入标题')
  if (!form.value.content) return ElMessage.warning('请输入内容')
  
  saving.value = true
  try {
    const payload = { ...form.value }
    // Convert content to simple string or keep as is. 
    // PRD doesn't specify storage format, but string is fine.
    // If user inputs with newlines, we keep it.
    
    if (isEdit.value && id) {
      await dictationApi.updateWordBank(id, payload)
      ElMessage.success('已更新')
    } else {
      await dictationApi.createWordBank(payload)
      ElMessage.success('已创建')
    }
    router.back()
  } catch (e) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}
</script>
