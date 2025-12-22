// Package logger provides a high-performance structured logging system with support for:
//   - Multiple log levels (ERROR, WARNING, INFO, DEBUG, TRACE)
//   - Colored output (Unix/Linux only)
//   - Timestamping with customizable formats
//   - JSON structured logging
//   - Zero-allocation object logging using sync.Pool
//   - Synchronous and asynchronous logging modes
//   - Named loggers for different components
//
// The logger supports three output formats:
//  1. Plain text logging with Logf(), Infof(), Errorf(), etc.
//  2. JSON logging with LogJSONf(), InfoJSONf(), etc.
//  3. Zero-allocation object logging with LogObjectf(), InfoObjectf(), etc.
//
// Example usage:
//
//	// Create a logger with configuration
//	logger := logger.NewLogger().
//	    WithName("api").
//	    WithLogLevel(logger.INFO).
//	    WithTimestamp().
//	    WithColoring()
//
//	// Plain text logging
//	logger.Infof("Server started on port %d", 8080)
//
//	// JSON logging
//	logger.InfoJSONf(map[string]any{"port": 8080, "env": "prod"}, "Server started")
//
//	// Zero-allocation object logging
//	logger.InfoObjectf("Request processed").
//	    AssignString("method", "GET").
//	    AssignInt("status", 200).
//	    AssignFloat64("duration", 0.125).
//	    Build()
package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

// loggerRegistry is a global registry that stores named logger instances.
// It is initialized by calling WithLoggerRegistry() and accessed via GetLogger().
var loggerRegistry map[string]*Logger

// GetLogger retrieves a logger instance from the registry by name.
// Returns nil if the registry is not initialized (via WithLoggerRegistry)
// or if no logger with the given name exists.
//
// Parameters:
//   - name: The name of the logger to retrieve
//
// Returns:
//   - The Logger instance if found, nil otherwise
//
// Example:
//
//	// Retrieve an existing logger
//	logger := logger.GetLogger("api")
//	if logger != nil {
//	    logger.Infof("Using existing logger")
//	}
func GetLogger(name string) *Logger {
	if loggerRegistry != nil {
		if logger, exists := loggerRegistry[name]; exists {
			return logger
		}
	}

	return NewLogger(name)
}

// Logger is a structured logger with support for multiple output formats and log levels.
// It provides thread-safe logging through mutex-protected output streams.
//
// The logger can operate in two modes:
//   - Synchronous: Logs are written immediately (default)
//   - Asynchronous: Logs are buffered in channels and written by a background goroutine
type Logger struct {
	// name is an optional identifier for this logger instance (e.g., "api", "database")
	name string
	// out is the standard output stream for INFO, DEBUG, TRACE logs
	out, err io.Writer
	// syncOut and syncErr provide thread-safe access to output streams
	syncOut, syncErr sync.Mutex
	// options holds the logger configuration
	options *Config
}

// jsonLog represents the structure of JSON-formatted log output.
// It matches common structured logging standards for easy parsing by log aggregators.
type jsonLog struct {
	// Source is the logger name (omitted if not set)
	Source string `json:"source,omitempty"`
	// Level is the log Level as a string (ERROR, WARNING, INFO, DEBUG, TRACE)
	Level string `json:"Level,omitempty"`
	// Timestamp is the log Timestamp in the configured format (omitted if timestamping disabled)
	Timestamp string `json:"Timestamp,omitempty"`
	// Message is the formatted log message
	Message string `json:"message,omitempty"`
	// Object contains structured data (can be any JSON-serializable value)
	Object any `json:"object,omitempty"`
}

