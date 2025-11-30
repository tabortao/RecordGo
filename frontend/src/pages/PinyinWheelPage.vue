<template>
  <div class="min-h-screen bg-green-50 dark:bg-green-900/20 flex flex-col relative overflow-hidden">
    <!-- 顶部导航 -->
    <div class="p-4 z-10 flex items-center justify-between">
      <div class="w-10 h-10 bg-white dark:bg-gray-800 rounded-full shadow flex items-center justify-center cursor-pointer hover:scale-105 transition" @click="router.back()">
        <el-icon :size="20" class="text-green-600"><ArrowLeft /></el-icon>
      </div>
      <h1 class="text-2xl font-bold text-green-700 dark:text-green-400 tracking-wider">拼音大转盘</h1>
      <div class="w-10"></div>
    </div>

    <!-- 游戏区域 -->
    <div class="flex-1 flex flex-col items-center justify-center relative -mt-10">
      
      <!-- 转盘容器 -->
      <div class="relative w-[320px] h-[320px] sm:w-[400px] sm:h-[400px] flex items-center justify-center">
        <!-- 指针 (固定在上方) -->
        <div class="absolute top-0 left-1/2 -translate-x-1/2 -translate-y-2 z-40">
          <div class="w-0 h-0 border-l-[12px] border-l-transparent border-r-[12px] border-r-transparent border-t-[24px] border-t-red-500 drop-shadow-md"></div>
        </div>

        <!-- 外圈：声母 -->
        <div 
          class="absolute inset-0 rounded-full border-4 border-green-200 bg-green-100 shadow-xl transition-transform duration-[3000ms] ease-out"
          :style="{ transform: `rotate(${rotations.initial}deg)` }"
        >
          <div v-for="(item, idx) in initials" :key="item"
            class="absolute top-0 left-1/2 -ml-[15px] w-[30px] h-1/2 origin-bottom pt-2 text-center font-bold text-green-800"
            :style="{ transform: `rotate(${idx * (360 / initials.length)}deg)` }"
          >
            <span class="block transform -rotate-90 mt-2">{{ item }}</span>
          </div>
        </div>

        <!-- 中圈：韵母 -->
        <div 
          class="absolute top-[15%] left-[15%] right-[15%] bottom-[15%] rounded-full border-4 border-yellow-200 bg-yellow-50 shadow-lg transition-transform duration-[3000ms] ease-out"
          :style="{ transform: `rotate(${rotations.final}deg)` }"
        >
          <div v-for="(item, idx) in finals" :key="item"
            class="absolute top-0 left-1/2 -ml-[15px] w-[30px] h-1/2 origin-bottom pt-2 text-center font-bold text-orange-700"
            :style="{ transform: `rotate(${idx * (360 / finals.length)}deg)` }"
          >
            <span class="block transform -rotate-90 mt-1 text-sm">{{ item }}</span>
          </div>
        </div>

        <!-- 内圈：声调 -->
        <div 
          class="absolute top-[35%] left-[35%] right-[35%] bottom-[35%] rounded-full border-4 border-pink-200 bg-pink-50 shadow-md transition-transform duration-[3000ms] ease-out flex items-center justify-center"
          :style="{ transform: `rotate(${rotations.tone}deg)` }"
        >
          <div v-for="(item, idx) in tones" :key="item"
            class="absolute top-0 left-1/2 -ml-[10px] w-[20px] h-1/2 origin-bottom pt-1 text-center font-bold text-pink-600"
            :style="{ transform: `rotate(${idx * 90}deg)` }"
          >
            <span class="block mt-1">{{ item }}声</span>
          </div>
          <!-- 中心装饰 -->
          <div class="w-12 h-12 rounded-full bg-white shadow-inner flex items-center justify-center z-10">
            <span class="text-2xl">🐼</span>
          </div>
        </div>
      </div>

      <!-- 结果展示区 -->
      <div class="mt-8 bg-white dark:bg-gray-800 rounded-2xl p-6 shadow-lg border-b-4 border-green-500 w-[90%] max-w-md text-center min-h-[160px] flex flex-col items-center justify-center transition-all" :class="{'scale-105': showResult}">
        <div v-if="isSpinning" class="text-green-600 animate-pulse font-bold text-xl">
          拼读中...
        </div>
        <div v-else-if="currentResult" class="space-y-3">
          <div class="flex items-center justify-center gap-2 text-2xl font-mono text-gray-500">
            <span>{{ currentResult.initial }}</span>
            <span>+</span>
            <span>{{ currentResult.final }}</span>
            <span>+</span>
            <span>{{ currentResult.tone }}声</span>
            <span>=</span>
            <span class="text-green-600 font-bold text-4xl">{{ currentResult.pinyinWithTone }}</span>
          </div>
          <div class="text-5xl font-bold text-gray-800 dark:text-white mt-2">{{ currentResult.word }}</div>
          <div class="text-gray-400 text-sm">{{ currentResult.example }}</div>
        </div>
        <div v-else class="text-gray-400">
          点击“转动”开始学习吧！
        </div>
      </div>

      <!-- 控制按钮 -->
      <div class="mt-8 flex gap-6">
        <button 
          class="bg-gradient-to-b from-green-400 to-green-600 text-white px-8 py-3 rounded-full shadow-lg border-b-4 border-green-700 active:border-b-0 active:translate-y-1 transition font-bold text-xl flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
          @click="spin"
          :disabled="isSpinning"
        >
          <el-icon><Refresh /></el-icon> 转动
        </button>
        
        <button 
          class="bg-gradient-to-b from-yellow-400 to-yellow-600 text-white px-8 py-3 rounded-full shadow-lg border-b-4 border-yellow-700 active:border-b-0 active:translate-y-1 transition font-bold text-xl flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
          @click="speak"
          :disabled="isSpinning || !currentResult"
        >
          <el-icon><Microphone /></el-icon> 朗读
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Refresh, Microphone } from '@element-plus/icons-vue'

