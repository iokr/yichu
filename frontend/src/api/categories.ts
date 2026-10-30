import { del, get, post, put } from './client'
import type { CategoryNode } from '@/types'

export function fetchCategories() {
  return get<CategoryNode[]>('/categories')
}

export interface CategoryPayload {
  name: string
  parent_id?: number
  sort?: number
}

export function createCategory(payload: CategoryPayload) {
  return post<CategoryNode>('/categories', payload)
}

export function updateCategory(id: number, payload: CategoryPayload) {
  return put<CategoryNode>(`/categories/${id}`, payload)
}

export function deleteCategory(id: number) {
  return del<null>(`/categories/${id}`)
}

export function sortCategories(items: { id: number; sort: number }[]) {
  return put<null>('/categories/sort', { items })
}
