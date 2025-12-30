import type { CustomTTSProfile } from '@/stores/appState'

export async function createCustomAudio(text: string, profile: CustomTTSProfile): Promise<HTMLAudioElement | null> {
  try {
    let rawUrl = profile.apiUrl.trim()
    if (!rawUrl.startsWith('http')) rawUrl = 'https://' + rawUrl

    const method = profile.method || 'GET'
    const headers: HeadersInit = {}
    if (profile.apiKey) headers['Authorization'] = `Bearer ${profile.apiKey}`

    let res: Response
    if (method === 'POST') {
      headers['Content-Type'] = 'application/json'
      const body = JSON.stringify({
        model: 'tts-1',
        input: text,
        voice: profile.voiceId || 'zh-CN-XiaoxiaoNeural',
        speed: profile.speed || 1.0,
        pitch: profile.pitch || 1.0,
        response_format: 'mp3'
      })
      res = await fetch(rawUrl, { method: 'POST', headers, body })
    } else {
      const urlObj = new URL(rawUrl)
      urlObj.searchParams.set('text', text)
      if (profile.voiceId) urlObj.searchParams.set('voice', profile.voiceId)
      
      const ratePercent = Math.round(((profile.speed || 1.0) - 1) * 100)
      urlObj.searchParams.set('rate', (ratePercent >= 0 ? '+' : '') + ratePercent + '%')
      
      const pitchDiff = Math.round(((profile.pitch || 1.0) - 1) * 20)
      urlObj.searchParams.set('pitch', (pitchDiff >= 0 ? '+' : '') + pitchDiff + 'Hz')
      
      res = await fetch(urlObj.toString(), { headers })
    }

    if (!res.ok) return null

    const blob = await res.blob()
    if (blob.size === 0) return null

    return new Audio(URL.createObjectURL(blob))
  } catch (e) {
    console.error('Custom TTS Error:', e)
    return null
  }
}
