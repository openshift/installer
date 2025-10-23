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
	"math/bits"
	"unsafe"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstr"
)

func createIPv6Section(segments []*AddressDivision) *IPv6AddressSection {
	return &IPv6AddressSection{
		ipAddressSectionInternal{
			addressSectionInternal{
				addressDivisionGroupingInternal{
					addressDivisionGroupingBase: addressDivisionGroupingBase{
						divisions: standardDivArray(segments),
						addrType:  ipv6Type,
						cache: &valueCache{
							stringCache: stringCache{
								ipv6StringCache: &ipv6StringCache{},
								ipStringCache:   &ipStringCache{},
							},
						},
					},
				},
			},
		},
	}
}

func newIPv6Section(segments []*AddressDivision) *IPv6AddressSection {
	return createIPv6Section(segments)
}

func newIPv6SectionParsed(segments []*AddressDivision, isMultiple bool) (res *IPv6AddressSection) {
	res = createIPv6Section(segments)
	res.isMult = isMultiple
	return
}

func newIPv6SectionFromMixed(segments []*AddressDivision) (res *IPv6AddressSection) {
	res = createIPv6Section(segments)
	res.initMultiple()
	return
}

func newPrefixedIPv6SectionParsed(segments []*AddressDivision, isMultiple bool, prefixLength PrefixLen, singleOnly bool) (res *IPv6AddressSection) {
	res = createIPv6Section(segments)
	res.isMult = isMultiple
	if prefixLength != nil {
		assignPrefix(prefixLength, segments, res.ToIP(), singleOnly, false, BitCount(len(segments)<<ipv6BitsToSegmentBitshift))
	}
	return
}

// NewIPv6Section constructs an IPv6 address or subnet section from the given segments.
func NewIPv6Section(segments []*IPv6AddressSegment) *IPv6AddressSection {
	return createIPv6SectionFromSegs(segments, nil)
}

// NewIPv6PrefixedSection constructs an IPv6 address or subnet section from the given segments and prefix length.
func NewIPv6PrefixedSection(segments []*IPv6AddressSegment, prefixLen PrefixLen) *IPv6AddressSection {
	return createIPv6SectionFromSegs(segments, prefixLen)
}

func createIPv6SectionFromSegs(orig []*IPv6AddressSegment, prefLen PrefixLen) (result *IPv6AddressSection) {
	divs, newPref, isMultiple := createDivisionsFromSegs(
		func(index int) *IPAddressSegment {
			return orig[index].ToIP()
		},
		len(orig),
		ipv6BitsToSegmentBitshift,
		IPv6BitsPerSegment,
		IPv6BytesPerSegment,
		IPv6MaxValuePerSegment,
		zeroIPv6Seg.ToIP(),
		zeroIPv6SegZeroPrefix.ToIP(),
		zeroIPv6SegPrefixBlock.ToIP(),
		prefLen)
	result = createIPv6Section(divs)
	result.prefixLength = newPref
	result.isMult = isMultiple
	return result
}

// NewIPv6SectionFromBigInt creates an IPv6 address section from the given big integer,
// returning an error if the value is too large for the given number of segments.
func NewIPv6SectionFromBigInt(val *big.Int, segmentCount int) (res *IPv6AddressSection, err addrerr.AddressValueError) {
	if val.Sign() < 0 {
		err = &addressValueError{
			addressError: addressError{key: "ipaddress.error.negative"},
		}
		return
	}
	return newIPv6SectionFromWords(val.Bits(), segmentCount, nil, false)
}

// NewIPv6SectionFromPrefixedBigInt creates an IPv6 address or prefix block section from the given big integer,
// returning an error if the value is too large for the given number of segments.
func NewIPv6SectionFromPrefixedBigInt(val *big.Int, segmentCount int, prefixLen PrefixLen) (res *IPv6AddressSection, err addrerr.AddressValueError) {
	if val.Sign() < 0 {
		err = &addressValueError{
			addressError: addressError{key: "ipaddress.error.negative"},
		}
		return
	}
	return newIPv6SectionFromWords(val.Bits(), segmentCount, prefixLen, false)
}

// NewIPv6SectionFromBytes constructs an IPv6 address from the given byte slice.
// The segment count is determined by the slice length, even if the segment count exceeds 8 segments.
func NewIPv6SectionFromBytes(bytes []byte) *IPv6AddressSection {
	res, _ := newIPv6SectionFromBytes(bytes, (len(bytes)+1)>>1, nil, false)
	return res
}

// NewIPv6SectionFromSegmentedBytes constructs an IPv6 address from the given byte slice.
// It allows you to specify the segment count for the supplied bytes.
// If the slice is too large for the given number of segments, an error is returned, although leading zeros are tolerated.
func NewIPv6SectionFromSegmentedBytes(bytes []byte, segmentCount int) (res *IPv6AddressSection, err addrerr.AddressValueError) {
	return newIPv6SectionFromBytes(bytes, segmentCount, nil, false)
}

// NewIPv6SectionFromPrefixedBytes constructs an IPv6 address or prefix block from the given byte slice and prefix length.
// It allows you to specify the segment count for the supplied bytes.
// If the slice is too large for the given number of segments, an error is returned, although leading zeros are tolerated.
func NewIPv6SectionFromPrefixedBytes(bytes []byte, segmentCount int, prefixLength PrefixLen) (res *IPv6AddressSection, err addrerr.AddressValueError) {
	return newIPv6SectionFromBytes(bytes, segmentCount, prefixLength, false)
}

