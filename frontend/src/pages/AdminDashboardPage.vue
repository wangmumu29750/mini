<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'

import {
  createStation,
  createTrain,
  deleteTrain,
  disableStation,
  fetchAdminStations,
  fetchAdminTrains,
  fetchInventories,
  fetchQuoteStats,
  fetchSellableStats,
  fetchTrainStops,
  flowInventory,
  saveInventory,
  saveTrainStops,
  updateStation,
  updateTrain,
  type SaveInventoryPayload,
  type SaveStationPayload,
  type SaveTrainPayload,
  type SaveTrainStopPayload,
} from '@/api/admin'
import AdminInventoryPanel from '@/components/admin/AdminInventoryPanel.vue'
import AdminStationsPanel from '@/components/admin/AdminStationsPanel.vue'
import AdminStopsPanel from '@/components/admin/AdminStopsPanel.vue'
import AdminSummaryCards from '@/components/admin/AdminSummaryCards.vue'
import AdminTrainsPanel from '@/components/admin/AdminTrainsPanel.vue'
import PageHeader from '@/components/PageHeader.vue'
import type { ApiErrorPayload } from '@/types/api'
import type { AdminStation, AdminTrain, Inventory, InventoryQuoteStats, SellableTrainStat, TrainStop } from '@/types/domain'

const tabs = [
  { key: 'stations', label: '站点' },
  { key: 'trains', label: '车次' },
  { key: 'stops', label: '经停' },
  { key: 'inventory', label: '票额' },
] as const

type AdminTab = (typeof tabs)[number]['key']

const activeTab = ref<AdminTab>('stations')
const loading = ref(false)
const saving = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const stationList = ref<AdminStation[]>([])
const activeStationTotal = ref(0)
const trains = ref<AdminTrain[]>([])
const inventories = ref<Inventory[]>([])
const stats = ref<SellableTrainStat[]>([])
const quoteStats = ref<InventoryQuoteStats | null>(null)
const stops = ref<TrainStop[]>([])

const selectedTrainId = ref<number | null>(null)
const selectedInventoryId = ref<number | null>(null)
const flowAction = ref('LOCK')
const flowQuantity = ref(1)

const stationForm = reactive<SaveStationPayload & { id?: number }>({
  code: '',
  name: '',
  city: '',
  status: 'ACTIVE',
})

const trainForm = reactive<SaveTrainPayload & { id?: number }>({
  trainNo: '',
  trainType: 'G',
  status: 'ACTIVE',
})

const inventoryForm = reactive<SaveInventoryPayload>({
  trainId: 0,
  travelDate: dateText(1),
  fromStationId: 0,
  toStationId: 0,
  seatClassCode: 'SECOND',
  priceCents: 0,
  totalCount: 0,
  availableCount: 0,
  lockedCount: 0,
  soldCount: 0,
  status: 'ACTIVE',
})

const stopDrafts = ref<SaveTrainStopPayload[]>([])

const totalSellableTrains = computed(() => stats.value.reduce((sum, item) => sum + item.trainCount, 0))
const lowestPrice = computed(() => {
  const prices = inventories.value.map((item) => item.priceCents).filter(Boolean)
  return quoteStats.value?.lowestPriceCents || (prices.length ? Math.min(...prices) : 0)
})
const selectedTrain = computed(() => trains.value.find((train) => train.id === selectedTrainId.value))

onMounted(loadAll)

async function loadAll() {
  loading.value = true
  errorMessage.value = ''
  successMessage.value = ''

  try {
    const [stationResult, trainResult, inventoryResult] = await Promise.all([
      fetchAdminStations({ page: 1, pageSize: 100 }),
      fetchAdminTrains({ page: 1, pageSize: 100 }),
      fetchInventories({ page: 1, pageSize: 100 }),
    ])

    stationList.value = stationResult.items
    activeStationTotal.value = stationResult.activeTotal
    trains.value = trainResult.items
    inventories.value = inventoryResult.items

    selectedTrainId.value ||= trains.value[0]?.id || null
    selectedInventoryId.value ||= inventories.value[0]?.id || null
    seedInventoryDefaults()

    await Promise.all([loadStops(), loadStats(), loadQuoteStats()])
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    loading.value = false
  }
}

