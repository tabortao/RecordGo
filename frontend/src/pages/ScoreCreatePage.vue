<template>
  <div class="min-h-screen bg-[#F5F7FA] dark:bg-gray-900 flex flex-col">
    <div class="bg-white/80 dark:bg-gray-900/80 backdrop-blur border-b border-gray-200 dark:border-gray-800 sticky top-0 z-20">
      <div class="px-4 py-3 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <el-icon :size="18" class="cursor-pointer text-gray-600 dark:text-gray-300" @click="router.back()"><ArrowLeft /></el-icon>
          <h2 class="font-semibold text-gray-800 dark:text-gray-100">新增成绩</h2>
        </div>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto px-4 py-4 sm:px-6 sm:py-6 pb-32">
      <div class="max-w-5xl mx-auto grid grid-cols-1 lg:grid-cols-12 gap-6">
        <div class="lg:col-span-7 space-y-6">
          <div class="rounded-3xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-6 shadow-sm space-y-6">
            <div class="flex items-center justify-between gap-3">
              <div class="text-base font-black text-gray-900 dark:text-white flex items-center gap-2">
                <div class="w-1 h-4 bg-indigo-500 rounded-full"></div>
                基本信息
              </div>
              <div class="text-xs text-gray-400 dark:text-gray-500 truncate max-w-[55%]">
                {{ examNamePreview }}
              </div>
            </div>
            <div class="grid grid-cols-1 gap-5 sm:grid-cols-2">
              <div class="space-y-2">
                <label class="text-sm text-gray-600 dark:text-gray-300">科目<span class="text-red-500">*</span></label>
                <el-select v-model="formSubject" placeholder="请选择科目" class="w-full custom-input">
                  <el-option v-for="s in subjectOptions" :key="s" :label="s" :value="s" />
                </el-select>
              </div>
              <div class="space-y-2">
                <label class="text-sm text-gray-600 dark:text-gray-300">考试名称</label>
                <el-input v-model="formExamName" placeholder="不填写将根据选项自动生成" class="custom-input" />
              </div>
              <div class="space-y-2">
                <label class="text-sm text-gray-600 dark:text-gray-300">考试类型</label>
                <el-select v-model="formExamType" placeholder="请选择考试类型" class="w-full custom-input">
                  <el-option v-for="t in examTypeOptions" :key="t" :label="t" :value="t" />
                </el-select>
              </div>
              <div class="space-y-2">
                <label class="text-sm text-gray-600 dark:text-gray-300">年级</label>
                <el-select v-model="formGrade" placeholder="请选择年级" class="w-full custom-input">
                  <el-option v-for="g in gradeOptions" :key="g" :label="g" :value="g" />
                </el-select>
              </div>
              <div class="space-y-2">
                <label class="text-sm text-gray-600 dark:text-gray-300">学期</label>
                <el-select v-model="formTerm" placeholder="请选择学期" class="w-full custom-input">
                  <el-option v-for="t in termOptions" :key="t" :label="t" :value="t" />
                </el-select>
              </div>
              <div class="space-y-2">
                <label class="text-sm text-gray-600 dark:text-gray-300">专题/章节</label>
                <el-input v-model="formTopic" placeholder="可填写专题或章节" class="custom-input" />
              </div>
              <div class="space-y-2">
                <label class="text-sm text-gray-600 dark:text-gray-300">难度</label>
                <el-select v-model="formDifficulty" placeholder="请选择难度" class="w-full custom-input">
                  <el-option v-for="d in difficultyOptions" :key="d" :label="d" :value="d" />
                </el-select>
              </div>
              <div class="space-y-2">
                <label class="text-sm text-gray-600 dark:text-gray-300">考试日期</label>
                <el-date-picker v-model="formExamDate" type="date" value-format="YYYY-MM-DD" placeholder="选择日期" class="!w-full custom-input" />
              </div>
            </div>
          </div>

          <div class="rounded-3xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-6 shadow-sm space-y-6">
            <div class="text-base font-black text-gray-900 dark:text-white flex items-center gap-2">
              <div class="w-1 h-4 bg-indigo-500 rounded-full"></div>
              成绩信息
            </div>
            <div class="grid grid-cols-1 gap-5 sm:grid-cols-2">
              <div class="space-y-2">
                <label class="text-sm text-gray-600 dark:text-gray-300">得分<span class="text-red-500">*</span></label>
                <el-input-number v-model="formScore" :min="0" class="!w-full custom-input-number" :controls="false" />
              </div>
              <div class="space-y-2">
                <label class="text-sm text-gray-600 dark:text-gray-300">满分<span class="text-red-500">*</span></label>
                <el-input-number v-model="formFullScore" :min="1" class="!w-full custom-input-number" :controls="false" />
              </div>
              <div class="space-y-2">
                <label class="text-sm text-gray-600 dark:text-gray-300">班级排名</label>
                <el-input-number v-model="formClassRank" :min="1" class="!w-full custom-input-number" :controls="false" />
              </div>
              <div class="space-y-2">
                <label class="text-sm text-gray-600 dark:text-gray-300">排名类型</label>
                <el-select v-model="formRankType" placeholder="请选择排名类型" class="w-full custom-input">
                  <el-option v-for="t in rankTypeOptions" :key="t" :label="t" :value="t" />
                </el-select>
              </div>
              <div class="space-y-2">
                <label class="text-sm text-gray-600 dark:text-gray-300">年级排名</label>
                <el-input-number v-model="formGradeRank" :min="1" class="!w-full custom-input-number" :controls="false" />
              </div>
              <div class="space-y-2">
                <label class="text-sm text-gray-600 dark:text-gray-300">班级平均分</label>
                <el-input-number v-model="formClassAvg" :min="0" class="!w-full custom-input-number" :controls="false" />
              </div>
              <div class="space-y-2 sm:col-span-2">
                <label class="text-sm text-gray-600 dark:text-gray-300">班级最高分</label>
                <el-input-number v-model="formClassHighest" :min="0" class="!w-full custom-input-number" :controls="false" />
              </div>
            </div>
            <div class="space-y-2">
              <label class="text-sm text-gray-600 dark:text-gray-300">备注</label>
              <el-input v-model="formRemark" type="textarea" :rows="3" placeholder="可填写补充说明" class="custom-input" resize="none" />
            </div>
          </div>
        </div>

        <div class="lg:col-span-5 space-y-6 lg:sticky lg:top-24 h-fit">
          <div class="rounded-3xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-6 shadow-sm overflow-hidden relative">
            <div class="absolute -right-12 -top-12 w-40 h-40 rounded-full bg-indigo-50 dark:bg-indigo-900/20 blur-2xl"></div>
            <div class="absolute -left-12 -bottom-12 w-40 h-40 rounded-full bg-emerald-50 dark:bg-emerald-900/20 blur-2xl"></div>
            <div class="relative">
              <div class="flex items-start justify-between gap-4">
                <div>
                  <div class="text-sm font-bold text-gray-500 dark:text-gray-400">成绩预览</div>
                  <div class="mt-1 text-xl font-black text-gray-900 dark:text-white tracking-tight">
                    {{ scorePreviewText }}
                  </div>
                </div>
                <div class="px-3 py-1 rounded-full text-xs font-bold" :class="rateBadgeClass">
                  {{ rateText }}
                </div>
              </div>

              <div class="mt-5">
                <div class="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400 mb-2">
                  <span>得分率</span>
                  <span class="font-bold">{{ ratePercentText }}</span>
                </div>
                <div class="h-2 rounded-full bg-gray-100 dark:bg-gray-700 overflow-hidden">
                  <div class="h-full rounded-full transition-all duration-700" :class="rateBarClass" :style="{ width: rateWidth }"></div>
                </div>
              </div>

              <div class="mt-5 grid grid-cols-2 gap-3 text-sm">
                <div class="rounded-2xl bg-gray-50 dark:bg-gray-900/40 border border-gray-100 dark:border-gray-700/60 p-3">
                  <div class="text-xs text-gray-400 mb-1">班级排名</div>
                  <div class="font-black text-gray-900 dark:text-white">{{ formClassRank ?? '-' }}</div>
                </div>
                <div class="rounded-2xl bg-gray-50 dark:bg-gray-900/40 border border-gray-100 dark:border-gray-700/60 p-3">
                  <div class="text-xs text-gray-400 mb-1">年级排名</div>
                  <div class="font-black text-gray-900 dark:text-white">{{ formGradeRank ?? '-' }}</div>
                </div>
              </div>
            </div>
          </div>

          <div class="rounded-3xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-6 shadow-sm space-y-4">
            <div class="flex items-center justify-between">
              <div class="text-base font-black text-gray-900 dark:text-white flex items-center gap-2">
                <div class="w-1 h-4 bg-indigo-500 rounded-full"></div>
                成绩照片
              </div>
              <button
                class="text-xs font-bold text-indigo-600 dark:text-indigo-300 hover:text-indigo-700 dark:hover:text-indigo-200 transition-colors"
                type="button"
                @click="triggerFileInput"
              >
                {{ photoPreview ? '更换' : '上传' }}
              </button>
            </div>

            <div
              class="group relative rounded-2xl border border-dashed border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-900/40 overflow-hidden cursor-pointer hover:border-indigo-400 dark:hover:border-indigo-500 transition-colors"
              @click="triggerFileInput"
            >
              <div v-if="photoPreview" class="aspect-[4/3] bg-gray-100 dark:bg-gray-950/30">
                <img :src="photoPreview" class="w-full h-full object-cover" />
              </div>
              <div v-else class="aspect-[4/3] flex flex-col items-center justify-center text-gray-400">
                <div class="w-12 h-12 rounded-full bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 flex items-center justify-center mb-2 shadow-sm">
                  <span class="text-lg font-black text-gray-700 dark:text-gray-200">+</span>
                </div>
                <div class="text-sm font-bold">点击上传</div>
                <div class="text-xs mt-1">将自动压缩为 webp</div>
              </div>

              <button
                v-if="photoPreview"
                class="absolute top-3 right-3 w-8 h-8 rounded-full bg-white/90 dark:bg-gray-800/90 backdrop-blur border border-gray-100 dark:border-gray-700 flex items-center justify-center text-gray-500 hover:text-rose-500 transition-colors"
                type="button"
                @click.stop="clearPhoto"
              >
                ×
              </button>
            </div>
            <input ref="fileInputRef" type="file" class="hidden" accept="image/*" @change="onSelectPhoto" />
          </div>
        </div>
      </div>
    </div>

    <div class="fixed bottom-0 left-0 right-0 z-20 bg-white/85 dark:bg-gray-900/85 backdrop-blur border-t border-gray-200 dark:border-gray-800">
      <div class="max-w-5xl mx-auto px-4 sm:px-6 py-3 pb-[calc(12px+env(safe-area-inset-bottom))] flex items-center justify-end gap-3">
        <el-button @click="router.back()" class="!rounded-xl px-6 !h-10">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submit" class="!rounded-xl px-6 !h-10 !bg-indigo-600 !border-indigo-600 hover:!bg-indigo-700">保存</el-button>
      </div>
    </div>
  </div>
