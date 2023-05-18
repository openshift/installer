// Copyright 2023 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package dcl

import (
	"strings"
)

// SelfLinkToName returns the element of a string after the last slash.
func SelfLinkToName(v *string) *string {
	if v == nil {
		return nil
	}
	val := *v
	comp := strings.Split(val, "/")
	ret := comp[len(comp)-1]
	return &ret
}

// SelfLinkToNameExpander returns the element of a string after the last slash.
// Return value also has error since the dcl template requires the expander to return error.
func SelfLinkToNameExpander(v *string) (*string, error) {
	return SelfLinkToName(v), nil
}

// SelfLinkToNameArrayExpander returns the last element of each string in a slice after the last slash.
// Return value also has error since the dcl template requires the expander to return error.
func SelfLinkToNameArrayExpander(v []string) ([]string, error) {
	r := make([]string, len(v))
	for i, w := range v {
		r[i] = *SelfLinkToName(&w)
	}
	return r, nil
}

// FalseToNil returns nil if the pointed-to boolean is 'false' - otherwise returns the pass-in pointer.
func FalseToNil(b *bool) (*bool, error) {
	if b != nil && *b == false {
		return nil, nil
	}
	return b, nil
}

// SelfLinkToNameArray returns a slice of the elements of a slice of strings after the last slash.
func SelfLinkToNameArray(v []string) []string {
	var a []string
	for _, vv := range v {
		ret := SelfLinkToName(&vv)
		if ret != nil {
			a = append(a, *ret)
		}
	}
	return a
}

// SelfLinkToNameWithPattern handles when the resource name can have `/` in it
// by matching the pattern.
func SelfLinkToNameWithPattern(v *string, pattern string) *string {
	if v == nil {
		return nil
	}
	regex, err := regexFromPattern(pattern)
	if err != nil {
		// Unable to compile regex, best guess return v
		return v
	}
	matches := regex.FindStringSubmatch(*v)
	if len(matches) == 0 {
		return v
	}
	return &matches[len(matches)-1]
}
