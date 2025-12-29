package number

import (
	"encoding/json"
	"math/big"
	"testing"
)

func TestBigNatural_Constructor_FromInt(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  string
	}{
		{"positive", 42, "42"},
		{"negative", -42, "-42"},
		{"zero", 0, "0"},
		{"max int32", 2147483647, "2147483647"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := BigNatural(tt.value)
			if bi == nil {
				t.Fatal("BigNatural returned nil")
			}
			if bi.String() != tt.want {
				t.Errorf("got %s, want %s", bi.String(), tt.want)
			}
		})
	}
}

func TestBigNatural_Constructor_FromUInt(t *testing.T) {
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
			bi := BigNatural(tt.value)
			if bi == nil {
				t.Fatal("BigNatural returned nil")
			}
			if bi.String() != tt.want {
				t.Errorf("got %s, want %s", bi.String(), tt.want)
			}
		})
	}
}

func TestBigNatural_Constructor_FromString(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"decimal", "123456789", "123456789"},
		{"negative", "-987654321", "-987654321"},
		{"large", "123456789012345678901234567890", "123456789012345678901234567890"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bi := BigNatural(tt.value)
			if bi == nil {
				t.Fatal("BigNatural returned nil")
			}
			if bi.String() != tt.want {
				t.Errorf("got %s, want %s", bi.String(), tt.want)
			}
		})
	}
}

func TestBigNatural_Constructor_FromStdBigInt(t *testing.T) {
	stdBigInt := new(big.Int)
	stdBigInt.SetInt64(42)

	bi := BigNatural(stdBigInt)
	if bi == nil {
		t.Fatal("BigNatural returned nil")
	}
	if bi.String() != "42" {
		t.Errorf("got %s, want 42", bi.String())
	}

	stdBigInt.SetInt64(100)
	if bi.String() == "100" {
		t.Error("modifying source affected BigNatural")
	}
}

func TestBigNatural_Constructor_FromBigInt(t *testing.T) {
	original := BigNatural(42)
	copied := BigNatural(original)

	if copied.String() != "42" {
		t.Errorf("got %s, want 42", copied.String())
	}

	copied.SetInt(100)
	if original.String() == "100" {
		t.Error("modifying copied affected original")
	}
}

func TestBigNatural_Constructor_FromBigFloat(t *testing.T) {
	bf := BigDecimal(42.7)
	bi := BigNatural(bf)

	if bi.String() != "42" {
		t.Errorf("got %s, want 42", bi.String())
	}
}

func TestBigNatural_Constructor_FromStdBigFloat(t *testing.T) {
	stdBigFloat := new(big.Float)
	stdBigFloat.SetFloat64(42.7)

	bi := BigNatural(stdBigFloat)
	if bi == nil {
		t.Fatal("BigNatural returned nil")
	}
	if bi.String() != "42" {
		t.Errorf("got %s, want 42", bi.String())
	}
}

