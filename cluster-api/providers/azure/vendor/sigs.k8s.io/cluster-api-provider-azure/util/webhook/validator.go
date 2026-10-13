/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhook

import (
	"reflect"
	"sort"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	unsetMessage     = "field is immutable, unable to set an empty value if it was already set"
	setMessage       = "field is immutable, unable to assign a value if it was already empty"
	immutableMessage = "field is immutable"
)

// ValidateImmutable validates equality across two values,
// and returns a meaningful error to indicate a changed value, a newly set value, or a newly unset value.
func ValidateImmutable(path *field.Path, oldVal, newVal any) *field.Error {
	if reflect.TypeOf(oldVal) != reflect.TypeOf(newVal) {
		return field.Invalid(path, newVal, "unexpected error")
	}
	if !reflect.ValueOf(oldVal).IsZero() {
		// Prevent modification if it was already set to some value
		if reflect.ValueOf(newVal).IsZero() {
			// unsetting the field is not allowed
			return field.Invalid(path, newVal, unsetMessage)
		}
		if !reflect.DeepEqual(oldVal, newVal) {
			// changing the field is not allowed
			return field.Invalid(path, newVal, immutableMessage)
		}
	} else if !reflect.ValueOf(newVal).IsZero() {
		return field.Invalid(path, newVal, setMessage)
	}

	return nil
}

// ValidateZeroTransition validates equality across two values, with only exception to allow
// the value to transition of a zero value.
func ValidateZeroTransition(path *field.Path, oldVal, newVal any) *field.Error {
	if reflect.ValueOf(newVal).IsZero() {
		// unsetting the field is allowed
		return nil
	}
	return ValidateImmutable(path, oldVal, newVal)
}

// EnsureStringSlicesAreEquivalent returns if two string slices have equal lengths,
// and that they have the exact same items; it does not enforce strict ordering of items.
func EnsureStringSlicesAreEquivalent(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	sort.Strings(a)
	sort.Strings(b)

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
