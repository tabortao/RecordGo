<template>
  <!-- 中文注释：独立编辑心愿页面，顶部返回图标；加载现有数据并支持修改图标 -->
  <div class="p-4 space-y-4">
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer text-gray-600" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" class="text-emerald-600"><Edit /></el-icon>
      <h2 class="font-semibold">编辑心愿</h2>
    </div>

    <div class="grid grid-cols-1 gap-4">
      <el-card shadow="never" class="rounded-lg">
        <el-form :model="form" label-width="90px">
          <el-form-item label="心愿图标">
            <div class="flex items-center gap-2">
              <img v-if="form.icon_preview || form.icon" :src="form.icon_preview || resolveIcon(form.icon)" class="w-10 h-10 rounded" @error="onIconError" />
              <el-upload :auto-upload="false" :show-file-list="false" accept="image/*" @change="onPickIcon">
                <el-button type="primary" size="small">更换图片</el-button>
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

    <div class="flex justify-end gap-2">
      <el-button @click="goBack">取消</el-button>
      <el-button type="primary" @click="submitForm">保存</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
// 中文注释：编辑心愿页面逻辑，加载现有心愿并提交更新；图标上传仍转换为 webp 并走后端 uploads
import { reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Edit } from '@element-plus/icons-vue'
import { useRoute } from 'vue-router'
import router from '@/router'
import { getWish, updateWish, uploadWishIcon, toWebp, normalizeUploadPath, type Wish } from '@/services/wishes'
import { getStaticBase } from '@/services/http'

const route = useRoute()
const userId = 1 // 中文注释：示例用户ID
function goBack() { router.back() }

type WishForm = Partial<Wish> & { icon_preview?: string }
const form = reactive<WishForm>({ name: '', content: '', need_coins: 1, unit: '次', exchange_amount: 1, icon: '', icon_preview: '' })

onMounted(async () => {
  const id = Number(route.params.id)
  const w = await getWish(id)
  Object.assign(form, w)
})

async function onPickIcon(fileEvent: any) {
  const raw: File | undefined = fileEvent?.raw || fileEvent?.target?.files?.[0] || fileEvent?.file
  if (!raw) return
  try { form.icon_preview = URL.createObjectURL(raw) } catch {}
  const webp = await toWebp(raw)
  try {
    const { path } = await uploadWishIcon(userId, webp)
    form.icon = normalizeUploadPath(path)
    try { form.icon_preview && URL.revokeObjectURL(form.icon_preview as any) } catch {}
    form.icon_preview = ''
  } catch (e) {
    // 失败则保留预览，不更新服务器路径
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

// 中文注释：解析图标路径，兼容内置与后端 uploads
function resolveIcon(icon?: string) {
  if (!icon) return new URL('../assets/wishs/领取记录.png', import.meta.url).href
  if (/\.(png|jpg|jpeg|webp)$/i.test(icon) && !icon.includes('/')) {
    return new URL(`../assets/wishs/${icon}`, import.meta.url).href
  }
  // 规范化路径并拼接基址
  const base = getStaticBase()
  const path = normalizeUploadPath(icon)
  // 中文注释：后端静态文件映射在 /api/uploads，这里需要拼接 /api 前缀
  return `${base}/api/${path}`
}

// 中文注释：心愿图标加载失败时采用占位图，提升容错体验
function onIconError(e: Event) {
  const img = e.target as HTMLImageElement
  try { img.src = new URL(`../assets/wishs/领取记录.png`, import.meta.url).href } catch {}
}
</script>

<style scoped>
/* 中文注释：使用 Tailwind 响应式布局，无额外样式 */
</style>