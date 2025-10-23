//
// Copyright 2020-2023 Sean C Foley
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

import "math/big"

var (
	// CountComparator compares by count first, then by value.
	CountComparator = AddressComparator{countComparator{}}

	// HighValueComparator compares by high value first, then low.
	HighValueComparator = AddressComparator{valueComparator{compareHighValue: true}}

	// LowValueComparator compares by low value first, then high.
	LowValueComparator = AddressComparator{valueComparator{}}

	// With the reverse comparators, ordering with the secondary values (higher or lower) follow a reverse ordering than the primary values (lower or higher)

	// ReverseHighValueComparator is like HighValueComparator but when comparing the low value, reverses the comparison.
	ReverseHighValueComparator = AddressComparator{valueComparator{compareHighValue: true, flipSecond: true}}

	// ReverseLowValueComparator is like LowValueComparator but when comparing the high value, reverses the comparison.
	ReverseLowValueComparator = AddressComparator{valueComparator{flipSecond: true}}
)

type componentComparator interface {
	compareSectionParts(one, two *AddressSection) int

	compareParts(one, two AddressDivisionSeries) int

	compareSegValues(oneUpper, oneLower, twoUpper, twoLower SegInt) int

	compareValues(oneUpper, oneLower, twoUpper, twoLower uint64) int

	compareLargeValues(oneUpper, oneLower, twoUpper, twoLower *big.Int) int
}

type groupingType int

const (
	ipv6sectype          groupingType = 7
	ipv4sectype          groupingType = 6
	ipsectype            groupingType = 5
	macsectype           groupingType = 4
	sectype              groupingType = 3
	ipv6v4groupingtype   groupingType = 2
	largegroupingtype    groupingType = -2
	standardgroupingtype groupingType = -3
	adaptivezerotype     groupingType = -4
	unknowntype          groupingType = -5
)

type divType int

const (
	ipv6segtype     divType = 6
	ipv4segtype     divType = 5
	ipsegtype       divType = 4
	macsegtype      divType = 3
	segtype         divType = 1
	standarddivtype divType = 0
	largedivtype    divType = -2
	unknowndivtype  divType = -3
)

type rangeType int

const (
	ipv6rangetype    rangeType = 2
	ipv4rangetype    rangeType = 1
	iprangetype      rangeType = 0
	unknownrangetype rangeType = -1
)

func checkSegmentType(div AddressSegmentType) (isNil bool, divType divType) {
	if isNil = div == nil; isNil {
		divType = unknowndivtype
	} else {
		seg := div.ToSegmentBase()
		if isNil = seg == nil; !isNil {
			if seg.IsIPv6() {
				divType = ipv6segtype
			} else if seg.IsIPv4() {
				divType = ipv4segtype
			} else if seg.IsMAC() {
				divType = macsegtype
			} else if seg.IsIP() {
				divType = ipsegtype
			} else {
				divType = segtype
			}
		} else {
			divType = unknowndivtype
		}
	}
	return
}

func checkDivisionType(genericDiv DivisionType) (isNil, isStandard bool, divType divType, standardDiv StandardDivisionType) {
	// Note: https://go.dev/play/p/4cHtDqDxpjp shows the behaviour or type-checking
	if standardDiv, isStandard = genericDiv.(StandardDivisionType); isStandard {
		div := standardDiv.ToDiv()
		if isNil = div == nil; !isNil {
			if div.IsIPv6() {
				divType = ipv6segtype
			} else if div.IsIPv4() {
				divType = ipv4segtype
			} else if div.IsMAC() {
				divType = macsegtype
			} else if div.IsIP() {
				divType = ipsegtype
			} else if div.IsSegmentBase() {
				divType = segtype
			} else {
				divType = standarddivtype
			}
		} else {
			divType = unknowndivtype
		}
	} else if largeDiv, isLarge := genericDiv.(*IPAddressLargeDivision); isLarge {
		if isNil = largeDiv.isNil(); !isNil {
			divType = largedivtype
		}
	} else {
		isNil = genericDiv == nil
		// it could still have some external type, so not a nil interface but a nil value with that type, but we have no way to know
		divType = unknowndivtype
	}
	return
}

func checkSectionType(sect AddressSectionType) (isNil bool, groupingType groupingType) {
	if isNil = sect == nil; isNil {
		groupingType = unknowntype
	} else {
		section := sect.ToSectionBase()
		if isNil = section == nil; !isNil {
			groupingType = checkSectionGroupingType(section)
		} else {
			groupingType = unknowntype
		}
	}
	return
}

func checkSectionGroupingType(section *AddressSection) (groupingType groupingType) {
	if section.IsIPv4() {
		groupingType = ipv4sectype
	} else if section.IsIPv6() {
		groupingType = ipv6sectype
	} else if section.IsIP() {
		groupingType = ipsectype
	} else if section.IsMAC() {
		groupingType = macsectype
	} else if section.IsAdaptiveZero() {
		// The zero grouping can represent a zero-length section of any address type.
		// This is necessary because sections and groupings have no init() method to ensure zero-sections are always assigned an address type.
		// We would need the zero grouping to be less than everything else or more than everything else for comparison consistency.
		// Empty sections or groupings that have an address type are not considered equal.  They can represent only one address type.
		// This is similar to the fact that a MAC section and an IPv4 section can be structurally identical but not equal due to the type.
		//
		// See IsAdaptiveZero() method for more details.
		groupingType = adaptivezerotype
	} else {
		groupingType = sectype
	}
	return
}

