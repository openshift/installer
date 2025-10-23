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
	"fmt"
	"math/big"
	"strconv"
	"unsafe"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstr"
)

var zeroSection = createSection(zeroDivs, nil, zeroType)

func createSection(segments []*AddressDivision, prefixLength PrefixLen, addrType addrType) *AddressSection {
	sect := &AddressSection{
		addressSectionInternal{
			addressDivisionGroupingInternal{
				addressDivisionGroupingBase: addressDivisionGroupingBase{
					divisions:    standardDivArray(segments),
					prefixLength: prefixLength,
					addrType:     addrType,
					cache:        &valueCache{},
				},
			},
		},
	}
	assignStringCache(&sect.addressDivisionGroupingBase, addrType)
	return sect
}

// callers to this function supply segments with prefix length consistent with the supplied prefix length
func createSectionMultiple(segments []*AddressDivision, prefixLength PrefixLen, addrType addrType, isMultiple bool) *AddressSection {
	result := createSection(segments, prefixLength, addrType)
	result.isMult = isMultiple
	return result
}

// callers to this function supply segments with prefix length consistent with the supplied prefix length
func createInitializedSection(segments []*AddressDivision, prefixLength PrefixLen, addrType addrType) *AddressSection {
	result := createSection(segments, prefixLength, addrType)
	result.initMultiple() // assigns isMultiple
	return result
}

// callers to this function supply segments with prefix length consistent with the supplied prefix length
func deriveAddressSectionPrefLen(from *AddressSection, segments []*AddressDivision, prefixLength PrefixLen) *AddressSection {
	result := createSection(segments, prefixLength, from.getAddrType())
	result.initMultiple() // assigns isMultiple
	return result
}

// callers to this function supply segments with prefix length consistent with the prefix length of this section
func deriveAddressSection(from *AddressSection, segments []*AddressDivision) (res *AddressSection) {
	return deriveAddressSectionPrefLen(from, segments, from.prefixLength)
}

func assignStringCache(section *addressDivisionGroupingBase, addrType addrType) {
	stringCache := &section.cache.stringCache
	if addrType.isIPv4() {
		stringCache.ipStringCache = &ipStringCache{}
		stringCache.ipv4StringCache = &ipv4StringCache{}
	} else if addrType.isIPv6() {
		stringCache.ipStringCache = &ipStringCache{}
		stringCache.ipv6StringCache = &ipv6StringCache{}
	} else if addrType.isMAC() {
		stringCache.macStringCache = &macStringCache{}
	}
}

// ////////////////////////////////////////////////////////////////
type addressSectionInternal struct {
	addressDivisionGroupingInternal
}

func (section *addressSectionInternal) initImplicitPrefLen(bitsPerSegment BitCount) {
	segCount := section.GetSegmentCount()
	if segCount != 0 {
		for i := segCount - 1; i >= 0; i-- {
			segment := section.GetSegment(i)
			minPref := segment.GetMinPrefixLenForBlock()
			if minPref > 0 {
				if minPref != bitsPerSegment || i != segCount-1 {
					section.prefixLength = getNetworkPrefixLen(bitsPerSegment, minPref, i)
				}
				return
			}
		}
		section.prefixLength = cacheBitCount(0)
	}
}

func (section *addressSectionInternal) initMultAndImplicitPrefLen(bitsPerSegment BitCount) {
	segCount := section.GetSegmentCount()
	if segCount != 0 {
		isMultiple := false
		isBlock := true
		for i := segCount - 1; i >= 0; i-- {
			segment := section.GetSegment(i)
			if isBlock {
				minPref := segment.GetMinPrefixLenForBlock()
				if minPref > 0 {
					if minPref != bitsPerSegment || i != segCount-1 {
						section.prefixLength = getNetworkPrefixLen(bitsPerSegment, minPref, i)
					}
					isBlock = false
					if isMultiple { // nothing left to do
						return
					}
				}
			}
			if !isMultiple && segment.isMultiple() {
				isMultiple = true
				section.isMult = true
				if !isBlock { // nothing left to do
					return
				}
			}
		}
		if isBlock {
			section.prefixLength = cacheBitCount(0)
		}
	}
}

func createDivisionsFromSegs(
	segProvider func(index int) *IPAddressSegment,
	segCount int,
	bitsToSegmentShift uint,
	bitsPerSegment BitCount,
	bytesPerSegment int,
	maxValuePerSegment SegInt,
	zeroSeg, zeroSegZeroPrefix, zeroSegPrefixBlock *IPAddressSegment,
	assignedPrefLen PrefixLen) (divs []*AddressDivision, newPref PrefixLen, isMultiple bool) {
	divs = make([]*AddressDivision, segCount)

	prefixedSegment := -1
	if assignedPrefLen != nil {
		p := assignedPrefLen.bitCount()
		if p < 0 {
			p = 0
			assignedPrefLen = cacheBitCount(p)
		} else {
			boundaryBits := BitCount(segCount << bitsToSegmentShift)
			if p > boundaryBits {
				p = boundaryBits
				assignedPrefLen = cacheBitCount(p)
			}
		}
		prefixedSegment = getNetworkSegmentIndex(p, bytesPerSegment, bitsPerSegment)
	}
	var previousSegPrefixed bool
	var lastSegment *IPAddressSegment
	for i := 0; i < segCount; i++ {
		segment := segProvider(i)
		if segment == nil {
			if previousSegPrefixed {
				divs[i] = zeroSegZeroPrefix.ToDiv()
			} else if i == prefixedSegment {
				newPref = cachePrefixLen(assignedPrefLen)
				segPref := getPrefixedSegmentPrefixLength(bitsPerSegment, assignedPrefLen.bitCount(), prefixedSegment)
				if i+1 < segCount && isPrefixSubnet(
					func(segmentIndex int) SegInt {
						seg := segProvider(segmentIndex + i + 1)
						if seg == nil {
							return 0
						}
						return seg.GetSegmentValue()
					},
					func(segmentIndex int) SegInt {
						seg := segProvider(segmentIndex + i + 1)
						if seg == nil {
							return 0
						}
						return seg.GetUpperSegmentValue()
					},
					segCount-(i+1), bytesPerSegment, bitsPerSegment, maxValuePerSegment, 0, zerosOnly) {
					divs[i] = zeroSeg.toPrefixedNetworkDivision(segPref)
					i++
					isMultiple = isMultiple || i < len(divs) || segPref.bitCount() < bitsPerSegment
					for ; i < len(divs); i++ {
						divs[i] = zeroSegPrefixBlock.ToDiv()
					}
					break
				} else {
					divs[i] = zeroSeg.toPrefixedNetworkDivision(segPref)
				}
			} else {
				divs[i] = zeroSeg.ToDiv() // nil segs are just zero
			}
		} else {
			// The final prefix length is the minimum amongst the assigned one and all of the segments' own prefixes
			segPrefix := segment.getDivisionPrefixLength()
			segIsPrefixed := segPrefix != nil
			if previousSegPrefixed {
				if !segIsPrefixed || segPrefix.bitCount() != 0 {
					divs[i] = createAddressDivision(
						segment.derivePrefixed(cacheBitCount(0))) // change seg prefix to 0
				} else {
					divs[i] = segment.ToDiv() // seg prefix is already 0
				}
			} else {
				// if a prefix length was supplied, we must check for prefix subnets
				var segPrefixSwitch bool
				var assignedSegPref PrefixLen
				if i == prefixedSegment || (prefixedSegment > 0 && segIsPrefixed) {
					// there exists an assigned prefix length
					assignedSegPref = getPrefixedSegmentPrefixLength(bitsPerSegment, assignedPrefLen.bitCount(), i)
					if segIsPrefixed {
						if assignedSegPref == nil || segPrefix.bitCount() < assignedSegPref.bitCount() {
							if segPrefix.bitCount() == 0 && i > 0 {
								// normalize boundaries by looking back
								if !lastSegment.IsPrefixed() {
									divs[i-1] = createAddressDivision(
										lastSegment.derivePrefixed(cacheBitCount(bitsPerSegment)))
								}
							}
							newPref = getNetworkPrefixLen(bitsPerSegment, segPrefix.bitCount(), i)
						} else {
							newPref = cachePrefixLen(assignedPrefLen)
							segPrefixSwitch = assignedSegPref.bitCount() < segPrefix.bitCount()
						}
					} else {
						newPref = cachePrefixLen(assignedPrefLen)
						segPrefixSwitch = true
					}
					if isPrefixSubnet(
						func(segmentIndex int) SegInt {
							seg := segProvider(segmentIndex)
							if seg == nil {
								return 0
							}
							return seg.GetSegmentValue()
						},
						func(segmentIndex int) SegInt {
							seg := segProvider(segmentIndex)
							if seg == nil {
								return 0
							}
							return seg.GetUpperSegmentValue()
						},
						segCount,
						bytesPerSegment,
						bitsPerSegment,
						maxValuePerSegment,
						newPref.bitCount(),
						zerosOnly) {

						divs[i] = segment.toPrefixedNetworkDivision(assignedSegPref)
						i++
						isMultiple = isMultiple || i < len(divs) || newPref.bitCount() < bitsPerSegment
						for ; i < len(divs); i++ {
							divs[i] = zeroSegPrefixBlock.ToDiv()
						}
						break
					}
					previousSegPrefixed = true
				} else if segIsPrefixed {
					if segPrefix.bitCount() == 0 && i > 0 {
						// normalize boundaries by looking back
						if !lastSegment.IsPrefixed() {
							divs[i-1] = createAddressDivision(lastSegment.derivePrefixed(cacheBitCount(bitsPerSegment)))
						}
					}
					newPref = getNetworkPrefixLen(bitsPerSegment, segPrefix.bitCount(), i)
					previousSegPrefixed = true
				}
				if segPrefixSwitch {
					divs[i] = createAddressDivision(segment.derivePrefixed(assignedSegPref)) // change seg prefix
				} else {
					divs[i] = segment.ToDiv()
				}
			}
			isMultiple = isMultiple || segment.isMultiple()
		}
		lastSegment = segment
	}
	return
}

func (section *addressSectionInternal) matchesTypeAndCount(other *AddressSection) (matches bool, count int) {
	count = section.GetDivisionCount()
	if count != other.GetDivisionCount() {
		return
	} else if section.getAddrType() != other.getAddrType() {
		return
	}
	matches = true
	return
}

func (section *addressSectionInternal) equal(otherT AddressSectionType) bool {
	if otherT == nil {
		return false
	}
	other := otherT.ToSectionBase()
	if other == nil {
		return false
	}
	matchesStructure, _ := section.matchesTypeAndCount(other)
	return matchesStructure && section.sameCountTypeEquals(other)
}

func (section *addressSectionInternal) sameCountTypeEquals(other *AddressSection) bool {
	count := section.GetSegmentCount()
	for i := count - 1; i >= 0; i-- {
		if !section.GetSegment(i).sameTypeEquals(other.GetSegment(i)) {
			return false
		}
	}
	return true
}

func (section *addressSectionInternal) sameCountTypeContains(other *AddressSection) bool {
	count := section.GetSegmentCount()
	for i := count - 1; i >= 0; i-- {
		if !section.GetSegment(i).sameTypeContains(other.GetSegment(i)) {
			return false
		}
	}
	return true
}

// GetBitsPerSegment returns the number of bits comprising each segment in this section.  Segments in the same address section are equal length.
func (section *addressSectionInternal) GetBitsPerSegment() BitCount {
	addrType := section.getAddrType()
	if addrType.isIPv4() {
		return IPv4BitsPerSegment
	} else if addrType.isIPv6() {
		return IPv6BitsPerSegment
	} else if addrType.isMAC() {
		return MACBitsPerSegment
	}
	if section.GetDivisionCount() == 0 {
		return 0
	}
	return section.getDivision(0).GetBitCount()
}

// GetBytesPerSegment returns the number of bytes comprising each segment in this section.  Segments in the same address section are equal length.
func (section *addressSectionInternal) GetBytesPerSegment() int {
	addrType := section.getAddrType()
	if addrType.isIPv4() {
		return IPv4BytesPerSegment
	} else if addrType.isIPv6() {
		return IPv6BytesPerSegment
	} else if addrType.isMAC() {
		return MACBytesPerSegment
	}
	if section.GetDivisionCount() == 0 {
		return 0
	}
	return section.getDivision(0).GetByteCount()
}

// GetSegment returns the segment at the given index.
// The first segment is at index 0.
// GetSegment will panic given a negative index or an index matching or larger than the segment count.
func (section *addressSectionInternal) GetSegment(index int) *AddressSegment {
	return section.getDivision(index).ToSegmentBase()
}

// GetGenericSegment returns the segment as an AddressSegmentType,
// allowing all segment types to be represented by a single type.
// The first segment is at index 0.
// GetGenericSegment will panic given a negative index or an index matching or larger than the segment count.
func (section *addressSectionInternal) GetGenericSegment(index int) AddressSegmentType {
	return section.GetSegment(index)
}

// GetSegmentCount returns the segment count.
func (section *addressSectionInternal) GetSegmentCount() int {
	return section.GetDivisionCount()
}

// ForEachSegment visits each segment in order from most-significant to least, the most significant with index 0, calling the given function for each, terminating early if the function returns true.
// Returns the number of visited segments.
func (section *addressSectionInternal) ForEachSegment(consumer func(segmentIndex int, segment *AddressSegment) (stop bool)) int {
	divArray := section.getDivArray()
	if divArray != nil {
		for i, div := range divArray {
			if consumer(i, div.ToSegmentBase()) {
				return i + 1
			}
		}
	}
	return len(divArray)
}

