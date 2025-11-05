<template>
  <!-- 中文注释：任务页面，包含统计、列表、创建/编辑、批量删除、番茄钟功能 -->
  <div class="p-4 space-y-4">
    <!-- 顶部统计栏 -->
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <!-- 中文注释：默认头像改为本地 assets/avatars/default.png，与“我的”页面保持一致 -->
          <el-avatar :size="36" :src="defaultAvatar" />
          <div class="font-semibold">今日统计</div>
        </div>
      <!-- 中文注释：右侧取消回收站，改为彩色统计图标与金币显示 -->
      <div class="flex items-center gap-3">
        <div class="flex items-center gap-1">
          <el-icon :size="20" style="color:#f59e0b"><Coin /></el-icon>
          <!-- 中文注释：总金币 = 已完成任务奖励金币总和 - 心愿兑换扣除金币 -->
          <span class="font-semibold">{{ totalCoins }}</span>
        </div>
        <el-icon :size="24" style="color:#ec4899"><DataAnalysis /></el-icon>
      </div>
    </div>

    <!-- 顶部统计：四项一行，不同颜色图标；下方单独大“统计”卡片居中显示 -->
    <div class="grid grid-cols-4 gap-2">
      <el-card shadow="never" class="stat-card">
        <div class="flex flex-col items-center">
          <el-icon style="color:#22c55e"><Clock /></el-icon>
          <div class="text-xs text-gray-500">日时长</div>
          <div class="font-semibold">{{ dayMinutes }} 分钟</div>
        </div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="flex flex-col items-center">
          <el-icon style="color:#3b82f6"><List /></el-icon>
          <div class="text-xs text-gray-500">任务数</div>
          <div class="font-semibold">{{ completedTasksCount }}/{{ tasks.length }}</div>
        </div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="flex flex-col items-center">
          <el-icon style="color:#f59e0b"><Coin /></el-icon>
          <div class="text-xs text-gray-500">日金币</div>
          <div class="font-semibold">{{ dayCoins }}</div>
        </div>
      </el-card>
      <el-card shadow="never" class="stat-card">
        <div class="flex flex-col items-center">
          <el-icon style="color:#14b8a6"><CircleCheck /></el-icon>
          <div class="text-xs text-gray-500">完成率</div>
          <div class="font-semibold">{{ completeRate }}%</div>
        </div>
      </el-card>
    </div>
    <!-- 已将大“统计”图标移动到顶部右侧，此处删除 -->

    <!-- 日期选择与周视图：移动到“今日任务”卡片上方 -->
    <div class="my-2">
      <WeekCalendar v-model:selectedDate="selectedDate" :task-count-map="taskCountMap" />
    </div>

    <!-- 任务列表卡片 -->
    <el-card shadow="never">
      <template #header>
        <div class="flex items-center justify-between">
          <span class="font-semibold">今日任务</span>
          <div class="space-x-2 flex items-center">
            <!-- 中文注释：筛选图标下拉菜单，点击选择“全部/已完成/待完成” -->
            <el-dropdown trigger="click" @command="onFilterCommand">
              <span class="el-dropdown-link">
                <el-icon class="cursor-pointer" :size="18" title="筛选"><Filter /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="全部">全部</el-dropdown-item>
                  <el-dropdown-item command="已完成">已完成</el-dropdown-item>
                  <el-dropdown-item command="待完成">待完成</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button type="primary" @click="openCreate">添加任务</el-button>
          </div>
        </div>
      </template>
      <!-- 日期选择与周视图已移动到卡片上方 -->
      <!-- 分类筛选：全部学科/语文/数学/英语 -->
      <div class="my-2">
        <el-radio-group v-model="categoryFilter" size="small">
          <el-radio-button label="全部任务" />
          <el-radio-button label="语文" />
          <el-radio-button label="数学" />
          <el-radio-button label="英语" />
        </el-radio-group>
      </div>

      <el-empty v-if="filteredTasks.length===0" description="暂无任务，快去添加吧~" />
      <div v-else class="space-y-4">
        <!-- 按分类分组显示 -->
        <div v-for="group in groupedTasks" :key="group.category" class="space-y-3">
          <div class="text-base font-semibold text-green-700">{{ group.category }}</div>
          <el-card v-for="t in group.items" :key="t.id" shadow="hover" class="relative">
            <!-- 中文注释：自定义圆形复选框，居中于第一行与第二行之间，略大，点击切换完成状态 -->
            <div class="absolute left-2 top-1/2 -translate-y-1/2">
              <div
                class="w-5 h-5 rounded-full border-2 flex items-center justify-center cursor-pointer"
                :class="t.status===2 ? 'bg-green-500 border-green-500 text-white' : 'border-gray-400 text-gray-400'"
                @click="() => onCheckComplete(t, t.status !== 2)"
                title="点击切换完成状态"
              >
                <el-icon :size="12">
                  <CircleCheck />
                </el-icon>
              </div>
            </div>
            <!-- 第一行：左侧任务名，右侧状态与番茄钟入口 + 菜单 -->
            <div class="flex items-center justify-between pl-10">
              <div class="flex items-center gap-3">
                <!-- 中文注释：番茄钟图标仅在未完成时显示，位于右侧“待完成”标签左侧，此处移除 -->
                <div class="font-semibold text-left" :class="{'text-gray-500': t.status === 2}">{{ t.name }}</div>
              </div>
              <div class="flex items-center gap-2">
                <template v-if="t.status !== 2">
                  <img src="@/assets/tomato.png" alt="番茄钟" class="w-4 h-4 cursor-pointer" @click="openTomato(t)" />
                  <el-tag type="danger" size="small">待完成</el-tag>
                </template>
                <template v-else>
                  <el-tag type="success" size="small">已完成</el-tag>
                </template>
                <el-dropdown trigger="click" @command="(cmd)=>onMenu(cmd, t)">
                  <span class="el-dropdown-link">
                    <el-icon class="cursor-pointer"><MoreFilled /></el-icon>
                  </span>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="tomato">
                        <el-icon class="mr-1"><Clock /></el-icon>番茄钟
                      </el-dropdown-item>
                      <el-dropdown-item command="edit">
                        <el-icon class="mr-1"><Edit /></el-icon>编辑
                      </el-dropdown-item>
                      <el-dropdown-item command="delete" style="color:#f56c6c">
                        <el-icon class="mr-1"><Delete /></el-icon>删除
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>

            <!-- 第二行：左侧备注/描述；右侧实际/计划/金币（实际精确到秒） -->
            <div class="flex items-center justify-between mt-1 pl-10">
              <div class="text-xs text-gray-500 truncate max-w-[60%] text-left">{{ t.remark || t.description }}</div>
              <div class="flex items-center gap-3 text-xs">
                <template v-if="t.status===2">
                  <div class="flex items-center gap-1 text-blue-600 text-xs" title="实际完成时间">
                    <el-icon :size="14"><Clock /></el-icon>
                    <span class="font-semibold">{{ formatHMS(actualSecondsLocal[t.id] ?? ((t.actual_minutes||0)*60)) }}</span>
                  </div>
                </template>
                <div class="flex items-center gap-1 text-green-600 text-xs" title="计划用时">
                  <el-icon :size="14"><List /></el-icon>
                  <span class="font-semibold">{{ t.plan_minutes || 0 }} 分</span>
                </div>
                <div class="flex items-center gap-1 text-amber-600 text-xs" title="金币">
                  <el-icon :size="14"><Coin /></el-icon>
                  <span class="font-semibold">{{ t.score || 0 }}</span>
                </div>
              </div>
            </div>
          </el-card>
        </div>
      </div>
    </el-card>

    <!-- 创建/编辑对话框 -->
    <el-dialog v-model="formVisible" :title="editing ? '编辑任务' : '创建任务'" width="520px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="88px">
        <el-form-item label="任务标题" prop="name">
          <el-input v-model="form.name" maxlength="128" show-word-limit />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="form.category" placeholder="选择分类">
            <el-option label="语文" value="语文" />
            <el-option label="数学" value="数学" />
            <el-option label="英语" value="英语" />
          </el-select>
        </el-form-item>
        <!-- 中文注释：字段文案统一为“金币” -->
        <el-form-item label="金币" prop="score">
          <el-input-number v-model="form.score" :min="-10" :max="10" />
        </el-form-item>
        <el-form-item label="计划时长" prop="plan_minutes">
          <!-- 中文注释：计划时长最小 1 分钟，满足用户需求 -->
          <el-input-number v-model="form.plan_minutes" :min="1" :max="240" />
        </el-form-item>
        <el-form-item label="开始日期" prop="start_date">
          <el-date-picker v-model="form.start_date" type="date" />
        </el-form-item>
        <el-form-item label="截止日期" prop="end_date">
          <el-date-picker v-model="form.end_date" type="date" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="formVisible=false">取消</el-button>
          <el-button type="primary" @click="submitForm">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 回收站对话框 -->
    <el-dialog v-model="recycleVisible" title="回收站" width="600px">
      <el-table :data="recycleList" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="任务" />
        <el-table-column prop="category" label="分类" width="100" />
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="restore([row.id])">恢复</el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="recycleVisible=false">关闭</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 番茄钟组件弹窗 -->
    <!-- 中文注释：番茄钟弹窗标题：🍅 番茄钟：任务名称 -->
    <el-dialog v-model="tomatoVisible" :title="`🍅 番茄钟：${currentTask?.name || '番茄钟'}`" width="520px" :show-close="true" :close-on-click-modal="!store.tomato.fixedTomatoPage">
      <!-- 中文注释：番茄钟默认使用任务计划时长，若无则 20 分钟；支持倒计时/正计时、预设与自定义 -->
      <TomatoTimer :work-minutes="currentTask?.plan_minutes || 20" :break-minutes="5" :task-name="currentTask?.name" :task-remark="currentTask?.remark || currentTask?.description" @complete="onTomatoComplete" />
    </el-dialog>

    <!-- 右下角绿色加号浮动按钮：创建任务 -->
  <el-button
      type="success"
      circle
      class="fixed bottom-20 right-6 shadow-lg"
      @click="openCreate"
      title="创建任务"
    >
      <!-- 中文注释：创建任务图标放大 0.5 倍（默认约 16px → 24px） -->
      <el-icon :size="24"><Plus /></el-icon>
    </el-button>

    <!-- 悬浮番茄钟：位于底部导航上方居中，水位填充动画展示进度 -->
    <div
      v-if="store.tomato.running && !store.tomato.fixedTomatoPage"
      class="fixed bottom-16 left-1/2 -translate-x-1/2 z-50"
    >
      <div class="w-14 h-14 rounded-full shadow-lg cursor-pointer overflow-hidden relative bg-blue-200 bg-opacity-60"
           @click="tomatoVisible=true">
        <div class="absolute bottom-0 left-0 right-0 bg-blue-500 bg-opacity-70 transition-all"
             :style="{ height: fillPercent + '%' }"></div>
        <!-- 中文注释：悬浮球中央显示 mm:ss，小号文字 -->
        <div class="absolute inset-0 flex items-center justify-center text-white font-semibold text-[10px]">{{ floatingTime }}</div>
      </div>
    </div>
  </div>
 </template>

