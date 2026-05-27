<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import { createOrder } from '@/api/orders'
import { fetchStations, searchTrains } from '@/api/trains'
import EmptyState from '@/components/EmptyState.vue'
import PageHeader from '@/components/PageHeader.vue'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notifications'
import type { ApiErrorPayload } from '@/types/api'
import type { SeatOption, Station, TrainSearchItem } from '@/types/domain'
import { formatDuration, formatMoney, formatTime } from '@/utils/format'

const router = useRouter()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const stations = ref<Station[]>([])
const trains = ref<TrainSearchItem[]>([])
const loading = ref(false)
const stationLoading = ref(false)
const bookingKey = ref('')
const errorMessage = ref('')
const successMessage = ref('')
const onlyHighSpeed = ref(false)
const sortMode = ref<'depart' | 'duration' | 'price'>('depart')
const sortDirection = ref<'asc' | 'desc'>('asc')

const tomorrow = new Date()
tomorrow.setDate(tomorrow.getDate() + 1)
const defaultTravelDate = tomorrow.toISOString().slice(0, 10)
const query = reactive({
  date: defaultTravelDate,
  fromStationId: '',
  toStationId: '',
})

const canSearch = computed(() => query.date && query.fromStationId && query.toStationId)
const quickDates = computed(() => {
  const base = new Date()
  return [0, 1, 2, 5].map((offset) => {
    const date = new Date(base)
    date.setDate(base.getDate() + offset)
    return {
      label: offset === 0 ? '今天' : offset === 1 ? '明天' : offset === 2 ? '后天' : '五天后',
      value: date.toISOString().slice(0, 10),
      text: `${date.getMonth() + 1}月${date.getDate()}日`,
    }
  })
})

const searchSummary = computed(() => {
  const from = stations.value.find((station) => String(station.id) === query.fromStationId)?.name || '出发站'
  const to = stations.value.find((station) => String(station.id) === query.toStationId)?.name || '到达站'
  return `${query.date} ${from} → ${to}`
})

const visibleTrains = computed(() => {
  let items = [...trains.value]
  if (onlyHighSpeed.value) {
    items = items.filter((item) => /^[GD]/.test(item.trainNo))
  }
  const direction = sortDirection.value === 'asc' ? 1 : -1
  if (sortMode.value === 'duration') {
    items.sort((a, b) => (a.durationMinutes - b.durationMinutes) * direction)
  } else if (sortMode.value === 'price') {
    items.sort((a, b) => (minPrice(a) - minPrice(b)) * direction)
  } else {
    items.sort((a, b) => a.departTime.localeCompare(b.departTime) * direction)
  }
  return items
})

