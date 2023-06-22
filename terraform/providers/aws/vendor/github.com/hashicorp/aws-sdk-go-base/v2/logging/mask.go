// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logging

import (
	"strings"
)

func partialMaskString(s string, first, last int) string {
	l := len(s)
	result := s[0:first]
	result += strings.Repeat("*", l-first-last)
	result += s[l-last:]
	return result
}