func checkGroupingType(series AddressDivisionSeries) (
	isNil bool, groupingType groupingType) {
	// Note: https://go.dev/play/p/4cHtDqDxpjp shows the behaviour or type-checking
	if sgrouping, isStandard := series.(StandardDivGroupingType); isStandard {
		group := sgrouping.ToDivGrouping()
		if isNil = group == nil; !isNil {
			if group.IsAdaptiveZero() {
				// The zero grouping can represent a zero-length section of any address type.
				// This is necessary because sections and groupings have no init() method to ensure zero-sections are always assigned an address type.
				// We would need the zero grouping to be less than everything else or more than everything else for comparison consistency.
				// Empty sections or groupings that have an address type are not considered equal.  They can represent only one address type.
				// This is similar to the fact that a MAC section and an IPv4 section can be structurally identical but not equal due to the type.
				//
				// See IsAdaptiveZero() method for more details.
				groupingType = adaptivezerotype
			} else if group.IsIPv6() {
				groupingType = ipv6sectype
			} else if group.IsMixedIPv6v4() {
				groupingType = ipv6v4groupingtype
			} else if group.IsIPv4() {
				groupingType = ipv4sectype
			} else if group.IsMAC() {
				groupingType = macsectype
			} else if group.IsIP() {
				// Currently the ipsectype result is impossible, because a zero IP section has type adaptivezerotype,
				// while a non-zero IP section can only be ipv6sectype or ipv4sectype
				groupingType = ipsectype
			} else if group.isAddressSection() {
				// Currently the sectype result is impossible, because a zero section has type adaptivezerotype,
				// while a non-zero IP section can only be ipv6sectype or ipv4sectype
				groupingType = sectype
			} else {
				groupingType = standardgroupingtype
			}
		} else {
			groupingType = unknowntype
		}
	} else if lgrouping, isLarge := series.(*IPAddressLargeDivisionGrouping); isLarge {
		if isNil = lgrouping.isNil(); !isNil {
			groupingType = largegroupingtype
		} else {
			groupingType = unknowntype
		}
	} else {
		isNil = series == nil
		// it could still have some external type, so not a nil interface but a nil value with that type, but we have no way to know
		groupingType = unknowntype
	}
	return
}

func checkRangeType(r IPAddressSeqRangeType) (isNil bool, rngType rangeType, rng *SequentialRange[*IPAddress]) {
	if isNil = r == nil; isNil {
		rngType = unknownrangetype
	} else {
		rng = r.ToIP()
		if isNil = rng == nil; !isNil {
			version := r.GetIPVersion()
			if version.IsIPv4() {
				rngType = ipv4rangetype
			} else if version.IsIPv6() {
				rngType = ipv6rangetype
			} else {
				rngType = iprangetype
			}
		} else {
			rngType = unknownrangetype
		}
	}
	return
}

// compareAddressLowerValues compares any two individual addresses (including different versions or address types)
// It returns a negative integer, zero, or a positive integer if address item one is less than, equal, or greater than address item two.
// It ignores IPv6 zone and panics on nil AddressType or nil addresses.
// If an address is not a single address, then the lower values are used.
func compareAddressLowerValues(one, two AddressType) int {
	oneAddr := one.ToAddressBase()
	twoAddr := two.ToAddressBase()
	return compareSegmentValues(false, oneAddr.GetSection(), twoAddr.GetSection())
}

// called only when it is known the sections have same segment count and segment bit size
func compareSegmentValues(compareUpper bool, one, two *AddressSection) (result int) {
	segCount := one.GetSegmentCount()
	for i := 0; i < segCount; i++ {
		segOne := one.GetSegment(i)
		segTwo := two.GetSegment(i)
		var s1, s2 SegInt
		if compareUpper {
			s1 = segOne.GetUpperSegmentValue()
			s2 = segTwo.GetUpperSegmentValue()
		} else {
			s1 = segOne.GetSegmentValue()
			s2 = segTwo.GetSegmentValue()
		}
		if s1 != s2 {
			if s1 > s2 {
				result = 1
			} else {
				result = -1
			}
			return
		}
	}
	return
}

// AddressComparator has methods to compare addresses, or sections, or division series, or segments, or divisions, or sequential ranges.
// AddressComparator also allows you to compare any two instances of any such address items, using the Compare method.
// The zero value acts like CountComparator, the default comparator.
type AddressComparator struct {
	componentComparator componentComparator
}

// CompareAddresses compares any two addresses (including different versions or address types)
// It returns a negative integer, zero, or a positive integer if address item one is less than, equal, or greater than address item two.
func (comp AddressComparator) CompareAddresses(one, two AddressType) int {
	if one == nil || one.ToAddressBase() == nil {
		if two == nil || two.ToAddressBase() == nil {
			return 0
		}
		return -1
	} else if two == nil || two.ToAddressBase() == nil {
		return 1
	}
	oneAddr := one.ToAddressBase()
	twoAddr := two.ToAddressBase()
	result := comp.CompareAddressSections(oneAddr.GetSection(), twoAddr.GetSection())
	if result == 0 {
		if oneIPv6 := oneAddr.ToIPv6(); oneIPv6 != nil {
			twoIPv6 := twoAddr.ToIPv6()
			oneZone := oneIPv6.zone
			twoZone := twoIPv6.zone
			if oneZone == twoZone {
				return 0
			} else if oneZone < twoZone {
				return -1
			}
			return 1
		}
	}
	return result
}

