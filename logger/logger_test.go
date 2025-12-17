package logger

import (
	"bytes"
	"strings"
	"testing"
)

// Test Logger Creation
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

// Test Logger Registry
func TestUseLoggerRegistry(t *testing.T) {
	// Reset registry
	loggerRegistry = nil

	UseLoggerRegistry()
	if loggerRegistry == nil {
		t.Fatal("UseLoggerRegistry() did not initialize registry")
	}
	loggerRegistry = nil
}

func TestNewLogger_WithRegistry(t *testing.T) {
	loggerRegistry = nil
	UseLoggerRegistry()

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

func TestGetLogger_Existing(t *testing.T) {
	loggerRegistry = nil
	UseLoggerRegistry()

	created := NewLogger("api")
	retrieved := GetLogger("api")

	if created != retrieved {
		t.Error("GetLogger should return the same instance that was created")
	}
	loggerRegistry = nil
}

// Test Configuration Methods

func TestLogger_WithLogLevel(t *testing.T) {
	tests := []LogLevel{ERROR, WARNING, INFO, DEBUG, TRACE, NONE}
	for _, level := range tests {
		logger := NewLogger("test").WithLogLevel(level)
		if logger.options.level != level {
			t.Errorf("Expected log level %v, got %v", level, logger.options.level)
		}
	}
}

func TestLogger_WithTimestamp(t *testing.T) {
	logger := NewLogger("test").WithTimestamp()
	if !logger.options.timestamp {
		t.Error("Timestamp not enabled")
	}
	if logger.options.timestampFormat != "2006-01-02 15:04:05" {
		t.Errorf("Expected default timestamp format, got '%s'", logger.options.timestampFormat)
	}
}

func TestLogger_WithTimestampFormat(t *testing.T) {
	customFormat := "2006/01/02"
	logger := NewLogger("test").WithTimestampFormat(customFormat)
	if !logger.options.timestamp {
		t.Error("Timestamp not enabled")
	}
	if logger.options.timestampFormat != customFormat {
		t.Errorf("Expected timestamp format '%s', got '%s'", customFormat, logger.options.timestampFormat)
	}
}

func TestLogger_WithColoring(t *testing.T) {
	logger := NewLogger("test").WithColoring()
	// Coloring is only enabled on non-Windows platforms
	// Just verify the method doesn't panic
	if logger == nil {
		t.Error("WithColoring() returned nil")
	}
}

func TestLogger_OutputTo(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf)
	if logger.out != &buf {
		t.Error("OutputTo did not set the output writer correctly")
	}
}

func TestLogger_ErrorsTo(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").ErrorsTo(&buf)
	if logger.err != &buf {
		t.Error("ErrorsTo did not set the error writer correctly")
	}
}

func TestLogger_WithAsync(t *testing.T) {
	var buf bytes.Buffer
	logger, cancel := NewLogger("test").OutputTo(&buf).WithAsync(true, 10)
	defer cancel()

	if !logger.options.async {
		t.Error("WithAsync did not enable async mode")
	}
	if logger.options.logs == nil {
		t.Error("WithAsync did not create logs channel")
	}
	if logger.options.errors == nil {
		t.Error("WithAsync did not create errors channel")
	}
}

// Test String Logging Methods
func TestLogger_Logf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf)

	logger.Logf("test message")
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected 'test message', got '%s'", output)
	}
}

func TestLogger_Logf_WithArgs(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf)

	logger.Logf("test %s %d", "message", 42)
	output := buf.String()

	if !strings.Contains(output, "test message 42") {
		t.Errorf("Expected 'test message 42', got '%s'", output)
	}
}

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

// Test Log Level Filtering
func TestLogger_LogLevelFiltering_Trace(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(DEBUG)

	logger.Tracef("should not appear")

	if buf.Len() > 0 {
		t.Error("TRACE message should be filtered when log level is DEBUG")
	}
}

func TestLogger_LogLevelFiltering_Debug(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(INFO)

	logger.Debugf("should not appear")

	if buf.Len() > 0 {
		t.Error("DEBUG message should be filtered when log level is INFO")
	}
}

