import { request } from '@/api/http'
import type { Order, PaymentResult } from '@/types/domain'

export interface CreateOrderPayload {
  trainId: number
  travelDate: string
  fromStationId: number
  toStationId: number
  seatClassCode: string
  idempotencyKey: string
}

export function createOrder(payload: CreateOrderPayload) {
  return request<Order>({
    method: 'POST',
    url: '/orders',
    data: payload,
  })
}

export function fetchOrders() {
  return request<Order[]>({
    method: 'GET',
    url: '/orders',
  })
}

export function payOrder(orderId: number) {
  return request<PaymentResult>({
    method: 'POST',
    url: `/orders/${orderId}/payments`,
    data: {
      channel: 'MOCK_BANK',
    },
  })
}

export function cancelOrder(orderId: number) {
  return request<Order>({
    method: 'POST',
    url: `/orders/${orderId}/cancel`,
  })
}
