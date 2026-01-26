<template>
  <SettingsShell
    title="任务备注"
    :subtitle="`任务ID：${taskId}${existingNotes.length ? ` · ${existingNotes.length} 条` : ''}`"
    tone="indigo"
    :decor="true"
  >
    <SettingsCard
      title="历史备注"
      :description="existingNotes.length ? '支持图片预览、音频播放与一键朗读' : '还没有备注，写下今天的收获吧'"
    >
      <div v-if="!existingNotes.length" class="py-4">
        <el-empty description="暂无备注" />
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="n in existingNotes"
          :key="n.id"
          class="rounded-2xl border border-gray-100/70 dark:border-gray-800/60 bg-white/70 dark:bg-gray-950/35 backdrop-blur p-3"
        >
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0 flex-1">
              <div class="text-[13px] text-gray-500 dark:text-gray-400">
                {{ formatTime(n.created_at) }}
              </div>
              <div class="mt-2 text-sm leading-relaxed whitespace-pre-wrap break-words">
                {{ n.text }}
              </div>
            </div>

            <div class="shrink-0 flex items-center gap-2">
              <el-button
                v-if="appState.speech.enabled"
                size="small"
                round
                class="!px-3"
                @click="speakNote(n.text)"
              >
                <el-icon class="mr-1"><Headset /></el-icon>
                朗读
              </el-button>
              <el-button size="small" round type="danger" class="!px-3" @click="removeNote(n.id)">删除</el-button>
            </div>
          </div>

          <div v-if="n.attachments?.length" class="mt-3">
            <div class="grid grid-cols-3 sm:grid-cols-4 gap-3">
              <template v-for="(att, idx) in n.attachments" :key="att.name + idx">
                <el-image
                  v-if="att.type === 'image'"
                  class="rounded-2xl overflow-hidden border border-gray-100/60 dark:border-gray-800/60 bg-white/60 dark:bg-gray-900/40"
                  :src="resolveUrl(att)"
                  :preview-src-list="imageList(n.attachments)"
                  :initial-index="imageIndex(n.attachments, att)"
                  fit="cover"
                  style="width:100%;aspect-ratio:1/1"
                />
                <div
                  v-else
                  class="col-span-3 sm:col-span-4 rounded-2xl border border-gray-100/60 dark:border-gray-800/60 bg-white/60 dark:bg-gray-900/40 px-3 py-2"
                >
                  <div class="flex items-center gap-2 text-xs text-gray-500 dark:text-gray-400">
                    <el-icon :size="16"><Microphone /></el-icon>
                    音频附件
                  </div>
                  <audio :src="resolveUrl(att)" controls preload="none" class="mt-2 w-full" />
                </div>
              </template>
            </div>
          </div>
        </div>
      </div>
    </SettingsCard>

    <SettingsCard title="新增备注" description="支持多行文本、图片附件与录音">
      <el-form label-position="top">
        <el-form-item label="备注内容">
          <el-input
            v-model="noteText"
            type="textarea"
            :rows="5"
            resize="none"
            placeholder="写下今天的进步、遇到的问题、下一步计划…"
          />
        </el-form-item>

        <el-form-item label="附件">
          <div class="flex flex-wrap items-center gap-2">
            <input ref="fileInput" type="file" accept="image/*,audio/*" class="hidden" multiple @change="onFileChange" />
            <el-button round type="primary" class="!px-4" @click="triggerFile">上传附件</el-button>
            <el-button round :type="recording ? 'danger' : 'success'" class="!px-4" @click="toggleRecord">
              <el-icon class="mr-1"><Microphone /></el-icon>{{ recording ? '停止录音' : '开始录音' }}
            </el-button>
          </div>
        </el-form-item>

        <el-form-item v-if="attachments.length" label="预览">
          <div class="grid grid-cols-3 sm:grid-cols-4 gap-3 w-full">
            <div v-for="(att, idx) in attachments" :key="att.name + idx" class="relative">
              <el-image
                v-if="att.type === 'image'"
                class="rounded-2xl overflow-hidden border border-gray-100/60 dark:border-gray-800/60 bg-white/60 dark:bg-gray-900/40"
                :src="att.url"
                :preview-src-list="previewList()"
                :initial-index="previewIndex(att)"
                fit="cover"
                style="width:100%;aspect-ratio:1/1"
              />
              <div
                v-else
                class="col-span-3 sm:col-span-4 rounded-2xl border border-gray-100/60 dark:border-gray-800/60 bg-white/60 dark:bg-gray-900/40 px-3 py-2"
              >
                <div class="flex items-center gap-2 text-xs text-gray-500 dark:text-gray-400">
                  <el-icon :size="16"><Microphone /></el-icon>
                  录音预览
                </div>
                <audio :src="att.url" controls preload="none" class="mt-2 w-full" />
              </div>

              <button
                type="button"
                class="absolute -top-2 -right-2 inline-flex h-8 w-8 items-center justify-center rounded-2xl border border-white/60 dark:border-gray-800/60 bg-white/85 dark:bg-gray-900/75 backdrop-blur text-gray-700 dark:text-gray-200 shadow-sm hover:bg-white dark:hover:bg-gray-900 transition-colors"
                @click="removeAttachment(idx)"
              >
                <el-icon :size="16"><Close /></el-icon>
              </button>
            </div>
          </div>
        </el-form-item>

        <div class="pt-2">
          <el-button round type="primary" class="w-full !h-11 !text-[15px]" @click="saveNote">保存备注</el-button>
        </div>
      </el-form>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
