<template>
  <div class="min-h-screen bg-gradient-to-b from-purple-50 to-pink-50 dark:from-gray-900 dark:to-gray-800 pb-20 font-sans">
    <div class="sticky top-0 z-20 bg-white/80 dark:bg-gray-800/80 backdrop-blur-md border-b border-purple-100 dark:border-gray-700 shadow-sm p-4 flex items-center justify-between">
       <div class="flex items-center gap-2">
           <el-icon :size="22" class="cursor-pointer text-gray-600 dark:text-gray-300 hover:text-purple-500 transition" @click="router.back()"><ArrowLeft /></el-icon>
           <span class="font-bold text-lg text-purple-700 dark:text-purple-300">诗词闯关</span>
       </div>
       <div class="bg-yellow-100 text-yellow-700 px-3 py-1 rounded-full text-xs font-bold flex items-center gap-1 shadow-sm border border-yellow-200">
           <el-icon><Coin /></el-icon> {{ score }} 分
       </div>
    </div>

    <div class="p-6 max-w-md mx-auto" v-if="!isPlaying">
        <!-- Level Select -->
        <div class="text-center mb-8 animate-fade-in-down">
            <h1 class="text-3xl font-black text-transparent bg-clip-text bg-gradient-to-r from-purple-600 to-pink-500 mb-2">挑战开始</h1>
            <p class="text-gray-500 text-sm">选择难度，挑战你的诗词储备量！</p>
        </div>

        <div class="space-y-4 animate-fade-in-up">
            <div 
                v-for="(level, idx) in levels" 
                :key="idx"
                class="bg-white dark:bg-gray-800 p-4 rounded-2xl shadow-sm border-2 border-transparent hover:border-purple-300 cursor-pointer transition-all transform hover:scale-105 active:scale-95 group relative overflow-hidden"
                @click="startGame(level.type)"
            >
                <div class="absolute right-0 top-0 w-24 h-24 bg-gradient-to-bl from-gray-50 to-transparent rounded-bl-full opacity-50 group-hover:from-purple-50 transition"></div>
                <div class="flex items-center gap-4 relative z-10">
                    <div class="w-14 h-14 rounded-2xl flex items-center justify-center text-3xl shadow-inner" :class="level.bgClass">
                        {{ level.icon }}
                    </div>
                    <div>
                        <div class="font-bold text-lg text-gray-800 dark:text-gray-100 group-hover:text-purple-600 transition">{{ level.title }}</div>
                        <div class="text-xs text-gray-500 dark:text-gray-400">{{ level.desc }}</div>
                    </div>
                </div>
            </div>
        </div>
        
        <!-- History Best -->
        <div class="mt-10 bg-white/50 dark:bg-gray-800/50 rounded-xl p-4 text-center border border-purple-50 dark:border-gray-700">
             <div class="text-xs text-gray-400 mb-1">历史最高分</div>
             <div class="text-2xl font-black text-purple-400">{{ bestScore }}</div>
        </div>
    </div>

    <!-- Game Interface -->
    <div v-else class="p-6 max-w-md mx-auto h-full flex flex-col min-h-[80vh]">
        <!-- Progress & Lives -->
        <div class="flex justify-between items-center mb-6">
            <div class="bg-white dark:bg-gray-800 px-3 py-1 rounded-full text-xs font-bold text-gray-500 shadow-sm border border-gray-100 dark:border-gray-700">
              第 {{ currentQuestionIndex + 1 }} / {{ totalQuestions }} 题
            </div>
            <div class="flex gap-1 bg-white dark:bg-gray-800 px-2 py-1 rounded-full shadow-sm border border-gray-100 dark:border-gray-700">
                <el-icon v-for="i in lives" :key="i" class="text-red-500 animate-pulse"><StarFilled /></el-icon>
                <el-icon v-for="i in (3 - lives)" :key="'lost-'+i" class="text-gray-300"><Star /></el-icon>
            </div>
        </div>

        <!-- Question Card -->
        <div class="bg-white dark:bg-gray-800 rounded-3xl p-8 shadow-xl shadow-purple-100/50 dark:shadow-none border border-purple-50 dark:border-purple-900/50 mb-6 flex-1 flex flex-col justify-center text-center relative overflow-hidden transition-all duration-300">
             <!-- Timer Bar -->
             <div class="absolute top-0 left-0 h-1.5 bg-gradient-to-r from-purple-400 to-pink-500 transition-all duration-1000 ease-linear" :style="{width: (timeLeft / 15 * 100) + '%'}"></div>
             
             <!-- Decorative -->
             <div class="absolute -bottom-10 -right-10 w-40 h-40 bg-purple-50 dark:bg-purple-900/20 rounded-full blur-3xl pointer-events-none"></div>

             <div class="text-xs text-purple-500 font-bold mb-6 tracking-widest uppercase bg-purple-50 dark:bg-purple-900/30 inline-block mx-auto px-3 py-1 rounded-full">{{ currentQuestion.typeLabel }}</div>
             
             <div class="text-xl md:text-2xl font-bold text-gray-800 dark:text-gray-100 mb-8 font-serif leading-relaxed px-2">
                 {{ currentQuestion.question }}
             </div>

             <div class="grid grid-cols-1 gap-3 w-full relative z-10">
                 <button 
                    v-for="(opt, idx) in currentQuestion.options" 
                    :key="idx"
                    class="py-4 px-6 rounded-2xl border-2 font-medium transition-all duration-200 text-sm md:text-base shadow-sm hover:shadow-md active:scale-98 flex items-center justify-center min-h-[60px]"
                    :class="getOptionClass(opt)"
                    @click="handleAnswer(opt)"
                    :disabled="answered"
                 >
                     {{ opt }}
                 </button>
             </div>
        </div>

        <div v-if="answered" class="text-center animate-fade-in-up pb-4">
            <div v-if="isCorrect" class="text-green-500 font-bold text-xl mb-4 flex items-center justify-center gap-2">
              <el-icon><CircleCheckFilled /></el-icon> 回答正确！
            </div>
            <div v-else class="text-red-500 font-bold text-lg mb-4 flex flex-col items-center">
               <div class="flex items-center gap-2 mb-1"><el-icon><CircleCloseFilled /></el-icon> 回答错误</div>
               <span class="text-sm font-normal text-gray-500">正确答案：{{ currentQuestion.answer }}</span>
            </div>
            
            <el-button 
              type="primary" 
              size="large" 
              round 
              class="w-full !h-12 !text-lg !rounded-full shadow-lg shadow-purple-200 hover:shadow-purple-300 transition-all bg-gradient-to-r from-purple-500 to-pink-500 border-none" 
              @click="nextQuestion"
            >
                {{ currentQuestionIndex < totalQuestions - 1 && lives > 0 ? '下一题' : '查看结果' }}
            </el-button>
        </div>
    </div>

    <!-- Result Modal -->
    <el-dialog 
      v-model="showResult" 
      title="挑战结束" 
      width="85%" 
      center 
      :show-close="false" 
      align-center
      class="rounded-3xl !p-0 overflow-hidden"
    >
        <div class="text-center py-8 px-4 bg-gradient-to-b from-white to-purple-50 dark:from-gray-800 dark:to-gray-900">
            <div class="text-7xl mb-6 animate-bounce">{{ score >= 80 ? '🏆' : (score >= 60 ? '👍' : '💪') }}</div>
            <div class="text-3xl font-black text-gray-800 dark:text-gray-100 mb-2">得分：<span class="text-purple-600">{{ score }}</span></div>
            <p class="text-gray-500 mb-8 bg-white dark:bg-gray-800 py-2 px-4 rounded-xl inline-block shadow-sm">
               共 {{ totalQuestions }} 题，答对 <span class="text-green-500 font-bold">{{ correctCount }}</span> 题
            </p>
            
            <div class="grid grid-cols-2 gap-4">
                <el-button size="large" round class="!h-12" @click="resetGame">返回首页</el-button>
                <el-button type="primary" size="large" round class="!h-12 bg-gradient-to-r from-purple-500 to-pink-500 border-none" @click="restart">再来一次</el-button>
            </div>
        </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { usePoetryStore } from '../stores/poetryStore'
