package glog

// LogLevel represents a log level. Use the package constants (TRACE, DEBUG, INFO, etc.).
// Order by severity: TRACE < DEBUG < INFO < WARN < ERROR, PANIC, FATAL.
type LogLevel struct {
	prefix string
	weight int
}

// TRACE is the lowest level; typically disabled in production.
var TRACE = LogLevel{
	prefix: "TRACE",
	weight: -2,
}

// DEBUG is for verbose development output.
var DEBUG = LogLevel{
	prefix: "DEBUG",
	weight: -1,
}

// INFO is the default level; general operational messages.
var INFO = LogLevel{
	prefix: " INFO",
	weight: 0,
}

// WARN is for recoverable or unexpected conditions.
var WARN = LogLevel{
	prefix: " WARN",
	weight: 1,
}

// ERROR is for errors; Error() also returns an error for chaining.
var ERROR = LogLevel{
	prefix: "ERROR",
	weight: 2,
}

// PANIC logs and then panics.
var PANIC = LogLevel{
	prefix: "PANIC",
	weight: 2,
}

// FATAL logs and then exits (os.Exit(1)).
var FATAL = LogLevel{
	prefix: "FATAL",
	weight: 2,
}

// String returns the level prefix (e.g. "DEBUG", " INFO").
func (l LogLevel) String() string {
	return l.prefix
}
