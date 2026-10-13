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

// Package cmp provides a set of comparison functions.
package cmp

import (
	"sort"

	"k8s.io/utils/ptr"
)

// ByPtrValue is a type to sort a slice of pointers to strings.
type ByPtrValue []*string

// Len returns the length of the slice.
func (s ByPtrValue) Len() int {
	return len(s)
}

// Swap swaps the elements with indexes i and j.
func (s ByPtrValue) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less returns true if the element with index i should sort before the element with index j.
func (s ByPtrValue) Less(i, j int) bool {
	return *s[i] < *s[j]
}

// Equals returns true if the two slices of pointers to strings are equal.
func Equals(slice1, slice2 []*string) bool {
	sort.Sort(ByPtrValue(slice1))
	sort.Sort(ByPtrValue(slice2))

	if len(slice1) == len(slice2) {
		for i, v := range slice1 {
			if !ptr.Equal(v, slice2[i]) {
				return false
			}
		}
	} else {
		return false
	}
	return true
}
