package logger

import (
	"runtime"
	"testing"
	"time"
)

// TestWithDefaultLogLevel tests that the default log level is applied to new loggers.
func TestWithDefaultLogLevel(t *testing.T) {
	// Save original defaultConfig to restore after test
	originalLevel := defaultConfig.Level
	defer func() { defaultConfig.Level = originalLevel }()

	tests := []LogLevel{ERROR, WARNING, INFO, DEBUG, TRACE, NONE}
	for _, level := range tests {
		t.Run(level.String(), func(t *testing.T) {
			WithDefaultLogLevel(level)
			if defaultConfig.Level != level {
				t.Errorf("Expected default level %v, got %v", level, defaultConfig.Level)
			}

			// Verify new loggers use this default
			logger := NewLogger("test")
			if logger.options.Level != level {
				t.Errorf("New logger should have level %v, got %v", level, logger.options.Level)
			}
		})
	}
}

// TestWithDefaultTimestamp tests that default timestamp settings are applied to new loggers.
func TestWithDefaultTimestamp(t *testing.T) {
	// Save original defaultConfig to restore after test
	originalTimestamp := defaultConfig.Timestamp
	originalFormat := defaultConfig.TimestampFormat
	defer func() {
		defaultConfig.Timestamp = originalTimestamp
		defaultConfig.TimestampFormat = originalFormat
	}()

	WithDefaultTimestamp()

	if !defaultConfig.Timestamp {
		t.Error("WithDefaultTimestamp should enable timestamp")
	}

	expectedFormat := "2006-01-02 15:04:05"
	if defaultConfig.TimestampFormat != expectedFormat {
		t.Errorf("Expected timestamp format '%s', got '%s'", expectedFormat, defaultConfig.TimestampFormat)
	}

	// Verify new loggers use this default
	logger := NewLogger("test")
	if !logger.options.Timestamp {
		t.Error("New logger should have timestamp enabled")
	}
	if logger.options.TimestampFormat != expectedFormat {
		t.Errorf("New logger should have timestamp format '%s', got '%s'", expectedFormat, logger.options.TimestampFormat)
	}
}

// TestWithDefaultTimestampFormat tests custom timestamp format configuration.
func TestWithDefaultTimestampFormat(t *testing.T) {
	// Save original defaultConfig to restore after test
	originalTimestamp := defaultConfig.Timestamp
	originalFormat := defaultConfig.TimestampFormat
	defer func() {
		defaultConfig.Timestamp = originalTimestamp
		defaultConfig.TimestampFormat = originalFormat
	}()

	customFormats := []string{
		time.RFC3339,
		time.RFC822,
		"2006/01/02",
		"15:04:05",
		"2006-01-02T15:04:05.000Z",
	}

	for _, format := range customFormats {
		t.Run(format, func(t *testing.T) {
			WithDefaultTimestampFormat(format)

			if !defaultConfig.Timestamp {
				t.Error("WithDefaultTimestampFormat should enable timestamp")
			}

			if defaultConfig.TimestampFormat != format {
				t.Errorf("Expected timestamp format '%s', got '%s'", format, defaultConfig.TimestampFormat)
			}

			// Verify new loggers use this default
			logger := NewLogger("test")
			if !logger.options.Timestamp {
				t.Error("New logger should have timestamp enabled")
			}
			if logger.options.TimestampFormat != format {
				t.Errorf("New logger should have timestamp format '%s', got '%s'", format, logger.options.TimestampFormat)
			}
		})
	}
}

// TestWithDefaultColoring tests that coloring is enabled on non-Windows platforms.
func TestWithDefaultColoring(t *testing.T) {
	// Save original defaultConfig to restore after test
	originalColoring := defaultConfig.Coloring
	defer func() { defaultConfig.Coloring = originalColoring }()

	WithDefaultColoring()

	// Coloring should only be enabled on non-Windows platforms
	if runtime.GOOS != "windows" {
		if !defaultConfig.Coloring {
			t.Error("WithDefaultColoring should enable coloring on non-Windows platforms")
		}

		// Verify new loggers use this default
		logger := NewLogger("test")
		if !logger.options.Coloring {
			t.Error("New logger should have coloring enabled on non-Windows platforms")
		}
	} else {
		if defaultConfig.Coloring {
			t.Error("WithDefaultColoring should not enable coloring on Windows")
		}
	}
}

