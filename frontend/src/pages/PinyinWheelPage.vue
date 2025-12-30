<template>
  <div class="min-h-screen bg-green-50 dark:bg-green-900/20 flex flex-col relative overflow-hidden">
    <!-- 顶部导航 -->
    <div class="p-4 z-10 flex items-center justify-between">
      <div class="w-10 h-10 bg-white dark:bg-gray-800 rounded-full shadow flex items-center justify-center cursor-pointer hover:scale-105 transition" @click="router.back()">
        <el-icon :size="20" class="text-green-600"><ArrowLeft /></el-icon>
      </div>
      <!-- 移除了之前的固定居中标题，改为在 Flex 布局中处理 -->
      <div class="w-10"></div>
    </div>

    <!-- 游戏区域 -->
    <div class="flex-1 w-full max-w-6xl mx-auto p-4 flex flex-col lg:flex-row items-center justify-center gap-8 lg:gap-16">
      
      <!-- 左侧：转盘区 -->
      <div class="relative flex-shrink-0">
        <div class="relative w-[90vw] h-[90vw] max-w-[500px] max-h-[500px] flex items-center justify-center aspect-square">
          <!-- 指针 (固定在上方) - 独立于转盘容器，绝对定位 -->
          <div class="absolute top-0 left-1/2 -translate-x-1/2 -translate-y-4 z-40 pointer-events-none">
            <div class="w-0 h-0 border-l-[16px] border-l-transparent border-r-[16px] border-r-transparent border-t-[32px] border-t-red-500 drop-shadow-md filter drop-shadow-lg"></div>
          </div>

          <!-- 外圈：声母 (放大尺寸) -->
          <div 
            class="absolute inset-0 rounded-full border-8 border-green-300 bg-green-100 shadow-xl transition-transform duration-[3000ms] ease-out"
            :style="{ transform: `rotate(${rotations.initial}deg)` }"
          >
            <div v-for="(item, idx) in initials" :key="item"
              class="absolute top-0 left-1/2 -ml-[30px] w-[60px] h-1/2 origin-bottom text-center"
              :style="{ transform: `rotate(${idx * (360 / initials.length)}deg)` }"
            >
              <div class="absolute top-[6%] left-1/2 -translate-x-1/2 flex flex-row items-center justify-center gap-0.5 sm:gap-1 min-w-max">
                <span class="text-lg sm:text-2xl font-bold text-green-800">{{ item }}</span>
                <span class="text-sm sm:text-lg">{{ initialsData[item]?.icon }}</span>
              </div>
            </div>
          </div>

          <!-- 中圈：韵母 -->
          <div 
            class="absolute top-[18%] left-[18%] right-[18%] bottom-[18%] rounded-full border-8 border-yellow-200 bg-yellow-50 shadow-lg transition-transform duration-[3000ms] ease-out"
            :style="{ transform: `rotate(${rotations.final}deg)` }"
          >
            <div v-for="(item, idx) in finals" :key="item"
              class="absolute top-0 left-1/2 -ml-[15px] w-[30px] h-1/2 origin-bottom pt-2 text-center font-bold text-orange-700"
              :style="{ transform: `rotate(${idx * (360 / finals.length)}deg)` }"
            >
              <span class="block transform -rotate-90 mt-1 text-sm sm:text-base">{{ item }}</span>
            </div>
          </div>

          <!-- 内圈：声调 -->
          <div 
            class="absolute top-[38%] left-[38%] right-[38%] bottom-[38%] rounded-full border-8 border-pink-200 bg-pink-50 shadow-md transition-transform duration-[3000ms] ease-out flex items-center justify-center"
            :style="{ transform: `rotate(${rotations.tone}deg)` }"
          >
            <div v-for="(item, idx) in tones" :key="item"
              class="absolute top-0 left-1/2 -ml-[10px] w-[20px] h-1/2 origin-bottom pt-1 text-center font-bold text-pink-600"
              :style="{ transform: `rotate(${idx * 90}deg)` }"
            >
              <span class="block mt-1 text-sm sm:text-base">{{ item }}声</span>
            </div>
            <!-- 中心装饰 -->
            <div class="w-16 h-16 sm:w-20 sm:h-20 rounded-full bg-white shadow-inner flex items-center justify-center z-10 border-4 border-pink-100">
              <span class="text-3xl sm:text-4xl">🐼</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧：结果与控制区 -->
      <div class="flex flex-col items-center gap-8 z-10 max-w-md w-full">
        
        <h1 class="text-3xl lg:text-4xl font-bold text-green-700 dark:text-green-400 tracking-wider drop-shadow-sm text-center">
          拼音大转盘
        </h1>

        <!-- 结果展示区 -->
        <div class="relative bg-white dark:bg-gray-800 rounded-3xl p-8 shadow-xl border-b-8 border-green-500 w-full text-center min-h-[220px] flex flex-col items-center justify-center transition-all transform hover:scale-[1.02] duration-300">
          <div v-if="isSpinning" class="flex flex-col items-center gap-3">
            <div class="w-12 h-12 border-4 border-green-200 border-t-green-500 rounded-full animate-spin"></div>
            <div class="text-green-600 font-bold text-xl">拼读中...</div>
          </div>
          <div v-else-if="currentResult" class="space-y-4 animate-fade-in">
            <!-- 朗读图标（右上角） -->
            <div class="absolute top-1 right-1">
              <el-icon :size="20" class="text-yellow-600 cursor-pointer hover:scale-110 transition" @click="speak"><Microphone /></el-icon>
            </div>
            <!-- 拼音公式 -->
            <div class="flex items-center justify-center gap-2 text-2xl sm:text-3xl font-mono text-gray-500 bg-gray-50 dark:bg-gray-700/50 px-4 py-2 rounded-full">
              <span class="text-green-700 font-bold">{{ currentResult.initial }}</span>
              <span>+</span>
              <span class="text-orange-600 font-bold">{{ currentResult.final }}</span>
              <span>+</span>
              <span class="text-pink-600 font-bold">{{ currentResult.tone }}声</span>
              <span>=</span>
              <span class="text-indigo-600 font-bold">{{ currentResult.pinyinWithTone }}</span>
            </div>
            
            <!-- 汉字展示 -->
            <div class="flex flex-col items-center">
              <div class="text-7xl sm:text-8xl font-black text-gray-800 dark:text-white drop-shadow-md my-2">
                {{ currentResult.word }}
              </div>
              <!-- 词组/例句 -->
              <div class="text-lg text-gray-500 bg-green-50 dark:bg-green-900/30 px-4 py-1 rounded-lg border border-green-100 dark:border-green-800">
                {{ currentResult.example }}
              </div>
            </div>
          </div>
          <div v-else class="text-gray-400 flex flex-col items-center gap-2">
            <span class="text-5xl mb-2">🎡</span>
            <span>点击“转动”开始学习吧！</span>
          </div>
        </div>

        <!-- 控制：改为 span 触发 -->
        <div class="flex justify-center gap-6 w-full">
          <span 
            class="inline-flex items-center justify-center select-none cursor-pointer bg-gradient-to-b from-green-400 to-green-600 text-white px-6 py-3 rounded-2xl shadow-lg border-b-4 border-green-700 active:border-b-0 active:translate-y-1 transition font-bold text-xl gap-2"
            @click="spin"
            :aria-disabled="isSpinning"
            :class="{ 'opacity-50 cursor-not-allowed': isSpinning }"
          >
            <el-icon :class="{ 'animate-spin': isSpinning }"><Refresh /></el-icon>
            {{ isSpinning ? '转动中...' : '转动' }}
          </span>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Refresh, Microphone } from '@element-plus/icons-vue'
