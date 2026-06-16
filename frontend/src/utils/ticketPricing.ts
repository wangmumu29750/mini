import type { TicketType } from '@/types/domain'

export function calculateTicketPricePreview(
  basePriceCents: number,
  trainType: string,
  seatClassCode: string,
  ticketType: TicketType,
) {
  if (ticketType === 'CHILD') {
    return Math.round(basePriceCents * 0.5)
  }

  if (ticketType === 'STUDENT') {
    if (!isSeatAllowedForTicketType(ticketType, seatClassCode)) {
      return null
    }
    return Math.round(basePriceCents * studentDiscountFactor(trainType))
  }

  return basePriceCents
}

export function isSeatAllowedForTicketType(ticketType: TicketType, seatClassCode: string) {
  if (ticketType !== 'STUDENT') {
    return true
  }
  return seatClassCode === 'SECOND'
}

function studentDiscountFactor(trainType: string) {
  return ['Z', 'T', 'K'].includes(trainType) ? 0.6 : 0.75
}
