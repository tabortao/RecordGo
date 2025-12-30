import http from './http'
import { presignUpload, putToURL } from './storage'
import { prepareUpload } from '@/utils/image'

// 中文注释：任务服务封装，统一使用 RESTful 接口并返回数据
export interface TaskItem {
  id: number
  user_id: number
  // 中文注释：重复任务系列 ID（同一系列的实例共享该 ID）；一次性任务为空
  series_id?: number
  name: string
  description: string
  category: string
  score: number
  plan_minutes: number
  actual_minutes: number
  status: number
  start_date: string
  end_date?: string
  remark?: string
  image_json?: string
  tomato_count?: number
}

export async function listTasks(params?: { status?: number; page?: number; page_size?: number }): Promise<{ items: TaskItem[]; total: number; page: number; page_size: number }>{
  const raw = await http.get('/tasks', { params }) as any
  const arr: any[] = Array.isArray(raw.items) ? raw.items : []
  const items: TaskItem[] = arr.map((x: any) => ({
    id: Number(x.id ?? x.ID ?? 0),
    user_id: Number(x.user_id ?? x.UserID ?? 0),
    series_id: x.series_id != null ? Number(x.series_id) : (x.SeriesID != null ? Number(x.SeriesID) : undefined),
    name: String(x.name ?? x.Name ?? ''),
    description: String(x.description ?? x.Description ?? ''),
    category: String(x.category ?? x.Category ?? ''),
    score: Number(x.score ?? x.Score ?? 0),
    plan_minutes: Number(x.plan_minutes ?? x.PlanMinutes ?? 0),
    actual_minutes: Number(x.actual_minutes ?? x.ActualMinutes ?? 0),
    status: Number(x.status ?? x.Status ?? 0),
    start_date: String(x.start_date ?? x.StartDate ?? ''),
    end_date: (x.end_date ?? x.EndDate) ? String(x.end_date ?? x.EndDate) : undefined,
    remark: String(x.remark ?? x.Remark ?? ''),
    image_json: String(x.image_json ?? x.ImageJSON ?? ''),
    tomato_count: Number(x.tomato_count ?? x.TomatoCount ?? 0),
    ...(x.repeat != null || x.Repeat != null ? { repeat: String(x.repeat ?? x.Repeat ?? 'none') } : {}),
    ...(Array.isArray(x.weekly_days ?? x.WeeklyDays) ? { weekly_days: (x.weekly_days ?? x.WeeklyDays) } : {})
  }))
  return { items, total: Number(raw.total ?? arr.length ?? 0), page: Number(raw.page ?? 1), page_size: Number(raw.page_size ?? (params?.page_size ?? 20)) }
}

export async function createTask(payload: any): Promise<TaskItem> {
  const norm = { ...payload }
  if (norm.start_date) {
    const d = new Date(norm.start_date)
    const y = d.getFullYear(), m = d.getMonth(), day = d.getDate()
    norm.start_date = new Date(Date.UTC(y, m, day, 0, 0, 0))
  }
  if (norm.end_date) {
    const d = new Date(norm.end_date)
    const y = d.getFullYear(), m = d.getMonth(), day = d.getDate()
    norm.end_date = new Date(Date.UTC(y, m, day, 0, 0, 0))
  }
  return (await http.post('/tasks', norm)) as any
}

export async function createTasksBatch(payload: any): Promise<{ items: TaskItem[]; count: number }> {
  return (await http.post('/tasks/batch', payload)) as any
}

export interface AITaskParseItem {
  name: string
  description: string
  category: string
  score: number
  plan_minutes: number
  start_date: string
  end_date?: string
  repeat_type: 'none'|'daily'|'weekdays'|'weekly'|'monthly'
  weekly_days: number[]
  confidence?: 'High' | 'Medium' | 'Low'
}

