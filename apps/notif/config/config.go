package config

import (
	"os"
)

var (
	DBUser     = os.Getenv("DB_USERNAME")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBHost     = os.Getenv("DB_HOST")
	DBPort     = os.Getenv("DB_PORT")
	DBName     = os.Getenv("DB_NAME")
	RabbitMQURL = os.Getenv("RABBITMQ_URL")
	ServicePort = os.Getenv("PORT")
)

// Loads default values if environment variables are missing
func Load() {
	if RabbitMQURL == "" {
		RabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}
	if ServicePort == "" {
		ServicePort = "8081"
	}
}
