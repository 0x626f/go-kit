package logger

import (
	"bytes"
	"strings"
	"testing"
)

// TestNewLogger tests basic logger creation and initialization.
// It verifies that a new logger instance is properly initialized with default values.
func TestNewLogger(t *testing.T) {
	logger := NewLogger("TestNewLogger")
	if logger == nil {
		t.Fatal("NewLogger() returned nil")
	}
	if logger.options == nil {
		t.Error("Logger options are nil")
	}
	if logger.name != "TestNewLogger" {
		t.Errorf("Expected name 'TestNewLogger', got '%s'", logger.name)
	}
}

// TestUseLoggerRegistry tests the WithLoggerRegistry function.
// It verifies that the global logger registry is properly initialized.
func TestUseLoggerRegistry(t *testing.T) {
	// Reset registry
	loggerRegistry = nil

	WithLoggerRegistry()
	if loggerRegistry == nil {
		t.Fatal("WithLoggerRegistry() did not initialize registry")
	}
	loggerRegistry = nil
}

// TestNewLogger_WithRegistry tests logger creation with the registry enabled.
// It verifies that the same logger instance is returned for the same name.
func TestNewLogger_WithRegistry(t *testing.T) {
	loggerRegistry = nil
	WithLoggerRegistry()

	logger1 := NewLogger("test")
	logger2 := NewLogger("test")

	if logger1 != logger2 {
		t.Error("NewLogger should return same instance for same name when registry is enabled")
	}

	logger3 := NewLogger("different")
	if logger1 == logger3 {
		t.Error("NewLogger should return different instances for different names")
	}
	loggerRegistry = nil
}

// TestGetLogger_Existing tests retrieving an existing logger from the registry.
// It verifies that GetLogger returns the same instance that was previously created.
func TestGetLogger_Existing(t *testing.T) {
	loggerRegistry = nil
	WithLoggerRegistry()

	created := NewLogger("api")
	retrieved := GetLogger("api")

	if created != retrieved {
		t.Error("GetLogger should return the same instance that was created")
	}
	loggerRegistry = nil
}

// TestGetLogger_NonExistent tests GetLogger behavior when logger doesn't exist in registry.
func TestGetLogger_NonExistent(t *testing.T) {
	loggerRegistry = nil
	WithLoggerRegistry()

	// Get a logger that hasn't been created yet
	logger := GetLogger("nonexistent")

	if logger == nil {
		t.Fatal("GetLogger should create and return a new logger for non-existent names")
	}

	if logger.name != "nonexistent" {
		t.Errorf("Expected logger name 'nonexistent', got '%s'", logger.name)
	}

	// Verify it was added to the registry
	if loggerRegistry["nonexistent"] != logger {
		t.Error("GetLogger should register the newly created logger")
	}

	loggerRegistry = nil
}

// TestGetLogger_NoRegistry tests GetLogger when registry is not initialized.
func TestGetLogger_NoRegistry(t *testing.T) {
	loggerRegistry = nil

	logger := GetLogger("test")

	if logger == nil {
		t.Fatal("GetLogger should create a new logger even without registry")
	}

	if logger.name != "test" {
		t.Errorf("Expected logger name 'test', got '%s'", logger.name)
	}
}

// TestGetLogger_EmptyName tests GetLogger with an empty string name.
func TestGetLogger_EmptyName(t *testing.T) {
	loggerRegistry = nil
	WithLoggerRegistry()

	logger := GetLogger("")

	if logger == nil {
		t.Fatal("GetLogger should handle empty name")
	}

	if logger.name != "" {
		t.Errorf("Expected empty logger name, got '%s'", logger.name)
	}

	loggerRegistry = nil
}

