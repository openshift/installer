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
	"math/big"
	"unsafe"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstr"
)

func createIPv4Section(segments []*AddressDivision) *IPv4AddressSection {
	return &IPv4AddressSection{
		ipAddressSectionInternal{
			addressSectionInternal{
				addressDivisionGroupingInternal{
					addressDivisionGroupingBase: addressDivisionGroupingBase{
						divisions: standardDivArray(segments),
						addrType:  ipv4Type,
						cache: &valueCache{
							stringCache: stringCache{
								ipStringCache:   &ipStringCache{},
								ipv4StringCache: &ipv4StringCache{},
							},
						},
					},
				},
			},
		},
	}
}

func newIPv4SectionParsed(segments []*AddressDivision, isMultiple bool) (res *IPv4AddressSection) {
	res = createIPv4Section(segments)
	res.isMult = isMultiple
	return
}

// this one is used by that parsing code when there are prefix lengths to be applied
func newPrefixedIPv4SectionParsed(segments []*AddressDivision, isMultiple bool, prefixLength PrefixLen, singleOnly bool) (res *IPv4AddressSection) {
	res = createIPv4Section(segments)
	res.isMult = isMultiple
	if prefixLength != nil {
		assignPrefix(prefixLength, segments, res.ToIP(), singleOnly, false, BitCount(len(segments)<<ipv4BitsToSegmentBitshift))
	}
	return
}

// NewIPv4Section constructs an IPv4 address or subnet section from the given segments.
func NewIPv4Section(segments []*IPv4AddressSegment) *IPv4AddressSection {
	return createIPv4SectionFromSegs(segments, nil)
}

// NewIPv4PrefixedSection constructs an IPv4 address or subnet section from the given segments and prefix length.
func NewIPv4PrefixedSection(segments []*IPv4AddressSegment, prefixLen PrefixLen) *IPv4AddressSection {
	return createIPv4SectionFromSegs(segments, prefixLen)
}

func createIPv4SectionFromSegs(orig []*IPv4AddressSegment, prefLen PrefixLen) (result *IPv4AddressSection) {
	divs, newPref, isMultiple := createDivisionsFromSegs(
		func(index int) *IPAddressSegment {
			return orig[index].ToIP()
		},
		len(orig),
		ipv4BitsToSegmentBitshift,
		IPv4BitsPerSegment,
		IPv4BytesPerSegment,
		IPv4MaxValuePerSegment,
		zeroIPv4Seg.ToIP(),
		zeroIPv4SegZeroPrefix.ToIP(),
		zeroIPv4SegPrefixBlock.ToIP(),
		prefLen)
	result = createIPv4Section(divs)
	result.prefixLength = newPref
	result.isMult = isMultiple
	return result
}

// NewIPv4SectionFromUint32 constructs an IPv4 address section of the given segment count from the given value.
func NewIPv4SectionFromUint32(value uint32, segmentCount int) (res *IPv4AddressSection) {
	return NewIPv4SectionFromPrefixedUint32(value, segmentCount, nil)
}

// NewIPv4SectionFromPrefixedUint32 constructs an IPv4 address or prefix block section of the given segment count from the given value and prefix length.
func NewIPv4SectionFromPrefixedUint32(value uint32, segmentCount int, prefixLength PrefixLen) (res *IPv4AddressSection) {
	if segmentCount < 0 {
		segmentCount = IPv4SegmentCount
	}
	segments := createSegmentsUint64(
		segmentCount,
		0,
		uint64(value),
		IPv4BytesPerSegment,
		IPv4BitsPerSegment,
		ipv4Network.getIPAddressCreator(),
		prefixLength)
	res = createIPv4Section(segments)
	if prefixLength != nil {
		assignPrefix(prefixLength, segments, res.ToIP(), false, false, BitCount(segmentCount<<ipv4BitsToSegmentBitshift))
	} else {
		res.cache.uint32Cache = &value
	}
	return
}

// NewIPv4SectionFromBytes constructs an IPv4 address section from the given byte slice.
// The segment count is determined by the slice length, even if the segment count exceeds 4 segments.
func NewIPv4SectionFromBytes(bytes []byte) *IPv4AddressSection {
	res, _ := newIPv4SectionFromBytes(bytes, len(bytes), nil, false)
	return res
}

// Useful if the byte array has leading zeros

// NewIPv4SectionFromSegmentedBytes constructs an IPv4 address section from the given byte slice.
// It allows you to specify the segment count for the supplied bytes.
// If the slice is too large for the given number of segments, an error is returned, although leading zeros are tolerated.
func NewIPv4SectionFromSegmentedBytes(bytes []byte, segmentCount int) (res *IPv4AddressSection, err addrerr.AddressValueError) {
	return newIPv4SectionFromBytes(bytes, segmentCount, nil, false)
}

// NewIPv4SectionFromPrefixedBytes constructs an IPv4 address or prefix block section from the given byte slice and prefix length.
// It allows you to specify the segment count for the supplied bytes.
// If the slice is too large for the given number of segments, an error is returned, although leading zeros are tolerated.
func NewIPv4SectionFromPrefixedBytes(bytes []byte, segmentCount int, prefixLength PrefixLen) (res *IPv4AddressSection, err addrerr.AddressValueError) {
	return newIPv4SectionFromBytes(bytes, segmentCount, prefixLength, false)
}

func newIPv4SectionFromBytes(bytes []byte, segmentCount int, prefixLength PrefixLen, singleOnly bool) (res *IPv4AddressSection, err addrerr.AddressValueError) {
	if segmentCount < 0 {
		segmentCount = len(bytes)
	}
	expectedByteCount := segmentCount
	segments, err := toSegments(
		bytes,
		segmentCount,
		IPv4BytesPerSegment,
		IPv4BitsPerSegment,
		ipv4Network.getIPAddressCreator(),
		prefixLength)
	if err == nil {
		res = createIPv4Section(segments)
		if prefixLength != nil {
			assignPrefix(prefixLength, segments, res.ToIP(), singleOnly, false, BitCount(segmentCount<<ipv4BitsToSegmentBitshift))
		}
		if expectedByteCount == len(bytes) && len(bytes) > 0 {
			bytes = clone(bytes)
			res.cache.bytesCache = &bytesCache{lowerBytes: bytes}
			if !res.isMult { // not a prefix block
				res.cache.bytesCache.upperBytes = bytes
			}
		}
	}
	return
}

// NewIPv4SectionFromVals constructs an IPv4 address section of the given segment count from the given values.
func NewIPv4SectionFromVals(vals IPv4SegmentValueProvider, segmentCount int) (res *IPv4AddressSection) {
	res = NewIPv4SectionFromPrefixedRange(vals, nil, segmentCount, nil)
	return
}

// NewIPv4SectionFromPrefixedVals constructs an IPv4 address or prefix block section of the given segment count from the given values and prefix length.
func NewIPv4SectionFromPrefixedVals(vals IPv4SegmentValueProvider, segmentCount int, prefixLength PrefixLen) (res *IPv4AddressSection) {
	return NewIPv4SectionFromPrefixedRange(vals, nil, segmentCount, prefixLength)
}

// NewIPv4SectionFromRange constructs an IPv4 subnet section of the given segment count from the given values.
func NewIPv4SectionFromRange(vals, upperVals IPv4SegmentValueProvider, segmentCount int) (res *IPv4AddressSection) {
	res = NewIPv4SectionFromPrefixedRange(vals, upperVals, segmentCount, nil)
	return
}

// NewIPv4SectionFromPrefixedRange constructs an IPv4 subnet section of the given segment count from the given values and prefix length.
func NewIPv4SectionFromPrefixedRange(vals, upperVals IPv4SegmentValueProvider, segmentCount int, prefixLength PrefixLen) (res *IPv4AddressSection) {
	return newIPv4SectionFromPrefixedSingle(vals, upperVals, segmentCount, prefixLength, false)
}