func newIPv6SectionFromBytes(bytes []byte, segmentCount int, prefixLength PrefixLen, singleOnly bool) (res *IPv6AddressSection, err addrerr.AddressValueError) {
	if segmentCount < 0 {
		segmentCount = (len(bytes) + 1) >> 1
	}
	expectedByteCount := segmentCount << 1
	segments, err := toSegments(
		bytes,
		segmentCount,
		IPv6BytesPerSegment,
		IPv6BitsPerSegment,
		ipv6Network.getIPAddressCreator(),
		prefixLength)
	if err == nil {
		res = createIPv6Section(segments)
		if prefixLength != nil {
			assignPrefix(prefixLength, segments, res.ToIP(), singleOnly, false, BitCount(segmentCount<<ipv6BitsToSegmentBitshift))
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

func newIPv6SectionFromWords(words []big.Word, segmentCount int, prefixLength PrefixLen, singleOnly bool) (res *IPv6AddressSection, err addrerr.AddressValueError) {
	if segmentCount < 0 {
		wordBitSize := bits.UintSize
		segmentCount = (len(words) * wordBitSize) >> 4
	}
	segments, err := toSegmentsFromWords(
		words,
		segmentCount,
		prefixLength)
	if err == nil {
		res = createIPv6Section(segments)
		if prefixLength != nil {
			assignPrefix(prefixLength, segments, res.ToIP(), singleOnly, false, BitCount(segmentCount<<ipv6BitsToSegmentBitshift))
		}
	}
	return
}

func toSegmentsFromWords(
	words []big.Word,
	segmentCount int,
	prefixLength PrefixLen) (segments []*AddressDivision, err addrerr.AddressValueError) {

	wordLen := len(words)
	wordBitSize := bits.UintSize
	segmentsPerWord := wordBitSize >> ipv6BitsToSegmentBitshift
	segments = createSegmentArray(segmentCount)
	var currentWord big.Word
	if wordLen > 0 {
		currentWord = words[0]
	}
	// start with little end
	for wordIndex, wordSegmentIndex, segmentIndex := 0, 0, segmentCount-1; ; segmentIndex-- {
		var value IPv6SegInt
		if wordIndex < wordLen {
			value = IPv6SegInt(currentWord)
			currentWord >>= uint(IPv6BitsPerSegment)
			wordSegmentIndex++
		}
		segmentPrefixLength := getSegmentPrefixLength(IPv6BitsPerSegment, prefixLength, segmentIndex)
		seg := NewIPv6PrefixedSegment(value, segmentPrefixLength)
		segments[segmentIndex] = seg.ToDiv()
		if wordSegmentIndex == segmentsPerWord {
			wordSegmentIndex = 0
			wordIndex++
			if wordIndex < wordLen {
				currentWord = words[wordIndex]
			}
		}
		if segmentIndex == 0 {
			// any remaining words should be zero
			var isErr bool
			if isErr = currentWord != 0; !isErr {
				for wordIndex++; wordIndex < wordLen; wordIndex++ {
					if isErr = words[wordIndex] != 0; isErr {
						break
					}
				}
			}
			if isErr {
				err = &addressValueError{
					addressError: addressError{key: "ipaddress.error.exceeds.size"},
					val:          int(words[wordIndex]),
				}
			}
			break
		}
	}
	return
}

// NewIPv6SectionFromUint64 constructs an IPv6 address section of the given segment count from the given values.
func NewIPv6SectionFromUint64(highBytes, lowBytes uint64, segmentCount int) (res *IPv6AddressSection) {
	return NewIPv6SectionFromPrefixedUint64(highBytes, lowBytes, segmentCount, nil)
}

// NewIPv6SectionFromPrefixedUint64 constructs an IPv6 address or prefix block section of the given segment count from the given values and prefix length.
func NewIPv6SectionFromPrefixedUint64(highBytes, lowBytes uint64, segmentCount int, prefixLength PrefixLen) (res *IPv6AddressSection) {
	if segmentCount < 0 {
		segmentCount = IPv6SegmentCount
	}
	segments := createSegmentsUint64(
		segmentCount,
		highBytes,
		lowBytes,
		IPv6BytesPerSegment,
		IPv6BitsPerSegment,
		ipv6Network.getIPAddressCreator(),
		prefixLength)
	res = createIPv6Section(segments)
	if prefixLength != nil {
		assignPrefix(prefixLength, segments, res.ToIP(), false, false, BitCount(segmentCount<<ipv6BitsToSegmentBitshift))
	} else {
		res.cache.uint128Cache = &uint128Cache{high: highBytes, low: lowBytes}
	}
	return
}

// NewIPv6SectionFromVals constructs an IPv6 address section of the given segment count from the given values.
func NewIPv6SectionFromVals(vals IPv6SegmentValueProvider, segmentCount int) (res *IPv6AddressSection) {
	res = NewIPv6SectionFromPrefixedRange(vals, nil, segmentCount, nil)
	return
}

// NewIPv6SectionFromPrefixedVals constructs an IPv6 address or prefix block section of the given segment count from the given values and prefix length.
func NewIPv6SectionFromPrefixedVals(vals IPv6SegmentValueProvider, segmentCount int, prefixLength PrefixLen) (res *IPv6AddressSection) {
	return NewIPv6SectionFromPrefixedRange(vals, nil, segmentCount, prefixLength)
}

// NewIPv6SectionFromRange constructs an IPv6 subnet section of the given segment count from the given values.
func NewIPv6SectionFromRange(vals, upperVals IPv6SegmentValueProvider, segmentCount int) (res *IPv6AddressSection) {
	res = NewIPv6SectionFromPrefixedRange(vals, upperVals, segmentCount, nil)
	return
}

// NewIPv6SectionFromPrefixedRange constructs an IPv6 subnet section of the given segment count from the given values and prefix length.
func NewIPv6SectionFromPrefixedRange(vals, upperVals IPv6SegmentValueProvider, segmentCount int, prefixLength PrefixLen) (res *IPv6AddressSection) {
	return newIPv6SectionFromPrefixedSingle(vals, upperVals, segmentCount, prefixLength, false)
}

func newIPv6SectionFromPrefixedSingle(vals, upperVals IPv6SegmentValueProvider, segmentCount int, prefixLength PrefixLen, singleOnly bool) (res *IPv6AddressSection) {
	if segmentCount < 0 {
		segmentCount = 0
	}
	segments, isMultiple := createSegments(
		WrapIPv6SegmentValueProvider(vals),
		WrapIPv6SegmentValueProvider(upperVals),
		segmentCount,
		IPv6BitsPerSegment,
		ipv6Network.getIPAddressCreator(),
		prefixLength)
	res = createIPv6Section(segments)
	res.isMult = isMultiple
	if prefixLength != nil {
		assignPrefix(prefixLength, segments, res.ToIP(), singleOnly, false, BitCount(segmentCount<<ipv6BitsToSegmentBitshift))
	}
	return
}

// NewIPv6SectionFromMAC constructs an IPv6 address section from a modified EUI-64 (Extended Unique Identifier) MAC address.
//
// If the supplied MAC address section is an 8-byte EUI-64, then it must match the required EUI-64 format of "xx-xx-ff-fe-xx-xx"
// with the "ff-fe" section in the middle.
//
// If the supplied MAC address section is a 6-byte MAC-48 or EUI-48, then the ff-fe pattern will be inserted when converting to IPv6.
//
// The constructor will toggle the MAC U/L (universal/local) bit as required with EUI-64.
//
// The error is IncompatibleAddressError when unable to join two MAC segments, at least one with ranged values, into an equivalent IPV6 segment range.
func NewIPv6SectionFromMAC(eui *MACAddress) (res *IPv6AddressSection, err addrerr.IncompatibleAddressError) {
	segments := createSegmentArray(4)
	if err = toIPv6SegmentsFromEUI(segments, 0, eui.GetSection(), nil); err != nil {
		return
	}
	res = createIPv6Section(segments)
	res.isMult = eui.isMultiple()
	return
}

// IPv6AddressSection represents a section of an IPv6 address comprising 0 to 8 IPv6 address segments.
// The zero values is a section with zero-segments.
type IPv6AddressSection struct {
	ipAddressSectionInternal
}

// containsSame returns whether this address section contains all address sections in the given address section collection of the same type.
func (addr *IPv6AddressSection) containsSame(other *IPv6AddressSection) bool {
	return addr.Contains(other)
}

// Contains returns whether this is same type and version as the given address section and whether it contains all values in the given section.
//
// Sections must also have the same number of segments to be comparable, otherwise false is returned.
func (section *IPv6AddressSection) Contains(other AddressSectionType) bool {
	if section == nil {
		return other == nil || other.ToSectionBase() == nil
	}
	return section.contains(other)
}

// Overlaps returns whether this is same type and version as the given address section and whether it overlaps the given section, both sections containing at least one individual section in common.
//
// Sections must also have the same number of segments to be comparable, otherwise false is returned.
func (section *IPv6AddressSection) Overlaps(other AddressSectionType) bool {
	if section == nil {
		return other == nil || other.ToSectionBase() == nil
	}
	return section.overlaps(other)
}

// Equal returns whether the given address section is equal to this address section.
// Two address sections are equal if they represent the same set of sections.
// They must match:
//   - type/version: IPv6
//   - segment count
//   - segment value ranges
//
// Prefix lengths are ignored.
func (section *IPv6AddressSection) Equal(other AddressSectionType) bool {
	if section == nil {
		return other == nil || other.ToSectionBase() == nil
	}
	return section.equal(other)
}

// Compare returns a negative integer, zero, or a positive integer if this address section is less than, equal, or greater than the given item.
// Any address item is comparable to any other.  All address items use CountComparator to compare.
func (section *IPv6AddressSection) Compare(item AddressItem) int {
	return CountComparator.Compare(section, item)
}

// CompareSize compares the counts of two items, the number of individual sections represented.
//
// Rather than calculating counts with GetCount, there can be more efficient ways of determining whether one item represents more individual items than another.
//
// CompareSize returns a positive integer if this address section has a larger count than the item given, zero if they are the same, or a negative integer if the other has a larger count.
func (section *IPv6AddressSection) CompareSize(other AddressItem) int {
	if section == nil {
		if isNilItem(other) {
			return 0
		}
		// we have size 0, other has size >= 1
		return -1
	}
	return section.compareSize(other)
}

// GetIPVersion returns IPv6, the IP version of this address section.
func (section *IPv6AddressSection) GetIPVersion() IPVersion {
	return IPv6
}

// GetBitsPerSegment returns the number of bits comprising each segment in this section.  Segments in the same address section are equal length.
func (section *IPv6AddressSection) GetBitsPerSegment() BitCount {
	return IPv6BitsPerSegment
}

// GetBytesPerSegment returns the number of bytes comprising each segment in this section.  Segments in the same address section are equal length.
func (section *IPv6AddressSection) GetBytesPerSegment() int {
	return IPv6BytesPerSegment
}

// GetCount returns the count of possible distinct values for this item.
// If not representing multiple values, the count is 1,
// unless this is a division grouping with no divisions, or an address section with no segments, in which case it is 0.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (section *IPv6AddressSection) GetCount() *big.Int {
	if section == nil {
		return bigZero()
	}
	return section.cacheCount(func() *big.Int {
		return count(func(index int) uint64 {
			return section.GetSegment(index).GetValueCount()
		}, section.GetSegmentCount(), 2, 0x7fffffffffff)
	})
}

func (section *IPv6AddressSection) getCachedCount() *big.Int {
	if section == nil {
		return bigZero()
	}
	return section.cachedCount(func() *big.Int {
		return count(func(index int) uint64 {
			return section.GetSegment(index).GetValueCount()
		}, section.GetSegmentCount(), 2, 0x7fffffffffff)
	})
}

// IsMultiple returns  whether this section represents multiple values.
func (section *IPv6AddressSection) IsMultiple() bool {
	return section != nil && section.isMultiple()
}

// IsPrefixed returns whether this section has an associated prefix length.
func (section *IPv6AddressSection) IsPrefixed() bool {
	return section != nil && section.isPrefixed()
}

// GetBlockCount returns the count of distinct values in the given number of initial (more significant) segments.
func (section *IPv6AddressSection) GetBlockCount(segments int) *big.Int {
	return section.calcCount(func() *big.Int {
		return count(func(index int) uint64 {
			return section.GetSegment(index).GetValueCount()
		}, segments, 2, 0x7fffffffffff)
	})
}

// GetPrefixCount returns the number of distinct prefix values in this item.
//
// The prefix length is given by GetPrefixLen.
//
// If this has a non-nil prefix length, returns the number of distinct prefix values.
//
// If this has a nil prefix length, returns the same value as GetCount.
func (section *IPv6AddressSection) GetPrefixCount() *big.Int {
	return section.cachePrefixCount(func() *big.Int {
		return section.GetPrefixCountLen(section.getPrefixLen().bitCount())
	})
}

// GetPrefixCountLen returns the number of distinct prefix values in this item for the given prefix length.
func (section *IPv6AddressSection) GetPrefixCountLen(prefixLen BitCount) *big.Int {
	if prefixLen <= 0 {
		return bigOne()
	} else if bc := section.GetBitCount(); prefixLen >= bc {
		return section.GetCount()
	}
	networkSegmentIndex := getNetworkSegmentIndex(prefixLen, section.GetBytesPerSegment(), section.GetBitsPerSegment())
	hostSegmentIndex := getHostSegmentIndex(prefixLen, section.GetBytesPerSegment(), section.GetBitsPerSegment())
	return section.calcCount(func() *big.Int {
		return count(func(index int) uint64 {
			if (networkSegmentIndex == hostSegmentIndex) && index == networkSegmentIndex {
				return section.GetSegment(index).GetPrefixValueCount()
			}
			return section.GetSegment(index).GetValueCount()
		},
			networkSegmentIndex+1,
			2,
			0x7fffffffffff)
	})
}

// GetSegment returns the segment at the given index.
// The first segment is at index 0.
// GetSegment will panic given a negative index or an index matching or larger than the segment count.
func (section *IPv6AddressSection) GetSegment(index int) *IPv6AddressSegment {
	return section.getDivision(index).ToIPv6()
}

// ForEachSegment visits each segment in order from most-significant to least, the most significant with index 0, calling the given function for each, terminating early if the function returns true.
// Returns the number of visited segments.
func (section *IPv6AddressSection) ForEachSegment(consumer func(segmentIndex int, segment *IPv6AddressSegment) (stop bool)) int {
	divArray := section.getDivArray()
	if divArray != nil {
		for i, div := range divArray {
			if consumer(i, div.ToIPv6()) {
				return i + 1
			}
		}
	}
	return len(divArray)
}

// GetTrailingSection gets the subsection from the series starting from the given index.
// The first segment is at index 0.
func (section *IPv6AddressSection) GetTrailingSection(index int) *IPv6AddressSection {
	return section.GetSubSection(index, section.GetSegmentCount())
}

// GetSubSection gets the subsection from the series starting from the given index and ending just before the give endIndex.
// The first segment is at index 0.
func (section *IPv6AddressSection) GetSubSection(index, endIndex int) *IPv6AddressSection {
	return section.getSubSection(index, endIndex).ToIPv6()
}

// GetNetworkSection returns a subsection containing the segments with the network bits of the section.
// The returned section will have only as many segments as needed as determined by the existing CIDR network prefix length.
//
// If this series has no CIDR prefix length, the returned network section will
// be the entire series as a prefixed section with prefix length matching the address bit length.
func (section *IPv6AddressSection) GetNetworkSection() *IPv6AddressSection {
	return section.getNetworkSection().ToIPv6()
}

// GetNetworkSectionLen returns a subsection containing the segments with the network of the address section, the prefix bits according to the given prefix length.
// The returned section will have only as many segments as needed to contain the network.
//
// The new section will be assigned the given prefix length,
// unless the existing prefix length is smaller, in which case the existing prefix length will be retained.
func (section *IPv6AddressSection) GetNetworkSectionLen(prefLen BitCount) *IPv6AddressSection {
	return section.getNetworkSectionLen(prefLen).ToIPv6()
}

// GetHostSection returns a subsection containing the segments with the host of the address section, the bits beyond the CIDR network prefix length.
// The returned section will have only as many segments as needed to contain the host.
//
// If this series has no prefix length, the returned host section will be the full section.
func (section *IPv6AddressSection) GetHostSection() *IPv6AddressSection {
	return section.getHostSection().ToIPv6()
}

// GetHostSectionLen returns a subsection containing the segments with the host of the address section, the bits beyond the given CIDR network prefix length.
// The returned section will have only as many segments as needed to contain the host.
// The returned section will have an assigned prefix length indicating the beginning of the host.
func (section *IPv6AddressSection) GetHostSectionLen(prefLen BitCount) *IPv6AddressSection {
	return section.getHostSectionLen(prefLen).ToIPv6()
}

// GetNetworkMask returns the network mask associated with the CIDR network prefix length of this address section.
// If this section has no prefix length, then the all-ones mask is returned.
func (section *IPv6AddressSection) GetNetworkMask() *IPv6AddressSection {
	return section.getNetworkMask(ipv6Network).ToIPv6()
}

// GetHostMask returns the host mask associated with the CIDR network prefix length of this address section.
// If this section has no prefix length, then the all-ones mask is returned.
func (section *IPv6AddressSection) GetHostMask() *IPv6AddressSection {
	return section.getHostMask(ipv6Network).ToIPv6()
}

// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
func (section *IPv6AddressSection) CopySubSegments(start, end int, segs []*IPv6AddressSegment) (count int) {
	start, end, targetStart := adjust1To1StartIndices(start, end, section.GetDivisionCount(), len(segs))
	segs = segs[targetStart:]
	return section.forEachSubDivision(start, end, func(index int, div *AddressDivision) {
		segs[index] = div.ToIPv6()
	}, len(segs))
}

// CopySegments copies the existing segments into the given slice,
// as much as can be fit into the slice, returning the number of segments copied.
func (section *IPv6AddressSection) CopySegments(segs []*IPv6AddressSegment) (count int) {
	return section.ForEachSegment(func(index int, seg *IPv6AddressSegment) (stop bool) {
		if stop = index >= len(segs); !stop {
			segs[index] = seg
		}
		return
	})
}

// GetSegments returns a slice with the address segments.  The returned slice is not backed by the same array as this section.
func (section *IPv6AddressSection) GetSegments() (res []*IPv6AddressSegment) {
	res = make([]*IPv6AddressSegment, section.GetSegmentCount())
	section.CopySegments(res)
	return
}

// Mask applies the given mask to all address sections represented by this secction, returning the result.
//
// If the sections do not have a comparable number of segments, an error is returned.
//
// If this represents multiple addresses, and applying the mask to all addresses creates a set of addresses
// that cannot be represented as a sequential range within each segment, then an error is returned.
func (section *IPv6AddressSection) Mask(other *IPv6AddressSection) (res *IPv6AddressSection, err addrerr.IncompatibleAddressError) {
	return section.maskPrefixed(other, true)
}

func (section *IPv6AddressSection) maskPrefixed(other *IPv6AddressSection, retainPrefix bool) (res *IPv6AddressSection, err addrerr.IncompatibleAddressError) {
	sec, err := section.mask(other.ToIP(), retainPrefix)
	if err == nil {
		res = sec.ToIPv6()
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
func (section *IPv6AddressSection) BitwiseOr(other *IPv6AddressSection) (res *IPv6AddressSection, err addrerr.IncompatibleAddressError) {
	return section.bitwiseOrPrefixed(other, true)
}

func (section *IPv6AddressSection) bitwiseOrPrefixed(other *IPv6AddressSection, retainPrefix bool) (res *IPv6AddressSection, err addrerr.IncompatibleAddressError) {
	sec, err := section.bitwiseOr(other.ToIP(), retainPrefix)
	if err == nil {
		res = sec.ToIPv6()
	}
	return
}

// MatchesWithMask applies the mask to this address section and then compares the result with the given address section,
// returning true if they match, false otherwise.  To match, both the given section and mask must have the same number of segments as this section.
func (section *IPv6AddressSection) MatchesWithMask(other *IPv6AddressSection, mask *IPv6AddressSection) bool {
	return section.matchesWithMask(other.ToIP(), mask.ToIP())
}

// Subtract subtracts the given subnet sections from this subnet section, returning an array of sections for the result (the subnet sections will not be contiguous so an array is required).
//
// Subtract  computes the subnet difference, the set of address sections in this address section but not in the provided section.
// This is also known as the relative complement of the given argument in this subnet section.
//
// This is set subtraction, not subtraction of values.
func (section *IPv6AddressSection) Subtract(other *IPv6AddressSection) (res []*IPv6AddressSection, err addrerr.SizeMismatchError) {
	sections, err := section.subtract(other.ToIP())
	if err == nil {
		res = cloneTo(sections, (*IPAddressSection).ToIPv6)
	}
	return
}

// Intersect returns the subnet sections whose individual sections are found in both this and the given subnet section argument, or nil if no such sections exist.
//
// This is also known as the conjunction of the two sets of address sections.
//
// If the two sections have different segment counts, an error is returned.
func (section *IPv6AddressSection) Intersect(other *IPv6AddressSection) (res *IPv6AddressSection, err addrerr.SizeMismatchError) {
	sec, err := section.intersect(other.ToIP())
	if err == nil {
		res = sec.ToIPv6()
	}
	return
}

// GetLower returns the section in the range with the lowest numeric value,
// which will be the same section if it represents a single value.
// For example, for "1::1:2-3:4:5-6", the section "1::1:2:4:5" is returned.
func (section *IPv6AddressSection) GetLower() *IPv6AddressSection {
	return section.getLower().ToIPv6()
}

// GetUpper returns the section in the range with the highest numeric value,
// which will be the same section if it represents a single value.
// For example, for "1::1:2-3:4:5-6", the section "1::1:3:4:6" is returned.
func (section *IPv6AddressSection) GetUpper() *IPv6AddressSection {
	return section.getUpper().ToIPv6()
}

// Uint64Values returns the lowest address in the address section range as a pair of uint64s.
func (section *IPv6AddressSection) Uint64Values() (high, low uint64) {
	cache := section.cache
	if cache == nil {
		return section.uint64Values()
	}
	res := (*uint128Cache)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache.uint128Cache))))
	if res == nil {
		val := uint128Cache{}
		val.high, val.low = section.uint64Values()
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.uint128Cache))
		atomicStorePointer(dataLoc, unsafe.Pointer(&val))
		return val.high, val.low
	}
	return res.high, res.low
}

