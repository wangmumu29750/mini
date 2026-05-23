<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import { createOrder } from '@/api/orders'
import { fetchStations, searchTrains } from '@/api/trains'
import EmptyState from '@/components/EmptyState.vue'
import PageHeader from '@/components/PageHeader.vue'
import { useAuthStore } from '@/stores/auth'
import type { ApiErrorPayload } from '@/types/api'
import type { SeatOption, Station, TrainSearchItem } from '@/types/domain'
import { formatDateTime, formatDuration, formatMoney } from '@/utils/format'

const router = useRouter()
const authStore = useAuthStore()

const stations = ref<Station[]>([])
const trains = ref<TrainSearchItem[]>([])
const loading = ref(false)
const stationLoading = ref(false)
const bookingKey = ref('')
const errorMessage = ref('')
const successMessage = ref('')

const tomorrow = new Date()
tomorrow.setDate(tomorrow.getDate() + 1)
const defaultTravelDate = tomorrow.toISOString().slice(0, 10)
const query = reactive({
  date: defaultTravelDate,
  fromStationId: '',
  toStationId: '',
})

const canSearch = computed(() => query.date && query.fromStationId && query.toStationId)

onMounted(async () => {
  stationLoading.value = true

  try {
    stations.value = await fetchStations()
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    stationLoading.value = false
  }
})

async function handleSearch() {
  if (!canSearch.value) {
    return
  }

  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    trains.value = await searchTrains(query)
  } catch (error) {
    trains.value = []
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    loading.value = false
  }
}

async function handleBook(train: TrainSearchItem, seat: SeatOption) {
  if (!authStore.isAuthenticated) {
    await router.push({ name: 'login', query: { redirect: '/' } })
    return
  }

  const key = `${train.trainId}-${train.travelDate}-${seat.seatClassCode}`
  bookingKey.value = key
  errorMessage.value = ''
  successMessage.value = ''

  try {
    await createOrder({
      trainId: train.trainId,
      travelDate: train.travelDate,
      fromStationId: train.fromStation.id,
      toStationId: train.toStation.id,
      seatClassCode: seat.seatClassCode,
      idempotencyKey: `${Date.now()}-${Math.random().toString(16).slice(2)}`,
    })
    successMessage.value = '订单已创建，请在订单页完成模拟支付。'
    await router.push({ name: 'orders' })
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    bookingKey.value = ''
  }
}
</script>

<template>
  <main>
    <PageHeader title="车次查询" description="按乘车日期、出发站和到达站查询可售车次与余票。" />

    <section class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
      <form class="rounded-lg border border-slate-200 bg-white p-4 shadow-subtle" @submit.prevent="handleSearch">
        <div class="grid gap-4 md:grid-cols-[1fr_1fr_1fr_auto] md:items-end">
          <label class="block">
            <span class="form-label">乘车日期</span>
            <input v-model="query.date" class="form-input mt-1" type="date" required />
          </label>

          <label class="block">
            <span class="form-label">出发站</span>
            <select v-model="query.fromStationId" class="form-input mt-1" :disabled="stationLoading" required>
              <option value="">请选择</option>
              <option v-for="station in stations" :key="station.id" :value="station.id">{{ station.name }}</option>
            </select>
          </label>

          <label class="block">
            <span class="form-label">到达站</span>
            <select v-model="query.toStationId" class="form-input mt-1" :disabled="stationLoading" required>
              <option value="">请选择</option>
              <option v-for="station in stations" :key="station.id" :value="station.id">{{ station.name }}</option>
            </select>
          </label>

          <button class="btn-primary h-10" type="submit" :disabled="loading || !canSearch">
            {{ loading ? '查询中...' : '查询车次' }}
          </button>
        </div>

        <p v-if="errorMessage" class="mt-4 rounded-md bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>
        <p v-if="successMessage" class="mt-4 rounded-md bg-emerald-50 px-3 py-2 text-sm text-emerald-700">{{ successMessage }}</p>
      </form>

      <div class="mt-6 overflow-hidden rounded-lg border border-slate-200 bg-white shadow-subtle">
        <div class="flex items-center justify-between border-b border-slate-200 px-4 py-3">
          <h2 class="text-base font-semibold text-slate-950">查询结果</h2>
          <span class="text-sm text-slate-500">{{ trains.length }} 趟车</span>
        </div>

        <div v-if="trains.length" class="overflow-x-auto">
          <table class="min-w-full divide-y divide-slate-200 text-sm">
            <thead class="bg-slate-50 text-left text-xs font-semibold uppercase text-slate-500">
              <tr>
                <th class="px-4 py-3">车次</th>
                <th class="px-4 py-3">出发</th>
                <th class="px-4 py-3">到达</th>
                <th class="px-4 py-3">历时</th>
                <th class="px-4 py-3">席别与余票</th>
                <th class="px-4 py-3 text-right">操作</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100">
              <tr v-for="train in trains" :key="`${train.trainId}-${train.travelDate}`">
                <td class="px-4 py-4">
                  <p class="font-semibold text-slate-950">{{ train.trainNo }}</p>
                  <p class="text-xs text-slate-500">{{ train.travelDate }}</p>
                </td>
                <td class="px-4 py-4">
                  <p class="font-medium text-slate-900">{{ train.fromStation.name }}</p>
                  <p class="text-xs text-slate-500">{{ formatDateTime(train.departTime) }}</p>
                </td>
                <td class="px-4 py-4">
                  <p class="font-medium text-slate-900">{{ train.toStation.name }}</p>
                  <p class="text-xs text-slate-500">{{ formatDateTime(train.arriveTime) }}</p>
                </td>
                <td class="px-4 py-4 text-slate-700">{{ formatDuration(train.durationMinutes) }}</td>
                <td class="px-4 py-4">
                  <div class="flex flex-wrap gap-2">
                    <span
                      v-for="seat in train.seatOptions"
                      :key="seat.seatClassCode"
                      class="rounded-md border border-slate-200 px-2 py-1 text-xs text-slate-700"
                    >
                      {{ seat.seatClassName }} / {{ formatMoney(seat.priceCents) }} / 余{{ seat.availableCount }}
                    </span>
                  </div>
                </td>
                <td class="px-4 py-4 text-right">
                  <div class="flex flex-wrap justify-end gap-2">
                    <button
                      v-for="seat in train.seatOptions"
                      :key="seat.seatClassCode"
                      class="btn-secondary"
                      type="button"
                      :disabled="seat.availableCount <= 0 || bookingKey === `${train.trainId}-${train.travelDate}-${seat.seatClassCode}`"
                      @click="handleBook(train, seat)"
                    >
                      {{ bookingKey === `${train.trainId}-${train.travelDate}-${seat.seatClassCode}` ? '提交中...' : `预订${seat.seatClassName}` }}
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <EmptyState
          v-else
          class="m-4"
          title="暂无查询结果"
          description="选择乘车日期、出发站和到达站后查询。"
        />
      </div>
    </section>
  </main>
</template>
