import http from './http'

export interface WordBank {
  id: number
  user_id: number
  title: string
  education_stage: string
  subject: string
  version: string
  grade: string
  content: string // JSON string or plain text
  created_at: string
}

export interface DictationSettings {
  user_id: number
  split_rule: 'space' | 'newline'
  play_mode: 'read' | 'dictate'
  order_mode: 'sequential' | 'random'
  repeat_count: number
  interval_seconds: number
  voice_type: string
  zh_voice_type?: string
  en_voice_type?: string
  default_education_stage?: string
  default_version?: string
  default_grade?: string
  speed: number
}

export interface Mistake {
  id: number
  word: string
  context: string
  created_at: string
}

export const dictationApi = {
  // Word Banks
  listWordBanks: () => http.get<WordBank[]>('/dictation/wordbanks'),
  createWordBank: (data: Partial<WordBank>) => http.post<WordBank>('/dictation/wordbanks', data),
  updateWordBank: (id: number, data: Partial<WordBank>) => http.put<WordBank>(`/dictation/wordbanks/${id}`, data),
  deleteWordBank: (id: number) => http.delete(`/dictation/wordbanks/${id}`),

  // Settings
  getSettings: () => http.get<DictationSettings>('/dictation/settings'),
  updateSettings: (data: Partial<DictationSettings>) => http.put<DictationSettings>('/dictation/settings', data),

  // Mistakes
  listMistakes: () => http.get<Mistake[]>('/dictation/mistakes'),
  addMistake: (data: Partial<Mistake>) => http.post<Mistake>('/dictation/mistakes', data),
  deleteMistake: (id: number) => http.delete(`/dictation/mistakes/${id}`),
  
  // History
  addHistory: (data: any) => http.post('/dictation/history', data),
  getStats: () => http.get('/dictation/stats')
}
