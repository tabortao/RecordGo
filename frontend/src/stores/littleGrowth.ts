import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import dayjs from 'dayjs'

export interface Tag {
  id: string
  name: string
  parentId?: string
  children?: Tag[]
  count?: number // For display
}

export interface GrowthRecord {
  id: string
  date: string // YYYY-MM-DD
  content: string
  images: string[]
  tags: string[] // Tag IDs
  audio?: string
  createdAt: number
}

export const useLittleGrowthStore = defineStore('littleGrowth', () => {
  // --- State ---
  const records = ref<GrowthRecord[]>([])
  const tags = ref<Tag[]>([
    {
      id: 't1',
      name: '生活',
      children: [
        { id: 't1-1', name: '生日', parentId: 't1', count: 0 },
        { id: 't1-2', name: '日常', parentId: 't1', count: 0 }
      ],
      count: 0
    },
    {
      id: 't2',
      name: '学习',
      children: [
        { id: 't2-1', name: '英语', parentId: 't2', count: 0 },
        { id: 't2-2', name: '钢琴', parentId: 't2', count: 0 }
      ],
      count: 0
    },
    {
      id: 't3',
      name: '旅行',
      children: [
        { id: 't3-1', name: '海边', parentId: 't3', count: 0 }
      ],
      count: 0
    }
  ])
  
  const activeFilterTagId = ref<string | null>(null)

  // --- Getters ---
  const flattenedTags = computed(() => {
    const list: Tag[] = []
    const traverse = (items: Tag[]) => {
      items.forEach(t => {
        list.push(t)
        if (t.children) traverse(t.children)
      })
    }
    traverse(tags.value)
    return list
  })

  const filteredRecords = computed(() => {
    if (!activeFilterTagId.value) {
      return records.value.sort((a, b) => b.createdAt - a.createdAt)
    }
    
    // Check if selected tag is parent or child
    const selectedTag = flattenedTags.value.find(t => t.id === activeFilterTagId.value)
    if (!selectedTag) return []

    const targetIds = new Set<string>()
    targetIds.add(selectedTag.id)
    if (selectedTag.children) {
      selectedTag.children.forEach(c => targetIds.add(c.id))
    }

    return records.value.filter(r => {
      return r.tags.some(tid => targetIds.has(tid))
    }).sort((a, b) => b.createdAt - a.createdAt)
  })

  // --- Actions ---
  function addRecord(record: Omit<GrowthRecord, 'id' | 'createdAt'>) {
    const newRecord: GrowthRecord = {
      ...record,
      id: `r${Date.now()}`,
      createdAt: Date.now()
    }
    records.value.unshift(newRecord)
    updateTagCounts()
  }

  function updateRecord(id: string, updates: Partial<GrowthRecord>) {
    const index = records.value.findIndex(r => r.id === id)
    if (index !== -1) {
      records.value[index] = { ...records.value[index], ...updates }
      updateTagCounts()
    }
  }

  function deleteRecord(id: string) {
    records.value = records.value.filter(r => r.id !== id)
    updateTagCounts()
  }

  function getRecordById(id: string) {
    return records.value.find(r => r.id === id)
  }

  function addTag(name: string, parentId?: string) {
    const newTag: Tag = {
      id: `t${Date.now()}`,
      name,
      parentId,
      count: 0,
      children: []
    }

    if (parentId) {
      const parent = tags.value.find(t => t.id === parentId)
      if (parent) {
        if (!parent.children) parent.children = []
        parent.children.push(newTag)
      } else {
        // Fallback to root if parent not found
        tags.value.push(newTag)
      }
    } else {
      tags.value.push(newTag)
    }
  }

  function updateTagCounts() {
    // Reset counts
    const counts = new Map<string, number>()
    
    records.value.forEach(r => {
      r.tags.forEach(tid => {
        counts.set(tid, (counts.get(tid) || 0) + 1)
        // Find parent and increment
        const tag = flattenedTags.value.find(t => t.id === tid)
        if (tag && tag.parentId) {
          counts.set(tag.parentId, (counts.get(tag.parentId) || 0) + 1)
        }
      })
    })

    const traverseAndUpdate = (items: Tag[]) => {
      items.forEach(t => {
        t.count = counts.get(t.id) || 0
        if (t.children) traverseAndUpdate(t.children)
      })
    }
    traverseAndUpdate(tags.value)
  }

  // Initial mock data
  if (records.value.length === 0) {
    addRecord({
      date: dayjs().format('YYYY-MM-DD'),
      content: '今天天气真好，带宝宝去公园玩了滑梯，他非常开心！#生活 #日常',
      images: [
        'https://picsum.photos/400/300?random=1',
        'https://picsum.photos/400/300?random=2',
        'https://picsum.photos/400/300?random=3'
      ],
      tags: ['t1', 't1-2']
    })
    addRecord({
      date: dayjs().subtract(1, 'day').format('YYYY-MM-DD'),
      content: '第一次上钢琴课，老师夸奖很有天赋。#学习 #钢琴',
      images: ['https://picsum.photos/400/300?random=4'],
      tags: ['t2', 't2-2']
    })
  }

  return {
    records,
    tags,
    activeFilterTagId,
    filteredRecords,
    flattenedTags,
    addRecord,
    updateRecord,
    deleteRecord,
    getRecordById,
    addTag
  }
}, {
  persist: true
})