</template>


<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue'
import { ArrowLeft } from '@element-plus/icons-vue'
import router from '@/router'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'
import { createScore } from '@/services/scores'
import { getStorageInfo, presignUpload, putToURL } from '@/services/storage'
import { toWebp } from '@/utils/image'
import { useAuth } from '@/stores/auth'

const auth = useAuth()
const saving = ref(false)
const formSubject = ref('语文')
const formExamName = ref('')
const formExamType = ref('')
const formGrade = ref('')
const formTerm = ref('')
const formTopic = ref('')
const formDifficulty = ref('')
const formExamDate = ref<string>(dayjs().format('YYYY-MM-DD'))
const formScore = ref<number | null>(null)
const formFullScore = ref<number | null>(100)
const formClassRank = ref<number | null>(null)
const formRankType = ref('')
const formGradeRank = ref<number | null>(null)
const formClassAvg = ref<number | null>(null)
const formClassHighest = ref<number | null>(null)
const formRemark = ref('')
const photoFile = ref<File | null>(null)
const photoPreview = ref('')
const fileInputRef = ref<HTMLInputElement | null>(null)
const storageInfo = ref<{ driver: string; public_base_url?: string } | null>(null)

const examNamePreview = computed(() => {
  const manual = formExamName.value.trim()
  if (manual) return manual
  const auto = `${formGrade.value}${formTerm.value}${formSubject.value}${formExamType.value}`
  return auto || '未命名考试'
})