async function loadStops() {
  if (!selectedTrainId.value) return

  stops.value = await fetchTrainStops(selectedTrainId.value)
  stopDrafts.value = stops.value.map((stop) => ({
    stationId: stop.station.id,
    stopOrder: stop.stopOrder,
    dayOffset: stop.dayOffset,
    arriveClock: stop.arriveClock,
    departClock: stop.departClock,
    mileage: stop.mileage,
  }))
}

async function loadStats() {
  const activeStations = stationList.value.filter((station) => station.status === 'ACTIVE')
  const from = activeStations[0]
  const to = activeStations[activeStations.length - 1]
  stats.value = from && to && from.id !== to.id ? await fetchSellableStats(from.id, to.id) : []
}

async function loadQuoteStats() {
  quoteStats.value = selectedTrainId.value ? await fetchQuoteStats(selectedTrainId.value) : null
}

async function handleSaveStation() {
  await runSave(async () => {
    if (stationForm.id) {
      await updateStation(stationForm.id, stationForm)
      successMessage.value = '站点已更新'
    } else {
      await createStation(stationForm)
      successMessage.value = '站点已新增'
    }
    resetStationForm()
    await loadAll()
  })
}

async function handleDisableStation(station: AdminStation) {
  await runSave(async () => {
    await disableStation(station.id)
    successMessage.value = '站点已停用'
    await loadAll()
  })
}

async function handleSaveTrain() {
  await runSave(async () => {
    if (trainForm.id) {
      await updateTrain(trainForm.id, trainForm)
      successMessage.value = '车次已更新'
    } else {
      await createTrain(trainForm)
      successMessage.value = '车次已新增'
    }
    resetTrainForm()
    await loadAll()
  })
}

async function handleDeleteTrain(train: AdminTrain) {
  await runSave(async () => {
    await deleteTrain(train.id)
    successMessage.value = '车次已停用'
    await loadAll()
  })
}

async function handleSaveStops() {
  if (!selectedTrainId.value) return

  await runSave(async () => {
    await saveTrainStops(selectedTrainId.value!, stopDrafts.value)
    successMessage.value = '经停数据已保存'
    await loadStops()
  })
}

async function handleSaveInventory() {
  await runSave(async () => {
    await saveInventory(inventoryForm)
    successMessage.value = '票额报价已保存'
    await loadAll()
  })
}

async function handleFlowInventory() {
  if (!selectedInventoryId.value) return

  await runSave(async () => {
    await flowInventory(selectedInventoryId.value!, flowAction.value, flowQuantity.value)
    successMessage.value = '票额流转已完成'
    await loadAll()
  })
}

async function runSave(action: () => Promise<void>) {
  saving.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    await action()
  } catch (error) {
    errorMessage.value = (error as ApiErrorPayload).message
  } finally {
    saving.value = false
  }
}

function editStation(station: AdminStation) {
  Object.assign(stationForm, {
    id: station.id,
    code: station.code,
    name: station.name,
    city: station.city,
    status: station.status,
  })
  activeTab.value = 'stations'
}

function editTrain(train: AdminTrain) {
  Object.assign(trainForm, {
    id: train.id,
    trainNo: train.trainNo,
    trainType: train.trainType,
    status: train.status,
  })
  selectedTrainId.value = train.id
  activeTab.value = 'trains'
}

function showStops(train: AdminTrain) {
  selectedTrainId.value = train.id
  activeTab.value = 'stops'
  loadStops()
}

function editInventory(item: Inventory) {
  Object.assign(inventoryForm, {
    trainId: item.trainId,
    travelDate: item.travelDate,
    fromStationId: item.fromStation.id,
    toStationId: item.toStation.id,
    seatClassCode: item.seatClassCode,
    priceCents: item.priceCents,
    totalCount: item.totalCount,
    availableCount: item.availableCount,
    lockedCount: item.lockedCount,
    soldCount: item.soldCount,
    status: item.status,
  })
  selectedInventoryId.value = item.id
  activeTab.value = 'inventory'
}

function addStopDraft() {
  stopDrafts.value.push({
    stationId: stationList.value[0]?.id || 0,
    stopOrder: stopDrafts.value.length + 1,
    dayOffset: 0,
    arriveClock: '',
    departClock: '',
    mileage: 0,
  })
}