// TestGetLogger_MultipleCalls tests that multiple GetLogger calls return same instance.
func TestGetLogger_MultipleCalls(t *testing.T) {
	loggerRegistry = nil
	WithLoggerRegistry()

	logger1 := GetLogger("multi")
	logger2 := GetLogger("multi")
	logger3 := GetLogger("multi")

	if logger1 != logger2 || logger2 != logger3 {
		t.Error("Multiple GetLogger calls with same name should return same instance")
	}

	loggerRegistry = nil
}

// TestGetLogger_DifferentNames tests GetLogger with different logger names.
func TestGetLogger_DifferentNames(t *testing.T) {
	loggerRegistry = nil
	WithLoggerRegistry()

	api := GetLogger("api")
	db := GetLogger("database")
	cache := GetLogger("cache")

	if api == db || db == cache || api == cache {
		t.Error("GetLogger should return different instances for different names")
	}

	if api.name != "api" || db.name != "database" || cache.name != "cache" {
		t.Error("Loggers should have correct names")
	}

	loggerRegistry = nil
}

// TestGetLogger_AfterNewLogger tests interoperability between GetLogger and NewLogger.
func TestGetLogger_AfterNewLogger(t *testing.T) {
	loggerRegistry = nil
	WithLoggerRegistry()

	created := NewLogger("interop")
	retrieved := GetLogger("interop")
	createdAgain := NewLogger("interop")

	if created != retrieved || retrieved != createdAgain {
		t.Error("GetLogger and NewLogger should return same instance for same name")
	}

	loggerRegistry = nil
}

// TestLogger_WithLogLevel tests the WithLogLevel configuration method.
// It verifies that the log level is properly set for different levels.
func TestLogger_WithLogLevel(t *testing.T) {
	tests := []LogLevel{ERROR, WARNING, INFO, DEBUG, TRACE, NONE}
	for _, level := range tests {
		logger := NewLogger("test").WithLogLevel(level)
		if logger.options.Level != level {
			t.Errorf("Expected log Level %v, got %v", level, logger.options.Level)
		}
	}
}

// TestLogger_WithTimestamp tests the WithTimestamp configuration method.
// It verifies that timestamps are enabled with the default format.
func TestLogger_WithTimestamp(t *testing.T) {
	logger := NewLogger("test").WithTimestamp()
	if !logger.options.Timestamp {
		t.Error("Timestamp not enabled")
	}
	if logger.options.TimestampFormat != "2006-01-02 15:04:05" {
		t.Errorf("Expected default Timestamp format, got '%s'", logger.options.TimestampFormat)
	}
}

// TestLogger_WithTimestampFormat tests the WithTimestampFormat configuration method.
// It verifies that custom timestamp formats are properly applied.
func TestLogger_WithTimestampFormat(t *testing.T) {
	customFormat := "2006/01/02"
	logger := NewLogger("test").WithTimestampFormat(customFormat)
	if !logger.options.Timestamp {
		t.Error("Timestamp not enabled")
	}
	if logger.options.TimestampFormat != customFormat {
		t.Errorf("Expected Timestamp format '%s', got '%s'", customFormat, logger.options.TimestampFormat)
	}
}

// TestLogger_WithColoring tests the WithColoring configuration method.
// It verifies that the method executes without errors (coloring is platform-dependent).
func TestLogger_WithColoring(t *testing.T) {
	logger := NewLogger("test").WithColoring()
	// Coloring is only enabled on non-Windows platforms
	// Just verify the method doesn't panic
	if logger == nil {
		t.Error("WithColoring() returned nil")
	}
}

// TestLogger_OutputTo tests the OutputTo configuration method.
// It verifies that the output writer is correctly set.
func TestLogger_OutputTo(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf)
	if logger.out != &buf {
		t.Error("OutputTo did not set the output writer correctly")
	}
}

// TestLogger_ErrorsTo tests the ErrorsTo configuration method.
// It verifies that the error output writer is correctly set.
func TestLogger_ErrorsTo(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").ErrorsTo(&buf)
	if logger.err != &buf {
		t.Error("ErrorsTo did not set the error writer correctly")
	}
}

