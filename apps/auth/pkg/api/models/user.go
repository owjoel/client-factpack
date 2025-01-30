package models

type SignUpReq struct {
	Email string `form:"email" example:"example@gmail.com"`
}

type ForgetPasswordReq struct {
	Username string `form:"username" example:"joel.ow.2022"`
}

type User struct {
	Username string
	Password string
}

type LoginReq struct {
	Username string `form:"username" example:"joel.ow.2022"`
	Password string `form:"password" example:"12345"`
}

type ConfirmForgetPasswordReq struct {
	Username string `form:"username" example:"joel.ow.2022"`
	Code     string `form:"code" example:"ABCDEF"`
	NewPassword string `form:"newPassword" example:"67890"`
}

type StatusRes struct {
	Status string `json:"status"`
}