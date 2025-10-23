//
// Copyright 2020-2024 Sean C Foley
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

import "sort"

type mergeableType[S any, T any] interface {
	*S

	AddressSegmentSeries

	IsSequential() bool
	IsSinglePrefixBlock() bool
	SequentialBlockIterator() Iterator[T]
	SpanWithPrefixBlocks() []T
	ToPrefixBlockLen(BitCount) T

	WithoutPrefixLen() T
	ToBlock(segmentIndex int, lower, upper SegInt) T
}

func getMergedPrefixBlocks[S any, T mergeableType[S, T]](sections []T) []T {
	singleElement, list := organizeSequentially(sections)
	if singleElement {
		return list
	}
	first := sections[0]
	bitCount := first.GetBitCount()
	bitsPerSegment := first.GetBitsPerSegment()
	bytesPerSegment := first.GetBytesPerSegment()

	// Now we see if we can match blocks or join them into larger blocks
	removedCount := 0
	listLen := len(list)
	j := listLen - 1
	i := j - 1
top:
	for j > 0 {
		item := list[i]
		otherItem := list[j]
		compare := ReverseHighValueComparator.CompareSeries(item, otherItem)
		// check for strict containment, case 1:
		// w   z
		//  x y

		if compare > 0 {
			removedCount++
			k := j + 1
			for k < listLen && list[k] == nil {
				k++
			}
			if k < listLen {
				list[j] = list[k]
				list[k] = nil
			} else {
				list[j] = nil
				j = i
				i--
			}
			continue
		}
		// non-strict containment, case 2:
		// w   z
		// w   z
		//
		// reverse containment, case 3:
		// w  y
		// w   z
		rcompare := ReverseLowValueComparator.CompareSeries(item, otherItem)
		if rcompare >= 0 {
			removedCount++
			list[i] = otherItem
			list[j] = nil
			j = i
			i--
			continue
		}
		// check for merge, case 4:
		// w   x
		//      y   z
		// where x and y adjacent, becoming:
		// w        z
		//
		prefixLen := item.GetPrefixLen()
		otherPrefixLen := otherItem.GetPrefixLen()
		if !prefixLen.Equal(otherPrefixLen) {
			j = i
			i--
			continue
		}
		var matchBitIndex BitCount
		if prefixLen == nil {
			matchBitIndex = bitCount - 1
		} else {
			matchBitIndex = prefixLen.bitCount() - 1
		}
		var lastMatchSegmentIndex, lastBitSegmentIndex int
		if matchBitIndex != 0 {
			lastMatchSegmentIndex = getNetworkSegmentIndex(matchBitIndex, bytesPerSegment, bitsPerSegment)
			lastBitSegmentIndex = getHostSegmentIndex(matchBitIndex, bytesPerSegment, bitsPerSegment)
		}
		itemSegment := item.GetGenericSegment(lastMatchSegmentIndex)
		otherItemSegment := otherItem.GetGenericSegment(lastMatchSegmentIndex)
		itemSegmentValue := itemSegment.GetSegmentValue()
		otherItemSegmentValue := otherItemSegment.GetSegmentValue()
		segmentLastBitIndex := bitsPerSegment - 1
		if lastBitSegmentIndex == lastMatchSegmentIndex {
			segmentBitToCheck := matchBitIndex % bitsPerSegment
			shift := segmentLastBitIndex - segmentBitToCheck
			itemSegmentValue >>= uint(shift)
			otherItemSegmentValue >>= uint(shift)
		} else {
			itemBitValue := item.GetGenericSegment(lastBitSegmentIndex).GetSegmentValue()
			otherItemBitalue := otherItem.GetGenericSegment(lastBitSegmentIndex).GetSegmentValue()

			//we will make space for the last bit so we can do a single comparison
			itemSegmentValue = (itemSegmentValue << 1) | (itemBitValue >> uint(segmentLastBitIndex))
			otherItemSegmentValue = (otherItemSegmentValue << 1) | (otherItemBitalue >> uint(segmentLastBitIndex))
		}
		if itemSegmentValue != otherItemSegmentValue {
			itemSegmentValue ^= 1 //the ^ 1 flips the first bit
			if itemSegmentValue != otherItemSegmentValue {
				//neither an exact match nor a match when flipping the bit, so move on
				j = i
				i--
				continue
			} //else we will merge these two into a single prefix block, presuming the initial segments match
		}
		//check initial segments
		for k := lastMatchSegmentIndex - 1; k >= 0; k-- {
			itemSegment = item.GetGenericSegment(k)
			otherItemSegment = otherItem.GetGenericSegment(k)
			val := itemSegment.GetSegmentValue()
			otherVal := otherItemSegment.GetSegmentValue()
			if val != otherVal {
				j = i
				i--
				continue top
			}
		}
		joinedItem := otherItem.ToPrefixBlockLen(matchBitIndex)
		list[i] = joinedItem
		removedCount++
		k := j + 1
		for k < listLen && list[k] == nil {
			k++
		}
		if k < listLen {
			list[j] = list[k]
			list[k] = nil
		} else {
			list[j] = nil
			j = i
			i--
		}
	}
	if removedCount > 0 {
		newSize := listLen - removedCount
		for k, l := 0, 0; k < newSize; k, l = k+1, l+1 {
			for list[l] == nil {
				l++
			}
			if k != l {
				list[k] = list[l]
			}
		}
		list = list[:newSize]
	}
	return list
}

