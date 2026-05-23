import { request } from '@/api/http'
import type { Ticket } from '@/types/domain'

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
