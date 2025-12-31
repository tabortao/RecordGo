<template>
  <div class="min-h-screen bg-[#FDF6F8] dark:bg-gray-900 pb-10 font-sans">
    <!-- Header -->
    <div class="sticky top-0 z-20 bg-white/80 dark:bg-gray-800/80 backdrop-blur-md border-b border-pink-100 dark:border-gray-700 flex items-center justify-between p-4 shadow-sm">
      <div class="flex items-center gap-3">
        <el-icon :size="22" class="cursor-pointer text-gray-600 dark:text-gray-300 hover:text-pink-500 transition" @click="router.back()"><ArrowLeft /></el-icon>
        <h1 class="text-xl font-bold text-gray-800 dark:text-white">古诗详情</h1>
      </div>
      <div class="flex gap-3">
         <div class="flex items-center gap-2">
            <el-button circle size="small" :type="showPinyin ? 'primary' : ''" @click="showPinyin = !showPinyin">
                <span class="text-xs font-bold" :class="showPinyin ? 'text-white' : 'text-gray-500'">拼</span>
            </el-button>
         </div>
         <el-icon :size="24" class="cursor-pointer transition hover:scale-110" :class="isCollected ? 'text-yellow-400' : 'text-gray-300'" @click="toggleCollection">
             <component :is="isCollected ? StarFilled : Star" />
         </el-icon>
      </div>
    </div>

    <div v-if="poem" class="p-6 max-w-4xl mx-auto">
       
       <!-- Main Card with Tabs -->
       <div class="bg-white dark:bg-gray-800 rounded-3xl shadow-xl shadow-pink-100/50 dark:shadow-none border border-pink-50 dark:border-gray-700 overflow-hidden">
          
          <!-- Poem Header (Title, Author) -->
          <div class="p-8 pb-4 text-center border-b border-gray-50 dark:border-gray-700">
             <div class="flex justify-between items-start mb-2">
                <div class="text-left">
                   <h2 class="text-3xl font-bold text-gray-800 dark:text-gray-100 tracking-wide mb-1">{{ poem.title_cns }}</h2>
                   <div class="text-sm text-gray-500 dark:text-gray-400 flex items-center gap-2">
                      <span>{{ poem.dynasty_cns }} · {{ poem.author_cns }}</span>
                   </div>
                </div>
                <!-- Status Tags -->
                 <div class="flex flex-col items-end gap-1">
                     <span class="bg-gray-100 dark:bg-gray-700 text-gray-500 text-xs px-2 py-0.5 rounded-full">{{ getPoemTag(poem) || '古诗' }}</span>
                     <span class="text-xs px-2 py-0.5 rounded-full" :class="studyRecord?.isMastered ? 'bg-green-100 text-green-600' : 'bg-blue-50 text-blue-500'">
                         {{ studyRecord?.isMastered ? '已背' : '未开始' }}
                     </span>
                 </div>
             </div>
          </div>

          <!-- Tabs -->
          <div class="flex border-b border-gray-100 dark:border-gray-700 bg-pink-50/30 dark:bg-gray-800">
            <button 
              v-for="tab in ['原文', '译文', '赏析', '趣事']" 
              :key="tab"
              @click="activeTab = tab as any"
              class="flex-1 py-3 text-sm font-medium transition-all relative"
              :class="activeTab === tab ? 'text-pink-600 bg-white dark:bg-gray-700 shadow-sm rounded-t-lg mx-1 mt-1' : 'text-gray-500 hover:text-gray-700 dark:text-gray-400'"
            >
              {{ tab }}
            </button>
          </div>
          
          <!-- Content Area -->
          <div class="p-8 min-h-[300px] flex flex-col justify-center">
             
             <!-- Tab: Original (Poem Content) -->
             <div v-if="activeTab === '原文'" class="space-y-6">
                <div class="flex flex-col items-center gap-6">
                   <div v-for="(line, idx) in formattedLines" :key="idx" class="flex flex-wrap justify-center gap-x-1 gap-y-2">
                      <!-- Character + Pinyin Pair -->
                      <div v-for="(char, cIdx) in line.chars" :key="cIdx" class="flex flex-col items-center">
                         <span v-if="showPinyin" class="text-xs text-gray-500 font-light h-4">{{ line.pinyins[cIdx] || '' }}</span>
                         <span class="text-2xl text-gray-800 dark:text-gray-200 font-serif">{{ char }}</span>
                      </div>
                      <!-- Add comma/punctuation handling if needed, but simple char split usually works for display -->
                   </div>
                </div>

                <!-- Audio Player -->
                <div class="flex flex-col items-center gap-2 mt-8">
                   <el-button 
                     @click="togglePlay"
                     class="!h-10 !px-4 !rounded-full shadow-sm"
                     :type="isPlaying ? 'primary' : 'default'"
                   >
                     <el-icon class="mr-1"><VideoPlay v-if="!isPlaying" /><VideoPause v-else /></el-icon>
                     朗读
                   </el-button>
                </div>
             </div>

             <!-- Other Tabs -->
             <div v-else>
               <div v-if="isLoadingAI" class="flex flex-col items-center justify-center py-12 text-gray-400">
                 <div class="w-12 h-12 border-4 border-pink-200 border-t-pink-500 rounded-full animate-spin mb-4"></div>
                 <span class="animate-pulse">AI 正在生成{{ activeTab }}...</span>
               </div>
               
               <div v-else-if="getContent(activeTab)" class="prose dark:prose-invert max-w-none text-gray-700 dark:text-gray-300 whitespace-pre-wrap leading-relaxed text-justify px-4">
                 {{ getContent(activeTab) }}
               </div>
               
               <div v-else class="text-center py-12">
                 <el-icon :size="40" class="text-gray-300 mb-4"><MagicStick /></el-icon>
                 <p class="text-gray-400 mb-6 text-sm">暂无{{ activeTab }}内容</p>
                 <el-button type="primary" plain round class="!px-6 !border-pink-300 !text-pink-500 hover:!bg-pink-50" @click="generateContent(activeTab)">
                   <el-icon class="mr-2"><MagicStick /></el-icon>
                   AI 生成{{ activeTab }}
                 </el-button>
               </div>
             </div>

          </div>
       </div>

       <!-- Footer Actions -->
       <div class="flex justify-between gap-4 mt-6">
           <el-button 
             class="flex-1 !h-12 !text-base !rounded-xl !border-none"
             color="#fce7f3" 
             style="color: #db2777"
             @click="router.push(`/homework/chinese/poetry/recite/${poem.id}`)"
           >
             <el-icon class="mr-2"><Reading /></el-icon>
             背诵打卡
           </el-button>
           <el-button 
             class="flex-1 !h-12 !text-base !rounded-xl !border-none bg-blue-100 text-blue-600 hover:bg-blue-200" 
             @click="router.push(`/homework/chinese/poetry/dictate/${poem.id}`)"
           >
             <el-icon class="mr-2"><EditPen /></el-icon>
             默写打卡
           </el-button>
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
import { ArrowLeft, VideoPlay, VideoPause, Loading, MagicStick, Star, StarFilled, Reading, EditPen } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
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

