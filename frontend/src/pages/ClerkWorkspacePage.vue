<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'

import { createClerkOrder } from '@/api/clerk'
import { fetchStations, searchTrains } from '@/api/trains'
import EmptyState from '@/components/EmptyState.vue'
import PageHeader from '@/components/PageHeader.vue'
import type { ApiErrorPayload } from '@/types/api'
import type { SeatOption, Station, TrainSearchItem } from '@/types/domain'
import { formatDuration, formatMoney, formatTime } from '@/utils/format'

const stations = ref<Station[]>([])
const trains = ref<TrainSearchItem[]>([])
const selectedTrain = ref<TrainSearchItem | null>(null)
const selectedSeat = ref<SeatOption | null>(null)
const loading = ref(false)
const submitting = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const tomorrow = new Date()
tomorrow.setDate(tomorrow.getDate() + 1)

const query = reactive({
  date: tomorrow.toISOString().slice(0, 10),
  fromStationId: '',
  toStationId: '',
})

const passenger = reactive({
  passengerName: '',
  idCardNo: '',
  phone: '',
  bankCardNo: '',
})

const canSearch = computed(() => query.date && query.fromStationId && query.toStationId)
const canSubmit = computed(() => selectedTrain.value && selectedSeat.value && passenger.passengerName && passenger.idCardNo && passenger.phone && passenger.bankCardNo)

onMounted(async () => {
  try {
    stations.value = await fetchStations()
    if (stations.value.length >= 2) {
      const first = stations.value[0]
      const last = stations.value[stations.value.length - 1]
      if (!first || !last) return
      query.fromStationId = String(first.id)
      query.toStationId = String(last.id)
      await handleSearch()
    }
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  }
})

async function handleSearch() {
  if (!canSearch.value) return
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''
  selectedTrain.value = null
  selectedSeat.value = null

  try {
    trains.value = await searchTrains(query)
  } catch (error) {
    trains.value = []
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    loading.value = false
  }
}

function chooseSeat(train: TrainSearchItem, seat: SeatOption) {
  selectedTrain.value = train
  selectedSeat.value = seat
  successMessage.value = ''
}

function decrementSeatAvailability(train: TrainSearchItem, seat: SeatOption) {
  trains.value = trains.value.map((item) => {
    if (item.trainId !== train.trainId || item.travelDate !== train.travelDate) {
      return item
    }
    return {
      ...item,
      seatOptions: item.seatOptions.map((option) => {
        if (option.seatClassCode !== seat.seatClassCode) {
          return option
        }
        return {
          ...option,
          availableCount: Math.max(option.availableCount - 1, 0),
        }
      }),
    }
  })
}

async function handleSubmit() {
  if (!selectedTrain.value || !selectedSeat.value) return
  submitting.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const order = await createClerkOrder({
      trainId: selectedTrain.value.trainId,
      travelDate: selectedTrain.value.travelDate,
      fromStationId: selectedTrain.value.fromStation.id,
      toStationId: selectedTrain.value.toStation.id,
      seatClassCode: selectedSeat.value.seatClassCode,
      idempotencyKey: `clerk-${Date.now()}-${Math.random().toString(16).slice(2)}`,
      ...passenger,
    })
    successMessage.value = `已为 ${order.passengerName} 创建订单 ${order.orderNo}，待模拟支付。`
    decrementSeatAvailability(selectedTrain.value, selectedSeat.value)
    selectedSeat.value = null
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    submitting.value = false
  }
}

function minPrice(train: TrainSearchItem) {
  return Math.min(...train.seatOptions.map((seat) => seat.priceCents))
}
</script>

