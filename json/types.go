package json

import (
	"math"
	"strconv"
)

// AppendString appends a JSON string value (quoted).
// The value is wrapped in double quotes as required by JSON specification.
//
// Parameters:
//   - value: The string to append
//
// Returns:
//   - The encoder for method chaining
//
// Example:
//
//	encoder.AppendKey("name").AppendString("John Doe")
//	// Produces: "name":"John Doe"
func (encoder *JSONEncoder) AppendString(value string) *JSONEncoder {
	encoder.buffer = append(append(append(encoder.buffer, '"'), value...), '"')
	return encoder
}

// AppendNewLine appends a newline character to the buffer.
// This is typically used to format JSON output for readability in log files.
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendNewLine() *JSONEncoder {
	encoder.buffer = append(encoder.buffer, '\n')
	return encoder
}

// AppendByte appends a single byte as a quoted string character.
// The byte is converted to its string representation (the character it represents)
// rather than its numeric value.
//
// Parameters:
//   - value: The byte to append
//
// Returns:
//   - The encoder for method chaining
//
// Example:
//
//	encoder.AppendKey("char").AppendByte('A')
//	// Produces: "char":"A"
func (encoder *JSONEncoder) AppendByte(value byte) *JSONEncoder {
	encoder.buffer = append(encoder.buffer, '"')
	encoder.buffer = append(encoder.buffer, value)
	encoder.buffer = append(encoder.buffer, '"')
	return encoder
}

// AppendBytes appends a byte slice as a quoted string.
// The bytes are interpreted as a string rather than a JSON array of numbers.
//
// Parameters:
//   - values: The byte slice to append
//
// Returns:
//   - The encoder for method chaining
//
// Example:
//
//	encoder.AppendKey("data").AppendBytes([]byte("hello"))
//	// Produces: "data":"hello"
func (encoder *JSONEncoder) AppendBytes(values []byte) *JSONEncoder {
	encoder.buffer = append(encoder.buffer, '"')
	encoder.buffer = append(encoder.buffer, values...)
	encoder.buffer = append(encoder.buffer, '"')
	return encoder
}

// AppendBool appends a JSON boolean value (true or false, unquoted).
//
// Parameters:
//   - value: The boolean value to append
//
// Returns:
//   - The encoder for method chaining
//
// Example:
//
//	encoder.AppendKey("active").AppendBool(true)
//	// Produces: "active":true
func (encoder *JSONEncoder) AppendBool(value bool) *JSONEncoder {
	encoder.buffer = strconv.AppendBool(encoder.buffer, value)
	return encoder
}

// AppendInt appends a JSON integer value in base 10.
//
// Parameters:
//   - value: The int value to append
//
// Returns:
//   - The encoder for method chaining
//
// Example:
//
//	encoder.AppendKey("count").AppendInt(42)
//	// Produces: "count":42
func (encoder *JSONEncoder) AppendInt(value int) *JSONEncoder {
	encoder.buffer = strconv.AppendInt(encoder.buffer, int64(value), 10)
	return encoder
}

// AppendInt8 appends a JSON integer value from an int8 in base 10.
//
// Parameters:
//   - value: The int8 value to append
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendInt8(value int8) *JSONEncoder {
	encoder.buffer = strconv.AppendInt(encoder.buffer, int64(value), 10)
	return encoder
}

// AppendInt16 appends a JSON integer value from an int16 in base 10.
//
// Parameters:
//   - value: The int16 value to append
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendInt16(value int16) *JSONEncoder {
	encoder.buffer = strconv.AppendInt(encoder.buffer, int64(value), 10)
	return encoder
}

// AppendInt32 appends a JSON integer value from an int32 in base 10.
//
// Parameters:
//   - value: The int32 value to append
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendInt32(value int32) *JSONEncoder {
	encoder.buffer = strconv.AppendInt(encoder.buffer, int64(value), 10)
	return encoder
}

// AppendInt64 appends a JSON integer value from an int64 in base 10.
//
// Parameters:
//   - value: The int64 value to append
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendInt64(value int64) *JSONEncoder {
	encoder.buffer = strconv.AppendInt(encoder.buffer, value, 10)
	return encoder
}

// AppendUInt appends a JSON unsigned integer value in base 10.
//
// Parameters:
//   - value: The uint value to append
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendUInt(value uint) *JSONEncoder {
	encoder.buffer = strconv.AppendUint(encoder.buffer, uint64(value), 10)
	return encoder
}

