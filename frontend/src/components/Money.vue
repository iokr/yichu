<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    value?: number | null
    display?: string
    masked?: boolean
    size?: 'sm' | 'md' | 'lg'
    tone?: 'default' | 'muted' | 'primary'
  }>(),
  { size: 'md', tone: 'default', masked: false },
)

const text = computed(() => {
  if (props.masked) return '¥ •••'
  if (props.display) return props.display
  if (props.value === null || props.value === undefined) return '—'
  return `¥${props.value.toFixed(2)}`
})

const classes = computed(() =>
  [
    'tabular whitespace-nowrap',
    props.size === 'sm' ? 'text-[12px]' : props.size === 'lg' ? 'text-[20px] font-semibold' : 'text-sm',
    props.tone === 'muted' ? 'text-ink-3' : props.tone === 'primary' ? 'text-primary' : 'text-ink',
  ].join(' '),
)
</script>

<template>
  <span :class="classes">{{ text }}</span>
</template>
