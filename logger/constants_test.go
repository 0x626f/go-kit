package logger

import (
	"testing"
)

// TestLogLevel_String tests the String method of LogLevel type.
// It verifies that each log level returns the correct string representation.
func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{ERROR, "ERROR"},
		{WARNING, "WARNING"},
		{INFO, "INFO"},
		{DEBUG, "DEBUG"},
		{TRACE, "TRACE"},
		{NONE, ""},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.level.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

// TestLogLevel_color tests the color method of LogLevel type.
// It verifies that each log level returns the correct ANSI color code.
func TestLogLevel_color(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected color
	}{
		{ERROR, colorRed},
		{WARNING, colorYellow},
		{INFO, colorGreen},
		{DEBUG, colorGrey},
		{TRACE, colorBlue},
		{NONE, []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			result := tt.level.color()
			if len(result) != len(tt.expected) {
				t.Errorf("Color length mismatch for %s: expected %d, got %d",
					tt.level.String(), len(tt.expected), len(result))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("Color mismatch for %s at index %d: expected %v, got %v",
						tt.level.String(), i, tt.expected[i], result[i])
				}
			}
		})
	}
}

// TestLogLevel_paint tests the paint method of LogLevel type.
// It verifies that log messages are correctly wrapped with ANSI color codes.
func TestLogLevel_paint(t *testing.T) {
	payload := []byte("test message")

	tests := []struct {
		level LogLevel
		name  string
	}{
		{ERROR, "ERROR"},
		{WARNING, "WARNING"},
		{INFO, "INFO"},
		{DEBUG, "DEBUG"},
		{TRACE, "TRACE"},
		{NONE, "NONE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.level.paint(payload)

			// For NONE, the color is empty, so it's color + payload + reset
			levelColor := tt.level.color()
			expectedLen := len(levelColor) + len(payload) + len(colorReset)
			if len(result) != expectedLen {
				t.Errorf("Expected painted length %d, got %d", expectedLen, len(result))
			}

			// If the Level has a color, should start with it
			if len(levelColor) > 0 {
				for i := range levelColor {
					if result[i] != levelColor[i] {
						t.Errorf("Paint should start with Level color")
						break
					}
				}
			}

			// Should end with color reset
			resetStart := len(result) - len(colorReset)
			for i := range colorReset {
				if result[resetStart+i] != colorReset[i] {
					t.Errorf("Paint should end with color reset")
					break
				}
			}
		})
	}
}

// TestColorConstants tests that all ANSI color constant values are properly defined.
// It ensures none of the color escape sequences are empty (except where expected).
func TestColorConstants(t *testing.T) {
	// Verify color constants are not empty (except for edge cases)
	if len(colorReset) == 0 {
		t.Error("colorReset should not be empty")
	}
	if len(colorRed) == 0 {
		t.Error("colorRed should not be empty")
	}
	if len(colorGreen) == 0 {
		t.Error("colorGreen should not be empty")
	}
	if len(colorYellow) == 0 {
		t.Error("colorYellow should not be empty")
	}
	if len(colorBlue) == 0 {
		t.Error("colorBlue should not be empty")
	}
	if len(colorGrey) == 0 {
		t.Error("colorGrey should not be empty")
	}
}

// TestLogLevel_Ordering tests that log levels are ordered correctly by severity.
// Lower numeric values should represent higher severity (ERROR < WARNING < INFO < DEBUG < TRACE < NONE).
func TestLogLevel_Ordering(t *testing.T) {
	// Verify log levels are in correct order
	if ERROR >= WARNING {
		t.Error("ERROR should be less than WARNING")
	}
	if WARNING >= INFO {
		t.Error("WARNING should be less than INFO")
	}
	if INFO >= DEBUG {
		t.Error("INFO should be less than DEBUG")
	}
	if DEBUG >= TRACE {
		t.Error("DEBUG should be less than TRACE")
	}
	if TRACE >= NONE {
		t.Error("TRACE should be less than NONE")
	}
}