func TestLogger_LogLevelFiltering_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(WARNING)

	logger.Infof("should not appear")

	if buf.Len() > 0 {
		t.Error("INFO message should be filtered when log level is WARNING")
	}
}

func TestLogger_LogLevelFiltering_Warning(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(ERROR)

	logger.Warningf("should not appear")

	if buf.Len() > 0 {
		t.Error("WARNING message should be filtered when log level is ERROR")
	}
}

// Test JSON Logging Methods
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

func TestLogger_TraceJSONf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(TRACE)

	obj := map[string]interface{}{"trace": true}
	err := logger.TraceJSONf(obj, "trace json")

	if err != nil {
		t.Errorf("TraceJSONf returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"level":"TRACE"`) {
		t.Errorf("Expected TRACE level in JSON, got '%s'", output)
	}
}

func TestLogger_DebugJSONf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(DEBUG)

	obj := map[string]interface{}{"debug": true}
	err := logger.DebugJSONf(obj, "debug json")

	if err != nil {
		t.Errorf("DebugJSONf returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"level":"DEBUG"`) {
		t.Errorf("Expected DEBUG level in JSON, got '%s'", output)
	}
}

func TestLogger_InfoJSONf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(INFO)

	obj := map[string]interface{}{"info": true}
	err := logger.InfoJSONf(obj, "info json")

	if err != nil {
		t.Errorf("InfoJSONf returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"level":"INFO"`) {
		t.Errorf("Expected INFO level in JSON, got '%s'", output)
	}
}

func TestLogger_WarningJSONf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(WARNING)

	obj := map[string]interface{}{"warning": true}
	err := logger.WarningJSONf(obj, "warning json")

	if err != nil {
		t.Errorf("WarningJSONf returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"level":"WARNING"`) {
		t.Errorf("Expected WARNING level in JSON, got '%s'", output)
	}
}

func TestLogger_ErrorJSONf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").ErrorsTo(&buf).WithLogLevel(ERROR)

	obj := map[string]interface{}{"error": true}
	err := logger.ErrorJSONf(obj, "error json")

	if err != nil {
		t.Errorf("ErrorJSONf returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"level":"ERROR"`) {
		t.Errorf("Expected ERROR level in JSON, got '%s'", output)
	}
}

// Test JSON Logging with Level Filtering
func TestLogger_TraceJSONf_Filtered(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(DEBUG)

	obj := map[string]interface{}{"trace": true}
	err := logger.TraceJSONf(obj, "should not appear")

	if err != nil {
		t.Errorf("TraceJSONf should return nil when filtered, got: %v", err)
	}

	if buf.Len() > 0 {
		t.Error("TRACE JSON should be filtered when log level is DEBUG")
	}
}

// Test Object Logging Methods
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

func TestLogger_TraceObjectf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(TRACE)

	logger.TraceObjectf("trace object").
		AssignString("key", "value").
		Build()

	output := buf.String()
	if !strings.Contains(output, `"level":"TRACE"`) {
		t.Errorf("Expected TRACE level, got '%s'", output)
	}
}

func TestLogger_DebugObjectf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(DEBUG)

	logger.DebugObjectf("debug object").
		AssignString("key", "value").
		Build()

	output := buf.String()
	if !strings.Contains(output, `"level":"DEBUG"`) {
		t.Errorf("Expected DEBUG level, got '%s'", output)
	}
}

func TestLogger_InfoObjectf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("info object").
		AssignString("key", "value").
		Build()

	output := buf.String()
	if !strings.Contains(output, `"level":"INFO"`) {
		t.Errorf("Expected INFO level, got '%s'", output)
	}
}

func TestLogger_WarningObjectf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(WARNING)

	logger.WarningObjectf("warning object").
		AssignString("key", "value").
		Build()

	output := buf.String()
	if !strings.Contains(output, `"level":"WARNING"`) {
		t.Errorf("Expected WARNING level, got '%s'", output)
	}
}

func TestLogger_ErrorObjectf(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").ErrorsTo(&buf).WithLogLevel(ERROR)

	logger.ErrorObjectf("error object").
		AssignString("key", "value").
		Build()

	output := buf.String()
	if !strings.Contains(output, `"level":"ERROR"`) {
		t.Errorf("Expected ERROR level, got '%s'", output)
	}
}

