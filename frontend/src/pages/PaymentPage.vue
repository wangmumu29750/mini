<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { fetchOrder, payOrder } from '@/api/orders'
import PageHeader from '@/components/PageHeader.vue'
import { useNotificationStore } from '@/stores/notifications'
import type { ApiErrorPayload } from '@/types/api'
import type { Order } from '@/types/domain'
import { formatDateTime, formatMoney } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const notificationStore = useNotificationStore()

const order = ref<Order | null>(null)
const loading = ref(false)
const paying = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const remainingSeconds = ref(600)
let timer: number | undefined

const orderId = computed(() => Number(route.params.orderId))
const expired = computed(() => remainingSeconds.value <= 0 || order.value?.status === 'CLOSED' || order.value?.status === 'CANCELLED')
const timeLeft = computed(() => {
  const minutes = Math.floor(Math.max(remainingSeconds.value, 0) / 60)
  const seconds = Math.max(remainingSeconds.value, 0) % 60
  return `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
})

onMounted(async () => {
  await loadOrder()
  startTimer()
})

onBeforeUnmount(() => {
  if (timer) window.clearInterval(timer)
})

async function loadOrder() {
  loading.value = true
  errorMessage.value = ''
  try {
    order.value = await fetchOrder(orderId.value)
    syncRemaining()
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    loading.value = false
  }
}

function syncRemaining() {
  if (!order.value?.payExpiresAt) {
    remainingSeconds.value = 600
    return
  }
  remainingSeconds.value = Math.max(Math.floor((new Date(order.value.payExpiresAt).getTime() - Date.now()) / 1000), 0)
}

function startTimer() {
  timer = window.setInterval(() => {
    syncRemaining()
    if (remainingSeconds.value <= 0 && timer) {
      window.clearInterval(timer)
    }
  }, 1000)
}

async function handlePay() {
  if (!order.value || expired.value) return
  paying.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const result = await payOrder(order.value.id)
    order.value = result.order
    successMessage.value = `支付成功，流水号：${result.paymentNo}`
    await notificationStore.refresh()
    window.setTimeout(() => router.push({ name: 'orders' }), 600)
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    paying.value = false
  }
}
</script>

<template>
  <main>
    <PageHeader title="支付收银台" description="订单提交后请在倒计时结束前完成模拟支付。" />

    <section class="mx-auto flex min-h-[60vh] max-w-3xl items-center px-4 py-8 sm:px-8">
      <div class="w-full rounded-lg border border-slate-200 bg-white p-8 text-center shadow-sm">
        <p v-if="errorMessage" class="mb-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>
        <p v-if="successMessage" class="mb-4 rounded-lg bg-emerald-50 px-3 py-2 text-sm text-emerald-700">{{ successMessage }}</p>

        <div v-if="loading" class="text-sm font-bold text-slate-500">加载订单中...</div>
        <template v-else-if="order">
          <div class="text-sm font-black text-slate-400">订单号</div>
          <div class="mt-2 text-xl font-black text-slate-950">{{ order.orderNo }}</div>

          <div class="mt-8 text-sm font-black text-slate-400">应付金额</div>
          <div class="mt-2 text-5xl font-black text-rose-600">{{ formatMoney(order.amountCents) }}</div>

          <div class="mt-8 rounded-lg bg-slate-50 p-4">
            <div class="text-sm font-bold text-slate-500">支付截止 {{ formatDateTime(order.payExpiresAt) }}</div>
            <div class="mt-2 text-3xl font-black" :class="expired ? 'text-red-600' : 'text-slate-950'">{{ timeLeft }}</div>
            <div v-if="expired" class="mt-2 text-sm font-bold text-red-600">订单已失效或不可支付</div>
          </div>

          <button class="btn-primary mt-8 w-full py-3 text-base" type="button" :disabled="paying || expired || order.status !== 'PENDING_PAYMENT'" @click="handlePay">
            {{ paying ? '支付中...' : '模拟支付' }}
          </button>
        </template>
      </div>
    </section>
  </main>
</template>
