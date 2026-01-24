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
        <h2 class="font-bold text-gray-800 dark:text-gray-100 text-lg">荣誉墙</h2>
      </div>
      <button 
        class="bg-gradient-to-r from-amber-500 via-orange-500 to-rose-500 text-white px-4 py-2 rounded-full text-sm font-black shadow-lg shadow-amber-200/60 dark:shadow-amber-900/20 ring-1 ring-white/30 dark:ring-white/10 hover:scale-105 active:scale-95 transition-all flex items-center gap-1 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-400"
        @click="openDialog"
      >
        <el-icon><Plus /></el-icon>
        <span>添加荣誉</span>
      </button>
    </div>

    <div class="flex-1 overflow-y-auto p-4 sm:p-6 pb-24">
      <div class="max-w-5xl mx-auto space-y-6">
        <!-- Profile Card -->
        <div class="bg-white dark:bg-gray-800 rounded-3xl p-6 shadow-sm border border-gray-100 dark:border-gray-700 flex flex-col sm:flex-row items-center sm:items-start gap-6 relative overflow-hidden">
          <!-- Decoration -->
          <div class="absolute top-0 right-0 w-32 h-32 bg-gradient-to-br from-amber-100 to-transparent dark:from-amber-900/20 rounded-bl-full -z-0 opacity-50"></div>
          
          <div class="relative z-10">
            <div class="relative">
              <el-avatar :size="80" :src="avatarSrc" class="border-4 border-white dark:border-gray-700 shadow-md" />
              <div class="absolute -bottom-1 -right-1 bg-amber-400 text-white p-1.5 rounded-full shadow-sm">
                <el-icon :size="16"><Trophy /></el-icon>
              </div>
            </div>
          </div>
          
          <div class="flex-1 text-center sm:text-left z-10">
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">{{ displayName }}</h1>
            <div class="flex flex-wrap items-center justify-center sm:justify-start gap-3 text-sm">
              <div class="px-3 py-1 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 font-medium">
                {{ ageText }}
              </div>
              <div class="px-3 py-1 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 font-medium">
                {{ genderText }}
              </div>
              <div class="px-3 py-1 rounded-full bg-amber-50 dark:bg-amber-900/30 text-amber-600 dark:text-amber-400 font-medium flex items-center gap-1">
                <el-icon><Medal /></el-icon>
                <span>共 {{ records.length }} 项荣誉</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="records.length === 0" class="flex flex-col items-center justify-center py-20 text-gray-400">
          <div class="w-24 h-24 bg-gray-100 dark:bg-gray-800 rounded-full flex items-center justify-center mb-4">
            <el-icon :size="40" class="text-gray-300 dark:text-gray-600"><Trophy /></el-icon>
          </div>
          <p>暂无荣誉记录，快去记录第一次高光时刻吧~</p>
        </div>

        <!-- Grid Layout -->
        <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          <div
            v-for="item in records"
            :key="item.id"
            class="group relative bg-white dark:bg-gray-800 rounded-3xl overflow-hidden shadow-sm hover:shadow-xl hover:-translate-y-1 transition-all duration-300 border border-gray-100 dark:border-gray-700 cursor-pointer flex flex-col"
            @click="openDetail(item.id)"
          >
            <!-- Image Area -->
            <div class="aspect-[4/3] bg-gray-100 dark:bg-gray-900 relative overflow-hidden">
              <img
                v-if="resolvePhoto(item.photo_url)"
                :src="resolvePhoto(item.photo_url)"
                class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110"
              />
              <div class="absolute inset-0 flex items-center justify-center bg-gray-100 dark:bg-gray-800 text-gray-300" v-else>
                <el-icon :size="32"><Picture /></el-icon>
              </div>
              <!-- Overlay Date -->
              <div class="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/60 to-transparent p-4 pt-12">
                <div class="text-white text-xs font-medium flex items-center gap-1">
                  <el-icon><Calendar /></el-icon>
                  <span>{{ formatDate(item.awarded_at) }}</span>
                </div>
              </div>
            </div>

            <!-- Content Area -->
            <div class="p-4 flex-1 flex flex-col">
              <h3 class="text-base font-bold text-gray-900 dark:text-white line-clamp-2 mb-2 leading-tight group-hover:text-amber-500 transition-colors">
                {{ item.title }}
              </h3>
              <div class="mt-auto flex items-center gap-2 text-xs text-gray-500 dark:text-gray-400">
                <el-icon><School /></el-icon>
                <span class="truncate">{{ item.issuer || '未知机构' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Dialog -->
    <el-dialog 
      v-model="showDialog" 
      title="添加荣誉" 
      width="90%" 
      class="max-w-md rounded-2xl overflow-hidden flex flex-col max-h-[90vh]"
      :show-close="false"
      align-center
    >
      <template #header>
        <div class="text-lg font-bold text-gray-900 dark:text-white flex items-center gap-2">
          <div class="w-1 bg-amber-500 h-4 rounded-full"></div>
          添加新荣誉
        </div>
      </template>
      
      <div class="space-y-5 py-2 overflow-y-auto flex-1 px-1">
        <!-- Photo Upload -->
        <div class="flex justify-center">
          <div class="relative group cursor-pointer" @click="triggerFileSelect">
            <div 
              class="w-32 h-32 rounded-2xl bg-gray-50 dark:bg-gray-700/50 border-2 border-dashed border-gray-200 dark:border-gray-600 flex flex-col items-center justify-center overflow-hidden transition-colors hover:border-amber-400 hover:bg-amber-50 dark:hover:bg-amber-900/10"
              :class="{'border-amber-500 bg-amber-50 dark:bg-amber-900/10': photoPreview}"
            >
              <img v-if="photoPreview" :src="photoPreview" class="w-full h-full object-cover" />
              <template v-else>
                <el-icon :size="28" class="text-gray-400 group-hover:text-amber-500 transition-colors"><Camera /></el-icon>
                <span class="text-xs text-gray-400 mt-2 font-medium group-hover:text-amber-500">上传照片</span>
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
            class="flex-1 py-2.5 rounded-xl bg-gradient-to-r from-amber-500 via-orange-500 to-rose-500 text-white font-black text-sm shadow-lg shadow-amber-200/70 dark:shadow-amber-900/20 ring-1 ring-white/20 dark:ring-white/10 hover:scale-[1.02] active:scale-95 transition-all disabled:opacity-60 disabled:hover:scale-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-amber-400"
            :disabled="saving"
            @click="submit"
          >
            {{ saving ? '保存中...' : '保存荣誉' }}
          </button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { ArrowLeft, Plus, Trophy, Medal, Picture, Calendar, School, Camera, Close } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import router from '@/router'
import { ElMessage } from 'element-plus'
import { useAuth } from '@/stores/auth'
import { createHonor, listHonors, type HonorRecord } from '@/services/honors'
import { getStorageInfo, presignUpload, presignView, putToURL } from '@/services/storage'
import { toWebp } from '@/utils/image'
import defaultAvatar from '@/assets/avatars/default.png'
import { getStaticBase } from '@/services/http'

const auth = useAuth()
const records = ref<HonorRecord[]>([])
const showDialog = ref(false)
const saving = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

const formTitle = ref('')
const formIssuer = ref('')
const formDate = ref<string>(dayjs().format('YYYY-MM-DD'))
const formRemark = ref('')
const photoFile = ref<File | null>(null)
const photoPreview = ref<string>('')

const displayName = computed(() => {
  const u = auth.user
  if (!u) return '未登录'
  const n = (u.nickname || '').trim()
  return n ? n : u.username
})

const ageText = computed(() => {
  const b = auth.user?.child_birthday || ''
  if (!b) return '未设置'
  const d = dayjs(b)
  if (!d.isValid()) return '未设置'
  const years = dayjs().diff(d, 'year')
  return years >= 0 ? `${years}岁` : '未设置'
})

const genderText = computed(() => {
  const g = auth.user?.child_gender || ''
  if (!g) return '未设置'
  if (g === 'male') return '男'
  if (g === 'female') return '女'
  return '未设置'
})

const photoResolvedMap = ref<Record<string, string>>({})
const storageInfo = ref<{ driver: string; public_base_url?: string } | null>(null)
const avatarSrc = ref<string>(defaultAvatar)

function resolvePhoto(url: string) {
  if (!url) return ''
  if (/^https?:\/\//i.test(url)) return url
  if (photoResolvedMap.value[url]) return photoResolvedMap.value[url]
  presignView(url).then((u) => { photoResolvedMap.value[url] = u }).catch(() => {})
  return ''
}

function formatDate(d: string) {
  return dayjs(d).format('YYYY-MM-DD')
}

function openDialog() {
  formTitle.value = ''
  formIssuer.value = ''
  formDate.value = dayjs().format('YYYY-MM-DD')
  formRemark.value = ''
  photoFile.value = null
  photoPreview.value = ''
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

async function ensureStorageInfo() {
  if (!storageInfo.value) {
    storageInfo.value = await getStorageInfo()
  }
  return storageInfo.value
}

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

async function load() {
  const resp = await listHonors()
  records.value = resp.items || []
}

function openDetail(id: number) {
  router.push(`/honors/${id}`)
}

async function submit() {
  if (!formTitle.value.trim() || !formDate.value) {
    ElMessage.error('请填写荣誉名称与日期')
    return
  }
  if (!photoFile.value) {
    ElMessage.error('请上传荣誉照片')
    return
  }
  saving.value = true
  try {
    const photoUrl = await uploadPhoto(photoFile.value)
    await createHonor({
      title: formTitle.value.trim(),
      issuer: formIssuer.value.trim(),
      awarded_at: formDate.value,
      remark: formRemark.value.trim(),
      photo_url: photoUrl
    })
    showDialog.value = false
    await load()
    ElMessage.success('添加成功')
  } catch (e: any) {
    ElMessage.error(e?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

async function updateAvatar() {
  const p = auth.user?.avatar_path
  if (!p) { avatarSrc.value = defaultAvatar; return }
  const s = String(p)
  if (/storage\/images\/avatars\/default\.png$/i.test(s) || /(^|\/)default\.png$/i.test(s)) { avatarSrc.value = defaultAvatar; return }
  if (/^https?:\/\//i.test(s)) { avatarSrc.value = s; return }
  const base = getStaticBase()
  if (/uploads\//i.test(s)) { avatarSrc.value = `${base}/api/${s.replace(/^\/+/, '')}`; return }
  try { avatarSrc.value = await presignView(s) } catch { avatarSrc.value = defaultAvatar }
}

onMounted(() => {
  updateAvatar()
  load()
})
watch(() => auth.user?.avatar_path, () => { updateAvatar() })
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
