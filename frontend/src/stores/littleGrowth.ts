import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import http from '@/services/http'

export interface Tag {
  id: string
  name: string
  parentId?: string
  children?: Tag[]
  count?: number // For display
}

export interface GrowthRecord {
  id: string // number in backend, string here for safety
  date: string // YYYY-MM-DD
  content: string
  images: string[] // URLs
  tags: string[] // Tag Names or IDs
  audio?: string
  created_at?: string
}

export const useLittleGrowthStore = defineStore('littleGrowth', () => {
  // --- State ---
  const records = ref<GrowthRecord[]>([])
  const activeFilterTagId = ref<string | null>(null)
  const loading = ref(false)

  // --- Getters ---
  // Dynamic tags based on records content
  const flattenedTags = computed(() => {
    // Extract all unique tags from records
    const tagCounts = new Map<string, number>()
    records.value.forEach(r => {
      // Parse tags if they are JSON string, or assuming they are array
      // Backend returns JSON string for tags? No, we should parse it in frontend service or here.
      // Let's assume the API returns parsed object or we parse it.
      // Actually, backend `GrowthRecord` struct has `Tags string`.
      // So we need to parse it.
      let rTags: string[] = []
      if (typeof r.tags === 'string') {
        try { rTags = JSON.parse(r.tags) } catch {}
      } else if (Array.isArray(r.tags)) {
        rTags = r.tags
      }

      rTags.forEach(t => {
        tagCounts.set(t, (tagCounts.get(t) || 0) + 1)
      })
    })

    // Convert to Tag objects
    const list: Tag[] = []
    tagCounts.forEach((count, name) => {
      list.push({ id: name, name, count }) // ID is name for simplicity
    })
    return list
  })

  // Hierarchical tags (mock hierarchy for now as backend doesn't store hierarchy yet, or flat)
  // User said: "Sidebar tags related to input tags".
  // We can just show a flat list or try to group them if we had logic.
  // For now, flat list of used tags.
  const tags = computed(() => flattenedTags.value)

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
      return rTags.includes(activeFilterTagId.value!)
    })
  })

  // --- Actions ---
  async function fetchRecords() {
    loading.value = true
    try {
      const res = await http.get('/little-growth/records')
      // Parse JSON fields
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
        
        // Backend returns relative path "uploads/...", we need full URL
        images = images.map(i => i.startsWith('http') ? i : `${import.meta.env.VITE_API_BASE}/api/${i}`)
        
        let audio = r.audio
        if (audio && !audio.startsWith('http')) {
            audio = `${import.meta.env.VITE_API_BASE}/api/${audio}`
        }

        let tags: string[] = []
        try { tags = JSON.parse(r.tags) } catch {}

        return {
          ...r,
          date: r.date.split('T')[0],
          images,
          audio,
          tags
        }
      })
      records.value = items
    } finally {
      loading.value = false
    }
  }

  async function createRecord(formData: FormData) {
    await http.post('/little-growth/records', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    await fetchRecords()
  }

  async function deleteRecord(id: string) {
    await http.delete(`/little-growth/records/${id}`)
    records.value = records.value.filter(r => String(r.id) !== String(id))
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
    createRecord,
    deleteRecord,
    getRecordById
  }
})
