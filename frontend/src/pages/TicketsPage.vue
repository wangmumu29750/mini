<script setup lang="ts">
import { onMounted, ref } from 'vue'

import { fetchTickets } from '@/api/tickets'
import EmptyState from '@/components/EmptyState.vue'
import PageHeader from '@/components/PageHeader.vue'
import type { ApiErrorPayload } from '@/types/api'
import type { Ticket } from '@/types/domain'
import { formatDateTime } from '@/utils/format'

const tickets = ref<Ticket[]>([])
const loading = ref(false)
const errorMessage = ref('')

onMounted(loadTickets)

async function loadTickets() {
  loading.value = true
  errorMessage.value = ''

  try {
    tickets.value = await fetchTickets()
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    loading.value = false
  }
}

function statusText(status: string) {
  const textMap: Record<string, string> = {
    ISSUED: '已出票',
  }
  return textMap[status] || status
}
</script>

<template>
  <main>
    <PageHeader title="我的车票" description="查看已出票车票，后续退票和改签会从这里进入。" />

    <section class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
      <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
        <div class="text-sm text-slate-500">{{ loading ? '加载中...' : `共 ${tickets.length} 张车票` }}</div>
        <button class="btn-secondary" type="button" :disabled="loading" @click="loadTickets">刷新</button>
      </div>

      <p v-if="errorMessage" class="mb-4 rounded-md bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>

      <div v-if="tickets.length" class="grid gap-4 lg:grid-cols-2">
        <article
          v-for="ticket in tickets"
          :key="ticket.id"
          class="rounded-lg border border-slate-200 bg-white p-4 shadow-subtle"
        >
          <div class="flex flex-wrap items-start justify-between gap-3 border-b border-slate-100 pb-3">
            <div>
              <p class="text-xs text-slate-500">票号 {{ ticket.ticketNo }}</p>
              <h2 class="mt-1 text-xl font-semibold text-slate-950">{{ ticket.trainNo }}</h2>
            </div>
            <span class="rounded-md bg-emerald-50 px-2 py-1 text-xs font-medium text-emerald-700">
              {{ statusText(ticket.status) }}
            </span>
          </div>

          <div class="mt-4 grid grid-cols-[1fr_auto_1fr] items-center gap-4">
            <div>
              <p class="text-lg font-semibold text-slate-950">{{ ticket.fromStation.name }}</p>
              <p class="mt-1 text-xs text-slate-500">{{ ticket.travelDate }}</p>
            </div>
            <div class="h-px min-w-12 bg-slate-300"></div>
            <div class="text-right">
              <p class="text-lg font-semibold text-slate-950">{{ ticket.toStation.name }}</p>
              <p class="mt-1 text-xs text-slate-500">{{ ticket.seatClassName }}</p>
            </div>
          </div>

          <dl class="mt-4 grid gap-3 rounded-md bg-slate-50 p-3 text-sm sm:grid-cols-2">
            <div>
              <dt class="text-xs text-slate-500">乘车人</dt>
              <dd class="mt-1 font-medium text-slate-900">{{ ticket.passengerName }}</dd>
            </div>
            <div>
              <dt class="text-xs text-slate-500">证件号</dt>
              <dd class="mt-1 font-medium text-slate-900">{{ ticket.idCardNoMasked }}</dd>
            </div>
            <div>
              <dt class="text-xs text-slate-500">订单 ID</dt>
              <dd class="mt-1 font-medium text-slate-900">{{ ticket.orderId }}</dd>
            </div>
            <div>
              <dt class="text-xs text-slate-500">出票时间</dt>
              <dd class="mt-1 font-medium text-slate-900">{{ formatDateTime(ticket.issuedAt) }}</dd>
            </div>
          </dl>
        </article>
      </div>

      <EmptyState
        v-else
        title="暂无车票"
        description="订单完成模拟支付后，已出票车票会显示在这里。"
      />
    </section>
  </main>
</template>