func newIPv4SectionFromPrefixedSingle(vals, upperVals IPv4SegmentValueProvider, segmentCount int, prefixLength PrefixLen, singleOnly bool) (res *IPv4AddressSection) {
	if segmentCount < 0 {
		segmentCount = 0
	}
	segments, isMultiple := createSegments(
		WrapIPv4SegmentValueProvider(vals),
		WrapIPv4SegmentValueProvider(upperVals),
		segmentCount,
		IPv4BitsPerSegment,
		ipv4Network.getIPAddressCreator(),
		prefixLength)
	res = createIPv4Section(segments)
	res.isMult = isMultiple
	if prefixLength != nil {
		assignPrefix(prefixLength, segments, res.ToIP(), singleOnly, false, BitCount(segmentCount<<ipv4BitsToSegmentBitshift))
	}
	return
}

// IPv4AddressSection represents a section of an IPv4 address comprising 0 to 4 IPv4 address segments.
// The zero values is a section with zero-segments.
type IPv4AddressSection struct {
	ipAddressSectionInternal
}

// containsSame returns whether this address section contains all address sections in the given address section collection of the same type.
func (addr *IPv4AddressSection) containsSame(other *IPv4AddressSection) bool {
	return addr.Contains(other)
}

// Contains returns whether this is same type and version as the given address section and whether it contains all values in the given section.
//
// Sections must also have the same number of segments to be comparable, otherwise false is returned.
func (section *IPv4AddressSection) Contains(other AddressSectionType) bool {
	if section == nil {
		return other == nil || other.ToSectionBase() == nil
	}
	return section.contains(other)
}

// Overlaps returns whether this is same type and version as the given address section and whether it overlaps the given section, both sections containing at least one individual section in common.
//
// Sections must also have the same number of segments to be comparable, otherwise false is returned.
func (section *IPv4AddressSection) Overlaps(other AddressSectionType) bool {
	if section == nil {
		return other == nil || other.ToSectionBase() == nil
	}
	return section.overlaps(other)
}

// Equal returns whether the given address section is equal to this address section.
// Two address sections are equal if they represent the same set of sections.
// They must match:
//   - type/version: IPv4
//   - segment counts
//   - segment value ranges
//
// Prefix lengths are ignored.
func (section *IPv4AddressSection) Equal(other AddressSectionType) bool {
	if section == nil {
		return other == nil || other.ToSectionBase() == nil
	}
	return section.equal(other)
}

// Compare returns a negative integer, zero, or a positive integer if this address section is less than, equal, or greater than the given item.
// Any address item is comparable to any other.  All address items use CountComparator to compare.
func (section *IPv4AddressSection) Compare(item AddressItem) int {
	return CountComparator.Compare(section, item)
}

// CompareSize compares the counts of two address sections or items, the number of individual sections or other items represented.
//
// Rather than calculating counts with GetCount, there can be more efficient ways of determining whether this section represents more individual address sections than another.
//
// CompareSize returns a positive integer if this address section has a larger count than the one given, zero if they are the same, or a negative integer if the other has a larger count.
func (section *IPv4AddressSection) CompareSize(other AddressItem) int {
	if section == nil {
		if isNilItem(other) {
			return 0
		}
		// we have size 0, other has size >= 1
		return -1
	}
	return section.compareSize(other)
}

// GetBitsPerSegment returns the number of bits comprising each segment in this section.  Segments in the same address section are equal length.
func (section *IPv4AddressSection) GetBitsPerSegment() BitCount {
	return IPv4BitsPerSegment
}

// GetBytesPerSegment returns the number of bytes comprising each segment in this section.  Segments in the same address section are equal length.
func (section *IPv4AddressSection) GetBytesPerSegment() int {
	return IPv4BytesPerSegment
}

// GetIPVersion returns IPv4, the IP version of this address section.
func (section *IPv4AddressSection) GetIPVersion() IPVersion {
	return IPv4
}

// IsMultiple returns  whether this section represents multiple values.
func (section *IPv4AddressSection) IsMultiple() bool {
	return section != nil && section.isMultiple()
}

// IsPrefixed returns whether this section has an associated prefix length.
func (section *IPv4AddressSection) IsPrefixed() bool {
	return section != nil && section.isPrefixed()
}

// GetCount returns the count of possible distinct values for this section.
// It is the same as GetIPv4Count but returns the value as a big integer instead of a uint64.
// If not representing multiple values, the count is 1,
// unless this is a division grouping with no divisions, or an address section with no segments, in which case it is 0.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (section *IPv4AddressSection) GetCount() *big.Int {
	if section == nil {
		return bigZero()
	}
	return section.cacheCount(func() *big.Int {
		return bigZero().SetUint64(section.getIPv4Count())
	})
}

func (section *IPv4AddressSection) getCachedCount() *big.Int {
	if section == nil {
		return bigZero()
	}
	return section.cachedCount(func() *big.Int {
		return bigZero().SetUint64(section.getIPv4Count())
	})
}

// GetIPv4Count returns the count of possible distinct values for this section.
// It is the same as GetCount but returns the value as a uint64 instead of a big integer.
// If not representing multiple values, the count is 1,
// unless this is a division grouping with no divisions, or an address section with no segments, in which case it is 0.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (section *IPv4AddressSection) GetIPv4Count() uint64 {
	if section == nil {
		return 0
	}
	return section.getCachedCount().Uint64()
}

func (section *IPv4AddressSection) getIPv4Count() uint64 {
	if !section.isMultiple() {
		return 1
	}
	return longCount(section.ToSectionBase(), section.GetSegmentCount())
}

// GetPrefixCount returns the number of distinct prefix values in this item.
//
// The prefix length is given by GetPrefixLen.
//
// If this has a non-nil prefix length, returns the number of distinct prefix values.
//
// If this has a nil prefix length, returns the same value as GetCount.
func (section *IPv4AddressSection) GetPrefixCount() *big.Int {
	return section.cachePrefixCount(func() *big.Int {
		return bigZero().SetUint64(section.getIPv4PrefixCount())
	})
}

// GetIPv4PrefixCount returns the number of distinct prefix values in this section.
// It is similar to GetPrefixCount but returns a uint64.
//
// The prefix length is given by GetPrefixLen.
//
// If this has a non-nil prefix length, returns the number of distinct prefix values.
//
// If this has a nil prefix length, returns the same value as GetIPv4Count.
func (section *IPv4AddressSection) GetIPv4PrefixCount() uint64 {
	return section.cacheUint64PrefixCount(func() uint64 {
		return section.getIPv4PrefixCount()
	})
}

func (section *IPv4AddressSection) getIPv4PrefixCount() uint64 {
	prefixLength := section.getPrefixLen()
	if prefixLength == nil {
		return section.GetIPv4Count()
	}
	return section.GetIPv4PrefixCountLen(prefixLength.bitCount())
}

// GetPrefixCountLen returns the number of distinct prefix values in this item for the given prefix length.
func (section *IPv4AddressSection) GetPrefixCountLen(prefixLen BitCount) *big.Int {
	if prefixLen <= 0 {
		return bigOne()
	} else if bc := section.GetBitCount(); prefixLen > bc {
		prefixLen = bc
	}
	return section.calcCount(func() *big.Int { return bigZero().SetUint64(section.GetIPv4PrefixCountLen(prefixLen)) })
}

// GetIPv4PrefixCountLen returns the number of distinct prefix values in this item for the given prefix length.
//
// It is the same as GetPrefixCountLen but returns a uint64, not a *big.Int.
func (section *IPv4AddressSection) GetIPv4PrefixCountLen(prefixLength BitCount) uint64 {
	if !section.isMultiple() {
		return 1
	} else if prefixLength >= section.GetBitCount() {
		return section.GetIPv4Count()
	} else if prefixLength < 0 {
		prefixLength = 0
	}
	return longPrefixCount(section.ToSectionBase(), prefixLength)
}

