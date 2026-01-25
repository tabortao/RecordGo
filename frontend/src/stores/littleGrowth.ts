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

export interface GrowthComment {
  id: string
  record_id: string
  user_id: string
  content: string
  created_at: string
  user?: { nickname?: string; username?: string; avatar_path?: string }
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
  is_favorite?: boolean // New
  visibility?: number // New: 0=Family, 1=Private
  comments?: GrowthComment[] // New
  user?: { nickname?: string; username?: string; avatar_path?: string } // New: Author info
}

export const useLittleGrowthStore = defineStore('littleGrowth', () => {
  // --- State ---
  const records = ref<GrowthRecord[]>([])
  const tags = ref<Tag[]>([]) // Backend tags
  const activeFilterTagId = ref<string | null>(null)
  const onlyFavorites = ref(false) // New: Filter state
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
    // Client-side filtering for tags (since we might have fetched a subset or all? 
    // Actually fetchRecords fetches page 1 by default. This is a limitation of current impl.
    // For now, we'll assume we display what we fetched.)
    
    let res = records.value
    if (activeFilterTagId.value) {
       res = res.filter(r => {
        let rTags: string[] = []
        if (typeof r.tags === 'string') {
          try { rTags = JSON.parse(r.tags) } catch {}
        } else {
          rTags = r.tags
        }
        return rTags.map(String).includes(String(activeFilterTagId.value))
      })
    }
    return res
  })

  // --- Actions ---
  async function fetchTags() {
    const res = await http.get('/little-growth/tags')
    tags.value = (res || []).map((t: any) => ({
      ...t,
      id: String(t.id),
      parentId: t.parent_id != null ? String(t.parent_id) : undefined
    }))
  }

  async function createTag(name: string, color?: string, parentId?: string) {
    const payload: any = { name, color }
    if (parentId) payload.parent_id = Number(parentId)
    const res = await http.post('/little-growth/tags', payload)
    const newTag = { ...res, id: String(res.id) }
    tags.value.push(newTag)
    return newTag
  }

  async function updateTag(id: string, name: string, color?: string, parentId?: string) {
    const payload: any = { name, color }
    if (parentId !== undefined) payload.parent_id = parentId ? Number(parentId) : null
    const res = await http.put(`/little-growth/tags/${id}`, payload)
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

  async function fetchRecords(params: { is_favorite?: boolean; silent?: boolean; page_size?: number } = {}) {
    if (!params.silent) loading.value = true
    try {
      if (!params.silent) await fetchTags() // Ensure tags are loaded
      
      // Construct query params
      const query: any = { page_size: Number(params.page_size || 100) }
      if (params.is_favorite) query.is_favorite = 'true'
      
      const res = await http.get('/little-growth/records', { params: query })
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
          id: String(r.id),
          date: r.date,
          images,
          audio,
          tags: rTags,
          comments: r.comments || []
        }
      })
      records.value = items
    } finally {
      if (!params.silent) loading.value = false
    }
  }

  async function fetchRecord(id: string) {
    const r = await http.get(`/little-growth/records/${id}`)
    // Process record similar to fetchRecords
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
    rTags = rTags.map(String)

    return {
        ...r,
        id: String(r.id),
        date: r.date,
        images,
        audio,
        tags: rTags,
        comments: r.comments || []
    } as GrowthRecord
  }

  async function createRecord(data: any) {
    // data contains images URLs, content, date, tags (IDs)
    await http.post('/little-growth/records', data)
    await fetchRecords({ is_favorite: onlyFavorites.value })
  }

  async function updateRecord(id: string, data: any) {
    await http.put(`/little-growth/records/${id}`, data)
    await fetchRecords({ is_favorite: onlyFavorites.value })
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
    await fetchRecords({ is_favorite: onlyFavorites.value })
  }
  
  async function toggleFavorite(id: string) {
    const res = await http.patch(`/little-growth/records/${id}/favorite`)
    const idx = records.value.findIndex(r => String(r.id) === String(id))
    if (idx !== -1) {
        records.value[idx].is_favorite = res.is_favorite
    }
    // If we are viewing only favorites and we unfavorite, we should probably remove it from list or re-fetch
    if (onlyFavorites.value && !res.is_favorite) {
        records.value = records.value.filter(r => String(r.id) !== String(id))
    }
  }
  
  async function addComment(id: string, content: string) {
    const res = await http.post(`/little-growth/records/${id}/comments`, { content })
    const idx = records.value.findIndex(r => String(r.id) === String(id))
    if (idx !== -1) {
        if (!records.value[idx].comments) records.value[idx].comments = []
        records.value[idx].comments!.push(res)
    }
  }

  function getRecordById(id: string) {
    return records.value.find(r => String(r.id) === String(id))
  }

  return {
    records,
    tags,
    activeFilterTagId,
    onlyFavorites,
    filteredRecords,
    flattenedTags,
    loading,
    fetchRecords,
    fetchTags,
    createTag,
    updateTag,
    deleteTag,
    createRecord,
    updateRecord,
    deleteRecord,
    togglePin,
    toggleFavorite,
    addComment,
    getRecordById,
    fetchRecord
  }
})
