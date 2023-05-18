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
	"regexp"
	"strings"
)

// Map from initialism -> TitleCase variant
// We can assume camelCase is the same as TitleCase except that we downcase the
// first segment
var initialisms = map[string]string{
	"ai":      "AI",
	"ip":      "IP",
	"os":      "OS",
	"ipv4":    "IPv4",
	"ipv6":    "IPv6",
	"oauth":   "OAuth",
	"oauth2":  "OAuth2",
	"tpu":     "TPU",
	"vpc":     "VPC",
	"v1beta1": "V1Beta1",
}

// SnakeToTitleCase converts a snake_case string to TitleCase / Go struct case.
func SnakeToTitleCase(s string) string {
	return strings.Join(SnakeToTitleParts(s), "")
}

// SnakeToJSONCase converts a snake_case string to jsonCase / camelCase, for use
// in JSON serialization.
func SnakeToJSONCase(s string) string {
	parts := SnakeToTitleParts(s)
	if len(parts) > 0 {
		parts[0] = strings.ToLower(parts[0])
	}

	return strings.Join(parts, "")
}

// SnakeToTitleParts returns the parts of a snake_case string titlecased as an
// array, taking into account common initialisms.
func SnakeToTitleParts(s string) []string {
	parts := []string{}
	segments := strings.Split(s, "_")
	for _, seg := range segments {
		if v, ok := initialisms[seg]; ok {
			parts = append(parts, v)
		} else {
			if len(seg) < 1 {
				continue
			}
			parts = append(parts, strings.ToUpper(seg[0:1])+seg[1:])
		}
	}

	return parts
}

// SnakeToTitleCasePath converts a resource path from snake to title case. For
// example: foo_bar.baz.qux -> FooBar.Baz.Qux
func SnakeToTitleCasePath(s, sep string) string {
	str := []string{}
	for _, p := range strings.Split(s, sep) {
		str = append(str, SnakeToTitleCase(p))
	}
	return strings.Join(str, sep)
}

// TitleToCamelCasePath converts a resource path from title case to lower title case.
// For example: FooBar.Baz.Qux -> fooBar.baz.qux
func TitleToCamelCasePath(s string) string {
	// Lowercase the first character and every character following a .
	parts := strings.Split(s, ".")
	for i, part := range parts {
		parts[i] = strings.ToLower(part[:1]) + part[1:]
	}
	return strings.Join(parts, ".")
}

// ProtoCamelCase converts a snake case name to a upper camel case name using the
// go protoc special rules: convert to camel case, except when
// the character following the underscore is a digit; e.g.,
// foo_bar_2 -> FooBar_2.
// From: http://google3/net/goa/codegen/names.go;l=14;rcl=294425921
func ProtoCamelCase(s string) string {
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	var b []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '.' && i+1 < len(s) && isASCIILower(s[i+1]):
			// Skip over '.' in ".{{lowercase}}".
		case c == '.':
			b = append(b, '_') // convert '.' to '_'
		case c == '_' && (i == 0 || s[i-1] == '.'):
			// Convert initial '_' to ensure we start with a capital letter.
			// Do the same for '_' after '.' to match historic behavior.
			b = append(b, 'X') // convert '_' to 'X'
		case c == '_' && i+1 < len(s) && isASCIILower(s[i+1]):
			// Skip over '_' in "_{{lowercase}}".
		case isASCIIDigit(c):
			b = append(b, c)
		default:
			// Assume we have a letter now - if not, it's a bogus identifier.
			// The next word is a sequence of characters that must start upper case.
			if isASCIILower(c) {
				c -= 'a' - 'A' // convert lowercase to uppercase
			}
			b = append(b, c)

			// Accept lower case sequence that follows.
			for ; i+1 < len(s) && isASCIILower(s[i+1]); i++ {
				b = append(b, s[i+1])
			}
		}
	}
	return string(b)
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// TitleToSnakeCase takes in a TitleCase string and returns a snake_case string.
func TitleToSnakeCase(s string) string {
	for k, v := range initialisms {
		kCap := strings.ToUpper(k[0:1]) + k[1:]
		s = strings.Replace(s, v, kCap, -1)
	}
	str := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(str, "${1}_${2}"))
}

// StringSliceContains returns true if the slice ss contains string s.
func StringSliceContains(s string, ss []string) bool {
	for _, st := range ss {
		if st == s {
			return true
		}
	}
	return false
}
