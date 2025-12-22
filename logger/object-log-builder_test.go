package logger

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// BenchmarkObjectLogBuilder_AssignString validates zero allocations
func BenchmarkObjectLogBuilder_AssignString(b *testing.B) {
	logger := NewLogger("BenchmarkObjectLogBuilder_AssignString").OutputTo(io.Discard).WithLogLevel(INFO)

	// Warm up the pool
	logger.InfoObjectf("warmup").AssignString("key", "value").Build()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.InfoObjectf("test message").
			AssignString("key1", "value1").
			AssignString("key2", "value2").
			Build()
	}
}

// BenchmarkObjectLogBuilder_AssignInt validates zero allocations
func BenchmarkObjectLogBuilder_AssignInt(b *testing.B) {
	logger := NewLogger("BenchmarkObjectLogBuilder_AssignInt").OutputTo(io.Discard).WithLogLevel(INFO)

	// Warm up the pool
	logger.InfoObjectf("warmup").AssignInt("key", 42).Build()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.InfoObjectf("test message").
			AssignInt("count", 123).
			AssignInt("total", 456).
			Build()
	}
}

// BenchmarkObjectLogBuilder_AssignBool validates zero allocations
func BenchmarkObjectLogBuilder_AssignBool(b *testing.B) {
	logger := NewLogger("BenchmarkObjectLogBuilder_AssignBool").OutputTo(io.Discard).WithLogLevel(INFO)

	// Warm up the pool
	logger.InfoObjectf("warmup").AssignBool("key", true).Build()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.InfoObjectf("test message").
			AssignBool("active", true).
			AssignBool("enabled", false).
			Build()
	}
}

// BenchmarkObjectLogBuilder_AssignFloat64 validates zero allocations
func BenchmarkObjectLogBuilder_AssignFloat64(b *testing.B) {
	logger := NewLogger("BenchmarkObjectLogBuilder_AssignFloat64").OutputTo(io.Discard).WithLogLevel(INFO)

	// Warm up the pool
	logger.InfoObjectf("warmup").AssignFloat64("key", 3.14).Build()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.InfoObjectf("test message").
			AssignFloat64("price", 99.99).
			AssignFloat64("tax", 8.75).
			Build()
	}
}

// BenchmarkObjectLogBuilder_AssignInt64 validates zero allocations
func BenchmarkObjectLogBuilder_AssignInt64(b *testing.B) {
	logger := NewLogger("BenchmarkObjectLogBuilder_AssignInt64").OutputTo(io.Discard).WithLogLevel(INFO)

	// Warm up the pool
	logger.InfoObjectf("warmup").AssignInt64("key", 123456789).Build()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.InfoObjectf("test message").
			AssignInt64("Timestamp", 1638720000000).
			AssignInt64("userId", 987654321).
			Build()
	}
}

// BenchmarkObjectLogBuilder_AssignUInt64 validates zero allocations
func BenchmarkObjectLogBuilder_AssignUInt64(b *testing.B) {
	logger := NewLogger("BenchmarkObjectLogBuilder_AssignUInt64").OutputTo(io.Discard).WithLogLevel(INFO)

	// Warm up the pool
	logger.InfoObjectf("warmup").AssignUInt64("key", 123456789).Build()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.InfoObjectf("test message").
			AssignUInt64("bytes", 1024000).
			AssignUInt64("count", 999999).
			Build()
	}
}

// BenchmarkObjectLogBuilder_AssignStringArray validates allocations
func BenchmarkObjectLogBuilder_AssignStringArray(b *testing.B) {
	logger := NewLogger("BenchmarkObjectLogBuilder_AssignStringArray").OutputTo(io.Discard).WithLogLevel(INFO)
	arr := []string{"foo", "bar", "baz"}

	// Warm up the pool
	logger.InfoObjectf("warmup").AssignStringArray("key", arr).Build()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.InfoObjectf("test message").
			AssignStringArray("tags", arr).
			Build()
	}
}

// BenchmarkObjectLogBuilder_AssignIntArray validates allocations
func BenchmarkObjectLogBuilder_AssignIntArray(b *testing.B) {
	logger := NewLogger("BenchmarkObjectLogBuilder_AssignIntArray").OutputTo(io.Discard).WithLogLevel(INFO)
	arr := []int{1, 2, 3, 4, 5}

	// Warm up the pool
	logger.InfoObjectf("warmup").AssignIntArray("key", arr).Build()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.InfoObjectf("test message").
			AssignIntArray("values", arr).
			Build()
	}
}

