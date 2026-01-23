import http from './http'

export type GrowthMetricType = 'height' | 'weight' | 'vision' | 'foot'

export interface GrowthMetricRecord {
  id: number
  user_id: number
  metric_type: GrowthMetricType
  record_date: string
  value: number
  created_at: string
  updated_at: string
}

export async function listGrowthRecords(type?: GrowthMetricType) {
  const params: any = {}
  if (type) params.type = type
  return await http.get('/growth/records', { params }) as { items: GrowthMetricRecord[]; total: number }
}

export async function createGrowthRecord(payload: { metric_type: GrowthMetricType; record_date: string; value: number }) {
  return await http.post('/growth/records', payload) as GrowthMetricRecord
}

export async function deleteGrowthRecord(id: number) {
  return await http.delete(`/growth/records/${id}`) as any
}
