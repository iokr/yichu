<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ChevronDown, ChevronUp, Loader2, Plus, Save, Tags, Trash2 } from 'lucide-vue-next'
import AppPageHeader from '@/components/AppPageHeader.vue'
import AppButton from '@/components/AppButton.vue'
import AppCard from '@/components/AppCard.vue'
import AppEmpty from '@/components/AppEmpty.vue'
import AppSkeleton from '@/components/AppSkeleton.vue'
import { createCategory, deleteCategory, fetchCategories, sortCategories, updateCategory } from '@/api/categories'
import { toErrorMessage } from '@/api/client'
import { useMetaStore } from '@/stores/meta'
import type { CategoryNode } from '@/types'

const meta = useMetaStore()

const categories = ref<CategoryNode[]>([])
const loading = ref(true)
const saving = ref(false)
const error = ref('')
const selectedId = ref(0)
const createParentId = ref<number | null>(null)
const form = reactive({ name: '', sort: 0 })

const flat = computed(() => categories.value.flatMap((top) => [top, ...(top.children ?? [])]))
const selected = computed(() => flat.value.find((cat) => cat.id === selectedId.value) || null)
const selectedParent = computed(() => {
  const cat = selected.value
  if (!cat || !cat.parent_id) return null
  return categories.value.find((top) => top.id === cat.parent_id) || null
})
const siblings = computed(() => {
  const cat = selected.value
  if (!cat) return []
  if (!cat.parent_id) return categories.value
  return selectedParent.value?.children ?? []
})
const isCreating = computed(() => createParentId.value !== null)

function select(cat: CategoryNode) {
  selectedId.value = cat.id
  createParentId.value = null
  form.name = cat.name
  form.sort = cat.sort
}

function startCreate(parentId: number) {
  createParentId.value = parentId
  selectedId.value = 0
  form.name = ''
  const parent = flat.value.find((cat) => cat.id === parentId)
  const count = parent?.children?.length ?? 0
  form.sort = (count + 1) * 10
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    categories.value = await fetchCategories()
    if (!selectedId.value && categories.value.length) {
      const firstChild = categories.value.flatMap((top) => top.children ?? [])[0]
      select(firstChild ?? categories.value[0])
    }
  } catch (err) {
    error.value = toErrorMessage(err)
  } finally {
    loading.value = false
  }
}

async function save() {
  if (!form.name.trim()) {
    ElMessage.warning('请填写分类名称')
    return
  }
  saving.value = true
  try {
    if (isCreating.value && createParentId.value !== null) {
      const created = await createCategory({ name: form.name.trim(), parent_id: createParentId.value, sort: form.sort })
      ElMessage.success('已新增分类')
      createParentId.value = null
      await load()
      if (created?.id) {
        const found = flat.value.find((cat) => cat.id === created.id)
        if (found) select(found)
      }
    } else if (selected.value) {
      await updateCategory(selected.value.id, { name: form.name.trim(), sort: form.sort })
      ElMessage.success('已保存')
      await load()
    }
  } catch (err) {
    ElMessage.error(toErrorMessage(err))
  } finally {
    saving.value = false
  }
}

async function remove() {
  const cat = selected.value
  if (!cat || cat.is_system || !cat.parent_id) return
  try {
    await ElMessageBox.confirm(`确定删除「${cat.name}」吗？分类下有关联衣物时将无法删除。`, '删除分类', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning',
    })
  } catch {
    return
  }
  try {
    await deleteCategory(cat.id)
    ElMessage.success('已删除')
    selectedId.value = 0
    await load()
  } catch (err) {
    ElMessage.error(toErrorMessage(err))
  }
}

async function move(direction: -1 | 1) {
  const cat = selected.value
  if (!cat) return
  const list = [...siblings.value]
  const index = list.findIndex((item) => item.id === cat.id)
  const target = index + direction
  if (index < 0 || target < 0 || target >= list.length) return
  const other = list[target]
  const payload = [
    { id: cat.id, sort: other.sort },
    { id: other.id, sort: cat.sort },
  ]
  saving.value = true
  try {
    await sortCategories(payload)
    const keepId = cat.id
    await load()
    const found = flat.value.find((item) => item.id === keepId)
    if (found) select(found)
  } catch (err) {
    ElMessage.error(toErrorMessage(err))
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  await load()
  meta.applyTitle('分类管理')
})
</script>

