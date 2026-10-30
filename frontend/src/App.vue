<script setup lang="ts">
import { RouterView } from 'vue-router'
import ReloadPrompt from '@/components/ReloadPrompt.vue'
</script>

<template>
  <RouterView v-slot="{ Component }">
    <Transition name="page" mode="out-in">
      <component :is="Component" />
    </Transition>
  </RouterView>
  <ReloadPrompt />
</template>
