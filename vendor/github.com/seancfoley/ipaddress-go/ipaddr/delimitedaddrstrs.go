//
// Copyright 2020-2022 Sean C Foley
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
//

package ipaddr

import "strings"

const SegmentValueDelimiter = ','

type DelimitedAddressString string

// CountDelimitedAddresses will count the possible combinations, given a string with comma delimiters separating segment elements.
// It is a counterpart to ParseDelimitedSegments, indicating the number of iterated elements from ParseDelimitedSegments.
//
// For example, given "1,2.3.4,5.6" this method will return 4 for the possible combinations: "1.3.4.6", "1.3.5.6", "2.3.4.6" and "2.3.5.6".
func (str DelimitedAddressString) CountDelimitedAddresses() int {
	segDelimitedCount := 0
	result := 1
	strlen := len(str)
	for i := 0; i < strlen; i++ {
		c := str[i]
		if isDelimitedBoundary(c) {
			if segDelimitedCount > 0 {
				result *= segDelimitedCount + 1
				segDelimitedCount = 0
			}
		} else if c == SegmentValueDelimiter {
			segDelimitedCount++
		}
	}
	if segDelimitedCount > 0 {
		result *= segDelimitedCount + 1
	}
	return result
}

func isDelimitedBoundary(c byte) bool {
	return c == IPv4SegmentSeparator ||
		c == IPv6SegmentSeparator ||
		c == RangeSeparator ||
		c == MacDashedSegmentRangeSeparator
}

// ParseDelimitedSegments will provide an iterator to iterate through the possible combinations, given a string with comma delimiters to denote segment elements,
//
// For example, given "1,2.3.4,5.6" this will iterate through "1.3.4.6", "1.3.5.6", "2.3.4.6" and "2.3.5.6"
//
// Another example: "1-2,3.4.5.6" will iterate through "1-2.4.5.6" and "1-3.4.5.6".
//
// This method will not validate strings.  Each string produced can be validated using an instance of [IPAddressString].
// Use CountDelimitedAddresses for the count of elements in the iterator.
func (str DelimitedAddressString) ParseDelimitedSegments() Iterator[string] {
	var parts [][]string
	var lastSegmentStartIndex, lastPartIndex, lastDelimiterIndex int
	anyDelimited := false
	var delimitedList []string
	strlen := len(str)
	s := string(str)
	for i := 0; i < strlen; i++ {
		c := str[i]
		if isDelimitedBoundary(c) { // end of segment or range boundary
			if len(delimitedList) > 0 {
				if parts == nil {
					parts = make([][]string, 0, IPv6SegmentCount)
				}
				parts, _ = addParts(s, parts, lastSegmentStartIndex, lastPartIndex, lastDelimiterIndex, delimitedList, i)
				lastPartIndex = i
				delimitedList = delimitedList[:0]
			}
			lastDelimiterIndex = i + 1
			lastSegmentStartIndex = lastDelimiterIndex
		} else if c == SegmentValueDelimiter {
			anyDelimited = true
			if delimitedList == nil {
				delimitedList = make([]string, 0, 4)
			}
			sub := str[lastDelimiterIndex:i]
			delimitedList = append(delimitedList, string(sub))
			lastDelimiterIndex = i + 1
		}
	}
	if anyDelimited {
		if len(delimitedList) > 0 {
			if parts == nil {
				parts = make([][]string, 0, IPv6SegmentCount)
			}
			parts, _ = addParts(s, parts, lastSegmentStartIndex, lastPartIndex, lastDelimiterIndex, delimitedList, len(str))
		} else {
			parts = append(parts, []string{s[lastPartIndex:]})
		}
		return newDelimitedStringsIterator(parts)
	}
	return newSingleStrIterator(s)
}

// ParseDelimitedIPAddrSegments will provide an iterator to iterate through the possible combinations, given a string with comma delimiters to denote segment elements.
func (str DelimitedAddressString) ParseDelimitedIPAddrSegments() Iterator[*IPAddressString] {
	return ipAddressStringIterator{str.ParseDelimitedSegments()}
}

func addParts(str string, parts [][]string, lastSegmentStartIndex, lastPartIndex,
	lastDelimiterIndex int, delimitedList []string, i int) (newParts [][]string, newDelimitedList []string) {
	sub := str[lastDelimiterIndex:i]
	delimitedList = append(delimitedList, sub)
	if lastPartIndex != lastSegmentStartIndex {
		parts = append(parts, []string{str[lastPartIndex:lastSegmentStartIndex]})
	}
	parts = append(parts, delimitedList)
	return parts, delimitedList
}

func newDelimitedStringsIterator(parts [][]string) Iterator[string] {
	partCount := len(parts)
	it := &delimitedStringsIterator{
		parts:      parts,
		variations: make([]Iterator[string], partCount),
		nextSet:    make([]string, partCount),
	}
	it.updateVariations(0)
	return it
}

type delimitedStringsIterator struct {
	parts      [][]string
	done       bool
	variations []Iterator[string]
	nextSet    []string
}

func (it *delimitedStringsIterator) updateVariations(start int) {
	variationLen := len(it.variations)
	variations := it.variations
	parts := it.parts
	nextSet := it.nextSet
	for i := start; i < variationLen; i++ {
		strSlice := parts[i]
		if len(strSlice) > 1 {
			variations[i] = newStrSliceIterator(strSlice)
		} else {
			variations[i] = newSingleStrIterator(strSlice[0])
		}
		nextSet[i] = variations[i].Next()
	}
}

func (it *delimitedStringsIterator) HasNext() bool {
	return !it.done
}

func (it *delimitedStringsIterator) Next() (res string) {
	if !it.done {
		result := strings.Builder{}
		nextSet := it.nextSet
		nextSetLen := len(nextSet)
		for i := 0; i < nextSetLen; i++ {
			result.WriteString(nextSet[i])
		}
		it.increment()
		res = result.String()
	}
	return
}

func (it *delimitedStringsIterator) increment() {
	variations := it.variations
	variationsLen := len(variations)
	nextSet := it.nextSet
	for j := variationsLen - 1; j >= 0; j-- {
		if variations[j].HasNext() {
			nextSet[j] = variations[j].Next()
			it.updateVariations(j + 1)
			return
		}
	}
	it.done = true
}

func newStrSliceIterator(strs []string) Iterator[string] {
	return &stringIterator{strs: strs}
}

type stringIterator struct {
	strs []string
}

func (it *stringIterator) HasNext() bool {
	return len(it.strs) > 0
}

func (it *stringIterator) Next() (res string) {
	if it.HasNext() {
		strs := it.strs
		res = strs[0]
		it.strs = strs[1:]
	}
	return
}

//
type singleStringIterator struct {
	str  string
	done bool
}

func (it *singleStringIterator) HasNext() bool {
	return !it.done
}

func (it *singleStringIterator) Next() (res string) {
	if it.HasNext() {
		it.done = true
		res = it.str
	}
	return
}

func newSingleStrIterator(str string) Iterator[string] {
	return &singleStringIterator{str: str}
}

//
type ipAddressStringIterator struct {
	Iterator[string]
}

func (iter ipAddressStringIterator) Next() *IPAddressString {
	if !iter.HasNext() {
		return nil
	}
	return NewIPAddressString(iter.Iterator.Next())
}
