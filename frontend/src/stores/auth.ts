// 中文注释：认证状态管理（Pinia），统一维护 token 与用户信息
import { defineStore } from 'pinia'
import { getProfile } from '@/services/user'

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
  phone?: string
  email?: string
  child_birthday?: string
  child_gender?: string
  must_change_password?: boolean
  // VIP 字段（由后端返回）
  last_login_time?: string | null
  is_vip?: boolean
  vip_expire_time?: string | null
  is_lifetime_vip?: boolean
}

export const useAuth = defineStore('auth', {
  state: () => ({
    // 中文注释：初始化时优先读取会话（不记住我），否则读取本地（记住我）
    token: ((): string => {
      const st = sessionStorage.getItem('auth_token')
      if (st) return st
      return localStorage.getItem('auth_token') || ''
    })(),
    user: ((): AuthUser | null => {
      // 中文注释：与 token 同步读取来源，避免来源不一致
      const su = sessionStorage.getItem('auth_user')
      const lu = localStorage.getItem('auth_user')
      const s = su || lu || ''
      try { return s ? JSON.parse(s) as AuthUser : null } catch { return null }
    })(),
    accounts: ((): { token: string; user: AuthUser; remember: boolean }[] => {
      try {
        const raw = localStorage.getItem('auth_accounts') || '[]'
        const arr = JSON.parse(raw)
        return Array.isArray(arr) ? arr.filter((a: any) => a && a.user && a.token) : []
      } catch { return [] }
    })()
  }),
  actions: {
    // 中文注释：登录后设置 token 与用户信息；remember=true 写入 localStorage，false 写入 sessionStorage
    setLogin(token: string, user: AuthUser, remember: boolean = true) {
      this.token = token
      this.user = user
      try {
        if (remember) {
          // 中文注释：记住我 → 本地持久化
          localStorage.setItem('auth_token', token)
          localStorage.setItem('auth_user', JSON.stringify(user))
          // 中文注释：清理会话存储，避免来源冲突
          sessionStorage.removeItem('auth_token')
          sessionStorage.removeItem('auth_user')
        } else {
          // 中文注释：不记住我 → 仅本次会话
          sessionStorage.setItem('auth_token', token)
          sessionStorage.setItem('auth_user', JSON.stringify(user))
          // 中文注释：清理本地存储，确保退出或刷新后需重新登录
          localStorage.removeItem('auth_token')
          localStorage.removeItem('auth_user')
        }
      } catch {}
      try {
        const idx = this.accounts.findIndex(a => a.user.id === user.id)
        const acc = { token, user, remember }
        if (idx >= 0) this.accounts.splice(idx, 1, acc)
        else this.accounts.push(acc)
        localStorage.setItem('auth_accounts', JSON.stringify(this.accounts))
      } catch {}
    },
    // 中文注释：退出登录，清除持久化并重置状态
    logout() {
      this.token = ''
      this.user = null
      try {
        localStorage.removeItem('auth_token')
        localStorage.removeItem('auth_user')
        sessionStorage.removeItem('auth_token')
        sessionStorage.removeItem('auth_user')
      } catch {}
    }
    ,
    switchAccount(userId: number) {
      const acc = this.accounts.find(a => a.user.id === userId)
      if (!acc) return
      this.setLogin(acc.token, acc.user, acc.remember)
    },
    removeAccount(userId: number) {
      this.accounts = this.accounts.filter(a => a.user.id !== userId)
      try { localStorage.setItem('auth_accounts', JSON.stringify(this.accounts)) } catch {}
    },
    // 中文注释：更新用户的部分字段（如昵称、头像路径等），并持久化
    updateUser(partial: Partial<AuthUser>) {
      if (!this.user) return
      this.user = { ...this.user, ...partial }
      try {
        // 中文注释：根据当前来源同步更新，两侧都写入以提升兼容性
        localStorage.setItem('auth_user', JSON.stringify(this.user))
        sessionStorage.setItem('auth_user', JSON.stringify(this.user))
      } catch {}
    },
    // 中文注释：刷新当前用户信息
    async refreshUser() {
      if (!this.token) return
      try {
        const u = await getProfile()
        if (u && u.id) {
          this.updateUser(u)
        }
      } catch (e) {
        // console.error('refresh user failed', e)
      }
    }
  }
})
