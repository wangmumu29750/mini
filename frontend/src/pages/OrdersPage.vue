<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { cancelOrder, fetchOrders, payOrder } from '@/api/orders'
import EmptyState from '@/components/EmptyState.vue'
import PageHeader from '@/components/PageHeader.vue'
import type { ApiErrorPayload } from '@/types/api'
import type { Order } from '@/types/domain'
import { formatDateTime, formatMoney } from '@/utils/format'

const orders = ref<Order[]>([])
const loading = ref(false)
const payingId = ref<number | null>(null)
const cancellingId = ref<number | null>(null)
const errorMessage = ref('')
const successMessage = ref('')

onMounted(loadOrders)

async function loadOrders() {
  loading.value = true
  errorMessage.value = ''

  try {
    orders.value = await fetchOrders()
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
    successMessage.value = `支付成功，流水号：${result.paymentNo}`
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    payingId.value = null
  }
}

async function handleCancel(order: Order) {
  cancellingId.value = order.id
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const result = await cancelOrder(order.id)
    orders.value = orders.value.map((item) => (item.id === order.id ? result : item))
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
    PAID: '已支付/已出票',
    CLOSED: '已关闭',
  }
  return textMap[status] || status
}
</script>

<template>
  <main>
    <PageHeader title="我的订单" description="查看购票订单、模拟支付并确认出票结果。" />

    <section class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
      <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
        <div class="text-sm text-slate-500">{{ loading ? '加载中...' : `共 ${orders.length} 个订单` }}</div>
        <button class="btn-secondary" type="button" :disabled="loading" @click="loadOrders">刷新</button>
      </div>

      <p v-if="errorMessage" class="mb-4 rounded-md bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>
      <p v-if="successMessage" class="mb-4 rounded-md bg-emerald-50 px-3 py-2 text-sm text-emerald-700">{{ successMessage }}</p>

      <div v-if="orders.length" class="overflow-hidden rounded-lg border border-slate-200 bg-white shadow-subtle">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-slate-200 text-sm">
            <thead class="bg-slate-50 text-left text-xs font-semibold uppercase text-slate-500">
              <tr>
                <th class="px-4 py-3">订单</th>
                <th class="px-4 py-3">行程</th>
                <th class="px-4 py-3">乘车人</th>
                <th class="px-4 py-3">金额</th>
                <th class="px-4 py-3">状态</th>
                <th class="px-4 py-3 text-right">操作</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="order in orders" :key="order.id">
                <td class="px-4 py-4">
                  <p class="font-semibold text-slate-950">{{ order.orderNo }}</p>
                  <p class="text-xs text-slate-500">{{ order.travelDate }} / {{ order.seatClassName }}</p>
                </td>
                <td class="px-4 py-4">
                  <p class="font-medium text-slate-900">{{ order.trainNo }}</p>
                  <p class="text-xs text-slate-500">{{ order.fromStation.name }} -> {{ order.toStation.name }}</p>
                </td>
                <td class="px-4 py-4 text-slate-700">{{ order.passengerName }}</td>
                <td class="px-4 py-4 font-medium text-slate-900">{{ formatMoney(order.amountCents) }}</td>
                <td class="px-4 py-4">
                  <span class="rounded-md bg-slate-100 px-2 py-1 text-xs font-medium text-slate-700">
                    {{ statusText(order.status) }}
                  </span>
                  <p v-if="order.status === 'PENDING_PAYMENT'" class="mt-2 text-xs text-slate-500">
                    支付截止 {{ formatDateTime(order.payExpiresAt) }}
                  </p>
                  <p v-if="order.ticketNo" class="mt-2 text-xs text-emerald-700">票号 {{ order.ticketNo }}</p>
                </td>
                <td class="px-4 py-4 text-right">
                  <div v-if="order.status === 'PENDING_PAYMENT'" class="flex flex-wrap justify-end gap-2">
                    <button
                      class="btn-primary"
                      type="button"
                      :disabled="payingId === order.id || cancellingId === order.id"
                      @click="handlePay(order)"
                    >
                      {{ payingId === order.id ? '支付中...' : '模拟支付' }}
                    </button>
                    <button
                      class="btn-secondary"
                      type="button"
                      :disabled="payingId === order.id || cancellingId === order.id"
                      @click="handleCancel(order)"
                    >
                      {{ cancellingId === order.id ? '取消中...' : '取消订单' }}
                    </button>
                  </div>
                  <span v-else class="text-xs text-slate-400">无需操作</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <EmptyState
        v-else
        title="暂无订单"
        description="从车次查询页选择席别预订后，订单会显示在这里。"
      />
    </section>
  </main>
</template>
