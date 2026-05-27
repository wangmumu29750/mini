package dto

type CreateOrderRequest struct {
	TrainID        uint64 `json:"trainId" binding:"required"`
	TravelDate     string `json:"travelDate" binding:"required"`
	FromStationID  uint64 `json:"fromStationId" binding:"required"`
	ToStationID    uint64 `json:"toStationId" binding:"required"`
	SeatClassCode  string `json:"seatClassCode" binding:"required"`
	IdempotencyKey string `json:"idempotencyKey"`
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
	ID            uint64          `json:"id"`
	OrderNo       string          `json:"orderNo"`
	TrainID       uint64          `json:"trainId"`
	TrainNo       string          `json:"trainNo"`
	TravelDate    string          `json:"travelDate"`
	FromStation   StationResponse `json:"fromStation"`
	ToStation     StationResponse `json:"toStation"`
	DepartTime    string          `json:"departTime,omitempty"`
	ArriveTime    string          `json:"arriveTime,omitempty"`
	SeatClassCode string          `json:"seatClassCode"`
	SeatClassName string          `json:"seatClassName"`
	PassengerName string          `json:"passengerName"`
	AmountCents   int64           `json:"amountCents"`
	Status        string          `json:"status"`
	PayExpiresAt  string          `json:"payExpiresAt"`
	PaidAt        *string         `json:"paidAt,omitempty"`
	TicketNo      string          `json:"ticketNo,omitempty"`
	TicketStatus  string          `json:"ticketStatus,omitempty"`
}

type PaymentResponse struct {
	PaymentNo string        `json:"paymentNo"`
	Order     OrderResponse `json:"order"`
}
