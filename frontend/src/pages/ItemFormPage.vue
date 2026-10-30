<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Check, Loader2, Plus, Tag as TagIcon, X } from 'lucide-vue-next'
import AppButton from '@/components/AppButton.vue'
import AppCard from '@/components/AppCard.vue'
import AppField from '@/components/AppField.vue'
import AppChip from '@/components/AppChip.vue'
import ImageUploader from '@/components/ImageUploader.vue'
import { createItem, getItem, listTags, updateItem } from '@/api/items'
import { fetchCategories } from '@/api/categories'
import { toErrorMessage } from '@/api/client'
import { useDict } from '@/composables/useDict'
import { useBreakpoint } from '@/composables/useBreakpoint'
import { useMetaStore } from '@/stores/meta'
import type { CategoryNode, ItemPayload, TagItem } from '@/types'

const route = useRoute()
const router = useRouter()
const dict = useDict()
const meta = useMetaStore()
const { isMobile } = useBreakpoint()

const id = computed(() => Number(route.params.id || 0))
const isEdit = computed(() => id.value > 0)

const form = reactive<ItemPayload>({
  name: '',
  category_id: null,
  sub_category_id: null,
  brand: '',
  color: [],
  size: '',
  season: [],
  occasion: [],
  purchase_date: '',
  purchase_price: null,
  purchase_channel: '',
  purchase_url: '',
  status: 'wearing',
  storage_location: '',
  remark: '',
  images: [],
  tags: [],
})

const categories = ref<CategoryNode[]>([])
const tagSuggestions = ref<TagItem[]>([])
const tagInput = ref('')
const loading = ref(false)
const saving = ref(false)
const error = ref('')
const fieldErrors = reactive<{ name?: string; category?: string }>({})

const subCategories = computed(() => {
  const top = categories.value.find((cat) => cat.id === form.category_id)
  return top?.children ?? []
})

const priceText = computed({
  get: () => (form.purchase_price === null ? '' : String(form.purchase_price)),
  set: (value: string) => {
    const parsed = Number(value)
    form.purchase_price = value === '' || Number.isNaN(parsed) ? null : parsed
  },
})

function toggle(list: string[], code: string) {
  const index = list.indexOf(code)
  if (index >= 0) list.splice(index, 1)
  else list.push(code)
}

function pickCategory(value: number | null) {
  form.category_id = value
  form.sub_category_id = null
}

// 深链直接打开表单时 history 里没有站内上一页，回退到详情/列表，避免退出应用。
function goBack() {
  const back = window.history.state?.back as string | null | undefined
  if (typeof back === 'string' && back.startsWith('/') && back !== route.fullPath) {
    router.back()
    return
  }
  void router.replace(isEdit.value ? `/items/${id.value}` : '/items')
}

function addTag() {
  const names = tagInput.value
    .split(/[,，\s]+/)
    .map((item) => item.trim())
    .filter(Boolean)
  for (const name of names) {
    if (name.length > 20) {
      ElMessage.warning('单个标签不超过 20 个字')
      continue
    }
    if (form.tags.length >= 20) {
      ElMessage.warning('最多 20 个标签')
      break
    }
    if (!form.tags.includes(name)) form.tags.push(name)
  }
  tagInput.value = ''
}

function removeTag(name: string) {
  form.tags = form.tags.filter((tag) => tag !== name)
}

function onTagKeydown(event: KeyboardEvent) {
  if (event.key === ',' || event.key === '，') {
    event.preventDefault()
    addTag()
  }
}

function suggestTag(name: string) {
  if (!form.tags.includes(name) && form.tags.length < 20) form.tags.push(name)
  tagInput.value = ''
}

function validate(): boolean {
  fieldErrors.name = undefined
  fieldErrors.category = undefined
  if (!form.name.trim()) fieldErrors.name = '请填写衣物名称'
  if (form.name.trim().length > 100) fieldErrors.name = '名称不超过 100 个字'
  if (!form.category_id) fieldErrors.category = '请选择一级分类'
  if (form.purchase_price !== null && form.purchase_price < 0) {
    ElMessage.warning('价格不能为负数')
    return false
  }
  return !fieldErrors.name && !fieldErrors.category
}

