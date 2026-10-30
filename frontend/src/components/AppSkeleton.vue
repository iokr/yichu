<script setup lang="ts">
// 多根节点：手动继承 class 等 attributes，避免 fallthrough 警告。
defineOptions({ inheritAttrs: false })

withDefaults(
  defineProps<{
    shape?: 'text' | 'block' | 'card' | 'circle'
    width?: string
    height?: string
    count?: number
  }>(),
  { shape: 'block', count: 1 },
)
</script>

<template>
  <template v-if="shape === 'card'">
    <div v-for="index in count" :key="index" class="card overflow-hidden" v-bind="$attrs">
      <div class="skeleton aspect-photo w-full rounded-none" />
      <div class="space-y-2 p-3">
        <div class="skeleton h-4 w-3/4" />
        <div class="skeleton h-3 w-1/2" />
      </div>
    </div>
  </template>
  <template v-else>
    <div
      v-for="index in count"
      :key="index"
      class="skeleton"
      v-bind="$attrs"
      :class="{
        'rounded-full': shape === 'circle',
        'h-4': shape === 'text' && !height,
      }"
      :style="{ width: width || (shape === 'text' ? '100%' : undefined), height: height || undefined }"
    />
  </template>
</template>