// AppendUInt8 appends a JSON unsigned integer value from a uint8 in base 10.
//
// Parameters:
//   - value: The uint8 value to append
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendUInt8(value uint8) *JSONEncoder {
	encoder.buffer = strconv.AppendUint(encoder.buffer, uint64(value), 10)
	return encoder
}

// AppendUInt16 appends a JSON unsigned integer value from a uint16 in base 10.
//
// Parameters:
//   - value: The uint16 value to append
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendUInt16(value uint16) *JSONEncoder {
	encoder.buffer = strconv.AppendUint(encoder.buffer, uint64(value), 10)
	return encoder
}

// AppendUInt32 appends a JSON unsigned integer value from a uint32 in base 10.
//
// Parameters:
//   - value: The uint32 value to append
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendUInt32(value uint32) *JSONEncoder {
	encoder.buffer = strconv.AppendUint(encoder.buffer, uint64(value), 10)
	return encoder
}

// AppendUInt64 appends a JSON unsigned integer value from a uint64 in base 10.
//
// Parameters:
//   - value: The uint64 value to append
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendUInt64(value uint64) *JSONEncoder {
	encoder.buffer = strconv.AppendUint(encoder.buffer, value, 10)
	return encoder
}

// AppendFloat32 appends a JSON floating-point value from a float32.
// Special values are handled as follows:
//   - NaN is encoded as the string "NaN"
//   - Positive infinity is encoded as the string "+Inf"
//   - Negative infinity is encoded as the string "-Inf"
//
// For normal values, the format is chosen automatically:
//   - Values in range [1e-6, 1e21) use fixed-point notation (e.g., 123.456)
//   - Values outside this range use scientific notation (e.g., 1.23e-7)
//
// Parameters:
//   - value: The float32 value to append
//
// Returns:
//   - The encoder for method chaining
//
// Example:
//
//	encoder.AppendKey("temp").AppendFloat32(98.6)
//	// Produces: "temp":98.6
//	encoder.AppendKey("tiny").AppendFloat32(1e-10)
//	// Produces: "tiny":1e-10
func (encoder *JSONEncoder) AppendFloat32(value float32) *JSONEncoder {
	switch {
	case math.IsNaN(float64(value)):
		encoder.buffer = append(encoder.buffer, `"NaN"`...)
	case math.IsInf(float64(value), 1):
		encoder.buffer = append(encoder.buffer, `"+Inf"`...)
	case math.IsInf(float64(value), -1):
		encoder.buffer = append(encoder.buffer, `"-Inf"`...)
	}

	format := byte('f')
	abs := math.Abs(float64(value))
	if abs != 0 && float32(abs) < 1e-6 || float32(abs) >= 1e21 {
		format = 'e'
	}

	encoder.buffer = strconv.AppendFloat(encoder.buffer, float64(value), format, -1, 32)
	return encoder
}

// AppendFloat64 appends a JSON floating-point value from a float64.
// Special values are handled as follows:
//   - NaN is encoded as the string "NaN"
//   - Positive infinity is encoded as the string "+Inf"
//   - Negative infinity is encoded as the string "-Inf"
//
// For normal values, the format is chosen automatically:
//   - Values in range [1e-6, 1e21) use fixed-point notation (e.g., 123.456)
//   - Values outside this range use scientific notation (e.g., 1.23e-7)
//
// Parameters:
//   - value: The float64 value to append
//
// Returns:
//   - The encoder for method chaining
//
// Example:
//
//	encoder.AppendKey("pi").AppendFloat64(3.14159265359)
//	// Produces: "pi":3.14159265359
func (encoder *JSONEncoder) AppendFloat64(value float64) *JSONEncoder {
	switch {
	case math.IsNaN(float64(value)):
		encoder.buffer = append(encoder.buffer, `"NaN"`...)
	case math.IsInf(float64(value), 1):
		encoder.buffer = append(encoder.buffer, `"+Inf"`...)
	case math.IsInf(float64(value), -1):
		encoder.buffer = append(encoder.buffer, `"-Inf"`...)
	}

	format := byte('f')
	abs := math.Abs(value)
	if abs != 0 && abs < 1e-6 || abs >= 1e21 {
		format = 'e'
	}

	encoder.buffer = strconv.AppendFloat(encoder.buffer, value, format, -1, 64)
	return encoder
}