onMounted(async () => {
  stationLoading.value = true

  try {
    stations.value = await fetchStations()
    if (stations.value.length >= 2) {
      const first = stations.value[0]
      const target = stations.value.find((station) => station.name === '上海虹桥') || stations.value[stations.value.length - 1]
      if (first && target) {
        query.fromStationId = String(first.id)
        query.toStationId = String(target.id)
      }
      await handleSearch()
    }
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
    decrementSeatAvailability(train, seat)
    successMessage.value = '订单已创建，请在订单页完成模拟支付。'
    await notificationStore.refresh()
    await router.push({ name: 'orders' })
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    bookingKey.value = ''
  }
}

function swapStations() {
  const from = query.fromStationId
  query.fromStationId = query.toStationId
  query.toStationId = from
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

async function selectQuickDate(date: string) {
  query.date = date
  if (canSearch.value) {
    await handleSearch()
  }
}

function minPrice(train: TrainSearchItem) {
  return Math.min(...train.seatOptions.map((seat) => seat.priceCents))
}

function trainType(trainNo: string) {
  if (trainNo.startsWith('G')) return 'G 高铁'
  if (trainNo.startsWith('D')) return 'D 动车'
  return `${trainNo.slice(0, 1)} 普客`
}

function sortLabel() {
  if (sortMode.value === 'duration') return sortDirection.value === 'asc' ? '历时最短优先' : '历时最长优先'
  if (sortMode.value === 'price') return sortDirection.value === 'asc' ? '票价最低优先' : '票价最高优先'
  return sortDirection.value === 'asc' ? '出发最早优先' : '出发最晚优先'
}
</script>

<template>
  <main>
    <PageHeader title="车次自主查询" description="按乘车日期、出发站和到达站查询可售车票，享极速便捷的网上订票体验。" />

    <section class="mx-auto max-w-7xl px-4 py-6 sm:px-8">
      <form class="rounded-lg border border-slate-200 bg-white p-6 shadow-sm" @submit.prevent="handleSearch">
        <div class="grid gap-5 xl:grid-cols-[1fr_auto_1fr_1fr_auto] xl:items-end">
          <label class="block">
            <span class="form-label">出发城市</span>
            <select v-model="query.fromStationId" class="form-input mt-2 h-14 text-lg" :disabled="stationLoading" required>
              <option value="">请选择</option>
              <option v-for="station in stations" :key="station.id" :value="station.id">{{ station.name }}</option>
            </select>
          </label>

          <button class="mt-7 h-14 w-14 rounded-full border border-slate-200 bg-slate-50 text-2xl text-slate-500 shadow-sm" type="button" @click="swapStations">
            ⇆
          </button>

          <label class="block">
            <span class="form-label">到达城市</span>
            <select v-model="query.toStationId" class="form-input mt-2 h-14 text-lg" :disabled="stationLoading" required>
              <option value="">请选择</option>
              <option v-for="station in stations" :key="station.id" :value="station.id">{{ station.name }}</option>
            </select>
          </label>

          <label class="block">
            <span class="form-label">出发日期</span>
            <input v-model="query.date" class="form-input mt-2 h-14 text-lg" type="date" required />
          </label>

          <button class="btn-primary h-14 px-10 text-lg" type="submit" :disabled="loading || !canSearch">
            {{ loading ? '查询中...' : '查询车次' }}
          </button>
        </div>

        <div class="mt-6 border-t border-slate-100 pt-5">
          <div class="flex flex-wrap items-center gap-3">
            <span class="text-sm font-bold text-slate-400">快捷日期:</span>
            <button
              v-for="item in quickDates"
              :key="item.value"
              class="rounded-lg px-5 py-2 text-sm font-bold transition"
              :class="query.date === item.value ? 'bg-emerald-50 text-emerald-600 shadow-sm ring-1 ring-emerald-100' : 'bg-slate-50 text-slate-500 hover:bg-slate-100'"
              type="button"
              :disabled="loading"
              @click="selectQuickDate(item.value)"
            >
              {{ item.label }} ({{ item.text }})
            </button>
          </div>
        </div>

        <p v-if="errorMessage" class="mt-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>
        <p v-if="successMessage" class="mt-4 rounded-lg bg-emerald-50 px-3 py-2 text-sm text-emerald-700">{{ successMessage }}</p>
      </form>

      <div class="sticky top-20 z-20 mt-6 rounded-lg border border-slate-200 bg-white/95 px-4 py-3 shadow-sm backdrop-blur">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <label class="inline-flex items-center gap-2 rounded-lg border border-slate-200 bg-white px-4 py-2 text-sm font-bold text-slate-600">
            <input v-model="onlyHighSpeed" type="checkbox" class="rounded border-slate-300 text-teal-600 focus:ring-teal-500" />
            只看高铁/动车 (G/D)
          </label>

          <div class="flex flex-wrap items-center gap-2 rounded-lg border border-slate-200 bg-white p-1">
            <span class="px-3 text-sm font-bold text-slate-400">排序方式:</span>
            <select v-model="sortMode" class="h-9 rounded-md border border-slate-200 bg-white px-3 text-sm font-bold text-slate-600">
              <option value="depart">出发时间</option>
              <option value="duration">历时</option>
              <option value="price">票价</option>
            </select>
            <button class="rounded-md bg-slate-50 px-3 py-1.5 text-sm font-bold text-slate-600 hover:bg-emerald-50 hover:text-emerald-600" type="button" @click="sortDirection = sortDirection === 'asc' ? 'desc' : 'asc'">
              {{ sortDirection === 'asc' ? '升序 ↑' : '降序 ↓' }}
            </button>
          </div>
        </div>
      </div>

      <div class="mt-7 flex flex-wrap items-center justify-between gap-2 text-base font-black text-slate-600">
        <span>共发现 {{ visibleTrains.length }} 趟可选车次</span>
        <span class="text-sm text-slate-400">{{ searchSummary }} · {{ sortLabel() }}</span>
      </div>

      <div v-if="visibleTrains.length" class="mt-5 space-y-5">
        <article
          v-for="train in visibleTrains"
          :key="`${train.trainId}-${train.travelDate}`"
          class="grid gap-6 rounded-lg border border-slate-200 bg-white p-6 shadow-sm transition hover:border-emerald-200 xl:grid-cols-[180px_180px_180px_1fr]"
        >
          <div>
            <div class="text-4xl font-black tracking-normal text-slate-950">{{ formatTime(train.departTime) }}</div>
            <div class="mt-2 text-lg font-black text-slate-700">● {{ train.fromStation.name }}</div>
          </div>

          <div class="flex flex-col items-center justify-center text-center">
            <div class="text-sm font-black text-slate-400">{{ formatDuration(train.durationMinutes) }}</div>
            <div class="my-2 h-px w-full bg-slate-200">
              <div class="mx-auto -mt-3 h-6 w-10 rounded-full border border-slate-200 bg-white text-slate-400">⇆</div>
            </div>
            <div class="text-sm font-black text-emerald-600">途经站 ⓘ</div>
          </div>

          <div>
            <div class="text-4xl font-black tracking-normal text-slate-950">{{ formatTime(train.arriveTime) }}</div>
            <div class="mt-2 text-lg font-black text-slate-700">● {{ train.toStation.name }}</div>
          </div>

          <div class="grid gap-4 lg:grid-cols-[120px_1fr]">
            <div class="self-center">
              <div class="inline-flex rounded-lg bg-slate-100 px-5 py-3 text-2xl font-black text-slate-700">{{ train.trainNo }}</div>
              <div class="mt-2 text-center text-sm font-bold text-slate-400">{{ trainType(train.trainNo) }}</div>
            </div>

            <div class="grid gap-3 sm:grid-cols-3">
              <button
                v-for="seat in train.seatOptions"
                :key="seat.seatClassCode"
                class="rounded-lg border border-slate-200 bg-slate-50 p-4 text-left transition hover:border-emerald-200 hover:bg-emerald-50 disabled:cursor-not-allowed disabled:opacity-60"
                type="button"
                :disabled="seat.availableCount <= 0 || bookingKey === `${train.trainId}-${train.travelDate}-${seat.seatClassCode}`"
                @click="handleBook(train, seat)"
              >
                <div class="text-sm font-black text-slate-500">{{ seat.seatClassName }}</div>
                <div class="mt-2 text-xl font-black text-slate-900">{{ formatMoney(seat.priceCents) }}</div>
                <div class="mt-1 text-sm font-bold" :class="seat.availableCount > 0 ? 'text-emerald-600' : 'text-slate-400'">
                  {{ seat.availableCount > 0 ? `余 ${seat.availableCount}` : '无票' }}
                </div>
              </button>
            </div>
          </div>
        </article>
      </div>

      <EmptyState v-else class="mt-5" title="暂无查询结果" description="选择乘车日期、出发站和到达站后查询。" />
    </section>
  </main>
</template>