// CompareAddressSections compares any two address sections (including from different versions or address types).
// It returns a negative integer, zero, or a positive integer if address item one is less than, equal, or greater than address item two.
func (comp AddressComparator) CompareAddressSections(one, two AddressSectionType) int {
	oneIsNil, oneGroupingType := checkSectionType(one)
	twoIsNil, twoGroupingType := checkSectionType(two)
	if oneIsNil {
		if twoIsNil {
			return 0
		}
		return -1
	} else if twoIsNil {
		return 1
	} else if result := oneGroupingType - twoGroupingType; result != 0 {
		return int(result)
	} else if result := int(one.GetBitCount() - two.GetBitCount()); result != 0 {
		return result
	}
	return comp.getCompComp().compareSectionParts(one.ToSectionBase(), two.ToSectionBase())
}

func (comp AddressComparator) getCompComp() componentComparator {
	compComp := comp.componentComparator
	if compComp == nil {
		return countComparator{}
	}
	return compComp
}

func unwrapWrapper(item AddressDivisionSeries) AddressDivisionSeries {
	if wrapper, ok := item.(ExtendedIPSegmentSeries); ok {
		return wrapper.Unwrap()
	}
	return item
}

// CompareSeries compares any two address division series (including from different versions or address types).
// It returns a negative integer, zero, or a positive integer if address item one is less than, equal, or greater than address item two.
func (comp AddressComparator) CompareSeries(one, two AddressDivisionSeries) int {
	one = unwrapWrapper(one)
	two = unwrapWrapper(two)
	if addrSeries1, ok := one.(AddressType); ok {
		if addrSeries2, ok := two.(AddressType); ok {
			return comp.CompareAddresses(addrSeries1, addrSeries2)
		}
		return 1
	} else if _, ok := two.(AddressType); ok {
		return -1
	}
	// at this point they must be both groupings if not nil
	if addrSection1, ok := one.(AddressSectionType); ok {
		if addrSection2, ok := two.(AddressSectionType); ok {
			return comp.CompareAddressSections(addrSection1, addrSection2)
		}
	}
	oneIsNil, oneGroupingType := checkGroupingType(one)
	twoIsNil, twoGroupingType := checkGroupingType(two)

	// All nils are equivalent.  We decided that a nil interface should be equivalent to an interface with a nil value (standard or large)
	// But if nil interface == nil standard, and nil interface == nil large, then nil standard == nil large, by transitive condition.
	// And in fact, if you attempted to categorize the 3 nil types, it would get quite confusing perhaps.
	if oneIsNil {
		if twoIsNil {
			return 0
		}
		return -1
	} else if twoIsNil {
		return 1
	} else if result := oneGroupingType - twoGroupingType; result != 0 {
		return int(result)
	} else if result := int(one.GetBitCount() - two.GetBitCount()); result != 0 {
		return result
	}
	return comp.getCompComp().compareParts(one, two)
}

// CompareSegments compares any two address segments (including from different versions or address types).
// It returns a negative integer, zero, or a positive integer if address item one is less than, equal, or greater than address item two.
func (comp AddressComparator) CompareSegments(one, two AddressSegmentType) int {
	oneIsNil, oneDivType := checkSegmentType(one)
	twoIsNil, twoDivType := checkSegmentType(two)
	// All nils are equivalent.  We decided that a nil interface should be equivalent to an interface with a nil value (standard or large)
	// But if nil interface == nil standard, and nil interface == nil large, then nil standard == nil large, by transitive condition.
	// And in fact, if you attempted to categorize the 3 nil types, it would get quite confusing perhaps.
	if oneIsNil {
		if twoIsNil {
			return 0
		}
		return -1
	} else if twoIsNil {
		return 1
	} else if result := oneDivType - twoDivType; result != 0 {
		return int(result)
	} else if result := int(one.GetBitCount() - two.GetBitCount()); result != 0 {
		return result
	}
	oneSeg := one.ToSegmentBase()
	twoSeg := two.ToSegmentBase()
	return comp.getCompComp().compareSegValues(oneSeg.GetUpperSegmentValue(), oneSeg.GetSegmentValue(),
		twoSeg.GetUpperSegmentValue(), twoSeg.GetSegmentValue())
}

