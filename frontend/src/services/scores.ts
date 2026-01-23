import http from './http'

export interface ScoreRecord {
  id: number
  user_id: number
  subject: string
  exam_name: string
  exam_type: string
  grade: string
  term: string
  topic: string
  difficulty: string
  exam_date: string
  score: number
  full_score: number
  class_rank: number | null
  rank_type: string
  grade_rank: number | null
  class_avg: number | null
  class_highest: number | null
  photo_url: string
  remark: string
  created_at: string
  updated_at: string
}

export async function listScores(params?: { subject?: string; start_date?: string; end_date?: string }): Promise<{ items: ScoreRecord[] }> {
  return await http.get('/scores', { params }) as any
}

export async function getScore(id: number): Promise<ScoreRecord> {
  return await http.get(`/scores/${id}`) as any
}

export async function createScore(payload: {
  subject: string
  exam_name: string
  exam_type?: string
  grade?: string
  term?: string
  topic?: string
  difficulty?: string
  exam_date: string
  score: number
  full_score: number
  class_rank?: number | null
  rank_type?: string
  grade_rank?: number | null
  class_avg?: number | null
  class_highest?: number | null
  photo_url?: string
  remark?: string
}): Promise<ScoreRecord> {
  return await http.post('/scores', payload) as any
}

export async function updateScore(id: number, payload: Partial<{
  subject: string
  exam_name: string
  exam_type: string
  grade: string
  term: string
  topic: string
  difficulty: string
  exam_date: string
  score: number
  full_score: number
  class_rank: number | null
  rank_type: string
  grade_rank: number | null
  class_avg: number | null
  class_highest: number | null
  photo_url: string
  remark: string
}>): Promise<ScoreRecord> {
  return await http.put(`/scores/${id}`, payload) as any
}

export async function deleteScore(id: number): Promise<any> {
  return await http.delete(`/scores/${id}`) as any
}
