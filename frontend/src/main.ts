import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { ElMessage } from 'element-plus'
import 'element-plus/theme-chalk/el-message.css'
// English comment: Ensure MessageBox styles are loaded for first-time usage
import 'element-plus/theme-chalk/el-message-box.css'
import dayjs from 'dayjs'
import 'dayjs/locale/zh-cn'
import 'element-plus/theme-chalk/dark/css-vars.css'
import './styles/index.css'
import './assets/main.css'
import './assets/delete-dialog.css'
import router from './router'
import App from './App.vue'
import http from './services/http'

// 中文注释：应用入口，注册 Pinia、路由
const app = createApp(App)
app.use(createPinia())
app.use(router)
dayjs.locale('zh-cn')
app.mount('#app')

// English comment: Prefetch core route components to avoid first-click delay
try {
  setTimeout(() => {
    Promise.all([
      import('./pages/TasksPage.vue'),
      import('./pages/WishesPage.vue'),
      import('./pages/HomeworkPage.vue'),
      import('./pages/MinePage.vue')
    ]).catch(() => {})
  }, 0)
} catch {}

function setThemeColor(color: string) {
  let meta = document.querySelector('meta[name="theme-color"]') as HTMLMetaElement | null
  if (!meta) {
    meta = document.createElement('meta')
    meta.name = 'theme-color'
    document.head.appendChild(meta)
  }
  meta.content = color
}

function currentPageColor(): string {
  const x = Math.max(0, Math.floor(window.innerWidth / 2))
  const y = 2
  let el = document.elementFromPoint(x, y) as HTMLElement | null
  const isTransparent = (c: string) => !c || c === 'transparent' || c === 'rgba(0, 0, 0, 0)'
  let color = ''
  while (el) {
    const c = getComputedStyle(el).backgroundColor
    if (!isTransparent(c)) { color = c; break }
    el = el.parentElement
  }
  if (!color) {
    const bodyC = getComputedStyle(document.body).backgroundColor
    color = isTransparent(bodyC) ? getComputedStyle(document.documentElement).backgroundColor : bodyC
  }
  return color || 'rgb(254, 254, 254)'
}

function updateThemeBar() { setThemeColor(currentPageColor()) }

const media = window.matchMedia('(prefers-color-scheme: dark)')
const apply = (v: boolean) => { document.documentElement.classList.toggle('dark', v); updateThemeBar() }
try {
  const raw = localStorage.getItem('appState')
  const mode = raw ? (JSON.parse(raw).themeMode || 'system') : 'system'
  if (mode === 'dark') {
    apply(true)
  } else if (mode === 'light') {
    apply(false)
  } else {
    apply(media.matches)
    media.addEventListener('change', (e) => apply(e.matches))
  }
} catch {
  apply(media.matches)
  media.addEventListener('change', (e) => apply(e.matches))
}

router.afterEach(() => { setTimeout(updateThemeBar, 0) })
// 路由加载进度：进入前开始，完成后结束
import { useAppState } from './stores/appState'
const store = useAppState()
router.beforeEach(() => { try { if (!store.hasLoadedOnce) { store.startPageLoading(); store.setPageProgress(15) } } catch {} return true })
router.afterEach(() => { try { if (!store.hasLoadedOnce) { store.setPageProgress(100); store.markLoadedOnce(); setTimeout(() => store.stopPageLoading(), 200) } } catch {} })

let ticking = false
function scheduleUpdate() {
  if (ticking) return
  ticking = true
  requestAnimationFrame(() => { ticking = false; updateThemeBar() })
}
window.addEventListener('scroll', scheduleUpdate, { passive: true })
window.addEventListener('resize', scheduleUpdate)

updateThemeBar()

try { (async () => { try { await http.get('/health') } catch { ElMessage.warning('后端API不可访问，请检查 VITE_API_BASE 或反向代理配置') } })() } catch {}

window.addEventListener('error', (e) => {
  const msg = String(e.message || '')
  if (/ResizeObserver/i.test(msg)) { console.warn('Ignored ResizeObserver warning:', msg); return }
  if (msg === 'Script error.') { console.warn('Ignored cross-origin script error'); return }
  try { ElMessage.error(`前端错误：${msg || '未知错误'}`) } catch {}
  console.error('GlobalError', e)
})
window.addEventListener('unhandledrejection', (e: PromiseRejectionEvent) => {
  try { ElMessage.error(`请求失败：${(e.reason && e.reason.message) || '未知错误'}`) } catch {}
  console.error('UnhandledRejection', e.reason)
})

try {
  const isProd = !(import.meta as any).env?.DEV
  const base = (import.meta as any).env?.VITE_API_BASE
  const host = window.location.hostname
  if (isProd && !base && !/^(localhost|127\.0\.0\.1)$/i.test(host)) {
    console.warn('生产环境未配置 VITE_API_BASE，API 将走相对路径 /api，需在托管平台配置反向代理或设置 VITE_API_BASE 指向后端')
  }
} catch {}
