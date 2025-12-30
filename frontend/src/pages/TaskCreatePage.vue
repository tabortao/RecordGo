<template>
  <!-- 中文注释：独立的创建任务页面，顶部带返回图标；布局响应式 -->
  <div class="p-4 space-y-4">
    <!-- 顶部栏：返回 + 标题 -->
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer text-gray-600" @click="goBack"><ArrowLeft /></el-icon>
      <div class="w-8 h-8 rounded-full flex items-center justify-center bg-emerald-100 dark:bg-emerald-900">
        <el-icon :size="20" class="text-emerald-600"><CirclePlusFilled /></el-icon>
      </div>
      <h2 class="font-semibold">创建任务</h2>
    </div>

    <el-tabs v-model="activeTab" class="w-full">
      <el-tab-pane label="普通创建" name="normal">
        <!-- 响应式网格：移动端单列，桌面端双列/三列（根据内容） -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <!-- 基础信息分区 -->
          <el-card shadow="never" class="rounded-lg">
            <el-form ref="formRefInstance" :model="form" :rules="rules" label-width="100px">
              <el-form-item prop="name" required>
                <template #label>
                  <div class="flex items-center gap-1"><el-icon><Edit /></el-icon><span>任务标题</span></div>
                </template>
                <el-input v-model="form.name" maxlength="128" show-word-limit />
              </el-form-item>
              <el-form-item label="任务描述" prop="description">
                <template #label>
                  <div class="flex items-center gap-1"><el-icon><List /></el-icon><span>任务描述</span></div>
                </template>
                <el-input v-model="form.description" type="textarea" />
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
                <el-select v-model="form.category" placeholder="选择分类" style="width: 100%">
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
                <el-input-number v-model="form.plan_minutes" :min="1" :max="240" />
              </el-form-item>
            </el-form>
          </el-card>

          <!-- 计划与重复分区 -->
          <el-card shadow="never" class="rounded-lg">
            <el-form :model="form" label-width="100px">
              <el-form-item prop="score">
                <template #label>
                  <div class="flex items-center gap-1"><el-icon><Coin /></el-icon><span>任务金币</span></div>
                </template>
                <el-input-number v-model="form.score" :min="-10" :max="10" />
              </el-form-item>
              <el-form-item prop="repeat_type">
                <template #label>
                  <div class="flex items-center gap-1"><el-icon><List /></el-icon><span>重复类型</span></div>
                </template>
                <el-select v-model="form.repeat_type" placeholder="选择重复类型" style="width: 100%">
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
        <div class="flex justify-end gap-2 mt-4">
          <el-button @click="goBack">取消</el-button>
          <el-button type="primary" @click="submitForm">确定</el-button>
        </div>
      </el-tab-pane>

      <!-- AI 智能创建 Tab -->
      <el-tab-pane label="AI智能创建" name="ai">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <el-card shadow="never" class="rounded-lg h-fit">
            <template #header>
              <div class="flex items-center justify-between">
                <span>输入需求</span>
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
               >
                  <el-button type="primary" plain :icon="Picture">上传图片识别</el-button>
                  <template #tip>
                    <div class="el-upload__tip">支持 jpg/png/webp，图片仅用于辅助识别</div>
                    <div class="text-[12px] text-[#999999] mt-2">AI智能创建任务正在测试，可能出现不准确的情况。如遇到问题，欢迎反馈给我们！</div>
                  </template>
               </el-upload>
            </div>
            <div class="mt-6 flex justify-end">
              <el-button type="primary" :loading="aiLoading" @click="handleAIParse" size="large">
                <el-icon class="mr-1"><MagicStick /></el-icon>开始智能识别
              </el-button>
            </div>
          </el-card>

          <!-- 识别结果 -->
          <div class="space-y-4">
             <div v-if="aiTasks.length > 0" class="flex justify-between items-center bg-gray-50 dark:bg-gray-800 p-3 rounded-lg">
                <h3 class="font-bold text-gray-700 dark:text-gray-200">识别结果 ({{ aiTasks.length }})</h3>
                <el-button type="success" @click="submitAITasks" :loading="aiSubmitting">
                   确认创建全部
                </el-button>
             </div>
             
             <el-empty v-if="aiTasks.length === 0 && !aiLoading" description="暂无识别结果，请在左侧输入并点击识别" />

             <div v-loading="aiLoading" class="space-y-3 max-h-[600px] overflow-y-auto pr-2">
                <el-card v-for="(task, idx) in aiTasks" :key="idx" shadow="hover" class="relative group">
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
  </div>
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
import { ArrowLeft, Plus, Edit, List, Clock, Coin, CirclePlusFilled, QuestionFilled, MagicStick, Delete, Picture } from '@element-plus/icons-vue'
import router from '@/router'
import TaskImageUploader from '@/components/TaskImageUploader.vue'
// 中文注释：补充导入 updateTask，用于在图片上传后写入 image_json，避免未定义错误
import { createTask, uploadTaskImage, updateTask, enqueueOfflineTask, parseTaskByAI, type AITaskParseItem } from '@/services/tasks'
import { prepareUpload } from '@/utils/image'

import { useAuth } from '@/stores/auth'
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

import { useAIStore } from '@/stores/ai'

// --- AI 智能创建逻辑 ---
const aiText = ref('')
const aiFileList = ref<any[]>([])
const aiImage = ref<File>()
const aiLoading = ref(false)
const aiTasks = ref<AITaskParseItem[]>([])
const aiSubmitting = ref(false)
const aiStore = useAIStore()

function handleAIImageChange(file: any) {
  aiImage.value = file.raw
  aiFileList.value = [file]
}

async function handleAIParse() {
  if (!aiText.value && !aiImage.value) {
    ElMessage.warning('请输入文本或上传图片')
    return
  }
  aiLoading.value = true
  try {
    const config = aiStore.$state
    const hasImage = !!aiImage.value
    const aiConfig = {
      url: hasImage ? (config.visionApiBaseUrl || config.apiBaseUrl) : config.apiBaseUrl,
      key: hasImage ? (config.visionApiKey || config.apiKey) : config.apiKey,
      model: hasImage ? (config.visionModel || config.model) : config.model
    }
    const categoryNames = categories.value.map(c => c.name)
    const res = await parseTaskByAI(aiText.value, aiImage.value, aiConfig, categoryNames)
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
</style>
