<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'

import { changeTicket, fetchChangeOptions, fetchTickets, refundTicket } from '@/api/tickets'
import EmptyState from '@/components/EmptyState.vue'
import PageHeader from '@/components/PageHeader.vue'
import { useNotificationStore } from '@/stores/notifications'
import type { ApiErrorPayload } from '@/types/api'
import type { Ticket, TrainSearchItem } from '@/types/domain'
import { formatDateTime, formatMoney, formatTime } from '@/utils/format'

const notificationStore = useNotificationStore()
const tickets = ref<Ticket[]>([])
const changeOptions = ref<TrainSearchItem[]>([])
const loading = ref(false)
const optionLoading = ref(false)
const refundingId = ref<number | null>(null)
const changingId = ref<number | null>(null)
const changeTicketId = ref<number | null>(null)
const errorMessage = ref('')
const successMessage = ref('')
const statusFilter = ref('ALL')

const changeForm = reactive({
  date: '',
  trainId: '',
  seatClassCode: '',
})

const selectedTicket = computed(() => tickets.value.find((ticket) => ticket.id === changeTicketId.value) || null)
const selectedChangeTrain = computed(() => changeOptions.value.find((train) => String(train.trainId) === changeForm.trainId) || null)
const selectedSeat = computed(() => selectedChangeTrain.value?.seatOptions.find((seat) => seat.seatClassCode === changeForm.seatClassCode) || null)
const filteredTickets = computed(() => {
  if (statusFilter.value === 'ALL') {
    return tickets.value
  }
  return tickets.value.filter((ticket) => ticket.status === statusFilter.value)
})
const activeTicketCount = computed(() => tickets.value.filter((ticket) => ticket.status === 'ISSUED').length)
const changedTicketCount = computed(() => tickets.value.filter((ticket) => ticket.status === 'CHANGED_OUT').length)
const refundedTicketCount = computed(() => tickets.value.filter((ticket) => ticket.status === 'REFUNDED').length)

onMounted(async () => {
  await loadTickets()
})

watch(
  () => changeForm.trainId,
  () => {
    const firstAvailableSeat = selectedChangeTrain.value?.seatOptions.find((seat) => seat.availableCount > 0)
    changeForm.seatClassCode = firstAvailableSeat?.seatClassCode || ''
  },
)

async function loadTickets() {
  loading.value = true
  errorMessage.value = ''

  try {
    tickets.value = await fetchTickets()
    notificationStore.activeTicketCount = activeTicketCount.value
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    loading.value = false
  }
}

async function handleRefund(ticket: Ticket) {
  if (!window.confirm(`确认退票 ${ticket.trainNo} ${ticket.fromStation.name} 到 ${ticket.toStation.name}？`)) {
    return
  }

  refundingId.value = ticket.id
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const result = await refundTicket(ticket.id)
    tickets.value = tickets.value.map((item) => (item.id === ticket.id ? result.ticket : item))
    await notificationStore.refresh()
    successMessage.value = `退票成功，退款流水号：${result.refundNo}`
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    refundingId.value = null
  }
}

async function openChange(ticket: Ticket) {
  changeTicketId.value = ticket.id
  changeForm.date = ticket.travelDate
  changeForm.trainId = ''
  changeForm.seatClassCode = ticket.seatClassCode
  changeOptions.value = []
  await loadChangeOptions(ticket)
}

async function loadChangeOptions(ticket = selectedTicket.value) {
  if (!ticket) return
  optionLoading.value = true
  errorMessage.value = ''

  try {
    const result = await fetchChangeOptions(ticket.id, changeForm.date)
    changeOptions.value = result.options
    const first = changeOptions.value.find((item) => item.seatOptions.some((seat) => seat.availableCount > 0))
    if (first) {
      changeForm.trainId = String(first.trainId)
      changeForm.seatClassCode = first.seatOptions.find((seat) => seat.availableCount > 0)?.seatClassCode || ''
    } else {
      changeForm.trainId = ''
      changeForm.seatClassCode = ''
    }
  } catch (error) {
    changeOptions.value = []
    changeForm.trainId = ''
    changeForm.seatClassCode = ''
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    optionLoading.value = false
  }
}

