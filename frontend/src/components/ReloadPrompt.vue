<script setup lang="ts">
import { useRegisterSW } from 'virtual:pwa-register/vue'
import { RefreshCw, X } from 'lucide-vue-next'

const { needRefresh, offlineReady, updateServiceWorker } = useRegisterSW({
  onRegisteredSW(_url, registration) {
    if (registration) {
      setInterval(() => {
        void registration.update()
      }, 60 * 60 * 1000)
    }
  },
})

function close() {
  needRefresh.value = false
  offlineReady.value = false
}

async function reload() {
  await updateServiceWorker(true)
}
</script>

<template>
  <div
    v-if="needRefresh || offlineReady"
    class="fixed inset-x-3 bottom-20 z-50 mx-auto flex max-w-md items-center gap-3 rounded-lg border border-line bg-surface p-3 shadow-md md:bottom-6"
  >
    <div class="flex-1 text-[13px] leading-[20px] text-ink">
      <template v-if="needRefresh">衣橱有新版本，点击刷新即可使用。</template>
      <template v-else>已缓存完成，离线也能打开首屏。</template>
    </div>
    <button
      v-if="needRefresh"
      type="button"
      class="flex h-8 items-center gap-1 rounded-md bg-primary px-3 text-[13px] text-white"
      @click="reload"
    >
      <RefreshCw :size="14" :stroke-width="1.5" /> 刷新
    </button>
    <button type="button" class="text-ink-3 hover:text-ink" @click="close">
      <X :size="16" :stroke-width="1.5" />
    </button>
  </div>
</template>
