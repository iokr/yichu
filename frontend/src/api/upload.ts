import { ApiError, http } from './client'
import type { ApiResponse, UploadResult } from '@/types'

/** 上传单张图片，返回压缩后的原图与缩略图地址。 */
export async function uploadImage(
  file: File,
  onProgress?: (percent: number) => void,
  signal?: AbortSignal,
): Promise<UploadResult> {
  const form = new FormData()
  form.append('file', file)
  const res = await http.post<ApiResponse<UploadResult>>('/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
    signal,
    onUploadProgress: (event) => {
      if (event.total) onProgress?.(Math.round((event.loaded / event.total) * 100))
    },
  })
  const body = res.data
  if (!body || body.code !== 0 || !body.data) {
    throw new ApiError(body?.message || '上传失败', body?.code ?? 5000, res.status)
  }
  return body.data
}