func getMergedSequentialBlocks[S any, T mergeableType[S, T]](sections []T) []T {
	singleElement, list := organizeSequentialMerge(sections)
	if singleElement {
		list[0] = list[0].WithoutPrefixLen()
		return list
	}
	removedCount := 0
	j := len(list) - 1
	i := j - 1
	ithRangeSegmentIndex, jthRangeSegmentIndex := -1, -1
top:
	for j > 0 {
		item := list[i]
		otherItem := list[j]
		compare := ReverseHighValueComparator.Compare(item, otherItem)
		// check for strict containment, case 1:
		// w   z
		//  x y
		if compare > 0 {
			removedCount++
			k := j + 1
			for k < len(list) && list[k] == nil {
				k++
			}
			if k < len(list) {
				list[j] = list[k]
				list[k] = nil
				jthRangeSegmentIndex = -1
			} else {
				list[j] = nil
				j = i
				i--
				jthRangeSegmentIndex = ithRangeSegmentIndex
				ithRangeSegmentIndex = -1
			}
			continue
		}
		// non-strict containment, case 2:
		// w   z
		// w   z
		//
		// reverse containment, case 3:
		// w  y
		// w   z
		rcompare := ReverseLowValueComparator.Compare(item, otherItem)
		if rcompare >= 0 {
			removedCount++
			list[i] = otherItem
			list[j] = nil
			j = i
			i--
			jthRangeSegmentIndex = ithRangeSegmentIndex
			ithRangeSegmentIndex = -1
			continue
		}

		//check for overlap

		if ithRangeSegmentIndex < 0 {
			ithRangeSegmentIndex = item.GetSequentialBlockIndex()
		}
		if jthRangeSegmentIndex < 0 {
			jthRangeSegmentIndex = otherItem.GetSequentialBlockIndex()
		}

		// check for overlap in the non-full range segment,
		// which must be the same segment in both, otherwise it cannot be overlap,
		// it can only be containment.
		// The one with the earlier range segment can only contain the other, there cannot be overlap.
		// eg 1.a-b.*.* and 1.2.3.* overlap in range segment 2 must have a <= 2 <= b and that means 1.a-b.*.* contains 1.2.3.*
		if ithRangeSegmentIndex != jthRangeSegmentIndex {
			j = i
			i--
			jthRangeSegmentIndex = ithRangeSegmentIndex
			ithRangeSegmentIndex = -1
			continue
		}

		rangeSegment := item.GetGenericSegment(ithRangeSegmentIndex)
		otherRangeSegment := otherItem.GetGenericSegment(ithRangeSegmentIndex)
		otherRangeItemValue := otherRangeSegment.GetSegmentValue()
		rangeItemUpperValue := rangeSegment.GetUpperSegmentValue()

		//check for overlapping range in the range segment
		if rangeItemUpperValue < otherRangeItemValue && rangeItemUpperValue+1 != otherRangeItemValue {
			j = i
			i--
			ithRangeSegmentIndex = -1
			continue
		}

		// now check all previous segments match
		for k := ithRangeSegmentIndex - 1; k >= 0; k-- {
			itemSegment := item.GetGenericSegment(k)
			otherItemSegment := otherItem.GetGenericSegment(k)
			val := itemSegment.GetSegmentValue()
			otherVal := otherItemSegment.GetSegmentValue()
			if val != otherVal {
				j = i
				i--
				ithRangeSegmentIndex = -1
				continue top
			}
		}
		upper := rangeItemUpperValue
		otherUpper := otherRangeSegment.GetUpperSegmentValue()
		if otherUpper > upper {
			upper = otherUpper
		}
		joinedItem := item.ToBlock(ithRangeSegmentIndex, rangeSegment.GetSegmentValue(), upper)
		list[i] = joinedItem
		if joinedItem.GetGenericSegment(ithRangeSegmentIndex).IsFullRange() {
			if ithRangeSegmentIndex == 0 {
				list = list[:1]
				list[i] = joinedItem
				return list
			}
			ithRangeSegmentIndex--
		}
		removedCount++
		k := j + 1
		for k < len(list) && list[k] == nil {
			k++
		}
		if k < len(list) {
			list[j] = list[k]
			list[k] = nil
			jthRangeSegmentIndex = -1
		} else {
			list[j] = nil
			j = i
			i--
			jthRangeSegmentIndex = ithRangeSegmentIndex
			ithRangeSegmentIndex = -1
		}
	}
	if removedCount > 0 {
		newSize := len(list) - removedCount
		for k, l := 0, 0; k < newSize; k, l = k+1, l+1 {
			for list[l] == nil {
				l++
			}
			list[k] = list[l].WithoutPrefixLen()
		}
		list = list[:newSize]
	} else {
		for n := 0; n < len(list); n++ {
			list[n] = list[n].WithoutPrefixLen()
		}
	}
	return list
}