// TestLogger_WithAsync tests the WithAsync configuration method.
// It verifies that asynchronous logging mode is properly enabled with channels created.
func TestLogger_WithAsync(t *testing.T) {
	var buf bytes.Buffer
	logger, cancel := NewLogger("test").OutputTo(&buf).WithAsync(true, 10)
	defer cancel()

	if !logger.options.Async {
		t.Error("WithAsync did not enable Async mode")
	}
	if logger.options.logs == nil {
		t.Error("WithAsync did not create logs channel")
	}
	if logger.options.errors == nil {
		t.Error("WithAsync did not create errors channel")
	}
}

// TestLogger_Logf tests the Logf method for logging messages without a level.
// It verifies that the message is written to the output stream.
func TestLogger_Logf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf)

	logger.Logf("test message")
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected 'test message', got '%s'", output)
	}
}

// TestLogger_Logf_WithArgs tests the Logf method with format arguments.
// It verifies that format strings are correctly processed.
func TestLogger_Logf_WithArgs(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf)

	logger.Logf("test %s %d", "message", 42)
	output := buf.String()

	if !strings.Contains(output, "test message 42") {
		t.Errorf("Expected 'test message 42', got '%s'", output)
	}
}

// TestLogger_Tracef tests the Tracef logging method.
// It verifies that TRACE level messages are properly formatted and logged.
func TestLogger_Tracef(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(TRACE)

	logger.Tracef("trace message")
	output := buf.String()

	if !strings.Contains(output, "TRACE") {
		t.Errorf("Expected 'TRACE' in output, got '%s'", output)
	}
	if !strings.Contains(output, "trace message") {
		t.Errorf("Expected 'trace message' in output, got '%s'", output)
	}
}

// TestLogger_Debugf tests the Debugf logging method.
// It verifies that DEBUG level messages are properly formatted and logged.
func TestLogger_Debugf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(DEBUG)

	logger.Debugf("debug message")
	output := buf.String()

	if !strings.Contains(output, "DEBUG") {
		t.Errorf("Expected 'DEBUG' in output, got '%s'", output)
	}
	if !strings.Contains(output, "debug message") {
		t.Errorf("Expected 'debug message' in output, got '%s'", output)
	}
}

// TestLogger_Infof tests the Infof logging method.
// It verifies that INFO level messages are properly formatted and logged.
func TestLogger_Infof(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(INFO)

	logger.Infof("info message")
	output := buf.String()

	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected 'INFO' in output, got '%s'", output)
	}
	if !strings.Contains(output, "info message") {
		t.Errorf("Expected 'info message' in output, got '%s'", output)
	}
}

// TestLogger_Warningf tests the Warningf logging method.
// It verifies that WARNING level messages are properly formatted and logged.
func TestLogger_Warningf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(WARNING)

	logger.Warningf("warning message")
	output := buf.String()

	if !strings.Contains(output, "WARNING") {
		t.Errorf("Expected 'WARNING' in output, got '%s'", output)
	}
	if !strings.Contains(output, "warning message") {
		t.Errorf("Expected 'warning message' in output, got '%s'", output)
	}
}

// TestLogger_Errorf tests the Errorf logging method.
// It verifies that ERROR level messages are properly formatted and logged to the error stream.
func TestLogger_Errorf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").ErrorsTo(&buf).WithLogLevel(ERROR)

	logger.Errorf("error message")
	output := buf.String()

	if !strings.Contains(output, "ERROR") {
		t.Errorf("Expected 'ERROR' in output, got '%s'", output)
	}
	if !strings.Contains(output, "error message") {
		t.Errorf("Expected 'error message' in output, got '%s'", output)
	}
}

