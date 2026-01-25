<template>
  <SettingsShell title="创建心愿" subtitle="自定义图标、金币与单位" :icon="Plus" tone="emerald" container-class="max-w-3xl" :decor="false">
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <div class="lg:col-span-1 space-y-4">
        <SettingsCard title="预览" description="创建后将展示在心愿列表中">
          <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-950/20 px-4 py-4">
            <div class="flex items-center gap-3">
              <div class="relative shrink-0">
                <div class="absolute -inset-2 rounded-2xl bg-emerald-200/35 dark:bg-emerald-500/12 blur-xl" />
                <div v-if="isEmojiIcon" class="relative w-14 h-14 rounded-3xl ring-1 ring-black/5 dark:ring-white/10 bg-white/70 dark:bg-gray-950/25 grid place-items-center text-[28px] leading-none">
                  {{ emojiChar }}
                </div>
                <img
                  v-else
                  :src="form.icon_preview || iconResolved"
                  class="relative w-14 h-14 rounded-3xl ring-1 ring-black/5 dark:ring-white/10"
                  @error="onIconError"
                />
              </div>
              <div class="min-w-0 flex-1">
                <div class="text-base font-extrabold tracking-tight text-gray-900 dark:text-gray-50 truncate">{{ previewTitle }}</div>
                <div class="mt-1 text-xs text-gray-500 dark:text-gray-400 truncate">{{ previewDesc }}</div>
              </div>
            </div>
            <div class="mt-4 grid grid-cols-3 gap-2">
              <div class="rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-950/20 px-3 py-2">
                <div class="text-[11px] font-extrabold tracking-[0.22em] text-gray-500 dark:text-gray-400">金币</div>
                <div class="mt-1 text-sm font-black text-emerald-700 dark:text-emerald-300 tabular-nums">{{ form.need_coins }}</div>
              </div>
              <div class="rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-950/20 px-3 py-2">
                <div class="text-[11px] font-extrabold tracking-[0.22em] text-gray-500 dark:text-gray-400">单位</div>
                <div class="mt-1 text-sm font-black text-gray-900 dark:text-gray-50">{{ form.unit }}</div>
              </div>
              <div class="rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-950/20 px-3 py-2">
                <div class="text-[11px] font-extrabold tracking-[0.22em] text-gray-500 dark:text-gray-400">数量</div>
                <div class="mt-1 text-sm font-black text-gray-900 dark:text-gray-50 tabular-nums">{{ form.exchange_amount }}</div>
              </div>
            </div>
          </div>
        </SettingsCard>

        <SettingsCard title="心愿图标" description="可上传图片或使用 Emoji 图标">
          <div class="flex items-center justify-between gap-3">
            <el-upload :auto-upload="false" :show-file-list="false" accept="image/*" @change="onPickIcon">
              <el-button type="primary" class="!rounded-2xl !font-extrabold">上传图标</el-button>
            </el-upload>
            <el-button plain class="!rounded-2xl !font-extrabold" @click="resetIcon">重置</el-button>
          </div>

          <div class="mt-5 rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/65 dark:bg-gray-950/15 px-4 py-4">
            <div class="flex items-center justify-between gap-3">
              <div>
                <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">Emoji 图标</div>
                <div class="mt-1 text-xs text-gray-500 dark:text-gray-400">点击即可使用，适合常见心愿</div>
              </div>
              <div v-if="isEmojiIcon" class="rounded-full border border-emerald-200/70 dark:border-emerald-900/40 bg-emerald-50/70 dark:bg-emerald-900/15 px-3 py-1 text-xs font-extrabold text-emerald-700 dark:text-emerald-200">
                已选择
              </div>
            </div>

            <div class="mt-4 flex flex-wrap gap-2">
              <button
                v-for="c in emojiCategories"
                :key="c.key"
                type="button"
                class="rounded-full border px-3 py-1 text-xs font-extrabold transition active:scale-[0.99]"
                :class="activeEmojiCategory === c.key
                  ? 'border-emerald-200/80 dark:border-emerald-900/45 bg-emerald-50/80 dark:bg-emerald-900/18 text-emerald-800 dark:text-emerald-200'
                  : 'border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-950/10 text-gray-600 dark:text-gray-300 hover:bg-white dark:hover:bg-gray-900/40'"
                @click="activeEmojiCategory = c.key"
              >
                {{ c.label }}
              </button>
            </div>

            <div class="mt-4 grid grid-cols-6 sm:grid-cols-8 gap-2">
              <button
                v-for="e in activeEmojiList"
                :key="e"
                type="button"
                class="h-10 w-10 rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/75 dark:bg-gray-950/15 hover:bg-white dark:hover:bg-gray-900/40 transition active:scale-[0.99] grid place-items-center text-[20px] leading-none"
                :class="form.icon === ('emoji:' + e) ? 'ring-2 ring-emerald-400/80 dark:ring-emerald-300/60' : ''"
                @click="pickEmojiIcon(e)"
              >
                {{ e }}
              </button>
            </div>
          </div>
        </SettingsCard>
      </div>

      <div class="lg:col-span-2 space-y-4">
        <SettingsCard title="基本信息" description="名称越具体，越容易坚持">
          <el-form :model="form" label-position="top" class="wish-form">
            <el-form-item label="心愿名称">
              <el-input v-model="form.name" placeholder="例如：周末看电影" size="large" />
            </el-form-item>
            <el-form-item label="心愿描述">
              <el-input type="textarea" v-model="form.content" :rows="4" placeholder="写一句更具体的目标，例如：本周末和家人一起看一场电影" />
            </el-form-item>
          </el-form>
        </SettingsCard>

        <SettingsCard title="兑换规则" description="设置所需金币、单位与可兑换数量">
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
            <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-950/20 px-4 py-4">
              <div class="text-[12px] font-extrabold tracking-wide text-gray-500 dark:text-gray-400">所需金币</div>
              <div class="mt-3">
                <el-input-number v-model="form.need_coins" :min="1" controls-position="right" class="w-full" size="large" />
              </div>
            </div>
            <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-950/20 px-4 py-4">
              <div class="text-[12px] font-extrabold tracking-wide text-gray-500 dark:text-gray-400">单位</div>
              <div class="mt-3">
                <el-select v-model="form.unit" class="w-full" size="large">
                  <el-option label="个" value="个" /><el-option label="次" value="次" />
                  <el-option label="分钟" value="分钟" /><el-option label="元" value="元" />
                </el-select>
              </div>
            </div>
            <div class="rounded-3xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-950/20 px-4 py-4">
              <div class="text-[12px] font-extrabold tracking-wide text-gray-500 dark:text-gray-400">兑换数量</div>
              <div class="mt-3">
                <el-input-number v-model="form.exchange_amount" :min="1" controls-position="right" class="w-full" size="large" />
              </div>
            </div>
          </div>
        </SettingsCard>

        <div class="rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/80 dark:bg-gray-900/65 backdrop-blur-xl shadow-sm px-4 py-4 flex justify-end gap-2">
          <el-button class="!rounded-2xl !font-extrabold" @click="goBack">取消</el-button>
          <el-button type="primary" class="!rounded-2xl !font-extrabold" @click="submitForm">确定</el-button>
        </div>
      </div>
    </div>
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

