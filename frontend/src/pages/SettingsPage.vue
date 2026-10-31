<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  DatabaseBackup,
  Download,
  Eye,
  EyeOff,
  LogOut,
  Monitor,
  Moon,
  RefreshCw,
  Sun,
} from 'lucide-vue-next'
import AppPageHeader from '@/components/AppPageHeader.vue'
import AppCard from '@/components/AppCard.vue'
import AppButton from '@/components/AppButton.vue'
import AppEmpty from '@/components/AppEmpty.vue'
import AppSkeleton from '@/components/AppSkeleton.vue'
import { downloadBackup, fetchSettings, updateSettings } from '@/api/settings'
import { toErrorMessage } from '@/api/client'
import { useAuthStore } from '@/stores/auth'
import { useMetaStore } from '@/stores/meta'
import { useDictStore } from '@/stores/dict'
import { useTheme, type ThemeMode } from '@/composables/useTheme'
import type { SettingsData } from '@/types'
import { formatDateTime } from '@/utils/format'

const router = useRouter()
const auth = useAuthStore()
const meta = useMetaStore()
const dict = useDictStore()
const { mode, setMode } = useTheme()

const settings = ref<SettingsData | null>(null)
const loading = ref(true)
const error = ref('')
const saving = ref(false)
const downloading = ref(false)
const clearing = ref(false)

const themes: { value: ThemeMode; label: string; icon: typeof Sun }[] = [
  { value: 'system', label: '跟随系统', icon: Monitor },
  { value: 'light', label: '浅色', icon: Sun },
  { value: 'dark', label: '深色', icon: Moon },
]

const hidePrice = computed({
  get: () => settings.value?.hide_price ?? false,
  set: (value: boolean) => {
    void toggleHidePrice(value)
  },
})

async function load() {
  loading.value = true
  error.value = ''
  try {
    settings.value = await fetchSettings()
    meta.setHidePrice(settings.value.hide_price)
  } catch (err) {
    error.value = toErrorMessage(err)
  } finally {
    loading.value = false
  }
}

async function toggleHidePrice(value: boolean) {
  if (!settings.value) return
  saving.value = true
  try {
    const next = await updateSettings({ hide_price: value })
    settings.value = next
    meta.setHidePrice(next.hide_price)
    ElMessage.success(value ? '已隐藏价格' : '已显示价格')
  } catch (err) {
    ElMessage.error(toErrorMessage(err))
  } finally {
    saving.value = false
  }
}

async function backup() {
  downloading.value = true
  try {
    const result = await downloadBackup()
    ElMessage.success(`已下载 ${result.filename}`)
  } catch (err) {
    ElMessage.error(toErrorMessage(err))
  } finally {
    downloading.value = false
  }
}

async function clearLocal() {
  clearing.value = true
  try {
    await dict.ensure(true)
    ElMessage.success('本地缓存已刷新')
  } catch (err) {
    ElMessage.error(toErrorMessage(err))
  } finally {
    clearing.value = false
  }
}

async function signOut() {
  try {
    await ElMessageBox.confirm('退出后需要重新登录，确定继续吗？', '退出登录', {
      confirmButtonText: '退出',
      cancelButtonText: '取消',
    })
  } catch {
    return
  }
  await auth.logout()
  await router.replace('/login')
}

onMounted(async () => {
  await load()
  meta.applyTitle('设置')
})
</script>

<template>
  <div class="page">
    <AppPageHeader title="设置" subtitle="外观、价格可见性与数据备份" />

    <div v-if="error" class="card mb-4 flex items-center justify-between p-4">
      <p class="text-sm text-danger">{{ error }}</p>
      <AppButton variant="secondary" size="sm" @click="load">重试</AppButton>
    </div>

    <div class="grid gap-4 lg:grid-cols-2">
      <AppCard padding="lg">
        <h2 class="mb-3 text-[15px] font-semibold">外观</h2>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="theme in themes"
            :key="theme.value"
            type="button"
            class="flex h-10 items-center gap-2 rounded-md border px-3 text-[13px] transition-colors"
            :class="
              mode === theme.value
                ? 'border-primary bg-primary-soft text-primary'
                : 'border-line bg-surface text-ink-2 hover:bg-surface-2'
            "
            @click="setMode(theme.value)"
          >
            <component :is="theme.icon" :size="15" :stroke-width="1.5" />
            {{ theme.label }}
          </button>
        </div>
        <p class="mt-3 text-[12px] text-ink-3">默认跟随系统，选择会保存在本设备。</p>
      </AppCard>

      <AppCard padding="lg">
        <div class="flex items-start justify-between gap-3">
          <div>
            <h2 class="text-[15px] font-semibold">价格可见性</h2>
            <p class="mt-1 text-[12px] leading-[20px] text-ink-3">
              关闭后服务端不再返回真实金额，接口与导出都会脱敏。
            </p>
          </div>
          <el-switch v-model="hidePrice" :loading="saving" />
        </div>
        <div class="mt-4 flex items-center gap-2 rounded-md bg-surface-2 px-3 py-2 text-[12px] text-ink-2">
          <component :is="hidePrice ? EyeOff : Eye" :size="14" :stroke-width="1.5" />
          {{ hidePrice ? '当前：价格已隐藏' : '当前：价格可见' }}
          <span v-if="settings?.price_toggle_require_password" class="text-ink-3">（需密码确认）</span>
        </div>
      </AppCard>

      <AppCard padding="lg">
        <h2 class="mb-3 text-[15px] font-semibold">数据</h2>
        <div class="flex flex-col gap-2">
          <AppButton variant="secondary" :loading="downloading" @click="backup">
            <Download :size="15" :stroke-width="1.5" /> 下载 JSON 备份
          </AppButton>
          <AppButton variant="ghost" :loading="clearing" @click="clearLocal">
            <RefreshCw :size="15" :stroke-width="1.5" /> 刷新本地缓存
          </AppButton>
          <p class="text-[12px] leading-[20px] text-ink-3">
            备份包含衣物、分类、标签与穿着记录，价格按当前可见性导出。数据库每日备份请使用 scripts/backup.sh。
          </p>
        </div>
      </AppCard>

      <AppCard padding="lg">
        <h2 class="mb-3 text-[15px] font-semibold">关于</h2>
        <AppSkeleton v-if="loading" height="96px" />
        <AppEmpty v-else-if="!settings" compact title="暂时无法读取设置" description="请稍后重试。" />
        <dl v-else class="grid grid-cols-2 gap-y-3 text-[13px]">
          <dt class="text-ink-3">产品名</dt>
          <dd class="text-right text-ink">{{ settings.app_name }}</dd>
          <dt class="text-ink-3">版本</dt>
          <dd class="text-right text-ink">v{{ settings.version }}</dd>
          <dt class="text-ink-3">构建时间</dt>
          <dd class="text-right text-ink">{{ formatDateTime(settings.build_time) }}</dd>
          <dt class="text-ink-3">当前账号</dt>
          <dd class="text-right text-ink">{{ auth.user?.nickname || auth.user?.username || '管理员' }}</dd>
        </dl>
        <div class="mt-4 flex flex-col gap-2">
          <AppButton variant="secondary" @click="router.push('/categories')">
            <DatabaseBackup :size="15" :stroke-width="1.5" /> 分类管理
          </AppButton>
          <AppButton variant="ghost" @click="signOut">
            <LogOut :size="15" :stroke-width="1.5" /> 退出登录
          </AppButton>
        </div>
      </AppCard>
    </div>
  </div>
</template>
