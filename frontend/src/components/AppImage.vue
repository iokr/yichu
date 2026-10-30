<script setup lang="ts">
import { ref } from 'vue'
import { ImageOff } from 'lucide-vue-next'

withDefaults(
  defineProps<{
    src?: string
    alt?: string
    ratio?: 'photo' | 'square' | 'auto'
    rounded?: 'md' | 'lg' | 'xl' | 'none'
    eager?: boolean
  }>(),
  { src: '', alt: '', ratio: 'photo', rounded: 'lg', eager: false },
)

const loaded = ref(false)
const failed = ref(false)
</script>

<template>
  <div
    class="relative overflow-hidden bg-surface-2"
    :class="[
      ratio === 'photo' ? 'aspect-photo' : ratio === 'square' ? 'aspect-square' : '',
      rounded === 'md' ? 'rounded-md' : rounded === 'lg' ? 'rounded-lg' : rounded === 'xl' ? 'rounded-xl' : '',
    ]"
  >
    <div v-if="!loaded && !failed" class="skeleton absolute inset-0 rounded-none" />
    <img
      v-if="src && !failed"
      :src="src"
      :alt="alt"
      class="h-full w-full object-cover transition-opacity duration-300"
      :class="loaded ? 'opacity-100' : 'opacity-0'"
      :loading="eager ? 'eager' : 'lazy'"
      decoding="async"
      @load="loaded = true"
      @error="failed = true"
    />
    <div v-else-if="failed" class="absolute inset-0 flex items-center justify-center text-ink-3">
      <ImageOff :size="20" :stroke-width="1.5" />
    </div>
    <slot />
  </div>
</template>
