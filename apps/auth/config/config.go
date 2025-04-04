package config

import (
	"os"
	"strconv"
	"strings"
)

/*
Configuration variables sourced from environment variables:
- Host: The hostname or IP where the service is running.
- ClientID: The Cognito User Pool Client ID used for authentication.
- ClientSecret: The secret associated with the Cognito User Pool Client.
- UserPoolId: The ID of the Cognito User Pool.
- AutoResetPassword: A temporary flag to enable or disable auto password reset.
*/
var (
	Host              = os.Getenv("HOST")
	AllowedDomains	  = GetDomains()
	ClientID          = os.Getenv("COGNITO_USERPOOL_CLIENT_ID")
	ClientSecret      = os.Getenv("COGNITO_USERPOOL_CLIENT_SECRET")
	UserPoolID        = os.Getenv("COGNITO_USERPOOL_ID")
	AwsRegion         = os.Getenv("AWS_REGION")
	AutoResetPassword = os.Getenv("AUTO_RESET_PASSWORD") // temp flag
)

func GetDomains() []string {
	return strings.Split(os.Getenv("ALLOWED_EMAIL_DOMAINS"), ",")
}

// GetPort returns the port number from the environment variable PORT.
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
