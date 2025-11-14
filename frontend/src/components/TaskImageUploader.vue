<template>
  <!-- 中文注释：任务图片上传组件（创建/编辑通用），支持多选、缩略图预览与上传进度 -->
  <div class="space-y-2">
    <div class="flex items-center gap-2">
      <el-button type="primary" size="small" @click="triggerPick">选择图片</el-button>
      <input ref="fileInput" type="file" accept="image/*" multiple class="hidden" @change="onPicked" />
    </div>

    <!-- 缩略图预览与进度条（下方显示） -->
    <div class="grid grid-cols-3 gap-3">
      <div v-for="(item, idx) in items" :key="item.key" class="border rounded p-2 relative">
        <!-- 中文注释：点击缩略图可放大预览，支持左右翻看；右上角提供删除图标 -->
        <img :src="item.url" alt="task image" class="w-full h-24 object-contain rounded" />
        <el-icon class="absolute top-1 right-1 cursor-pointer text-white bg-black/50 rounded-full p-0.5" :size="16" title="删除" @click="removeAt(idx)"><Close /></el-icon>
        <div v-if="item.uploading" class="absolute left-2 right-2 bottom-2">
          <el-progress :percentage="item.progress" :stroke-width="8" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 中文注释：属性定义——编辑模式下立即上传到后端；创建模式下仅本地暂存；支持是否使用/恢复草稿的开关
import { ref, reactive, watch, onMounted, computed, withDefaults } from 'vue'
import { Close } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { prepareUpload } from '@/utils/image'
import { uploadTaskImage, deleteTaskImage } from '@/services/tasks'
import { normalizeUploadPath } from '@/services/wishes'

type Item = { key: string; name: string; url: string; uploading: boolean; progress: number; serverPath?: string; localFile?: File }

const props = withDefaults(defineProps<{
  editing: boolean
  userId: number
  taskId?: number
  // 中文注释：最多可选图片数，默认 6 张
  limit?: number
  serverPaths: string[]
  localFiles: File[]
  // 中文注释：是否在创建模式下恢复本地草稿图片（默认关闭，避免新建任务显示上次图片）
  restoreDraft?: boolean
  // 中文注释：是否在创建模式下将选择的图片写入本地草稿（默认关闭）
  useDraftStorage?: boolean
}>(), { limit: 6, restoreDraft: false, useDraftStorage: false })

const emit = defineEmits<{
  (e: 'update:serverPaths', v: string[]): void
  (e: 'update:localFiles', v: File[]): void
  (e: 'added'): void
}>()

const fileInput = ref<HTMLInputElement | null>(null)
const items = reactive<Item[]>([])
const DRAFT_KEY = 'task_draft_images'
// 中文注释：预览列表（用于点击缩略图打开大图查看与翻看）
const previewUrls = computed(() => items.map(i => i.url))

function triggerPick() {
  fileInput.value?.click()
}

function pushItemPreview(name: string, url: string, uploading = false) {
  items.push({ key: `${name}-${Date.now()}-${Math.random()}`.replace(/\./g, ''), name, url, uploading, progress: 0 })
}

