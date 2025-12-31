// Package number provides arbitrary-precision numeric types and utilities.
package number

import (
	"encoding/json"
	"math/big"
)

// RoundingMode is an alias for big.RoundingMode, used for floating-point rounding operations.
type RoundingMode = big.RoundingMode

// Rounding mode constants from the math/big package with IEEE 754-2008 equivalents.
const (
	// ToNearestEven rounds to the nearest even value (IEEE 754-2008 roundTiesToEven).
	ToNearestEven RoundingMode = iota
	// ToNearestAway rounds to the nearest value, away from zero in case of a tie (IEEE 754-2008 roundTiesToAway).
	ToNearestAway
	// ToZero rounds toward zero (IEEE 754-2008 roundTowardZero).
	ToZero
	// AwayFromZero rounds away from zero (no IEEE 754-2008 equivalent).
	AwayFromZero
	// ToNegativeInf rounds toward negative infinity (IEEE 754-2008 roundTowardNegative).
	ToNegativeInf
	// ToPositiveInf rounds toward positive infinity (IEEE 754-2008 roundTowardPositive).
	ToPositiveInf
)

// BigFloat represents an arbitrary-precision floating-point number.
// It wraps the standard library's big.Float with additional functionality
// and implements the BigNumericField interface.
type BigFloat struct {
	// mutable controls whether operations can modify this instance.
	// If false, operations will create a new instance.
	mutable bool
	// value is the underlying big.Float value.
	value *big.Float
}

// BigDecimal creates a new BigFloat from various numeric types.
// This is the primary constructor for BigFloat objects.
//
// Type parameters:
//   - T: The source type, which can be a numeric primitive, string, or BigInt/BigFloat
//
// Returns a new mutable BigFloat instance initialized with the provided value.
// Returns nil if conversion fails.
func BigDecimal[T NumericField | string | BigInt | *BigInt | BigFloat | *BigFloat | *big.Float | *big.Int](value T) *BigFloat {

	obj := new(big.Float)

	switch v := any(value).(type) {
	case string:
		obj, _ = obj.SetString(v)
		break
	case *big.Float:
		obj.Set(v)
		break
	case *big.Int:
		obj.SetInt(v)
		break
	case *BigInt:
		obj = obj.SetInt(v.value)
		break
	case BigInt:
		obj = obj.SetInt(v.value)
		break
	case *BigFloat:
		obj = obj.Set(v.value)
		break
	case BigFloat:
		obj = obj.Set(v.value)
	case float32:
		obj.SetFloat64(float64(v))
		break
	case float64:
		obj.SetFloat64(v)
		break
	case int:
		obj.SetInt64(int64(v))
		break
	case int8:
		obj.SetInt64(int64(v))
		break
	case int16:
		obj.SetInt64(int64(v))
		break
	case int32:
		obj.SetInt64(int64(v))
		break
	case int64:
		obj.SetInt64(v)
		break
	case uint:
		obj.SetUint64(uint64(v))
		break
	case uint8:
		obj.SetUint64(uint64(v))
		break
	case uint16:
		obj.SetUint64(uint64(v))
		break
	case uint32:
		obj.SetUint64(uint64(v))
		break
	case uint64:
		obj.SetUint64(v)
		break
	}

	if obj == nil {
		return nil
	}

	return &BigFloat{
		mutable: true,
		value:   obj,
	}
}

// Value returns the underlying big.Float value.
// This provides direct access to the standard library's big.Float methods.
func (number *BigFloat) Value() *big.Float {
	return number.value
}

// IsMutable returns true if this instance can be modified by operations.
// This implements part of the BigNumericField interface.
func (number *BigFloat) IsMutable() bool {
	return number.mutable
}

// forwardBigFloat is a helper function that returns the receiver if it's mutable,
// or creates a new mutable copy if it's not.
// This ensures operations don't modify immutable numbers.
func forwardBigFloat(number *BigFloat) *BigFloat {
	if number.IsMutable() {
		return number
	}
	return BigDecimal(number)
}

// Mut marks this instance as mutable, allowing operations to modify it.
// Returns the receiver for method chaining.
func (number *BigFloat) Mut() *BigFloat {
	number.mutable = true
	return number
}

// UnMut marks this instance as immutable, preventing operations from modifying it.
// Operations will create new instances instead.
// Returns the receiver for method chaining.
func (number *BigFloat) UnMut() *BigFloat {
	number.mutable = false
	return number
}

// SetInt sets this BigFloat to the integer value if mutable.
// Returns the receiver for method chaining.
func (number *BigFloat) SetInt(arg int) *BigFloat {
	if number.IsMutable() {
		number.value.SetInt64(int64(arg))
	}
	return number
}

// SetUInt sets this BigFloat to the unsigned integer value if mutable.
// Returns the receiver for method chaining.
func (number *BigFloat) SetUInt(arg uint) *BigFloat {
	if number.IsMutable() {
		number.value.SetUint64(uint64(arg))
	}
	return number
}

// SetFloat sets this BigFloat to the floating-point value if mutable.
// Returns the receiver for method chaining.
func (number *BigFloat) SetFloat(arg float64) *BigFloat {
	if number.IsMutable() {
		number.value.SetFloat64(arg)
	}
	return number
}

