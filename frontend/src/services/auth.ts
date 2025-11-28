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
    // VIP 字段
    last_login_time?: string | null
    is_vip?: boolean
    vip_expire_time?: string | null
    is_lifetime_vip?: boolean
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

// 中文注释：子账号令牌登录（免密码，仅限子账号使用）
export async function apiTokenLogin(token: string) {
  return await http.post<LoginResp>('/auth/token-login', { token })
}
