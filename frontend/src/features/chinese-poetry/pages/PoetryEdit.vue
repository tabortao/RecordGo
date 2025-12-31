<template>
  <div class="min-h-screen bg-[#FDF6F8] dark:bg-gray-900 pb-20 font-sans">
    <!-- Header -->
    <div class="sticky top-0 z-20 bg-white/80 dark:bg-gray-800/80 backdrop-blur-md border-b border-pink-100 dark:border-gray-700 shadow-sm">
      <div class="flex items-center gap-3 p-4">
        <el-icon :size="22" class="cursor-pointer text-gray-600 dark:text-gray-300 hover:text-pink-500 transition" @click="router.back()"><ArrowLeft /></el-icon>
        <h1 class="text-xl font-bold text-gray-800 dark:text-white">{{ isEdit ? '编辑古诗' : '新增古诗' }}</h1>
      </div>
    </div>

    <div class="p-6 max-w-2xl mx-auto space-y-6">
      <div class="bg-white dark:bg-gray-800 rounded-3xl p-6 shadow-xl shadow-pink-100/50 dark:shadow-none border border-pink-50 dark:border-gray-700">
        <el-form :model="form" label-position="top" size="large" class="space-y-4">
          <el-form-item label="诗词标题">
            <el-input 
              v-model="form.title_cns" 
              placeholder="如：静夜思" 
              class="!text-lg"
            />
          </el-form-item>
          
          <div class="grid grid-cols-2 gap-6">
             <el-form-item label="作者">
               <el-input v-model="form.author_cns" placeholder="如：李白" />
             </el-form-item>
             <el-form-item label="朝代">
               <el-input v-model="form.dynasty_cns" placeholder="如：唐代" />
             </el-form-item>
          </div>
          
          <el-form-item label="适用学段">
             <el-select v-model="form.tags" multiple placeholder="选择学段" class="w-full">
                <el-option label="小学" value="primary" />
                <el-option label="初中" value="middle" />
                <el-option label="高中" value="high" />
             </el-select>
          </el-form-item>

          <el-form-item label="古诗内容 (每行一句)">
            <el-input 
              v-model="contentStr" 
              type="textarea" 
              :rows="8" 
              placeholder="床前明月光，&#10;疑是地上霜。" 
              class="font-serif text-lg tracking-wider"
            />
          </el-form-item>
          
          <el-form-item label="音频链接 (可选)">
            <el-input v-model="form.audio_url" placeholder="https://..." prefix-icon="Microphone" />
          </el-form-item>
          
          <div class="pt-4">
             <el-button 
               type="primary" 
               class="w-full !h-12 !text-lg !rounded-full shadow-lg shadow-pink-200 hover:shadow-pink-300 transition-all transform hover:-translate-y-0.5 bg-gradient-to-r from-pink-400 to-rose-400 border-none" 
               @click="save"
             >
               保存古诗
             </el-button>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowLeft } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { usePoetryStore } from '../stores/poetryStore'
import type { Poem } from '../types'

const router = useRouter()
const route = useRoute()
const store = usePoetryStore()

const isEdit = ref(false)
const contentStr = ref('')
const form = ref<Partial<Poem>>({
    title_cns: '',
    author_cns: '',
    dynasty_cns: '',
    tags: [],
    paragraphs_cns: [],
    audio_url: ''
})

onMounted(() => {
    store.init()
    const id = route.params.id
    if (id) {
        const poem = store.getPoemById(Number(id))
        if (poem && poem.isCustom) {
            isEdit.value = true
            form.value = { ...poem }
            contentStr.value = poem.paragraphs_cns.join('\n')
        }
    }
})

const save = () => {
    if (!form.value.title_cns || !form.value.author_cns || !contentStr.value) {
        ElMessage.warning('请填写必要信息')
        return
    }
    
    form.value.paragraphs_cns = contentStr.value.split('\n').filter(l => l.trim())
    
    if (isEdit.value) {
        store.updateCustomPoem(form.value.id!, form.value)
        ElMessage.success('更新成功')
    } else {
        store.addCustomPoem(form.value as Poem)
        ElMessage.success('添加成功')
    }
    router.back()
}
</script>
