<template>
  <!-- 中文注释：独立的编辑任务页面，顶部带返回图标；布局响应式 -->
  <div class="p-4 space-y-4" v-if="taskLoaded">
    <!-- 顶部栏：返回 + 标题 -->
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer text-gray-600" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" class="text-amber-600"><Edit /></el-icon>
      <h2 class="font-semibold">编辑任务</h2>
    </div>

    <!-- 响应式网格：移动端单列，桌面端双列 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <!-- 基础信息分区 -->
      <el-card shadow="never" class="rounded-lg">
        <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
          <el-form-item prop="name" required>
            <template #label><div class="flex items-center gap-1"><el-icon><Edit /></el-icon><span>任务标题</span></div></template>
            <el-input v-model="form.name" maxlength="128" show-word-limit />
          </el-form-item>
          <el-form-item label="任务描述" prop="description">
            <template #label><div class="flex items-center gap-1"><el-icon><List /></el-icon><span>任务描述</span></div></template>
            <el-input v-model="form.description" type="textarea" />
          </el-form-item>
          <el-form-item class="image-upload">
            <template #label><div class="flex items-center gap-1"><el-icon><Plus /></el-icon><span>任务图片</span></div></template>
            <TaskImageUploader
              :editing="true"
              :user-id="userId"
              :task-id="taskId"
              v-model:serverPaths="form.images"
              v-model:localFiles="form.local_images"
              @added="() => ElMessage.success('图片已添加')"
            />
          </el-form-item>
          <el-form-item prop="category" required>
            <template #label><div class="flex items-center gap-1"><el-icon><List /></el-icon><span>任务分类</span></div></template>
            <el-select v-model="form.category" placeholder="选择分类" style="width: 100%">
              <el-option label="语文" value="语文" />
              <el-option label="数学" value="数学" />
              <el-option label="英语" value="英语" />
            </el-select>
          </el-form-item>
          <el-form-item prop="plan_minutes" required>
            <template #label><div class="flex items-center gap-1"><el-icon><Clock /></el-icon><span>计划时长</span></div></template>
            <el-input-number v-model="form.plan_minutes" :min="1" :max="240" />
          </el-form-item>
        </el-form>
      </el-card>

      <!-- 计划与重复分区 -->
      <el-card shadow="never" class="rounded-lg">
        <el-form :model="form" label-width="100px">
          <el-form-item prop="score">
            <template #label><div class="flex items-center gap-1"><el-icon><Coin /></el-icon><span>任务金币</span></div></template>
            <el-input-number v-model="form.score" :min="-10" :max="10" />
          </el-form-item>
          <el-form-item prop="repeat_type">
            <template #label><div class="flex items-center gap-1"><el-icon><List /></el-icon><span>重复类型</span></div></template>
            <el-select v-model="form.repeat_type" placeholder="选择重复类型" style="width: 100%">
              <el-option label="无" value="none" />
              <el-option label="每天" value="daily" />
              <el-option label="每个工作日" value="weekdays" />
              <el-option label="每周" value="weekly" />
              <el-option label="每月" value="monthly" />
            </el-select>
          </el-form-item>
          <el-form-item v-if="form.repeat_type==='weekly'" prop="weekly_days">
            <template #label><div class="flex items-center gap-1"><el-icon><List /></el-icon><span>选择星期</span></div></template>
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
            <template #label><div class="flex items-center gap-1"><el-icon><Clock /></el-icon><span>开始日期</span></div></template>
            <el-date-picker v-model="form.start_date" type="date" />
          </el-form-item>
          <el-form-item prop="end_date">
            <template #label><div class="flex items-center gap-1"><el-icon><Clock /></el-icon><span>截止日期</span></div></template>
            <el-date-picker v-model="form.end_date" type="date" />
          </el-form-item>
        </el-form>
      </el-card>
    </div>

    <!-- 底部操作区 -->
    <div class="flex justify-end gap-2">
      <el-button @click="goBack">取消</el-button>
      <el-button type="primary" @click="submitForm">保存</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
// 中文注释：编辑任务页面逻辑，根据路由参数加载任务并更新保存
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Edit, List, Clock, Coin, Plus } from '@element-plus/icons-vue'
import router from '@/router'
import TaskImageUploader from '@/components/TaskImageUploader.vue'
import { getTask, updateTask, uploadTaskImage } from '@/services/tasks'
import { useRoute } from 'vue-router'
import { prepareUpload } from '@/utils/image'

const route = useRoute()
const taskId = Number(route.params.id)
const userId = 1 // 中文注释：示例用户ID
const taskLoaded = ref(false)
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

onMounted(async () => {
  try {
    const t = await getTask(taskId)
    // 中文注释：填充表单
    form.name = t.name
    form.description = t.description
    form.category = t.category
    form.score = t.score
    form.plan_minutes = t.plan_minutes
    form.start_date = new Date(t.start_date)
    form.end_date = t.end_date ? new Date(t.end_date) : undefined
    form.images = t.image_json ? (JSON.parse(t.image_json) as string[]) : []
    taskLoaded.value = true
  } catch (e: any) {
    ElMessage.error(e.message || '加载任务失败')
  }
})

async function submitForm() {
  try {
    const payload = {
      name: form.name,
      description: form.description,
      category: form.category,
      score: form.score,
      plan_minutes: form.plan_minutes,
      start_date: form.start_date,
      end_date: form.end_date
    }
    const t = await updateTask(taskId, payload)
    // 中文注释：如有新增本地图片，上传后将返回的相对路径写入 image_json，保证编辑页也能持久化图片列表
    if ((form.local_images || []).length > 0) {
      const newPaths: string[] = []
      for (const f of form.local_images) {
        try {
          const webp = await prepareUpload(f as File)
          const { path } = await uploadTaskImage(userId, webp, t.id)
          newPaths.push(path)
        } catch (err: any) {
          console.error('上传任务图片失败', { task_id: t.id, filename: (f as File)?.name, message: err?.message || err })
        }
      }
      // 将新的图片路径合并到当前表单状态，并同步写入后端任务的 image_json 字段
      if (newPaths.length > 0) {
        form.images = [...(form.images || []), ...newPaths]
        try {
          await updateTask(t.id, { image_json: JSON.stringify(form.images) })
        } catch (err: any) {
          console.error('更新任务图片列表失败', { task_id: t.id, message: err?.message || err, response: err?.response?.data })
        }
      }
    }
    ElMessage.success('保存成功')
    router.back()
  } catch (e: any) {
    ElMessage.error(e.message || '保存失败')
  }
}
</script>

<style scoped>
/* 中文注释：使用 Tailwind 进行响应式布局 */
</style>