// uint64Values returns the lowest address in the address range as a pair of uint64 values.
func (section *IPv6AddressSection) uint64Values() (high, low uint64) {
	segCount := section.GetSegmentCount()
	if segCount == 0 {
		return
	}
	arr := section.getDivArray()
	bitsPerSegment := section.GetBitsPerSegment()
	if segCount <= 4 {
		low = uint64(arr[0].getDivisionValue())
		for i := 1; i < segCount; i++ {
			low = (low << uint(bitsPerSegment)) | uint64(arr[i].getDivisionValue())
		}
	} else {
		high = uint64(arr[0].getDivisionValue())
		highCount := segCount - 4
		i := 1
		for ; i < highCount; i++ {
			high = (high << uint(bitsPerSegment)) | uint64(arr[i].getDivisionValue())
		}
		low = uint64(arr[i].getDivisionValue())
		for i++; i < segCount; i++ {
			low = (low << uint(bitsPerSegment)) | uint64(arr[i].getDivisionValue())
		}
	}
	return
}

// UpperUint64Values returns the highest address in the address section range as pair of uint64 values.
func (section *IPv6AddressSection) UpperUint64Values() (high, low uint64) {
	if !section.IsMultiple() {
		return section.Uint64Values()
	}
	segCount := section.GetSegmentCount()
	if segCount == 0 {
		return
	}
	arr := section.getDivArray()
	bitsPerSegment := section.GetBitsPerSegment()
	if segCount <= 4 {
		low = uint64(arr[0].getUpperDivisionValue())
		for i := 1; i < segCount; i++ {
			low = (low << uint(bitsPerSegment)) | uint64(arr[i].getUpperDivisionValue())
		}
	} else {
		high = uint64(arr[0].getUpperDivisionValue())
		highCount := segCount - 4
		i := 1
		for ; i < highCount; i++ {
			high = (high << uint(bitsPerSegment)) | uint64(arr[i].getUpperDivisionValue())
		}
		low = uint64(arr[i].getUpperDivisionValue())
		for i++; i < segCount; i++ {
			low = (low << uint(bitsPerSegment)) | uint64(arr[i].getUpperDivisionValue())
		}
	}
	return
}

// ToZeroHost converts the address section to one in which all individual address sections have a host of zero,
// the host being the bits following the prefix length.
// If the address section has no prefix length, then it returns an all-zero address section.
//
// The returned section will have the same prefix and prefix length.
//
// This returns an error if the section is a range of address sections which cannot be converted to a range in which all sections have zero hosts,
// because the conversion results in a segment that is not a sequential range of values.
func (section *IPv6AddressSection) ToZeroHost() (*IPv6AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.toZeroHost(false)
	return res.ToIPv6(), err
}

// ToZeroHostLen converts the address section to one in which all individual sections have a host of zero,
// the host being the bits following the given prefix length.
// If this address section has the same prefix length, then the returned one will too, otherwise the returned section will have no prefix length.
//
// This returns an error if the section is a range of which cannot be converted to a range in which all sections have zero hosts,
// because the conversion results in a segment that is not a sequential range of values.
func (section *IPv6AddressSection) ToZeroHostLen(prefixLength BitCount) (*IPv6AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.toZeroHostLen(prefixLength)
	return res.ToIPv6(), err
}

// ToZeroNetwork converts the address section to one in which all individual address sections have a network of zero,
// the network being the bits within the prefix length.
// If the address section has no prefix length, then it returns an all-zero address section.
//
// The returned address section will have the same prefix length.
func (section *IPv6AddressSection) ToZeroNetwork() *IPv6AddressSection {
	return section.toZeroNetwork().ToIPv6()
}

// ToMaxHost converts the address section to one in which all individual address sections have a host of all one-bits, the max value,
// the host being the bits following the prefix length.
// If the address section has no prefix length, then it returns an all-ones section, the max address section.
//
// The returned address section will have the same prefix and prefix length.
//
// This returns an error if the address section is a range of address sections which cannot be converted to a range in which all sections have max hosts,
// because the conversion results in a segment that is not a sequential range of values.
func (section *IPv6AddressSection) ToMaxHost() (*IPv6AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.toMaxHost()
	return res.ToIPv6(), err
}

// ToMaxHostLen converts the address section to one in which all individual address sections have a host of all one-bits, the max host,
// the host being the bits following the given prefix length.
// If this section has the same prefix length, then the resulting section will too, otherwise the resulting section will have no prefix length.
//
// For instance, the zero host of "1.2.3.4" for the prefix length of 16 is the address "1.2.255.255".
//
// This returns an error if the section is a range of address sections which cannot be converted to a range in which all address sections have max hosts,
// because the conversion results in a segment that is not a sequential range of values.
func (section *IPv6AddressSection) ToMaxHostLen(prefixLength BitCount) (*IPv6AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.toMaxHostLen(prefixLength)
	return res.ToIPv6(), err
}

// ToPrefixBlock returns the section with the same prefix as this section while the remaining bits span all values.
// The returned section will be the block of all sections with the same prefix.
//
// If this section has no prefix, this section is returned.
func (section *IPv6AddressSection) ToPrefixBlock() *IPv6AddressSection {
	return section.toPrefixBlock().ToIPv6()
}

// ToPrefixBlockLen returns the section with the same prefix of the given length as this section while the remaining bits span all values.
// The returned section will be the block of all sections with the same prefix.
func (section *IPv6AddressSection) ToPrefixBlockLen(prefLen BitCount) *IPv6AddressSection {
	return section.toPrefixBlockLen(prefLen).ToIPv6()
}