// GetIPv4BlockCount returns the count of distinct values in the given number of initial (more significant) segments.
// It is similar to GetBlockCount but returns a uint64 instead of a big integer.
func (section *IPv4AddressSection) GetIPv4BlockCount(segmentCount int) uint64 {
	if !section.isMultiple() {
		return 1
	}
	return longCount(section.ToSectionBase(), segmentCount)
}

// GetBlockCount returns the count of distinct values in the given number of initial (more significant) segments.
// It is similar to GetIPv4BlockCount but returns a big integer instead of a uint64.
func (section *IPv4AddressSection) GetBlockCount(segmentCount int) *big.Int {
	if segmentCount <= 0 {
		return bigOne()
	}
	return section.calcCount(func() *big.Int { return bigZero().SetUint64(section.GetIPv4BlockCount(segmentCount)) })
}

// GetSegment returns the segment at the given index.
// The first segment is at index 0.
// GetSegment will panic given a negative index or an index matching or larger than the segment count.
func (section *IPv4AddressSection) GetSegment(index int) *IPv4AddressSegment {
	return section.getDivision(index).ToIPv4()
}

// ForEachSegment visits each segment in order from most-significant to least, the most significant with index 0, calling the given function for each, terminating early if the function returns true.
// Returns the number of visited segments.
func (section *IPv4AddressSection) ForEachSegment(consumer func(segmentIndex int, segment *IPv4AddressSegment) (stop bool)) int {
	divArray := section.getDivArray()
	if divArray != nil {
		for i, div := range divArray {
			if consumer(i, div.ToIPv4()) {
				return i + 1
			}
		}
	}
	return len(divArray)
}

// GetTrailingSection gets the subsection from the series starting from the given index.
// The first segment is at index 0.
func (section *IPv4AddressSection) GetTrailingSection(index int) *IPv4AddressSection {
	return section.GetSubSection(index, section.GetSegmentCount())
}

// GetSubSection gets the subsection from the series starting from the given index and ending just before the give endIndex.
// The first segment is at index 0.
func (section *IPv4AddressSection) GetSubSection(index, endIndex int) *IPv4AddressSection {
	return section.getSubSection(index, endIndex).ToIPv4()
}

// GetNetworkSection returns a subsection containing the segments with the network bits of the section.
// The returned section will have only as many segments as needed as determined by the existing CIDR network prefix length.
//
// If this series has no CIDR prefix length, the returned network section will
// be the entire series as a prefixed section with prefix length matching the address bit length.
func (section *IPv4AddressSection) GetNetworkSection() *IPv4AddressSection {
	return section.getNetworkSection().ToIPv4()
}

// GetNetworkSectionLen returns a subsection containing the segments with the network of the section, the prefix bits according to the given prefix length.
// The returned section will have only as many segments as needed to contain the network.
//
// The new section will be assigned the given prefix length,
// unless the existing prefix length is smaller, in which case the existing prefix length will be retained.
func (section *IPv4AddressSection) GetNetworkSectionLen(prefLen BitCount) *IPv4AddressSection {
	return section.getNetworkSectionLen(prefLen).ToIPv4()
}

// GetHostSection returns a subsection containing the segments with the host of the address section, the bits beyond the CIDR network prefix length.
// The returned section will have only as many segments as needed to contain the host.
//
// If this series has no prefix length, the returned host section will be the full section.
func (section *IPv4AddressSection) GetHostSection() *IPv4AddressSection {
	return section.getHostSection().ToIPv4()
}

// GetHostSectionLen returns a subsection containing the segments with the host of the address section, the bits beyond the given CIDR network prefix length.
// The returned section will have only as many segments as needed to contain the host.
// The returned section will have an assigned prefix length indicating the beginning of the host.
func (section *IPv4AddressSection) GetHostSectionLen(prefLen BitCount) *IPv4AddressSection {
	return section.getHostSectionLen(prefLen).ToIPv4()
}

// GetNetworkMask returns the network mask associated with the CIDR network prefix length of this address section.
// If this section has no prefix length, then the all-ones mask is returned.
func (section *IPv4AddressSection) GetNetworkMask() *IPv4AddressSection {
	return section.getNetworkMask(ipv4Network).ToIPv4()
}

// GetHostMask returns the host mask associated with the CIDR network prefix length of this address section.
// If this section has no prefix length, then the all-ones mask is returned.
func (section *IPv4AddressSection) GetHostMask() *IPv4AddressSection {
	return section.getHostMask(ipv4Network).ToIPv4()
}

// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
func (section *IPv4AddressSection) CopySubSegments(start, end int, segs []*IPv4AddressSegment) (count int) {
	start, end, targetStart := adjust1To1StartIndices(start, end, section.GetDivisionCount(), len(segs))
	segs = segs[targetStart:]
	return section.forEachSubDivision(start, end, func(index int, div *AddressDivision) {
		segs[index] = div.ToIPv4()
	}, len(segs))
}

// CopySegments copies the existing segments into the given slice,
// as much as can be fit into the slice, returning the number of segments copied.
func (section *IPv4AddressSection) CopySegments(segs []*IPv4AddressSegment) (count int) {
	return section.ForEachSegment(func(index int, seg *IPv4AddressSegment) (stop bool) {
		if stop = index >= len(segs); !stop {
			segs[index] = seg
		}
		return
	})
}

// GetSegments returns a slice with the address segments.  The returned slice is not backed by the same array as this section.
func (section *IPv4AddressSection) GetSegments() (res []*IPv4AddressSegment) {
	res = make([]*IPv4AddressSegment, section.GetSegmentCount())
	section.CopySegments(res)
	return
}

// Mask applies the given mask to all address sections represented by this secction, returning the result.
//
// If the sections do not have a comparable number of segments, an error is returned.
//
// If this represents multiple addresses, and applying the mask to all addresses creates a set of addresses
// that cannot be represented as a sequential range within each segment, then an error is returned.
func (section *IPv4AddressSection) Mask(other *IPv4AddressSection) (res *IPv4AddressSection, err addrerr.IncompatibleAddressError) {
	return section.maskPrefixed(other, true)
}

func (section *IPv4AddressSection) maskPrefixed(other *IPv4AddressSection, retainPrefix bool) (res *IPv4AddressSection, err addrerr.IncompatibleAddressError) {
	sec, err := section.mask(other.ToIP(), retainPrefix)
	if err == nil {
		res = sec.ToIPv4()
	}
	return
}

// BitwiseOr does the bitwise disjunction with this address section, useful when subnetting.
// It is similar to Mask which does the bitwise conjunction.
//
// The operation is applied to all individual addresses and the result is returned.
//
// If this represents multiple address sections, and applying the operation to all sections creates a set of sections
// that cannot be represented as a sequential range within each segment, then an error is returned.
func (section *IPv4AddressSection) BitwiseOr(other *IPv4AddressSection) (res *IPv4AddressSection, err addrerr.IncompatibleAddressError) {
	return section.bitwiseOrPrefixed(other, true)
}

func (section *IPv4AddressSection) bitwiseOrPrefixed(other *IPv4AddressSection, retainPrefix bool) (res *IPv4AddressSection, err addrerr.IncompatibleAddressError) {
	sec, err := section.bitwiseOr(other.ToIP(), retainPrefix)
	if err == nil {
		res = sec.ToIPv4()
	}
	return
}

// MatchesWithMask applies the mask to this address section and then compares the result with the given address section,
// returning true if they match, false otherwise.  To match, both the given section and mask must have the same number of segments as this section.
func (section *IPv4AddressSection) MatchesWithMask(other *IPv4AddressSection, mask *IPv4AddressSection) bool {
	return section.matchesWithMask(other.ToIP(), mask.ToIP())
}

