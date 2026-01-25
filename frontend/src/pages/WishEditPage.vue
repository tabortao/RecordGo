<template>
  <SettingsShell title="编辑心愿" subtitle="更新图标、金币与单位" :icon="Edit" tone="emerald" container-class="max-w-3xl">
    <SettingsCard>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/50 backdrop-blur px-4 py-4">
          <div class="text-[12px] font-extrabold tracking-wide text-gray-500 dark:text-gray-400">心愿图标</div>
          <div class="mt-3 flex items-center gap-3">
            <div class="relative">
              <div class="absolute -inset-2 rounded-2xl bg-emerald-200/40 dark:bg-emerald-500/10 blur-xl" />
              <div v-if="isEmojiIcon" class="relative w-12 h-12 rounded-2xl ring-1 ring-black/5 dark:ring-white/10 bg-white/60 dark:bg-gray-950/25 grid place-items-center text-[24px] leading-none">
                {{ emojiChar }}
              </div>
              <img
                v-else-if="form.icon_preview || form.icon"
                :src="form.icon_preview || iconResolved"
                class="relative w-12 h-12 rounded-2xl ring-1 ring-black/5 dark:ring-white/10"
                @error="onIconError"
              />
              <div v-else class="relative w-12 h-12 rounded-2xl border border-gray-200 dark:border-gray-800 bg-white/60 dark:bg-gray-950/30" />
            </div>
            <el-upload :auto-upload="false" :show-file-list="false" accept="image/*" @change="onPickIcon">
              <el-button type="primary" class="!rounded-2xl !font-extrabold">更换图片</el-button>
            </el-upload>
          </div>
        </div>

        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/50 backdrop-blur px-4 py-4">
          <el-form :model="form" label-position="top" class="wish-form">
            <el-form-item label="心愿名称">
              <el-input v-model="form.name" placeholder="例如：看电影" />
            </el-form-item>
            <el-form-item label="心愿描述">
              <el-input type="textarea" v-model="form.content" :rows="3" placeholder="写一句更具体的目标" />
            </el-form-item>
          </el-form>
        </div>
      </div>

      <div class="mt-4 grid grid-cols-1 md:grid-cols-3 gap-4">
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/50 backdrop-blur px-4 py-4">
          <div class="text-[12px] font-extrabold tracking-wide text-gray-500 dark:text-gray-400">所需金币</div>
          <div class="mt-3">
            <el-input-number v-model="form.need_coins" :min="1" controls-position="right" class="w-full" />
          </div>
        </div>
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/50 backdrop-blur px-4 py-4">
          <div class="text-[12px] font-extrabold tracking-wide text-gray-500 dark:text-gray-400">单位</div>
          <div class="mt-3">
            <el-select v-model="form.unit" class="w-full">
              <el-option label="个" value="个" /><el-option label="次" value="次" />
              <el-option label="分钟" value="分钟" /><el-option label="元" value="元" />
            </el-select>
          </div>
        </div>
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/50 backdrop-blur px-4 py-4">
          <div class="text-[12px] font-extrabold tracking-wide text-gray-500 dark:text-gray-400">兑换数量</div>
          <div class="mt-3">
            <el-input-number v-model="form.exchange_amount" :min="1" controls-position="right" class="w-full" />
          </div>
        </div>
      </div>

      <div class="mt-5 flex justify-end gap-2">
        <el-button class="!rounded-2xl !font-extrabold" @click="goBack">取消</el-button>
        <el-button type="primary" class="!rounded-2xl !font-extrabold" @click="submitForm">保存</el-button>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { reactive, onMounted, ref, watch, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Edit } from '@element-plus/icons-vue'
import { useRoute } from 'vue-router'
import router from '@/router'
import { getWish, updateWish, uploadWishIcon, normalizeUploadPath, type Wish } from '@/services/wishes'
import { useAuth } from '@/stores/auth'
import { prepareUpload } from '@/utils/image'
import { getStaticBase } from '@/services/http'
import { presignView } from '@/services/storage'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