const activeTab = ref<'原文' | '译文' | '赏析' | '趣事'>('原文')
const isLoadingAI = ref(false)
const isPlaying = ref(false)
const showPinyin = ref(false) // Default to false
const audioPlayer = ref<HTMLAudioElement | null>(null)
const isCollected = computed(() => poem.value ? poetryStore.collections.includes(poem.value.id) : false)
const studyRecord = computed(() => poem.value ? poetryStore.getRecord(poem.value.id) : null)

// Helper to determine tag based on ID or content (Mock logic)
const getPoemTag = (p: Poem) => {
    if (p.cat_cns) return p.cat_cns
    return '古诗'
}

// Compute aligned lines for display
const formattedLines = computed(() => {
    if (!poem.value) return []
    
    return poem.value.paragraphs_cns.map((line, idx) => {
        const chars = line.split('')
        let pinyins: string[] = []
        
        if (poem.value?.paragraphs_py1 && poem.value.paragraphs_py1[idx]) {
            // Simple space split. 
            // Warning: This assumes pinyin string is perfectly space-separated matching chars.
            // Punctuation in pinyin string (e.g., "，") needs to be handled if it is separated by space.
            // Based on JSON check: "gū ... lái ，" (space before comma). 
            // "孤...来，" (comma is char).
            // So split(' ') should yield same length as split('').
            const pyLine = poem.value.paragraphs_py1[idx]
            pinyins = pyLine.split(' ').filter(s => s.trim() !== '') // remove empty splits
        }
        
        // Safety: ensure lengths match or pad/truncate
        // If mismatch, we might just show empty pinyin for some chars
        return {
            chars,
            pinyins
        }
    })
})

onMounted(() => {
  poetryStore.init()
  if (!poem.value) {
    if (!poetryStore.getPoemById(poemId.value)) {
       ElMessage.error('古诗不存在')
       router.back()
    }
  }
})

watch(activeTab, (newVal) => {
    if (newVal !== '原文') {
        checkAndGenerateContent(newVal)
    }
})

const checkAndGenerateContent = (type: string) => {
    if (type === '原文') return
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

onUnmounted(() => {
  stopAudio()
})

const togglePlay = async () => {
  if (isPlaying.value) {
    stopAudio()
  } else {
    if (!poem.value) return
    if (poem.value.audio_url) {
        playAudio(poem.value.audio_url)
        return
    }
    const text = poem.value.title_cns + '。' + poem.value.author_cns + '。' + poem.value.paragraphs_cns.join('。')
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
        } catch {
            ElMessage.error('语音生成失败')
        }
    } else {
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
