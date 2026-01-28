package glog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogLevel_String(t *testing.T) {
	assert.Equal(t, "TRACE", TRACE.String())
	assert.Equal(t, "DEBUG", DEBUG.String())
	assert.Equal(t, " INFO", INFO.String())
	assert.Equal(t, " WARN", WARN.String())
	assert.Equal(t, "ERROR", ERROR.String())
	assert.Equal(t, "PANIC", PANIC.String())
	assert.Equal(t, "FATAL", FATAL.String())
}
