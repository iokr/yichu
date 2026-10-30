<script setup lang="ts">
import type { DictOption } from '@/types'

const props = withDefaults(
  defineProps<{
    label?: string
    code?: string
    hex?: string
    active?: boolean
    dot?: boolean
    size?: 'sm' | 'md'
    disabled?: boolean
  }>(),
  { active: false, dot: false, size: 'md', disabled: false },
)

function swatch(option?: DictOption) {
  const hex = option?.hex ?? props.hex ?? ''
  if (option?.code === 'multi' || props.code === 'multi') {
    return 'linear-gradient(135deg,#E6A8B8,#D9C46A,#7C93AE)'
  }
  return hex
}
</script>

<template>
  <button
    type="button"
    :disabled="disabled"
    class="inline-flex items-center gap-1.5 whitespace-nowrap rounded-full border transition-all duration-150 ease-smooth disabled:cursor-not-allowed disabled:opacity-50"
    :class="[
      size === 'sm' ? 'h-7 px-2.5 text-[12px]' : 'h-7 px-3 text-[13px]',
      active
        ? 'border-primary bg-primary text-white'
        : 'border-transparent bg-surface-2 text-ink-2 hover:bg-primary-soft hover:text-ink',
    ]"
  >
    <span
      v-if="dot"
      class="h-2.5 w-2.5 shrink-0 rounded-full ring-1 ring-black/5"
      :style="{ background: swatch() || 'transparent' }"
    />
    <slot>{{ label }}</slot>
  </button>
</template>