// GetBitCount returns the number of bits in each value comprising this address item.
func (section *addressSectionInternal) GetBitCount() BitCount {
	divLen := section.GetDivisionCount()
	if divLen == 0 {
		return 0
	}
	return getSegmentsBitCount(section.getDivision(0).GetBitCount(), section.GetSegmentCount())
}

// GetByteCount returns the number of bytes required for each value comprising this address item.
func (section *addressSectionInternal) GetByteCount() int {
	return int((section.GetBitCount() + 7) >> 3)
}

// GetMaxSegmentValue returns the maximum possible segment value for this type of address.
//
// Note this is not the maximum of the range of segment values in this specific address,
// this is the maximum value of any segment for this address type and version, determined by the number of bits per segment.
func (section *addressSectionInternal) GetMaxSegmentValue() SegInt {
	addrType := section.getAddrType()
	if addrType.isIPv4() {
		return IPv4MaxValuePerSegment
	} else if addrType.isIPv6() {
		return IPv6MaxValuePerSegment
	} else if addrType.isMAC() {
		return MACMaxValuePerSegment
	}
	divLen := section.GetDivisionCount()
	if divLen == 0 {
		return 0
	}
	return section.GetSegment(0).GetMaxValue()
}

// TestBit returns true if the bit in the lower value of this section at the given index is 1, where index 0 refers to the least significant bit.
// In other words, it computes (bits & (1 << n)) != 0), using the lower value of this section.
// TestBit will panic if n < 0, or if it matches or exceeds the bit count of this item.
func (section *addressSectionInternal) TestBit(n BitCount) bool {
	return section.IsOneBit(section.GetBitCount() - (n + 1))
}

// IsOneBit returns true if the bit in the lower value of this section at the given index is 1, where index 0 refers to the most significant bit.
// IsOneBit will panic if bitIndex is less than zero, or if it is larger than the bit count of this item.
func (section *addressSectionInternal) IsOneBit(prefixBitIndex BitCount) bool {
	bitsPerSegment := section.GetBitsPerSegment()
	bytesPerSegment := section.GetBytesPerSegment()
	segment := section.GetSegment(getHostSegmentIndex(prefixBitIndex, bytesPerSegment, bitsPerSegment))
	segmentBitIndex := prefixBitIndex % bitsPerSegment
	return segment.IsOneBit(segmentBitIndex)
}

// Gets the subsection from the series starting from the given index and ending just before the give endIndex.
// The first segment is at index 0.
func (section *addressSectionInternal) getSubSection(index, endIndex int) *AddressSection {
	if index < 0 {
		index = 0
	}
	thisSegmentCount := section.GetSegmentCount()
	if endIndex > thisSegmentCount {
		endIndex = thisSegmentCount
	}
	segmentCount := endIndex - index
	if segmentCount <= 0 {
		if thisSegmentCount == 0 {
			return section.toAddressSection()
		}
		// we do not want an inconsistency where mac zero length can have prefix len zero while ip sections cannot
		return zeroSection
	}
	if index == 0 && endIndex == thisSegmentCount {
		return section.toAddressSection()
	}
	segs := section.getSubDivisions(index, endIndex)
	newPrefLen := section.getPrefixLen()
	if newPrefLen != nil {
		newPrefLen = getAdjustedPrefixLength(section.GetBitsPerSegment(), newPrefLen.bitCount(), index, endIndex)
	}
	addrType := section.getAddrType()
	if !section.isMultiple() {
		return createSection(segs, newPrefLen, addrType)
	}
	return deriveAddressSectionPrefLen(section.toAddressSection(), segs, newPrefLen)
}

func (section *addressSectionInternal) getLowestHighestSections() (lower, upper *AddressSection) {
	if !section.isMultiple() {
		lower = section.toAddressSection()
		upper = lower
		return
	}
	cache := section.cache
	if cache == nil {
		return section.createLowestHighestSections()
	}
	cached := (*groupingCache)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache.sectionCache))))
	if cached == nil {
		cached = &groupingCache{}
		cached.lower, cached.upper = section.createLowestHighestSections()
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.sectionCache))
		atomicStorePointer(dataLoc, unsafe.Pointer(cached))
	}
	lower = cached.lower
	upper = cached.upper
	return
}

func (section *addressSectionInternal) createLowestHighestSections() (lower, upper *AddressSection) {
	segmentCount := section.GetSegmentCount()
	lowSegs := createSegmentArray(segmentCount)
	var highSegs []*AddressDivision
	if section.isMultiple() {
		highSegs = createSegmentArray(segmentCount)
	}
	for i := 0; i < segmentCount; i++ {
		seg := section.GetSegment(i)
		lowSegs[i] = seg.GetLower().ToDiv()
		if highSegs != nil {
			highSegs[i] = seg.GetUpper().ToDiv()
		}
	}
	lower = deriveAddressSection(section.toAddressSection(), lowSegs)
	if highSegs == nil {
		upper = lower
	} else {
		upper = deriveAddressSection(section.toAddressSection(), highSegs)
	}
	return
}

// Returns the address created by converting this address to an address with a 0 as the first bit following the prefix, followed by all ones to the end, and with the prefix length then removed
// Returns the same address if it has no prefix length.
func (section *addressSectionInternal) toMaxLower() *AddressSection {
	return section.toAboveOrBelow(false)
}

// Returns the address created by converting this address to an address with a 1 as the first bit following the prefix, followed by all zeros to the end, and with the prefix length then removed
// Returns the same address if it has no prefix length
func (section *addressSectionInternal) toMinUpper() *AddressSection {
	return section.toAboveOrBelow(true)
}

func (section *addressSectionInternal) toAboveOrBelow(above bool) *AddressSection {
	prefLen := section.GetPrefixLen()
	if prefLen == nil {
		return section.toAddressSection()
	}
	prefBits := prefLen.Len()
	segmentCount := section.GetSegmentCount()
	if prefBits == section.GetBitCount() || segmentCount == 0 {
		return section.withoutPrefixLen()
	}
	segmentByteCount := section.GetBytesPerSegment()
	segmentBitCount := section.GetBitsPerSegment()
	newSegs := createSegmentArray(segmentCount)
	if prefBits > 0 {
		networkSegmentIndex := getNetworkSegmentIndex(prefBits, segmentByteCount, segmentBitCount)
		section.copySubDivisions(0, networkSegmentIndex, newSegs)
	}
	hostSegmentIndex := getHostSegmentIndex(prefBits, segmentByteCount, segmentBitCount)
	if hostSegmentIndex < segmentCount {
		oldSeg := section.getDivision(hostSegmentIndex)
		oldVal := oldSeg.getUpperSegmentValue()
		segPrefBits := getPrefixedSegmentPrefixLength(segmentBitCount, prefBits, hostSegmentIndex).bitCount()
		// 1 bit followed by zeros
		allOnes := ^SegInt(0)
		var newVal SegInt
		if above {
			hostBits := uint(segmentBitCount - segPrefBits)
			networkMask := allOnes << (hostBits - 1)
			hostMask := ^(allOnes << hostBits)
			newVal = (oldVal | hostMask) & networkMask
		} else {
			hostBits := uint(segmentBitCount - segPrefBits)
			networkMask := allOnes << hostBits
			hostMask := ^(allOnes<<hostBits - 1)
			newVal = (oldVal & networkMask) | hostMask
		}
		newSegs[hostSegmentIndex] = createAddressDivision(oldSeg.deriveNewSeg(newVal, nil))
		if j := hostSegmentIndex + 1; j < segmentCount {
			var endSeg *AddressDivision
			if above {
				endSeg = createAddressDivision(oldSeg.deriveNewSeg(0, nil))
			} else {
				maxSegVal := section.GetMaxSegmentValue()
				endSeg = createAddressDivision(oldSeg.deriveNewSeg(maxSegVal, nil))
			}
			newSegs[j] = endSeg
			for j++; j < segmentCount; j++ {
				newSegs[j] = endSeg
			}
		}
	}
	return deriveAddressSectionPrefLen(section.toAddressSection(), newSegs, nil)
}

func (section *addressSectionInternal) reverseSegments(segProducer func(int) (*AddressSegment, addrerr.IncompatibleAddressError)) (res *AddressSection, err addrerr.IncompatibleAddressError) {
	count := section.GetSegmentCount()
	if count == 0 { // case count == 1 we cannot exit early, we need to apply segProducer to each segment
		return section.withoutPrefixLen(), nil
	}
	newSegs := createSegmentArray(count)
	halfCount := count >> 1
	i := 0
	isSame := !section.isPrefixed() //when reversing, the prefix must go
	for j := count - 1; i < halfCount; i, j = i+1, j-1 {
		var newj, newi *AddressSegment
		if newj, err = segProducer(i); err != nil {
			return
		}
		if newi, err = segProducer(j); err != nil {
			return
		}
		origi := section.GetSegment(i)
		origj := section.GetSegment(j)
		newSegs[j] = newj.ToDiv()
		newSegs[i] = newi.ToDiv()
		if isSame &&
			!(segValsSame(newi.getSegmentValue(), origi.getSegmentValue(), newi.getUpperSegmentValue(), origi.getUpperSegmentValue()) &&
				segValsSame(newj.getSegmentValue(), origj.getSegmentValue(), newj.getUpperSegmentValue(), origj.getUpperSegmentValue())) {
			isSame = false
		}
	}
	if (count & 1) == 1 { //the count is odd, handle the middle one
		seg := section.getDivision(i)
		newSegs[i] = seg // gets segment i without prefix length
	}
	if isSame {
		res = section.toAddressSection()
		return
	}
	res = deriveAddressSectionPrefLen(section.toAddressSection(), newSegs, nil)
	return
}

func (section *addressSectionInternal) reverseBits(perByte bool) (res *AddressSection, err addrerr.IncompatibleAddressError) {
	if perByte {
		isSame := !section.isPrefixed() //when reversing, the prefix must go
		count := section.GetSegmentCount()
		newSegs := createSegmentArray(count)
		for i := 0; i < count; i++ {
			seg := section.GetSegment(i)
			var reversedSeg *AddressSegment
			reversedSeg, err = seg.ReverseBits(perByte)
			if err != nil {
				return
			}
			newSegs[i] = reversedSeg.ToDiv()
			if isSame && !segValsSame(seg.getSegmentValue(), reversedSeg.getSegmentValue(), seg.getUpperSegmentValue(), reversedSeg.getUpperSegmentValue()) {
				isSame = false
			}
		}
		if isSame {
			res = section.toAddressSection() //We can do this because for ipv6 startIndex stays the same and for mac startIndex and extended stays the same
			return
		}
		res = deriveAddressSectionPrefLen(section.toAddressSection(), newSegs, nil)
		return
	}
	return section.reverseSegments(
		func(i int) (*AddressSegment, addrerr.IncompatibleAddressError) {
			return section.GetSegment(i).ReverseBits(perByte)
		},
	)
}

func (section *addressSectionInternal) reverseBytes(perSegment bool) (res *AddressSection, err addrerr.IncompatibleAddressError) {
	if perSegment {
		isSame := !section.isPrefixed() //when reversing, the prefix must go
		count := section.GetSegmentCount()
		newSegs := createSegmentArray(count)
		for i := 0; i < count; i++ {
			seg := section.GetSegment(i)
			var reversedSeg *AddressSegment
			reversedSeg, err = seg.ReverseBytes()
			if err != nil {
				return
			}
			newSegs[i] = reversedSeg.ToDiv()
			if isSame && !segValsSame(seg.getSegmentValue(), reversedSeg.getSegmentValue(), seg.getUpperSegmentValue(), reversedSeg.getUpperSegmentValue()) {
				isSame = false
			}
		}
		if isSame {
			res = section.toAddressSection() //We can do this because for ipv6 startIndex stays the same and for mac startIndex and extended stays the same
			return
		}
		res = deriveAddressSectionPrefLen(section.toAddressSection(), newSegs, nil)
		return
	}
	return section.reverseSegments(
		func(i int) (*AddressSegment, addrerr.IncompatibleAddressError) {
			return section.GetSegment(i).ReverseBytes()
		},
	)
}

// callers to replace have ensures the component sections have consistent prefix lengths for the replacement
func (section *addressSectionInternal) replace(
	index,
	endIndex int,
	replacement *AddressSection,
	replacementStartIndex,
	replacementEndIndex int,
	prefixLen PrefixLen) *AddressSection {

	otherSegmentCount := replacementEndIndex - replacementStartIndex
	segmentCount := section.GetSegmentCount()
	totalSegmentCount := segmentCount + otherSegmentCount - (endIndex - index)
	segs := createSegmentArray(totalSegmentCount)
	sect := section.toAddressSection()
	sect.copySubDivisions(0, index, segs)
	if index < totalSegmentCount {
		replacement.copySubDivisions(replacementStartIndex, replacementEndIndex, segs[index:])
		if index+otherSegmentCount < totalSegmentCount {
			sect.copySubDivisions(endIndex, segmentCount, segs[index+otherSegmentCount:])
		}
	}
	addrType := sect.getAddrType()
	if addrType.isZeroSegments() { // zero-length section
		addrType = replacement.getAddrType()
	}
	return createInitializedSection(segs, prefixLen, addrType)
}