const rate = computed(() => {
  if (formScore.value === null || formFullScore.value === null || formFullScore.value <= 0) return null
  const r = formScore.value / formFullScore.value
  if (!Number.isFinite(r)) return null
  return Math.max(0, Math.min(1, r))
})

const scorePreviewText = computed(() => {
  if (formScore.value === null || formFullScore.value === null) return '填写分数后自动预览'
  return `${formScore.value} / ${formFullScore.value}`
})

const ratePercentText = computed(() => {
  if (rate.value === null) return '-'
  return `${Math.round(rate.value * 100)}%`
})

const rateWidth = computed(() => {
  if (rate.value === null) return '0%'
  return `${Math.round(rate.value * 100)}%`
})

const rateText = computed(() => {
  if (rate.value === null) return '未填写'
  if (rate.value >= 0.9) return '优秀'
  if (rate.value >= 0.8) return '良好'
  if (rate.value >= 0.6) return '合格'
  return '待提升'
})

const rateBadgeClass = computed(() => {
  if (rate.value === null) return 'bg-gray-100 text-gray-500 dark:bg-gray-700/60 dark:text-gray-300'
  if (rate.value >= 0.9) return 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300'
  if (rate.value >= 0.8) return 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-300'
  if (rate.value >= 0.6) return 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300'
  return 'bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-300'
})

