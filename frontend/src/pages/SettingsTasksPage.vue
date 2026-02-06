<template>
  <SettingsShell title="任务设置" subtitle="分类管理与显示偏好" :icon="List" tone="emerald" container-class="max-w-4xl">
    <SettingsCard
      v-if="isVIP"
      title="显示任务备注入口"
      description="关闭后：任务清单卡片不显示备注图标，操作菜单不显示“备注”按钮"
    >
      <template #right>
        <el-switch v-model="taskNotesEnabled" />
      </template>
    </SettingsCard>
    <SettingsCard v-else title="任务备注功能" description="该功能需要开通 VIP 会员才能使用">
      <div class="text-sm text-gray-700 dark:text-gray-200">添加微信：tabor2024，备注“任务家”，获取 VIP 资格</div>
    </SettingsCard>

    <SettingsCard
      title="任务自动排序"
      description="开启后：分类内已完成任务排在下方；当某分类全部完成时，该分类整体排在未完分类之后"
    >
      <template #right>
        <el-switch v-model="taskAutoSortEnabled" />
      </template>
    </SettingsCard>

    <SettingsCard title="任务分类管理" description="支持新增、编辑、删除与拖拽排序">
      <div class="space-y-4">
        <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:gap-3">
          <el-input v-model="newName" placeholder="输入分类名称，如：语文" class="w-full sm:max-w-[320px]" />
          <div class="flex items-center gap-3">
            <el-color-picker v-model="newColor" />
            <el-button type="primary" class="!rounded-xl" @click="addCategory">
              <el-icon class="mr-1"><Plus /></el-icon>添加分类
            </el-button>
          </div>
        </div>

        <div class="space-y-2">
          <div
            v-for="(c, idx) in sortable"
            :key="c.name"
            class="group flex items-center justify-between gap-3 rounded-2xl border border-gray-100 dark:border-gray-800 bg-white/50 dark:bg-gray-900/40 px-3 py-3 hover:bg-white/80 dark:hover:bg-gray-900/70 transition"
            draggable="true"
            @dragstart="onDragStart(idx)"
            @dragover="onDragOver"
            @drop="onDrop(idx)"
          >
            <div class="flex items-center gap-3 min-w-0">
              <el-icon :size="16" class="cursor-grab text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300"><Sort /></el-icon>
              <span class="inline-block w-4 h-4 rounded-lg ring-1 ring-black/5 dark:ring-white/10" :style="{ backgroundColor: c.color }"></span>
              <div class="flex items-center gap-2 min-w-0">
                <el-input v-model="editable[c.name].name" class="w-[160px] sm:w-[200px]" />
                <el-color-picker v-model="editable[c.name].color" />
              </div>
            </div>
            <div class="flex items-center gap-2 shrink-0">
              <el-button size="small" class="!rounded-xl" @click="saveCategory(c.name)">
                <el-icon class="mr-1"><Edit /></el-icon>保存
              </el-button>
              <el-button size="small" type="danger" class="!rounded-xl" @click="removeCategory(c.name)">
                <el-icon class="mr-1"><Delete /></el-icon>删除
              </el-button>
            </div>
          </div>
        </div>

        <div class="text-xs text-gray-500 dark:text-gray-400">提示：按住左侧拖拽图标，可调整分类顺序。</div>
      </div>
    </SettingsCard>
  </SettingsShell>
</template>

<script setup lang="ts">
// 中文注释：引入所需图标与 Color Picker 样式（避免样式加载异常）
import { List, Plus, Edit, Delete, Sort } from '@element-plus/icons-vue'
import 'element-plus/es/components/color-picker/style/css'
import { useAppState } from '@/stores/appState'
import { useAuth } from '@/stores/auth'
import { computed, reactive, ref, watchEffect, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useTaskCategories } from '@/stores/categories'
import type { TaskCategory } from '@/stores/categories'
import SettingsShell from '@/components/settings/SettingsShell.vue'
import SettingsCard from '@/components/settings/SettingsCard.vue'

// 中文注释：联动全局状态中的任务备注开关，使用计算属性实现双向绑定
const store = useAppState()
const auth = useAuth()
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
