// 中文注释：任务分类全局 Store，支持颜色设置与本地持久化
import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import http from '@/services/http'

// 中文注释：任务分类模型，加入 order 字段用于自定义排序（数值越小越靠前）
export type TaskCategory = { name: string; color: string; order: number }

const LS_KEY = 'recordgo.task_categories'

function loadInitial(): TaskCategory[] {
  try {
    const raw = localStorage.getItem(LS_KEY)
    if (raw) {
      const arr = JSON.parse(raw)
      if (Array.isArray(arr)) {
        return arr
          .filter(x => x && typeof x.name === 'string' && typeof x.color === 'string')
          .map((x: any, i: number) => ({ name: String(x.name), color: String(x.color), order: Number(x.order ?? i + 1) }))
          .sort((a: TaskCategory, b: TaskCategory) => a.order - b.order)
      }
    }
  } catch (_) {}
  // 中文注释：默认提供三类常用分类，并赋予直观颜色
  return [
    { name: '语文', color: '#ef4444', order: 1 }, // 红色
    { name: '数学', color: '#10b981', order: 2 }, // 绿色
    { name: '英语', color: '#3b82f6', order: 3 }  // 蓝色
  ]
}

export const useTaskCategories = defineStore('taskCategories', () => {
  const categories = ref<TaskCategory[]>(loadInitial())

  // 中文注释：持久化到本地存储
  watch(categories, (val) => {
    try { localStorage.setItem(LS_KEY, JSON.stringify(val)) } catch (_) {}
  }, { deep: true })

  // 中文注释：基础校验（名称与颜色）
  function isValidHex(color: string): boolean { return /^#([0-9a-fA-F]{6})$/.test(color) }
  function normalizeColor(color: string): string { return isValidHex(color) ? color : '#64748b' }
  function validateName(name: string): string { return String(name || '').trim().slice(0, 32) }

  function list(): TaskCategory[] { return [...categories.value].sort((a, b) => a.order - b.order) }
  function colorOf(name?: string | null): string {
    const n = String(name || '').trim()
    const found = categories.value.find(c => c.name === n)
    return found?.color || '#64748b' // 默认灰蓝
  }
  function orderOf(name?: string | null): number {
    const n = String(name || '').trim()
    const found = categories.value.find(c => c.name === n)
    return found?.order ?? Number.MAX_SAFE_INTEGER
  }
  function add(name: string, color: string): boolean {
    const n = validateName(name)
    if (!n) return false
    if (categories.value.some(c => c.name === n)) return false
    const nextOrder = (categories.value.reduce((m, c) => Math.max(m, c.order), 0) || 0) + 1
    categories.value.push({ name: n, color: normalizeColor(color || ''), order: nextOrder })
    void persistToServer()
    return true
  }
  function update(oldName: string, patch: Partial<TaskCategory>): boolean {
    const idx = categories.value.findIndex(c => c.name === oldName)
    if (idx < 0) return false
    const next = { ...categories.value[idx], ...patch }
    // 防止重名冲突
    const newName = validateName(String(next.name))
    if (newName && newName !== oldName && categories.value.some(c => c.name === newName)) return false
    next.name = newName || categories.value[idx].name
    next.color = normalizeColor(String(next.color || categories.value[idx].color))
    next.order = Number.isFinite(next.order) ? Number(next.order) : categories.value[idx].order
    categories.value.splice(idx, 1, next)
    void persistToServer()
    return true
  }
  function remove(name: string): boolean {
    const before = categories.value.length
    categories.value = categories.value.filter(c => c.name !== name)
    // 中文注释：删除后标准化顺序（保证连续）
    categories.value = categories.value
      .sort((a, b) => a.order - b.order)
      .map((c, i) => ({ ...c, order: i + 1 }))
    if (categories.value.length !== before) void persistToServer()
    return categories.value.length !== before
  }
  function setOrder(name: string, order: number): boolean {
    const idx = categories.value.findIndex(c => c.name === name)
    if (idx < 0) return false
    const o = Math.max(1, Math.floor(order))
    categories.value[idx].order = o
    // 中文注释：重新排序并规范化为 1..N
    categories.value = categories.value
      .sort((a, b) => a.order - b.order)
      .map((c, i) => ({ ...c, order: i + 1 }))
    // 中文注释：优先调用后端的“仅更新顺序”接口，失败则回退为整体持久化
    void persistOrder(name, o)
    return true
  }

  // ===== 后端同步（可选）：若后端提供 /task-categories 接口，则读写对接；失败时静默回退到本地 =====
  async function syncFromServer(): Promise<boolean> {
    try {
      const data = await http.get<any>('/task-categories')
      const arr: any[] = Array.isArray(data) ? data : (Array.isArray(data?.items) ? data.items : [])
      const normalized: TaskCategory[] = arr
        .filter(x => x && typeof x.name === 'string')
        .map((x, i) => ({
          name: validateName(String(x.name)),
          color: normalizeColor(String(x.color || '#64748b')),
          order: Number.isFinite(Number(x.order)) ? Number(x.order) : (i + 1)
        }))
        .sort((a, b) => a.order - b.order)
        .map((c, i) => ({ ...c, order: i + 1 }))
      if (normalized.length > 0) categories.value = normalized
      return true
    } catch {
      return false
    }
  }
  async function persistToServer(): Promise<boolean> {
    try {
      // 中文注释：兼容两种写法：直接数组或 {items:[]} 包裹
      const payload = { items: categories.value }
      await http.put('/task-categories', payload)
      return true
    } catch {
      return false
    }
  }
  // 中文注释：仅更新某个分类的顺序；后端若未提供该接口，将自动回退到整体更新
  async function persistOrder(name: string, order: number): Promise<boolean> {
    try {
      await http.patch('/task-categories/order', { name, order })
      return true
    } catch {
      return await persistToServer()
    }
  }

  return { categories, list, colorOf, orderOf, add, update, remove, setOrder, syncFromServer }
})