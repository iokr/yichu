<script setup lang="ts">
withDefaults(
  defineProps<{
    hover?: boolean
    padding?: 'none' | 'sm' | 'md' | 'lg'
  }>(),
  { hover: false, padding: 'md' },
)
</script>

<template>
  <div
    class="card"
    :class="[hover ? 'card-hover' : '', { 'p-3': padding === 'sm', 'p-4': padding === 'md', 'p-5': padding === 'lg' }]"
  >
    <slot />
  </div>
</template>
