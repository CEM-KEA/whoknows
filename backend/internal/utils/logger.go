package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func NewLogger(logLevel string, logFormat string) *logrus.Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel // Default to Info level if parsing fails
	}
	logger.SetLevel(level)

	// Set log output format
	switch logFormat {
	case "json", "JSON":
		logger.SetFormatter(&logrus.JSONFormatter{})
	default:
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	// Set output to stdout
	logger.SetOutput(os.Stdout)

	return logger
}

func InitGlobalLogger(logLevel, logFormat string) {
	Logger = NewLogger(logLevel, logFormat)
}

// LogDebug logs a debug message with specific fields
func LogDebug(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Debug(message)
}

// LogInfo logs an informational message with specific fields
func LogInfo(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Info(message)
}

// LogError logs an error with a specific context
func LogError(err error, message string, fields logrus.Fields) {
	Logger.WithFields(fields).WithError(err).Warn(message)
}

// LogWarn logs a warning message with specific fields
func LogWarn(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Warn(message)
}

// LogFatal logs a fatal message with specific fields
func LogFatal(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Fatal(message)
}

// LogPanic logs a panic message with specific fields
func LogPanic(message string, fields logrus.Fields) {
	Logger.WithFields(fields).Panic(message)
}
