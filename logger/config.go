package logger

import "runtime"

// Config holds configuration settings for a Logger instance.
// It controls log filtering, formatting, and Async behavior.
type Config struct {
	// Level is the minimum log Level to output (messages below this are discarded)
	Level LogLevel `env:"LOG_LEVEL"`

	// Timestamp enables/disables Timestamp prefixes on log messages
	Timestamp bool `env:"LOG_TIMESTAMP" default:"true"`
	// TimestampFormat is the Go time.Format string for Timestamp formatting
	TimestampFormat string `env:"LOG_TIMESTAMP_FORMAT" default:"2006-01-02 15:04:05"`
	// UseRegistry enables/disables logger registry
	UseRegistry bool `env:"LOG_USE_REGISTRY" default:"true"`
	// Coloring enables/disables ANSI color codes for log levels
	Coloring bool `env:"LOG_COLORING" default:"false"`
	// Async enables/disables asynchronous logging mode
	Async       bool `env:"LOG_ASYNC" default:"false"`
	AsyncBuffer int  `env:"LOG_ASYNC_BUFFER" default:"100"`
	// logs is the channel for buffering non-error log messages when Async is enabled
	logs, errors chan []byte
	// cancelAsync is used to shut down the Async logging goroutine
	cancelAsync chan struct{}
}

// defaultConfig contains default settings applied to all new Logger instances.
// These can be configured with WithGlobal* functions before creating loggers.
var defaultConfig = &Config{
	Level:           NONE,
	Timestamp:       false,
	TimestampFormat: "2006-01-02 15:04:05",
	UseRegistry:     true,
	Coloring:        false,
	Async:           false,
	AsyncBuffer:     100,
}

// WithDefaultLogLevel sets the default log Level for all newly created loggers.
// This does not affect existing logger instances.
//
// Parameters:
//   - Level: The minimum log Level to use for new loggers
//
// Example:
//
//	logger.WithDefaultLogLevel(logger.INFO)
//	logger1 := logger.NewLogger() // Uses INFO Level
//	logger2 := logger.NewLogger() // Also uses INFO Level
func WithDefaultLogLevel(level LogLevel) {
	defaultConfig.Level = level
}

// WithDefaultTimestamp enables Timestamp prefixes for all newly created loggers.
// Uses the default format "2006-01-02 15:04:05".
// This does not affect existing logger instances.
//
// Example:
//
//	logger.WithDefaultTimestamp()
//	logger := logger.NewLogger() // Will include timestamps
func WithDefaultTimestamp() {
	defaultConfig.Timestamp = true
	defaultConfig.TimestampFormat = "2006-01-02 15:04:05"
}

// WithDefaultTimestampFormat enables Timestamp prefixes with a custom format for all newly created loggers.
// The format string follows Go's time.Format conventions.
// This does not affect existing logger instances.
//
// Parameters:
//   - format: The Timestamp format string (e.g., time.RFC3339, "2006-01-02")
//
// Example:
//
//	logger.WithDefaultTimestampFormat(time.RFC3339)
//	logger := logger.NewLogger() // Timestamps use RFC3339 format
func WithDefaultTimestampFormat(format string) {
	defaultConfig.Timestamp = true
	defaultConfig.TimestampFormat = format
}

// WithDefaultColoring enables ANSI color codes for all newly created loggers.
// Colors are only enabled on non-Windows platforms.
// This does not affect existing logger instances.
//
// Example:
//
//	logger.WithDefaultColoring()
//	logger := logger.NewLogger() // Logs will be colored on Unix/Linux/macOS
func WithDefaultColoring() {
	if runtime.GOOS != "windows" {
		defaultConfig.Coloring = true
	}
}

func WithConfig(config *Config) {
	if config == nil {
		return
	}
	defaultConfig = config
}

// WithLoggerRegistry initializes the global logger registry.
// Once initialized, NewLogger() will automatically register loggers by name,
// and GetLogger() can be used to retrieve them.
//
// This enables centralized logger management where the same named logger
// instance is reused across the application. Call this once at application startup
// before creating any loggers if you want to use the registry feature.
//
// Example:
//
//	// Initialize registry at application startup
//	logger.WithLoggerRegistry()
//
//	// Create or get logger by name
//	apiLogger := logger.NewLogger("api")
//	dbLogger := logger.NewLogger("database")
//
//	// Later, retrieve the same instance
//	apiLogger2 := logger.GetLogger("api") // Returns the same instance as apiLogger
func WithLoggerRegistry() {
	defaultConfig.UseRegistry = true
	if loggerRegistry == nil {
		loggerRegistry = make(map[string]*Logger)
	}
}
