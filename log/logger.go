package log

// Individual log levels
const (
  LOG_LEVEL_UNKNOWN = iota
  LOG_LEVEL_DEBUG = iota
  LOG_LEVEL_INFO = iota
  LOG_LEVEL_WARN = iota
  LOG_LEVEL_ERROR = iota
  LOG_LEVEL_FATAL = iota
)

// Map of uppercase log level name strings to int log levels 
var StringLogLevel map[string]int = map[string]int{
  "DEBUG": LOG_LEVEL_DEBUG,
  "INFO": LOG_LEVEL_INFO,
  "WARN": LOG_LEVEL_WARN,
  "ERROR": LOG_LEVEL_ERROR,
  "FATAL": LOG_LEVEL_FATAL,
}

// Map of int log levels to string log level names
var LogLevelString map[int]string = map[int]string{
  LOG_LEVEL_DEBUG: "DEBUG",
  LOG_LEVEL_INFO: "INFO",
  LOG_LEVEL_WARN: "WARN",
  LOG_LEVEL_ERROR: "ERROR",
  LOG_LEVEL_FATAL: "FATAL",
}

type Logger interface {
  GetLogLevel() int
  SetLogLevel(loglevel int)

  Raw(message string, args... interface{})
  Fatal(message string, args... interface{})
  Error(message string, args... interface{})
  Warn(message string, args... interface{})
  Info(message string, args... interface{})
  Debug(message string, args... interface{})
}