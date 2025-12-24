<template>
  <div class="min-h-screen bg-white dark:bg-gray-900 flex flex-col">
    <!-- Header -->
    <div class="px-4 py-3 border-b border-gray-100 dark:border-gray-800 flex items-center justify-between sticky top-0 bg-white dark:bg-gray-900 z-20">
      <div class="flex items-center gap-2 cursor-pointer" @click="router.back()">
        <el-icon class="text-gray-600 dark:text-gray-300 text-lg"><ArrowLeft /></el-icon>
        <span class="text-gray-600 dark:text-gray-300">取消</span>
      </div>
      <h1 class="font-bold text-lg dark:text-white">{{ isEdit ? '编辑记录' : '新记录' }}</h1>
      <button 
        class="bg-purple-600 text-white px-4 py-1.5 rounded-full text-sm font-medium hover:bg-purple-700 transition-colors disabled:opacity-50"
        :disabled="(!form.content && previewImages.length === 0 && !audioUrl) || loading || uploading"
        @click="save"
      >
        {{ loading ? (isEdit ? '保存中...' : '发布中...') : (uploading ? '上传中...' : (isEdit ? '保存' : '发布')) }}
      </button>
    </div>

    <div class="flex-1 overflow-y-auto p-4">
      <div class="max-w-2xl mx-auto space-y-6">
        <!-- Date & Visibility -->
        <div class="flex items-center gap-2 justify-between">
          <el-date-picker
            v-model="form.date"
            type="datetime"
            placeholder="选择日期时间"
            format="YYYY年MM月DD日 HH:mm"
            :clearable="false"
            class="!w-52"
          />
          <el-radio-group v-model="form.visibility" size="small">
            <el-radio-button :label="0">家庭可见</el-radio-button>
            <el-radio-button :label="1">仅自己</el-radio-button>
          </el-radio-group>
        </div>

        <!-- Content -->
        <div class="relative">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="6"
            placeholder="记录这一刻的美好... 输入 # 可快速添加标签"
            class="text-lg dark:bg-gray-800"
            resize="none"
            @input="handleInput"
          />
          
          <!-- Smart Tag Suggestions -->
          <div v-if="showTagSuggestions" class="absolute top-full left-0 z-10 bg-white dark:bg-gray-800 shadow-lg border border-gray-100 dark:border-gray-700 rounded-lg mt-1 w-full max-h-48 overflow-y-auto">
            <div class="flex flex-wrap gap-2 p-2">
              <div 
                v-for="tag in suggestedTags" 
                :key="tag.id"
                class="px-3 py-1.5 rounded-full bg-gray-100 dark:bg-gray-700 hover:bg-purple-100 dark:hover:bg-purple-900/50 hover:text-purple-600 dark:hover:text-purple-300 cursor-pointer text-sm text-gray-700 dark:text-gray-300 transition-colors"
                @click="selectSuggestedTag(tag)"
              >
                {{ tag.name }}
              </div>
            </div>
            <div v-if="suggestedTags.length === 0" class="px-3 py-2 text-gray-400 text-xs">
              输入空格以创建新标签
            </div>
          </div>
        </div>

        <!-- Images -->
        <div>
          <div class="grid grid-cols-4 gap-2 mb-2">
            <div 
              v-for="(img, index) in previewImages" 
              :key="index"
              class="relative aspect-square rounded-xl overflow-hidden group bg-gray-100 dark:bg-gray-700 border border-gray-200 dark:border-gray-600"
            >
              <!-- Changed to object-contain for 100% display -->
              <img :src="img" class="w-full h-full object-contain" />
              <div class="absolute top-1 right-1 p-1 bg-black/50 rounded-full text-white opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer" @click="removeImage(index)">
                <el-icon><Close /></el-icon>
              </div>
            </div>
            
            <div 
              v-if="previewImages.length < 9"
              class="aspect-square rounded-xl border-2 border-dashed border-gray-200 dark:border-gray-600 flex flex-col items-center justify-center text-gray-400 hover:border-purple-400 hover:text-purple-500 transition-colors cursor-pointer bg-gray-50 dark:bg-gray-800"
              @click="triggerUpload"
            >
              <el-icon :size="20" v-if="!uploading"><Camera /></el-icon>
              <el-icon :size="20" v-else class="animate-spin"><Loading /></el-icon>
              <span class="text-[10px] mt-1">{{ uploading ? '上传中' : '添加照片' }}</span>
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

        <!-- Audio -->
        <div class="flex items-center gap-3 p-3 bg-gray-50 dark:bg-gray-800 rounded-xl">
            <div 
              class="w-10 h-10 rounded-full flex items-center justify-center cursor-pointer transition-colors"
              :class="isRecording ? 'bg-red-100 text-red-500 animate-pulse' : 'bg-white dark:bg-gray-700 text-purple-600 shadow-sm hover:bg-purple-50 dark:hover:bg-gray-600'"
              @click="toggleRecording"
            >
                <el-icon :size="20"><Microphone /></el-icon>
            </div>
            
            <div class="flex-1">
                <div v-if="isRecording" class="text-red-500 text-sm font-medium">正在录音 {{ recordingTime }}s</div>
                <div v-else-if="audioUrl" class="flex items-center gap-2">
                    <audio :src="audioUrl" controls class="h-8 w-full max-w-[200px]" />
                    <el-icon class="text-gray-400 cursor-pointer hover:text-red-500" @click="removeAudio"><Delete /></el-icon>
                </div>
                <div v-else class="text-gray-400 text-sm">点击录制音频或上传</div>
            </div>

            <el-upload
                action="#"
                :auto-upload="false"
                accept="audio/*"
                :show-file-list="false"
                :on-change="handleAudioUpload"
            >
                <el-button size="small" round :loading="uploading">上传</el-button>
            </el-upload>
        </div>

        <!-- Selected Tags -->
        <div class="flex flex-wrap gap-2">
          <div 
            v-for="tagId in form.tags" 
            :key="tagId"
            class="px-3 py-1 bg-purple-50 dark:bg-purple-900/30 text-purple-600 dark:text-purple-300 rounded-full text-sm flex items-center gap-1"
          >
            <span>{{ getTagName(tagId) }}</span>
            <el-icon class="cursor-pointer hover:text-purple-800 dark:hover:text-purple-200" @click="removeTag(tagId)"><Close /></el-icon>
          </div>
          <div
            v-for="pt in pendingTags"
            :key="pt.parentName ? (pt.parentName + '/' + pt.name) : ('__pending__' + pt.name)"
            class="px-3 py-1 bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 rounded-full text-sm flex items-center gap-1"
          >
            <span>待创建：{{ pt.parentName ? (pt.parentName + '/' + pt.name) : pt.name }}</span>
            <el-icon class="cursor-pointer hover:text-gray-800 dark:hover:text-gray-200" @click="removePending(pt)"><Close /></el-icon>
          </div>
          
          <el-popover placement="bottom" :width="240" trigger="click">
            <template #reference>
              <button class="px-3 py-1 border border-gray-200 dark:border-gray-700 text-gray-500 dark:text-gray-400 rounded-full text-sm hover:bg-gray-50 dark:hover:bg-gray-800 flex items-center gap-1">
                <el-icon><Plus /></el-icon> 标签
              </button>
            </template>
            <div class="space-y-2 max-h-64 overflow-y-auto dark:bg-gray-800 dark:text-gray-200 p-2">
              <div class="text-xs text-gray-500">选择已有标签</div>
              <div 
                v-for="tag in store.tags" 
                :key="tag.id"
                class="px-2 py-1.5 hover:bg-gray-50 dark:hover:bg-gray-700 cursor-pointer text-sm rounded"
                @click="addTagById(tag.id)"
              >
                {{ tag.name }}
              </div>
              <div class="pt-2 border-t border-gray-200 dark:border-gray-700"></div>
              <div class="text-xs text-gray-500">新建标签（支持 A/B）</div>
              <div class="flex items-center gap-2">
                <el-input v-model="newTagInput" placeholder="如：学习/语文 或 语文" size="small" />
                <el-button size="small" type="primary" @click="addNewTagDraft">添加</el-button>
              </div>
            </div>
          </el-popover>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Close, Camera, Plus, Microphone, Delete, Loading } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import { useLittleGrowthStore, type Tag } from '@/stores/littleGrowth'