func TestBigInt_Add(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want string
	}{
		{"positive", 10, 20, "30"},
		{"negative", -10, -20, "-30"},
		{"mixed", 10, -5, "5"},
		{"zero", 0, 42, "42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigNatural(tt.a)
			b := BigNatural(tt.b)
			result := a.Add(b)

			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigInt_Subtract(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want string
	}{
		{"positive", 30, 10, "20"},
		{"negative result", 10, 20, "-10"},
		{"negative operands", -10, -5, "-5"},
		{"zero result", 42, 42, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigNatural(tt.a)
			b := BigNatural(tt.b)
			result := a.Subtract(b)

			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigInt_Multiply(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want string
	}{
		{"positive", 6, 7, "42"},
		{"negative", -6, 7, "-42"},
		{"both negative", -6, -7, "42"},
		{"zero", 42, 0, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigNatural(tt.a)
			b := BigNatural(tt.b)
			result := a.Multiply(b)

			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigInt_Divide(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want string
	}{
		{"exact", 42, 6, "7"},
		{"truncate", 43, 6, "7"},
		{"negative", -42, 6, "-7"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigNatural(tt.a)
			b := BigNatural(tt.b)
			result := a.Divide(b)

			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigInt_Sqrt(t *testing.T) {
	tests := []struct {
		name string
		a    int
		want string
	}{
		{"perfect square", 49, "7"},
		{"truncate", 50, "7"},
		{"zero", 0, "0"},
		{"one", 1, "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigNatural(tt.a)
			result := a.Sqrt()

			if result == nil {
				t.Fatal("Sqrt returned nil")
			}
			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}

	t.Run("negative", func(t *testing.T) {
		a := BigNatural(-4)
		result := a.Sqrt()
		if result != nil {
			t.Error("Sqrt of negative should return nil")
		}
	})
}

func TestBigInt_Abs(t *testing.T) {
	tests := []struct {
		name string
		a    int
		want string
	}{
		{"positive", 42, "42"},
		{"negative", -42, "42"},
		{"zero", 0, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigNatural(tt.a)
			result := a.Abs()

			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigInt_Negate(t *testing.T) {
	tests := []struct {
		name string
		a    int
		want string
	}{
		{"positive", 42, "-42"},
		{"negative", -42, "42"},
		{"zero", 0, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigNatural(tt.a)
			result := a.Negate()

			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigInt_Exponent(t *testing.T) {
	tests := []struct {
		name string
		base int
		exp  int
		want string
	}{
		{"simple", 2, 3, "8"},
		{"zero exp", 42, 0, "1"},
		{"one exp", 42, 1, "42"},
		{"large", 10, 10, "10000000000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := BigNatural(tt.base)
			exp := BigNatural(tt.exp)
			result := base.Exponent(exp)

			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigInt_Mod(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want string
	}{
		{"simple", 10, 3, "1"},
		{"zero remainder", 12, 3, "0"},
		{"negative dividend", -10, 3, "2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigNatural(tt.a)
			b := BigNatural(tt.b)
			result := a.Mod(b)

			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigInt_Remainder(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want string
	}{
		{"simple", 10, 3, "1"},
		{"zero remainder", 12, 3, "0"},
		{"negative dividend", -10, 3, "-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigNatural(tt.a)
			b := BigNatural(tt.b)
			result := a.Remainder(b)

			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigInt_Sign(t *testing.T) {
	tests := []struct {
		name string
		a    int
		want int
	}{
		{"positive", 42, 1},
		{"negative", -42, -1},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigNatural(tt.a)
			result := a.Sign()

			if result != tt.want {
				t.Errorf("got %d, want %d", result, tt.want)
			}
		})
	}
}

func TestBigInt_GCD(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want string
	}{
		{"simple", 48, 18, "6"},
		{"coprime", 17, 19, "1"},
		{"one zero", 42, 0, "42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigNatural(tt.a)
			b := BigNatural(tt.b)
			result := a.GCD(b)

			if result.String() != tt.want {
				t.Errorf("got %s, want %s", result.String(), tt.want)
			}
		})
	}
}

func TestBigInt_Bitwise_And(t *testing.T) {
	a := BigNatural(12) // 1100
	b := BigNatural(10) // 1010
	result := a.And(b)  // 1000 = 8

	if result.String() != "8" {
		t.Errorf("got %s, want 8", result.String())
	}
}

func TestBigInt_Bitwise_Or(t *testing.T) {
	a := BigNatural(12) // 1100
	b := BigNatural(10) // 1010
	result := a.Or(b)   // 1110 = 14

	if result.String() != "14" {
		t.Errorf("got %s, want 14", result.String())
	}
}

func TestBigInt_Bitwise_Xor(t *testing.T) {
	a := BigNatural(12) // 1100
	b := BigNatural(10) // 1010
	result := a.Xor(b)  // 0110 = 6

	if result.String() != "6" {
		t.Errorf("got %s, want 6", result.String())
	}
}

func TestBigInt_Bitwise_Not(t *testing.T) {
	a := BigNatural(0)
	result := a.Not()

	if result.String() != "-1" {
		t.Errorf("got %s, want -1", result.String())
	}
}

func TestBigInt_Bitwise_AndNot(t *testing.T) {
	a := BigNatural(12)   // 1100
	b := BigNatural(10)   // 1010
	result := a.AndNot(b) // 0100 = 4

	if result.String() != "4" {
		t.Errorf("got %s, want 4", result.String())
	}
}

func TestBigInt_LeftShift(t *testing.T) {
	a := BigNatural(5)       // 101
	result := a.LeftShift(2) // 10100 = 20

	if result.String() != "20" {
		t.Errorf("got %s, want 20", result.String())
	}
}

func TestBigInt_RightShift(t *testing.T) {
	a := BigNatural(20)       // 10100
	result := a.RightShift(2) // 101 = 5

	if result.String() != "5" {
		t.Errorf("got %s, want 5", result.String())
	}
}

func TestBigInt_BitAt(t *testing.T) {
	a := BigNatural(5) // 101

	tests := []struct {
		index uint
		want  uint
	}{
		{0, 1},
		{1, 0},
		{2, 1},
		{3, 0},
	}

	for _, tt := range tests {
		result := a.BitAt(tt.index)
		if result != tt.want {
			t.Errorf("BitAt(%d) = %d, want %d", tt.index, result, tt.want)
		}
	}
}

func TestBigInt_BitLen(t *testing.T) {
	tests := []struct {
		name string
		a    int
		want int
	}{
		{"zero", 0, 0},
		{"one", 1, 1},
		{"two", 2, 2},
		{"255", 255, 8},
		{"256", 256, 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigNatural(tt.a)
			result := a.BitLen()

			if result != tt.want {
				t.Errorf("got %d, want %d", result, tt.want)
			}
		})
	}
}

func TestBigInt_Compare(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"less than", 5, 10, -1},
		{"equal", 10, 10, 0},
		{"greater than", 10, 5, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := BigNatural(tt.a)
			b := BigNatural(tt.b)
			result := a.Compare(b)

			if result != tt.want {
				t.Errorf("got %d, want %d", result, tt.want)
			}
		})
	}
}

func TestBigInt_Mutability(t *testing.T) {
	a := BigNatural(42)

	if !a.IsMutable() {
		t.Error("new BigInt should be mutable")
	}

	a.UnMut()
	if a.IsMutable() {
		t.Error("BigInt should be immutable after UnMut")
	}

	b := BigNatural(10)
	original := a.String()
	result := a.Add(b)

	if a.String() != original {
		t.Error("immutable BigInt was modified")
	}
	if result.String() != "52" {
		t.Errorf("got %s, want 52", result.String())
	}

	a.Mut()
	if !a.IsMutable() {
		t.Error("BigInt should be mutable after Mut")
	}
}

func TestBigInt_SetInt(t *testing.T) {
	a := BigNatural(0)
	a.SetInt(42)

	if a.String() != "42" {
		t.Errorf("got %s, want 42", a.String())
	}

	a.UnMut()
	a.SetInt(100)
	if a.String() != "42" {
		t.Error("immutable BigInt was modified by SetInt")
	}
}

func TestBigInt_SetUInt(t *testing.T) {
	a := BigNatural(0)
	a.SetUInt(42)

	if a.String() != "42" {
		t.Errorf("got %s, want 42", a.String())
	}
}

func TestBigInt_Copy(t *testing.T) {
	a := BigNatural(42)
	b := BigNatural(0)
	b.Copy(a)

	if b.String() != "42" {
		t.Errorf("got %s, want 42", b.String())
	}

	a.SetInt(100)
	if b.String() != "42" {
		t.Error("copy was affected by original modification")
	}
}

func TestBigInt_Bytes(t *testing.T) {
	a := BigNatural(255)
	bytes := a.Bytes()

	if len(bytes) != 1 || bytes[0] != 255 {
		t.Errorf("got %v, want [255]", bytes)
	}

	b := BigNatural(0)
	b.SetBytes(bytes)

	if b.String() != "255" {
		t.Errorf("got %s, want 255", b.String())
	}
}

func TestBigInt_BigDecimal(t *testing.T) {
	a := BigNatural(42)
	bf := a.BigDecimal()

	if bf == nil {
		t.Fatal("BigDecimal returned nil")
	}
	if bf.String() != "42" {
		t.Errorf("got %s, want 42", bf.String())
	}
}

func TestBigInt_MarshalJSON(t *testing.T) {
	a := BigNatural(42)
	data, err := json.Marshal(a)

	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	if string(data) != "42" {
		t.Errorf("got %s, want 42", string(data))
	}
}

func TestBigInt_UnmarshalJSON(t *testing.T) {
	var a BigInt
	err := json.Unmarshal([]byte("\"42\""), &a)

	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if a.String() != "42" {
		t.Errorf("got %s, want 42", a.String())
	}

	if a.IsMutable() {
		t.Error("unmarshaled BigInt should be immutable")
	}
}

func TestBigInt_MarshalText(t *testing.T) {
	a := BigNatural(42)
	data, err := a.MarshalText()

	if err != nil {
		t.Fatalf("MarshalText failed: %v", err)
	}

	if string(data) != "42" {
		t.Errorf("got %s, want 42", string(data))
	}
}

func TestBigInt_UnmarshalText(t *testing.T) {
	var a BigInt
	err := a.UnmarshalText([]byte("42"))

	if err != nil {
		t.Fatalf("UnmarshalText failed: %v", err)
	}

	if a.String() != "42" {
		t.Errorf("got %s, want 42", a.String())
	}

	if a.IsMutable() {
		t.Error("unmarshaled BigInt should be immutable")
	}
}

func TestBigInt_LargeNumbers(t *testing.T) {
	large := "123456789012345678901234567890"
	a := BigNatural(large)
	b := BigNatural(large)

	result := a.Add(b)
	expected := "246913578024691357802469135780"

	if result.String() != expected {
		t.Errorf("got %s, want %s", result.String(), expected)
	}
}

func TestBigInt_ChainedOperations(t *testing.T) {
	result := BigNatural(10).
		Add(BigNatural(5)).
		Multiply(BigNatural(2)).
		Subtract(BigNatural(10))

	if result.String() != "20" {
		t.Errorf("got %s, want 20", result.String())
	}
}