<script setup lang="ts">
// 中文注释：任务页面逻辑，统一使用服务层 API，实现表单校验与错误提示
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Plus, Clock, List, Coin, CircleCheck, MoreFilled, DataAnalysis, Edit, Delete, Filter } from '@element-plus/icons-vue'
import defaultAvatar from '@/assets/avatars/default.png'
import { useAppState } from '@/stores/appState'
import TomatoTimer from '@/components/TomatoTimer.vue'
import WeekCalendar from '@/components/WeekCalendar.vue'
import dayjs from 'dayjs'
import utc from 'dayjs/plugin/utc'
dayjs.extend(utc)
import { listTasks, createTask, updateTask, updateTaskStatus, deleteTask, completeTomato, listRecycleBin, restoreTasks, type TaskItem } from '@/services/tasks'

// 顶部统计占位（后续与后端联动）
const store = useAppState()
// 中文注释：总金币改为直接读取全局 store.coins（由后端任务完成/取消与心愿兑换实时更新），与心愿页保持一致
const totalCoins = computed(() => store.coins)
const completedTasksCount = computed(() => {
  return tasks.value.filter(t => t.status === 2).length
})
const dayCoins = ref(0)
const dayMinutes = ref(0)
const completeRate = ref(0)

// 列表与筛选
const tasks = ref<TaskItem[]>([])
const filter = ref<'全部' | '已完成' | '待完成'>('全部')
const categoryFilter = ref<'全部任务' | '语文' | '数学' | '英语'>('全部任务')
const selectedDate = ref<string>(dayjs().format('YYYY-MM-DD'))
const taskCountMap = computed<Record<string, number>>(() => {
  const map: Record<string, number> = {}
  for (const t of tasks.value) {
    const d = t.start_date ? dayjs.utc(t.start_date).local().format('YYYY-MM-DD') : ''
    if (!d) continue
    map[d] = (map[d] || 0) + 1
  }
  return map
})
const filteredTasks = computed(() => {
  let result = tasks.value
  if (filter.value === '已完成') result = result.filter((t) => t.status === 2)
  else if (filter.value === '待完成') result = result.filter((t) => t.status !== 2)
  if (categoryFilter.value !== '全部任务') result = result.filter((t) => (t.category || '') === categoryFilter.value)
  // 中文注释：日期过滤，选中日期在任务的开始/截止范围内（无截止则按开始日期）
  result = result.filter((t) => {
    const sd = t.start_date ? dayjs(t.start_date) : null
    const ed = t.end_date ? dayjs(t.end_date) : null
    const sel = dayjs(selectedDate.value)
    if (!sd) return false
    if (ed) return sel.isSame(sd, 'day') || (sel.isAfter(sd, 'day') && sel.isBefore(ed, 'day')) || sel.isSame(ed, 'day')
    return sel.isSame(sd, 'day')
  })
  return result
})

