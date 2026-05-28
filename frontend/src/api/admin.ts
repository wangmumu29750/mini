import { request } from '@/api/http'
import type {
  AdminStation,
  AdminStationList,
  AdminTrain,
  Inventory,
  InventoryFlowResult,
  InventoryQuoteStats,
  PageResult,
  SellableTrainStat,
  TrainStop,
} from '@/types/domain'

export interface ListParams {
  page?: number
  pageSize?: number
  status?: string
}

export interface SaveStationPayload {
  code: string
  name: string
  city: string
  status: string
}

export interface SaveTrainPayload {
  trainNo: string
  trainType: string
  status: string
}

export interface SaveTrainStopPayload {
  stationId: number
  stopOrder: number
  dayOffset: number
  arriveClock: string
  departClock: string
  mileage: number
}

export interface SaveInventoryPayload {
  trainId: number
  travelDate: string
  fromStationId: number
  toStationId: number
  seatClassCode: string
  priceCents: number
  totalCount: number
  availableCount: number
  lockedCount: number
  soldCount: number
  status: string
}

export function fetchAdminStations(params: ListParams = {}) {
  return request<AdminStationList>({
    method: 'GET',
    url: '/admin/stations',
    params,
  })
}

export function createStation(data: SaveStationPayload) {
  return request<AdminStation>({
    method: 'POST',
    url: '/admin/stations',
    data,
  })
}

export function updateStation(id: number, data: SaveStationPayload) {
  return request<AdminStation>({
    method: 'PUT',
    url: `/admin/stations/${id}`,
    data,
  })
}

export function disableStation(id: number) {
  return request<AdminStation>({
    method: 'DELETE',
    url: `/admin/stations/${id}`,
  })
}

export function fetchAdminTrains(params: ListParams & { trainNo?: string } = {}) {
  return request<PageResult<AdminTrain>>({
    method: 'GET',
    url: '/admin/trains',
    params,
  })
}

export function createTrain(data: SaveTrainPayload) {
  return request<AdminTrain>({
    method: 'POST',
    url: '/admin/trains',
    data,
  })
}

export function updateTrain(id: number, data: SaveTrainPayload) {
  return request<AdminTrain>({
    method: 'PUT',
    url: `/admin/trains/${id}`,
    data,
  })
}

export function deleteTrain(id: number) {
  return request<AdminTrain>({
    method: 'DELETE',
    url: `/admin/trains/${id}`,
  })
}

export function fetchTrainStops(trainId: number) {
  return request<TrainStop[]>({
    method: 'GET',
    url: `/admin/trains/${trainId}/stops`,
  })
}

export function saveTrainStops(trainId: number, stops: SaveTrainStopPayload[]) {
  return request<TrainStop[]>({
    method: 'PUT',
    url: `/admin/trains/${trainId}/stops`,
    data: { stops },
  })
}

export function fetchSellableStats(fromStationId: number, toStationId: number) {
  return request<SellableTrainStat[]>({
    method: 'GET',
    url: '/admin/trains/sellable-stats',
    params: { fromStationId, toStationId },
  })
}

export function fetchInventories(params: ListParams & { trainId?: number; seatClassCode?: string; date?: string } = {}) {
  return request<PageResult<Inventory>>({
    method: 'GET',
    url: '/admin/inventories',
    params,
  })
}

export function saveInventory(data: SaveInventoryPayload) {
  return request<Inventory>({
    method: 'PUT',
    url: '/admin/inventories',
    data,
  })
}

export function fetchQuoteStats(trainId: number, seatClassCode?: string) {
  return request<InventoryQuoteStats>({
    method: 'GET',
    url: '/admin/inventories/quote-stats',
    params: { trainId, seatClassCode },
  })
}

export function flowInventory(inventoryId: number, action: string, quantity: number) {
  return request<InventoryFlowResult>({
    method: 'POST',
    url: '/admin/inventories/flow',
    data: { inventoryId, action, quantity },
  })
}
