<template>
  <div class="min-h-screen bg-[#F5F7FA] dark:bg-gray-900 flex flex-col">
    <!-- Header -->
    <div class="px-4 py-3 bg-white/80 dark:bg-gray-900/80 backdrop-blur-md sticky top-0 z-30 flex items-center justify-between border-b border-gray-100 dark:border-gray-800 transition-colors">
      <div class="flex items-center gap-3">
        <div 
          class="w-8 h-8 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center cursor-pointer hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors text-gray-600 dark:text-gray-300"
          @click="router.back()"
        >
          <el-icon><ArrowLeft /></el-icon>
        </div>
        <h2 class="font-bold text-gray-800 dark:text-gray-100 text-lg">成绩详情</h2>
      </div>
      <div class="flex items-center gap-2">
        <button 
          class="w-8 h-8 rounded-full bg-gray-50 dark:bg-gray-800 flex items-center justify-center text-gray-600 dark:text-gray-400 hover:bg-indigo-50 hover:text-indigo-600 dark:hover:bg-indigo-900/30 dark:hover:text-indigo-400 transition-colors"
          @click="openEdit"
        >
          <el-icon><Edit /></el-icon>
        </button>
        <button 
          class="w-8 h-8 rounded-full bg-gray-50 dark:bg-gray-800 flex items-center justify-center text-gray-600 dark:text-gray-400 hover:bg-red-50 hover:text-red-600 dark:hover:bg-red-900/30 dark:hover:text-red-400 transition-colors"
          @click="onDelete"
        >
          <el-icon><Delete /></el-icon>
        </button>
      </div>
    </div>

    <div v-if="!record" class="flex-1 flex items-center justify-center">
      <div class="w-8 h-8 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <div v-else class="flex-1 overflow-y-auto p-4 sm:p-6 pb-24">
      <div class="max-w-3xl mx-auto space-y-6">
        
        <!-- Score Card -->
        <div class="bg-white dark:bg-gray-800 rounded-3xl p-8 shadow-sm border border-gray-100 dark:border-gray-700 relative overflow-hidden text-center">
          <div class="absolute top-0 left-0 w-full h-2 bg-gradient-to-r" :class="subjectGradientClass(record.subject)"></div>
          <div class="absolute -right-10 -top-10 w-40 h-40 bg-gray-50 dark:bg-gray-700/30 rounded-full blur-3xl -z-0"></div>
          <div class="absolute -left-10 -bottom-10 w-40 h-40 bg-gray-50 dark:bg-gray-700/30 rounded-full blur-3xl -z-0"></div>
          
          <div class="relative z-10">
            <div class="inline-flex items-center justify-center px-3 py-1 rounded-full text-xs font-bold mb-4" :class="subjectBadgeClass(record.subject)">
              {{ record.subject }}
            </div>
            
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">{{ record.exam_name }}</h1>
            <div class="text-sm text-gray-500 dark:text-gray-400 mb-8">{{ formatDate(record.exam_date) }}</div>
            
            <div class="flex items-baseline justify-center mb-2">
              <span class="text-6xl font-black text-gray-900 dark:text-white tracking-tighter" :class="scoreColorClass(record.score, record.full_score)">
                {{ record.score }}
              </span>
              <span class="text-xl text-gray-400 font-medium ml-2">/ {{ record.full_score }}</span>
            </div>
            
            <div class="w-full h-1.5 bg-gray-100 dark:bg-gray-700 rounded-full overflow-hidden max-w-xs mx-auto mt-4">
               <div class="h-full rounded-full transition-all duration-1000 ease-out" :class="subjectGradientClass(record.subject)" :style="{ width: `${(record.score / record.full_score) * 100}%` }"></div>
            </div>
          </div>
        </div>

        <!-- Details Grid -->
        <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
          <div class="bg-white dark:bg-gray-800 rounded-2xl p-4 border border-gray-100 dark:border-gray-700 shadow-sm flex flex-col items-center justify-center text-center group hover:border-indigo-100 dark:hover:border-indigo-900/30 transition-colors">
            <div class="w-10 h-10 rounded-full bg-amber-50 dark:bg-amber-900/20 text-amber-500 flex items-center justify-center mb-2 group-hover:scale-110 transition-transform">
              <el-icon :size="20"><Trophy /></el-icon>
            </div>
            <div class="text-xs text-gray-400 mb-1">班级排名</div>
            <div class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ record.class_rank || '-' }}</div>
          </div>
          
          <div class="bg-white dark:bg-gray-800 rounded-2xl p-4 border border-gray-100 dark:border-gray-700 shadow-sm flex flex-col items-center justify-center text-center group hover:border-indigo-100 dark:hover:border-indigo-900/30 transition-colors">
            <div class="w-10 h-10 rounded-full bg-blue-50 dark:bg-blue-900/20 text-blue-500 flex items-center justify-center mb-2 group-hover:scale-110 transition-transform">
              <el-icon :size="20"><DataLine /></el-icon>
            </div>
            <div class="text-xs text-gray-400 mb-1">班级均分</div>
            <div class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ record.class_avg || '-' }}</div>
          </div>

          <div class="bg-white dark:bg-gray-800 rounded-2xl p-4 border border-gray-100 dark:border-gray-700 shadow-sm flex flex-col items-center justify-center text-center group hover:border-indigo-100 dark:hover:border-indigo-900/30 transition-colors md:col-span-1 col-span-2">
             <div class="w-10 h-10 rounded-full bg-purple-50 dark:bg-purple-900/20 text-purple-500 flex items-center justify-center mb-2 group-hover:scale-110 transition-transform">
              <el-icon :size="20"><PieChart /></el-icon>
            </div>
            <div class="text-xs text-gray-400 mb-1">得分率</div>
            <div class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ Math.round((record.score / record.full_score) * 100) }}%</div>
          </div>
        </div>

        <!-- Remark -->
        <div v-if="record.remark" class="bg-white dark:bg-gray-800 rounded-2xl p-5 border border-gray-100 dark:border-gray-700 shadow-sm">
          <div class="flex items-center gap-2 mb-3 text-gray-900 dark:text-white font-bold text-sm">
            <el-icon class="text-indigo-500"><ChatLineSquare /></el-icon>
            <span>备注说明</span>
          </div>
          <div class="text-sm text-gray-600 dark:text-gray-300 leading-relaxed bg-gray-50 dark:bg-gray-900/50 p-3 rounded-xl">
            {{ record.remark }}
          </div>
        </div>

        <!-- Photo -->
        <div v-if="displayPhoto" class="bg-white dark:bg-gray-800 rounded-2xl p-1 border border-gray-100 dark:border-gray-700 shadow-sm">
           <div class="flex items-center justify-between px-4 py-3">
             <div class="flex items-center gap-2 text-gray-900 dark:text-white font-bold text-sm">
               <el-icon class="text-indigo-500"><Picture /></el-icon>
               <span>试卷/成绩单</span>
             </div>
           </div>
           <div class="rounded-xl overflow-hidden bg-gray-100 dark:bg-gray-900/50 flex items-center justify-center min-h-[200px]">
             <img 
               :src="displayPhoto" 
               class="w-full h-auto object-contain max-h-[80vh] cursor-zoom-in transition-transform hover:scale-[1.01]"
               @click="previewImage(displayPhoto)"
             />
           </div>
        </div>

      </div>
    </div>

    <!-- Edit Dialog (Reused existing logic) -->
    <el-dialog v-model="showDialog" title="编辑成绩" width="520px" class="rounded-2xl" destroy-on-close>
      <div class="space-y-4 py-2">
        <div class="space-y-1.5">
          <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">学科</label>
          <el-select v-model="formSubject" placeholder="请选择学科" class="w-full" size="large">
            <el-option v-for="s in subjectOptions" :key="s" :label="s" :value="s" />
          </el-select>
        </div>
        <div class="space-y-1.5">
          <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">考试名称</label>
          <el-input v-model="formExamName" placeholder="请输入考试名称" size="large" />
        </div>
        <div class="space-y-1.5">
          <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">考试日期</label>
          <el-date-picker v-model="formExamDate" type="date" value-format="YYYY-MM-DD" placeholder="选择日期" class="!w-full" size="large" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">分数</label>
            <el-input-number v-model="formScore" :min="0" class="!w-full" size="large" :controls="false" />
          </div>
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">满分</label>
            <el-input-number v-model="formFullScore" :min="1" class="!w-full" size="large" :controls="false" />
          </div>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">排名</label>
            <el-input-number v-model="formRank" :min="1" class="!w-full" size="large" :controls="false" placeholder="未设置" />
          </div>
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">班级均分</label>
            <el-input-number v-model="formClassAvg" :min="0" class="!w-full" size="large" :controls="false" placeholder="未设置" />
          </div>
        </div>
        <div class="space-y-1.5">
          <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">备注</label>
          <el-input v-model="formRemark" type="textarea" :rows="3" placeholder="可填写补充说明" resize="none" />
        </div>
        <div class="space-y-1.5">
          <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">成绩照片</label>
          <div class="flex items-center gap-4 p-3 border border-dashed border-gray-300 dark:border-gray-600 rounded-xl bg-gray-50 dark:bg-gray-800/50">
            <div class="w-20 h-20 rounded-lg overflow-hidden bg-gray-200 dark:bg-gray-700 flex items-center justify-center flex-shrink-0 relative group">
              <img v-if="photoPreview" :src="photoPreview" class="w-full h-full object-cover" />
              <el-icon v-else :size="24" class="text-gray-400"><Picture /></el-icon>
              
              <div class="absolute inset-0 bg-black/40 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer" @click="triggerFileInput">
                <span class="text-white text-xs">更换</span>
              </div>
            </div>
            <div class="flex-1">
               <button 
                 class="px-4 py-2 bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-lg text-sm font-medium hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors"
                 @click="triggerFileInput"
               >
                 选择图片
               </button>
               <input ref="fileInputRef" type="file" accept="image/*" class="hidden" @change="onSelectPhoto" />
               <div class="text-xs text-gray-400 mt-2">支持 jpg, png, heic 等格式，自动转换为 webp</div>
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-3 pt-2">
          <el-button @click="showDialog = false" class="!rounded-xl px-6">取消</el-button>
          <el-button type="primary" :loading="saving" @click="submit" class="!rounded-xl px-6 bg-indigo-600 border-indigo-600 hover:bg-indigo-700">保存</el-button>
        </div>
      </template>
    </el-dialog>

    <el-image-viewer
      v-if="showImageViewer"
      :url-list="[imageViewerUrl]"
      @close="showImageViewer = false"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ArrowLeft, Edit, Delete, Trophy, DataLine, PieChart, ChatLineSquare, Picture } from '@element-plus/icons-vue'