// 中文注释：任务备注页面逻辑，支持图片压缩为 webp 与音频录音保存为 wav
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useNotesStore, type TaskNote, type NoteAttachment } from '@/stores/notes'
import { useAuth } from '@/stores/auth'
import { uploadNoteImage, uploadTaskAudio } from '@/services/tasks'
import { normalizeUploadPath } from '@/services/wishes'
import { prepareUpload } from '@/utils/image'
import { Close, Headset, Microphone } from '@element-plus/icons-vue'
import { speak } from '@/utils/speech'
import { useAppState } from '@/stores/appState'
import { ElMessage } from 'element-plus'
import { getStaticBase } from '@/services/http'
import { presignView } from '@/services/storage'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

const route = useRoute()
const rawId = route.params.id
const taskId = Number(Array.isArray(rawId) ? rawId[0] : rawId)
const auth = useAuth()
const store = useNotesStore()
const appState = useAppState()

const noteText = ref('')
type TmpAttachment = { type: 'image' | 'audio'; name: string; url: string; file?: File; serverPath?: string }
const attachments = ref<TmpAttachment[]>([])
const fileInput = ref<HTMLInputElement | null>(null)

// 历史备注
const existingNotes = computed<TaskNote[]>(() => store.list(taskId))
onMounted(async () => { await resolveServerPaths(existingNotes.value) })
watch(() => JSON.stringify(existingNotes.value.map((n: TaskNote) => ((n.attachments||[]) as NoteAttachment[]).map((a: NoteAttachment) => a.serverPath || '').join(','))), async () => {
  await resolveServerPaths(existingNotes.value)
})
const resolvedMap = ref<Record<string,string>>({})
async function resolveServerPaths(notes: TaskNote[]) {
  const keys: string[] = []
  for (const n of notes) {
    for (const a of n.attachments) {
      const p = normalizeUploadPath(a.serverPath || '')
      if (p && !p.startsWith('uploads/')) keys.push(p)
    }
  }
  const uniques = Array.from(new Set(keys))
  await Promise.all(uniques.map(async (k) => { try { resolvedMap.value[k] = await presignView(k) } catch {} }))
}

function resolveUrl(att: NoteAttachment) {
  if (att.serverPath) {
    const base = getStaticBase()
    const rel = normalizeUploadPath(att.serverPath)
    if (rel.startsWith('uploads/')) return `${base}/api/${rel}`
    return resolvedMap.value[rel] || att.url
  }
  return att.url
}

function imageList(atts: NoteAttachment[]): string[] {
  return atts.filter((a: NoteAttachment) => a.type === 'image').map((a: NoteAttachment) => resolveUrl(a)).filter(Boolean)
}

function imageIndex(atts: NoteAttachment[], att: NoteAttachment): number {
  return atts.filter((a: NoteAttachment) => a.type === 'image').findIndex((a: NoteAttachment) => a.url === att.url || a.serverPath === att.serverPath)
}

function previewList(): string[] {
  return attachments.value.filter((a: TmpAttachment) => a.type === 'image').map((a: TmpAttachment) => a.url)
}

function previewIndex(att: { type: 'image' | 'audio'; url: string }): number {
  return attachments.value.filter((a: TmpAttachment) => a.type === 'image').findIndex((a: TmpAttachment) => a.url === att.url)
}

function formatTime(v: string) {
  const t = new Date(v)
  if (Number.isNaN(t.getTime())) return ''
  return t.toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
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
        const webp = await prepareUpload(f, 0.8)
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
        const { path } = await uploadNoteImage(auth.user.id, att.file, taskId)
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
  try { await resolveServerPaths([note]) } catch {}
  noteText.value = ''
  attachments.value = []
  ElMessage.success('备注已保存')
}

function removeNote(noteId: number) {
  store.remove(taskId, noteId)
  ElMessage.success('备注已删除')
}

// 中文注释：朗读备注内容，遵循全局朗读设置（语音/语速/音调）
async function speakNote(text: string) {
  try {
    if (!appState.speech.enabled) return
    const s = (text || '').trim()
    if (!s) return
    const ok = await speak(s, { voiceURI: appState.speech.voiceURI || undefined, rate: appState.speech.rate, pitch: appState.speech.pitch })
    if (!ok) ElMessage.warning('当前浏览器不支持朗读或语音不可用')
  } catch {
    // 忽略错误，避免影响页面交互
  }
}
</script>

<style scoped>
/* 中文注释：备注页面基础样式 */
</style>
