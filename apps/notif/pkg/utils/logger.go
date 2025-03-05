package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func InitLogger() {
	// Set log output to stdout
	Logger.Out = os.Stdout

	// Set log format to JSON for structured logging
	Logger.SetFormatter(&logrus.JSONFormatter{})

	// Set log level (Debug, Info, Warning, Error, Fatal, Panic)
	Logger.SetLevel(logrus.InfoLevel)
}
