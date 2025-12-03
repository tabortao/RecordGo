import imageCompression from 'browser-image-compression'

function isMobile(): boolean {
  if (typeof navigator === 'undefined') return false
  const ua = navigator.userAgent || ''
  return /(Android|iPhone|iPad|iPod)/i.test(ua)
}

function withTimeout<T>(p: Promise<T>, ms: number): Promise<T | undefined> {
  return Promise.race<T | undefined>([
    p,
    new Promise<undefined>((resolve) => setTimeout(() => resolve(undefined), ms))
  ])
}

async function loadImage(file: File): Promise<HTMLImageElement> {
  const url = URL.createObjectURL(file)
  return new Promise((resolve, reject) => {
    const img = new Image()
    img.onload = () => { URL.revokeObjectURL(url); resolve(img) }
    img.onerror = (e) => { URL.revokeObjectURL(url); reject(e) }
    img.src = url
  })
}

async function canvasCompress(file: File, maxDim: number, quality: number): Promise<Blob | undefined> {
  try {
    const img = await loadImage(file)
    const ratio = Math.min(1, Math.min(maxDim / (img.width || 1), maxDim / (img.height || 1)))
    const w = Math.max(1, Math.round((img.width || 1) * ratio))
    const h = Math.max(1, Math.round((img.height || 1) * ratio))
    const canvas = document.createElement('canvas')
    canvas.width = w
    canvas.height = h
    const ctx = canvas.getContext('2d')
    if (!ctx) return undefined
    ctx.drawImage(img, 0, 0, w, h)
    // 快速透明检测：缩放到 64x64 采样 alpha
    let hasAlpha = false
    try {
      const t = document.createElement('canvas')
      t.width = 64; t.height = 64
      const tc = t.getContext('2d')
      if (tc) {
        tc.drawImage(canvas, 0, 0, 64, 64)
        const data = tc.getImageData(0, 0, 64, 64).data
        for (let i = 3; i < data.length; i += 4) { if (data[i] < 255) { hasAlpha = true; break } }
      }
    } catch {}
    const preferType = hasAlpha ? 'image/png' : 'image/jpeg'
    const jpegQ = Math.max(0.5, Math.min(0.9, quality))
    const blob: Blob | null = await new Promise((resolve) => canvas.toBlob(resolve, preferType, jpegQ))
    // 作为备选，尝试 webp（若浏览器支持且更小）
    let webp: Blob | null = null
    try { webp = await new Promise((resolve) => canvas.toBlob(resolve, 'image/webp', jpegQ)) } catch {}
    const candidates = [blob, webp].filter(Boolean) as Blob[]
    if (!candidates.length) return undefined
    candidates.sort((a, b) => a.size - b.size)
    return candidates[0]
  } catch { return undefined }
}

async function compressBest(file: File, quality = 0.8): Promise<File> {
  if (!file.type.startsWith('image/')) return file
  const originalMB = Math.max(0.01, file.size / (1024 * 1024))
  const targetMB = Math.min(0.22, Math.max(0.06, originalMB * 0.25))
  const maxDim = isMobile() ? 1280 : 1920
  const opts: any = {
    maxSizeMB: targetMB,
    maxWidthOrHeight: maxDim,
    useWebWorker: true,
    initialQuality: Math.max(0.5, Math.min(0.85, quality))
  }
  const timeoutMs = isMobile() ? 4500 : 8000
  const tasks = [
    withTimeout(imageCompression(file, opts), timeoutMs),
    withTimeout(canvasCompress(file, maxDim, quality), timeoutMs)
  ]
  const settled = await Promise.all(tasks)
  const blobs = settled.filter(Boolean) as Blob[]
  if (!blobs.length) return file
  blobs.sort((a, b) => a.size - b.size)
  const best = blobs[0]
  const bestType = (best as any).type || file.type
  const ext = (() => {
    if (bestType.includes('jpeg')) return 'jpg'
    if (bestType.includes('png')) return 'png'
    if (bestType.includes('webp')) return 'webp'
    const raw = (file.name.split('.').pop() || '').toLowerCase()
    return raw || 'jpg'
  })()
  const name = replaceExt(file.name, ext)
  return new File([best], name, { type: bestType })
}

export async function toWebp(file: File, quality = 0.8): Promise<File> {
  return compressBest(file, quality)
}

function replaceExt(name: string, ext: string): string {
  const idx = name.lastIndexOf('.')
  return (idx > -1 ? name.substring(0, idx) : name) + (ext ? '.' + ext : '')
}

export async function prepareUpload(file: File, quality = 0.8): Promise<File> {
  try {
    if (!file.type.startsWith('image/')) return file
    const originalSize = file.size
    const converted = await compressBest(file, quality)
    const preferOriginalIfClose = quality >= 0.8
    if (!converted) return file
    if (preferOriginalIfClose ? converted.size >= Math.max(1, Math.floor(originalSize * 0.98)) : converted.size >= originalSize) return file
    return converted
  } catch {
    return file
  }
}
