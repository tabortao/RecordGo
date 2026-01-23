<template>
  <div class="min-h-screen bg-[#F5F7FA] dark:bg-gray-900">
    <div class="bg-white/80 dark:bg-gray-900/80 backdrop-blur border-b border-gray-200 dark:border-gray-800 sticky top-0 z-10">
      <div class="px-4 py-3 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <el-icon :size="18" class="cursor-pointer text-gray-600 dark:text-gray-300" @click="router.back()"><ArrowLeft /></el-icon>
          <h2 class="font-semibold text-gray-800 dark:text-gray-100">新增成绩</h2>
        </div>
      </div>
    </div>

    <div class="p-4 space-y-4">
      <div class="rounded-2xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4 shadow-sm space-y-4">
        <div class="text-sm font-semibold text-gray-800 dark:text-gray-100">基本信息</div>
        <div class="grid grid-cols-1 gap-3 md:grid-cols-2">
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">科目<span class="text-red-500">*</span></label>
            <el-select v-model="formSubject" placeholder="请选择科目" class="w-full">
              <el-option v-for="s in subjectOptions" :key="s" :label="s" :value="s" />
            </el-select>
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">考试名称</label>
            <el-input v-model="formExamName" placeholder="不填写将根据选项自动生成" />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">考试类型</label>
            <el-select v-model="formExamType" placeholder="请选择考试类型" class="w-full">
              <el-option v-for="t in examTypeOptions" :key="t" :label="t" :value="t" />
            </el-select>
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">年级</label>
            <el-select v-model="formGrade" placeholder="请选择年级" class="w-full">
              <el-option v-for="g in gradeOptions" :key="g" :label="g" :value="g" />
            </el-select>
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">学期</label>
            <el-select v-model="formTerm" placeholder="请选择学期" class="w-full">
              <el-option v-for="t in termOptions" :key="t" :label="t" :value="t" />
            </el-select>
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">专题/章节</label>
            <el-input v-model="formTopic" placeholder="可填写专题或章节" />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">难度</label>
            <el-select v-model="formDifficulty" placeholder="请选择难度" class="w-full">
              <el-option v-for="d in difficultyOptions" :key="d" :label="d" :value="d" />
            </el-select>
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">考试日期</label>
            <el-date-picker v-model="formExamDate" type="date" value-format="YYYY-MM-DD" placeholder="选择日期" class="w-full" />
          </div>
        </div>
      </div>

      <div class="rounded-2xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4 shadow-sm space-y-4">
        <div class="text-sm font-semibold text-gray-800 dark:text-gray-100">成绩信息</div>
        <div class="grid grid-cols-1 gap-3 md:grid-cols-2">
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">得分<span class="text-red-500">*</span></label>
            <el-input-number v-model="formScore" :min="0" class="w-full" />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">满分<span class="text-red-500">*</span></label>
            <el-input-number v-model="formFullScore" :min="1" class="w-full" />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">班级排名</label>
            <el-input-number v-model="formClassRank" :min="1" class="w-full" />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">排名类型</label>
            <el-select v-model="formRankType" placeholder="请选择排名类型" class="w-full">
              <el-option v-for="t in rankTypeOptions" :key="t" :label="t" :value="t" />
            </el-select>
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">年级排名</label>
            <el-input-number v-model="formGradeRank" :min="1" class="w-full" />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">班级平均分</label>
            <el-input-number v-model="formClassAvg" :min="0" class="w-full" />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">班级最高分</label>
            <el-input-number v-model="formClassHighest" :min="0" class="w-full" />
          </div>
        </div>
        <div class="space-y-2">
          <label class="text-sm text-gray-600 dark:text-gray-300">备注</label>
          <el-input v-model="formRemark" type="textarea" :rows="3" placeholder="可填写补充说明" />
        </div>
        <div class="space-y-2">
          <label class="text-sm text-gray-600 dark:text-gray-300">成绩照片</label>
          <div class="flex items-center gap-3">
            <div class="w-16 h-16 rounded-lg overflow-hidden bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
              <img v-if="photoPreview" :src="photoPreview" class="w-full h-full object-cover" />
            </div>
            <input type="file" accept="image/*" @change="onSelectPhoto" />
          </div>
        </div>
      </div>

      <div class="flex justify-end gap-2">
        <el-button @click="router.back()">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submit">保存</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
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
const storageInfo = ref<{ driver: string; public_base_url?: string } | null>(null)

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
</script>

<style scoped>
</style>
