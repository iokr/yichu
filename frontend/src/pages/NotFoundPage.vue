<script setup lang="ts">
import { useRouter } from 'vue-router'
import { Home } from 'lucide-vue-next'
import AppButton from '@/components/AppButton.vue'

const router = useRouter()
</script>

<template>
  <div class="flex min-h-screen flex-col items-center justify-center gap-4 bg-bg px-6 text-center">
    <p class="text-[64px] font-semibold leading-none text-primary/20">404</p>
    <p class="text-[18px] font-semibold text-ink">这个页面不存在</p>
    <p class="max-w-[300px] text-[13px] leading-[22px] text-ink-3">链接可能已经失效，回到首页继续逛逛吧。</p>
    <AppButton @click="router.replace('/')">
      <Home :size="16" :stroke-width="1.5" /> 回到首页
    </AppButton>
  </div>
</template>
