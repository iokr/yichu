<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import {
  BarChart3,
  ChevronLeft,
  LayoutDashboard,
  LogOut,
  Menu,
  Moon,
  Plus,
  Search,
  Settings,
  Shirt,
  Sun,
  Tags,
  WifiOff,
  X,
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useMetaStore } from '@/stores/meta'
import { useBreakpoint } from '@/composables/useBreakpoint'
import { useTheme } from '@/composables/useTheme'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const meta = useMetaStore()
const { isMobile, isDesktop } = useBreakpoint()
const { isDark, setMode } = useTheme()

const collapsed = ref(localStorage.getItem('yichu.sidebar') === '1')
const drawerOpen = ref(false)
const online = ref(navigator.onLine)

const navs = [
  { to: '/', label: '仪表盘', icon: LayoutDashboard, exact: true },
  { to: '/items', label: '衣物', icon: Shirt },
  { to: '/stats', label: '统计', icon: BarChart3 },
  { to: '/categories', label: '分类', icon: Tags },
  { to: '/settings', label: '设置', icon: Settings },
]

const mobileNavs = [
  { to: '/', label: '首页', icon: LayoutDashboard, exact: true },
  { to: '/items', label: '衣物', icon: Shirt },
  { to: '/stats', label: '统计', icon: BarChart3 },
  { to: '/settings', label: '我的', icon: Settings },
]

const pageTitle = computed(() => (route.meta.title as string) || meta.appName)
const showFab = computed(() => route.path === '/items')
// 全屏页面（新增/编辑衣物）自己带返回栏与底部保存条，隐藏手机顶栏与 Tab。
const isFullscreen = computed(() => route.meta.fullscreen === true)

function isActive(to: string, exact = false) {
  if (exact) return route.path === to
  return route.path === to || route.path.startsWith(`${to}/`)
}

function toggleCollapse() {
  collapsed.value = !collapsed.value
  localStorage.setItem('yichu.sidebar', collapsed.value ? '1' : '0')
}

function toggleTheme() {
  setMode(isDark.value ? 'light' : 'dark')
}

async function signOut() {
  await auth.logout()
  router.replace('/login')
}

function goSearch() {
  router.push({ path: '/items', query: { focus: 'search' } })
}

function onOnline() {
  online.value = true
}

function onOffline() {
  online.value = false
}

onMounted(() => {
  window.addEventListener('online', onOnline)
  window.addEventListener('offline', onOffline)
})

onBeforeUnmount(() => {
  window.removeEventListener('online', onOnline)
  window.removeEventListener('offline', onOffline)
})

watch(
  () => route.fullPath,
  () => {
    drawerOpen.value = false
  },
)
</script>

