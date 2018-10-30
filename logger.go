package memcacheha

// Logger defines what is expected of the passed in logger
type Logger interface {
	Error(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Info(message string, args ...interface{})
	Debug(message string, args ...interface{})
}

// scopedLogger wraps an existing logger, prefixing all logs with a string scope
type scopedLogger struct {
	l Logger
	prefix string
}

// newScopedLogger returns a new ScopedLogger with the given scope and base logger
func newScopedLogger(prefix string, logger Logger) *scopedLogger {
	return &scopedLogger{
		prefix: prefix,
		l: logger,
	}
}

// Error logs an ERROR message with the specified message and Printf-style arguments.
func (sl *scopedLogger) Error(message string, args ...interface{}) {
	sl.l.Error(sl.prefix+": "+message, args...)
}

// Warn logs a WARN message with the specified message and Printf-style arguments.
func (sl *scopedLogger) Warn(message string, args ...interface{}) {
	sl.l.Warn(sl.prefix+": "+message, args...)
}

// Info logs an INFO message with the specified message and Printf-style arguments.
func (sl *scopedLogger) Info(message string, args ...interface{}) {
	sl.l.Info(sl.prefix+": "+message, args...)
}

// Debug logs a DEBUG message with the specified message and Printf-style arguments.
func (sl *scopedLogger) Debug(message string, args ...interface{}) {
	sl.l.Debug(sl.prefix+": "+message, args...)
}
