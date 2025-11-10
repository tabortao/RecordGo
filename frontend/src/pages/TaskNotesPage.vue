<template>
  <!-- 中文注释：任务备注独立页面，顶部返回按钮，支持多条备注与附件（图片/音频） -->
  <div class="p-4">
    <!-- 顶部返回与标题 -->
    <div class="flex items-center gap-2 mb-3">
      <el-button link @click="goBack" class="p-0">
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <h2 class="font-semibold">任务备注</h2>
      <span class="text-xs text-gray-500">ID：{{ taskId }}</span>
    </div>

    <!-- 历史备注列表 -->
    <el-card shadow="never" class="mb-3" v-if="existingNotes.length">
      <div class="font-semibold mb-2">历史备注</div>
      <div class="space-y-3">
        <div v-for="n in existingNotes" :key="n.id" class="border rounded p-2">
          <div class="text-sm whitespace-pre-wrap">{{ n.text }}</div>
          <div class="mt-2 flex flex-wrap gap-3">
            <!-- 图片附件预览，可点击放大 -->
            <template v-for="(att, idx) in n.attachments" :key="att.name + idx">
              <el-image v-if="att.type==='image'" :src="resolveUrl(att)" :preview-src-list="imageList(n.attachments)" :initial-index="imageIndex(n.attachments, att)" fit="cover" style="width:96px;height:96px;border-radius:8px" />
              <div v-else class="flex items-center gap-2">
                <el-icon><Microphone /></el-icon>
                <audio :src="resolveUrl(att)" controls preload="metadata" />
              </div>
            </template>
          </div>
          <div class="text-right mt-2">
            <el-button size="small" type="danger" @click="removeNote(n.id)">删除备注</el-button>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 新建备注 -->
    <el-card shadow="never">
      <div class="font-semibold mb-2">新增备注</div>
      <el-form label-width="90px">
        <el-form-item label="备注内容">
          <!-- 中文注释：Element Plus 的 rows 期望 number，这里使用绑定语法避免字符串类型错误 -->
          <el-input v-model="noteText" type="textarea" :rows="4" placeholder="请输入备注内容，支持多行" />
        </el-form-item>

        <el-form-item label="附件">
          <div class="flex items-center gap-2">
            <input ref="fileInput" type="file" accept="image/*,audio/*" class="hidden" multiple @change="onFileChange" />
            <el-button type="primary" @click="triggerFile">上传附件</el-button>
            <el-button :type="recording ? 'danger' : 'success'" @click="toggleRecord">
              <el-icon class="mr-1"><Microphone /></el-icon>{{ recording ? '停止录音' : '开始录音' }}
            </el-button>
          </div>
          <!-- <div class="text-xs text-gray-500 mt-1">图片自动压缩并转为 webp；音频录音保存为 wav 格式</div> -->
        </el-form-item>

        <el-form-item label="预览区" v-if="attachments.length">
          <div class="flex flex-wrap gap-3">
            <div v-for="(att, idx) in attachments" :key="att.name + idx" class="relative">
              <el-image v-if="att.type==='image'" :src="att.url" :preview-src-list="previewList()" :initial-index="previewIndex(att)" fit="cover" style="width:96px;height:96px;border-radius:8px" />
              <div v-else class="w-[200px]">
                <audio :src="att.url" controls preload="metadata" class="w-full" />
              </div>
              <el-button size="small" type="danger" class="absolute -top-2 -right-2" @click="removeAttachment(idx)">删除</el-button>
            </div>
          </div>
        </el-form-item>

        <div class="text-right">
          <el-button type="primary" @click="saveNote">保存备注</el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
// 中文注释：任务备注页面逻辑，支持图片压缩为 webp 与音频录音保存为 wav
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useNotesStore, type TaskNote, type NoteAttachment } from '@/stores/notes'
import { useAuth } from '@/stores/auth'
import { uploadTaskImage, uploadTaskAudio } from '@/services/tasks'
import { toWebp } from '@/utils/image'
import { ArrowLeft, Microphone } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const rawId = route.params.id
const taskId = Number(Array.isArray(rawId) ? rawId[0] : rawId)
const auth = useAuth()
const store = useNotesStore()

