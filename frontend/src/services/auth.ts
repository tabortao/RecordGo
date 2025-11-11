// 中文注释：认证服务封装，统一调用后端登录/注册接口
import http from './http'

export interface LoginResp {
  token: string
  user: {
    id: number
    username: string
    nickname: string
    role: string
    permissions: string
    parent_id?: number | null
    coins: number
    tomatoes: number
    avatar_path: string
    phone?: string
    email?: string
  }
}

export async function apiLogin(username: string, password: string) {
  // 中文注释：调用后端登录接口，返回 {token,user}
  return await http.post<LoginResp>('/auth/login', { username, password })
}

export async function apiRegister(username: string, password: string, nickname?: string) {
  // 中文注释：调用后端注册接口
  return await http.post('/auth/register', { username, password, nickname })
}

