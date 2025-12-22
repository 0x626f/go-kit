package logger

import "strings"

// LogLevel represents the severity of a log message.
// Lower numeric values indicate higher severity.
// Levels in order from highest to lowest severity: ERROR < WARNING < INFO < DEBUG < TRACE < NONE
//
// When a logger's level is set to a particular value, only messages at that level or higher severity are logged.
// For example, if set to INFO, ERROR and WARNING messages are also logged, but DEBUG and TRACE are discarded.
type LogLevel uint8

// color is an internal type representing ANSI color escape sequences.
// These are byte sequences that terminals interpret as formatting commands.
type color []byte

const (
	// ERROR indicates a serious problem that should be addressed immediately.
	// Errors are written to stderr by default.
	ERROR LogLevel = iota
	// WARNING indicates a potentially harmful situation that should be reviewed.
	WARNING
	// INFO indicates general informational messages about application operation.
	INFO
	// DEBUG indicates detailed information useful for development and troubleshooting.
	DEBUG
	// TRACE indicates the most verbose level of diagnostic information.
	TRACE
	// NONE indicates no log level filtering - all messages are logged.
	// This is also used for messages that don't have an explicit level.
	NONE

	// space is the ASCII space character used as a delimiter in log formatting
	space = byte(' ')
	// ln is the ASCII newline character appended to all log messages
	ln = byte('\n')
)

var (
	// ANSI color escape sequences for terminal output
	// These codes work on Unix/Linux/macOS terminals but not Windows CMD

	// colorReset restores the terminal to default colors
	colorReset color = []byte{27, 91, 48, 109}
	// colorRed is used for ERROR level messages
	colorRed color = []byte{27, 91, 51, 49, 109}
	// colorGreen is used for INFO level messages
	colorGreen color = []byte{27, 91, 51, 50, 109}
	// colorYellow is used for WARNING level messages
	colorYellow color = []byte{27, 91, 51, 51, 109}
	// colorBlue is used for TRACE level messages
	colorBlue color = []byte{27, 91, 51, 52, 109}
	// colorGrey is used for DEBUG level messages
	colorGrey color = []byte{27, 91, 57, 48, 109}
)

// String returns the string representation of the log level.
//
// Returns:
//   - "ERROR", "WARNING", "INFO", "DEBUG", or "TRACE" for known levels
//   - "" (empty string) for NONE or unrecognized levels
//
// Example:
//
//	level := logger.INFO
//	fmt.Println(level.String()) // Output: "INFO"
func (level *LogLevel) String() string {
	switch *level {
	case ERROR:
		return "ERROR"
	case WARNING:
		return "WARNING"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	default:
		return ""
	}
}

// ParseLogLevel converts a string representation to a LogLevel.
// The comparison is case-insensitive, so "error", "ERROR", and "Error" all return ERROR.
//
// Parameters:
//   - level: A string representation of the log level
//
// Returns:
//   - The corresponding LogLevel constant
//   - NONE for unrecognized or empty strings
//
// Example:
//
//	level := logger.ParseLogLevel("info")    // Returns INFO
//	level := logger.ParseLogLevel("ERROR")   // Returns ERROR
//	level := logger.ParseLogLevel("invalid") // Returns NONE
func ParseLogLevel(level string) LogLevel {
	switch {
	case strings.EqualFold(level, "ERROR"):
		return ERROR
	case strings.EqualFold(level, "WARNING"):
		return WARNING
	case strings.EqualFold(level, "INFO"):
		return INFO
	case strings.EqualFold(level, "DEBUG"):
		return DEBUG
	case strings.EqualFold(level, "TRACE"):
		return TRACE
	default:
		return NONE
	}
}

func (level *LogLevel) UnmarshalText(text []byte) error {
	*level = ParseLogLevel(string(text))
	return nil
}

// color returns the ANSI color code for this log level.
// Used internally when coloring is enabled.
//
// Returns:
//   - The ANSI color escape sequence for this level
//   - Empty byte slice for NONE or unrecognized levels
func (level *LogLevel) color() color {
	switch *level {
	case ERROR:
		return colorRed
	case WARNING:
		return colorYellow
	case INFO:
		return colorGreen
	case DEBUG:
		return colorGrey
	case TRACE:
		return colorBlue
	default:
		return []byte{}
	}
}

// paint wraps the payload with ANSI color codes for this log level.
// The payload is prefixed with the level's color and suffixed with a reset code.
//
// Parameters:
//   - payload: The log message bytes to colorize
//
// Returns:
//   - The payload wrapped in ANSI color escape sequences
//
// Example output for ERROR level:
//
//	Input:  []byte("Error message")
//	Output: [ESC[31m]Error message[ESC[0m]  (displays in red)
func (level *LogLevel) paint(payload []byte) []byte {
	return append(append(level.color(), payload...), colorReset...)
}
