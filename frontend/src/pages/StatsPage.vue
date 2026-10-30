<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { EyeOff } from 'lucide-vue-next'
import AppPageHeader from '@/components/AppPageHeader.vue'
import AppCard from '@/components/AppCard.vue'
import AppEmpty from '@/components/AppEmpty.vue'
import AppImage from '@/components/AppImage.vue'
import AppSkeleton from '@/components/AppSkeleton.vue'
import AppButton from '@/components/AppButton.vue'
import ChartBox from '@/components/ChartBox.vue'
import Money from '@/components/Money.vue'
import {
  fetchByBrand,
  fetchByCategory,
  fetchByColor,
  fetchBySeason,
  fetchByYear,
  fetchCostPerWear,
  fetchOverview,
} from '@/api/stats'
import { toErrorMessage } from '@/api/client'
import { useTheme } from '@/composables/useTheme'
import { useMetaStore } from '@/stores/meta'
import type { EChartsOption } from 'echarts'
import type { CodeStat, CostPerWearItem, KeyStat, StatsOverview } from '@/types'
import { axisColors, chartPalette, tooltipStyle } from '@/utils/chart'
import { percentText } from '@/utils/format'

const meta = useMetaStore()
const { isDark } = useTheme()

const overview = ref<StatsOverview | null>(null)
const byCategory = ref<KeyStat[]>([])
const byBrand = ref<KeyStat[]>([])
const byYear = ref<KeyStat[]>([])
const bySeason = ref<CodeStat[]>([])
const byColor = ref<CodeStat[]>([])
const costPerWear = ref<CostPerWearItem[]>([])
const loading = ref(true)
const error = ref('')

const masked = computed(() => overview.value?.price_masked ?? false)

const cards = computed(() => {
  const o = overview.value
  return [
    { label: '衣物总数', value: `${o?.total_count ?? 0}`, unit: '件' },
    { label: '在穿', value: `${o?.wearing_count ?? 0}`, unit: '件' },
    { label: '闲置', value: `${o?.idle_count ?? 0}`, unit: '件' },
    { label: '累计穿着', value: `${o?.total_wear_count ?? 0}`, unit: '次' },
  ]
})

function baseGrid() {
  const colors = axisColors(isDark.value)
  return { colors, tooltip: tooltipStyle(isDark.value) }
}

const categoryOption = computed<EChartsOption>(() => {
  const { colors, tooltip } = baseGrid()
  const rows = byCategory.value.slice(0, 8)
  return {
    grid: { left: 8, right: 24, top: 8, bottom: 8, containLabel: true },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' }, ...tooltip },
    xAxis: {
      type: 'value',
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { lineStyle: { color: colors.split } },
      axisLabel: { color: colors.text, fontSize: 11 },
    },
    yAxis: {
      type: 'category',
      inverse: true,
      data: rows.map((row) => row.label),
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: colors.text, fontSize: 12 },
    },
    series: [
      {
        type: 'bar',
        barWidth: 14,
        data: rows.map((row, index) => ({
          value: row.count,
          itemStyle: { color: chartPalette[index % chartPalette.length], borderRadius: [4, 4, 4, 4] },
        })),
        label: { show: true, position: 'right', color: colors.text, fontSize: 11 },
      },
    ],
  }
})

const brandOption = computed<EChartsOption>(() => {
  const { colors, tooltip } = baseGrid()
  const rows = byBrand.value.slice(0, 8).reverse()
  return {
    grid: { left: 8, right: 24, top: 8, bottom: 8, containLabel: true },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' }, ...tooltip },
    xAxis: {
      type: 'value',
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { lineStyle: { color: colors.split } },
      axisLabel: { color: colors.text, fontSize: 11 },
    },
    yAxis: {
      type: 'category',
      data: rows.map((row) => row.label || '未填写'),
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: colors.text, fontSize: 12 },
    },
    series: [
      {
        type: 'bar',
        barWidth: 14,
        data: rows.map((row) => row.count),
        itemStyle: { color: chartPalette[1], borderRadius: [4, 4, 4, 4] },
        label: { show: true, position: 'right', color: colors.text, fontSize: 11 },
      },
    ],
  }
})

const yearOption = computed<EChartsOption>(() => {
  const { colors, tooltip } = baseGrid()
  const rows = byYear.value
  return {
    grid: { left: 8, right: 16, top: 24, bottom: 8, containLabel: true },
    tooltip: {
      trigger: 'axis',
      ...tooltip,
      formatter: '{b}：{c} 件',
    },
    xAxis: {
      type: 'category',
      data: rows.map((row) => row.label || row.key),
      axisLine: { lineStyle: { color: colors.line } },
      axisTick: { show: false },
      axisLabel: { color: colors.text, fontSize: 11 },
    },
    yAxis: {
      type: 'value',
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { lineStyle: { color: colors.split } },
      axisLabel: { color: colors.text, fontSize: 11 },
    },
    series: [
      {
        type: 'bar',
        barWidth: 22,
        data: rows.map((row) => row.count),
        itemStyle: { color: chartPalette[0], borderRadius: [6, 6, 0, 0] },
        label: { show: true, position: 'top', color: colors.text, fontSize: 11 },
      },
    ],
  }
})