export async function parseTaskByAI(text: string, image?: File, aiConfig?: { url: string; key: string; model: string }, categories?: string[]): Promise<{ tasks: AITaskParseItem[] }> {
  const formData = new FormData()
  formData.append('text', text)
  if (image) formData.append('image', image)
  if (aiConfig) {
    formData.append('ai_base_url', aiConfig.url)
    formData.append('ai_key', aiConfig.key)
    formData.append('ai_model', aiConfig.model)
  }
  if (categories && categories.length > 0) {
    formData.append('categories', categories.join(','))
  }
  // Increase timeout for AI Vision tasks (120s)
  const res = (await http.post('/ai/parse-task', formData, { timeout: 120000 })) as { tasks: AITaskParseItem[] }
  if (res.tasks) {
    res.tasks.forEach(t => {
      // Default duration to 20 minutes if not present
      if (!t.plan_minutes) t.plan_minutes = 20
    })
  }
  return res
}

export async function updateTask(id: number, payload: any): Promise<TaskItem> {
  const norm = { ...payload }
  if (norm.start_date) {
    const d = new Date(norm.start_date)
    const y = d.getFullYear(), m = d.getMonth(), day = d.getDate()
    norm.start_date = new Date(Date.UTC(y, m, day, 0, 0, 0))
  }
  if (norm.end_date) {
    const d = new Date(norm.end_date)
    const y = d.getFullYear(), m = d.getMonth(), day = d.getDate()
    norm.end_date = new Date(Date.UTC(y, m, day, 0, 0, 0))
  }
  return (await http.put(`/tasks/${id}`, norm)) as any
}

export async function updateTaskStatus(id: number, status: number, opts?: { allowByTomato?: boolean }): Promise<TaskItem> {
  const body: any = { status }
  if (opts?.allowByTomato) body.allow_by_tomato = true
  return (await http.patch(`/tasks/${id}/status`, body)) as any
}

// 中文注释：获取单个任务详情，独立页面加载使用
export async function getTask(id: number): Promise<TaskItem> {
  return (await http.get(`/tasks/${id}`)) as any
}

export async function deleteTask(id: number): Promise<{ id: number }> {
  return (await http.delete(`/tasks/${id}`)) as any
}

export async function batchDelete(ids: number[]): Promise<{ deleted: number }> {
  return (await http.delete('/tasks', { params: { ids: ids.join(',') } })) as any
}

export async function listRecycleBin(): Promise<TaskItem[]> {
  return (await http.get('/tasks/recycle-bin')) as any
}

export async function restoreTasks(ids: number[]): Promise<{ restored: number }> {
  return (await http.post('/tasks/recycle-bin/restore', null, { params: { ids: ids.join(',') } })) as any
}

export async function completeTomato(id: number, minutes = 25): Promise<TaskItem> {
  return (await http.post(`/tasks/${id}/tomato/complete`, { minutes })) as any
}

export async function listTaskOccurrences(params: { date?: string; start?: string; end?: string }): Promise<{ items: { task_id: number, date: string, status: number, minutes: number }[] }> {
  return (await http.get('/tasks/occurrences', { params })) as any
}
export async function completeOccurrence(id: number, payload: { date: string; minutes?: number }): Promise<{ task_id: number, date: string, status: number, minutes: number }> {
  return (await http.post(`/tasks/${id}/occurrences/complete`, payload)) as any
}
export async function uncompleteOccurrence(id: number, payload: { date: string }): Promise<{ task_id: number, date: string, status: number }> {
  return (await http.post(`/tasks/${id}/occurrences/uncomplete`, payload)) as any
}

export async function deleteOccurrence(id: number, payload: { date: string }): Promise<{ task_id: number, date: string, status: number }> {
  return (await http.post(`/tasks/${id}/occurrences/delete`, payload)) as any
}

// 上传任务图片（前端已转换为 webp）
// 中文注释：后端返回 { path }，相对路径 uploads/images/task_images/{用户id}/xxx.webp
export async function uploadTaskImage(
  userId: number,
  file: File,
  taskId: number,
  onProgress?: (percentage: number) => void
): Promise<{ path: string }> {
  const webp = await prepareUpload(file, 0.75)
  const sign = await presignUpload({ resource_type: 'task_image', user_id: userId, task_id: taskId, content_type: 'image/webp', ext: 'webp' })
  await putToURL(sign.upload_url, webp, sign.headers, onProgress)
  return { path: sign.object_key }
}

