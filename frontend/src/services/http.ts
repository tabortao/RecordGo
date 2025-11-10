import axios, { AxiosRequestConfig } from 'axios'

// 中文注释：Axios 实例，统一处理 RESTful 响应结构 {code,message,data}，并实现指数退避重试
// 开发环境优先使用 Vite 代理（baseURL='/api'，避免浏览器 CORS）；生产或独立部署使用 VITE_API_BASE
const base = (import.meta as any).env?.DEV ? '/api' : ((import.meta as any).env?.VITE_API_BASE || '/api')
const instance = axios.create({
  baseURL: base,
  timeout: 10000
})

// 请求拦截器：可附加 JWT 等
instance.interceptors.request.use((config) => {
  // 中文注释：此处可注入鉴权头，如从 localStorage 读取 token
  const token = localStorage.getItem('auth_token')
  if (token) {
    config.headers = config.headers || {}
    ;(config.headers as any).Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器：统一业务错误处理
instance.interceptors.response.use(
  (resp) => {
    const payload = resp.data
    if (payload && typeof payload.code !== 'undefined') {
      if (payload.code === 0) return payload.data
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
// 中文注释：补充 PATCH 方法（用于部分更新，如任务状态切换）
export async function patch<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
  return (await instance.patch(url, data, config)) as any
}
export async function del<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
  return (await instance.delete(url, config)) as any
}

// 中文注释：统一导出，包含 patch 与 delete（命名为 delete 以便直觉调用）
const http = { get, post, put, patch, delete: del }

export default http
