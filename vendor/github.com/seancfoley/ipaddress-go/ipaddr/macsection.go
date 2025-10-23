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

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstr"
)

func createMACSection(segments []*AddressDivision) *MACAddressSection {
	return &MACAddressSection{
		addressSectionInternal{
			addressDivisionGroupingInternal{
				addressDivisionGroupingBase: addressDivisionGroupingBase{
					divisions: standardDivArray(segments),
					addrType:  macType,
					cache: &valueCache{
						stringCache: stringCache{
							macStringCache: &macStringCache{},
						},
					},
				},
			},
		},
	}
}

// NewMACSection constructs a MAC address or address collection section from the given segments.
func NewMACSection(segments []*MACAddressSegment) *MACAddressSection {
	return createMACSectionFromSegs(segments)
}

func createMACSectionFromSegs(orig []*MACAddressSegment) *MACAddressSection {
	segCount := len(orig)
	newSegs := make([]*AddressDivision, segCount)
	var newPref PrefixLen
	isMultiple := false
	if segCount != 0 {
		isBlock := true
		for i := segCount - 1; i >= 0; i-- {
			segment := orig[i]
			if segment == nil {
				segment = zeroMACSeg
				if isBlock && i != segCount-1 {
					newPref = getNetworkPrefixLen(MACBitsPerSegment, MACBitsPerSegment, i)
					isBlock = false
				}
			} else {
				if isBlock {
					minPref := segment.GetMinPrefixLenForBlock()
					if minPref > 0 {
						if minPref != MACBitsPerSegment || i != segCount-1 {
							newPref = getNetworkPrefixLen(MACBitsPerSegment, minPref, i)
						}
						isBlock = false
					}
				}
				isMultiple = isMultiple || segment.isMultiple()
			}
			newSegs[i] = segment.ToDiv()
		}
		if isBlock {
			newPref = cacheBitCount(0)
		}
	}
	res := createMACSection(newSegs)
	res.isMult = isMultiple
	res.prefixLength = newPref
	return res
}

func newMACSectionParsed(segments []*AddressDivision, isMultiple bool) (res *MACAddressSection) {
	res = createMACSection(segments)
	res.initImplicitPrefLen(MACBitsPerSegment)
	res.isMult = isMultiple
	return
}

func newMACSectionEUI(segments []*AddressDivision) (res *MACAddressSection) {
	res = createMACSection(segments)
	res.initMultAndImplicitPrefLen(MACBitsPerSegment)
	return
}

// NewMACSectionFromBytes constructs a MAC address section from the given byte slice.
// The segment count is determined by the slice length, even if the segment count exceeds 8 segments.
func NewMACSectionFromBytes(bytes []byte, segmentCount int) (res *MACAddressSection, err addrerr.AddressValueError) {
	if segmentCount < 0 {
		segmentCount = len(bytes)
	}
	expectedByteCount := segmentCount
	segments, err := toSegments(
		bytes,
		segmentCount,
		MACBytesPerSegment,
		MACBitsPerSegment,
		macNetwork.getAddressCreator(),
		nil)
	if err == nil {
		// note prefix len is nil
		res = createMACSection(segments)
		if expectedByteCount == len(bytes) {
			bytes = clone(bytes)
			res.cache.bytesCache = &bytesCache{lowerBytes: bytes}
			if !res.isMult { // not a prefix block
				res.cache.bytesCache.upperBytes = bytes
			}
		}
	}
	return
}

// NewMACSectionFromUint64 constructs a MAC address section of the given segment count from the given value.
// The least significant bits of the given value will be used.
func NewMACSectionFromUint64(val uint64, segmentCount int) (res *MACAddressSection) {
	if segmentCount < 0 {
		segmentCount = MediaAccessControlSegmentCount
	}
	segments := createSegmentsUint64(
		segmentCount,
		0,
		val,
		MACBytesPerSegment,
		MACBitsPerSegment,
		macNetwork.getAddressCreator(),
		nil)
	// note prefix len is nil
	res = createMACSection(segments)
	return
}

// NewMACSectionFromVals constructs a MAC address section of the given segment count from the given values.
func NewMACSectionFromVals(vals MACSegmentValueProvider, segmentCount int) (res *MACAddressSection) {
	res = NewMACSectionFromRange(vals, nil, segmentCount)
	return
}