// ToBlock creates a new block of address sections by changing the segment at the given index to have the given lower and upper value,
// and changing the following segments to be full-range.
func (section *IPv6AddressSection) ToBlock(segmentIndex int, lower, upper SegInt) *IPv6AddressSection {
	return section.toBlock(segmentIndex, lower, upper).ToIPv6()
}

// WithoutPrefixLen provides the same address section but with no prefix length.  The values remain unchanged.
func (section *IPv6AddressSection) WithoutPrefixLen() *IPv6AddressSection {
	if !section.IsPrefixed() {
		return section
	}
	return section.withoutPrefixLen().ToIPv6()
}

// SetPrefixLen sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the address section.
// The provided prefix length will be adjusted to these boundaries if necessary.
func (section *IPv6AddressSection) SetPrefixLen(prefixLen BitCount) *IPv6AddressSection {
	return section.setPrefixLen(prefixLen).ToIPv6()
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
func (section *IPv6AddressSection) SetPrefixLenZeroed(prefixLen BitCount) (*IPv6AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.setPrefixLenZeroed(prefixLen)
	return res.ToIPv6(), err
}

// AdjustPrefixLen increases or decreases the prefix length by the given increment.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the address section.
//
// If this address section has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
func (section *IPv6AddressSection) AdjustPrefixLen(prefixLen BitCount) *IPv6AddressSection {
	return section.adjustPrefixLen(prefixLen).ToIPv6()
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
func (section *IPv6AddressSection) AdjustPrefixLenZeroed(prefixLen BitCount) (*IPv6AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.adjustPrefixLenZeroed(prefixLen)
	return res.ToIPv6(), err
}

// AssignPrefixForSingleBlock returns the equivalent prefix block that matches exactly the range of values in this address section.
// The returned block will have an assigned prefix length indicating the prefix length for the block.
//
// There may be no such address section - it is required that the range of values match the range of a prefix block.
// If there is no such address section, then nil is returned.
func (section *IPv6AddressSection) AssignPrefixForSingleBlock() *IPv6AddressSection {
	return section.assignPrefixForSingleBlock().ToIPv6()
}

// AssignMinPrefixForBlock returns an equivalent address section, assigned the smallest prefix length possible,
// such that the prefix block for that prefix length is in this address section.
//
// In other words, this method assigns a prefix length to this address section matching the largest prefix block in this address section.
func (section *IPv6AddressSection) AssignMinPrefixForBlock() *IPv6AddressSection {
	return section.assignMinPrefixForBlock().ToIPv6()
}

// Iterator provides an iterator to iterate through the individual address sections of this address section.
//
// When iterating, the prefix length is preserved.  Remove it using WithoutPrefixLen prior to iterating if you wish to drop it from all individual address sections.
//
// Call IsMultiple to determine if this instance represents multiple address sections, or GetCount for the count.
func (section *IPv6AddressSection) Iterator() Iterator[*IPv6AddressSection] {
	if section == nil {
		return ipv6SectionIterator{nilSectIterator()}
	}
	return ipv6SectionIterator{section.sectionIterator(nil)}
}

// PrefixIterator provides an iterator to iterate through the individual prefixes of this address section,
// each iterated element spanning the range of values for its prefix.
//
// It is similar to the prefix block iterator, except for possibly the first and last iterated elements, which might not be prefix blocks,
// instead constraining themselves to values from this address section.
//
// If the series has no prefix length, then this is equivalent to Iterator.
func (section *IPv6AddressSection) PrefixIterator() Iterator[*IPv6AddressSection] {
	return ipv6SectionIterator{section.prefixIterator(false)}
}

// PrefixBlockIterator provides an iterator to iterate through the individual prefix blocks, one for each prefix of this address section.
// Each iterated address section will be a prefix block with the same prefix length as this address section.
//
// If this address section has no prefix length, then this is equivalent to Iterator.
func (section *IPv6AddressSection) PrefixBlockIterator() Iterator[*IPv6AddressSection] {
	return ipv6SectionIterator{section.prefixIterator(true)}
}

// BlockIterator Iterates through the address sections that can be obtained by iterating through all the upper segments up to the given segment count.
// The segments following remain the same in all iterated sections.
func (section *IPv6AddressSection) BlockIterator(segmentCount int) Iterator[*IPv6AddressSection] {
	return ipv6SectionIterator{section.blockIterator(segmentCount)}
}

// SequentialBlockIterator iterates through the sequential address sections that make up this address section.
//
// Practically, this means finding the count of segments for which the segments that follow are not full range, and then using BlockIterator with that segment count.
//
// Use GetSequentialBlockCount to get the number of iterated elements.
func (section *IPv6AddressSection) SequentialBlockIterator() Iterator[*IPv6AddressSection] {
	return ipv6SectionIterator{section.sequentialBlockIterator()}
}

// GetZeroSegments returns the list of consecutive zero-segments.
// Each element in the list will be an segment index and a total segment count for which
// that count of consecutive segments starting from that index are all zero.
func (section *IPv6AddressSection) GetZeroSegments() SegmentSequenceList {
	return section.getZeroSegments(false)
}

// GetZeroRangeSegments returns the list of consecutive zero and zero prefix block segments.
// Each element in the list will be an segment index and a total segment count for which
// that count of consecutive segments starting from that index are all zero or a prefix block segment with lowest segment value zero.
func (section *IPv6AddressSection) GetZeroRangeSegments() SegmentSequenceList {
	if section.IsPrefixed() {
		return section.getZeroSegments(true)
	}
	return section.getZeroSegments(false)
}

// GetCompressIndexAndCount chooses a single segment to be compressed in an IPv6 string. If no segment could be chosen then count is 0.
// If options is nil, no segment will be chosen.  If createMixed is true, will assume the address string will be mixed IPv6/v4.
func (section *IPv6AddressSection) getCompressIndexAndCount(options addrstr.CompressOptions, createMixed bool) (maxIndex, maxCount int) {
	if options != nil {
		rangeSelection := options.GetCompressionChoiceOptions()
		var compressibleSegs SegmentSequenceList
		if rangeSelection.CompressHost() {
			compressibleSegs = section.GetZeroRangeSegments()
		} else {
			compressibleSegs = section.GetZeroSegments()
		}
		maxCount = 0
		segmentCount := section.GetSegmentCount()
		//compressMixed := createMixed && options.GetMixedCompressionOptions().compressMixed(section)
		compressMixed := createMixed && compressMixedSect(options.GetMixedCompressionOptions(), section)
		preferHost := rangeSelection == addrstr.HostPreferred
		preferMixed := createMixed && (rangeSelection == addrstr.MixedPreferred)
		for i := compressibleSegs.size() - 1; i >= 0; i-- {
			rng := compressibleSegs.getRange(i)
			index := rng.index
			count := rng.length
			if createMixed {
				//so here we shorten the range to exclude the mixed part if necessary
				mixedIndex := IPv6MixedOriginalSegmentCount
				if !compressMixed ||
					index > mixedIndex || index+count < segmentCount { //range does not include entire mixed part.  We never compress only part of a mixed part.
					//the compressible range must stop at the mixed part
					if val := mixedIndex - index; val < count {
						count = val
					}
				}
			}
			//select this range if is the longest
			if count > 0 && count >= maxCount && (options.CompressSingle() || count > 1) {
				maxIndex = index
				maxCount = count
			}
			if preferHost && section.IsPrefixed() &&
				(BitCount(index+count)*section.GetBitsPerSegment()) > section.getNetworkPrefixLen().bitCount() { //this range contains the host
				//Since we are going backwards, this means we select as the maximum any zero-segment that includes the host
				break
			}
			if preferMixed && index+count >= segmentCount { //this range contains the mixed section
				//Since we are going backwards, this means we select to compress the mixed segment
				break
			}
		}
	}
	return
}

func compressMixedSect(m addrstr.MixedCompressionOptions, addressSection *IPv6AddressSection) bool {
	switch m {
	case addrstr.AllowMixedCompression:
		return true
	case addrstr.NoMixedCompression:
		return false
	case addrstr.MixedCompressionNoHost:
		return !addressSection.IsPrefixed()
	case addrstr.MixedCompressionCoveredByHost:
		if addressSection.IsPrefixed() {
			mixedDistance := IPv6MixedOriginalSegmentCount
			mixedCount := addressSection.GetSegmentCount() - mixedDistance
			if mixedCount > 0 {
				return (BitCount(mixedDistance) * addressSection.GetBitsPerSegment()) >= addressSection.getNetworkPrefixLen().bitCount()
			}
		}
		return true
	default:
		return true
	}
}

func (section *IPv6AddressSection) getZeroSegments(includeRanges bool) SegmentSequenceList {
	divisionCount := section.GetSegmentCount()
	includeRanges = includeRanges && section.IsPrefixBlock() && section.GetPrefixLen().bitCount() < section.GetBitCount()
	var currentIndex, currentCount, rangeCount int
	var ranges [IPv6SegmentCount >> 1]SegmentSequence
	if includeRanges {
		bitsPerSegment := section.GetBitsPerSegment()
		networkIndex := getNetworkSegmentIndex(section.getPrefixLen().bitCount(), section.GetBytesPerSegment(), bitsPerSegment)
		i := 0
		for ; i <= networkIndex; i++ {
			division := section.GetSegment(i)
			isCompressible := division.IsZero() ||
				(includeRanges && division.IsPrefixed() && division.isSinglePrefixBlock(0, division.getUpperDivisionValue(), division.getDivisionPrefixLength().bitCount()))
			if isCompressible {
				currentCount++
				if currentCount == 1 {
					currentIndex = i
				}
			} else if currentCount > 0 {
				ranges[rangeCount] = SegmentSequence{index: currentIndex, length: currentCount}
				rangeCount++
				currentCount = 0
			}
		}
		if currentCount > 0 {
			// add all segments past the network segment index to the current sequence
			ranges[rangeCount] = SegmentSequence{index: currentIndex, length: currentCount + divisionCount - i}
			rangeCount++
		} else if i < divisionCount {
			// all segments past the network segment index are a new sequence
			ranges[rangeCount] = SegmentSequence{index: i, length: divisionCount - i}
			rangeCount++
		} // else the very last segment was a network segment, and a prefix block segment, but the lowest segment value is not zero, eg ::100/120
	} else {
		for i := 0; i < divisionCount; i++ {
			division := section.GetSegment(i)
			if division.IsZero() {
				currentCount++
				if currentCount == 1 {
					currentIndex = i
				}
			} else if currentCount > 0 {
				ranges[rangeCount] = SegmentSequence{index: currentIndex, length: currentCount}
				rangeCount++
				currentCount = 0
			}
		}
		if currentCount > 0 {
			ranges[rangeCount] = SegmentSequence{index: currentIndex, length: currentCount}
			rangeCount++
		} else if rangeCount == 0 {
			return SegmentSequenceList{}
		}
	}
	return SegmentSequenceList{ranges[:rangeCount]}
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
func (section *IPv6AddressSection) IncrementBoundary(increment int64) *IPv6AddressSection {
	return section.incrementBoundary(increment).ToIPv6()
}

func getIPv6MaxValue(segmentCount int) *big.Int {
	return bigZero().Set(ipv6MaxValues[segmentCount])
}

var ipv6MaxValues = []*big.Int{
	bigZero(),
	bigZero().SetUint64(IPv6MaxValuePerSegment),
	bigZero().SetUint64(0xffffffff),
	bigZero().SetUint64(0xffffffffffff),
	maxInt(4),
	maxInt(5),
	maxInt(6),
	maxInt(7),
	maxInt(8),
}

func maxInt(segCount int) *big.Int {
	res := bigZero().SetUint64(1)
	return res.Lsh(res, 16*uint(segCount)).Sub(res, bigOneConst())
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
func (section *IPv6AddressSection) Increment(increment int64) *IPv6AddressSection {
	if increment == 0 && !section.isMultiple() {
		return section
	}
	var bigIncrement big.Int
	bigIncrement.SetInt64(increment)
	if isOverflow := checkOverflowBig(increment, &bigIncrement, section.GetValue, section.GetUpperValue, section.GetCount, func() *big.Int { return getIPv6MaxValue(section.GetSegmentCount()) }, section.IsSequential); isOverflow {
		return nil
	}
	prefixLength := section.getPrefixLen()
	result := fastIncrement(
		section.ToSectionBase(),
		increment,
		ipv6Network.getIPAddressCreator(),
		section.getLower,
		section.getUpper,
		prefixLength)
	if result != nil {
		return result.ToIPv6()
	}
	return incrementBig(
		section.ToSectionBase(),
		increment,
		&bigIncrement,
		ipv6Network.getIPAddressCreator(),
		section.getLower,
		section.getUpper,
		prefixLength).ToIPv6()
}

// IncrementBig increments the address or subnet.  It is the same as Increment but allows for a larger increment value.
// See Increment for more details.
func (section *IPv6AddressSection) IncrementBig(bigIncrement *big.Int) *IPv6AddressSection {
	if bigIsZero(bigIncrement) && !section.IsMultiple() {
		return section
	}
	if isOverflow := checkOverflowBigger(bigIncrement, section.GetValue, section.GetUpperValue, section.GetCount, func() *big.Int { return getIPv6MaxValue(section.GetSegmentCount()) }, section.IsSequential); isOverflow {
		return nil
	}
	return incrementBigger(
		section.ToSectionBase(),
		bigIncrement,
		ipv6Network.getIPAddressCreator(),
		section.getLower,
		section.getUpper,
		section.getPrefixLen()).ToIPv6()
}

func low64IPv6(section *AddressSection) uint64 {
	_, low := section.ToIPv6().Uint64Values()
	return low
}

func low64UpperIPv6(section *AddressSection) uint64 {
	_, low := section.ToIPv6().UpperUint64Values()
	return low
}

func (section *IPv6AddressSection) enumerateAddr(other AddressSectionType) *big.Int {
	if otherSection := other.ToSectionBase(); otherSection.IsIPv6() {
		return enumerateBig(section.ToSectionBase(), otherSection, low64IPv6, low64UpperIPv6)
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
func (section *IPv6AddressSection) Enumerate(other AddressSectionType) *big.Int {
	if other != nil {
		if otherSection := other.ToSectionBase(); otherSection != nil {
			if matches, count := section.matchesTypeAndCount(otherSection); matches {
				if count < 4 {
					if val, ok := enumerateSmall(section.ToSectionBase(), otherSection, low64IPv6, low64UpperIPv6); ok {
						return big.NewInt(val)
					}
					return nil
				}
				return enumerateBig(section.ToSectionBase(), otherSection, low64IPv6, low64UpperIPv6)
			}
		}
	}
	return nil
}

// SpanWithPrefixBlocks returns an array of prefix blocks that spans the same set of individual address sections as this section.
//
// Unlike SpanWithPrefixBlocksTo, the result only includes blocks that are a part of this section.
func (section *IPv6AddressSection) SpanWithPrefixBlocks() []*IPv6AddressSection {
	if section.IsSequential() {
		if section.IsSinglePrefixBlock() {
			return []*IPv6AddressSection{section}
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
func (section *IPv6AddressSection) SpanWithPrefixBlocksTo(other *IPv6AddressSection) ([]*IPv6AddressSection, addrerr.SizeMismatchError) {
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
func (section *IPv6AddressSection) SpanWithSequentialBlocks() []*IPv6AddressSection {
	if section.IsSequential() {
		return []*IPv6AddressSection{section}
	}
	return spanWithSequentialBlocks(section)
}

// SpanWithSequentialBlocksTo produces the smallest slice of sequential block address sections that span from this section to the given section.
func (section *IPv6AddressSection) SpanWithSequentialBlocksTo(other *IPv6AddressSection) ([]*IPv6AddressSection, addrerr.SizeMismatchError) {
	if err := section.checkSegmentCount(other.ToIP()); err != nil {
		return nil, err
	}
	return getSpanningSequentialBlocks(section, other), nil
}

// CoverWithPrefixBlockTo returns the minimal-size prefix block section that covers all the address sections spanning from this to the given section.
//
// If the other section has a different segment count, an error is returned.
func (section *IPv6AddressSection) CoverWithPrefixBlockTo(other *IPv6AddressSection) (*IPv6AddressSection, addrerr.SizeMismatchError) {
	res, err := section.coverWithPrefixBlockTo(other.ToIP())
	return res.ToIPv6(), err
}

// CoverWithPrefixBlock returns the minimal-size prefix block that covers all the individual address sections in this section.
// The resulting block will have a larger count than this, unless this section is already a prefix block.
func (section *IPv6AddressSection) CoverWithPrefixBlock() *IPv6AddressSection {
	return section.coverWithPrefixBlock().ToIPv6()
}

func (section *IPv6AddressSection) checkSectionCounts(sections []*IPv6AddressSection) addrerr.SizeMismatchError {
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
func (section *IPv6AddressSection) MergeToSequentialBlocks(sections ...*IPv6AddressSection) ([]*IPv6AddressSection, addrerr.SizeMismatchError) {
	if err := section.checkSectionCounts(sections); err != nil {
		return nil, err
	}
	return getMergedSequentialBlocks(cloneSeries(section, sections)), nil
}

// MergeToPrefixBlocks merges this section with the list of sections to produce the smallest array of prefix blocks.
//
// The resulting slice is sorted from lowest value to highest, regardless of the size of each prefix block.
func (section *IPv6AddressSection) MergeToPrefixBlocks(sections ...*IPv6AddressSection) ([]*IPv6AddressSection, addrerr.SizeMismatchError) {
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
func (section *IPv6AddressSection) ReverseBits(perByte bool) (*IPv6AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.reverseBits(perByte)
	return res.ToIPv6(), err
}

// ReverseBytes returns a new section with the bytes reversed.  Any prefix length is dropped.
//
// If the bytes within a single segment cannot be reversed because the segment represents a range,
// and reversing the segment values results in a range that is not contiguous, then this returns an error.
//
// In practice this means that to be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
func (section *IPv6AddressSection) ReverseBytes() (*IPv6AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.reverseBytes(false)
	return res.ToIPv6(), err
}

// ReverseSegments returns a new section with the segments reversed.
func (section *IPv6AddressSection) ReverseSegments() *IPv6AddressSection {
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
	return res.ToIPv6()
}

// Append creates a new section by appending the given section to this section.
func (section *IPv6AddressSection) Append(other *IPv6AddressSection) *IPv6AddressSection {
	count := section.GetSegmentCount()
	return section.ReplaceLen(count, count, other, 0, other.GetSegmentCount())
}

// Insert creates a new section by inserting the given section into this section at the given index.
func (section *IPv6AddressSection) Insert(index int, other *IPv6AddressSection) *IPv6AddressSection {
	return section.insert(index, other.ToIP(), ipv6BitsToSegmentBitshift).ToIPv6()
}

// Replace replaces the segments of this section starting at the given index with the given replacement segments.
func (section *IPv6AddressSection) Replace(index int, replacement *IPv6AddressSection) *IPv6AddressSection {
	return section.ReplaceLen(index, index+replacement.GetSegmentCount(), replacement, 0, replacement.GetSegmentCount())
}

// ReplaceLen replaces the segments starting from startIndex and ending before endIndex with the segments starting at replacementStartIndex and
// ending before replacementEndIndex from the replacement section.
func (section *IPv6AddressSection) ReplaceLen(startIndex, endIndex int, replacement *IPv6AddressSection, replacementStartIndex, replacementEndIndex int) *IPv6AddressSection {
	return section.replaceLen(startIndex, endIndex, replacement.ToIP(), replacementStartIndex, replacementEndIndex, ipv6BitsToSegmentBitshift).ToIPv6()
}

// IsAdaptiveZero returns true if the division grouping was originally created as an implicitly zero-valued section or grouping (e.g. IPv4AddressSection{}),
// meaning it was not constructed using a constructor function.
// Such a grouping, which has no divisions or segments, is convertible to an implicitly zero-valued grouping of any type or version, whether IPv6, IPv4, MAC, or other.
// In other words, when a section or grouping is the zero-value, then it is equivalent and convertible to the zero value of any other section or grouping type.
func (section *IPv6AddressSection) IsAdaptiveZero() bool {
	return section != nil && section.matchesZeroGrouping()
}

var (
	compressAll            = new(addrstr.CompressOptionsBuilder).SetCompressSingle(true).SetCompressionChoiceOptions(addrstr.ZerosOrHost).ToOptions()
	compressMixed          = new(addrstr.CompressOptionsBuilder).SetCompressSingle(true).SetCompressionChoiceOptions(addrstr.MixedPreferred).ToOptions()
	compressAllNoSingles   = new(addrstr.CompressOptionsBuilder).SetCompressionChoiceOptions(addrstr.ZerosOrHost).ToOptions()
	compressHostPreferred  = new(addrstr.CompressOptionsBuilder).SetCompressSingle(true).SetCompressionChoiceOptions(addrstr.HostPreferred).ToOptions()
	compressZeros          = new(addrstr.CompressOptionsBuilder).SetCompressSingle(true).SetCompressionChoiceOptions(addrstr.ZerosCompression).ToOptions()
	compressZerosNoSingles = new(addrstr.CompressOptionsBuilder).SetCompressionChoiceOptions(addrstr.ZerosCompression).ToOptions()

	uncWildcards = new(addrstr.WildcardOptionsBuilder).SetWildcardOptions(addrstr.WildcardsNetworkOnly).SetWildcards(
		new(addrstr.WildcardsBuilder).SetRangeSeparator(IPv6UncRangeSeparatorStr).SetWildcard(SegmentWildcardStr).ToWildcards()).ToOptions()
	base85Wildcards = new(addrstr.WildcardsBuilder).SetRangeSeparator(AlternativeRangeSeparatorStr).ToWildcards()

	mixedParams         = new(addrstr.IPv6StringOptionsBuilder).SetMixed(true).SetCompressOptions(compressMixed).ToOptions()
	ipv6FullParams      = new(addrstr.IPv6StringOptionsBuilder).SetExpandedSegments(true).SetWildcardOptions(wildcardsRangeOnlyNetworkOnly).ToOptions()
	ipv6CanonicalParams = new(addrstr.IPv6StringOptionsBuilder).SetCompressOptions(compressAllNoSingles).ToOptions()
	uncParams           = new(addrstr.IPv6StringOptionsBuilder).SetSeparator(IPv6UncSegmentSeparator).SetZoneSeparator(IPv6UncZoneSeparatorStr).
				SetAddressSuffix(IPv6UncSuffix).SetWildcardOptions(uncWildcards).ToOptions()
	ipv6CompressedParams         = new(addrstr.IPv6StringOptionsBuilder).SetCompressOptions(compressAll).ToOptions()
	ipv6normalizedParams         = new(addrstr.IPv6StringOptionsBuilder).ToOptions()
	canonicalWildcardParams      = new(addrstr.IPv6StringOptionsBuilder).SetWildcardOptions(allWildcards).SetCompressOptions(compressZerosNoSingles).ToOptions()
	ipv6NormalizedWildcardParams = new(addrstr.IPv6StringOptionsBuilder).SetWildcardOptions(allWildcards).ToOptions()    //no compression
	ipv6SqlWildcardParams        = new(addrstr.IPv6StringOptionsBuilder).SetWildcardOptions(allSQLWildcards).ToOptions() //no compression
	wildcardCompressedParams     = new(addrstr.IPv6StringOptionsBuilder).SetWildcardOptions(allWildcards).SetCompressOptions(compressZeros).ToOptions()
	networkPrefixLengthParams    = new(addrstr.IPv6StringOptionsBuilder).SetCompressOptions(compressHostPreferred).ToOptions()

	ipv6ReverseDNSParams = new(addrstr.IPv6StringOptionsBuilder).SetReverse(true).SetAddressSuffix(IPv6ReverseDnsSuffix).
				SetSplitDigits(true).SetExpandedSegments(true).SetSeparator('.').ToOptions()
	base85Params = new(addrstr.IPStringOptionsBuilder).SetRadix(85).SetExpandedSegments(true).
			SetWildcards(base85Wildcards).SetZoneSeparator(IPv6AlternativeZoneSeparatorStr).ToOptions()
	ipv6SegmentedBinaryParams = new(addrstr.IPStringOptionsBuilder).SetRadix(2).SetSeparator(IPv6SegmentSeparator).SetSegmentStrPrefix(BinaryPrefix).
					SetExpandedSegments(true).ToOptions()
)

// String implements the [fmt.Stringer] interface, returning the normalized string provided by ToNormalizedString, or "<nil>" if the receiver is a nil pointer.
func (section *IPv6AddressSection) String() string {
	if section == nil {
		return nilString()
	}
	return section.toString()
}

// ToHexString writes this address section as a single hexadecimal value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0x" prefix.
//
// If a multiple-valued section cannot be written as a single prefix block or a range of two values, an error is returned.
func (section *IPv6AddressSection) ToHexString(with0xPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toHexString(with0xPrefix)
}

// ToOctalString writes this address section as a single octal value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0" prefix.
//
// If a multiple-valued section cannot be written as a single prefix block or a range of two values, an error is returned.
func (section *IPv6AddressSection) ToOctalString(with0Prefix bool) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toOctalString(with0Prefix)
}

// ToBinaryString writes this address section as a single binary value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0b" prefix.
//
// If a multiple-valued section cannot be written as a single prefix block or a range of two values, an error is returned.
func (section *IPv6AddressSection) ToBinaryString(with0bPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toBinaryString(with0bPrefix)
}

// ToBase85String creates the base 85 string, which is described by RFC 1924, "A Compact Representation of IPv6 Addresses".
// See https://www.rfc-editor.org/rfc/rfc1924.html
// It may be written as a range of two values if a range that is not a prefixed block.
//
// If a multiple-valued section cannot be written as a single prefix block or a range of two values, an error is returned.
func (section *IPv6AddressSection) ToBase85String() (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toBase85String(NoZone)
	}
	cacheField := &cache.base85String
	return cacheStrErr(cacheField,
		func() (string, addrerr.IncompatibleAddressError) {
			return section.toBase85String(NoZone)
		})
}

func (section *IPv6AddressSection) toBase85String(zone Zone) (string, addrerr.IncompatibleAddressError) {
	if isDual, err := section.isDualString(); err != nil {
		return "", err
	} else {
		var largeGrouping *IPAddressLargeDivisionGrouping
		if section.hasNoDivisions() {
			largeGrouping = NewIPAddressLargeDivGrouping(nil)
		} else {
			bytes := section.getBytes()
			prefLen := section.getNetworkPrefixLen()
			bitCount := section.GetBitCount()
			var div *IPAddressLargeDivision
			if isDual {
				div = NewIPAddressLargeRangePrefixDivision(bytes, section.getUpperBytes(), prefLen, bitCount, 85)
			} else {
				div = NewIPAddressLargePrefixDivision(bytes, prefLen, bitCount, 85)
			}
			largeGrouping = NewIPAddressLargeDivGrouping([]*IPAddressLargeDivision{div})
		}
		return toNormalizedIPZonedString(base85Params, largeGrouping, zone), nil
	}
}

// ToCanonicalString produces a canonical string for the address section.
//
// For IPv6, RFC 5952 describes canonical string representation.
// https://en.wikipedia.org/wiki/IPv6_address#Representation
// http://tools.ietf.org/html/rfc5952
//
// If this section has a prefix length, it will be included in the string.
func (section *IPv6AddressSection) ToCanonicalString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toCanonicalString(NoZone)
	}
	return cacheStr(&cache.canonicalString,
		func() string {
			return section.toCanonicalString(NoZone)
		})
}

// ToNormalizedString produces a normalized string for the address section.
//
// For IPv6, it differs from the canonical string.  Zero-segments are not compressed.
//
// If this section has a prefix length, it will be included in the string.
func (section *IPv6AddressSection) ToNormalizedString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toNormalizedString(NoZone)
	}
	return cacheStr(&cache.normalizedIPv6String,
		func() string {
			return section.toNormalizedString(NoZone)
		})
}

// ToCompressedString produces a short representation of this address section while remaining within the confines of standard representation(s) of the address.
//
// For IPv6, it differs from the canonical string.  It compresses the maximum number of zeros and/or host segments with the IPv6 compression notation '::'.
func (section *IPv6AddressSection) ToCompressedString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toCompressedString(NoZone)
	}
	return cacheStr(&cache.compressedIPv6String,
		func() string {
			return section.toCompressedString(NoZone)
		})
}

// This produces the mixed IPv6/IPv4 string.  It is the shortest such string (ie fully compressed).
// For some address sections with ranges of values in the IPv4 part of the address, there is no mixed string, and an error is returned.
func (section *IPv6AddressSection) toMixedString() (string, addrerr.IncompatibleAddressError) {
	cache := section.getStringCache()
	if cache == nil {
		return section.toMixedStringZoned(NoZone)
	}
	return cacheStrErr(&cache.mixedString,
		func() (string, addrerr.IncompatibleAddressError) {
			return section.toMixedStringZoned(NoZone)
		})
}

// ToNormalizedWildcardString produces a string similar to the normalized string but avoids the CIDR prefix length.
// CIDR addresses will be shown with wildcards and ranges (denoted by '*' and '-') instead of using the CIDR prefix notation.
func (section *IPv6AddressSection) ToNormalizedWildcardString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toNormalizedWildcardStringZoned(NoZone)
	}
	return cacheStr(&cache.normalizedWildcardString,
		func() string {
			return section.toNormalizedWildcardStringZoned(NoZone)
		})
}

// ToCanonicalWildcardString produces a string similar to the canonical string but avoids the CIDR prefix length.
// Address sections with a network prefix length will be shown with wildcards and ranges (denoted by '*' and '-') instead of using the CIDR prefix length notation.
// IPv6 sections will be compressed according to the canonical representation.
func (section *IPv6AddressSection) ToCanonicalWildcardString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toCanonicalWildcardStringZoned(NoZone)
	}
	return cacheStr(&cache.canonicalWildcardString,
		func() string {
			return section.toCanonicalWildcardStringZoned(NoZone)
		})
}

