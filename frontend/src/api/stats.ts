import { get } from './client'
import type { CodeStat, CostPerWearItem, KeyStat, StatsOverview } from '@/types'

export function fetchOverview() {
  return get<StatsOverview>('/stats/overview')
}

export function fetchByCategory() {
  return get<KeyStat[]>('/stats/by-category')
}

export function fetchByBrand() {
  return get<KeyStat[]>('/stats/by-brand')
}

export function fetchByYear() {
  return get<KeyStat[]>('/stats/by-year')
}

export function fetchBySeason() {
  return get<CodeStat[]>('/stats/by-season')
}

export function fetchByColor() {
  return get<CodeStat[]>('/stats/by-color')
}

export function fetchCostPerWear(limit = 20) {
  return get<CostPerWearItem[]>('/stats/cost-per-wear', { limit })
}
