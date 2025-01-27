package models

type SignUpRequest struct {
	Email string `form:"email"`
}

type ForgetPasswordRequest struct {
	Username string `form:"username"`
}

type User struct {
	Username string
	Password string
}

type LoginRequest struct {
	Username string `form:"username"`
	Password string `form:"password"`
}