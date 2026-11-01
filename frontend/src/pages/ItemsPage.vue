<script setup lang="ts">
import { computed, nextTick, onMounted, onUnmounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ArrowDownWideNarrow,
  CheckSquare,
  Columns3,
  Download,
  LayoutGrid,
  Pencil,
  Plus,
  RotateCcw,
  Search,
  SlidersHorizontal,
  Table2,
  Trash2,
  X,
} from 'lucide-vue-next'
import AppPageHeader from '@/components/AppPageHeader.vue'
import AppButton from '@/components/AppButton.vue'
import AppCard from '@/components/AppCard.vue'
import AppChip from '@/components/AppChip.vue'
import AppEmpty from '@/components/AppEmpty.vue'
import AppImage from '@/components/AppImage.vue'
import AppSkeleton from '@/components/AppSkeleton.vue'
import ItemCard from '@/components/ItemCard.vue'
import Money from '@/components/Money.vue'
import { batchItems, deleteItem, exportItems, listItems } from '@/api/items'
import { fetchCategories } from '@/api/categories'
import { toErrorMessage } from '@/api/client'
import { useDict } from '@/composables/useDict'
import { useBreakpoint } from '@/composables/useBreakpoint'
import { useMetaStore } from '@/stores/meta'
import type { CategoryNode, ItemListItem } from '@/types'

const route = useRoute()
const router = useRouter()
const dict = useDict()
const meta = useMetaStore()
const { isMobile, isDesktop } = useBreakpoint()

const searchInput = ref<HTMLInputElement | null>(null)
const showFilters = ref(false)
const items = ref<ItemListItem[]>([])
const total = ref(0)
const loading = ref(true)
const error = ref('')
const exporting = ref(false)
const categories = ref<CategoryNode[]>([])
const batchMode = ref(false)
const selected = ref<number[]>([])
const batchStatus = ref('idle')
const batchWorking = ref(false)

const VIEW_KEY = 'yichu.items.view'
const COLS_KEY = 'yichu.items.columns'

type ViewMode = 'grid' | 'table'

const columnOptions = [
  { key: 'category', label: '分类' },
  { key: 'brand', label: '品牌' },
  { key: 'color', label: '颜色' },
  { key: 'size', label: '尺码' },
  { key: 'status', label: '状态' },
  { key: 'purchase_date', label: '购买时间' },
  { key: 'purchase_price', label: '价格' },
  { key: 'wear_count', label: '穿着次数' },
]

const view = ref<ViewMode>(localStorage.getItem(VIEW_KEY) === 'table' ? 'table' : 'grid')
// 手机/平板始终用卡片视图，表格仅用于电脑端。
const effectiveView = computed<ViewMode>(() => (isDesktop.value ? view.value : 'grid'))
const columns = reactive(loadColumns())

const filters = reactive({
  keyword: '',
  category_id: undefined as number | undefined,
  tag_id: undefined as number | undefined,
  status: '',
  season: '',
  purchase_from: '',
  purchase_to: '',
  sort: 'created_at',
  order: 'desc' as 'asc' | 'desc',
  page: 1,
  size: 20,
})

// 购买日期拆成「起 / 止」两个紧凑输入，清空或选择后立即重新查询。
const purchaseFrom = computed<string>({
  get: () => filters.purchase_from,
  set: (value) => {
    filters.purchase_from = value ?? ''
    applyFilters()
  },
})

const purchaseTo = computed<string>({
  get: () => filters.purchase_to,
  set: (value) => {
    filters.purchase_to = value ?? ''
    applyFilters()
  },
})

const sortOptions = [
  { value: 'created_at', label: '最近添加' },
  { value: 'purchase_date', label: '购买时间' },
  { value: 'purchase_price', label: '价格' },
  { value: 'wear_count', label: '穿着次数' },
  { value: 'name', label: '名称' },
]

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / filters.size)))
const hasFilter = computed(
  () =>
    !!filters.keyword ||
    !!filters.category_id ||
    !!filters.tag_id ||
    !!filters.status ||
    !!filters.season ||
    !!filters.purchase_from ||
    !!filters.purchase_to ||
    filters.sort !== 'created_at' ||
    filters.order !== 'desc',
)

