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
    const user = await authStore.login(form)
    const redirect = String(route.query.redirect || '')
    if (redirect && redirect !== '/') {
      router.push(redirect)
      return
    }
    router.push(user.role === 'ADMIN' ? '/admin' : user.role === 'CLERK' ? '/clerk' : '/')
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form class="w-full max-w-md" @submit.prevent="handleSubmit">
    <h2 class="text-3xl font-bold tracking-normal text-[#173f5f]">账号登录</h2>
    <p class="mt-4 text-sm font-medium text-[#7b91a4]">使用旅客或管理员账号进入系统。</p>

    <div class="mt-10 space-y-6">
      <label class="block">
        <span class="text-sm font-semibold text-[#526d82]">用户名</span>
        <input
          v-model.trim="form.username"
          class="mt-3 block w-full rounded-lg border border-[#cdddf1] bg-[#e8f1ff] px-4 py-4 text-sm text-slate-900 shadow-none transition placeholder:text-[#9db1c5] hover:border-[#9bcaf2] focus:border-[#2da9ee] focus:ring-4 focus:ring-[#2da9ee]/15"
          autocomplete="username"
          placeholder="请输入用户名"
          required
        />
      </label>

      <label class="block">
        <span class="text-sm font-semibold text-[#526d82]">密码</span>
        <input
          v-model="form.password"
          class="mt-3 block w-full rounded-lg border border-[#cdddf1] bg-[#e8f1ff] px-4 py-4 text-sm text-slate-900 shadow-none transition placeholder:text-[#9db1c5] hover:border-[#9bcaf2] focus:border-[#2da9ee] focus:ring-4 focus:ring-[#2da9ee]/15"
          type="password"
          autocomplete="current-password"
          placeholder="请输入密码"
          required
        />
      </label>
    </div>

    <p v-if="errorMessage" class="mt-5 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>

    <button
      class="mt-8 inline-flex w-full items-center justify-center rounded-lg bg-[#2da9ee] px-4 py-4 text-sm font-bold text-white shadow-[0_16px_32px_rgb(45_169_238_/_0.28)] transition hover:bg-[#1b9ee8] disabled:cursor-not-allowed disabled:bg-[#bfd4e8]"
      type="submit"
      :disabled="loading"
    >
      {{ loading ? '登录中...' : '登录' }}
    </button>

    <p class="mt-6 text-center text-sm text-[#7b91a4]">
      还没有账号？
      <RouterLink class="font-semibold text-[#1a9bed]" to="/auth/register">立即注册</RouterLink>
    </p>
  </form>
</template>
