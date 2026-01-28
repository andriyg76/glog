# glogger

[![Go](https://github.com/andriyg76/glog/actions/workflows/go.yml/badge.svg)](https://github.com/andriyg76/glog/actions/workflows/go.yml)
git lo[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Leveled logging library for Go with configurable outputs, per-level routing, and composite loggers. Default minimum level is **INFO**. Level order (low to high): **TRACE** &lt; **DEBUG** &lt; **INFO** &lt; **WARN** &lt; **ERROR**, **PANIC**, **FATAL**.

## Installation

```bash
go get github.com/andriyg76/glog
```

## Usage

### Default logger (package-level)

Use the package-level functions; they write to the global default logger (stdout/stderr, level INFO).

```go
import "github.com/andriyg76/glog"

glog.Info("started")
glog.Debug("detail: %s", x)
glog.SetLevel(glog.DEBUG)
glog.Debug("now visible")
if glog.IsEnabled(glog.TRACE) {
    glog.Trace("trace message")
}
```

### Create your own logger

```go
logger := glog.Create(glog.DEBUG)
logger.Info("hello")
logger.Debug("detail")

// Custom writers (e.g. buffers, files)
buf := &bytes.Buffer{}
log := glog.NewWithWriters(buf, buf, glog.INFO)
log.Info("to buffer")
```

### Per-level output (LevelRouter)

Route different levels to different writers (e.g. debug to file, info to stdout).

```go
outputs := map[glog.LogLevel]io.Writer{
    glog.DEBUG: debugFile,
    glog.INFO:  os.Stdout,
    glog.WARN:  os.Stderr,
}
router := glog.NewLevelRouter(outputs, glog.DEBUG)
router.Info("to stdout")
router.Debug("to debugFile")

// Or set one level at a time
log := glog.NewWithWriters(os.Stdout, os.Stderr, glog.DEBUG)
if r, ok := log.(glog.LevelRouter); ok {
    r.SetOutputForLevel(glog.DEBUG, debugFile)
}
```

### Log to file

```go
// Default logger writes only to file (level defaults to INFO)
glog.ToFile("/var/log/app.log")
glog.ToFile("/var/log/app.log", glog.DEBUG)

// File + console (different levels)
glog.ToFileAndConsole("/var/log/app.log", glog.DEBUG, glog.INFO)
```

**Note:** `ToFile` and `ToFileAndConsole` do not return an error. If the file cannot be opened, the default logger is left unchanged (and `ToFileAndConsole` logs the error).

### Composite logger

Forward every log to multiple loggers (e.g. file and console).

```go
fileLog, _ := glog.Create(glog.DEBUG) // assume you set file output
consoleLog := glog.Create(glog.INFO)
composite := glog.Composite(fileLog, consoleLog)
composite.Info("goes to both")

// Set as default
glog.DefaultComposite(fileLog, consoleLog)
```

## API Reference

| Symbol | Description |
|--------|-------------|
| **Types** | |
| `Output` | Interface: `Printf(format, a...)`; used for level-specific writers. |
| `TraceLogger` | Trace, TraceLogger(), IsTrace(). |
| `DebugLogger` | Debug, IsDebug(), DebugLogger(). |
| `InfoLogger` | Info, IsInfo(). |
| `WarnLogger` | Warn, IsWarn(). |
| `ErrorLogger` | Error (returns error), IsError(). |
| `Logger` | Full interface: all level methods, Log, IsEnabled, GetOutput, Panic, Fatal. |
| `LevelSetter` | SetLevel(LogLevel). |
| `LevelRouter` | Logger + SetOutputForLevel, SetOutputs. |
| `LogLevel` | Level value; use constants TRACE, DEBUG, INFO, WARN, ERROR, PANIC, FATAL. |
| **Constructors** | |
| `Create(LogLevel)` | New Logger (stdout/stderr). |
| `NewWithWriters(out, err, LogLevel)` | Logger with custom writers. |
| `NewLevelRouter(outputs, level?)` | LevelRouter with optional per-level outputs. |
| **Default logger** | |
| `Default()` | Returns the global logger. |
| `SetLevel(LogLevel)` | Set default minimum level. |
| `SetWriters(out, err, LogLevel)` | Replace default with custom writers. |
| `SetOutputForLevel(level, out)` | Set output for one level (returns true if default is LevelRouter). |
| `SetOutputs(outputs)` | Set per-level outputs (returns true if default is LevelRouter). |
| `Trace/Debug/Info/Warn/Error(format, a...)` | Log at level; Error returns error. |
| `IsTrace/IsDebug/IsInfo/IsWarn/IsError()` | Report if level enabled. |
| `IsEnabled(LogLevel)` | Report if level enabled. |
| `Log(level, format, objs...)` | Log at given level. |
| `OutputLevel(level)` | Output that writes at that level. |
| `Panic/Fatal(format, a...)` | Log and panic / exit. |
| `ToFile(file, level?)` | Default logger appends to file; on failure default unchanged. |
| `ToFileAndConsole(file, fileLevel, consoleLevel)` | Default = file + console; on file failure default unchanged. |
| **Composite** | |
| `Composite(main, loggers...)` | Logger that forwards to main then each logger. |
| `DefaultComposite(main, loggers...)` | Set default to Composite(main, loggers...). |

## License

MIT License. Copyright (c) 2026 andriyg76. See [LICENSE](LICENSE) for full text.
