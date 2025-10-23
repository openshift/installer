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
)

type IPv6SegInt = uint16
type IPv6SegmentValueProvider func(segmentIndex int) IPv6SegInt

// WrapIPv6SegmentValueProvider converts the given IPv6SegmentValueProvider to a SegmentValueProvider.
func WrapIPv6SegmentValueProvider(f IPv6SegmentValueProvider) SegmentValueProvider {
	if f == nil {
		return nil
	}
	return func(segmentIndex int) SegInt {
		return SegInt(f(segmentIndex))
	}
}

// WrapSegmentValueProviderForIPv6 converts the given SegmentValueProvider to an IPv6SegmentValueProvider.
// Values that do not fit IPv6SegInt are truncated.
func WrapSegmentValueProviderForIPv6(f SegmentValueProvider) IPv6SegmentValueProvider {
	if f == nil {
		return nil
	}
	return func(segmentIndex int) IPv6SegInt {
		return IPv6SegInt(f(segmentIndex))
	}
}

const useIPv6SegmentCache = true

type ipv6SegmentValues struct {
	value      IPv6SegInt
	upperValue IPv6SegInt
	prefLen    PrefixLen
	cache      divCache
}

func (seg *ipv6SegmentValues) getAddrType() addrType {
	return ipv6Type
}

func (seg *ipv6SegmentValues) includesZero() bool {
	return seg.value == 0
}

func (seg *ipv6SegmentValues) includesMax() bool {
	return seg.upperValue == 0xffff
}

func (seg *ipv6SegmentValues) isMultiple() bool {
	return seg.value != seg.upperValue
}

func (seg *ipv6SegmentValues) getCount() *big.Int {
	return big.NewInt(int64(seg.upperValue-seg.value) + 1)
}

func (seg *ipv6SegmentValues) getBitCount() BitCount {
	return IPv6BitsPerSegment
}

func (seg *ipv6SegmentValues) getByteCount() int {
	return IPv6BytesPerSegment
}

func (seg *ipv6SegmentValues) getValue() *BigDivInt {
	return big.NewInt(int64(seg.value))
}

func (seg *ipv6SegmentValues) getUpperValue() *BigDivInt {
	return big.NewInt(int64(seg.upperValue))
}

func (seg *ipv6SegmentValues) getDivisionValue() DivInt {
	return DivInt(seg.value)
}

func (seg *ipv6SegmentValues) getUpperDivisionValue() DivInt {
	return DivInt(seg.upperValue)
}

func (seg *ipv6SegmentValues) getDivisionPrefixLength() PrefixLen {
	return seg.prefLen
}

func (seg *ipv6SegmentValues) getSegmentValue() SegInt {
	return SegInt(seg.value)
}

func (seg *ipv6SegmentValues) getUpperSegmentValue() SegInt {
	return SegInt(seg.upperValue)
}

func (seg *ipv6SegmentValues) calcBytesInternal() (bytes, upperBytes []byte) {
	bytes = []byte{byte(seg.value >> 8), byte(seg.value)}
	if seg.isMultiple() {
		upperBytes = []byte{byte(seg.upperValue >> 8), byte(seg.upperValue)}
	} else {
		upperBytes = bytes
	}
	return
}

func (seg *ipv6SegmentValues) bytesInternal(upper bool) []byte {
	var val IPv6SegInt
	if upper {
		val = seg.upperValue
	} else {
		val = seg.value
	}
	return []byte{byte(val >> 8), byte(val)}
}

func (seg *ipv6SegmentValues) deriveNew(val, upperVal DivInt, prefLen PrefixLen) divisionValues {
	return newIPv6SegmentPrefixedValues(IPv6SegInt(val), IPv6SegInt(upperVal), prefLen)
}

func (seg *ipv6SegmentValues) derivePrefixed(prefLen PrefixLen) divisionValues {
	return newIPv6SegmentPrefixedValues(seg.value, seg.upperValue, prefLen)
}

func (seg *ipv6SegmentValues) deriveNewSeg(val SegInt, prefLen PrefixLen) divisionValues {
	return newIPv6SegmentPrefixedVal(IPv6SegInt(val), prefLen)
}

func (seg *ipv6SegmentValues) deriveNewMultiSeg(val, upperVal SegInt, prefLen PrefixLen) divisionValues {
	return newIPv6SegmentPrefixedValues(IPv6SegInt(val), IPv6SegInt(upperVal), prefLen)
}

func (seg *ipv6SegmentValues) getCache() *divCache {
	return &seg.cache
}

var _ divisionValues = &ipv6SegmentValues{}

var zeroIPv6Seg = NewIPv6Segment(0)
var zeroIPv6SegZeroPrefix = NewIPv6PrefixedSegment(0, cacheBitCount(0))
var zeroIPv6SegPrefixBlock = NewIPv6RangePrefixedSegment(0, IPv6MaxValuePerSegment, cacheBitCount(0))

