<script setup lang="ts">
defineProps<{
  title: string
  subtitle?: string
}>()
</script>

<template>
  <header class="mb-4 flex items-start justify-between gap-3 md:mb-6">
    <div class="min-w-0">
      <h1 class="truncate text-[24px] font-semibold leading-[32px] text-ink">{{ title }}</h1>
      <p v-if="subtitle" class="mt-0.5 text-[13px] leading-[20px] text-ink-3">{{ subtitle }}</p>
    </div>
    <div v-if="$slots.actions" class="flex shrink-0 items-center gap-2">
      <slot name="actions" />
    </div>
  </header>
</template>
