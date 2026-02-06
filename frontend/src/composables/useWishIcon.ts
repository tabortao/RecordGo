import { ref } from 'vue'
import { getStaticBase } from '@/services/http'
import { presignView } from '@/services/storage'

export function isEmojiIcon(icon: unknown): icon is string {
  return typeof icon === 'string' && icon.startsWith('emoji:')
}

export function emojiChar(icon: unknown): string {
  if (!isEmojiIcon(icon)) return ''
  return String(icon).slice('emoji:'.length)
}

export function useWishIconResolver(options?: { fallbackBuiltin?: string }) {
  const fallbackBuiltin = options?.fallbackBuiltin || '零花钱.png'
  const iconResolvedMap = ref<Record<string, string>>({})

  function fallbackSrc() {
    return new URL(`../assets/wishs/${fallbackBuiltin}`, import.meta.url).href
  }

  function resolveIconSrc(icon: string | undefined) {
    if (!icon) return fallbackSrc()
    if (isEmojiIcon(icon)) return fallbackSrc()
    if (/\.(png|jpg|jpeg|webp)$/i.test(icon) && !icon.includes('/')) {
      return new URL(`../assets/wishs/${icon}`, import.meta.url).href
    }
    const base = getStaticBase()
    const path = String(icon).replace(/^\/+/, '')
    if (path.startsWith('uploads/')) return `${base}/api/${path}`
    return iconResolvedMap.value[path] || fallbackSrc()
  }

  function resolveBuiltinByName(name: string, ext = 'png') {
    const file = `${name}.${ext}`
    try {
      return new URL(`../assets/wishs/${file}`, import.meta.url).href
    } catch {
      return fallbackSrc()
    }
  }

  async function prefetchPresigned(icons: Array<string | undefined>) {
    const targets = icons
      .map(i => String(i || ''))
      .filter(p => !!p && !p.startsWith('emoji:') && !p.startsWith('uploads/') && p.includes('/'))

    const uniq = Array.from(new Set(targets))
    await Promise.all(uniq.map(async (k) => {
      if (iconResolvedMap.value[k]) return
      try {
        iconResolvedMap.value[k] = await presignView(k)
      } catch {
        return
      }
    }))
  }

  return { iconResolvedMap, resolveIconSrc, resolveBuiltinByName, prefetchPresigned }
}

