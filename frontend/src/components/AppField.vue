<script setup lang="ts">
withDefaults(
  defineProps<{
    label?: string
    hint?: string
    error?: string
    required?: boolean
  }>(),
  { required: false },
)
</script>

<template>
  <div class="flex flex-col gap-1.5">
    <label v-if="label" class="flex items-center gap-1 text-[13px] font-medium text-ink-2">
      {{ label }}
      <span v-if="required" class="text-danger">*</span>
    </label>
    <slot />
    <p v-if="error" class="text-[12px] leading-[18px] text-danger">{{ error }}</p>
    <p v-else-if="hint" class="text-[12px] leading-[18px] text-ink-3">{{ hint }}</p>
  </div>
</template>
