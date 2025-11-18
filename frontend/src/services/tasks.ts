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

export async function listTasks(params?: { status?: number; page?: number; page_size?: number }): Promise<{ items: TaskItem[]; total: number; page: number; page_size: number }>
{
  // 中文注释：返回值显式类型，避免 AxiosResponse 推断造成的 TS 报错
  return (await http.get('/tasks', { params })) as any
}

export async function createTask(payload: any): Promise<TaskItem> {
  return (await http.post('/tasks', payload)) as any
}

export async function updateTask(id: number, payload: any): Promise<TaskItem> {
  return (await http.put(`/tasks/${id}`, payload)) as any
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

// 上传任务图片（前端已转换为 webp）
// 中文注释：后端返回 { path }，相对路径 uploads/images/task_images/{用户id}/xxx.webp
export async function uploadTaskImage(
  userId: number,
  file: File,
  taskId: number,
  onProgress?: (percentage: number) => void
): Promise<{ path: string }> {
  const webp = await prepareUpload(file)
  const sign = await presignUpload({ resource_type: 'task_image', user_id: userId, task_id: taskId, content_type: 'image/webp', ext: 'webp' })
  await putToURL(sign.upload_url, webp, sign.headers, onProgress)
  return { path: sign.object_key }
}

// 删除单个任务图片（物理文件 + 数据库 image_json 更新）
// 中文注释：后端期望传递 path 相对路径（uploads/images/task_images/...）
export async function deleteTaskImage(taskId: number, path: string): Promise<{ images: string[] }> {
  // Axios 的 delete 支持 data 字段传递 JSON 请求体
  return (await http.delete(`/tasks/${taskId}/image`, { data: { path } })) as any
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
