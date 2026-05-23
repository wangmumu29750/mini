<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { cancelOrder, fetchOrders, payOrder } from '@/api/orders'
import EmptyState from '@/components/EmptyState.vue'
import PageHeader from '@/components/PageHeader.vue'
import { useNotificationStore } from '@/stores/notifications'
import type { ApiErrorPayload } from '@/types/api'
import type { Order } from '@/types/domain'
import { formatDateTime, formatMoney, formatTime } from '@/utils/format'

const notificationStore = useNotificationStore()
const orders = ref<Order[]>([])
const loading = ref(false)
const payingId = ref<number | null>(null)
const cancellingId = ref<number | null>(null)
const errorMessage = ref('')
const successMessage = ref('')
const statusFilter = ref('ALL')

const filteredOrders = computed(() => {
  if (statusFilter.value === 'ALL') {
    return orders.value
  }
  return orders.value.filter((order) => order.status === statusFilter.value)
})

const pendingCount = computed(() => orders.value.filter((order) => order.status === 'PENDING_PAYMENT').length)
const paidCount = computed(() => orders.value.filter((order) => order.status === 'PAID').length)

onMounted(loadOrders)

async function loadOrders() {
  loading.value = true
  errorMessage.value = ''

  try {
    orders.value = await fetchOrders()
    notificationStore.pendingOrderCount = pendingCount.value
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    loading.value = false
  }
}

async function handlePay(order: Order) {
  payingId.value = order.id
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const result = await payOrder(order.id)
    orders.value = orders.value.map((item) => (item.id === order.id ? result.order : item))
    await notificationStore.refresh()
    successMessage.value = `支付成功，流水号：${result.paymentNo}`
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    payingId.value = null
  }
}

async function handleCancel(order: Order) {
  if (!window.confirm(`确认取消订单 ${order.orderNo}？锁定座席会释放。`)) {
    return
  }

  cancellingId.value = order.id
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const result = await cancelOrder(order.id)
    orders.value = orders.value.map((item) => (item.id === order.id ? result : item))
    await notificationStore.refresh()
    successMessage.value = '订单已取消，锁定余票已释放。'
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    cancellingId.value = null
  }
}

function statusText(status: string) {
  const textMap: Record<string, string> = {
    PENDING_PAYMENT: '待支付',
    CANCELLED: '已取消',
    PAID: '已支付',
    CLOSED: '已关闭',
  }
  return textMap[status] || status
}

function statusClass(status: string) {
  if (status === 'PENDING_PAYMENT') return 'bg-amber-50 text-amber-600 ring-amber-100'
  if (status === 'PAID') return 'bg-emerald-50 text-emerald-600 ring-emerald-100'
  return 'bg-slate-100 text-slate-500 ring-slate-200'
}
</script>

