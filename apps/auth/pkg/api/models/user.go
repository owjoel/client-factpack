package models

import "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

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

type SetNewPasswordReq struct {
	Username string `form:"username"`
	NewPassword string `form:"newPassword"`
	Session string
}

type VerifyMFAReq struct {
	Code string `form:"code"`
	Session string
}

type AuthenticationRes struct {
	Result types.AuthenticationResultType
}

type SignInMFAReq struct {
	Username string `form:"username"`
	Code string `form:"code"`
	Session string
}

type StatusRes struct {
	Status string `json:"status"`
}