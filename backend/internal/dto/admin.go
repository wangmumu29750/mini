package dto

type PageResponse[T any] struct {
	Items    []T   `json:"items"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

type StationListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Status   string `form:"status"`
}

type AdminStationResponse struct {
	ID        uint64 `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	City      string `json:"city"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type StationListResponse struct {
	Items       []AdminStationResponse `json:"items"`
	Page        int                    `json:"page"`
	PageSize    int                    `json:"pageSize"`
	Total       int64                  `json:"total"`
	ActiveTotal int64                  `json:"activeTotal"`
}

type SaveStationRequest struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	City   string `json:"city" binding:"required"`
	Status string `json:"status"`
}

type AdminTrainResponse struct {
	ID             uint64   `json:"id"`
	TrainNo        string   `json:"trainNo"`
	TrainType      string   `json:"trainType"`
	SeatClassCodes []string `json:"seatClassCodes"`
	Status         string   `json:"status"`
	StopCount      int64    `json:"stopCount"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}

type AdminTrainListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Status   string `form:"status"`
	TrainNo  string `form:"trainNo"`
}

type SaveTrainRequest struct {
	TrainNo   string `json:"trainNo" binding:"required"`
	TrainType string `json:"trainType" binding:"required"`
	Status    string `json:"status"`
}

type SellableTrainStatsQuery struct {
	FromStationID uint64 `form:"fromStationId" binding:"required"`
	ToStationID   uint64 `form:"toStationId" binding:"required"`
}

type SellableTrainStatItem struct {
	Date       string `json:"date"`
	TrainCount int64  `json:"trainCount"`
}

type TrainStopResponse struct {
	ID          uint64          `json:"id"`
	TrainID     uint64          `json:"trainId"`
	Station     StationResponse `json:"station"`
	StopOrder   int             `json:"stopOrder"`
	DayOffset   int             `json:"dayOffset"`
	ArriveClock string          `json:"arriveClock"`
	DepartClock string          `json:"departClock"`
	Mileage     int             `json:"mileage"`
}

type SaveTrainStopItem struct {
	StationID   uint64 `json:"stationId" binding:"required"`
	StopOrder   int    `json:"stopOrder" binding:"required"`
	DayOffset   int    `json:"dayOffset"`
	ArriveClock string `json:"arriveClock"`
	DepartClock string `json:"departClock"`
	Mileage     int    `json:"mileage"`
}

type SaveTrainStopsRequest struct {
	Stops []SaveTrainStopItem `json:"stops" binding:"required"`
}

type InventoryResponse struct {
	ID             uint64          `json:"id"`
	TrainID        uint64          `json:"trainId"`
	TrainNo        string          `json:"trainNo"`
	TrainType      string          `json:"trainType"`
	TravelDate     string          `json:"travelDate"`
	FromStation    StationResponse `json:"fromStation"`
	ToStation      StationResponse `json:"toStation"`
	SeatClassCode  string          `json:"seatClassCode"`
	SeatClassName  string          `json:"seatClassName"`
	PriceCents     int64           `json:"priceCents"`
	TotalCount     int             `json:"totalCount"`
	AvailableCount int             `json:"availableCount"`
	LockedCount    int             `json:"lockedCount"`
	SoldCount      int             `json:"soldCount"`
	Status         string          `json:"status"`
	UpdatedAt      string          `json:"updatedAt"`
}

type InventoryListQuery struct {
	Page          int    `form:"page"`
	PageSize      int    `form:"pageSize"`
	TrainID       uint64 `form:"trainId"`
	SeatClassCode string `form:"seatClassCode"`
	Date          string `form:"date"`
}

type SaveInventoryRequest struct {
	TrainID        uint64 `json:"trainId" binding:"required"`
	TravelDate     string `json:"travelDate" binding:"required"`
	FromStationID  uint64 `json:"fromStationId" binding:"required"`
	ToStationID    uint64 `json:"toStationId" binding:"required"`
	SeatClassCode  string `json:"seatClassCode" binding:"required"`
	PriceCents     int64  `json:"priceCents" binding:"required"`
	TotalCount     int    `json:"totalCount" binding:"required"`
	AvailableCount int    `json:"availableCount" binding:"required"`
	LockedCount    int    `json:"lockedCount"`
	SoldCount      int    `json:"soldCount"`
	Status         string `json:"status"`
}

type InventoryQuoteStatsQuery struct {
	TrainID       uint64 `form:"trainId" binding:"required"`
	SeatClassCode string `form:"seatClassCode"`
}

type InventoryQuoteStatsResponse struct {
	TrainID       uint64 `json:"trainId"`
	SeatClassCode string `json:"seatClassCode,omitempty"`
	QuoteCount    int64  `json:"quoteCount"`
	LowestPrice   int64  `json:"lowestPriceCents"`
}

type InventoryFlowRequest struct {
	InventoryID uint64 `json:"inventoryId" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required"`
}

type InventoryFlowResponse struct {
	Inventory   InventoryResponse `json:"inventory"`
	LowestPrice int64             `json:"lowestPriceCents"`
}
