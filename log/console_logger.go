package log

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// ConsoleLogger is a logger that logs to STDOUT (console)
type ConsoleLogger struct {
	LogLevel int

	mutex *sync.Mutex
}

// NewConsoleLogger returns a new ConsoleLogger with the loglevel supplied
func NewConsoleLogger(logLevel string) *ConsoleLogger {
	logger := &ConsoleLogger{
		LogLevel: parseLogLevel(strings.ToLower(strings.Trim(logLevel, " "))),
		mutex:    &sync.Mutex{},
	}
	if logger.LogLevel == LOG_LEVEL_UNKNOWN {
		logger.Warn("Cannot parse log level '%s', assuming debug", logLevel)
		logger.LogLevel = LOG_LEVEL_DEBUG
	}
	return logger
}

// GetLogLevel returns the current log level
func (consoleLogger *ConsoleLogger) GetLogLevel() int { return consoleLogger.LogLevel }

// SetLogLevel sets the current log level
func (consoleLogger *ConsoleLogger) SetLogLevel(loglevel int) { consoleLogger.LogLevel = loglevel }

// Raw logs a Raw message (-----) with the specified message and Printf-style arguments.
func (consoleLogger *ConsoleLogger) Raw(message string, args ...interface{}) {
	consoleLogger.printLog("-----", message, args...)
}

// Fatal logs a FATAL message with the specified message and Printf-style arguments.
func (consoleLogger *ConsoleLogger) Fatal(message string, args ...interface{}) {
	consoleLogger.printLog("FATAL", message, args...)
}

// Error logs an ERROR message with the specified message and Printf-style arguments.
func (consoleLogger *ConsoleLogger) Error(message string, args ...interface{}) {
	consoleLogger.printLog("ERROR", message, args...)
}

// Warn logs a WARN message with the specified message and Printf-style arguments.
func (consoleLogger *ConsoleLogger) Warn(message string, args ...interface{}) {
	if consoleLogger.LogLevel > LOG_LEVEL_WARN {
		return
	}
	consoleLogger.printLog("WARN ", message, args...)
}

// Info logs an INFO message with the specified message and Printf-style arguments.
func (consoleLogger *ConsoleLogger) Info(message string, args ...interface{}) {
	if consoleLogger.LogLevel > LOG_LEVEL_INFO {
		return
	}
	consoleLogger.printLog("INFO ", message, args...)
}

// Debug logs a DEBUG message with the specified message and Printf-style arguments.
func (consoleLogger *ConsoleLogger) Debug(message string, args ...interface{}) {
	if consoleLogger.LogLevel > LOG_LEVEL_DEBUG {
		return
	}
	consoleLogger.printLog("DEBUG", message, args...)
}

func (consoleLogger *ConsoleLogger) printLog(level string, message string, args ...interface{}) {
	consoleLogger.mutex.Lock()
	defer consoleLogger.mutex.Unlock()

	fmt.Printf("%s [%s] ", GetTimeUTCString(), level)
	fmt.Printf(message, args...)
	fmt.Print("\n")
}

// GetTimeUTCString returns the current time in UTC (Zulu), in RFC3339 format.
func GetTimeUTCString() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// Parse the given log level string into a log level for use with SetLogLevel()
func parseLogLevel(logLevel string) int {
	level, found := StringLogLevel[strings.ToUpper(logLevel)]
	if !found {
		return LOG_LEVEL_UNKNOWN
	}
	return level
}