// IPv6AddressSegment represents a segment of an IPv6 address.
// An IPv6 segment contains a single value or a range of sequential values, a prefix length, and it has bit length of 16 bits.
//
// Like strings, segments are immutable, which also makes them concurrency-safe.
//
// See AddressSegment for more details regarding segments.
type IPv6AddressSegment struct {
	ipAddressSegmentInternal
}

func (seg *IPv6AddressSegment) init() *IPv6AddressSegment {
	if seg.divisionValues == nil {
		return zeroIPv6Seg
	}
	return seg
}

// GetIPv6SegmentValue returns the lower value.  Same as GetSegmentValue but returned as a IPv6SegInt.
func (seg *IPv6AddressSegment) GetIPv6SegmentValue() IPv6SegInt {
	return IPv6SegInt(seg.GetSegmentValue())
}

// GetIPv6UpperSegmentValue returns the lower value.  Same as GetUpperSegmentValue but returned as a IPv6SegInt.
func (seg *IPv6AddressSegment) GetIPv6UpperSegmentValue() IPv6SegInt {
	return IPv6SegInt(seg.GetUpperSegmentValue())
}

// Contains returns whether this is same type and version as the given segment and whether it contains all values in the given segment.
func (seg *IPv6AddressSegment) Contains(other AddressSegmentType) bool {
	if seg == nil {
		return other == nil || other.ToSegmentBase() == nil
	}
	return seg.init().contains(other)
}

// Overlaps returns whether this is same type and version as the given segment and whether it overlaps with the values in the given segment.
func (seg *IPv6AddressSegment) Overlaps(other AddressSegmentType) bool {
	if seg == nil {
		return other == nil || other.ToSegmentBase() == nil
	}
	return seg.init().overlaps(other)
}

// Equal returns whether the given segment is equal to this segment.
// Two segments are equal if they match:
//   - type/version: IPv6
//   - value range
//
// Prefix lengths are ignored.
func (seg *IPv6AddressSegment) Equal(other AddressSegmentType) bool {
	if seg == nil {
		return other == nil || other.ToDiv() == nil
	}
	return seg.init().equal(other)
}

// PrefixContains returns whether the prefix values in the prefix of the given segment are also prefix values in this segment.
// It returns whether the prefix of this segment contains the prefix of the given segment.
func (seg *IPv6AddressSegment) PrefixContains(other AddressSegmentType, prefixLength BitCount) bool {
	return seg.init().ipAddressSegmentInternal.PrefixContains(other, prefixLength)
}

// PrefixEqual returns whether the prefix bits of this segment match the same bits of the given segment.
// It returns whether the two segments share the same range of prefix values using the given prefix length.
func (seg *IPv6AddressSegment) PrefixEqual(other AddressSegmentType, prefixLength BitCount) bool {
	return seg.init().ipAddressSegmentInternal.PrefixEqual(other, prefixLength)
}

// Compare returns a negative integer, zero, or a positive integer if this address segment is less than, equal, or greater than the given item.
// Any address item is comparable to any other.  All address items use CountComparator to compare.
func (seg *IPv6AddressSegment) Compare(item AddressItem) int {
	if seg != nil {
		seg = seg.init()
	}
	return CountComparator.Compare(seg, item)
}

// CompareSize compares the counts of two items, the number of individual values within.
//
// Rather than calculating counts with GetCount, there can be more efficient ways of determining whether one represents more individual values than another.
//
// CompareSize returns a positive integer if this segment has a larger count than the one given, zero if they are the same, or a negative integer if the other has a larger count.
func (seg *IPv6AddressSegment) CompareSize(other AddressItem) int {
	if seg == nil {
		if isNilItem(other) {
			return 0
		}
		// we have size 0, other has size >= 1
		return -1
	}
	return seg.init().compareSize(other)
}

// GetBitCount returns the number of bits in each value comprising this address item, which is 16.
func (seg *IPv6AddressSegment) GetBitCount() BitCount {
	return IPv6BitsPerSegment
}

// GetByteCount returns the number of bytes required for each value comprising this address item, which is 2.
func (seg *IPv6AddressSegment) GetByteCount() int {
	return IPv6BytesPerSegment
}

// GetMaxValue gets the maximum possible value for this type or version of segment, determined by the number of bits.
//
// For the highest range value of this particular segment, use GetUpperSegmentValue.
func (seg *IPv6AddressSegment) GetMaxValue() IPv6SegInt {
	return 0xffff
}

// GetLower returns a segment representing just the lowest value in the range, which will be the same segment if it represents a single value.
func (seg *IPv6AddressSegment) GetLower() *IPv6AddressSegment {
	return seg.init().getLower().ToIPv6()
}

// GetUpper returns a segment representing just the highest value in the range, which will be the same segment if it represents a single value.
func (seg *IPv6AddressSegment) GetUpper() *IPv6AddressSegment {
	return seg.init().getUpper().ToIPv6()
}

// IsMultiple returns whether this segment represents multiple values.
func (seg *IPv6AddressSegment) IsMultiple() bool {
	return seg != nil && seg.isMultiple()
}

