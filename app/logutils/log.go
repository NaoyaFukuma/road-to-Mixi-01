package logutils

import (
	"log"
	"os"
)

const (
	InfoColor    = "\033[1;34m"
	NoticeColor  = "\033[1;36m"
	WarningColor = "\033[1;33m"
	ErrorColor   = "\033[1;31m"
	DebugColor   = "\033[0;36m"
	ResetColor   = "\033[0m"
)

func InitLog() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

// Info logs a message at the info severity level.
func Info(msg string) {
	log.Print(InfoColor + "INFO: " + ResetColor + msg)
}

// Notice logs a message at the notice severity level.
func Notice(msg string) {
	log.Print(NoticeColor, "NOTICE: "+ResetColor+msg)
}

// Warning logs a message at the warning severity level.
func Warning(msg string) {
	log.Print(WarningColor, "WARNING: "+ResetColor+msg)
}

// Error logs a message at the error severity level.
func Error(msg string) {
	log.Print(ErrorColor, "ERROR: "+ResetColor+msg)
}