// Subtract subtracts the given subnet sections from this subnet section, returning an array of sections for the result (the subnet sections will not be contiguous so an array is required).
//
// Subtract  computes the subnet difference, the set of address sections in this address section but not in the provided section.
// This is also known as the relative complement of the given argument in this subnet section.
//
// This is set subtraction, not subtraction of values.
func (section *IPv4AddressSection) Subtract(other *IPv4AddressSection) (res []*IPv4AddressSection, err addrerr.SizeMismatchError) {
	sections, err := section.subtract(other.ToIP())
	if err == nil {
		res = cloneTo(sections, (*IPAddressSection).ToIPv4)
	}
	return
}

// Intersect returns the subnet sections whose individual sections are found in both this and the given subnet section argument, or nil if no such sections exist.
//
// This is also known as the conjunction of the two sets of address sections.
//
// If the two sections have different segment counts, an error is returned.
func (section *IPv4AddressSection) Intersect(other *IPv4AddressSection) (res *IPv4AddressSection, err addrerr.SizeMismatchError) {
	sec, err := section.intersect(other.ToIP())
	if err == nil {
		res = sec.ToIPv4()
	}
	return
}

// GetLower returns the section in the range with the lowest numeric value,
// which will be the same section if it represents a single value.
// For example, for "1.2-3.4.5-6", the section "1.2.4.5" is returned.
func (section *IPv4AddressSection) GetLower() *IPv4AddressSection {
	return section.getLower().ToIPv4()
}

// GetUpper returns the section in the range with the highest numeric value,
// which will be the same section if it represents a single value.
// For example, for "1.2-3.4.5-6", the section "1.3.4.6" is returned.
func (section *IPv4AddressSection) GetUpper() *IPv4AddressSection {
	return section.getUpper().ToIPv4()
}

// Uint32Value returns the lowest address in the address section range as a uint32.
func (section *IPv4AddressSection) Uint32Value() uint32 {
	cache := section.cache
	if cache == nil {
		return section.uint32Value()
	}
	res := (*uint32)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache.uint32Cache))))
	if res == nil {
		val := section.uint32Value()
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.uint32Cache))
		atomicStorePointer(dataLoc, unsafe.Pointer(&val))
		return val
	}
	return *res
}

// Uint32Value returns the lowest address in the address section range as a uint32.
func (section *IPv4AddressSection) uint32Value() uint32 {
	segCount := section.GetSegmentCount()
	if segCount == 0 {
		return 0
	}
	arr := section.getDivArray()
	val := uint32(arr[0].getDivisionValue())
	bitsPerSegment := section.GetBitsPerSegment()
	for i := 1; i < segCount; i++ {
		val = (val << uint(bitsPerSegment)) | uint32(arr[i].getDivisionValue())
	}
	return val
}

// UpperUint32Value returns the highest address in the address section range as a uint32.
func (section *IPv4AddressSection) UpperUint32Value() uint32 {
	if !section.IsMultiple() {
		return section.Uint32Value()
	}
	segCount := section.GetSegmentCount()
	if segCount == 0 {
		return 0
	}
	arr := section.getDivArray()
	val := uint32(arr[0].getUpperDivisionValue())
	bitsPerSegment := section.GetBitsPerSegment()
	for i := 1; i < segCount; i++ {
		val = (val << uint(bitsPerSegment)) | uint32(arr[i].getUpperDivisionValue())
	}
	return val
}

// ToZeroHost converts the address section to one in which all individual address sections have a host of zero,
// the host being the bits following the prefix length.
// If the address section has no prefix length, then it returns an all-zero address section.
//
// The returned section will have the same prefix and prefix length.
//
// This returns an error if the section is a range of address sections which cannot be converted to a range in which all sections have zero hosts,
// because the conversion results in a segment that is not a sequential range of values.
func (section *IPv4AddressSection) ToZeroHost() (*IPv4AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.toZeroHost(false)
	return res.ToIPv4(), err
}

// ToZeroHostLen converts the address section to one in which all individual sections have a host of zero,
// the host being the bits following the given prefix length.
// If this address section has the same prefix length, then the returned one will too, otherwise the returned section will have no prefix length.
//
// This returns an error if the section is a range of which cannot be converted to a range in which all sections have zero hosts,
// because the conversion results in a segment that is not a sequential range of values.
func (section *IPv4AddressSection) ToZeroHostLen(prefixLength BitCount) (*IPv4AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.toZeroHostLen(prefixLength)
	return res.ToIPv4(), err
}

// ToZeroNetwork converts the address section to one in which all individual address sections have a network of zero,
// the network being the bits within the prefix length.
// If the address section has no prefix length, then it returns an all-zero address section.
//
// The returned address section will have the same prefix length.
func (section *IPv4AddressSection) ToZeroNetwork() *IPv4AddressSection {
	return section.toZeroNetwork().ToIPv4()
}

// ToMaxHost converts the address section to one in which all individual address sections have a host of all one-bits, the max value,
// the host being the bits following the prefix length.
// If the address section has no prefix length, then it returns an all-ones section, the max address section.
//
// The returned address section will have the same prefix and prefix length.
//
// This returns an error if the address section is a range of address sections which cannot be converted to a range in which all sections have max hosts,
// because the conversion results in a segment that is not a sequential range of values.
func (section *IPv4AddressSection) ToMaxHost() (*IPv4AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.toMaxHost()
	return res.ToIPv4(), err
}

// ToMaxHostLen converts the address section to one in which all individual address sections have a host of all one-bits, the max host,
// the host being the bits following the given prefix length.
// If this section has the same prefix length, then the resulting section will too, otherwise the resulting section will have no prefix length.
//
// This returns an error if the section is a range of address sections which cannot be converted to a range in which all address sections have max hosts,
// because the conversion results in a segment that is not a sequential range of values.
func (section *IPv4AddressSection) ToMaxHostLen(prefixLength BitCount) (*IPv4AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.toMaxHostLen(prefixLength)
	return res.ToIPv4(), err
}

// ToPrefixBlock returns the section with the same prefix as this section while the remaining bits span all values.
// The returned section will be the block of all sections with the same prefix.
//
// If this section has no prefix, this section is returned.
func (section *IPv4AddressSection) ToPrefixBlock() *IPv4AddressSection {
	return section.toPrefixBlock().ToIPv4()
}

// ToPrefixBlockLen returns the section with the same prefix of the given length as this section while the remaining bits span all values.
// The returned section will be the block of all sections with the same prefix.
func (section *IPv4AddressSection) ToPrefixBlockLen(prefLen BitCount) *IPv4AddressSection {
	return section.toPrefixBlockLen(prefLen).ToIPv4()
}

// ToBlock creates a new block of address sections by changing the segment at the given index to have the given lower and upper value,
// and changing the following segments to be full-range.
func (section *IPv4AddressSection) ToBlock(segmentIndex int, lower, upper SegInt) *IPv4AddressSection {
	return section.toBlock(segmentIndex, lower, upper).ToIPv4()
}

// WithoutPrefixLen provides the same address section but with no prefix length.  The values remain unchanged.
func (section *IPv4AddressSection) WithoutPrefixLen() *IPv4AddressSection {
	if !section.IsPrefixed() {
		return section
	}
	return section.withoutPrefixLen().ToIPv4()
}