const seasonOption = computed<EChartsOption>(() => {
  const { colors, tooltip } = baseGrid()
  const rows = bySeason.value.filter((row) => row.count > 0)
  return {
    color: chartPalette,
    tooltip: { trigger: 'item', ...tooltip, formatter: '{b}：{c} 件（{d}%）' },
    legend: {
      bottom: 0,
      icon: 'circle',
      itemWidth: 8,
      itemHeight: 8,
      textStyle: { color: colors.text, fontSize: 12 },
    },
    series: [
      {
        type: 'pie',
        radius: ['48%', '70%'],
        center: ['50%', '44%'],
        itemStyle: { borderWidth: 2, borderColor: isDark.value ? '#1C1B1A' : '#FFFFFF' },
        label: { show: false },
        data: rows.map((row) => ({ name: row.label, value: row.count })),
      },
    ],
  }
})

const colorOption = computed<EChartsOption>(() => {
  const { colors, tooltip } = baseGrid()
  const rows = byColor.value.filter((row) => row.count > 0).slice(0, 10)
  return {
    grid: { left: 8, right: 24, top: 8, bottom: 8, containLabel: true },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' }, ...tooltip },
    xAxis: {
      type: 'value',
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { lineStyle: { color: colors.split } },
      axisLabel: { color: colors.text, fontSize: 11 },
    },
    yAxis: {
      type: 'category',
      inverse: true,
      data: rows.map((row) => row.label),
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: colors.text, fontSize: 12 },
    },
    series: [
      {
        type: 'bar',
        barWidth: 12,
        data: rows.map((row) => ({
          value: row.count,
          itemStyle: {
            color: row.code === 'multi' ? chartPalette[2] : row.hex || chartPalette[3],
            borderRadius: [4, 4, 4, 4],
            borderColor: isDark.value ? '#2E2C2A' : '#EAE8E5',
            borderWidth: 1,
          },
        })),
        label: { show: true, position: 'right', color: colors.text, fontSize: 11 },
      },
    ],
  }
})

async function load() {
  loading.value = true
  error.value = ''
  try {
    const [ov, cat, brand, year, season, color, cost] = await Promise.all([
      fetchOverview(),
      fetchByCategory(),
      fetchByBrand(),
      fetchByYear(),
      fetchBySeason(),
      fetchByColor(),
      fetchCostPerWear(10).catch(() => [] as CostPerWearItem[]),
    ])
    overview.value = ov
    byCategory.value = cat
    byBrand.value = brand
    byYear.value = year
    bySeason.value = season
    byColor.value = color
    costPerWear.value = cost
  } catch (err) {
    error.value = toErrorMessage(err)
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await load()
  meta.applyTitle('统计')
})
</script>

