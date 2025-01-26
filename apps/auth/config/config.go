package config

import (
	"os"
	"strconv"
)

var (
	ClientId = os.Getenv("AWS_COGNITO_CLIENT_ID")
	ClientSecret = os.Getenv("AWS_COGNITO_CLIENT_SECRET")
	UserPoolId = os.Getenv("COGNITO_USERPOOL_ID")
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