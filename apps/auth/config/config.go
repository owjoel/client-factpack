package config

import (
	"os"
	"strconv"
)

var (
	Host = os.Getenv("HOST")
	ClientId = os.Getenv("COGNITO_USERPOOL_CLIENT_ID")
	ClientSecret = os.Getenv("COGNITO_USERPOOL_CLIENT_SECRET")
	UserPoolId = os.Getenv("COGNITO_USERPOOL_ID")
	AutoResetPassword = os.Getenv("AUTO_RESET_PASSWORD") // temp flag
)

func GetPort(defaultPort int) int {
	_port, exist := os.LookupEnv("PORT")
	if !exist {
		return defaultPort
	}
	port, err := strconv.Atoi(_port)
	if err != nil {
		return defaultPort
	}
	return port
}