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

package cmp

import (
	"sort"

	"k8s.io/utils/pointer"
)

type ByPtrValue []*string

func (s ByPtrValue) Len() int {
	return len(s)
}

func (s ByPtrValue) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByPtrValue) Less(i, j int) bool {
	return *s[i] < *s[j]
}

func Equals(slice1, slice2 []*string) bool {
	sort.Sort(ByPtrValue(slice1))
	sort.Sort(ByPtrValue(slice2))

	if len(slice1) == len(slice2) {
		for i, v := range slice1 {
			if !pointer.StringEqual(v, slice2[i]) {
				return false
			}
		}
	} else {
		return false
	}
	return true
}
