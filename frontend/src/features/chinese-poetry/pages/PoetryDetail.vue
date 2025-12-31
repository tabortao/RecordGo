<template>
  <div class="min-h-screen bg-[#FDF6F8] dark:bg-gray-900 pb-10 font-sans">
    <!-- Header -->
    <div class="sticky top-0 z-20 bg-white/80 dark:bg-gray-800/80 backdrop-blur-md border-b border-pink-100 dark:border-gray-700 flex items-center justify-between p-4 shadow-sm">
      <div class="flex items-center gap-3">
        <el-icon :size="22" class="cursor-pointer text-gray-600 dark:text-gray-300 hover:text-pink-500 transition" @click="router.back()"><ArrowLeft /></el-icon>
        <h1 class="text-xl font-bold text-gray-800 dark:text-white">古诗详情</h1>
      </div>
      <div class="flex gap-4">
         <el-icon :size="24" class="cursor-pointer transition hover:scale-110" :class="isCollected ? 'text-yellow-400' : 'text-gray-300'" @click="toggleCollection">
             <component :is="isCollected ? StarFilled : Star" />
         </el-icon>
         <el-icon v-if="poem?.isCustom" :size="22" class="cursor-pointer text-gray-400 hover:text-blue-500 transition hover:scale-110" @click="router.push(`/homework/chinese/poetry/edit/${poem.id}`)"><Edit /></el-icon>
         <el-icon v-if="poem?.isCustom" :size="22" class="cursor-pointer text-gray-400 hover:text-red-500 transition hover:scale-110" @click="handleDelete"><Delete /></el-icon>
      </div>
    </div>

    <div v-if="poem" class="p-6 space-y-8 max-w-3xl mx-auto">
       <!-- Poem Content Card -->
       <div class="bg-white dark:bg-gray-800 rounded-3xl p-8 shadow-xl shadow-pink-100/50 dark:shadow-none border border-pink-50 dark:border-gray-700 text-center relative overflow-hidden">
         <!-- Decorative background elements -->
         <div class="absolute top-0 left-0 w-32 h-32 bg-pink-100/30 rounded-full blur-3xl -translate-x-10 -translate-y-10 pointer-events-none"></div>
         <div class="absolute bottom-0 right-0 w-40 h-40 bg-purple-100/30 rounded-full blur-3xl translate-x-10 translate-y-10 pointer-events-none"></div>

         <h2 class="text-3xl font-bold text-gray-800 dark:text-gray-100 mb-3 tracking-wide">{{ poem.title_cns }}</h2>
         <div class="text-sm text-gray-500 dark:text-gray-400 mb-10 flex justify-center items-center gap-3">
           <span class="bg-pink-50 dark:bg-gray-700 text-pink-600 dark:text-pink-300 px-3 py-1 rounded-full text-xs font-medium">{{ poem.dynasty_cns }}</span>
           <span class="font-medium text-gray-600 dark:text-gray-300">{{ poem.author_cns }}</span>
         </div>
         
         <div class="space-y-6 text-2xl leading-loose text-gray-700 dark:text-gray-200 font-serif tracking-widest mb-10">
           <p v-for="(line, idx) in poem.paragraphs_cns" :key="idx" class="hover:text-pink-600 transition-colors duration-300">{{ line }}</p>
         </div>

         <!-- Audio Player Controls -->
         <div class="flex flex-col items-center gap-3 mb-8">
           <button 
             @click="togglePlay"
             class="w-16 h-16 rounded-full flex items-center justify-center transition-all transform hover:scale-105 active:scale-95 shadow-lg"
             :class="isPlaying ? 'bg-gradient-to-br from-pink-400 to-rose-500 text-white shadow-pink-200' : 'bg-white text-pink-500 border-2 border-pink-100 hover:border-pink-300 hover:shadow-pink-100'"
           >
             <el-icon :size="32" v-if="!isPlaying" class="ml-1"><VideoPlay /></el-icon>
             <el-icon :size="32" v-else class="animate-pulse"><VideoPause /></el-icon>
           </button>
           <span class="text-xs text-gray-400 font-medium">点击朗读</span>
         </div>
         
         <div class="flex justify-center gap-6">
             <el-button 
               class="!h-12 !px-8 !text-base !rounded-full shadow-lg shadow-purple-100 hover:shadow-purple-200 transition-all transform hover:-translate-y-0.5 bg-gradient-to-r from-purple-400 to-indigo-400 border-none text-white" 
               @click="router.push(`/homework/chinese/poetry/recite/${poem.id}`)"
             >
               背诵打卡
             </el-button>
             <el-button 
               class="!h-12 !px-8 !text-base !rounded-full shadow-lg shadow-green-100 hover:shadow-green-200 transition-all transform hover:-translate-y-0.5 bg-gradient-to-r from-emerald-400 to-teal-400 border-none text-white" 
               @click="router.push(`/homework/chinese/poetry/dictate/${poem.id}`)"
             >
               默写打卡
             </el-button>
         </div>
       </div>

       <!-- Study Stats Section (New) -->
       <div v-if="studyRecord" class="grid grid-cols-3 gap-4">
          <div class="bg-white dark:bg-gray-800 rounded-2xl p-4 shadow-sm border border-gray-50 dark:border-gray-700 text-center">
             <div class="text-2xl font-bold text-gray-800 dark:text-white mb-1">{{ studyRecord.timesStudied }}</div>
             <div class="text-xs text-gray-400">学习次数</div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-2xl p-4 shadow-sm border border-gray-50 dark:border-gray-700 text-center">
             <div class="text-2xl font-bold mb-1" :class="studyRecord.isMastered ? 'text-green-500' : 'text-gray-400'">
                {{ studyRecord.isMastered ? '已掌握' : '学习中' }}
             </div>
             <div class="text-xs text-gray-400">当前状态</div>
          </div>
           <div class="bg-white dark:bg-gray-800 rounded-2xl p-4 shadow-sm border border-gray-50 dark:border-gray-700 text-center">
             <div class="text-2xl font-bold text-blue-500 mb-1">{{ lastDictationScore }}%</div>
             <div class="text-xs text-gray-400">最近默写</div>
          </div>
       </div>

       <!-- Tabs for Appreciation/Translation -->
       <div class="bg-white dark:bg-gray-800 rounded-3xl shadow-sm border border-pink-50 dark:border-gray-700 overflow-hidden min-h-[300px]">
         <div class="flex border-b border-gray-100 dark:border-gray-700">
           <button 
             v-for="tab in ['译文', '赏析', '趣事']" 
             :key="tab"
             @click="activeTab = tab as any"
             class="flex-1 py-4 text-sm font-medium transition-all relative"
             :class="activeTab === tab ? 'text-pink-600 bg-pink-50/50 dark:bg-pink-900/20' : 'text-gray-500 hover:text-gray-700 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700/50'"
           >
             {{ tab }}
             <div v-if="activeTab === tab" class="absolute bottom-0 left-1/2 -translate-x-1/2 w-8 h-1 bg-pink-500 rounded-t-full"></div>
           </button>
         </div>
         
         <div class="p-8">
           <div v-if="isLoadingAI" class="flex flex-col items-center justify-center py-12 text-gray-400">
             <div class="w-12 h-12 border-4 border-pink-200 border-t-pink-500 rounded-full animate-spin mb-4"></div>
             <span class="animate-pulse">AI 正在生成{{ activeTab }}...</span>
           </div>
           
           <div v-else-if="getContent(activeTab)" class="prose dark:prose-invert max-w-none text-gray-700 dark:text-gray-300 whitespace-pre-wrap leading-relaxed text-justify">
             {{ getContent(activeTab) }}
           </div>
           
           <div v-else class="text-center py-12">
             <div class="bg-gray-50 dark:bg-gray-700/50 w-20 h-20 rounded-full flex items-center justify-center mx-auto mb-4 text-gray-300">
                <el-icon :size="40"><MagicStick /></el-icon>
             </div>
             <p class="text-gray-400 mb-6 text-sm">暂无{{ activeTab }}内容，点击下方按钮生成</p>
             <el-button type="primary" plain round class="!px-6 !border-pink-300 !text-pink-500 hover:!bg-pink-50 hover:!border-pink-400" @click="generateContent(activeTab)">
               <el-icon class="mr-2"><MagicStick /></el-icon>
               AI 生成{{ activeTab }}
             </el-button>
           </div>
         </div>
       </div>
    </div>
    
    <div v-else class="p-20 text-center text-gray-500">
        <el-icon :size="40" class="is-loading mb-4 text-pink-300"><Loading /></el-icon>
        <p>古诗加载中...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, VideoPlay, VideoPause, Loading, MagicStick, Star, StarFilled, Edit, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Poem } from '../types'