// GetCount returns the count of possible distinct values for this item.
// If not representing multiple values, the count is 1.
//
// For instance, a segment with the value range of 3-7 has count 5.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (seg *IPv6AddressSegment) GetCount() *big.Int {
	if seg == nil {
		return bigZero()
	}
	return seg.getCount()
}

// GetPrefixCountLen returns the count of the number of distinct prefix values for the given prefix length in the range of values of this segment.
func (seg *IPv6AddressSegment) GetPrefixCountLen(segmentPrefixLength BitCount) *big.Int {
	return seg.init().ipAddressSegmentInternal.GetPrefixCountLen(segmentPrefixLength)
}

// GetPrefixValueCountLen returns the same value as GetPrefixCountLen as an integer.
func (seg *IPv6AddressSegment) GetPrefixValueCountLen(segmentPrefixLength BitCount) SegIntCount {
	return seg.init().ipAddressSegmentInternal.GetPrefixValueCountLen(segmentPrefixLength)
}

// IsOneBit returns true if the bit in the lower value of this segment at the given index is 1, where index 0 is the most significant bit.
func (seg *IPv6AddressSegment) IsOneBit(segmentBitIndex BitCount) bool {
	return seg.init().ipAddressSegmentInternal.IsOneBit(segmentBitIndex)
}

// Bytes returns the lowest value in the address segment range as a byte slice.
func (seg *IPv6AddressSegment) Bytes() []byte {
	return seg.init().ipAddressSegmentInternal.Bytes()
}

// UpperBytes returns the highest value in the address segment range as a byte slice.
func (seg *IPv6AddressSegment) UpperBytes() []byte {
	return seg.init().ipAddressSegmentInternal.UpperBytes()
}

// CopyBytes copies the lowest value in the address segment range into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (seg *IPv6AddressSegment) CopyBytes(bytes []byte) []byte {
	return seg.init().ipAddressSegmentInternal.CopyBytes(bytes)
}

// CopyUpperBytes copies the highest value in the address segment range into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (seg *IPv6AddressSegment) CopyUpperBytes(bytes []byte) []byte {
	return seg.init().ipAddressSegmentInternal.CopyUpperBytes(bytes)
}

// GetPrefixValueCount returns the count of prefixes in this segment for its prefix length, or the total count if it has no prefix length.
func (seg *IPv6AddressSegment) GetPrefixValueCount() SegIntCount {
	return seg.init().ipAddressSegmentInternal.GetPrefixValueCount()
}

// MatchesWithPrefixMask applies the network mask of the given bit-length to this segment and then compares the result with the given value masked by the same mask,
// returning true if the resulting range matches the given single value.
func (seg *IPv6AddressSegment) MatchesWithPrefixMask(value IPv6SegInt, networkBits BitCount) bool {
	return seg.init().ipAddressSegmentInternal.MatchesWithPrefixMask(SegInt(value), networkBits)
}

// GetBlockMaskPrefixLen returns the prefix length if this address segment is equivalent to the mask for a CIDR prefix block.
// Otherwise, it returns nil.
// A CIDR network mask is a segment with all ones in the network bits and then all zeros in the host bits.
// A CIDR host mask is a segment with all zeros in the network bits and then all ones in the host bits.
// The prefix length is the bit-length of the network bits.
//
// Also, keep in mind that the prefix length returned by this method is not equivalent to the prefix length of this segment.
// The prefix length returned here indicates the whether the value of this segment can be used as a mask for the network and host
// bits of any other segment.  Therefore, the two values can be different values, or one can be nil while the other is not.
//
// This method applies only to the lower value of the range if this segment represents multiple values.
func (seg *IPv6AddressSegment) GetBlockMaskPrefixLen(network bool) PrefixLen {
	return seg.init().ipAddressSegmentInternal.GetBlockMaskPrefixLen(network)
}

// GetTrailingBitCount returns the number of consecutive trailing one or zero bits.
// If ones is true, returns the number of consecutive trailing zero bits.
// Otherwise, returns the number of consecutive trailing one bits.
//
// This method applies only to the lower value of the range if this segment represents multiple values.
func (seg *IPv6AddressSegment) GetTrailingBitCount(ones bool) BitCount {
	return seg.init().ipAddressSegmentInternal.GetTrailingBitCount(ones)
}

// GetLeadingBitCount returns the number of consecutive leading one or zero bits.
// If ones is true, returns the number of consecutive leading one bits.
// Otherwise, returns the number of consecutive leading zero bits.
//
// This method applies only to the lower value of the range if this segment represents multiple values.
func (seg *IPv6AddressSegment) GetLeadingBitCount(ones bool) BitCount {
	return seg.init().ipAddressSegmentInternal.GetLeadingBitCount(ones)
}

// ToPrefixedNetworkSegment returns a segment with the network bits matching this segment but the host bits converted to zero.
// The new segment will be assigned the given prefix length.
func (seg *IPv6AddressSegment) ToPrefixedNetworkSegment(segmentPrefixLength PrefixLen) *IPv6AddressSegment {
	return seg.init().toPrefixedNetworkDivision(segmentPrefixLength).ToIPv6()
}

