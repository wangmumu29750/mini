package dto

type TicketResponse struct {
	ID             uint64          `json:"id"`
	TicketNo       string          `json:"ticketNo"`
	OrderID        uint64          `json:"orderId"`
	TrainID        uint64          `json:"trainId"`
	TrainNo        string          `json:"trainNo"`
	TravelDate     string          `json:"travelDate"`
	FromStation    StationResponse `json:"fromStation"`
	ToStation      StationResponse `json:"toStation"`
	DepartTime     string          `json:"departTime,omitempty"`
	ArriveTime     string          `json:"arriveTime,omitempty"`
	SeatClassCode  string          `json:"seatClassCode"`
	SeatClassName  string          `json:"seatClassName"`
	TicketType     string          `json:"ticketType"`
	RealPriceCents int64           `json:"realPriceCents"`
	CoachNo        string          `json:"coachNo"`
	SeatNo         string          `json:"seatNo"`
	PassengerName  string          `json:"passengerName"`
	IDCardNoMasked string          `json:"idCardNoMasked"`
	Status         string          `json:"status"`
	IssuedAt       string          `json:"issuedAt"`
	RefundedAt     *string         `json:"refundedAt,omitempty"`
}

type RefundTicketRequest struct {
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotencyKey"`
}

type RefundTicketResponse struct {
	RefundNo          string         `json:"refundNo"`
	RefundAmountCents int64          `json:"refundAmountCents"`
	FeeCents          int64          `json:"feeCents"`
	Ticket            TicketResponse `json:"ticket"`
}

type ChangeTicketRequest struct {
	NewTrainID       uint64 `json:"newTrainId" binding:"required"`
	NewTravelDate    string `json:"newTravelDate" binding:"required"`
	NewSeatClassCode string `json:"newSeatClassCode" binding:"required"`
	IdempotencyKey   string `json:"idempotencyKey"`
}

type ChangeTicketResponse struct {
	ChangeNo        string         `json:"changeNo"`
	PriceDiffCents  int64          `json:"priceDiffCents"`
	FeeCents        int64          `json:"feeCents"`
	SettlementCents int64          `json:"settlementCents"`
	PaymentNo       string         `json:"paymentNo,omitempty"`
	RefundNo        string         `json:"refundNo,omitempty"`
	OldTicket       TicketResponse `json:"oldTicket"`
	NewTicket       TicketResponse `json:"newTicket"`
}

type ChangeOptionsQuery struct {
	Date string `form:"date" binding:"required"`
}

type ChangeOptionsResponse struct {
	OriginalTicket TicketResponse            `json:"originalTicket"`
	Options        []TrainSearchItemResponse `json:"options"`
}