<template>
  <main>
    <PageHeader title="订单交易流水" description="管理您已提交的购票预订。未支付订单会保留座席15分钟，请尽快完成结账。" />

    <section class="mx-auto max-w-7xl px-4 py-6 sm:px-8">
      <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
        <div class="flex flex-wrap items-center gap-2">
          <button class="rounded-lg px-3 py-2 text-sm font-black" :class="statusFilter === 'ALL' ? 'bg-emerald-50 text-emerald-600' : 'bg-white text-slate-500 ring-1 ring-slate-200'" type="button" @click="statusFilter = 'ALL'">
            全部 {{ orders.length }}
          </button>
          <button class="rounded-lg px-3 py-2 text-sm font-black" :class="statusFilter === 'PENDING_PAYMENT' ? 'bg-amber-50 text-amber-600' : 'bg-white text-slate-500 ring-1 ring-slate-200'" type="button" @click="statusFilter = 'PENDING_PAYMENT'">
            待支付 {{ pendingCount }}
          </button>
          <button class="rounded-lg px-3 py-2 text-sm font-black" :class="statusFilter === 'PAID' ? 'bg-emerald-50 text-emerald-600' : 'bg-white text-slate-500 ring-1 ring-slate-200'" type="button" @click="statusFilter = 'PAID'">
            已出票 {{ paidCount }}
          </button>
        </div>
        <div class="flex items-center gap-3">
          <div class="text-sm font-bold text-slate-400">{{ loading ? '加载中...' : `当前显示 ${filteredOrders.length} 个订单` }}</div>
          <button class="btn-secondary" type="button" :disabled="loading" @click="loadOrders">刷新</button>
        </div>
      </div>

      <p v-if="errorMessage" class="mb-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>
      <p v-if="successMessage" class="mb-4 rounded-lg bg-emerald-50 px-3 py-2 text-sm text-emerald-700">{{ successMessage }}</p>

      <div v-if="filteredOrders.length" class="space-y-6">
        <article v-for="order in filteredOrders" :key="order.id" class="overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm">
          <div class="flex flex-wrap items-center justify-between gap-3 bg-slate-50 px-6 py-4">
            <div class="flex flex-wrap items-center gap-5 text-base font-bold text-slate-500">
              <span>订单号: <strong class="text-slate-800">{{ order.orderNo }}</strong></span>
              <span class="hidden h-6 w-px bg-slate-200 sm:block"></span>
              <span>下单时间: {{ order.paidAt ? formatDateTime(order.paidAt) : formatDateTime(order.payExpiresAt) }}</span>
            </div>
            <span class="rounded-lg px-3 py-1.5 text-sm font-black ring-1" :class="statusClass(order.status)">
              {{ statusText(order.status) }}
            </span>
          </div>

          <div class="grid gap-6 px-6 py-7 xl:grid-cols-[1.1fr_1fr_220px]">
            <div class="border-slate-100 xl:border-r">
              <div class="flex items-center gap-3">
                <span class="rounded-md bg-slate-100 px-3 py-1 text-sm font-black text-slate-600">{{ order.trainNo }}</span>
                <span class="text-sm font-bold text-slate-400">{{ order.travelDate }}</span>
              </div>
              <div class="mt-5 flex items-center gap-5">
                <div>
                  <div class="text-2xl font-black text-slate-950">{{ order.fromStation.name }}</div>
                  <div class="mt-2 text-sm font-bold text-slate-500">出发时刻: {{ formatTime(order.departTime) }}</div>
                </div>
                <div class="text-3xl font-black text-slate-300">→</div>
                <div>
                  <div class="text-2xl font-black text-slate-950">{{ order.toStation.name }}</div>
                  <div class="mt-2 text-sm font-bold text-slate-500">到达时刻: {{ formatTime(order.arriveTime) }}</div>
                </div>
              </div>
            </div>

            <div class="border-slate-100 xl:border-r">
              <div class="text-base font-black text-slate-400">乘车旅伴与席位</div>
              <div class="mt-4 max-w-sm rounded-lg border border-slate-100 bg-slate-50 p-4">
                <div class="text-lg font-black text-slate-800">{{ order.passengerName }}</div>
                <div class="mt-2 inline-flex rounded-md bg-emerald-50 px-2 py-1 text-xs font-black text-emerald-600">
                  {{ order.seatClassName }} · {{ order.ticketNo ? '已出票' : '待出票' }}
                </div>
                <div v-if="order.ticketNo" class="mt-2 text-xs font-bold text-slate-400">票号 {{ order.ticketNo }}</div>
              </div>
            </div>

            <div class="flex flex-col items-end justify-between gap-5">
              <div class="text-right">
                <div class="text-base font-black text-slate-400">订单金额:</div>
                <div class="mt-2 text-4xl font-black text-rose-600">{{ formatMoney(order.amountCents) }}</div>
                <div v-if="order.status === 'PENDING_PAYMENT'" class="mt-2 text-xs font-bold text-slate-400">
                  截止 {{ formatDateTime(order.payExpiresAt) }}
                </div>
              </div>

              <div v-if="order.status === 'PENDING_PAYMENT'" class="flex gap-3">
                <button class="btn-secondary" type="button" :disabled="payingId === order.id || cancellingId === order.id" @click="handleCancel(order)">
                  {{ cancellingId === order.id ? '取消中...' : '取消' }}
                </button>
                <button class="btn-primary bg-orange-500 hover:bg-orange-600" type="button" :disabled="payingId === order.id || cancellingId === order.id" @click="handlePay(order)">
                  {{ payingId === order.id ? '支付中...' : '立即支付' }}
                </button>
              </div>
              <span v-else class="rounded-lg bg-emerald-50 px-4 py-2 text-sm font-black text-emerald-600">
                {{ order.status === 'PAID' ? '交易已开票' : statusText(order.status) }}
              </span>
            </div>
          </div>
        </article>
      </div>

      <EmptyState v-else title="暂无订单" description="从车次查询页选择席别预订后，订单会显示在这里。" />
    </section>
  </main>
</template>