// BenchmarkObjectLogBuilder_Mixed validates zero allocations with mixed types
func BenchmarkObjectLogBuilder_Mixed(b *testing.B) {
	logger := NewLogger("BenchmarkObjectLogBuilder_Mixed").OutputTo(io.Discard).WithLogLevel(INFO)

	// Warm up the pool
	logger.InfoObjectf("warmup").
		AssignString("name", "test").
		AssignInt("count", 42).
		AssignBool("active", true).
		AssignFloat64("price", 99.99).
		Build()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.InfoObjectf("test message").
			AssignString("user", "john").
			AssignInt("id", 12345).
			AssignBool("admin", true).
			AssignFloat64("balance", 1234.56).
			AssignInt64("Timestamp", 1638720000000).
			Build()
	}
}

// BenchmarkObjectLogBuilder_NestedObjects validates allocations with nested objects
func BenchmarkObjectLogBuilder_NestedObjects(b *testing.B) {
	logger := NewLogger("BenchmarkObjectLogBuilder_NestedObjects").OutputTo(io.Discard).WithLogLevel(INFO)

	// Warm up the pool
	logger.InfoObjectf("warmup").
		NestedStart("user").
		AssignString("name", "test").
		AssignInt("age", 25).
		NestedEnd().
		Build()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.InfoObjectf("test message").
			NestedStart("user").
			AssignString("name", "john").
			AssignInt("age", 30).
			NestedEnd().
			AssignBool("success", true).
			Build()
	}
}

// BenchmarkObjectLogBuilder_AllTypes validates allocations with all primitive types
func BenchmarkObjectLogBuilder_AllTypes(b *testing.B) {
	logger := NewLogger("BenchmarkObjectLogBuilder_AllTypes").OutputTo(io.Discard).WithLogLevel(INFO)

	// Warm up the pool
	logger.InfoObjectf("warmup").
		AssignByte("byte", 'a').
		AssignInt8("int8", 8).
		AssignInt16("int16", 16).
		AssignInt32("int32", 32).
		AssignInt64("int64", 64).
		AssignUInt("uint", 100).
		AssignUInt8("uint8", 8).
		AssignUInt16("uint16", 16).
		AssignUInt32("uint32", 32).
		AssignUInt64("uint64", 64).
		AssignFloat32("float32", 3.14).
		AssignFloat64("float64", 3.14159).
		Build()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.InfoObjectf("test message").
			AssignByte("byte", 'x').
			AssignInt8("int8", 127).
			AssignInt16("int16", 32767).
			AssignInt32("int32", 2147483647).
			AssignInt64("int64", 9223372036854775807).
			AssignUInt("uint", 4294967295).
			AssignUInt8("uint8", 255).
			AssignUInt16("uint16", 65535).
			AssignUInt32("uint32", 4294967295).
			AssignUInt64("uint64", 18446744073709551615).
			AssignFloat32("float32", 1.234).
			AssignFloat64("float64", 1.23456789).
			Build()
	}
}

// Functional Tests

// Test ObjectLogBuilder creation
func TestNewObjectLogBuilder(t *testing.T) {
	logger := NewLogger("TestNewObjectLogBuilder").OutputTo(io.Discard).WithLogLevel(INFO)
	builder := logger.InfoObjectf("test")

	if builder == nil {
		t.Fatal("newObjectLogBuilder returned nil")
	}
	if builder.logger != logger {
		t.Error("Builder logger reference is incorrect")
	}
	if builder.level != INFO {
		t.Error("Builder Level is incorrect")
	}
}

// Test AssignString
func TestObjectLogBuilder_AssignString(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignString").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test message").
		AssignString("name", "John Doe").
		AssignString("city", "New York").
		Build()

	output := buf.String()
	if !strings.Contains(output, `"name":"John Doe"`) {
		t.Errorf("Expected name field in output, got: %s", output)
	}
	if !strings.Contains(output, `"city":"New York"`) {
		t.Errorf("Expected city field in output, got: %s", output)
	}
}

// Test AssignByte
func TestObjectLogBuilder_AssignByte(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignByte").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").AssignByte("char", 'A').Build()

	output := buf.String()
	if !strings.Contains(output, `"char":"A"`) {
		t.Errorf("Expected char:\"A\" in output, got: %s", output)
	}
}

