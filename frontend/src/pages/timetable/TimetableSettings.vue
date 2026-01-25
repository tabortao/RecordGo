<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950 pb-20 relative overflow-hidden">
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <div class="absolute -top-24 -right-24 h-72 w-72 rounded-full bg-sky-300/35 dark:bg-sky-500/16 blur-3xl" />
      <div class="absolute -bottom-40 -left-28 h-80 w-80 rounded-full bg-amber-200/35 dark:bg-amber-500/14 blur-3xl" />
      <div class="absolute inset-0 bg-[radial-gradient(1200px_circle_at_20%_-20%,rgba(255,255,255,.65),transparent_55%),radial-gradient(900px_circle_at_80%_0%,rgba(255,255,255,.45),transparent_55%)] dark:bg-[radial-gradient(1200px_circle_at_20%_-20%,rgba(255,255,255,.07),transparent_55%),radial-gradient(900px_circle_at_80%_0%,rgba(255,255,255,.06),transparent_55%)]" />
    </div>

    <div class="sticky top-0 z-20 px-4 pt-4">
      <div class="rounded-3xl border border-white/50 dark:border-gray-800/60 bg-white/75 dark:bg-gray-900/70 backdrop-blur-xl shadow-sm px-3 py-3 flex items-center gap-3">
        <button
          type="button"
          class="inline-flex h-10 w-10 items-center justify-center rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/70 dark:bg-gray-900/60 text-gray-600 dark:text-gray-300 hover:bg-white hover:text-gray-900 dark:hover:text-gray-50 transition-colors"
          @click="goBack"
        >
          <el-icon :size="20"><ArrowLeft /></el-icon>
        </button>
        <div class="min-w-0">
          <h1 class="font-extrabold tracking-tight text-gray-900 dark:text-gray-50 text-base">课表设置</h1>
          <div class="mt-0.5 text-xs text-gray-500 dark:text-gray-400">配置年级学期、课程库与课表编排</div>
        </div>
      </div>
    </div>

    <div class="relative z-10 p-4 space-y-6">
        <!-- 基础配置 -->
        <el-card shadow="never" class="pretty-card">
            <template #header>
                <div class="font-bold">基础配置</div>
            </template>
            <el-form label-position="left" label-width="100px">
                <el-form-item label="年级">
                    <el-select v-model="form.current_grade" @change="onConfigChange">
                        <el-option v-for="g in grades" :key="g" :label="g" :value="g" />
                    </el-select>
                </el-form-item>
                <el-form-item label="学期">
                    <el-radio-group v-model="form.current_semester" @change="onConfigChange">
                        <el-radio-button label="上学期" />
                        <el-radio-button label="下学期" />
                    </el-radio-group>
                </el-form-item>
                <el-form-item label="周六">
                    <el-switch v-model="form.show_saturday" @change="onConfigChange" />
                </el-form-item>
                <el-form-item label="周日">
                    <el-switch v-model="form.show_sunday" @change="onConfigChange" />
                </el-form-item>
            </el-form>
        </el-card>

        <!-- 课程库管理 -->
        <el-card shadow="never" class="pretty-card">
            <template #header>
                <div class="flex justify-between items-center">
                    <div class="font-bold">课程库</div>
                    <el-button type="primary" link @click="openAddCourseDialog">添加课程</el-button>
                </div>
            </template>
            
            <!-- 配色方案 -->
            <div class="mb-4">
                <div class="text-sm text-gray-500 mb-2">一键配色</div>
                <div class="flex flex-wrap gap-2">
                    <div v-for="palette in palettes" :key="palette.name" 
                         class="cursor-pointer border border-gray-200 dark:border-gray-700 rounded-lg p-1 hover:border-primary transition-colors flex flex-col gap-1 items-center"
                         @click="applyPalette(palette)">
                        <div class="flex gap-0.5">
                            <div v-for="c in palette.colors" :key="c" class="w-3 h-3 rounded-full" :style="{ backgroundColor: c }"></div>
                        </div>
                        <span class="text-[10px] text-gray-500">{{ palette.name }}</span>
                    </div>
                </div>
            </div>

            <div class="flex flex-wrap gap-2">
                <div v-for="course in courses" :key="course.id" 
                     class="group relative px-3 py-1 rounded-full text-sm font-medium shadow-sm flex items-center gap-1 cursor-default text-gray-800"
                     :style="{ backgroundColor: getCourseColor(course), color: isLightColor(getCourseColor(course)) ? '#000' : '#fff' }">
                     {{ course.name }}
                     
                     <!-- 自定义课程操作按钮 -->
                     <div v-if="course.user_id" class="hidden group-hover:flex items-center gap-1 ml-1 pl-1 border-l border-white/30">
                        <el-icon class="cursor-pointer hover:scale-110" @click.stop="openEditCourseDialog(course)"><Edit /></el-icon>
                        <el-icon class="cursor-pointer hover:scale-110" @click.stop="handleDeleteCourse(course)"><Delete /></el-icon>
                     </div>
                </div>
            </div>
        </el-card>

        <!-- 课表编排 -->
        <el-card shadow="never" class="pretty-card overflow-visible">
            <template #header>
                <div class="flex justify-between items-center">
                    <div class="font-bold">课表编排 ({{ form.current_grade }} {{ form.current_semester }})</div>
                    <div class="text-xs text-gray-400">修改自动保存</div>
                </div>
            </template>
            
            <div class="overflow-x-auto -mx-4 px-4">
                 <div class="min-w-[320px]">
                    <div class="grid gap-1 mb-1" :style="{ gridTemplateColumns: gridColumns }">
                        <div class="w-16"></div>
                        <div v-for="day in days" :key="day.value" class="text-center py-1 text-xs font-semibold text-gray-500">
                            {{ day.label }}
                        </div>
                    </div>
                    <div v-for="period in periods" :key="period" class="grid gap-1 mb-1" :style="{ gridTemplateColumns: gridColumns }">
                        <div class="flex flex-col items-center justify-center text-xs text-gray-400 w-16 border-r pr-1 cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700 rounded" @click="editPeriodTime(period)">
                            <span class="font-bold">{{ period }}</span>
                            <span class="text-[10px] scale-90 text-gray-300" v-if="getPeriodTime(period)">{{ getPeriodTime(period) }}</span>
                            <span class="text-[10px] scale-90 text-blue-400" v-else>设置时间</span>
                        </div>
                        <div 
                            v-for="day in days" 
                            :key="`${day.value}-${period}`"
                            class="h-12 rounded border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 flex items-center justify-center text-xs font-bold cursor-pointer hover:border-primary transition-colors"
                            :style="getCellStyle(day.value, period)"
                            @click="openCourseSelector(day.value, period)"
                        >
                             {{ getCourseName(day.value, period) }}
                        </div>
                    </div>
                 </div>
            </div>
        </el-card>
    </div>

    <!-- 课程选择弹窗 -->
    <el-dialog v-model="selectorVisible" title="选择课程" width="90%" class="pretty-dialog" append-to-body>
        <div class="grid grid-cols-4 gap-3">
            <div v-for="course in courses" :key="course.id"
                 class="h-10 rounded-lg flex items-center justify-center text-sm font-medium shadow-sm cursor-pointer active:scale-95 transition-transform"
                 :style="{ backgroundColor: getCourseColor(course), color: isLightColor(getCourseColor(course)) ? '#000' : '#fff' }"
                 @click="selectCourse(course)"
            >
                {{ course.name }}
            </div>
            <div class="h-10 rounded-lg border border-dashed border-gray-300 flex items-center justify-center text-gray-400 text-sm cursor-pointer" @click="clearCell">
                空
            </div>
        </div>
    </el-dialog>

    <!-- 添加/编辑课程弹窗 -->
    <el-dialog v-model="showCourseDialog" :title="editingCourseId ? '编辑课程' : '添加新课程'" width="80%" class="pretty-dialog" append-to-body>
        <el-form>
            <el-form-item label="名称">
                <el-input v-model="courseForm.name" placeholder="例如: 编程" />
            </el-form-item>
            <el-form-item label="颜色">
                <el-color-picker v-model="courseForm.color" />
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="showCourseDialog = false">取消</el-button>
            <el-button type="primary" @click="saveCourse" :loading="savingCourse">确定</el-button>
        </template>
    </el-dialog>

    <!-- 课程时间设置弹窗 -->
    <el-dialog v-model="timeDialogVisible" title="设置课程时间" width="80%" class="pretty-dialog" append-to-body>
        <div v-if="editingPeriod" class="space-y-4">
            <div class="text-center font-bold text-lg mb-4">第 {{ editingPeriod.period }} 节</div>
            <div class="flex items-center gap-2 justify-center">
                <el-time-select
                    v-model="editingPeriod.start_time"
                    start="06:00"
                    step="00:05"
                    end="22:00"
                    placeholder="开始"
                    class="w-32"
                />
                <span>-</span>
                <el-time-select
                    v-model="editingPeriod.end_time"
                    start="06:00"
                    step="00:05"
                    end="22:00"
                    placeholder="结束"
                    class="w-32"
                />
            </div>
        </div>
        <template #footer>
            <el-button @click="timeDialogVisible = false">取消</el-button>
            <el-button type="primary" @click="savePeriodTime">确定</el-button>
        </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft, Edit, Delete } from '@element-plus/icons-vue'
