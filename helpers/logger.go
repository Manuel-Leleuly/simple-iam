package helpers

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func NewLogger() *log.Logger {
	logger := log.New()

	logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = log.InfoLevel
	}

	logger.SetLevel(logLevel)
	logger.SetFormatter(&log.JSONFormatter{})

	return logger
}
