import { useAppState } from '@/stores/appState'
import { createCustomAudio } from './customTTS'

let currentAudio: HTMLAudioElement | null = null

// 中文注释：Web Speech API 封装，提供语音列表与朗读函数

export interface SpeakOptions {
  voiceURI?: string | null
  rate?: number
  pitch?: number
  lang?: string
  onEnd?: () => void
  onError?: (err: any) => void
}

// 中文注释：获取可用语音列表（可能需要等待 voiceschanged 事件）
export function listVoices(): SpeechSynthesisVoice[] {
  try {
    return window.speechSynthesis?.getVoices?.() || []
  } catch {
    return []
  }
}

// 中文注释：根据 URI 查找语音对象，找不到则返回 undefined
export function findVoiceByURI(uri?: string | null): SpeechSynthesisVoice | undefined {
  if (!uri) return undefined
  const voices = listVoices()
  return voices.find(v => v.voiceURI === uri)
}

// 中文注释：选择一个合适的中文语音（首选 zh-CN），若无则返回 undefined
export function pickChineseVoice(): SpeechSynthesisVoice | undefined {
  const voices = listVoices()
  const zh = voices.find(v => /^zh(?:-|_)?/i.test(v.lang))
  return zh || voices.find(v => /Chinese|中文/i.test(v.name))
}

// 中文注释：朗读文本；支持指定语音、语速与音调；异常时返回 false
export async function speak(text: string, opts: SpeakOptions = {}): Promise<boolean> {
  stop() // Stop previous

  try {
    // 尝试使用自定义 TTS
    try {
        const store = useAppState()
        if (store.speech.engine === 'custom' && store.speech.activeCustomId) {
            const profile = store.speech.customProfiles.find(p => p.id === store.speech.activeCustomId)
            if (profile) {
                const audio = await createCustomAudio(text, profile)
                if (audio) {
                    currentAudio = audio
                    audio.onended = () => {
                      opts.onEnd?.()
                      if (currentAudio === audio) currentAudio = null
                    }
                    audio.onerror = (e) => {
                      opts.onError?.(e)
                      if (currentAudio === audio) currentAudio = null
                    }
                    await audio.play()
                    return true
                }
            }
        }
    } catch (e) {
        // Pinia 可能未初始化（极少情况），忽略并降级
    }

    if (!('speechSynthesis' in window)) return false
    const utter = new SpeechSynthesisUtterance(text)
    // 语音
    const chosen = opts.voiceURI ? findVoiceByURI(opts.voiceURI) : undefined
    const fallback = !chosen ? pickChineseVoice() : undefined
    // 中文注释：Element 类型为 SpeechSynthesisVoice | null，空时置为 null
    utter.voice = chosen ?? fallback ?? null
    // 语言（若选择中文语音则自动设为 zh-CN）
    utter.lang = opts.lang || utter.voice?.lang || 'zh-CN'
    // 语速与音调
    utter.rate = typeof opts.rate === 'number' ? Math.min(2, Math.max(0.5, opts.rate)) : 1
    utter.pitch = typeof opts.pitch === 'number' ? Math.min(2, Math.max(0, opts.pitch)) : 1
    
    utter.onend = () => opts.onEnd?.()
    utter.onerror = (e) => opts.onError?.(e)
    
    window.speechSynthesis.speak(utter)
    return true
  } catch {
    return false
  }
}

// 中文注释：停止朗读
export function stop() {
  if (currentAudio) {
    currentAudio.pause()
    currentAudio = null
  }
  try { window.speechSynthesis?.cancel?.() } catch {}
}