import { useTimetableStore } from '@/stores/timetable'
import { storeToRefs } from 'pinia'
import { ElMessage, ElMessageBox } from 'element-plus'
import { timetableApi, type Course, type TimetableItem, type PeriodSetting } from '@/api/timetable'

const router = useRouter()
const store = useTimetableStore()
const { config, courses, timetable: originalTimetable } = storeToRefs(store)

const form = ref({
    current_grade: '一年级',
    current_semester: '上学期',
    show_saturday: false,
    show_sunday: false,
    period_settings_json: '',
    course_colors_json: '',
    background_emojis: ''
})

const palettes = [
    { name: '晨曦微光', colors: ['#FAD6A5', '#E8C3A2', '#9EC9D0', '#7BA7BC', '#4A707A'] },
    { name: '森林学者', colors: ['#E9E5D6', '#A3B18A', '#588157', '#3A5A40', '#344E41'] },
    { name: '中性都市', colors: ['#F5F5F5', '#CFCFCF', '#9E9E9E', '#616161', '#3D3D3D'] },
    { name: '克莱因幻想', colors: ['#FFFFFF', '#F0F0F0', '#0D4C92', '#00B4D8', '#0077B6'] },
    { name: '莫兰迪课堂', colors: ['#F7F1E5', '#E7D8C9', '#B8A897', '#A3B4C4', '#7C8B9C'] },
    { name: '大地探险家', colors: ['#EDE6D6', '#D4A574', '#B05B3B', '#6A994E', '#386641'] },
    { name: '星际穿梭', colors: ['#0F0F1C', '#2D3047', '#419D78', '#E0E0E0', '#FF9F1C'] },
    { name: '新中式雅韵', colors: ['#F8F4E9', '#DED0B6', '#BBAB8C', '#7A918D', '#556B5E'] },
    { name: '纯净冰山', colors: ['#FFFFFF', '#E3F2FD', '#90CAF9', '#42A5F5', '#1E88E5'] },
    { name: '低饱和对比', colors: ['#FFF8F0', '#FFD6C2', '#C4E4FF', '#A2D2A1', '#D9B8C4'] }
]