async function handleChange() {
  const ticket = selectedTicket.value
  const train = selectedChangeTrain.value
  if (!ticket || !train || !selectedSeat.value) {
    return
  }

  changingId.value = ticket.id
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const result = await changeTicket(ticket.id, {
      newTrainId: train.trainId,
      newTravelDate: train.travelDate,
      newSeatClassCode: selectedSeat.value.seatClassCode,
      idempotencyKey: `${Date.now()}-${Math.random().toString(16).slice(2)}`,
    })
    tickets.value = tickets.value.map((item) => (item.id === ticket.id ? result.oldTicket : item))
    tickets.value.unshift(result.newTicket)
    await notificationStore.refresh()
    successMessage.value = `改签成功，流水号：${result.changeNo}，差价 ${formatMoney(result.priceDiffCents)}`
    closeChange()
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    changingId.value = null
  }
}

function closeChange() {
  changeTicketId.value = null
  changeOptions.value = []
}

function statusText(status: string) {
  const textMap: Record<string, string> = {
    ISSUED: '票价支付已结',
    REFUNDED: '已退票',
    CHANGED_OUT: '已改签',
  }
  return textMap[status] || status
}

function canOperate(ticket: Ticket) {
  return ticket.status === 'ISSUED'
}
</script>

<template>
  <main>
    <PageHeader title="全防伪电子乘车证" description="检票时请直接出示对应的电子验票二维码。中途站可凭身份证原件和行程单验证进站上车。" />

    <section class="mx-auto max-w-7xl px-4 py-6 sm:px-8">
      <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
        <div class="flex flex-wrap items-center gap-2">
          <button class="rounded-lg px-3 py-2 text-sm font-black" :class="statusFilter === 'ALL' ? 'bg-emerald-50 text-emerald-600' : 'bg-white text-slate-500 ring-1 ring-slate-200'" type="button" @click="statusFilter = 'ALL'">
            全部 {{ tickets.length }}
          </button>
          <button class="rounded-lg px-3 py-2 text-sm font-black" :class="statusFilter === 'ISSUED' ? 'bg-emerald-50 text-emerald-600' : 'bg-white text-slate-500 ring-1 ring-slate-200'" type="button" @click="statusFilter = 'ISSUED'">
            可用 {{ activeTicketCount }}
          </button>
          <button class="rounded-lg px-3 py-2 text-sm font-black" :class="statusFilter === 'CHANGED_OUT' ? 'bg-slate-100 text-slate-700' : 'bg-white text-slate-500 ring-1 ring-slate-200'" type="button" @click="statusFilter = 'CHANGED_OUT'">
            已改签 {{ changedTicketCount }}
          </button>
          <button class="rounded-lg px-3 py-2 text-sm font-black" :class="statusFilter === 'REFUNDED' ? 'bg-rose-50 text-rose-600' : 'bg-white text-slate-500 ring-1 ring-slate-200'" type="button" @click="statusFilter = 'REFUNDED'">
            已退票 {{ refundedTicketCount }}
          </button>
        </div>
        <div class="flex items-center gap-3">
          <div class="text-sm font-bold text-slate-400">{{ loading ? '加载中...' : `当前显示 ${filteredTickets.length} 张车票` }}</div>
          <button class="btn-secondary" type="button" :disabled="loading" @click="loadTickets">刷新</button>
        </div>
      </div>

      <p v-if="errorMessage" class="mb-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>
      <p v-if="successMessage" class="mb-4 rounded-lg bg-emerald-50 px-3 py-2 text-sm text-emerald-700">{{ successMessage }}</p>

      <div v-if="filteredTickets.length" class="grid gap-6 xl:grid-cols-2">
        <article v-for="ticket in filteredTickets" :key="ticket.id" class="overflow-hidden rounded-lg border border-slate-200 bg-white shadow-sm">
          <div class="bg-emerald-50/60 p-6">
            <div class="flex flex-wrap items-center justify-between gap-3">
              <div class="flex items-center gap-3">
                <span class="text-2xl text-emerald-600">▣</span>
                <span class="text-lg font-black text-slate-800">{{ ticket.trainNo }}</span>
                <span class="rounded-md bg-emerald-50 px-2 py-1 text-xs font-black text-emerald-600">中国铁路</span>
              </div>
              <span class="text-sm font-black text-slate-400">票号: {{ ticket.ticketNo }}</span>
            </div>

            <div class="mt-8 grid grid-cols-[1fr_auto_1fr] items-center gap-4">
              <div>
                <div class="text-4xl font-black text-slate-950">{{ formatTime(ticket.departTime) }}</div>
                <div class="mt-3 text-2xl font-black text-slate-800">{{ ticket.fromStation.name }}</div>
                <div class="mt-2 text-sm font-bold text-slate-400">{{ ticket.travelDate }}</div>
              </div>
              <div class="text-center">
                <div class="text-sm font-black text-slate-400">电子客票</div>
                <div class="my-2 h-px w-24 bg-slate-200"></div>
                <div class="rounded-md bg-emerald-50 px-2 py-1 text-xs font-black text-emerald-600">{{ ticket.seatClassName }}</div>
              </div>
              <div class="text-right">
                <div class="text-4xl font-black text-slate-950">{{ formatTime(ticket.arriveTime) }}</div>
                <div class="mt-3 text-2xl font-black text-slate-800">{{ ticket.toStation.name }}</div>
                <div class="mt-2 text-sm font-bold text-slate-400">{{ statusText(ticket.status) }}</div>
              </div>
            </div>
          </div>

          <div class="relative p-6">
            <div class="absolute -left-4 top-0 h-8 w-8 -translate-y-1/2 rounded-full border border-slate-200 bg-[#f8fafc]"></div>
            <div class="absolute -right-4 top-0 h-8 w-8 -translate-y-1/2 rounded-full border border-slate-200 bg-[#f8fafc]"></div>

            <div class="grid gap-5 sm:grid-cols-[1fr_auto]">
              <div class="flex items-center gap-4">
                <div class="flex h-14 w-14 items-center justify-center rounded-full border border-slate-300 text-2xl text-slate-500">♙</div>
                <div>
                  <div class="text-xl font-black text-slate-900">
                    {{ ticket.passengerName }}
                    <span class="ml-2 rounded-md bg-slate-100 px-2 py-1 text-xs text-slate-500">成人</span>
                  </div>
                  <div class="mt-2 text-sm font-bold text-slate-400">ID: {{ ticket.idCardNoMasked }}</div>
                </div>
              </div>

              <div class="text-right">
                <div class="text-sm font-black text-slate-400">座位号</div>
                <div class="mt-2 rounded-lg bg-emerald-50 px-4 py-2 text-2xl font-black text-emerald-600">
                  {{ ticket.coachNo }}车 {{ ticket.seatNo }}
                </div>
              </div>
            </div>

            <div class="mt-6 flex flex-wrap items-end justify-between gap-4 border-t border-slate-100 pt-5">
              <div>
                <div class="h-10 w-36 bg-[repeating-linear-gradient(90deg,#94a3b8_0,#94a3b8_2px,transparent_2px,transparent_5px)] opacity-60"></div>
                <div class="mt-1 text-center text-xs font-bold text-slate-300">{{ ticket.ticketNo }}</div>
              </div>
              <div class="flex flex-wrap gap-3">
                <button class="btn-secondary" type="button" :disabled="!canOperate(ticket) || refundingId === ticket.id" @click="handleRefund(ticket)">
                  {{ refundingId === ticket.id ? '退票中...' : '退票' }}
                </button>
                <button class="btn-primary bg-slate-950 hover:bg-slate-800" type="button" :disabled="!canOperate(ticket)" @click="openChange(ticket)">
                  验签二维码 / 改签
                </button>
              </div>
            </div>
          </div>
        </article>
      </div>

      <EmptyState v-else title="暂无车票" description="订单完成模拟支付后，已出票车票会显示在这里。" />
    </section>

    <div v-if="selectedTicket" class="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/40 p-4">
      <div class="w-full max-w-2xl rounded-lg bg-white p-6 shadow-xl">
        <div class="flex items-start justify-between gap-4">
          <div>
            <h2 class="text-2xl font-black text-slate-950">办理改签</h2>
            <p class="mt-1 text-sm font-medium text-slate-500">
              原票 {{ selectedTicket.trainNo }} {{ selectedTicket.fromStation.name }} → {{ selectedTicket.toStation.name }}
            </p>
          </div>
          <button class="rounded-lg px-3 py-2 text-slate-400 hover:bg-slate-100" type="button" @click="closeChange">关闭</button>
        </div>

        <div class="mt-5 grid gap-4 sm:grid-cols-[1fr_auto]">
          <label>
            <span class="form-label">新乘车日期</span>
            <input v-model="changeForm.date" class="form-input mt-2 h-11" type="date" />
          </label>
          <button class="btn-secondary mt-7 h-11" type="button" :disabled="optionLoading" @click="loadChangeOptions()">
            {{ optionLoading ? '查询中...' : '查询可改签车次' }}
          </button>
        </div>

        <div class="mt-5 grid gap-4 sm:grid-cols-2">
          <label>
            <span class="form-label">新车次</span>
            <select v-model="changeForm.trainId" class="form-input mt-2 h-11">
              <option value="">请选择</option>
              <option v-for="train in changeOptions" :key="`${train.trainId}-${train.travelDate}`" :value="String(train.trainId)">
                {{ train.trainNo }} {{ formatTime(train.departTime) }} → {{ formatTime(train.arriveTime) }}
              </option>
            </select>
          </label>

          <label>
            <span class="form-label">席别</span>
            <select v-model="changeForm.seatClassCode" class="form-input mt-2 h-11" :disabled="!selectedChangeTrain">
              <option value="">请选择</option>
              <option v-for="seat in selectedChangeTrain?.seatOptions || []" :key="seat.seatClassCode" :value="seat.seatClassCode" :disabled="seat.availableCount <= 0">
                {{ seat.seatClassName }} / {{ formatMoney(seat.priceCents) }} / 余{{ seat.availableCount }}
              </option>
            </select>
          </label>
        </div>

        <div v-if="!optionLoading && !changeOptions.length" class="mt-5 rounded-lg bg-amber-50 p-4 text-sm font-bold text-amber-700">
          褰撳墠鏃ユ湡娌℃湁鍙敼绛剧殑杞︽锛岃鏇存崲鏂颁箻杞︽棩鏈熷悗鍐嶆煡璇€?
        </div>

        <div v-if="selectedChangeTrain && selectedSeat" class="mt-5 rounded-lg bg-slate-50 p-4 text-sm font-bold text-slate-500">
          新行程：{{ selectedChangeTrain.trainNo }} {{ selectedChangeTrain.fromStation.name }} {{ formatTime(selectedChangeTrain.departTime) }}
          → {{ selectedChangeTrain.toStation.name }} {{ formatTime(selectedChangeTrain.arriveTime) }}，{{ selectedSeat.seatClassName }} {{ formatMoney(selectedSeat.priceCents) }}
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <button class="btn-secondary" type="button" @click="closeChange">取消</button>
          <button class="btn-primary" type="button" :disabled="!selectedChangeTrain || !selectedSeat || changingId === selectedTicket.id" @click="handleChange">
            {{ changingId === selectedTicket.id ? '改签中...' : '确认改签' }}
          </button>
        </div>
      </div>
    </div>
  </main>
</template>