// ToNetworkSegment returns a segment with the network bits matching this segment but the host bits converted to zero.
// The new segment will have no assigned prefix length.
func (seg *IPv6AddressSegment) ToNetworkSegment(segmentPrefixLength PrefixLen) *IPv6AddressSegment {
	return seg.init().toNetworkDivision(segmentPrefixLength, false).ToIPv6()
}

// ToPrefixedHostSegment returns a segment with the host bits matching this segment but the network bits converted to zero.
// The new segment will be assigned the given prefix length.
func (seg *IPv6AddressSegment) ToPrefixedHostSegment(segmentPrefixLength PrefixLen) *IPv6AddressSegment {
	return seg.init().toPrefixedHostDivision(segmentPrefixLength).ToIPv6()
}

// ToHostSegment returns a segment with the host bits matching this segment but the network bits converted to zero.
// The new segment will have no assigned prefix length.
func (seg *IPv6AddressSegment) ToHostSegment(segmentPrefixLength PrefixLen) *IPv6AddressSegment {
	return seg.init().toHostDivision(segmentPrefixLength, false).ToIPv6()
}

// Iterator provides an iterator to iterate through the individual address segments of this address segment.
//
// When iterating, the prefix length is preserved.  Remove it using WithoutPrefixLen prior to iterating if you wish to drop it from all individual address segments.
//
// Call IsMultiple to determine if this instance represents multiple address segments, or GetValueCount for the count.
func (seg *IPv6AddressSegment) Iterator() Iterator[*IPv6AddressSegment] {
	if seg == nil {
		return ipv6SegmentIterator{nilSegIterator()}
	}
	return ipv6SegmentIterator{seg.init().iterator()}
}

// PrefixBlockIterator provides an iterator to iterate through the individual prefix blocks, one for each prefix of this address segment.
// Each iterated address segment will be a prefix block with the same prefix length as this address segment.
//
// If this address segment has no prefix length, then this is equivalent to Iterator.
func (seg *IPv6AddressSegment) PrefixBlockIterator() Iterator[*IPv6AddressSegment] {
	return ipv6SegmentIterator{seg.init().prefixBlockIterator()}
}

// PrefixedBlockIterator provides an iterator to iterate through the individual prefix blocks of the given prefix length in this segment,
// one for each prefix of this address or subnet.
//
// It is similar to PrefixBlockIterator except that this method allows you to specify the prefix length.
func (seg *IPv6AddressSegment) PrefixedBlockIterator(segmentPrefixLen BitCount) Iterator[*IPv6AddressSegment] {
	return ipv6SegmentIterator{seg.init().prefixedBlockIterator(segmentPrefixLen)}
}

// PrefixIterator provides an iterator to iterate through the individual prefixes of this segment,
// each iterated element spanning the range of values for its prefix.
//
// It is similar to the prefix block iterator, except for possibly the first and last iterated elements, which might not be prefix blocks,
// instead constraining themselves to values from this segment.
//
// If this address segment has no prefix length, then this is equivalent to Iterator.
func (seg *IPv6AddressSegment) PrefixIterator() Iterator[*IPv6AddressSegment] {
	return ipv6SegmentIterator{seg.init().prefixIterator()}
}

// IsPrefixed returns whether this segment has an associated prefix length.
func (seg *IPv6AddressSegment) IsPrefixed() bool {
	return seg != nil && seg.isPrefixed()
}

// WithoutPrefixLen returns a segment with the same value range but without a prefix length.
func (seg *IPv6AddressSegment) WithoutPrefixLen() *IPv6AddressSegment {
	if !seg.IsPrefixed() {
		return seg
	}
	return seg.withoutPrefixLen().ToIPv6()
}

// Converts this IPv6 address segment into smaller segments,
// copying them into the given array starting at the given index.
//
// If a segment does not fit into the array because the segment index in the array is out of bounds of the array,
// then it is not copied.
//
// It is used to create both IPv4 and MAC segments.
func (seg *IPv6AddressSegment) visitSplitSegments(creator func(index int, value, upperValue SegInt, prefLen PrefixLen)) addrerr.IncompatibleAddressError {
	if seg.isMultiple() {
		return seg.visitSplitSegmentsMultiple(creator)
	} else {
		index := 0
		bitSizeSplit := IPv6BitsPerSegment >> 1
		myPrefix := seg.GetSegmentPrefixLen()
		val := seg.highByte()
		highPrefixBits := getSegmentPrefixLength(bitSizeSplit, myPrefix, 0)
		creator(index, val, val, highPrefixBits)
		index++
		val = seg.lowByte()
		lowPrefixBits := getSegmentPrefixLength(bitSizeSplit, myPrefix, 1)
		creator(index, val, val, lowPrefixBits)
		return nil
	}
}

