import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import http from '@/services/http'

export interface Tag {
  id: string
  name: string
  color?: string
  parentId?: string
  children?: Tag[]
  count?: number
}

export interface GrowthRecord {
  id: string
  date: string
  content: string
  images: string[]
  tags: string[] // Tag IDs
  audio?: string
  created_at?: string
  is_pinned?: boolean
}

export const useLittleGrowthStore = defineStore('littleGrowth', () => {
  // --- State ---
  const records = ref<GrowthRecord[]>([])
  const tags = ref<Tag[]>([]) // Backend tags
  const activeFilterTagId = ref<string | null>(null)
  const loading = ref(false)

  // --- Getters ---
  const flattenedTags = computed(() => {
    // Map tags to include counts from records
    // Note: records.tags are IDs now.
    const tagCounts = new Map<string, number>()
    records.value.forEach(r => {
      let rTags: string[] = []
      if (typeof r.tags === 'string') {
        try { rTags = JSON.parse(r.tags) } catch {}
      } else if (Array.isArray(r.tags)) {
        rTags = r.tags
      }

      rTags.forEach(tid => {
        tagCounts.set(String(tid), (tagCounts.get(String(tid)) || 0) + 1)
      })
    })

    return tags.value.map(t => ({
      ...t,
      id: String(t.id), // Ensure string ID
      count: tagCounts.get(String(t.id)) || 0
    }))
  })

  const filteredRecords = computed(() => {
    if (!activeFilterTagId.value) {
      return records.value
    }
    return records.value.filter(r => {
      let rTags: string[] = []
      if (typeof r.tags === 'string') {
        try { rTags = JSON.parse(r.tags) } catch {}
      } else {
        rTags = r.tags
      }
      // Ensure comparison as strings
      return rTags.map(String).includes(String(activeFilterTagId.value))
    })
  })

  // --- Actions ---
  async function fetchTags() {
    const res = await http.get('/little-growth/tags')
    // Ensure IDs are strings
    tags.value = (res || []).map((t: any) => ({ ...t, id: String(t.id) }))
  }

  async function createTag(name: string, color?: string) {
    const res = await http.post('/little-growth/tags', { name, color })
    const newTag = { ...res, id: String(res.id) }
    tags.value.push(newTag)
    return newTag
  }

  async function updateTag(id: string, name: string, color?: string) {
    const res = await http.put(`/little-growth/tags/${id}`, { name, color })
    const idx = tags.value.findIndex(t => String(t.id) === String(id))
    if (idx !== -1) {
      tags.value[idx] = { ...res, id: String(res.id) }
    }
  }

  async function deleteTag(id: string) {
    await http.delete(`/little-growth/tags/${id}`)
    tags.value = tags.value.filter(t => String(t.id) !== String(id))
    // Also update records locally to reflect tag removal
    records.value.forEach(r => {
        let rTags: string[] = []
        if (Array.isArray(r.tags)) rTags = r.tags
        // remove id
        const newTags = rTags.filter(tid => String(tid) !== String(id))
        if (newTags.length !== rTags.length) {
            r.tags = newTags
        }
    })
  }

  async function fetchRecords() {
    loading.value = true
    try {
      await fetchTags() // Ensure tags are loaded
      const res = await http.get('/little-growth/records')
      const items = (res.items || []).map((r: any) => {
        let images: string[] = []
        try { 
          if (typeof r.images === 'string') {
            images = JSON.parse(r.images)
            if (!Array.isArray(images)) images = []
          } else if (Array.isArray(r.images)) {
            images = r.images
          }
        } catch {}
        
        images = images.map(i => i.startsWith('http') ? i : `${import.meta.env.VITE_API_BASE}/api/${i}`)
        
        let audio = r.audio
        if (audio && !audio.startsWith('http')) {
            audio = `${import.meta.env.VITE_API_BASE}/api/${audio}`
        }

        let rTags: string[] = []
        try { rTags = JSON.parse(r.tags) } catch {}
        // Ensure tags are array of strings
        rTags = rTags.map(String)

        return {
          ...r,
          date: r.date.split('T')[0],
          images,
          audio,
          tags: rTags
        }
      })
      records.value = items
    } finally {
      loading.value = false
    }
  }

  async function createRecord(data: any) {
    // data contains images URLs, content, date, tags (IDs)
    await http.post('/little-growth/records', data)
    await fetchRecords()
  }

  async function deleteRecord(id: string) {
    await http.delete(`/little-growth/records/${id}`)
    records.value = records.value.filter(r => String(r.id) !== String(id))
  }

  async function togglePin(id: string) {
    const res = await http.patch(`/little-growth/records/${id}/pin`)
    const idx = records.value.findIndex(r => String(r.id) === String(id))
    if (idx !== -1) {
        records.value[idx].is_pinned = res.is_pinned
    }
    // Re-fetch to sort correctly or sort locally? Backend sort is better.
    await fetchRecords()
  }

  function getRecordById(id: string) {
    return records.value.find(r => String(r.id) === String(id))
  }

  return {
    records,
    tags,
    activeFilterTagId,
    filteredRecords,
    flattenedTags,
    loading,
    fetchRecords,
    fetchTags,
    createTag,
    updateTag,
    deleteTag,
    createRecord,
    deleteRecord,
    togglePin,
    getRecordById
  }
})