import { initials, finals, tones, pinyinDict, initialsData } from '@/utils/pinyinData'

const router = useRouter()

// 状态
const isSpinning = ref(false)
const rotations = reactive({ initial: 0, final: 0, tone: 0 })
const currentResult = ref<{
  initial: string
  final: string
  tone: number
  pinyinWithTone: string
  word: string
  example: string
} | null>(null)

// 查找下一个有效组合
function findValidCombination(): { initialIdx: number, finalIdx: number, toneIdx: number, data: any } {
  // 随机尝试 200 次，找不到就返回默认
  for (let i = 0; i < 200; i++) {
    const iIdx = Math.floor(Math.random() * initials.length)
    const fIdx = Math.floor(Math.random() * finals.length)
    const tIdx = Math.floor(Math.random() * tones.length)

    const initial = initials[iIdx]
    const final = finals[fIdx]
    const tone = tones[tIdx]

    let key = initial + final

    // 特殊拼写规则处理
    // j, q, x, y 遇到 ü 时，写成 u (ju, qu, xu, yu)
    if (['j', 'q', 'x', 'y'].includes(initial)) {
      if (final === 'ü') {
        key = initial + 'u'
      } else if (final === 'u') {
        // j, q, x, y 不能与 u (乌) 相拼，直接视为无效
        continue
      }
    }

    if (pinyinDict[key] && pinyinDict[key][tone]) {
      return { initialIdx: iIdx, finalIdx: fIdx, toneIdx: tIdx, data: pinyinDict[key][tone] }
    }
  }
  // Fallback (ba 1)
  return { initialIdx: 0, finalIdx: 0, toneIdx: 0, data: pinyinDict['ba'][1] }
}

