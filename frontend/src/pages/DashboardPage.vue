<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { ArrowRight, BarChart3, Plus, Shirt, Sparkles } from 'lucide-vue-next'
import AppPageHeader from '@/components/AppPageHeader.vue'
import AppCard from '@/components/AppCard.vue'
import AppSkeleton from '@/components/AppSkeleton.vue'
import AppEmpty from '@/components/AppEmpty.vue'
import AppButton from '@/components/AppButton.vue'
import ItemCard from '@/components/ItemCard.vue'
import ChartBox from '@/components/ChartBox.vue'
import Money from '@/components/Money.vue'
import { fetchDashboard } from '@/api/items'
import { toErrorMessage } from '@/api/client'
import type { EChartsOption } from 'echarts'
import type { DashboardData } from '@/types'
import { useTheme } from '@/composables/useTheme'
import { useMetaStore } from '@/stores/meta'
import { axisColors, chartPalette, tooltipStyle } from '@/utils/chart'
import { percentText } from '@/utils/format'

const meta = useMetaStore()
const { isDark } = useTheme()

const data = ref<DashboardData | null>(null)
const loading = ref(true)
const error = ref('')

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return '夜深了'
  if (hour < 11) return '早上好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  return '晚上好'
})

const todayText = computed(() =>
  new Date().toLocaleDateString('zh-CN', { month: 'long', day: 'numeric', weekday: 'long' }),
)

const cards = computed(() => {
  const d = data.value
  return [
    { label: '衣物总数', value: d?.total_count ?? 0, suffix: '件', to: '/items' },
    { label: '在穿', value: d?.wearing_count ?? 0, suffix: '件', to: '/items?status=wearing' },
    { label: '闲置', value: d?.idle_count ?? 0, suffix: '件', to: '/items?status=idle' },
    { label: '累计穿着', value: d?.total_wear_count ?? 0, suffix: '次', to: '/items' },
  ]
})

const categoryOption = computed<EChartsOption>(() => {
  const colors = axisColors(isDark.value)
  const stats = (data.value?.category_stats ?? []).filter((item) => item.count > 0)
  return {
    color: chartPalette,
    tooltip: {
      trigger: 'item',
      ...tooltipStyle(isDark.value),
      formatter: '{b}：{c} 件（{d}%）',
    },
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
        radius: ['52%', '72%'],
        center: ['50%', '44%'],
        avoidLabelOverlap: true,
        itemStyle: { borderWidth: 2, borderColor: isDark.value ? '#1C1B1A' : '#FFFFFF' },
        label: { show: false },
        emphasis: { scale: true, scaleSize: 4 },
        data: stats.map((item) => ({ name: item.label, value: item.count })),
      },
    ],
  }
})

const statusRows = computed(() => data.value?.status_stats ?? [])

onMounted(load)