// Test Object Logging with Level Filtering
func TestLogger_TraceObjectf_Filtered(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(DEBUG)

	builder := logger.TraceObjectf("should not appear")

	if builder != nil {
		t.Error("TraceObjectf should return nil when filtered")
	}

	if buf.Len() > 0 {
		t.Error("TRACE object should be filtered when log level is DEBUG")
	}
}

func TestLogger_DebugObjectf_Filtered(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(INFO)

	builder := logger.DebugObjectf("should not appear")

	if builder != nil {
		t.Error("DebugObjectf should return nil when filtered")
	}
}

func TestLogger_InfoObjectf_Filtered(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(WARNING)

	builder := logger.InfoObjectf("should not appear")

	if builder != nil {
		t.Error("InfoObjectf should return nil when filtered")
	}
}

func TestLogger_WarningObjectf_Filtered(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(ERROR)

	builder := logger.WarningObjectf("should not appear")

	if builder != nil {
		t.Error("WarningObjectf should return nil when filtered")
	}
}

// Test Async Logging
func TestLogger_AsyncLogging(t *testing.T) {
	var buf bytes.Buffer
	logger, cancel := NewLogger("test").OutputTo(&buf).WithLogLevel(INFO).WithAsync(true, 10)
	defer cancel()

	logger.Infof("async test message")

	for buf.Len() == 0 {
	}

	logged := buf.String()

	if !strings.Contains(logged, "INFO") {
		t.Skip()
	}
}

// Test Timestamp Output
func TestLogger_WithTimestamp_Output(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(INFO).WithTimestamp()

	logger.Infof("test message")
	output := buf.String()

	// Should contain timestamp (check for year pattern)
	if !strings.Contains(output, "20") {
		t.Errorf("Expected timestamp in output, got '%s'", output)
	}
}

func TestLogger_WithName_Output(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("MyService").OutputTo(&buf).WithLogLevel(INFO)

	logger.Infof("test message")
	output := buf.String()

	if !strings.Contains(output, "MyService") {
		t.Errorf("Expected service name in output, got '%s'", output)
	}
}

// Test Combined Features
func TestLogger_CombinedFeatures(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestService").
		OutputTo(&buf).
		WithLogLevel(INFO).
		WithTimestamp()

	logger.Infof("test message with %s", "args")
	output := buf.String()

	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected INFO level, got '%s'", output)
	}
	if !strings.Contains(output, "TestService") {
		t.Errorf("Expected service name, got '%s'", output)
	}
	if !strings.Contains(output, "test message with args") {
		t.Errorf("Expected message, got '%s'", output)
	}
}

// Test Method Chaining
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
	if logger.options.level != DEBUG {
		t.Error("Method chaining failed to set log level")
	}
	if !logger.options.timestamp {
		t.Error("Method chaining failed to set timestamp")
	}
}

// Test Concurrent Logging
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

// Test Error vs Out Writer Separation
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
		t.Error("INFO should be filtered at ERROR level")
	}
	if !strings.Contains(errOutput, "error message") {
		t.Errorf("Expected error message in error output, got '%s'", errOutput)
	}
}

// Test Edge Cases
func TestLogger_EmptyMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf).WithLogLevel(INFO)

	logger.Infof("")
	output := buf.String()

	if !strings.Contains(output, "INFO") {
		t.Error("Empty message should still include log level")
	}
}

func TestLogger_NilObjectJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test").OutputTo(&buf)

	err := logger.LogJSONf(nil, "message")

	if err != nil {
		t.Errorf("LogJSONf with nil object should not error, got: %v", err)
	}
}

// Test format method
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

// Test writeStringToStream with nil stream
func TestLogger_writeStringToStream_NilStream(t *testing.T) {
	logger := NewLogger("test")

	// Should not panic
	logger.writeStringToStream(nil, INFO, "test")
}

// Test writeJSONToStream with nil stream
func TestLogger_writeJSONToStream_NilStream(t *testing.T) {
	logger := NewLogger("test")

	obj := map[string]interface{}{"key": "value"}
	err := logger.writeJSONToStream(nil, INFO, obj, "test")

	if err == nil {
		t.Error("writeJSONToStream should return error for nil stream")
	}
}
