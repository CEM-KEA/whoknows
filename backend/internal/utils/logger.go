package utils

import (
	"os"
	"strings"

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

// ObfuscateSensitiveFields obfuscates sensitive fields in log messages
func ObfuscateSensitiveFields(fields logrus.Fields) logrus.Fields {
	obfuscatedFields := make(logrus.Fields)
	for key, value := range fields {
		if key == "jwt" {
			strValue, ok := value.(string)
			if ok && len(strValue) > 10 {
				obfuscatedFields[key] = strValue[:5] + "*****" + strValue[len(strValue)-5:]
			} else {
				obfuscatedFields[key] = value
			}
		} else {
			obfuscatedFields[key] = value
		}
	}
	return obfuscatedFields
}

// cleanFields removes newlines, carriage returns, and other potentially malicious characters from string fields
func cleanFields(fields logrus.Fields) logrus.Fields {
	for key, value := range fields {
		if str, ok := value.(string); ok {
			str = strings.ReplaceAll(str, "\n", "")
			str = strings.ReplaceAll(str, "\r", "")
			str = strings.ReplaceAll(str, "\t", "")
			str = strings.ReplaceAll(str, "\b", "")
			str = strings.ReplaceAll(str, "\f", "")
			str = strings.ReplaceAll(str, "\\", "\\\\")
			str = strings.ReplaceAll(str, "\"", "\\\"")
			str = strings.ReplaceAll(str, "'", "\\'")
			fields[key] = str
		}
	}
	return fields
}

// LogDebug logs a debug message with specific fields
func LogDebug(message string, fields logrus.Fields) {
	Logger.WithFields(cleanFields(fields)).Debug(message)
}

// LogInfo logs an informational message with specific fields
func LogInfo(message string, fields logrus.Fields) {
	Logger.WithFields(cleanFields(fields)).Info(message)
}

// LogError logs an error with a specific context
func LogError(err error, message string, fields logrus.Fields) {
	Logger.WithFields(cleanFields(fields)).WithError(err).Warn(message)
}

// LogWarn logs a warning message with specific fields
func LogWarn(message string, fields logrus.Fields) {
	Logger.WithFields(cleanFields(fields)).Warn(message)
}

// LogFatal logs a fatal message with specific fields
func LogFatal(message string, fields logrus.Fields) {
	Logger.WithFields(cleanFields(fields)).Fatal(message)
}

// LogPanic logs a panic message with specific fields
func LogPanic(message string, fields logrus.Fields) {
	Logger.WithFields(cleanFields(fields)).Panic(message)
}
