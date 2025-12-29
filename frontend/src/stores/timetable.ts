import { defineStore } from 'pinia'
import { ref } from 'vue'
import { timetableApi, type TimetableConfig, type Course, type TimetableItem } from '@/api/timetable'

export const useTimetableStore = defineStore('timetable', () => {
  const config = ref<TimetableConfig | null>(null)
  const courses = ref<Course[]>([])
  const timetable = ref<TimetableItem[]>([])
  const loading = ref(false)

  const fetchConfig = async (userId?: number) => {
    try {
      const res = await timetableApi.getConfig(userId)
      config.value = res
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
    try {
      const res = await timetableApi.getTimetable(grade, semester, userId)
      timetable.value = res
    } catch (error) {
      console.error('Failed to fetch timetable', error)
    } finally {
      loading.value = false
    }
  }

  return {
    config,
    courses,
    timetable,
    loading,
    fetchConfig,
    fetchCourses,
    fetchTimetable
  }
})
