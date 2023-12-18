/*
Copyright 2017 The Kubernetes Authors.

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

package util

import (
	"fmt"
	"strings"
)

// Filter filters a list for a string.
func Filter(list []string, strToFilter string) (newList []string) {
	for _, item := range list {
		if item != strToFilter {
			newList = append(newList, item)
		}
	}
	return
}

// Contains returns true if a list contains a string.
func Contains(list []string, strToSearch string) bool {
	for _, item := range list {
		if item == strToSearch {
			return true
		}
	}
	return false
}

// MergeCommaSeparatedKeyValuePairs merges multiple comma separated lists of key=value pairs into a single, comma-separated, list
// of key=value pairs. If a key is present in multiple lists, the value from the last list is used.
func MergeCommaSeparatedKeyValuePairs(lists ...string) string {
	merged := make(map[string]string)
	for _, list := range lists {
		for _, kv := range strings.Split(list, ",") {
			kv := strings.Split(kv, "=")
			if len(kv) != 2 {
				// ignore invalid key=value pairs
				continue
			}
			merged[kv[0]] = kv[1]
		}
	}
	// convert the map back to a comma separated list
	var result []string
	for k, v := range merged {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(result, ",")
}
