package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DBUser      string
	DBPassword  string
	DBHost      string
	DBPort      string
	DBName      string
	RabbitMQURL string
	Port        string
)

// Loads default values if environment variables are missing
func Load() {
	err := godotenv.Load("config/.env") // path to your .env file
	if err != nil {
		log.Println("No .env file found or unable to load .env")
	}

	// Now load the env variables
	DBUser = os.Getenv("DB_USERNAME")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBHost = os.Getenv("DB_HOST")
	DBPort = os.Getenv("DB_PORT")
	DBName = os.Getenv("DB_NAME")
	RabbitMQURL = os.Getenv("RABBITMQ_URL")
	Port = os.Getenv("PORT")

	// Set defaults
	if RabbitMQURL == "" {
		RabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}
	if Port == "" {
		Port = "8082"
	}
}
