<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, RouterView, useRouter } from 'vue-router'

import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const router = useRouter()

const navItems = computed(() => [
  { to: '/', label: '车次查询', show: true },
  { to: '/orders', label: '我的订单', show: authStore.isAuthenticated },
  { to: '/tickets', label: '我的车票', show: authStore.isAuthenticated },
  { to: '/admin', label: '管理后台', show: authStore.role === 'ADMIN' },
])

async function handleLogout() {
  await authStore.logout()
  router.push('/')
}
</script>

<template>
  <div class="min-h-screen bg-rail-surface">
    <header class="border-b border-slate-200 bg-white">
      <div class="mx-auto flex max-w-7xl flex-col gap-4 px-4 py-4 sm:px-6 lg:px-8">
        <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <RouterLink to="/" class="flex items-center gap-3">
            <span class="flex h-10 w-10 items-center justify-center rounded-md bg-teal-700 text-base font-bold text-white">
              M
            </span>
            <span>
              <span class="block text-lg font-bold text-slate-950">Mini-12306</span>
              <span class="block text-xs text-slate-500">在线车票服务系统</span>
            </span>
          </RouterLink>

          <div class="flex flex-wrap items-center gap-2">
            <template v-if="authStore.isAuthenticated">
              <span class="rounded-md bg-slate-100 px-3 py-2 text-sm text-slate-700">
                {{ authStore.user?.username }} · {{ authStore.user?.role }}
              </span>
              <button class="btn-secondary" type="button" @click="handleLogout">退出</button>
            </template>
            <template v-else>
              <RouterLink class="btn-secondary" to="/auth/login">登录</RouterLink>
              <RouterLink class="btn-primary" to="/auth/register">注册</RouterLink>
            </template>
          </div>
        </div>

        <nav class="flex gap-2 overflow-x-auto">
          <RouterLink
            v-for="item in navItems.filter((item) => item.show)"
            :key="item.to"
            :to="item.to"
            class="whitespace-nowrap rounded-md px-3 py-2 text-sm font-medium text-slate-600 hover:bg-slate-100 hover:text-slate-950"
            active-class="bg-teal-50 text-teal-800"
          >
            {{ item.label }}
          </RouterLink>
        </nav>
      </div>
    </header>

    <RouterView />
  </div>
</template>

