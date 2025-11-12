// 中文注释：子账号服务封装，统一对接后端 /api/subaccounts 与令牌生成接口
import http from './http'

// 中文注释：子账号响应结构体（与后端 handlers.childResp 一致）
export interface ChildAccount {
  id: number
  nickname: string
  avatar_path: string
  permissions: string
  login_token?: string
}

// 列出子账号
export async function listSubAccounts(): Promise<ChildAccount[]> {
  return await http.get('/subaccounts') as any
}

// 创建子账号（用户名后端自动生成）
export async function createSubAccount(payload: { nickname: string; permissions: string; avatar_path?: string }): Promise<ChildAccount> {
  return await http.post('/subaccounts', payload) as any
}

// 更新子账号
export async function updateSubAccount(id: number, payload: Partial<{ nickname: string; permissions: string; avatar_path: string }>): Promise<ChildAccount> {
  return await http.put(`/subaccounts/${id}`, payload) as any
}

// 删除子账号
export async function deleteSubAccount(id: number): Promise<{ id: number }> {
  return await http.delete(`/subaccounts/${id}`) as any
}

// 生成或刷新子账号登录令牌
// 生成或刷新子账号登录令牌（支持可选的有效期秒数）
export async function generateChildToken(id: number, opts?: { expiresInSeconds?: number }): Promise<{ token: string; expires_at?: string }> {
  const payload = opts?.expiresInSeconds != null ? { expires_in_seconds: opts.expiresInSeconds } : {}
  return await http.post(`/subaccounts/${id}/token`, payload) as any
}