import { useRoute } from 'vue-router'
import router from '@/router'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteScore, getScore, updateScore, type ScoreRecord } from '@/services/scores'
import { getStorageInfo, presignUpload, presignView, putToURL } from '@/services/storage'
import { toWebp } from '@/utils/image'
import { useAuth } from '@/stores/auth'

const auth = useAuth()
const route = useRoute()
const record = ref<ScoreRecord | null>(null)
const displayPhoto = ref('')

const showDialog = ref(false)
const saving = ref(false)
const fileInputRef = ref<HTMLInputElement>()

const formSubject = ref('')
const formExamName = ref('')
const formExamDate = ref<string>(dayjs().format('YYYY-MM-DD'))
const formScore = ref<number | null>(null)
const formFullScore = ref<number | null>(null)
const formRank = ref<number | null>(null)
const formClassAvg = ref<number | null>(null)
const formRemark = ref('')
const photoFile = ref<File | null>(null)
const photoPreview = ref<string>('')
const storageInfo = ref<{ driver: string; public_base_url?: string } | null>(null)

const showImageViewer = ref(false)
const imageViewerUrl = ref('')

const subjectOptions = [
  '语文', '数学', '英语', '科学', '体育', '物理', '化学', '政治', '历史', '地理', '综合', '自定义'
]

