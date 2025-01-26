package models

type SignUpRequest struct {
	Email string `form:"email"`
}

type User struct {
	Username string
	Password string
}