// NewLogger creates a new Logger instance with the specified name and default configuration.
// By default, the logger:
//   - Writes to stdout for logs and stderr for errors
//   - Uses the global log Level (NONE by default, which logs everything)
//   - Has no Timestamp, no Coloring, and operates synchronously
//
// If the logger registry is enabled (via WithLoggerRegistry), this function will:
//   - Return the existing logger instance if one with the given name already exists
//   - Create and register a new logger if no instance with this name exists
//
// Parameters:
//   - name: The identifier for this logger (e.g., "api", "database", "cache")
//     The name appears in log output and is used for registry lookups
//
// Returns:
//   - A pointer to the Logger instance (new or existing)
//
// Example:
//
//	// Without registry (default behavior)
//	logger := logger.NewLogger("app")
//	logger.Infof("Application started")
//
//	// With registry (returns same instance for same name)
//	logger.WithLoggerRegistry()
//	apiLogger := logger.NewLogger("api")
//	sameLogger := logger.NewLogger("api") // Returns apiLogger
func NewLogger(name string, config ...*Config) *Logger {
	if loggerRegistry != nil {
		if logger, exists := loggerRegistry[name]; exists {
			return logger
		}
	}

	localConfig := defaultConfig
	if len(config) > 0 {
		localConfig = config[0]
	}

	logger := &Logger{
		name: name,
		out:  os.Stdout,
		err:  os.Stderr,
		options: &Config{
			Level:           localConfig.Level,
			Coloring:        localConfig.Coloring,
			Timestamp:       localConfig.Timestamp,
			TimestampFormat: localConfig.TimestampFormat,
			Async:           localConfig.Async,
			AsyncBuffer:     localConfig.AsyncBuffer,
		},
	}

	logger.WithAsync(logger.options.Async, logger.options.AsyncBuffer)

	if loggerRegistry != nil {
		loggerRegistry[name] = logger
	}

	return logger
}

// formatMessage formats a log message with optional Timestamp, Level, and logger name.
// Returns the formatted message as a byte slice ready for writing to output.
//
// Parameters:
//   - Level: The log Level (used for filtering and formatting)
//   - msg: The message format string
//   - args: Optional format arguments for msg
//
// Returns:
//   - Formatted message bytes with newline
func (logger *Logger) formatMessage(level LogLevel, msg string, args ...any) []byte {
	message := logger.format(msg, args...)
	if level == NONE {
		return append(message, ln)
	}

	var payload []byte

	if logger.options.Timestamp {
		timestamp := time.Now().Format(logger.options.TimestampFormat)
		payload = append(payload, []byte(timestamp)...)
	}

	if len(payload) > 0 {
		payload = append(payload, space)
	}
	payload = append(payload, []byte(level.String())...)

	if len(logger.name) > 0 {
		payload = append(payload, space)
		payload = append(payload, '[')
		payload = append(payload, []byte(logger.name)...)
		payload = append(payload, ']')
	}

	payload = append(payload, ':')
	payload = append(payload, space)
	payload = append(payload, message...)

	if logger.options.Coloring {
		payload = level.paint(payload)
	}

	return append(payload, ln)
}

// formatJSONMessage formats a log entry as JSON with message and structured object data.
//
// Parameters:
//   - Level: The log Level
//   - object: Structured data to include in the JSON output
//   - msg: The message format string
//   - args: Optional format arguments for msg
//
// Returns:
//   - JSON-formatted log bytes with newline
//   - Error if JSON marshaling fails
func (logger *Logger) formatJSONMessage(level LogLevel, object any, msg string, args ...any) ([]byte, error) {
	log := jsonLog{
		Level:   level.String(),
		Message: fmt.Sprintf(msg, args...),
		Object:  object,
	}

	if level == NONE {
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		return append(raw, ln), nil
	}

	if logger.options.Timestamp {
		log.Timestamp = time.Now().Format(logger.options.TimestampFormat)
	}

	if len(logger.name) > 0 {
		log.Source = logger.name
	}

	raw, err := json.Marshal(log)
	if err != nil {
		return nil, err
	}

	return append(raw, ln), nil
}

// format is an internal helper that formats a message with optional arguments.
//
// Parameters:
//   - msg: The message format string
//   - args: Optional format arguments
//
// Returns:
//   - Formatted message as bytes
func (logger *Logger) format(msg string, args ...any) []byte {
	if len(args) == 0 {
		return []byte(msg)
	}
	return []byte(fmt.Sprintf(msg, args...))
}

// writeStringToStream writes a formatted text log message to the specified stream.
//
// Parameters:
//   - stream: The output writer
//   - Level: The log Level
//   - msg: The message format string
//   - args: Optional format arguments
func (logger *Logger) writeStringToStream(stream io.Writer, level LogLevel, msg string, args ...any) {
	if stream == nil {
		return
	}
	_, _ = stream.Write(logger.formatMessage(level, msg, args...))
}