const scoreId = computed(() => Number(route.params.id || 0))

function formatDate(d: string) {
  return dayjs(d).format('YYYY年MM月DD日')
}

function subjectGradientClass(subject: string) {
    const map: Record<string, string> = {
    语文: 'bg-gradient-to-r from-amber-400 to-orange-500',
    数学: 'bg-gradient-to-r from-blue-400 to-indigo-500',
    英语: 'bg-gradient-to-r from-emerald-400 to-teal-500',
    科学: 'bg-gradient-to-r from-cyan-400 to-sky-500',
    物理: 'bg-gradient-to-r from-purple-400 to-violet-500',
    化学: 'bg-gradient-to-r from-rose-400 to-pink-500',
    生物: 'bg-gradient-to-r from-lime-400 to-green-500',
    政治: 'bg-gradient-to-r from-red-400 to-rose-500',
    历史: 'bg-gradient-to-r from-yellow-400 to-amber-500',
    地理: 'bg-gradient-to-r from-teal-400 to-emerald-500',
  }
  return map[subject] || 'bg-gradient-to-r from-gray-400 to-gray-500'
}

function subjectBadgeClass(subject: string) {
  const map: Record<string, string> = {
    语文: 'bg-amber-100 text-amber-700 dark:bg-amber-500/20 dark:text-amber-300',
    数学: 'bg-blue-100 text-blue-700 dark:bg-blue-500/20 dark:text-blue-300',
    英语: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/20 dark:text-emerald-300',
    科学: 'bg-cyan-100 text-cyan-700 dark:bg-cyan-500/20 dark:text-cyan-300',
    物理: 'bg-purple-100 text-purple-700 dark:bg-purple-500/20 dark:text-purple-300',
    化学: 'bg-rose-100 text-rose-700 dark:bg-rose-500/20 dark:text-rose-300',
    生物: 'bg-lime-100 text-lime-700 dark:bg-lime-500/20 dark:text-lime-300',
    政治: 'bg-red-100 text-red-700 dark:bg-red-500/20 dark:text-red-300',
    历史: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-500/20 dark:text-yellow-300',
    地理: 'bg-teal-100 text-teal-700 dark:bg-teal-500/20 dark:text-teal-300',
  }
  return map[subject] || 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
}

