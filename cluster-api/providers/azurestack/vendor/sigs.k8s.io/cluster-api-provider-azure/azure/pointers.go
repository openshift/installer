/*
Copyright 2023 The Kubernetes Authors.

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

package azure

import "k8s.io/utils/ptr"

// StringSlice returns a string slice value for the passed string slice pointer. It returns a nil
// slice if the pointer is nil.
func StringSlice(s *[]string) []string {
	if s != nil {
		return *s
	}
	return nil
}

// PtrSlice returns a slice of pointers from a pointer to a slice. It returns nil if the
// pointer is nil or the slice pointed to is empty.
func PtrSlice[T any](p *[]T) []*T {
	if p == nil || len(*p) == 0 {
		return nil
	}
	s := make([]*T, 0, len(*p))
	for _, v := range *p {
		s = append(s, ptr.To(v))
	}
	return s
}

// AliasOrNil returns a pointer to a string-derived type from a passed string pointer,
// or nil if the pointer is nil or an empty string.
func AliasOrNil[T ~string](s *string) *T {
	if s == nil || *s == "" {
		return nil
	}
	return ptr.To(T(*s))
}

// StringMapPtr converts a map[string]string into a map[string]*string. It returns nil if the map is nil.
func StringMapPtr(m map[string]string) map[string]*string {
	if m == nil {
		return nil
	}
	msp := make(map[string]*string, len(m))
	for k, v := range m {
		msp[k] = ptr.To(v)
	}
	return msp
}