// NewMACSectionFromRange constructs a MAC address collection section of the given segment count from the given values.
func NewMACSectionFromRange(vals, upperVals MACSegmentValueProvider, segmentCount int) (res *MACAddressSection) {
	if segmentCount < 0 {
		segmentCount = 0
	}
	segments, isMultiple := createSegments(
		WrapMACSegmentValueProvider(vals),
		WrapMACSegmentValueProvider(upperVals),
		segmentCount,
		MACBitsPerSegment,
		macNetwork.getAddressCreator(),
		nil)
	res = createMACSection(segments)
	if isMultiple {
		res.initImplicitPrefLen(MACBitsPerSegment)
		res.isMult = true
	}
	return
}

// MACAddressSection is a section of a MACAddress.
//
// It is a series of 0 to 8 individual MAC address segments.
type MACAddressSection struct {
	addressSectionInternal
}

// containsSame returns whether this address section contains all address sections in the given address section collection of the same type.
func (addr *MACAddressSection) containsSame(other *MACAddressSection) bool {
	return addr.Contains(other)
}

// Contains returns whether this is same type and version as the given address section and whether it contains all values in the given section.
//
// Sections must also have the same number of segments to be comparable, otherwise false is returned.
func (section *MACAddressSection) Contains(other AddressSectionType) bool {
	if section == nil {
		return other == nil || other.ToSectionBase() == nil
	}
	return section.contains(other)
}

// Overlaps returns whether this is same type and version as the given address section and whether it overlaps the given section, both sections containing at least one individual section in common.
//
// Sections must also have the same number of segments to be comparable, otherwise false is returned.
func (section *MACAddressSection) Overlaps(other AddressSectionType) bool {
	if section == nil {
		return other == nil || other.ToSectionBase() == nil
	}
	return section.overlaps(other)
}

// Equal returns whether the given address section is equal to this address section.
// Two address sections are equal if they represent the same set of sections.
// They must match:
//   - type/version: MAC
//   - segment counts
//   - segment value ranges
//
// Prefix lengths are ignored.
func (section *MACAddressSection) Equal(other AddressSectionType) bool {
	if section == nil {
		return other == nil || other.ToSectionBase() == nil
	}
	return section.equal(other)
}

// Compare returns a negative integer, zero, or a positive integer if this address section is less than, equal, or greater than the given item.
// Any address item is comparable to any other.  All address items use CountComparator to compare.
func (section *MACAddressSection) Compare(item AddressItem) int {
	return CountComparator.Compare(section, item)
}

// CompareSize compares the counts of two items, the number of individual items represented.
//
// Rather than calculating counts with GetCount, there can be more efficient ways of determining whether this section represents more individual address sections than another.
//
// CompareSize returns a positive integer if this address section has a larger count than the one given, zero if they are the same, or a negative integer if the other has a larger count.
func (section *MACAddressSection) CompareSize(other AddressItem) int {
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
func (section *MACAddressSection) GetBitsPerSegment() BitCount {
	return MACBitsPerSegment
}

// GetBytesPerSegment returns the number of bytes comprising each segment in this section.  Segments in the same address section are equal length.
func (section *MACAddressSection) GetBytesPerSegment() int {
	return MACBytesPerSegment
}

// GetCount returns the count of possible distinct values for this item.
// If not representing multiple values, the count is 1,
// unless this is a division grouping with no divisions, or an address section with no segments, in which case it is 0.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (section *MACAddressSection) GetCount() *big.Int {
	if section == nil {
		return bigZero()
	}
	return section.cacheCount(func() *big.Int {
		return count(func(index int) uint64 {
			return section.GetSegment(index).GetValueCount()
		}, section.GetSegmentCount(), 6, 0x7fffffffffffff)
	})
}

func (section *MACAddressSection) getCachedCount() *big.Int {
	if section == nil {
		return bigZero()
	}
	return section.cachedCount(func() *big.Int {
		return count(func(index int) uint64 {
			return section.GetSegment(index).GetValueCount()
		}, section.GetSegmentCount(), 6, 0x7fffffffffffff)
	})
}