import { prepareUpload } from '@/utils/image'
import http from '@/services/http'

const route = useRoute()
const router = useRouter()
const store = useLittleGrowthStore()

const isEdit = computed(() => !!route.params.id)
const loading = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

const form = ref({
  date: dayjs().toDate(),
  content: '',
  tags: [] as string[], // IDs
  visibility: 0
})

// Files management
const previewImages = ref<string[]>([]) 
const uploading = ref(false)

const audioUrl = ref<string>('')

// Recording
const isRecording = ref(false)
const recordingTime = ref(0)
let mediaRecorder: MediaRecorder | null = null
let timer: number | null = null
let chunks: Blob[] = []

onMounted(async () => {
  await store.fetchTags()
  if (isEdit.value) {
    try {
      const record = await store.fetchRecord(route.params.id as string)
      if (record) {
        form.value = {
          date: dayjs(record.date).toDate(),
          content: record.content,
          tags: [...record.tags],
          visibility: (record as any).visibility || 0
        }
        previewImages.value = [...(record.images || [])]
        if (record.audio) audioUrl.value = record.audio
      }
    } catch (e) {
      console.error(e)
    }
  }
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
  if (mediaRecorder && isRecording.value) {
    mediaRecorder.stop()
  }
})

// --- Tags ---
const pendingTags = ref<{ name: string; parentName?: string }[]>([])
const addTagById = (id: string) => {
  if (!form.value.tags.includes(id)) {
    form.value.tags.push(id)
  }
}

