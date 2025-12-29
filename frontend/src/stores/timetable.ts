import { defineStore } from 'pinia'
import { ref } from 'vue'
import { timetableApi, type TimetableConfig, type Course, type TimetableItem } from '@/api/timetable'

export const useTimetableStore = defineStore('timetable', () => {
  const config = ref<TimetableConfig | null>(null)
  const courses = ref<Course[]>([])
  const timetable = ref<TimetableItem[]>([])
  const loading = ref(false)

  const fetchConfig = async (userId?: number) => {
    // 尝试从缓存读取
    const cacheKey = `recordgo_timetable_config_${userId || 'default'}`
    const cached = localStorage.getItem(cacheKey)
    if (cached) {
      try {
        config.value = JSON.parse(cached)
      } catch (e) {
        console.error('Failed to parse cached config', e)
      }
    }

    try {
      const res = await timetableApi.getConfig(userId)
      config.value = res
      // 更新缓存
      localStorage.setItem(cacheKey, JSON.stringify(res))
    } catch (error) {
      console.error('Failed to fetch config', error)
    }
  }

  const fetchCourses = async (userId?: number) => {
    try {
      const res = await timetableApi.getCourses(userId)
      courses.value = res
    } catch (error) {
      console.error('Failed to fetch courses', error)
    }
  }

  const fetchTimetable = async (grade: string, semester: string, userId?: number) => {
    loading.value = true
    
    // 尝试从缓存读取
    const cacheKey = `recordgo_timetable_${grade}_${semester}_${userId || 'default'}`
    const cached = localStorage.getItem(cacheKey)
    if (cached) {
      try {
        timetable.value = JSON.parse(cached)
      } catch (e) {
        console.error('Failed to parse cached timetable', e)
      }
    }

    try {
      const res = await timetableApi.getTimetable(grade, semester, userId)
      timetable.value = res
      // 更新缓存
      localStorage.setItem(cacheKey, JSON.stringify(res))
    } catch (error) {
      console.error('Failed to fetch timetable', error)
    } finally {
      loading.value = false
    }
  }

  // 更新配置并同步缓存
  const updateConfig = async (newConfig: TimetableConfig, userId?: number) => {
      const cacheKey = `recordgo_timetable_config_${userId || 'default'}`
      config.value = newConfig
      localStorage.setItem(cacheKey, JSON.stringify(newConfig))
  }

  return {
    config,
    courses,
    timetable,
    loading,
    fetchConfig,
    fetchCourses,
    fetchTimetable,
    updateConfig
  }
})