// IsMultiple returns whether this section represents multiple values.
func (section *MACAddressSection) IsMultiple() bool {
	return section != nil && section.isMultiple()
}

// IsPrefixed returns whether this section has an associated prefix length.
func (section *MACAddressSection) IsPrefixed() bool {
	return section != nil && section.isPrefixed()
}

// GetPrefixCount returns the number of distinct prefix values in this item.
//
// The prefix length is given by GetPrefixLen.
//
// If this has a non-nil prefix length, returns the number of distinct prefix values.
//
// If this has a nil prefix length, returns the same value as GetCount.
func (section *MACAddressSection) GetPrefixCount() *big.Int {
	return section.cachePrefixCount(func() *big.Int {
		return section.GetPrefixCountLen(section.getPrefixLen().bitCount())
	})
}

// GetPrefixCountLen returns the number of distinct prefix values in this item for the given prefix length.
func (section *MACAddressSection) GetPrefixCountLen(prefixLen BitCount) *big.Int {
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
				segmentPrefixLength := getPrefixedSegmentPrefixLength(section.GetBitsPerSegment(), prefixLen, index)
				return getPrefixValueCount(section.GetSegment(index).ToSegmentBase(), segmentPrefixLength.bitCount())
			}
			return section.GetSegment(index).GetValueCount()
		}, networkSegmentIndex+1, 6, 0x7fffffffffffff)
	})
}

// GetBlockCount returns the count of distinct values in the given number of initial (more significant) segments.
func (section *MACAddressSection) GetBlockCount(segments int) *big.Int {
	return section.calcCount(func() *big.Int {
		return count(func(index int) uint64 {
			return section.GetSegment(index).GetValueCount()
		},
			segments, 6, 0x7fffffffffffff)
	})
}

// WithoutPrefixLen provides the same address section but with no prefix length.  The values remain unchanged.
func (section *MACAddressSection) WithoutPrefixLen() *MACAddressSection {
	if !section.IsPrefixed() {
		return section
	}
	return section.withoutPrefixLen().ToMAC()
}

// SetPrefixLen sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the address section.
// The provided prefix length will be adjusted to these boundaries if necessary.
func (section *MACAddressSection) SetPrefixLen(prefixLen BitCount) *MACAddressSection {
	return section.setPrefixLen(prefixLen).ToMAC()
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
func (section *MACAddressSection) SetPrefixLenZeroed(prefixLen BitCount) (*MACAddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.setPrefixLenZeroed(prefixLen)
	return res.ToMAC(), err
}

// AdjustPrefixLen increases or decreases the prefix length by the given increment.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the address section.
//
// If this address section has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
func (section *MACAddressSection) AdjustPrefixLen(prefixLen BitCount) *AddressSection {
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
func (section *MACAddressSection) AdjustPrefixLenZeroed(prefixLen BitCount) (*AddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.adjustPrefixLenZeroed(prefixLen)
	return res.ToSectionBase(), err
}

// AssignPrefixForSingleBlock returns the equivalent prefix block that matches exactly the range of values in this address section.
// The returned block will have an assigned prefix length indicating the prefix length for the block.
//
// There may be no such address section - it is required that the range of values match the range of a prefix block.
// If there is no such address section, then nil is returned.
func (section *MACAddressSection) AssignPrefixForSingleBlock() *MACAddressSection {
	return section.assignPrefixForSingleBlock().ToMAC()
}

// AssignMinPrefixForBlock returns an equivalent address section, assigned the smallest prefix length possible,
// such that the prefix block for that prefix length is in this address section.
//
// In other words, this method assigns a prefix length to this address section matching the largest prefix block in this address section.
func (section *MACAddressSection) AssignMinPrefixForBlock() *MACAddressSection {
	return section.assignMinPrefixForBlock().ToMAC()
}

// GetSegment returns the segment at the given index.
// The first segment is at index 0.
// GetSegment will panic given a negative index or an index matching or larger than the segment count.
func (section *MACAddressSection) GetSegment(index int) *MACAddressSegment {
	return section.getDivision(index).ToMAC()
}

// ForEachSegment visits each segment in order from most-significant to least, the most significant with index 0, calling the given function for each, terminating early if the function returns true.
// Returns the number of visited segments.
func (section *MACAddressSection) ForEachSegment(consumer func(segmentIndex int, segment *MACAddressSegment) (stop bool)) int {
	divArray := section.getDivArray()
	if divArray != nil {
		for i, div := range divArray {
			if consumer(i, div.ToMAC()) {
				return i + 1
			}
		}
	}
	return len(divArray)
}

// ToDivGrouping converts to an AddressDivisionGrouping, a polymorphic type usable with all address sections and division groupings.
// Afterwards, you can convert back with ToMAC.
//
// ToDivGrouping can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *MACAddressSection) ToDivGrouping() *AddressDivisionGrouping {
	return section.ToSectionBase().ToDivGrouping()
}

// ToSectionBase converts to an AddressSection, a polymorphic type usable with all address sections.
// Afterwards, you can convert back with ToMAC.
//
// ToSectionBase can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (section *MACAddressSection) ToSectionBase() *AddressSection {
	return (*AddressSection)(section)
}