function buildPayload(): ItemPayload {
  return {
    ...form,
    name: form.name.trim(),
    brand: form.brand.trim(),
    size: form.size.trim(),
    purchase_channel: form.purchase_channel.trim(),
    purchase_url: form.purchase_url.trim(),
    storage_location: form.storage_location.trim(),
    remark: form.remark.trim(),
    images: form.images.map((image) => ({ url: image.url, thumb_url: image.thumb_url })),
    tags: [...form.tags],
  }
}

async function save() {
  if (!validate()) return
  saving.value = true
  error.value = ''
  try {
    const payload = buildPayload()
    if (isEdit.value) {
      await updateItem(id.value, payload)
      ElMessage.success('已保存')
      await router.replace(`/items/${id.value}`)
    } else {
      const created = await createItem(payload)
      ElMessage.success('已添加')
      await router.replace(`/items/${created.id}`)
    }
  } catch (err) {
    error.value = toErrorMessage(err)
  } finally {
    saving.value = false
  }
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    await dict.ensure()
    const [cats, tags] = await Promise.all([fetchCategories(), listTags().catch(() => [] as TagItem[])])
    categories.value = cats
    tagSuggestions.value = tags
    if (isEdit.value) {
      const detail = await getItem(id.value)
      form.name = detail.name
      form.category_id = detail.category_id
      form.sub_category_id = detail.sub_category_id
      form.brand = detail.brand
      form.color = detail.color.map((item) => item.code)
      form.size = detail.size
      form.season = detail.season.map((item) => item.code)
      form.occasion = detail.occasion.map((item) => item.code)
      form.purchase_date = detail.purchase_date || ''
      form.purchase_price = detail.purchase_price
      form.purchase_channel = detail.purchase_channel
      form.purchase_url = detail.purchase_url
      form.status = detail.status.code || 'wearing'
      form.storage_location = detail.storage_location
      form.remark = detail.remark
      form.images = detail.images.map((image) => ({ url: image.url, thumb_url: image.thumb_url }))
      form.tags = detail.tags.map((tag) => tag.name)
      meta.applyTitle(`编辑 ${detail.name}`)
    } else {
      form.status = dict.defaultStatus
      meta.applyTitle('新增衣物')
    }
  } catch (err) {
    error.value = toErrorMessage(err)
  } finally {
    loading.value = false
  }
}

watch(
  () => form.category_id,
  () => {
    if (form.sub_category_id && !subCategories.value.some((cat) => cat.id === form.sub_category_id)) {
      form.sub_category_id = null
    }
  },
)

onMounted(() => {
  window.scrollTo({ top: 0 })
  void load()
})
</script>

