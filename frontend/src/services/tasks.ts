import http from './http'

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
  const form = new FormData()
  form.append('user_id', String(userId))
  form.append('task_id', String(taskId))
  // 中文注释：后端期望字段名为 image（兼容 file）
  // 中文注释：显式附带文件名，提升后端解析稳定性（中文文件名也可正常传递）
  form.append('image', file, (file as any).name || 'image.webp')
  // 中文注释：为最大兼容性，额外附带 file 键（服务端会择一使用）
  form.append('file', file, (file as any).name || 'image.webp')
  // 中文注释：不手动设置 Content-Type，浏览器自动带 boundary，避免连接被重置；支持上传进度回调
  const resp = await http.post('/upload/task-image', form, {
    timeout: 30000,
    onUploadProgress: (e: any) => {
      try {
        const total = e?.total || 0
        const loaded = e?.loaded || 0
        const p = total > 0 ? Math.round((loaded * 100) / total) : 0
        onProgress && onProgress(p)
      } catch {}
    }
  } as any)
  return resp as { path: string }
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
  const form = new FormData()
  form.append('user_id', String(userId))
  form.append('task_id', String(taskId))
  // 中文注释：优先使用 audio 字段，兼容 file 字段
  form.append('audio', file, (file as any).name || 'audio.wav')
  form.append('file', file, (file as any).name || 'audio.wav')
  const resp = await http.post('/upload/task-audio', form, {
    timeout: 30000,
    onUploadProgress: (e: any) => {
      try {
        const total = e?.total || 0
        const loaded = e?.loaded || 0
        const p = total > 0 ? Math.round((loaded * 100) / total) : 0
        onProgress && onProgress(p)
      } catch {}
    }
  } as any)
  return resp as { path: string }
}
