<template>
  <div class="min-h-screen bg-[#F5F7FA] dark:bg-gray-900">
    <div class="bg-white/80 dark:bg-gray-900/80 backdrop-blur border-b border-gray-200 dark:border-gray-800 sticky top-0 z-10">
      <div class="px-4 py-3 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <el-icon :size="18" class="cursor-pointer text-gray-600 dark:text-gray-300" @click="router.back()"><ArrowLeft /></el-icon>
          <h2 class="font-semibold text-gray-800 dark:text-gray-100">荣誉详情</h2>
        </div>
        <div class="flex items-center gap-2">
          <el-button size="small" circle @click="openEdit">
            <el-icon><Edit /></el-icon>
          </el-button>
          <el-button size="small" type="danger" circle @click="onDelete">
            <el-icon><Delete /></el-icon>
          </el-button>
        </div>
      </div>
    </div>

    <div class="p-4 space-y-4">
      <div v-if="!record" class="text-sm text-gray-500 dark:text-gray-400">加载中...</div>
      <div v-else class="rounded-2xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4 shadow-sm space-y-3">
        <div class="w-full max-h-[60vh] min-h-[220px] rounded-xl overflow-hidden bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
          <img v-if="displayPhoto" :src="displayPhoto" class="max-h-[60vh] w-auto max-w-full object-contain" />
        </div>
        <div class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ record.title }}</div>
        <div class="text-sm text-gray-600 dark:text-gray-300">颁发机构：{{ record.issuer || '未填写' }}</div>
        <div class="text-sm text-gray-600 dark:text-gray-300">颁发日期：{{ formatDate(record.awarded_at) }}</div>
        <div class="text-sm text-gray-500 dark:text-gray-400">备注：{{ record.remark || '无' }}</div>
      </div>
    </div>

    <el-dialog v-model="showDialog" title="编辑荣誉" width="420px">
      <div class="space-y-3">
        <div class="space-y-2">
          <label class="text-sm text-gray-600 dark:text-gray-300">荣誉名称</label>
          <el-input v-model="formTitle" placeholder="请输入荣誉名称" />
        </div>
        <div class="space-y-2">
          <label class="text-sm text-gray-600 dark:text-gray-300">颁发机构</label>
          <el-input v-model="formIssuer" placeholder="请输入颁发机构" />
        </div>
        <div class="space-y-2">
          <label class="text-sm text-gray-600 dark:text-gray-300">颁发日期</label>
          <el-date-picker v-model="formDate" type="date" value-format="YYYY-MM-DD" placeholder="选择日期" class="w-full" />
        </div>
        <div class="space-y-2">
          <label class="text-sm text-gray-600 dark:text-gray-300">备注</label>
          <el-input v-model="formRemark" type="textarea" :rows="3" placeholder="可填写补充说明" />
        </div>
        <div class="space-y-2">
          <label class="text-sm text-gray-600 dark:text-gray-300">荣誉照片</label>
          <div class="flex items-center gap-3">
            <div class="w-16 h-16 rounded-lg overflow-hidden bg-gray-100 dark:bg-gray-700">
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
// 中文注释：荣誉详情与编辑逻辑
import { computed, onMounted, ref } from 'vue'
import { ArrowLeft, Delete, Edit } from '@element-plus/icons-vue'
import { useRoute } from 'vue-router'
import router from '@/router'
import dayjs from 'dayjs'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteHonor, getHonor, updateHonor, type HonorRecord } from '@/services/honors'
import { getStorageInfo, presignUpload, presignView, putToURL } from '@/services/storage'
import { toWebp } from '@/utils/image'
import { useAuth } from '@/stores/auth'

const auth = useAuth()
const route = useRoute()
const record = ref<HonorRecord | null>(null)
const displayPhoto = ref('')

const showDialog = ref(false)
const saving = ref(false)

const formTitle = ref('')
const formIssuer = ref('')
const formDate = ref<string>(dayjs().format('YYYY-MM-DD'))
const formRemark = ref('')
const photoFile = ref<File | null>(null)
const photoPreview = ref<string>('')
const storageInfo = ref<{ driver: string; public_base_url?: string } | null>(null)

const honorId = computed(() => Number(route.params.id || 0))

function formatDate(d: string) {
  return dayjs(d).format('YYYY-MM-DD')
}

async function ensureStorageInfo() {
  if (!storageInfo.value) {
    storageInfo.value = await getStorageInfo()
  }
  return storageInfo.value
}

// 中文注释：解析照片 URL，非 http 地址走预签名
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
  if (!honorId.value) return
  const resp = await getHonor(honorId.value)
  record.value = resp
  formTitle.value = resp.title
  formIssuer.value = resp.issuer
  formDate.value = formatDate(resp.awarded_at)
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

// 中文注释：上传荣誉照片并返回可存储的 URL
async function uploadPhoto(file: File): Promise<string> {
  const info = await ensureStorageInfo()
  const ext = (file.name.split('.').pop() || 'webp').toLowerCase()
  const contentType = file.type || 'image/webp'
  const userId = auth.user?.id || 0
  const resp = await presignUpload({ resource_type: 'honor_image', user_id: userId, content_type: contentType, ext })
  await putToURL(resp.upload_url, file, resp.headers)
  if (info?.public_base_url) {
    return `${info.public_base_url.replace(/\/+$/, '')}/${resp.object_key}`
  }
  return resp.object_key
}

async function submit() {
  if (!record.value) return
  if (!formTitle.value.trim() || !formDate.value) {
    ElMessage.error('请填写荣誉名称与日期')
    return
  }
  saving.value = true
  try {
    let photoUrl = record.value.photo_url
    if (photoFile.value) {
      photoUrl = await uploadPhoto(photoFile.value)
    }
    const resp = await updateHonor(record.value.id, {
      title: formTitle.value.trim(),
      issuer: formIssuer.value.trim(),
      awarded_at: formDate.value,
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
    await ElMessageBox.confirm('确认删除该荣誉吗？', '提示', { type: 'warning' })
  } catch {
    return
  }
  await deleteHonor(record.value.id)
  ElMessage.success('已删除')
  router.back()
}

onMounted(load)
</script>

<style scoped>
</style>
