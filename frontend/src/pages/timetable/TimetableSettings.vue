<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 pb-20">
    <!-- Header -->
    <div class="bg-white dark:bg-gray-800 shadow-sm px-4 py-3 flex items-center gap-3 sticky top-0 z-20">
        <el-icon :size="20" class="cursor-pointer" @click="goBack"><ArrowLeft /></el-icon>
        <h1 class="font-bold text-lg">课表设置</h1>
    </div>

    <div class="p-4 space-y-6">
        <!-- 基础配置 -->
        <el-card shadow="never" class="rounded-xl">
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
                <el-form-item label="背景Emoji">
                    <el-input v-model="form.background_emojis" placeholder="输入Emoji，用逗号分隔，例如: 🌟,🎈,🐱" @change="onConfigChange" />
                </el-form-item>
            </el-form>
        </el-card>

        <!-- 课程库管理 -->
        <el-card shadow="never" class="rounded-xl">
            <template #header>
                <div class="flex justify-between items-center">
                    <div class="font-bold">课程库</div>
                    <el-button type="primary" link @click="openAddCourseDialog">添加课程</el-button>
                </div>
            </template>
            <div class="flex flex-wrap gap-2">
                <div v-for="course in courses" :key="course.id" 
                     class="group relative px-3 py-1 rounded-full text-sm font-medium text-white shadow-sm flex items-center gap-1 cursor-default"
                     :style="{ backgroundColor: course.color }">
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
        <el-card shadow="never" class="rounded-xl overflow-visible">
            <template #header>
                <div class="flex justify-between items-center">
                    <div class="font-bold">课表编排 ({{ form.current_grade }} {{ form.current_semester }})</div>
                    <el-button type="primary" size="small" @click="saveTimetable" :loading="saving">保存编排</el-button>
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
    <el-dialog v-model="selectorVisible" title="选择课程" width="90%" class="rounded-xl" append-to-body>
        <div class="grid grid-cols-4 gap-3">
            <div v-for="course in courses" :key="course.id"
                 class="h-10 rounded-lg flex items-center justify-center text-white text-sm font-medium shadow-sm cursor-pointer active:scale-95 transition-transform"
                 :style="{ backgroundColor: course.color }"
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
    <el-dialog v-model="showCourseDialog" :title="editingCourseId ? '编辑课程' : '添加新课程'" width="80%" append-to-body>
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
    <el-dialog v-model="timeDialogVisible" title="设置课程时间" width="80%" append-to-body>
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
    background_emojis: ''
})

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
            background_emojis: config.value.background_emojis || ''
        }
        try {
            if (form.value.period_settings_json) {
                periodSettings.value = JSON.parse(form.value.period_settings_json)
            }
        } catch (e) {
            console.error("Failed to parse period settings", e)
        }
    }
    await loadTimetableForEdit()
})

// 监听配置变化，自动保存配置并重新加载课表
async function onConfigChange() {
    // 保存配置
    await timetableApi.updateConfig({
        ...form.value,
        period_settings_json: JSON.stringify(periodSettings.value)
    })
    // 刷新 store config
    if (store.config) {
        store.config = { 
            ...store.config, 
            ...form.value, 
            period_settings_json: JSON.stringify(periodSettings.value) 
        }
    }
    
    // 重新加载对应课表
    await loadTimetableForEdit()
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
            return {
                backgroundColor: hexToRgba(course.color, 0.5),
                color: '#000',
                border: 'none'
            }
        }
    }
    return {}
}

function hexToRgba(hex: string, alpha: number) {
    let c: any;
    if(/^#([A-Fa-f0-9]{3}){1,2}$/.test(hex)){
        c= hex.substring(1).split('');
        if(c.length== 3){
            c= [c[0], c[0], c[1], c[1], c[2], c[2]];
        }
        c= '0x'+c.join('');
        return 'rgba('+[(c>>16)&255, (c>>8)&255, c&255].join(',')+','+alpha+')';
    }
    return hex; // Fallback
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
}

function clearCell() {
    if (!currentCell.value) return
    const { day, period } = currentCell.value
    const idx = editTimetable.value.findIndex(t => t.day_of_week === day && t.period === period)
    if (idx > -1) {
        editTimetable.value.splice(idx, 1)
    }
    selectorVisible.value = false
}

async function saveTimetable() {
    saving.value = true
    try {
        await timetableApi.saveTimetable({
            grade: form.value.current_grade,
            semester: form.value.current_semester,
            items: editTimetable.value
        })
        ElMessage.success('保存成功')
        // 刷新 Store
        await store.fetchTimetable(form.value.current_grade, form.value.current_semester)
    } catch {
        ElMessage.error('保存失败')
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

function goBack() {
    router.back()
}
</script>