<template>
  <div class="page pb-32 md:pb-8">
    <div
      class="-mx-3 sticky top-0 z-20 mb-4 flex items-center gap-2 bg-bg/95 px-3 py-2 backdrop-blur md:static md:mx-0 md:mb-6 md:bg-transparent md:px-0 md:py-0 md:backdrop-blur-none"
    >
      <button
        type="button"
        class="flex h-10 w-10 items-center justify-center rounded-md text-ink-2 hover:bg-surface-2"
        @click="goBack"
      >
        <ArrowLeft :size="18" :stroke-width="1.5" />
      </button>
      <h1 class="text-[20px] font-semibold leading-7">{{ isEdit ? '编辑衣物' : '新增衣物' }}</h1>
      <div class="ml-auto hidden md:block">
        <AppButton :loading="saving" :disabled="loading" @click="save">
          <Check :size="16" :stroke-width="1.5" /> 保存
        </AppButton>
      </div>
    </div>

    <div v-if="error" class="card mb-4 flex items-center justify-between p-4">
      <p class="text-sm text-danger">{{ error }}</p>
      <AppButton variant="secondary" size="sm" @click="load">重试</AppButton>
    </div>

    <div class="grid gap-4 lg:grid-cols-[1.05fr_0.95fr]">
      <div class="flex flex-col gap-4">
        <AppCard padding="lg">
          <h2 class="mb-3 text-[15px] font-semibold">图片</h2>
          <ImageUploader v-model="form.images" :max="5" />
        </AppCard>

        <AppCard padding="lg">
          <h2 class="mb-4 section-title !text-[15px]">基本信息</h2>
          <div class="flex flex-col gap-4">
            <AppField label="名称" required :error="fieldErrors.name">
              <input
                v-model="form.name"
                type="text"
                maxlength="100"
                placeholder="例如：白色针织开衫"
                class="h-10 w-full rounded-md bg-surface-2 px-3 text-sm outline-none placeholder:text-ink-3 focus:shadow-[0_0_0_2px_var(--color-primary)_inset]"
              />
            </AppField>

            <div class="grid grid-cols-2 gap-3">
              <AppField label="一级分类" required :error="fieldErrors.category">
                <el-select
                  :model-value="form.category_id ?? undefined"
                  placeholder="选择分类"
                  class="w-full"
                  @update:model-value="pickCategory"
                >
                  <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
                </el-select>
              </AppField>
              <AppField label="二级分类">
                <el-select
                  v-model="form.sub_category_id"
                  placeholder="可选"
                  clearable
                  class="w-full"
                  :disabled="!subCategories.length"
                >
                  <el-option v-for="cat in subCategories" :key="cat.id" :label="cat.name" :value="cat.id" />
                </el-select>
              </AppField>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <AppField label="品牌">
                <input
                  v-model="form.brand"
                  type="text"
                  maxlength="50"
                  placeholder="选填"
                  class="h-10 w-full rounded-md bg-surface-2 px-3 text-sm outline-none placeholder:text-ink-3 focus:shadow-[0_0_0_2px_var(--color-primary)_inset]"
                />
              </AppField>
              <AppField label="尺码">
                <input
                  v-model="form.size"
                  type="text"
                  maxlength="20"
                  placeholder="S / M / L / 38"
                  class="h-10 w-full rounded-md bg-surface-2 px-3 text-sm outline-none placeholder:text-ink-3 focus:shadow-[0_0_0_2px_var(--color-primary)_inset]"
                />
              </AppField>
            </div>

            <AppField label="颜色" hint="可多选，用于筛选与色点展示">
              <div class="flex flex-wrap gap-1.5">
                <AppChip
                  v-for="opt in dict.colors"
                  :key="opt.code"
                  dot
                  :label="opt.label"
                  :code="opt.code"
                  :hex="opt.hex"
                  :active="form.color.includes(opt.code)"
                  @click="toggle(form.color, opt.code)"
                />
              </div>
            </AppField>

            <AppField label="状态">
              <div class="flex flex-wrap gap-1.5">
                <AppChip
                  v-for="opt in dict.statuses"
                  :key="opt.code"
                  :label="opt.label"
                  :active="form.status === opt.code"
                  @click="form.status = opt.code"
                />
              </div>
            </AppField>
          </div>
        </AppCard>

        <AppCard padding="lg">
          <h2 class="mb-4 section-title !text-[15px]">购买信息</h2>
          <div class="flex flex-col gap-4">
            <div class="grid grid-cols-2 gap-3">
              <AppField label="购买日期">
                <el-date-picker
                  v-model="form.purchase_date"
                  type="date"
                  value-format="YYYY-MM-DD"
                  placeholder="选择日期"
                  class="!w-full"
                />
              </AppField>
              <AppField label="价格" :hint="meta.hidePrice ? '价格可见性已关闭，保存后仍可统计' : '单位：元'">
                <div class="relative">
                  <span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-ink-3">¥</span>
                  <input
                    v-model="priceText"
                    type="number"
                    inputmode="decimal"
                    min="0"
                    step="0.01"
                    placeholder="0.00"
                    class="tabular h-10 w-full rounded-md bg-surface-2 pl-7 pr-3 text-sm outline-none placeholder:text-ink-3 focus:shadow-[0_0_0_2px_var(--color-primary)_inset]"
                  />
                </div>
              </AppField>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <AppField label="购买渠道">
                <input
                  v-model="form.purchase_channel"
                  type="text"
                  maxlength="50"
                  placeholder="淘宝 / 线下…"
                  class="h-10 w-full rounded-md bg-surface-2 px-3 text-sm outline-none placeholder:text-ink-3 focus:shadow-[0_0_0_2px_var(--color-primary)_inset]"
                />
              </AppField>
              <AppField label="存放位置">
                <input
                  v-model="form.storage_location"
                  type="text"
                  maxlength="50"
                  placeholder="衣柜第二层…"
                  class="h-10 w-full rounded-md bg-surface-2 px-3 text-sm outline-none placeholder:text-ink-3 focus:shadow-[0_0_0_2px_var(--color-primary)_inset]"
                />
              </AppField>
            </div>
            <AppField label="订单链接">
              <input
                v-model="form.purchase_url"
                type="url"
                maxlength="500"
                placeholder="https://"
                class="h-10 w-full rounded-md bg-surface-2 px-3 text-sm outline-none placeholder:text-ink-3 focus:shadow-[0_0_0_2px_var(--color-primary)_inset]"
              />
            </AppField>
          </div>
        </AppCard>
      </div>

      <div class="flex flex-col gap-4">
        <AppCard padding="lg">
          <h2 class="mb-4 section-title !text-[15px]">穿着场景</h2>
          <div class="flex flex-col gap-4">
            <AppField label="季节">
              <div class="flex flex-wrap gap-1.5">
                <AppChip
                  v-for="opt in dict.seasons"
                  :key="opt.code"
                  :label="opt.label"
                  :active="form.season.includes(opt.code)"
                  @click="toggle(form.season, opt.code)"
                />
              </div>
            </AppField>
            <AppField label="场合">
              <div class="flex flex-wrap gap-1.5">
                <AppChip
                  v-for="opt in dict.occasions"
                  :key="opt.code"
                  :label="opt.label"
                  :active="form.occasion.includes(opt.code)"
                  @click="toggle(form.occasion, opt.code)"
                />
              </div>
            </AppField>
          </div>
        </AppCard>

        <AppCard padding="lg">
          <h2 class="mb-4 section-title !text-[15px]">标签与备注</h2>
          <div class="flex flex-col gap-4">
            <AppField label="标签" hint="回车或逗号分隔，最多 20 个">
              <div class="flex flex-wrap items-center gap-1.5 rounded-md bg-surface-2 p-2">
                <span
                  v-for="tag in form.tags"
                  :key="tag"
                  class="inline-flex h-7 items-center gap-1 rounded-full bg-primary-soft px-2.5 text-[12px] text-primary"
                >
                  {{ tag }}
                  <button type="button" class="text-primary/70 hover:text-primary" @click="removeTag(tag)">
                    <X :size="12" :stroke-width="2" />
                  </button>
                </span>
                <input
                  v-model="tagInput"
                  type="text"
                  placeholder="添加标签"
                  class="h-7 min-w-[80px] flex-1 bg-transparent px-1 text-sm outline-none placeholder:text-ink-3"
                  @keydown.enter.prevent="addTag"
                  @keydown="onTagKeydown"
                  @blur="addTag"
                />
              </div>
            </AppField>
            <div v-if="tagSuggestions.length" class="flex flex-wrap gap-1.5">
              <button
                v-for="tag in tagSuggestions.filter((item) => !form.tags.includes(item.name)).slice(0, 12)"
                :key="tag.id"
                type="button"
                class="inline-flex h-7 items-center gap-1 rounded-full bg-surface-2 px-2.5 text-[12px] text-ink-3 hover:text-primary"
                @click="suggestTag(tag.name)"
              >
                <TagIcon :size="11" :stroke-width="1.5" /> {{ tag.name }}
              </button>
            </div>
            <AppField label="备注">
              <textarea
                v-model="form.remark"
                rows="4"
                maxlength="500"
                placeholder="材质、搭配、保养提醒…"
                class="w-full resize-none rounded-md bg-surface-2 p-3 text-sm leading-[22px] outline-none placeholder:text-ink-3 focus:shadow-[0_0_0_2px_var(--color-primary)_inset]"
              />
            </AppField>
          </div>
        </AppCard>
      </div>
    </div>

    <div
      class="fixed inset-x-0 bottom-0 z-30 border-t border-line bg-surface p-3 pb-[calc(12px+env(safe-area-inset-bottom))] md:hidden"
    >
      <AppButton block size="lg" :loading="saving" :disabled="loading" @click="save">
        <Loader2 v-if="saving" :size="16" class="animate-spin" />
        <Check v-else :size="16" :stroke-width="1.5" />
        {{ isEdit ? '保存修改' : '保存衣物' }}
      </AppButton>
    </div>

    <div v-if="loading" class="fixed inset-0 z-20 flex items-center justify-center bg-bg/60">
      <p class="flex items-center gap-2 text-sm text-ink-2">
        <Plus :size="16" class="animate-pulse" /> 正在加载…
      </p>
    </div>
  </div>
</template>
