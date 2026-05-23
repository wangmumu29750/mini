import { defineStore } from 'pinia'

import { fetchOrders } from '@/api/orders'
import { fetchTickets } from '@/api/tickets'

interface NotificationState {
  pendingOrderCount: number
  activeTicketCount: number
  loading: boolean
  loaded: boolean
}

export const useNotificationStore = defineStore('notifications', {
  state: (): NotificationState => ({
    pendingOrderCount: 0,
    activeTicketCount: 0,
    loading: false,
    loaded: false,
  }),
  actions: {
    reset() {
      this.pendingOrderCount = 0
      this.activeTicketCount = 0
      this.loading = false
      this.loaded = false
    },
    async refresh() {
      if (this.loading) {
        return
      }

      this.loading = true
      try {
        const [orders, tickets] = await Promise.all([fetchOrders(), fetchTickets()])
        this.pendingOrderCount = orders.filter((order) => order.status === 'PENDING_PAYMENT').length
        this.activeTicketCount = tickets.filter((ticket) => ticket.status === 'ISSUED').length
        this.loaded = true
      } finally {
        this.loading = false
      }
    },
  },
})
