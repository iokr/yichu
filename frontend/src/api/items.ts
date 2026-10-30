import { del, download, get, post, put } from './client'
import type {
  DashboardData,
  ItemDetail,
  ItemListItem,
  ItemPayload,
  ItemQueryParams,
  PageResult,
  TagItem,
  WearLogItem,
} from '@/types'

export function fetchDashboard() {
  return get<DashboardData>('/dashboard')
}

export function listItems(params: ItemQueryParams) {
  return get<PageResult<ItemListItem>>('/items', params as Record<string, unknown>)
}

export function getItem(id: number) {
  return get<ItemDetail>(`/items/${id}`)
}

export function createItem(payload: ItemPayload) {
  return post<ItemDetail>('/items', payload)
}

export function updateItem(id: number, payload: ItemPayload) {
  return put<ItemDetail>(`/items/${id}`, payload)
}

export function deleteItem(id: number) {
  return del<{ id: number }>(`/items/${id}`)
}

export function batchItems(ids: number[], action: 'delete' | 'status', status?: string) {
  return post<{ affected: number }>('/items/batch', { ids, action, status })
}

export function exportItems(params: ItemQueryParams, format: 'csv' | 'xlsx') {
  const suffix = format === 'xlsx' ? 'xlsx' : 'csv'
  return download('/items/export', { ...params, format: suffix }, `yichu-items.${suffix}`)
}

export function listWearLogs(itemId: number) {
  return get<WearLogItem[]>(`/items/${itemId}/wear-logs`)
}

export function addWearLog(itemId: number, payload: { worn_date?: string; occasion?: string; remark?: string }) {
  return post<WearLogItem>(`/items/${itemId}/wear-logs`, payload)
}

export function listTags() {
  return get<TagItem[]>('/tags')
}
