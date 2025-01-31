package config

import (
	"os"
	"strconv"
)

var (
	DBUser     = os.Getenv("DB_USERNAME")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBHost     = os.Getenv("DB_HOST")
	DBPort     = os.Getenv("DB_PORT")
	DBName     = os.Getenv("DB_NAME")
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
