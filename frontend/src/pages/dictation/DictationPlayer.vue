<template>
  <div class="min-h-screen bg-gray-900 text-white flex flex-col items-center justify-center relative overflow-hidden">
    <!-- Top Bar -->
    <div class="absolute top-0 left-0 w-full p-4 flex justify-between items-center z-10">
      <el-icon :size="24" class="cursor-pointer" @click="goBack"><ArrowLeft /></el-icon>
      <div class="text-sm opacity-70">{{ settings.play_mode === 'read' ? '朗读模式' : '听写模式' }}</div>
      <div class="w-6"></div>
    </div>

    <!-- Main Display -->
    <div class="flex-1 flex flex-col items-center justify-center w-full px-6 space-y-8">
      <div class="font-bold text-center transition-all duration-500 px-2 break-words max-w-full" :class="[{'opacity-0 blur-sm': !showWord && isPlaying}, wordClass]">
        {{ displayWord || '准备开始' }}
      </div>
      
      <div class="text-xl opacity-60" v-if="playlist.length > 0">
        {{ Math.min(currentIndex + 1, playlist.length) }} / {{ playlist.length }}
      </div>

      <div class="text-sm opacity-50 mt-2">
        已进行: {{ formatTime(elapsedTime) }}
      </div>

      <!-- Progress -->
      <div class="w-full max-w-md h-2 bg-gray-700 rounded-full overflow-hidden mt-8">
        <div class="h-full bg-blue-500 transition-all duration-300" :style="{ width: progressPercent + '%' }"></div>
      </div>
    </div>

    <!-- Controls -->
    <div class="pb-12 w-full max-w-md px-6 flex flex-col gap-6">
      
      <!-- Tools Row -->
      <div class="flex justify-center gap-4">
        <el-button 
          :type="settings.order_mode === 'random' ? 'warning' : 'info'" 
          plain 
          round 
          size="default" 
          @click="toggleOrder"
        >
          <el-icon class="mr-1"><Sort /></el-icon>
          {{ settings.order_mode === 'random' ? '乱序播放' : '顺序播放' }}
        </el-button>

        <el-button 
          type="danger" 
          plain 
          round 
          size="default" 
          v-if="currentWord"
          @click="markMistake"
        >
          标记难词
        </el-button>
      </div>

      <div class="flex items-center justify-between">
        <el-icon :size="36" class="cursor-pointer hover:text-blue-400 transition" @click="prev"><ArrowLeftBold /></el-icon>
        
        <div 
          class="w-20 h-20 rounded-full bg-blue-600 flex items-center justify-center cursor-pointer hover:bg-blue-500 transition shadow-lg shadow-blue-600/40"
          @click="togglePlay"
        >
          <el-icon :size="40" v-if="!isPlaying"><VideoPlay /></el-icon>
          <el-icon :size="40" v-else><VideoPause /></el-icon>
        </div>

        <el-icon :size="36" class="cursor-pointer hover:text-blue-400 transition" @click="next"><ArrowRightBold /></el-icon>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, ArrowLeftBold, ArrowRightBold, VideoPlay, VideoPause, Sort } from '@element-plus/icons-vue'
import { dictationApi, type DictationSettings } from '@/services/dictation'
import { ElMessage } from 'element-plus'
import { speak as speakTTS, stop as stopTTS } from '@/utils/speech'

const router = useRouter()

// State
const content = ref('')
const playlist = ref<string[]>([])
const originalPlaylist = ref<string[]>([])
const currentIndex = ref(0)
const isPlaying = ref(false)
const startTime = ref(Date.now())
const elapsedTime = ref(0)
const elapsedTimer = ref<any>(null)
const localMistakeCount = ref(0)
const settings = ref<DictationSettings>({
  user_id: 0,
  split_rule: 'newline',
  play_mode: 'dictate',
  order_mode: 'sequential',
  repeat_count: 1,
  interval_seconds: 3,
  voice_type: '',
  zh_voice_type: '',
  en_voice_type: '',
  speed: 1.0
})

// Logic State
const currentRepeat = ref(0)
const timer = ref<any>(null)
const isWaiting = ref(false)

const currentWord = computed(() => playlist.value[currentIndex.value] || '')

const displayWord = computed(() => {
  const w = currentWord.value
  if (!w) return ''
  if (w.startsWith('http://') || w.startsWith('https://')) {
    try {
      const decoded = decodeURIComponent(w)
      const parts = decoded.split('/')
      return parts[parts.length - 1]
    } catch {
      return w
    }
  }
  return w
})

const wordClass = computed(() => {
  const len = displayWord.value.length
  if (len > 30) return 'text-lg'
  if (len > 20) return 'text-xl'
  if (len > 10) return 'text-3xl'
  return 'text-6xl'
})

const progressPercent = computed(() => {
  if (playlist.value.length === 0) return 0
  return ((currentIndex.value) / playlist.value.length) * 100
})
// In Dictation mode, maybe we hide the word? 
// PRD says "Dictation Mode... for practice". Usually you don't see the word.
// But for "Read Mode", you see it.
// Let's hide it in Dictation Mode unless paused? Or just hide it.
// Let's allow user to peek if they pause.
const showWord = computed(() => {
  if (settings.value.play_mode === 'read') return true
  return !isPlaying.value // Show when paused/stopped
})

