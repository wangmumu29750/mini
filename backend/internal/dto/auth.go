package dto

type RegisterRequest struct {
	Username   string `json:"username" binding:"required,min=3,max=64"`
	Password   string `json:"password" binding:"required,min=6,max=72"`
	RealName   string `json:"realName" binding:"required,min=2,max=64"`
	IDCardNo   string `json:"idCardNo" binding:"required,min=6,max=32"`
	Phone      string `json:"phone" binding:"required,min=6,max=20"`
	BankCardNo string `json:"bankCardNo" binding:"required,min=12,max=32"`
	PassengerType string `json:"passengerType" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

type CurrentUserResponse struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type AuthResponse struct {
	AccessToken string              `json:"accessToken"`
	User        CurrentUserResponse `json:"user"`
}

type PassengerProfileRequest struct {
	RealName      string `json:"realName" binding:"required,min=2,max=64"`
	IDCardNo      string `json:"idCardNo" binding:"required,min=6,max=32"`
	Phone         string `json:"phone" binding:"required,min=6,max=20"`
	BankCardNo    string `json:"bankCardNo" binding:"required,min=12,max=32"`
	PassengerType string `json:"passengerType" binding:"required"`
}