// SetRoundingMode sets the rounding mode for this BigFloat if mutable.
// Returns the receiver for method chaining.
func (number *BigFloat) SetRoundingMode(mode RoundingMode) *BigFloat {
	if number.IsMutable() {
		number.value.SetMode(mode)
	}
	return number
}

// SetPrecision sets the precision of this BigFloat in bits if mutable.
// The precision is the number of bits used to store the mantissa.
// Returns the receiver for method chaining.
func (number *BigFloat) SetPrecision(bits uint) *BigFloat {
	if number.IsMutable() {
		number.value.SetPrec(bits)
	}
	return number
}

// Copy sets this BigFloat to the value of the argument if mutable.
// Returns the receiver for method chaining.
func (number *BigFloat) Copy(arg *BigFloat) *BigFloat {
	if number.IsMutable() {
		number.value.Set(arg.value)
	}
	return number
}

// Add adds the argument to this BigFloat and returns the result.
// If this instance is immutable, a new instance is created.
// Implements part of the BigNumericField interface.
func (number *BigFloat) Add(arg *BigFloat) *BigFloat {
	obj := forwardBigFloat(number)
	obj.value.Add(obj.value, arg.value)
	return obj
}

// Subtract subtracts the argument from this BigFloat and returns the result.
// If this instance is immutable, a new instance is created.
// Implements part of the BigNumericField interface.
func (number *BigFloat) Subtract(arg *BigFloat) *BigFloat {
	obj := forwardBigFloat(number)
	obj.value.Sub(obj.value, arg.value)
	return obj
}

// Multiply multiplies this BigFloat by the argument and returns the result.
// If this instance is immutable, a new instance is created.
// Implements part of the BigNumericField interface.
func (number *BigFloat) Multiply(arg *BigFloat) *BigFloat {
	obj := forwardBigFloat(number)
	obj.value.Mul(obj.value, arg.value)
	return obj
}

// Divide divides this BigFloat by the argument and returns the result.
// If this instance is immutable, a new instance is created.
// Implements part of the BigNumericField interface.
func (number *BigFloat) Divide(arg *BigFloat) *BigFloat {
	obj := forwardBigFloat(number)
	obj.value.Quo(obj.value, arg.value)
	return obj
}

// Sqrt calculates the square root of this BigFloat and returns the result.
// If this instance is immutable, a new instance is created.
// Implements part of the BigNumericField interface.
func (number *BigFloat) Sqrt() *BigFloat {
	obj := forwardBigFloat(number)
	obj.value.Sqrt(obj.value)
	return obj
}

// Abs calculates the absolute value of this BigFloat and returns the result.
// If this instance is immutable, a new instance is created.
// Implements part of the BigNumericField interface.
func (number *BigFloat) Abs() *BigFloat {
	obj := forwardBigFloat(number)
	obj.value.Abs(obj.value)
	return obj
}

// Negate negates this BigFloat and returns the result.
// If this instance is immutable, a new instance is created.
// Implements part of the BigNumericField interface.
func (number *BigFloat) Negate() *BigFloat {
	obj := forwardBigFloat(number)
	obj.value.Neg(obj.value)
	return obj
}

// Sign returns:
//   - 1 if this BigFloat is greater than zero
//   - 0 if this BigFloat is zero
//   - -1 if this BigFloat is less than zero
//
// Implements part of the BigNumericField interface.
func (number *BigFloat) Sign() int {
	return number.value.Sign()
}

// Compare compares this BigFloat with the argument and returns:
//   - 1 if this BigFloat is greater than the argument
//   - 0 if they are equal
//   - -1 if this BigFloat is less than the argument
//
// Implements the Comparable interface.
func (number *BigFloat) Compare(arg *BigFloat) int {
	return number.value.Cmp(arg.value)
}

// ToFloat converts this BigFloat to a float64 value.
// Note that this may lose precision for very large numbers.
// The conversion error is silently ignored.
//
// Returns the float64 approximation of this BigFloat.
func (number *BigFloat) ToFloat() float64 {
	value, _ := number.value.Float64()
	return value
}

// String returns the string representation of this BigFloat.
// Implements the fmt.Stringer interface.
func (number *BigFloat) String() string {
	return number.value.String()
}

// MarshalText implements the encoding.TextMarshaler interface.
// This allows BigFloat to be serialized to text formats like XML.
func (number *BigFloat) MarshalText() (text []byte, err error) {
	return number.value.MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// This allows BigFloat to be deserialized from text formats like XML.
// The resulting BigFloat will be immutable.
func (number *BigFloat) UnmarshalText(text []byte) error {
	number.mutable = false
	number.value = new(big.Float)
	return number.value.UnmarshalText(text)
}

// MarshalJSON implements the json.Marshaler interface.
// This allows BigFloat to be serialized to JSON.
func (number *BigFloat) MarshalJSON() (bytes []byte, err error) {
	return json.Marshal(number.value.Text('g', -1))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// This allows BigFloat to be deserialized from JSON.
// The resulting BigFloat will be immutable.
func (number *BigFloat) UnmarshalJSON(bytes []byte) error {
	number.mutable = false
	number.value = new(big.Float)

	var literal string
	if err := json.Unmarshal(bytes, &literal); err != nil {
		return err
	}

	if literal == "null" {
		return nil
	}

	number.value.SetString(literal)

	return nil
}

func (number *BigFloat) BigNatural() *BigInt {
	return BigNatural(number)
}
