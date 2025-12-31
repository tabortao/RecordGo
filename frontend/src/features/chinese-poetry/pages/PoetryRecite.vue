<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 flex flex-col">
     <!-- Header -->
    <div class="bg-white dark:bg-gray-800 p-4 flex items-center justify-between border-b border-gray-200 dark:border-gray-700">
      <el-icon :size="24" class="cursor-pointer" @click="router.back()"><Close /></el-icon>
      <div class="font-bold text-lg">背诵挑战</div>
      <div class="w-6"></div>
    </div>

    <!-- Content -->
    <div class="flex-1 flex flex-col items-center justify-center p-6 space-y-8" v-if="poem">
       <div class="text-center space-y-2">
           <h2 class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ poem.title_cns }}</h2>
           <div class="text-gray-500">{{ poem.author_cns }}</div>
       </div>

       <!-- Interaction Area -->
       <div class="w-full max-w-md bg-white dark:bg-gray-800 rounded-2xl p-6 shadow-sm min-h-[300px] flex flex-col items-center justify-center relative">
           
           <!-- Mode: Full Hide / Partial Hide -->
           <div class="space-y-4 text-xl leading-loose text-center font-serif" v-if="mode !== 'full_hide'">
               <p v-for="(line, idx) in poem.paragraphs_cns" :key="idx">
                   <span v-if="mode === 'show'">{{ line }}</span>
                   <span v-else-if="mode === 'partial_hide'">
                       {{ idx % 2 === 0 ? line : '_______' }}
                   </span>
                   <span v-else-if="mode === 'first_char'">
                       {{ line[0] }}______
                   </span>
               </p>
           </div>
           
           <div v-else class="text-gray-400 italic">
               (请开始背诵)
           </div>

           <!-- Recording Animation -->
           <div v-if="isRecording" class="absolute inset-0 bg-black/5 flex flex-col items-center justify-center rounded-2xl backdrop-blur-[2px]">
               <div class="animate-pulse text-pink-500 font-bold text-xl mb-4">正在录音...</div>
               <div class="flex gap-1 h-8 items-end">
                   <div class="w-1 bg-pink-500 animate-[bounce_1s_infinite] h-4"></div>
                   <div class="w-1 bg-pink-500 animate-[bounce_1.2s_infinite] h-6"></div>
                   <div class="w-1 bg-pink-500 animate-[bounce_0.8s_infinite] h-3"></div>
               </div>
           </div>
       </div>

       <!-- Controls -->
       <div class="flex flex-col items-center gap-6 w-full max-w-md relative">
           <!-- Modes -->
           <div class="flex gap-2 bg-gray-100 dark:bg-gray-700 p-1 rounded-lg">
               <button @click="mode = 'show'" :class="getModeClass('show')">看全文</button>
               <button @click="mode = 'partial_hide'" :class="getModeClass('partial_hide')">填空</button>
               <button @click="mode = 'first_char'" :class="getModeClass('first_char')">首字</button>
               <button @click="mode = 'full_hide'" :class="getModeClass('full_hide')">背诵</button>
           </div>

           <!-- Record Button -->
           <button 
             @touchstart.prevent="startRecord" 
             @touchend.prevent="stopRecord"
             @mousedown.prevent="startRecord"
             @mouseup.prevent="stopRecord"
             class="w-20 h-20 rounded-full bg-pink-500 shadow-lg shadow-pink-200 flex items-center justify-center active:scale-95 transition"
           >
               <el-icon :size="40" class="text-white"><Microphone /></el-icon>
           </button>
           <div class="text-sm text-gray-500">按住 说话</div>
           
           <!-- Action Buttons (Bottom Left/Right) -->
           <div class="absolute bottom-0 left-0">
               <el-button type="primary" plain round size="default" @click="handleComplete">完成打卡</el-button>
           </div>
           <div class="absolute bottom-0 right-0">
               <el-button type="success" plain round size="default" @click="handleMastered">记住了</el-button>
           </div>
       </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Close, Microphone } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Poem } from '../types'
import { usePoetryStore } from '../stores/poetryStore'

const route = useRoute()
const router = useRouter()
const store = usePoetryStore()

const poem = ref<Poem | null>(null)
const mode = ref('show')
const isRecording = ref(false)

onMounted(() => {
    store.init()
    const id = Number(route.params.id)
    poem.value = store.getPoemById(id) || null
})

const getModeClass = (m: string) => {
    const base = "px-3 py-1.5 rounded-md text-sm transition "
    if (mode.value === m) return base + "bg-white dark:bg-gray-600 shadow-sm font-bold text-gray-800 dark:text-gray-100"
    return base + "text-gray-500"
}

const startRecord = () => {
    isRecording.value = true
    // Mock recording start
}

const stopRecord = () => {
    isRecording.value = false
    // Mock recording end & verify
    setTimeout(() => {
        // Simple mock result
        const score = Math.floor(Math.random() * 20) + 80
        ElMessageBox.alert(`背诵完成！\n准确度：${score}%`, '挑战结果', {
            confirmButtonText: '完成打卡',
            callback: () => {
                if (poem.value) {
                    store.recordStudy(poem.value.id, score >= 90)
                    ElMessage.success('打卡成功！')
                    router.back()
                }
            }
        })
    }, 500)
}

const handleComplete = () => {
    if (poem.value) {
        store.recordStudy(poem.value.id, false)
        ElMessage.success('打卡成功，学习次数 +1')
        router.back()
    }
}

const handleMastered = () => {
    if (poem.value) {
        store.recordStudy(poem.value.id, true)
        ElMessage.success('太棒了！已标记为熟练背诵')
        router.back()
    }
}
</script>
