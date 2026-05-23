import { request } from '@/api/http'
import type { ChangeOptionsResult, ChangeResult, ChangeTicketPayload, RefundResult, Ticket } from '@/types/domain'

export function fetchTickets() {
  return request<Ticket[]>({
    method: 'GET',
    url: '/tickets',
  })
}

export function fetchTicket(ticketId: number) {
  return request<Ticket>({
    method: 'GET',
    url: `/tickets/${ticketId}`,
  })
}

export function refundTicket(ticketId: number, reason = '行程变更') {
  return request<RefundResult>({
    method: 'POST',
    url: `/tickets/${ticketId}/refund`,
    data: {
      reason,
      idempotencyKey: `${Date.now()}-${Math.random().toString(16).slice(2)}`,
    },
  })
}

export function changeTicket(ticketId: number, payload: ChangeTicketPayload) {
  return request<ChangeResult>({
    method: 'POST',
    url: `/tickets/${ticketId}/change`,
    data: payload,
  })
}

export function fetchChangeOptions(ticketId: number, date: string) {
  return request<ChangeOptionsResult>({
    method: 'GET',
    url: `/tickets/${ticketId}/change-options`,
    params: { date },
  })
}