import { ArrowLeft, Coin, Star, StarFilled, CircleCheckFilled, CircleCloseFilled } from '@element-plus/icons-vue'
import { sampleSize, shuffle, random } from 'lodash-es'
import type { Poem } from '../types'

const router = useRouter()
const store = usePoetryStore()

// State
const isPlaying = ref(false)
const score = ref(0)
const lives = ref(3)
const currentQuestionIndex = ref(0)
const totalQuestions = 10
const timeLeft = ref(15)
const timer = ref<any>(null)
const answered = ref(false)
const isCorrect = ref(false)
const showResult = ref(false)
const correctCount = ref(0)
const currentMode = ref('')

// Levels Config
const levels = [
    { type: 'easy', title: '初入江湖', desc: '看诗句猜作者或标题', icon: '🐣', bgClass: 'bg-green-100 text-green-600' },
    { type: 'medium', title: '小有所成', desc: '根据上句接下句', icon: '🐯', bgClass: 'bg-blue-100 text-blue-600' },
    { type: 'hard', title: '一代宗师', desc: '补全诗句中的缺字', icon: '🐉', bgClass: 'bg-purple-100 text-purple-600' }
]

interface Question {
    typeLabel: string
    question: string
    options: string[]
    answer: string
}

const questions = ref<Question[]>([])
const currentQuestion = computed(() => questions.value[currentQuestionIndex.value] || { typeLabel: '', question: '', options: [], answer: '' })