// 移除未使用的 addTag 函数

const removeTag = (id: string) => {
  form.value.tags = form.value.tags.filter(t => t !== id)
}

const getTagName = (id: string) => {
  return store.tags.find(t => t.id === id)?.name || '未知标签'
}

// Input handling for #tags
const showTagSuggestions = ref(false)
const tagSearchQuery = ref('')

const handleInput = async (val: string) => {
  const match = val.match(/#([^\s]*)$/)
  if (match) {
    showTagSuggestions.value = true
    tagSearchQuery.value = match[1]
    return
  }
  const spaceMatch = val.match(/#([^\s]+)\s$/)
  if (spaceMatch) {
    const raw = spaceMatch[1]
    const [p, c] = raw.split('/')
    if (c) {
      if (!pendingTags.value.some((pt) => pt.name === c && pt.parentName === p)) pendingTags.value.push({ name: c, parentName: p })
    } else {
      if (!pendingTags.value.some((pt) => pt.name === p)) pendingTags.value.push({ name: p })
    }
    form.value.content = form.value.content.replace(/#\S+\s$/, '')
    showTagSuggestions.value = false
  } else {
    showTagSuggestions.value = false
  }
}

const suggestedTags = computed(() => {
  const query = tagSearchQuery.value.toLowerCase()
  return store.tags.filter(t => t.name.toLowerCase().includes(query)).slice(0, 5)
})

const selectSuggestedTag = (tag: Tag) => {
  form.value.content = form.value.content.replace(/#\S*$/, '')
  addTagById(tag.id)
  showTagSuggestions.value = false
}

function removePending(pt: { name: string; parentName?: string }) {
  pendingTags.value = pendingTags.value.filter(x => !(x.name === pt.name && x.parentName === pt.parentName))
}

const newTagInput = ref('')
function addNewTagDraft() {
  const raw = (newTagInput.value || '').trim()
  if (!raw) return
  const [p, c] = raw.split('/')
  if (c) {
    if (!pendingTags.value.some((x) => x.name === c && x.parentName === p)) pendingTags.value.push({ name: c, parentName: p })
  } else {
    if (!pendingTags.value.some((x) => x.name === p)) pendingTags.value.push({ name: p })
  }
  newTagInput.value = ''
}

// --- Images ---
const triggerUpload = () => {
  fileInput.value?.click()
}

const handleFileChange = async (e: Event) => {
  const files = (e.target as HTMLInputElement).files
  if (!files) return
  
  uploading.value = true
  const fileList = Array.from(files)
  
  for (const file of fileList) {
    if (previewImages.value.length >= 9) break
    
    try {
      const webp = await prepareUpload(file, 0.75)
      const formData = new FormData()
      formData.append("file", webp, (webp as any).name || file.name)
      
      const res = await http.post("/upload/growth-file", formData, {
           headers: { "Content-Type": "multipart/form-data" }
      })
      
      let url = res.url
      if (!url.startsWith('http')) {
           url = `${import.meta.env.VITE_API_BASE}/api/${url}`
      }
      previewImages.value.push(url)
    } catch (err) {
      console.error("Upload failed", err)
      ElMessage.error(`图片 ${file.name} 上传失败`)
    }
  }
  
  uploading.value = false
  if (fileInput.value) fileInput.value.value = ''
}

const removeImage = (index: number) => {
  previewImages.value.splice(index, 1)
}

// --- Audio ---
const toggleRecording = async () => {
  if (isRecording.value) {
    mediaRecorder?.stop()
    isRecording.value = false
    if (timer) clearInterval(timer)
  } else {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      mediaRecorder = new MediaRecorder(stream)
      chunks = []
      
      mediaRecorder.ondataavailable = (e) => chunks.push(e.data)
      mediaRecorder.onstop = async () => {
        const blob = new Blob(chunks, { type: 'audio/m4a' })
        await uploadAudio(blob)
        stream.getTracks().forEach(track => track.stop())
      }
      
      mediaRecorder.start()
      isRecording.value = true
      recordingTime.value = 0
      timer = window.setInterval(() => recordingTime.value++, 1000)
    } catch (e) {
      ElMessage.error('无法访问麦克风')
    }
  }
}

const uploadAudio = async (blob: Blob) => {
    uploading.value = true
    try {
        const formData = new FormData()
        formData.append("file", blob, `record_${Date.now()}.m4a`)
        const res = await http.post("/upload/growth-file", formData, {
             headers: { "Content-Type": "multipart/form-data" }
        })
        let url = res.url
        if (!url.startsWith('http')) {
             url = `${import.meta.env.VITE_API_BASE}/api/${url}`
        }
        audioUrl.value = url
    } catch (e) {
        ElMessage.error('音频上传失败')
    } finally {
        uploading.value = false
    }
}

const handleAudioUpload = async (file: any) => {
  if (file.raw) {
    await uploadAudio(file.raw)
  }
}

const removeAudio = () => {
  audioUrl.value = ''
}

const save = async () => {
  if (!form.value.content && previewImages.value.length === 0 && !audioUrl.value) return
  
  loading.value = true
  try {
    // 解析并创建待创建标签，得到最终 IDs
    const resolvedIds: string[] = []
    const emojiSeeds = ['🚀','⭐','🍎','📚','🎯','🌟','💡','📝','🏃','🎵','🌈','🧠','🧩','🔬','🎨']
    const randomEmoji = () => emojiSeeds[Math.floor(Math.random() * emojiSeeds.length)]
    for (const pt of pendingTags.value) {
      let parentId: string | undefined
      if (pt.parentName) {
        const existP = store.tags.find(t => t.name === pt.parentName)
        if (existP) parentId = existP.id
        else {
          const createdP = await store.createTag(`${randomEmoji()}${pt.parentName}`)
          parentId = createdP.id
        }
      }
      const exist = store.tags.find(t => t.name === pt.name && (!parentId || t.parentId === parentId))
      if (exist) resolvedIds.push(exist.id)
      else {
        const created = await store.createTag(`${randomEmoji()}${pt.name}`, undefined, parentId)
        resolvedIds.push(created.id)
      }
    }

    const data = {
      ...form.value,
      images: previewImages.value,
      audio: audioUrl.value,
      tags: Array.from(new Set([...(form.value.tags || []), ...resolvedIds]))
    }

    if (isEdit.value) {
       await store.updateRecord(route.params.id as string, data)
    } else {
      await store.createRecord(data)
    }
    
    pendingTags.value = []
    ElMessage.success('发布成功')
    router.back()
  } catch (e) {
    console.error(e)
    ElMessage.error('发布失败')
  } finally {
    loading.value = false
  }
}
</script>
