package console

import "runtime"

type consoleOptions struct {
	level LogLevel

	timestamp       bool
	timestampFormat string

	coloring bool
}

var globalConsoleOptions = &consoleOptions{
	level: NONE,
}

func WithGlobalLogLevel(level LogLevel) {
	globalConsoleOptions.level = level
}

func WithGlobalTimestamp() {
	globalConsoleOptions.timestamp = true
	globalConsoleOptions.timestampFormat = "2006-01-02 15:04:05"
}

func WithGlobalTimestampFormat(format string) {
	globalConsoleOptions.timestamp = true
	globalConsoleOptions.timestampFormat = format
}

func WithGlobalColoring() {
	if runtime.GOOS != "windows" {
		globalConsoleOptions.coloring = true
	}
}
