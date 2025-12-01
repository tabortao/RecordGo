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
      <div class="text-6xl font-bold text-center transition-all duration-500" :class="{'opacity-0 blur-sm': !showWord && isPlaying}">
        {{ currentWord || '准备开始' }}
      </div>
      
      <div class="text-xl opacity-60" v-if="playlist.length > 0">
        {{ currentIndex + 1 }} / {{ playlist.length }}
      </div>

      <!-- Progress -->
      <div class="w-full max-w-md h-2 bg-gray-700 rounded-full overflow-hidden">
        <div class="h-full bg-blue-500 transition-all duration-300" :style="{ width: progressPercent + '%' }"></div>
      </div>
    </div>

    <!-- Controls -->
    <div class="pb-12 w-full max-w-md px-6 flex flex-col gap-6">
      
      <!-- Mark Mistake Button (only visible when playing/paused on a word) -->
      <div class="flex justify-center">
        <el-button 
          type="danger" 
          plain 
          round 
          size="large" 
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
import { ArrowLeft, ArrowLeftBold, ArrowRightBold, VideoPlay, VideoPause } from '@element-plus/icons-vue'
import { dictationApi, type DictationSettings } from '@/services/dictation'
import { ElMessage } from 'element-plus'

const router = useRouter()

// State
const content = ref('')
const playlist = ref<string[]>([])
const currentIndex = ref(0)
const isPlaying = ref(false)
const settings = ref<DictationSettings>({
  user_id: 0,
  split_rule: 'newline',
  play_mode: 'dictate',
  order_mode: 'sequential',
  repeat_count: 1,
  interval_seconds: 3,
  voice_type: '',
  speed: 1.0
})

// Logic State
const currentRepeat = ref(0)
const timer = ref<any>(null)
const isWaiting = ref(false)

const currentWord = computed(() => playlist.value[currentIndex.value] || '')
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

  if (settings.value.order_mode === 'random') {
    // Shuffle
    for (let i = items.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [items[i], items[j]] = [items[j], items[i]];
    }
  }
  playlist.value = items
})

onUnmounted(() => {
  stop()
})

function goBack() {
  stop()
  router.back()
}

function stop() {
  isPlaying.value = false
  window.speechSynthesis.cancel()
  clearTimeout(timer.value)
}

function togglePlay() {
  if (isPlaying.value) {
    stop()
  } else {
    isPlaying.value = true
    playLoop()
  }
}

function playLoop() {
  if (!isPlaying.value) return
  if (currentIndex.value >= playlist.value.length) {
    isPlaying.value = false
    ElMessage.success('听写完成')
    // Save history
    dictationApi.addHistory({
      content_snapshot: content.value.slice(0, 200),
      duration_seconds: 0, // TODO: track time
      mistake_count: 0 // TODO: track mistakes count
    })
    return
  }

  const word = playlist.value[currentIndex.value]
  speak(word, () => {
    // After speak
    currentRepeat.value++
    if (currentRepeat.value < settings.value.repeat_count) {
       // Repeat same word immediately? Usually there is a small gap.
       timer.value = setTimeout(() => {
         playLoop() // Recursively call playLoop but without incrementing index
       }, 500)
    } else {
      // Move to next word after interval
      currentRepeat.value = 0
      const gap = settings.value.play_mode === 'dictate' ? settings.value.interval_seconds * 1000 : 500
      isWaiting.value = true
      timer.value = setTimeout(() => {
        isWaiting.value = false
        currentIndex.value++
        playLoop()
      }, gap)
    }
  })
}

function speak(text: string, onEnd: () => void) {
  window.speechSynthesis.cancel()
  const u = new SpeechSynthesisUtterance(text)
  u.rate = settings.value.speed
  if (settings.value.voice_type) {
    const v = window.speechSynthesis.getVoices().find(x => x.voiceURI === settings.value.voice_type)
    if (v) u.voice = v
  }
  if (!u.voice) {
      // Try to match language if no specific voice
    if (/[\u4e00-\u9fa5]/.test(text)) {
        u.lang = 'zh-CN'
    } else {
        u.lang = 'en-US'
    }
  }
  
  u.onend = onEnd
  u.onerror = onEnd // handle error as end
  window.speechSynthesis.speak(u)
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
  try {
    await dictationApi.addMistake({ word: currentWord.value, context: '听写练习' })
    ElMessage.success(`已将“${currentWord.value}”加入错题本`)
  } catch {}
}
</script>
