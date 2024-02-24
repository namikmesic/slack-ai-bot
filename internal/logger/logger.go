package logger

import (
	"log"
	"os"
	// Assuming `slog` is the structured logging library
	// Replace with the actual import path
)

// New creates a new logger instance.
func New(environment string) *log.Logger {
	logger := log.New(os.Stdout, environment, log.Lshortfile|log.LstdFlags)
	return logger
}
