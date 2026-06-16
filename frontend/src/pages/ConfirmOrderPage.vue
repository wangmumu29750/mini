<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { createOrder, fetchPassengers } from '@/api/orders'
import { searchTrains } from '@/api/trains'
import EmptyState from '@/components/EmptyState.vue'
import PageHeader from '@/components/PageHeader.vue'
import type { ApiErrorPayload } from '@/types/api'
import type { PassengerSummary, SeatOption, TicketType, TrainSearchItem } from '@/types/domain'
import { formatDuration, formatMoney, formatTime } from '@/utils/format'

type PassengerDraft = {
  passengerId: number
  selected: boolean
  seatType: string
  ticketType: TicketType
}

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const submitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')
const train = ref<TrainSearchItem | null>(null)
const passengers = ref<PassengerSummary[]>([])
const drafts = ref<PassengerDraft[]>([])

const query = computed(() => ({
  trainId: Number(route.query.trainId),
  travelDate: String(route.query.travelDate || ''),
  fromStationId: Number(route.query.fromStationId),
  toStationId: Number(route.query.toStationId),
  seatType: String(route.query.seatType || 'SECOND'),
}))

const selectedDrafts = computed(() => drafts.value.filter((item) => item.selected))
const totalPrice = computed(() => selectedDrafts.value.reduce((sum, item) => sum + displayPrice(item), 0))

onMounted(loadPage)

async function loadPage() {
  loading.value = true
  errorMessage.value = ''
  try {
    const [trainOptions, passengerItems] = await Promise.all([
      searchTrains({
        date: query.value.travelDate,
        fromStationId: query.value.fromStationId,
        toStationId: query.value.toStationId,
      }),
      fetchPassengers(),
    ])
    train.value = trainOptions.find((item) => item.trainId === query.value.trainId) || null
    passengers.value = passengerItems
    drafts.value = passengerItems.map((passenger, index) => ({
      passengerId: passenger.id,
      selected: index === 0,
      seatType: query.value.seatType,
      ticketType: defaultTicketType(passenger.passengerType),
    }))
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    loading.value = false
  }
}

async function submitOrder() {
  if (!train.value || selectedDrafts.value.length === 0) {
    errorMessage.value = '请至少选择一位乘车人'
    return
  }

  submitting.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const order = await createOrder({
      trainId: query.value.trainId,
      travelDate: query.value.travelDate,
      fromStationId: query.value.fromStationId,
      toStationId: query.value.toStationId,
      passengers: selectedDrafts.value.map((item) => ({
        passengerId: item.passengerId,
        seatType: item.seatType,
        ticketType: item.ticketType,
      })),
      idempotencyKey: `${Date.now()}-${Math.random().toString(16).slice(2)}`,
    })
    successMessage.value = '订单已创建，正在前往支付收银台。'
    await router.push({ name: 'payment', params: { orderId: order.id } })
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    submitting.value = false
  }
}

function seatOptions() {
  return train.value?.seatOptions || []
}

function seatName(code: string) {
  return seatOptions().find((item) => item.seatClassCode === code)?.seatClassName || code
}

function seatBasePrice(code: string) {
  return seatOptions().find((item) => item.seatClassCode === code)?.priceCents || 0
}

function displayPrice(draft: PassengerDraft) {
  const base = seatBasePrice(draft.seatType)
  if (draft.ticketType === 'STUDENT') {
    return Math.round(base * studentDiscountFactor())
  }
  if (draft.ticketType === 'CHILD') return Math.round(base * 0.5)
  return base
}

function studentDiscountFactor() {
  const type = train.value?.trainType || train.value?.trainNo.slice(0, 1) || ''
  return ['Z', 'T', 'K'].includes(type) ? 0.6 : 0.75
}

function passengerName(id: number) {
  return passengers.value.find((item) => item.id === id)?.realName || '乘车人'
}

function defaultTicketType(value: string): TicketType {
  if (value === 'STUDENT') return 'STUDENT'
  if (value === 'CHILD') return 'CHILD'
  return 'ADULT'
}

function ticketTypeName(value: string) {
  if (value === 'STUDENT') return '学生票'
  if (value === 'CHILD') return '儿童票'
  return '成人票'
}

function ticketTypeOptions(draft: PassengerDraft) {
  const options = [
    { value: 'ADULT', label: '成人票' },
    { value: 'STUDENT', label: '学生票' },
  ]
  return options
}
</script>