function scoreColorClass(score: number, full: number) {
  const rate = score / full
  if (rate >= 0.9) return 'text-emerald-500 dark:text-emerald-400'
  if (rate >= 0.8) return 'text-indigo-500 dark:text-indigo-400'
  if (rate >= 0.6) return 'text-amber-500 dark:text-amber-400'
  return 'text-rose-500 dark:text-rose-400'
}

async function ensureStorageInfo() {
  if (!storageInfo.value) {
    storageInfo.value = await getStorageInfo()
  }
  return storageInfo.value
}

async function resolvePhoto(url: string) {
  if (!url) {
    displayPhoto.value = ''
    return
  }
  if (/^https?:\/\//i.test(url)) {
    displayPhoto.value = url
    return
  }
  try {
    displayPhoto.value = await presignView(url)
  } catch {
    displayPhoto.value = ''
  }
}

async function load() {
  if (!scoreId.value) return
  const resp = await getScore(scoreId.value)
  record.value = resp
  formSubject.value = resp.subject
  formExamName.value = resp.exam_name
  formExamDate.value = dayjs(resp.exam_date).format('YYYY-MM-DD')
  formScore.value = resp.score
  formFullScore.value = resp.full_score
  formRank.value = resp.class_rank ?? null
  formClassAvg.value = resp.class_avg ?? null
  formRemark.value = resp.remark
  photoPreview.value = ''
  await resolvePhoto(resp.photo_url)
}

function openEdit() {
  if (!record.value) return
  showDialog.value = true
}

function triggerFileInput() {
  fileInputRef.value?.click()
}

async function onSelectPhoto(e: Event) {
  const input = e.target as HTMLInputElement
  const f = input.files && input.files[0]
  if (!f) return
  const converted = await toWebp(f)
  photoFile.value = converted
  photoPreview.value = URL.createObjectURL(converted)
}

async function uploadPhoto(file: File): Promise<string> {
  const info = await ensureStorageInfo()
  const ext = (file.name.split('.').pop() || 'webp').toLowerCase()
  const contentType = file.type || 'image/webp'
  const userId = auth.user?.id || 0
  const resp = await presignUpload({ resource_type: 'score_image', user_id: userId, content_type: contentType, ext })
  await putToURL(resp.upload_url, file, resp.headers)
  if (info?.public_base_url) {
    return `${info.public_base_url.replace(/\/+$/, '')}/${resp.object_key}`
  }
  return resp.object_key
}

async function submit() {
  if (!record.value) return
  if (!formSubject.value || !formExamName.value.trim() || !formExamDate.value) {
    ElMessage.error('请填写学科、考试名称与日期')
    return
  }
  if (formScore.value === null || formFullScore.value === null) {
    ElMessage.error('请填写分数与满分')
    return
  }
  if (formScore.value < 0 || formScore.value > formFullScore.value) {
    ElMessage.error('分数范围错误')
    return
  }
  saving.value = true
  try {
    let photoUrl = record.value.photo_url
    if (photoFile.value) {
      photoUrl = await uploadPhoto(photoFile.value)
    }
    const resp = await updateScore(record.value.id, {
      subject: formSubject.value,
      exam_name: formExamName.value.trim(),
      exam_date: formExamDate.value,
      score: formScore.value,
      full_score: formFullScore.value,
      class_rank: formRank.value,
      class_avg: formClassAvg.value,
      remark: formRemark.value.trim(),
      photo_url: photoUrl
    })
    record.value = resp
    showDialog.value = false
    await resolvePhoto(resp.photo_url)
  } catch (e: any) {
    ElMessage.error(e?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

async function onDelete() {
  if (!record.value) return
  try {
    await ElMessageBox.confirm('确认删除该成绩吗？', '提示', { type: 'warning' })
  } catch {
    return
  }
  await deleteScore(record.value.id)
  ElMessage.success('已删除')
  router.back()
}

function previewImage(url: string) {
  imageViewerUrl.value = url
  showImageViewer.value = true
}

onMounted(load)
</script>

<style scoped>
/* Custom Scrollbar for image viewer if needed */
</style>
