package glog

import (
	"io"
	"sync/atomic"
)

var defaultLogger atomic.Value

func init() {
	defaultLogger.Store(create(INFO))
}

func Default() Logger {
	return defaultLogger.Load().(Logger)
}

func setDefault(logger Logger) {
	defaultLogger.Store(logger)
}

func SetLevel(logLevel LogLevel) {
	if setter, ok := Default().(LevelSetter); ok {
		setter.SetLevel(logLevel)
		return
	}
	setDefault(create(logLevel))
}

func SetWriters(out io.Writer, err io.Writer, logLevel LogLevel) {
	setDefault(NewWithWriters(out, err, logLevel))
}

func SetOutputForLevel(logLevel LogLevel, out io.Writer) bool {
	if router, ok := Default().(LevelRouter); ok {
		router.SetOutputForLevel(logLevel, out)
		return true
	}
	return false
}

func SetOutputs(outputs map[LogLevel]io.Writer) bool {
	if router, ok := Default().(LevelRouter); ok {
		router.SetOutputs(outputs)
		return true
	}
	return false
}

func Trace(format string, a ...interface{}) {
	Default().Trace(format, a...)
}

func IsTrace() bool {
	return Default().IsTrace()
}

func Debug(format string, a ...interface{}) {
	Default().Debug(format, a...)
}

func IsDebug() bool {
	return Default().IsDebug()
}

func Info(format string, a ...interface{}) {
	Default().Info(format, a...)
}

func IsInfo() bool {
	return Default().IsInfo()
}

func Warn(format string, a ...interface{}) {
	Default().Warn(format, a...)
}

func IsWarn() bool {
	return Default().IsWarn()
}

func Error(format string, a ...interface{}) error {
	return Default().Error(format, a...)
}

func IsError() bool {
	return Default().IsError()
}

func IsEnabled(logLevel LogLevel) bool {
	return Default().IsEnabled(logLevel)
}

func Panic(format string, a ...interface{}) {
	Default().Panic(format, a...)
}

func Fatal(format string, a ...interface{}) {
	Default().Fatal(format, a...)
}

func Log(level LogLevel, a string, objs ...interface{}) {
	Default().Log(level, a, objs...)
}

func OutputLevel(level LogLevel) Output {
	return Default().GetOutput(level)
}

func ToFile(file string, level ...LogLevel) {
	log, error := createFileLogger(file, level...)
	if error == nil {
		setDefault(log)
	}
}

func ToFileAndConsole(file string, fileLevel LogLevel, consoleLevel LogLevel) {
	log, error := createFileLogger(file, fileLevel)
	if error == nil {
		_ = Error("Can't create file loger for composite logger")
	}

	setDefault(composite{
		chain: []Logger{log, create(consoleLevel)},
	})
}
