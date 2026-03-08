package tools

import "reflect"

// TypeName returns the fully qualified name of type T as a string.
//
// It uses reflection to obtain the type name at runtime and is primarily
// intended for logging, error messages, and generic type identification
// during initialization.
//
// The function is safe for any type, including pointers and interfaces.
// For pointer types, it returns the name of the underlying type without
// the pointer prefix.
//
// Example:
//
//	name := tools.TypeName[time.Time]() // returns "time.Time"
//	ptrName := tools.TypeName[*int]()   // returns "int"
func TypeName[T any]() string {
	return reflect.TypeOf((*T)(nil)).Elem().String()
}
