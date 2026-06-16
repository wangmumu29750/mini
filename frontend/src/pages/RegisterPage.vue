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
  passengerType: 'ADULT',
})
const loading = ref(false)
const errorMessage = ref('')

function fillMockProfile() {
  const seed = String(Date.now()).slice(-5)
  form.username = `passenger${seed}`
  form.password = 'Password123'
  form.realName = '赵铁柱'
  form.phone = `138${seed.padStart(8, '0')}`
  form.idCardNo = `11010119900101${seed.slice(-4)}`
  form.bankCardNo = `6222020202020${seed.padStart(6, '0')}`
  form.passengerType = 'ADULT'
}

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
  <div class="w-full max-w-lg">
    <div class="grid rounded-2xl bg-[#edf5fc] p-1 text-center text-sm font-semibold">
      <div class="grid grid-cols-2">
        <RouterLink class="rounded-xl px-4 py-3 text-[#8ba0b4] transition hover:text-[#238fd1]" to="/auth/login">
          账号登录
        </RouterLink>
        <span class="rounded-xl bg-white px-4 py-3 text-[#238fd1] shadow-[0_4px_14px_rgb(45_137_190_/_0.14)]">旅客注册</span>
      </div>
    </div>

    <form class="mt-8" @submit.prevent="handleSubmit">
      <div class="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <h2 class="text-3xl font-bold tracking-normal text-[#173f5f]">旅客注册</h2>
          <p class="mt-3 text-sm font-medium text-[#7b91a4]">填写您的法定验证字段以模拟实名制购票审核。</p>
        </div>
        <button
          class="inline-flex shrink-0 items-center justify-center rounded-lg border border-[#cde7fb] bg-[#e7f5ff] px-4 py-2 text-sm font-semibold text-[#238fd1] shadow-subtle transition hover:bg-[#d9efff]"
          type="button"
          @click="fillMockProfile"
        >
          一键生成模拟信息
        </button>
      </div>

      <div class="mt-6 grid gap-4 sm:grid-cols-2">
        <label class="block">
          <span class="text-sm font-semibold text-[#526d82]">用户名</span>
          <input
            v-model.trim="form.username"
            class="mt-2 block w-full rounded-lg border border-[#cdddf1] bg-[#e8f1ff] px-4 py-3 text-sm text-slate-900 shadow-none transition placeholder:text-[#9db1c5] hover:border-[#9bcaf2] focus:border-[#2da9ee] focus:ring-4 focus:ring-[#2da9ee]/15"
            autocomplete="username"
            placeholder="英文字母"
            required
          />
        </label>
        <label class="block">
          <span class="text-sm font-semibold text-[#526d82]">登录密码</span>
          <input
            v-model="form.password"
            class="mt-2 block w-full rounded-lg border border-[#cdddf1] bg-[#e8f1ff] px-4 py-3 text-sm text-slate-900 shadow-none transition placeholder:text-[#9db1c5] hover:border-[#9bcaf2] focus:border-[#2da9ee] focus:ring-4 focus:ring-[#2da9ee]/15"
            type="password"
            autocomplete="new-password"
            placeholder="防伪密码"
            minlength="8"
            required
          />
        </label>
        <label class="block">
          <span class="text-sm font-semibold text-[#526d82]">真实姓名</span>
          <input
            v-model.trim="form.realName"
            class="mt-2 block w-full rounded-lg border border-[#cdddf1] bg-[#e8f1ff] px-4 py-3 text-sm text-slate-900 shadow-none transition placeholder:text-[#9db1c5] hover:border-[#9bcaf2] focus:border-[#2da9ee] focus:ring-4 focus:ring-[#2da9ee]/15"
            autocomplete="name"
            placeholder="例如：赵铁柱"
            required
          />
        </label>
        <label class="block">
          <span class="text-sm font-semibold text-[#526d82]">手机号码</span>
          <input
            v-model.trim="form.phone"
            class="mt-2 block w-full rounded-lg border border-[#cdddf1] bg-[#e8f1ff] px-4 py-3 text-sm text-slate-900 shadow-none transition placeholder:text-[#9db1c5] hover:border-[#9bcaf2] focus:border-[#2da9ee] focus:ring-4 focus:ring-[#2da9ee]/15"
            autocomplete="tel"
            inputmode="numeric"
            pattern="1[3-9][0-9]{9}"
            placeholder="11位合法号"
            required
          />
        </label>
        <label class="block">
          <span class="text-sm font-semibold text-[#526d82]">乘车人类型</span>
          <select
            v-model="form.passengerType"
            class="mt-2 block w-full rounded-lg border border-[#cdddf1] bg-[#e8f1ff] px-4 py-3 text-sm text-slate-900 shadow-none transition hover:border-[#9bcaf2] focus:border-[#2da9ee] focus:ring-4 focus:ring-[#2da9ee]/15"
          >
            <option value="ADULT">成人</option>
            <option value="STUDENT">学生</option>
            <option value="CHILD">儿童</option>
          </select>
        </label>
        <label class="block sm:col-span-2">
          <span class="text-sm font-semibold text-[#526d82]">身份证号码（18位）</span>
          <input
            v-model.trim="form.idCardNo"
            class="mt-2 block w-full rounded-lg border border-[#cdddf1] bg-[#e8f1ff] px-4 py-3 text-sm text-slate-900 shadow-none transition placeholder:text-[#9db1c5] hover:border-[#9bcaf2] focus:border-[#2da9ee] focus:ring-4 focus:ring-[#2da9ee]/15"
            inputmode="text"
            pattern="[0-9]{17}[0-9Xx]"
            placeholder="请输入符合合法格式的18位身份证号码"
            required
          />
        </label>
        <label class="block sm:col-span-2">
          <span class="text-sm font-semibold text-[#526d82]">关联储蓄卡 / 银行卡号</span>
          <input
            v-model.trim="form.bankCardNo"
            class="mt-2 block w-full rounded-lg border border-[#cdddf1] bg-[#e8f1ff] px-4 py-3 text-sm text-slate-900 shadow-none transition placeholder:text-[#9db1c5] hover:border-[#9bcaf2] focus:border-[#2da9ee] focus:ring-4 focus:ring-[#2da9ee]/15"
            inputmode="numeric"
            pattern="[0-9]{12,24}"
            placeholder="用于充值及款项自动退还：62220..."
            required
          />
        </label>
      </div>

      <p class="mt-6 rounded-lg border border-[#d7e8f4] bg-[#f2f8fc] px-4 py-3 text-sm leading-6 text-[#7b91a4]">
        本平台为仿真测试环境，数据仅缓存在本地localStorage。请勿在此填写任何外部真实强密码或其他隐私凭证。
      </p>

      <p v-if="errorMessage" class="mt-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>

      <button
        class="mt-6 inline-flex w-full items-center justify-center rounded-lg bg-[#2da9ee] px-4 py-4 text-sm font-bold text-white shadow-[0_16px_32px_rgb(45_169_238_/_0.28)] transition hover:bg-[#1b9ee8] disabled:cursor-not-allowed disabled:bg-[#bfd4e8]"
        type="submit"
        :disabled="loading"
      >
        {{ loading ? '注册中...' : '校验实名信息并注册' }}
      </button>
    </form>
  </div>
</template>