const courseColors = ref<Record<string, string>>({})

function getCourseColor(course: Course) {
    return courseColors.value[course.name] || course.color
}

// 简单的判断颜色深浅
function isLightColor(color: string) {
    if (!color) return true
    const hex = color.replace('#', '')
    const r = parseInt(hex.substring(0, 2), 16)
    const g = parseInt(hex.substring(2, 4), 16)
    const b = parseInt(hex.substring(4, 6), 16)
    const brightness = (r * 299 + g * 587 + b * 114) / 1000
    return brightness > 155 // 阈值可调
}

function applyPalette(palette: { name: string, colors: string[] }) {
    const colors = palette.colors
    const newColors: Record<string, string> = {}
    
    // 给当前所有课程分配颜色
    courses.value.forEach((course, index) => {
        newColors[course.name] = colors[index % colors.length]
    })
    
    courseColors.value = newColors
    form.value.course_colors_json = JSON.stringify(newColors)
    onConfigChange()
    ElMessage.success(`已应用“${palette.name}”配色`)
}

const grades = ['一年级', '二年级', '三年级', '四年级', '五年级', '六年级', '初一', '初二', '初三', '高一', '高二', '高三']
const periods = Array.from({ length: 10 }, (_, i) => i + 1)

// 编辑中的课表数据
const editTimetable = ref<Partial<TimetableItem>[]>([])
const saving = ref(false)