func (seg *IPv6AddressSegment) splitSegValues() (highLower, highUpper, lowLower, lowUpper SegInt, err addrerr.IncompatibleAddressError) {
	val := seg.GetSegmentValue()
	upperVal := seg.GetUpperSegmentValue()
	highLower = highByteIpv6(val)
	highUpper = highByteIpv6(upperVal)
	lowLower = lowByteIpv6(val)
	lowUpper = lowByteIpv6(upperVal)
	if (highLower != highUpper) && (lowLower != 0 || lowUpper != 0xff) {
		err = &incompatibleAddressError{addressError{key: "ipaddress.error.splitSeg"}}
	}
	return
}

// Used to create both IPv4 and MAC segments
func (seg *IPv6AddressSegment) visitSplitSegmentsMultiple(creator func(index int, value, upperValue SegInt, prefLen PrefixLen)) addrerr.IncompatibleAddressError {
	myPrefix := seg.GetSegmentPrefixLen()
	bitSizeSplit := BitCount(IPv6BitsPerSegment >> 1)
	highLower, highUpper, lowLower, lowUpper, err := seg.splitSegValues()
	if err != nil {
		return err
	}
	highPrefixBits := getSegmentPrefixLength(bitSizeSplit, myPrefix, 0)
	lowPrefixBits := getSegmentPrefixLength(bitSizeSplit, myPrefix, 1)
	creator(0, highLower, highUpper, highPrefixBits)
	creator(1, lowLower, lowUpper, lowPrefixBits)
	return nil
}

func (seg *IPv6AddressSegment) highByte() SegInt {
	return highByteIpv6(seg.GetSegmentValue())
}

func (seg *IPv6AddressSegment) lowByte() SegInt {
	return lowByteIpv6(seg.GetSegmentValue())
}

func highByteIpv6(value SegInt) SegInt {
	return value >> 8
}

func lowByteIpv6(value SegInt) SegInt {
	return value & 0xff
}

// Converts this IPv6 address segment into smaller segments,
// copying them into the given array starting at the given index.
//
// If a segment does not fit into the array because the segment index in the array is out of bounds of the array,
// then it is not copied.
func (seg *IPv6AddressSegment) getSplitSegments(segs []*IPv4AddressSegment, startIndex int) addrerr.IncompatibleAddressError {
	return seg.visitSplitSegments(func(index int, value, upperValue SegInt, prefLen PrefixLen) {
		if ind := startIndex + index; ind < len(segs) {
			segs[ind] = NewIPv4RangePrefixedSegment(IPv4SegInt(value), IPv4SegInt(upperValue), prefLen)
		}
	})
}

func (seg *IPv6AddressSegment) splitIntoIPv4Segments(segs []*AddressDivision, startIndex int) addrerr.IncompatibleAddressError {
	return seg.visitSplitSegments(func(index int, value, upperValue SegInt, prefLen PrefixLen) {
		if ind := startIndex + index; ind < len(segs) {
			segs[ind] = NewIPv4RangePrefixedSegment(IPv4SegInt(value), IPv4SegInt(upperValue), prefLen).ToDiv()
		}
	})
}

func (seg *IPv6AddressSegment) splitIntoMACSegments(segs []*AddressDivision, startIndex int) addrerr.IncompatibleAddressError {
	return seg.visitSplitSegments(func(index int, value, upperValue SegInt, prefLen PrefixLen) {
		if ind := startIndex + index; ind < len(segs) {
			segs[ind] = NewMACRangeSegment(MACSegInt(value), MACSegInt(upperValue)).ToDiv()
		}
	})
}

// ReverseBits returns a segment with the bits reversed.
//
// If this segment represents a range of values that cannot be reversed, then this returns an error.
//
// To be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
// Otherwise the result is not contiguous and thus cannot be represented by a sequential range of values.
//
// If perByte is true, the bits are reversed within each byte, otherwise all the bits are reversed.
//
// If perByte is true, the bits are reversed within each byte, otherwise all the bits are reversed.
func (seg *IPv6AddressSegment) ReverseBits(perByte bool) (res *IPv6AddressSegment, err addrerr.IncompatibleAddressError) {
	if seg.divisionValues == nil {
		res = seg
		return
	}
	if seg.isMultiple() {
		var addrSeg *AddressSegment
		addrSeg, err = seg.reverseMultiValSeg(perByte)
		res = addrSeg.ToIPv6()
		return
	}
	oldVal := IPv6SegInt(seg.GetSegmentValue())
	val := IPv6SegInt(reverseUint16(uint16(oldVal)))
	if perByte {
		val = ((val & 0xff) << 8) | (val >> 8)
	}
	if oldVal == val && !seg.isPrefixed() {
		res = seg
	} else {
		res = NewIPv6Segment(val)
	}
	return
}