// Replaces segments starting from startIndex and ending before endIndex with the segments starting at replacementStartIndex and
// ending before replacementEndIndex from the replacement section.
func (section *addressSectionInternal) replaceLen(startIndex, endIndex int, replacement *AddressSection, replacementStartIndex, replacementEndIndex int, segmentToBitsShift uint) *AddressSection {
	segmentCount := section.GetSegmentCount()
	startIndex, endIndex, replacementStartIndex, replacementEndIndex =
		adjustIndices(startIndex, endIndex, segmentCount, replacementStartIndex, replacementEndIndex, replacement.GetSegmentCount())

	replacedCount := endIndex - startIndex
	replacementCount := replacementEndIndex - replacementStartIndex

	// unlike ipvx, sections of zero length with 0 prefix are still considered to be applying their prefix during replacement,
	// because you can have zero length prefixes when there are no bits in the section
	prefixLength := section.getPrefixLen()
	if replacementCount == 0 && replacedCount == 0 {
		if prefixLength != nil {
			prefLen := prefixLength.bitCount()
			if prefLen <= BitCount(startIndex<<segmentToBitsShift) {
				return section.toAddressSection()
			} else {
				replacementPrefisLength := replacement.getPrefixLen()
				if replacementPrefisLength == nil {
					return section.toAddressSection()
				} else if replacementPrefisLength.bitCount() > BitCount(replacementStartIndex<<segmentToBitsShift) {
					return section.toAddressSection()
				}
			}
		} else {
			replacementPrefisLength := replacement.getPrefixLen()
			if replacementPrefisLength == nil {
				return section.toAddressSection()
			} else if replacementPrefisLength.bitCount() > BitCount(replacementStartIndex<<segmentToBitsShift) {
				return section.toAddressSection()
			}
		}
	} else if segmentCount == replacedCount {
		if prefixLength == nil || prefixLength.bitCount() > 0 {
			return replacement
		} else {
			replacementPrefisLength := replacement.getPrefixLen()
			if replacementPrefisLength != nil && replacementPrefisLength.bitCount() == 0 { // prefix length is 0
				return replacement
			}
		}
	}

	startBits := BitCount(startIndex << segmentToBitsShift)
	var newPrefixLength PrefixLen
	if prefixLength != nil && prefixLength.bitCount() <= startBits {
		newPrefixLength = prefixLength
	} else {
		replacementPrefLen := replacement.getPrefixLen()
		if replacementPrefLen != nil && replacementPrefLen.bitCount() <= BitCount(replacementEndIndex<<segmentToBitsShift) {
			var replacementPrefixLen BitCount
			replacementStartBits := BitCount(replacementStartIndex << segmentToBitsShift)
			if replacementPrefLen.bitCount() > replacementStartBits {
				replacementPrefixLen = replacementPrefLen.bitCount() - replacementStartBits
			}
			newPrefixLength = cacheBitCount(startBits + replacementPrefixLen)
		} else if prefixLength != nil {
			replacementBits := BitCount(replacementCount << segmentToBitsShift)
			var endPrefixBits BitCount
			endIndexBits := BitCount(endIndex << segmentToBitsShift)
			if prefixLength.bitCount() > endIndexBits {
				endPrefixBits = prefixLength.bitCount() - endIndexBits
			}
			newPrefixLength = cacheBitCount(startBits + replacementBits + endPrefixBits)
		} else {
			newPrefixLength = nil
		}
	}
	result := section.replace(startIndex, endIndex, replacement, replacementStartIndex, replacementEndIndex, newPrefixLength)
	return result
}

func (section *addressSectionInternal) toPrefixBlock() *AddressSection {
	prefixLength := section.getPrefixLen()
	if prefixLength == nil {
		return section.toAddressSection()
	}
	return section.toPrefixBlockLen(prefixLength.bitCount())
}

func (section *addressSectionInternal) toPrefixBlockLen(prefLen BitCount) *AddressSection {
	prefLen = checkSubnet(section, prefLen)
	segCount := section.GetSegmentCount()
	if segCount == 0 {
		return section.toAddressSection()
	}
	segmentByteCount := section.GetBytesPerSegment()
	segmentBitCount := section.GetBitsPerSegment()
	existingPrefixLength := section.getPrefixLen()
	prefixMatches := existingPrefixLength != nil && existingPrefixLength.bitCount() == prefLen
	if prefixMatches {
		prefixedSegmentIndex := getHostSegmentIndex(prefLen, segmentByteCount, segmentBitCount)
		if prefixedSegmentIndex >= segCount {
			return section.toAddressSection()
		}
		segPrefLength := getPrefixedSegmentPrefixLength(segmentBitCount, prefLen, prefixedSegmentIndex).bitCount()
		seg := section.GetSegment(prefixedSegmentIndex)
		if seg.containsPrefixBlock(segPrefLength) {
			i := prefixedSegmentIndex + 1
			for ; i < segCount; i++ {
				seg = section.GetSegment(i)
				if !seg.IsFullRange() {
					break
				}
			}
			if i == segCount {
				return section.toAddressSection()
			}
		}
	}
	prefixedSegmentIndex := 0
	newSegs := createSegmentArray(segCount)
	if prefLen > 0 {
		prefixedSegmentIndex = getNetworkSegmentIndex(prefLen, segmentByteCount, segmentBitCount)
		section.copySubDivisions(0, prefixedSegmentIndex, newSegs)
	}
	for i := prefixedSegmentIndex; i < segCount; i++ {
		segPrefLength := getPrefixedSegmentPrefixLength(segmentBitCount, prefLen, i)
		oldSeg := section.getDivision(i)
		newSegs[i] = oldSeg.toPrefixedNetworkDivision(segPrefLength)
	}
	return createSectionMultiple(newSegs, cacheBitCount(prefLen), section.getAddrType(), section.isMultiple() || prefLen < section.GetBitCount())
}

func (section *addressSectionInternal) toBlock(segmentIndex int, lower, upper SegInt) *AddressSection {
	segCount := section.GetSegmentCount()
	i := segmentIndex
	if i < 0 {
		i = 0
	}
	maxSegVal := section.GetMaxSegmentValue()
	for ; i < segCount; i++ {
		seg := section.GetSegment(segmentIndex)
		var lowerVal, upperVal SegInt
		if i == segmentIndex {
			lowerVal, upperVal = lower, upper
		} else {
			upperVal = maxSegVal
		}
		if !segsSame(nil, seg.getDivisionPrefixLength(), lowerVal, seg.GetSegmentValue(), upperVal, seg.GetUpperSegmentValue()) {
			newSegs := createSegmentArray(segCount)
			section.copySubDivisions(0, i, newSegs)
			newSeg := createAddressDivision(seg.deriveNewMultiSeg(lowerVal, upperVal, nil))
			newSegs[i] = newSeg
			var allSeg *AddressDivision
			if j := i + 1; j < segCount {
				if i == segmentIndex {
					allSeg = createAddressDivision(seg.deriveNewMultiSeg(0, maxSegVal, nil))
				} else {
					allSeg = newSeg
				}
				newSegs[j] = allSeg
				for j++; j < segCount; j++ {
					newSegs[j] = allSeg
				}
			}
			return createSectionMultiple(newSegs, nil, section.getAddrType(),
				segmentIndex < segCount-1 || lower != upper)
		}
	}
	return section.toAddressSection()
}

func (section *addressSectionInternal) withoutPrefixLen() *AddressSection {
	if !section.isPrefixed() {
		return section.toAddressSection()
	}
	if sect := section.toIPAddressSection(); sect != nil {
		return sect.withoutPrefixLen().ToSectionBase()
	}
	return createSectionMultiple(section.getDivisionsInternal(), nil, section.getAddrType(), section.isMultiple())
}

func (section *addressSectionInternal) getAdjustedPrefix(adjustment BitCount) BitCount {
	prefix := section.getPrefixLen()
	bitCount := section.GetBitCount()
	var result BitCount
	if prefix == nil {
		if adjustment > 0 { // start from 0
			if adjustment > bitCount {
				result = bitCount
			} else {
				result = adjustment
			}
		} else { // start from end
			if -adjustment < bitCount {
				result = bitCount + adjustment
			}
		}
	} else {
		result = prefix.bitCount() + adjustment
		if result > bitCount {
			result = bitCount
		} else if result < 0 {
			result = 0
		}
	}
	return result
}

func (section *addressSectionInternal) adjustPrefixLen(adjustment BitCount) *AddressSection {
	// no zeroing
	res, _ := section.adjustPrefixLength(adjustment, false)
	return res
}

func (section *addressSectionInternal) adjustPrefixLenZeroed(adjustment BitCount) (*AddressSection, addrerr.IncompatibleAddressError) {
	return section.adjustPrefixLength(adjustment, true)
}

func (section *addressSectionInternal) adjustPrefixLength(adjustment BitCount, withZeros bool) (*AddressSection, addrerr.IncompatibleAddressError) {
	if adjustment == 0 && section.isPrefixed() {
		return section.toAddressSection(), nil
	}
	prefix := section.getAdjustedPrefix(adjustment)
	return section.setPrefixLength(prefix, withZeros)
}

func (section *addressSectionInternal) setPrefixLen(prefixLen BitCount) *AddressSection {
	// no zeroing
	res, _ := section.setPrefixLength(prefixLen, false)
	return res
}

func (section *addressSectionInternal) setPrefixLenZeroed(prefixLen BitCount) (*AddressSection, addrerr.IncompatibleAddressError) {
	return section.setPrefixLength(prefixLen, true)
}

func (section *addressSectionInternal) setPrefixLength(
	networkPrefixLength BitCount,
	withZeros bool,
) (res *AddressSection, err addrerr.IncompatibleAddressError) {
	existingPrefixLength := section.getPrefixLen()
	if existingPrefixLength != nil && networkPrefixLength == existingPrefixLength.bitCount() {
		res = section.toAddressSection()
		return
	}
	segmentCount := section.GetSegmentCount()
	var appliedPrefixLen PrefixLen // purposely nil when there are no segments
	verifyMask := false
	var startIndex int
	var segmentMaskProducer func(int) SegInt
	if segmentCount != 0 {
		maxVal := section.GetMaxSegmentValue()
		appliedPrefixLen = cacheBitCount(networkPrefixLength)
		var minPrefIndex, maxPrefIndex int
		var minPrefLen, maxPrefLen BitCount
		bitsPerSegment := section.GetBitsPerSegment()
		bytesPerSegment := section.GetBytesPerSegment()
		prefIndex := getNetworkSegmentIndex(networkPrefixLength, bytesPerSegment, bitsPerSegment)
		if existingPrefixLength != nil {
			verifyMask = true
			existingPrefLen := existingPrefixLength.bitCount()
			existingPrefIndex := getNetworkSegmentIndex(existingPrefLen, bytesPerSegment, bitsPerSegment) // can be -1 if existingPrefLen is 0
			if prefIndex > existingPrefIndex {
				maxPrefIndex = prefIndex
				minPrefIndex = existingPrefIndex
			} else {
				maxPrefIndex = existingPrefIndex
				minPrefIndex = prefIndex
			}
			if withZeros {
				if networkPrefixLength < existingPrefLen {
					minPrefLen = networkPrefixLength
					maxPrefLen = existingPrefLen
				} else {
					minPrefLen = existingPrefLen
					maxPrefLen = networkPrefixLength
				}
				startIndex = minPrefIndex
				segmentMaskProducer = func(i int) SegInt {
					if i >= minPrefIndex {
						if i <= maxPrefIndex {
							minSegPrefLen := getPrefixedSegmentPrefixLength(bitsPerSegment, minPrefLen, i).bitCount()
							minMask := maxVal << uint(bitsPerSegment-minSegPrefLen)
							maxSegPrefLen := getPrefixedSegmentPrefixLength(bitsPerSegment, maxPrefLen, i)
							if maxSegPrefLen != nil {
								maxMask := maxVal << uint(bitsPerSegment-maxSegPrefLen.bitCount())
								return minMask | ^maxMask
							}
							return minMask
						}
					}
					return maxVal
				}
			} else {
				startIndex = minPrefIndex
			}
		} else {
			startIndex = prefIndex
		}
		if segmentMaskProducer == nil {
			segmentMaskProducer = func(i int) SegInt {
				return maxVal
			}
		}
	}
	if startIndex < 0 {
		startIndex = 0
	}

	return section.getSubnetSegments(
		startIndex,
		appliedPrefixLen,
		verifyMask,
		func(i int) *AddressDivision {
			return section.getDivision(i)
		},
		segmentMaskProducer,
	)
}

func (section *addressSectionInternal) assignPrefixForSingleBlock() *AddressSection {
	newPrefix := section.GetPrefixLenForSingleBlock()
	if newPrefix == nil {
		return nil
	}
	newSect := section.setPrefixLen(newPrefix.bitCount())
	cache := newSect.cache
	if cache != nil {
		// no atomic writes required since we created this new section in here
		cache.isSinglePrefixBlock = &trueVal
		cache.equivalentPrefix = cachePrefix(newPrefix.bitCount())
		cache.minPrefix = newPrefix
	}
	return newSect
}

// Constructs an equivalent address section with the smallest CIDR prefix possible (largest network),
// such that the range of values are a set of subnet blocks for that prefix.
func (section *addressSectionInternal) assignMinPrefixForBlock() *AddressSection {
	return section.setPrefixLen(section.GetMinPrefixLenForBlock())
}

// PrefixEqual determines if the given section matches this section up to the prefix length of this section.
// It returns whether the argument section has the same address section prefix values as this.
//
// All prefix bits of this section must be present in the other section to be comparable, otherwise false is returned.
func (section *addressSectionInternal) PrefixEqual(other AddressSectionType) (res bool) {
	o := other.ToSectionBase()
	if section.toAddressSection() == o {
		return true
	} else if section.getAddrType() != o.getAddrType() {
		return
	}
	return section.prefixContains(o, false)
}

