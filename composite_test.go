package glog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
