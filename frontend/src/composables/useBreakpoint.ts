import { computed, onUnmounted, ref } from 'vue'

export type Device = 'mobile' | 'tablet' | 'desktop'

function detect(width: number): Device {
  if (width < 768) return 'mobile'
  if (width < 1200) return 'tablet'
  return 'desktop'
}

const width = ref(typeof window === 'undefined' ? 1440 : window.innerWidth)
const device = ref<Device>(detect(width.value))

let listeners = 0
let timer: number | undefined

function onResize() {
  window.clearTimeout(timer)
  timer = window.setTimeout(() => {
    width.value = window.innerWidth
    device.value = detect(width.value)
  }, 80)
}

/** 响应式断点：手机 <768，平板 768~1199，电脑 ≥1200（禁止 UA 判断）。 */
export function useBreakpoint() {
  listeners += 1
  if (listeners === 1) window.addEventListener('resize', onResize, { passive: true })

  onUnmounted(() => {
    listeners -= 1
    if (listeners === 0) window.removeEventListener('resize', onResize)
  })

  return {
    width: computed(() => width.value),
    device: computed(() => device.value),
    isMobile: computed(() => device.value === 'mobile'),
    isTablet: computed(() => device.value === 'tablet'),
    isDesktop: computed(() => device.value === 'desktop'),
  }
}
