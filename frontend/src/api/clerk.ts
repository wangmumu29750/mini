import { request } from '@/api/http'
import type { ClerkCreateOrderPayload, Order } from '@/types/domain'

export function createClerkOrder(payload: ClerkCreateOrderPayload) {
  return request<Order>({
    method: 'POST',
    url: '/clerk/orders',
    data: payload,
  })
}
