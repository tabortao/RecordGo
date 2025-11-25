<template>
  <!-- 中文注释：独立创建心愿页面，顶部返回图标；响应式布局 -->
  <div class="p-4 space-y-4">
    <!-- 顶部栏：返回 + 标题 -->
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer text-gray-600" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" class="text-emerald-600"><Plus /></el-icon>
      <h2 class="font-semibold">创建心愿</h2>
    </div>

    <!-- 响应式网格：移动端单列，桌面端双列/多列 -->
    <div class="grid grid-cols-1 gap-4">
      <el-card shadow="never" class="rounded-lg">
        <el-form :model="form" label-width="90px">
          <el-form-item label="心愿图标">
            <div class="flex items-center gap-2">
              <img v-if="form.icon_preview || form.icon" :src="form.icon_preview || iconResolved" class="w-10 h-10 rounded" @error="onIconError" />
              <el-upload :auto-upload="false" :show-file-list="false" accept="image/*" @change="onPickIcon">
                <el-button type="primary" size="small">选择图片</el-button>
              </el-upload>
            </div>
          </el-form-item>
          <el-form-item label="心愿名称"><el-input v-model="form.name" /></el-form-item>
          <el-form-item label="心愿描述"><el-input type="textarea" v-model="form.content" /></el-form-item>
          <el-form-item label="所需金币"><el-input-number v-model="form.need_coins" :min="1" /></el-form-item>
          <el-form-item label="单位">
            <el-select v-model="form.unit" style="width: 100%">
              <el-option label="个" value="个" /><el-option label="次" value="次" />
              <el-option label="分钟" value="分钟" /><el-option label="元" value="元" />
            </el-select>
          </el-form-item>
          <el-form-item label="兑换数量"><el-input-number v-model="form.exchange_amount" :min="1" /></el-form-item>
        </el-form>
      </el-card>
    </div>

    <!-- 底部操作区 -->
    <div class="flex justify-end gap-2">
      <el-button @click="goBack">取消</el-button>
      <el-button type="primary" @click="submitForm">确定</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, computed, ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Plus } from '@element-plus/icons-vue'
import router from '@/router'
import { createWish, uploadWishIcon, normalizeUploadPath } from '@/services/wishes'
import { prepareUpload } from '@/utils/image'
import { presignView } from '@/services/storage'
import { getStaticBase } from '@/services/http'
import { useAuth } from '@/stores/auth'

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
    form.icon = path
    try { form.icon_preview && URL.revokeObjectURL(form.icon_preview as any) } catch {}
    form.icon_preview = ''
    await updateIconResolved()
  } catch (e) {
    try { form.icon_preview = URL.createObjectURL(raw) } catch {}
    form.icon = raw.name
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
/* 中文注释：使用 Tailwind 进行响应式布局，不额外定义样式 */
</style>