import { usePoetryStore } from '../stores/poetryStore'
import { useAppState } from '@/stores/appState'
import { createCustomAudio } from '@/utils/customTTS'

const route = useRoute()
const router = useRouter()
const appState = useAppState()
const poetryStore = usePoetryStore()

const poemId = computed(() => parseInt(route.params.id as string))
const poem = computed<Poem | undefined>(() => poetryStore.getPoemById(poemId.value))

const activeTab = ref<'译文' | '赏析' | '趣事'>('译文')
const isLoadingAI = ref(false)
const isPlaying = ref(false)
const audioPlayer = ref<HTMLAudioElement | null>(null)
const isCollected = computed(() => poem.value ? poetryStore.collections.includes(poem.value.id) : false)
const studyRecord = computed(() => poem.value ? poetryStore.getRecord(poem.value.id) : null)
const lastDictationScore = computed(() => {
    if (!poem.value) return 0
    const history = poetryStore.getDictationHistory(poem.value.id)
    return history.length > 0 ? history[0].score : 0
})

onMounted(() => {
  poetryStore.init()
  if (!poem.value) {
    // Retry once in case init is slow, though init is synchronous for localStorage
    if (!poetryStore.getPoemById(poemId.value)) {
       ElMessage.error('古诗不存在')
       router.back()
    }
  }
  // Check initial tab content
  checkAndGenerateContent(activeTab.value)
})