<template>
  <div class="page">
    <AppPageHeader title="分类管理" subtitle="一级分类固定，二级分类可自由增删与排序">
      <template #actions>
        <AppButton size="sm" variant="secondary" @click="startCreate(categories[0]?.id ?? 0)" :disabled="!categories.length">
          <Plus :size="14" :stroke-width="1.5" /> 新增二级
        </AppButton>
      </template>
    </AppPageHeader>

    <div v-if="error" class="card mb-4 flex items-center justify-between p-4">
      <p class="text-sm text-danger">{{ error }}</p>
      <AppButton variant="secondary" size="sm" @click="load">重试</AppButton>
    </div>

    <div class="grid gap-4 lg:grid-cols-[340px_1fr]">
      <AppCard padding="sm" class="h-fit">
        <AppSkeleton v-if="loading" height="320px" />
        <div v-else class="flex flex-col">
          <div v-for="top in categories" :key="top.id" class="py-1">
            <div
              class="flex h-10 cursor-pointer items-center gap-2 rounded-md px-2 text-sm transition-colors"
              :class="selectedId === top.id && !isCreating ? 'bg-primary-soft text-primary' : 'text-ink hover:bg-surface-2'"
              @click="select(top)"
            >
              <Tags :size="15" :stroke-width="1.5" />
              <span class="flex-1 truncate font-medium">{{ top.name }}</span>
              <span class="text-[11px] text-ink-3">{{ top.item_count }}</span>
              <button
                type="button"
                class="flex h-7 w-7 items-center justify-center rounded text-ink-3 hover:bg-surface-2 hover:text-primary"
                title="新增二级分类"
                @click.stop="startCreate(top.id)"
              >
                <Plus :size="14" :stroke-width="1.5" />
              </button>
            </div>

            <div
              v-for="child in top.children"
              :key="child.id"
              class="ml-5 flex h-9 cursor-pointer items-center gap-2 rounded-md pl-3 pr-2 text-[13px] transition-colors"
              :class="[
                selectedId === child.id && !isCreating ? 'bg-primary-soft text-primary' : 'text-ink-2 hover:bg-surface-2',
                selectedId === child.id && !isCreating ? 'border-l-[3px] border-primary pl-[9px]' : 'border-l-[3px] border-transparent',
              ]"
              @click="select(child)"
            >
              <span class="flex-1 truncate">{{ child.name }}</span>
              <span class="text-[11px] text-ink-3">{{ child.item_count }}</span>
            </div>
          </div>
        </div>
      </AppCard>

      <AppCard padding="lg" class="h-fit">
        <template v-if="isCreating">
          <h2 class="mb-4 text-[15px] font-semibold">新增二级分类</h2>
          <div class="flex flex-col gap-4">
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] text-ink-2">上级分类</label>
              <el-select v-model="createParentId" class="w-full">
                <el-option v-for="top in categories" :key="top.id" :label="top.name" :value="top.id" />
              </el-select>
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] text-ink-2">名称</label>
              <input
                v-model="form.name"
                type="text"
                maxlength="50"
                placeholder="例如：毛衣"
                class="h-10 w-full rounded-md bg-surface-2 px-3 text-sm outline-none placeholder:text-ink-3 focus:shadow-[0_0_0_2px_var(--color-primary)_inset]"
              />
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] text-ink-2">排序（越小越靠前）</label>
              <input
                v-model.number="form.sort"
                type="number"
                min="0"
                class="tabular h-10 w-full rounded-md bg-surface-2 px-3 text-sm outline-none focus:shadow-[0_0_0_2px_var(--color-primary)_inset]"
              />
            </div>
            <div class="flex gap-2">
              <AppButton :loading="saving" @click="save">
                <Loader2 v-if="saving" :size="14" class="animate-spin" />
                <Save v-else :size="14" :stroke-width="1.5" /> 保存
              </AppButton>
              <AppButton variant="ghost" @click="createParentId = null">取消</AppButton>
            </div>
          </div>
        </template>

        <template v-else-if="selected">
          <div class="mb-4 flex items-center justify-between">
            <h2 class="text-[15px] font-semibold">
              {{ selected.parent_id ? '编辑二级分类' : '编辑一级分类' }}
            </h2>
            <span class="text-[12px] text-ink-3">{{ selected.item_count }} 件衣物</span>
          </div>
          <div class="flex max-w-md flex-col gap-4">
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] text-ink-2">名称</label>
              <input
                v-model="form.name"
                type="text"
                maxlength="50"
                class="h-10 w-full rounded-md bg-surface-2 px-3 text-sm outline-none focus:shadow-[0_0_0_2px_var(--color-primary)_inset]"
              />
            </div>
            <div class="flex flex-col gap-1.5">
              <label class="text-[13px] text-ink-2">排序（越小越靠前）</label>
              <input
                v-model.number="form.sort"
                type="number"
                min="0"
                class="tabular h-10 w-full rounded-md bg-surface-2 px-3 text-sm outline-none focus:shadow-[0_0_0_2px_var(--color-primary)_inset]"
              />
            </div>
            <div class="flex flex-wrap gap-2">
              <AppButton :loading="saving" @click="save">
                <Loader2 v-if="saving" :size="14" class="animate-spin" />
                <Save v-else :size="14" :stroke-width="1.5" /> 保存
              </AppButton>
              <AppButton variant="secondary" :disabled="saving" @click="move(-1)">
                <ChevronUp :size="14" :stroke-width="1.5" /> 上移
              </AppButton>
              <AppButton variant="secondary" :disabled="saving" @click="move(1)">
                <ChevronDown :size="14" :stroke-width="1.5" /> 下移
              </AppButton>
              <AppButton v-if="selected.parent_id" variant="ghost" @click="remove">
                <Trash2 :size="14" :stroke-width="1.5" /> 删除
              </AppButton>
            </div>
            <p class="text-[12px] leading-[20px] text-ink-3">
              一级分类不可删除，可改名；二级分类在有关联衣物时禁止删除。
            </p>
          </div>
        </template>

        <AppEmpty v-else title="选择一个分类" description="在左侧列表中选择要编辑的分类。" compact />
      </AppCard>
    </div>
  </div>
</template>