async function onPicked(e: Event) {
  const input = e.target as HTMLInputElement
  const files = Array.from(input.files || [])
  if (!files.length) return
  const remain = (props.limit ?? 6) - items.length
  if (files.length > remain) {
    ElMessage.warning(`最多选择${props.limit ?? 6}张，已为你截取前${remain}张`)
  }
  const pick = files.slice(0, Math.max(0, remain))

  for (const f of pick) {
    // 大小限制：< 5MB（与后端一致）
    if (f.size > 5 * 1024 * 1024) {
      ElMessage.error(`${f.name} 过大（需小于5MB）`)
      continue
    }
    try {
      const webp = await prepareUpload(f)
      const url = URL.createObjectURL(webp)
      if (props.editing) {
        // 编辑模式：立即上传，显示进度
        const item: Item = { key: `${webp.name}-${Date.now()}`, name: webp.name, url, uploading: true, progress: 0, localFile: webp }
        items.push(item)
        try {
          if (!props.taskId) throw new Error('缺少任务ID')
          const resp = await uploadTaskImage(props.userId, webp, props.taskId, (p) => { item.progress = p })
          item.uploading = false
          item.progress = 100
          item.serverPath = resp.path
          // 同步 v-model：后端返回相对路径，前端统一记录
          emit('update:serverPaths', [...props.serverPaths, resp.path])
          ElMessage.success('图片上传成功')
        } catch (err: any) {
          item.uploading = false
          ElMessage.error(err?.message || '图片上传失败')
        }
  } else {
    // 创建模式：仅本地暂存，记录到 v-model 与草稿
    emit('update:localFiles', [...props.localFiles, webp])
    pushItemPreview(webp.name, url, false)
    // 中文注释：按需写入草稿（默认关闭，避免下次创建页显示历史图片）
    if (props.useDraftStorage) {
      try {
        const reader = new FileReader()
        reader.onload = () => {
          try {
            const raw = localStorage.getItem(DRAFT_KEY)
            const arr: any[] = raw ? JSON.parse(raw) : []
            arr.push({ name: webp.name, type: webp.type || 'image/webp', dataURL: String(reader.result) })
            localStorage.setItem(DRAFT_KEY, JSON.stringify(arr))
          } catch {}
        }
        reader.readAsDataURL(webp)
      } catch {}
    }
    // 中文注释：创建模式添加完成，由父组件统一提示一次，此处不弹窗
  }
    } catch (err: any) {
      ElMessage.error(err?.message || '处理图片失败')
    }
  }
  // 清空选择，避免重复选择同一文件不触发 change
  if (fileInput.value) fileInput.value.value = ''
  // 中文注释：批量选择完成后仅触发一次 added 事件（创建模式）
  if (!props.editing) emit('added')
}

 function removeAt(idx: number) {
  const item = items[idx]
  if (!item) return
  items.splice(idx, 1)
  if (props.editing) {
    // 编辑模式：调用后端删除接口，物理文件 + 数据库 image_json
    if (item.serverPath && props.taskId) {
      deleteTaskImage(props.taskId, item.serverPath)
        .then((resp) => {
          emit('update:serverPaths', resp.images || props.serverPaths.filter((x) => x !== item.serverPath))
          ElMessage.success('已删除服务器图片')
        })
        .catch((err: any) => {
          // 后端删除失败时，至少移除前端列表，避免卡死；并提示错误
          emit('update:serverPaths', props.serverPaths.filter((x) => x !== item.serverPath))
          ElMessage.error(err?.message || '删除服务器图片失败')
        })
    } else if (item.serverPath) {
      // 缺少任务ID时的兜底处理：仅移除前端列表
      emit('update:serverPaths', props.serverPaths.filter((x) => x !== item.serverPath))
    }
  } else {
    // 创建模式：移除本地文件与草稿
    if (item.localFile) {
      emit('update:localFiles', props.localFiles.filter((x) => (x as any).name !== item.localFile!.name))
    } else {
      // 通过名称匹配移除
      emit('update:localFiles', props.localFiles.filter((x) => (x as any).name !== item.name))
    }
    // 中文注释：仅在启用草稿存储时同步清理本地草稿
    if (props.useDraftStorage) {
      try {
        const raw = localStorage.getItem(DRAFT_KEY)
        if (raw) {
          const arr: { name: string; type: string; dataURL: string }[] = JSON.parse(raw)
          const next = arr.filter((x) => x.name !== item.name)
          localStorage.setItem(DRAFT_KEY, JSON.stringify(next))
        }
      } catch {}
    }
  }
}

// 初始化：编辑模式显示已有服务端图片；创建模式恢复草稿
onMounted(() => {
  // 先显示已上传图片
  if (props.serverPaths?.length) {
    for (const raw of props.serverPaths) {
      const p = normalizeUploadPath(raw)
      const base = (import.meta as any).env.VITE_API_BASE || ''
      const full = `${base.replace(/\/+$/, '')}/api/${p}`
      items.push({ key: `${p}-${Date.now()}`, name: p.split('/').pop() || 'image', url: full, uploading: false, progress: 100, serverPath: p })
    }
  }
  // 创建模式：按需恢复草稿缩略图（默认关闭）
  if (!props.editing && props.restoreDraft) {
    try {
      const raw = localStorage.getItem(DRAFT_KEY)
      if (raw) {
        const arr: { name: string; type: string; dataURL: string }[] = JSON.parse(raw)
        for (const x of arr) {
          items.push({ key: `${x.name}-${Date.now()}`, name: x.name, url: x.dataURL, uploading: false, progress: 0 })
        }
      }
    } catch {}
  }
})

// 同步：当 v-model 的 serverPaths 或 localFiles 外部变化时，保持列表一致（简单重建）
watch(() => [props.serverPaths.length, props.localFiles.length, props.editing], () => {
  const base = (import.meta as any).env.VITE_API_BASE || ''
  const next: Item[] = []
  if (props.editing) {
    for (const raw of props.serverPaths) {
      const p = normalizeUploadPath(raw)
      const full = `${base.replace(/\/+$/, '')}/api/${p}`
      next.push({ key: `${p}`, name: p.split('/').pop() || 'image', url: full, uploading: false, progress: 100, serverPath: p })
    }
  } else {
    for (const f of props.localFiles) {
      const url = URL.createObjectURL(f)
      next.push({ key: `${(f as any).name}`, name: (f as any).name, url, uploading: false, progress: 0, localFile: f })
    }
  }
  items.splice(0, items.length, ...next)
})
</script>

<style scoped>
</style>