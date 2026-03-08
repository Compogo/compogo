package tools

import "reflect"

// Pointer returns a pointer to the given value.
// Useful for creating pointers in one line, e.g.: tools.Pointer("hello")
func Pointer[T any](value T) *T {
	return &value
}

// PointerOrNil returns a pointer to the value if it's not the zero value of its type.
// If the value is the zero value, returns nil.
// Uses reflect.ValueOf(value).IsZero() for comparison.
func PointerOrNil[T any](value T) *T {
	if reflect.ValueOf(value).IsZero() {
		return nil
	}

	return &value
}

// PointerValue safely returns the value from a pointer.
// If the pointer is nil, returns the zero value of type T.
func PointerValue[T any](value *T) T {
	if value == nil {
		return *new(T)
	}

	return *value
}
