<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowLeft,
  CalendarCheck,
  ExternalLink,
  Loader2,
  Pencil,
  Shirt,
  Trash2,
} from 'lucide-vue-next'
import AppButton from '@/components/AppButton.vue'
import AppCard from '@/components/AppCard.vue'
import AppChip from '@/components/AppChip.vue'
import AppEmpty from '@/components/AppEmpty.vue'
import AppImage from '@/components/AppImage.vue'
import AppSkeleton from '@/components/AppSkeleton.vue'
import Money from '@/components/Money.vue'
import { addWearLog, deleteItem, getItem, listWearLogs } from '@/api/items'
import { toErrorMessage } from '@/api/client'
import { useDict } from '@/composables/useDict'
import { useMetaStore } from '@/stores/meta'
import type { ItemDetail, WearLogItem } from '@/types'
import { formatDate, formatDateTime, relativeDays, today, wearText } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const dict = useDict()
const meta = useMetaStore()

const id = computed(() => Number(route.params.id || 0))
const detail = ref<ItemDetail | null>(null)
const logs = ref<WearLogItem[]>([])
const loading = ref(true)
const error = ref('')
const activeImage = ref(0)
const wearOpen = ref(false)
const wearSaving = ref(false)
const wearForm = reactive({ worn_date: today(), occasion: 'daily', remark: '' })

const gallery = computed(() => detail.value?.images ?? [])
const currentImage = computed(() => gallery.value[activeImage.value]?.url || detail.value?.cover || '')

const infoRows = computed(() => {
  const d = detail.value
  if (!d) return []
  return [
    { label: '分类', value: [d.category, d.sub_category].filter(Boolean).join(' / ') || '未填写' },
    { label: '品牌', value: d.brand || '未填写' },
    { label: '尺码', value: d.size || '未填写' },
    { label: '状态', value: d.status?.label || '未填写' },
    { label: '购买日期', value: formatDate(d.purchase_date) },
    { label: '购买渠道', value: d.purchase_channel || '未填写' },
    { label: '存放位置', value: d.storage_location || '未填写' },
    { label: '穿着次数', value: wearText(d.wear_count) },
    { label: '最近穿着', value: relativeDays(d.last_worn_at) },
  ]
})

async function load() {
  loading.value = true
  error.value = ''
  try {
    await dict.ensure()
    const [item, wearLogs] = await Promise.all([getItem(id.value), listWearLogs(id.value).catch(() => [])])
    detail.value = item
    logs.value = wearLogs
    activeImage.value = 0
    meta.applyTitle(item.name)
  } catch (err) {
    error.value = toErrorMessage(err)
  } finally {
    loading.value = false
  }
}

async function checkWear() {
  wearSaving.value = true
  try {
    await addWearLog(id.value, {
      worn_date: wearForm.worn_date,
      occasion: wearForm.occasion,
      remark: wearForm.remark.trim(),
    })
    ElMessage.success('已记录今天穿过它')
    wearOpen.value = false
    wearForm.remark = ''
    await load()
  } catch (err) {
    ElMessage.error(toErrorMessage(err))
  } finally {
    wearSaving.value = false
  }
}

async function remove() {
  try {
    await ElMessageBox.confirm('删除后不可恢复，确定删除这件衣物吗？', '删除衣物', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    })
  } catch {
    return
  }
  try {
    await deleteItem(id.value)
    ElMessage.success('已删除')
    await router.replace('/items')
  } catch (err) {
    ElMessage.error(toErrorMessage(err))
  }
}

onMounted(() => {
  window.scrollTo({ top: 0 })
  void load()
})
</script>

