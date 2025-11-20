import axios, { AxiosRequestConfig } from 'axios'

// 中文注释：Axios 实例，统一处理 RESTful 响应结构 {code,message,data}，并实现指数退避重试
// 开发环境优先使用 Vite 代理（baseURL='/api'，避免浏览器 CORS）；生产或独立部署使用 VITE_API_BASE
let base = (import.meta as any).env?.DEV ? '/api' : ((import.meta as any).env?.VITE_API_BASE || '/api')
if (base && base !== '/api' && /^https?:\/\//i.test(String(base))) {
  base = String(base).replace(/\/+$/, '')
  if (!/\/api$/i.test(base)) base = base + '/api'
}
const instance = axios.create({
  baseURL: base,
  timeout: 10000
})

// 请求拦截器：可附加 JWT 等
let __envLogged = false
instance.interceptors.request.use((config) => {
  // 中文注释：为兼容“未勾选记住我”的会话，优先读取 sessionStorage，其次读取 localStorage
  const token = (sessionStorage.getItem('auth_token') || localStorage.getItem('auth_token'))
  if (token) {
    config.headers = config.headers || {}
    ;(config.headers as any).Authorization = `Bearer ${token}`
  }
  if (!__envLogged) {
    __envLogged = true
    try { console.info('[API] baseURL=', instance.defaults.baseURL, 'host=', location.host) } catch {}
  }
  return config
})

// 响应拦截器：统一业务错误处理
instance.interceptors.response.use(
  async (resp) => {
    const payload = resp.data
    if (typeof payload === 'string' && /<html[\s\S]*<\/html>/i.test(payload)) {
      return Promise.reject(new Error(`API响应异常：收到HTML内容。请在生产环境设置VITE_API_BASE指向后端，如 https://your-api-host`))
    }
    if (payload && typeof payload.code !== 'undefined') {
      if (payload.code === 0) {
        // 中文注释：统一在成功响应中同步最新金币，避免页面各处手动维护导致不一致
        try {
          const data = payload.data
          let latestCoins: any = undefined
          if (data && typeof data === 'object') {
            if (Object.prototype.hasOwnProperty.call(data, 'user_coins')) {
              latestCoins = (data as any).user_coins
            } else if (Object.prototype.hasOwnProperty.call(data, 'coins')) {
              latestCoins = (data as any).coins
            } else if ((data as any).user && typeof (data as any).user === 'object') {
              const u = (data as any).user
              if (Object.prototype.hasOwnProperty.call(u, 'coins')) latestCoins = (u as any).coins
            }
          }
          const n = Number(latestCoins)
          if (Number.isFinite(n)) {
            // 中文注释：为避免循环依赖，采用按需动态导入 Store 并在此处写入
            const { useAppState } = await import('@/stores/appState')
            const { useAuth } = await import('@/stores/auth')
            const appState = useAppState()
            const auth = useAuth()
            appState.setCoins(n)
            auth.updateUser({ coins: n })
          }
        } catch (_) {}
        return payload.data
      }
      // 中文注释：业务错误统一抛出，交由上层捕获
      return Promise.reject(new Error(payload.message || '业务错误'))
    }
    return resp.data
  },
  async (error) => {
    // 中文注释：网络错误重试（最多 3 次，指数退避）
    const config: any = error.config
    config.__retryCount = config.__retryCount || 0
    if (config.__retryCount < 3) {
      const delay = Math.pow(2, config.__retryCount) * 500
      config.__retryCount++
      await new Promise((r) => setTimeout(r, delay))
      return instance(config)
    }
    // 中文注释：如果后端返回未登录（如 401 或业务码），可跳转登录
    try {
      const status = error?.response?.status
      const payload = error?.response?.data
      const code = typeof payload?.code === 'number' ? payload.code : -1
      if (status === 401 || code == 40100) {
        // 清理本地 token 并跳转登录
        localStorage.removeItem('auth_token')
        localStorage.removeItem('auth_user')
        // 为避免循环依赖，这里不直接导入路由；页面调用处可根据返回错误重定向
      }
    } catch (_) {}
    try {
      const data = error?.response?.data
      if (typeof data === 'string' && /<html[\s\S]*<\/html>/i.test(data)) {
        return Promise.reject(new Error(`API响应异常：收到HTML内容。可能是静态托管未配置后端代理，请设置VITE_API_BASE`))
      }
    } catch {}
    return Promise.reject(error)
  }
)

// 中文注释：对外暴露轻量封装，显式返回 Promise<T>
// 说明：这样可以避免在调用处出现 AxiosResponse 的类型与业务数据类型不匹配的报错
export async function get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
  return (await instance.get(url, config)) as any
}
export async function post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
  return (await instance.post(url, data, config)) as any
}
export async function put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
  return (await instance.put(url, data, config)) as any
}
export async function putExternal<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
  const c = { ...(config || {}) }
  const client = axios.create({ timeout: c.timeout || 10000 })
  return (await client.put(url, data, c)) as any
}

// 中文注释：静态资源基址（用于拼接 /api/uploads/... 的完整地址）
// 说明：
// - 开发环境使用 Vite 代理，Axios baseURL 为 '/api'，此时静态基址应为空字符串，最终形成 '/api/uploads/...'
// - 生产环境若配置了 VITE_API_BASE（可为 'https://host' 或 'https://host/api'），需要安全地去掉末尾的 '/api' 前缀，避免出现重复 '/api/api/...'
// - 若未配置 VITE_API_BASE，则回退为相对路径（空字符串），依赖前端托管平台的反向代理规则
export function getStaticBase(): string {
  let b = instance.defaults.baseURL || '/api'
  if (!b) return ''
  b = String(b).trim()
  // 相对路径 '/api'：静态基址置空
  if (b === '/api') return ''
  // 去除尾部斜杠与可能存在的 '/api' 后缀
  b = b.replace(/\/+$/, '')
  b = b.replace(/\/api$/i, '')
  return b
}
// 中文注释：补充 PATCH 方法（用于部分更新，如任务状态切换）
export async function patch<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
  return (await instance.patch(url, data, config)) as any
}
export async function del<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
  return (await instance.delete(url, config)) as any
}

// 中文注释：统一导出，包含 patch 与 delete（命名为 delete 以便直觉调用）
const http = { get, post, put, patch, delete: del, putExternal }

export default http