// Best Score
const bestScore = computed(() => {
    if (store.challengeRecords.length === 0) return 0
    return Math.max(...store.challengeRecords.map(r => r.score))
})

onMounted(() => {
    store.init()
})

onUnmounted(() => {
    clearInterval(timer.value)
})

const startGame = (mode: string) => {
    currentMode.value = mode
    score.value = 0
    lives.value = 3
    correctCount.value = 0
    currentQuestionIndex.value = 0
    generateQuestions(mode)
    isPlaying.value = true
    startTimer()
}

const generateQuestions = (mode: string) => {
    const pool = store.poems
    if (pool.length < 4) {
        // Fallback if not enough poems
        questions.value = [] // Should probably show error or handle gracefully
        return
    }

    const selectedPoems = sampleSize(pool, totalQuestions)
    questions.value = selectedPoems.map(poem => createQuestion(poem, mode, pool))
}

const createQuestion = (poem: Poem, mode: string, pool: Poem[]): Question => {
    if (mode === 'easy') {
        // Guess Title or Author
        const isTitle = Math.random() > 0.5
        const questionLine = poem.paragraphs_cns[0] || ''
        const answer = isTitle ? poem.title_cns : poem.author_cns
        
        // Distractors
        const others = sampleSize(pool.filter(p => p.id !== poem.id), 3)
        const distractors = others.map(p => isTitle ? p.title_cns : p.author_cns)
        
        return {
            typeLabel: isTitle ? '猜诗名' : '猜作者',
            question: `"${questionLine}" 出自哪里？`,
            options: shuffle([answer, ...distractors]),
            answer
        }
    } else if (mode === 'medium') {
        // Next Line
        // Pick a pair of lines if possible, usually lines are even. 
        // paragraphs_cns might be ["床前明月光，", "疑是地上霜。"]
        // We need to clean punctuation for display maybe?
        const lineCount = poem.paragraphs_cns.length
        const idx = random(0, lineCount - 2) // Ensure there is a next line
        const qLine = poem.paragraphs_cns[idx]
        const aLine = poem.paragraphs_cns[idx + 1]
        
        // Distractors: Random lines from other poems
        const others = sampleSize(pool.filter(p => p.id !== poem.id), 3)
        const distractors = others.map(p => p.paragraphs_cns[0]) // Just take first line of others
        
        return {
            typeLabel: '接下句',
            question: `${qLine}`,
            options: shuffle([aLine, ...distractors]),
            answer: aLine
        }
    } else {
        // Hard: Fill in missing character
        // Pick a line
        const line = sampleSize(poem.paragraphs_cns, 1)[0]
        // Pick a char to hide (not punctuation)
        const chars = line.split('').filter(c => !['，', '。', '？', '！'].includes(c))
        const hiddenChar = sampleSize(chars, 1)[0]
        const questionText = line.replace(hiddenChar, '___')
        
        // Distractors: Random chars
        // A simple way is to pick random chars from other lines
        const allChars = pool.flatMap(p => p.paragraphs_cns.join('').split('')).filter(c => !['，', '。', '？', '！'].includes(c))
        const distractors = sampleSize(allChars.filter(c => c !== hiddenChar), 3)
        
        return {
            typeLabel: '填缺字',
            question: questionText,
            options: shuffle([hiddenChar, ...distractors]),
            answer: hiddenChar
        }
    }
}

