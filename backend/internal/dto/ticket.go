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
	SeatClassCode  string          `json:"seatClassCode"`
	SeatClassName  string          `json:"seatClassName"`
	PassengerName  string          `json:"passengerName"`
	IDCardNoMasked string          `json:"idCardNoMasked"`
	Status         string          `json:"status"`
	IssuedAt       string          `json:"issuedAt"`
}
