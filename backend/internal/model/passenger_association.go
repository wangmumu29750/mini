package model

type PassengerAssociation struct {
	BaseModel

	OwnerUserID        uint64        `gorm:"not null;index;uniqueIndex:idx_passenger_association_owner_profile,priority:1" json:"ownerUserId"`
	PassengerProfileID uint64        `gorm:"not null;index;uniqueIndex:idx_passenger_association_owner_profile,priority:2" json:"passengerProfileId"`
	PassengerType      PassengerType `gorm:"size:20;not null;default:ADULT" json:"passengerType"`
}