let keywordTimer: number | undefined
let internalNav = false

function syncFromQuery() {
  const q = route.query
  filters.keyword = typeof q.keyword === 'string' ? q.keyword : ''
  filters.category_id = q.category_id ? Number(q.category_id) : undefined
  filters.tag_id = q.tag_id ? Number(q.tag_id) : undefined
  filters.status = typeof q.status === 'string' ? q.status : ''
  filters.season = typeof q.season === 'string' ? q.season : ''
  filters.purchase_from = typeof q.purchase_from === 'string' ? q.purchase_from : ''
  filters.purchase_to = typeof q.purchase_to === 'string' ? q.purchase_to : ''
  filters.sort = typeof q.sort === 'string' ? q.sort : 'created_at'
  filters.order = q.order === 'asc' ? 'asc' : 'desc'
  filters.page = q.page ? Math.max(1, Number(q.page)) : 1
}

function syncToQuery() {
  const query: Record<string, string> = {}
  if (filters.keyword) query.keyword = filters.keyword
  if (filters.category_id) query.category_id = String(filters.category_id)
  if (filters.tag_id) query.tag_id = String(filters.tag_id)
  if (filters.status) query.status = filters.status
  if (filters.season) query.season = filters.season
  if (filters.purchase_from) query.purchase_from = filters.purchase_from
  if (filters.purchase_to) query.purchase_to = filters.purchase_to
  if (filters.sort !== 'created_at') query.sort = filters.sort
  if (filters.order !== 'desc') query.order = filters.order
  if (filters.page > 1) query.page = String(filters.page)
  internalNav = true
  void router.replace({ path: '/items', query }).finally(() => {
    internalNav = false
  })
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    const result = await listItems({
      keyword: filters.keyword.trim() || undefined,
      category_id: filters.category_id,
      tag_id: filters.tag_id,
      status: filters.status || undefined,
      season: filters.season || undefined,
      purchase_from: filters.purchase_from || undefined,
      purchase_to: filters.purchase_to || undefined,
      sort: filters.sort,
      order: filters.order,
      page: filters.page,
      size: filters.size,
    })
    items.value = result.list ?? []
    total.value = result.total ?? 0
  } catch (err) {
    error.value = toErrorMessage(err)
    items.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

async function loadCategories() {
  try {
    categories.value = await fetchCategories()
  } catch {
    categories.value = []
  }
}

function applyFilters() {
  filters.page = 1
  syncToQuery()
  void load()
}

function resetFilters() {
  filters.keyword = ''
  filters.category_id = undefined
  filters.tag_id = undefined
  filters.status = ''
  filters.season = ''
  filters.purchase_from = ''
  filters.purchase_to = ''
  filters.sort = 'created_at'
  filters.order = 'desc'
  applyFilters()
}

function pickCategory(id?: number) {
  filters.category_id = id
  applyFilters()
}

function changePage(page: number) {
  if (page < 1 || page > totalPages.value) return
  filters.page = page
  syncToQuery()
  void load()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function loadColumns(): Record<string, boolean> {
  const fallback: Record<string, boolean> = {}
  for (const option of columnOptions) fallback[option.key] = true
  try {
    const raw = localStorage.getItem(COLS_KEY)
    if (!raw) return fallback
    const parsed = JSON.parse(raw) as Record<string, boolean>
    return { ...fallback, ...parsed }
  } catch {
    return fallback
  }
}

function setView(mode: ViewMode) {
  view.value = mode
  localStorage.setItem(VIEW_KEY, mode)
}

function onSelectionChange(rows: ItemListItem[]) {
  selected.value = rows.map((row) => row.id)
}

async function removeOne(item: ItemListItem) {
  try {
    await ElMessageBox.confirm(`确定删除「${item.name}」吗？删除后不可恢复。`, '删除衣物', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    })
  } catch {
    return
  }
  try {
    await deleteItem(item.id)
    ElMessage.success('已删除')
    await load()
  } catch (err) {
    ElMessage.error(toErrorMessage(err))
  }
}

// 电脑端快捷键：`/` 聚焦搜索，`N` 新增，`Esc` 关闭当前浮层。
function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    if (batchMode.value) {
      toggleBatchMode()
      return
    }
    if (showFilters.value) {
      showFilters.value = false
      return
    }
  }
  const target = event.target as HTMLElement | null
  const typing = !!target && (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable)
  if (typing) return
  if (event.key === '/') {
    event.preventDefault()
    searchInput.value?.focus()
    return
  }
  if ((event.key === 'n' || event.key === 'N') && !event.metaKey && !event.ctrlKey && !event.altKey) {
    event.preventDefault()
    void router.push('/items/new')
  }
}

