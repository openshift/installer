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

import (
	"math"
	"math/big"
)

// returns true for overflow
func checkOverflow( // used by IPv4 and MAC
	increment int64,
	lowerValue,
	upperValue,
	counter,
	maxValue func() uint64,
	isSequential func() bool) bool {
	if increment < 0 {
		return lowerValue() < uint64(-increment)
	} else if increment != 0 {
		if isSequential() {
			maxVal := maxValue()
			return uint64(increment) > maxVal-lowerValue()
		} else {
			countMinus1 := counter() - 1
			uIncrement := uint64(increment)
			if uIncrement > countMinus1 {
				if countMinus1 > 0 {
					uIncrement -= countMinus1
				}
				room := maxValue() - upperValue()
				return uIncrement > room
			}
		}
	}
	return false
}

func checkOverflowBig( // used by IPv6
	increment int64,
	bigIncrement *big.Int,
	lowerValue,
	upperValue,
	count,
	maxValue func() *big.Int,
	isSequential func() bool) bool {
	return checkOverflowB(increment < 0, bigIncrement, lowerValue, upperValue, count, maxValue, isSequential)
}

func checkOverflowBigger( // used by IPv6
	bigIncrement *big.Int,
	lowerValue,
	upperValue,
	count,
	maxValue func() *big.Int,
	isSequential func() bool) bool {
	return checkOverflowB(bigIsNegative(bigIncrement), bigIncrement, lowerValue, upperValue, count, maxValue, isSequential)
}

func checkOverflowB(
	incrementIsNegative bool,
	bigIncrement *big.Int,
	lowerValue,
	upperValue,
	counter,
	maxValue func() *big.Int,
	isSequential func() bool) bool {
	if incrementIsNegative {
		return lowerValue().CmpAbs(bigZero().Neg(bigIncrement)) < 0
	} else if !bigIsZero(bigIncrement) {
		if isSequential() {
			maxVal := maxValue()
			return bigIncrement.CmpAbs(maxVal.Sub(maxVal, lowerValue())) > 0
		} else {
			count := counter()
			if bigIncrement.CmpAbs(count) >= 0 {
				count.Sub(count, bigOneConst())
				count.Sub(bigIncrement, count)
				maxVal := maxValue()
				maxVal.Sub(maxVal, upperValue())
				return count.CmpAbs(maxVal) > 0
			}
		}
	}
	return false
}

// Handles the cases in which we can use longs rather than BigInteger
func fastIncrement( // used by IPv6
	section *AddressSection,
	inc int64,
	creator addressSegmentCreator,
	lowerProducer,
	upperProducer func() *AddressSection,
	prefixLength PrefixLen) *AddressSection {
	if inc >= 0 {
		countMinus1 := section.GetCount()
		countMinus1.Sub(countMinus1, bigOneConst())
		uincrement := uint64(inc)
		var maxUint64 big.Int
		maxUint64.SetUint64(math.MaxUint64)
		if countMinus1.CmpAbs(&maxUint64) <= 0 {
			longCountMinus1 := countMinus1.Uint64()
			if longCountMinus1 >= uincrement {
				if longCountMinus1 == uincrement {
					return upperProducer()
				}
				return incrementRange(section, inc, lowerProducer, prefixLength)
			}
			upperValue := section.GetUpperValue()
			if upperValue.CmpAbs(&maxUint64) <= 0 {
				return increment(
					section,
					inc,
					creator,
					func() uint64 { return longCountMinus1 + 1 },
					func() uint64 { return section.GetUpperValue().Uint64() },
					func() uint64 { return upperValue.Uint64() },
					lowerProducer,
					upperProducer,
					prefixLength)
			}
		}
	} else {
		var maxUint64 big.Int
		maxUint64.SetUint64(math.MaxUint64)
		value := section.GetValue()
		if value.CmpAbs(&maxUint64) <= 0 {
			return add(lowerProducer(), value.Uint64(), inc, creator, prefixLength)
		}
	}
	return nil
}