watch(activeTab, (newVal) => {
    checkAndGenerateContent(newVal)
})

const checkAndGenerateContent = (type: string) => {
    if (!getContent(type) && !isLoadingAI.value) {
        generateContent(type as any)
    }
}

const getContent = (type: string) => {
  if (!poem.value) return ''
  switch (type) {
    case '译文': return poem.value.translation
    case '赏析': return poem.value.appreciation
    case '趣事': return poem.value.fun_fact
    default: return ''
  }
}

const generateContent = async (_type: '译文' | '赏析' | '趣事') => {
  if (!poem.value) return
  
  isLoadingAI.value = true
  try {
    // In a real app, we would pass the 'type' to the generator.
    // Our mock generator currently generates all fields at once.
    await poetryStore.generateAIContent(poem.value.id)
    ElMessage.success('AI 内容生成成功')
  } catch (e: any) {
    ElMessage.error(e.message || '生成失败')
  } finally {
    isLoadingAI.value = false
  }
}

const toggleCollection = () => {
    if (poem.value) {
        poetryStore.toggleCollection(poem.value.id)
        ElMessage.success(isCollected.value ? '已收藏' : '已取消收藏')
    }
}

const handleDelete = () => {
    if (!poem.value || !poem.value.isCustom) return
    ElMessageBox.confirm('确定要删除这首古诗吗？此操作不可恢复。', '删除确认', {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
        poetryStore.deleteCustomPoem(poem.value!.id)
        ElMessage.success('删除成功')
        router.back()
    })
}

onUnmounted(() => {
  stopAudio()
})

const togglePlay = async () => {
  if (isPlaying.value) {
    stopAudio()
  } else {
    if (!poem.value) return
    
    // Check if custom audio URL exists
    if (poem.value.audio_url) {
        playAudio(poem.value.audio_url)
        return
    }

    // Use TTS
    const text = poem.value.title_cns + '。' + poem.value.author_cns + '。' + poem.value.paragraphs_cns.join('。')
    
    // Check for custom TTS profile
    let profile = null
    if (appState.speech.engine === 'custom' && appState.speech.activeCustomId) {
        profile = appState.speech.customProfiles.find((p: any) => p.id === appState.speech.activeCustomId)
    }

    if (profile) {
        try {
            const audio = await createCustomAudio(text, profile)
            if (audio) {
                if (audioPlayer.value) {
                    audioPlayer.value.pause()
                }
                audioPlayer.value = audio
                audio.onended = () => isPlaying.value = false
                audio.play()
                isPlaying.value = true
            }
        } catch (e) {
            ElMessage.error('语音生成失败')
        }
    } else {
        // Browser native fallback
        const ut = new SpeechSynthesisUtterance(text)
        ut.lang = 'zh-CN'
        ut.onend = () => isPlaying.value = false
        window.speechSynthesis.speak(ut)
        isPlaying.value = true
    }
  }
}

const playAudio = (url: string) => {
    if (audioPlayer.value) {
        audioPlayer.value.src = url
        audioPlayer.value.play()
        isPlaying.value = true
        audioPlayer.value.onended = () => isPlaying.value = false
    } else {
        const audio = new Audio(url)
        audio.onended = () => isPlaying.value = false
        audio.play()
        audioPlayer.value = audio
        isPlaying.value = true
    }
}

const stopAudio = () => {
    if (audioPlayer.value) {
        audioPlayer.value.pause()
        audioPlayer.value.currentTime = 0
    }
    window.speechSynthesis.cancel()
    isPlaying.value = false
}
</script>
