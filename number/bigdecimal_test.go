package number

import (
	"encoding/json"
	"math/big"
	"testing"
)

func TestBigDecimal_Constructor_FromInt(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  string
	}{
		{"positive", 42, "42"},
		{"negative", -42, "-42"},
		{"zero", 0, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := BigDecimal(tt.value)
			if bf == nil {
				t.Fatal("BigDecimal returned nil")
			}
			if bf.String() != tt.want {
				t.Errorf("got %s, want %s", bf.String(), tt.want)
			}
		})
	}
}

func TestBigDecimal_Constructor_FromUInt(t *testing.T) {
	tests := []struct {
		name  string
		value uint
		want  string
	}{
		{"small", 42, "42"},
		{"zero", 0, "0"},
		{"large", 4294967295, "4294967295"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := BigDecimal(tt.value)
			if bf == nil {
				t.Fatal("BigDecimal returned nil")
			}
			if bf.String() != tt.want {
				t.Errorf("got %s, want %s", bf.String(), tt.want)
			}
		})
	}
}

func TestBigDecimal_Constructor_FromFloat(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  string
	}{
		{"integer float", 42.0, "42"},
		{"decimal", 3.14, "3.14"},
		{"negative", -2.5, "-2.5"},
		{"zero", 0.0, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := BigDecimal(tt.value)
			if bf == nil {
				t.Fatal("BigDecimal returned nil")
			}
			if bf.String() != tt.want {
				t.Errorf("got %s, want %s", bf.String(), tt.want)
			}
		})
	}
}

func TestBigDecimal_Constructor_FromString(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"integer", "123", "123"},
		{"decimal", "123.456", "123.456"},
		{"negative", "-987.654", "-987.654"},
		{"scientific", "1.23e10", "1.23e+10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := BigDecimal(tt.value)
			if bf == nil {
				t.Fatal("BigDecimal returned nil")
			}
			if bf.String() != tt.want {
				t.Errorf("got %s, want %s", bf.String(), tt.want)
			}
		})
	}
}

func TestBigDecimal_Constructor_FromBigInt(t *testing.T) {
	bi := BigNatural(42)
	bf := BigDecimal(bi)

	if bf.String() != "42" {
		t.Errorf("got %s, want 42", bf.String())
	}
}

func TestBigDecimal_Constructor_FromStdBigInt(t *testing.T) {
	stdBigInt := new(big.Int)
	stdBigInt.SetInt64(42)

	bf := BigDecimal(stdBigInt)
	if bf == nil {
		t.Fatal("BigDecimal returned nil")
	}
	if bf.String() != "42" {
		t.Errorf("got %s, want 42", bf.String())
	}

	stdBigInt.SetInt64(100)
	if bf.String() == "100" {
		t.Error("modifying source affected BigDecimal")
	}
}

func TestBigDecimal_Constructor_FromBigFloat(t *testing.T) {
	original := BigDecimal(3.14)
	copied := BigDecimal(original)

	if copied.String() != "3.14" {
		t.Errorf("got %s, want 3.14", copied.String())
	}

	copied.SetFloat(2.71)
	if original.String() == "2.71" {
		t.Error("modifying copied affected original")
	}
}

func TestBigDecimal_Constructor_FromStdBigFloat(t *testing.T) {
	stdBigFloat := new(big.Float)
	stdBigFloat.SetFloat64(3.14)

	bf := BigDecimal(stdBigFloat)
	if bf == nil {
		t.Fatal("BigDecimal returned nil")
	}
	if bf.String() != "3.14" {
		t.Errorf("got %s, want 3.14", bf.String())
	}

	stdBigFloat.SetFloat64(2.71)
	if bf.String() == "2.71" {
		t.Error("modifying source affected BigDecimal")
	}
}

