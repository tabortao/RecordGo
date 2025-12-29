import http from '@/services/http'

export interface Course {
  id: number
  user_id?: number
  name: string
  color: string
}

export interface TimetableConfig {
  id: number
  user_id: number
  show_saturday: boolean
  show_sunday: boolean
  current_grade: string
  current_semester: string
  period_settings_json?: string
  background_emojis?: string
}

export interface PeriodSetting {
  period: number
  start_time: string
  end_time: string
}

export interface TimetableItem {
  id: number
  user_id: number
  grade: string
  semester: string
  day_of_week: number // 1-7
  period: number
  course_id: number
  course?: Course
}

export const timetableApi = {
  // 获取配置
  getConfig(userId?: number) {
    return http.get<TimetableConfig>('/timetable/config', { params: { user_id: userId } })
  },
  // 更新配置
  updateConfig(data: Partial<TimetableConfig>) {
    return http.put<TimetableConfig>('/timetable/config', data)
  },
  // 获取课程列表
  getCourses(userId?: number) {
    return http.get<Course[]>('/timetable/courses', { params: { user_id: userId } })
  },
  // 添加自定义课程
  addCourse(data: { name: string; color: string; user_id?: number }) {
    return http.post<Course>('/timetable/courses', data)
  },
  // 更新自定义课程
  updateCourse(id: number, data: { name: string; color: string }) {
    return http.put<Course>(`/timetable/courses/${id}`, data)
  },
  // 删除自定义课程
  deleteCourse(id: number) {
    return http.delete(`/timetable/courses/${id}`)
  },
  // 获取课表数据
  getTimetable(grade: string, semester: string, userId?: number) {
    return http.get<TimetableItem[]>('/timetable/data', { params: { grade, semester, user_id: userId } })
  },
  // 保存课表数据
  saveTimetable(data: { user_id?: number; grade: string; semester: string; items: Partial<TimetableItem>[] }) {
    return http.post('/timetable/data', data)
  }
}
