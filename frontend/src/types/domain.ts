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

export interface ClerkCreateOrderPayload {
  trainId: number
  travelDate: string
  fromStationId: number
  toStationId: number
  seatClassCode: string
  idempotencyKey: string
  passengerName: string
  idCardNo: string
  phone: string
  bankCardNo: string
}

export interface SystemSetting {
  key: string
  value: string
  valueType: 'INT' | 'BOOL' | 'STRING' | string
  description: string
}

export interface Order {
  id: number
  orderNo: string
  trainId: number
  trainNo: string
  travelDate: string
  fromStation: Station
  toStation: Station
  departTime?: string
  arriveTime?: string
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
  departTime?: string
  arriveTime?: string
  seatClassCode: string
  seatClassName: string
  coachNo: string
  seatNo: string
  passengerName: string
  idCardNoMasked: string
  status: 'ISSUED' | 'REFUNDED' | 'CHANGED_OUT' | string
  issuedAt: string
  refundedAt?: string
}

export interface RefundResult {
  refundNo: string
  ticket: Ticket
}

export interface ChangeTicketPayload {
  newTrainId: number
  newTravelDate: string
  newSeatClassCode: string
  idempotencyKey: string
}

export interface ChangeOptionsResult {
  originalTicket: Ticket
  options: TrainSearchItem[]
}

export interface ChangeResult {
  changeNo: string
  priceDiffCents: number
  oldTicket: Ticket
  newTicket: Ticket
}