const rateBarClass = computed(() => {
  if (rate.value === null) return 'bg-gray-300 dark:bg-gray-600'
  if (rate.value >= 0.9) return 'bg-gradient-to-r from-emerald-500 to-teal-500'
  if (rate.value >= 0.8) return 'bg-gradient-to-r from-indigo-500 to-violet-500'
  if (rate.value >= 0.6) return 'bg-gradient-to-r from-amber-500 to-orange-500'
  return 'bg-gradient-to-r from-rose-500 to-pink-500'
})

const subjectOptions = [
  '语文', '数学', '英语', '物理', '化学', '生物', '历史', '地理', '政治', '信息技术', '科学'
]
const examTypeOptions = ['月考', '期中考试', '期末考试', '模拟考试', '单元测试', '周测', '随堂测验']
const gradeOptions = ['一年级', '二年级', '三年级', '四年级', '五年级', '六年级', '初一', '初二', '初三', '高一', '高二', '高三']
const termOptions = ['上学期', '下学期']
const difficultyOptions = ['简单', '中等', '困难', '较难']
const rankTypeOptions = ['年级排名', '校排名', '区排名', '市排名']

async function ensureStorageInfo() {
  if (!storageInfo.value) {
    storageInfo.value = await getStorageInfo()
  }
  return storageInfo.value
}

function triggerFileInput() {
  fileInputRef.value?.click()
}

function revokePreviewUrl() {
  if (photoPreview.value) {
    URL.revokeObjectURL(photoPreview.value)
  }
}

function clearPhoto() {
  revokePreviewUrl()
  photoFile.value = null
  photoPreview.value = ''
  if (fileInputRef.value) fileInputRef.value.value = ''
}