<template>
  <main>
    <PageHeader title="售票员工作台" description="面向窗口售票场景，售票员可替旅客查询车次并录入实名信息创建订单。" />

    <section class="mx-auto grid max-w-7xl gap-6 px-4 py-6 sm:px-8 xl:grid-cols-[1fr_380px]">
      <div class="space-y-6">
        <form class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm" @submit.prevent="handleSearch">
          <div class="grid gap-4 md:grid-cols-[1fr_1fr_1fr_auto] md:items-end">
            <label>
              <span class="form-label">出发站</span>
              <select v-model="query.fromStationId" class="form-input mt-2 h-12" required>
                <option value="">请选择</option>
                <option v-for="station in stations" :key="station.id" :value="station.id">{{ station.name }}</option>
              </select>
            </label>
            <label>
              <span class="form-label">到达站</span>
              <select v-model="query.toStationId" class="form-input mt-2 h-12" required>
                <option value="">请选择</option>
                <option v-for="station in stations" :key="station.id" :value="station.id">{{ station.name }}</option>
              </select>
            </label>
            <label>
              <span class="form-label">乘车日期</span>
              <input v-model="query.date" class="form-input mt-2 h-12" type="date" required />
            </label>
            <button class="btn-primary h-12 px-8" type="submit" :disabled="loading || !canSearch">{{ loading ? '查询中...' : '查询' }}</button>
          </div>
        </form>

        <p v-if="errorMessage" class="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>
        <p v-if="successMessage" class="rounded-lg bg-emerald-50 px-3 py-2 text-sm text-emerald-700">{{ successMessage }}</p>

        <div v-if="trains.length" class="space-y-4">
          <article v-for="train in trains" :key="`${train.trainId}-${train.travelDate}`" class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
            <div class="flex flex-wrap items-center justify-between gap-3 border-b border-slate-100 pb-4">
              <div>
                <div class="text-2xl font-black text-slate-950">{{ train.trainNo }}</div>
                <div class="mt-1 text-sm font-bold text-slate-500">
                  {{ train.fromStation.name }} {{ formatTime(train.departTime) }} 至 {{ train.toStation.name }} {{ formatTime(train.arriveTime) }}
                </div>
              </div>
              <div class="text-right">
                <div class="text-sm font-bold text-slate-400">{{ formatDuration(train.durationMinutes) }}</div>
                <div class="text-lg font-black text-rose-600">起 {{ formatMoney(minPrice(train)) }}</div>
              </div>
            </div>

            <div class="mt-4 grid gap-3 sm:grid-cols-3">
              <button
                v-for="seat in train.seatOptions"
                :key="seat.seatClassCode"
                class="rounded-lg border p-4 text-left transition disabled:cursor-not-allowed disabled:opacity-50"
                :class="selectedTrain?.trainId === train.trainId && selectedSeat?.seatClassCode === seat.seatClassCode ? 'border-teal-500 bg-teal-50' : 'border-slate-200 bg-slate-50 hover:border-teal-200'"
                type="button"
                :disabled="seat.availableCount <= 0"
                @click="chooseSeat(train, seat)"
              >
                <div class="font-black text-slate-800">{{ seat.seatClassName }}</div>
                <div class="mt-2 text-xl font-black text-slate-950">{{ formatMoney(seat.priceCents) }}</div>
                <div class="mt-1 text-sm font-bold text-emerald-600">余 {{ seat.availableCount }}</div>
              </button>
            </div>
          </article>
        </div>
        <EmptyState v-else title="暂无可售车次" description="请调整站点或日期后重新查询。" />
      </div>

      <aside class="h-fit rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
        <h2 class="text-lg font-black text-slate-950">旅客实名信息</h2>
        <div class="mt-4 space-y-4">
          <label class="block">
            <span class="form-label">姓名</span>
            <input v-model.trim="passenger.passengerName" class="form-input mt-2 h-11" placeholder="张三" />
          </label>
          <label class="block">
            <span class="form-label">身份证号</span>
            <input v-model.trim="passenger.idCardNo" class="form-input mt-2 h-11" placeholder="110101199001011234" />
          </label>
          <label class="block">
            <span class="form-label">手机号</span>
            <input v-model.trim="passenger.phone" class="form-input mt-2 h-11" placeholder="13800138000" />
          </label>
          <label class="block">
            <span class="form-label">银行卡号</span>
            <input v-model.trim="passenger.bankCardNo" class="form-input mt-2 h-11" placeholder="6222020202020202020" />
          </label>
        </div>

        <div class="mt-5 rounded-lg bg-slate-50 p-4 text-sm font-bold text-slate-600">
          <template v-if="selectedTrain && selectedSeat">
            已选 {{ selectedTrain.trainNo }} / {{ selectedSeat.seatClassName }} / {{ formatMoney(selectedSeat.priceCents) }}
          </template>
          <template v-else>请选择车次和席别</template>
        </div>

        <button class="btn-primary mt-5 h-12 w-full" type="button" :disabled="submitting || !canSubmit" @click="handleSubmit">
          {{ submitting ? '提交中...' : '创建售票订单' }}
        </button>
      </aside>
    </section>
  </main>
</template>