// 课程选择器状态
const selectorVisible = ref(false)
const currentCell = ref<{day: number, period: number} | null>(null)

// 课程管理状态
const showCourseDialog = ref(false)
const editingCourseId = ref<number | null>(null)
const courseForm = ref({ name: '', color: '#409EFF' })
const savingCourse = ref(false)

// 时间设置状态
const timeDialogVisible = ref(false)
const editingPeriod = ref<PeriodSetting | null>(null)
const periodSettings = ref<PeriodSetting[]>([])

const days = computed(() => {
  const allDays = [
    { label: '周一', value: 1 },
    { label: '周二', value: 2 },
    { label: '周三', value: 3 },
    { label: '周四', value: 4 },
    { label: '周五', value: 5 },
  ]
  if (form.value.show_saturday) allDays.push({ label: '周六', value: 6 })
  if (form.value.show_sunday) allDays.push({ label: '周日', value: 7 })
  return allDays
})

const gridColumns = computed(() => {
  return `4rem repeat(${days.value.length}, 1fr)`
})

// 初始化
onMounted(async () => {
    await store.fetchConfig()
    await store.fetchCourses()
    if (config.value) {
        form.value = { 
            ...config.value, 
            period_settings_json: config.value.period_settings_json || '',
            course_colors_json: config.value.course_colors_json || '',
            background_emojis: config.value.background_emojis || ''
        }
        try {
            if (form.value.period_settings_json) {
                periodSettings.value = JSON.parse(form.value.period_settings_json)
            }
            if (form.value.course_colors_json) {
                courseColors.value = JSON.parse(form.value.course_colors_json)
            }
        } catch (e) {
            console.error("Failed to parse settings", e)
        }
    }
    await loadTimetableForEdit()
})

// 监听配置变化，自动保存配置并重新加载课表
async function onConfigChange() {
    const newConfig = {
        ...form.value,
        period_settings_json: JSON.stringify(periodSettings.value)
    }
    // 保存配置到后端
    await timetableApi.updateConfig(newConfig)
    
    // 更新 Store 和本地缓存
    if (store.config) {
        await store.updateConfig({
            ...store.config,
            ...newConfig
        })
    }
    
    // 重新加载对应课表
    await loadTimetableForEdit()
}

function goBack() {
    router.back()
}

async function loadTimetableForEdit() {
    await store.fetchTimetable(form.value.current_grade, form.value.current_semester)
    // 深度拷贝到编辑副本
    editTimetable.value = JSON.parse(JSON.stringify(originalTimetable.value))
}

function getCourse(day: number, period: number) {
    return editTimetable.value.find(t => t.day_of_week === day && t.period === period)
}

function getCourseName(day: number, period: number) {
    const item = getCourse(day, period)
    if (item && item.course_id) {
        const course = courses.value.find(c => c.id === item.course_id)
        return course ? course.name : ''
    }
    return ''
}

function getCellStyle(day: number, period: number) {
    const item = getCourse(day, period)
    if (item && item.course_id) {
        const course = courses.value.find(c => c.id === item.course_id)
        if (course) {
            const color = getCourseColor(course)
            return {
                backgroundColor: color, // 移除透明度，保持颜色一致
                color: isLightColor(color) ? '#000' : '#fff',
                border: 'none'
            }
        }
    }
    return {}
}

function openCourseSelector(day: number, period: number) {
    currentCell.value = { day, period }
    selectorVisible.value = true
}

function selectCourse(course: Course) {
    if (!currentCell.value) return
    const { day, period } = currentCell.value
    
    // 移除旧的
    const idx = editTimetable.value.findIndex(t => t.day_of_week === day && t.period === period)
    if (idx > -1) {
        editTimetable.value.splice(idx, 1)
    }
    
    // 添加新的
    editTimetable.value.push({
        day_of_week: day,
        period: period,
        course_id: course.id,
    })
    
    selectorVisible.value = false
    saveTimetable(true)
}