// ToSegmentedBinaryString writes this address section as segments of binary values preceded by the "0b" prefix.
func (section *IPv6AddressSection) ToSegmentedBinaryString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toSegmentedBinaryStringZoned(NoZone)
	}
	return cacheStr(&cache.segmentedBinaryString,
		func() string {
			return section.toSegmentedBinaryStringZoned(NoZone)
		})
}

// ToSQLWildcardString create a string similar to that from toNormalizedWildcardString except that
// it uses SQL wildcards.  It uses '%' instead of '*' and also uses the wildcard '_'.
func (section *IPv6AddressSection) ToSQLWildcardString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toSQLWildcardStringZoned(NoZone)
	}
	return cacheStr(&cache.sqlWildcardString,
		func() string {
			return section.toSQLWildcardStringZoned(NoZone)
		})
}

// ToFullString produces a string with no compressed segments and all segments of full length with leading zeros,
// which is 4 characters for IPv6 segments.
func (section *IPv6AddressSection) ToFullString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toFullStringZoned(NoZone)
	}
	return cacheStr(&cache.fullString,
		func() string {
			return section.toFullStringZoned(NoZone)
		})
}

// ToReverseDNSString generates the reverse-DNS lookup string,
// returning an error if this address section is a multiple-valued section for which the range cannot be represented.
// For "2001:db8::567:89ab" it is "b.a.9.8.7.6.5.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa".
func (section *IPv6AddressSection) ToReverseDNSString() (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toReverseDNSStringZoned(NoZone)
	}
	return cacheStrErr(&cache.reverseDNSString,
		func() (string, addrerr.IncompatibleAddressError) {
			return section.toReverseDNSStringZoned(NoZone)
		})
}

