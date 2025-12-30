import http from '@/services/http'

export type OCRModelType = 'PP-OCRv5' | 'PaddleOCR-VL' | 'PP-StructureV3'

export interface OCRConfig {
  apiUrl: string
  token: string
}

async function sendProxyRequest(targetUrl: string, token: string, payload: any): Promise<any> {
  // Use backend proxy to avoid CORS issues with third-party OCR APIs
  return http.post('/proxy/ocr', {
    target_url: targetUrl,
    token: token,
    payload: payload
  })
}

export async function recognizeImage(
  file: File,
  model: OCRModelType,
  config: OCRConfig
): Promise<string> {
  // Convert file to base64
  const base64 = await fileToBase64(file)
  // Remove data URL prefix if present (e.g., "data:image/png;base64,")
  const fileData = base64.includes(',') ? base64.split(',')[1] : base64

  // Build payload based on PRD requirements
  const requiredPayload = {
    file: fileData,
    fileType: 1, // 1 for images
  }

  let optionalPayload = {}
  
  if (model === 'PP-OCRv5') {
    optionalPayload = {
        useDocOrientationClassify: false,
        useDocUnwarping: false,
        useTextlineOrientation: false,
    }
  } else if (model === 'PaddleOCR-VL') {
    optionalPayload = {
        useDocOrientationClassify: false,
        useDocUnwarping: false,
        useChartRecognition: false,
    }
  } else {
    // PP-StructureV3
    optionalPayload = {
        useDocOrientationClassify: false,
        useDocUnwarping: false,
        useTextlineOrientation: false,
        useChartRecognition: false,
    }
  }

  const payload = { ...requiredPayload, ...optionalPayload }

  try {
    // Use proxy request instead of direct axios call
    const data: any = await sendProxyRequest(config.apiUrl, config.token, payload)
    
    // The proxy returns the raw response from OCR service
    console.log('[OCR] Raw Response Data:', JSON.stringify(data, null, 2))
    const result = data.result
    
    if (!result) {
        console.error('[OCR] Invalid response structure: missing "result" field', data)
        throw new Error('Invalid response from OCR service')
    }

    if (model === 'PP-OCRv5') {
       // result["ocrResults"] is array
       console.log('[OCR] PP-OCRv5 Results:', result.ocrResults)
       if (result.ocrResults && Array.isArray(result.ocrResults)) {
           const text = result.ocrResults.map((res: any) => {
               // Handle structured prunedResult (contains rec_texts)
               if (res.prunedResult && typeof res.prunedResult === 'object' && Array.isArray(res.prunedResult.rec_texts)) {
                   return res.prunedResult.rec_texts.join('\n')
               }
               // Fallback: direct string or stringify
               return typeof res.prunedResult === 'string' ? res.prunedResult : ''
           }).filter((t: string) => t).join('\n')
           
           if (!text) console.warn('[OCR] PP-OCRv5 extracted text is empty')
           return text
       }
       console.warn('[OCR] PP-OCRv5 result.ocrResults is missing or not an array', result)
       return ''
    } else {
       // PaddleOCR-VL and PP-StructureV3 use layoutParsingResults -> markdown -> text
       if (result.layoutParsingResults && Array.isArray(result.layoutParsingResults)) {
         return result.layoutParsingResults.map((res: any) => res.markdown?.text || '').join('\n\n')
       }
       return ''
    }
  } catch (error) {
    console.error('OCR Failed:', error)
    throw error
  }
}

export async function testConnection(
  model: OCRModelType,
  config: OCRConfig
): Promise<boolean> {
  // 1x1 pixel transparent PNG
  const fileData = 'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFhAJ/wlseKgAAAABJRU5ErkJggg=='
  
  const requiredPayload = {
    file: fileData,
    fileType: 1,
  }

  let optionalPayload = {}
  
  if (model === 'PP-OCRv5') {
    optionalPayload = {
        useDocOrientationClassify: false,
        useDocUnwarping: false,
        useTextlineOrientation: false,
    }
  } else if (model === 'PaddleOCR-VL') {
    optionalPayload = {
        useDocOrientationClassify: false,
        useDocUnwarping: false,
        useChartRecognition: false,
    }
  } else {
    optionalPayload = {
        useDocOrientationClassify: false,
        useDocUnwarping: false,
        useTextlineOrientation: false,
        useChartRecognition: false,
    }
  }

  const payload = { ...requiredPayload, ...optionalPayload }

  try {
    // Use proxy request
    await sendProxyRequest(config.apiUrl, config.token, payload)
    return true
  } catch (error) {
    console.error('OCR Connection Test Failed:', error)
    throw error
  }
}

function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.readAsDataURL(file)
    reader.onload = () => resolve(reader.result as string)
    reader.onerror = error => reject(error)
  })
}
