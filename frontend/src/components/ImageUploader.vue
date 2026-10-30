<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { Camera, ImagePlus, RotateCcw, Star, X } from 'lucide-vue-next'
import { useUpload } from '@/composables/useUpload'
import type { ItemImagePayload } from '@/types'

const props = withDefaults(
  defineProps<{
    modelValue: ItemImagePayload[]
    max?: number
    capture?: boolean
  }>(),
  { max: 5, capture: true },
)

const emit = defineEmits<{ (e: 'update:modelValue', value: ItemImagePayload[]): void }>()

const { tasks, error, add, cancel, retry, remove, reset } = useUpload(props.max)

const galleryInput = ref<HTMLInputElement | null>(null)
const cameraInput = ref<HTMLInputElement | null>(null)
const dragging = ref(false)

function openGallery() {
  galleryInput.value?.click()
}

function openCamera() {
  cameraInput.value?.click()
}

async function handleFiles(files: FileList | File[] | null) {
  if (!files || (files instanceof FileList && files.length === 0)) return
  const created = await add(files)
  const next = [...props.modelValue]
  for (const task of created) {
    if (task.status === 'done' && task.result) {
      next.push({ url: task.result.url, thumb_url: task.result.thumb_url })
      remove(task.id)
    }
  }
  if (next.length !== props.modelValue.length) emit('update:modelValue', next.slice(0, props.max))
}

async function onInput(event: Event) {
  const input = event.target as HTMLInputElement
  await handleFiles(input.files)
  input.value = ''
}

function onDrop(event: DragEvent) {
  dragging.value = false
  const files = event.dataTransfer?.files
  if (files?.length) void handleFiles(files)
}

function onPaste(event: ClipboardEvent) {
  const files = Array.from(event.clipboardData?.files || [])
  if (files.length) {
    event.preventDefault()
    void handleFiles(files)
  }
}

function removeImage(index: number) {
  const next = props.modelValue.filter((_, i) => i !== index)
  emit('update:modelValue', next)
}

function makeCover(index: number) {
  if (index <= 0) return
  const next = [...props.modelValue]
  const [picked] = next.splice(index, 1)
  next.unshift(picked)
  emit('update:modelValue', next)
}

onMounted(() => {
  window.addEventListener('paste', onPaste)
})

onBeforeUnmount(() => {
  window.removeEventListener('paste', onPaste)
  reset()
})
</script>

<template>
  <div class="flex flex-col gap-2">
    <div
      class="grid grid-cols-3 gap-2 md:grid-cols-5"
      @dragover.prevent="dragging = true"
      @dragleave.prevent="dragging = false"
      @drop.prevent="onDrop"
    >
      <div
        v-for="(image, index) in modelValue"
        :key="image.url"
        class="group relative aspect-photo overflow-hidden rounded-lg bg-surface-2"
      >
        <img :src="image.thumb_url || image.url" :alt="`图片 ${index + 1}`" class="h-full w-full object-cover" />
        <span
          v-if="index === 0"
          class="absolute left-1.5 top-1.5 rounded-full bg-primary/90 px-1.5 py-0.5 text-[10px] text-white"
        >
          封面
        </span>
        <div class="absolute inset-x-0 bottom-0 flex justify-between gap-1 bg-gradient-to-t from-black/50 to-transparent p-1.5">
          <button
            type="button"
            class="flex h-7 w-7 items-center justify-center rounded-full bg-white/85 text-ink-2 transition-colors hover:text-primary"
            title="设为封面"
            @click="makeCover(index)"
          >
            <Star :size="14" :stroke-width="1.5" />
          </button>
          <button
            type="button"
            class="flex h-7 w-7 items-center justify-center rounded-full bg-white/85 text-ink-2 transition-colors hover:text-danger"
            title="删除"
            @click="removeImage(index)"
          >
            <X :size="14" :stroke-width="1.5" />
          </button>
        </div>
      </div>

      <div
        v-for="task in tasks"
        :key="task.id"
        class="relative aspect-photo overflow-hidden rounded-lg bg-surface-2"
      >
        <img :src="task.preview" alt="上传中" class="h-full w-full object-cover opacity-70" />
        <div class="absolute inset-0 flex flex-col items-center justify-center gap-1 bg-black/35 px-2 text-white">
          <template v-if="task.status === 'uploading'">
            <p class="tabular text-[12px]">{{ task.percent }}%</p>
            <div class="h-1 w-4/5 overflow-hidden rounded-full bg-white/30">
              <div class="h-full bg-white transition-all" :style="{ width: `${task.percent}%` }" />
            </div>
            <button type="button" class="mt-1 text-[11px] underline" @click="cancel(task.id)">取消</button>
          </template>
          <template v-else>
            <p class="line-clamp-2 text-center text-[11px]">{{ task.error || '上传失败' }}</p>
            <button type="button" class="mt-1 flex items-center gap-1 text-[11px] underline" @click="retry(task.id)">
              <RotateCcw :size="12" /> 重试
            </button>
          </template>
        </div>
      </div>

      <button
        v-if="modelValue.length + tasks.length < max"
        type="button"
        class="flex aspect-photo flex-col items-center justify-center gap-1 rounded-lg border border-dashed border-line bg-surface-2/60 text-ink-3 transition-colors hover:border-primary hover:text-primary"
        @click="openGallery"
      >
        <ImagePlus :size="22" :stroke-width="1.5" />
        <span class="text-[12px]">上传图片</span>
        <span class="text-[10px]">{{ modelValue.length + tasks.length }}/{{ max }}</span>
      </button>

      <button
        v-if="capture && modelValue.length + tasks.length < max"
        type="button"
        class="flex aspect-photo flex-col items-center justify-center gap-1 rounded-lg border border-dashed border-line bg-surface-2/60 text-ink-3 transition-colors hover:border-primary hover:text-primary md:hidden"
        @click="openCamera"
      >
        <Camera :size="22" :stroke-width="1.5" />
        <span class="text-[12px]">拍照</span>
      </button>
    </div>

    <div
      v-if="dragging"
      class="rounded-lg border border-dashed border-primary bg-primary-soft px-3 py-2 text-[12px] text-primary"
    >
      松开即可上传
    </div>
    <p v-if="error" class="text-[12px] text-danger">{{ error }}</p>
    <p class="text-[12px] text-ink-3">支持 jpg / png / webp，单张 ≤ 10MB，最多 {{ max }} 张；桌面端可直接拖拽或粘贴。</p>

    <input ref="galleryInput" type="file" accept="image/*" multiple class="hidden" @change="onInput" />
    <input
      ref="cameraInput"
      type="file"
      accept="image/*"
      capture="environment"
      class="hidden"
      @change="onInput"
    />
  </div>
</template>
