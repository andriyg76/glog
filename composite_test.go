package glog

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposite_ForwardsToAllLoggers(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	log1 := NewWithWriters(&buf1, &buf1, INFO)
	log2 := NewWithWriters(&buf2, &buf2, INFO)
	comp := Composite(log1, log2)

	comp.Info("forwarded")

	assert.Contains(t, buf1.String(), "forwarded")
	assert.Contains(t, buf2.String(), "forwarded")
}

func TestDefaultComposite_DefaultWritesToAll(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	log1 := NewWithWriters(&buf1, &buf1, INFO)
	log2 := NewWithWriters(&buf2, &buf2, INFO)
	defer SetWriters(os.Stdout, os.Stderr, INFO)

	DefaultComposite(log1, log2)
	Info("multi")

	assert.Contains(t, buf1.String(), "multi")
	assert.Contains(t, buf2.String(), "multi")
}

func TestComposite_IsEnabled(t *testing.T) {
	debug := create(DEBUG)
	warn := create(WARN)
	fatal := create(FATAL)
	log := Composite(warn, debug, fatal)

	assert.True(t, log.IsEnabled(FATAL))

	assert.True(t, log.IsEnabled(INFO))
	assert.True(t, log.IsInfo())

	assert.True(t, log.IsEnabled(DEBUG))
	assert.True(t, log.IsDebug())

	assert.False(t, log.IsEnabled(TRACE))
	assert.False(t, log.IsTrace())

	assert.True(t, debug.IsEnabled(DEBUG))
	assert.False(t, warn.IsEnabled(DEBUG))
	assert.False(t, fatal.IsEnabled(DEBUG))

	assert.False(t, debug.IsEnabled(TRACE))
	assert.False(t, warn.IsEnabled(TRACE))
	assert.False(t, fatal.IsEnabled(TRACE))

	assert.True(t, debug.IsEnabled(FATAL))
	assert.True(t, warn.IsEnabled(FATAL))
	assert.True(t, fatal.IsEnabled(FATAL))
}
