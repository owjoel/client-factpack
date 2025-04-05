package config

import (
	"os"
	"strconv"
	"strings"
)

var (
	DBUser     = clean(os.Getenv("DB_USERNAME"))
	DBPassword = clean(os.Getenv("DB_PASSWORD"))
	DBHost     = clean(os.Getenv("DB_HOST"))
	DBPort     = clean(os.Getenv("DB_PORT"))
	DBName     = clean(os.Getenv("DB_NAME"))

	MongoURI            = clean(os.Getenv("MONGO_URI"))
	PrefectAPIURL       = clean(os.Getenv("PREFECT_API_URL"))
	PrefectAPIKey       = clean(os.Getenv("PREFECT_API_KEY"))
	PrefectScrapeFlowID = clean(os.Getenv("PREFECT_SCRAPE_FLOW_ID"))

	ClientID     = os.Getenv("COGNITO_USERPOOL_CLIENT_ID")
	ClientSecret = os.Getenv("COGNITO_USERPOOL_CLIENT_SECRET")
	UserPoolID   = os.Getenv("COGNITO_USERPOOL_ID")
	AwsRegion    = os.Getenv("AWS_REGION")
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

func GetVersion() string {
	version, exist := os.LookupEnv("VERSION")
	if !exist {
		return "0.0.1"
	}
	return version
}

func clean(s string) string {
	return strings.Trim(s, "\r\n\t ")
}