// PrefixContains returns whether the prefix values in the given address section
// are prefix values in this address section, using the prefix length of this section.
// If this address section has no prefix length, the entire address is compared.
//
// It returns whether the prefix of this address contains all values of the same prefix length in the given address.
//
// All prefix bits of this section must be present in the other section to be comparable.
func (section *addressSectionInternal) PrefixContains(other AddressSectionType) (res bool) {
	o := other.ToSectionBase()
	if section.toAddressSection() == o {
		return true
	} else if section.getAddrType() != o.getAddrType() {
		return
	}
	return section.prefixContains(o, true)
}

func (section *addressSectionInternal) prefixContains(other *AddressSection, contains bool) (res bool) {
	prefixLength := section.getPrefixLen()
	var prefixedSection int
	if prefixLength == nil {
		prefixedSection = section.GetSegmentCount()
		if prefixedSection > other.GetSegmentCount() {
			return
		}
	} else {
		prefLen := prefixLength.bitCount()
		prefixedSection = getNetworkSegmentIndex(prefLen, section.GetBytesPerSegment(), section.GetBitsPerSegment())
		if prefixedSection >= 0 {
			if prefixedSection >= other.GetSegmentCount() {
				return
			}
			one := section.GetSegment(prefixedSection)
			two := other.GetSegment(prefixedSection)
			segPrefixLength := getPrefixedSegmentPrefixLength(one.getBitCount(), prefLen, prefixedSection)
			if contains {
				if !one.PrefixContains(two, segPrefixLength.bitCount()) {
					return
				}
			} else {
				if !one.PrefixEqual(two, segPrefixLength.bitCount()) {
					return
				}
			}
		}
	}
	for prefixedSection--; prefixedSection >= 0; prefixedSection-- {
		one := section.GetSegment(prefixedSection)
		two := other.GetSegment(prefixedSection)
		if contains {
			if !one.Contains(two) {
				return
			}
		} else {
			if !one.equalsSegment(two) {
				return
			}
		}
	}
	return true
}

func (section *addressSectionInternal) contains(other AddressSectionType) bool {
	if other == nil {
		return true
	}
	otherSection := other.ToSectionBase()
	if section.toAddressSection() == otherSection || otherSection == nil {
		return true
	}
	//check if they are comparable first
	matches, count := section.matchesTypeAndCount(otherSection)
	if !matches {
		return false
	}
	for i := count - 1; i >= 0; i-- {
		if !section.GetSegment(i).sameTypeContains(otherSection.GetSegment(i)) {
			return false
		}
	}
	return true
}

func (section *addressSectionInternal) overlaps(other AddressSectionType) bool {
	if other == nil {
		return true
	}
	otherSection := other.ToSectionBase()
	if section.toAddressSection() == otherSection || otherSection == nil {
		return true
	}
	//check if they are comparable first
	matches, count := section.matchesTypeAndCount(otherSection)
	if !matches {
		return false
	}
	for i := count - 1; i >= 0; i-- {
		if !section.GetSegment(i).sameTypeOverlaps(otherSection.GetSegment(i)) {
			return false
		}
	}
	return true
}

func low64(section *AddressSection) (res uint64) {
	return low64Bytes(section.Bytes())
}

func low64Upper(section *AddressSection) (res uint64) {
	return low64Bytes(section.UpperBytes())
}

func low64Bytes(bts []byte) (res uint64) {
	byteLen := len(bts)
	if byteLen > 8 {
		bts = bts[byteLen-8:]
	} else if byteLen > 0 {
		res = uint64(bts[0])
		for i := 1; i < byteLen; i++ {
			res = (res << 8) | uint64(bts[i])
		}
	}
	return
}

func (section *addressSectionInternal) enumerate(other AddressSectionType) *big.Int {
	if sect := section.toIPv4AddressSection(); sect != nil {
		return sect.Enumerate(other)
	} else if sect := section.toIPv6AddressSection(); sect != nil {
		return sect.Enumerate(other)
	} else if sect := section.toMACAddressSection(); sect != nil {
		return sect.Enumerate(other)
	}
	if other != nil {
		if otherSection := other.ToSectionBase(); otherSection != nil {
			if matches, _ := section.matchesTypeAndCount(otherSection); matches {
				return enumerateBig(section.toAddressSection(), otherSection, low64, low64Upper)
			}
		}
	}
	return nil
}

func (section *addressSectionInternal) getStringCache() *stringCache {
	if section.hasNoDivisions() {
		return &zeroStringCache
	}
	cache := section.cache
	if cache == nil {
		return nil
	}
	return &cache.stringCache
}

func (section *addressSectionInternal) getLower() *AddressSection {
	lower, _ := section.getLowestHighestSections()
	return lower
}

func (section *addressSectionInternal) getUpper() *AddressSection {
	_, upper := section.getLowestHighestSections()
	return upper
}

func (section *addressSectionInternal) incrementBoundary(increment int64) *AddressSection {
	if increment <= 0 {
		if increment == 0 {
			return section.toAddressSection()
		}
		return section.getLower().increment(increment)
	}
	return section.getUpper().increment(increment)
}

func (section *addressSectionInternal) increment(increment int64) *AddressSection {
	if sect := section.toIPv4AddressSection(); sect != nil {
		return sect.Increment(increment).ToSectionBase()
	} else if sect := section.toIPv6AddressSection(); sect != nil {
		return sect.Increment(increment).ToSectionBase()
	} else if sect := section.toMACAddressSection(); sect != nil {
		return sect.Increment(increment).ToSectionBase()
	}
	return nil
}

var (
	otherOctalPrefix = "0o"
	otherHexPrefix   = "0X"

	//decimalParams            = new(StringOptionsBuilder).SetRadix(10).SetExpandedSegments(true).ToOptions()
	hexParams                  = new(addrstr.StringOptionsBuilder).SetRadix(16).SetHasSeparator(false).SetExpandedSegments(true).ToOptions()
	hexUppercaseParams         = new(addrstr.StringOptionsBuilder).SetRadix(16).SetHasSeparator(false).SetExpandedSegments(true).SetUppercase(true).ToOptions()
	hexPrefixedParams          = new(addrstr.StringOptionsBuilder).SetRadix(16).SetHasSeparator(false).SetExpandedSegments(true).SetAddressLabel(HexPrefix).ToOptions()
	hexPrefixedUppercaseParams = new(addrstr.StringOptionsBuilder).SetRadix(16).SetHasSeparator(false).SetExpandedSegments(true).SetAddressLabel(HexPrefix).SetUppercase(true).ToOptions()
	octalParams                = new(addrstr.StringOptionsBuilder).SetRadix(8).SetHasSeparator(false).SetExpandedSegments(true).ToOptions()
	octalPrefixedParams        = new(addrstr.StringOptionsBuilder).SetRadix(8).SetHasSeparator(false).SetExpandedSegments(true).SetAddressLabel(OctalPrefix).ToOptions()
	octal0oPrefixedParams      = new(addrstr.StringOptionsBuilder).SetRadix(8).SetHasSeparator(false).SetExpandedSegments(true).SetAddressLabel(otherOctalPrefix).ToOptions()
	binaryParams               = new(addrstr.StringOptionsBuilder).SetRadix(2).SetHasSeparator(false).SetExpandedSegments(true).ToOptions()
	binaryPrefixedParams       = new(addrstr.StringOptionsBuilder).SetRadix(2).SetHasSeparator(false).SetExpandedSegments(true).SetAddressLabel(BinaryPrefix).ToOptions()
	decimalParams              = new(addrstr.StringOptionsBuilder).SetRadix(10).SetHasSeparator(false).SetExpandedSegments(true).ToOptions()
)

// Format is intentionally the only method with non-pointer receivers.  It is not intended to be called directly, it is intended for use by the fmt package.
// When called by a function in the fmt package, nil values are detected before this method is called, avoiding a panic when calling this method.

// Format implements [fmt.Formatter] interface. It accepts the formats
//   - 'v' for the default address and section format (either the normalized or canonical string),
//   - 's' (string) for the same,
//   - 'b' (binary), 'o' (octal with 0 prefix), 'O' (octal with 0o prefix),
//   - 'd' (decimal), 'x' (lowercase hexadecimal), and
//   - 'X' (uppercase hexadecimal).
//
// Also supported are some of fmt's format flags for integral types.
// Sign control is not supported since addresses and sections are never negative.
// '#' for an alternate format is supported, which adds a leading zero for octal, and for hexadecimal it adds
// a leading "0x" or "0X" for "%#x" and "%#X" respectively.
// Also supported is specification of minimum digits precision, output field width,
// space or zero padding, and '-' for left or right justification.
func (section addressSectionInternal) Format(state fmt.State, verb rune) {
	section.format(state, verb, NoZone, false)
}

func (section *addressSectionInternal) format(state fmt.State, verb rune, zone Zone, useCanonical bool) {
	var str string
	var err error
	var isStringFormat bool

	_, hasPrecision := state.Precision()
	_, hasWidth := state.Width()
	useDefaultStr := !hasPrecision && !hasWidth
	switch verb {
	case 's', 'v', 'q':
		isStringFormat = true
		if useCanonical {
			if zone != NoZone {
				str = section.toAddressSection().ToIPv6().toCanonicalString(zone)
			} else {
				str = section.toCanonicalString()
			}
		} else {
			if zone != NoZone {
				str = section.toAddressSection().ToIPv6().toNormalizedString(zone)
			} else {
				str = section.toNormalizedString()
			}
		}
		if verb == 'q' && useDefaultStr {
			if state.Flag('#') && (zone == NoZone || strconv.CanBackquote(string(zone))) {
				str = "`" + str + "`"
			} else if zone == NoZone {
				str = `"` + str + `"`
			} else {
				str = strconv.Quote(str) // zones should not have special characters, but you cannot be sure
			}
		}
	case 'x':
		useDefaultStr = useDefaultStr && zone == NoZone
		str, err = section.toHexString(useDefaultStr && state.Flag('#'))
	case 'X':
		useDefaultStr = useDefaultStr && zone == NoZone
		if useDefaultStr && state.Flag('#') {
			str, err = section.toLongStringZoned(NoZone, hexPrefixedUppercaseParams)
		} else {

			str, err = section.toLongStringZoned(NoZone, hexUppercaseParams)
		}
	case 'b':
		useDefaultStr = useDefaultStr && zone == NoZone
		str, err = section.toBinaryString(useDefaultStr && state.Flag('#'))
	case 'o':
		useDefaultStr = useDefaultStr && zone == NoZone
		str, err = section.toOctalString(useDefaultStr && state.Flag('#'))
	case 'O':
		useDefaultStr = useDefaultStr && zone == NoZone
		if useDefaultStr {
			str, err = section.toLongOctalStringZoned(NoZone, octal0oPrefixedParams)
		} else {
			str, err = section.toLongOctalStringZoned(NoZone, octalParams)
		}
	case 'd':
		useDefaultStr = useDefaultStr && zone == NoZone
		str, err = section.toDecimalStringZoned(NoZone)
	default:
		// format not supported
		_, _ = fmt.Fprintf(state, "%%!%c(address=%s)", verb, section.toString())
		return
	}
	if err != nil { // could not produce an octal, binary, hex or decimal string, so use string format instead
		isStringFormat = true
		if useCanonical {
			str = section.toCanonicalString()
		} else {
			str = section.toNormalizedString()
		}
	}
	if useDefaultStr {
		_, _ = state.Write([]byte(str))
	} else if isStringFormat {
		section.writeStrFmt(state, verb, str, zone)
	} else {
		section.writeNumberFmt(state, verb, str, zone)
	}
}

func (section addressSectionInternal) writeStrFmt(state fmt.State, verb rune, str string, zone Zone) {
	if precision, hasPrecision := state.Precision(); hasPrecision && len(str) > precision {
		str = str[:precision]
	}
	if verb == 'q' {
		if state.Flag('#') && (zone == NoZone || strconv.CanBackquote(string(zone))) {
			str = "`" + str + "`"
		} else if zone == NoZone {
			str = `"` + str + `"`
		} else {
			str = strconv.Quote(str) // zones should not have special characters, but you cannot be sure
		}
	}
	var leftPaddingCount, rightPaddingCount int
	if width, hasWidth := state.Width(); hasWidth && len(str) < width { // padding required
		paddingCount := width - len(str)
		if state.Flag('-') {
			// right padding with spaces (takes precedence over '0' flag)
			rightPaddingCount = paddingCount
		} else {
			// left padding with spaces
			leftPaddingCount = paddingCount
		}
	}
	// left padding/str/right padding
	writeBytes(state, ' ', leftPaddingCount)
	_, _ = state.Write([]byte(str))
	writeBytes(state, ' ', rightPaddingCount)
}

