import axios, { type AxiosRequestConfig, type InternalAxiosRequestConfig } from 'axios'
import type { ApiResponse, TokenResponse } from '@/types'

const TOKEN_KEY = 'yichu.token'
const REFRESH_KEY = 'yichu.refresh'
const USER_KEY = 'yichu.user'

export class ApiError extends Error {
  code: number
  status: number

  constructor(message: string, code = 5000, status = 500) {
    super(message)
    this.name = 'ApiError'
    this.code = code
    this.status = status
  }
}

export const session = {
  token: () => localStorage.getItem(TOKEN_KEY) || '',
  refreshToken: () => localStorage.getItem(REFRESH_KEY) || '',
  user: <T>(): T | null => {
    const raw = localStorage.getItem(USER_KEY)
    if (!raw) return null
    try {
      return JSON.parse(raw) as T
    } catch {
      return null
    }
  },
  save(data: TokenResponse) {
    localStorage.setItem(TOKEN_KEY, data.token)
    localStorage.setItem(REFRESH_KEY, data.refresh_token)
    localStorage.setItem(USER_KEY, JSON.stringify(data.user))
    window.dispatchEvent(new CustomEvent('yichu:session'))
  },
  clear() {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_KEY)
    localStorage.removeItem(USER_KEY)
    window.dispatchEvent(new CustomEvent('yichu:session'))
  },
}

export function notifyUnauthorized() {
  window.dispatchEvent(new CustomEvent('yichu:unauthorized'))
}

const http = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' },
})

http.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = session.token()
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

let refreshing: Promise<string> | null = null

function refreshAccessToken(): Promise<string> {
  if (!refreshing) {
    const refreshToken = session.refreshToken()
    if (!refreshToken) {
      session.clear()
      notifyUnauthorized()
      return Promise.reject(new ApiError('登录已过期，请重新登录', 1002, 401))
    }
    refreshing = axios
      .post<ApiResponse<TokenResponse>>('/api/v1/auth/refresh', { refresh_token: refreshToken })
      .then((res) => {
        if (res.data?.code !== 0 || !res.data.data?.token) {
          throw new ApiError(res.data?.message || '登录已过期', 1002, 401)
        }
        session.save(res.data.data)
        return res.data.data.token
      })
      .catch((err) => {
        session.clear()
        notifyUnauthorized()
        throw err
      })
      .finally(() => {
        refreshing = null
      })
  }
  return refreshing
}

http.interceptors.response.use(
  (res) => res,
  async (error) => {
    const config = error.config as (InternalAxiosRequestConfig & { _retry?: boolean }) | undefined
    const status: number = error.response?.status ?? 0
    const body = error.response?.data as ApiResponse<unknown> | undefined
    const url: string = config?.url ?? ''
    const isAuthCall = url.includes('/auth/login') || url.includes('/auth/refresh')
    const unauthorized = status === 401 || body?.code === 1002

    if (unauthorized && !isAuthCall && config && !config._retry) {
      config._retry = true
      try {
        const token = await refreshAccessToken()
        config.headers.Authorization = `Bearer ${token}`
        return http(config)
      } catch {
        return Promise.reject(new ApiError('登录已过期，请重新登录', 1002, 401))
      }
    }
    return Promise.reject(error)
  },
)

function normalizeError(error: unknown): ApiError {
  if (error instanceof ApiError) return error
  if (axios.isAxiosError(error)) {
    const body = error.response?.data as ApiResponse<unknown> | undefined
    if (body?.message) {
      return new ApiError(body.message, body.code ?? 5000, error.response?.status ?? 500)
    }
    if (error.code === 'ECONNABORTED') return new ApiError('请求超时，请稍后重试', 5000, 408)
    if (!error.response) return new ApiError('网络不可用，请联网后重试', 5000, 0)
    return new ApiError(`请求失败（${error.response.status}）`, 5000, error.response.status)
  }
  if (error instanceof Error) return new ApiError(error.message)
  return new ApiError('请求失败，请稍后重试')
}

export async function request<T>(config: AxiosRequestConfig): Promise<T> {
  try {
    const res = await http.request<ApiResponse<T>>(config)
    const body = res.data
    if (body && typeof body.code === 'number' && body.code !== 0) {
      throw new ApiError(body.message || '请求失败', body.code, res.status)
    }
    return (body?.data ?? (null as unknown as T)) as T
  } catch (error) {
    throw normalizeError(error)
  }
}

export function get<T>(url: string, params?: Record<string, unknown>) {
  return request<T>({ method: 'GET', url, params })
}

export function post<T>(url: string, data?: unknown) {
  return request<T>({ method: 'POST', url, data })
}

export function put<T>(url: string, data?: unknown) {
  return request<T>({ method: 'PUT', url, data })
}

export function del<T>(url: string, data?: unknown) {
  return request<T>({ method: 'DELETE', url, data })
}

/** 下载文件（导出 / 备份），返回响应头供调用方判断脱敏状态。 */
export async function download(
  url: string,
  params: Record<string, unknown> | undefined,
  fallbackName: string,
): Promise<{ filename: string; masked: boolean }> {
  try {
    const res = await http.request<Blob>({ method: 'GET', url, params, responseType: 'blob' })
    const disposition = (res.headers['content-disposition'] as string | undefined) ?? ''
    const matched = /filename="?([^";]+)"?/.exec(disposition)
    const filename = matched?.[1] ?? fallbackName
    const masked = String(res.headers['x-price-masked'] ?? 'false') === 'true'
    const blobUrl = URL.createObjectURL(res.data)
    const link = document.createElement('a')
    link.href = blobUrl
    link.download = filename
    document.body.appendChild(link)
    link.click()
    link.remove()
    setTimeout(() => URL.revokeObjectURL(blobUrl), 1000)
    return { filename, masked }
  } catch (error) {
    throw normalizeError(error)
  }
}

export function toErrorMessage(error: unknown): string {
  return normalizeError(error).message
}

export { http }