// 上传任务备注图片（与任务图片区分，受 VIP 限制）
export async function uploadNoteImage(
  userId: number,
  file: File,
  taskId: number,
  onProgress?: (percentage: number) => void
): Promise<{ path: string }> {
  const webp = await prepareUpload(file, 0.75)
  const sign = await presignUpload({ resource_type: 'task_attachment_img', user_id: userId, task_id: taskId, content_type: 'image/webp', ext: 'webp' })
  await putToURL(sign.upload_url, webp, sign.headers, onProgress)
  return { path: sign.object_key }
}

// 删除单个任务图片（物理文件 + 数据库 image_json 更新）
// 中文注释：后端期望传递 path 相对路径（uploads/images/task_images/...）
export async function deleteTaskImage(taskId: number, path: string): Promise<{ images: string[] }> {
  // Axios 的 delete 支持 data 字段传递 JSON 请求体
  return (await http.delete(`/tasks/${taskId}/image`, { data: { path } })) as any
}

const OFFLINE_KEY = 'recordgo_offline_tasks'
type OfflineCreateEntry = {
  name: string
  description: string
  category: string
  score: number
  plan_minutes: number
  start_date: string
  end_date?: string
  repeat_type: 'none'|'daily'|'weekdays'|'weekly'|'monthly'
  weekly_days: number[]
  images: { name: string; type: string; dataURL: string }[]
}

function readOfflineQueue(): OfflineCreateEntry[] {
  try {
    const raw = localStorage.getItem(OFFLINE_KEY) || '[]'
    const arr = JSON.parse(raw)
    return Array.isArray(arr) ? arr : []
  } catch { return [] }
}
function writeOfflineQueue(arr: OfflineCreateEntry[]) {
  try { localStorage.setItem(OFFLINE_KEY, JSON.stringify(arr)) } catch {}
}
export function enqueueOfflineTask(entry: OfflineCreateEntry) {
  const arr = readOfflineQueue()
  arr.push(entry)
  writeOfflineQueue(arr)
}
function dataURLToFile(dataURL: string, name: string, type: string): File {
  const arr = dataURL.split(',')
  const bstr = atob(arr[1] || '')
  const len = bstr.length
  const u8 = new Uint8Array(len)
  for (let i = 0; i < len; i++) u8[i] = bstr.charCodeAt(i)
  return new File([u8], name, { type })
}
export async function syncOfflineTasks(userId: number): Promise<{ synced: number }> {
  const queue = readOfflineQueue()
  if (!queue.length) return { synced: 0 }
  let synced = 0
  const rest: OfflineCreateEntry[] = []
  for (const q of queue) {
    try {
      const s = new Date(q.start_date)
      const e = q.end_date ? new Date(q.end_date) : undefined
      const t = await createTask({
        user_id: userId,
        name: q.name,
        description: q.description,
        category: q.category,
        score: q.score,
        plan_minutes: q.plan_minutes,
        start_date: s,
        end_date: e,
        repeat_type: q.repeat_type,
        weekly_days: q.weekly_days || []
      })
      if (q.images && q.images.length) {
        const paths: string[] = []
        for (const img of q.images) {
          const file = dataURLToFile(img.dataURL, img.name, img.type)
          const webp = await prepareUpload(file, 0.75)
          const { path } = await uploadTaskImage(userId, webp, t.id)
          paths.push(path)
        }
        if (paths.length) await updateTask(t.id, { image_json: JSON.stringify(paths) })
      }
      synced++
    } catch {
      rest.push(q)
    }
  }
  writeOfflineQueue(rest)
  return { synced }
}

// 上传任务音频（wav/mp3），保存到与图片相同的目录
// 中文注释：后端返回 { path }，相对路径 uploads/images/task_images/{用户id}/{任务id}/xxx.wav
export async function uploadTaskAudio(
  userId: number,
  file: File,
  taskId: number,
  onProgress?: (percentage: number) => void
): Promise<{ path: string }> {
  const name = String((file as any).name || '')
  const ext = name.toLowerCase().endsWith('.wav') ? 'wav' : (name.toLowerCase().endsWith('.mp3') ? 'mp3' : 'mp3')
  const ct = ext === 'wav' ? 'audio/wav' : 'audio/mpeg'
  const sign = await presignUpload({ resource_type: 'task_attachment_audio', user_id: userId, task_id: taskId, content_type: ct, ext })
  await putToURL(sign.upload_url, file, sign.headers, onProgress)
  return { path: sign.object_key }
}