const noteText = ref('')
type TmpAttachment = { type: 'image' | 'audio'; name: string; url: string; file?: File; serverPath?: string }
const attachments = ref<TmpAttachment[]>([])
const fileInput = ref<HTMLInputElement | null>(null)

// 顶部返回
function goBack() { router.back() }

// 历史备注
const existingNotes = computed(() => store.list(taskId))

function resolveUrl(att: NoteAttachment) {
  // 若有后端相对路径则转为可访问的静态资源路径
  if (att.serverPath) return `/api/${att.serverPath}`
  return att.url
}

function imageList(atts: NoteAttachment[]): string[] {
  return atts.filter(a => a.type==='image').map(a => resolveUrl(a))
}

function imageIndex(atts: NoteAttachment[], att: NoteAttachment): number {
  return atts.filter(a => a.type==='image').findIndex(a => a.url === att.url || a.serverPath === att.serverPath)
}

function previewList(): string[] {
  return attachments.value.filter(a => a.type==='image').map(a => a.url)
}

function previewIndex(att: { type: 'image' | 'audio'; url: string }): number {
  return attachments.value.filter(a => a.type==='image').findIndex(a => a.url === att.url)
}

// 触发文件选择
function triggerFile() { fileInput.value?.click() }

// 选择文件（图片/音频）
async function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const files = Array.from(input.files || [])
  for (const f of files) {
    if (f.type.startsWith('image/')) {
      try {
        const webp = await toWebp(f)
        const url = URL.createObjectURL(webp)
        attachments.value.push({ type: 'image', name: webp.name, url, file: webp })
      } catch (err: any) {
        ElMessage.warning(`图片转换失败，已回退原格式：${err?.message || err}`)
        const url = URL.createObjectURL(f)
        attachments.value.push({ type: 'image', name: f.name, url, file: f })
      }
    } else if (f.type.startsWith('audio/')) {
      const url = URL.createObjectURL(f)
      attachments.value.push({ type: 'audio', name: f.name, url, file: f })
    }
  }
  input.value = ''
}

// 录音逻辑：采集 PCM 并编码为 16-bit WAV
let mediaStream: MediaStream | null = null
let audioCtx: AudioContext | null = null
let sourceNode: MediaStreamAudioSourceNode | null = null
let processor: ScriptProcessorNode | null = null
let pcmBuffers: Float32Array[] = []
let sampleRate = 44100
const recording = ref(false)

function toggleRecord() {
  if (recording.value) {
    stopRecord()
  } else {
    startRecord().catch((e) => ElMessage.error(`录音失败：${e?.message || e}`))
  }
}

async function startRecord() {
  mediaStream = await navigator.mediaDevices.getUserMedia({ audio: true })
  audioCtx = new (window.AudioContext || (window as any).webkitAudioContext)()
  sampleRate = audioCtx.sampleRate
  sourceNode = audioCtx.createMediaStreamSource(mediaStream)
  processor = audioCtx.createScriptProcessor(4096, 1, 1)
  pcmBuffers = []
  processor.onaudioprocess = (ev) => {
    const data = ev.inputBuffer.getChannelData(0)
    pcmBuffers.push(new Float32Array(data))
  }
  sourceNode.connect(processor)
  processor.connect(audioCtx.destination)
  recording.value = true
}

function stopRecord() {
  if (!recording.value) return
  processor?.disconnect()
  sourceNode?.disconnect()
  audioCtx?.close()
  mediaStream?.getTracks().forEach(t => t.stop())
  const wavBlob = encodeWav(pcmBuffers, sampleRate)
  const url = URL.createObjectURL(wavBlob)
  // 中文注释：将 Blob 转为 File，以便后续上传到后端
  const file = new File([wavBlob], `rec_${Date.now()}.wav`, { type: 'audio/wav' })
  attachments.value.push({ type: 'audio', name: file.name, url, file })
  recording.value = false
}