// TestWithConfig tests setting a complete custom configuration.
func TestWithConfig(t *testing.T) {
	// Save original defaultConfig to restore after test
	originalConfig := defaultConfig
	defer func() { defaultConfig = originalConfig }()

	customConfig := &Config{
		Level:           DEBUG,
		Timestamp:       true,
		TimestampFormat: time.RFC3339,
		UseRegistry:     false,
		Coloring:        true,
		Async:           true,
		AsyncBuffer:     500,
	}

	WithConfig(customConfig)

	if defaultConfig != customConfig {
		t.Error("WithConfig should replace the default config")
	}

	// Verify all fields were set
	if defaultConfig.Level != DEBUG {
		t.Errorf("Expected level DEBUG, got %v", defaultConfig.Level)
	}
	if !defaultConfig.Timestamp {
		t.Error("Expected timestamp to be true")
	}
	if defaultConfig.TimestampFormat != time.RFC3339 {
		t.Errorf("Expected timestamp format %s, got %s", time.RFC3339, defaultConfig.TimestampFormat)
	}
	if defaultConfig.UseRegistry {
		t.Error("Expected UseRegistry to be false")
	}
	if !defaultConfig.Coloring {
		t.Error("Expected coloring to be true")
	}
	if !defaultConfig.Async {
		t.Error("Expected async to be true")
	}
	if defaultConfig.AsyncBuffer != 500 {
		t.Errorf("Expected async buffer 500, got %d", defaultConfig.AsyncBuffer)
	}
}

// TestWithConfig_Nil tests that nil config is handled safely.
func TestWithConfig_Nil(t *testing.T) {
	// Save original defaultConfig to restore after test
	originalConfig := defaultConfig
	defer func() { defaultConfig = originalConfig }()

	WithConfig(nil)

	// Default config should remain unchanged
	if defaultConfig != originalConfig {
		t.Error("WithConfig(nil) should not change default config")
	}
}

// TestWithLoggerRegistry tests logger registry initialization and usage.
func TestWithLoggerRegistry(t *testing.T) {
	// Reset registry
	loggerRegistry = nil

	WithLoggerRegistry()

	if loggerRegistry == nil {
		t.Fatal("WithLoggerRegistry should initialize the registry")
	}

	if !defaultConfig.UseRegistry {
		t.Error("WithLoggerRegistry should set UseRegistry to true")
	}

	// Test that registry works
	logger1 := NewLogger("test1")
	logger2 := NewLogger("test2")
	logger1Again := NewLogger("test1")

	if logger1 != logger1Again {
		t.Error("NewLogger should return same instance for same name when registry is enabled")
	}

	if logger1 == logger2 {
		t.Error("NewLogger should return different instances for different names")
	}

	// Verify loggers are in registry
	if loggerRegistry["test1"] != logger1 {
		t.Error("Logger should be registered in registry")
	}
	if loggerRegistry["test2"] != logger2 {
		t.Error("Logger should be registered in registry")
	}

	// Clean up
	loggerRegistry = nil
}

// TestWithLoggerRegistry_AlreadyInitialized tests calling WithLoggerRegistry multiple times.
func TestWithLoggerRegistry_AlreadyInitialized(t *testing.T) {
	// Reset registry
	loggerRegistry = nil

	// First call
	WithLoggerRegistry()
	logger1 := NewLogger("test")

	// Second call should not reset registry
	WithLoggerRegistry()
	logger2 := GetLogger("test")

	if logger1 != logger2 {
		t.Error("WithLoggerRegistry should not reset existing registry")
	}

	// Clean up
	loggerRegistry = nil
}

// TestConfig_DefaultValues tests that the default config has expected initial values.
func TestConfig_DefaultValues(t *testing.T) {
	// Create a fresh config to test defaults
	config := &Config{
		Level:           NONE,
		Timestamp:       false,
		TimestampFormat: "2006-01-02 15:04:05",
		UseRegistry:     true,
		Coloring:        false,
		Async:           false,
		AsyncBuffer:     100,
	}

	// These should match the defaults in config.go
	if config.Level != NONE {
		t.Errorf("Default level should be NONE, got %v", config.Level)
	}
	if config.Timestamp {
		t.Error("Default timestamp should be false")
	}
	if config.TimestampFormat != "2006-01-02 15:04:05" {
		t.Errorf("Default timestamp format should be '2006-01-02 15:04:05', got '%s'", config.TimestampFormat)
	}
	if !config.UseRegistry {
		t.Error("Default UseRegistry should be true")
	}
	if config.Coloring {
		t.Error("Default coloring should be false")
	}
	if config.Async {
		t.Error("Default async should be false")
	}
	if config.AsyncBuffer != 100 {
		t.Errorf("Default async buffer should be 100, got %d", config.AsyncBuffer)
	}
}