func (section addressSectionInternal) writeNumberFmt(state fmt.State, verb rune, str string, zone Zone) {
	var prefix string
	if verb == 'O' {
		prefix = otherOctalPrefix // "0o"
	} else if state.Flag('#') {
		switch verb {
		case 'x':
			prefix = HexPrefix
		case 'X':
			prefix = otherHexPrefix
		case 'b':
			prefix = BinaryPrefix
		case 'o':
			prefix = OctalPrefix
		}
	}
	isMulti := section.isMultiple()
	var addrStr, secondStr string
	var separator byte
	if isMulti {
		separatorIndex := len(str) >> 1
		addrStr = str[:separatorIndex]
		separator = str[separatorIndex]
		secondStr = str[separatorIndex+1:]
	} else {
		addrStr = str
	}
	precision, hasPrecision := state.Precision()
	width, hasWidth := state.Width()
	usePrecision := hasPrecision
	if section.hasNoDivisions() {
		usePrecision = false
		prefix = ""
	}
	for {
		var zeroCount, leftPaddingCount, rightPaddingCount int
		if usePrecision {
			if len(addrStr) > precision {
				frontChar := addrStr[0]
				if frontChar == '0' {
					i := 1
					// eliminate leading zeros to match the precision (all the way to nothing)
					for len(addrStr) > precision+i {
						frontChar = addrStr[i]
						if frontChar != '0' {
							break
						}
						i++
					}
					addrStr = addrStr[i:]
				}
			} else if len(addrStr) < precision {
				// expand to match the precision
				zeroCount = precision - len(addrStr)
			}
		}
		length := len(prefix) + zeroCount + len(addrStr)
		zoneRequired := len(zone) > 0
		if zoneRequired {
			length += len(zone) + 1
		}
		if hasWidth && length < width { // padding required
			paddingCount := width - length
			if state.Flag('-') {
				// right padding with spaces (takes precedence over '0' flag)
				rightPaddingCount = paddingCount
			} else if state.Flag('0') && !hasPrecision {
				// left padding with zeros
				zeroCount = paddingCount
			} else {
				// left padding with spaces
				leftPaddingCount = paddingCount
			}
		}

		// left padding/prefix/zeros/str/right padding
		writeBytes(state, ' ', leftPaddingCount)
		writeStr(state, prefix, 1)
		writeBytes(state, '0', zeroCount)
		_, _ = state.Write([]byte(addrStr))
		if zoneRequired {
			_, _ = state.Write([]byte{IPv6ZoneSeparator})
			_, _ = state.Write([]byte(zone))
		}
		writeBytes(state, ' ', rightPaddingCount)

		if !isMulti {
			break
		}
		addrStr = secondStr
		isMulti = false
		_, _ = state.Write([]byte{separator})
	}
}

func writeStr(state fmt.State, str string, count int) {
	if count > 0 && len(str) > 0 {
		bytes := []byte(str)
		for ; count > 0; count-- {
			_, _ = state.Write(bytes)
		}
	}
}

func writeBytes(state fmt.State, b byte, count int) {
	if count > 0 {
		bytes := make([]byte, count)
		for i := range bytes {
			bytes[i] = b
		}
		_, _ = state.Write(bytes)
	}
}

func (section *addressSectionInternal) toCanonicalString() string {
	if sect := section.toIPv4AddressSection(); sect != nil {
		return sect.ToCanonicalString()
	} else if sect := section.toIPv6AddressSection(); sect != nil {
		return sect.ToCanonicalString()
	} else if sect := section.toMACAddressSection(); sect != nil {
		return sect.ToCanonicalString()
	}
	// zero section
	return nilSection()
}

func (section *addressSectionInternal) toNormalizedString() string {
	if sect := section.toIPv4AddressSection(); sect != nil {
		return sect.ToNormalizedString()
	} else if sect := section.toIPv6AddressSection(); sect != nil {
		return sect.ToNormalizedString()
	} else if sect := section.toMACAddressSection(); sect != nil {
		return sect.ToNormalizedString()
	}
	return nilSection()
}

func (section *addressSectionInternal) toNormalizedWildcardString() string {
	if sect := section.toIPv4AddressSection(); sect != nil {
		return sect.ToNormalizedWildcardString()
	} else if sect := section.toIPv6AddressSection(); sect != nil {
		return sect.ToNormalizedWildcardString()
	} else if sect := section.toMACAddressSection(); sect != nil {
		return sect.ToNormalizedWildcardString()
	}
	return nilSection()
}

func (section *addressSectionInternal) toCompressedString() string {
	if sect := section.toIPv4AddressSection(); sect != nil {
		return sect.ToCompressedString()
	} else if sect := section.toIPv6AddressSection(); sect != nil {
		return sect.ToCompressedString()
	} else if sect := section.toMACAddressSection(); sect != nil {
		return sect.ToCompressedString()
	}
	return nilSection()
}

func (section *addressSectionInternal) toDecimalStringZoned(zone Zone) (string, addrerr.IncompatibleAddressError) {
	if isDual, err := section.isDualString(); err != nil {
		return "", err
	} else {
		var largeGrouping *IPAddressLargeDivisionGrouping
		if section.hasNoDivisions() {
			largeGrouping = NewIPAddressLargeDivGrouping(nil)
		} else {
			bytes := section.getBytes()
			prefLen := section.getPrefixLen()
			bitCount := section.GetBitCount()
			var div *IPAddressLargeDivision
			if isDual {
				div = NewIPAddressLargeRangePrefixDivision(bytes, section.getUpperBytes(), prefLen, bitCount, 10)
			} else {
				div = NewIPAddressLargePrefixDivision(bytes, prefLen, bitCount, 10)
			}
			largeGrouping = NewIPAddressLargeDivGrouping([]*IPAddressLargeDivision{div})
		}
		return toNormalizedZonedString(decimalParams, largeGrouping, zone), nil
	}
}

func (section *addressSectionInternal) toHexString(with0xPrefix bool) (string, addrerr.IncompatibleAddressError) {
	cache := section.getStringCache()
	if cache == nil {
		return section.toHexStringZoned(with0xPrefix, NoZone)
	}
	var cacheField **string
	if with0xPrefix {
		cacheField = &cache.hexStringPrefixed
	} else {
		cacheField = &cache.hexString
	}
	return cacheStrErr(cacheField,
		func() (string, addrerr.IncompatibleAddressError) {
			return section.toHexStringZoned(with0xPrefix, NoZone)
		})
}

func (section *addressSectionInternal) toHexStringZoned(with0xPrefix bool, zone Zone) (string, addrerr.IncompatibleAddressError) {
	if with0xPrefix {
		return section.toLongStringZoned(zone, hexPrefixedParams)
	}
	return section.toLongStringZoned(zone, hexParams)
}

func (section *addressSectionInternal) toOctalString(with0Prefix bool) (string, addrerr.IncompatibleAddressError) {
	cache := section.getStringCache()
	if cache == nil {
		return section.toOctalStringZoned(with0Prefix, NoZone)
	}
	var cacheField **string
	if with0Prefix {
		cacheField = &cache.octalStringPrefixed
	} else {
		cacheField = &cache.octalString
	}
	return cacheStrErr(cacheField,
		func() (string, addrerr.IncompatibleAddressError) {
			return section.toOctalStringZoned(with0Prefix, NoZone)
		})
}

func (section *addressSectionInternal) toOctalStringZoned(with0Prefix bool, zone Zone) (string, addrerr.IncompatibleAddressError) {
	var opts addrstr.StringOptions
	if with0Prefix {
		opts = octalPrefixedParams
	} else {
		opts = octalParams
	}
	return section.toLongOctalStringZoned(zone, opts)
}

func (section *addressSectionInternal) toLongOctalStringZoned(zone Zone, opts addrstr.StringOptions) (string, addrerr.IncompatibleAddressError) {
	if isDual, err := section.isDualString(); err != nil {
		return "", err
	} else if isDual {
		lowerDivs, _ := section.getLower().createNewDivisions(3)
		upperDivs, _ := section.getUpper().createNewDivisions(3)
		lowerPart := createInitializedGrouping(lowerDivs, nil)
		upperPart := createInitializedGrouping(upperDivs, nil)
		return toNormalizedStringRange(toParams(opts), lowerPart, upperPart, zone), nil
	}
	divs, _ := section.createNewDivisions(3)
	part := createInitializedGrouping(divs, nil)
	return toParams(opts).toZonedString(part, zone), nil
}

func (section *addressSectionInternal) toBinaryString(with0bPrefix bool) (string, addrerr.IncompatibleAddressError) {
	cache := section.getStringCache()
	if cache == nil {
		return section.toBinaryStringZoned(with0bPrefix, NoZone)
	}
	var cacheField **string
	if with0bPrefix {
		cacheField = &cache.binaryStringPrefixed
	} else {
		cacheField = &cache.binaryString
	}
	return cacheStrErr(cacheField,
		func() (string, addrerr.IncompatibleAddressError) {
			return section.toBinaryStringZoned(with0bPrefix, NoZone)
		})
}

func (section *addressSectionInternal) toBinaryStringZoned(with0bPrefix bool, zone Zone) (string, addrerr.IncompatibleAddressError) {
	if with0bPrefix {
		return section.toLongStringZoned(zone, binaryPrefixedParams)
	}
	return section.toLongStringZoned(zone, binaryParams)
}

func (section *addressSectionInternal) toLongStringZoned(zone Zone, params addrstr.StringOptions) (string, addrerr.IncompatibleAddressError) {
	if isDual, err := section.isDualString(); err != nil {
		return "", err
	} else if isDual {
		sect := section.toAddressSection()
		return toNormalizedStringRange(toParams(params), sect.GetLower(), sect.GetUpper(), zone), nil
	}
	return section.toCustomStringZoned(params, zone), nil
}

func (section *addressSectionInternal) toCustomString(stringOptions addrstr.StringOptions) string {
	return toNormalizedString(stringOptions, section.toAddressSection())
}

func (section *addressSectionInternal) toCustomStringZoned(stringOptions addrstr.StringOptions, zone Zone) string {
	return toNormalizedZonedString(stringOptions, section.toAddressSection(), zone)
}

func (section *addressSectionInternal) isDualString() (bool, addrerr.IncompatibleAddressError) {
	count := section.GetSegmentCount()
	if section.isMultiple() {
		//at this point we know we will return true, but we determine now if we must return addrerr.IncompatibleAddressError
		for i := 0; i < count; i++ {
			division := section.GetSegment(i)
			if division.isMultiple() {
				isLastFull := true
				for j := count - 1; j >= i; j-- {
					division = section.GetSegment(j)
					if division.isMultiple() {
						if !isLastFull {
							return false, &incompatibleAddressError{addressError{key: "ipaddress.error.segmentMismatch"}}
						}
						isLastFull = division.IsFullRange()
					} else {
						isLastFull = false
					}
				}
				return true, nil
			}
		}
	}
	return false, nil
}

// used by iterator() and nonZeroHostIterator() in section types
func (section *addressSectionInternal) sectionIterator(excludeFunc func([]*AddressDivision) bool) Iterator[*AddressSection] {
	useOriginal := !section.isMultiple()
	var original = section.toAddressSection()
	var iterator Iterator[[]*AddressDivision]
	if useOriginal {
		if excludeFunc != nil && excludeFunc(section.getDivisionsInternal()) {
			original = nil // the single-valued iterator starts out empty
		}
	} else {
		iterator = allSegmentsIterator(
			section.GetSegmentCount(),
			nil,
			func(index int) Iterator[*AddressSegment] { return section.GetSegment(index).iterator() },
			excludeFunc)
	}
	return sectIterator(
		useOriginal,
		original,
		false,
		iterator)
}

func (section *addressSectionInternal) prefixIterator(isBlockIterator bool) Iterator[*AddressSection] {
	prefLen := section.prefixLength
	if prefLen == nil {
		return section.sectionIterator(nil)
	}
	prefLength := prefLen.bitCount()
	var useOriginal bool
	if isBlockIterator {
		useOriginal = section.IsSinglePrefixBlock()
	} else {
		useOriginal = section.GetPrefixCount().CmpAbs(bigOneConst()) == 0
	}
	bitsPerSeg := section.GetBitsPerSegment()
	bytesPerSeg := section.GetBytesPerSegment()
	networkSegIndex := getNetworkSegmentIndex(prefLength, bytesPerSeg, bitsPerSeg)
	hostSegIndex := getHostSegmentIndex(prefLength, bytesPerSeg, bitsPerSeg)
	segCount := section.GetSegmentCount()
	var iterator Iterator[[]*AddressDivision]
	if !useOriginal {
		var hostSegIteratorProducer func(index int) Iterator[*AddressSegment]
		if isBlockIterator {
			hostSegIteratorProducer = func(index int) Iterator[*AddressSegment] {
				return section.GetSegment(index).prefixBlockIterator()
			}
		} else {
			hostSegIteratorProducer = func(index int) Iterator[*AddressSegment] {
				return section.GetSegment(index).prefixIterator()
			}
		}
		iterator = segmentsIterator(
			segCount,
			nil, //when no prefix we defer to other iterator, when there is one we use the whole original section in the encompassing iterator and not just the original segments
			func(index int) Iterator[*AddressSegment] { return section.GetSegment(index).iterator() },
			nil,
			networkSegIndex,
			hostSegIndex,
			hostSegIteratorProducer)
	}
	if isBlockIterator {
		return sectIterator(
			useOriginal,
			section.toAddressSection(),
			prefLength < section.GetBitCount(),
			iterator)
	}
	return prefixSectIterator(
		useOriginal,
		section.toAddressSection(),
		iterator)
}