// Wrap wraps this address section, returning a WrappedAddressSection, an implementation of ExtendedSegmentSeries,
// which can be used to write code that works with both addresses and address sections.
func (section *MACAddressSection) Wrap() WrappedAddressSection {
	return wrapSection(section.ToSectionBase())
}

// GetTrailingSection gets the subsection from the series starting from the given index.
// The first segment is at index 0.
func (section *MACAddressSection) GetTrailingSection(index int) *MACAddressSection {
	return section.GetSubSection(index, section.GetSegmentCount())
}

// GetSubSection gets the subsection from the series starting from the given index and ending just before the give endIndex.
// The first segment is at index 0.
func (section *MACAddressSection) GetSubSection(index, endIndex int) *MACAddressSection {
	return section.getSubSection(index, endIndex).ToMAC()
}

// CopySegments copies the existing segments into the given slice,
// as much as can be fit into the slice, returning the number of segments copied.
func (section *MACAddressSection) CopySegments(segs []*MACAddressSegment) (count int) {
	return section.ForEachSegment(func(index int, seg *MACAddressSegment) (stop bool) {
		if stop = index >= len(segs); !stop {
			segs[index] = seg
		}
		return
	})
}

// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
func (section *MACAddressSection) CopySubSegments(start, end int, segs []*MACAddressSegment) (count int) {
	start, end, targetStart := adjust1To1StartIndices(start, end, section.GetDivisionCount(), len(segs))
	segs = segs[targetStart:]
	return section.forEachSubDivision(start, end, func(index int, div *AddressDivision) {
		segs[index] = div.ToMAC()
	}, len(segs))
}

// GetSegments returns a slice with the address segments.  The returned slice is not backed by the same array as this section.
func (section *MACAddressSection) GetSegments() (res []*MACAddressSegment) {
	res = make([]*MACAddressSegment, section.GetSegmentCount())
	section.CopySegments(res)
	return
}

// GetLower returns the section in the range with the lowest numeric value,
// which will be the same section if it represents a single value.
// For example, for "1:1:1:2-3:4:5-6", the series "1:1:1:2:4:5" is returned.
func (section *MACAddressSection) GetLower() *MACAddressSection {
	return section.getLower().ToMAC()
}

// GetUpper returns the section in the range with the highest numeric value,
// which will be the same section if it represents a single value.
// For example, for "1:1:1:2-3:4:5-6", the series "1:1:1:3:4:6" is returned.
func (section *MACAddressSection) GetUpper() *MACAddressSection {
	return section.getUpper().ToMAC()
}

// Uint64Value returns the lowest individual address section in the address section collection as a uint64.
func (section *MACAddressSection) Uint64Value() uint64 {
	return section.getLongValue(true)
}

// UpperUint64Value returns the highest individual address section in the address section collection as a uint64.
func (section *MACAddressSection) UpperUint64Value() uint64 {
	return section.getLongValue(false)
}

func (section *MACAddressSection) getLongValue(lower bool) (result uint64) {
	segCount := section.GetSegmentCount()
	if segCount == 0 {
		return
	}
	seg := section.GetSegment(0)
	if lower {
		result = uint64(seg.GetSegmentValue())
	} else {
		result = uint64(seg.GetUpperSegmentValue())
	}
	bitsPerSegment := section.GetBitsPerSegment()
	for i := 1; i < segCount; i++ {
		result = result << uint(bitsPerSegment)
		seg = section.GetSegment(i)
		if lower {
			result |= uint64(seg.GetSegmentValue())
		} else {
			result |= uint64(seg.GetUpperSegmentValue())
		}
	}
	return
}

