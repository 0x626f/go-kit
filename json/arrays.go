package json

// AppendStringArray appends a JSON array of string values.
// Each element is quoted according to JSON string rules.
//
// Parameters:
//   - values: The slice of strings to encode as a JSON array
//
// Returns:
//   - The encoder for method chaining
//
// Example:
//
//	encoder.AppendKey("tags").AppendStringArray([]string{"urgent", "production", "error"})
//	// Produces: "tags":["urgent","production","error"]
func (encoder *JSONEncoder) AppendStringArray(values []string) *JSONEncoder {
	encoder.AppendArrayStart()

	for index := range values {

		encoder.AppendString(values[index])

		if index == len(values)-1 {
			break
		}
		encoder.AppendDelimiter()
	}

	return encoder.AppendArrayEnd()
}

// AppendByteArray appends a JSON array of byte values as characters.
// Each byte is encoded as a single-character string, not as a numeric value.
//
// Parameters:
//   - values: The slice of bytes to encode as a JSON array
//
// Returns:
//   - The encoder for method chaining
//
// Example:
//
//	encoder.AppendKey("chars").AppendByteArray([]byte{'A', 'B', 'C'})
//	// Produces: "chars":["A","B","C"]
func (encoder *JSONEncoder) AppendByteArray(values []byte) *JSONEncoder {
	encoder.AppendArrayStart()

	for index := range values {

		encoder.AppendByte(values[index])

		if index == len(values)-1 {
			break
		}
		encoder.AppendDelimiter()
	}

	return encoder.AppendArrayEnd()
}

// AppendBoolArray appends a JSON array of boolean values.
//
// Parameters:
//   - values: The slice of booleans to encode as a JSON array
//
// Returns:
//   - The encoder for method chaining
//
// Example:
//
//	encoder.AppendKey("flags").AppendBoolArray([]bool{true, false, true})
//	// Produces: "flags":[true,false,true]
func (encoder *JSONEncoder) AppendBoolArray(values []bool) *JSONEncoder {
	encoder.AppendArrayStart()

	for index := range values {

		encoder.AppendBool(values[index])

		if index == len(values)-1 {
			break
		}
		encoder.AppendDelimiter()
	}

	return encoder.AppendArrayEnd()
}

// AppendIntArray appends a JSON array of integer values.
//
// Parameters:
//   - values: The slice of int values to encode as a JSON array
//
// Returns:
//   - The encoder for method chaining
//
// Example:
//
//	encoder.AppendKey("scores").AppendIntArray([]int{95, 87, 92})
//	// Produces: "scores":[95,87,92]
func (encoder *JSONEncoder) AppendIntArray(values []int) *JSONEncoder {
	encoder.AppendArrayStart()

	for index := range values {

		encoder.AppendInt(values[index])

		if index == len(values)-1 {
			break
		}
		encoder.AppendDelimiter()
	}

	return encoder.AppendArrayEnd()
}

// AppendInt8Array appends a JSON array of int8 values.
//
// Parameters:
//   - values: The slice of int8 values to encode as a JSON array
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendInt8Array(values []int8) *JSONEncoder {
	encoder.AppendArrayStart()

	for index := range values {

		encoder.AppendInt8(values[index])

		if index == len(values)-1 {
			break
		}
		encoder.AppendDelimiter()
	}

	return encoder.AppendArrayEnd()
}

// AppendInt16Array appends a JSON array of int16 values.
//
// Parameters:
//   - values: The slice of int16 values to encode as a JSON array
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendInt16Array(values []int16) *JSONEncoder {
	encoder.AppendArrayStart()

	for index := range values {

		encoder.AppendInt16(values[index])

		if index == len(values)-1 {
			break
		}
		encoder.AppendDelimiter()
	}

	return encoder.AppendArrayEnd()
}

// AppendInt32Array appends a JSON array of int32 values.
//
// Parameters:
//   - values: The slice of int32 values to encode as a JSON array
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendInt32Array(values []int32) *JSONEncoder {
	encoder.AppendArrayStart()

	for index := range values {

		encoder.AppendInt32(values[index])

		if index == len(values)-1 {
			break
		}
		encoder.AppendDelimiter()
	}

	return encoder.AppendArrayEnd()
}

