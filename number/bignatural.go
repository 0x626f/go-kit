// Package number provides arbitrary-precision numeric types and utilities.
package number

import (
	"math/big"
	"strings"
)

// BigInt represents an arbitrary-precision integer number.
// It wraps the standard library's big.Int with additional functionality
// and implements the BigNumericField interface.
type BigInt struct {
	// Mutable controls whether operations can modify this instance.
	// If false, operations will create a new instance.
	mutable bool
	// value is the underlying big.Int value.
	value *big.Int
}

// BigNatural creates a new BigInt from various numeric types.
// This is the primary constructor for BigInt objects.
//
// Type parameters:
//   - T: The source type, which can be a numeric primitive, string, or BigInt/BigFloat
//
// String values can be in different bases:
//   - "0x..." for hexadecimal (base 16)
//   - "0o..." for octal (base 8)
//   - Otherwise decimal (base 10) is assumed
//
// Returns a new mutable BigInt instance initialized with the provided value.
// Returns nil if conversion fails.
func BigNatural[T NumericField | string | BigInt | *BigInt | BigFloat | *BigFloat | *big.Int | *big.Float](value T) *BigInt {

	obj := new(big.Int)

	switch v := any(value).(type) {
	case string:
		if strings.HasPrefix(v, "0x") {
			obj, _ = obj.SetString(strings.Replace(v, "0x", "", 1), 16)
		} else if strings.HasPrefix(v, "0o") {
			obj, _ = obj.SetString(strings.Replace(v, "0o", "", 1), 8)
		} else {
			obj, _ = obj.SetString(v, 10)
		}
		break
	case *big.Int:
		obj.Set(v)
		break
	case *big.Float:
		v.Int(obj)
		break
	case *BigInt:
		obj = obj.Set(v.value)
		break
	case BigInt:
		obj = obj.Set(v.value)
		break
	case *BigFloat:
		v.value.Int(obj)
		break
	case BigFloat:
		obj, _ = v.value.Int(obj)
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
	case float32:
		obj.SetInt64(int64(v))
		break
	case float64:
		obj.SetInt64(int64(v))
		break
	}

	if obj == nil {
		return nil
	}

	return &BigInt{
		mutable: true,
		value:   obj,
	}
}

// Value returns the underlying big.Int value.
// This provides direct access to the standard library's big.Int methods.
func (number *BigInt) Value() *big.Int {
	return number.value
}

// IsMutable returns true if this instance can be modified by operations.
// This implements part of the BigNumericField interface.
func (number *BigInt) IsMutable() bool {
	return number.mutable
}

// forwardBigInt is a helper function that returns the receiver if it's mutable,
// or creates a new mutable copy if it's not.
// This ensures operations don't modify immutable numbers.
func forwardBigInt(number *BigInt) *BigInt {
	if number.IsMutable() {
		return number
	}
	return BigNatural(int(number.value.Int64()))
}

// Mut marks this instance as mutable, allowing operations to modify it.
// Returns the receiver for method chaining.
func (number *BigInt) Mut() *BigInt {
	number.mutable = true
	return number
}

// UnMut marks this instance as immutable, preventing operations from modifying it.
// Operations will create new instances instead.
// Returns the receiver for method chaining.
func (number *BigInt) UnMut() *BigInt {
	number.mutable = false
	return number
}

// SetInt sets this BigInt to the integer value if mutable.
// Returns the receiver for method chaining.
func (number *BigInt) SetInt(arg int) *BigInt {
	if number.IsMutable() {
		number.value.SetInt64(int64(arg))
	}
	return number
}

// SetUInt sets this BigInt to the unsigned integer value if mutable.
// Returns the receiver for method chaining.
func (number *BigInt) SetUInt(arg uint) *BigInt {
	if number.IsMutable() {
		number.value.SetUint64(uint64(arg))
	}
	return number
}

// Copy sets this BigInt to the value of the argument if mutable.
// Returns the receiver for method chaining.
func (number *BigInt) Copy(arg *BigInt) *BigInt {
	if number.IsMutable() {
		number.value.Set(arg.value)
	}
	return number
}

