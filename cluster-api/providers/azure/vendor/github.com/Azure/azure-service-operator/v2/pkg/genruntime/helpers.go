/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package genruntime

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

// NewObjectFromExemplar creates a new client.Object with the same GVK as the provided client.Object.
// The supplied client.Object is not changed and the returned client.Object is empty.
func NewObjectFromExemplar(obj client.Object, scheme *runtime.Scheme) (client.Object, error) {
	gvk, err := apiutil.GVKForObject(obj, scheme)
	if err != nil {
		return nil, err
	}

	// Create a fresh destination to deserialize to
	newObj, err := scheme.New(gvk)
	if err != nil {
		return nil, err
	}

	// Ensure GVK is populated
	newObj.GetObjectKind().SetGroupVersionKind(gvk)

	return newObj.(client.Object), nil
}

// InterleaveStrSlice interleaves the elements of the two provided slices. The resulting slice looks like:
// []{<element 1 from a>, <element 1 from b>, <element 2 from a>, <element 2 from b>...}. If one slice is longer than
// the other, the elements are interleaved until the shorter slice is out of elements, at which point all remaining
// elements are from the longer slice.
func InterleaveStrSlice(a []string, b []string) []string {
	smallestLen := MinInt(len(a), len(b))
	var larger []string
	if len(a) == smallestLen {
		larger = b
	} else {
		larger = a
	}

	var result []string

	for i := 0; i < smallestLen; i++ {
		result = append(result, a[i])
		result = append(result, b[i])
	}

	result = append(result, larger[smallestLen:]...)

	return result
}

// MinInt returns the minimum of the two provided ints.
// The fact that this doesn't exist in the Go standard library is depressing.
func MinInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

// ARMSpecNames returns a slice of names from the given ARMResourceSpec slice.
func ARMSpecNames(specs []ARMResourceSpec) []string {
	result := make([]string, len(specs))
	for ix := range specs {
		result[ix] = specs[ix].GetName()
	}

	return result
}

// ARMSpecNames returns a slice of names from the given ARMResourceSpec slice.
func RawNames(specs []any) []string {
	result := make([]string, len(specs))
	for ix := range specs {
		m, ok := specs[ix].(map[string]any)
		if !ok {
			result[ix] = "<UNKNOWN>"
			continue
		}

		name, ok := m["name"]
		if !ok {
			result[ix] = "<UNKNOWN>"
			continue
		}

		nameStr, ok := name.(string)
		if !ok {
			result[ix] = "<UNKNOWN>"
			continue
		}

		result[ix] = nameStr
	}

	return result
}