const router = useRouter()

// --- 数据定义 ---
const initials = ['b', 'p', 'm', 'f', 'd', 't', 'n', 'l', 'g', 'k', 'h', 'j', 'q', 'x', 'zh', 'ch', 'sh', 'r', 'z', 'c', 's', 'y', 'w']
const finals = ['a', 'o', 'e', 'i', 'u', 'ü', 'ai', 'ei', 'ui', 'ao', 'ou', 'iu', 'ie', 'üe', 'er', 'an', 'en', 'in', 'un', 'ün', 'ang', 'eng', 'ing', 'ong']
const tones = [1, 2, 3, 4]

// 简单的拼音数据字典 (Valid combinations)
// 格式: key = "initial+final", value = { 1: "word", ... }
const pinyinDict: Record<string, Record<number, { word: string, pinyin: string }>> = {
  'ba': { 1: { word: '八', pinyin: 'bā' }, 2: { word: '拔', pinyin: 'bá' }, 3: { word: '把', pinyin: 'bǎ' }, 4: { word: '爸', pinyin: 'bà' } },
  'po': { 1: { word: '坡', pinyin: 'pō' }, 2: { word: '婆', pinyin: 'pó' }, 3: { word: '叵', pinyin: 'pǒ' }, 4: { word: '破', pinyin: 'pò' } },
  'mi': { 1: { word: '咪', pinyin: 'mī' }, 2: { word: '迷', pinyin: 'mí' }, 3: { word: '米', pinyin: 'mǐ' }, 4: { word: '密', pinyin: 'mì' } },
  'fu': { 1: { word: '夫', pinyin: 'fū' }, 2: { word: '福', pinyin: 'fú' }, 3: { word: '府', pinyin: 'fǔ' }, 4: { word: '富', pinyin: 'fù' } },
  'da': { 1: { word: '搭', pinyin: 'dā' }, 2: { word: '达', pinyin: 'dá' }, 3: { word: '打', pinyin: 'dǎ' }, 4: { word: '大', pinyin: 'dà' } },
  'tu': { 1: { word: '秃', pinyin: 'tū' }, 2: { word: '图', pinyin: 'tú' }, 3: { word: '土', pinyin: 'tǔ' }, 4: { word: '兔', pinyin: 'tù' } },
  'nü': { 3: { word: '女', pinyin: 'nǚ' } }, 
  'lü': { 4: { word: '绿', pinyin: 'lǜ' }, 3: { word: '旅', pinyin: 'lǚ' } },
  'ge': { 1: { word: '歌', pinyin: 'gē' }, 2: { word: '格', pinyin: 'gé' }, 3: { word: '葛', pinyin: 'gě' }, 4: { word: '个', pinyin: 'gè' } },
  'ka': { 1: { word: '咖', pinyin: 'kā' }, 3: { word: '卡', pinyin: 'kǎ' } },
  'he': { 1: { word: '喝', pinyin: 'hē' }, 2: { word: '河', pinyin: 'hé' }, 4: { word: '贺', pinyin: 'hè' } },
  'ji': { 1: { word: '鸡', pinyin: 'jī' }, 2: { word: '急', pinyin: 'jí' }, 3: { word: '几', pinyin: 'jǐ' }, 4: { word: '寄', pinyin: 'jì' } },
  'qi': { 1: { word: '七', pinyin: 'qī' }, 2: { word: '齐', pinyin: 'qí' }, 3: { word: '起', pinyin: 'qǐ' }, 4: { word: '气', pinyin: 'qì' } },
  'xi': { 1: { word: '西', pinyin: 'xī' }, 2: { word: '习', pinyin: 'xí' }, 3: { word: '洗', pinyin: 'xǐ' }, 4: { word: '戏', pinyin: 'xì' } },
  'zhu': { 1: { word: '猪', pinyin: 'zhū' }, 2: { word: '竹', pinyin: 'zhú' }, 3: { word: '主', pinyin: 'zhǔ' }, 4: { word: '住', pinyin: 'zhù' } },
  'chi': { 1: { word: '吃', pinyin: 'chī' }, 2: { word: '迟', pinyin: 'chí' }, 3: { word: '尺', pinyin: 'chǐ' }, 4: { word: '翅', pinyin: 'chì' } },
  'shu': { 1: { word: '书', pinyin: 'shū' }, 2: { word: '熟', pinyin: 'shú' }, 3: { word: '鼠', pinyin: 'shǔ' }, 4: { word: '树', pinyin: 'shù' } },
  'ri': { 4: { word: '日', pinyin: 'rì' } },
  'za': { 1: { word: '扎', pinyin: 'zā' }, 2: { word: '杂', pinyin: 'zá' } },
  'ci': { 1: { word: '刺', pinyin: 'cī' }, 2: { word: '词', pinyin: 'cí' }, 3: { word: '此', pinyin: 'cǐ' }, 4: { word: '次', pinyin: 'cì' } },
  'si': { 1: { word: '丝', pinyin: 'sī' }, 3: { word: '死', pinyin: 'sǐ' }, 4: { word: '四', pinyin: 'sì' } },
  'ya': { 1: { word: '鸭', pinyin: 'yā' }, 2: { word: '牙', pinyin: 'yá' }, 3: { word: '雅', pinyin: 'yǎ' }, 4: { word: '亚', pinyin: 'yà' } },
  'wo': { 1: { word: '窝', pinyin: 'wō' }, 3: { word: '我', pinyin: 'wǒ' }, 4: { word: '握', pinyin: 'wò' } },
  'ma': { 1: { word: '妈', pinyin: 'mā' }, 2: { word: '麻', pinyin: 'má' }, 3: { word: '马', pinyin: 'mǎ' }, 4: { word: '骂', pinyin: 'mà' } },
}

