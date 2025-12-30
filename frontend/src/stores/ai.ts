import { defineStore } from 'pinia'

export interface ProviderConfig {
  apiKey: string
  apiBaseUrl: string
  model: string
  visionApiKey: string
  visionApiBaseUrl: string
  visionModel: string
}

export interface OCRConfig {
  apiUrl: string
  token: string
}

export interface OCRSettings {
  activeModel: 'PP-OCRv5' | 'PaddleOCR-VL' | 'PP-StructureV3'
  configs: {
    'PP-OCRv5': OCRConfig
    'PaddleOCR-VL': OCRConfig
    'PP-StructureV3': OCRConfig
  }
}

export interface AISettings {
  apiBaseUrl: string
  apiKey: string
  model: string
  visionApiBaseUrl: string
  visionApiKey: string
  visionModel: string
  activeProvider: string
  savedProviderConfigs: Record<string, ProviderConfig>
  ocr: OCRSettings
}

export const PROVIDERS = {
  GOOGLE: {
    name: 'Google Gemini',
    url: '',
    model: 'gemini-2.5-flash',
    visionModel: 'gemini-2.5-flash'
  },
  ZHIPU: {
    name: '智谱 AI (Zhipu)',
    url: 'https://open.bigmodel.cn/api/paas/v4/',
    model: 'glm-4.5-flash',
    visionModel: 'glm-4.6v-flash'
  },
  DEEPSEEK: {
    name: 'DeepSeek (Official)',
    url: 'https://api.deepseek.com',
    model: 'deepseek-chat',
    visionModel: ''
  },
  SILICON: {
    name: 'SiliconFlow (硅基流动)',
    url: 'https://api.siliconflow.cn/v1',
    model: 'deepseek-ai/DeepSeek-V3',
    visionModel: 'zai-org/GLM-4.6V'
  },
  CUSTOM: {
    name: '自定义 / OpenAI 兼容',
    url: '',
    model: '',
    visionModel: ''
  }
}

const DEFAULT_STATE: AISettings = {
  apiBaseUrl: '',
  apiKey: '',
  model: 'gemini-2.5-flash',
  visionApiBaseUrl: '',
  visionApiKey: '',
  visionModel: 'gemini-2.5-flash',
  activeProvider: 'GOOGLE',
  savedProviderConfigs: {},
  ocr: {
    activeModel: 'PP-OCRv5',
    configs: {
      'PP-OCRv5': { apiUrl: '', token: '' },
      'PaddleOCR-VL': { apiUrl: '', token: '' },
      'PP-StructureV3': { apiUrl: '', token: '' }
    }
  }
}

export const useAIStore = defineStore('ai', {
  state: (): AISettings => {
    const raw = localStorage.getItem('aiSettings')
    if (!raw) return DEFAULT_STATE
    try {
      const persisted = JSON.parse(raw)
      return { ...DEFAULT_STATE, ...persisted }
    } catch {
      return DEFAULT_STATE
    }
  },
  actions: {
    updateSettings(partial: Partial<AISettings>) {
      this.$patch(partial)
      this.persist()
    },
    updateProviderConfig(provider: string, config: ProviderConfig) {
      this.savedProviderConfigs[provider] = config
      this.persist()
    },
    persist() {
      localStorage.setItem('aiSettings', JSON.stringify(this.$state))
    },
    switchProvider(newProviderKey: string) {
      // Save current to cache
      const currentCache: ProviderConfig = {
        apiKey: this.apiKey,
        apiBaseUrl: this.apiBaseUrl,
        model: this.model,
        visionApiKey: this.visionApiKey,
        visionApiBaseUrl: this.visionApiBaseUrl,
        visionModel: this.visionModel
      }
      this.savedProviderConfigs[this.activeProvider] = currentCache

      // Load new from cache or defaults
      const cached = this.savedProviderConfigs[newProviderKey]
      // @ts-expect-error PROVIDERS index signature mismatch
      const defaultPreset = PROVIDERS[newProviderKey]

      let newState: Partial<AISettings> = {
        activeProvider: newProviderKey
      }

      if (cached && (cached.apiKey || cached.apiBaseUrl)) {
        newState = {
          ...newState,
          apiKey: cached.apiKey || '',
          apiBaseUrl: cached.apiBaseUrl || defaultPreset.url,
          model: cached.model || defaultPreset.model,
          visionApiKey: cached.visionApiKey || '',
          visionApiBaseUrl: cached.visionApiBaseUrl || '',
          visionModel: cached.visionModel || defaultPreset.visionModel || ''
        }
      } else {
        newState = {
          ...newState,
          apiKey: '',
          apiBaseUrl: defaultPreset.url,
          model: defaultPreset.model,
          visionApiKey: '',
          visionApiBaseUrl: defaultPreset.url,
          visionModel: defaultPreset.visionModel || ''
        }
        if (newProviderKey === 'GOOGLE') {
            newState.visionApiBaseUrl = ''
        }
      }
      
      this.updateSettings(newState)
    }
  }
})
