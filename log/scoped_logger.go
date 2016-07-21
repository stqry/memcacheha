package log

// ScopedLogger wraps an existing logger, prefixing all logs with a string scope
type ScopedLogger struct {
	Logger Logger
	Scope  string
}

// Return a new ScopedLogger with the given scope and base logger
func NewScopedLogger(name string, logger Logger) *ScopedLogger {
	return &ScopedLogger{
		Scope:  name,
		Logger: logger,
	}
}

// Get the current logging level
func (me *ScopedLogger) GetLogLevel() int { return me.Logger.GetLogLevel() }

// Set the current logging level
func (me *ScopedLogger) SetLogLevel(loglevel int) { me.Logger.SetLogLevel(loglevel) }

// Logs a Raw message (-----) with the specified message and Printf-style arguments.
func (me *ScopedLogger) Raw(message string, args ...interface{}) {
	me.Logger.Raw(me.Scope+": "+message, args...)
}

// Logs a FATAL message with the specified message and Printf-style arguments.
func (me *ScopedLogger) Fatal(message string, args ...interface{}) {
	me.Logger.Fatal(me.Scope+": "+message, args...)
}

// Logs an ERROR message with the specified message and Printf-style arguments.
func (me *ScopedLogger) Error(message string, args ...interface{}) {
	me.Logger.Error(me.Scope+": "+message, args...)
}

// Logs a WARN message with the specified message and Printf-style arguments.
func (me *ScopedLogger) Warn(message string, args ...interface{}) {
	me.Logger.Warn(me.Scope+": "+message, args...)
}

// Logs an INFO message with the specified message and Printf-style arguments.
func (me *ScopedLogger) Info(message string, args ...interface{}) {
	me.Logger.Info(me.Scope+": "+message, args...)
}

// Logs a DEBUG message with the specified message and Printf-style arguments.
func (me *ScopedLogger) Debug(message string, args ...interface{}) {
	me.Logger.Debug(me.Scope+": "+message, args...)
}
