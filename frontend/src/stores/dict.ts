import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { fetchDicts } from '@/api/meta'
import type { CodeStat, DictOption, Dicts } from '@/types'

const CACHE_KEY = 'yichu.dicts'
const CACHE_TTL = 7 * 24 * 3600 * 1000

interface CachedDicts {
  savedAt: number
  data: Dicts
}

function readCache(): Dicts | null {
  try {
    const raw = localStorage.getItem(CACHE_KEY)
    if (!raw) return null
    const parsed = JSON.parse(raw) as CachedDicts
    if (!parsed?.savedAt || Date.now() - parsed.savedAt > CACHE_TTL) return null
    return parsed.data
  } catch {
    return null
  }
}

const EMPTY: Dicts = {
  colors: [],
  seasons: [],
  occasions: [],
  statuses: [],
  defaults: { status: 'wearing' },
}

export const useDictStore = defineStore('dict', () => {
  const cached = readCache()
  const dicts = ref<Dicts>(cached ?? { ...EMPTY })
  const loaded = ref(!!cached)
  let pending: Promise<Dicts> | null = null

  const colors = computed(() => dicts.value.colors)
  const seasons = computed(() => dicts.value.seasons)
  const occasions = computed(() => dicts.value.occasions)
  const statuses = computed(() => dicts.value.statuses)
  const defaultStatus = computed(() => dicts.value.defaults?.status || 'wearing')

  const colorMap = computed(() => new Map(colors.value.map((c) => [c.code, c])))
  const seasonMap = computed(() => new Map(seasons.value.map((c) => [c.code, c])))
  const occasionMap = computed(() => new Map(occasions.value.map((c) => [c.code, c])))
  const statusMap = computed(() => new Map(statuses.value.map((c) => [c.code, c])))

  async function ensure(force = false): Promise<Dicts> {
    if (loaded.value && !force) return dicts.value
    if (!pending) {
      pending = fetchDicts()
        .then((data) => {
          dicts.value = { ...EMPTY, ...data }
          loaded.value = true
          localStorage.setItem(CACHE_KEY, JSON.stringify({ savedAt: Date.now(), data: dicts.value }))
          return dicts.value
        })
        .finally(() => {
          pending = null
        })
    }
    return pending
  }

  function colorOf(code: string): DictOption {
    return colorMap.value.get(code) ?? { code, label: code }
  }

  function seasonOf(code: string): DictOption {
    return seasonMap.value.get(code) ?? { code, label: code }
  }

  function occasionOf(code: string): DictOption {
    return occasionMap.value.get(code) ?? { code, label: code }
  }

  function statusOf(code: string): DictOption {
    return statusMap.value.get(code) ?? { code, label: code }
  }

  function codeStatLabel(stat: CodeStat): string {
    return stat.label || stat.code
  }

  return {
    dicts,
    loaded,
    colors,
    seasons,
    occasions,
    statuses,
    defaultStatus,
    ensure,
    colorOf,
    seasonOf,
    occasionOf,
    statusOf,
    codeStatLabel,
  }
})
