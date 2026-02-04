<template>
  <!-- 中文注释：独立的创建任务页面，顶部带返回图标；布局响应式 -->
  <SettingsShell title="创建任务" subtitle="普通创建 · AI 智能创建" :icon="CirclePlusFilled" tone="emerald" container-class="max-w-5xl" back-to="/tasks" :decor="false">
    <div class="rounded-3xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 shadow-sm px-5 py-5">
      <div class="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4">
        <div class="min-w-0">
          <div class="text-[18px] font-extrabold tracking-tight text-gray-900 dark:text-gray-50">两种方式，快速创建</div>
          <div class="mt-1 text-sm leading-relaxed text-gray-600 dark:text-gray-300">
            普通创建适合精确填写；AI 智能创建适合把一段计划批量转成任务。
          </div>
        </div>
        <div class="grid grid-cols-3 gap-2">
          <div class="rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-950/10 px-3 py-2">
            <div class="text-[11px] font-extrabold tracking-[0.22em] text-gray-500 dark:text-gray-400">TIP</div>
            <div class="mt-1 text-xs font-extrabold text-gray-900 dark:text-gray-50">先定分类</div>
          </div>
          <div class="rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-950/10 px-3 py-2">
            <div class="text-[11px] font-extrabold tracking-[0.22em] text-gray-500 dark:text-gray-400">TIP</div>
            <div class="mt-1 text-xs font-extrabold text-gray-900 dark:text-gray-50">再定时长</div>
          </div>
          <div class="rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-950/10 px-3 py-2">
            <div class="text-[11px] font-extrabold tracking-[0.22em] text-gray-500 dark:text-gray-400">TIP</div>
            <div class="mt-1 text-xs font-extrabold text-gray-900 dark:text-gray-50">奖励要具体</div>
          </div>
        </div>
      </div>
    </div>

    <el-tabs v-model="activeTab" class="w-full task-create-tabs">
      <el-tab-pane name="normal">
        <template #label>
          <div class="flex items-center gap-2">
            <el-icon :size="16"><Edit /></el-icon>
            <span>普通创建</span>
          </div>
        </template>
        <!-- 响应式网格：移动端单列，桌面端双列/三列（根据内容） -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <!-- 基础信息分区 -->
          <el-card shadow="never" class="pretty-card">
            <template #header>
              <div class="flex items-center justify-between gap-4">
                <div class="flex items-center gap-2">
                  <div class="h-9 w-9 rounded-2xl bg-emerald-50 dark:bg-emerald-900/25 text-emerald-700 dark:text-emerald-300 grid place-items-center">
                    <el-icon :size="18"><Edit /></el-icon>
                  </div>
                  <div class="min-w-0">
                    <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">基础信息</div>
                    <div class="mt-0.5 text-xs text-gray-500 dark:text-gray-400">标题 · 描述 · 图片 · 分类 · 时长</div>
                  </div>
                </div>
              </div>
            </template>
            <el-form ref="formRefInstance" :model="form" :rules="rules" label-width="100px">
              <el-form-item prop="name" required>
                <template #label>
                  <div class="flex items-center gap-1"><el-icon><Edit /></el-icon><span>任务标题</span></div>
                </template>
                <el-input v-model="form.name" maxlength="128" show-word-limit placeholder="例如：背诵古诗 1 首 / 数学口算 20 题" size="large" />
              </el-form-item>
              <el-form-item label="任务描述" prop="description">
                <template #label>
                  <div class="flex items-center gap-1"><el-icon><List /></el-icon><span>任务描述</span></div>
                </template>
                <el-input v-model="form.description" type="textarea" :rows="4" placeholder="写下要点：怎么做、做到什么程度、需要多久" />
              </el-form-item>
              <el-form-item class="image-upload">
                <template #label>
                  <div class="flex items-center gap-1"><el-icon><Plus /></el-icon><span>任务图片</span></div>
                </template>
                <TaskImageUploader
                  :editing="false"
                  :user-id="userId"
                  :task-id="undefined"
                  :restore-draft="false"
                  :use-draft-storage="false"
                  v-model:serverPaths="form.images"
                  v-model:localFiles="form.local_images"
                  @added="() => ElMessage.success('图片已添加')"
                />
              </el-form-item>
              <el-form-item prop="category" required>
                <template #label>
                  <div class="flex items-center gap-1"><el-icon><List /></el-icon><span>任务分类</span></div>
                </template>
                <el-select v-model="form.category" placeholder="选择分类" style="width: 100%" size="large">
                  <el-option v-for="c in getCategories()" :key="c.name" :label="c.name" :value="c.name">
                    <span class="inline-block w-2 h-2 rounded mr-2" :style="{ backgroundColor: c.color }"></span>
                    <span>{{ c.name }}</span>
                  </el-option>
                </el-select>
              </el-form-item>
              <el-form-item prop="plan_minutes" required>
                <template #label>
                  <div class="flex items-center gap-1"><el-icon><Clock /></el-icon><span>计划时长</span></div>
                </template>
                <el-input-number v-model="form.plan_minutes" :min="1" :max="240" size="large" />
              </el-form-item>
            </el-form>
          </el-card>

          <!-- 计划与重复分区 -->
          <el-card shadow="never" class="pretty-card">
            <template #header>
              <div class="flex items-center justify-between gap-4">
                <div class="flex items-center gap-2">
                  <div class="h-9 w-9 rounded-2xl bg-sky-50 dark:bg-sky-900/25 text-sky-700 dark:text-sky-300 grid place-items-center">
                    <el-icon :size="18"><Clock /></el-icon>
                  </div>
                  <div class="min-w-0">
                    <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">计划与重复</div>
                    <div class="mt-0.5 text-xs text-gray-500 dark:text-gray-400">金币 · 重复 · 日期范围</div>
                  </div>
                </div>
              </div>
            </template>
            <el-form :model="form" label-width="100px">
              <el-form-item prop="score">
                <template #label>
                  <div class="flex items-center gap-1"><el-icon><Coin /></el-icon><span>任务金币</span></div>
                </template>
                <el-input-number v-model="form.score" :min="-10" :max="10" size="large" />
              </el-form-item>
              <el-form-item prop="daily_max_checkins">
                <template #label>
                  <div class="flex items-center gap-1"><el-icon><List /></el-icon><span>每日上限</span></div>
                </template>
                <el-input-number v-model="form.daily_max_checkins" :min="1" :max="20" size="large" />
              </el-form-item>
              <el-form-item prop="repeat_type">
                <template #label>
                  <div class="flex items-center gap-1"><el-icon><List /></el-icon><span>重复类型</span></div>
                </template>
                <el-select v-model="form.repeat_type" placeholder="选择重复类型" style="width: 100%" size="large">
                  <el-option label="无" value="none" />
                  <el-option label="每天" value="daily" />
                  <el-option label="每个工作日" value="weekdays" />
                  <el-option label="每周" value="weekly" />
                  <el-option label="每月" value="monthly" />
                </el-select>
              </el-form-item>
              <el-form-item v-if="form.repeat_type==='weekly'" prop="weekly_days">
                <template #label>
                  <div class="flex items-center gap-1"><el-icon><List /></el-icon><span>选择星期</span></div>
                </template>
                <el-checkbox-group v-model="form.weekly_days">
                  <el-checkbox :label="1">周一</el-checkbox>
                  <el-checkbox :label="2">周二</el-checkbox>
                  <el-checkbox :label="3">周三</el-checkbox>
                  <el-checkbox :label="4">周四</el-checkbox>
                  <el-checkbox :label="5">周五</el-checkbox>
                  <el-checkbox :label="6">周六</el-checkbox>
                  <el-checkbox :label="7">周日</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              <el-form-item prop="start_date" required>
                <template #label>
                  <div class="flex items-center gap-1"><el-icon><Clock /></el-icon><span>开始日期</span></div>
                </template>
                <el-date-picker v-model="form.start_date" type="date" :editable="false" :clearable="false" />
              </el-form-item>
              <el-form-item prop="end_date">
                <template #label>
                  <div class="flex items-center gap-1"><el-icon><Clock /></el-icon><span>截止日期</span></div>
                </template>
                <el-date-picker v-model="form.end_date" type="date" :editable="false" :clearable="false" :disabled="form.repeat_type==='none'" />
              </el-form-item>
            </el-form>
          </el-card>
        </div>

        <!-- 底部操作区 -->
        <div class="mt-5 rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/80 dark:bg-gray-900/65 backdrop-blur-xl shadow-sm px-4 py-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
          <div class="text-xs text-gray-500 dark:text-gray-400 leading-relaxed">
            提示：重复任务建议设置截止日期，避免任务无限延长。
          </div>
          <div class="flex justify-end gap-2">
            <el-button class="!rounded-2xl !font-extrabold" @click="goBack">取消</el-button>
            <el-button type="primary" class="!rounded-2xl !font-extrabold" @click="submitForm">确定</el-button>
          </div>
        </div>
      </el-tab-pane>

      <!-- AI 智能创建 Tab -->
      <el-tab-pane name="ai">
        <template #label>
          <div class="flex items-center gap-2">
            <el-icon :size="16"><MagicStick /></el-icon>
            <span>AI智能创建</span>
          </div>
        </template>
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <el-card shadow="never" class="pretty-card h-fit">
            <template #header>
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <div class="h-9 w-9 rounded-2xl bg-indigo-50 dark:bg-indigo-900/25 text-indigo-700 dark:text-indigo-300 grid place-items-center">
                    <el-icon :size="18"><MagicStick /></el-icon>
                  </div>
                  <div class="min-w-0">
                    <div class="text-sm font-extrabold text-gray-900 dark:text-gray-50">输入需求</div>
                    <div class="mt-0.5 text-xs text-gray-500 dark:text-gray-400">支持自然语言与图片辅助</div>
                  </div>
                </div>
                <el-tooltip content="支持输入自然语言或上传图片，AI将自动提取任务要素" placement="top">
                  <el-icon><QuestionFilled /></el-icon>
                </el-tooltip>
              </div>
            </template>
            <el-input 
              v-model="aiText" 
              type="textarea" 
              :rows="8"
              placeholder="请输入学习计划或任务描述，例如：&#10;1. 每天早上7点背单词30分钟&#10;2. 周末下午跑步5公里"
            />
            <div class="mt-4">
               <el-upload
                  action="#"
                  :auto-upload="false"
                  :on-change="handleAIImageChange"
                  :on-remove="() => { aiImage = undefined }"
                  :limit="1"
                  list-type="picture"
                  :file-list="aiFileList"
                  drag
               >
                  <div class="rounded-3xl border border-dashed border-gray-200 dark:border-gray-700 bg-white/60 dark:bg-gray-950/20 px-4 py-5 text-center">
                    <div class="mx-auto h-11 w-11 rounded-3xl bg-sky-50 dark:bg-sky-900/25 text-sky-700 dark:text-sky-300 grid place-items-center">
                      <el-icon :size="20"><Picture /></el-icon>
                    </div>
                    <div class="mt-3 text-sm font-extrabold text-gray-900 dark:text-gray-50">拖拽图片到这里</div>
                    <div class="mt-1 text-xs text-gray-500 dark:text-gray-400">或点击上传（jpg/png/webp）</div>
                    <div class="mt-4">
                      <el-button type="primary" plain class="!rounded-2xl !font-extrabold">选择图片</el-button>
                    </div>
                  </div>
                  <template #tip>
                    <div class="mt-2 text-[12px] text-gray-500 dark:text-gray-400 leading-5">图片仅用于辅助识别；识别结果可在右侧逐条修改后再创建。</div>
                  </template>
               </el-upload>
            </div>
            <div class="mt-6 flex justify-end">
              <el-button type="primary" class="!rounded-2xl !font-extrabold" :loading="aiLoading" @click="handleAIParse" size="large">
                <el-icon class="mr-1"><MagicStick /></el-icon>开始智能识别
              </el-button>
            </div>
          </el-card>

          <!-- 识别结果 -->
          <div class="space-y-4">
             <div v-if="aiTasks.length > 0" class="flex justify-between items-center bg-white/70 dark:bg-gray-900/55 backdrop-blur border border-gray-100 dark:border-gray-800 p-3 rounded-2xl">
                <h3 class="font-extrabold text-gray-900 dark:text-gray-50">识别结果 ({{ aiTasks.length }})</h3>
                <el-button type="success" class="!rounded-2xl !font-extrabold" @click="submitAITasks" :loading="aiSubmitting">
                   确认创建全部
                </el-button>
             </div>
             
             <el-empty v-if="aiTasks.length === 0 && !aiLoading" description="暂无识别结果，请在左侧输入并点击识别" />

             <div v-loading="aiLoading" class="space-y-3 max-h-[600px] overflow-y-auto pr-2">
                <el-card v-for="(task, idx) in aiTasks" :key="idx" shadow="hover" class="pretty-card relative group">
                   <div class="absolute top-2 right-2 z-10 opacity-0 group-hover:opacity-100 transition-opacity flex gap-2">
                      <el-tag v-if="task.confidence" size="small" :type="task.confidence === 'High' ? 'success' : 'warning'">
                        置信度: {{ task.confidence === 'High' ? '高' : '中' }}
                      </el-tag>
                      <el-button type="danger" circle size="small" :icon="Delete" @click="removeAITask(idx)" />
                   </div>
                   <el-form :model="task" label-width="80px" size="default">
                      <el-form-item label="任务标题" required>
                         <el-input v-model="task.name" />
                      </el-form-item>
                      <el-form-item label="任务分类" required>
                         <el-select v-model="task.category" style="width: 100%">
                            <el-option v-for="c in getCategories()" :key="c.name" :label="c.name" :value="c.name" />
                         </el-select>
                      </el-form-item>
                      <div class="flex gap-2">
                         <el-form-item label="计划时长" class="flex-1" required>
                            <el-input-number v-model="task.plan_minutes" :min="1" :max="240" style="width: 100%" />
                         </el-form-item>
                         <el-form-item label="任务金币" class="flex-1">
                            <el-input-number v-model="task.score" :min="-10" :max="10" style="width: 100%" />
                         </el-form-item>
                      </div>
                      <el-form-item label="重复类型">
                         <el-select v-model="task.repeat_type" style="width: 100%">
                            <el-option label="无" value="none" />
                            <el-option label="每天" value="daily" />
                            <el-option label="每个工作日" value="weekdays" />
                            <el-option label="每周" value="weekly" />
                            <el-option label="每月" value="monthly" />
                         </el-select>
                      </el-form-item>
                      <el-form-item v-if="task.repeat_type==='weekly'" label="选择星期">
                        <el-checkbox-group v-model="task.weekly_days">
                          <el-checkbox :label="1">周一</el-checkbox>
                          <el-checkbox :label="2">周二</el-checkbox>
                          <el-checkbox :label="3">周三</el-checkbox>
                          <el-checkbox :label="4">周四</el-checkbox>
                          <el-checkbox :label="5">周五</el-checkbox>
                          <el-checkbox :label="6">周六</el-checkbox>
                          <el-checkbox :label="7">周日</el-checkbox>
                        </el-checkbox-group>
                      </el-form-item>
                      <div class="flex gap-2">
                        <el-form-item label="开始日期" class="flex-1" required>
                           <el-date-picker v-model="task.start_date" type="date" :editable="false" :clearable="false" style="width: 100%" />
                        </el-form-item>
                        <el-form-item label="截止日期" class="flex-1">
                           <el-date-picker v-model="task.end_date" type="date" :editable="false" :clearable="false" :disabled="task.repeat_type==='none'" style="width: 100%" />
                        </el-form-item>
                      </div>
                   </el-form>
                </el-card>
             </div>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </SettingsShell>
