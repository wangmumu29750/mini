<script setup lang="ts">
import { reactive, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'

import { useAuthStore } from '@/stores/auth'
import type { ApiErrorPayload } from '@/types/api'

const authStore = useAuthStore()
const router = useRouter()

const form = reactive({
  username: '',
  password: '',
  realName: '',
  idCardNo: '',
  phone: '',
  bankCardNo: '',
})
const loading = ref(false)
const errorMessage = ref('')

async function handleSubmit() {
  errorMessage.value = ''
  loading.value = true

  try {
    await authStore.register(form)
    router.push('/')
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <form class="rounded-lg border border-slate-200 bg-white p-6 shadow-subtle" @submit.prevent="handleSubmit">
    <h2 class="text-xl font-bold text-slate-950">旅客注册</h2>
    <p class="mt-1 text-sm text-slate-500">实名信息仅用于课程系统内的模拟校验。</p>

    <div class="mt-6 grid gap-4 sm:grid-cols-2">
      <label class="block sm:col-span-2">
        <span class="form-label">用户名</span>
        <input v-model.trim="form.username" class="form-input mt-1" autocomplete="username" required />
      </label>
      <label class="block sm:col-span-2">
        <span class="form-label">密码</span>
        <input v-model="form.password" class="form-input mt-1" type="password" autocomplete="new-password" required />
      </label>
      <label class="block">
        <span class="form-label">真实姓名</span>
        <input v-model.trim="form.realName" class="form-input mt-1" required />
      </label>
      <label class="block">
        <span class="form-label">手机号</span>
        <input v-model.trim="form.phone" class="form-input mt-1" required />
      </label>
      <label class="block sm:col-span-2">
        <span class="form-label">身份证号</span>
        <input v-model.trim="form.idCardNo" class="form-input mt-1" required />
      </label>
      <label class="block sm:col-span-2">
        <span class="form-label">银行卡号</span>
        <input v-model.trim="form.bankCardNo" class="form-input mt-1" required />
      </label>
    </div>

    <p v-if="errorMessage" class="mt-4 rounded-md bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>

    <button class="btn-primary mt-6 w-full" type="submit" :disabled="loading">
      {{ loading ? '注册中...' : '注册并登录' }}
    </button>

    <p class="mt-4 text-center text-sm text-slate-500">
      已有账号？
      <RouterLink class="font-semibold text-teal-800" to="/auth/login">去登录</RouterLink>
    </p>
  </form>
</template>