<template>
  <div class="page">
    <div class="mb-4 flex items-center gap-2">
      <button
        type="button"
        class="flex h-10 w-10 items-center justify-center rounded-md text-ink-2 hover:bg-surface-2"
        @click="router.back()"
      >
        <ArrowLeft :size="18" :stroke-width="1.5" />
      </button>
      <h1 class="min-w-0 flex-1 truncate text-[20px] font-semibold leading-7">{{ detail?.name || '衣物详情' }}</h1>
      <div class="flex items-center gap-2">
        <RouterLink v-if="detail" :to="`/items/${detail.id}/edit`">
          <AppButton variant="secondary" size="sm">
            <Pencil :size="14" :stroke-width="1.5" /> 编辑
          </AppButton>
        </RouterLink>
        <AppButton v-if="detail" variant="ghost" size="sm" icon @click="remove">
          <Trash2 :size="16" :stroke-width="1.5" />
        </AppButton>
      </div>
    </div>

    <div v-if="error" class="card mb-4 flex items-center justify-between p-4">
      <p class="text-sm text-danger">{{ error }}</p>
      <AppButton variant="secondary" size="sm" @click="load">重试</AppButton>
    </div>

    <div v-if="loading" class="grid gap-4 lg:grid-cols-[0.9fr_1.1fr]">
      <AppSkeleton height="420px" />
      <AppSkeleton height="420px" />
    </div>

    <div v-else-if="detail" class="grid gap-4 lg:grid-cols-[0.9fr_1.1fr]">
      <div class="flex flex-col gap-3">
        <AppCard padding="none" class="overflow-hidden">
          <AppImage :src="currentImage" :alt="detail.name" rounded="none" :eager="true" />
        </AppCard>
        <div v-if="gallery.length > 1" class="no-scrollbar flex gap-2 overflow-x-auto">
          <button
            v-for="(image, index) in gallery"
            :key="image.id"
            type="button"
            class="h-16 w-12 shrink-0 overflow-hidden rounded-md border-2 transition-colors"
            :class="index === activeImage ? 'border-primary' : 'border-transparent'"
            @click="activeImage = index"
          >
            <img :src="image.thumb_url || image.url" :alt="`缩略图 ${index + 1}`" class="h-full w-full object-cover" />
          </button>
        </div>
        <div class="flex flex-wrap gap-1.5">
          <AppChip v-for="color in detail.color" :key="color.code" dot :label="color.label" :hex="color.hex" />
        </div>
      </div>

      <div class="flex flex-col gap-4">
        <AppCard padding="lg">
          <div class="mb-4 flex items-start justify-between gap-3">
            <div>
              <p class="text-[12px] text-ink-3">{{ detail.category }}{{ detail.sub_category ? ` · ${detail.sub_category}` : '' }}</p>
              <h2 class="mt-1 text-[20px] font-semibold leading-7">{{ detail.name }}</h2>
            </div>
            <div class="text-right">
              <p class="text-[12px] text-ink-3">购买价格</p>
              <Money
                size="lg"
                class="mt-1 block"
                :value="detail.purchase_price"
                :display="detail.purchase_price_display"
                :masked="detail.price_masked"
              />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-x-4 gap-y-4 md:grid-cols-3">
            <div v-for="row in infoRows" :key="row.label" class="min-w-0">
              <p class="label">{{ row.label }}</p>
              <p class="val mt-0.5 truncate">{{ row.value }}</p>
            </div>
          </div>

          <div v-if="detail.season.length || detail.occasion.length" class="mt-4 flex flex-wrap gap-1.5">
            <span
              v-for="season in detail.season"
              :key="`s-${season.code}`"
              class="rounded-full bg-surface-2 px-2.5 py-0.5 text-[12px] text-ink-2"
            >
              {{ season.label }}
            </span>
            <span
              v-for="occasion in detail.occasion"
              :key="`o-${occasion.code}`"
              class="rounded-full bg-surface-2 px-2.5 py-0.5 text-[12px] text-ink-2"
            >
              {{ occasion.label }}
            </span>
          </div>

          <div v-if="detail.tags.length" class="mt-4 flex flex-wrap gap-1.5">
            <RouterLink
              v-for="tag in detail.tags"
              :key="tag.id"
              :to="`/items?tag_id=${tag.id}`"
              class="rounded-full bg-primary-soft px-2.5 py-0.5 text-[12px] text-primary"
            >
              #{{ tag.name }}
            </RouterLink>
          </div>

          <p v-if="detail.remark" class="mt-4 whitespace-pre-wrap rounded-md bg-surface-2 p-3 text-[13px] leading-[22px] text-ink-2">
            {{ detail.remark }}
          </p>

          <div class="mt-4 flex flex-wrap items-center gap-3">
            <a
              v-if="detail.purchase_url"
              :href="detail.purchase_url"
              target="_blank"
              rel="noopener noreferrer"
              class="flex items-center gap-1 text-[13px] text-primary hover:underline"
            >
              订单链接 <ExternalLink :size="13" :stroke-width="1.5" />
            </a>
            <span class="text-[12px] text-ink-3">添加于 {{ formatDateTime(detail.created_at) }}</span>
            <span v-if="detail.cost_per_wear_display" class="text-[12px] text-ink-3">
              单次成本
              <Money
                size="sm"
                tone="muted"
                :value="detail.cost_per_wear"
                :display="detail.cost_per_wear_display"
                :masked="detail.price_masked"
              />
            </span>
          </div>
        </AppCard>

        <AppCard padding="lg">
          <div class="mb-3 flex items-center justify-between">
            <h2 class="text-[15px] font-semibold">穿着记录</h2>
            <AppButton size="sm" @click="wearOpen = true">
              <CalendarCheck :size="14" :stroke-width="1.5" /> 打卡
            </AppButton>
          </div>
          <AppEmpty
            v-if="!logs.length"
            compact
            title="还没有穿着记录"
            description="今天穿了它？打个卡，单次成本会自动计算。"
          />
          <ul v-else class="flex flex-col">
            <li
              v-for="(log, index) in logs"
              :key="log.id"
              class="flex items-start gap-3 py-2.5"
              :class="index !== logs.length - 1 ? 'border-b border-line' : ''"
            >
              <span class="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary-soft text-primary">
                <Shirt :size="14" :stroke-width="1.5" />
              </span>
              <div class="min-w-0 flex-1">
                <p class="text-sm text-ink">{{ formatDate(log.worn_date) }}</p>
                <p class="text-[12px] text-ink-3">
                  {{ log.occasion?.label || '未标注场合' }}
                  <span v-if="log.remark">· {{ log.remark }}</span>
                </p>
              </div>
            </li>
          </ul>
        </AppCard>
      </div>
    </div>

    <AppCard v-else-if="!loading" padding="none">
      <AppEmpty title="这件衣物不存在" description="它可能已经被删除了。" />
    </AppCard>

    <el-dialog v-model="wearOpen" title="穿着打卡" width="380px">
      <div class="flex flex-col gap-3">
        <div class="flex flex-col gap-1.5">
          <label class="text-[13px] text-ink-2">穿着日期</label>
          <el-date-picker v-model="wearForm.worn_date" type="date" value-format="YYYY-MM-DD" class="!w-full" />
        </div>
        <div class="flex flex-col gap-1.5">
          <label class="text-[13px] text-ink-2">场合</label>
          <el-select v-model="wearForm.occasion" class="w-full">
            <el-option v-for="opt in dict.occasions" :key="opt.code" :label="opt.label" :value="opt.code" />
          </el-select>
        </div>
        <div class="flex flex-col gap-1.5">
          <label class="text-[13px] text-ink-2">备注</label>
          <input
            v-model="wearForm.remark"
            type="text"
            maxlength="100"
            placeholder="选填"
            class="h-10 w-full rounded-md bg-surface-2 px-3 text-sm outline-none placeholder:text-ink-3"
          />
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <AppButton variant="ghost" @click="wearOpen = false">取消</AppButton>
          <AppButton :loading="wearSaving" @click="checkWear">
            <Loader2 v-if="wearSaving" :size="14" class="animate-spin" /> 记录
          </AppButton>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
