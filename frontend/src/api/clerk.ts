import { request } from '@/api/http'
import type { ClerkCreateOrderPayload, PaymentResult } from '@/types/domain'

export function createClerkOrder(payload: ClerkCreateOrderPayload) {
  return request<PaymentResult>({
    method: 'POST',
    url: '/clerk/orders',
    data: payload,
  })
}
