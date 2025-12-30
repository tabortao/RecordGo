import { describe, it, expect, vi, beforeEach } from 'vitest'
import { createCustomAudio } from '../utils/customTTS'
import type { CustomTTSProfile } from '../stores/appState'

// Mock Browser APIs
global.URL.createObjectURL = vi.fn(() => 'blob:url')
global.Audio = vi.fn(() => ({
  play: vi.fn(),
  pause: vi.fn(),
  src: '',
})) as any

describe('createCustomAudio', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    global.fetch = vi.fn()
  })

  it('should create audio from valid GET profile', async () => {
    const profile: CustomTTSProfile = {
      id: '1',
      name: 'Test',
      apiUrl: 'https://api.test/tts',
      apiKey: 'key',
      voiceId: 'voice-1',
      speed: 1.0,
      pitch: 1.0,
      method: 'GET'
    }

    const mockBlob = { size: 100 }
    ;(global.fetch as any).mockResolvedValue({
      ok: true,
      blob: () => Promise.resolve(mockBlob)
    })

    const audio = await createCustomAudio('hello', profile)
    expect(audio).toBeTruthy()
    expect(global.fetch).toHaveBeenCalledWith(expect.stringContaining('https://api.test/tts'), expect.any(Object))
    
    // Check params
    const callUrl = (global.fetch as any).mock.calls[0][0]
    const url = new URL(callUrl)
    expect(url.searchParams.get('text')).toBe('hello')
    expect(url.searchParams.get('voice')).toBe('voice-1')
    
    // Check Headers
    const opts = (global.fetch as any).mock.calls[0][1]
    expect(opts.headers['Authorization']).toBe('Bearer key')
  })

  it('should create audio from valid POST profile', async () => {
    const profile: CustomTTSProfile = {
      id: '2',
      name: 'OpenAI',
      apiUrl: 'https://api.openai.com/v1/audio/speech',
      apiKey: 'sk-xxx',
      voiceId: 'alloy',
      speed: 1.0,
      pitch: 1.0,
      method: 'POST'
    }

    const mockBlob = { size: 100 }
    ;(global.fetch as any).mockResolvedValue({
      ok: true,
      blob: () => Promise.resolve(mockBlob)
    })

    const audio = await createCustomAudio('hello', profile)
    expect(audio).toBeTruthy()
    
    const [url, opts] = (global.fetch as any).mock.calls[0]
    expect(url).toBe('https://api.openai.com/v1/audio/speech')
    expect(opts.method).toBe('POST')
    expect(JSON.parse(opts.body)).toMatchObject({
        input: 'hello',
        voice: 'alloy',
        model: 'tts-1'
    })
    expect(opts.headers['Content-Type']).toBe('application/json')
  })

  it('should return null on fetch error', async () => {
    const profile: CustomTTSProfile = {
      id: '3',
      name: 'Fail',
      apiUrl: 'https://fail.test',
      apiKey: '',
      voiceId: '',
      speed: 1.0,
      pitch: 1.0
    }
    ;(global.fetch as any).mockResolvedValue({ ok: false })
    const audio = await createCustomAudio('hello', profile)
    expect(audio).toBeNull()
  })

  it('should return null on empty blob', async () => {
    const profile: CustomTTSProfile = {
      id: '4',
      name: 'Empty',
      apiUrl: 'https://empty.test',
      apiKey: '',
      voiceId: '',
      speed: 1.0,
      pitch: 1.0
    }
    ;(global.fetch as any).mockResolvedValue({
      ok: true,
      blob: () => Promise.resolve({ size: 0 })
    })
    const audio = await createCustomAudio('hello', profile)
    expect(audio).toBeNull()
  })
})
