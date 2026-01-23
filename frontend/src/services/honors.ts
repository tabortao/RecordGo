import http from './http'

// 中文注释：荣誉记录类型
export interface HonorRecord {
  id: number
  user_id: number
  title: string
  issuer: string
  awarded_at: string
  photo_url: string
  remark: string
  created_at: string
  updated_at: string
}

// 中文注释：获取荣誉列表
export async function listHonors(): Promise<{ items: HonorRecord[] }> {
  return await http.get('/honors') as any
}

// 中文注释：获取荣誉详情
export async function getHonor(id: number): Promise<HonorRecord> {
  return await http.get(`/honors/${id}`) as any
}

// 中文注释：创建荣誉
export async function createHonor(payload: { title: string; issuer: string; awarded_at: string; photo_url: string; remark: string }): Promise<HonorRecord> {
  return await http.post('/honors', payload) as any
}

// 中文注释：更新荣誉
export async function updateHonor(id: number, payload: Partial<{ title: string; issuer: string; awarded_at: string; photo_url: string; remark: string }>): Promise<HonorRecord> {
  return await http.put(`/honors/${id}`, payload) as any
}

// 中文注释：删除荣誉
export async function deleteHonor(id: number): Promise<any> {
  return await http.delete(`/honors/${id}`) as any
}
