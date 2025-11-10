// 中文注释：任务分类全局 Store，支持颜色设置与本地持久化
import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export type TaskCategory = { name: string; color: string }

const LS_KEY = 'recordgo.task_categories'

function loadInitial(): TaskCategory[] {
  try {
    const raw = localStorage.getItem(LS_KEY)
    if (raw) {
      const arr = JSON.parse(raw)
      if (Array.isArray(arr)) return arr.filter(x => x && typeof x.name === 'string' && typeof x.color === 'string')
    }
  } catch (_) {}
  // 中文注释：默认提供三类常用分类，并赋予直观颜色
  return [
    { name: '语文', color: '#ef4444' }, // 红色
    { name: '数学', color: '#10b981' }, // 绿色
    { name: '英语', color: '#3b82f6' }  // 蓝色
  ]
}

export const useTaskCategories = defineStore('taskCategories', () => {
  const categories = ref<TaskCategory[]>(loadInitial())

  // 中文注释：持久化到本地存储
  watch(categories, (val) => {
    try { localStorage.setItem(LS_KEY, JSON.stringify(val)) } catch (_) {}
  }, { deep: true })

  function list(): TaskCategory[] { return categories.value }
  function colorOf(name?: string | null): string {
    const n = String(name || '').trim()
    const found = categories.value.find(c => c.name === n)
    return found?.color || '#64748b' // 默认灰蓝
  }
  function add(name: string, color: string) {
    const n = String(name || '').trim()
    if (!n) return
    if (categories.value.some(c => c.name === n)) return
    categories.value.push({ name: n, color: color || '#64748b' })
  }
  function update(oldName: string, patch: Partial<TaskCategory>) {
    const idx = categories.value.findIndex(c => c.name === oldName)
    if (idx < 0) return
    const next = { ...categories.value[idx], ...patch }
    // 防止重名冲突
    if (patch.name && patch.name !== oldName && categories.value.some(c => c.name === patch.name)) return
    categories.value.splice(idx, 1, next)
  }
  function remove(name: string) {
    categories.value = categories.value.filter(c => c.name !== name)
  }

  return { categories, list, colorOf, add, update, remove }
})