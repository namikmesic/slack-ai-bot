package logger

import (
	"log"
	"os"
)

// For now very rudimentary logger
func New() *log.Logger {
	logger := log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags)

	return logger
}
