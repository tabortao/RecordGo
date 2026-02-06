<template>
  <SettingsShell title="创建心愿" subtitle="自定义图标、金币与单位" :icon="Plus" tone="emerald" container-class="max-w-6xl" :decor="false">
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 lg:gap-10 items-start">
      <div class="lg:col-span-5 xl:col-span-4 space-y-6 lg:sticky lg:top-24 transition-all duration-300">
        <SettingsCard title="预览" description="创建后将展示在心愿列表中">
          <div class="rounded-[28px] border border-gray-100 dark:border-gray-800 bg-gradient-to-br from-white/80 to-white/40 dark:from-gray-900/40 dark:to-gray-900/20 p-5 shadow-sm">
            <div class="flex items-center gap-4">
              <div class="relative shrink-0">
                <div class="absolute -inset-4 rounded-full bg-emerald-400/20 dark:bg-emerald-500/10 blur-2xl" />
                <div v-if="isEmojiIcon" class="relative w-16 h-16 rounded-[20px] shadow-sm ring-1 ring-black/5 dark:ring-white/10 bg-white dark:bg-gray-800 grid place-items-center text-[32px] leading-none transition-transform duration-300 hover:scale-105">
                  {{ emojiChar }}
                </div>
                <img
                  v-else
                  :src="form.icon_preview || iconResolved"
                  class="relative w-16 h-16 rounded-[20px] shadow-sm ring-1 ring-black/5 dark:ring-white/10 object-cover bg-white dark:bg-gray-800 transition-transform duration-300 hover:scale-105"
                  @error="onIconError"
                />
              </div>
              <div class="min-w-0 flex-1 space-y-1">
                <div class="text-lg font-black tracking-tight text-gray-900 dark:text-white truncate">{{ previewTitle }}</div>
                <div class="text-xs font-medium text-gray-500 dark:text-gray-400 truncate">{{ previewDesc }}</div>
              </div>
            </div>
            <div class="mt-6 grid grid-cols-3 gap-3">
              <div class="group rounded-2xl bg-gray-50/80 dark:bg-gray-800/50 p-3 text-center transition-colors hover:bg-emerald-50/50 dark:hover:bg-emerald-900/10">
                <div class="text-[10px] font-bold uppercase tracking-wider text-gray-400 group-hover:text-emerald-600/70 dark:group-hover:text-emerald-400/70">金币</div>
                <div class="mt-1 text-base font-black text-emerald-600 dark:text-emerald-400 tabular-nums">{{ form.need_coins }}</div>
              </div>
              <div class="group rounded-2xl bg-gray-50/80 dark:bg-gray-800/50 p-3 text-center transition-colors hover:bg-blue-50/50 dark:hover:bg-blue-900/10">
                <div class="text-[10px] font-bold uppercase tracking-wider text-gray-400 group-hover:text-blue-600/70 dark:group-hover:text-blue-400/70">单位</div>
                <div class="mt-1 text-base font-black text-gray-900 dark:text-gray-100">{{ form.unit }}</div>
              </div>
              <div class="group rounded-2xl bg-gray-50/80 dark:bg-gray-800/50 p-3 text-center transition-colors hover:bg-purple-50/50 dark:hover:bg-purple-900/10">
                <div class="text-[10px] font-bold uppercase tracking-wider text-gray-400 group-hover:text-purple-600/70 dark:group-hover:text-purple-400/70">数量</div>
                <div class="mt-1 text-base font-black text-gray-900 dark:text-gray-100 tabular-nums">{{ form.exchange_amount }}</div>
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

          <div class="mt-6 rounded-[24px] border border-gray-100 dark:border-gray-800 bg-white/50 dark:bg-gray-950/20 p-5">
            <div class="flex items-center justify-between gap-3 mb-5">
              <div>
                <div class="text-sm font-black text-gray-900 dark:text-gray-50">选择 Emoji</div>
                <div class="mt-1 text-xs font-medium text-gray-400 dark:text-gray-500">点击即可使用，快速美化心愿</div>
              </div>
              <div v-if="isEmojiIcon" class="rounded-full border border-emerald-500/20 bg-emerald-500/10 px-3 py-1 text-[10px] font-bold text-emerald-600 dark:text-emerald-400">
                已选中
              </div>
            </div>

            <!-- 分类标签 -->
            <div class="flex flex-wrap gap-2 pb-2 -mx-1 px-1">
              <button
                v-for="c in emojiCategories"
                :key="c.key"
                type="button"
                class="shrink-0 rounded-full border px-4 py-1.5 text-xs font-bold transition-all duration-200 active:scale-95"
                :class="activeEmojiCategory === c.key
                  ? 'border-emerald-500/30 bg-emerald-500/10 text-emerald-700 dark:text-emerald-300 shadow-sm shadow-emerald-500/10'
                  : 'border-transparent bg-gray-100/50 dark:bg-gray-800/50 text-gray-500 dark:text-gray-400 hover:bg-white dark:hover:bg-gray-800'"
                @click="activeEmojiCategory = c.key"
              >
                {{ c.label }}
              </button>
            </div>

            <!-- Emoji 网格 -->
            <div class="mt-4 max-h-[320px] overflow-y-auto pr-2 grid grid-cols-6 sm:grid-cols-8 md:grid-cols-6 lg:grid-cols-5 xl:grid-cols-6 gap-3 custom-scrollbar">
              <button
                v-for="e in activeEmojiList"
                :key="e"
                type="button"
                class="aspect-square rounded-2xl border border-transparent hover:border-gray-200 dark:hover:border-gray-700 bg-white/60 dark:bg-gray-800/40 hover:bg-white dark:hover:bg-gray-800 hover:shadow-sm hover:-translate-y-0.5 transition-all duration-200 grid place-items-center text-2xl leading-none"
                :class="form.icon === ('emoji:' + e) ? '!border-emerald-500/50 !bg-emerald-500/10 ring-2 ring-emerald-500/20' : ''"
                @click="pickEmojiIcon(e)"
              >
                {{ e }}
              </button>
            </div>
          </div>
        </SettingsCard>
      </div>

      <div class="lg:col-span-7 xl:col-span-8 space-y-6">
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
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 lg:gap-6">
            <div class="rounded-[24px] border border-gray-100 dark:border-gray-800 bg-white/50 dark:bg-gray-950/20 p-5 transition-colors hover:bg-white dark:hover:bg-gray-900/40 flex flex-col">
              <div class="text-xs font-bold uppercase tracking-wider text-gray-400 dark:text-gray-500 mb-3">所需金币</div>
              <el-input-number v-model="form.need_coins" :min="1" controls-position="right" class="w-full !rounded-xl mt-auto" size="large" />
            </div>
            <div class="rounded-[24px] border border-gray-100 dark:border-gray-800 bg-white/50 dark:bg-gray-950/20 p-5 transition-colors hover:bg-white dark:hover:bg-gray-900/40 flex flex-col">
              <div class="text-xs font-bold uppercase tracking-wider text-gray-400 dark:text-gray-500 mb-3">单位</div>
              <el-select v-model="form.unit" class="w-full mt-auto" size="large">
                <el-option label="个" value="个" /><el-option label="次" value="次" />
                <el-option label="分钟" value="分钟" /><el-option label="元" value="元" />
              </el-select>
            </div>
            <div class="rounded-[24px] border border-gray-100 dark:border-gray-800 bg-white/50 dark:bg-gray-950/20 p-5 transition-colors hover:bg-white dark:hover:bg-gray-900/40 flex flex-col">
              <div class="text-xs font-bold uppercase tracking-wider text-gray-400 dark:text-gray-500 mb-3">兑换数量</div>
              <el-input-number v-model="form.exchange_amount" :min="1" controls-position="right" class="w-full mt-auto" size="large" />
            </div>
          </div>
        </SettingsCard>

        <div class="sticky bottom-6 rounded-[24px] border border-white/60 dark:border-gray-800/60 bg-white/80 dark:bg-gray-900/80 backdrop-blur-xl shadow-lg shadow-gray-200/20 dark:shadow-black/20 p-4 flex justify-end gap-3 z-10">
          <el-button class="!rounded-xl !h-10 !px-6 !font-bold !text-gray-600 dark:!text-gray-300" @click="goBack">取消</el-button>
          <el-button type="primary" class="!rounded-xl !h-10 !px-8 !font-bold shadow-md shadow-emerald-500/20" @click="submitForm">确定创建</el-button>
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
    label: '游戏',
    emojis: ['🎮', '🕹️', '🎲', '♟️', '🧩', '🃏', '🏆', '🏅', '⚔️', '🛡️', '🏎️', '🏁']
  },
  {
    key: 'food',
    label: '美食',
    emojis: ['🍔', '🍟', '🍕', '🌭', '🍿', '🍣', '🍜', '🍰', '🍩', '🍪', '🧋', '🍫']
  },
  {
    key: 'life',
    label: '娱乐',
    emojis: ['🎬', '🎤', '🎧', '🎵', '🎸', '🥳', '🎡', '🎢', '🎠', '🏖️', '🏕️', '🧘']
  },
  {
    key: 'digital',
    label: '数码',
    emojis: ['📱', '💻', '🖥️', '⌨️', '🖱️', '🎧', '🎮', '📷', '🎥', '📺', '⌚', '🔋']
  },
  {
    key: 'study',
    label: '学习',
    emojis: ['📚', '📖', '📝', '✏️', '🖊️', '📓', '📒', '📔', '🧠', '🔬', '🧪', '🎓']
  },
  {
    key: 'reward',
    label: '奖励',
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
  box-shadow: none;
}

.dark :deep(.wish-form .el-input__inner),
.dark :deep(.wish-form .el-textarea__inner) {
  color: rgb(249 250 251);
}

.dark :deep(.wish-form .el-select__wrapper .el-select__selected-item) {
  color: rgb(249 250 251);
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

.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: rgba(156, 163, 175, 0.3);
  border-radius: 4px;
}
.dark .custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: rgba(156, 163, 175, 0.2);
}
</style>
