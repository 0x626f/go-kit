package logger

import (
	"github.com/0x626f/go-kit/json"
	"sync"
	"time"
)

// builderPool is a sync.Pool that recycles ObjectLogBuilder instances to achieve zero allocations.
// When Build() is called, the builder is cleaned and returned to this pool for reuse.
// This is a critical optimization for high-throughput logging scenarios where creating
// new builder instances would cause significant GC pressure.
var builderPool = &sync.Pool{
	New: func() any {
		return &ObjectLogBuilder{
			json: json.NewJSONEncoder(),
		}
	},
}

// ObjectLogBuilder provides a fluent interface for constructing structured JSON logs with zero allocations.
// It uses a pooled JSON encoder to build log output incrementally without heap allocations.
//
// The builder must be finalized with Build() to emit the log and return the builder to the pool.
// Do NOT reuse a builder after calling Build().
//
// Example usage:
//
//	logger.InfoObjectf("HTTP request").
//	    AssignString("method", "GET").
//	    AssignString("path", "/api/users").
//	    AssignInt("status", 200).
//	    AssignFloat64("duration_ms", 45.2).
//	    Build()
//
// Output JSON:
//
//	{
//	  "level": "INFO",
//	  "message": "HTTP request",
//	  "object": {
//	    "method": "GET",
//	    "path": "/api/users",
//	    "status": 200,
//	    "duration_ms": 45.2
//	  }
//	}
type ObjectLogBuilder struct {
	// logger is the parent Logger instance that will emit the log
	logger *Logger
	// json is the pooled JSON encoder used to build the output
	json *json.JSONEncoder
	// level is the log level for this message
	level LogLevel
}

// newObjectLogBuilder retrieves a builder from the pool and initializes it for a new log entry.
// This is an internal function called by Logger's *Objectf() methods.
//
// Parameters:
//   - logger: The Logger instance that will emit this log
//   - level: The log level for this entry
//   - msg: The log message bytes (already formatted)
//
// Returns:
//   - A pooled ObjectLogBuilder ready for field assignment
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

// Build finalizes the JSON log, emits it to the logger's output, and returns the builder to the pool.
// This method MUST be called to complete the log entry.
// After calling Build(), do not use this builder instance again.
//
// Example:
//
//	logger.InfoObjectf("Task completed").
//	    AssignString("task_id", "12345").
//	    AssignInt("duration_sec", 42).
//	    Build() // Emits the log and recycles the builder
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

// AssignString adds a string field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The string value
//
// Returns:
//   - The builder for method chaining
//
// Example:
//
//	builder.AssignString("username", "alice")
//	// Produces: "username":"alice"
func (builder *ObjectLogBuilder) AssignString(name string, value string) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendString(value)
	return builder
}

// AssignByte adds a byte field to the log object as a character string.
// The byte is encoded as its character representation, not its numeric value.
//
// Parameters:
//   - name: The field name
//   - value: The byte value
//
// Returns:
//   - The builder for method chaining
//
// Example:
//
//	builder.AssignByte("initial", 'A')
//	// Produces: "initial":"A"
func (builder *ObjectLogBuilder) AssignByte(name string, value byte) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendByte(value)
	return builder
}

// AssignBool adds a boolean field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The boolean value
//
// Returns:
//   - The builder for method chaining
//
// Example:
//
//	builder.AssignBool("success", true)
//	// Produces: "success":true
func (builder *ObjectLogBuilder) AssignBool(name string, value bool) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendBool(value)
	return builder
}

// AssignInt adds an integer field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The int value
//
// Returns:
//   - The builder for method chaining
//
// Example:
//
//	builder.AssignInt("count", 42)
//	// Produces: "count":42
func (builder *ObjectLogBuilder) AssignInt(name string, value int) *ObjectLogBuilder {
	builder.json.AppendDelimiter()
	builder.json.AppendKey(name)
	builder.json.AppendInt(value)

	return builder
}

// AssignInt8 adds an int8 field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The int8 value
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignInt8(name string, value int8) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendInt8(value)
	return builder
}

// AssignInt16 adds an int16 field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The int16 value
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignInt16(name string, value int16) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendInt16(value)
	return builder
}

// AssignInt32 adds an int32 field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The int32 value
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignInt32(name string, value int32) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendInt32(value)
	return builder
}

// AssignInt64 adds an int64 field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The int64 value
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignInt64(name string, value int64) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendInt64(value)
	return builder
}

// AssignUInt adds an unsigned integer field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The uint value
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignUInt(name string, value uint) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt(value)
	return builder
}

// AssignUInt8 adds a uint8 field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The uint8 value
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignUInt8(name string, value uint8) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt8(value)
	return builder
}

// AssignUInt16 adds a uint16 field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The uint16 value
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignUInt16(name string, value uint16) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt16(value)
	return builder
}

// AssignUInt32 adds a uint32 field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The uint32 value
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignUInt32(name string, value uint32) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt32(value)
	return builder
}

// AssignUInt64 adds a uint64 field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The uint64 value
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignUInt64(name string, value uint64) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt64(value)
	return builder
}

// AssignFloat32 adds a float32 field to the log object.
// Special values (NaN, Inf) are encoded as strings.
//
// Parameters:
//   - name: The field name
//   - value: The float32 value
//
// Returns:
//   - The builder for method chaining
//
// Example:
//
//	builder.AssignFloat32("temperature", 98.6)
//	// Produces: "temperature":98.6
func (builder *ObjectLogBuilder) AssignFloat32(name string, value float32) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendFloat32(value)
	return builder
}

