/*
Copyright (c) 2021 Red Hat, Inc.

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

// This file contains functions useful for printing lists in a format that is easy to read for
// humans in log files.

package logging

import (
	"fmt"
	"strings"
)

// All generates a human readable representation of the given list of strings, for use in a log
// file. It puts quotes around each item, separates the first items with commas and the last with
// the word 'and'.
func All(items []string) string {
	return list(items, "and")
}

// any generates a human readable representation of the given list of strings, for use in a log
// file. It puts quotes around each item, separates the first items with commas and the last with
// the word 'or'.
func Any(items []string) string {
	return list(items, "or")
}

func list(items []string, conjunction string) string {
	count := len(items)
	if count == 0 {
		return ""
	}
	quoted := make([]string, len(items))
	for i, item := range items {
		quoted[i] = fmt.Sprintf("'%s'", item)
	}
	if count == 1 {
		return quoted[0]
	}
	head := quoted[0 : count-1]
	tail := quoted[count-1]
	return fmt.Sprintf("%s %s %s", strings.Join(head, ", "), conjunction, tail)
}
