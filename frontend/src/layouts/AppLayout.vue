<script setup lang="ts">
import { computed, onMounted, watch } from 'vue'
import { RouterLink, RouterView, useRouter } from 'vue-router'

import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notifications'

const authStore = useAuthStore()
const notificationStore = useNotificationStore()
const router = useRouter()

const navItems = computed(() => [
  { to: '/', label: '车次查询', icon: '⌕', show: true },
  { to: '/orders', label: '我的订单', icon: '▤', show: authStore.isAuthenticated, badge: notificationStore.pendingOrderCount },
  { to: '/tickets', label: '我的车票', icon: '▦', show: authStore.isAuthenticated, badge: notificationStore.activeTicketCount },
  { to: '/admin', label: '管理后台', icon: '⚙', show: authStore.role === 'ADMIN' },
])

onMounted(() => {
  if (authStore.isAuthenticated) {
    notificationStore.refresh().catch(() => notificationStore.reset())
  }
})

watch(
  () => authStore.isAuthenticated,
  (isAuthenticated) => {
    if (isAuthenticated) {
      notificationStore.refresh().catch(() => notificationStore.reset())
    } else {
      notificationStore.reset()
    }
  },
)

async function handleLogout() {
  await authStore.logout()
  notificationStore.reset()
  router.push('/')
}
</script>

<template>
  <div class="min-h-screen bg-[#f8fafc] text-slate-900">
    <aside class="fixed inset-y-0 left-0 z-40 hidden w-72 border-r border-slate-200 bg-white px-7 py-8 shadow-sm lg:block">
      <RouterLink to="/" class="flex items-center gap-4">
        <span class="flex h-14 w-14 items-center justify-center rounded-lg bg-teal-500 text-2xl font-bold text-white shadow-lg shadow-teal-200">
          ▣
        </span>
        <span>
          <span class="block text-2xl font-black tracking-normal text-slate-950">Mini-12306</span>
          <span class="block text-sm text-slate-400">在线车票服务系统</span>
        </span>
      </RouterLink>

      <nav class="mt-12 space-y-2">
        <RouterLink
          v-for="item in navItems.filter((item) => item.show)"
          :key="item.to"
          :to="item.to"
          class="group flex items-center gap-3 rounded-lg px-4 py-3 text-base font-bold text-slate-500 transition hover:bg-teal-50 hover:text-teal-600"
          active-class="bg-teal-50 text-teal-600 shadow-sm"
        >
          <span class="text-2xl leading-none">{{ item.icon }}</span>
          <span>{{ item.label }}</span>
          <span v-if="item.badge" class="ml-auto rounded-full bg-rose-500 px-2 py-0.5 text-xs text-white">{{ item.badge }}</span>
        </RouterLink>
      </nav>
    </aside>

    <div class="lg:pl-72">
      <header class="sticky top-0 z-30 border-b border-slate-200 bg-white/95 shadow-sm backdrop-blur">
        <div class="flex min-h-20 items-center justify-between gap-4 px-4 sm:px-8">
          <RouterLink to="/" class="flex items-center gap-3 lg:hidden">
            <span class="flex h-10 w-10 items-center justify-center rounded-lg bg-teal-500 font-bold text-white">▣</span>
            <span class="font-black text-slate-950">Mini-12306</span>
          </RouterLink>

          <nav class="hidden gap-2 overflow-x-auto lg:hidden">
            <RouterLink
              v-for="item in navItems.filter((item) => item.show)"
              :key="item.to"
              :to="item.to"
              class="whitespace-nowrap rounded-lg px-3 py-2 text-sm font-semibold text-slate-500"
              active-class="bg-teal-50 text-teal-700"
            >
              {{ item.label }}
            </RouterLink>
          </nav>

          <div class="ml-auto flex items-center gap-3">
            <template v-if="authStore.isAuthenticated">
              <div class="hidden text-right sm:block">
                <div class="text-base font-black text-slate-950">
                  {{ authStore.user?.username }}
                  <span class="rounded-md bg-slate-100 px-2 py-0.5 text-xs text-slate-500">{{ authStore.user?.role }}</span>
                </div>
                <div class="mt-1 rounded-md bg-emerald-50 px-2 py-0.5 text-xs font-bold text-emerald-600">PASSENGER 乘客账户</div>
              </div>
              <div class="relative flex h-12 w-12 items-center justify-center rounded-full border border-slate-200 bg-white text-2xl shadow-sm">
                ♙
                <span class="absolute bottom-1 right-1 h-2.5 w-2.5 rounded-full bg-emerald-500 ring-2 ring-white"></span>
              </div>
              <button class="btn-secondary h-11" type="button" @click="handleLogout">退出</button>
            </template>
            <template v-else>
              <RouterLink class="btn-secondary" to="/auth/login">登录</RouterLink>
              <RouterLink class="btn-primary" to="/auth/register">注册</RouterLink>
            </template>
          </div>
        </div>

        <nav class="flex gap-2 overflow-x-auto border-t border-slate-100 px-4 py-2 sm:px-8 lg:hidden">
          <RouterLink
            v-for="item in navItems.filter((item) => item.show)"
            :key="item.to"
            :to="item.to"
            class="whitespace-nowrap rounded-lg px-3 py-2 text-sm font-semibold text-slate-500"
            active-class="bg-teal-50 text-teal-700"
          >
            {{ item.label }}
          </RouterLink>
        </nav>
      </header>

      <RouterView />
    </div>
  </div>
</template>
