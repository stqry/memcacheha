package log

import(
  "fmt"
  "strings"
  "sync"
  "time"
)

// This is a type to provide ConsoleLogger to the system
type ConsoleLogger struct {
	LogLevel int

  mutex *sync.Mutex
}

// The function creates a New ConsoleLogger with the loglevel supplied
func NewConsoleLogger(logLevel string) *ConsoleLogger {
  logger := &ConsoleLogger{
    LogLevel: parseLogLevel(strings.ToLower(strings.Trim(logLevel," "))),
    mutex: &sync.Mutex{},
  }
  if logger.LogLevel == LOG_LEVEL_UNKNOWN {
    logger.Warn("Cannot parse log level '%s', assuming debug", logLevel)
    logger.LogLevel = LOG_LEVEL_DEBUG
  }
  return logger
}

func (me *ConsoleLogger) GetLogLevel() int { return me.LogLevel }
func (me *ConsoleLogger) SetLogLevel(loglevel int) { me.LogLevel = loglevel }

// Logs a Raw message (-----) with the specified message and Printf-style arguments.
func (me *ConsoleLogger) Raw(message string, args... interface{}) {
  me.printLog("-----",message, args...)
}

// Logs a FATAL message with the specified message and Printf-style arguments.
func (me *ConsoleLogger) Fatal(message string, args... interface{}) {
  me.printLog("FATAL",message, args...)
}

// Logs an ERROR message with the specified message and Printf-style arguments.
func (me *ConsoleLogger) Error(message string, args... interface{}) {
  me.printLog("ERROR",message, args...)
}

// Logs a WARN message with the specified message and Printf-style arguments.
func (me *ConsoleLogger) Warn(message string, args... interface{}) {
  if me.LogLevel > LOG_LEVEL_WARN { return }
  me.printLog("WARN ",message, args...)
}

// Logs an INFO message with the specified message and Printf-style arguments.
func (me *ConsoleLogger) Info(message string, args... interface{}) {
  if me.LogLevel > LOG_LEVEL_INFO { return }
  me.printLog("INFO ",message, args...)
}

// Logs a DEBUG message with the specified message and Printf-style arguments.
func (me *ConsoleLogger) Debug(message string, args... interface{}) {
  if me.LogLevel > LOG_LEVEL_DEBUG { return }
  me.printLog("DEBUG",message, args...)
}

func (me *ConsoleLogger) printLog(level string, message string, args... interface{}) {
  me.mutex.Lock()
  defer me.mutex.Unlock()

  fmt.Printf("%s [%s] ", GetTimeUTCString(), level)
  fmt.Printf(message, args...)
  fmt.Print("\n")
}


func GetTimeUTCString() string {
  return time.Now().UTC().Format(time.RFC3339)
}

func parseLogLevel(logLevel string) int {
  level, found := StringLogLevel[strings.ToUpper(logLevel)]

  if !found { return LOG_LEVEL_UNKNOWN }
  return level
}