function toggleSelect(id: number) {
  const index = selected.value.indexOf(id)
  if (index >= 0) selected.value.splice(index, 1)
  else selected.value.push(id)
}

function toggleBatchMode() {
  batchMode.value = !batchMode.value
  if (!batchMode.value) selected.value = []
}

function selectAll() {
  if (selected.value.length === items.value.length) selected.value = []
  else selected.value = items.value.map((item) => item.id)
}

async function removeSelected() {
  if (!selected.value.length) return
  try {
    await ElMessageBox.confirm(`确定删除选中的 ${selected.value.length} 件衣物吗？删除后不可恢复。`, '批量删除', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    })
  } catch {
    return
  }
  batchWorking.value = true
  try {
    await batchItems([...selected.value], 'delete')
    ElMessage.success('已删除')
    selected.value = []
    batchMode.value = false
    await load()
  } catch (err) {
    ElMessage.error(toErrorMessage(err))
  } finally {
    batchWorking.value = false
  }
}

async function updateSelectedStatus() {
  if (!selected.value.length) return
  batchWorking.value = true
  try {
    await batchItems([...selected.value], 'status', batchStatus.value)
    ElMessage.success('状态已更新')
    selected.value = []
    batchMode.value = false
    await load()
  } catch (err) {
    ElMessage.error(toErrorMessage(err))
  } finally {
    batchWorking.value = false
  }
}

async function doExport(format: 'csv' | 'xlsx') {
  exporting.value = true
  try {
    const result = await exportItems(
      {
        keyword: filters.keyword.trim() || undefined,
        category_id: filters.category_id,
        tag_id: filters.tag_id,
        status: filters.status || undefined,
        season: filters.season || undefined,
        purchase_from: filters.purchase_from || undefined,
        purchase_to: filters.purchase_to || undefined,
        sort: filters.sort,
        order: filters.order,
      },
      format,
    )
    ElMessage.success(result.masked ? '已导出，价格列已隐藏' : `已导出 ${result.filename}`)
  } catch (err) {
    ElMessage.error(toErrorMessage(err))
  } finally {
    exporting.value = false
  }
}

watch(
  () => filters.keyword,
  () => {
    window.clearTimeout(keywordTimer)
    keywordTimer = window.setTimeout(() => applyFilters(), 320)
  },
)

watch(
  () => route.query,
  () => {
    if (internalNav) return
    syncFromQuery()
    void load()
  },
)

watch(
  columns,
  (value) => localStorage.setItem(COLS_KEY, JSON.stringify(value)),
  { deep: true },
)

onMounted(async () => {
  window.addEventListener('keydown', onKeydown)
  syncFromQuery()
  await Promise.all([load(), loadCategories()])
  meta.applyTitle('衣物')
  if (route.query.focus === 'search') {
    await nextTick()
    searchInput.value?.focus()
  }
})

onUnmounted(() => window.removeEventListener('keydown', onKeydown))
</script>