// ReverseBytes returns a segment with the bytes reversed.
//
// If this segment represents a range of values that cannot be reversed, then this returns an error.
//
// To be reversible, a range must include all values except possibly the largest and/or smallest, which reverse to themselves.
// Otherwise the result is not contiguous and thus cannot be represented by a sequential range of values.
func (seg *IPv6AddressSegment) ReverseBytes() (res *IPv6AddressSegment, err addrerr.IncompatibleAddressError) {
	if seg.divisionValues == nil {
		res = seg
		return
	}
	if seg.isMultiple() {
		var addrSeg *AddressSegment
		addrSeg, err = seg.reverseMultiValSeg(false)
		res = addrSeg.ToIPv6()
		return
	}
	oldVal := IPv6SegInt(seg.GetSegmentValue())
	val := IPv6SegInt(reverseUint16(uint16(oldVal)))
	if oldVal == val && !seg.isPrefixed() {
		res = seg
	} else {
		res = NewIPv6Segment(val)
	}
	return
}

// ToDiv converts to an AddressDivision, a polymorphic type usable with all address segments and divisions.
// Afterwards, you can convert back with ToIPv6.
//
// ToDiv can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (seg *IPv6AddressSegment) ToDiv() *AddressDivision {
	return seg.ToIP().ToDiv()
}

// ToSegmentBase converts to an AddressSegment, a polymorphic type usable with all address segments.
// Afterwards, you can convert back with ToIPv6.
//
// ToSegmentBase can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (seg *IPv6AddressSegment) ToSegmentBase() *AddressSegment {
	return seg.ToIP().ToSegmentBase()
}

// ToIP converts to an IPAddressSegment, a polymorphic type usable with all IP address segments.
// Afterwards, you can convert back with ToIPv6.
//
// ToIP can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (seg *IPv6AddressSegment) ToIP() *IPAddressSegment {
	if seg == nil {
		return nil
	}
	return (*IPAddressSegment)(seg.init())
}

// GetString produces a normalized string to represent the segment.
// If the segment is a CIDR network prefix block for its prefix length, then the string contains only the lower value of the block range.
// Otherwise, the explicit range will be printed.
//
// The string returned is useful in the context of creating strings for address sections or full addresses,
// in which case the radix and bit-length can be deduced from the context.
// The String method produces strings more appropriate when no context is provided.
func (seg *IPv6AddressSegment) GetString() string {
	if seg == nil {
		return nilString()
	}
	return seg.init().getString()
}

// GetWildcardString produces a normalized string to represent the segment, favouring wildcards and range characters while ignoring any network prefix length.
// The explicit range of a range-valued segment will be printed.
//
// The string returned is useful in the context of creating strings for address sections or full addresses,
// in which case the radix and the bit-length can be deduced from the context.
// The String method produces strings more appropriate when no context is provided.
func (seg *IPv6AddressSegment) GetWildcardString() string {
	if seg == nil {
		return nilString()
	}
	return seg.init().getWildcardString()
}

// String produces a string that is useful when a segment is provided with no context.  It uses the hexadecimal radix with the string prefix for hex ("0x").
// GetWildcardString is more appropriate in context with other segments or divisions.  It does not use a string prefix and uses '*' for full-range segments.
// GetString is more appropriate in context with prefix lengths, it uses zeros instead of wildcards with full prefix block ranges alongside prefix lengths.
func (seg *IPv6AddressSegment) String() string {
	if seg == nil {
		return nilString()
	}
	return seg.init().toString()
}

// NewIPv6Segment constructs a segment of an IPv6 address with the given value.
func NewIPv6Segment(val IPv6SegInt) *IPv6AddressSegment {
	return newIPv6Segment(newIPv6SegmentVal(val))
}

// NewIPv6RangeSegment constructs a segment of an IPv6 subnet with the given range of sequential values.
func NewIPv6RangeSegment(val, upperVal IPv6SegInt) *IPv6AddressSegment {
	return newIPv6Segment(newIPv6SegmentPrefixedValues(val, upperVal, nil))
}

// NewIPv6PrefixedSegment constructs a segment of an IPv6 address with the given value and assigned prefix length.
func NewIPv6PrefixedSegment(val IPv6SegInt, prefixLen PrefixLen) *IPv6AddressSegment {
	return newIPv6Segment(newIPv6SegmentPrefixedVal(val, prefixLen))
}

// NewIPv6RangePrefixedSegment constructs a segment of an IPv6 subnet with the given range of sequential values and assigned prefix length.
func NewIPv6RangePrefixedSegment(val, upperVal IPv6SegInt, prefixLen PrefixLen) *IPv6AddressSegment {
	return newIPv6Segment(newIPv6SegmentPrefixedValues(val, upperVal, prefixLen))
}

func newIPv6Segment(vals *ipv6SegmentValues) *IPv6AddressSegment {
	return &IPv6AddressSegment{
		ipAddressSegmentInternal{
			addressSegmentInternal{
				addressDivisionInternal{
					addressDivisionBase{vals},
				},
			},
		},
	}
}

type ipv6DivsBlock struct {
	block []ipv6SegmentValues
}

type ipv6DivsPartition struct {
	block []*ipv6DivsBlock
}