// ToPrefixBlock returns the section with the same prefix as this section while the remaining bits span all values.
// The returned section will be the block of all sections with the same prefix.
//
// If this section has no prefix, this section is returned.
func (section *MACAddressSection) ToPrefixBlock() *MACAddressSection {
	return section.toPrefixBlock().ToMAC()
}

// ToPrefixBlockLen returns the section with the same prefix of the given length as this section while the remaining bits span all values.
// The returned section will be the block of all sections with the same prefix.
func (section *MACAddressSection) ToPrefixBlockLen(prefLen BitCount) *MACAddressSection {
	return section.toPrefixBlockLen(prefLen).ToMAC()
}

// ToBlock creates a new block of address sections by changing the segment at the given index to have the given lower and upper value,
// and changing the following segments to be full-range.
func (section *MACAddressSection) ToBlock(segmentIndex int, lower, upper SegInt) *MACAddressSection {
	return section.toBlock(segmentIndex, lower, upper).ToMAC()
}

// Iterator provides an iterator to iterate through the individual address sections of this address section.
//
// When iterating, the prefix length is preserved.  Remove it using WithoutPrefixLen prior to iterating if you wish to drop it from all individual address sections.
//
// Call IsMultiple to determine if this instance represents multiple address sections, or GetCount for the count.
func (section *MACAddressSection) Iterator() Iterator[*MACAddressSection] {
	if section == nil {
		return macSectionIterator{nilSectIterator()}
	}
	return macSectionIterator{section.sectionIterator(nil)}
}

// PrefixIterator provides an iterator to iterate through the individual prefixes of this address section,
// each iterated element spanning the range of values for its prefix.
//
// It is similar to the prefix block iterator, except for possibly the first and last iterated elements, which might not be prefix blocks,
// instead constraining themselves to values from this address section.
//
// If the series has no prefix length, then this is equivalent to Iterator.
func (section *MACAddressSection) PrefixIterator() Iterator[*MACAddressSection] {
	return macSectionIterator{section.prefixIterator(false)}
}

