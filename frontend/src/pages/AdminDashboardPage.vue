<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import { fetchStations, searchTrains } from '@/api/trains'
import EmptyState from '@/components/EmptyState.vue'
import PageHeader from '@/components/PageHeader.vue'
import type { ApiErrorPayload } from '@/types/api'
import type { Station, TrainSearchItem } from '@/types/domain'
import { formatMoney, formatTime } from '@/utils/format'

const stations = ref<Station[]>([])
const todayTrains = ref<TrainSearchItem[]>([])
const tomorrowTrains = ref<TrainSearchItem[]>([])
const loading = ref(false)
const errorMessage = ref('')

const today = dateText(0)
const tomorrow = dateText(1)

const totalTrainCount = computed(() => todayTrains.value.length + tomorrowTrains.value.length)
const routeName = computed(() => {
  if (stations.value.length < 2) return '暂无可查线路'
  return `${stations.value[0]?.name} → ${stations.value[stations.value.length - 1]?.name}`
})
const seatCount = computed(() => [...todayTrains.value, ...tomorrowTrains.value].reduce((sum, train) => sum + train.seatOptions.length, 0))
const lowestPrice = computed(() => {
  const prices = [...todayTrains.value, ...tomorrowTrains.value].flatMap((train) => train.seatOptions.map((seat) => seat.priceCents))
  return prices.length ? Math.min(...prices) : 0
})

const modules = [
  { title: '站点管理', status: '后端管理接口未接入', description: '当前可读取旅客端站点列表，新增、编辑、停用站点仍需 /admin/stations。' },
  { title: '车次管理', status: '后端管理接口未接入', description: '当前可通过查询接口核对可售车次，维护车次仍需 /admin/trains。' },
  { title: '经停管理', status: '后端管理接口未接入', description: '经停数据已参与查询和出票时间计算，覆盖保存经停仍需 /admin/trains/{id}/stops。' },
  { title: '票额管理', status: '后端管理接口未接入', description: '票额会在购票、支付、退票、改签中流转，后台编辑仍需 /admin/inventories。' },
]

onMounted(loadDashboard)

async function loadDashboard() {
  loading.value = true
  errorMessage.value = ''

  try {
    stations.value = await fetchStations()
    const fromStation = stations.value[0]
    const toStation = stations.value[stations.value.length - 1]
    if (fromStation && toStation && fromStation.id !== toStation.id) {
      const fromStationId = fromStation.id
      const toStationId = toStation.id
      const [todayResult, tomorrowResult] = await Promise.all([
        searchTrains({ date: today, fromStationId, toStationId }),
        searchTrains({ date: tomorrow, fromStationId, toStationId }),
      ])
      todayTrains.value = todayResult
      tomorrowTrains.value = tomorrowResult
    }
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    loading.value = false
  }
}

function dateText(offset: number) {
  const date = new Date()
  date.setDate(date.getDate() + offset)
  return date.toISOString().slice(0, 10)
}

function minPrice(train: TrainSearchItem) {
  return Math.min(...train.seatOptions.map((seat) => seat.priceCents))
}
</script>

<template>
  <main>
    <PageHeader title="管理后台" description="查看基础数据接入情况，维护接口接入前先用于核对演示数据是否可查询、可售。" />

    <section class="mx-auto max-w-7xl px-4 py-6 sm:px-8">
      <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
        <div class="text-sm font-bold text-slate-400">{{ loading ? '加载中...' : `当前核对线路：${routeName}` }}</div>
        <button class="btn-secondary" type="button" :disabled="loading" @click="loadDashboard">
          {{ loading ? '刷新中...' : '刷新数据' }}
        </button>
      </div>

      <p v-if="errorMessage" class="mb-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>

      <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <article class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
          <div class="text-sm font-black text-slate-400">启用站点</div>
          <div class="mt-3 text-4xl font-black text-slate-950">{{ stations.length }}</div>
          <div class="mt-2 text-sm font-bold text-slate-500">来自 /stations</div>
        </article>
        <article class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
          <div class="text-sm font-black text-slate-400">两日可售车次</div>
          <div class="mt-3 text-4xl font-black text-slate-950">{{ totalTrainCount }}</div>
          <div class="mt-2 text-sm font-bold text-slate-500">今日 + 明日查询结果</div>
        </article>
        <article class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
          <div class="text-sm font-black text-slate-400">席别报价项</div>
          <div class="mt-3 text-4xl font-black text-slate-950">{{ seatCount }}</div>
          <div class="mt-2 text-sm font-bold text-slate-500">按可售车次席别统计</div>
        </article>
        <article class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
          <div class="text-sm font-black text-slate-400">最低票价</div>
          <div class="mt-3 text-4xl font-black text-rose-600">{{ lowestPrice ? formatMoney(lowestPrice) : '-' }}</div>
          <div class="mt-2 text-sm font-bold text-slate-500">用于核对票额价格</div>
        </article>
      </div>

      <div class="mt-6 grid gap-6 xl:grid-cols-[1fr_360px]">
        <section class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
          <div class="mb-4 flex items-center justify-between">
            <h2 class="text-lg font-black text-slate-950">可售车次抽检</h2>
            <span class="text-sm font-bold text-slate-400">{{ today }} / {{ tomorrow }}</span>
          </div>

          <div v-if="[...todayTrains, ...tomorrowTrains].length" class="space-y-3">
            <div
              v-for="train in [...todayTrains, ...tomorrowTrains].slice(0, 8)"
              :key="`${train.trainId}-${train.travelDate}`"
              class="grid gap-3 rounded-lg border border-slate-100 bg-slate-50 p-4 sm:grid-cols-[100px_1fr_120px]"
            >
              <div>
                <div class="text-lg font-black text-slate-900">{{ train.trainNo }}</div>
                <div class="text-xs font-bold text-slate-400">{{ train.travelDate }}</div>
              </div>
              <div class="text-sm font-bold text-slate-600">
                {{ train.fromStation.name }} {{ formatTime(train.departTime) }} → {{ train.toStation.name }} {{ formatTime(train.arriveTime) }}
              </div>
              <div class="text-right text-sm font-black text-rose-600">起 {{ formatMoney(minPrice(train)) }}</div>
            </div>
          </div>

          <EmptyState v-else title="暂无可抽检车次" description="请确认种子数据已写入，并选择有库存的出发站和到达站。" />
        </section>

        <section class="space-y-4">
          <article v-for="item in modules" :key="item.title" class="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
            <div class="flex items-start justify-between gap-3">
              <h2 class="text-base font-black text-slate-950">{{ item.title }}</h2>
              <span class="rounded-md bg-amber-50 px-2 py-1 text-xs font-black text-amber-600">{{ item.status }}</span>
            </div>
            <p class="mt-3 text-sm font-medium leading-6 text-slate-500">{{ item.description }}</p>
          </article>
        </section>
      </div>
    </section>
  </main>
</template>