// TestNewLogger_UsesDefaultConfig tests that new loggers inherit from defaultConfig.
func TestNewLogger_UsesDefaultConfig(t *testing.T) {
	// Save original defaultConfig to restore after test
	originalConfig := &Config{
		Level:           defaultConfig.Level,
		Timestamp:       defaultConfig.Timestamp,
		TimestampFormat: defaultConfig.TimestampFormat,
		UseRegistry:     defaultConfig.UseRegistry,
		Coloring:        defaultConfig.Coloring,
		Async:           defaultConfig.Async,
		AsyncBuffer:     defaultConfig.AsyncBuffer,
	}
	defer func() { defaultConfig = originalConfig }()

	// Set custom defaults
	defaultConfig.Level = WARNING
	defaultConfig.Timestamp = true
	defaultConfig.TimestampFormat = time.RFC3339
	defaultConfig.Coloring = true

	logger := NewLogger("test")

	if logger.options.Level != WARNING {
		t.Errorf("Logger should inherit level WARNING, got %v", logger.options.Level)
	}
	if !logger.options.Timestamp {
		t.Error("Logger should inherit timestamp=true")
	}
	if logger.options.TimestampFormat != time.RFC3339 {
		t.Errorf("Logger should inherit timestamp format %s, got %s", time.RFC3339, logger.options.TimestampFormat)
	}
	if !logger.options.Coloring {
		t.Error("Logger should inherit coloring=true")
	}
}

// TestConfig_AsyncChannels tests that async channels are properly initialized.
func TestConfig_AsyncChannels(t *testing.T) {
	config := &Config{
		Async:       true,
		AsyncBuffer: 50,
	}

	// Initially channels should be nil
	if config.logs != nil {
		t.Error("logs channel should be nil before initialization")
	}
	if config.errors != nil {
		t.Error("errors channel should be nil before initialization")
	}
	if config.cancelAsync != nil {
		t.Error("cancelAsync channel should be nil before initialization")
	}
}

// TestDefaultConfig_Independence tests that modifying a logger's config doesn't affect default config.
func TestDefaultConfig_Independence(t *testing.T) {
	// Save original defaultConfig to restore after test
	originalLevel := defaultConfig.Level
	defer func() { defaultConfig.Level = originalLevel }()

	defaultConfig.Level = INFO

	logger := NewLogger("test")
	logger.WithLogLevel(ERROR)

	// Default config should remain unchanged
	if defaultConfig.Level != INFO {
		t.Errorf("Default config should remain INFO, got %v", defaultConfig.Level)
	}

	// Logger should have ERROR
	if logger.options.Level != ERROR {
		t.Errorf("Logger should have ERROR level, got %v", logger.options.Level)
	}
}

// TestConfig_MultipleDefaults tests setting multiple default configurations in sequence.
func TestConfig_MultipleDefaults(t *testing.T) {
	// Save original defaultConfig to restore after test
	originalConfig := &Config{
		Level:           defaultConfig.Level,
		Timestamp:       defaultConfig.Timestamp,
		TimestampFormat: defaultConfig.TimestampFormat,
		Coloring:        defaultConfig.Coloring,
	}
	defer func() {
		defaultConfig.Level = originalConfig.Level
		defaultConfig.Timestamp = originalConfig.Timestamp
		defaultConfig.TimestampFormat = originalConfig.TimestampFormat
		defaultConfig.Coloring = originalConfig.Coloring
	}()

	// Set multiple defaults
	WithDefaultLogLevel(DEBUG)
	WithDefaultTimestamp()
	WithDefaultColoring()

	// Verify all were applied
	if defaultConfig.Level != DEBUG {
		t.Errorf("Expected level DEBUG, got %v", defaultConfig.Level)
	}
	if !defaultConfig.Timestamp {
		t.Error("Expected timestamp to be true")
	}

	expectedOnNonWindows := runtime.GOOS != "windows"
	if defaultConfig.Coloring != expectedOnNonWindows {
		t.Errorf("Expected coloring %v, got %v", expectedOnNonWindows, defaultConfig.Coloring)
	}

	// Verify new logger inherits all defaults
	logger := NewLogger("test")
	if logger.options.Level != DEBUG {
		t.Errorf("Logger should have level DEBUG, got %v", logger.options.Level)
	}
	if !logger.options.Timestamp {
		t.Error("Logger should have timestamp enabled")
	}
	if logger.options.Coloring != expectedOnNonWindows {
		t.Errorf("Logger should have coloring %v, got %v", expectedOnNonWindows, logger.options.Coloring)
	}
}
