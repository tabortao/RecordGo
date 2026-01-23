<template>
  <div class="min-h-screen bg-[#F5F7FA] dark:bg-gray-900">
    <div class="bg-white/80 dark:bg-gray-900/80 backdrop-blur border-b border-gray-200 dark:border-gray-800 sticky top-0 z-10">
      <div class="px-4 py-3 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <el-icon :size="18" class="cursor-pointer text-gray-600 dark:text-gray-300" @click="router.back()"><ArrowLeft /></el-icon>
          <h2 class="font-semibold text-gray-800 dark:text-gray-100">荣誉</h2>
        </div>
        <el-button type="primary" size="small" @click="openDialog">添加荣誉</el-button>
      </div>
    </div>

    <div class="p-4 space-y-4">
      <div class="rounded-2xl bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 p-4 shadow-sm">
        <div class="flex items-center gap-3">
          <el-avatar :size="44" :src="avatarSrc" />
          <div class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ displayName }}</div>
        </div>
        <div class="mt-2 flex items-center gap-4 text-sm text-gray-600 dark:text-gray-300">
          <div class="flex items-center gap-1">
            <span class="text-gray-500 dark:text-gray-400">年龄</span>
            <span class="text-gray-800 dark:text-gray-100">{{ ageText }}</span>
          </div>
          <div class="flex items-center gap-1">
            <span class="text-gray-500 dark:text-gray-400">性别</span>
            <span class="text-gray-800 dark:text-gray-100">{{ genderText }}</span>
          </div>
        </div>
      </div>

      <div v-if="records.length === 0" class="text-sm text-gray-500 dark:text-gray-400">暂无荣誉记录</div>
      <div v-else class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
        <div
          v-for="item in records"
          :key="item.id"
          class="group relative overflow-hidden rounded-2xl border border-gray-100 bg-white shadow-sm transition hover:shadow-md dark:border-gray-700 dark:bg-gray-800 cursor-pointer"
          @click="openDetail(item.id)"
        >
          <div class="flex h-52 items-center justify-center bg-gradient-to-br from-amber-50 via-white to-amber-100 dark:from-gray-800 dark:via-gray-800 dark:to-gray-900">
            <img
              v-if="resolvePhoto(item.photo_url)"
              :src="resolvePhoto(item.photo_url)"
              class="max-h-full max-w-full object-contain"
            />
          </div>
          <div class="flex flex-col items-center gap-1 px-3 py-2 text-sm text-gray-700 dark:text-gray-200">
            <div class="w-full truncate text-center font-semibold">{{ item.title }}</div>
            <div class="text-xs text-gray-500 dark:text-gray-400">{{ formatDate(item.awarded_at) }}</div>
          </div>
        </div>
      </div>
    </div>

    <el-dialog v-model="showDialog" title="添加荣誉" width="420px">
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
import { computed, onMounted, ref, watch } from 'vue'
import { ArrowLeft } from '@element-plus/icons-vue'
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
</style>
