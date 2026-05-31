package dto

type CreateOrderRequest struct {
	TrainID        uint64              `json:"trainId" binding:"required"`
	TravelDate     string              `json:"travelDate" binding:"required"`
	FromStationID  uint64              `json:"fromStationId" binding:"required"`
	ToStationID    uint64              `json:"toStationId" binding:"required"`
	Passengers     []OrderPassengerDTO `json:"passengers" binding:"required,min=1,dive"`
	IdempotencyKey string              `json:"idempotencyKey"`
}

type OrderPassengerDTO struct {
	PassengerID uint64 `json:"passengerId" binding:"required"`
	SeatType    string `json:"seatType" binding:"required"`
	TicketType  string `json:"ticketType" binding:"required"`
}

type PassengerSummaryResponse struct {
	ID             uint64 `json:"id"`
	RealName       string `json:"realName"`
	IDCardNoMasked string `json:"idCardNoMasked"`
	PassengerType  string `json:"passengerType"`
}

type ClerkCreateOrderRequest struct {
	CreateOrderRequest
	PassengerName string `json:"passengerName" binding:"required,min=2,max=64"`
	IDCardNo      string `json:"idCardNo" binding:"required,min=6,max=32"`
	Phone         string `json:"phone" binding:"required,min=6,max=20"`
	BankCardNo    string `json:"bankCardNo" binding:"required,min=12,max=32"`
}

type PayOrderRequest struct {
	Channel string `json:"channel"`
}

type OrderResponse struct {
	ID            uint64            `json:"id"`
	OrderNo       string            `json:"orderNo"`
	TrainID       uint64            `json:"trainId"`
	TrainNo       string            `json:"trainNo"`
	TravelDate    string            `json:"travelDate"`
	FromStation   StationResponse   `json:"fromStation"`
	ToStation     StationResponse   `json:"toStation"`
	DepartTime    string            `json:"departTime,omitempty"`
	ArriveTime    string            `json:"arriveTime,omitempty"`
	SeatClassCode string            `json:"seatClassCode"`
	SeatClassName string            `json:"seatClassName"`
	PassengerName string            `json:"passengerName"`
	AmountCents   int64             `json:"amountCents"`
	ItemCount     int               `json:"itemCount"`
	Status        string            `json:"status"`
	PayExpiresAt  string            `json:"payExpiresAt"`
	PaidAt        *string           `json:"paidAt,omitempty"`
	TicketNo      string            `json:"ticketNo,omitempty"`
	TicketStatus  string            `json:"ticketStatus,omitempty"`
	Tickets       []OrderTicketItem `json:"tickets,omitempty"`
}

type OrderTicketItem struct {
	PassengerName  string `json:"passengerName"`
	SeatType       string `json:"seatType"`
	TicketType     string `json:"ticketType"`
	RealPriceCents int64  `json:"realPriceCents"`
	TicketNo       string `json:"ticketNo,omitempty"`
	Status         string `json:"status,omitempty"`
}

type PassengerListResponse struct {
	Items []PassengerSummaryResponse `json:"items"`
}

type PaymentResponse struct {
	PaymentNo string        `json:"paymentNo"`
	Order     OrderResponse `json:"order"`
}
