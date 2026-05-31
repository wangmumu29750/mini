package dto

type StationResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type SeatOptionResponse struct {
	SeatClassCode  string `json:"seatClassCode"`
	SeatClassName  string `json:"seatClassName"`
	PriceCents     int64  `json:"priceCents"`
	AvailableCount int    `json:"availableCount"`
}

type TrainSearchItemResponse struct {
	TrainID         uint64               `json:"trainId"`
	TrainNo         string               `json:"trainNo"`
	TrainType       string               `json:"trainType"`
	TravelDate      string               `json:"travelDate"`
	FromStation     StationResponse      `json:"fromStation"`
	ToStation       StationResponse      `json:"toStation"`
	DepartTime      string               `json:"departTime"`
	ArriveTime      string               `json:"arriveTime"`
	DurationMinutes int                  `json:"durationMinutes"`
	SeatOptions     []SeatOptionResponse `json:"seatOptions"`
}

type TrainSearchQuery struct {
	Date          string `form:"date" binding:"required"`
	FromStationID uint64 `form:"fromStationId" binding:"required"`
	ToStationID   uint64 `form:"toStationId" binding:"required"`
}