// CompareDivisions compares any two address divisions (including from different versions or address types).
// It returns a negative integer, zero, or a positive integer if address item one is less than, equal, or greater than address item two.
func (comp AddressComparator) CompareDivisions(one, two DivisionType) int {
	if addrSeg1, ok := one.(AddressSegmentType); ok {
		if addrSeg2, ok := two.(AddressSegmentType); ok {
			return comp.CompareSegments(addrSeg1, addrSeg2)
		}
	}
	oneIsNil, oneIsStandard, oneDivType, oneStandardDiv := checkDivisionType(one)
	twoIsNil, twoIsStandard, twoDivType, twoStandardDiv := checkDivisionType(two)
	// All nils are equivalent.  We decided that a nil interface should be equivalent to an interface with a nil value (standard or large)
	// But if nil interface == nil standard, and nil interface == nil large, then nil standard == nil large, by transitive condition.
	// And in fact, if you attempted to categorize the 3 nil types, it would get quite confusing perhaps.
	if oneIsNil {
		if twoIsNil {
			return 0
		}
		return -1
	} else if twoIsNil {
		return 1
	} else if result := oneDivType - twoDivType; result != 0 {
		return int(result)
	} else if result := int(one.GetBitCount() - two.GetBitCount()); result != 0 {
		return result
	}
	compComp := comp.getCompComp()
	if oneIsStandard {
		if twoIsStandard {
			div1, div2 := oneStandardDiv.ToDiv(), twoStandardDiv.ToDiv()
			return compComp.compareValues(div1.GetUpperDivisionValue(), div1.GetDivisionValue(), div2.GetUpperDivisionValue(), div2.GetDivisionValue())
		}
	}
	return compComp.compareLargeValues(one.GetUpperValue(), one.GetValue(), two.GetUpperValue(), two.GetValue())
}

// CompareRanges compares any two IP address sequential ranges (including from different IP versions).
// It returns a negative integer, zero, or a positive integer if address item one is less than, equal, or greater than address item two.
func (comp AddressComparator) CompareRanges(one, two IPAddressSeqRangeType) int {
	oneIsNil, r1Type, r1 := checkRangeType(one)
	twoIsNil, r2Type, r2 := checkRangeType(two)
	if oneIsNil {
		if twoIsNil {
			return 0
		}
		return -1
	} else if twoIsNil {
		return 1
	}
	result := r1Type - r2Type
	if result != 0 {
		return int(result)
	}
	compComp := comp.getCompComp()
	if r1Type == ipv4rangetype { // avoid using the large values
		r1ipv4 := r1.ToIPv4()
		r2ipv4 := r2.ToIPv4()
		return compComp.compareValues(uint64(r1ipv4.GetUpper().Uint32Value()), uint64(r1ipv4.GetLower().Uint32Value()), uint64(r2ipv4.GetUpper().Uint32Value()), uint64(r2ipv4.GetLower().Uint32Value()))
	}
	return compComp.compareLargeValues(one.GetUpperValue(), one.GetValue(), two.GetUpperValue(), two.GetValue())
}

// Compare returns a negative integer, zero, or a positive integer if address item one is less than, equal, or greater than address item two.
// Any address item is comparable to any other.
func (comp AddressComparator) Compare(one, two AddressItem) int {
	if one == nil {
		if two == nil {
			return 0
		}
		return -1
	} else if two == nil {
		return 1
	}

	if divSeries1, ok := one.(AddressDivisionSeries); ok {
		if divSeries2, ok := two.(AddressDivisionSeries); ok {
			return comp.CompareSeries(divSeries1, divSeries2)
		} else {
			return 1
		}
	} else if div1, ok := one.(DivisionType); ok {
		if div2, ok := two.(DivisionType); ok {
			return comp.CompareDivisions(div1, div2)
		} else {
			return -1
		}
	} else if rng1, ok := one.(IPAddressSeqRangeType); ok {
		if rng2, ok := two.(IPAddressSeqRangeType); ok {
			return comp.CompareRanges(rng1, rng2)
		} else if _, ok := two.(AddressDivisionSeries); ok {
			return -1
		}
		return 1
	}
	// we've covered all known address items for 'one', so check 'two'
	if _, ok := two.(AddressDivisionSeries); ok {
		return -1
	} else if _, ok := two.(DivisionType); ok {
		return 1
		//} else if _, ok := two.(IPAddressSeqRangeType); ok {
		//	return -1
	} else if _, ok := two.(IPAddressSeqRangeType); ok {
		return -1
	}
	// neither are a known AddressItem type
	res := int(one.GetBitCount() - two.GetBitCount())
	if res == 0 {
		return res
	}
	return comp.getCompComp().compareLargeValues(one.GetUpperValue(), one.GetValue(), two.GetUpperValue(), two.GetValue())
}

type valueComparator struct {
	compareHighValue, flipSecond bool
}

func (comp valueComparator) compareSectionParts(one, two *AddressSection) int {
	compareHigh := comp.compareHighValue
	for {
		result := compareSegmentValues(compareHigh, one, two)
		if result != 0 {
			if comp.flipSecond && compareHigh != comp.compareHighValue {
				result = -result
			}
			return result
		}
		compareHigh = !compareHigh
		if compareHigh == comp.compareHighValue {
			break
		}
	}
	return 0
}

