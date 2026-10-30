<script setup lang="ts">
import { computed } from 'vue'

type Variant = 'primary' | 'secondary' | 'ghost' | 'danger'
type Size = 'sm' | 'md' | 'lg'

const props = withDefaults(
  defineProps<{
    variant?: Variant
    size?: Size
    type?: 'button' | 'submit' | 'reset'
    block?: boolean
    loading?: boolean
    disabled?: boolean
    icon?: boolean
  }>(),
  { variant: 'primary', size: 'md', type: 'button', block: false, loading: false, disabled: false, icon: false },
)

const classes = computed(() => {
  const variant = {
    primary: 'bg-primary text-white hover:bg-primary-hover active:scale-[.98]',
    secondary: 'bg-surface border border-line text-primary hover:bg-primary-soft active:scale-[.98]',
    ghost: 'bg-transparent text-ink-2 hover:bg-surface-2 active:scale-[.98]',
    danger: 'bg-danger text-white hover:opacity-90 active:scale-[.98]',
  }[props.variant]
  const size = {
    sm: 'h-8 px-3 text-[13px]',
    md: 'h-10 px-4 text-sm md:h-10',
    lg: 'h-11 px-5 text-sm',
  }[props.size]
  return [
    'inline-flex items-center justify-center gap-2 rounded-md font-medium transition-all duration-150 ease-smooth select-none',
    'disabled:cursor-not-allowed disabled:opacity-50',
    props.icon ? '!px-2.5' : '',
    props.block ? 'w-full' : '',
    variant,
    size,
  ].join(' ')
})
</script>

<template>
  <button :type="type" :class="classes" :disabled="disabled || loading">
    <span
      v-if="loading"
      class="h-3.5 w-3.5 animate-spin rounded-full border-2 border-current border-t-transparent"
    />
    <slot />
  </button>
</template>
