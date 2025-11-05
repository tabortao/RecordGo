import http from './http'

// 中文注释：任务服务封装，统一使用 RESTful 接口并返回数据
export interface TaskItem {
  id: number
  user_id: number
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

export async function updateTaskStatus(id: number, status: number): Promise<TaskItem> {
  return (await http.patch(`/tasks/${id}/status`, { status })) as any
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
