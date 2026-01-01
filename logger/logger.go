package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// Log is the global logger instance
// It is initialized on package import with sensible defaults
var Log zerolog.Logger

func init() {
	// Use human-readable console output for development
	output := zerolog.ConsoleWriter{Out: os.Stderr}

	// Set global log level from environment variable (default: info)
	logLevel := os.Getenv("LOG_LEVEL")
	var level zerolog.Level
	switch logLevel {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn", "warning":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	default:
		level = zerolog.InfoLevel
	}

	Log = zerolog.New(output).Level(level).With().
		Timestamp().
		Logger()
}
