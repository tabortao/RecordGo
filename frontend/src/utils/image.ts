// 中文注释：浏览器端图片压缩与 webp 转换工具，失败时回退到原格式
export async function toWebp(file: File, quality = 0.8): Promise<File> {
  // 非图片直接返回
  if (!file.type.startsWith('image/')) return file

  const bitmap = await createImageBitmap(file).catch(() => null)
  if (!bitmap) return file

  const canvas = document.createElement('canvas')
  canvas.width = bitmap.width
  canvas.height = bitmap.height
  const ctx = canvas.getContext('2d')
  if (!ctx) return file
  ctx.drawImage(bitmap, 0, 0)

  // 尝试导出 webp 数据
  const blob: Blob | null = await new Promise((resolve) =>
    canvas.toBlob((b) => resolve(b), 'image/webp', quality)
  )

  if (!blob) return file

  return new File([blob], replaceExt(file.name, 'webp'), { type: 'image/webp' })
}

// 中文注释：替换文件扩展名为指定后缀
function replaceExt(name: string, ext: string): string {
  const idx = name.lastIndexOf('.')
  return (idx > -1 ? name.substring(0, idx) : name) + '.' + ext
}

// 中文注释：示例上传流程（调用接口前转换 webp）
export async function prepareUpload(file: File): Promise<File> {
  try {
    const webp = await toWebp(file)
    return webp
  } catch {
    // 转换失败回退原文件
    return file
  }
}