// PrefixBlockIterator provides an iterator to iterate through the individual prefix blocks, one for each prefix of this address section.
// Each iterated address section will be a prefix block with the same prefix length as this address section.
//
// If this address section has no prefix length, then this is equivalent to Iterator.
func (section *MACAddressSection) PrefixBlockIterator() Iterator[*MACAddressSection] {
	return macSectionIterator{section.prefixIterator(true)}
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
func (section *MACAddressSection) IncrementBoundary(increment int64) *MACAddressSection {
	return section.incrementBoundary(increment).ToMAC()
}

// IsAdaptiveZero returns true if the division grouping was originally created as an implicitly zero-valued section or grouping (e.g. IPv4AddressSection{}),
// meaning it was not constructed using a constructor function.
// Such a grouping, which has no divisions or segments, is convertible to an implicitly zero-valued grouping of any type or version, whether IPv6, IPv4, MAC, or other.
// In other words, when a section or grouping is the zero-value, then it is equivalent and convertible to the zero value of any other section or grouping type.
func (section *MACAddressSection) IsAdaptiveZero() bool {
	return section != nil && section.matchesZeroGrouping()
}

func getMacMaxValueLong(segmentCount int) uint64 {
	return macMaxValues[segmentCount]
}

var macMaxValues = []uint64{
	0,
	MACMaxValuePerSegment,
	0xffff,
	0xffffff,
	0xffffffff,
	0xffffffffff,
	0xffffffffffff,
	0xffffffffffffff,
	0xffffffffffffffff}

func (section *MACAddressSection) getMacCountLong() uint64 {
	return section.getCachedCount().Uint64()
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
func (section *MACAddressSection) Increment(incrementVal int64) *MACAddressSection {
	if incrementVal == 0 && !section.isMultiple() {
		return section
	}
	segCount := section.GetSegmentCount()
	if isOverflow := checkOverflow(incrementVal, section.Uint64Value, section.UpperUint64Value, section.getMacCountLong, func() uint64 { return getMacMaxValueLong(segCount) }, section.IsSequential); isOverflow {
		return nil
	}
	return increment(
		section.ToSectionBase(),
		incrementVal,
		macNetwork.getAddressCreator(),
		section.getMacCountLong,
		section.Uint64Value,
		section.UpperUint64Value,
		section.addressSectionInternal.getLower,
		section.addressSectionInternal.getUpper,
		section.getPrefixLen()).ToMAC()
}

func low64MAC(section *AddressSection) uint64 {
	return section.ToMAC().Uint64Value()
}

func low64UpperMAC(section *AddressSection) uint64 {
	return section.ToMAC().UpperUint64Value()
}

func (section *MACAddressSection) enumerateAddr(other AddressSectionType) *big.Int {
	if otherSection := other.ToSectionBase(); otherSection != nil {
		if matches, count := section.matchesTypeAndCount(otherSection); matches {
			if count < 8 {
				if val, ok := enumerateSmall(section.ToSectionBase(), otherSection, low64MAC, low64UpperMAC); ok {
					return big.NewInt(val)
				}
				return nil
			}
			return enumerateBig(section.ToSectionBase(), otherSection, low64MAC, low64UpperMAC)
		}
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
// If the given address section does not have the same type, then nil is returned.
//
// Sections must also have the same number of segments to be comparable, otherwise nil is returned.
func (section *MACAddressSection) Enumerate(other AddressSectionType) *big.Int {
	if other != nil {
		return section.enumerateAddr(other)
	}
	return nil
}

// ReverseBits returns a new section with the bits reversed.  Any prefix length is dropped.
//
// If the bits within a single segment cannot be reversed because the segment represents a range,
// and reversing the segment values results in a range that is not contiguous, this returns an error.
//
// In practice this means that to be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
//
// If perByte is true, the bits are reversed within each byte, otherwise all the bits are reversed.
func (section *MACAddressSection) ReverseBits(perByte bool) (*MACAddressSection, addrerr.IncompatibleAddressError) {
	res, err := section.reverseBits(perByte)
	return res.ToMAC(), err
}

// ReverseBytes returns a new section with the bytes reversed.  Any prefix length is dropped.
func (section *MACAddressSection) ReverseBytes() *MACAddressSection {
	return section.ReverseSegments()
}

// ReverseSegments returns a new section with the segments reversed.
func (section *MACAddressSection) ReverseSegments() *MACAddressSection {
	if section.GetSegmentCount() <= 1 {
		if section.IsPrefixed() {
			return section.WithoutPrefixLen()
		}
		return section
	}
	res, _ := section.reverseSegments(
		func(i int) (*AddressSegment, addrerr.IncompatibleAddressError) {
			return section.GetSegment(i).ToSegmentBase(), nil
		},
	)
	return res.ToMAC()
}

// Append creates a new section by appending the given section to this section.
func (section *MACAddressSection) Append(other *MACAddressSection) *MACAddressSection {
	count := section.GetSegmentCount()
	return section.ReplaceLen(count, count, other, 0, other.GetSegmentCount())
}

// Insert creates a new section by inserting the given section into this section at the given index.
func (section *MACAddressSection) Insert(index int, other *MACAddressSection) *MACAddressSection {
	return section.ReplaceLen(index, index, other, 0, other.GetSegmentCount())
}

// Replace replaces the segments of this section starting at the given index with the given replacement segments.
func (section *MACAddressSection) Replace(index int, replacement *MACAddressSection) *MACAddressSection {
	return section.ReplaceLen(index, index+replacement.GetSegmentCount(), replacement, 0, replacement.GetSegmentCount())
}

// ReplaceLen replaces segments starting from startIndex and ending before endIndex with the segments starting at replacementStartIndex and
// ending before replacementEndIndex from the replacement section.
func (section *MACAddressSection) ReplaceLen(startIndex, endIndex int, replacement *MACAddressSection, replacementStartIndex, replacementEndIndex int) *MACAddressSection {
	return section.replaceLen(startIndex, endIndex, replacement.ToSectionBase(), replacementStartIndex, replacementEndIndex, macBitsToSegmentBitshift).ToMAC()
}

var (
	canonicalWildcards = new(addrstr.WildcardsBuilder).SetRangeSeparator(MacDashedSegmentRangeSeparatorStr).SetWildcard(SegmentWildcardStr).ToWildcards()

	macNormalizedParams  = new(addrstr.MACStringOptionsBuilder).SetExpandedSegments(true).ToOptions()
	macCanonicalParams   = new(addrstr.MACStringOptionsBuilder).SetSeparator(MACDashSegmentSeparator).SetExpandedSegments(true).SetWildcards(canonicalWildcards).ToOptions()
	macCompressedParams  = new(addrstr.MACStringOptionsBuilder).ToOptions()
	dottedParams         = new(addrstr.MACStringOptionsBuilder).SetSeparator(MacDottedSegmentSeparator).SetExpandedSegments(true).ToOptions()
	spaceDelimitedParams = new(addrstr.MACStringOptionsBuilder).SetSeparator(MacSpaceSegmentSeparator).SetExpandedSegments(true).ToOptions()
)

// ToHexString writes this address section as a single hexadecimal value (possibly two values if a range),
// the number of digits according to the bit count, with or without a preceding "0x" prefix.
//
// If a multiple-valued section cannot be written as a range of two values, an error is returned.
func (section *MACAddressSection) ToHexString(with0xPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toHexString(with0xPrefix)
}

// ToOctalString writes this address section as a single octal value (possibly two values if a range),
// the number of digits according to the bit count, with or without a preceding "0" prefix.
//
// If a multiple-valued section cannot be written as a single prefix block or a range of two values, an error is returned.
func (section *MACAddressSection) ToOctalString(with0Prefix bool) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toOctalString(with0Prefix)
}

// ToBinaryString writes this address section as a single binary value (possibly two values if a range),
// the number of digits according to the bit count, with or without a preceding "0b" prefix.
//
// If a multiple-valued section cannot be written as a range of two values, an error is returned.
func (section *MACAddressSection) ToBinaryString(with0bPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	return section.toBinaryString(with0bPrefix)
}

// ToCanonicalString produces a canonical string for the address section.
//
// For MAC, it uses the canonical standardized IEEE 802 MAC address representation of xx-xx-xx-xx-xx-xx.  An example is "01-23-45-67-89-ab".
// For range segments, '|' is used: "11-22-33|44-55-66".
func (section *MACAddressSection) ToCanonicalString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toCustomString(macCanonicalParams)
	}
	return cacheStr(&cache.canonicalString,
		func() string {
			return section.toCustomString(macCanonicalParams)
		})
}

// ToNormalizedString produces a normalized string for the address section.
//
// For MAC, it differs from the canonical string.  It uses the most common representation of MAC addresses: "xx:xx:xx:xx:xx:xx".  An example is "01:23:45:67:89:ab".
// For range segments, '-' is used: "11:22:33-44:55:66".
func (section *MACAddressSection) ToNormalizedString() string {
	if section == nil {
		return nilString()
	}
	cch := section.getStringCache()
	if cch == nil {
		return section.toCustomString(macNormalizedParams)
	}
	strp := &cch.normalizedMACString
	return cacheStr(strp,
		func() string {
			return section.toCustomString(macNormalizedParams)
		})
}

// ToNormalizedWildcardString produces the normalized string.
func (section *MACAddressSection) ToNormalizedWildcardString() string {
	return section.ToNormalizedString()
}

// ToCompressedString produces a short representation of this address section while remaining within the confines of standard representation(s) of the address.
//
// For MAC, it differs from the canonical string.  It produces a shorter string for the address that has no leading zeros.
func (section *MACAddressSection) ToCompressedString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toCustomString(macCompressedParams)
	}
	return cacheStr(&cache.compressedMACString,
		func() string {
			return section.toCustomString(macCompressedParams)
		})
}

