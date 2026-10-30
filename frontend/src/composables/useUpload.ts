import { ref } from 'vue'
import { uploadImage } from '@/api/upload'
import { toErrorMessage } from '@/api/client'
import type { UploadResult } from '@/types'

export interface UploadTask {
  id: string
  file: File
  preview: string
  percent: number
  status: 'uploading' | 'done' | 'error'
  error?: string
  result?: UploadResult
  controller?: AbortController
}

const MAX_SIZE = 10 * 1024 * 1024
const ALLOWED = ['image/jpeg', 'image/png', 'image/webp']

let seq = 0

/** 图片上传：校验类型/大小、显示进度、支持取消与重试。 */
export function useUpload(maxCount = 5) {
  const tasks = ref<UploadTask[]>([])
  const error = ref('')

  function validate(file: File, current: number): string {
    if (!ALLOWED.includes(file.type)) return '仅支持 jpg / png / webp 图片'
    if (file.size > MAX_SIZE) return '单张图片不能超过 10MB'
    if (current >= maxCount) return `最多上传 ${maxCount} 张图片`
    return ''
  }

  async function start(task: UploadTask) {
    task.status = 'uploading'
    task.percent = 0
    task.error = ''
    task.controller = new AbortController()
    try {
      const result = await uploadImage(
        task.file,
        (percent) => {
          task.percent = percent
        },
        task.controller.signal,
      )
      task.result = result
      task.status = 'done'
      task.percent = 100
    } catch (err) {
      task.status = 'error'
      task.error = toErrorMessage(err)
    }
  }

  async function add(files: FileList | File[]): Promise<UploadTask[]> {
    error.value = ''
    const list = Array.from(files)
    const created: UploadTask[] = []
    for (const file of list) {
      const message = validate(file, tasks.value.length)
      if (message) {
        error.value = message
        continue
      }
      const task: UploadTask = {
        id: `upload-${++seq}`,
        file,
        preview: URL.createObjectURL(file),
        percent: 0,
        status: 'uploading',
      }
      tasks.value.push(task)
      created.push(task)
    }
    await Promise.all(created.map((task) => start(task)))
    return created
  }

  function cancel(id: string) {
    const task = tasks.value.find((item) => item.id === id)
    if (!task) return
    task.controller?.abort()
    remove(id)
  }

  function retry(id: string) {
    const task = tasks.value.find((item) => item.id === id)
    if (!task) return
    void start(task)
  }

  function remove(id: string) {
    const index = tasks.value.findIndex((item) => item.id === id)
    if (index < 0) return
    URL.revokeObjectURL(tasks.value[index].preview)
    tasks.value.splice(index, 1)
  }

  function reset() {
    tasks.value.forEach((task) => URL.revokeObjectURL(task.preview))
    tasks.value = []
    error.value = ''
  }

  return { tasks, error, add, cancel, retry, remove, reset, maxCount }
}
