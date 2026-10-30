import { post } from './client'
import type { TokenResponse } from '@/types'

export interface LoginPayload {
  username: string
  password: string
}

export function login(payload: LoginPayload) {
  return post<TokenResponse>('/auth/login', payload)
}

export function logout() {
  return post<null>('/auth/logout')
}