var (
	allRangeValsIPv6 = &ipv6SegmentValues{
		upperValue: IPv6MaxValuePerSegment,
		cache: divCache{
			isSinglePrefBlock: &falseVal,
		},
	}
	allPrefixedCacheIPv6 = makePrefixCacheIPv6()

	// single-valued no-prefix cache.
	// there are 0x10000 (ie 0xffff + 1 or 64k) possible segment values in IPv6.  We break the cache into 0x100 blocks of size 0x100
	segmentCacheIPv6 = make([]*ipv6DivsBlock, (IPv6MaxValuePerSegment>>8)+1)

	// single-valued cache for each prefix.
	segmentPrefixCacheIPv6 = make([]*ipv6DivsPartition, IPv6BitsPerSegment+1) // for each prefix, all segment values, 0x100 blocks of size 0x100

	// prefix-block cache: all the prefix blocks for each prefix.
	// for each prefix, all prefix blocks.
	// For a given prefix, you shift left by 8 bits for the blocks of size 0x100, the remaining bits to the left are the number of blocks.
	//
	// For prefix of size 8, 1 block of size 0x100
	// For prefix of size < 8, 1 block of size (1 << prefix)
	// For prefix of size > 8, (1 << (prefix - 8)) blocks of size 0x100.
	//
	// So, you start with the prefix to get the right ipv6DivsPartition.
	// Then, you use the formula above to look up the block index.
	// For the first two above, the whole prefix finds the index into the single block.
	// For the third, the 8 rightmost bits in the prefix give the index into the block of size ff,
	// while the leftmost bits in the prefix select that block.
	prefixBlocksCacheIPv6 = make([]*ipv6DivsPartition, IPv6BitsPerSegment+1)
)

func makePrefixCacheIPv6() (allPrefixedCacheIPv6 []ipv6SegmentValues) {
	if useIPv6SegmentCache {
		allPrefixedCacheIPv6 = make([]ipv6SegmentValues, IPv6BitsPerSegment+1)
		for i := range allPrefixedCacheIPv6 {
			vals := &allPrefixedCacheIPv6[i]
			vals.upperValue = IPv6MaxValuePerSegment
			vals.prefLen = cacheBitCount(i)
			vals.cache.isSinglePrefBlock = &falseVal
		}
		allPrefixedCacheIPv6[0].cache.isSinglePrefBlock = &trueVal
	}
	return
}

func newIPv6SegmentVal(value IPv6SegInt) *ipv6SegmentValues {
	if useIPv6SegmentCache {
		cache := segmentCacheIPv6
		blockIndex := value >> 8 // divide by 0x100
		firstBlockVal := blockIndex << 8
		resultIndex := value - firstBlockVal // mod 0x100
		block := (*ipv6DivsBlock)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache[blockIndex]))))
		//block := cache[blockIndex]
		if block == nil {
			block = &ipv6DivsBlock{make([]ipv6SegmentValues, 0x100)}
			vals := block.block
			for i := range vals {
				item := &vals[i]
				itemVal := firstBlockVal | IPv6SegInt(i)
				item.value = itemVal
				item.upperValue = itemVal
				item.cache.isSinglePrefBlock = &falseVal
			}
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache[blockIndex]))
			atomicStorePointer(dataLoc, unsafe.Pointer(block))
		}
		result := &block.block[resultIndex]
		return result
	}
	return &ipv6SegmentValues{
		value:      value,
		upperValue: value,
		cache: divCache{
			isSinglePrefBlock: &falseVal,
		},
	}
}

func newIPv6SegmentPrefixedVal(value IPv6SegInt, prefLen PrefixLen) (result *ipv6SegmentValues) {
	if prefLen == nil {
		return newIPv6SegmentVal(value)
	}
	prefixIndex := prefLen.bitCount()
	if prefixIndex < 0 {
		prefixIndex = 0
	} else if prefixIndex > IPv6BitsPerSegment {
		prefixIndex = IPv6BitsPerSegment
	}
	prefLen = cacheBitCount(prefixIndex) // this ensures we use the prefix length cache for all segments
	if useIPv6SegmentCache {
		cache := segmentPrefixCacheIPv6
		prefixCache := (*ipv6DivsPartition)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache[prefixIndex]))))
		if prefixCache == nil {
			prefixCache = &ipv6DivsPartition{make([]*ipv6DivsBlock, (IPv6MaxValuePerSegment>>8)+1)}
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache[prefixIndex]))
			atomicStorePointer(dataLoc, unsafe.Pointer(prefixCache))
		}
		blockIndex := value >> 8 // divide by 0x100
		firstBlockVal := blockIndex << 8
		resultIndex := value - (firstBlockVal) // mod 0x100
		blockCache := (*ipv6DivsBlock)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&prefixCache.block[blockIndex]))))
		if blockCache == nil {
			blockCache = &ipv6DivsBlock{make([]ipv6SegmentValues, (IPv6MaxValuePerSegment>>8)+1)}
			vals := blockCache.block
			var isSinglePrefBlock *bool
			if prefixIndex == IPv6BitsPerSegment {
				isSinglePrefBlock = &trueVal
			} else {
				isSinglePrefBlock = &falseVal
			}
			for i := range vals {
				item := &vals[i]
				itemVal := firstBlockVal | IPv6SegInt(i)
				item.value = itemVal
				item.upperValue = itemVal
				item.prefLen = prefLen
				item.cache.isSinglePrefBlock = isSinglePrefBlock
			}
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&prefixCache.block[blockIndex]))
			atomicStorePointer(dataLoc, unsafe.Pointer(blockCache))
		}
		result := &blockCache.block[resultIndex]
		return result
	}
	var isSinglePrefBlock *bool
	if prefixIndex == IPv6BitsPerSegment {
		isSinglePrefBlock = &trueVal
	} else {
		isSinglePrefBlock = &falseVal
	}
	return &ipv6SegmentValues{
		value:      value,
		upperValue: value,
		prefLen:    prefLen,
		cache: divCache{
			isSinglePrefBlock: isSinglePrefBlock,
		},
	}
}