function spin() {
  if (isSpinning.value) return
  isSpinning.value = true
  currentResult.value = null

  // 1. 确定目标组合
  const target = findValidCombination()
  
  // 2. 计算旋转角度
  const degPerInitial = 360 / initials.length
  const degPerFinal = 360 / finals.length
  const degPerTone = 90

  // 目标角度计算：
  // 指针在上方 0 度。
  // CSS 旋转是顺时针正向。
  // 如果索引 i 在 0 位置（顶部），旋转容器 0 度。
  // 如果索引 i 在 1 位置（偏右），旋转容器 -1 * step 使其回到顶部。
  // 为了实现顺时针旋转效果，我们需要增加正向的角度。
  // 目标角度 = (圈数 * 360) + (360 - (idx * step)) 
  // 解释：idx * step 是该项相对于容器 0 度的偏移。要让它指向上方，容器应该旋转 -idx * step。
  // 也就是 360 - idx * step。
  
  const k = Math.ceil(rotations.initial / 360) + 5
  rotations.initial = k * 360 + (360 - (target.initialIdx * degPerInitial))

  const k2 = Math.ceil(rotations.final / 360) + 6
  rotations.final = k2 * 360 + (360 - (target.finalIdx * degPerFinal))

  const k3 = Math.ceil(rotations.tone / 360) + 8
  rotations.tone = k3 * 360 + (360 - (target.toneIdx * degPerTone))

  // 4. 动画结束后显示结果
  setTimeout(() => {
    isSpinning.value = false
    currentResult.value = {
      initial: initials[target.initialIdx],
      final: finals[target.finalIdx],
      tone: tones[target.toneIdx],
      pinyinWithTone: target.data.pinyin,
      word: target.data.word,
      example: target.data.example || target.data.word
    }
    playDing()
  }, 3000)
}

function playDing() {
  try {
    const ctx = new (window.AudioContext || (window as any).webkitAudioContext)()
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()
    osc.connect(gain)
    gain.connect(ctx.destination)
    osc.type = 'sine'
    // 叮~ (高频衰减)
    osc.frequency.setValueAtTime(1200, ctx.currentTime) 
    gain.gain.setValueAtTime(0.1, ctx.currentTime)
    gain.gain.exponentialRampToValueAtTime(0.001, ctx.currentTime + 0.8)
    osc.start()
    osc.stop(ctx.currentTime + 0.8)
  } catch (e) {}
}

import { speak as speakText } from '@/utils/speech'
import { useAppState } from '@/stores/appState'

const appState = useAppState()

async function speak() {
  if (!currentResult.value) return
  // 朗读：词组 + 逐字
  const text = `${currentResult.value.word}。${currentResult.value.example}`
  
  await speakText(text, {
    voiceURI: appState.speech.voiceURI || undefined,
    rate: appState.speech.rate || 0.9,
    pitch: appState.speech.pitch || 1.0
  })
}
</script>

<style scoped>
::-webkit-scrollbar { display: none; }

@keyframes fade-in {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fade-in {
  animation: fade-in 0.5s ease-out forwards;
}
</style>