func (section *addressSectionInternal) blockIterator(segmentCount int) Iterator[*AddressSection] {
	if segmentCount < 0 {
		segmentCount = 0
	}
	allSegsCount := section.GetSegmentCount()
	if segmentCount >= allSegsCount {
		return section.sectionIterator(nil)
	}
	useOriginal := !section.isMultipleTo(segmentCount)
	var iterator Iterator[[]*AddressDivision]
	if !useOriginal {
		var hostSegIteratorProducer func(index int) Iterator[*AddressSegment]
		hostSegIteratorProducer = func(index int) Iterator[*AddressSegment] {
			return section.GetSegment(index).identityIterator()
		}
		segIteratorProducer := func(index int) Iterator[*AddressSegment] {
			return section.GetSegment(index).iterator()
		}
		iterator = segmentsIterator(
			allSegsCount,
			nil, //when no prefix we defer to other iterator, when there is one we use the whole original section in the encompassing iterator and not just the original segments
			segIteratorProducer,
			nil,
			segmentCount-1,
			segmentCount,
			hostSegIteratorProducer)
	}
	return sectIterator(
		useOriginal,
		section.toAddressSection(),
		section.isMultipleFrom(segmentCount),
		iterator)
}

func (section *addressSectionInternal) sequentialBlockIterator() Iterator[*AddressSection] {
	return section.blockIterator(section.GetSequentialBlockIndex())
}

// GetSequentialBlockCount provides the count of elements from the sequential block iterator, the minimal number of sequential address sections that comprise this address section.
func (section *addressSectionInternal) GetSequentialBlockCount() *big.Int {
	sequentialSegCount := section.GetSequentialBlockIndex()
	return section.GetPrefixCountLen(BitCount(sequentialSegCount) * section.GetBitsPerSegment())
}

func (section *addressSectionInternal) isMultipleTo(segmentCount int) bool {
	for i := 0; i < segmentCount; i++ {
		if section.GetSegment(i).isMultiple() {
			return true
		}
	}
	return false
}

func (section *addressSectionInternal) isMultipleFrom(segmentCount int) bool {
	segTotal := section.GetSegmentCount()
	for i := segmentCount; i < segTotal; i++ {
		if section.GetSegment(i).isMultiple() {
			return true
		}
	}
	return false
}

func (section *addressSectionInternal) getSubnetSegments( // called by methods to adjust/remove/set prefix length, masking methods, zero host and zero network methods
	startIndex int,
	networkPrefixLength PrefixLen,
	verifyMask bool,
	segProducer func(int) *AddressDivision,
	segmentMaskProducer func(int) SegInt,
) (res *AddressSection, err addrerr.IncompatibleAddressError) {
	networkPrefixLength = checkPrefLen(networkPrefixLength, section.GetBitCount())
	bitsPerSegment := section.GetBitsPerSegment()
	count := section.GetSegmentCount()
	for i := startIndex; i < count; i++ {
		segmentPrefixLength := getSegmentPrefixLength(bitsPerSegment, networkPrefixLength, i)
		seg := segProducer(i)
		//note that the mask can represent a range (for example a CIDR mask),
		//but we use the lowest value (maskSegment.value) in the range when masking (ie we discard the range)
		maskValue := segmentMaskProducer(i)
		origValue, origUpperValue := seg.getSegmentValue(), seg.getUpperSegmentValue()
		value, upperValue := origValue, origUpperValue
		if verifyMask {
			mask64 := uint64(maskValue)
			val64 := uint64(value)
			upperVal64 := uint64(upperValue)
			masker := MaskRange(val64, upperVal64, mask64, seg.GetMaxValue())
			if !masker.IsSequential() {
				err = &incompatibleAddressError{addressError{key: "ipaddress.error.maskMismatch"}}
				return
			}
			value = SegInt(masker.GetMaskedLower(val64, mask64))
			upperValue = SegInt(masker.GetMaskedUpper(upperVal64, mask64))
		} else {
			value &= maskValue
			upperValue &= maskValue
		}
		if !segsSame(segmentPrefixLength, seg.getDivisionPrefixLength(), value, origValue, upperValue, origUpperValue) {
			newSegments := createSegmentArray(count)
			section.copySubDivisions(0, i, newSegments)
			newSegments[i] = createAddressDivision(seg.deriveNewMultiSeg(value, upperValue, segmentPrefixLength))
			for i++; i < count; i++ {
				segmentPrefixLength = getSegmentPrefixLength(bitsPerSegment, networkPrefixLength, i)
				seg = segProducer(i)
				maskValue = segmentMaskProducer(i)
				origValue, origUpperValue = seg.getSegmentValue(), seg.getUpperSegmentValue()
				value, upperValue = origValue, origUpperValue
				if verifyMask {
					mask64 := uint64(maskValue)
					val64 := uint64(value)
					upperVal64 := uint64(upperValue)
					masker := MaskRange(val64, upperVal64, mask64, seg.GetMaxValue())
					if !masker.IsSequential() {
						err = &incompatibleAddressError{addressError{key: "ipaddress.error.maskMismatch"}}
						return
					}
					value = SegInt(masker.GetMaskedLower(val64, mask64))
					upperValue = SegInt(masker.GetMaskedUpper(upperVal64, mask64))
				} else {
					value &= maskValue
					upperValue &= maskValue
				}
				if !segsSame(segmentPrefixLength, seg.getDivisionPrefixLength(), value, origValue, upperValue, origUpperValue) {
					newSegments[i] = createAddressDivision(seg.deriveNewMultiSeg(value, upperValue, segmentPrefixLength))
				} else {
					newSegments[i] = seg
				}
			}
			res = deriveAddressSectionPrefLen(section.toAddressSection(), newSegments, networkPrefixLength)
			return
		}
	}
	res = section.toAddressSection()
	return
}

func (section *addressSectionInternal) toAddressSection() *AddressSection {
	return (*AddressSection)(unsafe.Pointer(section))
}

func (section *addressSectionInternal) toIPAddressSection() *IPAddressSection {
	return section.toAddressSection().ToIP()
}

func (section *addressSectionInternal) toIPv4AddressSection() *IPv4AddressSection {
	return section.toAddressSection().ToIPv4()
}

func (section *addressSectionInternal) toIPv6AddressSection() *IPv6AddressSection {
	return section.toAddressSection().ToIPv6()
}

func (section *addressSectionInternal) toMACAddressSection() *MACAddressSection {
	return section.toAddressSection().ToMAC()
}

//// only needed for godoc / pkgsite

// IsZero returns whether this section matches exactly the value of zero.
func (section *addressSectionInternal) IsZero() bool {
	return section.addressDivisionGroupingInternal.IsZero()
}

// IncludesZero returns whether this section includes the value of zero within its range.
func (section *addressSectionInternal) IncludesZero() bool {
	return section.addressDivisionGroupingInternal.IncludesZero()
}

// IsMax returns whether this section matches exactly the maximum possible value, the value whose bits are all ones.
func (section *addressSectionInternal) IsMax() bool {
	return section.addressDivisionGroupingInternal.IsMax()
}

// IncludesMax returns whether this section includes the max value, the value whose bits are all ones, within its range.
func (section *addressSectionInternal) IncludesMax() bool {
	return section.addressDivisionGroupingInternal.IncludesMax()
}

// IsFullRange returns whether this address item represents all possible values attainable by an address item of this type.
//
// This is true if and only if both IncludesZero and IncludesMax return true.
func (section *addressSectionInternal) IsFullRange() bool {
	return section.addressDivisionGroupingInternal.IsFullRange()
}

// GetSequentialBlockIndex gets the minimal segment index for which all following segments are full-range blocks.
//
// The segment at this index is not a full-range block itself, unless all segments are full-range.
// The segment at this index and all following segments form a sequential range.
// For the full address section to be sequential, the preceding segments must be single-valued.
func (section *addressSectionInternal) GetSequentialBlockIndex() int {
	return section.addressDivisionGroupingInternal.GetSequentialBlockIndex()
}

// GetPrefixLen returns the prefix length, or nil if there is no prefix length.
//
// A prefix length indicates the number of bits in the initial part of the address item that comprises the prefix.
//
// A prefix is a part of the address item that is not specific to that address but common amongst a group of such items, such as a CIDR prefix block subnet.
func (section *addressSectionInternal) GetPrefixLen() PrefixLen {
	return section.addressDivisionGroupingInternal.GetPrefixLen()
}

// ContainsPrefixBlock returns whether the values of this item contains the block of values for the given prefix length.
//
// Unlike ContainsSinglePrefixBlock, whether there are multiple prefix values in this item for the given prefix length makes no difference.
//
// Use GetMinPrefixLenForBlock to determine the smallest prefix length for which this method returns true.
func (section *addressSectionInternal) ContainsPrefixBlock(prefixLen BitCount) bool {
	prefixLen = checkSubnet(section, prefixLen)
	divCount := section.GetSegmentCount()
	bitsPerSegment := section.GetBitsPerSegment()
	i := getHostSegmentIndex(prefixLen, section.GetBytesPerSegment(), bitsPerSegment)
	if i < divCount {
		div := section.GetSegment(i)
		segmentPrefixLength := getPrefixedSegmentPrefixLength(bitsPerSegment, prefixLen, i)
		if !div.ContainsPrefixBlock(segmentPrefixLength.bitCount()) {
			return false
		}
		for i++; i < divCount; i++ {
			div = section.GetSegment(i)
			if !div.IsFullRange() {
				return false
			}
		}
	}
	return true
}

// ContainsSinglePrefixBlock returns whether the values of this grouping contains a single prefix block for the given prefix length.
//
// This means there is only one prefix of the given length in this item, and this item contains the prefix block for that given prefix.
//
// Use GetPrefixLenForSingleBlock to determine whether there is a prefix length for which this method returns true.
func (section *addressSectionInternal) ContainsSinglePrefixBlock(prefixLen BitCount) bool {
	return section.addressDivisionGroupingInternal.ContainsSinglePrefixBlock(prefixLen)
}

// IsPrefixBlock returns whether this address segment series has a prefix length and includes the block associated with its prefix length.
// If the prefix length matches the bit count, this returns true.
//
// This is different from ContainsPrefixBlock in that this method returns
// false if the series has no prefix length or a prefix length that differs from a prefix length for which ContainsPrefixBlock returns true.
func (section *addressSectionInternal) IsPrefixBlock() bool {
	prefLen := section.getPrefixLen()
	return prefLen != nil && section.ContainsPrefixBlock(prefLen.bitCount())
}

// IsSinglePrefixBlock returns whether the range matches the block of values for a single prefix identified by the prefix length of this address.
// This is similar to IsPrefixBlock except that it returns false when the subnet has multiple prefixes.
//
// What distinguishes this method from ContainsSinglePrefixBlock is that this method returns
// false if the series does not have a prefix length assigned to it,
// or a prefix length that differs from a prefix length for which ContainsSinglePrefixBlock returns true.
//
// It is similar to IsPrefixBlock but returns false when there are multiple prefixes.
func (section *addressSectionInternal) IsSinglePrefixBlock() bool {
	return section.addressDivisionGroupingInternal.IsSinglePrefixBlock()
}

// GetMinPrefixLenForBlock returns the smallest prefix length such that this section includes the block of all values for that prefix length.
//
// If the entire range can be described this way, then this method returns the same value as GetPrefixLenForSingleBlock.
//
// There may be a single prefix, or multiple possible prefix values in this item for the returned prefix length.
// Use GetPrefixLenForSingleBlock to avoid the case of multiple prefix values.
//
// If this section represents a single value, this returns the bit count.
func (section *addressSectionInternal) GetMinPrefixLenForBlock() BitCount {
	return section.addressDivisionGroupingInternal.GetMinPrefixLenForBlock()
}

// GetPrefixLenForSingleBlock returns a prefix length for which the range of this address section matches the block of addresses for that prefix.
//
// If no such prefix exists, GetPrefixLenForSingleBlock returns nil.
//
// If this address section represents a single value, returns the bit length.
func (section *addressSectionInternal) GetPrefixLenForSingleBlock() PrefixLen {
	return section.addressDivisionGroupingInternal.GetPrefixLenForSingleBlock()
}

// GetValue returns the lowest individual address section in this address section as an integer value.
func (section *addressSectionInternal) GetValue() *big.Int {
	return section.addressDivisionGroupingInternal.GetValue()
}

// GetUpperValue returns the highest individual address section in this address section as an integer value.
func (section *addressSectionInternal) GetUpperValue() *big.Int {
	return section.addressDivisionGroupingInternal.GetUpperValue()
}

// Bytes returns the lowest individual address section in this address section as a byte slice.
func (section *addressSectionInternal) Bytes() []byte {
	return section.addressDivisionGroupingInternal.Bytes()
}

// UpperBytes returns the highest individual address section in this address section as a byte slice.
func (section *addressSectionInternal) UpperBytes() []byte {
	return section.addressDivisionGroupingInternal.UpperBytes()
}

// CopyBytes copies the value of the lowest individual address section in the section into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (section *addressSectionInternal) CopyBytes(bytes []byte) []byte {
	return section.addressDivisionGroupingInternal.CopyBytes(bytes)
}

// CopyUpperBytes copies the value of the highest individual address section in the section into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (section *addressSectionInternal) CopyUpperBytes(bytes []byte) []byte {
	return section.addressDivisionGroupingInternal.CopyUpperBytes(bytes)
}

// IsSequential returns  whether the section represents a range of values that are sequential.
//
// Generally, this means that any segment covering a range of values must be followed by segment that are full range, covering all values.
func (section *addressSectionInternal) IsSequential() bool {
	return section.addressDivisionGroupingInternal.IsSequential()
}

