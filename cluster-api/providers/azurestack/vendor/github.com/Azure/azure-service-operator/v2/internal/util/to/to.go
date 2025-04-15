// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package to

// Value extracts a value from a pointer. If the pointer is nil, the default value is returned.
func Value[T any](v *T) T {
	if v == nil {
		var t T
		return t
	}

	return *v
}

// Ptr returns a pointer to the provided value.
func Ptr[T any](v T) *T {
	return &v
}
