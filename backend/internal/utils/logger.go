package utils

import (
	"os"
	"strings"
	"html"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// NewLogger creates a new instance of a logrus.Logger with the specified log level and format.
// The log level can be any valid logrus log level (e.g., "debug", "info", "warn", "error").
// If the provided log level is invalid, it defaults to "info" level.
// The log format can be either "json" or "JSON" for JSON formatted logs, or any other value for text formatted logs.
// The logger output is set to os.Stdout.
//
// Parameters:
//   - logLevel: The desired log level as a string.
//   - logFormat: The desired log format as a string.
//
// Returns:
//   - *logrus.Logger: A pointer to the configured logrus.Logger instance.
func NewLogger(logLevel string, logFormat string) *logrus.Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel // Default to Info level if parsing fails
	}
	logger.SetLevel(level)

	switch logFormat {
	case "json", "JSON":
		logger.SetFormatter(&logrus.JSONFormatter{})
	default:
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	logger.SetOutput(os.Stdout)

	return logger
}

func InitGlobalLogger(logLevel, logFormat string) {
	Logger = NewLogger(logLevel, logFormat)
}


// ObfuscateSensitiveFields takes a map of logrus fields and obfuscates the value of the "jwt" field if it exists.
// If the "jwt" field is a string and its length is greater than 10, it replaces the middle part of the string with asterisks.
// The first and last 5 characters of the "jwt" string are preserved.
// All other fields are returned unchanged.
//
// Parameters:
//   fields (logrus.Fields): A map of logrus fields to be processed.
//
// Returns:
//   logrus.Fields: A new map with the "jwt" field obfuscated if applicable.
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


// cleanFields sanitizes the values of the provided logrus.Fields map by removing
// certain control characters (newline, carriage return, tab, backspace, form feed)
// and escaping HTML special characters, backslashes, double quotes, and single quotes.
// This ensures that the log fields are safe for logging and do not contain any
// potentially harmful or malformed data.
//
// Parameters:
//   fields (logrus.Fields): A map of log fields to be sanitized.
//
// Returns:
//   logrus.Fields: A new map with sanitized log field values.
func cleanFields(fields logrus.Fields) logrus.Fields {
	for key, value := range fields {
		if str, ok := value.(string); ok {
			str = strings.ReplaceAll(str, "\n", "")
			str = strings.ReplaceAll(str, "\r", "")
			str = strings.ReplaceAll(str, "\t", "")
			str = strings.ReplaceAll(str, "\b", "")
			str = strings.ReplaceAll(str, "\f", "")
			str = html.EscapeString(str)
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