function encodeWav(buffers: Float32Array[], sr: number): Blob {
  // 合并 PCM
  let length = buffers.reduce((sum, b) => sum + b.length, 0)
  const pcm = new Float32Array(length)
  let offset = 0
  for (const b of buffers) { pcm.set(b, offset); offset += b.length }
  // 转为 16-bit PCM
  const wavBuffer = new ArrayBuffer(44 + pcm.length * 2)
  const view = new DataView(wavBuffer)
  // 写 WAV 头
  writeString(view, 0, 'RIFF')
  view.setUint32(4, 36 + pcm.length * 2, true)
  writeString(view, 8, 'WAVE')
  writeString(view, 12, 'fmt ')
  view.setUint32(16, 16, true) // PCM 子块大小
  view.setUint16(20, 1, true)  // 编码格式 PCM
  view.setUint16(22, 1, true)  // 声道数 1
  view.setUint32(24, sr, true) // 采样率
  view.setUint32(28, sr * 2, true) // 字节率
  view.setUint16(32, 2, true)  // BlockAlign
  view.setUint16(34, 16, true) // 位深 16
  writeString(view, 36, 'data')
  view.setUint32(40, pcm.length * 2, true)
  // PCM 数据
  let pos = 44
  for (let i = 0; i < pcm.length; i++) {
    let s = Math.max(-1, Math.min(1, pcm[i]))
    view.setInt16(pos, s < 0 ? s * 0x8000 : s * 0x7FFF, true)
    pos += 2
  }
  return new Blob([view], { type: 'audio/wav' })
}

function writeString(view: DataView, offset: number, str: string) {
  for (let i = 0; i < str.length; i++) {
    view.setUint8(offset + i, str.charCodeAt(i))
  }
}

function removeAttachment(idx: number) {
  attachments.value.splice(idx, 1)
}

async function saveNote() {
  const text = noteText.value.trim()
  if (!text && attachments.value.length === 0) {
    return ElMessage.warning('请填写备注或添加附件')
  }
  if (!auth.user) {
    ElMessage.error('请先登录后再保存备注')
    return
  }
  // 图片附件上传到后端（沿用任务图片接口）；音频上传到 /upload/task-audio；均保存 serverPath
  const finalAtts: { type: 'image' | 'audio'; name: string; url: string; serverPath?: string }[] = []
  for (const att of attachments.value) {
    if (att.type === 'image') {
      try {
        if (att.file) {
          const { path } = await uploadTaskImage(auth.user.id, att.file, taskId)
          finalAtts.push({ type: 'image', name: att.name, url: att.url, serverPath: path })
        } else {
          finalAtts.push({ type: 'image', name: att.name, url: att.url })
        }
      } catch (e: any) {
        ElMessage.warning(`图片上传失败：${e?.message || e}`)
        finalAtts.push({ type: 'image', name: att.name, url: att.url })
      }
    } else {
      try {
        if (att.file) {
          const { path } = await uploadTaskAudio(auth.user.id, att.file, taskId)
          // 中文注释：保存 serverPath；前端历史列表播放时使用静态地址 /api/{path}
          finalAtts.push({ type: 'audio', name: att.name, url: att.url, serverPath: path })
        } else {
          finalAtts.push({ type: 'audio', name: att.name, url: att.url })
        }
      } catch (e: any) {
        ElMessage.warning(`音频上传失败：${e?.message || e}`)
        finalAtts.push({ type: 'audio', name: att.name, url: att.url })
      }
    }
  }

  const note: TaskNote = {
    id: Date.now(),
    task_id: taskId,
    text,
    attachments: finalAtts,
    created_at: new Date().toISOString()
  }
  store.add(note)
  noteText.value = ''
  attachments.value = []
  ElMessage.success('备注已保存')
}

function removeNote(noteId: number) {
  store.remove(taskId, noteId)
  ElMessage.success('备注已删除')
}
</script>

<style scoped>
/* 中文注释：备注页面基础样式 */
</style>