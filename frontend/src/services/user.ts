// 中文注释：用户资料与安全相关服务封装
import http from './http'
import { presignUpload, putToURL } from './storage'
import { prepareUpload } from '@/utils/image'

// 中文注释：获取当前用户最新资料
export async function getProfile() {
  return await http.get('/user/profile') as any
}

// 中文注释：更新昵称（与后端字段保持一致）
export async function updateNickname(nickname: string) {
  // 中文注释：保留旧函数，实际委托到 updateProfile
  return await updateProfile({ nickname })
}

// 中文注释：更新用户资料（昵称/电话/邮箱），字段可选
export async function updateProfile(payload: { nickname?: string, phone?: string, email?: string, child_birthday?: string, child_gender?: string }) {
  return await http.put('/user/profile', payload) as any
}

// 中文注释：修改密码
export async function changePassword(old_password: string, new_password: string) {
  return await http.post('/auth/change-password', { old_password, new_password }) as any
}

// 中文注释：上传用户头像（前端需先压缩并转换为 webp）
export async function uploadAvatar(userId: number, file: File) {
  const webp = await prepareUpload(file)
  const sign = await presignUpload({ resource_type: 'avatar', user_id: userId, content_type: 'image/webp', ext: 'webp' })
  await putToURL(sign.upload_url, webp, sign.headers)
  const saved = await http.post('/upload/avatar/object', { object_key: sign.object_key } as any)
  return saved as { path: string }
}