// ToPrefixLenString returns a string with a CIDR network prefix length if this address has a network prefix length.
// For IPv6, a zero host section will be compressed with "::". For IPv4 the string is equivalent to the canonical string.
func (section *IPv6AddressSection) ToPrefixLenString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toPrefixLenStringZoned(NoZone)
	}
	return cacheStr(&cache.networkPrefixLengthString,
		func() string {
			return section.toPrefixLenStringZoned(NoZone)
		})
}

// ToSubnetString produces a string with specific formats for subnets.
// The subnet string looks like "1.2.*.*" or "1:2::/16".
//
// In the case of IPv6, when a network prefix has been supplied, the prefix will be shown and the host section will be compressed with "::".
func (section *IPv6AddressSection) ToSubnetString() string {
	if section == nil {
		return nilString()
	}
	return section.ToPrefixLenString()
}

// ToCompressedWildcardString produces a string similar to ToNormalizedWildcardString, avoiding the CIDR prefix, but with full IPv6 segment compression as well, including single zero-segments.
func (section *IPv6AddressSection) ToCompressedWildcardString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toCompressedWildcardStringZoned(NoZone)
	}
	return cacheStr(&cache.compressedWildcardString,
		func() string {
			return section.toCompressedWildcardStringZoned(NoZone)
		})
}

func (section *IPv6AddressSection) toCanonicalString(zone Zone) string {
	return section.toNormalizedZonedString(ipv6CanonicalParams, zone)
}

func (section *IPv6AddressSection) toNormalizedString(zone Zone) string {
	return section.toNormalizedZonedString(ipv6normalizedParams, zone)
}

func (section *IPv6AddressSection) toCompressedString(zone Zone) string {
	return section.toNormalizedZonedString(ipv6CompressedParams, zone)
}

