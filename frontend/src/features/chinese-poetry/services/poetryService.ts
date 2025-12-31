import { useAIStore } from '@/stores/ai'
import type { Poem } from '../types'
import http from '@/services/http'

export class PoetryService {
  async getPoetryData() {
    return http.get('/poetry/data')
  }

  async savePoetryData(data: string) {
    return http.post('/poetry/data', { data })
  }

  async generateContent(poem: Poem, type: '译文' | '赏析' | '趣事'): Promise<string> {
    const aiStore = useAIStore()
    
    // Check if configuration exists
    if (!aiStore.apiKey && !aiStore.apiBaseUrl) {
      throw new Error('请先在“我的-设置-AI模型设置”中配置AI模型')
    }

    const prompt = `你是一个专业的古诗词辅导老师。请为古诗《${poem.title_cns}》（作者：${poem.author_cns}，朝代：${poem.dynasty_cns}）生成【${type}】。
    
    古诗内容：
    ${poem.paragraphs_cns.join('\n')}
    
    要求：
    1. ${type === '译文' ? '提供准确、优美的现代汉语翻译，逐句对应，通俗易懂。' : ''}
    2. ${type === '赏析' ? '从意象、修辞、情感、写作手法等方面进行深入浅出的赏析，适合学生阅读。' : ''}
    3. ${type === '趣事' ? '讲述与这首诗或作者相关的历史趣事、背景故事或传说，语言生动有趣。' : ''}
    4. 输出格式清晰，直接输出内容，不要包含"好的"、"当然"等客套话。`

    return this.callAI(prompt, aiStore)
  }

  private async callAI(prompt: string, settings: any): Promise<string> {
    const payload = {
        model: settings.model,
        messages: [
            { role: "system", content: "You are a helpful assistant specialized in Chinese poetry." },
            { role: "user", content: prompt }
        ],
        temperature: 0.7
    }

    let url = settings.apiBaseUrl.replace(/\/+$/, '')
    if (!url.endsWith('/v1') && !url.includes('google')) {
       if (!url.includes('/chat/completions')) {
           url = `${url}/chat/completions`
       }
    } else if (url.endsWith('/v1')) {
        url = `${url}/chat/completions`
    }

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${settings.apiKey}`
            },
            body: JSON.stringify(payload)
        })

        if (!response.ok) {
            const errText = await response.text()
            console.error('AI Call Failed:', response.status, errText)
            throw new Error(`AI请求失败 (${response.status}): ${errText}`)
        }
        
        const data = await response.json()
        if (data.choices && data.choices.length > 0) {
            return data.choices[0].message.content
        }
        throw new Error('AI返回数据格式异常')
    } catch (e: any) {
        console.error('AI Service Error:', e)
        throw e
    }
  }
}

export const poetryService = new PoetryService()
