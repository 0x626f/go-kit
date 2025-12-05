package encoder

import (
	"math"
	"strconv"
)

func (encoder *JSONEncoder) AppendString(value string) *JSONEncoder {
	encoder.buffer = append(append(append(encoder.buffer, '"'), value...), '"')
	return encoder
}

func (encoder *JSONEncoder) AppendNewLine() *JSONEncoder {
	encoder.buffer = append(encoder.buffer, '\n')
	return encoder
}

func (encoder *JSONEncoder) AppendByte(value byte) *JSONEncoder {
	encoder.buffer = append(encoder.buffer, '"')
	encoder.buffer = append(encoder.buffer, value)
	encoder.buffer = append(encoder.buffer, '"')
	return encoder
}

func (encoder *JSONEncoder) AppendBytes(values []byte) *JSONEncoder {
	encoder.buffer = append(encoder.buffer, '"')
	encoder.buffer = append(encoder.buffer, values...)
	encoder.buffer = append(encoder.buffer, '"')
	return encoder
}

func (encoder *JSONEncoder) AppendBool(value bool) *JSONEncoder {
	encoder.buffer = strconv.AppendBool(encoder.buffer, value)
	return encoder
}

func (encoder *JSONEncoder) AppendInt(value int) *JSONEncoder {
	encoder.buffer = strconv.AppendInt(encoder.buffer, int64(value), 10)
	return encoder
}

func (encoder *JSONEncoder) AppendInt8(value int8) *JSONEncoder {
	encoder.buffer = strconv.AppendInt(encoder.buffer, int64(value), 10)
	return encoder
}

func (encoder *JSONEncoder) AppendInt16(value int16) *JSONEncoder {
	encoder.buffer = strconv.AppendInt(encoder.buffer, int64(value), 10)
	return encoder
}

func (encoder *JSONEncoder) AppendInt32(value int32) *JSONEncoder {
	encoder.buffer = strconv.AppendInt(encoder.buffer, int64(value), 10)
	return encoder
}

func (encoder *JSONEncoder) AppendInt64(value int64) *JSONEncoder {
	encoder.buffer = strconv.AppendInt(encoder.buffer, value, 10)
	return encoder
}

func (encoder *JSONEncoder) AppendUInt(value uint) *JSONEncoder {
	encoder.buffer = strconv.AppendUint(encoder.buffer, uint64(value), 10)
	return encoder
}

func (encoder *JSONEncoder) AppendUInt8(value uint8) *JSONEncoder {
	encoder.buffer = strconv.AppendUint(encoder.buffer, uint64(value), 10)
	return encoder
}

func (encoder *JSONEncoder) AppendUInt16(value uint16) *JSONEncoder {
	encoder.buffer = strconv.AppendUint(encoder.buffer, uint64(value), 10)
	return encoder
}

func (encoder *JSONEncoder) AppendUInt32(value uint32) *JSONEncoder {
	encoder.buffer = strconv.AppendUint(encoder.buffer, uint64(value), 10)
	return encoder
}

func (encoder *JSONEncoder) AppendUInt64(value uint64) *JSONEncoder {
	encoder.buffer = strconv.AppendUint(encoder.buffer, value, 10)
	return encoder
}

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