// this does not handle overflow, overflow should be checked before calling this
func increment( // used by IPv4 and MAC, but also IPv6 addresses with prefix ::/64
	section *AddressSection,
	increment int64,
	creator addressSegmentCreator,
	counter,
	lowerValue,
	upperValue func() uint64,
	lowerProducer,
	upperProducer func() *AddressSection,
	prefixLength PrefixLen) *AddressSection {
	if !section.isMultiple() {
		return add(section, lowerValue(), increment, creator, prefixLength)
	}
	isDecrement := increment <= 0
	if isDecrement {
		//we know lowerValue + increment >= 0 because we already did an overflow check
		return add(lowerProducer(), lowerValue(), increment, creator, prefixLength)
	}
	uIncrement := uint64(increment)
	countMinus1 := counter() - 1
	if countMinus1 >= uIncrement {
		if countMinus1 == uIncrement {
			return upperProducer()
		}
		return incrementRange(section, increment, lowerProducer, prefixLength)
	}
	upperVal := upperValue()
	if uIncrement <= math.MaxUint64-upperVal {
		return add(upperProducer(), upperVal, int64(uIncrement-countMinus1), creator, prefixLength)
	}
	return addBig(upperProducer(), bigZero().SetUint64(uIncrement-countMinus1), creator, prefixLength)
}

// this does not handle overflow, overflow should be checked before calling this
func incrementBig( // used by MAC and IPv6
	section *AddressSection,
	increment int64,
	bigIncrement *big.Int,
	creator addressSegmentCreator,
	lowerProducer,
	upperProducer func() *AddressSection,
	prefixLength PrefixLen) *AddressSection {
	if !section.isMultiple() {
		return addBig(section, bigIncrement, creator, prefixLength)
	}
	isDecrement := increment <= 0
	if isDecrement {
		return addBig(lowerProducer(), bigIncrement, creator, prefixLength)
	}
	count := section.GetCount()
	incrementPlus1 := bigZero().Add(bigIncrement, bigOneConst())
	countCompare := count.CmpAbs(incrementPlus1)
	if countCompare <= 0 {
		if countCompare == 0 {
			return upperProducer()
		}
		return addBig(upperProducer(), incrementPlus1.Sub(incrementPlus1, count), creator, prefixLength) //
	}
	return incrementRange(section, increment, lowerProducer, prefixLength)
}

// this does not handle overflow, overflow should be checked before calling this
func incrementBigger( // used by MAC and IPv6
	section *AddressSection,
	bigIncrement *big.Int,
	creator addressSegmentCreator,
	lowerProducer,
	upperProducer func() *AddressSection,
	prefixLength PrefixLen) *AddressSection {
	if !section.isMultiple() {
		return addBig(section, bigIncrement, creator, prefixLength)
	} else if bigIsNonPositive(bigIncrement) {
		return addBig(lowerProducer(), bigIncrement, creator, prefixLength)
	}
	count := section.GetCount()
	incrementPlus1 := bigZero().Add(bigIncrement, bigOneConst())
	countCompare := count.CmpAbs(incrementPlus1)
	if countCompare <= 0 {
		if countCompare == 0 {
			return upperProducer()
		}
		return addBig(upperProducer(), incrementPlus1.Sub(incrementPlus1, count), creator, prefixLength)
	}
	return incrementRangeBig(section, bigIncrement, lowerProducer, prefixLength)
}

// rangeIncrement the positive value of the number of increments through the range (0 means take lower or upper value in range)
func incrementRange(
	section *AddressSection,
	increment int64,
	lowerProducer func() *AddressSection,
	prefixLength PrefixLen) *AddressSection {
	if increment == 0 {
		return lowerProducer()
	}
	segCount := section.GetSegmentCount()

	var newSegments []*AddressDivision
	if segCount > 0 {
		newSegments = make([]*AddressDivision, segCount)
		i := segCount - 1
		seg := section.GetSegment(i)
		bitCount := seg.getBitCount()

		for {
			segRange := seg.GetValueCount()
			segRange64 := int64(segRange)
			var revolutions int64
			var remainder SegInt
			var newSegment *AddressDivision
			segPrefixLength := getSegmentPrefixLength(section.GetBitsPerSegment(), prefixLength, i)
			if bitCount == IPv6BitsPerSegment && segRange == IPv6MaxValuePerSegment+1 {
				revolutions = increment >> IPv6BitsPerSegment
				remainder = SegInt(increment & IPv6MaxValuePerSegment)
				newSegment = createAddressDivision(seg.deriveNewMultiSeg(remainder, remainder, segPrefixLength))
			} else if bitCount == IPv4BitsPerSegment && segRange == IPv4MaxValuePerSegment+1 {
				revolutions = increment >> IPv4BitsPerSegment
				remainder = SegInt(increment & IPv4MaxValuePerSegment)
				newSegment = createAddressDivision(seg.deriveNewMultiSeg(remainder, remainder, segPrefixLength))
			} else if segRange == 1 {
				revolutions = increment
				val := seg.getSegmentValue()
				newSegment = createAddressDivision(seg.deriveNewMultiSeg(val, val, segPrefixLength))
			} else {
				revolutions = increment / segRange64
				remainder = SegInt(increment % segRange64)
				val := seg.getSegmentValue() + SegInt(remainder)
				newSegment = createAddressDivision(seg.deriveNewMultiSeg(val, val, segPrefixLength))
			}
			newSegments[i] = newSegment
			if revolutions == 0 {
				for i--; i >= 0; i-- {
					original := section.GetSegment(i)
					val := original.getSegmentValue()
					segPrefixLength = getSegmentPrefixLength(section.GetBitsPerSegment(), prefixLength, i)
					newSegment = createAddressDivision(seg.deriveNewMultiSeg(val, val, segPrefixLength))
					newSegments[i] = newSegment
				}
				break
			}
			if i--; i < 0 {
				break
			}
			increment = revolutions
			seg = section.GetSegment(i)
		}
	} else {
		newSegments = section.getDivisionsInternal()
	}
	return createSection(newSegments, prefixLength, section.getAddrType())
}

