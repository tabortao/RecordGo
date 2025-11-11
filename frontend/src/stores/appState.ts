import { defineStore } from 'pinia'

// 中文注释：应用全局状态，包含权限、番茄钟、金币缓存等，遵循一致性与状态管理要求
export interface TomatoState {
  running: boolean
  mode: 'countdown' | 'countup'
  durationMinutes: number
  remainingSeconds: number
  fixedTomatoPage: boolean
  showFloating: boolean
  currentTaskId: number | null
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
  permissions: Permissions
  tomato: TomatoState
  speech: SpeechSettings
  // 中文注释：任务备注入口开关（默认开启，关闭后任务卡片与菜单不显示备注）
  taskNotesEnabled: boolean
  // 中文注释：任务自动排序开关（默认开启）。开启后：
  // 1）分类内已完成任务排在下方；2）全部完成的分类排在未完分类之后。
  taskAutoSortEnabled: boolean
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
    fixedTomatoPage: true,
    showFloating: false,
    currentTaskId: null
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
  taskAutoSortEnabled: true
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