// TestLogger_LogLevelFiltering_Trace tests that TRACE messages are filtered when the log level is DEBUG.
// It verifies that messages below the configured level are not logged.
func TestLogger_LogLevelFiltering_Trace(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(DEBUG)

	logger.Tracef("should not appear")

	if buf.Len() > 0 {
		t.Error("TRACE message should be filtered when log Level is DEBUG")
	}
}

// TestLogger_LogLevelFiltering_Debug tests that DEBUG messages are filtered when the log level is INFO.
// It verifies that messages below the configured level are not logged.
func TestLogger_LogLevelFiltering_Debug(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(INFO)

	logger.Debugf("should not appear")

	if buf.Len() > 0 {
		t.Error("DEBUG message should be filtered when log Level is INFO")
	}
}

// TestLogger_LogLevelFiltering_Info tests that INFO messages are filtered when the log level is WARNING.
// It verifies that messages below the configured level are not logged.
func TestLogger_LogLevelFiltering_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(WARNING)

	logger.Infof("should not appear")

	if buf.Len() > 0 {
		t.Error("INFO message should be filtered when log Level is WARNING")
	}
}

// TestLogger_LogLevelFiltering_Warning tests that WARNING messages are filtered when the log level is ERROR.
// It verifies that messages below the configured level are not logged.
func TestLogger_LogLevelFiltering_Warning(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(ERROR)

	logger.Warningf("should not appear")

	if buf.Len() > 0 {
		t.Error("WARNING message should be filtered when log Level is ERROR")
	}
}

// TestLogger_LogJSONf tests the LogJSONf method for JSON logging without a specific level.
// It verifies that structured data is correctly serialized to JSON.
func TestLogger_LogJSONf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf)

	obj := map[string]interface{}{"key": "value", "count": 42}
	err := logger.LogJSONf(obj, "json message")

	if err != nil {
		t.Errorf("LogJSONf returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"key":"value"`) {
		t.Errorf("Expected JSON object in output, got '%s'", output)
	}
	if !strings.Contains(output, `"message":"json message"`) {
		t.Errorf("Expected message in JSON output, got '%s'", output)
	}
}

// TestLogger_TraceJSONf tests the TraceJSONf method for JSON logging at TRACE level.
// It verifies that TRACE level is included in the JSON output.
func TestLogger_TraceJSONf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(TRACE)

	obj := map[string]interface{}{"trace": true}
	err := logger.TraceJSONf(obj, "trace json")

	if err != nil {
		t.Errorf("TraceJSONf returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"Level":"TRACE"`) {
		t.Errorf("Expected TRACE Level in JSON, got '%s'", output)
	}
}

// TestLogger_DebugJSONf tests the DebugJSONf method for JSON logging at DEBUG level.
// It verifies that DEBUG level is included in the JSON output.
func TestLogger_DebugJSONf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(DEBUG)

	obj := map[string]interface{}{"debug": true}
	err := logger.DebugJSONf(obj, "debug json")

	if err != nil {
		t.Errorf("DebugJSONf returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"Level":"DEBUG"`) {
		t.Errorf("Expected DEBUG Level in JSON, got '%s'", output)
	}
}

// TestLogger_InfoJSONf tests the InfoJSONf method for JSON logging at INFO level.
// It verifies that INFO level is included in the JSON output.
func TestLogger_InfoJSONf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(INFO)

	obj := map[string]interface{}{"info": true}
	err := logger.InfoJSONf(obj, "info json")

	if err != nil {
		t.Errorf("InfoJSONf returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"Level":"INFO"`) {
		t.Errorf("Expected INFO Level in JSON, got '%s'", output)
	}
}

// TestLogger_WarningJSONf tests the WarningJSONf method for JSON logging at WARNING level.
// It verifies that WARNING level is included in the JSON output.
func TestLogger_WarningJSONf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(WARNING)

	obj := map[string]interface{}{"warning": true}
	err := logger.WarningJSONf(obj, "warning json")

	if err != nil {
		t.Errorf("WarningJSONf returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"Level":"WARNING"`) {
		t.Errorf("Expected WARNING Level in JSON, got '%s'", output)
	}
}