async function load() {
  loading.value = true
  error.value = ''
  try {
    data.value = await fetchDashboard()
    meta.applyTitle('仪表盘')
  } catch (err) {
    error.value = toErrorMessage(err)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page">
    <AppPageHeader :title="`${greeting}，今天穿什么？`" :subtitle="todayText">
      <template #actions>
        <RouterLink to="/items/new" class="hidden md:block">
          <AppButton>
            <Plus :size="16" :stroke-width="1.5" /> 新增衣物
          </AppButton>
        </RouterLink>
      </template>
    </AppPageHeader>

    <div v-if="error" class="card mb-4 flex items-center justify-between p-4">
      <p class="text-sm text-danger">{{ error }}</p>
      <AppButton variant="secondary" size="sm" @click="load">重试</AppButton>
    </div>

    <div class="grid grid-cols-2 gap-3 lg:grid-cols-4">
      <template v-if="loading">
        <AppSkeleton v-for="index in 4" :key="index" height="92px" />
      </template>
      <RouterLink v-for="card in cards" v-else :key="card.label" :to="card.to">
        <div class="card card-hover flex h-full flex-col justify-between p-4">
          <p class="text-[12px] text-ink-3">{{ card.label }}</p>
          <p class="mt-2 flex items-baseline gap-1">
            <span class="tabular text-[28px] font-semibold leading-9 text-ink">{{ card.value }}</span>
            <span class="text-[12px] text-ink-3">{{ card.suffix }}</span>
          </p>
        </div>
      </RouterLink>
    </div>

    <div class="mt-4 grid gap-4 lg:grid-cols-3">
      <AppCard class="lg:col-span-2" padding="lg">
        <div class="mb-2 flex items-center justify-between">
          <h2 class="section-title">分类分布</h2>
          <RouterLink to="/stats" class="flex items-center gap-1 text-[12px] text-ink-3 hover:text-primary">
            查看统计 <ArrowRight :size="12" :stroke-width="1.5" />
          </RouterLink>
        </div>
        <AppSkeleton v-if="loading" height="260px" />
        <AppEmpty
          v-else-if="!data?.category_stats?.length"
          compact
          title="还没有分类数据"
          description="添加第一件衣物后，这里会出现分类占比。"
        />
        <ChartBox v-else :option="categoryOption" :height="260" />
      </AppCard>

      <div class="flex flex-col gap-4">
        <AppCard padding="lg">
          <p class="text-[12px] text-ink-3">{{ data?.year ?? new Date().getFullYear() }} 年支出</p>
          <p class="mt-2">
            <Money
              size="lg"
              :value="data?.year_amount ?? null"
              :display="data?.year_amount_display"
              :masked="data?.price_masked"
            />
          </p>
          <div class="divider my-3" />
          <p class="text-[12px] text-ink-3">累计支出</p>
          <p class="mt-1">
            <Money
              :value="data?.total_amount ?? null"
              :display="data?.total_amount_display"
              :masked="data?.price_masked"
            />
          </p>
        </AppCard>

        <AppCard padding="lg" class="flex-1">
          <h2 class="mb-3 text-[15px] font-semibold">状态分布</h2>
          <div class="flex flex-col gap-2.5">
            <div v-for="stat in statusRows" :key="stat.code" class="flex items-center gap-2">
              <span class="w-12 shrink-0 text-[12px] text-ink-2">{{ stat.label }}</span>
              <div class="h-1.5 flex-1 overflow-hidden rounded-full bg-surface-2">
                <div class="h-full rounded-full bg-primary" :style="{ width: `${Math.min(stat.percent, 100)}%` }" />
              </div>
              <span class="tabular w-10 shrink-0 text-right text-[12px] text-ink-3">{{ stat.count }}</span>
            </div>
          </div>
        </AppCard>
      </div>
    </div>

    <section class="mt-6">
      <div class="mb-3 flex items-center justify-between">
        <h2 class="section-title">最近添加</h2>
        <RouterLink to="/items" class="flex items-center gap-1 text-[12px] text-ink-3 hover:text-primary">
          全部衣物 <ArrowRight :size="12" :stroke-width="1.5" />
        </RouterLink>
      </div>

      <div v-if="loading" class="grid grid-cols-2 gap-3 md:grid-cols-4 xl:grid-cols-6">
        <AppSkeleton v-for="index in 6" :key="index" shape="card" />
      </div>
      <AppCard v-else-if="!data?.recent?.length" padding="none">
        <AppEmpty title="衣橱还空着" description="拍一张照片，记下第一件衣服。">
          <template #action>
            <RouterLink to="/items/new">
              <AppButton>
                <Sparkles :size="16" :stroke-width="1.5" /> 添加第一件
              </AppButton>
            </RouterLink>
          </template>
        </AppEmpty>
      </AppCard>
      <div v-else class="grid grid-cols-2 gap-3 md:grid-cols-4 xl:grid-cols-6">
        <ItemCard v-for="item in data.recent" :key="item.id" :item="item" />
      </div>
    </section>

    <section class="mt-6 grid gap-3 md:hidden">
      <RouterLink to="/items" class="card flex items-center justify-between p-4">
        <span class="flex items-center gap-2 text-sm text-ink"><Shirt :size="16" :stroke-width="1.5" /> 浏览全部衣物</span>
        <ArrowRight :size="16" :stroke-width="1.5" class="text-ink-3" />
      </RouterLink>
      <RouterLink to="/stats" class="card flex items-center justify-between p-4">
        <span class="flex items-center gap-2 text-sm text-ink">
          <BarChart3 :size="16" :stroke-width="1.5" /> 查看穿用统计
        </span>
        <ArrowRight :size="16" :stroke-width="1.5" class="text-ink-3" />
      </RouterLink>
    </section>

    <p v-if="data && data.total_count > 0" class="mt-4 text-center text-[12px] text-ink-3">
      在穿的 {{ data.wearing_count }} 件衣物，占全部 {{ percentText((data.wearing_count / Math.max(data.total_count, 1)) * 100) }}
    </p>
  </div>
</template>
