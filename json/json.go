// Package encoder provides efficient JSON encoding functionality for structured logging.
// It implements a zero-allocation JSON encoder that builds JSON incrementally
// using a byte buffer, designed specifically for high-performance logging scenarios.
package json

// JSONEncoder is a specialized JSON encoder optimized for logging.
// It builds JSON incrementally in a preallocated buffer to minimize allocations
// and improve performance in high-throughput logging scenarios.
//
// The encoder supports:
//   - Objects and arrays
//   - Primitive types (strings, numbers, booleans)
//   - Arrays of primitives
//   - Nested structures
//
// Usage pattern:
//  1. Create encoder with NewJSONEncoder()
//  2. Build JSON using Append* methods
//  3. Retrieve data with Data()
//  4. Clear buffer with Clear() for reuse
type JSONEncoder struct {
	// buffer holds the JSON being built
	buffer Buffer
}

// Buffer is a byte slice used to build JSON content incrementally.
// It's preallocated to reduce memory allocations during encoding.
type Buffer []byte

// NewJSONEncoder creates a new JSON encoder with a preallocated 500-byte buffer.
// The buffer automatically grows if needed, but the initial size is chosen to
// accommodate typical log entries without reallocation.
//
// Returns:
//   - A pointer to a new JSONEncoder ready for use
//
// Example:
//
//	encoder := encoder.NewJSONEncoder()
//	encoder.AppendObjectStart()
//	encoder.AppendKey("level").AppendString("INFO")
//	encoder.AppendDelimiter()
//	encoder.AppendKey("message").AppendString("User logged in")
//	encoder.AppendObjectEnd()
//	jsonData := encoder.Data()
func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{
		buffer: make(Buffer, 0, 500),
	}
}

// Data returns the current JSON buffer content.
// The returned buffer is a reference to the internal buffer and should not be modified.
//
// Returns:
//   - The byte slice containing the JSON data
func (encoder *JSONEncoder) Data() Buffer {
	return encoder.buffer
}

// Clear resets the buffer to empty while preserving its capacity.
// This allows the encoder to be reused without reallocating memory.
// This is essential for zero-allocation logging when using a sync.Pool.
func (encoder *JSONEncoder) Clear() {
	encoder.buffer = encoder.buffer[:0]
}

// AppendObjectStart appends an opening curly brace '{' to begin a JSON object.
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendObjectStart() *JSONEncoder {
	encoder.buffer = append(encoder.buffer, '{')
	return encoder
}

// AppendObject appends the provided object to end a JSON object.
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendObject(object []byte) *JSONEncoder {
	encoder.buffer = append(encoder.buffer, object...)
	return encoder
}

// AppendObjectEnd appends a closing curly brace '}' to end a JSON object.
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendObjectEnd() *JSONEncoder {
	encoder.buffer = append(encoder.buffer, '}')
	return encoder
}

// AppendArrayStart appends an opening square bracket '[' to begin a JSON array.
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendArrayStart() *JSONEncoder {
	encoder.buffer = append(encoder.buffer, '[')
	return encoder
}

// AppendArrayEnd appends a closing square bracket ']' to end a JSON array.
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendArrayEnd() *JSONEncoder {
	encoder.buffer = append(encoder.buffer, ']')
	return encoder
}

// AppendDelimiter appends a comma separator between JSON elements.
// It intelligently avoids adding a comma if the previous character is already
// a comma or an opening brace (to prevent invalid JSON).
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendDelimiter() *JSONEncoder {
	prev := encoder.buffer[len(encoder.buffer)-1]
	if prev == ',' || prev == '{' {
		return encoder
	}
	encoder.buffer = append(encoder.buffer, ',')
	return encoder
}

// AppendNil appends a JSON null value.
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendNil() *JSONEncoder {
	return encoder.AppendString("null")
}

// AppendKey appends a JSON object key (quoted string followed by colon).
//
// Parameters:
//   - key: The key name to append
//
// Returns:
//   - The encoder for method chaining
//
// Example:
//
//	encoder.AppendKey("username") // Appends: "username":
func (encoder *JSONEncoder) AppendKey(key string) *JSONEncoder {
	encoder.buffer = append(append(encoder.buffer, "\""+key+"\""...), ':')
	return encoder
}