func newIPv6SegmentPrefixedValues(value, upperValue IPv6SegInt, prefLen PrefixLen) *ipv6SegmentValues {
	var isSinglePrefBlock *bool
	if prefLen == nil {
		if value == upperValue {
			return newIPv6SegmentVal(value)
		} else if value > upperValue {
			value, upperValue = upperValue, value
		}
		if useIPv6SegmentCache && value == 0 && upperValue == IPv6MaxValuePerSegment {
			return allRangeValsIPv6
		}
		isSinglePrefBlock = &falseVal
	} else {
		if value == upperValue {
			return newIPv6SegmentPrefixedVal(value, prefLen)
		} else if value > upperValue {
			value, upperValue = upperValue, value
		}
		prefixIndex := prefLen.bitCount()
		if prefixIndex < 0 {
			prefixIndex = 0
		} else if prefixIndex > IPv6BitsPerSegment {
			prefixIndex = IPv6BitsPerSegment
		}
		prefLen = cacheBitCount(prefixIndex) // this ensures we use the prefix length cache for all segments
		if useIPv6SegmentCache {
			shiftBits := uint(IPv6BitsPerSegment - prefixIndex)
			nmask := ^IPv6SegInt(0) << shiftBits
			prefixBlockLower := value & nmask
			hmask := ^nmask
			prefixBlockUpper := value | hmask
			if value == prefixBlockLower && upperValue == prefixBlockUpper {
				// cache is the prefix block for any prefix length
				cache := prefixBlocksCacheIPv6
				prefixCache := (*ipv6DivsPartition)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache[prefixIndex]))))
				if prefixCache == nil {
					if prefixIndex <= 8 { // 1 block of size (1 << prefix)
						prefixCache = &ipv6DivsPartition{make([]*ipv6DivsBlock, 1)}
					} else { // (1 << (prefix - 8)) blocks of size 0x100.
						prefixCache = &ipv6DivsPartition{make([]*ipv6DivsBlock, 1<<uint(prefixIndex-8))}
					}
					dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cache[prefixIndex]))
					atomicStorePointer(dataLoc, unsafe.Pointer(prefixCache))
				}
				valueIndex := value >> shiftBits
				blockIndex := valueIndex >> 8 // divide by 0x100
				firstBlockVal := blockIndex << 8
				resultIndex := valueIndex - (firstBlockVal) // mod 0x100
				blockCache := (*ipv6DivsBlock)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&prefixCache.block[blockIndex]))))
				if blockCache == nil {
					if prefixIndex <= 8 { // 1 block of size (1 << prefix)
						blockCache = &ipv6DivsBlock{make([]ipv6SegmentValues, 1<<uint(prefixIndex))}
					} else { // (1 << (prefix - 8)) blocks of size 0x100.
						blockCache = &ipv6DivsBlock{make([]ipv6SegmentValues, 1<<8)}
					}
					vals := blockCache.block
					for i := range vals {
						item := &vals[i]
						itemVal := (firstBlockVal | IPv6SegInt(i)) << shiftBits
						item.value = itemVal
						item.upperValue = itemVal | hmask
						item.prefLen = prefLen
						item.cache.isSinglePrefBlock = &trueVal
					}
					dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&prefixCache.block[blockIndex]))
					atomicStorePointer(dataLoc, unsafe.Pointer(blockCache))
				}
				result := &blockCache.block[resultIndex]
				return result
			}
			if value == 0 {
				// cache is 0-0xffff for any prefix length
				if upperValue == IPv6MaxValuePerSegment {
					result := &allPrefixedCacheIPv6[prefixIndex]
					return result
				}
			}
			isSinglePrefBlock = &falseVal
		}
	}
	return &ipv6SegmentValues{
		value:      value,
		upperValue: upperValue,
		prefLen:    prefLen,
		cache: divCache{
			isSinglePrefBlock: isSinglePrefBlock,
		},
	}
}
