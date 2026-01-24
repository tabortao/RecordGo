<template>
  <div class="min-h-screen bg-[#F5F7FA] dark:bg-gray-900 flex flex-col relative">
    <!-- Header -->
    <div class="px-4 py-3 bg-white/80 dark:bg-gray-900/80 backdrop-blur-md sticky top-0 z-30 flex items-center justify-between border-b border-gray-100 dark:border-gray-800 transition-colors">
      <div 
        class="w-8 h-8 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center cursor-pointer hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors text-gray-600 dark:text-gray-300"
        @click="router.back()"
      >
        <el-icon><ArrowLeft /></el-icon>
      </div>
      
      <span class="font-bold text-gray-800 dark:text-gray-100 text-base">{{ isEdit ? '编辑美好回忆' : '记录美好瞬间' }}</span>
      
      <button 
        class="bg-purple-600 text-white px-5 py-1.5 rounded-full text-sm font-bold shadow-lg shadow-purple-200 dark:shadow-purple-900/20 hover:bg-purple-700 hover:scale-105 active:scale-95 transition-all disabled:opacity-50 disabled:hover:scale-100"
        :disabled="(!form.content && previewImages.length === 0 && !audioUrl) || loading || uploading"
        @click="save"
      >
        {{ loading ? (isEdit ? '保存中...' : '发布中...') : (uploading ? '上传中...' : (isEdit ? '保存' : '发布')) }}
      </button>
    </div>

    <div class="flex-1 overflow-y-auto p-4 sm:p-6 pb-24">
      <div class="max-w-2xl mx-auto space-y-6">
        
        <!-- Meta Card (Date & Visibility) -->
        <div class="bg-white dark:bg-gray-800 rounded-3xl p-2 shadow-sm border border-gray-100 dark:border-gray-700 flex flex-wrap items-center gap-2">
           <div class="flex-1 min-w-[200px] relative group">
              <el-date-picker
                v-model="form.date"
                type="datetime"
                placeholder="选择日期时间"
                format="YYYY年MM月DD日 HH:mm"
                :clearable="false"
                class="!w-full custom-datepicker"
                :prefix-icon="Calendar"
              />
           </div>
           <div class="bg-gray-100 dark:bg-gray-700 h-8 w-[1px] mx-1 hidden sm:block"></div>
           <div class="flex p-1 bg-gray-100 dark:bg-gray-700 rounded-2xl">
             <button 
               v-for="opt in [{l: 0, t: '家庭可见'}, {l: 1, t: '仅自己'}]" 
               :key="opt.l"
               class="px-3 py-1.5 rounded-xl text-xs font-medium transition-all"
               :class="form.visibility === opt.l ? 'bg-white dark:bg-gray-600 text-purple-600 shadow-sm' : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'"
               @click="form.visibility = opt.l"
             >
               {{ opt.t }}
             </button>
           </div>
        </div>

        <!-- Main Editor Card -->
        <div class="bg-white dark:bg-gray-800 rounded-3xl p-5 shadow-sm border border-gray-100 dark:border-gray-700 relative group focus-within:ring-2 ring-purple-100 dark:ring-purple-900/30 transition-all">
          <textarea
            v-model="form.content"
            rows="6"
            placeholder="今天发生了什么有趣的事？(输入 # 可快速添加标签)"
            class="w-full text-lg text-gray-800 dark:text-gray-200 placeholder-gray-400 bg-transparent border-none outline-none resize-none leading-relaxed custom-scrollbar"
            @input="(e) => handleInput((e.target as HTMLTextAreaElement).value)"
          ></textarea>
          
          <!-- Smart Tag Suggestions -->
          <transition name="fade-slide">
            <div v-if="showTagSuggestions" class="absolute top-full left-0 z-20 bg-white dark:bg-gray-800 shadow-xl border border-gray-100 dark:border-gray-700 rounded-2xl mt-2 w-full overflow-hidden">
              <div class="p-3 bg-gray-50 dark:bg-gray-700/50 text-xs font-bold text-gray-500 dark:text-gray-400">标签建议</div>
              <div class="max-h-48 overflow-y-auto p-2">
                <div class="flex flex-wrap gap-2">
                  <div 
                    v-for="tag in suggestedTags" 
                    :key="tag.id"
                    class="px-3 py-1.5 rounded-xl bg-purple-50 dark:bg-purple-900/20 text-purple-600 dark:text-purple-300 hover:bg-purple-100 dark:hover:bg-purple-900/40 cursor-pointer text-sm transition-colors flex items-center gap-1"
                    @click="selectSuggestedTag(tag)"
                  >
                    <span>#</span> {{ tag.name }}
                  </div>
                </div>
                <div v-if="suggestedTags.length === 0" class="px-3 py-4 text-center text-gray-400 text-xs">
                  输入空格以创建新标签
                </div>
              </div>
            </div>
          </transition>
        </div>

        <!-- Media Area -->
        <div class="space-y-4">
          <!-- Images Grid -->
          <div class="grid grid-cols-3 sm:grid-cols-4 gap-3">
             <div 
               v-for="(img, index) in previewImages" 
               :key="index"
               class="relative aspect-square rounded-2xl overflow-hidden group shadow-sm border border-gray-100 dark:border-gray-700"
             >
               <el-image :src="img" fit="cover" class="w-full h-full transition-transform duration-500 group-hover:scale-110" />
               <div class="absolute inset-0 bg-black/20 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                  <div 
                    class="bg-white/90 dark:bg-gray-800/90 text-red-500 p-2 rounded-full cursor-pointer hover:scale-110 transition-transform shadow-sm"
                    @click="removeImage(index)"
                  >
                    <el-icon><Close /></el-icon>
                  </div>
               </div>
             </div>
             
             <div 
               v-if="previewImages.length < 9"
               class="aspect-square rounded-2xl border-2 border-dashed border-gray-200 dark:border-gray-700 flex flex-col items-center justify-center text-gray-400 hover:border-purple-400 hover:text-purple-500 hover:bg-purple-50 dark:hover:bg-purple-900/10 transition-all cursor-pointer bg-white dark:bg-gray-800"
               @click="triggerUpload"
             >
               <el-icon :size="24" v-if="!uploading" class="mb-1"><Camera /></el-icon>
               <el-icon :size="24" v-else class="animate-spin mb-1"><Loading /></el-icon>
               <span class="text-xs font-medium">{{ uploading ? '上传中' : '添加照片' }}</span>
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

          <!-- Audio Recorder -->
          <div class="bg-white dark:bg-gray-800 rounded-2xl p-4 shadow-sm border border-gray-100 dark:border-gray-700 flex items-center gap-4">
              <div 
                class="w-12 h-12 rounded-full flex items-center justify-center cursor-pointer transition-all active:scale-95 shadow-sm"
                :class="isRecording ? 'bg-red-50 text-red-500 ring-4 ring-red-100 dark:ring-red-900/30 animate-pulse' : 'bg-purple-50 text-purple-600 hover:bg-purple-100 dark:bg-purple-900/20 dark:text-purple-300'"
                @click="toggleRecording"
              >
                  <el-icon :size="22"><Microphone /></el-icon>
              </div>
              
              <div class="flex-1 min-w-0">
                  <div v-if="isRecording" class="flex items-center gap-2">
                    <div class="w-2 h-2 rounded-full bg-red-500 animate-ping"></div>
                    <span class="text-red-500 font-bold tracking-wide">{{ formatRecordingTime(recordingTime) }}</span>
                    <span class="text-xs text-gray-400 ml-2">正在录音...</span>
                  </div>
                  <div v-else-if="audioUrl" class="flex items-center gap-3">
                      <div class="flex-1 bg-gray-100 dark:bg-gray-700 rounded-full px-3 py-1 flex items-center gap-2">
                         <el-icon class="text-purple-500"><VideoPlay /></el-icon>
                         <audio :src="audioUrl" controls class="h-8 w-full bg-transparent custom-audio" />
                      </div>
                      <div class="w-8 h-8 rounded-full bg-gray-100 dark:bg-gray-700 flex items-center justify-center text-gray-500 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 cursor-pointer transition-colors" @click="removeAudio">
                        <el-icon><Delete /></el-icon>
                      </div>
                  </div>
                  <div v-else class="text-gray-400 text-sm">点击麦克风录制语音备注</div>
              </div>

              <!-- Hidden Upload for manual selection if needed -->
              <el-upload
                  v-if="!isRecording && !audioUrl"
                  action="#"
                  :auto-upload="false"
                  accept="audio/*"
                  :show-file-list="false"
                  :on-change="handleAudioUpload"
              >
                  <div class="text-xs text-purple-600 dark:text-purple-400 font-medium cursor-pointer hover:underline px-2">上传文件</div>
              </el-upload>
          </div>
        </div>

        <!-- Tags List -->
        <div class="flex flex-wrap gap-2">
          <div 
            v-for="tagId in form.tags" 
            :key="tagId"
            class="px-3 py-1.5 bg-purple-50 dark:bg-purple-900/30 text-purple-600 dark:text-purple-300 rounded-xl text-sm font-medium flex items-center gap-1.5 border border-purple-100 dark:border-purple-800/50"
          >
            <span># {{ getTagName(tagId) }}</span>
            <el-icon class="cursor-pointer hover:text-purple-800 dark:hover:text-purple-100" @click="removeTag(tagId)"><Close /></el-icon>
          </div>
          
          <div
            v-for="pt in pendingTags"
            :key="pt.parentName ? (pt.parentName + '/' + pt.name) : ('__pending__' + pt.name)"
            class="px-3 py-1.5 bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 rounded-xl text-sm font-medium flex items-center gap-1.5 border border-gray-200 dark:border-gray-600"
          >
            <span># 待创建：{{ pt.parentName ? (pt.parentName + '/' + pt.name) : pt.name }}</span>
            <el-icon class="cursor-pointer hover:text-gray-800 dark:hover:text-gray-100" @click="removePending(pt)"><Close /></el-icon>
          </div>
          
          <el-popover placement="bottom" :width="280" trigger="click" popper-class="custom-popover">
            <template #reference>
              <button class="px-3 py-1.5 border border-dashed border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-400 rounded-xl text-sm hover:border-purple-400 hover:text-purple-500 dark:hover:border-purple-400 dark:hover:text-purple-400 flex items-center gap-1 transition-colors">
                <el-icon><Plus /></el-icon> 添加标签
              </button>
            </template>
            <div class="p-1">
              <div class="text-xs font-bold text-gray-500 mb-2 px-1">选择已有标签</div>
              <div class="flex flex-wrap gap-1.5 max-h-48 overflow-y-auto custom-scrollbar mb-3">
                <div 
                  v-for="tag in store.tags" 
                  :key="tag.id"
                  class="px-2.5 py-1 bg-gray-50 dark:bg-gray-700/50 hover:bg-purple-50 dark:hover:bg-purple-900/30 hover:text-purple-600 dark:hover:text-purple-300 cursor-pointer text-sm rounded-lg transition-colors border border-transparent hover:border-purple-100 dark:hover:border-purple-800"
                  @click="addTagById(tag.id)"
                >
                  {{ tag.name }}
                </div>
              </div>
              
              <div class="pt-2 border-t border-gray-100 dark:border-gray-700">
                <div class="text-xs font-bold text-gray-500 mb-2 px-1">新建标签</div>
                <div class="flex items-center gap-2">
                  <el-input v-model="newTagInput" placeholder="如：学习/语文" size="small" class="flex-1" />
                  <el-button size="small" type="primary" color="#9333ea" @click="addNewTagDraft">添加</el-button>
                </div>
                <div class="text-[10px] text-gray-400 mt-1 px-1">支持 "父/子" 格式快速创建分类标签</div>
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
import { ArrowLeft, Close, Camera, Plus, Microphone, Delete, Loading, Calendar, VideoPlay } from '@element-plus/icons-vue'
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
  form.value.content = val
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
  return store.tags.filter(t => t.name.toLowerCase().includes(query)).slice(0, 10)
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
const formatRecordingTime = (seconds: number) => {
    const m = Math.floor(seconds / 60)
    const s = seconds % 60
    return `${m}:${s.toString().padStart(2, '0')}`
}

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
      date: dayjs(form.value.date).format('YYYY-MM-DD HH:mm:ss'), // 修复日期问题：明确格式化
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

<style scoped>
.custom-datepicker :deep(.el-input__wrapper) {
  box-shadow: none !important;
  background-color: transparent;
  padding-left: 0;
}
.custom-datepicker :deep(.el-input__inner) {
  font-weight: 600;
  color: #4b5563;
  font-size: 0.95rem;
}
.dark .custom-datepicker :deep(.el-input__inner) {
  color: #e5e7eb;
}

.custom-audio {
  filter: drop-shadow(0 1px 2px rgba(0,0,0,0.1));
}
.custom-audio::-webkit-media-controls-panel {
  background-color: transparent;
}

.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.2s ease;
}
.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
