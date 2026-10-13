package utils

import (
	"fmt"
	"time"
)

// time.Time isn't really a primitive, but it's only pointer is to Location which would be the same one after the copy
type primitive interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64 | ~string | ~bool | time.Time
}

// New takes any type of primitive and returns a pointer to a copy that is allocated in the heap.
// This is needed to allow setting values (hard coded or struct members) by reference to the api structs.
func New[T primitive](value T) *T {
	return &value
}

// NewFormat is equivalent to Sprintf, but it returns a pointer to a string allocated in the heap.
func NewFormat(format string, a ...interface{}) *string {
	value := fmt.Sprintf(format, a...)
	return &value
}

// NewByteArray takes the given byte array and returns a pointer to a copy that is allocated in the heap.
func NewByteArray(value []byte) []byte {
	if len(value) > 0 {
		ret := make([]byte, len(value))
		copy(ret, value)
		return ret
	}
	return nil
}

// NewStringArray takes a string array and returns an array whose string elements are allocated in the heap.
func NewStringArray(array []string) []*string {
	if len(array) > 0 {
		ret := make([]*string, len(array))
		for idx, str := range array {
			ret[idx] = New(str)
		}
		return ret
	}
	return nil
}