<template>
  <div class="page">
    <AppPageHeader title="衣物" :subtitle="total > 0 ? `共 ${total} 件` : '记录你的每一件衣服'">
      <template #actions>
        <div class="hidden items-center gap-2 md:flex">
          <el-dropdown trigger="click" :disabled="exporting">
            <AppButton variant="secondary" size="sm" :loading="exporting">
              <Download :size="15" :stroke-width="1.5" /> 导出
            </AppButton>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="doExport('csv')">导出 CSV</el-dropdown-item>
                <el-dropdown-item @click="doExport('xlsx')">导出 Excel</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <AppButton size="sm" @click="router.push('/items/new')">
            <Plus :size="15" :stroke-width="1.5" /> 新增
          </AppButton>
        </div>
        <div v-if="isDesktop" class="flex items-center gap-1 rounded-md bg-surface-2 p-0.5">
          <button
            type="button"
            class="flex h-7 w-7 items-center justify-center rounded"
            :class="view === 'grid' ? 'bg-surface text-ink shadow-sm' : 'text-ink-3'"
            title="卡片视图"
            @click="setView('grid')"
          >
            <LayoutGrid :size="15" :stroke-width="1.5" />
          </button>
          <button
            type="button"
            class="flex h-7 w-7 items-center justify-center rounded"
            :class="view === 'table' ? 'bg-surface text-ink shadow-sm' : 'text-ink-3'"
            title="表格视图"
            @click="setView('table')"
          >
            <Table2 :size="15" :stroke-width="1.5" />
          </button>
        </div>
        <el-popover v-if="isDesktop && effectiveView === 'table'" trigger="click" :width="160" placement="bottom-end">
          <template #reference>
            <AppButton variant="ghost" size="sm" icon title="列设置">
              <Columns3 :size="16" :stroke-width="1.5" />
            </AppButton>
          </template>
          <div class="flex flex-col gap-1.5 py-1">
            <label
              v-for="option in columnOptions"
              :key="option.key"
              class="flex cursor-pointer items-center gap-2 text-[12px] text-ink-2"
            >
              <el-checkbox v-model="columns[option.key]" size="small" />
              {{ option.label }}
            </label>
          </div>
        </el-popover>
        <AppButton v-if="!isMobile" variant="ghost" size="sm" @click="toggleBatchMode">
          <component :is="batchMode ? X : CheckSquare" :size="15" :stroke-width="1.5" />
          {{ batchMode ? '取消' : '批量' }}
        </AppButton>
      </template>
    </AppPageHeader>

    <div class="card mb-4 flex flex-col gap-3 p-3 md:p-4">
      <div class="flex items-center gap-2">
        <div class="relative flex-1">
          <Search
            :size="16"
            :stroke-width="1.5"
            class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-ink-3"
          />
          <input
            ref="searchInput"
            v-model="filters.keyword"
            type="search"
            placeholder="搜索名称、品牌…"
            class="h-10 w-full rounded-md bg-surface-2 pl-9 pr-3 text-sm text-ink outline-none placeholder:text-ink-3 focus:shadow-[0_0_0_2px_var(--color-primary)_inset]"
          />
        </div>
        <AppButton variant="ghost" size="sm" icon class="md:hidden" @click="showFilters = !showFilters">
          <SlidersHorizontal :size="18" :stroke-width="1.5" />
        </AppButton>
        <AppButton variant="ghost" size="sm" icon v-if="hasFilter" class="hidden md:inline-flex" @click="resetFilters">
          <RotateCcw :size="16" :stroke-width="1.5" />
        </AppButton>
      </div>

      <div class="no-scrollbar -mx-1 flex gap-1.5 overflow-x-auto px-1">
        <AppChip :active="!filters.category_id" @click="pickCategory(undefined)">全部</AppChip>
        <AppChip
          v-for="cat in categories"
          :key="cat.id"
          :active="filters.category_id === cat.id"
          @click="pickCategory(cat.id)"
        >
          {{ cat.name }}
        </AppChip>
        <AppChip v-if="filters.tag_id" active @click="filters.tag_id = undefined; applyFilters()">
          #标签筛选 ×
        </AppChip>
      </div>

      <div class="flex flex-wrap items-center gap-2" :class="showFilters || !isMobile ? '' : 'hidden md:flex'">
        <el-select v-model="filters.status" placeholder="状态" clearable size="default" class="!w-[110px]" @change="applyFilters">
          <el-option v-for="opt in dict.statuses" :key="opt.code" :label="opt.label" :value="opt.code" />
        </el-select>
        <el-select v-model="filters.season" placeholder="季节" clearable size="default" class="!w-[110px]" @change="applyFilters">
          <el-option v-for="opt in dict.seasons" :key="opt.code" :label="opt.label" :value="opt.code" />
        </el-select>
        <el-date-picker
          v-model="purchaseFrom"
          type="date"
          value-format="YYYY-MM-DD"
          placeholder="购买起"
          size="default"
          class="filter-date !w-[126px]"
        />
        <span class="px-1 text-[12px] text-ink-3">至</span>
        <el-date-picker
          v-model="purchaseTo"
          type="date"
          value-format="YYYY-MM-DD"
          placeholder="购买止"
          size="default"
          class="filter-date !w-[126px]"
        />
        <el-select v-model="filters.sort" size="default" class="!w-[130px]" @change="applyFilters">
          <el-option v-for="opt in sortOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
        </el-select>
        <AppButton variant="ghost" size="sm" icon @click="filters.order = filters.order === 'desc' ? 'asc' : 'desc'; applyFilters()">
          <ArrowDownWideNarrow :size="16" :stroke-width="1.5" :class="filters.order === 'asc' ? 'rotate-180' : ''" />
        </AppButton>
        <AppButton v-if="isMobile" variant="ghost" size="sm" @click="resetFilters">重置</AppButton>
      </div>
    </div>

    <div v-if="error" class="card mb-4 flex items-center justify-between p-4">
      <p class="text-sm text-danger">{{ error }}</p>
      <AppButton variant="secondary" size="sm" @click="load">重试</AppButton>
    </div>

    <div v-if="loading" class="grid grid-cols-2 gap-3 md:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-6">
      <AppSkeleton v-for="index in 8" :key="index" shape="card" />
    </div>

    <AppCard v-else-if="!items.length" padding="none">
      <AppEmpty
        :title="hasFilter ? '没有找到符合条件的衣物' : '衣橱还空着'"
        :description="hasFilter ? '换个关键词或筛选条件再试试。' : '从第一件衣物开始，慢慢记录下来。'"
      >
        <template #action>
          <AppButton v-if="hasFilter" variant="secondary" @click="resetFilters">清除筛选</AppButton>
          <RouterLink v-else to="/items/new">
            <AppButton>
              <Plus :size="16" :stroke-width="1.5" /> 添加衣物
            </AppButton>
          </RouterLink>
        </template>
      </AppEmpty>
    </AppCard>

    <template v-else>
      <div v-if="effectiveView === 'grid'" class="grid grid-cols-2 gap-3 md:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-6">
        <ItemCard
          v-for="item in items"
          :key="item.id"
          :item="item"
          :selectable="batchMode"
          :selected="selected.includes(item.id)"
          @select="toggleSelect"
        />
      </div>

      <AppCard v-else padding="none" class="overflow-hidden">
        <el-table :data="items" row-key="id" class="items-table" @selection-change="onSelectionChange">
          <el-table-column v-if="batchMode" type="selection" width="46" />
          <el-table-column label="图片" width="72">
            <template #default="{ row }">
              <RouterLink :to="`/items/${row.id}`" class="block h-12 w-9 overflow-hidden rounded">
                <AppImage :src="row.cover" :alt="row.name" rounded="none" ratio="auto" class="h-full" />
              </RouterLink>
            </template>
          </el-table-column>
          <el-table-column label="名称" min-width="190" show-overflow-tooltip>
            <template #default="{ row }">
              <RouterLink :to="`/items/${row.id}`" class="text-[13px] font-medium text-ink hover:text-primary">
                {{ row.name }}
              </RouterLink>
            </template>
          </el-table-column>
          <el-table-column v-if="columns.category" label="分类" min-width="120" show-overflow-tooltip>
            <template #default="{ row }">
              <span class="text-[12px] text-ink-2">{{ row.category }}</span>
              <span v-if="row.sub_category" class="text-[12px] text-ink-3"> / {{ row.sub_category }}</span>
            </template>
          </el-table-column>
          <el-table-column v-if="columns.brand" label="品牌" min-width="100" show-overflow-tooltip>
            <template #default="{ row }">
              <span class="text-[12px] text-ink-2">{{ row.brand || '—' }}</span>
            </template>
          </el-table-column>
          <el-table-column v-if="columns.color" label="颜色" min-width="104">
            <template #default="{ row }">
              <span class="flex flex-wrap gap-1">
                <span
                  v-for="code in row.color"
                  :key="code.code"
                  class="h-3 w-3 rounded-full border border-line"
                  :style="{ backgroundColor: code.hex || 'transparent' }"
                  :title="code.label"
                />
              </span>
            </template>
          </el-table-column>
          <el-table-column v-if="columns.size" label="尺码" width="78">
            <template #default="{ row }">
              <span class="text-[12px] text-ink-2">{{ row.size || '—' }}</span>
            </template>
          </el-table-column>
          <el-table-column v-if="columns.status" label="状态" width="88">
            <template #default="{ row }">
              <span class="rounded-full bg-surface-2 px-2 py-0.5 text-[12px] text-ink-2">{{ row.status.label }}</span>
            </template>
          </el-table-column>
          <el-table-column v-if="columns.purchase_date" label="购买时间" width="112">
            <template #default="{ row }">
              <span class="tabular text-[12px] text-ink-2">{{ row.purchase_date || '—' }}</span>
            </template>
          </el-table-column>
          <el-table-column v-if="columns.purchase_price" label="价格" width="112" align="right">
            <template #default="{ row }">
              <Money :value="row.purchase_price" :display="row.purchase_price_display" :masked="row.price_masked" size="sm" />
            </template>
          </el-table-column>
          <el-table-column v-if="columns.wear_count" label="穿着" width="82" align="right">
            <template #default="{ row }">
              <span class="tabular text-[12px] text-ink-2">{{ row.wear_count }} 次</span>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="104" fixed="right">
            <template #default="{ row }">
              <div class="flex items-center gap-1">
                <RouterLink :to="`/items/${row.id}/edit`">
                  <AppButton variant="ghost" size="sm" icon title="编辑">
                    <Pencil :size="14" :stroke-width="1.5" />
                  </AppButton>
                </RouterLink>
                <AppButton variant="ghost" size="sm" icon title="删除" @click="removeOne(row)">
                  <Trash2 :size="14" :stroke-width="1.5" />
                </AppButton>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </AppCard>

      <div class="mt-6 flex items-center justify-between gap-2">
        <AppButton v-if="batchMode && effectiveView === 'grid'" variant="secondary" size="sm" @click="selectAll">
          {{ selected.length === items.length ? '取消全选' : '全选本页' }}
        </AppButton>
        <span v-else class="text-[12px] text-ink-3">第 {{ filters.page }} / {{ totalPages }} 页</span>
        <div class="flex items-center gap-2">
          <AppButton variant="secondary" size="sm" :disabled="filters.page <= 1" @click="changePage(filters.page - 1)">
            上一页
          </AppButton>
          <AppButton
            variant="secondary"
            size="sm"
            :disabled="filters.page >= totalPages"
            @click="changePage(filters.page + 1)"
          >
            下一页
          </AppButton>
        </div>
      </div>
    </template>

    <div
      v-if="batchMode && selected.length"
      class="fixed inset-x-0 bottom-14 z-30 mx-auto flex max-w-3xl items-center gap-2 border-t border-line bg-surface p-3 shadow-md md:bottom-4 md:rounded-lg md:border"
    >
      <span class="text-[13px] text-ink-2">已选 {{ selected.length }} 件</span>
      <div class="ml-auto flex items-center gap-2">
        <el-select v-model="batchStatus" size="small" class="!w-[100px]">
          <el-option v-for="opt in dict.statuses" :key="opt.code" :label="opt.label" :value="opt.code" />
        </el-select>
        <AppButton size="sm" variant="secondary" :loading="batchWorking" @click="updateSelectedStatus">改状态</AppButton>
        <AppButton size="sm" variant="danger" :loading="batchWorking" @click="removeSelected">
          <Trash2 :size="14" :stroke-width="1.5" /> 删除
        </AppButton>
      </div>
    </div>
  </div>
</template>
