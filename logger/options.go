package logger

import "runtime"

// loggerOptions holds configuration settings for a Logger instance.
// It controls log filtering, formatting, and async behavior.
type loggerOptions struct {
	// level is the minimum log level to output (messages below this are discarded)
	level LogLevel

	// timestamp enables/disables timestamp prefixes on log messages
	timestamp bool
	// timestampFormat is the Go time.Format string for timestamp formatting
	timestampFormat string

	// coloring enables/disables ANSI color codes for log levels
	coloring bool

	// async enables/disables asynchronous logging mode
	async bool
	// logs is the channel for buffering non-error log messages when async is enabled
	logs, errors chan []byte
	// cancelAsync is used to shut down the async logging goroutine
	cancelAsync chan struct{}
}

// globalLoggerOptions contains default settings applied to all new Logger instances.
// These can be configured with WithGlobal* functions before creating loggers.
var globalLoggerOptions = &loggerOptions{
	level: NONE,
}

// WithGlobalLogLevel sets the default log level for all newly created loggers.
// This does not affect existing logger instances.
//
// Parameters:
//   - level: The minimum log level to use for new loggers
//
// Example:
//
//	logger.WithGlobalLogLevel(logger.INFO)
//	logger1 := logger.NewLogger() // Uses INFO level
//	logger2 := logger.NewLogger() // Also uses INFO level
func WithGlobalLogLevel(level LogLevel) {
	globalLoggerOptions.level = level
}

// WithGlobalTimestamp enables timestamp prefixes for all newly created loggers.
// Uses the default format "2006-01-02 15:04:05".
// This does not affect existing logger instances.
//
// Example:
//
//	logger.WithGlobalTimestamp()
//	logger := logger.NewLogger() // Will include timestamps
func WithGlobalTimestamp() {
	globalLoggerOptions.timestamp = true
	globalLoggerOptions.timestampFormat = "2006-01-02 15:04:05"
}

// WithGlobalTimestampFormat enables timestamp prefixes with a custom format for all newly created loggers.
// The format string follows Go's time.Format conventions.
// This does not affect existing logger instances.
//
// Parameters:
//   - format: The timestamp format string (e.g., time.RFC3339, "2006-01-02")
//
// Example:
//
//	logger.WithGlobalTimestampFormat(time.RFC3339)
//	logger := logger.NewLogger() // Timestamps use RFC3339 format
func WithGlobalTimestampFormat(format string) {
	globalLoggerOptions.timestamp = true
	globalLoggerOptions.timestampFormat = format
}

// WithGlobalColoring enables ANSI color codes for all newly created loggers.
// Colors are only enabled on non-Windows platforms.
// This does not affect existing logger instances.
//
// Example:
//
//	logger.WithGlobalColoring()
//	logger := logger.NewLogger() // Logs will be colored on Unix/Linux/macOS
func WithGlobalColoring() {
	if runtime.GOOS != "windows" {
		globalLoggerOptions.coloring = true
	}
}