// AppendInt64Array appends a JSON array of int64 values.
//
// Parameters:
//   - values: The slice of int64 values to encode as a JSON array
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendInt64Array(values []int64) *JSONEncoder {
	encoder.AppendArrayStart()

	for index := range values {

		encoder.AppendInt64(values[index])

		if index == len(values)-1 {
			break
		}
		encoder.AppendDelimiter()
	}

	return encoder.AppendArrayEnd()
}

// AppendUIntArray appends a JSON array of unsigned integer values.
//
// Parameters:
//   - values: The slice of uint values to encode as a JSON array
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendUIntArray(values []uint) *JSONEncoder {
	encoder.AppendArrayStart()

	for index := range values {

		encoder.AppendUInt(values[index])

		if index == len(values)-1 {
			break
		}
		encoder.AppendDelimiter()
	}

	return encoder.AppendArrayEnd()
}

// AppendUInt8Array appends a JSON array of uint8 values.
//
// Parameters:
//   - values: The slice of uint8 values to encode as a JSON array
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendUInt8Array(values []uint8) *JSONEncoder {
	encoder.AppendArrayStart()

	for index := range values {

		encoder.AppendUInt8(values[index])

		if index == len(values)-1 {
			break
		}
		encoder.AppendDelimiter()
	}

	return encoder.AppendArrayEnd()
}

// AppendUInt16Array appends a JSON array of uint16 values.
//
// Parameters:
//   - values: The slice of uint16 values to encode as a JSON array
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendUInt16Array(values []uint16) *JSONEncoder {
	encoder.AppendArrayStart()

	for index := range values {

		encoder.AppendUInt16(values[index])

		if index == len(values)-1 {
			break
		}
		encoder.AppendDelimiter()
	}

	return encoder.AppendArrayEnd()
}

// AppendUInt32Array appends a JSON array of uint32 values.
//
// Parameters:
//   - values: The slice of uint32 values to encode as a JSON array
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendUInt32Array(values []uint32) *JSONEncoder {
	encoder.AppendArrayStart()

	for index := range values {

		encoder.AppendUInt32(values[index])

		if index == len(values)-1 {
			break
		}
		encoder.AppendDelimiter()
	}

	return encoder.AppendArrayEnd()
}

// AppendUInt64Array appends a JSON array of uint64 values.
//
// Parameters:
//   - values: The slice of uint64 values to encode as a JSON array
//
// Returns:
//   - The encoder for method chaining
func (encoder *JSONEncoder) AppendUInt64Array(values []uint64) *JSONEncoder {
	encoder.AppendArrayStart()

	for index := range values {

		encoder.AppendUInt64(values[index])

		if index == len(values)-1 {
			break
		}
		encoder.AppendDelimiter()
	}

	return encoder.AppendArrayEnd()
}

// AppendFloat32Array appends a JSON array of float32 values.
// Special values (NaN, Inf) are handled according to AppendFloat32 behavior.
//
// Parameters:
//   - values: The slice of float32 values to encode as a JSON array
//
// Returns:
//   - The encoder for method chaining
//
// Example:
//
//	encoder.AppendKey("temperatures").AppendFloat32Array([]float32{98.6, 99.1, 97.8})
//	// Produces: "temperatures":[98.6,99.1,97.8]
func (encoder *JSONEncoder) AppendFloat32Array(values []float32) *JSONEncoder {
	encoder.AppendArrayStart()

	for index := range values {

		encoder.AppendFloat32(values[index])

		if index == len(values)-1 {
			break
		}
		encoder.AppendDelimiter()
	}

	return encoder.AppendArrayEnd()
}

// AppendFloat64Array appends a JSON array of float64 values.
// Special values (NaN, Inf) are handled according to AppendFloat64 behavior.
//
// Parameters:
//   - values: The slice of float64 values to encode as a JSON array
//
// Returns:
//   - The encoder for method chaining
//
// Example:
//
//	encoder.AppendKey("metrics").AppendFloat64Array([]float64{3.14159, 2.71828, 1.41421})
//	// Produces: "metrics":[3.14159,2.71828,1.41421]
func (encoder *JSONEncoder) AppendFloat64Array(values []float64) *JSONEncoder {
	encoder.AppendArrayStart()

	for index := range values {

		encoder.AppendFloat64(values[index])

		if index == len(values)-1 {
			break
		}
		encoder.AppendDelimiter()
	}

	return encoder.AppendArrayEnd()
}
