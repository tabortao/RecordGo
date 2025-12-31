<template>
  <div class="h-full flex flex-col bg-gray-50 dark:bg-gray-900 p-4 overflow-y-auto">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center gap-3">
        <el-button circle :icon="ArrowLeft" @click="router.back()" />
        <h2 class="text-xl font-bold dark:text-white">默写挑战</h2>
      </div>
    </div>

    <div v-if="poem" class="max-w-2xl mx-auto w-full bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm">
      <div class="text-center mb-8">
        <h1 class="text-2xl font-bold mb-2 dark:text-white">{{ poem.title_cns }}</h1>
        <p class="text-gray-500 text-sm">[{{ poem.dynasty_cns }}] {{ poem.author_cns }}</p>
      </div>

      <!-- Input Area -->
      <div v-if="!result" class="space-y-6">
        <div v-for="(_line, index) in poem.paragraphs_cns" :key="index" class="flex flex-col gap-2">
          <label class="text-xs text-gray-400 ml-1">第 {{ index + 1 }} 句</label>
          <el-input 
            v-model="inputs[index]" 
            :placeholder="`请输入第 ${index + 1} 句`"
            size="large"
            @keydown.enter.prevent="focusNext(index)"
            :ref="el => setInputRef(el, index)"
          />
        </div>

        <div class="pt-6 flex justify-center">
          <el-button type="primary" size="large" class="w-48 !rounded-full" @click="submitDictation" :disabled="!canSubmit">
            提交默写
          </el-button>
        </div>
      </div>

      <!-- Result Area -->
      <div v-else class="space-y-6">
        <div class="text-center">
          <div class="text-6xl font-bold mb-2 transition-colors" :class="getScoreColor(result.accuracy)">
            {{ result.accuracy }}<span class="text-2xl">%</span>
          </div>
          <p class="text-gray-500">准确率</p>
        </div>

        <div class="space-y-4">
          <div v-for="(lineResult, index) in checkResults" :key="index" class="p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
             <div class="flex flex-wrap gap-1 text-lg mb-1">
                <span v-for="(char, cIdx) in lineResult.diffs" :key="cIdx" 
                  class="w-8 h-8 flex items-center justify-center rounded border"
                  :class="{
                    'border-green-200 bg-green-50 text-green-600 dark:bg-green-900/30 dark:border-green-800 dark:text-green-400': char.status === 'correct',
                    'border-red-200 bg-red-50 text-red-600 dark:bg-red-900/30 dark:border-red-800 dark:text-red-400': char.status === 'wrong',
                    'border-gray-200 bg-gray-100 text-gray-400 line-through dark:bg-gray-600 dark:border-gray-500': char.status === 'extra',
                    'border-yellow-200 bg-yellow-50 text-yellow-600 dark:bg-yellow-900/30 dark:border-yellow-800 dark:text-yellow-400': char.status === 'missing'
                  }">
                  {{ char.char }}
                </span>
             </div>
             <div v-if="lineResult.hasError" class="mt-2 text-sm text-gray-500 flex items-center gap-2">
               <el-tag size="small" type="info">原文</el-tag>
               <span>{{ poem.paragraphs_cns[index] }}</span>
             </div>
          </div>
        </div>

        <div class="flex justify-center gap-4 pt-4">
          <el-button @click="reset" size="large">再试一次</el-button>
          <el-button type="primary" size="large" @click="router.back()">完成</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePoetryStore } from '../stores/poetryStore'
import { ArrowLeft } from '@element-plus/icons-vue'
import { v4 as uuidv4 } from 'uuid'
import { ElMessage } from 'element-plus'
import type { DictationResult } from '../types'

const route = useRoute()
const router = useRouter()
const store = usePoetryStore()

const poem = computed(() => store.getPoemById(Number(route.params.id)))
const inputs = ref<string[]>([])
const result = ref<DictationResult | null>(null)
const inputRefs = ref<any[]>([])

// Diff structures
interface CharDiff {
  char: string
  status: 'correct' | 'wrong' | 'missing' | 'extra'
  expected?: string
}
interface LineResult {
  diffs: CharDiff[]
  hasError: boolean
}
const checkResults = ref<LineResult[]>([])

onMounted(() => {
  if (!poem.value) {
    ElMessage.error('未找到该诗词')
    router.back()
    return
  }
  inputs.value = new Array(poem.value.paragraphs_cns.length).fill('')
  // Auto focus first input
  nextTick(() => {
    if (inputRefs.value[0]) {
      inputRefs.value[0].focus()
    }
  })
})

const setInputRef = (el: any, index: number) => {
  if (el) {
    // el is likely an Element Plus component wrapper, try to find input
    inputRefs.value[index] = el
  }
}

const focusNext = (index: number) => {
  if (index < inputs.value.length - 1) {
    const nextInput = inputRefs.value[index + 1]
    if (nextInput) nextInput.focus()
  }
}

const canSubmit = computed(() => inputs.value.some(s => s.trim().length > 0))

const submitDictation = () => {
  if (!poem.value) return

  let totalChars = 0
  let correctChars = 0
  const errors: any[] = []
  const lines: LineResult[] = []

  poem.value.paragraphs_cns.forEach((expectedLine: string, lineIdx: number) => {
    // Remove punctuation for comparison? Usually dictation includes punctuation or ignores it?
    // Let's strictly compare including punctuation for now, but trim spaces.
    const expected = expectedLine.trim()
    const actual = (inputs.value[lineIdx] || '').trim()
    
    const maxLen = Math.max(expected.length, actual.length)
    const lineDiffs: CharDiff[] = []
    let lineHasError = false

    for (let i = 0; i < maxLen; i++) {
      const expChar = expected[i]
      const actChar = actual[i]
      
      if (expChar && actChar) {
        if (expChar === actChar) {
          lineDiffs.push({ char: actChar, status: 'correct' })
          correctChars++
        } else {
          lineDiffs.push({ char: actChar, status: 'wrong', expected: expChar })
          errors.push({ lineIndex: lineIdx, charIndex: i, expected: expChar, actual: actChar, type: 'wrong' })
          lineHasError = true
        }
      } else if (expChar && !actChar) {
        lineDiffs.push({ char: expChar, status: 'missing' })
        errors.push({ lineIndex: lineIdx, charIndex: i, expected: expChar, actual: '', type: 'missing' })
        lineHasError = true
      } else if (!expChar && actChar) {
        lineDiffs.push({ char: actChar, status: 'extra' })
        errors.push({ lineIndex: lineIdx, charIndex: i, expected: '', actual: actChar, type: 'extra' })
        lineHasError = true
      }
    }
    
    totalChars += expected.length
    lines.push({ diffs: lineDiffs, hasError: lineHasError })
  })

  checkResults.value = lines
  const accuracy = totalChars > 0 ? Math.round((correctChars / totalChars) * 100) : 0

  const newResult: DictationResult = {
    id: uuidv4(),
    poemId: poem.value.id,
    date: Date.now(),
    accuracy,
    input: [...inputs.value],
    errors
  }

  store.saveDictationResult(newResult)
  result.value = newResult
}

const reset = () => {
  result.value = null
  checkResults.value = []
  inputs.value = new Array(poem.value?.paragraphs_cns.length || 0).fill('')
  nextTick(() => {
    if (inputRefs.value[0]) inputRefs.value[0].focus()
  })
}

const getScoreColor = (score: number) => {
  if (score >= 90) return 'text-green-500'
  if (score >= 60) return 'text-yellow-500'
  return 'text-red-500'
}
</script>