// SetPrefixLen sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the address section.
// The provided prefix length will be adjusted to these boundaries if necessary.
func (section *IPv4AddressSection) SetPrefixLen(prefixLen BitCount) *IPv4AddressSection {
	return section.setPrefixLen(prefixLen).ToIPv4()
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
func (section *IPv4AddressSection) SetPrefixLenZeroed(prefixLen BitCount) (*IPv4AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.setPrefixLenZeroed(prefixLen)
	return res.ToIPv4(), err
}

// AdjustPrefixLen increases or decreases the prefix length by the given increment.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the address section.
//
// If this address section has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
func (section *IPv4AddressSection) AdjustPrefixLen(prefixLen BitCount) *IPv4AddressSection {
	return section.adjustPrefixLen(prefixLen).ToIPv4()
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
func (section *IPv4AddressSection) AdjustPrefixLenZeroed(prefixLen BitCount) (*IPv4AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.adjustPrefixLenZeroed(prefixLen)
	return res.ToIPv4(), err
}

// AssignPrefixForSingleBlock returns the equivalent prefix block that matches exactly the range of values in this address section.
// The returned block will have an assigned prefix length indicating the prefix length for the block.
//
// There may be no such address section - it is required that the range of values match the range of a prefix block.
// If there is no such address section, then nil is returned.
func (section *IPv4AddressSection) AssignPrefixForSingleBlock() *IPv4AddressSection {
	return section.assignPrefixForSingleBlock().ToIPv4()
}

// AssignMinPrefixForBlock returns an equivalent address section, assigned the smallest prefix length possible,
// such that the prefix block for that prefix length is in this address section.
//
// In other words, this method assigns a prefix length to this address section matching the largest prefix block in this address section.
func (section *IPv4AddressSection) AssignMinPrefixForBlock() *IPv4AddressSection {
	return section.assignMinPrefixForBlock().ToIPv4()
}

// Iterator provides an iterator to iterate through the individual address sections of this address section.
//
// When iterating, the prefix length is preserved.  Remove it using WithoutPrefixLen prior to iterating if you wish to drop it from all individual address sections.
//
// Call IsMultiple to determine if this instance represents multiple address sections, or GetCount for the count.
func (section *IPv4AddressSection) Iterator() Iterator[*IPv4AddressSection] {
	if section == nil {
		return ipv4SectionIterator{nilSectIterator()}
	}
	return ipv4SectionIterator{section.sectionIterator(nil)}
}

// PrefixIterator provides an iterator to iterate through the individual prefixes of this address section,
// each iterated element spanning the range of values for its prefix.
//
// It is similar to the prefix block iterator, except for possibly the first and last iterated elements, which might not be prefix blocks,
// instead constraining themselves to values from this address section.
//
// If the series has no prefix length, then this is equivalent to Iterator.
func (section *IPv4AddressSection) PrefixIterator() Iterator[*IPv4AddressSection] {
	return ipv4SectionIterator{section.prefixIterator(false)}
}

// PrefixBlockIterator provides an iterator to iterate through the individual prefix blocks, one for each prefix of this address section.
// Each iterated address section will be a prefix block with the same prefix length as this address section.
//
// If this address section has no prefix length, then this is equivalent to Iterator.
func (section *IPv4AddressSection) PrefixBlockIterator() Iterator[*IPv4AddressSection] {
	return ipv4SectionIterator{section.prefixIterator(true)}
}

// BlockIterator Iterates through the address sections that can be obtained by iterating through all the upper segments up to the given segment count.
// The segments following remain the same in all iterated sections.
func (section *IPv4AddressSection) BlockIterator(segmentCount int) Iterator[*IPv4AddressSection] {
	return ipv4SectionIterator{section.blockIterator(segmentCount)}
}

// SequentialBlockIterator iterates through the sequential address sections that make up this address section.
//
// Practically, this means finding the count of segments for which the segments that follow are not full range, and then using BlockIterator with that segment count.
//
// Use GetSequentialBlockCount to get the number of iterated elements.
func (section *IPv4AddressSection) SequentialBlockIterator() Iterator[*IPv4AddressSection] {
	return ipv4SectionIterator{section.sequentialBlockIterator()}
}

// ToDivGrouping converts to an AddressDivisionGrouping, a polymorphic type usable with all address sections and division groupings.
// Afterwards, you can convert back with ToIPv4.
//
// ToDivGrouping can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *IPv4AddressSection) ToDivGrouping() *AddressDivisionGrouping {
	return section.ToSectionBase().ToDivGrouping()
}

// ToSectionBase converts to an AddressSection, a polymorphic type usable with all address sections.
// Afterwards, you can convert back with ToIPv4.
//
// ToSectionBase can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *IPv4AddressSection) ToSectionBase() *AddressSection {
	return section.ToIP().ToSectionBase()
}

// ToIP converts to an IPAddressSection, a polymorphic type usable with all IP address sections.
//
// ToIP can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *IPv4AddressSection) ToIP() *IPAddressSection {
	return (*IPAddressSection)(section)
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
func (section *IPv4AddressSection) IncrementBoundary(increment int64) *IPv4AddressSection {
	return section.incrementBoundary(increment).ToIPv4()
}