// Add adds the argument to this BigInt and returns the result.
// If this instance is immutable, a new instance is created.
// Implements part of the BigNumericField interface.
func (number *BigInt) Add(arg *BigInt) *BigInt {
	obj := forwardBigInt(number)
	obj.value.Add(obj.value, arg.value)
	return obj
}

// Substruct subtracts the argument from this BigInt and returns the result.
// If this instance is immutable, a new instance is created.
// Implements part of the BigNumericField interface.
func (number *BigInt) Subtract(arg *BigInt) *BigInt {
	obj := forwardBigInt(number)
	obj.value.Sub(obj.value, arg.value)
	return obj
}

// Multiply multiplies this BigInt by the argument and returns the result.
// If this instance is immutable, a new instance is created.
// Implements part of the BigNumericField interface.
func (number *BigInt) Multiply(arg *BigInt) *BigInt {
	obj := forwardBigInt(number)
	obj.value.Mul(obj.value, arg.value)
	return obj
}

// Divide divides this BigInt by the argument and returns the result (truncated division).
// If this instance is immutable, a new instance is created.
// Implements part of the BigNumericField interface.
func (number *BigInt) Divide(arg *BigInt) *BigInt {
	obj := forwardBigInt(number)
	obj.value.Quo(obj.value, arg.value)
	return obj
}

// Sqrt calculates the integer square root of this BigInt and returns the result.
// Returns nil if this BigInt is negative.
// If this instance is immutable, a new instance is created.
// Implements part of the BigNumericField interface.
func (number *BigInt) Sqrt() *BigInt {
	if number.value.Sign() < 0 {
		return nil
	}

	obj := forwardBigInt(number)
	obj.value.Sqrt(obj.value)

	return obj
}

// Abs calculates the absolute value of this BigInt and returns the result.
// If this instance is immutable, a new instance is created.
// Implements part of the BigNumericField interface.
func (number *BigInt) Abs() *BigInt {
	obj := forwardBigInt(number)
	obj.value.Abs(obj.value)
	return obj
}

// Negate negates this BigInt and returns the result.
// If this instance is immutable, a new instance is created.
// Implements part of the BigNumericField interface.
func (number *BigInt) Negate() *BigInt {
	obj := forwardBigInt(number)
	obj.value.Neg(obj.value)
	return obj
}

// Exponent raises this BigInt to the power of the argument and returns the result.
// If this instance is immutable, a new instance is created.
func (number *BigInt) Exponent(arg *BigInt) *BigInt {
	obj := forwardBigInt(number)
	obj.value.Exp(obj.value, arg.value, nil)
	return obj
}

// Mod computes the modulus of this BigInt by the argument and returns the result.
// The result will have the same sign as the argument (divisor).
// If this instance is immutable, a new instance is created.
func (number *BigInt) Mod(arg *BigInt) *BigInt {
	obj := forwardBigInt(number)
	obj.value.Mod(obj.value, arg.value)
	return obj
}

// Remainder computes the remainder of dividing this BigInt by the argument and returns the result.
// The result will have the same sign as this BigInt (dividend).
// If this instance is immutable, a new instance is created.
func (number *BigInt) Remainder(arg *BigInt) *BigInt {
	obj := forwardBigInt(number)
	obj.value.Rem(number.value, arg.value)
	return obj
}

// Sign returns:
//   - 1 if this BigInt is greater than zero
//   - 0 if this BigInt is zero
//   - -1 if this BigInt is less than zero
//
// Implements part of the BigNumericField interface.
func (number *BigInt) Sign() int {
	return number.value.Sign()
}

// GCD computes the greatest common divisor of this BigInt and the argument.
// Returns a new BigInt containing the result.
func (number *BigInt) GCD(arg *BigInt) *BigInt {
	obj := BigNatural(0)
	obj.value.GCD(nil, nil, number.value, arg.value)
	return obj
}

// Bytes returns the absolute value of this BigInt as a big-endian byte slice.
func (number *BigInt) Bytes() []byte {
	return number.value.Bytes()
}

// SetBytes sets this BigInt to the value represented by the given big-endian byte slice.
// Returns the receiver for method chaining.
func (number *BigInt) SetBytes(bytes []byte) *BigInt {
	if number.IsMutable() {
		number.value.SetBytes(bytes)
	}
	return number
}

