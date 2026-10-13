/*
Copyright 2020 The Kubernetes Authors.

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

// Package tristate provides a helper for working with bool pointers.
package tristate

// withDefault evaluates a pointer to a bool with a default value.
func withDefault(def bool, b *bool) bool {
	if b == nil {
		return def
	}
	return *b
}

// EqualWithDefault compares two bool pointers using a default value.
func EqualWithDefault(def bool, a *bool, b *bool) bool {
	return withDefault(def, a) == withDefault(def, b)
}
