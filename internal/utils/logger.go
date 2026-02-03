package utils

import (
	"fmt"
	"time"
)

// LogLevel represents logging levels
type LogLevel int

const (
	LogLevelError LogLevel = iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

// Logger provides structured logging with context
type Logger struct {
	context string
}

// NewLogger creates a new logger with context
func NewLogger(context string) *Logger {
	return &Logger{context: context}
}

// getTimestamp returns current time in Vietnam timezone
func getTimestamp() string {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	return time.Now().In(loc).Format("02/01/2006 15:04:05")
}

// formatMessage formats log message with timestamp and context
func (l *Logger) formatMessage(level, message string) string {
	return fmt.Sprintf("[%s] [%s] [%s] %s", getTimestamp(), level, l.context, message)
}

// Debug logs debug message
func (l *Logger) Debug(message string, args ...interface{}) {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	fmt.Println(l.formatMessage("DEBUG", message))
}

// Info logs info message
func (l *Logger) Info(message string, args ...interface{}) {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	fmt.Println(l.formatMessage("INFO", message))
}

// Warn logs warning message
func (l *Logger) Warn(message string, args ...interface{}) {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	fmt.Println(l.formatMessage("WARN", message))
}

// Error logs error message
func (l *Logger) Error(message string, args ...interface{}) {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	fmt.Println(l.formatMessage("ERROR", message))
}
