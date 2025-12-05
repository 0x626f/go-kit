package logger

import "runtime"

type loggerOptions struct {
	level LogLevel

	timestamp       bool
	timestampFormat string

	coloring bool

	async        bool
	logs, errors chan []byte
	cancelAsync  chan struct{}
}

var globalLoggerOptions = &loggerOptions{
	level: NONE,
}

func WithGlobalLogLevel(level LogLevel) {
	globalLoggerOptions.level = level
}

func WithGlobalTimestamp() {
	globalLoggerOptions.timestamp = true
	globalLoggerOptions.timestampFormat = "2006-01-02 15:04:05"
}

func WithGlobalTimestampFormat(format string) {
	globalLoggerOptions.timestamp = true
	globalLoggerOptions.timestampFormat = format
}

func WithGlobalColoring() {
	if runtime.GOOS != "windows" {
		globalLoggerOptions.coloring = true
	}
}