<template>
  <div class="page">
    <AppPageHeader title="统计" subtitle="看清衣橱的构成与花费">
      <template #actions>
        <span v-if="masked" class="flex items-center gap-1 text-[12px] text-ink-3">
          <EyeOff :size="13" :stroke-width="1.5" /> 价格已隐藏
        </span>
      </template>
    </AppPageHeader>

    <div v-if="error" class="card mb-4 flex items-center justify-between p-4">
      <p class="text-sm text-danger">{{ error }}</p>
      <AppButton variant="secondary" size="sm" @click="load">重试</AppButton>
    </div>

    <div class="grid grid-cols-2 gap-3 md:grid-cols-4">
      <template v-if="loading">
        <AppSkeleton v-for="index in 4" :key="index" height="84px" />
      </template>
      <div v-for="card in cards" v-else :key="card.label" class="card p-4">
        <p class="text-[12px] text-ink-3">{{ card.label }}</p>
        <p class="mt-1 flex items-baseline gap-1">
          <span class="tabular text-[24px] font-semibold leading-8">{{ card.value }}</span>
          <span class="text-[12px] text-ink-3">{{ card.unit }}</span>
        </p>
      </div>
    </div>

    <div class="mt-4 grid gap-4 md:grid-cols-2">
      <AppCard padding="lg">
        <div class="mb-2 flex items-baseline justify-between">
          <h2 class="text-[15px] font-semibold">总花费</h2>
          <span class="text-[12px] text-ink-3">均价</span>
        </div>
        <div class="flex items-end justify-between">
          <Money
            size="lg"
            :value="overview?.total_amount ?? null"
            :display="overview?.total_amount_display"
            :masked="masked"
          />
          <Money
            size="sm"
            tone="muted"
            :value="overview?.avg_price ?? null"
            :display="overview?.avg_price_display"
            :masked="masked"
          />
        </div>
      </AppCard>

      <AppCard padding="lg">
        <div class="mb-3 flex items-center justify-between">
          <h2 class="text-[15px] font-semibold">在穿比例</h2>
          <span class="text-[12px] text-ink-3">
            {{ overview ? percentText((overview.wearing_count / Math.max(overview.total_count, 1)) * 100) : '—' }}
          </span>
        </div>
        <div class="h-2 overflow-hidden rounded-full bg-surface-2">
          <div
            class="h-full rounded-full bg-primary transition-all"
            :style="{
              width: overview ? `${Math.min((overview.wearing_count / Math.max(overview.total_count, 1)) * 100, 100)}%` : '0%',
            }"
          />
        </div>
        <div class="mt-3 flex items-center justify-between text-[12px] text-ink-3">
          <span>在穿 {{ overview?.wearing_count ?? 0 }}</span>
          <span>闲置 {{ overview?.idle_count ?? 0 }}</span>
        </div>
      </AppCard>
    </div>

    <div class="mt-4 grid gap-4 lg:grid-cols-2">
      <AppCard padding="lg">
        <h2 class="mb-2 text-[15px] font-semibold">分类分布</h2>
        <AppSkeleton v-if="loading" height="260px" />
        <AppEmpty v-else-if="!byCategory.length" compact title="暂无数据" description="添加衣物后即可查看。" />
        <ChartBox v-else :option="categoryOption" :height="260" />
      </AppCard>

      <AppCard padding="lg">
        <h2 class="mb-2 text-[15px] font-semibold">品牌 TOP 8</h2>
        <AppSkeleton v-if="loading" height="260px" />
        <AppEmpty v-else-if="!byBrand.length" compact title="暂无品牌数据" description="在衣物里填写品牌后会出现在这里。" />
        <ChartBox v-else :option="brandOption" :height="260" />
      </AppCard>

      <AppCard padding="lg">
        <h2 class="mb-2 text-[15px] font-semibold">年度购买</h2>
        <AppSkeleton v-if="loading" height="260px" />
        <AppEmpty v-else-if="!byYear.length" compact title="暂无购买年份" description="填写购买日期后即可统计。" />
        <ChartBox v-else :option="yearOption" :height="260" />
        <p v-if="masked" class="mt-2 text-[12px] text-ink-3">价格隐藏时，年度维度按件数展示。</p>
      </AppCard>

      <AppCard padding="lg">
        <h2 class="mb-2 text-[15px] font-semibold">季节分布</h2>
        <AppSkeleton v-if="loading" height="260px" />
        <AppEmpty v-else-if="!bySeason.some((row) => row.count > 0)" compact title="暂无季节数据" description="给衣物标注季节后即可统计。" />
        <ChartBox v-else :option="seasonOption" :height="260" />
      </AppCard>
    </div>

    <div class="mt-4 grid gap-4 lg:grid-cols-2">
      <AppCard padding="lg">
        <h2 class="mb-2 text-[15px] font-semibold">颜色分布</h2>
        <AppSkeleton v-if="loading" height="260px" />
        <AppEmpty v-else-if="!byColor.some((row) => row.count > 0)" compact title="暂无颜色数据" description="给衣物选颜色后即可统计。" />
        <ChartBox v-else :option="colorOption" :height="260" />
      </AppCard>

      <AppCard padding="lg">
        <div class="mb-2 flex items-center justify-between">
          <h2 class="text-[15px] font-semibold">单次穿着成本</h2>
          <span class="text-[12px] text-ink-3">越低越划算</span>
        </div>
        <div v-if="loading">
          <AppSkeleton v-for="index in 4" :key="index" height="48px" class="mb-2" />
        </div>
        <AppEmpty
          v-else-if="masked"
          compact
          title="价格已隐藏"
          description="开启价格可见性后，才能查看单次成本排行。"
        />
        <AppEmpty v-else-if="!costPerWear.length" compact title="暂无排行" description="记录穿着次数后即可计算。" />
        <ul v-else class="flex flex-col">
          <li
            v-for="(row, index) in costPerWear"
            :key="row.id"
            class="flex items-center gap-3 py-2"
            :class="index !== costPerWear.length - 1 ? 'border-b border-line' : ''"
          >
            <span class="tabular w-4 text-[12px] text-ink-3">{{ index + 1 }}</span>
            <RouterLink :to="`/items/${row.id}`" class="h-10 w-8 shrink-0 overflow-hidden rounded">
              <AppImage :src="row.cover" :alt="row.name" rounded="none" ratio="auto" class="h-full" />
            </RouterLink>
            <div class="min-w-0 flex-1">
              <p class="truncate text-[13px] text-ink">{{ row.name }}</p>
              <p class="text-[11px] text-ink-3">穿过 {{ row.wear_count }} 次</p>
            </div>
            <Money
              size="sm"
              :value="row.cost_per_wear"
              :display="row.cost_per_wear_display"
              :masked="row.price_masked"
            />
          </li>
        </ul>
      </AppCard>
    </div>
  </div>
</template>