onMounted(async () => {
  const stored = localStorage.getItem('dictation_current_content')
  if (!stored) {
    ElMessage.warning('没有听写内容')
    router.back()
    return
  }
  content.value = stored

  // Load settings
  try {
      const res = await dictationApi.getSettings()
    const data = (res as any).data || (res as any)
    if (data) settings.value = { ...settings.value, ...data }
  } catch {}

  // Prepare playlist
  let items: string[] = []
  if (settings.value.split_rule === 'space') {
    items = content.value.split(/[\s\n]+/).filter(t => t.trim())
  } else {
    items = content.value.split('\n').filter(t => t.trim())
  }
  originalPlaylist.value = [...items]

  if (settings.value.order_mode === 'random') {
    // Shuffle
    shuffle(items)
  }
  playlist.value = items
})

onUnmounted(() => {
  stop()
})

function shuffle(array: string[]) {
  for (let i = array.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [array[i], array[j]] = [array[j], array[i]];
  }
}

function toggleOrder() {
  stop()
  if (settings.value.order_mode === 'sequential') {
    settings.value.order_mode = 'random'
    const temp = [...originalPlaylist.value]
    shuffle(temp)
    playlist.value = temp
    ElMessage.info('已切换为乱序播放')
  } else {
    settings.value.order_mode = 'sequential'
    playlist.value = [...originalPlaylist.value]
    ElMessage.info('已切换为顺序播放')
  }
  currentIndex.value = 0
  currentRepeat.value = 0
}

function goBack() {
  stop()
  router.back()
}

function stop() {
  isPlaying.value = false
  stopTTS()
  clearTimeout(timer.value)
  clearInterval(elapsedTimer.value)
}

function togglePlay() {
  if (isPlaying.value) {
    stop()
  } else {
    isPlaying.value = true
    startTime.value = Date.now() - elapsedTime.value * 1000
    elapsedTimer.value = setInterval(() => {
      elapsedTime.value = Math.floor((Date.now() - startTime.value) / 1000)
    }, 1000)
    playLoop()
  }
}

function formatTime(sec: number) {
  const m = Math.floor(sec / 60)
  const s = sec % 60
  return `${m}分${s}秒`
}

function playLoop() {
  if (!isPlaying.value) return
  if (currentIndex.value >= playlist.value.length) {
    isPlaying.value = false
    ElMessage.success('听写完成')
    // Save history
    const duration = Math.floor((Date.now() - startTime.value) / 1000)
    dictationApi.addHistory({
      content_snapshot: content.value.slice(0, 200),
      duration_seconds: duration,
      mistake_count: localMistakeCount.value
    })
    return
  }

  const word = playlist.value[currentIndex.value]
  speak(word, () => {
    // After speak
    currentRepeat.value++
    if (currentRepeat.value < settings.value.repeat_count) {
      const gap = (Number(settings.value.interval_seconds) || 0) * 1000
      timer.value = setTimeout(() => {
        playLoop()
      }, gap)
    } else {
      currentRepeat.value = 0
      const gap = (Number(settings.value.interval_seconds) || 0) * 1000
      isWaiting.value = true
      timer.value = setTimeout(() => {
        isWaiting.value = false
        currentIndex.value++
        playLoop()
      }, gap)
    }
  })
}

const audioRef = ref<HTMLAudioElement | null>(null)

async function speak(text: string, onEnd: () => void) {
  stopTTS()
  
  // Check if text is a URL
  if (text.startsWith('http://') || text.startsWith('https://')) {
    if (!audioRef.value) {
      audioRef.value = new Audio(text)
    } else {
      audioRef.value.src = text
    }
    audioRef.value.playbackRate = settings.value.speed
    audioRef.value.onended = onEnd
    audioRef.value.onerror = () => {
      ElMessage.error('音频播放失败')
      onEnd()
    }
    audioRef.value.play().catch(() => {
      ElMessage.error('音频播放被阻止')
      onEnd()
    })
    return
  }

  const isChinese = /[\u4e00-\u9fa5]/.test(text)
  let targetVoiceURI = isChinese ? settings.value.zh_voice_type : settings.value.en_voice_type
  if (!targetVoiceURI) targetVoiceURI = settings.value.voice_type

  await speakTTS(text, {
    voiceURI: targetVoiceURI || undefined,
    rate: settings.value.speed,
    lang: isChinese ? 'zh-CN' : 'en-US',
    onEnd: onEnd,
    onError: onEnd
  })
}

function prev() {
  stop()
  if (currentIndex.value > 0) {
    currentIndex.value--
    currentRepeat.value = 0
  }
}

function next() {
  stop()
  if (currentIndex.value < playlist.value.length - 1) {
    currentIndex.value++
    currentRepeat.value = 0
  }
}

async function markMistake() {
  if (!currentWord.value) return
  localMistakeCount.value++
  try {
    await dictationApi.addMistake({ word: currentWord.value, context: '听写练习' })
    ElMessage.success(`已将“${currentWord.value}”加入难点收藏`)
  } catch {}
}
</script>