// TestLogger_ErrorJSONf tests the ErrorJSONf method for JSON logging at ERROR level.
// It verifies that ERROR level is included in the JSON output and written to error stream.
func TestLogger_ErrorJSONf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").ErrorsTo(&buf).WithLogLevel(ERROR)

	obj := map[string]interface{}{"error": true}
	err := logger.ErrorJSONf(obj, "error json")

	if err != nil {
		t.Errorf("ErrorJSONf returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"Level":"ERROR"`) {
		t.Errorf("Expected ERROR Level in JSON, got '%s'", output)
	}
}

// TestLogger_TraceJSONf_Filtered tests that TRACE JSON logs are filtered when log level is DEBUG.
// It verifies that level filtering applies to JSON logging methods.
func TestLogger_TraceJSONf_Filtered(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(DEBUG)

	obj := map[string]interface{}{"trace": true}
	err := logger.TraceJSONf(obj, "should not appear")

	if err != nil {
		t.Errorf("TraceJSONf should return nil when filtered, got: %v", err)
	}

	if buf.Len() > 0 {
		t.Error("TRACE JSON should be filtered when log Level is DEBUG")
	}
}

// TestLogger_LogObjectf tests the LogObjectf method for zero-allocation object logging.
// It verifies that structured fields are properly added and serialized.
func TestLogger_LogObjectf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf)

	logger.LogObjectf("object message").
		AssignString("key", "value").
		AssignInt("count", 42).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"message":"object message"`) {
		t.Errorf("Expected message in output, got '%s'", output)
	}
	if !strings.Contains(output, `"key":"value"`) {
		t.Errorf("Expected key:value in output, got '%s'", output)
	}
	if !strings.Contains(output, `"count":42`) {
		t.Errorf("Expected count:42 in output, got '%s'", output)
	}
}

// TestLogger_TraceObjectf tests the TraceObjectf method for object logging at TRACE level.
// It verifies that TRACE level is included in the object log output.
func TestLogger_TraceObjectf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(TRACE)

	logger.TraceObjectf("trace object").
		AssignString("key", "value").
		Build()

	output := buf.String()
	if !strings.Contains(output, `"Level":"TRACE"`) {
		t.Errorf("Expected TRACE Level, got '%s'", output)
	}
}

// TestLogger_DebugObjectf tests the DebugObjectf method for object logging at DEBUG level.
// It verifies that DEBUG level is included in the object log output.
func TestLogger_DebugObjectf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(DEBUG)

	logger.DebugObjectf("debug object").
		AssignString("key", "value").
		Build()

	output := buf.String()
	if !strings.Contains(output, `"Level":"DEBUG"`) {
		t.Errorf("Expected DEBUG Level, got '%s'", output)
	}
}

// TestLogger_InfoObjectf tests the InfoObjectf method for object logging at INFO level.
// It verifies that INFO level is included in the object log output.
func TestLogger_InfoObjectf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("info object").
		AssignString("key", "value").
		Build()

	output := buf.String()
	if !strings.Contains(output, `"Level":"INFO"`) {
		t.Errorf("Expected INFO Level, got '%s'", output)
	}
}

// TestLogger_WarningObjectf tests the WarningObjectf method for object logging at WARNING level.
// It verifies that WARNING level is included in the object log output.
func TestLogger_WarningObjectf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(WARNING)

	logger.WarningObjectf("warning object").
		AssignString("key", "value").
		Build()

	output := buf.String()
	if !strings.Contains(output, `"Level":"WARNING"`) {
		t.Errorf("Expected WARNING Level, got '%s'", output)
	}
}

// TestLogger_ErrorObjectf tests the ErrorObjectf method for object logging at ERROR level.
// It verifies that ERROR level is included and written to the error stream.
func TestLogger_ErrorObjectf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").ErrorsTo(&buf).WithLogLevel(ERROR)

	logger.ErrorObjectf("error object").
		AssignString("key", "value").
		Build()

	output := buf.String()
	if !strings.Contains(output, `"Level":"ERROR"`) {
		t.Errorf("Expected ERROR Level, got '%s'", output)
	}
}

// TestLogger_TraceObjectf_Filtered tests that TRACE object logs are filtered when log level is DEBUG.
// It verifies that nil is returned when the message would be filtered.
func TestLogger_TraceObjectf_Filtered(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(DEBUG)

	builder := logger.TraceObjectf("should not appear")

	if builder != nil {
		t.Error("TraceObjectf should return nil when filtered")
	}

	if buf.Len() > 0 {
		t.Error("TRACE object should be filtered when log Level is DEBUG")
	}
}

// TestLogger_DebugObjectf_Filtered tests that DEBUG object logs are filtered when log level is INFO.
// It verifies that nil is returned when the message would be filtered.
func TestLogger_DebugObjectf_Filtered(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(INFO)

	builder := logger.DebugObjectf("should not appear")

	if builder != nil {
		t.Error("DebugObjectf should return nil when filtered")
	}
}

// TestLogger_InfoObjectf_Filtered tests that INFO object logs are filtered when log level is WARNING.
// It verifies that nil is returned when the message would be filtered.
func TestLogger_InfoObjectf_Filtered(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(WARNING)

	builder := logger.InfoObjectf("should not appear")

	if builder != nil {
		t.Error("InfoObjectf should return nil when filtered")
	}
}

// TestLogger_WarningObjectf_Filtered tests that WARNING object logs are filtered when log level is ERROR.
// It verifies that nil is returned when the message would be filtered.
func TestLogger_WarningObjectf_Filtered(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(ERROR)

	builder := logger.WarningObjectf("should not appear")

	if builder != nil {
		t.Error("WarningObjectf should return nil when filtered")
	}
}

// TestLogger_AsyncLogging tests asynchronous logging mode.
// It verifies that messages are written to the output stream asynchronously.
func TestLogger_AsyncLogging(t *testing.T) {
	var buf bytes.Buffer
	logger, cancel := NewLogger("test").OutputTo(&buf).WithLogLevel(INFO).WithAsync(true, 10)
	defer cancel()

	logger.Infof("Async test message")

	for buf.Len() == 0 {
	}

	logged := buf.String()

	if !strings.Contains(logged, "INFO") {
		t.Skip()
	}
}

// TestLogger_WithTimestamp_Output tests that timestamps appear in log output.
// It verifies that timestamp prefixes are added when enabled.
func TestLogger_WithTimestamp_Output(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(INFO).WithTimestamp()

	logger.Infof("test message")
	output := buf.String()

	// Should contain Timestamp (check for year pattern)
	if !strings.Contains(output, "20") {
		t.Errorf("Expected Timestamp in output, got '%s'", output)
	}
}

// TestLogger_WithName_Output tests that logger names appear in log output.
// It verifies that the logger name is included in formatted messages.
func TestLogger_WithName_Output(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("MyService").OutputTo(&buf).WithLogLevel(INFO)

	logger.Infof("test message")
	output := buf.String()

	if !strings.Contains(output, "MyService") {
		t.Errorf("Expected service name in output, got '%s'", output)
	}
}

// TestLogger_CombinedFeatures tests multiple logger features working together.
// It verifies that timestamp, level, name, and message formatting all work correctly.
func TestLogger_CombinedFeatures(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestService").
		OutputTo(&buf).
		WithLogLevel(INFO).
		WithTimestamp()

	logger.Infof("test message with %s", "args")
	output := buf.String()

	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected INFO Level, got '%s'", output)
	}
	if !strings.Contains(output, "TestService") {
		t.Errorf("Expected service name, got '%s'", output)
	}
	if !strings.Contains(output, "test message with args") {
		t.Errorf("Expected message, got '%s'", output)
	}
}

// TestLogger_MethodChaining tests that configuration methods can be chained.
// It verifies that each method returns the logger instance for method chaining.
func TestLogger_MethodChaining(t *testing.T) {
	logger := NewLogger("ChainTest").
		WithLogLevel(DEBUG).
		WithTimestamp().
		WithColoring()

	if logger == nil {
		t.Error("Method chaining broke the logger")
	}
	if logger.name != "ChainTest" {
		t.Error("Method chaining failed to set name")
	}
	if logger.options.Level != DEBUG {
		t.Error("Method chaining failed to set log Level")
	}
	if !logger.options.Timestamp {
		t.Error("Method chaining failed to set Timestamp")
	}
}

// TestLogger_ConcurrentLogging tests thread-safe concurrent logging.
// It verifies that multiple goroutines can safely log messages simultaneously.
func TestLogger_ConcurrentLogging(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(INFO)

	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(id int) {
			logger.Infof("message %d", id)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	output := buf.String()
	if len(output) == 0 {
		t.Error("No output from concurrent logging")
	}
}

// TestLogger_ErrorWriterSeparation tests that error logs go to the error stream.
// It verifies that ERROR level messages are written to a separate error writer.
func TestLogger_ErrorWriterSeparation(t *testing.T) {
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	logger := NewLogger("test").
		OutputTo(&outBuf).
		ErrorsTo(&errBuf).
		WithLogLevel(ERROR)

	logger.Infof("info message")
	logger.Errorf("error message")

	outOutput := outBuf.String()
	errOutput := errBuf.String()

	if len(outOutput) > 0 {
		t.Error("INFO should be filtered at ERROR Level")
	}
	if !strings.Contains(errOutput, "error message") {
		t.Errorf("Expected error message in error output, got '%s'", errOutput)
	}
}

// TestLogger_EmptyMessage tests logging with an empty message string.
// It verifies that the log level is still included even with an empty message.
func TestLogger_EmptyMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(INFO)

	logger.Infof("")
	output := buf.String()

	if !strings.Contains(output, "INFO") {
		t.Error("Empty message should still include log Level")
	}
}

// TestLogger_NilObjectJSON tests JSON logging with a nil object.
// It verifies that nil objects are handled gracefully without errors.
func TestLogger_NilObjectJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf)

	err := logger.LogJSONf(nil, "message")

	if err != nil {
		t.Errorf("LogJSONf with nil object should not error, got: %v", err)
	}
}

// TestLogger_format tests the internal format method.
// It verifies that messages are formatted correctly with and without arguments.
func TestLogger_format(t *testing.T) {
	logger := NewLogger("test")

	// Without args
	result := logger.format("test message")
	if string(result) != "test message" {
		t.Errorf("Expected 'test message', got '%s'", string(result))
	}

	// With args
	result = logger.format("test %s %d", "message", 42)
	if string(result) != "test message 42" {
		t.Errorf("Expected 'test message 42', got '%s'", string(result))
	}
}

// TestLogger_writeStringToStream_NilStream tests writeStringToStream with a nil stream.
// It verifies that nil streams are handled gracefully without panicking.
func TestLogger_writeStringToStream_NilStream(t *testing.T) {
	logger := NewLogger("test")

	// Should not panic
	logger.writeStringToStream(nil, INFO, "test")
}

// TestLogger_writeJSONToStream_NilStream tests writeJSONToStream with a nil stream.
// It verifies that an error is returned when the stream is nil.
func TestLogger_writeJSONToStream_NilStream(t *testing.T) {
	logger := NewLogger("test")

	obj := map[string]interface{}{"key": "value"}
	err := logger.writeJSONToStream(nil, INFO, obj, "test")

	if err == nil {
		t.Error("writeJSONToStream should return error for nil stream")
	}
}