// BitLen returns the length of this BigInt in bits.
// The bit length is the minimum number of bits needed to represent this BigInt.
func (number *BigInt) BitLen() int {
	return number.value.BitLen()
}

// LeftShift performs a left shift operation (multiplies by 2^n) and returns the result.
// If this instance is immutable, a new instance is created.
func (number *BigInt) LeftShift(n uint) *BigInt {
	obj := forwardBigInt(number)
	obj.value.Lsh(obj.value, n)
	return obj
}

// RightShift performs a right shift operation (divides by 2^n) and returns the result.
// If this instance is immutable, a new instance is created.
func (number *BigInt) RightShift(n uint) *BigInt {
	obj := forwardBigInt(number)
	obj.value.Rsh(obj.value, n)
	return obj
}

// BitAt returns the bit at the specified index (0 or 1).
// Index 0 is the least significant bit.
func (number *BigInt) BitAt(index uint) uint {
	return number.value.Bit(int(index))
}

// And performs a bitwise AND operation with the argument and returns the result.
// If this instance is immutable, a new instance is created.
func (number *BigInt) And(arg *BigInt) *BigInt {
	obj := forwardBigInt(number)
	obj.value.And(obj.value, arg.value)
	return obj
}

// AndNot performs a bitwise AND NOT operation (AND with the complement of arg)
// and returns the result.
// If this instance is immutable, a new instance is created.
func (number *BigInt) AndNot(arg *BigInt) *BigInt {
	obj := forwardBigInt(number)
	obj.value.AndNot(obj.value, arg.value)
	return obj
}

// Or performs a bitwise OR operation with the argument and returns the result.
// If this instance is immutable, a new instance is created.
func (number *BigInt) Or(arg *BigInt) *BigInt {
	obj := forwardBigInt(number)
	obj.value.Or(obj.value, arg.value)
	return obj
}

// Xor performs a bitwise XOR operation with the argument and returns the result.
// If this instance is immutable, a new instance is created.
func (number *BigInt) Xor(arg *BigInt) *BigInt {
	obj := forwardBigInt(number)
	obj.value.Xor(obj.value, arg.value)
	return obj
}

// Not performs a bitwise NOT operation (complement) and returns the result.
// If this instance is immutable, a new instance is created.
func (number *BigInt) Not() *BigInt {
	obj := forwardBigInt(number)
	obj.value.Not(obj.value)
	return obj
}

// Compare compares this BigInt with the argument and returns:
//   - 1 if this BigInt is greater than the argument
//   - 0 if they are equal
//   - -1 if this BigInt is less than the argument
//
// Implements the Comparable interface.
func (number *BigInt) Compare(arg *BigInt) int {
	return number.value.Cmp(arg.value)
}

// BigDecimal converts this BigInt to a BigFloat.
// Returns a new BigFloat containing the exact integer value.
func (number *BigInt) BigDecimal() *BigFloat {
	return BigDecimal(number)
}

// String returns the string representation of this BigInt.
// Implements the fmt.Stringer interface.
func (number *BigInt) String() string {
	return number.value.String()
}

// MarshalText implements the encoding.TextMarshaler interface.
// This allows BigInt to be serialized to text formats like XML.
func (number *BigInt) MarshalText() (text []byte, err error) {
	return number.value.MarshalText()
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// This allows BigInt to be deserialized from text formats like XML.
// The resulting BigInt will be immutable.
func (number *BigInt) UnmarshalText(text []byte) error {
	number.mutable = false
	number.value = new(big.Int)
	return number.value.UnmarshalText(text)
}

// MarshalJSON implements the json.Marshaler interface.
// This allows BigInt to be serialized to JSON.
func (number *BigInt) MarshalJSON() (text []byte, err error) {
	return number.value.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// This allows BigInt to be deserialized from JSON.
// The resulting BigInt will be immutable.
func (number *BigInt) UnmarshalJSON(text []byte) error {
	number.mutable = false
	number.value = new(big.Int)
	return number.value.UnmarshalJSON(text)
}