async function onSelectPhoto(e: Event) {
  const input = e.target as HTMLInputElement
  const f = input.files && input.files[0]
  if (!f) return
  const converted = await toWebp(f)
  revokePreviewUrl()
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
  if (!formSubject.value) {
    ElMessage.error('请选择科目')
    return
  }
  
  // 如果考试名称为空，则自动生成
  let examName = formExamName.value.trim()
  if (!examName) {
    examName = `${formGrade.value}${formTerm.value}${formSubject.value}${formExamType.value}`
    // 清理可能出现的空字符串
    if (examName === '') {
      examName = '未命名考试'
    }
  }

  if (!formExamDate.value) {
    ElMessage.error('请选择考试日期')
    return
  }
  if (formScore.value === null || formFullScore.value === null) {
    ElMessage.error('请填写得分与满分')
    return
  }
  if (formScore.value < 0 || formScore.value > formFullScore.value) {
    ElMessage.error('分数范围错误')
    return
  }
  saving.value = true
  try {
    let photoUrl = ''
    if (photoFile.value) {
      photoUrl = await uploadPhoto(photoFile.value)
    }
    await createScore({
      subject: formSubject.value,
      exam_name: examName,
      exam_type: formExamType.value,
      grade: formGrade.value,
      term: formTerm.value,
      topic: formTopic.value.trim(),
      difficulty: formDifficulty.value,
      exam_date: formExamDate.value,
      score: formScore.value,
      full_score: formFullScore.value,
      class_rank: formClassRank.value,
      rank_type: formRankType.value,
      grade_rank: formGradeRank.value,
      class_avg: formClassAvg.value,
      class_highest: formClassHighest.value,
      remark: formRemark.value.trim(),
      photo_url: photoUrl
    })
    router.replace('/grades')
  } catch (e: any) {
    ElMessage.error(e?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

onBeforeUnmount(() => {
  revokePreviewUrl()
})
</script>

<style scoped>
.custom-input :deep(.el-input__wrapper) {
  background-color: #f9fafb;
  border-radius: 0.75rem;
  box-shadow: none !important;
  border: 1px solid #e5e7eb;
  padding: 8px 12px;
  transition: all 0.2s;
}
.dark .custom-input :deep(.el-input__wrapper) {
  background-color: #374151;
  border-color: #4b5563;
}
.custom-input :deep(.el-input__wrapper:hover),
.custom-input :deep(.el-input__wrapper.is-focus) {
  border-color: #6366f1;
  background-color: #fff;
}
.dark .custom-input :deep(.el-input__wrapper:hover),
.dark .custom-input :deep(.el-input__wrapper.is-focus) {
  border-color: #6366f1;
  background-color: #1f2937;
}

.custom-input :deep(.el-textarea__inner) {
  background-color: #f9fafb;
  border-radius: 0.75rem;
  box-shadow: none !important;
  border: 1px solid #e5e7eb;
  padding: 12px;
}
.dark .custom-input :deep(.el-textarea__inner) {
  background-color: #374151;
  border-color: #4b5563;
  color: #fff;
}
.custom-input :deep(.el-textarea__inner:hover),
.custom-input :deep(.el-textarea__inner:focus) {
  border-color: #6366f1;
  background-color: #fff;
}
.dark .custom-input :deep(.el-textarea__inner:hover),
.dark .custom-input :deep(.el-textarea__inner:focus) {
  border-color: #6366f1;
  background-color: #1f2937;
}

.custom-input-number :deep(.el-input__wrapper) {
  background-color: #f9fafb;
  border-radius: 0.75rem;
  box-shadow: none !important;
  border: 1px solid #e5e7eb;
}
.dark .custom-input-number :deep(.el-input__wrapper) {
  background-color: #374151;
  border-color: #4b5563;
}
.custom-input-number :deep(.el-input__wrapper:hover),
.custom-input-number :deep(.el-input__wrapper.is-focus) {
  border-color: #6366f1;
}
.dark .custom-input-number :deep(.el-input__wrapper:hover),
.dark .custom-input-number :deep(.el-input__wrapper.is-focus) {
  border-color: #6366f1;
}
</style>
