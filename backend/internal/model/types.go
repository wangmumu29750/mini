package model

type UserRole string

const (
	UserRolePassenger UserRole = "PASSENGER"
	UserRoleAdmin     UserRole = "ADMIN"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "ACTIVE"
	UserStatusDisabled UserStatus = "DISABLED"
)

type VerificationStatus string

const (
	VerificationStatusPending  VerificationStatus = "PENDING"
	VerificationStatusVerified VerificationStatus = "VERIFIED"
	VerificationStatusFailed   VerificationStatus = "FAILED"
)

type TrainStatus string

const (
	TrainStatusActive   TrainStatus = "ACTIVE"
	TrainStatusDisabled TrainStatus = "DISABLED"
)

type StationStatus string

const (
	StationStatusActive   StationStatus = "ACTIVE"
	StationStatusDisabled StationStatus = "DISABLED"
)

type InventoryStatus string

const (
	InventoryStatusActive InventoryStatus = "ACTIVE"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusClosed         OrderStatus = "CLOSED"
)

type PaymentStatus string

const (
	PaymentStatusSuccess PaymentStatus = "SUCCESS"
)

type TicketStatus string

const (
	TicketStatusIssued TicketStatus = "ISSUED"
)