func TestBigFloat_Add(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		b    float64
		want string
	}{
		{"positive", 1.5, 2.5, "4"},
		{"negative", -1.5, -2.5, "-4"},
		{"mixed", 1.5, -0.5, "1"},
		{"zero", 0.0, 3.14, "3.14"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigDecimal(tt.a)
			b := BigDecimal(tt.b)
			result := a.Add(b)

			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigFloat_Subtract(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		b    float64
		want string
	}{
		{"positive", 5.5, 2.5, "3"},
		{"negative result", 2.5, 5.5, "-3"},
		{"zero result", 3.14, 3.14, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigDecimal(tt.a)
			b := BigDecimal(tt.b)
			result := a.Subtract(b)

			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigFloat_Multiply(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		b    float64
		want string
	}{
		{"positive", 2.5, 4.0, "10"},
		{"negative", -2.5, 4.0, "-10"},
		{"both negative", -2.5, -4.0, "10"},
		{"zero", 3.14, 0.0, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigDecimal(tt.a)
			b := BigDecimal(tt.b)
			result := a.Multiply(b)

			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigFloat_Divide(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		b    float64
		want string
	}{
		{"exact", 10.0, 2.0, "5"},
		{"decimal", 10.0, 4.0, "2.5"},
		{"negative", -10.0, 2.0, "-5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigDecimal(tt.a)
			b := BigDecimal(tt.b)
			result := a.Divide(b)

			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigFloat_Sqrt(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		want string
	}{
		{"perfect square", 49.0, "7"},
		{"decimal", 2.0, "1.414213562"},
		{"zero", 0.0, "0"},
		{"one", 1.0, "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigDecimal(tt.a)
			result := a.Sqrt()

			if result == nil {
				t.Fatal("Sqrt returned nil")
			}
			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigFloat_Abs(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		want string
	}{
		{"positive", 3.14, "3.14"},
		{"negative", -3.14, "3.14"},
		{"zero", 0.0, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigDecimal(tt.a)
			result := a.Abs()

			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigFloat_Negate(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		want string
	}{
		{"positive", 3.14, "-3.14"},
		{"negative", -3.14, "3.14"},
		{"zero", 0.0, "-0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigDecimal(tt.a)
			result := a.Negate()

			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigFloat_Sign(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		want int
	}{
		{"positive", 3.14, 1},
		{"negative", -3.14, -1},
		{"zero", 0.0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigDecimal(tt.a)
			result := a.Sign()

			if result != tt.want {
				t.Errorf("got %d, want %d", result, tt.want)
			}
		})
	}
}

func TestBigFloat_Compare(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		b    float64
		want int
	}{
		{"less than", 3.14, 10.5, -1},
		{"equal", 10.5, 10.5, 0},
		{"greater than", 10.5, 3.14, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigDecimal(tt.a)
			b := BigDecimal(tt.b)
			result := a.Compare(b)

			if result != tt.want {
				t.Errorf("got %d, want %d", result, tt.want)
			}
		})
	}
}

func TestBigFloat_Mutability(t *testing.T) {
	a := BigDecimal(3.14)

	if !a.IsMutable() {
		t.Error("new BigFloat should be mutable")
	}

	a.UnMut()
	if a.IsMutable() {
		t.Error("BigFloat should be immutable after UnMut")
	}

	b := BigDecimal(1.0)
	original := a.String()
	result := a.Add(b)

	if a.String() != original {
		t.Error("immutable BigFloat was modified")
	}
	if result.String() != "4.14" {
		t.Errorf("got %s, want 4.14", result.String())
	}

	a.Mut()
	if !a.IsMutable() {
		t.Error("BigFloat should be mutable after Mut")
	}
}

func TestBigFloat_SetInt(t *testing.T) {
	a := BigDecimal(0.0)
	a.SetInt(42)

	if a.String() != "42" {
		t.Errorf("got %s, want 42", a.String())
	}

	a.UnMut()
	a.SetInt(100)
	if a.String() != "42" {
		t.Error("immutable BigFloat was modified by SetInt")
	}
}

func TestBigFloat_SetUInt(t *testing.T) {
	a := BigDecimal(0.0)
	a.SetUInt(42)

	if a.String() != "42" {
		t.Errorf("got %s, want 42", a.String())
	}
}

func TestBigFloat_SetFloat(t *testing.T) {
	a := BigDecimal(0.0)
	a.SetFloat(3.14)

	if a.String() != "3.14" {
		t.Errorf("got %s, want 3.14", a.String())
	}

	a.UnMut()
	a.SetFloat(2.71)
	if a.String() != "3.14" {
		t.Error("immutable BigFloat was modified by SetFloat")
	}
}

func TestBigFloat_SetPrecision(t *testing.T) {
	a := BigDecimal(0.0)
	a.SetPrecision(128)

	if a.Value().Prec() != 128 {
		t.Errorf("got precision %d, want 128", a.Value().Prec())
	}
}

func TestBigFloat_SetRoundingMode(t *testing.T) {
	a := BigDecimal(0.0)
	a.SetRoundingMode(ToNearestEven)

	if a.Value().Mode() != ToNearestEven {
		t.Errorf("got mode %v, want ToNearestEven", a.Value().Mode())
	}
}

func TestBigFloat_Copy(t *testing.T) {
	a := BigDecimal(3.14)
	b := BigDecimal(0.0)
	b.Copy(a)

	if b.String() != "3.14" {
		t.Errorf("got %s, want 3.14", b.String())
	}

	a.SetFloat(2.71)
	if b.String() != "3.14" {
		t.Error("copy was affected by original modification")
	}
}

func TestBigFloat_ToFloat(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{"positive", 3.14},
		{"negative", -2.71},
		{"zero", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := BigDecimal(tt.value)
			result := bf.ToFloat()

			if result != tt.value {
				t.Errorf("got %f, want %f", result, tt.value)
			}
		})
	}
}

func TestBigFloat_MarshalJSON(t *testing.T) {
	a := BigDecimal(3.14)
	data, err := json.Marshal(a)

	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	if string(data) != `"3.14"` {
		t.Errorf("got %s, want \"3.14\"", string(data))
	}
}

func TestBigFloat_UnmarshalJSON(t *testing.T) {
	var a BigFloat
	err := json.Unmarshal([]byte("\"3.14\""), &a)

	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if a.String() != "3.14" {
		t.Errorf("got %s, want 3.14", a.String())
	}

	if a.IsMutable() {
		t.Error("unmarshaled BigFloat should be immutable")
	}
}

func TestBigFloat_UnmarshalJSON_Null(t *testing.T) {
	var a BigFloat
	err := json.Unmarshal([]byte("null"), &a)

	if err != nil {
		t.Fatalf("UnmarshalJSON with null failed: %v", err)
	}
}

func TestBigFloat_MarshalText(t *testing.T) {
	a := BigDecimal(3.14)
	data, err := a.MarshalText()

	if err != nil {
		t.Fatalf("MarshalText failed: %v", err)
	}

	if string(data) != "3.14" {
		t.Errorf("got %s, want 3.14", string(data))
	}
}

func TestBigFloat_UnmarshalText(t *testing.T) {
	var a BigFloat
	err := a.UnmarshalText([]byte("3.14"))

	if err != nil {
		t.Fatalf("UnmarshalText failed: %v", err)
	}

	if a.String() != "3.14" {
		t.Errorf("got %s, want 3.14", a.String())
	}

	if a.IsMutable() {
		t.Error("unmarshaled BigFloat should be immutable")
	}
}

func TestBigFloat_HighPrecision(t *testing.T) {
	a := BigDecimal(0.0)
	a.SetPrecision(256)
	a.SetFloat(3.14159265358979323846)

	result := a.String()
	if len(result) < 10 {
		t.Errorf("expected high precision result, got %s", result)
	}
}

func TestBigFloat_RoundingModes(t *testing.T) {
	modes := []RoundingMode{
		ToNearestEven,
		ToNearestAway,
		ToZero,
		AwayFromZero,
		ToNegativeInf,
		ToPositiveInf,
	}

	for _, mode := range modes {
		a := BigDecimal(0.0)
		a.SetRoundingMode(mode)

		if a.Value().Mode() != mode {
			t.Errorf("rounding mode not set correctly: got %v, want %v", a.Value().Mode(), mode)
		}
	}
}

func TestBigFloat_ChainedOperations(t *testing.T) {
	result := BigDecimal(10.0).
		Add(BigDecimal(5.0)).
		Multiply(BigDecimal(2.0)).
		Subtract(BigDecimal(10.0))

	if result.String() != "20" {
		t.Errorf("got %s, want 20", result.String())
	}
}

func TestBigFloat_LargeNumbers(t *testing.T) {
	large := "123456789012345678901234567890.123456789"
	a := BigDecimal(large)
	b := BigDecimal(large)

	result := a.Add(b)
	expected := "2.46913578e+29"

	if result.String() != expected {
		t.Errorf("got %s, want %s", result.String(), expected)
	}
}

func TestBigFloat_VerySmallNumbers(t *testing.T) {
	small := "0.000000000000000001"
	a := BigDecimal(small)
	b := BigDecimal(small)

	result := a.Add(b)
	expected := "2e-18"

	if result.String() != expected {
		t.Errorf("got %s, want %s", result.String(), expected)
	}
}

func TestBigFloat_DivisionBySmallNumber(t *testing.T) {
	a := BigDecimal(1.0)
	b := BigDecimal(3.0)
	result := a.Divide(b)

	resultStr := result.String()
	if len(resultStr) < 5 {
		t.Errorf("division result should have precision, got %s", resultStr)
	}
}