// Test AssignBool
func TestObjectLogBuilder_AssignBool(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignBool").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		AssignBool("active", true).
		AssignBool("disabled", false).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"active":true`) {
		t.Errorf("Expected active:true in output, got: %s", output)
	}
	if !strings.Contains(output, `"disabled":false`) {
		t.Errorf("Expected disabled:false in output, got: %s", output)
	}
}

// Test AssignInt
func TestObjectLogBuilder_AssignInt(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignInt").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").AssignInt("count", 42).Build()

	output := buf.String()
	if !strings.Contains(output, `"count":42`) {
		t.Errorf("Expected count:42 in output, got: %s", output)
	}
}

// Test AssignInt8
func TestObjectLogBuilder_AssignInt8(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignInt8").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").AssignInt8("value", 127).Build()

	output := buf.String()
	if !strings.Contains(output, `"value":127`) {
		t.Errorf("Expected value:127 in output, got: %s", output)
	}
}

// Test AssignInt16
func TestObjectLogBuilder_AssignInt16(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignInt16").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").AssignInt16("value", 32767).Build()

	output := buf.String()
	if !strings.Contains(output, `"value":32767`) {
		t.Errorf("Expected value:32767 in output, got: %s", output)
	}
}

// Test AssignInt32
func TestObjectLogBuilder_AssignInt32(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignInt32").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").AssignInt32("value", 2147483647).Build()

	output := buf.String()
	if !strings.Contains(output, `"value":2147483647`) {
		t.Errorf("Expected value:2147483647 in output, got: %s", output)
	}
}

// Test AssignInt64
func TestObjectLogBuilder_AssignInt64(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignInt64").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").AssignInt64("value", 9223372036854775807).Build()

	output := buf.String()
	if !strings.Contains(output, `"value":9223372036854775807`) {
		t.Errorf("Expected value:9223372036854775807 in output, got: %s", output)
	}
}

// Test AssignUInt
func TestObjectLogBuilder_AssignUInt(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignUInt").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").AssignUInt("value", 4294967295).Build()

	output := buf.String()
	if !strings.Contains(output, `"value":4294967295`) {
		t.Errorf("Expected value:4294967295 in output, got: %s", output)
	}
}

// Test AssignUInt8
func TestObjectLogBuilder_AssignUInt8(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignUInt8").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").AssignUInt8("value", 255).Build()

	output := buf.String()
	if !strings.Contains(output, `"value":255`) {
		t.Errorf("Expected value:255 in output, got: %s", output)
	}
}

// Test AssignUInt16
func TestObjectLogBuilder_AssignUInt16(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignUInt16").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").AssignUInt16("value", 65535).Build()

	output := buf.String()
	if !strings.Contains(output, `"value":65535`) {
		t.Errorf("Expected value:65535 in output, got: %s", output)
	}
}

// Test AssignUInt32
func TestObjectLogBuilder_AssignUInt32(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignUInt32").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").AssignUInt32("value", 4294967295).Build()

	output := buf.String()
	if !strings.Contains(output, `"value":4294967295`) {
		t.Errorf("Expected value:4294967295 in output, got: %s", output)
	}
}

// Test AssignUInt64
func TestObjectLogBuilder_AssignUInt64(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignUInt64").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").AssignUInt64("value", 18446744073709551615).Build()

	output := buf.String()
	if !strings.Contains(output, `"value":18446744073709551615`) {
		t.Errorf("Expected value:18446744073709551615 in output, got: %s", output)
	}
}

// Test AssignFloat32
func TestObjectLogBuilder_AssignFloat32(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignFloat32").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").AssignFloat32("value", 3.14).Build()

	output := buf.String()
	if !strings.Contains(output, `"value":3.14`) {
		t.Errorf("Expected value:3.14 in output, got: %s", output)
	}
}

// Test AssignFloat64
func TestObjectLogBuilder_AssignFloat64(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignFloat64").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").AssignFloat64("value", 3.141592653589793).Build()

	output := buf.String()
	if !strings.Contains(output, `"value":3.141592653589793`) {
		t.Errorf("Expected value:3.141592653589793 in output, got: %s", output)
	}
}

// Test AssignStringArray
func TestObjectLogBuilder_AssignStringArray(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignStringArray").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		AssignStringArray("tags", []string{"go", "logger", "test"}).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"tags":["go","logger","test"]`) {
		t.Errorf("Expected tags array in output, got: %s", output)
	}
}

// Test AssignByteArray
func TestObjectLogBuilder_AssignByteArray(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignByteArray").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		AssignByteArray("bytes", []byte{65, 66, 67}).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"bytes":["A","B","C"]`) {
		t.Errorf("Expected bytes array as strings in output, got: %s", output)
	}
}