</template>

<script setup lang="ts">
// 中文注释：创建任务页面逻辑，支持普通创建与AI智能创建
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElForm } from 'element-plus'
import { computed } from 'vue'
import { useTaskCategories } from '@/stores/categories'

// 中文注释：联动任务分类设置，创建页下拉选项与颜色一致
const cats = useTaskCategories()
const categories = computed(() => cats.list())
onMounted(async () => { try { await cats.syncFromServer() } catch {} })
import { Plus, Edit, List, Clock, Coin, CirclePlusFilled, QuestionFilled, MagicStick, Delete, Picture } from '@element-plus/icons-vue'
import router from '@/router'
import TaskImageUploader from '@/components/TaskImageUploader.vue'
import SettingsShell from '@/components/settings/SettingsShell.vue'
// 中文注释：补充导入 updateTask，用于在图片上传后写入 image_json，避免未定义错误
import { createTask, uploadTaskImage, updateTask, enqueueOfflineTask, parseTaskByAI, type AITaskParseItem } from '@/services/tasks'
import { prepareUpload } from '@/utils/image'

import { useAuth } from '@/stores/auth'
import { useAIStore } from '@/stores/ai'
const auth = useAuth()
const userId = auth.user?.id ?? 0
function goBack() { router.back() }
function getCategories() { return categories.value }

