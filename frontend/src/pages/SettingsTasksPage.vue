<template>
  <div class="p-4 space-y-4">
    <div class="flex items-center gap-2">
      <el-icon :size="18" class="cursor-pointer" style="color:#64748b" @click="goBack"><ArrowLeft /></el-icon>
      <el-icon :size="18" style="color:#10b981"><List /></el-icon>
      <h2 class="font-semibold">任务设置</h2>
    </div>
    <!-- 中文注释：任务备注入口开关（默认开启）；关闭后任务卡片不显示备注图标，菜单不显示备注按钮 -->
    <el-card shadow="never" v-if="isVIP">
      <div class="flex items-center justify-between">
        <div>
          <div class="font-medium">显示任务备注入口</div>
          <div class="text-xs text-gray-500 mt-1">关闭后：任务清单卡片不显示备注图标，操作菜单不显示“备注”按钮</div>
        </div>
        <el-switch v-model="taskNotesEnabled" />
      </div>
    </el-card>
    <el-card shadow="never" v-else>
      <div>
        <div class="font-medium">任务备注功能</div>
        <div class="text-sm text-gray-600 mt-1">该功能需要开通VIP会员才能使用</div>
        <div class="text-xs text-gray-500 mt-2">添加微信：tabor2024，备注“任务家”，获取VIP资格</div>
      </div>
    </el-card>
    <!-- 中文注释：任务自动排序开关（默认开启）；开启后分类内已完成排下方，全部完成的分类排到未完分类之后 -->
    <el-card shadow="never">
      <div class="flex items-center justify-between">
        <div>
          <div class="font-medium">任务自动排序</div>
          <div class="text-xs text-gray-500 mt-1">开启后：分类内已完成任务排在下方；当某分类全部完成时，该分类整体排在未完分类之后</div>
        </div>
        <el-switch v-model="taskAutoSortEnabled" />
      </div>
    </el-card>
    <!-- 分类管理（移动自“任务分类设置”页面）：新增、编辑、删除、拖拽排序 -->
    <el-card shadow="never">
      <div class="font-medium mb-3">任务分类管理</div>
      <!-- 新增分类 -->
      <div class="flex items-center gap-3 mb-4">
        <el-input v-model="newName" placeholder="输入分类名称，如：语文" style="max-width: 300px" />
        <el-color-picker v-model="newColor" />
        <el-button type="primary" @click="addCategory"><el-icon class="mr-1"><Plus /></el-icon>添加分类</el-button>
      </div>
      <div class="text-xs text-gray-500 mb-2">中文注释：支持添加、修改、删除分类，并设置颜色；拖拽排序与列表展示同步。</div>

      <!-- 分类列表（可编辑名称与颜色；支持拖拽排序） -->
      <div class="space-y-2">
        <div
          v-for="(c, idx) in sortable"
          :key="c.name"
          class="flex items-center justify-between px-2 py-2 rounded border border-gray-200"
          draggable="true"
          @dragstart="onDragStart(idx)"
          @dragover="onDragOver"
          @drop="onDrop(idx)"
        >
          <div class="flex items-center gap-3">
            <!-- 中文注释：拖拽把手 -->
            <el-icon :size="16" class="cursor-grab text-gray-500"><Sort /></el-icon>
            <span class="inline-block w-4 h-4 rounded" :style="{ backgroundColor: c.color }"></span>
            <el-input v-model="editable[c.name].name" style="width: 160px" />
            <el-color-picker v-model="editable[c.name].color" />
          </div>
          <div class="flex items-center gap-2">
            <el-button size="small" @click="saveCategory(c.name)"><el-icon class="mr-1"><Edit /></el-icon>保存</el-button>
            <el-button size="small" type="danger" @click="removeCategory(c.name)"><el-icon class="mr-1"><Delete /></el-icon>删除</el-button>
          </div>
        </div>
      </div>
    </el-card>
    <el-card shadow="never">
      <div class="text-gray-600">其他配置项待完善：任务默认提醒、重复规则等。</div>
    </el-card>
    <!-- 取消底部返回按钮；统一使用标题左侧返回图标 -->
    <div class="mt-2">
      <el-button v-if="isAdmin" type="primary" @click="router.push('/admin')">用户管理</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
// 中文注释：引入所需图标与 Color Picker 样式（避免样式加载异常）
import { List, ArrowLeft, Plus, Edit, Delete, Sort } from '@element-plus/icons-vue'
import 'element-plus/es/components/color-picker/style/css'
import router from '@/router'
import { useAppState } from '@/stores/appState'
import { useAuth } from '@/stores/auth'
import { computed, reactive, ref, watchEffect, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useTaskCategories } from '@/stores/categories'
import type { TaskCategory } from '@/stores/categories'
function goBack() { router.back() }

// 中文注释：联动全局状态中的任务备注开关，使用计算属性实现双向绑定
const store = useAppState()
const auth = useAuth()
const isAdmin = computed(() => Number(auth.user?.id || 0) === 1)
const isVIP = computed(() => {
  const u = auth.user
  if (!u) return false
  const lifetime = !!(u as any).is_lifetime_vip
  const vip = !!(u as any).is_vip
  const expire = (u as any).vip_expire_time ? new Date((u as any).vip_expire_time as string) : null
  const valid = lifetime || (vip && !!expire && expire.getTime() > Date.now())
  return valid
})
const taskNotesEnabled = computed({
  get: () => store.taskNotesEnabled,
  set: (v: boolean) => store.setTaskNotesEnabled(v)
})
const taskAutoSortEnabled = computed({
  get: () => store.taskAutoSortEnabled,
  set: (v: boolean) => store.setTaskAutoSortEnabled(v)
})

// ===== 分类管理逻辑（与原“任务分类设置”保持一致） =====
const cats = useTaskCategories()
const categories = computed<TaskCategory[]>(() => cats.list())
onMounted(async () => { try { await cats.syncFromServer() } catch {} })
// 可排序显示数组
const sortable = ref<{ name: string; color: string }[]>([])
watchEffect(() => {
  sortable.value = categories.value.map((c: TaskCategory) => ({ name: c.name, color: c.color }))
})
// 新增分类输入与提交
const newName = ref('')
const newColor = ref('#64748b')
function addCategory() {
  const name = newName.value.trim()
  if (!name) return
  cats.add(name, newColor.value || '#64748b')
  newName.value = ''
  newColor.value = '#64748b'
}
// 编辑缓冲区
const editable = reactive<Record<string, { name: string; color: string }>>({})
watchEffect(() => {
  const next: Record<string, { name: string; color: string }> = {}
  for (const c of categories.value as TaskCategory[]) next[c.name] = { name: c.name, color: c.color }
  Object.assign(editable, next)
})
function saveCategory(oldName: string) {
  const patch = editable[oldName]
  if (!patch) return
  cats.update(oldName, { name: (patch.name || '').trim(), color: patch.color })
}
function removeCategory(name: string) { cats.remove(name) }
// 拖拽排序
const draggingIndex = ref<number | null>(null)
function onDragStart(idx: number) { draggingIndex.value = idx }
function onDragOver(e: DragEvent) { e.preventDefault() }
function onDrop(targetIdx: number) {
  const from = draggingIndex.value
  draggingIndex.value = null
  if (from === null || from === targetIdx) return
  const arr = [...sortable.value]
  const [moved] = arr.splice(from, 1)
  arr.splice(targetIdx, 0, moved)
  sortable.value = arr
  for (let i = 0; i < arr.length; i++) {
    cats.setOrder(arr[i].name, i + 1)
  }
  ElMessage.success('已更新分类顺序')
}
</script>

<style scoped>
</style>
