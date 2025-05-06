package logger

import (
	"github.com/sirupsen/logrus"
)

// InitLogger initializes and returns a new instance of a Logrus logger.
//
// The logger is configured to use JSON formatting and is set to the Info level by default.
//
// Returns:
//   - *logrus.Logger: A configured Logrus logger instance.
func InitLogger() *logrus.Logger {
	// Create a new logger instance
	logger := logrus.New()

	// Set log format to JSON
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Set log level to Info (adjustable based on environment)
	logger.SetLevel(logrus.InfoLevel)

	// Return the configured logger instance
	return logger
}
