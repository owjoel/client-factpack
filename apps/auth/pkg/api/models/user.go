package models

import "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

// User represents a user.
type User struct {
	Username string
	Password string
}

// SignUpReq represents the request payload for user sign up.
type SignUpReq struct {
	Email    string `form:"email" example:"example@gmail.com"`
	Password string `form:"password" binding:"required"`
}

// ForgetPasswordReq represents the request payload for forgotten passwords.
type ForgetPasswordReq struct {
	Username string `form:"username" example:"joel.ow.2022"`
}

// LoginReq represents the request payload for user login.
type LoginReq struct {
	Username string `form:"username" example:"joel.ow.2022"`
	Password string `form:"password" example:"12345"`
}

// LoginRes represents the response payload for user login.
type LoginRes struct {
	Challenge string
	Session   string
}

// ConfirmForgetPasswordReq represents the request payload for confirming a forgotten password.
type ConfirmForgetPasswordReq struct {
	Username    string `form:"username" example:"joel.ow.2022"`
	Code        string `form:"code" example:"ABCDEF"`
	NewPassword string `form:"newPassword" example:"67890"`
}

// SetNewPasswordReq represents the request payload for setting a new password.
type SetNewPasswordReq struct {
	Username    string `form:"username" example:"joel.ow.2022"`
	NewPassword string `form:"newPassword" example:"ABCDEF"`
	Session     string
}

// SetNewPasswordRes represents the response payload for setting a new password.
type SetNewPasswordRes struct {
	Challenge string
	Session   string
}

// VerifyMFAReq represents the request payload for verifying MFA.
type VerifyMFAReq struct {
	Code    string `form:"code"`
	Session string
}

// AuthenticationRes represents the response payload for user authentication.
type AuthenticationRes struct {
	Result    types.AuthenticationResultType
	Challenge string
}

// AuthChallengeRes represents the response payload for an authentication challenge.
type AuthChallengeRes struct {
	Challenge string `json:"challenge" example:"SOFTWARE_TOKEN_MFA" enums:"NEW_PASSWORD_REQUIRED,MFA_SETUP,SOFTWARE_TOKEN_MFA"`
}

// AssociateTokenRes represents the response payload for associating a software token.
type AssociateTokenRes struct {
	Token   string
	Session string
}

// SetupMFARes represents the response payload for setting up MFA.
type SetupMFARes struct {
	Token string `json:"token"`
}

// SignInMFAReq represents the request payload for signing in with MFA.
type SignInMFAReq struct {
	Username string `form:"username"`
	Code     string `form:"code"`
	Session  string
}

// StatusRes represents the response payload for status messages.
type StatusRes struct {
	Status string `json:"status"`
}