func (comp valueComparator) compareParts(oneSeries, twoSeries AddressDivisionSeries) int {
	sizeResult := int(oneSeries.GetBitCount() - twoSeries.GetBitCount())
	if sizeResult != 0 {
		return sizeResult
	}
	result := compareDivBitCounts(oneSeries, twoSeries)
	if result != 0 {
		return result
	}
	compareHigh := comp.compareHighValue
	var one, two *AddressDivisionGrouping
	if o, ok := oneSeries.(StandardDivGroupingType); ok {
		if t, ok := twoSeries.(StandardDivGroupingType); ok {
			one = o.ToDivGrouping()
			two = t.ToDivGrouping()
		}
	}
	oneSeriesByteCount := oneSeries.GetByteCount()
	twoSeriesByteCount := twoSeries.GetByteCount()
	oneBytes := make([]byte, oneSeriesByteCount)
	twoBytes := make([]byte, twoSeriesByteCount)
	for {
		var oneByteCount, twoByteCount, oneByteIndex, twoByteIndex, oneIndex, twoIndex int
		var oneBitCount, twoBitCount, oneTotalBitCount, twoTotalBitCount BitCount
		var oneValue, twoValue uint64
		for oneIndex < oneSeries.GetDivisionCount() || twoIndex < twoSeries.GetDivisionCount() {
			if one != nil {
				if oneBitCount == 0 {
					oneCombo := one.GetDivision(oneIndex)
					oneIndex++
					oneBitCount = oneCombo.GetBitCount()
					if compareHigh {
						oneValue = oneCombo.GetUpperDivisionValue()
					} else {
						oneValue = oneCombo.GetDivisionValue()
					}
				}
				if twoBitCount == 0 {
					twoCombo := two.GetDivision(twoIndex)
					twoIndex++
					twoBitCount = twoCombo.GetBitCount()
					if compareHigh {
						twoValue = twoCombo.GetUpperDivisionValue()
					} else {
						twoValue = twoCombo.GetDivisionValue()
					}
				}
			} else {
				if oneBitCount == 0 {
					if oneByteCount == 0 {
						oneCombo := oneSeries.GetGenericDivision(oneIndex)
						oneIndex++
						if compareHigh {
							oneBytes = oneCombo.CopyUpperBytes(oneBytes)
						} else {
							oneBytes = oneCombo.CopyBytes(oneBytes)
						}
						oneTotalBitCount = oneCombo.GetBitCount()
						oneByteCount = oneCombo.GetByteCount()
						oneByteIndex = 0
					}
					//put some or all of the bytes into a long
					count := 8
					oneValue = 0
					if count < oneByteCount {
						oneBitCount = BitCount(count) << 3
						oneTotalBitCount -= oneBitCount
						oneByteCount -= count
						for count > 0 {
							count--

							oneValue = (oneValue << 8) | uint64(oneBytes[oneByteIndex])
							oneByteIndex++
						}
					} else {
						shortCount := oneByteCount - 1
						lastBitsCount := oneTotalBitCount - (BitCount(shortCount) << 3)
						for shortCount > 0 {
							shortCount--
							oneValue = (oneValue << 8) | uint64(oneBytes[oneByteIndex])
							oneByteIndex++
						}
						oneValue = (oneValue << uint64(lastBitsCount)) | uint64(oneBytes[oneByteIndex]>>uint64(8-lastBitsCount))
						oneByteIndex++
						oneBitCount = oneTotalBitCount
						oneTotalBitCount = 0
						oneByteCount = 0
					}
				}
				if twoBitCount == 0 {
					if twoByteCount == 0 {
						twoCombo := twoSeries.GetGenericDivision(twoIndex)
						twoIndex++
						if compareHigh {
							twoBytes = twoCombo.CopyUpperBytes(twoBytes)
						} else {
							twoBytes = twoCombo.CopyBytes(twoBytes)
						}
						twoTotalBitCount = twoCombo.GetBitCount()
						twoByteCount = twoCombo.GetByteCount()
						twoByteIndex = 0
					}
					//put some or all of the bytes into a long
					count := 8
					twoValue = 0
					if count < twoByteCount {
						twoBitCount = BitCount(count) << 3
						twoTotalBitCount -= twoBitCount
						twoByteCount -= count
						for count > 0 {
							count--

							twoValue = (twoValue << 8) | uint64(twoBytes[twoByteIndex])
							twoByteIndex++
						}
					} else {
						shortCount := twoByteCount - 1
						lastBitsCount := twoTotalBitCount - (BitCount(shortCount) << 3)
						for shortCount > 0 {
							shortCount--

							twoValue = (twoValue << 8) | uint64(twoBytes[twoByteIndex])
							twoByteIndex++
						}
						twoValue = (twoValue << uint(lastBitsCount)) | uint64(twoBytes[twoByteIndex]>>uint(8-lastBitsCount))
						twoByteIndex++
						twoBitCount = twoTotalBitCount
						twoTotalBitCount = 0
						twoByteCount = 0
					}
				}
			}
			oneResultValue := oneValue
			twoResultValue := twoValue
			if twoBitCount == oneBitCount {
				//no adjustment required, compare the values straight up
				oneBitCount = 0
				twoBitCount = 0
			} else {
				diffBits := twoBitCount - oneBitCount
				if diffBits > 0 {
					twoResultValue >>= uint(diffBits)
					twoValue &= ^(^uint64(0) << uint(diffBits))
					twoBitCount = diffBits
					oneBitCount = 0
				} else {
					diffBits = -diffBits
					oneResultValue >>= uint(diffBits)
					oneValue &= ^(^uint64(0) << uint(diffBits))
					oneBitCount = diffBits
					twoBitCount = 0
				}
			}
			if oneResultValue != twoResultValue {
				if comp.flipSecond && compareHigh != comp.compareHighValue {
					if oneResultValue > twoResultValue {
						return -1
					}
					return 1
				}
				if oneResultValue > twoResultValue {
					return 1
				}
				return -1
			}
		}
		compareHigh = !compareHigh
		if compareHigh == comp.compareHighValue {
			break
		}
	}
	return 0
}

