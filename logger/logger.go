package logger

import (
	"log"
	"os"
	"time"
)

var (
	InfoLogger  *log.Logger
	DebugLogger *log.Logger
)

func init() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	InfoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(logFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

func Debug(format string, v ...interface{}) {
	DebugLogger.Printf(format, v...)
}

func LogRequest(method, path string, status int, duration time.Duration) {
	Info("Request: %s %s - Status: %d - Duration: %v", method, path, status, duration)
}

func LogError(err error, context string) {
	Info("Error in %s: %v", context, err)
}
