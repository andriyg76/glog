package glog

import (
	"io"
	"sync/atomic"
)

type defaultHolder struct{ Logger }

var defaultLogger atomic.Value

func init() {
	defaultLogger.Store(defaultHolder{Logger: create(INFO)})
}

// Default returns the global default logger (initial level INFO, stdout/stderr).
func Default() Logger {
	return defaultLogger.Load().(defaultHolder).Logger
}

func setDefault(logger Logger) {
	defaultLogger.Store(defaultHolder{Logger: logger})
}

// SetLevel sets the minimum level of the default logger. If it does not implement LevelSetter, replaces the default with a new logger.
func SetLevel(logLevel LogLevel) {
	if setter, ok := Default().(LevelSetter); ok {
		setter.SetLevel(logLevel)
		return
	}
	setDefault(create(logLevel))
}

// SetWriters replaces the default logger with one that writes to out (info and below) and err (warn and above).
func SetWriters(out io.Writer, err io.Writer, logLevel LogLevel) {
	setDefault(NewWithWriters(out, err, logLevel))
}

// SetOutputForLevel sets a dedicated output for the given level on the default logger. Returns true only if the default is a LevelRouter.
func SetOutputForLevel(logLevel LogLevel, out io.Writer) bool {
	if router, ok := Default().(LevelRouter); ok {
		router.SetOutputForLevel(logLevel, out)
		return true
	}
	return false
}

// SetOutputs sets per-level outputs on the default logger. Returns true only if the default is a LevelRouter.
func SetOutputs(outputs map[LogLevel]io.Writer) bool {
	if router, ok := Default().(LevelRouter); ok {
		router.SetOutputs(outputs)
		return true
	}
	return false
}

// Trace logs at TRACE level using the default logger.
func Trace(format string, a ...interface{}) {
	Default().Trace(format, a...)
}

// IsTrace reports whether the default logger is enabled for TRACE.
func IsTrace() bool {
	return Default().IsTrace()
}

// Debug logs at DEBUG level using the default logger.
func Debug(format string, a ...interface{}) {
	Default().Debug(format, a...)
}

// IsDebug reports whether the default logger is enabled for DEBUG.
func IsDebug() bool {
	return Default().IsDebug()
}

// Info logs at INFO level using the default logger.
func Info(format string, a ...interface{}) {
	Default().Info(format, a...)
}

// IsInfo reports whether the default logger is enabled for INFO.
func IsInfo() bool {
	return Default().IsInfo()
}

// Warn logs at WARN level using the default logger.
func Warn(format string, a ...interface{}) {
	Default().Warn(format, a...)
}

// IsWarn reports whether the default logger is enabled for WARN.
func IsWarn() bool {
	return Default().IsWarn()
}

// Error logs at ERROR level using the default logger and returns an error for chaining.
func Error(format string, a ...interface{}) error {
	return Default().Error(format, a...)
}

// IsError reports whether the default logger is enabled for ERROR.
func IsError() bool {
	return Default().IsError()
}

// IsEnabled reports whether the default logger is enabled for the given level.
func IsEnabled(logLevel LogLevel) bool {
	return Default().IsEnabled(logLevel)
}

// Panic logs at PANIC level and then panics.
func Panic(format string, a ...interface{}) {
	Default().Panic(format, a...)
}

// Fatal logs at FATAL level and then exits the program (os.Exit(1)).
func Fatal(format string, a ...interface{}) {
	Default().Fatal(format, a...)
}

// Log writes to the default logger at the given level with the format and args.
func Log(level LogLevel, a string, objs ...interface{}) {
	Default().Log(level, a, objs...)
}

// OutputLevel returns an Output that writes at the given level to the default logger.
func OutputLevel(level LogLevel) Output {
	return Default().GetOutput(level)
}

// ToFile switches the default logger to append to the given file. Level defaults to INFO. On open failure, the default logger is unchanged.
func ToFile(file string, level ...LogLevel) {
	log, error := createFileLogger(file, level...)
	if error == nil {
		setDefault(log)
	}
}

// ToFileAndConsole sets the default logger to a composite: file (at fileLevel) and console (at consoleLevel). On file open failure, logs the error and leaves the default unchanged.
func ToFileAndConsole(file string, fileLevel LogLevel, consoleLevel LogLevel) {
	log, err := createFileLogger(file, fileLevel)
	if err != nil {
		_ = Error("Can't create file logger for composite logger: %v", err)
		return
	}
	setDefault(composite{
		chain: []Logger{log, create(consoleLevel)},
	})
}
