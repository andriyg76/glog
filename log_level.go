package glog

type LogLevel struct {
	prefix string
	weight int
}

var TRACE = LogLevel{
	prefix: "[trace]",
	weight: -2,
}

var DEBUG = LogLevel{
	prefix: "[debug]",
	weight: -1,
}

var INFO = LogLevel{
	prefix: "[info ]",
	weight: 0,
}

var WARN = LogLevel{
	prefix: "[warn ]",
	weight: 1,
}

var ERROR = LogLevel{
	prefix: "[error]",
	weight: 2,
}

var PANIC = LogLevel{
	prefix: "[trace]",
	weight: 2,
}

var FATAL = LogLevel{
	prefix: "[fatal]",
	weight: 2,
}

func (l LogLevel) String() string {
	return l.prefix
}
