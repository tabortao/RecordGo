import http from './http'
import { presignUpload, putToURL } from './storage'
import { prepareUpload } from '@/utils/image'

// 中文注释：心愿服务，统一封装心愿相关的 API 调用，保证字段与后端一致
export interface Wish {
  id: number
  user_id: number
  name: string
  content: string
  icon: string
  need_coins: number
  exchange_amount: number
  unit: string
  exchanged: number
  built_in: boolean
}

export interface WishRecord {
  id: number
  user_id: number
  wish_id: number
  wish_name: string
  coins_used: number
  amount: number
  unit: string
  status: string
  // 中文注释：兑换备注（可选），由后端在兑换时保存
  remark?: string
  created_at: string
}

// 中文注释：兑换记录列表响应结构体（与后端返回一致）
export interface WishRecordsResp {
  items: WishRecord[]
  total: number
  page: number
  page_size: number
}

// 中文注释：规范化后端返回的上传路径为相对路径（uploads/...），并统一斜杠
export function normalizeUploadPath(p?: string): string {
  if (!p) return ''
  let s = String(p).trim().replace(/\\/g, '/')
  // 抽取从 uploads 开始的相对路径
  const idx = s.indexOf('uploads/')
  if (idx >= 0) s = s.slice(idx)
  // 去掉 storage/ 前缀
  s = s.replace(/^storage\//, '')
  // 清理多余前导斜杠
  s = s.replace(/^\/+/, '')
  return s
}

// 列表心愿
export async function listWishes(userId: number): Promise<Wish[]> {
  // 中文注释：按用户ID查询心愿列表；后端返回模型字段为大写驼峰（如 NeedCoins），此处统一映射为前端使用的下划线命名，确保页面显示与表单绑定不为空
  const raw = await http.get(`/wishes`, { params: { user_id: userId } } as any)
  // 中文注释：若后端直接返回数组，则进行字段归一化；否则直接返回
  if (Array.isArray(raw)) {
    return raw.map((x: any) => ({
      id: Number(x.ID ?? x.id),
      user_id: Number(x.UserID ?? x.user_id),
      name: x.Name ?? x.name ?? '',
      content: x.Content ?? x.content ?? '',
      icon: (() => {
        const raw = x.Icon ?? x.icon
        if (!raw) return ''
        const str = String(raw)
        // 内置图标文件名（无斜杠）保留原样；否则规范化为 uploads 相对路径
        if (!str.includes('/') && !str.includes('\\')) return str
        return normalizeUploadPath(str)
      })(),
      need_coins: Number(x.NeedCoins ?? x.need_coins ?? 1),
      exchange_amount: Number(x.ExchangeAmount ?? x.exchange_amount ?? 1),
      unit: x.Unit ?? x.unit ?? '次',
      exchanged: Number(x.Exchanged ?? x.exchanged ?? 0),
      built_in: Boolean(x.BuiltIn ?? x.built_in ?? false)
    }))
  }
  return raw as Wish[]
}

// 创建心愿
export async function createWish(payload: Omit<Wish, 'id' | 'exchanged' | 'built_in'> & { user_id: number }) {
  return await http.post('/wishes', payload)
}

// 更新心愿
export async function updateWish(id: number, payload: Partial<Wish>) {
  return await http.put(`/wishes/${id}`, payload)
}

// 删除心愿
export async function deleteWish(id: number) {
  return await http.delete(`/wishes/${id}`)
}

// 兑换心愿
export async function exchangeWish(id: number, userId: number, count = 1, remark?: string) {
  // 中文注释：支持一次兑换多份；可选备注字段 remark 便于记录具体兑换情况（后端可忽略该字段）
  const payload: any = { user_id: userId, count }
  if (remark && remark.trim()) payload.remark = remark.trim()
  const resp = await http.post(`/wishes/${id}/exchange`, payload)
  // 中文注释：后端返回 { wish, user_coins, record }，此处明确类型，避免 AxiosResponse 类型干扰
  return resp as { wish: Wish; user_coins: number; record: WishRecord }
}

// 兑换记录
export async function listWishRecords(userId: number, page = 1, pageSize = 10, opts?: { start?: string; end?: string }) {
  const params: any = { user_id: userId, page, page_size: pageSize }
  if (opts?.start) params.start = opts.start
  if (opts?.end) params.end = opts.end
  const resp = await http.get('/wishes/records', { params } as any)
  // 中文注释：兼容后端可能返回的大小写字段，统一映射为前端使用的下划线命名
  const raw = resp as any
  const items = Array.isArray(raw.items) ? raw.items.map((x: any) => ({
    id: Number(x.id ?? x.ID),
    user_id: Number(x.user_id ?? x.UserID),
    wish_id: Number(x.wish_id ?? x.WishID),
    wish_name: String(x.wish_name ?? x.WishName ?? ''),
    coins_used: Number(x.coins_used ?? x.CoinsUsed ?? 0),
    amount: Number(x.amount ?? x.Amount ?? 0),
    unit: String(x.unit ?? x.Unit ?? ''),
    status: String(x.status ?? x.Status ?? ''),
    remark: String(x.remark ?? x.Remark ?? ''),
    created_at: String(x.created_at ?? x.CreatedAt ?? '')
  })) : []
  return { items, total: Number(raw.total ?? 0), page: Number(raw.page ?? page), page_size: Number(raw.page_size ?? pageSize) } as WishRecordsResp
}

// 上传心愿图标（前端先转换为 webp）
export async function uploadWishIcon(userId: number, file: File) {
  const webp = await prepareUpload(file, 0.8)
  const sign = await presignUpload({ resource_type: 'wish_image', user_id: userId, content_type: 'image/webp', ext: 'webp' })
  await putToURL(sign.upload_url, webp, sign.headers)
  return { path: sign.object_key }
}

// 中文注释：获取单个心愿详情（编辑页面使用）
export async function getWish(id: number): Promise<Wish> {
  const w = await http.get(`/wishes/${id}`)
  // 规范化返回的 icon 字段
  const icon = normalizeUploadPath((w as any)?.Icon ?? (w as any)?.icon)
  return { ...(w as any), icon }
}

// 统一到 utils/image.prepareUpload
