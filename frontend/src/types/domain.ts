export interface Station {
  id: number
  name: string
}

export interface SeatOption {
  seatClassCode: string
  seatClassName: string
  priceCents: number
  availableCount: number
}

export interface TrainSearchItem {
  trainId: number
  trainNo: string
  travelDate: string
  fromStation: Station
  toStation: Station
  departTime: string
  arriveTime: string
  durationMinutes: number
  seatOptions: SeatOption[]
}

export interface Order {
  id: number
  orderNo: string
  trainId: number
  trainNo: string
  travelDate: string
  fromStation: Station
  toStation: Station
  seatClassCode: string
  seatClassName: string
  passengerName: string
  amountCents: number
  status: 'PENDING_PAYMENT' | 'CANCELLED' | 'PAID' | 'CLOSED' | string
  payExpiresAt: string
  paidAt?: string
  ticketNo?: string
  ticketStatus?: string
}

export interface PaymentResult {
  paymentNo: string
  order: Order
}

export interface Ticket {
  id: number
  ticketNo: string
  orderId: number
  trainId: number
  trainNo: string
  travelDate: string
  fromStation: Station
  toStation: Station
  seatClassCode: string
  seatClassName: string
  passengerName: string
  idCardNoMasked: string
  status: 'ISSUED' | string
  issuedAt: string
}