const previewTitle = computed(() => (form.name || '').trim() || '未命名心愿')
const previewDesc = computed(() => (form.content || '').trim() || '写一句更具体的目标，会更有动力')

const emojiCategories = [
  {
    key: 'game',
    label: '玩家游戏',
    emojis: ['🎮', '🕹️', '🎲', '♟️', '🧩', '🃏', '🏆', '🏅', '⚔️', '🛡️', '🏎️', '🏁']
  },
  {
    key: 'food',
    label: '美食零食',
    emojis: ['🍔', '🍟', '🍕', '🌭', '🍿', '🍣', '🍜', '🍰', '🍩', '🍪', '🧋', '🍫']
  },
  {
    key: 'life',
    label: '生活娱乐',
    emojis: ['🎬', '🎤', '🎧', '🎵', '🎸', '🥳', '🎡', '🎢', '🎠', '🏖️', '🏕️', '🧘']
  },
  {
    key: 'digital',
    label: '电子产品',
    emojis: ['📱', '💻', '🖥️', '⌨️', '🖱️', '🎧', '🎮', '📷', '🎥', '📺', '⌚', '🔋']
  },
  {
    key: 'study',
    label: '书籍学习',
    emojis: ['📚', '📖', '📝', '✏️', '🖊️', '📓', '📒', '📔', '🧠', '🔬', '🧪', '🎓']
  },
  {
    key: 'reward',
    label: '特权奖励',
    emojis: ['👑', '✨', '🎁', '🎟️', '🎫', '💎', '🍦', '🛍️', '🧸', '🎉', '🌟', '🏖️']
  },
  {
    key: 'other',
    label: '其他',
    emojis: ['🚀', '🌈', '🌙', '☁️', '🎈', '🧿', '🌻', '🐶', '🐱', '🎧', '🧃', '🎀']
  }
]
const activeEmojiCategory = ref<string>(emojiCategories[0].key)
const activeEmojiList = computed(() => emojiCategories.find(c => c.key === activeEmojiCategory.value)?.emojis || [])

const isEmojiIcon = computed(() => typeof form.icon === 'string' && form.icon.startsWith('emoji:'))
const emojiChar = computed(() => (isEmojiIcon.value ? String(form.icon).slice('emoji:'.length) : ''))

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

function pickEmojiIcon(emoji: string) {
  try { form.icon_preview && URL.revokeObjectURL(form.icon_preview as any) } catch {}
  form.icon_preview = ''
  form.icon = `emoji:${emoji}`
}

function resetIcon() {
  try { form.icon_preview && URL.revokeObjectURL(form.icon_preview as any) } catch {}
  form.icon_preview = ''
  form.icon = ''
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

:deep(.wish-form .el-input__wrapper),
:deep(.wish-form .el-textarea__inner),
:deep(.wish-form .el-select__wrapper) {
  border-radius: 16px;
  border: 1px solid rgb(203 213 225);
  background: rgb(255 255 255 / 0.92);
}

.dark :deep(.wish-form .el-input__wrapper),
.dark :deep(.wish-form .el-textarea__inner),
.dark :deep(.wish-form .el-select__wrapper) {
  border: 1px solid rgb(51 65 85);
  background: rgb(2 6 23 / 0.18);
}

:deep(.wish-form .el-input__wrapper),
:deep(.wish-form .el-select__wrapper),
:deep(.wish-form .el-textarea__inner) {
  box-shadow: none;
}

:deep(.wish-form .el-input__wrapper.is-focus),
:deep(.wish-form .el-select__wrapper.is-focused),
:deep(.wish-form .el-textarea__inner:focus) {
  border-color: rgb(16 185 129);
  box-shadow: 0 0 0 4px rgb(16 185 129 / 0.16);
}

.dark :deep(.wish-form .el-input__wrapper.is-focus),
.dark :deep(.wish-form .el-select__wrapper.is-focused),
.dark :deep(.wish-form .el-textarea__inner:focus) {
  border-color: rgb(16 185 129);
  box-shadow: 0 0 0 4px rgb(16 185 129 / 0.14);
}
</style>
