export interface Station {
  id: number
  name: string
  code?: string
  city?: string
  status?: string
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

export interface PageResult<T> {
  items: T[]
  page: number
  pageSize: number
  total: number
}

export interface AdminStation extends Station {
  code: string
  city: string
  status: 'ACTIVE' | 'DISABLED' | string
  createdAt: string
  updatedAt: string
}

export interface AdminStationList extends PageResult<AdminStation> {
  activeTotal: number
}

export interface AdminTrain {
  id: number
  trainNo: string
  trainType: string
  status: 'ACTIVE' | 'DISABLED' | string
  stopCount: number
  createdAt: string
  updatedAt: string
}

export interface TrainStop {
  id: number
  trainId: number
  station: Station
  stopOrder: number
  dayOffset: number
  arriveClock: string
  departClock: string
  mileage: number
}

export interface Inventory {
  id: number
  trainId: number
  trainNo: string
  travelDate: string
  fromStation: Station
  toStation: Station
  seatClassCode: string
  seatClassName: string
  priceCents: number
  totalCount: number
  availableCount: number
  lockedCount: number
  soldCount: number
  status: string
  updatedAt: string
}

export interface SellableTrainStat {
  date: string
  trainCount: number
}

export interface InventoryQuoteStats {
  trainId: number
  seatClassCode?: string
  quoteCount: number
  lowestPriceCents: number
}

export interface InventoryFlowResult {
  inventory: Inventory
  lowestPriceCents: number
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
