import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { fetchMeta } from '@/api/meta'
import type { Meta } from '@/types'

const FALLBACK: Meta = {
  name: '衣橱',
  version: 'dev',
  buildTime: 'unknown',
  hidePrice: false,
}

export const useMetaStore = defineStore('meta', () => {
  const meta = ref<Meta>({ ...FALLBACK })
  const loaded = ref(false)
  const loading = ref(false)

  const appName = computed(() => meta.value.name || FALLBACK.name)
  const hidePrice = computed(() => meta.value.hidePrice)

  async function load(force = false) {
    if (loaded.value && !force) return meta.value
    loading.value = true
    try {
      const data = await fetchMeta()
      if (data) meta.value = { ...FALLBACK, ...data }
      loaded.value = true
      applyTitle()
      return meta.value
    } catch {
      // meta 失败不阻塞首屏，使用兜底文案
      loaded.value = true
      return meta.value
    } finally {
      loading.value = false
    }
  }

  function applyTitle(pageTitle?: string) {
    const base = appName.value
    document.title = pageTitle ? `${pageTitle} · ${base}` : base
  }

  function setHidePrice(value: boolean) {
    meta.value.hidePrice = value
  }

  return { meta, loaded, loading, appName, hidePrice, load, applyTitle, setHidePrice }
})
