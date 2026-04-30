package pointer_helper

import (
	"reflect"
	"time"
)

// StringPtr returns a pointer to the given string.
// This is useful for assigning string literals to pointer variables
// without creating an intermediate variable.
func StringPtr(val string) *string {
	return &val
}

// IntPtr returns a pointer to the given int.
// This is useful for assigning int literals to pointer variables
// without creating an intermediate variable.
func IntPtr(val int) *int {
	return &val
}

// FloatPtr returns a pointer to the given float64.
// This is useful for assigning float64 literals to pointer variables
// without creating an intermediate variable.
func FloatPtr(val float64) *float64 {
	return &val
}

// BoolPtr returns a pointer to the given bool.
// This is useful for assigning bool literals to pointer variables
// without creating an intermediate variable.
func BoolPtr(val bool) *bool {
	return &val
}

// TimePtr returns a pointer to the given time.Time.
// This is useful for assigning time.Time literals to pointer variables
// without creating an intermediate variable.
func TimePtr(val time.Time) *time.Time {
	return &val
}

/*
TypePtr returns a pointer to the given value of a specified type.
The function supports types: string, int, float64, time.Time, and bool.

T: The type of the value, constrained to string, int, float64, time.Time, or bool.
val: The value to obtain a pointer for.

Returns a pointer to the provided value.
*/
func TypePtr[T string | int | float64 | time.Time | bool](val T) *T {
	return &val
}

/*
DeTypePtr dereferences a pointer to a value of type T and returns the value.
If the pointer is nil, it returns the zero value for the type T.

Type Parameters:
  - T: The type of the value, which can be string, int, float64, time.Time, or bool.

Parameters:
  - val: A pointer to a value of type T.

Returns:
  - The dereferenced value of type T, or the zero value if the pointer is nil.
*/
func DeTypePtr[T string | int | float64 | time.Time | bool](val *T) T {
	if !IsNilOrZeroVal(val) {
		return *val
	}

	rf := reflect.ValueOf(val)
	return reflect.Zero(rf.Type().Elem()).Interface().(T)
}

/*
IsNilOrZeroVal checks if the provided pointer value is either nil or points to a zero value.
It supports pointers to string, int, float64, time.Time, and bool types.

Parameters:
  - val: A pointer to one of the supported types.

Returns:
  - true if the pointer is nil or points to a zero value; otherwise, false.
*/
func IsNilOrZeroVal[T string | int | float64 | time.Time | bool](val *T) bool {
	if val == nil {
		return true
	}

	return reflect.ValueOf(val).Elem().IsZero()
}
