import { download, get, put } from './client'
import type { SettingsData } from '@/types'

export function fetchSettings() {
  return get<SettingsData>('/settings')
}

export function updateSettings(payload: { hide_price?: boolean }) {
  return put<SettingsData>('/settings', payload)
}

export function downloadBackup() {
  return download('/backup', undefined, 'yichu-backup.json')
}
