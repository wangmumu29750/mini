package model

type UserRole string

const (
	UserRolePassenger UserRole = "PASSENGER"
	UserRoleClerk     UserRole = "CLERK"
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

type PassengerType string

const (
	PassengerTypeAdult   PassengerType = "ADULT"
	PassengerTypeStudent PassengerType = "STUDENT"
	PassengerTypeChild   PassengerType = "CHILD"
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

type RefundStatus string

const (
	RefundStatusSuccess RefundStatus = "SUCCESS"
)

type ChangeStatus string

const (
	ChangeStatusSuccess ChangeStatus = "SUCCESS"
)

type TicketStatus string

const (
	TicketStatusIssued     TicketStatus = "ISSUED"
	TicketStatusRefunded   TicketStatus = "REFUNDED"
	TicketStatusChangedOut TicketStatus = "CHANGED_OUT"
)

type TicketType string

const (
	TicketTypeAdult   TicketType = "ADULT"
	TicketTypeStudent TicketType = "STUDENT"
	TicketTypeChild   TicketType = "CHILD"
)

type SeatType string

const (
	SeatTypeSecond   SeatType = "SECOND"
	SeatTypeFirst    SeatType = "FIRST"
	SeatTypeBusiness SeatType = "BUSINESS"
)
