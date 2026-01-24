<template>
  <div class="min-h-screen bg-[#F5F7FA] dark:bg-gray-900 flex flex-col relative">
    <!-- Header -->
    <div class="px-4 py-3 bg-white/80 dark:bg-gray-900/80 backdrop-blur-md sticky top-0 z-30 flex items-center justify-between border-b border-gray-100 dark:border-gray-800 transition-colors">
      <div 
        class="w-8 h-8 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center cursor-pointer hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors text-gray-600 dark:text-gray-300"
        @click="router.back()"
      >
        <el-icon><ArrowLeft /></el-icon>
      </div>
      
      <div class="flex items-center gap-2">
        <button 
          class="w-8 h-8 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center text-gray-500 hover:text-amber-500 hover:bg-amber-50 dark:hover:bg-amber-900/20 transition-all"
          @click="openEdit"
        >
          <el-icon><Edit /></el-icon>
        </button>
        <button 
          class="w-8 h-8 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center text-gray-500 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 transition-all"
          @click="onDelete"
        >
          <el-icon><Delete /></el-icon>
        </button>
      </div>
    </div>

    <div class="flex-1 overflow-y-auto p-4 sm:p-6 pb-24">
      <div v-if="!record" class="flex flex-col items-center justify-center h-64 text-gray-400">
        <el-icon class="animate-spin text-2xl mb-2"><Loading /></el-icon>
        <span>加载中...</span>
      </div>
      
      <div v-else class="max-w-2xl mx-auto space-y-6">
        <!-- Main Image Card -->
        <div class="bg-white dark:bg-gray-800 rounded-3xl overflow-hidden shadow-lg border border-gray-100 dark:border-gray-700 group relative">
          <div class="aspect-[4/3] bg-gray-100 dark:bg-gray-900 flex items-center justify-center relative">
            <img 
              v-if="displayPhoto" 
              :src="displayPhoto" 
              class="w-full h-full object-contain cursor-zoom-in transition-transform duration-500 group-hover:scale-105"
              @click="previewImage(displayPhoto)"
            />
            <el-icon v-else :size="48" class="text-gray-300 dark:text-gray-700"><Picture /></el-icon>
            
            <!-- Date Badge -->
            <div class="absolute top-4 left-4 bg-white/90 dark:bg-gray-800/90 backdrop-blur-sm px-3 py-1.5 rounded-full text-xs font-bold text-gray-700 dark:text-gray-200 shadow-sm flex items-center gap-1.5">
              <el-icon class="text-amber-500"><Calendar /></el-icon>
              {{ formatDate(record.awarded_at) }}
            </div>
          </div>
          
          <div class="p-6">
            <div class="flex items-start justify-between gap-4 mb-4">
              <div>
                <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-2 leading-tight">{{ record.title }}</h1>
                <div class="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
                  <el-icon><School /></el-icon>
                  <span>颁发机构：{{ record.issuer || '未填写' }}</span>
                </div>
              </div>
              <div class="w-12 h-12 rounded-full bg-amber-50 dark:bg-amber-900/20 flex items-center justify-center text-amber-500 flex-shrink-0">
                <el-icon :size="24"><Trophy /></el-icon>
              </div>
            </div>
            
            <div v-if="record.remark" class="bg-gray-50 dark:bg-gray-700/50 rounded-xl p-4 text-gray-600 dark:text-gray-300 text-sm leading-relaxed border border-gray-100 dark:border-gray-700/50">
              <div class="font-bold text-xs text-gray-400 mb-1 flex items-center gap-1">
                <el-icon><ChatDotSquare /></el-icon> 备注
              </div>
              {{ record.remark }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Dialog -->
    <el-dialog 
      v-model="showDialog" 
      title="编辑荣誉" 
      width="90%" 
      class="max-w-md rounded-2xl overflow-hidden"
      :show-close="false"
      align-center
    >
      <template #header>
        <div class="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
          <div class="w-1 bg-amber-500 h-4 rounded-full"></div>
          编辑荣誉信息
        </div>
      </template>
      
      <div class="space-y-5 py-2">
        <!-- Photo Upload -->
        <div class="flex justify-center">
          <div class="relative group cursor-pointer" @click="triggerFileSelect">
            <div 
              class="w-32 h-32 rounded-2xl bg-gray-50 dark:bg-gray-700/50 border-2 border-dashed border-gray-200 dark:border-gray-600 flex flex-col items-center justify-center overflow-hidden transition-colors hover:border-amber-400 hover:bg-amber-50 dark:hover:bg-amber-900/10"
              :class="{'border-amber-500 bg-amber-50 dark:bg-amber-900/10': photoPreview || displayPhoto}"
            >
              <img v-if="photoPreview" :src="photoPreview" class="w-full h-full object-cover" />
              <img v-else-if="displayPhoto && !photoFile" :src="displayPhoto" class="w-full h-full object-cover opacity-80" />
              <template v-else>
                <el-icon :size="28" class="text-gray-400 group-hover:text-amber-500 transition-colors"><Camera /></el-icon>
                <span class="text-xs text-gray-400 mt-2 font-medium group-hover:text-amber-500">更换照片</span>
              </template>
            </div>
            <div v-if="photoPreview" class="absolute -top-2 -right-2 bg-white dark:bg-gray-800 text-red-500 rounded-full p-1 shadow-md border border-gray-100 dark:border-gray-700 z-10" @click.stop="clearPhoto">
              <el-icon><Close /></el-icon>
            </div>
          </div>
          <input type="file" ref="fileInput" class="hidden" accept="image/*" @change="onSelectPhoto" />
        </div>

        <div class="space-y-4">
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">荣誉名称</label>
            <el-input 
              v-model="formTitle" 
              placeholder="例如：三好学生奖状" 
              class="custom-input"
            />
          </div>
          
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">颁发机构</label>
              <el-input 
                v-model="formIssuer" 
                placeholder="例如：xx学校" 
                class="custom-input"
              />
            </div>
            <div class="space-y-1.5">
              <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">颁发日期</label>
              <el-date-picker 
                v-model="formDate" 
                type="date" 
                value-format="YYYY-MM-DD" 
                placeholder="选择日期" 
                class="!w-full custom-input" 
                :clearable="false"
              />
            </div>
          </div>

          <div class="space-y-1.5">
            <label class="text-xs font-bold text-gray-500 dark:text-gray-400 ml-1">备注说明</label>
            <el-input 
              v-model="formRemark" 
              type="textarea" 
              :rows="3" 
              placeholder="记录一下当时的心情吧..." 
              class="custom-input"
              resize="none"
            />
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex gap-3 pt-2">
          <button 
            class="flex-1 py-2.5 rounded-xl bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 font-bold text-sm hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
            @click="showDialog = false"
          >
            取消
          </button>
          <button 
            class="flex-1 py-2.5 rounded-xl bg-gradient-to-r from-amber-400 to-amber-500 text-white font-bold text-sm shadow-lg shadow-amber-200 dark:shadow-amber-900/20 hover:scale-[1.02] active:scale-95 transition-all disabled:opacity-50 disabled:hover:scale-100"
            :disabled="saving"
            @click="submit"
          >
            {{ saving ? '保存中...' : '保存修改' }}
          </button>
        </div>
      </template>
    </el-dialog>

    <!-- Image Preview -->
    <el-image-viewer 
      v-if="showImageViewer" 
      :url-list="[displayPhoto]" 
      @close="showImageViewer = false"
      :hide-on-click-modal="true"
    />
  </div>
</template>

<script setup lang="ts">
// 中文注释：荣誉详情与编辑逻辑
import { computed, onMounted, ref } from 'vue'
import { ArrowLeft, Delete, Edit, Picture, Calendar, School, Trophy, ChatDotSquare, Camera, Close, Loading } from '@element-plus/icons-vue'
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
const showImageViewer = ref(false)

const showDialog = ref(false)
const saving = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

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

function previewImage(url: string) {
  if (url) showImageViewer.value = true
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

function triggerFileSelect() {
  fileInput.value?.click()
}

function clearPhoto() {
  photoFile.value = null
  photoPreview.value = ''
  if (fileInput.value) fileInput.value.value = ''
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
    ElMessage.success('修改成功')
  } catch (e: any) {
    ElMessage.error(e?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

async function onDelete() {
  if (!record.value) return
  try {
    await ElMessageBox.confirm('确认删除该荣誉吗？', '提示', { type: 'warning', cancelButtonText: '取消', confirmButtonText: '删除', confirmButtonClass: 'el-button--danger' })
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
  border-color: #f59e0b;
  background-color: #fff;
}
.dark .custom-input :deep(.el-input__wrapper:hover),
.dark .custom-input :deep(.el-input__wrapper.is-focus) {
  border-color: #f59e0b;
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
  border-color: #f59e0b;
  background-color: #fff;
}
.dark .custom-input :deep(.el-textarea__inner:hover),
.dark .custom-input :deep(.el-textarea__inner:focus) {
  border-color: #f59e0b;
  background-color: #1f2937;
}
</style>