// GetLeadingBitCount returns the number of consecutive leading one or zero bits.
// If ones is true, returns the number of consecutive leading one bits.
// Otherwise, returns the number of consecutive leading zero bits.
//
// This method applies to the lower value of the range if this section represents multiple values.
func (section *addressSectionInternal) GetLeadingBitCount(ones bool) BitCount {
	count := section.GetSegmentCount()
	if count == 0 {
		return 0
	}
	var front SegInt
	if ones {
		front = section.GetSegment(0).GetMaxValue()
	}
	var prefixLen BitCount
	for i := 0; i < count; i++ {
		seg := section.GetSegment(i)
		value := seg.getSegmentValue()
		if value != front {
			return prefixLen + seg.GetLeadingBitCount(ones)
		}
		prefixLen += seg.getBitCount()
	}
	return prefixLen
}

// GetTrailingBitCount returns the number of consecutive trailing one or zero bits.
// If ones is true, returns the number of consecutive trailing zero bits.
// Otherwise, returns the number of consecutive trailing one bits.
//
// This method applies to the lower value of the range if this section represents multiple values.
func (section *addressSectionInternal) GetTrailingBitCount(ones bool) BitCount {
	count := section.GetSegmentCount()
	if count == 0 {
		return 0
	}
	var back SegInt
	if ones {
		back = section.GetSegment(0).GetMaxValue()
	}
	var bitLen BitCount
	for i := count - 1; i >= 0; i-- {
		seg := section.GetSegment(i)
		value := seg.getSegmentValue()
		if value != back {
			return bitLen + seg.GetTrailingBitCount(ones)
		}
		bitLen += seg.getBitCount()
	}
	return bitLen
}

//// end needed for godoc / pkgsite

// An AddressSection is section of an address, containing a certain number of consecutive segments.
//
// It is a series of individual address segments.  Each segment has equal bit-length.  Each address is backed by an address section that contains all the segments of the address.
//
// AddressSection instances are immutable.  This also makes them concurrency-safe.
//
// Most operations that can be performed on Address instances can also be performed on AddressSection instances and vice-versa.
type AddressSection struct {
	addressSectionInternal
}

// Contains returns whether this is same type and version as the given address section and whether it contains all values in the given section.
//
// Sections must also have the same number of segments to be comparable, otherwise false is returned.
func (section *AddressSection) Contains(other AddressSectionType) bool {
	if section == nil {
		return other == nil || other.ToSectionBase() == nil
	}
	return section.contains(other)
}

// Overlaps returns whether this is same type and version as the given address section and whether it overlaps the given section, both sections containing at least one individual section in common.
//
// Sections must also have the same number of segments to be comparable, otherwise false is returned.
func (section *AddressSection) Overlaps(other AddressSectionType) bool {
	if section == nil {
		return other == nil || other.ToSectionBase() == nil
	}
	return section.overlaps(other)
}

// Equal returns whether the given address section is equal to this address section.
// Two address sections are equal if they represent the same set of sections.
// They must match:
//   - type/version (IPv4, IPv6, MAC, etc)
//   - segment counts
//   - bits per segment
//   - segment value ranges
//
// Prefix lengths are ignored.
func (section *AddressSection) Equal(other AddressSectionType) bool {
	if section == nil {
		return other == nil || other.ToSectionBase() == nil
	}
	return section.equal(other)
}

// Compare returns a negative integer, zero, or a positive integer if this address section is less than, equal, or greater than the given item.
// Any address item is comparable to any other.  All address items use CountComparator to compare.
func (section *AddressSection) Compare(item AddressItem) int {
	return CountComparator.Compare(section, item)
}

// CompareSize compares the counts of two address sections, the number of individual sections represented.
//
// Rather than calculating counts with GetCount, there can be more efficient ways of determining whether one section represents more individual address sections than another.
//
// CompareSize returns a positive integer if this address section has a larger count than the one given, zero if they are the same, or a negative integer if the other has a larger count.
func (section *AddressSection) CompareSize(other AddressItem) int {
	if section == nil {
		if isNilItem(other) {
			return 0
		}
		// we have size 0, other has size >= 1
		return -1
	}
	return section.compareSize(other)
}

// GetCount returns the count of possible distinct values for this item.
// If not representing multiple values, the count is 1,
// unless this is a division grouping with no divisions, or an address section with no segments, in which case it is 0.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (section *AddressSection) GetCount() *big.Int {
	if section == nil {
		return bigZero()
	} else if sect := section.ToIPv4(); sect != nil {
		return sect.GetCount()
	} else if sect := section.ToIPv6(); sect != nil {
		return sect.GetCount()
	} else if sect := section.ToMAC(); sect != nil {
		return sect.GetCount()
	}
	return section.addressDivisionGroupingBase.getCount()
}

// IsMultiple returns whether this section represents multiple values.
func (section *AddressSection) IsMultiple() bool {
	return section != nil && section.isMultiple()
}

// GetPrefixCount returns the number of distinct prefix values in this item.
//
// The prefix length is given by GetPrefixLen.
//
// If this has a non-nil prefix length, returns the number of distinct prefix values.
//
// If this has a nil prefix length, returns the same value as GetCount.
func (section *AddressSection) GetPrefixCount() *big.Int {
	if sect := section.ToIPv4(); sect != nil {
		return sect.GetPrefixCount()
	} else if sect := section.ToIPv6(); sect != nil {
		return sect.GetPrefixCount()
	} else if sect := section.ToMAC(); sect != nil {
		return sect.GetPrefixCount()
	}
	return section.addressDivisionGroupingBase.GetPrefixCount()
}

// GetPrefixCountLen returns the number of distinct prefix values in this item for the given prefix length.
func (section *AddressSection) GetPrefixCountLen(prefixLen BitCount) *big.Int {
	if sect := section.ToIPv4(); sect != nil {
		return sect.GetPrefixCountLen(prefixLen)
	} else if sect := section.ToIPv6(); sect != nil {
		return sect.GetPrefixCountLen(prefixLen)
	} else if sect := section.ToMAC(); sect != nil {
		return sect.GetPrefixCountLen(prefixLen)
	}
	return section.addressDivisionGroupingBase.GetPrefixCountLen(prefixLen)
}

// GetBlockCount returns the count of distinct values in the given number of initial (more significant) segments.
func (section *AddressSection) GetBlockCount(segments int) *big.Int {
	if sect := section.ToIPv4(); sect != nil {
		return sect.GetBlockCount(segments)
	} else if sect := section.ToIPv6(); sect != nil {
		return sect.GetBlockCount(segments)
	} else if sect := section.ToMAC(); sect != nil {
		return sect.GetBlockCount(segments)
	}
	return section.addressDivisionGroupingBase.GetBlockCount(segments)
}

// GetTrailingSection gets the subsection from the series starting from the given index.
// The first segment is at index 0.
func (section *AddressSection) GetTrailingSection(index int) *AddressSection {
	return section.getSubSection(index, section.GetSegmentCount())
}

// GetSubSection gets the subsection from the series starting from the given index and ending just before the give endIndex.
// The first segment is at index 0.
func (section *AddressSection) GetSubSection(index, endIndex int) *AddressSection {
	return section.getSubSection(index, endIndex)
}

// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
func (section *AddressSection) CopySubSegments(start, end int, segs []*AddressSegment) (count int) {
	start, end, targetStart := adjust1To1StartIndices(start, end, section.GetDivisionCount(), len(segs))
	segs = segs[targetStart:]
	return section.forEachSubDivision(start, end, func(index int, div *AddressDivision) {
		segs[index] = div.ToSegmentBase()
	}, len(segs))
}

// CopySegments copies the existing segments into the given slice,
// as much as can be fit into the slice, returning the number of segments copied.
func (section *AddressSection) CopySegments(segs []*AddressSegment) (count int) {
	return section.ForEachSegment(func(index int, seg *AddressSegment) (stop bool) {
		if stop = index >= len(segs); !stop {
			segs[index] = seg
		}
		return
	})
}

// GetSegments returns a slice with the address segments.  The returned slice is not backed by the same array as this section.
func (section *AddressSection) GetSegments() (res []*AddressSegment) {
	res = make([]*AddressSegment, section.GetSegmentCount())
	section.CopySegments(res)
	return
}

// GetLower returns the section in the range with the lowest numeric value,
// which will be the same section if it represents a single value.
// For example, for "1.2-3.4.5-6", the section "1.2.4.5" is returned.
func (section *AddressSection) GetLower() *AddressSection {
	return section.getLower()
}

// GetUpper returns the section in the range with the highest numeric value,
// which will be the same section if it represents a single value.
// For example, for "1.2-3.4.5-6", the section "1.3.4.6" is returned.
func (section *AddressSection) GetUpper() *AddressSection {
	return section.getUpper()
}

// IsPrefixed returns whether this section has an associated prefix length.
func (section *AddressSection) IsPrefixed() bool {
	return section != nil && section.isPrefixed()
}

// ToPrefixBlock returns the section with the same prefix as this section while the remaining bits span all values.
// The returned section will be the block of all sections with the same prefix.
//
// If this section has no prefix, this section is returned.
func (section *AddressSection) ToPrefixBlock() *AddressSection {
	return section.toPrefixBlock()
}

// ToPrefixBlockLen returns the section with the same prefix of the given length as this section while the remaining bits span all values.
// The returned section will be the block of all sections with the same prefix.
func (section *AddressSection) ToPrefixBlockLen(prefLen BitCount) *AddressSection {
	return section.toPrefixBlockLen(prefLen)
}

// WithoutPrefixLen provides the same address section but with no prefix length.  The values remain unchanged.
func (section *AddressSection) WithoutPrefixLen() *AddressSection {
	if !section.IsPrefixed() {
		return section
	}
	return section.withoutPrefixLen()
}

// SetPrefixLen sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the address section.
// The provided prefix length will be adjusted to these boundaries if necessary.
func (section *AddressSection) SetPrefixLen(prefixLen BitCount) *AddressSection {
	return section.setPrefixLen(prefixLen)
}

// SetPrefixLenZeroed sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the address section.
// The provided prefix length will be adjusted to these boundaries if necessary.
//
// If this address section has a prefix length, and the prefix length is increased when setting the new prefix length, the bits moved within the prefix become zero.
// If this address section has a prefix length, and the prefix length is decreased when setting the new prefix length, the bits moved outside the prefix become zero.
//
// In other words, bits that move from one side of the prefix length to the other (bits moved into the prefix or outside the prefix) are zeroed.
//
// If the result cannot be zeroed because zeroing out bits results in a non-contiguous segment, an error is returned.
func (section *AddressSection) SetPrefixLenZeroed(prefixLen BitCount) (*AddressSection, addrerr.IncompatibleAddressError) {
	return section.setPrefixLenZeroed(prefixLen)
}

// AdjustPrefixLen increases or decreases the prefix length by the given increment.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the address section.
//
// If this address section has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
func (section *AddressSection) AdjustPrefixLen(prefixLen BitCount) *AddressSection {
	return section.adjustPrefixLen(prefixLen).ToSectionBase()
}

// AdjustPrefixLenZeroed increases or decreases the prefix length by the given increment while zeroing out the bits that have moved into or outside the prefix.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the address section.
//
// If this address section has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
//
// When prefix length is increased, the bits moved within the prefix become zero.
// When a prefix length is decreased, the bits moved outside the prefix become zero.
//
// If the result cannot be zeroed because zeroing out bits results in a non-contiguous segment, an error is returned.
func (section *AddressSection) AdjustPrefixLenZeroed(prefixLen BitCount) (*AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.adjustPrefixLenZeroed(prefixLen)
	return res.ToSectionBase(), err
}

// AssignPrefixForSingleBlock returns the equivalent prefix block that matches exactly the range of values in this address section.
// The returned block will have an assigned prefix length indicating the prefix length for the block.
//
// There may be no such address section - it is required that the range of values match the range of a prefix block.
// If there is no such address section, then nil is returned.
func (section *AddressSection) AssignPrefixForSingleBlock() *AddressSection {
	return section.assignPrefixForSingleBlock()
}

// AssignMinPrefixForBlock returns an equivalent address section, assigned the smallest prefix length possible,
// such that the prefix block for that prefix length is in this address section.
//
// In other words, this method assigns a prefix length to this address section matching the largest prefix block in this address section.
func (section *AddressSection) AssignMinPrefixForBlock() *AddressSection {
	return section.assignMinPrefixForBlock()
}

// ToBlock creates a new block of address sections by changing the segment at the given index to have the given lower and upper value,
// and changing the following segments to be full-range.
func (section *AddressSection) ToBlock(segmentIndex int, lower, upper SegInt) *AddressSection {
	return section.toBlock(segmentIndex, lower, upper)
}

// IsAdaptiveZero returns true if the division grouping was originally created as an implicitly zero-valued section or grouping (e.g. IPv4AddressSection{}),
// meaning it was not constructed using a constructor function.
// Such a grouping, which has no divisions or segments, is convertible to an implicitly zero-valued grouping of any type or version, whether IPv6, IPv4, MAC, or other.
// In other words, when a section or grouping is the zero-value, then it is equivalent and convertible to the zero value of any other section or grouping type.
func (section *AddressSection) IsAdaptiveZero() bool {
	return section != nil && section.matchesZeroGrouping()
}

// IsIP returns true if this address section originated as an IPv4 or IPv6 section, or a zero-length IP section.  If so, use ToIP to convert back to the IP-specific type.
func (section *AddressSection) IsIP() bool {
	return section != nil && section.matchesIPSectionType()
}

// IsIPv4 returns true if this address section originated as an IPv4 section.  If so, use ToIPv4 to convert back to the IPv4-specific type.
func (section *AddressSection) IsIPv4() bool {
	return section != nil && section.matchesIPv4SectionType()
}

// IsIPv6 returns true if this address section originated as an IPv6 section.  If so, use ToIPv6 to convert back to the IPv6-specific type.
func (section *AddressSection) IsIPv6() bool {
	return section != nil && section.matchesIPv6SectionType()
}

