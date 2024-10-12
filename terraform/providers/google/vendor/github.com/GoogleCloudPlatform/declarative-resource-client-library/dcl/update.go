// Copyright 2024 Google LLC. All Rights Reserved.
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
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"bitbucket.org/creachadair/stringset"
)

// UpdateMask creates a Update Mask string according to https://google.aip.dev/161
func UpdateMask(ds []*FieldDiff) string {
	var ss []string
	for _, v := range ds {
		ss = append(ss, convertUpdateMaskVal(v.FieldName))
	}

	dupesRemoved := stringset.New(ss...).Elements()

	// Sorting the entries is optional, but makes it easier to read + test.
	sort.Strings(dupesRemoved)
	return strings.Join(dupesRemoved, ",")
}

func titleCaseToCamelCase(s string) string {
	r, n := utf8.DecodeRuneInString(s)
	p := string(unicode.ToLower(r))
	p = p + s[n:]
	return p
}

// Diffs come in the form Http.AuthInfo.Password
// Needs to be in the form http.authInfo.password
func convertUpdateMaskVal(s string) string {
	r := regexp.MustCompile(`\[\d\]`)

	// camelCase string (right now, it's in TitleCase).
	parts := strings.Split(s, ".")
	var p []string
	for _, q := range parts {
		if r.MatchString(q) {
			// Indexing into a repeated field.
			bareFieldName := r.ReplaceAllString(q, "")
			p = append(p, titleCaseToCamelCase(bareFieldName))

			// Repeated fields cannot be intermediary in a field mask, so we
			// must terminate the field mask here.
			break
		} else {
			p = append(p, titleCaseToCamelCase(q))
		}
	}

	// * notation should only be used if this is not the last field.
	// Example: res.array.* should be res.array, but res.array.*.bar means "update only bar in all my array fields"
	if p[len(p)-1] == "*" {
		p = p[0 : len(p)-1]
	}

	return strings.Join(p, ".")
}

// TopLevelUpdateMask returns only the top-level fields.
func TopLevelUpdateMask(ds []*FieldDiff) string {
	var ss []string
	for _, v := range ds {
		part := strings.Split(v.FieldName, ".")[0]
		ss = append(ss, convertUpdateMaskVal(part))
	}

	dupesRemoved := stringset.New(ss...).Elements()

	// Sorting the entries is optional, but makes it easier to read + test.
	sort.Strings(dupesRemoved)
	return strings.Join(dupesRemoved, ",")
}

// SnakeCaseUpdateMask returns the update mask, but all fields are snake case.
func SnakeCaseUpdateMask(ds []*FieldDiff) string {
	var ss []string
	for _, v := range ds {
		ss = append(ss, TitleToSnakeCase(convertUpdateMaskVal(v.FieldName)))
	}
	dupesRemoved := stringset.New(ss...).Elements()

	// Sorting the entries is optional, but makes it easier to read + test.
	sort.Strings(dupesRemoved)
	return strings.Join(dupesRemoved, ",")
}

// UpdateMaskWithPrefix returns a Standard Update Mask with a prefix attached.
func UpdateMaskWithPrefix(ds []*FieldDiff, prefix string) string {
	um := UpdateMask(ds)
	parts := strings.Split(um, ",")

	var ss []string

	for _, part := range parts {
		ss = append(ss, fmt.Sprintf("%s.%s", prefix, part))
	}

	return strings.Join(ss, ",")
}