<template>
  <div class="min-h-screen bg-bg text-ink">
    <div
      v-if="!online"
      class="sticky top-0 z-40 flex items-center justify-center gap-1.5 bg-warning/15 py-1 text-[12px] text-warning"
    >
      <WifiOff :size="13" :stroke-width="1.5" /> 当前离线，数据可能不是最新
    </div>

    <div class="flex">
      <!-- 电脑：侧边栏 -->
      <aside
        v-if="isDesktop"
        class="sticky top-0 hidden h-screen shrink-0 flex-col border-r border-line bg-surface transition-all duration-200 ease-smooth md:flex"
        :class="collapsed ? 'w-[72px]' : 'w-[228px]'"
      >
        <div class="flex h-16 items-center gap-2 px-4" :class="collapsed ? 'justify-center' : ''">
          <span class="flex h-9 w-9 items-center justify-center rounded-lg bg-primary-soft text-primary">
            <Shirt :size="20" :stroke-width="1.5" />
          </span>
          <span v-if="!collapsed" class="truncate text-[16px] font-semibold">{{ meta.appName }}</span>
        </div>
        <nav class="flex flex-1 flex-col gap-1 px-2 py-2">
          <RouterLink
            v-for="nav in navs"
            :key="nav.to"
            :to="nav.to"
            class="flex h-10 items-center gap-3 rounded-md px-3 text-sm transition-colors"
            :class="
              isActive(nav.to, nav.exact)
                ? 'bg-primary-soft font-medium text-primary'
                : 'text-ink-2 hover:bg-surface-2 hover:text-ink'
            "
          >
            <component :is="nav.icon" :size="18" :stroke-width="1.5" class="shrink-0" />
            <span v-if="!collapsed" class="truncate">{{ nav.label }}</span>
          </RouterLink>
        </nav>
        <div class="flex flex-col gap-1 border-t border-line p-2">
          <button
            type="button"
            class="flex h-10 items-center gap-3 rounded-md px-3 text-sm text-ink-2 transition-colors hover:bg-surface-2 hover:text-ink"
            @click="toggleTheme"
          >
            <component :is="isDark ? Sun : Moon" :size="18" :stroke-width="1.5" class="shrink-0" />
            <span v-if="!collapsed">{{ isDark ? '浅色模式' : '深色模式' }}</span>
          </button>
          <button
            type="button"
            class="flex h-10 items-center gap-3 rounded-md px-3 text-sm text-ink-2 transition-colors hover:bg-surface-2 hover:text-ink"
            @click="signOut"
          >
            <LogOut :size="18" :stroke-width="1.5" class="shrink-0" />
            <span v-if="!collapsed">退出登录</span>
          </button>
          <button
            type="button"
            class="flex h-9 items-center justify-center rounded-md text-ink-3 transition-colors hover:bg-surface-2 hover:text-ink"
            @click="toggleCollapse"
          >
            <ChevronLeft :size="16" :stroke-width="1.5" :class="collapsed ? 'rotate-180' : ''" />
          </button>
        </div>
      </aside>

      <!-- 平板：图标栏 -->
      <aside
        v-else-if="!isMobile"
        class="sticky top-0 flex h-screen w-[68px] shrink-0 flex-col items-center gap-1 border-r border-line bg-surface py-4"
      >
        <span class="mb-2 flex h-9 w-9 items-center justify-center rounded-lg bg-primary-soft text-primary">
          <Shirt :size="20" :stroke-width="1.5" />
        </span>
        <RouterLink
          v-for="nav in navs"
          :key="nav.to"
          :to="nav.to"
          class="flex h-10 w-10 items-center justify-center rounded-md transition-colors"
          :class="isActive(nav.to, nav.exact) ? 'bg-primary-soft text-primary' : 'text-ink-2 hover:bg-surface-2'"
          :title="nav.label"
        >
          <component :is="nav.icon" :size="19" :stroke-width="1.5" />
        </RouterLink>
        <div class="mt-auto flex flex-col items-center gap-1">
          <button
            type="button"
            class="flex h-10 w-10 items-center justify-center rounded-md text-ink-2 hover:bg-surface-2"
            @click="toggleTheme"
          >
            <component :is="isDark ? Sun : Moon" :size="19" :stroke-width="1.5" />
          </button>
          <button
            type="button"
            class="flex h-10 w-10 items-center justify-center rounded-md text-ink-2 hover:bg-surface-2"
            @click="signOut"
          >
            <LogOut :size="19" :stroke-width="1.5" />
          </button>
        </div>
      </aside>

      <div class="flex min-w-0 flex-1 flex-col">
        <!-- 手机：顶栏 -->
        <header
          v-if="isMobile && !isFullscreen"
          class="sticky top-0 z-30 flex h-14 items-center justify-between gap-2 border-b border-line bg-surface/95 px-3 backdrop-blur"
        >
          <button type="button" class="flex h-10 w-10 items-center justify-center text-ink-2" @click="drawerOpen = true">
            <Menu :size="20" :stroke-width="1.5" />
          </button>
          <span class="truncate text-[15px] font-medium">{{ pageTitle }}</span>
          <button type="button" class="flex h-10 w-10 items-center justify-center text-ink-2" @click="goSearch">
            <Search :size="20" :stroke-width="1.5" />
          </button>
        </header>

        <main class="flex-1" :class="isMobile && !isFullscreen ? 'pb-16' : 'md:pb-0'">
          <RouterView v-slot="{ Component }">
            <Transition name="page" mode="out-in">
              <component :is="Component" />
            </Transition>
          </RouterView>
        </main>
      </div>
    </div>

    <!-- 手机：底部 Tab -->
    <nav
      v-if="isMobile && !isFullscreen"
      class="fixed inset-x-0 bottom-0 z-30 grid grid-cols-4 border-t border-line bg-surface pb-[env(safe-area-inset-bottom)]"
    >
      <RouterLink
        v-for="nav in mobileNavs"
        :key="nav.to"
        :to="nav.to"
        class="flex h-14 flex-col items-center justify-center gap-0.5 text-[11px]"
        :class="isActive(nav.to, nav.exact) ? 'text-primary' : 'text-ink-3'"
      >
        <component :is="nav.icon" :size="20" :stroke-width="1.5" />
        {{ nav.label }}
      </RouterLink>
    </nav>

    <!-- 手机：新增 FAB -->
    <RouterLink
      v-if="showFab && isMobile"
      to="/items/new"
      class="fixed bottom-[72px] right-4 z-30 flex h-14 w-14 items-center justify-center rounded-full bg-primary text-white shadow-md"
    >
      <Plus :size="24" :stroke-width="1.5" />
    </RouterLink>

    <!-- 手机：侧滑菜单 -->
    <div
      v-if="isMobile && drawerOpen"
      class="fixed inset-0 z-40 bg-black/30 backdrop-blur-[2px]"
      @click="drawerOpen = false"
    >
      <div class="h-full w-[78%] max-w-[300px] bg-surface p-4" @click.stop>
        <div class="mb-4 flex items-center justify-between">
          <span class="text-[16px] font-semibold">{{ meta.appName }}</span>
          <button type="button" class="flex h-9 w-9 items-center justify-center text-ink-3" @click="drawerOpen = false">
            <X :size="18" :stroke-width="1.5" />
          </button>
        </div>
        <div class="flex flex-col gap-1">
          <RouterLink
            v-for="nav in navs"
            :key="nav.to"
            :to="nav.to"
            class="flex h-11 items-center gap-3 rounded-md px-3 text-sm"
            :class="isActive(nav.to, nav.exact) ? 'bg-primary-soft text-primary' : 'text-ink-2'"
          >
            <component :is="nav.icon" :size="18" :stroke-width="1.5" />
            {{ nav.label }}
          </RouterLink>
        </div>
        <div class="mt-6 border-t border-line pt-3">
          <p class="px-3 text-[12px] text-ink-3">
            {{ auth.user?.nickname || auth.user?.username || '管理员' }} · v{{ meta.meta.version }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