// TestConstants tests that package-level constants (space, ln) have correct values.
// These constants are used internally for formatting log messages.
func TestConstants(t *testing.T) {
	if space != byte(' ') {
		t.Errorf("Expected space to be ' ', got %v", space)
	}
	if ln != byte('\n') {
		t.Errorf("Expected ln to be '\\n', got %v", ln)
	}
}

// TestParseLogLevel tests the ParseLogLevel function with various inputs.
// It verifies case-insensitive parsing and proper handling of invalid inputs.
func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected LogLevel
	}{
		{"Uppercase ERROR", "ERROR", ERROR},
		{"Lowercase error", "error", ERROR},
		{"Mixed case Error", "Error", ERROR},
		{"Uppercase WARNING", "WARNING", WARNING},
		{"Lowercase warning", "warning", WARNING},
		{"Mixed case Warning", "Warning", WARNING},
		{"Uppercase INFO", "INFO", INFO},
		{"Lowercase info", "info", INFO},
		{"Mixed case Info", "Info", INFO},
		{"Uppercase DEBUG", "DEBUG", DEBUG},
		{"Lowercase debug", "debug", DEBUG},
		{"Mixed case Debug", "Debug", DEBUG},
		{"Uppercase TRACE", "TRACE", TRACE},
		{"Lowercase trace", "trace", TRACE},
		{"Mixed case Trace", "Trace", TRACE},
		{"Invalid Level", "invalid", NONE},
		{"Empty string", "", NONE},
		{"Random string", "xyz123", NONE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseLogLevel(tt.input)
			if result != tt.expected {
				t.Errorf("ParseLogLevel(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestLogLevel_UnmarshalText tests the UnmarshalText implementation for LogLevel.
// This is used for parsing log levels from environment variables and configuration files.
func TestLogLevel_UnmarshalText(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected LogLevel
	}{
		{"Uppercase ERROR", []byte("ERROR"), ERROR},
		{"Lowercase error", []byte("error"), ERROR},
		{"Mixed case Error", []byte("Error"), ERROR},
		{"Uppercase WARNING", []byte("WARNING"), WARNING},
		{"Lowercase warning", []byte("warning"), WARNING},
		{"Mixed case Warning", []byte("Warning"), WARNING},
		{"Uppercase INFO", []byte("INFO"), INFO},
		{"Lowercase info", []byte("info"), INFO},
		{"Mixed case Info", []byte("Info"), INFO},
		{"Uppercase DEBUG", []byte("DEBUG"), DEBUG},
		{"Lowercase debug", []byte("debug"), DEBUG},
		{"Mixed case Debug", []byte("Debug"), DEBUG},
		{"Uppercase TRACE", []byte("TRACE"), TRACE},
		{"Lowercase trace", []byte("trace"), TRACE},
		{"Mixed case Trace", []byte("Trace"), TRACE},
		{"Invalid Level", []byte("invalid"), NONE},
		{"Empty bytes", []byte(""), NONE},
		{"Random bytes", []byte("xyz123"), NONE},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var level LogLevel
			err := level.UnmarshalText(tt.input)
			if err != nil {
				t.Errorf("UnmarshalText(%q) returned error: %v", tt.input, err)
			}
			if level != tt.expected {
				t.Errorf("UnmarshalText(%q) = %v, expected %v", tt.input, level, tt.expected)
			}
		})
	}
}

// TestLogLevel_UnmarshalText_AllLevels verifies that all valid log levels can be unmarshaled correctly.
func TestLogLevel_UnmarshalText_AllLevels(t *testing.T) {
	levels := []LogLevel{ERROR, WARNING, INFO, DEBUG, TRACE}

	for _, expectedLevel := range levels {
		t.Run(expectedLevel.String(), func(t *testing.T) {
			var level LogLevel
			err := level.UnmarshalText([]byte(expectedLevel.String()))
			if err != nil {
				t.Errorf("UnmarshalText failed for %s: %v", expectedLevel.String(), err)
			}
			if level != expectedLevel {
				t.Errorf("Expected level %v, got %v", expectedLevel, level)
			}
		})
	}
}
