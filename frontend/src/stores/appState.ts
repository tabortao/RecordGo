import { defineStore } from 'pinia'

// 中文注释：应用全局状态，包含权限、番茄钟、金币缓存等，遵循一致性与状态管理要求
export interface TomatoState {
  running: boolean
  mode: 'countdown' | 'countup'
  durationMinutes: number
  remainingSeconds: number
  fixedTomatoPage: boolean
}

export interface Permissions {
  view_only: boolean
}

export interface AppState {
  coins: number
  // 中文注释：心愿兑换扣除的金币累计，用于总金币计算（主/子账号共享金币池）
  wishDeductedCoins: number
  permissions: Permissions
  tomato: TomatoState
}

const DEFAULT_STATE: AppState = {
  coins: 0,
  wishDeductedCoins: 0,
  permissions: { view_only: false },
  tomato: {
    running: false,
    mode: 'countdown',
    durationMinutes: 20,
    remainingSeconds: 20 * 60,
    fixedTomatoPage: false
  }
}

export const useAppState = defineStore('appState', {
  state: (): AppState => {
    // 中文注释：从本地缓存恢复状态，保证刷新后状态持久化
    const raw = localStorage.getItem('appState')
    return raw ? JSON.parse(raw) as AppState : DEFAULT_STATE
  },
  actions: {
    setCoins(v: number) {
      this.coins = v
      this.persist()
    },
    setWishDeducted(v: number) {
      // 中文注释：设置心愿兑换扣减金币（累计值）
      this.wishDeductedCoins = v
      this.persist()
    },
    setPermissions(p: Permissions) {
      this.permissions = p
      this.persist()
    },
    updateTomato(partial: Partial<TomatoState>) {
      this.tomato = { ...this.tomato, ...partial }
      this.persist()
    },
    reset() {
      Object.assign(this, DEFAULT_STATE)
      this.persist()
    },
    persist() {
      // 中文注释：持久化到本地缓存
      localStorage.setItem('appState', JSON.stringify(this.$state))
    }
  }
})
