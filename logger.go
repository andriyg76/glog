package glog

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
)

type Output interface {
	Printf(format string, a ...interface{})
}

type TraceLogger interface {
	Trace(format string, a ...interface{})
	TraceLogger() Output
	IsTrace() bool
}

type DebugLogger interface {
	Debug(format string, a ...interface{})
	IsDebug() bool
	DebugLogger() Output
}

type InfoLogger interface {
	Info(format string, a ...interface{})
	IsInfo() bool
}

type WarnLogger interface {
	Warn(format string, a ...interface{})
	IsWarn() bool
}

type ErrorLogger interface {
	Error(format string, a ...interface{}) error
	IsError() bool
}

type Logger interface {
	DebugLogger
	TraceLogger
	WarnLogger
	InfoLogger
	ErrorLogger

	// Log formats according to a format specifier
	Log(LogLevel LogLevel, format string, a ...interface{})
	IsEnabled(logLevel LogLevel) bool
	GetOutput(LogLevel LogLevel) Output

	Panic(format string, a ...interface{})
	Fatal(format string, a ...interface{})
}

type LevelSetter interface {
	SetLevel(logLevel LogLevel)
}

type LevelRouter interface {
	Logger
	SetOutputForLevel(logLevel LogLevel, out io.Writer)
	SetOutputs(outputs map[LogLevel]io.Writer)
}

type logger struct {
	level  *int32
	out    *log.Logger
	err    *log.Logger
	fatalf func(format string, a ...interface{})
	router *outputRouter
}

type outputRouter struct {
	mu      sync.RWMutex
	outputs map[LogLevel]Output
}

func newOutputRouter() *outputRouter {
	return &outputRouter{
		outputs: make(map[LogLevel]Output),
	}
}

func (r *outputRouter) SetOutput(level LogLevel, out Output) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if out == nil {
		delete(r.outputs, level)
		return
	}

	r.outputs[level] = out
}

func (r *outputRouter) SetOutputs(outputs map[LogLevel]Output) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.outputs = make(map[LogLevel]Output, len(outputs))
	for level, out := range outputs {
		if out != nil {
			r.outputs[level] = out
		}
	}
}

func (r *outputRouter) OutputFor(level LogLevel) (Output, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out, ok := r.outputs[level]
	return out, ok
}

type discardWriter struct{}

func (discardWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

var discardWriterInstance io.Writer = discardWriter{}

func newLevelPointer(logLevel LogLevel) *int32 {
	level := int32(logLevel.weight)
	return &level
}

func newStdLogger(writer io.Writer) *log.Logger {
	if writer == nil {
		writer = discardWriterInstance
	}
	return log.New(writer, "", log.LstdFlags)
}

func outputFromWriter(writer io.Writer) Output {
	if writer == nil {
		return nil
	}
	return newStdLogger(writer)
}

func outputsFromWriters(outputs map[LogLevel]io.Writer) map[LogLevel]Output {
	if len(outputs) == 0 {
		return map[LogLevel]Output{}
	}
	converted := make(map[LogLevel]Output, len(outputs))
	for level, writer := range outputs {
		if writer != nil {
			converted[level] = newStdLogger(writer)
		}
	}
	return converted
}
func (l logger) IsDebug() bool {
	return l.IsEnabled(DEBUG)
}

func (l logger) IsTrace() bool {
	return l.IsEnabled(TRACE)
}

func (l logger) IsWarn() bool {
	return l.IsEnabled(WARN)
}

func (l logger) IsInfo() bool {
	return l.IsEnabled(INFO)
}

func (l logger) IsError() bool {
	return l.IsEnabled(ERROR)
}

func (l logger) IsEnabled(logLevel LogLevel) bool {
	return int32(logLevel.weight) >= l.currentLevelWeight()
}

func (l logger) currentLevelWeight() int32 {
	if l.level == nil {
		return int32(INFO.weight)
	}
	return atomic.LoadInt32(l.level)
}

func create(logLevel LogLevel) logger {
	return logger{
		level:  newLevelPointer(logLevel),
		err:    _stderr,
		out:    _stdout,
		fatalf: _stderr.Fatalf,
		router: newOutputRouter(),
	}
}

func createWithWriters(out io.Writer, err io.Writer, logLevel LogLevel) logger {
	outLogger := newStdLogger(out)
	errLogger := newStdLogger(err)
	return logger{
		level:  newLevelPointer(logLevel),
		err:    errLogger,
		out:    outLogger,
		fatalf: errLogger.Fatalf,
		router: newOutputRouter(),
	}
}

func createFileLogger(file string, level ...LogLevel) (logger, error) {
	logLevel := INFO
	if len(level) > 0 {
		logLevel = level[0]
	}

	instance := create(logLevel)
	openFile, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return instance, Error("Error creating file %s output: %s", file, err)
	}
	w := newStdLogger(openFile)
	instance.err = w
	instance.out = w
	instance.fatalf = w.Fatalf
	return instance, nil
}