func (section *IPv6AddressSection) toMixedStringZoned(zone Zone) (string, addrerr.IncompatibleAddressError) {
	return section.toNormalizedMixedZonedString(mixedParams, zone)
}

func (section *IPv6AddressSection) toNormalizedWildcardStringZoned(zone Zone) string {
	return section.toNormalizedZonedString(ipv6NormalizedWildcardParams, zone)
}

func (section *IPv6AddressSection) toCanonicalWildcardStringZoned(zone Zone) string {
	return section.toNormalizedZonedString(canonicalWildcardParams, zone)
}

func (section *IPv6AddressSection) toSegmentedBinaryStringZoned(zone Zone) string {
	return section.ipAddressSectionInternal.toCustomZonedString(ipv6SegmentedBinaryParams, zone)
}

func (section *IPv6AddressSection) toSQLWildcardStringZoned(zone Zone) string {
	return section.toNormalizedZonedString(ipv6SqlWildcardParams, zone)
}

func (section *IPv6AddressSection) toFullStringZoned(zone Zone) string {
	return section.toNormalizedZonedString(ipv6FullParams, zone)
}

func (section *IPv6AddressSection) toReverseDNSStringZoned(zone Zone) (string, addrerr.IncompatibleAddressError) {
	return section.toNormalizedSplitZonedString(ipv6ReverseDNSParams, zone)
}

func (section *IPv6AddressSection) toPrefixLenStringZoned(zone Zone) string {
	return section.toNormalizedZonedString(networkPrefixLengthParams, zone)
}

func (section *IPv6AddressSection) toCompressedWildcardStringZoned(zone Zone) string {
	return section.toNormalizedZonedString(wildcardCompressedParams, zone)
}

// ToCustomString creates a customized string from this address section according to the given string option parameters.
//
// Errors can result from split digits with ranged values, or mixed IPv4/v6 with ranged values, when the segment ranges are incompatible.
func (section *IPv6AddressSection) ToCustomString(stringOptions addrstr.IPv6StringOptions) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toCustomString(stringOptions, NoZone)
}

func (section *IPv6AddressSection) toCustomString(stringOptions addrstr.IPv6StringOptions, zone Zone) (string, addrerr.IncompatibleAddressError) {
	if stringOptions.IsMixed() {
		return section.toNormalizedMixedZonedString(stringOptions, zone)
	} else if stringOptions.IsSplitDigits() {
		return section.toNormalizedSplitZonedString(stringOptions, zone)
	}
	return section.toNormalizedZonedString(stringOptions, zone), nil
}

func (section *IPv6AddressSection) toNormalizedMixedZonedString(options addrstr.IPv6StringOptions, zone Zone) (string, addrerr.IncompatibleAddressError) {
	stringParams := from(options, section)
	if stringParams.nextUncompressedIndex <= IPv6MixedOriginalSegmentCount { //the mixed section is not compressed
		mixedParams := &ipv6v4MixedParams{
			ipv6Params: stringParams,
			ipv4Params: toIPParams(options.GetIPv4Opts()),
		}
		return section.toNormalizedMixedString(mixedParams, zone)
	}
	// the mixed section is compressed
	return stringParams.toZonedString(section, zone), nil
}

func (section *IPv6AddressSection) toNormalizedZonedString(options addrstr.IPv6StringOptions, zone Zone) string {
	return from(options, section).toZonedString(section, zone)
}

func (section *IPv6AddressSection) toNormalizedSplitZonedString(options addrstr.IPv6StringOptions, zone Zone) (string, addrerr.IncompatibleAddressError) {
	return from(options, section).toZonedSplitString(section, zone)
}

func (section *IPv6AddressSection) toNormalizedMixedString(mixedParams *ipv6v4MixedParams, zone Zone) (string, addrerr.IncompatibleAddressError) {
	mixed, err := section.getMixedAddressGrouping()
	if err != nil {
		return "", err
	}
	return mixedParams.toZonedString(mixed, zone), nil
}

// GetSegmentStrings returns a slice with the string for each segment being the string that is normalized with wildcards.
func (section *IPv6AddressSection) GetSegmentStrings() []string {
	if section == nil {
		return nil
	}
	return section.getSegmentStrings()
}

// ToDivGrouping converts to an AddressDivisionGrouping, a polymorphic type usable with all address sections and division groupings.
// Afterwards, you can convert back with ToIPv6.
//
// ToDivGrouping can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *IPv6AddressSection) ToDivGrouping() *AddressDivisionGrouping {
	return section.ToSectionBase().ToDivGrouping()
}

// ToSectionBase converts to an AddressSection, a polymorphic type usable with all address sections.
// Afterwards, you can convert back with ToIPv6.
//
// ToSectionBase can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *IPv6AddressSection) ToSectionBase() *AddressSection {
	return section.ToIP().ToSectionBase()
}

// ToIP converts to an IPAddressSection, a polymorphic type usable with all IP address sections.
//
// ToIP can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *IPv6AddressSection) ToIP() *IPAddressSection {
	return (*IPAddressSection)(section)
}

func (section *IPv6AddressSection) getMixedAddressGrouping() (*IPv6v4MixedAddressGrouping, addrerr.IncompatibleAddressError) {
	cache := section.cache
	var sect *IPv6v4MixedAddressGrouping
	var mCache *mixedCache
	if cache != nil {
		mCache = (*mixedCache)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache.mixed))))
		if mCache != nil {
			sect = mCache.defaultMixedAddressSection
		}
	}
	if sect == nil {
		mixedSect, err := section.createEmbeddedIPv4AddressSection()
		if err != nil {
			return nil, err
		}
		sect = newIPv6v4MixedGrouping(
			section.createNonMixedSection(),
			mixedSect,
		)
		if cache != nil {
			mixed := &mixedCache{
				defaultMixedAddressSection: sect,
				embeddedIPv6Section:        sect.GetIPv6AddressSection(),
				embeddedIPv4Section:        sect.GetIPv4AddressSection(),
			}
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache.mixed))
			atomicStorePointer(dataLoc, unsafe.Pointer(mixed))
		}
	}
	return sect, nil
}

// Gets the IPv4 section corresponding to the lowest (least-significant) 4 bytes in the original address,
// which will correspond to between 0 and 4 bytes in this address.  Many IPv4 to IPv6 mapping schemes (but not all) use these 4 bytes for a mapped IPv4 address.
func (section *IPv6AddressSection) getEmbeddedIPv4AddressSection() (*IPv4AddressSection, addrerr.IncompatibleAddressError) {
	cache := section.cache
	if cache == nil {
		return section.createEmbeddedIPv4AddressSection()
	}
	sect, err := section.getMixedAddressGrouping()
	if err != nil {
		return nil, err
	}
	return sect.GetIPv4AddressSection(), nil
}

// GetIPv4AddressSection produces an IPv4 address section from a sequence of bytes in this IPv6 address section.
func (section *IPv6AddressSection) GetIPv4AddressSection(startByteIndex, endByteIndex int) (*IPv4AddressSection, addrerr.IncompatibleAddressError) {
	if startByteIndex == IPv6MixedOriginalSegmentCount<<1 && endByteIndex == (section.GetSegmentCount()<<1) {
		return section.getEmbeddedIPv4AddressSection()
	}
	segments := make([]*AddressDivision, endByteIndex-startByteIndex)
	i := startByteIndex
	j := 0
	bytesPerSegment := section.GetBytesPerSegment()
	if i%bytesPerSegment == 1 {
		ipv6Segment := section.GetSegment(i >> 1)
		i++
		if err := ipv6Segment.splitIntoIPv4Segments(segments, j-1); err != nil {
			return nil, err
		}
		j++
	}
	for ; i < endByteIndex; i, j = i+bytesPerSegment, j+bytesPerSegment {
		ipv6Segment := section.GetSegment(i >> 1)
		if err := ipv6Segment.splitIntoIPv4Segments(segments, j); err != nil {
			return nil, err
		}
	}
	res := createIPv4Section(segments)
	res.initMultAndPrefLen()
	return res, nil
}

func (section *IPv6AddressSection) createNonMixedSection() *EmbeddedIPv6AddressSection {
	nonMixedCount := IPv6MixedOriginalSegmentCount
	mixedCount := section.GetSegmentCount() - nonMixedCount
	var result *IPv6AddressSection
	if mixedCount <= 0 {
		result = section
	} else {
		nonMixed := make([]*AddressDivision, nonMixedCount)
		section.copySubDivisions(0, nonMixedCount, nonMixed)
		result = createIPv6Section(nonMixed)
		result.initMultAndPrefLen()
	}
	return &EmbeddedIPv6AddressSection{
		embeddedIPv6AddressSection: embeddedIPv6AddressSection{*result},
		encompassingSection:        section,
	}
}

type embeddedIPv6AddressSection struct {
	IPv6AddressSection
}

// EmbeddedIPv6AddressSection represents the initial IPv6 section of an IPv6v4MixedAddressGrouping.
type EmbeddedIPv6AddressSection struct {
	embeddedIPv6AddressSection
	encompassingSection *IPv6AddressSection
}

// IsPrefixBlock returns whether this address segment series has a prefix length and includes the block associated with its prefix length.
// If the prefix length matches the bit count, this returns true.
//
// This is different from ContainsPrefixBlock in that this method returns
// false if the series has no prefix length, or a prefix length that differs from a prefix length for which ContainsPrefixBlock returns true.
func (section *EmbeddedIPv6AddressSection) IsPrefixBlock() bool {
	ipv6Sect := section.encompassingSection
	if ipv6Sect == nil {
		ipv6Sect = zeroIPv6AddressSection
	}
	return ipv6Sect.IsPrefixBlock()
}

func (section *IPv6AddressSection) createEmbeddedIPv4AddressSection() (sect *IPv4AddressSection, err addrerr.IncompatibleAddressError) {
	nonMixedCount := IPv6MixedOriginalSegmentCount
	segCount := section.GetSegmentCount()
	mixedCount := segCount - nonMixedCount
	lastIndex := segCount - 1
	var mixed []*AddressDivision
	if mixedCount == 0 {
		mixed = []*AddressDivision{}
	} else if mixedCount == 1 {
		mixed = make([]*AddressDivision, section.GetBytesPerSegment())
		last := section.GetSegment(lastIndex)
		if err := last.splitIntoIPv4Segments(mixed, 0); err != nil {
			return nil, err
		}
	} else {
		bytesPerSeg := section.GetBytesPerSegment()
		mixed = make([]*AddressDivision, bytesPerSeg<<1)
		low := section.GetSegment(lastIndex)
		high := section.GetSegment(lastIndex - 1)
		if err := high.splitIntoIPv4Segments(mixed, 0); err != nil {
			return nil, err
		}
		if err := low.splitIntoIPv4Segments(mixed, bytesPerSeg); err != nil {
			return nil, err
		}
	}
	sect = createIPv4Section(mixed)
	sect.initMultAndPrefLen()
	return
}

