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

package slice

import "k8s.io/utils/ptr"

// Contains tells whether a Contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// ToPtrs takes a slice of values and returns a slice of pointers to the same values.
func ToPtrs[T any](in []T) []*T {
	if in == nil {
		return nil
	}
	ptrs := make([]*T, 0, len(in))
	for i := range in {
		ptrs = append(ptrs, ptr.To(in[i]))
	}
	return ptrs
}
