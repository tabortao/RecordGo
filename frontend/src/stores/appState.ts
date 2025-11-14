import { defineStore } from 'pinia'

// 中文注释：应用全局状态，包含权限、番茄钟、金币缓存等，遵循一致性与状态管理要求
export interface TomatoState {
  running: boolean
  mode: 'countdown' | 'countup'
  runningMode?: 'countdown' | 'countup' | null
  durationMinutes: number
  remainingSeconds: number
  fixedTomatoPage: boolean
  showFloating: boolean
  currentTaskId: number | null
  // 中文注释：为支持锁屏/标签页不可见时继续计时，记录开始与结束的时间戳（毫秒）
  startAtMs?: number | null
  endAtMs?: number | null
  // 中文注释：是否启用保持常亮（Wake Lock），移动端默认开启
  keepAwakeEnabled: boolean
}

export interface Permissions {
  view_only: boolean
}

// 中文注释：朗读（TTS）设置，支持语音、语速、音调与开关
export interface SpeechSettings {
  enabled: boolean
  voiceURI?: string | null
  rate: number // 语速（0.5~2.0，默认1）
  pitch: number // 音调（0~2，默认1）
}

export interface AppState {
  coins: number
  // 中文注释：心愿兑换扣除的金币累计，用于总金币计算（主/子账号共享金币池）
  wishDeductedCoins: number
  // 中文注释：最近金币变动历史（最多保留20条），记录每次变更的增量与时间
  coinsHistory: { delta: number; newTotal: number; at: string }[]
  permissions: Permissions
  tomato: TomatoState
  speech: SpeechSettings
  // 中文注释：任务备注入口开关（默认开启，关闭后任务卡片与菜单不显示备注）
  taskNotesEnabled: boolean
  // 中文注释：任务自动排序开关（默认开启）。开启后：
  // 1）分类内已完成任务排在下方；2）全部完成的分类排在未完分类之后。
  taskAutoSortEnabled: boolean
  // 中文注释：主题外观模式（system/dark/light），用于控制深浅色
  themeMode: 'system' | 'dark' | 'light'
}

const DEFAULT_STATE: AppState = {
  coins: 0,
  wishDeductedCoins: 0,
  coinsHistory: [],
  permissions: { view_only: false },
  tomato: {
    running: false,
    mode: 'countdown',
    runningMode: null,
    durationMinutes: 20,
    remainingSeconds: 20 * 60,
    fixedTomatoPage: true,
    showFloating: false,
    currentTaskId: null,
    startAtMs: null,
    endAtMs: null,
    keepAwakeEnabled: true
  },
  speech: {
    enabled: true,
    voiceURI: null,
    rate: 1,
    pitch: 1
  },
  // 中文注释：任务备注入口默认开启
  taskNotesEnabled: true
  ,
  // 中文注释：任务自动排序默认开启
  taskAutoSortEnabled: true,
  // 中文注释：主题外观默认跟随系统
  themeMode: 'system'
}

export const useAppState = defineStore('appState', {
  state: (): AppState => {
    // 中文注释：从本地缓存恢复状态；合并默认值，保证新增字段有默认值
    const raw = localStorage.getItem('appState')
    if (!raw) return DEFAULT_STATE
    try {
      const persisted = JSON.parse(raw) as Partial<AppState>
      return {
        ...DEFAULT_STATE,
        ...persisted,
        tomato: { ...DEFAULT_STATE.tomato, ...(persisted.tomato || {}) },
        speech: { ...DEFAULT_STATE.speech, ...(persisted.speech || {}) }
      }
    } catch {
      return DEFAULT_STATE
    }
  },
  actions: {
    setCoins(v: number) {
      // 中文注释：记录金币变动历史（增量=新值-旧值），时间为 ISO 字符串
      const prev = Number(this.coins || 0)
      const next = Number(v || 0)
      this.coins = next
      const delta = next - prev
      try {
        const record = { delta, newTotal: next, at: new Date().toISOString() }
        const arr = Array.isArray(this.coinsHistory) ? [...this.coinsHistory] : []
        arr.unshift(record)
        this.coinsHistory = arr.slice(0, 20)
      } catch { this.coinsHistory = [] }
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
    // 中文注释：更新朗读设置（语音、语速、音调、开关），并持久化
    updateSpeech(partial: Partial<SpeechSettings>) {
      this.speech = { ...this.speech, ...partial }
      this.persist()
    },
    // 中文注释：开启/关闭任务备注入口，并持久化
    setTaskNotesEnabled(enabled: boolean) {
      this.taskNotesEnabled = enabled
      this.persist()
    },
    // 中文注释：开启/关闭任务自动排序，并持久化
    setTaskAutoSortEnabled(enabled: boolean) {
      this.taskAutoSortEnabled = enabled
      this.persist()
    },
    // 中文注释：设置主题外观（系统/深色/浅色），并立即应用到页面
    setThemeMode(mode: 'system'|'dark'|'light') {
      this.themeMode = mode
      try {
        const isDark = mode === 'dark' ? true : mode === 'light' ? false : (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches)
        document.documentElement.classList.toggle('dark', !!isDark)
      } catch {}
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