// Test AssignBoolArray
func TestObjectLogBuilder_AssignBoolArray(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignBoolArray").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		AssignBoolArray("flags", []bool{true, false, true}).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"flags":[true,false,true]`) {
		t.Errorf("Expected flags array in output, got: %s", output)
	}
}

// Test AssignIntArray
func TestObjectLogBuilder_AssignIntArray(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignIntArray").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		AssignIntArray("numbers", []int{1, 2, 3, 4, 5}).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"numbers":[1,2,3,4,5]`) {
		t.Errorf("Expected numbers array in output, got: %s", output)
	}
}

// Test AssignInt8Array
func TestObjectLogBuilder_AssignInt8Array(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignInt8Array").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		AssignInt8Array("values", []int8{1, 2, 3}).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"values":[1,2,3]`) {
		t.Errorf("Expected values array in output, got: %s", output)
	}
}

// Test AssignInt16Array
func TestObjectLogBuilder_AssignInt16Array(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignInt16Array").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		AssignInt16Array("values", []int16{100, 200, 300}).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"values":[100,200,300]`) {
		t.Errorf("Expected values array in output, got: %s", output)
	}
}

// Test AssignInt32Array
func TestObjectLogBuilder_AssignInt32Array(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignInt32Array").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		AssignInt32Array("values", []int32{1000, 2000, 3000}).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"values":[1000,2000,3000]`) {
		t.Errorf("Expected values array in output, got: %s", output)
	}
}

// Test AssignIn64Array (note the typo in the method name)
func TestObjectLogBuilder_AssignIn64Array(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignIn64Array").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		AssignIn64Array("values", []int64{10000, 20000, 30000}).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"values":[10000,20000,30000]`) {
		t.Errorf("Expected values array in output, got: %s", output)
	}
}

// Test AssignUIntArray
func TestObjectLogBuilder_AssignUIntArray(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignUIntArray").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		AssignUIntArray("values", []uint{1, 2, 3}).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"values":[1,2,3]`) {
		t.Errorf("Expected values array in output, got: %s", output)
	}
}

// Test AssignUInt8Array
func TestObjectLogBuilder_AssignUInt8Array(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignUInt8Array").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		AssignUInt8Array("values", []uint8{10, 20, 30}).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"values":[10,20,30]`) {
		t.Errorf("Expected values array in output, got: %s", output)
	}
}

// Test AssignUInt16Array
func TestObjectLogBuilder_AssignUInt16Array(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignUInt16Array").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		AssignUInt16Array("values", []uint16{100, 200, 300}).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"values":[100,200,300]`) {
		t.Errorf("Expected values array in output, got: %s", output)
	}
}

// Test AssignUInt32Array
func TestObjectLogBuilder_AssignUInt32Array(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignUInt32Array").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		AssignUInt32Array("values", []uint32{1000, 2000, 3000}).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"values":[1000,2000,3000]`) {
		t.Errorf("Expected values array in output, got: %s", output)
	}
}

// Test AssignUInt64Array
func TestObjectLogBuilder_AssignUInt64Array(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignUInt64Array").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		AssignUInt64Array("values", []uint64{10000, 20000, 30000}).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"values":[10000,20000,30000]`) {
		t.Errorf("Expected values array in output, got: %s", output)
	}
}

// Test AssignFloat32Array
func TestObjectLogBuilder_AssignFloat32Array(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignFloat32Array").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		AssignFloat32Array("values", []float32{1.1, 2.2, 3.3}).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"values":[1.1,2.2,3.3`) {
		t.Errorf("Expected values array in output, got: %s", output)
	}
}

// Test AssignFloat64Array
func TestObjectLogBuilder_AssignFloat64Array(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_AssignFloat64Array").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		AssignFloat64Array("values", []float64{1.11, 2.22, 3.33}).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"values":[1.11,2.22,3.33`) {
		t.Errorf("Expected values array in output, got: %s", output)
	}
}

// Test NestedStart and NestedEnd
func TestObjectLogBuilder_NestedObjects(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_NestedObjects").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("test").
		NestedStart("user").
		AssignString("name", "John").
		AssignInt("age", 30).
		NestedEnd().
		AssignBool("verified", true).
		Build()

	output := buf.String()
	if !strings.Contains(output, `"user":{`) {
		t.Errorf("Expected nested user object start in output, got: %s", output)
	}
	if !strings.Contains(output, `"name":"John"`) {
		t.Errorf("Expected name in nested object, got: %s", output)
	}
	if !strings.Contains(output, `"age":30`) {
		t.Errorf("Expected age in nested object, got: %s", output)
	}
	if !strings.Contains(output, `"verified":true`) {
		t.Errorf("Expected verified field, got: %s", output)
	}
}

