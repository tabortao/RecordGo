// 中文注释：用户资料与安全相关服务封装
import http from './http'

// 中文注释：更新昵称（与后端字段保持一致）
export async function updateNickname(nickname: string) {
  return await http.put('/user/profile', { nickname }) as any
}

// 中文注释：修改密码
export async function changePassword(old_password: string, new_password: string) {
  return await http.post('/auth/change-password', { old_password, new_password }) as any
}

// 中文注释：上传用户头像（前端需先压缩并转换为 webp）
export async function uploadAvatar(userId: number, file: File) {
  const form = new FormData()
  form.append('user_id', String(userId))
  form.append('image', file)
  // 后端应返回 { path } 为相对路径：uploads/images/avatars/{用户id}/xxx.webp
  return await http.post('/upload/avatar', form, { timeout: 30000 } as any) as { path: string }
}