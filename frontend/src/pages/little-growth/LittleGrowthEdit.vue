<template>
  <div class="min-h-screen bg-white flex flex-col">
    <!-- Header -->
    <div class="px-4 py-3 border-b border-gray-100 flex items-center justify-between sticky top-0 bg-white z-20">
      <div class="flex items-center gap-2" @click="router.back()">
        <el-icon class="text-gray-600 text-lg"><ArrowLeft /></el-icon>
        <span class="text-gray-600">取消</span>
      </div>
      <h1 class="font-bold text-lg">{{ isEdit ? '编辑记录' : '新记录' }}</h1>
      <button 
        class="bg-purple-600 text-white px-4 py-1.5 rounded-full text-sm font-medium hover:bg-purple-700 transition-colors disabled:opacity-50"
        :disabled="!form.content && form.images.length === 0"
        @click="save"
      >
        发布
      </button>
    </div>

    <div class="flex-1 overflow-y-auto p-4">
      <div class="max-w-2xl mx-auto space-y-6">
        <!-- Date -->
        <div class="flex items-center gap-2">
          <el-date-picker
            v-model="form.date"
            type="date"
            placeholder="选择日期"
            format="YYYY年MM月DD日"
            value-format="YYYY-MM-DD"
            :clearable="false"
            class="!w-40"
          />
        </div>

        <!-- Content -->
        <div class="relative">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="6"
            placeholder="记录这一刻的美好... 输入 # 可快速添加标签"
            class="text-lg"
            resize="none"
            @input="handleInput"
          />
          
          <!-- Smart Tag Suggestions -->
          <div v-if="showTagSuggestions" class="absolute top-full left-0 z-10 bg-white shadow-lg border border-gray-100 rounded-lg mt-1 w-64 max-h-48 overflow-y-auto">
            <div 
              v-for="tag in suggestedTags" 
              :key="tag.id"
              class="px-3 py-2 hover:bg-purple-50 cursor-pointer text-sm text-gray-700"
              @click="selectSuggestedTag(tag)"
            >
              # {{ tag.name }}
            </div>
            <div v-if="suggestedTags.length === 0" class="px-3 py-2 text-gray-400 text-xs">
              输入空格以创建新标签
            </div>
          </div>
        </div>

        <!-- Images -->
        <div>
          <div class="grid grid-cols-3 gap-2 mb-2">
            <div 
              v-for="(img, index) in form.images" 
              :key="index"
              class="relative aspect-square rounded-xl overflow-hidden group bg-gray-100"
            >
              <img :src="img" class="w-full h-full object-cover" />
              <div class="absolute top-1 right-1 p-1 bg-black/50 rounded-full text-white opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer" @click="removeImage(index)">
                <el-icon><Close /></el-icon>
              </div>
            </div>
            
            <div 
              v-if="form.images.length < 9"
              class="aspect-square rounded-xl border-2 border-dashed border-gray-200 flex flex-col items-center justify-center text-gray-400 hover:border-purple-400 hover:text-purple-500 transition-colors cursor-pointer bg-gray-50"
              @click="triggerUpload"
            >
              <el-icon :size="24"><Camera /></el-icon>
              <span class="text-xs mt-1">添加照片</span>
            </div>
          </div>
          <input 
            type="file" 
            ref="fileInput" 
            class="hidden" 
            accept="image/*" 
            multiple 
            @change="handleFileChange" 
          />
        </div>

        <!-- Selected Tags -->
        <div class="flex flex-wrap gap-2">
          <div 
            v-for="tagId in form.tags" 
            :key="tagId"
            class="px-3 py-1 bg-purple-50 text-purple-600 rounded-full text-sm flex items-center gap-1"
          >
            <span>#{{ getTagName(tagId) }}</span>
            <el-icon class="cursor-pointer hover:text-purple-800" @click="removeTag(tagId)"><Close /></el-icon>
          </div>
          
          <el-popover placement="bottom" :width="200" trigger="click">
            <template #reference>
              <button class="px-3 py-1 border border-gray-200 text-gray-500 rounded-full text-sm hover:bg-gray-50 flex items-center gap-1">
                <el-icon><Plus /></el-icon> 标签
              </button>
            </template>
            <div class="max-h-60 overflow-y-auto">
              <div 
                v-for="tag in store.flattenedTags" 
                :key="tag.id"
                class="px-2 py-1.5 hover:bg-gray-50 cursor-pointer text-sm rounded"
                @click="addTag(tag.id)"
              >
                {{ tag.name }}
              </div>
            </div>
          </el-popover>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Close, Camera, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import { useLittleGrowthStore, type Tag } from '@/stores/littleGrowth'