// AssignFloat64 adds a float64 field to the log object.
// Special values (NaN, Inf) are encoded as strings.
//
// Parameters:
//   - name: The field name
//   - value: The float64 value
//
// Returns:
//   - The builder for method chaining
//
// Example:
//
//	builder.AssignFloat64("duration_ms", 125.437)
//	// Produces: "duration_ms":125.437
func (builder *ObjectLogBuilder) AssignFloat64(name string, value float64) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendFloat64(value)
	return builder
}

// AssignStringArray adds a string array field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The string slice
//
// Returns:
//   - The builder for method chaining
//
// Example:
//
//	builder.AssignStringArray("tags", []string{"urgent", "bug", "security"})
//	// Produces: "tags":["urgent","bug","security"]
func (builder *ObjectLogBuilder) AssignStringArray(name string, value []string) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendStringArray(value)
	return builder
}

// AssignByteArray adds a byte array field to the log object.
// Each byte is encoded as a character string, not as a numeric value.
//
// Parameters:
//   - name: The field name
//   - value: The byte slice
//
// Returns:
//   - The builder for method chaining
//
// Example:
//
//	builder.AssignByteArray("chars", []byte{'A', 'B', 'C'})
//	// Produces: "chars":["A","B","C"]
func (builder *ObjectLogBuilder) AssignByteArray(name string, value []byte) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendByteArray(value)
	return builder
}

// AssignBoolArray adds a boolean array field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The bool slice
//
// Returns:
//   - The builder for method chaining
//
// Example:
//
//	builder.AssignBoolArray("flags", []bool{true, false, true})
//	// Produces: "flags":[true,false,true]
func (builder *ObjectLogBuilder) AssignBoolArray(name string, value []bool) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendBoolArray(value)
	return builder
}

// AssignIntArray adds an integer array field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The int slice
//
// Returns:
//   - The builder for method chaining
//
// Example:
//
//	builder.AssignIntArray("scores", []int{95, 87, 92})
//	// Produces: "scores":[95,87,92]
func (builder *ObjectLogBuilder) AssignIntArray(name string, value []int) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendIntArray(value)
	return builder
}

// AssignInt8Array adds an int8 array field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The int8 slice
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignInt8Array(name string, value []int8) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendInt8Array(value)
	return builder
}

// AssignInt16Array adds an int16 array field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The int16 slice
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignInt16Array(name string, value []int16) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendInt16Array(value)
	return builder
}

// AssignInt32Array adds an int32 array field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The int32 slice
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignInt32Array(name string, value []int32) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendInt32Array(value)
	return builder
}

// AssignIn64Array adds an int64 array field to the log object.
// Note: This method has a typo in its name (should be AssignInt64Array).
//
// Parameters:
//   - name: The field name
//   - value: The int64 slice
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignIn64Array(name string, value []int64) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendInt64Array(value)
	return builder
}

// AssignUIntArray adds a uint array field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The uint slice
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignUIntArray(name string, value []uint) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUIntArray(value)
	return builder
}

// AssignUInt8Array adds a uint8 array field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The uint8 slice
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignUInt8Array(name string, value []uint8) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt8Array(value)
	return builder
}

// AssignUInt16Array adds a uint16 array field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The uint16 slice
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignUInt16Array(name string, value []uint16) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt16Array(value)
	return builder
}

// AssignUInt32Array adds a uint32 array field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The uint32 slice
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignUInt32Array(name string, value []uint32) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt32Array(value)
	return builder
}

// AssignUInt64Array adds a uint64 array field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The uint64 slice
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) AssignUInt64Array(name string, value []uint64) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendUInt64Array(value)
	return builder
}

// AssignFloat32Array adds a float32 array field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The float32 slice
//
// Returns:
//   - The builder for method chaining
//
// Example:
//
//	builder.AssignFloat32Array("temps", []float32{98.6, 99.1, 97.8})
//	// Produces: "temps":[98.6,99.1,97.8]
func (builder *ObjectLogBuilder) AssignFloat32Array(name string, value []float32) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendFloat32Array(value)
	return builder
}

// AssignFloat64Array adds a float64 array field to the log object.
//
// Parameters:
//   - name: The field name
//   - value: The float64 slice
//
// Returns:
//   - The builder for method chaining
//
// Example:
//
//	builder.AssignFloat64Array("metrics", []float64{0.95, 0.87, 0.92})
//	// Produces: "metrics":[0.95,0.87,0.92]
func (builder *ObjectLogBuilder) AssignFloat64Array(name string, value []float64) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendFloat64Array(value)
	return builder
}

// NestedStart begins a nested object field.
// Must be paired with NestedEnd() to close the nested object.
//
// Parameters:
//   - name: The field name for the nested object
//
// Returns:
//   - The builder for method chaining
//
// Example:
//
//	builder.NestedStart("request").
//	    AssignString("method", "GET").
//	    AssignInt("status", 200).
//	    NestedEnd()
//	// Produces: "request":{"method":"GET","status":200}
func (builder *ObjectLogBuilder) NestedStart(name string) *ObjectLogBuilder {
	builder.json.AppendDelimiter().AppendKey(name).AppendObjectStart()
	return builder
}

// NestedEnd closes a nested object started with NestedStart().
//
// Returns:
//   - The builder for method chaining
func (builder *ObjectLogBuilder) NestedEnd() *ObjectLogBuilder {
	builder.json.AppendObjectEnd()
	return builder
}