const startTimer = () => {
    clearInterval(timer.value)
    timeLeft.value = 15
    timer.value = setInterval(() => {
        timeLeft.value--
        if (timeLeft.value <= 0) {
            handleTimeOut()
        }
    }, 1000)
}

const handleTimeOut = () => {
    clearInterval(timer.value)
    answered.value = true
    isCorrect.value = false
    lives.value--
    // Auto next after delay if desired, or wait for user
    if (lives.value <= 0) {
        // Game Over logic handled in nextQuestion check
    }
}

const handleAnswer = (opt: string) => {
    if (answered.value) return
    clearInterval(timer.value)
    answered.value = true
    
    if (opt === currentQuestion.value.answer) {
        isCorrect.value = true
        score.value += 10
        correctCount.value++
    } else {
        isCorrect.value = false
        lives.value--
    }
}

const nextQuestion = () => {
    if (lives.value <= 0) {
        endGame()
        return
    }
    
    if (currentQuestionIndex.value < totalQuestions - 1) {
        currentQuestionIndex.value++
        answered.value = false
        isCorrect.value = false
        startTimer()
    } else {
        endGame()
    }
}

const endGame = () => {
    clearInterval(timer.value)
    store.saveChallengeResult(score.value, currentMode.value)
    showResult.value = true
}

const getOptionClass = (opt: string) => {
    if (!answered.value) return 'bg-white dark:bg-gray-700 border-gray-100 dark:border-gray-600 hover:border-purple-300 dark:hover:border-purple-500 text-gray-700 dark:text-gray-200'
    
    if (opt === currentQuestion.value.answer) {
        return 'bg-green-100 dark:bg-green-900/50 border-green-500 text-green-700 dark:text-green-300'
    }
    
    // If this option was selected and wrong
    // Note: We don't track which one was clicked in this simple version unless we add a state.
    // But typically we show the correct one green, and if the user clicked wrong, we might want to show that as red.
    // For simplicity, just highlight correct answer green. 
    // To highlight wrong selection, we'd need `selectedOption` state.
    
    return 'bg-white dark:bg-gray-700 border-gray-100 dark:border-gray-600 opacity-50'
}

const resetGame = () => {
    showResult.value = false
    isPlaying.value = false
}

const restart = () => {
    showResult.value = false
    startGame(currentMode.value)
}
</script>

<style scoped>
.animate-fade-in-down {
    animation: fadeInDown 0.5s ease-out;
}
.animate-fade-in-up {
    animation: fadeInUp 0.5s ease-out;
}

@keyframes fadeInDown {
    from { opacity: 0; transform: translateY(-20px); }
    to { opacity: 1; transform: translateY(0); }
}

@keyframes fadeInUp {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
}
</style>
