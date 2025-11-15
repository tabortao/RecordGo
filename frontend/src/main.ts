import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import './styles/index.css'
import router from './router'
import App from './App.vue'

// 中文注释：应用入口，注册 Pinia、路由与 Element Plus
const app = createApp(App)
app.use(createPinia())
app.use(router)
app.use(ElementPlus)
app.mount('#app')

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
  const path = router.currentRoute.value.path || ''
  if (path.startsWith('/tomato')) return 'rgb(48, 48, 46)'
  const isDark = document.documentElement.classList.contains('dark')
  return isDark ? 'rgb(17, 24, 39)' : 'rgb(254, 254, 254)'
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

router.afterEach(() => updateThemeBar())

updateThemeBar()