import http from './http'

export async function getStorageInfo(): Promise<{ driver: string, public_base_url?: string }> {
  return await http.get('/storage/info') as any
}

export async function presignUpload(payload: { resource_type: string, user_id: number, task_id?: number, content_type: string, ext: string }): Promise<{ upload_url: string, object_key: string, headers?: Record<string,string>, expires: number }> {
  return await http.post('/storage/presign', payload) as any
}

export async function presignView(object_key: string): Promise<string> {
  const r = await http.get(`/storage/presign-view?object_key=${encodeURIComponent(object_key)}`) as any
  return r.url as string
}

export async function putToURL(url: string, file: File, headers?: Record<string,string>, onProgress?: (p: number) => void): Promise<void> {
  await http.putExternal(url, file, { headers: headers || {}, timeout: 30000, onUploadProgress: (e: any) => {
    const total = e?.total || 0
    const loaded = e?.loaded || 0
    const p = total > 0 ? Math.round((loaded * 100) / total) : 0
    onProgress && onProgress(p)
  } })
}