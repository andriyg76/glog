package glog

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault_NonNull(t *testing.T) {
	assert.NotNil(t, Default())
}

func TestSetLevel_ChangesDefault(t *testing.T) {
	defer SetWriters(os.Stdout, os.Stderr, INFO)

	SetLevel(DEBUG)
	assert.True(t, IsDebug())
	assert.True(t, IsInfo())

	SetLevel(INFO)
	assert.False(t, IsDebug())
	assert.True(t, IsInfo())
}

func TestSetWriters_LogsToBuffers(t *testing.T) {
	var out, err bytes.Buffer
	defer SetWriters(os.Stdout, os.Stderr, INFO)

	SetWriters(&out, &err, INFO)
	Info("hello")
	Warn("warn")

	assert.Contains(t, out.String(), "hello")
	assert.Contains(t, err.String(), "warn")
}

func TestSetOutputForLevel_ReturnsTrueAndRoutes(t *testing.T) {
	var out, err, debugOut bytes.Buffer
	defer SetWriters(os.Stdout, os.Stderr, INFO)

	SetWriters(&out, &err, DEBUG)
	ok := SetOutputForLevel(DEBUG, &debugOut)
	assert.True(t, ok)

	Debug("debug msg")
	Info("info msg")

	assert.Contains(t, debugOut.String(), "DEBUG")
	assert.Contains(t, debugOut.String(), "debug msg")
	assert.Contains(t, out.String(), "info msg")
	assert.NotContains(t, out.String(), "debug msg")
}

func TestSetOutputs_RoutesByLevel(t *testing.T) {
	var infoBuf, errBuf bytes.Buffer
	defer SetWriters(os.Stdout, os.Stderr, INFO)

	wm := map[LogLevel]io.Writer{
		INFO:  &infoBuf,
		WARN:  &errBuf,
		ERROR: &errBuf,
	}
	SetWriters(&infoBuf, &errBuf, DEBUG)
	ok := SetOutputs(wm)
	assert.True(t, ok)

	Info("info only")
	Warn("warn only")

	assert.Contains(t, infoBuf.String(), "info only")
	assert.Contains(t, errBuf.String(), "warn only")
}

func TestToFile_WritesToFile(t *testing.T) {
	fpath := filepath.Join(os.TempDir(), "glog_test_tofile.txt")
	defer func() {
		SetWriters(os.Stdout, os.Stderr, INFO)
		_ = os.Remove(fpath)
	}()

	ToFile(fpath, INFO)
	Info("file message")

	data, err := os.ReadFile(fpath)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "file message")
}

func TestToFile_InvalidPath_LeavesDefaultUnchanged(t *testing.T) {
	var out, err bytes.Buffer
	SetWriters(&out, &err, INFO)
	defer SetWriters(os.Stdout, os.Stderr, INFO)

	ToFile("/nonexistent/path/xyz/123", INFO)
	Info("still default")

	assert.Contains(t, out.String(), "still default")
}

func TestToFileAndConsole_Success(t *testing.T) {
	fpath := filepath.Join(os.TempDir(), "glog_test_combo.txt")
	defer func() {
		SetWriters(os.Stdout, os.Stderr, INFO)
		_ = os.Remove(fpath)
	}()

	ToFileAndConsole(fpath, INFO, INFO)
	Info("combo message")

	data, ferr := os.ReadFile(fpath)
	assert.NoError(t, ferr)
	assert.Contains(t, string(data), "combo message")
}

func TestToFileAndConsole_InvalidPath_LeavesDefaultUnchanged(t *testing.T) {
	var out, err bytes.Buffer
	SetWriters(&out, &err, INFO)
	defer SetWriters(os.Stdout, os.Stderr, INFO)

	ToFileAndConsole("/nonexistent/path/xyz", INFO, INFO)
	Info("unchanged default")

	assert.Contains(t, out.String(), "unchanged default")
}