const activeTab = ref('normal')

// --- 普通创建逻辑 ---
type FormModel = {
  name: string
  description: string
  category: string
  score: number
  daily_max_checkins: number
  plan_minutes: number
  start_date: Date
  end_date?: Date
  images: string[]
  local_images: File[]
  repeat_type: 'none'|'daily'|'weekdays'|'weekly'|'monthly'
  weekly_days: number[]
}
const formRefInstance = ref<InstanceType<typeof ElForm>>()
const form = reactive<FormModel>({
  name: '', description: '', category: '语文', score: 1, plan_minutes: 20,
  daily_max_checkins: 1,
  start_date: new Date(), end_date: undefined,
  images: [], local_images: [], repeat_type: 'none', weekly_days: []
})
const rules = {
  name: [{ required: true, message: '请输入任务标题', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  plan_minutes: [{ required: true, message: '请输入计划时长', trigger: 'change' }],
  start_date: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
}

async function submitForm() {
  if (!formRefInstance.value) return
  await formRefInstance.value.validate(async (valid: boolean) => {
    if (!valid) return
    try {
    const start = new Date(form.start_date)
    const end = form.end_date ? new Date(form.end_date) : undefined

    const created: number[] = []
    const webps: File[] = []
    if ((form.local_images || []).length > 0) {
      const files = [...form.local_images]
      const conv = await Promise.all(files.map(async (f: any) => {
        try { return await prepareUpload(f as File, 0.75) } catch { return f as File }
      }))
      for (const w of conv) webps.push(w as File)
    }
    const t = await createTask({
      user_id: userId,
      name: form.name,
      description: form.description,
      category: form.category,
      score: form.score,
      daily_max_checkins: form.daily_max_checkins,
      plan_minutes: form.plan_minutes,
      start_date: start,
      end_date: end,
      repeat_type: form.repeat_type,
      weekly_days: form.weekly_days || []
    })
    const batchTasks = [t]
    created.push(t.id)

    if (webps.length > 0 && batchTasks.length > 0) {
      for (const t of batchTasks) {
        const paths: string[] = []
        const limit = 3
        for (let j = 0; j < webps.length; j += limit) {
          const chunk = webps.slice(j, j + limit)
          const results = await Promise.all(chunk.map(async (wf: File) => {
            try { const { path } = await uploadTaskImage(userId, wf, t.id); return path } catch { return '' }
          }))
          for (const p of results) { if (p) paths.push(p) }
        }
        if (paths.length > 0) { try { await updateTask(t.id, { image_json: JSON.stringify(paths) }) } catch {} }
      }
    }
    ElMessage.success(`创建成功${created.length>1?`（${created.length}条）`:''}`)
    router.back()
  } catch (e: any) {
    // 离线逻辑...
    try {
      const images: { name: string; type: string; dataURL: string }[] = []
      for (const f of form.local_images) {
        const reader = new FileReader()
        await new Promise<void>((resolve) => {
          reader.onload = () => { images.push({ name: (f as any).name || 'image.webp', type: 'image/webp', dataURL: String(reader.result || '') }); resolve() }
          reader.readAsDataURL(f as any)
        })
      }
      enqueueOfflineTask({
        name: form.name,
        description: form.description,
        category: form.category,
        score: form.score,
        daily_max_checkins: form.daily_max_checkins,
        plan_minutes: form.plan_minutes,
        start_date: new Date(form.start_date).toISOString(),
        end_date: form.end_date ? new Date(form.end_date).toISOString() : undefined,
        repeat_type: form.repeat_type,
        weekly_days: form.weekly_days || [],
        images
      })
      ElMessage.success('网络异常，任务已保存到本地，恢复网络后将自动同步')
      router.back()
    } catch {
      ElMessage.error(e.message || '创建失败')
    }
  }
  })
}

import { recognizeImage } from '@/services/ocr'

// --- AI 智能创建逻辑 ---
const aiText = ref('')
const aiFileList = ref<any[]>([])
const aiImage = ref<File>()
const isCompressing = ref(false)
const aiLoading = ref(false)
const aiLoadingText = ref('开始智能识别')
const aiTasks = ref<AITaskParseItem[]>([])
const aiSubmitting = ref(false)
const aiStore = useAIStore()

async function handleAIImageChange(file: any) {
  aiImage.value = file.raw
  aiFileList.value = [file]
  // 自动压缩
  isCompressing.value = true
  try {
    const compressed = await prepareUpload(file.raw, 0.6)
    aiImage.value = compressed
  } catch (e) {
    console.error('Auto compression failed', e)
  } finally {
    isCompressing.value = false
  }
}

async function handleAIParse() {
  if (!aiText.value && !aiImage.value && !isCompressing.value) {
    ElMessage.warning('请输入文本或上传图片')
    return
  }
  
  if (isCompressing.value) {
    aiLoadingText.value = '正在压缩图片...'
    aiLoading.value = true
    while (isCompressing.value) {
      await new Promise(resolve => setTimeout(resolve, 200))
    }
    // 压缩完成后，如果用户清空了图片，需要重新检查
    if (!aiText.value && !aiImage.value) {
      aiLoading.value = false
      aiLoadingText.value = '开始智能识别'
      ElMessage.warning('请输入文本或上传图片')
      return
    }
  }
  
  aiLoading.value = true
  aiLoadingText.value = '准备中...'
  
  let finalText = aiText.value
  let useVisionModel = !!aiImage.value

  // OCR 优先逻辑
  if (aiImage.value) {
    const ocrSettings = aiStore.ocr
    const activeOCR = ocrSettings?.activeModel
    const ocrConfig = activeOCR ? ocrSettings.configs[activeOCR] : null

    if (ocrConfig && ocrConfig.apiUrl && ocrConfig.token) {
      try {
        ElMessage.info('正在使用 OCR 识别图片内容...')
        // 图片已在上传时自动压缩
        const imageToOCR = aiImage.value

        aiLoadingText.value = '正在识别文字...'
        const ocrResult = await recognizeImage(imageToOCR!, activeOCR, ocrConfig)
        if (ocrResult && ocrResult.trim().length > 0) {
          finalText = (finalText ? finalText + '\n\n' : '') + `[图片OCR识别结果]:\n${ocrResult}`
          useVisionModel = false // OCR 成功，转为纯文本模式
          ElMessage.success('OCR 识别成功，正在进行智能分析...')
        } else {
          console.warn('OCR result empty, falling back to Vision Model')
        }
      } catch (e) {
        console.error('OCR Error, falling back to Vision Model', e)
        ElMessage.warning('OCR 服务调用失败，自动回退至 AI 视觉模型')
      }
    } else {
      // 检查 AI 视觉模型配置
      const config = aiStore.$state
      // PRD 要求：若AI视觉模型也未配置，则显示提示
      if (!config.visionApiKey && !config.apiKey && !config.visionApiBaseUrl) {
         ElMessage.error('请先配置AI视觉模型或设置OCR服务的API_URL和TOKEN')
         aiLoading.value = false
         aiLoadingText.value = '开始智能识别'
         return
      }
    }
  }

  try {
    aiLoadingText.value = '正在智能分析...'
    const config = aiStore.$state
    const hasImage = useVisionModel
    const aiConfig = {
      url: hasImage ? (config.visionApiBaseUrl || config.apiBaseUrl) : config.apiBaseUrl,
      key: hasImage ? (config.visionApiKey || config.apiKey) : config.apiKey,
      model: hasImage ? (config.visionModel || config.model) : config.model
    }
    const categoryNames = categories.value.map(c => c.name)
    // 如果 OCR 成功，imageToPass 为 undefined，parseTaskByAI 将只处理文本
    // 如果 OCR 失败或未配置，且 useVisionModel 为真，则传递图片
    // 图片已在上传时自动压缩
    const imageToPass = useVisionModel ? aiImage.value : undefined
    
    console.log('[AI Task Create] Final Text sent to AI:', finalText)

    const res = await parseTaskByAI(finalText, imageToPass, aiConfig, categoryNames)
    if (res.tasks && res.tasks.length > 0) {
      aiTasks.value = res.tasks
      ElMessage.success(`成功识别 ${res.tasks.length} 个任务`)
    } else {
      ElMessage.info('未识别到任务，请尝试调整输入')
    }
  } catch (e: any) {
    ElMessage.error(e.message || 'AI 识别失败')
  } finally {
    aiLoading.value = false
    aiLoadingText.value = '开始智能识别'
  }
}

function removeAITask(index: number) {
  aiTasks.value.splice(index, 1)
}

async function submitAITasks() {
  if (aiTasks.value.length === 0) return
  aiSubmitting.value = true
  let successCount = 0
  const failedTasks: AITaskParseItem[] = []

  try {
    // 串行创建，避免并发过高
    for (const task of aiTasks.value) {
      try {
        await createTask({
          user_id: userId,
          name: task.name,
          description: task.description,
          category: task.category,
          score: task.score,
          daily_max_checkins: 1,
          plan_minutes: task.plan_minutes,
          start_date: new Date(task.start_date || new Date()),
          end_date: task.end_date ? new Date(task.end_date) : undefined,
          repeat_type: task.repeat_type || 'none',
          weekly_days: task.weekly_days || []
        })
        successCount++
      } catch (e) {
        console.error('Failed to create task:', task.name, e)
        failedTasks.push(task)
      }
    }
    if (failedTasks.length === 0) {
      ElMessage.success(`成功创建 ${successCount} 个任务`)
      router.back()
    } else {
      aiTasks.value = failedTasks
      ElMessage.warning(`成功创建 ${successCount} 个任务，失败 ${failedTasks.length} 个，请检查剩余任务`)
    }
  } catch (e: any) {
    ElMessage.error(`创建过程异常：${e.message}`)
  } finally {
    aiSubmitting.value = false
  }
}

</script>

<style scoped>
:deep(input.el-input__inner),
:deep(textarea.el-textarea__inner) {
  font-size: 16px;
}

:deep(.pretty-card .el-card__header) {
  border-bottom: 1px solid rgb(255 255 255 / 0.45);
  padding: 16px 16px 12px;
}

.dark :deep(.pretty-card .el-card__header) {
  border-bottom: 1px solid rgb(31 41 55 / 0.55);
}

:deep(.pretty-card .el-card__body) {
  padding: 16px;
}

:deep(.task-create-tabs .el-tabs__header) {
  margin: 14px 0 12px;
}

:deep(.task-create-tabs .el-tabs__content) {
  padding-top: 6px;
}

:deep(.task-create-tabs .el-tabs__nav-wrap::after) {
  height: 0;
}

:deep(.task-create-tabs .el-tabs__nav) {
  background: rgb(255 255 255 / 0.75);
  border: 1px solid rgb(255 255 255 / 0.5);
  backdrop-filter: blur(16px);
  border-radius: 9999px;
  padding: 4px;
  box-shadow: 0 10px 28px rgb(0 0 0 / 0.06);
}

.dark :deep(.task-create-tabs .el-tabs__nav) {
  background: rgb(17 24 39 / 0.65);
  border: 1px solid rgb(31 41 55 / 0.6);
}

:deep(.task-create-tabs .el-tabs__item) {
  height: 36px;
  line-height: 36px;
  font-size: 13px;
  font-weight: 800;
  border-radius: 9999px;
  padding: 0 16px;
  box-sizing: border-box;
  color: rgb(55 65 81);
  transition: background-color 180ms ease, color 180ms ease, box-shadow 220ms ease, transform 220ms ease;
}

:deep(.task-create-tabs .el-tabs__item:hover) {
  background: rgb(255 255 255 / 0.55);
}

.dark :deep(.task-create-tabs .el-tabs__item) {
  color: rgb(229 231 235);
}

.dark :deep(.task-create-tabs .el-tabs__item:hover) {
  background: rgb(31 41 55 / 0.45);
}

:deep(.task-create-tabs .el-tabs__item.is-active) {
  color: rgb(4 120 87);
  background: rgb(16 185 129 / 0.16);
  border: 1px solid rgb(16 185 129 / 0.28);
  box-shadow: 0 10px 24px rgb(16 185 129 / 0.12);
  transform: translateY(-0.5px);
}

.dark :deep(.task-create-tabs .el-tabs__item.is-active) {
  color: rgb(167 243 208);
  background: rgb(16 185 129 / 0.12);
  border: 1px solid rgb(16 185 129 / 0.2);
  box-shadow: 0 10px 24px rgb(0 0 0 / 0.35);
  transform: translateY(-0.5px);
}

:deep(.task-create-tabs .el-tabs__active-bar) {
  display: none;
}

:deep(.pretty-card.el-card) {
  border-radius: 24px;
  border: 1px solid rgb(255 255 255 / 0.5);
  background: rgb(255 255 255 / 0.8);
  backdrop-filter: blur(16px);
  box-shadow: 0 1px 2px rgb(0 0 0 / 0.04);
}

.dark :deep(.pretty-card.el-card) {
  border: 1px solid rgb(31 41 55 / 0.6);
  background: rgb(17 24 39 / 0.65);
}

:deep(.pretty-card .el-input__wrapper),
:deep(.pretty-card .el-textarea__inner),
:deep(.pretty-card .el-select__wrapper) {
  border-radius: 16px;
  border: 1px solid rgb(203 213 225);
  background: rgb(255 255 255 / 0.92);
}

.dark :deep(.pretty-card .el-input__wrapper),
.dark :deep(.pretty-card .el-textarea__inner),
.dark :deep(.pretty-card .el-select__wrapper) {
  border: 1px solid rgb(51 65 85);
  background: rgb(2 6 23 / 0.18);
}

:deep(.pretty-card .el-input__wrapper),
:deep(.pretty-card .el-select__wrapper) {
  box-shadow: none;
}

:deep(.pretty-card .el-textarea__inner) {
  box-shadow: none;
}

:deep(.pretty-card .el-input__wrapper.is-focus),
:deep(.pretty-card .el-select__wrapper.is-focused) {
  box-shadow: 0 0 0 4px rgb(16 185 129 / 0.16);
  border-color: rgb(16 185 129);
}

.dark :deep(.pretty-card .el-input__wrapper.is-focus),
.dark :deep(.pretty-card .el-select__wrapper.is-focused) {
  box-shadow: 0 0 0 4px rgb(16 185 129 / 0.14);
  border-color: rgb(16 185 129);
}

:deep(.pretty-card .el-textarea__inner:focus) {
  box-shadow: 0 0 0 4px rgb(16 185 129 / 0.16);
  border-color: rgb(16 185 129);
}

.dark :deep(.pretty-card .el-textarea__inner:focus) {
  box-shadow: 0 0 0 4px rgb(16 185 129 / 0.14);
  border-color: rgb(16 185 129);
}
</style>
