// Package glog provides a leveled logging library with configurable outputs,
// level routing, and composite loggers. Default log level is INFO.
// Level order (low to high): TRACE < DEBUG < INFO < WARN < ERROR, PANIC, FATAL.
//
// Use the package-level functions (Info, Debug, Trace, etc.) with the default
// logger, or create your own logger with Create, NewWithWriters, or NewLevelRouter.
package glog
