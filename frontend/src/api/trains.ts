import { request } from '@/api/http'
import type { Station, TrainSearchItem } from '@/types/domain'

export interface TrainSearchParams {
  date: string
  fromStationId: number | string
  toStationId: number | string
}

export function fetchStations() {
  return request<Station[]>({
    method: 'GET',
    url: '/stations',
  })
}

export function searchTrains(params: TrainSearchParams) {
  return request<TrainSearchItem[]>({
    method: 'GET',
    url: '/trains/search',
    params,
  })
}

