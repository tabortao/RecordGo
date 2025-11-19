<template>
  <!-- 中文注释：独立的创建任务页面，顶部带返回图标；布局响应式 -->
  <div class="p-4 space-y-4">
    <!-- 顶部栏：返回 + 标题 -->
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer text-gray-600" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="27" class="text-emerald-600"><CirclePlusFilled /></el-icon>
      <h2 class="font-semibold">创建任务</h2>
    </div>

    <!-- 响应式网格：移动端单列，桌面端双列/三列（根据内容） -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <!-- 基础信息分区 -->
      <el-card shadow="never" class="rounded-lg">
        <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
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
              <el-option v-for="c in categories" :key="c.name" :label="c.name" :value="c.name">
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
    <div class="flex justify-end gap-2">
      <el-button @click="goBack">取消</el-button>
      <el-button type="primary" @click="submitForm">确定</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
// 中文注释：创建任务页面逻辑，与列表页表单一致但独立路由展示
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { computed } from 'vue'
import { useTaskCategories } from '@/stores/categories'

// 中文注释：联动任务分类设置，创建页下拉选项与颜色一致
const cats = useTaskCategories()
const categories = computed(() => cats.list())
import { ArrowLeft, Plus, Edit, List, Clock, Coin, CirclePlusFilled } from '@element-plus/icons-vue'
import router from '@/router'
import TaskImageUploader from '@/components/TaskImageUploader.vue'
// 中文注释：补充导入 updateTask，用于在图片上传后写入 image_json，避免未定义错误
import { createTask, uploadTaskImage, updateTask, enqueueOfflineTask } from '@/services/tasks'
import { prepareUpload } from '@/utils/image'

import { useAuth } from '@/stores/auth'
const auth = useAuth()
const userId = auth.user?.id ?? 0
function goBack() { router.back() }

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
const formRef = ref()
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
  try {
    const start = new Date(form.start_date)
    const end = form.end_date ? new Date(form.end_date) : undefined
    // 改为仅创建一条任务，重复规则保存在后端字段，不再为每一天生成多条任务

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
}
</script>

<style scoped>
:deep(input.el-input__inner),
:deep(textarea.el-textarea__inner) {
  font-size: 16px;
}
</style>