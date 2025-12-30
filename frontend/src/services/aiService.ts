interface AISettings {
    apiBaseUrl: string
    apiKey: string
    model: string
    visionApiBaseUrl: string
    visionApiKey: string
    visionModel: string
}

export async function testConnection(config: AISettings): Promise<boolean> {
    const settings = config
    const systemPrompt = "You are a helpful assistant."
    const userPrompt = "Hello, reply with 'OK'."
    
    // Minimal OpenAI Compatible Payload
    const payload = {
        model: settings.model,
        messages: [
            { role: "system", content: systemPrompt },
            { role: "user", content: userPrompt }
        ],
        max_tokens: 10
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
            console.error('Test Connection Failed:', response.status, await response.text())
            return false
        }
        
        const data = await response.json()
        if (data.choices && data.choices.length > 0) {
            return true
        }
        return false
    } catch (e) {
        console.error('Test Connection Error:', e)
        return false
    }
}

export async function testVisionConnection(config: AISettings): Promise<boolean> {
    const baseUrl = config.visionApiBaseUrl || config.apiBaseUrl
    const apiKey = config.visionApiKey || config.apiKey
    const model = config.visionModel || config.model
    
    // Try text-only handshake first as it is lighter and supported by most multimodal models
    const payload = {
        model: model,
        messages: [
            {
                role: "user",
                content: "Hello, just checking connection. Reply 'OK'."
            }
        ],
        max_tokens: 10
    }

    let url = baseUrl.replace(/\/+$/, '')
    if (!url.endsWith('/v1') && !url.includes('google')) {
         if (!url.includes('/chat/completions')) url = `${url}/chat/completions`
    } else if (url.endsWith('/v1')) {
         url = `${url}/chat/completions`
    }

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${apiKey}`
            },
            body: JSON.stringify(payload)
        })

        if (!response.ok) {
            console.error('Test Vision Connection Failed:', response.status, await response.text())
            return false
        }
        
        const data = await response.json()
        if (data.choices && data.choices.length > 0) {
            return true
        }
        return false
    } catch (e) {
        console.error('Test Vision Connection Error:', e)
        return false
    }
}