func getIPv4MaxValueLong(segmentCount int) uint64 {
	return macMaxValues[segmentCount]
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
func (section *IPv4AddressSection) Increment(inc int64) *IPv4AddressSection {
	if inc == 0 && !section.isMultiple() {
		return section
	}
	lowerValueFunc := func() uint64 {
		return uint64(section.Uint32Value())
	}
	upperValueFunc := func() uint64 {
		return uint64(section.UpperUint32Value())
	}
	if isOverflow := checkOverflow(inc, lowerValueFunc, upperValueFunc, section.GetIPv4Count, func() uint64 { return getIPv4MaxValueLong(section.GetSegmentCount()) }, section.IsSequential); isOverflow {
		return nil
	}

	return increment(
		section.ToSectionBase(),
		inc,
		ipv4Network.getIPAddressCreator(),
		section.GetIPv4Count,
		lowerValueFunc,
		upperValueFunc,
		section.getLower,
		section.getUpper,
		section.getPrefixLen()).ToIPv4()
}

func low64IPv4(section *AddressSection) uint64 {
	return uint64(section.ToIPv4().Uint32Value())
}

func low64UpperIPv4(section *AddressSection) uint64 {
	return uint64(section.ToIPv4().UpperUint32Value())
}

func (section *IPv4AddressSection) enumerateAddrIPv4(other AddressSectionType) (val int64, ok bool) {
	if otherSection := other.ToSectionBase(); otherSection.IsIPv4() {
		return enumerateSmall(section.ToSectionBase(), otherSection, low64IPv4, low64UpperIPv4)
	}
	return
}

// Enumerate indicates where an individual address section sits relative to the address section range ordering.
//
// Determines how many address section elements of a range precede the given address section element, if the address section is in the range.
// If above the range, it is the distance to the upper boundary added to the range count less one, and if below the range, the distance to the lower boundary.
//
// In other words, if the given address section is not in the range but above it, returns the number of address sections preceding the address from the upper range boundary,
// added to one less than the total number of range address sections.  If the given address section is not in the subnet but below it, returns the number of address sections following the address section to the lower subnet boundary.
//
// If the argument is not in the range, but neither above nor below the range, then 0 is returned and ok is false.
//
// Enumerate returns 0 and false when the argument is multi-valued. The argument must be an individual address section.
//
// When this is also an individual address section, the returned value is the distance (difference) between the two address section values.
//
// If the given address section does not have the same version or type, then ok is false.
//
// Sections must also have the same number of segments to be comparable, otherwise ok is false.
func (section *IPv4AddressSection) EnumerateIPv4(other AddressSectionType) (val int64, ok bool) {
	if other != nil {
		if otherSection := other.ToSectionBase(); otherSection != nil && otherSection.IsIPv4() {
			if matches, count := section.matchesTypeAndCount(otherSection); matches && count <= 8 {
				return enumerateSmall(section.ToSectionBase(), otherSection, low64IPv4, low64UpperIPv4)
			}
		}
	}
	return
}

func (section *IPv4AddressSection) enumerateAddr(other AddressSectionType) *big.Int {
	if val, ok := section.enumerateAddrIPv4(other); ok {
		return big.NewInt(val)
	}
	return nil
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
func (section *IPv4AddressSection) Enumerate(other AddressSectionType) *big.Int {
	if other != nil {
		if otherSection := other.ToSectionBase(); otherSection != nil {
			if matches, count := section.matchesTypeAndCount(otherSection); matches {
				if count <= 8 {
					if val, ok := enumerateSmall(section.ToSectionBase(), otherSection, low64IPv4, low64UpperIPv4); ok {
						return big.NewInt(val)
					}
					return nil
				}
				return enumerateBig(section.ToSectionBase(), otherSection, low64IPv4, low64UpperIPv4)
			}
		}
	}
	return nil
}

// SpanWithPrefixBlocks returns an array of prefix blocks that spans the same set of individual address sections as this section.
//
// Unlike SpanWithPrefixBlocksTo, the result only includes blocks that are a part of this section.
func (section *IPv4AddressSection) SpanWithPrefixBlocks() []*IPv4AddressSection {
	if section.IsSequential() {
		if section.IsSinglePrefixBlock() {
			return []*IPv4AddressSection{section}
		}
		return getSpanningPrefixBlocks(section, section)
	}
	return spanWithPrefixBlocks(section)
}

// SpanWithPrefixBlocksTo returns the smallest slice of prefix block subnet sections that span from this section to the given section.
//
// If the given section has a different segment count, an error is returned.
//
// The resulting slice is sorted from lowest address value to highest, regardless of the size of each prefix block.
func (section *IPv4AddressSection) SpanWithPrefixBlocksTo(other *IPv4AddressSection) ([]*IPv4AddressSection, addrerr.SizeMismatchError) {
	if err := section.checkSegmentCount(other.ToIP()); err != nil {
		return nil, err
	}
	return getSpanningPrefixBlocks(section, other), nil
}

// SpanWithSequentialBlocks produces the smallest slice of sequential blocks that cover the same set of sections as this.
//
// This slice can be shorter than that produced by SpanWithPrefixBlocks and is never longer.
//
// Unlike SpanWithSequentialBlocksTo, this method only includes values that are a part of this section.
func (section *IPv4AddressSection) SpanWithSequentialBlocks() []*IPv4AddressSection {
	if section.IsSequential() {
		return []*IPv4AddressSection{section}
	}
	return spanWithSequentialBlocks(section)
}

// SpanWithSequentialBlocksTo produces the smallest slice of sequential block address sections that span from this section to the given section.
func (section *IPv4AddressSection) SpanWithSequentialBlocksTo(other *IPv4AddressSection) ([]*IPv4AddressSection, addrerr.SizeMismatchError) {
	if err := section.checkSegmentCount(other.ToIP()); err != nil {
		return nil, err
	}
	return getSpanningSequentialBlocks(section, other), nil
}

// CoverWithPrefixBlockTo returns the minimal-size prefix block section that covers all the address sections spanning from this to the given section.
//
// If the other section has a different segment count, an error is returned.
func (section *IPv4AddressSection) CoverWithPrefixBlockTo(other *IPv4AddressSection) (*IPv4AddressSection, addrerr.SizeMismatchError) {
	res, err := section.coverWithPrefixBlockTo(other.ToIP())
	return res.ToIPv4(), err
}

// CoverWithPrefixBlock returns the minimal-size prefix block that covers all the individual address sections in this section.
// The resulting block will have a larger count than this, unless this section is already a prefix block.
func (section *IPv4AddressSection) CoverWithPrefixBlock() *IPv4AddressSection {
	return section.coverWithPrefixBlock().ToIPv4()
}

func (section *IPv4AddressSection) checkSectionCounts(sections []*IPv4AddressSection) addrerr.SizeMismatchError {
	segCount := section.GetSegmentCount()
	length := len(sections)
	for i := 0; i < length; i++ {
		section2 := sections[i]
		if section2 == nil {
			continue
		}
		if section2.GetSegmentCount() != segCount {
			return &sizeMismatchError{incompatibleAddressError{addressError{key: "ipaddress.error.sizeMismatch"}}}
		}
	}
	return nil
}

// MergeToSequentialBlocks merges this with the list of sections to produce the smallest array of sequential blocks.
//
// The resulting slice is sorted from lowest address value to highest, regardless of the size of each prefix block.
func (section *IPv4AddressSection) MergeToSequentialBlocks(sections ...*IPv4AddressSection) ([]*IPv4AddressSection, addrerr.SizeMismatchError) {
	if err := section.checkSectionCounts(sections); err != nil {
		return nil, err
	}
	return getMergedSequentialBlocks(cloneSeries(section, sections)), nil
}

// MergeToPrefixBlocks merges this section with the list of sections to produce the smallest array of prefix blocks.
//
// The resulting slice is sorted from lowest value to highest, regardless of the size of each prefix block.
func (section *IPv4AddressSection) MergeToPrefixBlocks(sections ...*IPv4AddressSection) ([]*IPv4AddressSection, addrerr.SizeMismatchError) {
	if err := section.checkSectionCounts(sections); err != nil {
		return nil, err
	}
	return getMergedPrefixBlocks(cloneSeries(section, sections)), nil
}

// ReverseBits returns a new section with the bits reversed.  Any prefix length is dropped.
//
// If the bits within a single segment cannot be reversed because the segment represents a range,
// and reversing the segment values results in a range that is not contiguous, this returns an error.
//
// In practice this means that to be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
//
// If perByte is true, the bits are reversed within each byte, otherwise all the bits are reversed.
func (section *IPv4AddressSection) ReverseBits(perByte bool) (*IPv4AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.reverseBits(perByte)
	return res.ToIPv4(), err
}

// ReverseBytes returns a new section with the bytes reversed.  Any prefix length is dropped.
func (section *IPv4AddressSection) ReverseBytes() *IPv4AddressSection {
	return section.ReverseSegments()
}

// ReverseSegments returns a new section with the segments reversed.
func (section *IPv4AddressSection) ReverseSegments() *IPv4AddressSection {
	if section.GetSegmentCount() <= 1 {
		if section.IsPrefixed() {
			return section.WithoutPrefixLen()
		}
		return section
	}
	res, _ := section.reverseSegments(
		func(i int) (*AddressSegment, addrerr.IncompatibleAddressError) {
			return section.GetSegment(i).WithoutPrefixLen().ToSegmentBase(), nil
		},
	)
	return res.ToIPv4()
}

// Append creates a new section by appending the given section to this section.
func (section *IPv4AddressSection) Append(other *IPv4AddressSection) *IPv4AddressSection {
	count := section.GetSegmentCount()
	return section.ReplaceLen(count, count, other, 0, other.GetSegmentCount())
}

// Insert creates a new section by inserting the given section into this section at the given index.
func (section *IPv4AddressSection) Insert(index int, other *IPv4AddressSection) *IPv4AddressSection {
	return section.insert(index, other.ToIP(), ipv4BitsToSegmentBitshift).ToIPv4()
}

// Replace replaces the segments of this section starting at the given index with the given replacement segments.
func (section *IPv4AddressSection) Replace(index int, replacement *IPv4AddressSection) *IPv4AddressSection {
	return section.ReplaceLen(index, index+replacement.GetSegmentCount(), replacement, 0, replacement.GetSegmentCount())
}

// ReplaceLen replaces segments starting from startIndex and ending before endIndex with the segments starting at replacementStartIndex and
// ending before replacementEndIndex from the replacement section.
func (section *IPv4AddressSection) ReplaceLen(startIndex, endIndex int, replacement *IPv4AddressSection, replacementStartIndex, replacementEndIndex int) *IPv4AddressSection {
	return section.replaceLen(startIndex, endIndex, replacement.ToIP(), replacementStartIndex, replacementEndIndex, ipv4BitsToSegmentBitshift).ToIPv4()
}

// IsAdaptiveZero returns true if the division grouping was originally created as an implicitly zero-valued section or grouping (e.g. IPv4AddressSection{}),
// meaning it was not constructed using a constructor function.
// Such a grouping, which has no divisions or segments, is convertible to an implicitly zero-valued grouping of any type or version, whether IPv6, IPv4, MAC, or other.
// In other words, when a section or grouping is the zero-value, then it is equivalent and convertible to the zero value of any other section or grouping type.
func (section *IPv4AddressSection) IsAdaptiveZero() bool {
	return section != nil && section.matchesZeroGrouping()
}

var (
	ipv4CanonicalParams          = new(addrstr.IPv4StringOptionsBuilder).ToOptions()
	ipv4FullParams               = new(addrstr.IPv4StringOptionsBuilder).SetExpandedSegments(true).SetWildcardOptions(wildcardsRangeOnlyNetworkOnly).ToOptions()
	ipv4NormalizedWildcardParams = new(addrstr.IPv4StringOptionsBuilder).SetWildcardOptions(allWildcards).ToOptions()
	ipv4SqlWildcardParams        = new(addrstr.IPv4StringOptionsBuilder).SetWildcardOptions(allSQLWildcards).ToOptions()

	inetAtonOctalParams       = new(addrstr.IPv4StringOptionsBuilder).SetRadix(Inet_aton_radix_octal.GetRadix()).SetSegmentStrPrefix(Inet_aton_radix_octal.GetSegmentStrPrefix()).ToOptions()
	inetAtonHexParams         = new(addrstr.IPv4StringOptionsBuilder).SetRadix(Inet_aton_radix_hex.GetRadix()).SetSegmentStrPrefix(Inet_aton_radix_hex.GetSegmentStrPrefix()).ToOptions()
	ipv4ReverseDNSParams      = new(addrstr.IPv4StringOptionsBuilder).SetWildcardOptions(allWildcards).SetReverse(true).SetAddressSuffix(IPv4ReverseDnsSuffix).ToOptions()
	ipv4SegmentedBinaryParams = new(addrstr.IPStringOptionsBuilder).SetRadix(2).SetSeparator(IPv4SegmentSeparator).SetSegmentStrPrefix(BinaryPrefix).ToOptions()
)

// ToHexString writes this address section as a single hexadecimal value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0x" prefix.
//
// If a multiple-valued section cannot be written as a single prefix block or a range of two values, an error is returned.
func (section *IPv4AddressSection) ToHexString(with0xPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toHexString(with0xPrefix)
}

// ToOctalString writes this address section as a single octal value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0" prefix.
//
// If a multiple-valued section cannot be written as a single prefix block or a range of two values, an error is returned.
func (section *IPv4AddressSection) ToOctalString(with0Prefix bool) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toOctalString(with0Prefix)
}

// ToBinaryString writes this address section as a single binary value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0b" prefix.
//
// If a multiple-valued section cannot be written as a single prefix block or a range of two values, an error is returned.
func (section *IPv4AddressSection) ToBinaryString(with0bPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toBinaryString(with0bPrefix)
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
// If this section has a prefix length, it will be included in the string.
func (section *IPv4AddressSection) ToCanonicalString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toNormalizedString(ipv4CanonicalParams)
	}
	return cacheStr(&cache.canonicalString,
		func() string {
			return section.toNormalizedString(ipv4CanonicalParams)
		})
}

