package model

type PassengerProfile struct {
	BaseModel

	UserID         uint64             `gorm:"not null;uniqueIndex" json:"userId"`
	RealName       string             `gorm:"size:64;not null" json:"realName"`
	IDCardNo       string             `gorm:"size:32;not null;uniqueIndex" json:"idCardNo"`
	Phone          string             `gorm:"size:20;not null;uniqueIndex" json:"phone"`
	BankCardNo     string             `gorm:"size:32;not null" json:"bankCardNo"`
	VerifiedStatus VerificationStatus `gorm:"size:20;not null;default:VERIFIED;index" json:"verifiedStatus"`
}
