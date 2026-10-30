<script setup lang="ts">
import { Shirt } from 'lucide-vue-next'

withDefaults(
  defineProps<{
    title?: string
    description?: string
    compact?: boolean
  }>(),
  { title: '这里还空着', description: '从第一件衣物开始，慢慢记录下来。', compact: false },
)
</script>

<template>
  <div class="flex flex-col items-center justify-center text-center" :class="compact ? 'gap-2 py-8' : 'gap-3 py-16'">
    <div class="flex h-14 w-14 items-center justify-center rounded-full bg-primary-soft text-primary/60">
      <Shirt :size="26" :stroke-width="1.5" />
    </div>
    <p class="text-sm font-medium text-ink">{{ title }}</p>
    <p class="max-w-[280px] text-[13px] leading-[20px] text-ink-3">{{ description }}</p>
    <div v-if="$slots.action" class="mt-2">
      <slot name="action" />
    </div>
  </div>
</template>