// Test method chaining
func TestObjectLogBuilder_MethodChaining(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_MethodChaining").OutputTo(&buf).WithLogLevel(INFO)

	builder := logger.InfoObjectf("test").
		AssignString("field1", "value1").
		AssignInt("field2", 123).
		AssignBool("field3", true)

	if builder == nil {
		t.Error("Method chaining broke the builder")
	}

	builder.Build()

	output := buf.String()
	if !strings.Contains(output, `"field1":"value1"`) {
		t.Error("Method chaining failed to set field1")
	}
	if !strings.Contains(output, `"field2":123`) {
		t.Error("Method chaining failed to set field2")
	}
	if !strings.Contains(output, `"field3":true`) {
		t.Error("Method chaining failed to set field3")
	}
}

// Test builder pool reuse
func TestObjectLogBuilder_PoolReuse(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_PoolReuse").OutputTo(&buf).WithLogLevel(INFO)

	// First build
	logger.InfoObjectf("message1").AssignString("field", "value1").Build()
	output1 := buf.String()

	buf.Reset()

	// Second build (should reuse from pool)
	logger.InfoObjectf("message2").AssignString("field", "value2").Build()
	output2 := buf.String()

	if !strings.Contains(output1, "message1") {
		t.Error("First build failed")
	}
	if !strings.Contains(output2, "message2") {
		t.Error("Second build failed")
	}
	if strings.Contains(output2, "message1") {
		t.Error("Pool did not clear previous data")
	}
}

// Test empty message
func TestObjectLogBuilder_EmptyMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_EmptyMessage").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("").AssignString("key", "value").Build()

	output := buf.String()
	if !strings.Contains(output, `"key":"value"`) {
		t.Errorf("Expected key:value in output, got: %s", output)
	}
}

// Test complex nested structure
func TestObjectLogBuilder_ComplexNested(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("TestObjectLogBuilder_ComplexNested").OutputTo(&buf).WithLogLevel(INFO)

	logger.InfoObjectf("complex test").
		AssignString("topLevel", "value").
		NestedStart("level1").
		AssignString("field1", "value1").
		NestedStart("level2").
		AssignString("field2", "value2").
		AssignInt("number", 42).
		NestedEnd().
		AssignBool("flag", true).
		NestedEnd().
		AssignString("end", "done").
		Build()

	output := buf.String()
	if !strings.Contains(output, `"topLevel":"value"`) {
		t.Error("Missing top Level field")
	}
	if !strings.Contains(output, `"level1":{`) {
		t.Error("Missing level1 nested object")
	}
	if !strings.Contains(output, `"level2":{`) {
		t.Error("Missing level2 nested object")
	}
	if !strings.Contains(output, `"number":42`) {
		t.Error("Missing nested number field")
	}
}

// Test all log levels with object builder
func TestObjectLogBuilder_AllLogLevels(t *testing.T) {
	tests := []struct {
		level    LogLevel
		name     string
		buildFn  func(*Logger) *ObjectLogBuilder
		expected string
	}{
		{TRACE, "TRACE", func(l *Logger) *ObjectLogBuilder { return l.TraceObjectf("test") }, "TRACE"},
		{DEBUG, "DEBUG", func(l *Logger) *ObjectLogBuilder { return l.DebugObjectf("test") }, "DEBUG"},
		{INFO, "INFO", func(l *Logger) *ObjectLogBuilder { return l.InfoObjectf("test") }, "INFO"},
		{WARNING, "WARNING", func(l *Logger) *ObjectLogBuilder { return l.WarningObjectf("test") }, "WARNING"},
		{ERROR, "ERROR", func(l *Logger) *ObjectLogBuilder { return l.ErrorObjectf("test") }, "ERROR"},
		{NONE, "NONE", func(l *Logger) *ObjectLogBuilder { return l.LogObjectf("test") }, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger("TestObjectLogBuilder_AllLogLevels_" + tt.name).OutputTo(&buf).ErrorsTo(&buf).WithLogLevel(tt.level)

			builder := tt.buildFn(logger)
			if builder != nil {
				builder.AssignString("key", "value").Build()

				output := buf.String()
				if tt.expected != "" && !strings.Contains(output, `"Level":"`+tt.expected+`"`) {
					t.Errorf("Expected Level %s in output, got: %s", tt.expected, output)
				}
			}
		})
	}
}