// 中文注释：按分类分组，便于移动端展示与筛选
const groupedTasks = computed(() => {
  const map = new Map<string, TaskItem[]>()
  for (const t of filteredTasks.value) {
    const cat = t.category || '未分类'
    if (!map.has(cat)) map.set(cat, [])
    map.get(cat)!.push(t)
  }
  return Array.from(map.entries()).map(([category, items]) => ({ category, items }))
})

// 中文注释：本地保存每个任务的实际秒数，便于前端精确展示（后端以分钟存储）
const actualSecondsLocal = reactive<Record<number, number>>({})
// 中文注释：格式化为“x时x分x秒”，按需显示（仅显示非零单位；全为零则显示0秒）
function formatHMS(sec: number) {
  const h = Math.floor(sec / 3600)
  const m = Math.floor((sec % 3600) / 60)
  const s = Math.floor(sec % 60)
  const parts: string[] = []
  if (h > 0) parts.push(`${h}时`)
  if (m > 0) parts.push(`${m}分`)
  if (s > 0 || parts.length === 0) parts.push(`${s}秒`)
  return parts.join('')
}

// 选择与批量删除：移动端不再提供批量删除，这里移除选择逻辑

// 创建/编辑表单
const formVisible = ref(false)
const recycleVisible = ref(false)
const tomatoVisible = ref(false)
const editing = ref(false)
const currentTask = ref<TaskItem | null>(null)
const formRef = ref<FormInstance>()
const form = reactive<any>({ name: '', description: '', category: '语文', score: 1, plan_minutes: 20, start_date: new Date(), end_date: undefined })
const rules: FormRules = {
  name: [{ required: true, message: '请输入任务标题', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  score: [{ required: true, message: '请输入金币', trigger: 'change' }],
  plan_minutes: [{ required: true, message: '请输入计划时长', trigger: 'change' }],
  start_date: [{ required: true, message: '请选择开始日期', trigger: 'change' }]
}

// 中文注释：移除未使用的状态文案/样式函数，简化页面逻辑

async function fetchTasks() {
  try {
    const res = await listTasks()
    tasks.value = res.items || []
    // 简单统计
    dayMinutes.value = tasks.value.reduce((sum, t) => sum + (t.actual_minutes || 0), 0)
    dayCoins.value = tasks.value.filter((t) => t.status === 2).reduce((sum, t) => sum + (t.score || 0), 0)
    completeRate.value = tasks.value.length ? Math.round((tasks.value.filter((t) => t.status === 2).length / tasks.value.length) * 100) : 0
    // 中文注释：同步更新全局 coins（考虑心愿扣减），心愿页面读取该值作为可用金币
    store.setCoins(totalCoins.value)
  } catch (e: any) {
    ElMessage.error(`任务列表加载失败：${e.message || e}`)
  }
}

function openCreate() {
  editing.value = false
  Object.assign(form, { name: '', description: '', category: '语文', score: 1, plan_minutes: 20, start_date: new Date(), end_date: undefined })
  formVisible.value = true
}

function openEdit(t: TaskItem) {
  editing.value = true
  Object.assign(form, { name: t.name, description: t.description, category: t.category, score: t.score, plan_minutes: t.plan_minutes, start_date: new Date(t.start_date), end_date: t.end_date ? new Date(t.end_date) : undefined })
  currentTask.value = t
  formVisible.value = true
}

async function submitForm() {
  try {
    await formRef.value?.validate()
  } catch { return }
  try {
    if (editing.value && currentTask.value) {
      await updateTask(currentTask.value.id, {
        name: form.name,
        description: form.description,
        category: form.category,
        score: form.score,
        plan_minutes: form.plan_minutes,
        start_date: form.start_date,
        end_date: form.end_date
      })
      ElMessage.success('任务已更新')
    } else {
      await createTask({
        user_id: 1, // 中文注释：演示用，后续接入登录用户 ID
        name: form.name,
        description: form.description,
        category: form.category,
        score: form.score,
        plan_minutes: form.plan_minutes,
        start_date: form.start_date,
        end_date: form.end_date
      })
      ElMessage.success('任务已创建')
    }
    formVisible.value = false
    await fetchTasks()
  } catch (e: any) {
    ElMessage.error(`提交失败：${e.message || e}`)
  }
}

// 勾选即完成：只允许从未完成 -> 已完成，不提供取消
async function onCheckComplete(t: TaskItem, checked: boolean) {
  try {
    if (checked) {
      // 中文注释：勾选为完成：按计划时长计入实际，并标记为已完成
      const planM = t.plan_minutes || 20
      await updateTask(t.id, { actual_minutes: planM })
      const resp: any = await updateTaskStatus(t.id, 2)
      t.status = 2
      t.actual_minutes = (t.actual_minutes || 0) + planM
      actualSecondsLocal[t.id] = planM * 60
      if (resp && typeof resp.user_coins !== 'undefined') store.setCoins(Number(resp.user_coins))
      ElMessage.success('已标记为完成（按计划时长计）')
    } else {
      // 中文注释：取消完成：标记为未完成，并从日金币与总金币中扣除该任务金币
      const resp: any = await updateTaskStatus(t.id, 0)
      t.status = 0
      if (resp && typeof resp.user_coins !== 'undefined') store.setCoins(Number(resp.user_coins))
      ElMessage.success('已取消完成，金币已扣除')
    }
    // 统一刷新统计
    dayMinutes.value = tasks.value.reduce((sum, x) => sum + (x.actual_minutes || 0), 0)
    dayCoins.value = tasks.value.filter((x) => x.status === 2).reduce((sum, x) => sum + (x.score || 0), 0)
    completeRate.value = tasks.value.length ? Math.round((tasks.value.filter((x) => x.status === 2).length / tasks.value.length) * 100) : 0
  } catch (e: any) {
    ElMessage.error(`状态变更失败：${e.message || e}`)
  }
}
// 取消切换状态功能：保留空函数避免引用错误（模板已移除）

function confirmDelete(t: TaskItem) {
  ElMessageBox.confirm(`确认删除任务「${t.name}」？删除后可在回收站恢复。`, '提示', { type: 'warning' })
    .then(async () => {
      try {
        await deleteTask(t.id)
        ElMessage.success('已删除，可在回收站恢复')
        await fetchTasks()
      } catch (e: any) {
        ElMessage.error(`删除失败：${e.message || e}`)
      }
    })
    .catch(() => {})
}



const recycleList = ref<TaskItem[]>([])
async function restore(ids: number[]) {
  try {
    await restoreTasks(ids)
    ElMessage.success('已恢复')
    recycleList.value = await listRecycleBin()
    await fetchTasks()
  } catch (e: any) {
    ElMessage.error(`恢复失败：${e.message || e}`)
  }
}

function openTomato(t: TaskItem) {
  currentTask.value = t
  // 中文注释：点击番茄钟按钮默认倒计时模式，时间为任务设定时间
  const m = t.plan_minutes || 20
  store.updateTomato({ mode: 'countdown', durationMinutes: m, remainingSeconds: m * 60 })
  tomatoVisible.value = true
}

async function onTomatoComplete(seconds?: number) {
  if (!currentTask.value) return
  try {
    // 中文注释：按实际秒数精确展示，后端按分钟上报（四舍五入）；无秒数则按计划时长
    const usedSec = Math.max(1, seconds || (currentTask.value.plan_minutes || 20) * 60)
    const reportMinutes = Math.max(1, Math.round(usedSec / 60))
    const updated = await completeTomato(currentTask.value.id, reportMinutes)
    const idx = tasks.value.findIndex((x) => x.id === currentTask.value!.id)
    if (idx >= 0) tasks.value[idx] = updated
    actualSecondsLocal[currentTask.value.id] = usedSec
    // 完成后标记任务为已完成
    await updateTaskStatus(currentTask.value.id, 2)
    if (idx >= 0) tasks.value[idx].status = 2
    dayMinutes.value = tasks.value.reduce((sum, x) => sum + (x.actual_minutes || 0), 0)
    ElMessage.success('番茄钟完成，数据已记录')
    tomatoVisible.value = false
  } catch (e: any) {
    ElMessage.error(`番茄上报失败：${e.message || e}`)
  }
}

onMounted(() => {
  fetchTasks()
})

// 菜单命令统一处理
function onMenu(cmd: string, t: TaskItem) {
  if (cmd === 'tomato') return openTomato(t)
  if (cmd === 'edit') return openEdit(t)
  if (cmd === 'delete') return confirmDelete(t)
}

// 中文注释：筛选图标下拉菜单命令处理，更新状态筛选条件
function onFilterCommand(cmd: '全部' | '已完成' | '待完成') {
  filter.value = cmd
}

// 中文注释：悬浮球填充百分比（正计时用已用时 / 目标时长，倒计时用剩余时间）
const fillPercent = computed(() => {
  const dur = store.tomato.durationMinutes * 60
  const sec = store.tomato.remainingSeconds
  if (store.tomato.mode === 'countup') {
    return Math.min(100, Math.round((sec / dur) * 100))
  }
  return Math.min(100, Math.round(((dur - sec) / dur) * 100))
})

// 中文注释：悬浮球显示的时间文本（mm:ss），倒计时显示剩余，正计时显示累计
const floatingTime = computed(() => {
  const sec = Math.max(0, store.tomato.remainingSeconds || 0)
  const mm = String(Math.floor(sec / 60)).padStart(2, '0')
  const ss = String(sec % 60).padStart(2, '0')
  return `${mm}:${ss}`
})
// 中文注释：移除未使用的函数（toggleStatus、openRecycle），消除编译器警告
</script>

<style scoped>
/* 中文注释：基本页面样式，响应式栅格布局已通过 Tailwind 实现 */
/* 中文注释：分类筛选按钮的文字加粗，增强可读性 */
.el-radio-button__inner { font-weight: 600; }
</style>
