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

import "math/bits"

func getCoveringPrefixBlock[S any, T spannableType[S, T]](first, other T) T {
	result := checkPrefixBlockContainment(first, other)
	if result != nil {
		return result
	}
	return applyOperatorToLowerUpper(first, other, false, coverWithPrefixBlockWrapped[S, T])[0]
}

func coverWithPrefixBlockWrapped[S any, T spannableType[S, T]](lower, upper T) []T {
	return []T{coverWithPrefixBlock(lower, upper)}
}

func coverWithPrefixBlock[S any, T spannableType[S, T]](lower, upper T) T {
	segCount := lower.GetSegmentCount()
	bitsPerSegment := lower.GetBitsPerSegment()
	var currentSegment int
	var previousSegmentBits BitCount
	for ; currentSegment < segCount; currentSegment++ {
		lowerSeg := lower.GetGenericSegment(currentSegment)
		upperSeg := upper.GetGenericSegment(currentSegment)
		var lowerValue, upperValue SegInt
		lowerValue = lowerSeg.GetSegmentValue() //these are single addresses, so lower or upper value no different here
		upperValue = upperSeg.GetSegmentValue()
		differing := lowerValue ^ upperValue
		if differing != 0 {
			highestDifferingBitInRange := BitCount(bits.LeadingZeros32(differing)) - (SegIntSize - bitsPerSegment)
			differingBitPrefixLen := highestDifferingBitInRange + previousSegmentBits
			return lower.ToPrefixBlockLen(differingBitPrefixLen)
		}
		previousSegmentBits += bitsPerSegment
	}
	//all bits match, it's just a single address
	return lower.ToPrefixBlockLen(lower.GetBitCount())
}
