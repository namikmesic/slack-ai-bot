package logger

import (
	"log"
	"os"
)

// New creates a new logger instance.
func New(environment string) *log.Logger {
	prefix := "[" + environment + "] "
	logger := log.New(os.Stdout, prefix, log.Lshortfile|log.LstdFlags)

	return logger
}
