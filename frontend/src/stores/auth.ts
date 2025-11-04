// 中文注释：认证状态管理（Pinia），统一维护 token 与用户信息
import { defineStore } from 'pinia'

export interface AuthUser {
  id: number
  username: string
  nickname: string
  role: string
  permissions: string
  parent_id?: number | null
  coins: number
  tomatoes: number
  avatar_path: string
}

export const useAuth = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('auth_token') || '',
    user: ((): AuthUser | null => {
      const s = localStorage.getItem('auth_user')
      try { return s ? JSON.parse(s) as AuthUser : null } catch { return null }
    })()
  }),
  actions: {
    // 中文注释：登录后设置 token 与用户信息，并持久化
    setLogin(token: string, user: AuthUser) {
      this.token = token
      this.user = user
      localStorage.setItem('auth_token', token)
      localStorage.setItem('auth_user', JSON.stringify(user))
    },
    // 中文注释：退出登录，清除持久化并重置状态
    logout() {
      this.token = ''
      this.user = null
      localStorage.removeItem('auth_token')
      localStorage.removeItem('auth_user')
    }
  }
})

