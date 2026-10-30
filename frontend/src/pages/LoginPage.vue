<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Eye, EyeOff, Loader2, Shirt } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useMetaStore } from '@/stores/meta'
import { toErrorMessage } from '@/api/client'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const meta = useMetaStore()

const username = ref('')
const password = ref('')
const showPassword = ref(false)
const error = ref('')
const submitting = ref(false)

const canSubmit = computed(() => username.value.trim().length > 0 && password.value.length > 0 && !submitting.value)

onMounted(() => {
  void meta.load()
})

async function submit() {
  if (!canSubmit.value) return
  error.value = ''
  submitting.value = true
  try {
    await auth.login(username.value.trim(), password.value)
    ElMessage.success('欢迎回来')
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/'
    await router.replace(redirect)
  } catch (err) {
    error.value = toErrorMessage(err)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="flex min-h-screen bg-bg">
    <section class="relative hidden w-[45%] items-center justify-center overflow-hidden bg-primary-soft lg:flex">
      <div class="absolute -left-24 top-1/4 h-72 w-72 rounded-full bg-primary/10" />
      <div class="absolute -bottom-24 right-0 h-80 w-80 rounded-full bg-accent/10" />
      <div class="relative flex flex-col items-center gap-4 px-8 text-center">
        <span class="flex h-16 w-16 items-center justify-center rounded-2xl bg-surface text-primary shadow-sm">
          <Shirt :size="32" :stroke-width="1.5" />
        </span>
        <h1 class="text-[32px] font-semibold leading-10 text-ink">{{ meta.appName }}</h1>
        <p class="max-w-[300px] text-sm leading-[22px] text-ink-2">
          把衣服一件件记下来，穿什么就不再是每天早晨的难题。
        </p>
      </div>
    </section>

    <section class="flex flex-1 flex-col px-6 pb-10 pt-16 lg:items-center lg:justify-center lg:pt-0">
      <div class="mx-auto w-full max-w-[380px]">
        <div class="mb-8 lg:hidden">
          <span class="mb-3 flex h-12 w-12 items-center justify-center rounded-xl bg-primary-soft text-primary">
            <Shirt :size="24" :stroke-width="1.5" />
          </span>
          <h1 class="text-[24px] font-semibold leading-8">{{ meta.appName }}</h1>
          <p class="mt-1 text-[13px] text-ink-3">记录每一件衣服的故事</p>
        </div>

        <form class="flex flex-col gap-4" @submit.prevent="submit">
          <div class="flex flex-col gap-1.5">
            <label class="text-[13px] font-medium text-ink-2">用户名</label>
            <input
              v-model="username"
              type="text"
              autocomplete="username"
              placeholder="请输入用户名"
              class="h-11 w-full rounded-md border-none bg-surface-2 px-3 text-sm text-ink outline-none transition-shadow placeholder:text-ink-3 focus:shadow-[0_0_0_2px_var(--color-primary)_inset]"
            />
          </div>

          <div class="flex flex-col gap-1.5">
            <label class="text-[13px] font-medium text-ink-2">密码</label>
            <div class="relative">
              <input
                v-model="password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                placeholder="请输入密码"
                class="h-11 w-full rounded-md border-none bg-surface-2 px-3 pr-11 text-sm text-ink outline-none transition-shadow placeholder:text-ink-3 focus:shadow-[0_0_0_2px_var(--color-primary)_inset]"
              />
              <button
                type="button"
                class="absolute right-2 top-1/2 flex h-8 w-8 -translate-y-1/2 items-center justify-center text-ink-3"
                @click="showPassword = !showPassword"
              >
                <component :is="showPassword ? EyeOff : Eye" :size="16" :stroke-width="1.5" />
              </button>
            </div>
          </div>

          <p v-if="error" class="rounded-md bg-danger/10 px-3 py-2 text-[13px] text-danger">{{ error }}</p>

          <button
            type="submit"
            class="mt-2 flex h-11 w-full items-center justify-center gap-2 rounded-md bg-primary text-sm font-medium text-white transition-colors hover:bg-primary-hover active:scale-[.99] disabled:cursor-not-allowed disabled:opacity-60"
            :disabled="!canSubmit"
          >
            <Loader2 v-if="submitting" :size="16" class="animate-spin" />
            {{ submitting ? '登录中…' : '登录' }}
          </button>
        </form>

        <p class="mt-6 text-center text-[12px] text-ink-3">首次使用请用配置的管理员账号登录，并及时修改密码。</p>
      </div>
    </section>
  </div>
</template>
