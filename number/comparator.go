// Package number provides arbitrary-precision numeric types and utilities.
package number

// NumericComparator compares two values of any numeric primitive type and returns their relationship.
// This is a generic comparator function that works with any type satisfying the NumericField constraint.
//
// Type parameters:
//   - T: The numeric type being compared, which must satisfy the NumericField constraint
//
// Returns:
//   - 1 if arg0 is greater than arg1
//   - 0 if arg0 is equal to arg1
//   - -1 if arg0 is less than arg1
//
// This implements the abstract.Comparator function signature for built-in numeric types.
func NumericComparator[T NumericField](arg0, arg1 T) int {
	if arg0 > arg1 {
		return 1
	} else if arg0 < arg1 {
		return -1
	}
	return 0
}

// BigNaturalComparator compares two BigInt instances and returns their relationship.
// This is a convenience function that delegates to the Compare method of BigInt.
//
// Returns:
//   - 1 if arg0 is greater than arg1
//   - 0 if arg0 is equal to arg1
//   - -1 if arg0 is less than arg1
//
// This implements the abstract.Comparator function signature for BigInt types.
func BigNaturalComparator(arg0, arg1 *BigInt) int {
	return arg0.Compare(arg1)
}

// BigDecimalComparator compares two BigFloat instances and returns their relationship.
// This is a convenience function that delegates to the Compare method of BigFloat.
//
// Returns:
//   - 1 if arg0 is greater than arg1
//   - 0 if arg0 is equal to arg1
//   - -1 if arg0 is less than arg1
//
// This implements the abstract.Comparator function signature for BigFloat types.
func BigDecimalComparator(arg0, arg1 *BigFloat) int {
	return arg0.Compare(arg1)
}
