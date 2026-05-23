<script setup lang="ts">
import { reactive, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'

import { useAuthStore } from '@/stores/auth'
import type { ApiErrorPayload } from '@/types/api'

const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()

const form = reactive({
  username: '',
  password: '',
})
const loading = ref(false)
const errorMessage = ref('')

async function handleSubmit() {
  errorMessage.value = ''
  loading.value = true

  try {
    await authStore.login(form)
    router.push(String(route.query.redirect || '/'))
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form class="rounded-lg border border-slate-200 bg-white p-6 shadow-subtle" @submit.prevent="handleSubmit">
    <h2 class="text-xl font-bold text-slate-950">账号登录</h2>
    <p class="mt-1 text-sm text-slate-500">使用旅客或管理员账号进入系统。</p>

    <div class="mt-6 space-y-4">
      <label class="block">
        <span class="form-label">用户名</span>
        <input v-model.trim="form.username" class="form-input mt-1" autocomplete="username" required />
      </label>

      <label class="block">
        <span class="form-label">密码</span>
        <input
          v-model="form.password"
          class="form-input mt-1"
          type="password"
          autocomplete="current-password"
          required
        />
      </label>
    </div>

    <p v-if="errorMessage" class="mt-4 rounded-md bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>

    <button class="btn-primary mt-6 w-full" type="submit" :disabled="loading">
      {{ loading ? '登录中...' : '登录' }}
    </button>

    <p class="mt-4 text-center text-sm text-slate-500">
      还没有账号？
      <RouterLink class="font-semibold text-teal-800" to="/auth/register">立即注册</RouterLink>
    </p>
  </form>
</template>

