package web

type RegisterRequest struct {
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
	Name         string `json:"name" validate:"required"`
	MobileNumber string `json:"mobileNumber" validate:"required"`
}

type VerifyRegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Otp      string `json:"otp" validate:"required"`
}