// ToDottedString produces the dotted hexadecimal format "aaaa.bbbb.cccc".
func (section *MACAddressSection) ToDottedString() (string, addrerr.IncompatibleAddressError) {
	if section == nil {
		return nilString(), nil
	}
	dottedGrouping, err := section.GetDottedGrouping()
	if err != nil {
		return "", err
	}
	cache := section.getStringCache()
	if cache == nil {
		return toNormalizedString(dottedParams, dottedGrouping), nil
	}
	return cacheStrErr(&cache.dottedString,
		func() (string, addrerr.IncompatibleAddressError) {
			return toNormalizedString(dottedParams, dottedGrouping), nil
		})
}

// GetDottedGrouping returns an AddressDivisionGrouping which organizes the address section into segments of bit-length 16, rather than the more typical 8 bits per segment.
//
// If this represents a collection of MAC addresses, this returns an error when unable to join two address segments,
// the first with a range of values, into a division of the larger bit-length that represents the same set of values.
func (section *MACAddressSection) GetDottedGrouping() (*AddressDivisionGrouping, addrerr.IncompatibleAddressError) {
	segmentCount := section.GetSegmentCount()
	var newSegs []*AddressDivision
	newSegmentBitCount := section.GetBitsPerSegment() << 1
	var segIndex, newSegIndex int
	newSegmentCount := (segmentCount + 1) >> 1
	newSegs = make([]*AddressDivision, newSegmentCount)
	bitsPerSeg := section.GetBitsPerSegment()
	for segIndex+1 < segmentCount {
		segment1 := section.GetSegment(segIndex)
		segIndex++
		segment2 := section.GetSegment(segIndex)
		segIndex++
		if segment1.isMultiple() && !segment2.IsFullRange() {
			return nil, &incompatibleAddressError{addressError{key: "ipaddress.error.invalid.joined.ranges"}}
		}
		val := (segment1.GetSegmentValue() << uint(bitsPerSeg)) | segment2.GetSegmentValue()
		upperVal := (segment1.GetUpperSegmentValue() << uint(bitsPerSeg)) | segment2.GetUpperSegmentValue()
		vals := newRangeDivision(DivInt(val), DivInt(upperVal), newSegmentBitCount)
		newSegs[newSegIndex] = createAddressDivision(vals)
		newSegIndex++
	}
	if segIndex < segmentCount {
		segment := section.GetSegment(segIndex)
		val := segment.GetSegmentValue() << uint(bitsPerSeg)
		upperVal := segment.GetUpperSegmentValue() << uint(bitsPerSeg)
		vals := newRangeDivision(DivInt(val), DivInt(upperVal), newSegmentBitCount)
		newSegs[newSegIndex] = createAddressDivision(vals)
	}
	grouping := createInitializedGrouping(newSegs, section.getPrefixLen())
	return grouping, nil
}

