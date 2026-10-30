import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { session } from '@/api/client'
import AppLayout from '@/layouts/AppLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/pages/LoginPage.vue'),
    meta: { public: true, title: '登录' },
  },
  {
    path: '/',
    component: AppLayout,
    children: [
      {
        path: '',
        name: 'dashboard',
        component: () => import('@/pages/DashboardPage.vue'),
        meta: { title: '仪表盘' },
      },
      {
        path: 'items',
        name: 'items',
        component: () => import('@/pages/ItemsPage.vue'),
        meta: { title: '衣物' },
      },
      {
        path: 'items/new',
        name: 'item-new',
        component: () => import('@/pages/ItemFormPage.vue'),
        meta: { title: '新增衣物', fullscreen: true },
      },
      {
        path: 'items/:id(\\d+)',
        name: 'item-detail',
        component: () => import('@/pages/ItemDetailPage.vue'),
        meta: { title: '衣物详情' },
      },
      {
        path: 'items/:id(\\d+)/edit',
        name: 'item-edit',
        component: () => import('@/pages/ItemFormPage.vue'),
        meta: { title: '编辑衣物', fullscreen: true },
      },
      {
        path: 'categories',
        name: 'categories',
        component: () => import('@/pages/CategoriesPage.vue'),
        meta: { title: '分类管理' },
      },
      {
        path: 'stats',
        name: 'stats',
        component: () => import('@/pages/StatsPage.vue'),
        meta: { title: '统计' },
      },
      {
        path: 'settings',
        name: 'settings',
        component: () => import('@/pages/SettingsPage.vue'),
        meta: { title: '设置' },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/pages/NotFoundPage.vue'),
    meta: { public: true, title: '页面不存在' },
  },
]

const router = createRouter({
  history: createWebHistory('/'),
  routes,
  scrollBehavior(_to, _from, saved) {
    return saved ?? { top: 0 }
  },
})

router.beforeEach((to) => {
  const isPublic = to.meta.public === true
  const loggedIn = !!session.token()

  if (!isPublic && !loggedIn) {
    return { path: '/login', query: to.fullPath === '/' ? {} : { redirect: to.fullPath } }
  }
  if (to.path === '/login' && loggedIn) {
    return { path: '/' }
  }
  return true
})

router.afterEach((to) => {
  const title = to.meta.title as string | undefined
  document.title = title ? `${title} · 衣橱` : '衣橱'
})

export default router
