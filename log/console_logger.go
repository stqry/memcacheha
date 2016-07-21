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
func (me *ConsoleLogger) GetLogLevel() int { return me.LogLevel }

// SetLogLevel sets the current log level
func (me *ConsoleLogger) SetLogLevel(loglevel int) { me.LogLevel = loglevel }

// Raw logs a Raw message (-----) with the specified message and Printf-style arguments.
func (me *ConsoleLogger) Raw(message string, args ...interface{}) {
	me.printLog("-----", message, args...)
}

// Fatal logs a FATAL message with the specified message and Printf-style arguments.
func (me *ConsoleLogger) Fatal(message string, args ...interface{}) {
	me.printLog("FATAL", message, args...)
}

// Error logs an ERROR message with the specified message and Printf-style arguments.
func (me *ConsoleLogger) Error(message string, args ...interface{}) {
	me.printLog("ERROR", message, args...)
}

// Warn logs a WARN message with the specified message and Printf-style arguments.
func (me *ConsoleLogger) Warn(message string, args ...interface{}) {
	if me.LogLevel > LOG_LEVEL_WARN {
		return
	}
	me.printLog("WARN ", message, args...)
}

// Info logs an INFO message with the specified message and Printf-style arguments.
func (me *ConsoleLogger) Info(message string, args ...interface{}) {
	if me.LogLevel > LOG_LEVEL_INFO {
		return
	}
	me.printLog("INFO ", message, args...)
}

// Debug logs a DEBUG message with the specified message and Printf-style arguments.
func (me *ConsoleLogger) Debug(message string, args ...interface{}) {
	if me.LogLevel > LOG_LEVEL_DEBUG {
		return
	}
	me.printLog("DEBUG", message, args...)
}

func (me *ConsoleLogger) printLog(level string, message string, args ...interface{}) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

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