// ToNormalizedString produces a normalized string for the address section.
//
// For IPv4, it is the same as the canonical string.
//
// If this section has a prefix length, it will be included in the string.
func (section *IPv4AddressSection) ToNormalizedString() string {
	if section == nil {
		return nilString()
	}
	return section.ToCanonicalString()
}

// ToCompressedString produces a short representation of this address section while remaining within the confines of standard representation(s) of the address.
//
// For IPv4, it is the same as the canonical string.
func (section *IPv4AddressSection) ToCompressedString() string {
	if section == nil {
		return nilString()
	}
	return section.ToCanonicalString()
}

// ToNormalizedWildcardString produces a string similar to the normalized string but avoids the CIDR prefix length.
// CIDR addresses will be shown with wildcards and ranges (denoted by '*' and '-') instead of using the CIDR prefix notation.
func (section *IPv4AddressSection) ToNormalizedWildcardString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toNormalizedString(ipv4NormalizedWildcardParams)
	}
	return cacheStr(&cache.normalizedWildcardString,
		func() string {
			return section.toNormalizedString(ipv4NormalizedWildcardParams)
		})
}

// ToCanonicalWildcardString produces a string similar to the canonical string but avoids the CIDR prefix length.
// Address sections with a network prefix length will be shown with wildcards and ranges (denoted by '*' and '-') instead of using the CIDR prefix length notation.
// For IPv4 it is the same as ToNormalizedWildcardString.
func (section *IPv4AddressSection) ToCanonicalWildcardString() string {
	if section == nil {
		return nilString()
	}
	return section.ToNormalizedWildcardString()
}

// ToSegmentedBinaryString writes this address section as segments of binary values preceded by the "0b" prefix.
func (section *IPv4AddressSection) ToSegmentedBinaryString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toNormalizedString(ipv4SegmentedBinaryParams)
	}
	return cacheStr(&cache.segmentedBinaryString,
		func() string {
			return section.toNormalizedString(ipv4SegmentedBinaryParams)
		})
}

// ToSQLWildcardString create a string similar to that from toNormalizedWildcardString except that
// it uses SQL wildcards.  It uses '%' instead of '*' and also uses the wildcard '_'.
func (section *IPv4AddressSection) ToSQLWildcardString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toNormalizedString(ipv4SqlWildcardParams)
	}
	return cacheStr(&cache.sqlWildcardString,
		func() string {
			return section.toNormalizedString(ipv4SqlWildcardParams)
		})
}

// ToFullString produces a string with no compressed segments and all segments of full length with leading zeros,
// which is 3 characters for IPv4 segments.
func (section *IPv4AddressSection) ToFullString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toNormalizedString(ipv4FullParams)
	}
	return cacheStr(&cache.fullString,
		func() string {
			return section.toNormalizedString(ipv4FullParams)
		})
}

// ToReverseDNSString generates the reverse-DNS lookup string.
// For IPV4, the error is always nil.
// For "8.255.4.4" it is "4.4.255.8.in-addr.arpa".
func (section *IPv4AddressSection) ToReverseDNSString() (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toNormalizedString(ipv4ReverseDNSParams), nil
	}
	return cacheStr(&cache.reverseDNSString,
		func() string {
			return section.toNormalizedString(ipv4ReverseDNSParams)
		}), nil
}

// ToPrefixLenString returns a string with a CIDR network prefix length if this address has a network prefix length.
// For IPv4 the string is equivalent to the canonical string.
func (section *IPv4AddressSection) ToPrefixLenString() string {
	if section == nil {
		return nilString()
	}
	return section.ToCanonicalString()
}

// ToSubnetString produces a string with specific formats for subnets.
// The subnet string looks like "1.2.*.*" or "1:2::/16".
//
// In the case of IPv4, this means that wildcards are used instead of a network prefix when a network prefix has been supplied.
func (section *IPv4AddressSection) ToSubnetString() string {
	if section == nil {
		return nilString()
	}
	return section.ToNormalizedWildcardString()
}

// ToCompressedWildcardString produces a string similar to ToNormalizedWildcardString, and in fact
// for IPv4 it is the same as ToNormalizedWildcardString.
func (section *IPv4AddressSection) ToCompressedWildcardString() string {
	if section == nil {
		return nilString()
	}
	return section.ToNormalizedWildcardString()
}