func (comp valueComparator) compareSegValues(oneUpper, oneLower, twoUpper, twoLower SegInt) int {
	if comp.compareHighValue {
		if oneUpper == twoUpper {
			if oneLower == twoLower {
				return 0
			} else if oneLower > twoLower {
				if !comp.flipSecond {
					return 1
				}
			}
		} else if oneUpper > twoUpper {
			return 1
		}
	} else {
		if oneLower == twoLower {
			if oneUpper == twoUpper {
				return 0
			} else if oneUpper > twoUpper {
				if !comp.flipSecond {
					return 1
				}
			}
		} else if oneLower > twoLower {
			return 1
		}
	}
	return -1
}

func (comp valueComparator) compareValues(oneUpper, oneLower, twoUpper, twoLower uint64) int {
	if comp.compareHighValue {
		if oneUpper == twoUpper {
			if oneLower == twoLower {
				return 0
			} else if oneLower > twoLower {
				if !comp.flipSecond {
					return 1
				}
			}
		} else if oneUpper > twoUpper {
			return 1
		}
	} else {
		if oneLower == twoLower {
			if oneUpper == twoUpper {
				return 0
			} else if oneUpper > twoUpper {
				if !comp.flipSecond {
					return 1
				}
			}
		} else if oneLower > twoLower {
			return 1
		}
	}
	return -1
}

func (comp valueComparator) compareLargeValues(oneUpper, oneLower, twoUpper, twoLower *big.Int) int {
	var result int
	if comp.compareHighValue {
		result = oneUpper.CmpAbs(twoUpper)
		if result == 0 {
			result = oneLower.CmpAbs(twoLower)
			if comp.flipSecond {
				result = -result
			}
		}
	} else {
		result = oneLower.CmpAbs(twoLower)
		if result == 0 {
			result = oneUpper.CmpAbs(twoUpper)
			if comp.flipSecond {
				result = -result
			}
		}
	}
	return result
}

type countComparator struct{}

func (comp countComparator) compareSectionParts(one, two *AddressSection) int {
	result := compareCount(one, two)
	if result == 0 {
		result = comp.compareEqualSizedSections(one, two)
	}
	return result
}

func (comp countComparator) compareEqualSizedSections(one, two *AddressSection) int {
	segCount := one.GetSegmentCount()
	for i := 0; i < segCount; i++ {
		segOne := one.GetSegment(i)
		segTwo := two.GetSegment(i)
		oneUpper := segOne.GetUpperSegmentValue()
		twoUpper := segTwo.GetUpperSegmentValue()
		oneLower := segOne.GetSegmentValue()
		twoLower := segTwo.GetSegmentValue()
		result := comp.compareSegValues(oneUpper, oneLower, twoUpper, twoLower)
		if result != 0 {
			return result
		}
	}
	return 0
}

func (comp countComparator) compareParts(one, two AddressDivisionSeries) int {
	result := int(one.GetBitCount() - two.GetBitCount())
	if result == 0 {
		result = compareCount(one, two)
		if result == 0 {
			result = comp.compareDivisionGroupings(one, two)
		}
	}
	return result
}