<template>
  <main>
    <PageHeader title="确认订单" description="核对车次信息，选择乘车人与票种。前端金额仅用于核对，最终金额以后端计价为准。" />

    <section class="mx-auto max-w-7xl px-4 py-6 sm:px-8">
      <p v-if="errorMessage" class="mb-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>
      <p v-if="successMessage" class="mb-4 rounded-lg bg-emerald-50 px-3 py-2 text-sm text-emerald-700">{{ successMessage }}</p>

      <div v-if="loading" class="rounded-lg border border-slate-200 bg-white p-6 text-sm font-bold text-slate-500">加载订单信息中...</div>

      <EmptyState v-else-if="!train" title="未找到车次" description="请返回车次查询页重新选择车次。" />

      <div v-else class="grid gap-6 xl:grid-cols-[1fr_360px]">
        <div class="space-y-6">
          <section class="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
            <div class="flex flex-wrap items-center justify-between gap-4">
              <div>
                <div class="text-3xl font-black text-slate-950">{{ train.trainNo }}</div>
                <div class="mt-2 text-sm font-bold text-slate-400">{{ train.travelDate }} · {{ formatDuration(train.durationMinutes) }}</div>
              </div>
              <div class="flex items-center gap-5">
                <div>
                  <div class="text-2xl font-black text-slate-950">{{ formatTime(train.departTime) }}</div>
                  <div class="mt-1 text-sm font-bold text-slate-500">{{ train.fromStation.name }}</div>
                </div>
                <div class="text-2xl font-black text-slate-300">→</div>
                <div>
                  <div class="text-2xl font-black text-slate-950">{{ formatTime(train.arriveTime) }}</div>
                  <div class="mt-1 text-sm font-bold text-slate-500">{{ train.toStation.name }}</div>
                </div>
              </div>
            </div>
          </section>

          <section class="rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
            <div class="mb-5 flex items-center justify-between">
              <h2 class="text-lg font-black text-slate-950">乘车人</h2>
              <span class="text-sm font-bold text-slate-400">已选 {{ selectedDrafts.length }} 人</span>
            </div>

            <EmptyState v-if="!drafts.length" title="暂无乘车人" description="当前账号还没有可选实名乘车人。" />

            <div v-else class="space-y-4">
              <article v-for="draft in drafts" :key="draft.passengerId" class="grid gap-4 rounded-lg border border-slate-100 bg-slate-50 p-4 lg:grid-cols-[220px_1fr_1fr_140px] lg:items-center">
                <label class="flex items-center gap-3">
                  <input v-model="draft.selected" class="rounded border-slate-300 text-teal-600 focus:ring-teal-500" type="checkbox" />
                  <span>
                    <span class="block text-base font-black text-slate-800">{{ passengerName(draft.passengerId) }}</span>
                    <span class="text-xs font-bold text-slate-400">{{ passengers.find((item) => item.id === draft.passengerId)?.idCardNoMasked }}</span>
                  </span>
                </label>

                <select v-model="draft.seatType" class="form-input" :disabled="!draft.selected">
                  <option v-for="seat in seatOptions()" :key="seat.seatClassCode" :value="seat.seatClassCode">
                    {{ seat.seatClassName }} · {{ formatMoney(seat.priceCents) }}
                  </option>
                </select>

                <select v-model="draft.ticketType" class="form-input" :disabled="!draft.selected">
                  <option v-for="item in ticketTypeOptions(draft)" :key="item.value" :value="item.value">
                    {{ item.label }}
                  </option>
                </select>

                <div class="text-right">
                  <div class="text-xs font-bold text-slate-400">{{ seatName(draft.seatType) }} · {{ ticketTypeName(draft.ticketType) }}</div>
                  <div class="mt-1 text-xl font-black text-rose-600">{{ formatMoney(displayPrice(draft)) }}</div>
                </div>
              </article>
            </div>
          </section>
        </div>

        <aside class="h-fit rounded-lg border border-slate-200 bg-white p-6 shadow-sm">
          <h2 class="text-lg font-black text-slate-950">订单核对</h2>
          <div class="mt-5 space-y-3 text-sm font-bold text-slate-500">
            <div class="flex justify-between"><span>车次</span><span>{{ train.trainNo }}</span></div>
            <div class="flex justify-between"><span>乘车日期</span><span>{{ train.travelDate }}</span></div>
            <div class="flex justify-between"><span>人数</span><span>{{ selectedDrafts.length }}</span></div>
          </div>
          <div class="mt-6 border-t border-slate-100 pt-5">
            <div class="text-sm font-black text-slate-400">预计总价</div>
            <div class="mt-2 text-4xl font-black text-rose-600">{{ formatMoney(totalPrice) }}</div>
            <p class="mt-2 text-xs font-bold text-slate-400">最终支付金额以后端创建订单返回为准。</p>
          </div>
          <button class="btn-primary mt-6 w-full py-3" type="button" :disabled="submitting || selectedDrafts.length === 0" @click="submitOrder">
            {{ submitting ? '提交中...' : '提交订单' }}
          </button>
        </aside>
      </div>
    </section>
  </main>
</template>