func organizeSequentially[S any, T mergeableType[S, T]](sections []T) (singleElement bool, list []T) {
	var sequentialList []T
	length := len(sections)
	for i := 0; i < length; i++ {
		section := sections[i]
		if section == nil {
			continue
		}
		if !section.IsSequential() {
			if sequentialList == nil {
				sequentialList = make([]T, 0, length)
				for j := 0; j < i; j++ {
					series := sections[j]
					if series != nil {
						sequentialList = append(sequentialList, series)
					}
				}
			}
			iterator := section.SequentialBlockIterator()
			for iterator.HasNext() {
				sequentialList = append(sequentialList, iterator.Next())
			}
		} else if sequentialList != nil {
			sequentialList = append(sequentialList, section)
		}
	}
	if sequentialList == nil {
		sequentialList = sections
	}
	sequentialLen := len(sequentialList)
	for j := 0; j < sequentialLen; j++ {
		series := sequentialList[j]
		if series.IsSinglePrefixBlock() {
			list = append(list, series)
		} else {
			span := series.SpanWithPrefixBlocks()
			list = append(list, span...)
		}
	}
	if len(list) <= 1 {
		return true, list
	}
	sort.Slice(list, func(i, j int) bool {
		return LowValueComparator.CompareSeries(list[i], list[j]) < 0
	})
	return false, list
}

func organizeSequentialMerge[S any, T mergeableType[S, T]](sections []T) (singleElement bool, list []T) {
	for i := 0; i < len(sections); i++ {
		section := sections[i]
		if section == nil {
			continue
		}
		if section.IsSequential() {
			list = append(list, section)
		} else {
			iterator := section.SequentialBlockIterator()
			for iterator.HasNext() {
				list = append(list, iterator.Next())
			}
		}
	}
	if len(list) == 1 {
		singleElement = true
		return
	}
	sort.Slice(list, func(i, j int) bool {
		return LowValueComparator.CompareSeries(list[i], list[j]) < 0
	})
	return
}