// ToSpaceDelimitedString produces a string delimited by spaces: "aa bb cc dd ee ff".
func (section *MACAddressSection) ToSpaceDelimitedString() string {
	if section == nil {
		return nilString()
	}
	cache := section.getStringCache()
	if cache == nil {
		return section.toCustomString(spaceDelimitedParams)
	}
	return cacheStr(&cache.spaceDelimitedString,
		func() string {
			return section.toCustomString(spaceDelimitedParams)
		})
}

// ToDashedString produces a string delimited by dashes: "aa-bb-cc-dd-ee-ff".
// For range segments, '|' is used: "11-22-33|44-55-66".
// It returns the same string as ToCanonicalString.
func (section *MACAddressSection) ToDashedString() string {
	if section == nil {
		return nilString()
	}
	return section.ToCanonicalString()
}

// ToColonDelimitedString produces a string delimited by colons: "aa:bb:cc:dd:ee:ff".
// For range segments, '-' is used: "11:22:33-44:55:66".
// It returns the same string as ToNormalizedString.
func (section *MACAddressSection) ToColonDelimitedString() string {
	if section == nil {
		return nilString()
	}
	return section.ToNormalizedString()
}

// String implements the [fmt.Stringer] interface, returning the normalized string provided by ToNormalizedString, or "<nil>" if the receiver is a nil pointer.
func (section *MACAddressSection) String() string {
	if section == nil {
		return nilString()
	}
	return section.toString()
}

// GetSegmentStrings returns a slice with the string for each segment being the string that is normalized with wildcards.
func (section *MACAddressSection) GetSegmentStrings() []string {
	if section == nil {
		return nil
	}
	return section.getSegmentStrings()
}
