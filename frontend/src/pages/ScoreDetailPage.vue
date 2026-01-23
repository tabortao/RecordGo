<template>
  <div class="min-h-screen bg-[#F5F7FA] dark:bg-gray-900">
    <div class="bg-white/80 dark:bg-gray-900/80 backdrop-blur border-b border-gray-200 dark:border-gray-800 sticky top-0 z-10">
      <div class="px-4 py-3 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <el-icon :size="18" class="cursor-pointer text-gray-600 dark:text-gray-300" @click="router.back()"><ArrowLeft /></el-icon>
          <h2 class="font-semibold text-gray-800 dark:text-gray-100">成绩详情</h2>
        </div>
        <div class="flex items-center gap-2">
          <el-button size="small" @click="openEdit">编辑</el-button>
          <el-button size="small" type="danger" @click="onDelete">删除</el-button>
        </div>
      </div>
    </div>

    <div class="p-4 space-y-4">
      <div v-if="!record" class="text-sm text-gray-500 dark:text-gray-400">加载中...</div>
      <div v-else class="rounded-2xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4 shadow-sm space-y-3">
        <div class="w-full max-h-[60vh] min-h-[220px] rounded-xl overflow-hidden bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
          <img v-if="displayPhoto" :src="displayPhoto" class="max-h-[60vh] w-auto max-w-full object-contain" />
        </div>
        <div class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ record.exam_name }}</div>
        <div class="text-sm text-gray-600 dark:text-gray-300">学科：{{ record.subject }}</div>
        <div class="text-sm text-gray-600 dark:text-gray-300">考试日期：{{ formatDate(record.exam_date) }}</div>
        <div class="text-sm text-gray-600 dark:text-gray-300">分数：{{ record.score }} / {{ record.full_score }}</div>
        <div class="text-sm text-gray-600 dark:text-gray-300">排名：{{ record.rank ?? '-' }}</div>
        <div class="text-sm text-gray-600 dark:text-gray-300">班级均分：{{ record.class_avg ?? '-' }}</div>
        <div class="text-sm text-gray-500 dark:text-gray-400">备注：{{ record.remark || '无' }}</div>
      </div>
    </div>

    <el-dialog v-model="showDialog" title="编辑成绩" width="520px">
      <div class="space-y-3">
        <div class="space-y-2">
          <label class="text-sm text-gray-600 dark:text-gray-300">学科</label>
          <el-select v-model="formSubject" placeholder="请选择学科" class="w-full">
            <el-option v-for="s in subjectOptions" :key="s" :label="s" :value="s" />
          </el-select>
        </div>
        <div class="space-y-2">
          <label class="text-sm text-gray-600 dark:text-gray-300">考试名称</label>
          <el-input v-model="formExamName" placeholder="请输入考试名称" />
        </div>
        <div class="space-y-2">
          <label class="text-sm text-gray-600 dark:text-gray-300">考试日期</label>
          <el-date-picker v-model="formExamDate" type="date" value-format="YYYY-MM-DD" placeholder="选择日期" class="w-full" />
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">分数</label>
            <el-input-number v-model="formScore" :min="0" class="w-full" />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">满分</label>
            <el-input-number v-model="formFullScore" :min="1" class="w-full" />
          </div>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">排名</label>
            <el-input-number v-model="formRank" :min="1" class="w-full" />
          </div>
          <div class="space-y-2">
            <label class="text-sm text-gray-600 dark:text-gray-300">班级均分</label>
            <el-input-number v-model="formClassAvg" :min="0" class="w-full" />
          </div>
        </div>
        <div class="space-y-2">
          <label class="text-sm text-gray-600 dark:text-gray-300">备注</label>
          <el-input v-model="formRemark" type="textarea" :rows="3" placeholder="可填写补充说明" />
        </div>
        <div class="space-y-2">
          <label class="text-sm text-gray-600 dark:text-gray-300">成绩照片（可选）</label>
          <div class="flex items-center gap-3">
            <div class="w-16 h-16 rounded-lg overflow-hidden bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
              <img v-if="photoPreview" :src="photoPreview" class="w-full h-full object-cover" />
            </div>
            <input type="file" accept="image/*" @change="onSelectPhoto" />
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="showDialog = false">取消</el-button>
          <el-button type="primary" :loading="saving" @click="submit">保存</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ArrowLeft } from '@element-plus/icons-vue'
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

const subjectOptions = [
  '语文', '数学', '英语', '科学', '体育', '物理', '化学', '政治', '历史', '地理', '综合', '自定义'
]

const scoreId = computed(() => Number(route.params.id || 0))

function formatDate(d: string) {
  return dayjs(d).format('YYYY-MM-DD')
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
  formExamDate.value = formatDate(resp.exam_date)
  formScore.value = resp.score
  formFullScore.value = resp.full_score
  formRank.value = resp.rank ?? null
  formClassAvg.value = resp.class_avg ?? null
  formRemark.value = resp.remark
  photoPreview.value = ''
  await resolvePhoto(resp.photo_url)
}

function openEdit() {
  if (!record.value) return
  showDialog.value = true
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
      rank: formRank.value,
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

onMounted(load)
</script>

<style scoped>
</style>
