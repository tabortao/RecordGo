import { defineStore } from 'pinia'
import { useLocalStorage } from '@vueuse/core'

// 中文注释：备注附件类型定义（图片或音频）
export type AttachmentType = 'image' | 'audio'

export interface NoteAttachment {
  type: AttachmentType
  name: string
  url: string // 预览/播放用 URL（可为 objectURL 或相对路径）
  serverPath?: string // 后端返回的相对路径（图片走现有上传接口）
}

export interface TaskNote {
  id: number
  task_id: number
  text: string
  attachments: NoteAttachment[]
  created_at: string
}

// 中文注释：任务备注持久化存储（本地存储），键为 task_id
export const useNotesStore = defineStore('notes', {
  state: () => ({
    // 中文注释：通过 useLocalStorage 持久化，Pinia 会自动对 Ref 做解包，故在 actions 中直接使用 this.notesMap 即可
    notesMap: useLocalStorage<Record<number, TaskNote[]>>('recordgo_task_notes', {})
  }),
  actions: {
    list(taskId: number): TaskNote[] {
      return (this.notesMap as unknown as Record<number, TaskNote[]>)[taskId] || []
    },
    add(note: TaskNote) {
      const map = this.notesMap as unknown as Record<number, TaskNote[]>
      const arr = map[note.task_id] || []
      arr.push(note)
      map[note.task_id] = arr
    },
    remove(taskId: number, noteId: number) {
      const map = this.notesMap as unknown as Record<number, TaskNote[]>
      const arr = map[taskId] || []
      map[taskId] = arr.filter((n: TaskNote) => n.id !== noteId)
    }
  }
})