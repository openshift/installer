// Copyright 2015 Red Hat, Inc.
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

package util

import (
	"fmt"
	"strings"
)

func Highlight(srcb []byte, offset int64) (int, int, string) {
	src := string(srcb)

	startOfLine := strings.LastIndex(src[:offset], "\n") + 1
	if startOfLine == len(src) {
		startOfLine--
	}

	endOfLine := strings.Index(src[offset:], "\n")
	if endOfLine == -1 {
		endOfLine = len(src) - 1
	} else {
		endOfLine += int(offset)
	}

	line := strings.Count(src[:offset], "\n") + 1
	col := int(offset) - startOfLine + 1

	lineh := strings.Replace(src[startOfLine:endOfLine], "\t", " ", -1)
	highlight := fmt.Sprintf("%s\n%s^\n", lineh, strings.Repeat(" ", col-1))
	return line, col, highlight
}
