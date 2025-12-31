import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import dayjs from 'dayjs'
import type { Poem, StudyRecord, DailyStats, DictationResult } from '../types'
import rawPoems from '../data/chinese_poems.json'

// Ebbinghaus intervals in days: 1, 2, 4, 7, 15, 30
const REVIEW_INTERVALS = [1, 2, 4, 7, 15, 30]

export const usePoetryStore = defineStore('poetry', () => {
  // State
  const poems = ref<Poem[]>([])
  const studyRecords = ref<Record<number, StudyRecord>>({})
  const dailyStats = ref<Record<string, DailyStats>>({})
  const collections = ref<number[]>([])
  const customPoems = ref<Poem[]>([])
  const dictationRecords = ref<DictationResult[]>([])
  const poemUpdates = ref<Record<number, Partial<Poem>>>({}) // Store AI content for system poems
  const challengeRecords = ref<{ date: number, score: number, mode: string }[]>([])

  // Init
  const init = () => {
    // Load local storage
    const savedRecords = localStorage.getItem('poetry_records')
    if (savedRecords) studyRecords.value = JSON.parse(savedRecords)

    const savedStats = localStorage.getItem('poetry_stats')
    if (savedStats) dailyStats.value = JSON.parse(savedStats)

    const savedCollections = localStorage.getItem('poetry_collections')
    if (savedCollections) collections.value = JSON.parse(savedCollections)
    
    const savedCustom = localStorage.getItem('poetry_custom')
    if (savedCustom) customPoems.value = JSON.parse(savedCustom)

    const savedDictations = localStorage.getItem('poetry_dictations')
    if (savedDictations) dictationRecords.value = JSON.parse(savedDictations)

    const savedUpdates = localStorage.getItem('poetry_updates')
    if (savedUpdates) poemUpdates.value = JSON.parse(savedUpdates)

    const savedChallenges = localStorage.getItem('poetry_challenges')
    if (savedChallenges) challengeRecords.value = JSON.parse(savedChallenges)

    // Merge raw poems with custom poems
    // We need to ensure IDs don't clash. Raw poems use index as ID currently.
    // Let's assume raw poems have IDs 0 to N. Custom poems start from 100000.
    
    const loadedRaw = (rawPoems as any[]).map((p, idx) => {
      const id = p.id || idx
      const updates = poemUpdates.value[id] || {}
      return {
        ...p,
        id, 
        tags: classifyPoem(p),
        ...updates
      }
    })

    poems.value = [...loadedRaw, ...customPoems.value]
  }

  // Simple classification logic (Mock)
  // In reality, this should come from a mapping file as per PRD
  const classifyPoem = (p: any): string[] => {
    const tags: string[] = []
    // Mock logic based on keywords or authors commonly in textbooks
    const primaryAuthors = ['李白', '杜甫', '白居易', '王维', '孟浩然', '骆宾王']
    if (primaryAuthors.includes(p.author_cns) && p.paragraphs_cns.length <= 4) {
      tags.push('primary')
    } else if (p.dynasty_cns === '宋代' || p.dynasty_cns === '宋') {
      tags.push('middle')
    } else {
      tags.push('high')
    }
    return tags
  }

  // Persist helpers
  const saveRecords = () => localStorage.setItem('poetry_records', JSON.stringify(studyRecords.value))
  const saveStats = () => localStorage.setItem('poetry_stats', JSON.stringify(dailyStats.value))
  const saveCollections = () => localStorage.setItem('poetry_collections', JSON.stringify(collections.value))
  const saveCustom = () => localStorage.setItem('poetry_custom', JSON.stringify(customPoems.value))
  const saveDictations = () => localStorage.setItem('poetry_dictations', JSON.stringify(dictationRecords.value))
  const saveUpdates = () => localStorage.setItem('poetry_updates', JSON.stringify(poemUpdates.value))
  const saveChallenges = () => localStorage.setItem('poetry_challenges', JSON.stringify(challengeRecords.value))

  // Actions
  const addCustomPoem = (poem: Omit<Poem, 'id'>) => {
    const newId = 100000 + customPoems.value.length + 1
    const newPoem = { ...poem, id: newId, isCustom: true } as Poem
    customPoems.value.push(newPoem)
    poems.value.push(newPoem)
    saveCustom()
  }
  
  const updateCustomPoem = (id: number, updates: Partial<Poem>) => {
      const idx = customPoems.value.findIndex(p => p.id === id)
      if (idx !== -1) {
          customPoems.value[idx] = { ...customPoems.value[idx], ...updates }
          saveCustom()
          
          // Update in main list
          const mainIdx = poems.value.findIndex(p => p.id === id)
          if (mainIdx !== -1) {
              poems.value[mainIdx] = { ...poems.value[mainIdx], ...updates }
          }
      } else {
        // System poem update (e.g. AI content)
        poemUpdates.value[id] = { ...(poemUpdates.value[id] || {}), ...updates }
        saveUpdates()
        const mainIdx = poems.value.findIndex(p => p.id === id)
        if (mainIdx !== -1) {
            poems.value[mainIdx] = { ...poems.value[mainIdx], ...updates }
        }
      }
  }

  const deleteCustomPoem = (id: number) => {
    customPoems.value = customPoems.value.filter(p => p.id !== id)
    poems.value = poems.value.filter(p => p.id !== id)
    saveCustom()
    // Also clean records
    delete studyRecords.value[id]
    saveRecords()
  }

  const toggleCollection = (id: number) => {
    if (collections.value.includes(id)) {
      collections.value = collections.value.filter(cid => cid !== id)
    } else {
      collections.value.push(id)
    }
    saveCollections()
  }

  const recordStudy = (id: number, isMasteredInput: boolean = false) => {
    const now = Date.now()
    const today = dayjs().format('YYYY-MM-DD')
    
    // Update Daily Stats
    if (!dailyStats.value[today]) {
      dailyStats.value[today] = { date: today, studiedCount: 0, details: [] }
    }
    if (!dailyStats.value[today].details.includes(id)) {
      dailyStats.value[today].studiedCount++
      dailyStats.value[today].details.push(id)
      saveStats()
    }

    // Update Poem Record
    if (!studyRecords.value[id]) {
      studyRecords.value[id] = {
        poemId: id,
        lastStudyTime: now,
        nextReviewTime: now + REVIEW_INTERVALS[0] * 24 * 60 * 60 * 1000,
        stage: 0,
        timesStudied: 1,
        isMastered: isMasteredInput,
        history: [now]
      }
    } else {
      const record = studyRecords.value[id]
      record.lastStudyTime = now
      record.timesStudied++
      record.history.push(now)
      if (isMasteredInput) record.isMastered = true

      // Advance stage logic
      if (record.stage < REVIEW_INTERVALS.length - 1) {
        record.stage++
      }
      record.nextReviewTime = now + REVIEW_INTERVALS[record.stage] * 24 * 60 * 60 * 1000
    }
    saveRecords()
  }

  const saveDictationResult = (result: DictationResult) => {
    dictationRecords.value.push(result)
    saveDictations()
    // Also mark as studied
    recordStudy(result.poemId, result.accuracy >= 90)
  }

  const generateAIContent = async (id: number) => {
    const poem = getPoemById(id)
    if (!poem) return

    // Mock AI generation
    return new Promise<void>((resolve) => {
      setTimeout(() => {
        const aiData = {
          translation: `【AI生成译文】这是关于《${poem.title_cns}》的现代汉语翻译...（此处为模拟生成内容）`,
          appreciation: `【AI生成赏析】这首诗表达了作者${poem.author_cns}在${poem.dynasty_cns}时期的思想感情...（此处为模拟生成内容）`,
          fun_fact: `【AI生成趣事】关于这首诗有一个有趣的传说...（此处为模拟生成内容）`
        }
        updateCustomPoem(id, aiData)
        resolve()
      }, 1500)
    })
  }

  const saveChallengeResult = (score: number, mode: string) => {
    challengeRecords.value.push({ date: Date.now(), score, mode })
    saveChallenges()
  }

  // Getters
  const getPoemById = (id: number) => poems.value.find(p => p.id === id)
  
  const dailyReviewList = computed(() => {
    const now = Date.now()
    return poems.value.filter(p => {
      const record = studyRecords.value[p.id]
      return record && record.nextReviewTime <= now && !record.isMastered
    })
  })
  
  const dailyNewList = computed(() => {
      // Suggest 3 new poems per day from unstudied
      // Filter out studied
      const studiedIds = Object.keys(studyRecords.value).map(Number)
      return poems.value.filter(p => !studiedIds.includes(p.id)).slice(0, 3)
  })

  const getRecord = (id: number) => studyRecords.value[id]
  const getDictationHistory = (id: number) => dictationRecords.value.filter(d => d.poemId === id).sort((a, b) => b.date - a.date)

  const dailyRecommendPoem = computed(() => {
    if (poems.value.length === 0) return null
    // Pick based on date to be consistent for the day
    const dayOfYear = dayjs().dayOfYear()
    const index = dayOfYear % poems.value.length
    return poems.value[index]
  })

  return {
    poems,
    studyRecords,
    dailyStats,
    collections,
    dictationRecords,
    init,
    addCustomPoem,
    updateCustomPoem,
    deleteCustomPoem,
    toggleCollection,
    recordStudy,
    saveDictationResult,
    generateAIContent,
    getPoemById,
    dailyReviewList,
    dailyNewList,
    getRecord,
    getDictationHistory,
    dailyRecommendPoem,
    saveChallengeResult,
    challengeRecords
  }
})
