<script setup lang="ts">
import { computed } from 'vue'
import { Check } from 'lucide-vue-next'
import AppImage from './AppImage.vue'
import Money from './Money.vue'
import type { ItemListItem } from '@/types'

const props = withDefaults(
  defineProps<{
    item: ItemListItem
    selectable?: boolean
    selected?: boolean
  }>(),
  { selectable: false, selected: false },
)

defineEmits<{ (e: 'select', id: number): void }>()

const meta = computed(() => {
  const parts = [props.item.category]
  if (props.item.sub_category) parts.push(props.item.sub_category)
  if (props.item.purchase_date) parts.push(props.item.purchase_date.slice(0, 7).replace('-', '/'))
  return parts.filter(Boolean).join(' · ')
})
</script>

<template>
  <div class="group relative">
    <RouterLink
      :to="`/items/${item.id}`"
      class="card card-hover block overflow-hidden"
      :class="selected ? 'ring-2 ring-primary ring-offset-1 ring-offset-bg' : ''"
    >
      <div class="relative">
        <AppImage :src="item.cover" :alt="item.name" rounded="none" ratio="photo" />
        <span
          v-if="item.status && item.status.code !== 'wearing'"
          class="absolute left-2 top-2 rounded-full bg-black/45 px-2 py-0.5 text-[11px] text-white backdrop-blur-sm"
        >
          {{ item.status.label }}
        </span>
      </div>
      <div class="flex items-start justify-between gap-2 p-3">
        <div class="min-w-0">
          <p class="truncate text-sm font-medium text-ink">{{ item.name }}</p>
          <p class="mt-0.5 truncate text-[12px] text-ink-3">{{ meta || '未填充分类' }}</p>
        </div>
        <Money
          size="sm"
          tone="muted"
          :value="item.purchase_price"
          :display="item.purchase_price_display"
          :masked="item.price_masked"
        />
      </div>
    </RouterLink>

    <button
      v-if="selectable"
      type="button"
      class="absolute right-2 top-2 flex h-6 w-6 items-center justify-center rounded-full border transition-colors"
      :class="selected ? 'border-primary bg-primary text-white' : 'border-white/70 bg-black/25 text-transparent'"
      @click.prevent.stop="$emit('select', item.id)"
    >
      <Check :size="14" :stroke-width="2" />
    </button>
  </div>
</template>
