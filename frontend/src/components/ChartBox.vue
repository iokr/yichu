<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import * as echarts from 'echarts'
import { useTheme } from '@/composables/useTheme'

const props = withDefaults(
  defineProps<{
    option: echarts.EChartsOption
    height?: number
  }>(),
  { height: 280 },
)

const el = ref<HTMLDivElement | null>(null)
const { isDark } = useTheme()

let chart: echarts.ECharts | null = null
let observer: ResizeObserver | null = null

function render() {
  if (!el.value) return
  chart?.dispose()
  chart = echarts.init(el.value, isDark.value ? 'dark' : undefined, { renderer: 'canvas' })
  chart.setOption(props.option)
}

onMounted(() => {
  render()
  if (el.value) {
    observer = new ResizeObserver(() => chart?.resize())
    observer.observe(el.value)
  }
})

watch(
  () => props.option,
  (option) => chart?.setOption(option, true),
  { deep: true },
)

watch(isDark, () => render())

onBeforeUnmount(() => {
  observer?.disconnect()
  chart?.dispose()
  chart = null
})
</script>

<template>
  <div ref="el" :style="{ height: `${height}px`, width: '100%' }" />
</template>