function removeStopDraft(index: number) {
  stopDrafts.value.splice(index, 1)
}

function resetStationForm() {
  Object.assign(stationForm, { id: undefined, code: '', name: '', city: '', status: 'ACTIVE' })
}

function resetTrainForm() {
  Object.assign(trainForm, { id: undefined, trainNo: '', trainType: 'G', status: 'ACTIVE' })
}

function seedInventoryDefaults() {
  const inventory = inventories.value[0]
  if (inventory && !inventoryForm.trainId) {
    editInventory(inventory)
    return
  }

  const train = trains.value[0]
  const activeStations = stationList.value.filter((station) => station.status === 'ACTIVE')
  inventoryForm.trainId ||= train?.id || 0
  inventoryForm.fromStationId ||= activeStations[0]?.id || 0
  inventoryForm.toStationId ||= activeStations[activeStations.length - 1]?.id || 0
}

function dateText(offset: number) {
  const date = new Date()
  date.setDate(date.getDate() + offset)
  return date.toISOString().slice(0, 10)
}
</script>

<template>
  <main>
    <PageHeader title="管理后台" description="维护站点、车次、经停与票额报价，库存流转由后端服务层事务处理。" />

    <section class="mx-auto max-w-7xl px-4 py-6 sm:px-8">
      <div class="mb-4 flex flex-wrap items-center justify-between gap-3">
        <div class="text-sm font-bold text-slate-400">{{ loading ? '加载中...' : '管理员基础数据维护' }}</div>
        <button class="btn-secondary" type="button" :disabled="loading || saving" @click="loadAll">刷新</button>
      </div>

      <p v-if="errorMessage" class="mb-4 rounded-lg bg-red-50 px-3 py-2 text-sm text-red-700">{{ errorMessage }}</p>
      <p v-if="successMessage" class="mb-4 rounded-lg bg-emerald-50 px-3 py-2 text-sm text-emerald-700">{{ successMessage }}</p>

      <AdminSummaryCards
        :active-station-total="activeStationTotal"
        :stats="stats"
        :quote-stats="quoteStats"
        :inventory-count="inventories.length"
        :total-sellable-trains="totalSellableTrains"
        :lowest-price="lowestPrice"
        :selected-train="selectedTrain"
      />

      <div class="mt-6 overflow-x-auto border-b border-slate-200">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          class="mr-2 border-b-2 px-4 py-3 text-sm font-black"
          :class="activeTab === tab.key ? 'border-teal-600 text-teal-700' : 'border-transparent text-slate-500'"
          type="button"
          @click="activeTab = tab.key"
        >
          {{ tab.label }}
        </button>
      </div>

      <AdminStationsPanel
        v-if="activeTab === 'stations'"
        :stations="stationList"
        :form="stationForm"
        :saving="saving"
        @save="handleSaveStation"
        @reset="resetStationForm"
        @edit="editStation"
        @disable="handleDisableStation"
      />

      <AdminTrainsPanel
        v-else-if="activeTab === 'trains'"
        :trains="trains"
        :form="trainForm"
        :saving="saving"
        @save="handleSaveTrain"
        @reset="resetTrainForm"
        @edit="editTrain"
        @stops="showStops"
        @delete="handleDeleteTrain"
      />

      <AdminStopsPanel
        v-else-if="activeTab === 'stops'"
        v-model:selected-train-id="selectedTrainId"
        :trains="trains"
        :stations="stationList"
        :stop-drafts="stopDrafts"
        :saving="saving"
        @load="loadStops"
        @add="addStopDraft"
        @remove="removeStopDraft"
        @save="handleSaveStops"
      />

      <AdminInventoryPanel
        v-else
        v-model:selected-inventory-id="selectedInventoryId"
        v-model:flow-action="flowAction"
        v-model:flow-quantity="flowQuantity"
        :trains="trains"
        :stations="stationList"
        :inventories="inventories"
        :form="inventoryForm"
        :saving="saving"
        @save="handleSaveInventory"
        @flow="handleFlowInventory"
        @edit="editInventory"
      />
    </section>
  </main>
</template>
