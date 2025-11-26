/*
Copyright 2025 The Kubernetes Authors.

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

// Package webhook implements reusable validation functions for webhooks.
package webhook

import (
	"reflect"

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

// ValidateNonNegative checks if a numeric value is non-negative,
// and returns a meaningul error to indicate that it is.
func ValidateNonNegative[T int | int32 | int64](path *field.Path, value *T) *field.Error {
	if value != nil && *value < 0 {
		return field.Invalid(path, value, "must be non-negative")
	}

	return nil
}
