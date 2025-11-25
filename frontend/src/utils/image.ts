// 中文注释：统一使用 image-conversion 进行高效压缩与 webp 转换
import { compressAccurately } from 'image-conversion'

export async function toWebp(file: File): Promise<File> {
  if (!file.type.startsWith('image/')) return file
  if (file.type === 'image/webp') return file
  const originalKB = Math.max(1, Math.round(file.size / 1024))
  const targetKB = Math.min(220, Math.max(60, Math.round(originalKB * 0.2)))
  const blob = await compressAccurately(file, { size: targetKB, accuracy: 0.92, type: 'image/webp' }).catch(() => null)
  if (!blob) return file
  return new File([blob], replaceExt(file.name, 'webp'), { type: 'image/webp' })
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
    const converted = await toWebp(file)
    const preferOriginalIfClose = quality >= 0.8
    if (!converted) return file
    if (preferOriginalIfClose ? converted.size >= Math.max(1, Math.floor(originalSize * 0.98)) : converted.size >= originalSize) return file
    return converted
  } catch {
    return file
  }
}