function clearCell() {
    if (!currentCell.value) return
    const { day, period } = currentCell.value
    const idx = editTimetable.value.findIndex(t => t.day_of_week === day && t.period === period)
    if (idx > -1) {
        editTimetable.value.splice(idx, 1)
    }
    selectorVisible.value = false
    saveTimetable(true)
}

async function saveTimetable(silent = false) {
    saving.value = true
    try {
        await timetableApi.saveTimetable({
            grade: form.value.current_grade,
            semester: form.value.current_semester,
            items: editTimetable.value
        })
        if (!silent) ElMessage.success('保存成功')
        // 刷新 Store
        await store.fetchTimetable(form.value.current_grade, form.value.current_semester)
    } catch {
        if (!silent) ElMessage.error('保存失败')
    } finally {
        saving.value = false
    }
}

function openAddCourseDialog() {
    editingCourseId.value = null
    courseForm.value = { name: '', color: '#409EFF' }
    showCourseDialog.value = true
}

function openEditCourseDialog(course: Course) {
    editingCourseId.value = course.id
    courseForm.value = { name: course.name, color: course.color }
    showCourseDialog.value = true
}

async function handleDeleteCourse(course: Course) {
    try {
        await ElMessageBox.confirm(`确定要删除课程“${course.name}”吗？此操作不可恢复。`, '提示', {
            confirmButtonText: '删除',
            cancelButtonText: '取消',
            type: 'warning'
        })
        
        await timetableApi.deleteCourse(course.id)
        ElMessage.success('删除成功')
        await store.fetchCourses()
    } catch (e: any) {
        if (e !== 'cancel') {
             ElMessage.error('删除失败: ' + (e.message || '未知错误'))
        }
    }
}

async function saveCourse() {
    if (!courseForm.value.name) return
    savingCourse.value = true
    try {
        if (editingCourseId.value) {
             await timetableApi.updateCourse(editingCourseId.value, courseForm.value)
             ElMessage.success('更新成功')
        } else {
             await timetableApi.addCourse(courseForm.value)
             ElMessage.success('添加成功')
        }
        await store.fetchCourses()
        showCourseDialog.value = false
    } catch (e: any) {
        ElMessage.error((editingCourseId.value ? '更新' : '添加') + '失败: ' + (e.message || '未知错误'))
        console.error(e)
    } finally {
        savingCourse.value = false
    }
}

function editPeriodTime(period: number) {
    const existing = periodSettings.value.find(p => p.period === period)
    editingPeriod.value = existing ? { ...existing } : { period, start_time: '', end_time: '' }
    timeDialogVisible.value = true
}

function getPeriodTime(period: number) {
    const s = periodSettings.value.find(p => p.period === period)
    if (s && s.start_time && s.end_time) {
        return `${s.start_time}-${s.end_time}`
    }
    return ''
}

async function savePeriodTime() {
    if (!editingPeriod.value) return
    
    const idx = periodSettings.value.findIndex(p => p.period === editingPeriod.value!.period)
    if (idx > -1) {
        periodSettings.value[idx] = editingPeriod.value
    } else {
        periodSettings.value.push(editingPeriod.value)
    }
    
    timeDialogVisible.value = false
    
    // 立即保存配置
    await onConfigChange()
}
</script>

<style scoped>
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

:deep(.pretty-card .el-card__header) {
  border-bottom: 1px solid rgb(243 244 246 / 0.8);
}

.dark :deep(.pretty-card .el-card__header) {
  border-bottom: 1px solid rgb(31 41 55 / 0.6);
}

:deep(.pretty-dialog .el-dialog) {
  border-radius: 24px;
  overflow: hidden;
  max-width: calc(100vw - 24px);
}

.dark :deep(.pretty-dialog .el-dialog) {
  background: rgb(17 24 39 / 0.92);
  border: 1px solid rgb(31 41 55 / 0.9);
}

:deep(.pretty-dialog .el-dialog__header) {
  padding: 16px 16px 12px;
  margin-right: 0;
}

:deep(.pretty-dialog .el-dialog__body) {
  padding: 16px;
}

:deep(.pretty-dialog .el-dialog__footer) {
  padding: 12px 16px 16px;
}
</style>
