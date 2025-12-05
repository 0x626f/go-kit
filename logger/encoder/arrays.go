package encoder

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
