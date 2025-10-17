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

package v1beta1

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestTags_Merge(t *testing.T) {
	tests := []struct {
		name     string
		other    Tags
		expected Tags
	}{
		{
			name:  "nil other",
			other: nil,
			expected: Tags{
				"a": "b",
				"c": "d",
			},
		},
		{
			name:  "empty other",
			other: Tags{},
			expected: Tags{
				"a": "b",
				"c": "d",
			},
		},
		{
			name: "disjoint",
			other: Tags{
				"1": "2",
				"3": "4",
			},
			expected: Tags{
				"a": "b",
				"c": "d",
				"1": "2",
				"3": "4",
			},
		},
		{
			name: "overlapping, other wins",
			other: Tags{
				"1": "2",
				"3": "4",
				"a": "hello",
			},
			expected: Tags{
				"a": "hello",
				"c": "d",
				"1": "2",
				"3": "4",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)
			tags := Tags{
				"a": "b",
				"c": "d",
			}

			tags.Merge(tc.other)
			g.Expect(tags).To(Equal(tc.expected))
		})
	}
}