func (comp countComparator) compareDivisionGroupings(oneSeries, twoSeries AddressDivisionSeries) int {
	var one, two *AddressDivisionGrouping
	if o, ok := oneSeries.(StandardDivGroupingType); ok {
		if t, ok := twoSeries.(StandardDivGroupingType); ok {
			one = o.ToDivGrouping()
			two = t.ToDivGrouping()
		}
	}
	result := compareDivBitCounts(oneSeries, twoSeries)
	if result != 0 {
		return result
	}

	oneSeriesByteCount := oneSeries.GetByteCount()
	twoSeriesByteCount := twoSeries.GetByteCount()

	oneUpperBytes := make([]byte, oneSeriesByteCount)
	oneLowerBytes := make([]byte, oneSeriesByteCount)
	twoUpperBytes := make([]byte, twoSeriesByteCount)
	twoLowerBytes := make([]byte, twoSeriesByteCount)

	var oneByteCount, twoByteCount, oneByteIndex, twoByteIndex, oneIndex, twoIndex int
	var oneBitCount, twoBitCount, oneTotalBitCount, twoTotalBitCount BitCount
	var oneUpper, oneLower, twoUpper, twoLower uint64
	for oneIndex < oneSeries.GetDivisionCount() || twoIndex < twoSeries.GetDivisionCount() {
		if one != nil {
			if oneBitCount == 0 {
				oneCombo := one.getDivision(oneIndex)
				oneIndex++
				oneBitCount = oneCombo.GetBitCount()
				oneUpper = oneCombo.GetUpperDivisionValue()
				oneLower = oneCombo.GetDivisionValue()
			}
			if twoBitCount == 0 {
				twoCombo := two.getDivision(twoIndex)
				twoIndex++
				twoBitCount = twoCombo.GetBitCount()
				twoUpper = twoCombo.GetUpperDivisionValue()
				twoLower = twoCombo.GetDivisionValue()
			}
		} else {
			if oneBitCount == 0 {
				if oneByteCount == 0 {
					oneCombo := oneSeries.GetGenericDivision(oneIndex)
					oneIndex++
					oneUpperBytes = oneCombo.CopyUpperBytes(oneUpperBytes)
					oneLowerBytes = oneCombo.CopyBytes(oneLowerBytes)
					oneTotalBitCount = oneCombo.GetBitCount()
					oneByteCount = oneCombo.GetByteCount()
					oneByteIndex = 0
				}
				//put some or all of the bytes into a uint64
				count := 8
				oneUpper = 0
				oneLower = 0
				if count < oneByteCount {
					oneBitCount = BitCount(count << 3)
					oneTotalBitCount -= oneBitCount
					oneByteCount -= count
					for count > 0 {
						count--
						upperByte := oneUpperBytes[oneByteIndex]
						lowerByte := oneLowerBytes[oneByteIndex]
						oneByteIndex++
						oneUpper = (oneUpper << 8) | uint64(upperByte)
						oneLower = (oneLower << 8) | uint64(lowerByte)
					}
				} else {
					shortCount := oneByteCount - 1
					lastBitsCount := oneTotalBitCount - (BitCount(shortCount) << 3)
					for shortCount > 0 {
						shortCount--
						upperByte := oneUpperBytes[oneByteIndex]
						lowerByte := oneLowerBytes[oneByteIndex]
						oneByteIndex++
						oneUpper = (oneUpper << 8) | uint64(upperByte)
						oneLower = (oneLower << 8) | uint64(lowerByte)
					}
					upperByte := oneUpperBytes[oneByteIndex]
					lowerByte := oneLowerBytes[oneByteIndex]
					oneByteIndex++
					oneUpper = (oneUpper << uint(lastBitsCount)) | uint64(upperByte>>uint(8-lastBitsCount))
					oneLower = (oneLower << uint(lastBitsCount)) | uint64(lowerByte>>uint(8-lastBitsCount))
					oneBitCount = oneTotalBitCount
					oneTotalBitCount = 0
					oneByteCount = 0
				}
			}
			if twoBitCount == 0 {
				if twoByteCount == 0 {
					twoCombo := twoSeries.GetGenericDivision(twoIndex)
					twoIndex++
					twoUpperBytes = twoCombo.CopyUpperBytes(twoUpperBytes)
					twoLowerBytes = twoCombo.CopyBytes(twoLowerBytes)
					twoTotalBitCount = twoCombo.GetBitCount()
					twoByteCount = twoCombo.GetByteCount()
					twoByteIndex = 0
				}
				//put some or all of the bytes into a long
				count := 8
				twoUpper = 0
				twoLower = 0
				if count < twoByteCount {
					twoBitCount = BitCount(count << 3)
					twoTotalBitCount -= twoBitCount
					twoByteCount -= count
					for count > 0 {
						count--
						upperByte := twoUpperBytes[twoByteIndex]
						lowerByte := twoLowerBytes[twoByteIndex]
						twoByteIndex++
						twoUpper = (twoUpper << 8) | uint64(upperByte)
						twoLower = (twoLower << 8) | uint64(lowerByte)
					}
				} else {
					shortCount := twoByteCount - 1
					lastBitsCount := twoTotalBitCount - (BitCount(shortCount) << 3)
					for shortCount > 0 {
						shortCount--
						upperByte := twoUpperBytes[twoByteIndex]
						lowerByte := twoLowerBytes[twoByteIndex]
						twoByteIndex++
						twoUpper = (twoUpper << 8) | uint64(upperByte)
						twoLower = (twoLower << 8) | uint64(lowerByte)
					}
					upperByte := twoUpperBytes[twoByteIndex]
					lowerByte := twoLowerBytes[twoByteIndex]
					twoByteIndex++
					twoUpper = (twoUpper << uint(lastBitsCount)) | uint64(upperByte>>uint(8-lastBitsCount))
					twoLower = (twoLower << uint(lastBitsCount)) | uint64(lowerByte>>uint(8-lastBitsCount))
					twoBitCount = twoTotalBitCount
					twoTotalBitCount = 0
					twoByteCount = 0
				}
			}
		}
		oneResultUpper := oneUpper
		oneResultLower := oneLower
		twoResultUpper := twoUpper
		twoResultLower := twoLower
		if twoBitCount == oneBitCount {
			//no adjustment required, compare the values straight up
			oneBitCount = 0
			twoBitCount = 0
		} else {
			diffBits := twoBitCount - oneBitCount
			if diffBits > 0 {
				twoResultUpper >>= uint(diffBits) //look at the high bits only (we are comparing left to right, high to low)
				twoResultLower >>= uint(diffBits)
				mask := ^(^uint64(0) << uint(diffBits))
				twoUpper &= mask
				twoLower &= mask
				twoBitCount = diffBits
				oneBitCount = 0
			} else {
				diffBits = -diffBits
				oneResultUpper >>= uint(diffBits)
				oneResultLower >>= uint(diffBits)
				mask := ^(^uint64(0) << uint(diffBits))
				oneUpper &= mask
				oneLower &= mask
				oneBitCount = diffBits
				twoBitCount = 0
			}
		}
		result := comp.compareValues(oneResultUpper, oneResultLower, twoResultUpper, twoResultLower)
		if result != 0 {
			return result
		}
	}
	return 0
}

func (countComparator) compareSegValues(oneUpper, oneLower, twoUpper, twoLower SegInt) int {
	size1 := oneUpper - oneLower
	size2 := twoUpper - twoLower
	if size1 == size2 {
		//the size of the range is the same, so just compare either upper or lower values
		if oneLower == twoLower {
			return 0
		} else if oneLower > twoLower {
			return 1
		}
	} else if size1 > size2 {
		return 1
	}
	return -1
}