// IsMAC returns true if this address section originated as a MAC section.  If so, use ToMAC to convert back to the MAC-specific type.
func (section *AddressSection) IsMAC() bool {
	return section != nil && section.matchesMACSectionType()
}

// ToDivGrouping converts to an AddressDivisionGrouping, a polymorphic type usable with all address sections and division groupings.
// Afterwards, you can convert back with ToSectionBase.
//
// ToDivGrouping can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *AddressSection) ToDivGrouping() *AddressDivisionGrouping {
	return (*AddressDivisionGrouping)(unsafe.Pointer(section))
}

// ToIP converts to an IPAddressSection if this address section originated as an IPv4 or IPv6 section, or an implicitly zero-valued IP section.
// If not, ToIP returns nil.
//
// ToIP can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *AddressSection) ToIP() *IPAddressSection {
	if section.IsIP() {
		return (*IPAddressSection)(unsafe.Pointer(section))
	}
	return nil
}

// ToIPv6 converts to an IPv6AddressSection if this section originated as an IPv6 section.
// If not, ToIPv6 returns nil.
//
// ToIPv6 can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *AddressSection) ToIPv6() *IPv6AddressSection {
	if section.IsIPv6() {
		return (*IPv6AddressSection)(unsafe.Pointer(section))
	}
	return nil
}

// ToIPv4 converts to an IPv4AddressSection if this section originated as an IPv4 section.
// If not, ToIPv4 returns nil.
//
// ToIPv4 can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *AddressSection) ToIPv4() *IPv4AddressSection {
	if section.IsIPv4() {
		return (*IPv4AddressSection)(unsafe.Pointer(section))
	}
	return nil
}

// ToMAC converts to a MACAddressSection if this section originated as a MAC section.
// If not, ToMAC returns nil.
//
// ToMAC can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *AddressSection) ToMAC() *MACAddressSection {
	if section.IsMAC() {
		return (*MACAddressSection)(section)
	}
	return nil
}

// ToSectionBase is an identity method.
//
// ToSectionBase can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *AddressSection) ToSectionBase() *AddressSection {
	return section
}

// Wrap wraps this address section, returning a WrappedAddressSection, an implementation of ExtendedSegmentSeries,
// which can be used to write code that works with both addresses and address sections.
func (section *AddressSection) Wrap() WrappedAddressSection {
	return wrapSection(section)
}

// Iterator provides an iterator to iterate through the individual address sections of this address section.
//
// When iterating, the prefix length is preserved.  Remove it using WithoutPrefixLen prior to iterating if you wish to drop it from all individual address sections.
//
// Call IsMultiple to determine if this instance represents multiple address sections, or GetCount for the count.
func (section *AddressSection) Iterator() Iterator[*AddressSection] {
	if section == nil {
		return nilSectIterator()
	}
	return section.sectionIterator(nil)
}

// PrefixIterator provides an iterator to iterate through the individual prefixes of this address section,
// each iterated element spanning the range of values for its prefix.
//
// It is similar to the prefix block iterator, except for possibly the first and last iterated elements, which might not be prefix blocks,
// instead constraining themselves to values from this address section.
//
// If the series has no prefix length, then this is equivalent to Iterator.
func (section *AddressSection) PrefixIterator() Iterator[*AddressSection] {
	return section.prefixIterator(false)
}

// PrefixBlockIterator provides an iterator to iterate through the individual prefix blocks, one for each prefix of this address section.
// Each iterated address section will be a prefix block with the same prefix length as this address section.
//
// If this address section has no prefix length, then this is equivalent to Iterator.
func (section *AddressSection) PrefixBlockIterator() Iterator[*AddressSection] {
	return section.prefixIterator(true)
}

// IncrementBoundary returns the item that is the given increment from the range boundaries of this item.
//
// If the given increment is positive, adds the value to the highest (GetUpper) in the range to produce a new item.
// If the given increment is negative, adds the value to the lowest (GetLower) in the range to produce a new item.
// If the increment is zero, returns this.
//
// If this represents just a single value, this item is simply incremented by the given increment value, positive or negative.
//
// On overflow or underflow, IncrementBoundary returns nil.
func (section *AddressSection) IncrementBoundary(increment int64) *AddressSection {
	return section.incrementBoundary(increment)
}

// Increment returns the item that is the given increment upwards into the range,
// with the increment of 0 returning the first in the range.
//
// If the increment i matches or exceeds the range count c, then i - c + 1
// is added to the upper item of the range.
// An increment matching the count gives you the item just above the highest in the range.
//
// If the increment is negative, it is added to the lowest of the range.
// To get the item just below the lowest of the range, use the increment -1.
//
// If this represents just a single value, the item is simply incremented by the given increment, positive or negative.
//
// If this item represents multiple values, a positive increment i is equivalent i + 1 values from the iterator and beyond.
// For instance, a increment of 0 is the first value from the iterator, an increment of 1 is the second value from the iterator, and so on.
// An increment of a negative value added to the count is equivalent to the same number of iterator values preceding the last value of the iterator.
// For instance, an increment of count - 1 is the last value from the iterator, an increment of count - 2 is the second last value, and so on.
//
// On overflow or underflow, Increment returns nil.
func (section *AddressSection) Increment(increment int64) *AddressSection {
	return section.increment(increment)
}

// Enumerate indicates where an individual address section sits relative to the address section range ordering.
//
// Determines how many address section elements of a range precede the given address section element, if the address section is in the range.
// If above the range, it is the distance to the upper boundary added to the range count less one, and if below the range, the distance to the lower boundary.
//
// In other words, if the given address section is not in the range but above it, returns the number of address sections preceding the address from the upper range boundary,
// added to one less than the total number of range address sections.  If the given address section is not in the subnet but below it, returns the number of address sections following the address section to the lower subnet boundary.
//
// If the argument is not in the range, but neither above nor below the range, then nil is returned.
//
// Enumerate returns nil when the argument is multi-valued. The argument must be an individual address section.
//
// When this is also an individual address section, the returned value is the distance (difference) between the two address section values.
//
// If the given address section does not have the same version or type, then nil is returned.
//
// Sections must also have the same number of segments to be comparable, otherwise nil is returned.
func (section *AddressSection) Enumerate(other AddressSectionType) *big.Int {
	return section.enumerate(other)
}

// ReverseBits returns a new section with the bits reversed.  Any prefix length is dropped.
//
// If the bits within a single segment cannot be reversed because the segment represents a range,
// and reversing the segment values results in a range that is not contiguous, this returns an error.
//
// In practice this means that to be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
//
// If perByte is true, the bits are reversed within each byte, otherwise all the bits are reversed.
func (section *AddressSection) ReverseBits(perByte bool) (*AddressSection, addrerr.IncompatibleAddressError) {
	return section.reverseBits(perByte)
}

// ReverseBytes returns a new section with the bytes reversed.  Any prefix length is dropped.
//
// If each segment is more than 1 byte long, and the bytes within a single segment cannot be reversed because the segment represents a range,
// and reversing the segment values results in a range that is not contiguous, then this returns an error.
//
// In practice this means that to be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
func (section *AddressSection) ReverseBytes() (*AddressSection, addrerr.IncompatibleAddressError) {
	return section.reverseBytes(false)
}

// ReverseSegments returns a new section with the segments reversed.
func (section *AddressSection) ReverseSegments() *AddressSection {
	if section.GetSegmentCount() <= 1 {
		if section.IsPrefixed() {
			return section.WithoutPrefixLen()
		}
		return section
	}
	res, _ := section.reverseSegments(
		func(i int) (*AddressSegment, addrerr.IncompatibleAddressError) {
			return section.GetSegment(i).withoutPrefixLen(), nil
		},
	)
	return res
}

// String implements the [fmt.Stringer] interface, returning the normalized string provided by ToNormalizedString, or "<nil>" if the receiver is a nil pointer.
func (section *AddressSection) String() string {
	if section == nil {
		return nilString()
	}
	return section.toString()
}

// ToCanonicalString produces a canonical string for the address section.
//
// For IPv4, dotted octet format, also known as dotted decimal format, is used.
// https://datatracker.ietf.org/doc/html/draft-main-ipaddr-text-rep-00#section-2.1
//
// For IPv6, RFC 5952 describes canonical string representation.
// https://en.wikipedia.org/wiki/IPv6_address#Representation
// http://tools.ietf.org/html/rfc5952
//
// For MAC, it uses the canonical standardized IEEE 802 MAC address representation of xx-xx-xx-xx-xx-xx.  An example is "01-23-45-67-89-ab".
// For range segments, '|' is used: "11-22-33|44-55-66".
func (section *AddressSection) ToCanonicalString() string {
	if section == nil {
		return nilString()
	}
	return section.toCanonicalString()
}

// ToNormalizedString produces a normalized string for the address section.
//
// For IPv4, it is the same as the canonical string.
//
// For IPv6, it differs from the canonical string.  Zero-segments are not compressed.
//
// For MAC, it differs from the canonical string.  It uses the most common representation of MAC addresses: "xx:xx:xx:xx:xx:xx".  An example is "01:23:45:67:89:ab".
// For range segments, '-' is used: "11:22:33-44:55:66".
func (section *AddressSection) ToNormalizedString() string {
	if section == nil {
		return nilString()
	}
	return section.toNormalizedString()
}

// ToNormalizedWildcardString produces a string similar to the normalized string but for IP address sections it avoids the CIDR prefix length.
// Multiple-valued segments will be shown with wildcards and ranges (denoted by '*' and '-') instead of using the CIDR prefix notation.
func (section *AddressSection) ToNormalizedWildcardString() string {
	if section == nil {
		return nilString()
	}
	return section.toNormalizedWildcardString()
}

// ToCompressedString produces a short representation of this address section while remaining within the confines of standard representation(s) of the address.
//
// For IPv4, it is the same as the canonical string.
//
// For IPv6, it differs from the canonical string.  It compresses the maximum number of zeros and/or host segments with the IPv6 compression notation '::'.
//
// For MAC, it differs from the canonical string.  It produces a shorter string for the address that has no leading zeros.
func (section *AddressSection) ToCompressedString() string {
	if section == nil {
		return nilString()
	}
	return section.toCompressedString()
}

// ToHexString writes this address section as a single hexadecimal value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0x" prefix.
//
// If a multiple-valued section cannot be written as a single prefix block or a range of two values, an error is returned.
func (section *AddressSection) ToHexString(with0xPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toHexString(with0xPrefix)
}

// ToOctalString writes this address section as a single octal value (possibly two values if a range),
// the number of digits according to the bit count, with or without a preceding "0" prefix.
//
// If a multiple-valued section cannot be written as a single prefix block or a range of two values, an error is returned.
func (section *AddressSection) ToOctalString(with0Prefix bool) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toOctalString(with0Prefix)
}

// ToBinaryString writes this address section as a single binary value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0b" prefix.
//
// If a multiple-valued section cannot be written as a single prefix block or a range of two values, an error is returned.
func (section *AddressSection) ToBinaryString(with0bPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toBinaryString(with0bPrefix)
}

// ToCustomString creates a customized string from this address section according to the given string option parameters.
func (section *AddressSection) ToCustomString(stringOptions addrstr.StringOptions) string {
	if section == nil {
		return nilString()
	}
	return section.toCustomString(stringOptions)
}

// GetSegmentStrings returns a slice with the string for each segment being the string that is normalized with wildcards.
func (section *AddressSection) GetSegmentStrings() []string {
	if section == nil {
		return nil
	}
	return section.getSegmentStrings()
}

func seriesValsSame(one, two AddressSegmentSeries) bool {
	if one == two {
		return true
	}
	count := one.GetDivisionCount()
	if count != two.GetDivisionCount() {
		panic(two)
	}
	for i := count - 1; i >= 0; i-- { // reverse order since less significant segments more likely to differ
		oneSeg := one.GetGenericSegment(i)
		twoSeg := two.GetGenericSegment(i)
		if !segValsSame(oneSeg.GetSegmentValue(), twoSeg.GetSegmentValue(),
			oneSeg.GetUpperSegmentValue(), twoSeg.GetUpperSegmentValue()) {
			return false
		}
	}
	return true
}

func toSegments(
	bytes []byte,
	segmentCount int,
	bytesPerSegment int,
	bitsPerSegment BitCount,
	creator addressSegmentCreator,
	assignedPrefixLength PrefixLen) (segments []*AddressDivision, err addrerr.AddressValueError) {

	segments = createSegmentArray(segmentCount)
	byteIndex, segmentIndex := len(bytes), segmentCount-1
	for ; segmentIndex >= 0; segmentIndex-- {
		var value SegInt
		k := byteIndex - bytesPerSegment
		if k < 0 {
			k = 0
		}
		for j := k; j < byteIndex; j++ {
			byteValue := bytes[j]
			value <<= 8
			value |= SegInt(byteValue)
		}
		byteIndex = k
		segmentPrefixLength := getSegmentPrefixLength(bitsPerSegment, assignedPrefixLength, segmentIndex)
		seg := creator.createSegment(value, value, segmentPrefixLength)
		segments[segmentIndex] = seg
	}
	// any remaining bytes should be zero
	for byteIndex--; byteIndex >= 0; byteIndex-- {
		if bytes[byteIndex] != 0 {
			err = &addressValueError{
				addressError: addressError{key: "ipaddress.error.exceeds.size"},
				val:          int(bytes[byteIndex]),
			}
			break
		}
	}
	return
}
