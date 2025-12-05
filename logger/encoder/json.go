package encoder

type JSONEncoder struct {
	buffer Buffer
}

type Buffer []byte

func NewJSONEncoder() *JSONEncoder {
	return &JSONEncoder{
		buffer: make(Buffer, 0, 500),
	}
}

func (encoder *JSONEncoder) Data() Buffer {
	return encoder.buffer
}

func (encoder *JSONEncoder) Clear() {
	encoder.buffer = encoder.buffer[:0]
}

func (encoder *JSONEncoder) AppendObjectStart() *JSONEncoder {
	encoder.buffer = append(encoder.buffer, '{')
	return encoder
}

func (encoder *JSONEncoder) AppendObjectEnd() *JSONEncoder {
	encoder.buffer = append(encoder.buffer, '}')
	return encoder
}

func (encoder *JSONEncoder) AppendArrayStart() *JSONEncoder {
	encoder.buffer = append(encoder.buffer, '[')
	return encoder
}

func (encoder *JSONEncoder) AppendArrayEnd() *JSONEncoder {
	encoder.buffer = append(encoder.buffer, ']')
	return encoder
}

func (encoder *JSONEncoder) AppendDelimiter() *JSONEncoder {
	prev := encoder.buffer[len(encoder.buffer)-1]
	if prev == ',' || prev == '{' {
		return encoder
	}
	encoder.buffer = append(encoder.buffer, ',')
	return encoder
}

func (encoder *JSONEncoder) AppendNil() *JSONEncoder {
	return encoder.AppendString("null")
}

func (encoder *JSONEncoder) AppendKey(key string) *JSONEncoder {
	encoder.buffer = append(append(encoder.buffer, "\""+key+"\""...), ':')
	return encoder
}