func createMixedAddressGrouping(divisions []*AddressDivision, mixedCache *mixedCache) *IPv6v4MixedAddressGrouping {
	grouping := &IPv6v4MixedAddressGrouping{
		addressDivisionGroupingInternal: addressDivisionGroupingInternal{
			addressDivisionGroupingBase: addressDivisionGroupingBase{
				divisions: standardDivArray(divisions),
				addrType:  ipv6v4MixedType,
				cache:     &valueCache{mixed: mixedCache},
			},
		},
	}
	ipv6Section := mixedCache.embeddedIPv6Section
	ipv4Section := mixedCache.embeddedIPv4Section
	grouping.isMult = ipv6Section.isMultiple() || ipv4Section.isMultiple()
	if ipv6Section.IsPrefixed() {
		grouping.prefixLength = ipv6Section.getPrefixLen()
	} else if ipv4Section.IsPrefixed() {
		grouping.prefixLength = cacheBitCount(ipv6Section.GetBitCount() + ipv4Section.getPrefixLen().bitCount())
	}
	return grouping
}

func newIPv6v4MixedGrouping(ipv6Section *EmbeddedIPv6AddressSection, ipv4Section *IPv4AddressSection) *IPv6v4MixedAddressGrouping {
	ipv6Len := ipv6Section.GetSegmentCount()
	ipv4Len := ipv4Section.GetSegmentCount()
	allSegs := make([]*AddressDivision, ipv6Len+ipv4Len)
	ipv6Section.copySubDivisions(0, ipv6Len, allSegs)
	ipv4Section.copySubDivisions(0, ipv4Len, allSegs[ipv6Len:])
	grouping := createMixedAddressGrouping(allSegs, &mixedCache{
		embeddedIPv6Section: ipv6Section,
		embeddedIPv4Section: ipv4Section,
	})
	return grouping
}

// IPv6v4MixedAddressGrouping has divisions which are a mix of IPv6 and IPv4 divisions.
// It has an initial IPv6 section followed by an IPv4 section.
type IPv6v4MixedAddressGrouping struct {
	addressDivisionGroupingInternal
}

// Compare returns a negative integer, zero, or a positive integer if this address division grouping is less than, equal, or greater than the given item.
// Any address item is comparable to any other.  All address items use CountComparator to compare.
func (grouping *IPv6v4MixedAddressGrouping) Compare(item AddressItem) int {
	return CountComparator.Compare(grouping, item)
}

// CompareSize compares the counts of two items, the number of individual items represented in each.
//
// Rather than calculating counts with GetCount, there can be more efficient ways of determining whether this grouping represents more individual address groupings than another item.
//
// CompareSize returns a positive integer if this address division grouping has a larger count than the item given, zero if they are the same, or a negative integer if the other has a larger count.
func (grouping *IPv6v4MixedAddressGrouping) CompareSize(other AddressItem) int {
	if grouping == nil {
		if isNilItem(other) {
			return 0
		}
		// we have size 0, other has size >= 1
		return -1
	}
	return grouping.compareSize(other)
}

// GetCount returns the count of possible distinct values for this item.
// If not representing multiple values, the count is 1,
// unless this is a division grouping with no divisions, or an address section with no segments, in which case it is 0.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (grouping *IPv6v4MixedAddressGrouping) GetCount() *big.Int {
	if grouping == nil {
		return bigZero()
	}
	cnt := grouping.GetIPv6AddressSection().GetCount()
	return cnt.Mul(cnt, grouping.GetIPv4AddressSection().GetCount())
}

// IsMultiple returns  whether this grouping represents multiple values.
func (grouping *IPv6v4MixedAddressGrouping) IsMultiple() bool {
	return grouping != nil && grouping.isMultiple()
}

// IsPrefixed returns whether this grouping has an associated prefix length.
func (grouping *IPv6v4MixedAddressGrouping) IsPrefixed() bool {
	return grouping != nil && grouping.isPrefixed()
}

// IsAdaptiveZero returns true if the division grouping was originally created as an implicitly zero-valued section or grouping (e.g. IPv4AddressSection{}),
// meaning it was not constructed using a constructor function.
// Such a grouping, which has no divisions or segments, is convertible to an implicitly zero-valued grouping of any type or version, whether IPv6, IPv4, MAC, or other.
// In other words, when a section or grouping is the zero-value, then it is equivalent and convertible to the zero value of any other section or grouping type.
func (grouping *IPv6v4MixedAddressGrouping) IsAdaptiveZero() bool {
	return grouping != nil && grouping.matchesZeroGrouping()
}

// ToDivGrouping converts to an AddressDivisionGrouping, a polymorphic type usable with all address sections and division groupings.
// Afterwards, you can convert back with ToMixedIPv6v4.
//
// ToDivGrouping can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (grouping *IPv6v4MixedAddressGrouping) ToDivGrouping() *AddressDivisionGrouping {
	return (*AddressDivisionGrouping)(grouping)
}

// by using cached sections for zero values, we will return the same section
// pointer for repeated calls to the same zero-valued containing section

var (
	zeroEmbeddedIPv6AddressSection = &EmbeddedIPv6AddressSection{}
	zeroIPv4AddressSection         = &IPv4AddressSection{}
	zeroIPv6AddressSection         = &IPv6AddressSection{}
)

// GetIPv6AddressSection returns the initial IPv6 section of the grouping.
func (grouping *IPv6v4MixedAddressGrouping) GetIPv6AddressSection() *EmbeddedIPv6AddressSection {
	if grouping == nil {
		return nil
	}
	cache := grouping.cache
	if cache == nil { // zero-valued
		return zeroEmbeddedIPv6AddressSection
	}
	return cache.mixed.embeddedIPv6Section
}

// GetIPv4AddressSection returns the ending IPv4 section of the grouping.
func (grouping *IPv6v4MixedAddressGrouping) GetIPv4AddressSection() *IPv4AddressSection {
	if grouping == nil {
		return nil
	}
	cache := grouping.cache
	if cache == nil { // zero-valued
		return zeroIPv4AddressSection
	}
	return cache.mixed.embeddedIPv4Section
}

// String implements the [fmt.Stringer] interface,
// as a slice string with each division converted to a string by String ( ie "[ div0 div1 ...]"),
// or "<nil>" if the receiver is a nil pointer.
func (grouping *IPv6v4MixedAddressGrouping) String() string {
	if grouping == nil {
		return nilString()
	}
	// used to use grouping.toString() but decided to use a mixed string instead
	parms := mixedParams
	ipv6Sect := grouping.GetIPv6AddressSection().encompassingSection
	if ipv6Sect == nil {
		ipv6Sect = zeroIPv6AddressSection
	}
	// using ipv6Sect here instead of grouping.GetIPv6AddressSection() affects whether compression is used
	stringParams := from(parms, ipv6Sect) // I guess this would have to be on the ipv6 section - but no compress host
	params := &ipv6v4MixedParams{
		ipv6Params: stringParams,
		ipv4Params: toIPParams(parms.GetIPv4Opts()),
	}
	// using grouping.GetIPv6AddressSection() here affects the results of IsPrefixBlock to account for the IPv4 section that follows the IPv6 section
	result := params.toZonedString(grouping, NoZone)
	return result
}

// Format is intentionally the only method with non-pointer receivers.  It is not intended to be called directly, it is intended for use by the fmt package.
// When called by a function in the fmt package, nil values are detected before this method is called, avoiding a panic when calling this method.

// Format implements [fmt.Formatter] interface. It accepts the formats
//   - 'v' for the default address and section format (either the normalized or canonical string),
//   - 's' (string) for the same
//   - 'q' for a quoted string
func (grouping IPv6v4MixedAddressGrouping) Format(state fmt.State, verb rune) {
	var str string
	switch verb {
	case 's', 'v', 'q':
		str = grouping.String()
		if verb == 'q' {
			if state.Flag('#') {
				str = "`" + str + "`"
			} else {
				str = `"` + str + `"`
			}
		}
	default:
		// format not supported
		_, _ = fmt.Fprintf(state, "%%!%c(address=%s)", verb, grouping.String())
		return
	}
	_, _ = state.Write([]byte(str))
}

var ffMACSeg, feMACSeg = NewMACSegment(0xff), NewMACSegment(0xfe)

func toIPv6SegmentsFromEUI(
	segments []*AddressDivision,
	ipv6StartIndex int, // the index into the IPv6 segment array to put the MAC-based IPv6 segments
	eui *MACAddressSection, // must be full 6 or 8 mac sections
	prefixLength PrefixLen) addrerr.IncompatibleAddressError {
	euiSegmentIndex := 0
	var seg3, seg4 *MACAddressSegment
	var err addrerr.IncompatibleAddressError
	seg0 := eui.GetSegment(euiSegmentIndex)
	euiSegmentIndex++
	seg1 := eui.GetSegment(euiSegmentIndex)
	euiSegmentIndex++
	seg2 := eui.GetSegment(euiSegmentIndex)
	euiSegmentIndex++
	isExtended := eui.GetSegmentCount() == ExtendedUniqueIdentifier64SegmentCount
	if isExtended {
		seg3 = eui.GetSegment(euiSegmentIndex)
		euiSegmentIndex++
		if !seg3.matches(0xff) {
			return &incompatibleAddressError{addressError{key: "ipaddress.mac.error.not.eui.convertible"}}
		}
		seg4 = eui.GetSegment(euiSegmentIndex)
		euiSegmentIndex++
		if !seg4.matches(0xfe) {
			return &incompatibleAddressError{addressError{key: "ipaddress.mac.error.not.eui.convertible"}}
		}
	} else {
		seg3 = ffMACSeg
		seg4 = feMACSeg
	}
	seg5 := eui.GetSegment(euiSegmentIndex)
	euiSegmentIndex++
	seg6 := eui.GetSegment(euiSegmentIndex)
	euiSegmentIndex++
	seg7 := eui.GetSegment(euiSegmentIndex)
	var currentPrefix PrefixLen
	if prefixLength != nil {
		//since the prefix comes from the ipv6 section and not the MAC section, any segment prefix for the MAC section is 0 or nil
		//prefixes across segments have the pattern: nil, nil, ..., nil, 0-16, 0, 0, ..., 0
		//So if the overall prefix is 0, then the prefix of every segment is 0
		currentPrefix = cacheBitCount(0)
	}
	var seg *IPv6AddressSegment
	if seg, err = seg0.JoinAndFlip2ndBit(seg1, currentPrefix); /* only this first one gets the flipped bit */ err == nil {
		segments[ipv6StartIndex] = seg.ToDiv()
		ipv6StartIndex++
		if seg, err = seg2.Join(seg3, currentPrefix); err == nil {
			segments[ipv6StartIndex] = seg.ToDiv()
			ipv6StartIndex++
			if seg, err = seg4.Join(seg5, currentPrefix); err == nil {
				segments[ipv6StartIndex] = seg.ToDiv()
				ipv6StartIndex++
				if seg, err = seg6.Join(seg7, currentPrefix); err == nil {
					segments[ipv6StartIndex] = seg.ToDiv()
					return nil
				}
			}
		}
	}
	return err
}

// SegmentSequence represents a sequence of consecutive segments with the given length starting from the given segment index.
type SegmentSequence struct {
	index, length int
}

// SegmentSequenceList represents a list of SegmentSequence instances.
type SegmentSequenceList struct {
	ranges []SegmentSequence
}

func (list SegmentSequenceList) size() int {
	return len(list.ranges)
}

func (list SegmentSequenceList) getRange(index int) SegmentSequence {
	return list.ranges[index]
}
