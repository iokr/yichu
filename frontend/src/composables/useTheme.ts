import { computed, ref, watch } from 'vue'

export type ThemeMode = 'system' | 'light' | 'dark'

const STORAGE_KEY = 'yichu.theme'
const media = window.matchMedia('(prefers-color-scheme: dark)')

function readMode(): ThemeMode {
  const raw = localStorage.getItem(STORAGE_KEY)
  return raw === 'light' || raw === 'dark' ? raw : 'system'
}

const mode = ref<ThemeMode>(readMode())
const systemDark = ref(media.matches)

media.addEventListener('change', (event) => {
  systemDark.value = event.matches
})

const resolved = computed<'light' | 'dark'>(() => {
  if (mode.value === 'system') return systemDark.value ? 'dark' : 'light'
  return mode.value
})

watch(
  resolved,
  (value) => {
    document.documentElement.classList.toggle('dark', value === 'dark')
    document
      .querySelector('meta[name="theme-color"]')
      ?.setAttribute('content', value === 'dark' ? '#141312' : '#A67C89')
  },
  { immediate: true },
)

watch(mode, (value) => {
  if (value === 'system') localStorage.removeItem(STORAGE_KEY)
  else localStorage.setItem(STORAGE_KEY, value)
})

/** 深色模式：跟随系统 + 手动切换，选择存 localStorage。 */
export function useTheme() {
  return {
    mode,
    resolved,
    isDark: computed(() => resolved.value === 'dark'),
    setMode(value: ThemeMode) {
      mode.value = value
    },
  }
}