func (countComparator) compareValues(oneUpper, oneLower, twoUpper, twoLower uint64) int {
	size1 := oneUpper - oneLower
	size2 := twoUpper - twoLower
	if size1 == size2 {
		//the size of the range is the same, so just compare either upper or lower values
		if oneLower == twoLower {
			return 0
		} else if oneLower > twoLower {
			return 1
		}
	} else if size1 > size2 {
		return 1
	}
	return -1
}

func (countComparator) compareLargeValues(oneUpper, oneLower, twoUpper, twoLower *big.Int) (result int) {
	oneUpper.Sub(oneUpper, oneLower)
	twoUpper.Sub(twoUpper, twoLower)
	result = oneUpper.CmpAbs(twoUpper)
	if result == 0 {
		//the size of the range is the same, so just compare either upper or lower values
		result = oneLower.CmpAbs(twoLower)
	}
	return
}

func compareDivBitCounts(oneSeries, twoSeries AddressDivisionSeries) int {
	//when this is called we know the two series have the same bit-size, we want to check that the divisions
	//also have the same bit size (which of course also implies that there are the same number of divisions)
	count := oneSeries.GetDivisionCount()
	result := count - twoSeries.GetDivisionCount()
	if result == 0 {
		for i := 0; i < count; i++ {
			result = int(oneSeries.GetGenericDivision(i).GetBitCount() - twoSeries.GetGenericDivision(i).GetBitCount())
			if result != 0 {
				break
			}
		}
	}
	return result
}

func isNilItem(item AddressItem) bool {
	if divSeries, ok := item.(AddressDivisionSeries); ok {
		if addr, ok := divSeries.(AddressType); ok {
			return addr.ToAddressBase() == nil
		} else if grouping, ok := divSeries.(StandardDivGroupingType); ok {
			return grouping.ToDivGrouping() == nil
		} else if largeGrouping, ok := divSeries.(*IPAddressLargeDivisionGrouping); ok {
			return largeGrouping.isNil()
		} // else a type external to this library, which we cannot test for nil
		//} else if rng, ok := item.(IPAddressSeqRangeType); ok {
		//	return rng.ToIP() == nil
	} else if rng, ok := item.(IPAddressSeqRangeType); ok {
		return rng.ToIP() == nil
	} else if div, ok := item.(DivisionType); ok {
		if sdiv, ok := div.(StandardDivisionType); ok {
			return sdiv.ToDiv() == nil
		} else if ldiv, ok := div.(*IPAddressLargeDivision); ok {
			return ldiv.isNil()
		} // else a type external to this library, which we cannot test for nil
	}
	return item == nil
}

// Note: never called with an address instance, never called with an instance of AddressType
func compareCount(one, two AddressItem) int {
	if !one.IsMultiple() {
		if two.IsMultiple() {
			return -1
		}
		return 0
	} else if !two.IsMultiple() {
		return 1
	}
	b1, u1 := getCount(one)
	b2, u2 := getCount(two)
	if b1 == nil {
		if b2 != nil {
			if b2.IsUint64() {
				u2 = b2.Uint64()
			} else {
				return -1
			}
		}
		if u1 < u2 {
			return -1
		} else if u1 == u2 {
			return 0
		}
		return 1
	} else if b2 == nil {
		if b1.IsUint64() {
			u1 = b1.Uint64()
			if u1 < u2 {
				return -1
			} else if u1 == u2 {
				return 0
			}
		}
		return 1
	}
	return b1.CmpAbs(b2)
}

// Note: never called with an address instance, never called with an instance of AddressType
func getCount(item AddressItem) (b *big.Int, u uint64) {
	if sect, ok := item.(StandardDivGroupingType); ok {
		grouping := sect.ToDivGrouping()
		if grouping != nil {
			b = grouping.getCachedCount()
		}
		//} else if rng, ok := item.(IPAddressSeqRangeType); ok {
		//	//ipRange := rng.ToIP()
		//	//if ipRange != nil {
		//	b = rng.GetCount()
		//	//b = rng.getCachedCount(false)
		//	//}
	} else if rng, ok := item.(IPAddressSeqRangeType); ok {
		//ipRange := rng.ToIP()
		//if ipRange != nil {
		b = rng.GetCount()
		//b = rng.getCachedCount(false)
		//}
	} else if div, ok := item.(StandardDivisionType); ok {
		base := div.ToDiv()
		if base != nil {
			if segBase := base.ToSegmentBase(); segBase != nil {
				u = uint64((segBase.getUpperSegmentValue() - base.getSegmentValue()) + 1)
			} else {
				r := base.getUpperDivisionValue() - base.getDivisionValue()
				if r == 0xffffffffffffffff {
					b = bigZero().SetUint64(0xffffffffffffffff)
					b.Add(b, bigOneConst())
					return
				}
				u = r + 1
			}
		}
	} else if lgrouping, ok := item.(*IPAddressLargeDivisionGrouping); ok {
		if lgrouping != nil {
			b = lgrouping.getCachedCount()
		}
	} else if ldiv, ok := item.(*IPAddressLargeDivision); ok {
		if ldiv != nil {
			b = ldiv.getCount()
		}
	} else {
		b = item.GetCount()
	}
	return
}