// 状态
const isSpinning = ref(false)
const showResult = ref(false)
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
  // 随机尝试 100 次，找不到就返回默认
  for (let i = 0; i < 100; i++) {
    const iIdx = Math.floor(Math.random() * initials.length)
    const fIdx = Math.floor(Math.random() * finals.length)
    const tIdx = Math.floor(Math.random() * tones.length) // 0-3 representing tones 1-4

    const initial = initials[iIdx]
    const final = finals[fIdx]
    const tone = tones[tIdx]

    const key = initial + final
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
  showResult.value = false
  currentResult.value = null

  // 1. 确定目标组合
  const target = findValidCombination()
  
  // 2. 计算旋转角度
  // 每个扇区角度
  const degPerInitial = 360 / initials.length
  const degPerFinal = 360 / finals.length
  const degPerTone = 90

  // 目标角度 = (圈数 * 360) - (索引 * 单个角度) 
  // 指针在 12 点方向 (0度)，所以我们需要让目标扇区转到 0 度。
  // 初始状态 0 度对应索引 0。
  // 要让索引 N 转到 0 度，需要旋转 -N * deg。为了正向旋转效果，加多圈。
  
  const extraSpins = 5 * 360 // 至少转 5 圈
  
  // 随机偏移一点点让它看起来自然 (可选，这里为了精准对齐先不加随机偏移)
  rotations.initial += extraSpins + (360 - target.initialIdx * degPerInitial) % 360
  // 累加旋转，保证一直往一个方向转
  // 修正：直接设为累加值
  // 让我们计算增量：
  // current % 360 是当前位置。
  // 目标是 targetPos. 
  // newRotation = current + extra + (targetPos - currentPos)
  
  // 简单做法：
  // 每次都在当前基础上增加至少 5 圈 + 目标偏移
  // 目标角度： 让 index 处在顶部。
  // CSS transform rotate(X deg). 0 deg is top. 
  // items are placed at rotate(idx * step). Item 0 is at 0 deg.
  // To bring Item N (at N*step) to top (0 deg), we must rotate container by -N*step.
  
  // Initial
  const targetRotInitial = -(target.initialIdx * degPerInitial)
  
  // 简化逻辑：随便转几圈，最后停在目标
  // 实际上不需要精确计算当前，只需要算最终目标角度 = k * 360 - index * step
  // 为了保证动画顺滑，k 必须比当前 k 大
  const k = Math.ceil(rotations.initial / 360) + 5
  rotations.initial = k * 360 - (target.initialIdx * degPerInitial)

  const k2 = Math.ceil(rotations.final / 360) + 6 // 中圈多转一圈
  rotations.final = k2 * 360 - (target.finalIdx * degPerFinal)

  const k3 = Math.ceil(rotations.tone / 360) + 7 // 内圈再多转
  rotations.tone = k3 * 360 - (target.toneIdx * degPerTone)

  // 3. 播放音效 (模拟)
  // playSpinSound() 

  // 4. 动画结束后显示结果
  setTimeout(() => {
    isSpinning.value = false
    showResult.value = true
    currentResult.value = {
      initial: initials[target.initialIdx],
      final: finals[target.finalIdx],
      tone: tones[target.toneIdx],
      pinyinWithTone: target.data.pinyin,
      word: target.data.word,
      example: target.data.word
    }
    // 播放提示音
    playDing()
  }, 3000)
}

function playDing() {
  try {
    // 简单的 Web Audio API 提示音
    const ctx = new (window.AudioContext || (window as any).webkitAudioContext)()
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()
    osc.connect(gain)
    gain.connect(ctx.destination)
    osc.type = 'sine'
    osc.frequency.setValueAtTime(523.25, ctx.currentTime) // C5
    osc.frequency.exponentialRampToValueAtTime(880, ctx.currentTime + 0.1)
    gain.gain.setValueAtTime(0.1, ctx.currentTime)
    gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.5)
    osc.start()
    osc.stop(ctx.currentTime + 0.5)
  } catch (e) {}
}

function speak() {
  if (!currentResult.value) return
  const u = new SpeechSynthesisUtterance(currentResult.value.word)
  u.lang = 'zh-CN'
  u.rate = 0.8
  window.speechSynthesis.speak(u)
}
</script>

<style scoped>
/* 隐藏滚动条 */
::-webkit-scrollbar { display: none; }
</style>
