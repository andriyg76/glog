package glog

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

type ioWriterMock struct {
	string string
	log    *log.Logger
}

func (f *ioWriterMock) Write(buff []byte) (int, error) {
	f.string = string(buff)
	f.log.Print(f.string)
	return len(buff), nil
}

var so, se = ioWriterMock{"", _stdout}, ioWriterMock{"", _stderr}

func init() {
	_stdout = log.New(&so, "", log.LstdFlags)
	_stderr = log.New(&se, "", log.LstdFlags)
}

func TestFatal(t *testing.T) {
	_logger := Create(TRACE)

	var err string
	logger := _logger.(logger)
	logger.fatalf = func(format string, a ...interface{}) {
		err = fmt.Sprintf(format, a...)
	}

	logger.Fatal("Fatal %s", "fatal")

	assert.NotEqual(t, err, "")

	_logger.Info("Created fatal: %s", err)
	assert.Equal(t, "FATAL Fatal fatal", err)
}

func TestPanic(t *testing.T) {
	_logger := Create(TRACE)

	var r interface{}
	func() {
		defer func() {
			r = recover()
		}()
		_logger.Panic("Panic %s", "panic")
	}()

	assert.NotNil(t, r)
	_logger.Info("Recover from panic: %s", r)
	assert.Equal(t, "PANIC Panic panic", r)
}

func TestLevels(t *testing.T) {
	str := "trace"
	checkLog(t, TRACE, str, TRACE, str, "")
	checkLog(t, TRACE, str, DEBUG, str, "")
	checkLog(t, TRACE, str, WARN, "", str)

	str = "debug"
	checkLog(t, DEBUG, str, TRACE, "", "")
	checkLog(t, DEBUG, str, DEBUG, str, "")
	checkLog(t, DEBUG, str, WARN, "", str)

	str = "info"
	checkLog(t, INFO, str, DEBUG, "", "")
	checkLog(t, INFO, str, INFO, str, "")
	checkLog(t, INFO, str, WARN, "", str)

	str = "warn"
	checkLog(t, WARN, str, INFO, "", "")
	checkLog(t, WARN, str, WARN, "", str)
	checkLog(t, WARN, str, ERROR, "", str)

	str = "error"
	checkLog(t, ERROR, str, INFO, "", "")
	checkLog(t, ERROR, str, WARN, "", "")
	checkLog(t, ERROR, str, ERROR, "", str)
}

func TestSetLevel(t *testing.T) {
	logger := Create(INFO)
	setter, ok := logger.(LevelSetter)
	assert.True(t, ok)

	assert.False(t, logger.IsDebug())
	setter.SetLevel(DEBUG)
	assert.True(t, logger.IsDebug())
}

func TestOutputRoutingByLevel(t *testing.T) {
	var infoOut bytes.Buffer
	var debugOut bytes.Buffer

	logger := NewWithWriters(&infoOut, &infoOut, DEBUG)
	router, ok := logger.(LevelRouter)
	assert.True(t, ok)

	router.SetOutputForLevel(DEBUG, &debugOut)
	router.SetOutputForLevel(TRACE, &debugOut)

	logger.Debug("debug message")
	logger.Info("info message")

	assert.Contains(t, debugOut.String(), "DEBUG debug message")
	assert.Contains(t, infoOut.String(), "INFO info message")
	assert.NotContains(t, infoOut.String(), "DEBUG debug message")
}

func checkLog(t *testing.T, ll LogLevel, str string, pl LogLevel, std_out string, stderr_out string) {
	_logger := Create(ll)

	se.string = ""
	so.string = ""

	_logger.Log(pl, str)
	assertOut(t, pl, so.string, std_out, "stdout")
	assertOut(t, pl, se.string, stderr_out, "stderr")

	if ll.weight <= pl.weight {
		assert.Truef(t, _logger.IsEnabled(pl), "%s should be enabled for %s", pl, ll)
	}
	if ll.weight > pl.weight {
		assert.Truef(t, !_logger.IsEnabled(pl), "%s should not be enabled for %s", pl, ll)
	}
}

func assertOut(t *testing.T, pl LogLevel, actual, expected, stream string) {
	if actual == "" && expected == "" {
		return
	}

	if actual != "" && expected != "" {
		return
	}

	assert.Fail(t, "", "Unexpected %s output expected %s: [%#v] actual: [%#v]", stream, pl, expected, actual)
}
