import { useDictStore } from '@/stores/dict'

/** 字典访问（Pinia 缓存 + localStorage 7 天），首次调用自动拉取。 */
export function useDict() {
  const store = useDictStore()
  void store.ensure()
  return store
}