func incrementRangeBig(
	section *AddressSection,
	increment *big.Int,
	lowerProducer func() *AddressSection,
	prefixLength PrefixLen) *AddressSection {
	if bigIsZero(increment) {
		return lowerProducer()
	}
	segCount := section.GetSegmentCount()
	newSegments := make([]*AddressDivision, segCount)
	for i := segCount - 1; i >= 0; i-- {
		seg := section.GetSegment(i)
		segRange := seg.GetValueCount()
		var revolutions, remainder big.Int
		revolutions.DivMod(increment, bigZero().SetUint64(segRange), &remainder)
		val := seg.getSegmentValue() + SegInt(remainder.Uint64())
		segPrefixLength := getSegmentPrefixLength(section.GetBitsPerSegment(), prefixLength, i)
		newSegment := createAddressDivision(seg.deriveNewMultiSeg(val, val, segPrefixLength))
		newSegments[i] = newSegment
		if bigIsZero(&revolutions) {
			for i--; i >= 0; i-- {
				original := section.GetSegment(i)
				val = original.getSegmentValue()
				segPrefixLength = getSegmentPrefixLength(section.GetBitsPerSegment(), prefixLength, i)
				newSegment = createAddressDivision(seg.deriveNewMultiSeg(val, val, segPrefixLength))
				newSegments[i] = newSegment
			}
			break
		} else {
			increment = &revolutions
		}
	}
	return createSection(newSegments, prefixLength, section.getAddrType())
}

// this does not handle overflow, overflow should be checked before calling this
func addBig(section *AddressSection, increment *big.Int, creator addressSegmentCreator, prefixLength PrefixLen) *AddressSection {
	segCount := section.GetSegmentCount()
	fullValue := section.GetValue()
	fullValue.Add(fullValue, increment)
	expectedByteCount := section.GetByteCount()
	bytes := fullValue.Bytes() // could use FillBytes but that only came with 1.15
	segments, _ := toSegments(
		bytes,
		segCount,
		section.GetBytesPerSegment(),
		section.GetBitsPerSegment(),
		creator,
		prefixLength)
	res := createSection(segments, prefixLength, section.getAddrType())
	if expectedByteCount == len(bytes) && res.cache != nil {
		res.cache.bytesCache = &bytesCache{
			lowerBytes: bytes,
			upperBytes: bytes,
		}
	}
	return res
}

func add(section *AddressSection, fullValue uint64, increment int64, creator addressSegmentCreator, prefixLength PrefixLen) *AddressSection {
	var highBytes, lowBytes uint64
	if increment < 0 {
		lowBytes = fullValue - uint64(-increment)
	} else {
		space := math.MaxUint64 - fullValue
		uIncrement := uint64(increment)
		if uIncrement > space {
			// val := inc + fullValue
			// val = (math.MaxUint64 + 1)  + (fullValue - (MaxUint64 - inc + 1))
			// val = (highBytes of 1) + (loweBytes of fullValue - (MaxUint64 - inc + 1)))
			highBytes = 1
			lowBytes = fullValue - (math.MaxUint64 - (uIncrement - 1))
		} else {
			lowBytes = fullValue + uint64(increment)
		}
	}
	segCount := section.GetSegmentCount()
	newSegs := createSegmentsUint64(
		segCount,
		highBytes,
		lowBytes,
		section.GetBytesPerSegment(),
		section.GetBitsPerSegment(),
		creator,
		prefixLength)
	return createSection(newSegs, prefixLength, section.getAddrType())
}

