<template>
  <SettingsShell title="创建心愿" subtitle="自定义图标、金币与单位" :icon="Plus" tone="emerald" container-class="max-w-3xl" :decor="false">
    <SettingsCard>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/50 backdrop-blur px-4 py-4">
          <div class="text-[12px] font-extrabold tracking-wide text-gray-500 dark:text-gray-400">心愿图标</div>
          <div class="mt-3 flex items-center gap-3">
            <div class="relative">
              <div class="absolute -inset-2 rounded-2xl bg-emerald-200/40 dark:bg-emerald-500/10 blur-xl" />
              <img
                v-if="form.icon_preview || form.icon"
                :src="form.icon_preview || iconResolved"
                class="relative w-12 h-12 rounded-2xl ring-1 ring-black/5 dark:ring-white/10"
                @error="onIconError"
              />
              <div v-else class="relative w-12 h-12 rounded-2xl border border-gray-200 dark:border-gray-800 bg-white/60 dark:bg-gray-950/30" />
            </div>
            <el-upload :auto-upload="false" :show-file-list="false" accept="image/*" @change="onPickIcon">
              <el-button type="primary" class="!rounded-2xl !font-extrabold">选择图片</el-button>
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
        <el-button type="primary" class="!rounded-2xl !font-extrabold" @click="submitForm">确定</el-button>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
import { reactive, computed, ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import router from '@/router'
import { createWish, uploadWishIcon, normalizeUploadPath } from '@/services/wishes'
import { prepareUpload } from '@/utils/image'
import { presignView } from '@/services/storage'
import { getStaticBase } from '@/services/http'
import { useAuth } from '@/stores/auth'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

const auth = useAuth()
const userId = computed(() => auth.user?.id ?? 0)
function goBack() { router.back() }

type WishForm = { user_id: number; name: string; content: string; icon?: string; icon_preview?: string; need_coins: number; exchange_amount: number; unit: string }
const form = reactive<WishForm>({ user_id: userId.value, name: '', content: '', icon: '', icon_preview: '', need_coins: 1, exchange_amount: 1, unit: '次' })
const iconResolved = ref('')
async function updateIconResolved() {
  const icon = form.icon
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
onMounted(updateIconResolved)
watch(() => form.icon, async () => { await updateIconResolved() })

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
      await updateIconResolved()
    } catch (e) {
      // 中文注释：上传失败时仅保留本地预览，不回填文件名到 icon，避免列表错误显示为内置图标
      try { form.icon_preview = URL.createObjectURL(raw) } catch {}
    }
  }

async function submitForm() {
  try {
    const uid = userId.value
    if (!uid) { ElMessage.error('未登录或令牌无效'); return }
    form.user_id = uid
    await createWish({ user_id: uid, name: form.name, content: form.content, icon: form.icon || '', need_coins: form.need_coins, exchange_amount: form.exchange_amount, unit: form.unit })
    ElMessage.success('创建成功')
    router.back()
  } catch (e: any) {
    ElMessage.error(e.message || '创建失败')
  }
}

// 解析已改为 iconResolved（异步字符串），移除旧的 resolveIcon

// 中文注释：上传后预览失败时回退占位图，避免破图
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
