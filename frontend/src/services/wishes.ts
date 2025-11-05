import http from './http'

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
  created_at: string
}

// 中文注释：兑换记录列表响应结构体（与后端返回一致）
export interface WishRecordsResp {
  items: WishRecord[]
  total: number
  page: number
  page_size: number
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
      icon: x.Icon ?? x.icon ?? '',
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
export async function exchangeWish(id: number, userId: number, count = 1) {
  // 中文注释：支持一次兑换多份，后端将按 count 扣减金币与累计次数
  const resp = await http.post(`/wishes/${id}/exchange`, { user_id: userId, count })
  // 中文注释：后端返回 { wish, user_coins, record }，此处明确类型，避免 AxiosResponse 类型干扰
  return resp as { wish: Wish; user_coins: number; record: WishRecord }
}

// 兑换记录
export async function listWishRecords(userId: number, page = 1, pageSize = 10) {
  const resp = await http.get('/wishes/records', { params: { user_id: userId, page, page_size: pageSize } } as any)
  return resp as WishRecordsResp
}

// 上传心愿图标（前端先转换为 webp）
export async function uploadWishIcon(userId: number, file: File) {
  const form = new FormData()
  form.append('user_id', String(userId))
  form.append('image', file)
  const resp = await http.post('/upload/wish-icon', form, { timeout: 30000 } as any)
  // 中文注释：后端返回 { path }，为相对路径 uploads/images/wish/{用户id}/xxx.webp
  return resp as { path: string }
}

// 工具：将图片转换为 webp（质量 0.8），失败则返回原文件
export async function toWebp(file: File, quality = 0.8): Promise<File> {
  return new Promise((resolve) => {
    const img = new Image()
    img.onload = () => {
      const canvas = document.createElement('canvas')
      canvas.width = img.width
      canvas.height = img.height
      const ctx = canvas.getContext('2d')!
      ctx.drawImage(img, 0, 0)
      canvas.toBlob((blob) => {
        if (blob) {
          const webp = new File([blob], file.name.replace(/\.[^.]+$/, '.webp'), { type: 'image/webp' })
          resolve(webp)
        } else {
          resolve(file)
        }
      }, 'image/webp', quality)
    }
    img.onerror = () => resolve(file)
    img.src = URL.createObjectURL(file)
  })
}
