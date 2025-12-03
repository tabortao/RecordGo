import imageCompression from 'browser-image-compression'

export async function toWebp(file: File, quality = 0.8): Promise<File> {
  if (!file.type.startsWith('image/')) return file
  if (file.type === 'image/webp') return file
  const originalMB = Math.max(0.01, file.size / (1024 * 1024))
  const targetMB = Math.min(0.22, Math.max(0.06, originalMB * 0.2))
  const opts: any = {
    maxSizeMB: targetMB,
    maxWidthOrHeight: 1920,
    useWebWorker: true,
    fileType: 'image/webp',
    initialQuality: Math.max(0.6, Math.min(0.9, quality))
  }
  try {
    const out = await imageCompression(file, opts)
    if (!out) return file
    const name = replaceExt(file.name, 'webp')
    return new File([out], name, { type: 'image/webp' })
  } catch {
    return file
  }
}

// 中文注释：替换文件扩展名为指定后缀
function replaceExt(name: string, ext: string): string {
  const idx = name.lastIndexOf('.')
  return (idx > -1 ? name.substring(0, idx) : name) + '.' + ext
}

// 中文注释：示例上传流程（调用接口前转换 webp）
export async function prepareUpload(file: File, quality = 0.8): Promise<File> {
  try {
    if (!file.type.startsWith('image/')) return file
    if (file.type === 'image/webp') return file
    const originalSize = file.size
    const converted = await toWebp(file, quality)
    const preferOriginalIfClose = quality >= 0.8
    if (!converted) return file
    if (preferOriginalIfClose ? converted.size >= Math.max(1, Math.floor(originalSize * 0.98)) : converted.size >= originalSize) return file
    return converted
  } catch {
    return file
  }
}