// When this is called, we know the series are comparable and not nil.
// Callers must ensure that the section values can fit into an int64.
func enumerateSmall(section, otherSection *AddressSection, low64, low64Upper func(*AddressSection) uint64) (val int64, ok bool) {
	if otherSection.isMultiple() {
		return
	} else if section == otherSection {
		ok = true // val is already zero, no need to set it
		return
	}
	return enumerateSmallImpl(section, otherSection, low64, low64Upper)
}

// low64 returns the lower value's lower 64 bits as a uint64
// low64Upper returns the upper value's lower 64 bits as a uint64
// When this is called, we know the series are comparable and not nil.
// Callers must ensure that the section values can fit into an int64.
func enumerateSmallImpl(section, otherSection *AddressSection, low64, low64Upper func(*AddressSection) uint64) (val int64, exists bool) {
	if section.isMultiple() {
		if !section.IsSequential() {
			if compareSegmentValues(true, section, otherSection) < 0 {
				result := (int64(low64(otherSection)) - int64(low64Upper(section))) + int64(section.getCachedCount().Uint64()-1)
				return result, true
			} else if compareSegmentValues(false, section, otherSection) <= 0 {
				var total uint64
				var cumulativeSize uint64 = 1
				for i := section.GetSegmentCount() - 1; ; i-- {
					segment, otherSegment := section.GetSegment(i), otherSection.GetSegment(i)
					otherValue := otherSegment.GetSegmentValue()
					segValue := segment.GetSegmentValue()
					if otherValue < segValue || otherValue > segment.getUpperSegmentValue() {
						return
					}
					total += cumulativeSize * uint64(otherValue-segValue)
					if i == 0 {
						return int64(total), true
					}
					cumulativeSize *= segment.GetValueCount()
				}
			}
		}
	}
	return int64(low64(otherSection)) - int64(low64(section)), true
}

func enumerateBig(section, otherSection *AddressSection, low64, low64Upper func(*AddressSection) uint64) *big.Int {
	if otherSection.isMultiple() {
		return nil
	} else if otherSection == section { // both the same individual address
		return bigZero()
	}
	// If the initial segments beyond 64 bits match, which is probably the case for most subnets,
	// then we can just use long values to calculate
	initialSegsMatch := true
	bitsPerSegment := section.GetBitsPerSegment()
	segmentCount := section.GetSegmentCount()
	totalBits := getTotalBits(segmentCount, section.GetBytesPerSegment(), bitsPerSegment)
	i := 0
	for totalBits > 64 {
		seg, otherSeg := section.GetSegment(i), otherSection.GetSegment(i)
		if !seg.Matches(otherSeg.GetSegmentValue()) {
			initialSegsMatch = false
			break
		}
		totalBits -= bitsPerSegment
		i++
	}
	if initialSegsMatch {
		if totalBits == 64 {
			seg, otherSeg := section.GetSegment(i), otherSection.GetSegment(i)
			var mask SegInt = (1 << bitsPerSegment) >> 1
			if seg.MatchesWithMask(otherSeg.getSegmentValue()&mask, mask) {
				// we can use uint64 to calculate
				if result, ok := enumerateSmallImpl(section, otherSection, low64, low64Upper); ok {
					return big.NewInt(result)
				}
				return nil
			}
		} else {
			// we can use uint64 to calculate
			if result, ok := enumerateSmallImpl(section, otherSection, low64, low64Upper); ok {
				return big.NewInt(result)
			}
			return nil
		}
	}
	// we use big ints to calculate
	if section.isMultiple() {
		if !section.IsSequential() {
			if section.getUpper().Compare(otherSection) < 0 {
				val := otherSection.GetValue()
				val.Sub(val, section.GetUpperValue())
				val.Add(val, section.getCachedCount())
				return val.Sub(val, bigOneConst())
			} else if section.getLower().Compare(otherSection) <= 0 {
				total := bigZero()
				cumulativeSize := bigOne()
				for j := section.GetSegmentCount() - 1; ; j-- {
					segment, otherSegment := section.GetSegment(j), otherSection.GetSegment(j)
					otherValue := otherSegment.GetSegmentValue()
					segValue := segment.GetSegmentValue()
					if otherValue < segValue || otherValue > segment.getUpperSegmentValue() {
						return nil
					}
					diff := bigZero().SetUint64(uint64(otherValue - segValue))
					diff.Mul(diff, cumulativeSize)
					total.Add(total, diff)
					if j == 0 {
						return total
					}
					cumulativeSize.Mul(cumulativeSize, segment.getCount())
				}
			}
		}
	}
	val := otherSection.GetValue()
	return val.Sub(val, section.GetValue())
}
