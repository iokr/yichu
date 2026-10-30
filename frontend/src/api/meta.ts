import { get } from './client'
import type { Dicts, Meta } from '@/types'

export function fetchMeta() {
  return get<Meta>('/meta')
}

export function fetchDicts() {
  return get<Dicts>('/dicts')
}

export function fetchVersion() {
  return get<{ version: string; buildTime: string }>('/version')
}
