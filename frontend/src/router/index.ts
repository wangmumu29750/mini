import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

import { useAuthStore } from '@/stores/auth'
import type { UserRole } from '@/types/auth'

declare module 'vue-router' {
  interface RouteMeta {
    requiresAuth?: boolean
    roles?: UserRole[]
  }
}

const routes: RouteRecordRaw[] = [
  {
    path: '/auth',
    component: () => import('@/layouts/AuthLayout.vue'),
    children: [
      {
        path: 'login',
        name: 'login',
        component: () => import('@/pages/LoginPage.vue'),
      },
      {
        path: 'register',
        name: 'register',
        component: () => import('@/pages/RegisterPage.vue'),
      },
    ],
  },
  {
    path: '/',
    component: () => import('@/layouts/AppLayout.vue'),
    children: [
      {
        path: '',
        name: 'train-search',
        component: () => import('@/pages/TrainSearchPage.vue'),
      },
      {
        path: 'orders',
        name: 'orders',
        component: () => import('@/pages/OrdersPage.vue'),
        meta: { requiresAuth: true, roles: ['PASSENGER'] },
      },
      {
        path: 'tickets',
        name: 'tickets',
        component: () => import('@/pages/TicketsPage.vue'),
        meta: { requiresAuth: true, roles: ['PASSENGER'] },
      },
      {
        path: 'clerk',
        name: 'clerk',
        component: () => import('@/pages/ClerkWorkspacePage.vue'),
        meta: { requiresAuth: true, roles: ['CLERK', 'ADMIN'] },
      },
      {
        path: 'admin',
        name: 'admin',
        component: () => import('@/pages/AdminDashboardPage.vue'),
        meta: { requiresAuth: true, roles: ['ADMIN'] },
      },
      {
        path: 'admin/settings',
        name: 'admin-settings',
        component: () => import('@/pages/AdminSettingsPage.vue'),
        meta: { requiresAuth: true, roles: ['ADMIN'] },
      },
      {
        path: 'forbidden',
        name: 'forbidden',
        component: () => import('@/pages/ForbiddenPage.vue'),
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

function homeRouteForRole(role: UserRole | null) {
  if (role === 'ADMIN') return { name: 'admin' }
  if (role === 'CLERK') return { name: 'clerk' }
  return { name: 'train-search' }
}

router.beforeEach((to) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return {
      name: 'login',
      query: { redirect: to.fullPath },
    }
  }

  if (to.meta.roles?.length && authStore.role && !to.meta.roles.includes(authStore.role)) {
    return { name: 'forbidden' }
  }

  if (to.name === 'train-search' && (authStore.role === 'ADMIN' || authStore.role === 'CLERK')) {
    return homeRouteForRole(authStore.role)
  }

  if (to.name === 'login' && authStore.isAuthenticated) {
    return homeRouteForRole(authStore.role)
  }

  return true
})

export default router