// writeJSONToStream writes a JSON-formatted log message to the specified stream.
//
// Parameters:
//   - stream: The output writer
//   - Level: The log Level
//   - object: Structured data to include
//   - msg: The message format string
//   - args: Optional format arguments
//
// Returns:
//   - Error if stream is nil or JSON marshaling fails
func (logger *Logger) writeJSONToStream(stream io.Writer, level LogLevel, object any, msg string, args ...any) error {
	if stream == nil {
		return errors.New("nil stream or object")
	}

	payload, err := logger.formatJSONMessage(level, object, msg, args...)

	if err != nil {
		return err
	}

	_, err = stream.Write(payload)

	return err
}

// writeByLevel writes data to the appropriate output stream based on log Level.
// ERROR logs go to stderr, all others go to stdout.
// This method handles locking for thread safety.
//
// Parameters:
//   - Level: The log Level (determines output stream)
//   - data: The formatted log data to write
func (logger *Logger) writeByLevel(level LogLevel, data []byte) {
	if level == ERROR {
		logger.syncErr.Lock()
		defer logger.syncErr.Unlock()
		_, _ = logger.err.Write(data)
		return
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()
	_, _ = logger.out.Write(data)
}

// sendToChannelByLevel sends log data to the appropriate Async channel based on log Level.
// Used when Async logging is enabled.
//
// Parameters:
//   - Level: The log Level (determines which channel to use)
//   - data: The formatted log data to send
func (logger *Logger) sendToChannelByLevel(level LogLevel, data []byte) {
	if level == ERROR {
		logger.options.errors <- data
		return
	}

	logger.options.logs <- data
}

// OutputTo sets the output stream for INFO, DEBUG, TRACE, and WARNING logs.
//
// Parameters:
//   - target: The io.Writer to use for output (e.g., os.Stdout, file, buffer)
//
// Returns:
//   - The logger for method chaining
//
// Example:
//
//	file, _ := os.Create("app.log")
//	logger := logger.NewLogger().OutputTo(file)
func (logger *Logger) OutputTo(target io.Writer) *Logger {
	if target != nil {
		logger.out = target
	}
	return logger
}

// ErrorsTo sets the output stream for ERROR logs.
//
// Parameters:
//   - target: The io.Writer to use for error output
//
// Returns:
//   - The logger for method chaining
//
// Example:
//
//	file, _ := os.Create("errors.log")
//	logger := logger.NewLogger().ErrorsTo(file)
func (logger *Logger) ErrorsTo(target io.Writer) *Logger {
	if target != nil {
		logger.err = target
	}
	return logger
}

// WithAsync enables asynchronous logging mode.
// When enabled, log messages are sent to buffered channels and written by a background goroutine.
// This improves performance by preventing I/O operations from blocking the caller.
//
// Parameters:
//   - option: If true, enables Async mode; if false, no effect
//   - capacity: The buffer size for the Async channels
//
// Returns:
//   - The logger for method chaining
//   - A cancel function that must be called to shut down the Async goroutine
//
// Example:
//
//	logger, cancel := logger.NewLogger().WithAsync(true, 100)
//	defer cancel()
//	logger.Infof("Async log message")
func (logger *Logger) WithAsync(option bool, capacity int) (*Logger, func()) {
	cancel := func() {}
	if option {
		if logger.options.logs != nil {
			close(logger.options.logs)
		}

		if logger.options.errors != nil {
			close(logger.options.errors)
		}

		if logger.options.cancelAsync != nil {
			close(logger.options.cancelAsync)
		}

		logger.options.Async = true
		logger.options.logs, logger.options.errors = make(chan []byte, capacity), make(chan []byte, capacity)
		logger.options.cancelAsync = make(chan struct{})

		cancel = func() {
			if logger.options.cancelAsync != nil {
				close(logger.options.cancelAsync)
			}
		}

		go func(logs, errs chan []byte, cancel chan struct{}) {
			for {
				select {
				case log := <-logs:
					logger.syncOut.Lock()
					_, _ = logger.out.Write(log)
					logger.syncOut.Unlock()
				case err := <-errs:
					logger.syncErr.Lock()
					_, _ = logger.err.Write(err)
					logger.syncErr.Unlock()
				case <-cancel:
					return
				}
			}
		}(logger.options.logs, logger.options.errors, logger.options.cancelAsync)

	}
	return logger, cancel
}

// WithLogLevel sets the minimum log Level for this logger.
// Logs below this Level are discarded.
// Levels in order: ERROR < WARNING < INFO < DEBUG < TRACE < NONE
//
// Parameters:
//   - Level: The minimum log Level to output
//
// Returns:
//   - The logger for method chaining
//
// Example:
//
//	logger := logger.NewLogger().WithLogLevel(logger.INFO)
//	logger.Debugf("This won't be logged") // DEBUG < INFO
//	logger.Infof("This will be logged")   // INFO >= INFO
func (logger *Logger) WithLogLevel(level LogLevel) *Logger {
	logger.options.Level = level
	return logger
}

// WithTimestamp enables Timestamp prefixes on log messages.
// Uses the default format "2006-01-02 15:04:05".
//
// Returns:
//   - The logger for method chaining
//
// Example:
//
//	logger := logger.NewLogger().WithTimestamp()
//	logger.Infof("Message") // Output: "2025-12-06 10:30:45 INFO Message"
func (logger *Logger) WithTimestamp() *Logger {
	logger.options.Timestamp = true
	logger.options.TimestampFormat = "2006-01-02 15:04:05"
	return logger
}

// WithTimestampFormat enables Timestamp prefixes with a custom format.
// The format string follows Go's time.Format conventions.
//
// Parameters:
//   - format: The Timestamp format string (e.g., "2006-01-02", time.RFC3339)
//
// Returns:
//   - The logger for method chaining
//
// Example:
//
//	logger := logger.NewLogger().WithTimestampFormat(time.RFC3339)
//	logger.Infof("Message") // Output: "2025-12-06T10:30:45Z INFO Message"
func (logger *Logger) WithTimestampFormat(format string) *Logger {
	logger.options.Timestamp = true
	logger.options.TimestampFormat = format
	return logger
}

// WithColoring enables ANSI color codes for different log levels.
// Colors are only enabled on non-Windows platforms.
//   - ERROR: Red
//   - WARNING: Yellow
//   - INFO: Green
//   - DEBUG: Grey
//   - TRACE: Blue
//
// Returns:
//   - The logger for method chaining
//
// Example:
//
//	logger := logger.NewLogger().WithColoring()
//	logger.Errorf("Error message") // Displayed in red on Unix terminals
func (logger *Logger) WithColoring() *Logger {
	if runtime.GOOS != "windows" {
		logger.options.Coloring = true
	}
	return logger
}

// Logf logs a message without a log Level prefix.
// This message is always logged regardless of the configured log Level.
//
// Parameters:
//   - msg: The message format string
//   - args: Optional format arguments
//
// Example:
//
//	logger.Logf("Server started on port %d", 8080)
func (logger *Logger) Logf(msg string, args ...any) {
	if logger.options.Async {
		logger.options.logs <- logger.formatMessage(NONE, msg, args...)
		return
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	logger.writeStringToStream(logger.out, NONE, msg, args...)
}

// Tracef logs a message at TRACE Level.
// TRACE is the most verbose Level, typically used for detailed diagnostic information.
//
// Parameters:
//   - msg: The message format string
//   - args: Optional format arguments
func (logger *Logger) Tracef(msg string, args ...any) {
	if logger.options.Level < TRACE {
		return
	}

	if logger.options.Async {
		logger.options.logs <- logger.formatMessage(TRACE, msg, args...)
		return
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	logger.writeStringToStream(logger.out, TRACE, msg, args...)
}

// Debugf logs a message at DEBUG Level.
// DEBUG logs are typically used for development and troubleshooting.
//
// Parameters:
//   - msg: The message format string
//   - args: Optional format arguments
func (logger *Logger) Debugf(msg string, args ...any) {
	if logger.options.Level < DEBUG {
		return
	}

	if logger.options.Async {
		logger.options.logs <- logger.formatMessage(DEBUG, msg, args...)
		return
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	logger.writeStringToStream(logger.out, DEBUG, msg, args...)
}

// Infof logs a message at INFO Level.
// INFO logs are for general informational messages about application operation.
//
// Parameters:
//   - msg: The message format string
//   - args: Optional format arguments
func (logger *Logger) Infof(msg string, args ...any) {
	if logger.options.Level < INFO {
		return
	}

	if logger.options.Async {
		logger.options.logs <- logger.formatMessage(INFO, msg, args...)
		return
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	logger.writeStringToStream(logger.out, INFO, msg, args...)
}

// Warningf logs a message at WARNING Level.
// WARNING logs indicate potentially harmful situations that should be reviewed.
//
// Parameters:
//   - msg: The message format string
//   - args: Optional format arguments
func (logger *Logger) Warningf(msg string, args ...any) {
	if logger.options.Level < WARNING {
		return
	}

	if logger.options.Async {
		logger.options.logs <- logger.formatMessage(WARNING, msg, args...)
		return
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	logger.writeStringToStream(logger.out, WARNING, msg, args...)
}

// Errorf logs a message at ERROR Level.
// ERROR logs indicate serious problems that should be addressed.
// These logs are written to the error stream (stderr by default).
//
// Parameters:
//   - msg: The message format string
//   - args: Optional format arguments
func (logger *Logger) Errorf(msg string, args ...any) {
	if logger.options.Async {
		logger.options.errors <- logger.formatMessage(ERROR, msg, args...)
		return
	}

	logger.syncErr.Lock()
	defer logger.syncErr.Unlock()

	logger.writeStringToStream(logger.err, ERROR, msg, args...)
}

// LogJSONf logs a message with structured JSON data at no specific Level.
// The object is marshaled to JSON and included in the log output.
//
// Parameters:
//   - object: Any JSON-serializable value to include in the log
//   - msg: The message format string
//   - args: Optional format arguments
//
// Returns:
//   - Error if JSON marshaling fails
//
// Example:
//
//	logger.LogJSONf(map[string]any{"user": "alice", "action": "login"}, "User activity")
func (logger *Logger) LogJSONf(object any, msg string, args ...any) error {
	if logger.options.Async {
		data, err := logger.formatJSONMessage(NONE, object, msg, args...)
		if err != nil {
			return err
		}
		logger.options.logs <- data
		return nil
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	return logger.writeJSONToStream(logger.out, NONE, object, msg, args...)
}

// TraceJSONf logs a message with structured JSON data at TRACE Level.
//
// Parameters:
//   - object: Any JSON-serializable value
//   - msg: The message format string
//   - args: Optional format arguments
//
// Returns:
//   - Error if JSON marshaling fails
func (logger *Logger) TraceJSONf(object any, msg string, args ...any) error {
	if logger.options.Level < TRACE {
		return nil
	}

	if logger.options.Async {
		data, err := logger.formatJSONMessage(TRACE, object, msg, args...)
		if err != nil {
			return err
		}
		logger.options.logs <- data
		return nil
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	return logger.writeJSONToStream(logger.out, TRACE, object, msg, args...)
}

// DebugJSONf logs a message with structured JSON data at DEBUG Level.
//
// Parameters:
//   - object: Any JSON-serializable value
//   - msg: The message format string
//   - args: Optional format arguments
//
// Returns:
//   - Error if JSON marshaling fails
func (logger *Logger) DebugJSONf(object any, msg string, args ...any) error {
	if logger.options.Level < DEBUG {
		return nil
	}

	if logger.options.Async {
		data, err := logger.formatJSONMessage(DEBUG, object, msg, args...)
		if err != nil {
			return err
		}
		logger.options.logs <- data
		return nil
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	return logger.writeJSONToStream(logger.out, DEBUG, object, msg, args...)
}

// InfoJSONf logs a message with structured JSON data at INFO Level.
//
// Parameters:
//   - object: Any JSON-serializable value
//   - msg: The message format string
//   - args: Optional format arguments
//
// Returns:
//   - Error if JSON marshaling fails
func (logger *Logger) InfoJSONf(object any, msg string, args ...any) error {
	if logger.options.Level < INFO {
		return nil
	}

	if logger.options.Async {
		data, err := logger.formatJSONMessage(INFO, object, msg, args...)
		if err != nil {
			return err
		}
		logger.options.logs <- data
		return nil
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	return logger.writeJSONToStream(logger.out, INFO, object, msg, args...)
}

// WarningJSONf logs a message with structured JSON data at WARNING Level.
//
// Parameters:
//   - object: Any JSON-serializable value
//   - msg: The message format string
//   - args: Optional format arguments
//
// Returns:
//   - Error if JSON marshaling fails
func (logger *Logger) WarningJSONf(object any, msg string, args ...any) error {
	if logger.options.Level < WARNING {
		return nil
	}

	if logger.options.Async {
		data, err := logger.formatJSONMessage(WARNING, object, msg, args...)
		if err != nil {
			return err
		}
		logger.options.logs <- data
		return nil
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	return logger.writeJSONToStream(logger.out, WARNING, object, msg, args...)
}

// ErrorJSONf logs a message with structured JSON data at ERROR Level.
//
// Parameters:
//   - object: Any JSON-serializable value
//   - msg: The message format string
//   - args: Optional format arguments
//
// Returns:
//   - Error if JSON marshaling fails
func (logger *Logger) ErrorJSONf(object any, msg string, args ...any) error {
	if logger.options.Async {
		data, err := logger.formatJSONMessage(ERROR, object, msg, args...)
		if err != nil {
			return err
		}
		logger.options.errors <- data
		return nil
	}

	logger.syncErr.Lock()
	defer logger.syncErr.Unlock()

	return logger.writeJSONToStream(logger.err, ERROR, object, msg, args...)
}

// LogObjectf creates a zero-allocation object log builder at no specific Level.
// Use the builder's Assign* methods to add fields, then call Build() to emit the log.
//
// Parameters:
//   - msg: The message format string
//   - args: Optional format arguments
//
// Returns:
//   - An ObjectLogBuilder for constructing the structured log
//
// Example:
//
//	logger.LogObjectf("Request completed").
//	    AssignString("method", "GET").
//	    AssignInt("status", 200).
//	    Build()
func (logger *Logger) LogObjectf(msg string, args ...any) *ObjectLogBuilder {
	return newObjectLogBuilder(logger, NONE, logger.format(msg, args...))
}

// TraceObjectf creates a zero-allocation object log builder at TRACE Level.
//
// Parameters:
//   - msg: The message format string
//   - args: Optional format arguments
//
// Returns:
//   - An ObjectLogBuilder, or nil if log Level filtering discards this log
func (logger *Logger) TraceObjectf(msg string, args ...any) *ObjectLogBuilder {
	if logger.options.Level < TRACE {
		return nil
	}

	return newObjectLogBuilder(logger, TRACE, logger.format(msg, args...))

}

// DebugObjectf creates a zero-allocation object log builder at DEBUG Level.
//
// Parameters:
//   - msg: The message format string
//   - args: Optional format arguments
//
// Returns:
//   - An ObjectLogBuilder, or nil if log Level filtering discards this log
func (logger *Logger) DebugObjectf(msg string, args ...any) *ObjectLogBuilder {
	if logger.options.Level < DEBUG {
		return nil
	}

	return newObjectLogBuilder(logger, DEBUG, logger.format(msg, args...))

}

// InfoObjectf creates a zero-allocation object log builder at INFO Level.
//
// Parameters:
//   - msg: The message format string
//   - args: Optional format arguments
//
// Returns:
//   - An ObjectLogBuilder, or nil if log Level filtering discards this log
func (logger *Logger) InfoObjectf(msg string, args ...any) *ObjectLogBuilder {
	if logger.options.Level < INFO {
		return nil
	}

	return newObjectLogBuilder(logger, INFO, logger.format(msg, args...))

}

// WarningObjectf creates a zero-allocation object log builder at WARNING Level.
//
// Parameters:
//   - msg: The message format string
//   - args: Optional format arguments
//
// Returns:
//   - An ObjectLogBuilder, or nil if log Level filtering discards this log
func (logger *Logger) WarningObjectf(msg string, args ...any) *ObjectLogBuilder {
	if logger.options.Level < WARNING {
		return nil
	}

	return newObjectLogBuilder(logger, WARNING, logger.format(msg, args...))

}

// ErrorObjectf creates a zero-allocation object log builder at ERROR Level.
//
// Parameters:
//   - msg: The message format string
//   - args: Optional format arguments
//
// Returns:
//   - An ObjectLogBuilder for constructing the error log
func (logger *Logger) ErrorObjectf(msg string, args ...any) *ObjectLogBuilder {
	return newObjectLogBuilder(logger, ERROR, logger.format(msg, args...))

}