const route = useRoute()
const router = useRouter()
const store = useLittleGrowthStore()

const isEdit = computed(() => !!route.params.id)
const fileInput = ref<HTMLInputElement | null>(null)

const form = ref({
  date: dayjs().format('YYYY-MM-DD'),
  content: '',
  images: [] as string[],
  tags: [] as string[]
})

const showTagSuggestions = ref(false)
const tagSearchQuery = ref('')

onMounted(() => {
  if (isEdit.value) {
    const record = store.getRecordById(route.params.id as string)
    if (record) {
      form.value = {
        date: record.date,
        content: record.content,
        images: [...record.images],
        tags: [...record.tags]
      }
    }
  }
})

// Tag Helpers
const getTagName = (id: string) => {
  const t = store.flattenedTags.find(x => x.id === id)
  return t ? t.name : id
}

const addTag = (id: string) => {
  if (!form.value.tags.includes(id)) {
    form.value.tags.push(id)
  }
}

const removeTag = (id: string) => {
  form.value.tags = form.value.tags.filter(t => t !== id)
}

// Input handling for #tags
const handleInput = (val: string) => {
  const lastChar = val.slice(-1)
  
  // Simple detection logic
  const match = val.match(/#(\S*)$/)
  if (match) {
    showTagSuggestions.value = true
    tagSearchQuery.value = match[1]
  } else {
    showTagSuggestions.value = false
  }
}

const suggestedTags = computed(() => {
  const query = tagSearchQuery.value.toLowerCase()
  return store.flattenedTags.filter(t => t.name.toLowerCase().includes(query)).slice(0, 5)
})

const selectSuggestedTag = (tag: Tag) => {
  // Replace #query with nothing or keep it? 
  // Usually we replace "#query" with "#tag " or just add tag ID and remove text?
  // PRD says: "Input tags should be highlighted".
  // Since I can't highlight easily in textarea, I will just add it to the tag list below and remove the #text from textarea to keep it clean, OR keep #text in textarea.
  // Let's keep it clean: Remove the #part and add to tag list.
  
  form.value.content = form.value.content.replace(/#\S*$/, '')
  addTag(tag.id)
  showTagSuggestions.value = false
}

// Image Upload
const triggerUpload = () => {
  fileInput.value?.click()
}

const handleFileChange = (e: Event) => {
  const files = (e.target as HTMLInputElement).files
  if (!files) return
  
  // Simple mock upload: read as data URL
  Array.from(files).forEach(file => {
    if (form.value.images.length >= 9) return
    
    const reader = new FileReader()
    reader.onload = (e) => {
      if (e.target?.result) {
        form.value.images.push(e.target.result as string)
      }
    }
    reader.readAsDataURL(file)
  })
  
  // Reset input
  if (fileInput.value) fileInput.value.value = ''
}

const removeImage = (index: number) => {
  form.value.images.splice(index, 1)
}

// Save
const save = () => {
  if (isEdit.value) {
    store.updateRecord(route.params.id as string, {
      ...form.value
    })
    ElMessage.success('更新成功')
  } else {
    store.addRecord({
      ...form.value
    })
    ElMessage.success('发布成功')
  }
  router.back()
}
</script>
