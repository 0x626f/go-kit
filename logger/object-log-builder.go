package logger

import (
	"github.com/0x626f/go-kit/logger/encoder"
	"sync"
	"time"
)

var builderPool = &sync.Pool{
	New: func() any {
		return &ObjectLogBuilder{
			json: encoder.NewJSONEncoder(),
		}
	},
}

type ObjectLogBuilder struct {
	logger *Logger
	json   *encoder.JSONEncoder
	level  LogLevel
}

func newObjectLogBuilder(logger *Logger, level LogLevel, msg []byte) *ObjectLogBuilder {
	instance := builderPool.Get().(*ObjectLogBuilder)
	instance.logger = logger
	instance.level = level

	instance.json.AppendObjectStart()

	// insert log level
	instance.json.AppendKey("level").AppendString(level.String()).AppendDelimiter()
	// insert timestamp
	if instance.logger.options.timestamp {
		timestamp := time.Now().Format(logger.options.timestampFormat)
		instance.json.AppendKey("timestamp").AppendString(timestamp).AppendDelimiter()
	}

	// insert a source of the log
	if len(logger.name) > 0 {
		instance.json.AppendKey("source").AppendString(logger.name).AppendDelimiter()
	}

	// insert log message
	if len(msg) > 0 {
		instance.json.AppendKey("message").AppendBytes(msg).AppendDelimiter()
	}

	instance.json.AppendKey("object").AppendObjectStart()
	return instance
}

func (builder *ObjectLogBuilder) Build() {
	builder.json.AppendObjectEnd().AppendObjectEnd().AppendNewLine()

	if builder.logger.options.async {
		builder.logger.sendToChannelByLevel(builder.level, builder.json.Data())
	} else {
		builder.logger.writeByLevel(builder.level, builder.json.Data())
	}

	builder.json.Clear()
	builderPool.Put(builder)
}

func (builder *ObjectLogBuilder) AssignString(name string, value string) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendString(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignByte(name string, value byte) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendByte(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignBool(name string, value bool) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendBool(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignInt(name string, value int) *ObjectLogBuilder {
	builder.json.AppendDelimiter()
	builder.json.AppendKey(name)
	builder.json.AppendInt(value)

	return builder
}

func (builder *ObjectLogBuilder) AssignInt8(name string, value int8) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendInt8(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignInt16(name string, value int16) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendInt16(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignInt32(name string, value int32) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendInt32(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignInt64(name string, value int64) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendInt64(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignUInt(name string, value uint) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignUInt8(name string, value uint8) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt8(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignUInt16(name string, value uint16) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt16(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignUInt32(name string, value uint32) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt32(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignUInt64(name string, value uint64) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt64(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignFloat32(name string, value float32) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendFloat32(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignFloat64(name string, value float64) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendFloat64(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignStringArray(name string, value []string) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendStringArray(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignByteArray(name string, value []byte) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendByteArray(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignBoolArray(name string, value []bool) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendBoolArray(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignIntArray(name string, value []int) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendIntArray(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignInt8Array(name string, value []int8) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendInt8Array(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignInt16Array(name string, value []int16) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendInt16Array(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignInt32Array(name string, value []int32) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendInt32Array(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignIn64Array(name string, value []int64) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendInt64Array(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignUIntArray(name string, value []uint) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUIntArray(value)
	return builder
}
func (builder *ObjectLogBuilder) AssignUInt8Array(name string, value []uint8) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt8Array(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignUInt16Array(name string, value []uint16) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt16Array(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignUInt32Array(name string, value []uint32) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt32Array(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignUInt64Array(name string, value []uint64) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt64Array(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignFloat32Array(name string, value []float32) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendFloat32Array(value)
	return builder
}

func (builder *ObjectLogBuilder) AssignFloat64Array(name string, value []float64) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendFloat64Array(value)
	return builder
}

func (builder *ObjectLogBuilder) NestedStart(name string) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendObjectStart()
	return builder
}

func (builder *ObjectLogBuilder) NestedEnd() *ObjectLogBuilder {
	builder.json.AppendObjectEnd()
	return builder
}
