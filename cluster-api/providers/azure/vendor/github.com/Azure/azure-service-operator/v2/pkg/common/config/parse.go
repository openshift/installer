// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package config

import (
	"strings"
)

// ParseCommaCollection splits a comma-separated string into a slice
// of strings with spaces trimmed.
func ParseCommaCollection(fromEnv string) []string {
	if len(strings.TrimSpace(fromEnv)) == 0 {
		return nil
	}
	items := strings.Split(fromEnv, ",")
	// Remove any whitespace used to separate items.
	for i, item := range items {
		items[i] = strings.TrimSpace(item)
	}
	return items
}