// ToInetAtonString returns a string with a format that is styled from the inet_aton routine.
// The string can have an octal or hexadecimal radix rather than decimal.
// When using octal, the octal segments each have a leading zero prefix of "0", and when using hex, a prefix of "0x".
//
// The allowable radices are 8, 10, and 16.  Any other radix causes a panic.
func (section *IPv4AddressSection) ToInetAtonString(radix Inet_aton_radix) string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if radix == Inet_aton_radix_octal {
		if cache == nil {
			return section.toNormalizedString(inetAtonOctalParams)
		}
		return cacheStr(&cache.inetAtonOctalString,
			func() string {
				return section.toNormalizedString(inetAtonOctalParams)
			})
	} else if radix == Inet_aton_radix_hex {
		if cache == nil {
			return section.toNormalizedString(inetAtonHexParams)
		}
		return cacheStr(&cache.inetAtonHexString,
			func() string {
				return section.toNormalizedString(inetAtonHexParams)
			})
	} else if radix == Inet_aton_radix_decimal {
		return section.ToCanonicalString()
	} else {
		panic(invalidRadix)
	}
}

// ToInetAtonJoinedString returns a string with a format that is styled from the inet_aton routine.
// The string can have an octal or hexadecimal radix rather than decimal,
// and can have less than the typical four IPv4 segments by joining the least significant segments together,
// resulting in a string which just 1, 2 or 3 divisions.
//
// When using octal, the octal segments each have a leading zero prefix of "0", and when using hex, a prefix of "0x".
//
// If this represents a subnet section, this returns an error when unable to join two or more segments
// into a division of a larger bit-length that represents the same set of values.
//
// The allowable radices are 8, 10, and 16.  Any other radix causes a panic.
func (section *IPv4AddressSection) ToInetAtonJoinedString(radix Inet_aton_radix, joinedCount int) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	if joinedCount <= 0 {
		return section.ToInetAtonString(radix), nil
	}
	var stringParams addrstr.IPStringOptions
	if radix == Inet_aton_radix_octal {
		stringParams = inetAtonOctalParams
	} else if radix == Inet_aton_radix_hex {
		stringParams = inetAtonHexParams
	} else if radix == Inet_aton_radix_decimal {
		stringParams = ipv4CanonicalParams
	} else {
		panic(invalidRadix)
	}
	return section.ToNormalizedJoinedString(stringParams, joinedCount)
}

// ToNormalizedJoinedString returns a string with a format that is styled from the inet_aton routine.
// The string can have less than the typical four IPv4 segments by joining the least significant segments together,
// resulting in a string which just 1, 2 or 3 divisions.
//
// The method accepts an argument of string options as well, allowing callers to customize the string in other ways as well.
//
// If this represents a subnet section, this returns an error when unable to join two or more segments
// into a division of a larger bit-length that represents the same set of values.
func (section *IPv4AddressSection) ToNormalizedJoinedString(stringParams addrstr.IPStringOptions, joinedCount int) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	if joinedCount <= 0 || section.GetSegmentCount() <= 1 {
		return section.toNormalizedString(stringParams), nil
	}
	equivalentPart, err := section.ToJoinedSegments(joinedCount) // AddressDivisionSeries
	if err != nil {
		return "", err
	}
	return toNormalizedIPString(stringParams, equivalentPart), nil
}

// ToJoinedSegments returns an AddressDivisionSeries which organizes the address section by joining the least significant segments together.
// If joined count is not a positive number, or this section has less than 2 segments, then this returns the original receiver section.
// Otherwise this returns an AddressDivisionGrouping in which the last division is the division created by joining two or more segments.
//
// If this represents a subnet section, this returns an error when unable to join address segments,
// one of the first with a range of values, into a division of the larger bit-length that represents the same set of values.
func (section *IPv4AddressSection) ToJoinedSegments(joinCount int) (AddressDivisionSeries, addrerr.IncompatibleAddressError) {
	thisCount := section.GetSegmentCount()
	if joinCount <= 0 || thisCount <= 1 {
		return section, nil
	}
	var totalCount int
	if joinCount >= thisCount {
		joinCount = thisCount - 1
		totalCount = 1
	} else {
		totalCount = thisCount - joinCount
	}
	joinedSegment, err := section.joinSegments(joinCount) //IPv4JoinedSegments
	if err != nil {
		return nil, err
	}
	notJoinedCount := totalCount - 1
	segs := make([]*AddressDivision, totalCount)
	section.copySubDivisions(0, notJoinedCount, segs)
	segs[notJoinedCount] = joinedSegment
	equivalentPart := createInitializedGrouping(segs, section.getPrefixLen())
	return equivalentPart, nil
}

func (section *IPv4AddressSection) joinSegments(joinCount int) (*AddressDivision, addrerr.IncompatibleAddressError) {
	var lower, upper DivInt
	var prefix PrefixLen
	var networkPrefixLength BitCount

	var firstRange *IPv4AddressSegment
	firstJoinedIndex := section.GetSegmentCount() - 1 - joinCount
	bitsPerSeg := section.GetBitsPerSegment()
	for j := 0; j <= joinCount; j++ {
		thisSeg := section.GetSegment(firstJoinedIndex + j)
		if firstRange != nil {
			if !thisSeg.IsFullRange() {
				return nil, &incompatibleAddressError{addressError{key: "ipaddress.error.invalidMixedRange"}}
			}
		} else if thisSeg.isMultiple() {
			firstRange = thisSeg
		}
		lower = (lower << uint(bitsPerSeg)) | DivInt(thisSeg.getSegmentValue())
		upper = (upper << uint(bitsPerSeg)) | DivInt(thisSeg.getUpperSegmentValue())
		if prefix == nil {
			thisSegPrefix := thisSeg.getDivisionPrefixLength()
			if thisSegPrefix != nil {
				prefix = cacheBitCount(networkPrefixLength + thisSegPrefix.bitCount())
			} else {
				networkPrefixLength += thisSeg.getBitCount()
			}
		}
	}
	return newRangePrefixDivision(lower, upper, prefix, (BitCount(joinCount)+1)<<3), nil
}

func (section *IPv4AddressSection) toNormalizedString(stringOptions addrstr.IPStringOptions) string {
	return toNormalizedIPString(stringOptions, section)
}

// String implements the [fmt.Stringer] interface, returning the normalized string provided by ToNormalizedString, or "<nil>" if the receiver is a nil pointer.
func (section *IPv4AddressSection) String() string {
	if section == nil {
		return nilString()
	}
	return section.toString()
}

// GetSegmentStrings returns a slice with the string for each segment being the string that is normalized with wildcards.
func (section *IPv4AddressSection) GetSegmentStrings() []string {
	if section == nil {
		return nil
	}
	return section.getSegmentStrings()
}

// Inet_aton_radix represents a radix for printing an address string.
type Inet_aton_radix int

// GetRadix converts the radix to an int.
func (rad Inet_aton_radix) GetRadix() int {
	return int(rad)
}

// GetSegmentStrPrefix returns the string prefix used to identify the radix.
func (rad Inet_aton_radix) GetSegmentStrPrefix() string {
	if rad == Inet_aton_radix_octal {
		return OctalPrefix
	} else if rad == Inet_aton_radix_hex {
		return HexPrefix
	}
	return ""
}

// String returns the name of the radix.
func (rad Inet_aton_radix) String() string {
	if rad == Inet_aton_radix_octal {
		return "octal"
	} else if rad == Inet_aton_radix_hex {
		return "hexadecimal"
	}
	return "decimal"
}

const (
	Inet_aton_radix_octal   Inet_aton_radix = 8
	Inet_aton_radix_hex     Inet_aton_radix = 16
	Inet_aton_radix_decimal Inet_aton_radix = 10
)