func Create(logLevel LogLevel) Logger {
	return create(logLevel)
}

func NewWithWriters(out io.Writer, err io.Writer, logLevel LogLevel) Logger {
	return createWithWriters(out, err, logLevel)
}

func NewLevelRouter(outputs map[LogLevel]io.Writer, level ...LogLevel) LevelRouter {
	logLevel := INFO
	if len(level) > 0 {
		logLevel = level[0]
	}
	instance := create(logLevel)
	instance.SetOutputs(outputs)
	return instance
}

var _stdout = log.New(os.Stdout, "", log.LstdFlags)
var _stderr = log.New(os.Stderr, "", log.LstdFlags)

type dumbLogger struct{}

func (d dumbLogger) Printf(format string, objs ...interface{}) {}

var dumbLoggerInstance = dumbLogger{}

type loggerWithLevel struct {
	logLevel *LogLevel
	logger   *logger
}

func (l loggerWithLevel) Printf(format string, objs ...interface{}) {
	l.logger.Log(*l.logLevel, format, objs...)
}

func (l logger) GetOutput(logLevel LogLevel) Output {
	return loggerWithLevel{
		logLevel: &logLevel,
		logger:   &l,
	}
}

func (l logger) TraceLogger() Output {
	return l.GetOutput(TRACE)
}

func (l logger) SetLevel(logLevel LogLevel) {
	if l.level == nil {
		return
	}
	atomic.StoreInt32(l.level, int32(logLevel.weight))
}

func (l logger) SetOutputForLevel(logLevel LogLevel, out io.Writer) {
	if l.router == nil {
		return
	}
	l.router.SetOutput(logLevel, outputFromWriter(out))
}

func (l logger) SetOutputs(outputs map[LogLevel]io.Writer) {
	if l.router == nil {
		return
	}
	l.router.SetOutputs(outputsFromWriters(outputs))
}

func (l logger) Log(logLevel LogLevel, format string, objs ...interface{}) {
	logFormat := logLevel.prefix + " " + format

	if logLevel == PANIC {
		if out, ok := l.outputForLevel(logLevel); ok {
			if panicOut, ok := out.(interface {
				Panicf(format string, a ...interface{})
			}); ok {
				panicOut.Panicf(logFormat, objs...)
				return
			}
		}
		l.err.Panicf(logFormat, objs...)
		return
	}
	if logLevel == FATAL {
		if out, ok := l.outputForLevel(logLevel); ok {
			if fatalOut, ok := out.(interface {
				Fatalf(format string, a ...interface{})
			}); ok {
				fatalOut.Fatalf(logFormat, objs...)
				return
			}
		}
		l.fatalf(logFormat, objs...)
		return
	}

	if !l.IsEnabled(logLevel) {
		return
	}

	if out, ok := l.outputForLevel(logLevel); ok {
		out.Printf(logFormat, objs...)
		return
	}

	if logLevel.weight >= WARN.weight {
		l.err.Printf(logFormat, objs...)
		return
	}

	l.out.Printf(logFormat, objs...)
}

func (l logger) outputForLevel(logLevel LogLevel) (Output, bool) {
	if l.router == nil {
		return nil, false
	}
	return l.router.OutputFor(logLevel)
}

func (l logger) Debug(format string, objs ...interface{}) {
	l.Log(DEBUG, format, objs...)
}

func (l logger) DebugLogger() Output {
	return l.GetOutput(DEBUG)
}

func (l logger) Trace(format string, objs ...interface{}) {
	l.Log(TRACE, format, objs...)
}

func (l logger) Info(format string, objs ...interface{}) {
	l.Log(INFO, format, objs...)
}

func (l logger) Warn(format string, objs ...interface{}) {
	l.Log(WARN, format, objs...)
}

func (l logger) Error(format string, objs ...interface{}) error {
	err := fmt.Errorf(format, objs...)
	l.Log(ERROR, "%s", err)
	return err
}

func (l logger) Panic(format string, objs ...interface{}) {
	l.Log(PANIC, format, objs...)
}

func (l logger) Fatal(format string, objs ...interface{}) {
	l.Log(FATAL, format, objs...)
}