const route = useRoute()
const auth = useAuth()
const userId = computed(() => auth.user?.id ?? 0)
function goBack() { router.back() }

type WishForm = Partial<Wish> & { icon_preview?: string }
const form = reactive<WishForm>({ name: '', content: '', need_coins: 1, unit: '次', exchange_amount: 1, icon: '', icon_preview: '' })
const iconResolved = ref('')

const isEmojiIcon = computed(() => typeof form.icon === 'string' && form.icon.startsWith('emoji:'))
const emojiChar = computed(() => (isEmojiIcon.value ? String(form.icon).slice('emoji:'.length) : ''))

onMounted(async () => {
  const id = Number(route.params.id)
  const w = await getWish(id)
  Object.assign(form, w)
  await updateIconResolved()
})

async function onPickIcon(fileEvent: any) {
  const raw: File | undefined = fileEvent?.raw || fileEvent?.target?.files?.[0] || fileEvent?.file
  if (!raw) return
  const webp = await prepareUpload(raw, 0.8)
  try { form.icon_preview = URL.createObjectURL(webp) } catch {}
  try {
    const { path } = await uploadWishIcon(userId.value, webp)
    form.icon = normalizeUploadPath(path)
    try { form.icon_preview && URL.revokeObjectURL(form.icon_preview as any) } catch {}
    form.icon_preview = ''
  } catch (e) {
    // 失败则保留预览，不更新服务器路径
    try { form.icon_preview = URL.createObjectURL(raw) } catch {}
  }
}

async function submitForm() {
  try {
    const id = Number(route.params.id)
    const payload: Partial<Wish> = {
      name: form.name!,
      content: form.content!,
      need_coins: Number(form.need_coins || 1),
      unit: form.unit!,
      exchange_amount: Number(form.exchange_amount || 1),
      icon: form.icon || ''
    }
    await updateWish(id, payload)
    ElMessage.success('保存成功')
    router.back()
  } catch (e: any) {
    ElMessage.error(e.message || '保存失败')
  }
}

// 中文注释：解析并设置图标展示地址（异步），模板中仅绑定字符串
async function updateIconResolved() {
  const icon = form.icon
  if (typeof icon === 'string' && icon.startsWith('emoji:')) { iconResolved.value = new URL('../assets/wishs/领取记录.png', import.meta.url).href; return }
  if (!icon) { iconResolved.value = new URL('../assets/wishs/领取记录.png', import.meta.url).href; return }
  if (/\.(png|jpg|jpeg|webp)$/i.test(icon) && !icon.includes('/')) {
    iconResolved.value = new URL(`../assets/wishs/${icon}`, import.meta.url).href
    return
  }
  const base = getStaticBase()
  const path = normalizeUploadPath(icon)
  if (path.startsWith('uploads/')) { iconResolved.value = `${base}/api/${path}`; return }
  try { iconResolved.value = await presignView(path) } catch { iconResolved.value = new URL('../assets/wishs/领取记录.png', import.meta.url).href }
}

watch(() => form.icon, async () => { await updateIconResolved() })

// 中文注释：心愿图标加载失败时采用占位图，提升容错体验
function onIconError(e: Event) {
  const img = e.target as HTMLImageElement
  try { img.src = new URL(`../assets/wishs/领取记录.png`, import.meta.url).href } catch {}
}
</script>

<style scoped>
:deep(.wish-form .el-form-item__label) {
  font-size: 12px;
  font-weight: 800;
  color: rgb(107 114 128);
}

.dark :deep(.wish-form .el-form-item__label) {
  color: rgb(156 163 175);
}

:deep(.wish-form .el-input__inner),
:deep(.wish-form .el-textarea__inner),
:deep(.wish-form .el-input-number__decrease),
:deep(.wish-form .el-input-number__increase) {
  font-size: 14px;
}
</style>
