package log

// ScopedLogger wraps an existing logger, prefixing all logs with a string scope
type ScopedLogger struct {
	Logger Logger
	Scope  string
}

// NewScopedLogger returns a new ScopedLogger with the given scope and base logger
func NewScopedLogger(name string, logger Logger) *ScopedLogger {
	return &ScopedLogger{
		Scope:  name,
		Logger: logger,
	}
}

// GetLogLevel returns the current logging level
func (scopedLogger *ScopedLogger) GetLogLevel() int { return scopedLogger.Logger.GetLogLevel() }

// SetLogLevel sets the current logging level
func (scopedLogger *ScopedLogger) SetLogLevel(loglevel int) { scopedLogger.Logger.SetLogLevel(loglevel) }

// Raw logs a Raw message (-----) with the specified message and Printf-style arguments.
func (scopedLogger *ScopedLogger) Raw(message string, args ...interface{}) {
	scopedLogger.Logger.Raw(scopedLogger.Scope+": "+message, args...)
}

// Fatal logs a FATAL message with the specified message and Printf-style arguments.
func (scopedLogger *ScopedLogger) Fatal(message string, args ...interface{}) {
	scopedLogger.Logger.Fatal(scopedLogger.Scope+": "+message, args...)
}

// Error logs an ERROR message with the specified message and Printf-style arguments.
func (scopedLogger *ScopedLogger) Error(message string, args ...interface{}) {
	scopedLogger.Logger.Error(scopedLogger.Scope+": "+message, args...)
}

// Warn logs a WARN message with the specified message and Printf-style arguments.
func (scopedLogger *ScopedLogger) Warn(message string, args ...interface{}) {
	scopedLogger.Logger.Warn(scopedLogger.Scope+": "+message, args...)
}

// Info logs an INFO message with the specified message and Printf-style arguments.
func (scopedLogger *ScopedLogger) Info(message string, args ...interface{}) {
	scopedLogger.Logger.Info(scopedLogger.Scope+": "+message, args...)
}

// Debug logs a DEBUG message with the specified message and Printf-style arguments.
func (scopedLogger *ScopedLogger) Debug(message string, args ...interface{}) {
	scopedLogger.Logger.Debug(scopedLogger.Scope+": "+message, args...)
}
