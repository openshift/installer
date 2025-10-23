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
	"net"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstr"
)

const (
	MACBitsPerSegment           = 8
	MACBytesPerSegment          = 1
	MACDefaultTextualRadix      = 16
	MACMaxValuePerSegment       = 0xff
	MACMaxValuePerDottedSegment = 0xffff

	MediaAccessControlSegmentCount         = 6
	MediaAccessControlDottedSegmentCount   = 3
	MediaAccessControlDotted64SegmentCount = 4
	ExtendedUniqueIdentifier48SegmentCount = MediaAccessControlSegmentCount
	ExtendedUniqueIdentifier64SegmentCount = 8

	MACOrganizationalUniqueIdentifierSegmentCount = 3

	MACSegmentMaxChars = 2

	MACDashSegmentSeparator   = '-'
	MACColonSegmentSeparator  = ':'
	MacSpaceSegmentSeparator  = ' '
	MacDottedSegmentSeparator = '.'

	MacDashedSegmentRangeSeparator    = '|'
	MacDashedSegmentRangeSeparatorStr = "|"

	macBitsToSegmentBitshift = 3
)

func newMACAddress(section *MACAddressSection) *MACAddress {
	return createAddress(section.ToSectionBase(), NoZone).ToMAC()
}

// NewMACAddress constructs a MAC address or address collection from the given segments.
func NewMACAddress(section *MACAddressSection) (*MACAddress, addrerr.AddressValueError) {
	segCount := section.GetSegmentCount()
	if segCount != MediaAccessControlSegmentCount && segCount != ExtendedUniqueIdentifier64SegmentCount {
		return nil, &addressValueError{
			addressError: addressError{key: "ipaddress.error.invalid.size"},
			val:          segCount,
		}
	}
	return createAddress(section.ToSectionBase(), NoZone).ToMAC(), nil
}

// NewMACAddressFromBytes constructs a MAC address from the given byte slice.
// An error is returned when the byte slice has too many bytes to match the maximum MAC segment count of 8.
// There should be 8 bytes or less, although extra leading zeros are tolerated.
func NewMACAddressFromBytes(bytes net.HardwareAddr) (*MACAddress, addrerr.AddressValueError) {
	section, err := createMACSectionFromBytes(bytes)
	if err != nil {
		return nil, err
	}
	segCount := section.GetSegmentCount()
	if segCount != MediaAccessControlSegmentCount && segCount != ExtendedUniqueIdentifier64SegmentCount {
		return nil, &addressValueError{
			addressError: addressError{key: "ipaddress.error.invalid.size"},
			val:          segCount,
		}
	}
	return createAddress(section.ToSectionBase(), NoZone).ToMAC(), nil
}

// NewMACAddressFromUint64Ext constructs a 6 or 8-byte MAC address from the given value.
// If isExtended is true, it is an 8-byte address, 6 otherwise.
// If 6 bytes, then the bytes are taken from the lower 48 bits of the uint64.
func NewMACAddressFromUint64Ext(val uint64, isExtended bool) *MACAddress {
	section := NewMACSectionFromUint64(val, getMacSegCount(isExtended))
	return createAddress(section.ToSectionBase(), NoZone).ToMAC()
}

// NewMACAddressFromSegs constructs a MAC address or address collection from the given segments.
// If the given slice does not have either 6 or 8 segments, an error is returned.
func NewMACAddressFromSegs(segments []*MACAddressSegment) (*MACAddress, addrerr.AddressValueError) {
	segsLen := len(segments)
	if segsLen != MediaAccessControlSegmentCount && segsLen != ExtendedUniqueIdentifier64SegmentCount {
		return nil, &addressValueError{val: segsLen, addressError: addressError{key: "ipaddress.error.mac.invalid.segment.count"}}
	}
	section := NewMACSection(segments)
	return createAddress(section.ToSectionBase(), NoZone).ToMAC(), nil
}

// NewMACAddressFromVals constructs a 6-byte MAC address from the given values.
func NewMACAddressFromVals(vals MACSegmentValueProvider) (addr *MACAddress) {
	return NewMACAddressFromValsExt(vals, false)
}

// NewMACAddressFromValsExt constructs a 6 or 8-byte MAC address from the given values.
// If isExtended is true, it will be 8 bytes.
func NewMACAddressFromValsExt(vals MACSegmentValueProvider, isExtended bool) (addr *MACAddress) {
	section := NewMACSectionFromVals(vals, getMacSegCount(isExtended))
	addr = newMACAddress(section)
	return
}

// NewMACAddressFromRange constructs a 6-byte MAC address collection from the given values.
func NewMACAddressFromRange(vals, upperVals MACSegmentValueProvider) (addr *MACAddress) {
	return NewMACAddressFromRangeExt(vals, upperVals, false)
}

// NewMACAddressFromRangeExt constructs a 6 or 8-byte MAC address collection from the given values.
// If isExtended is true, it will be 8 bytes.
func NewMACAddressFromRangeExt(vals, upperVals MACSegmentValueProvider, isExtended bool) (addr *MACAddress) {
	section := NewMACSectionFromRange(vals, upperVals, getMacSegCount(isExtended))
	addr = newMACAddress(section)
	return
}

func createMACSectionFromBytes(bytes []byte) (*MACAddressSection, addrerr.AddressValueError) {
	var segCount int
	length := len(bytes)
	//We round down the bytes to 6 bytes if we can.  Otherwise, we round up.
	if length < ExtendedUniqueIdentifier64SegmentCount {
		segCount = MediaAccessControlSegmentCount
		if length > MediaAccessControlSegmentCount {
			for i := 0; ; i++ {
				if bytes[i] != 0 {
					segCount = ExtendedUniqueIdentifier64SegmentCount
					break
				}
				length--
				if length <= MediaAccessControlSegmentCount {
					break
				}
			}
		}
	} else {
		segCount = ExtendedUniqueIdentifier64SegmentCount
	}
	return NewMACSectionFromBytes(bytes, segCount)
}

func getMacSegCount(isExtended bool) (segmentCount int) {
	if isExtended {
		segmentCount = ExtendedUniqueIdentifier64SegmentCount
	} else {
		segmentCount = MediaAccessControlSegmentCount
	}
	return
}

var zeroMAC = createMACZero(false)
var macAll = zeroMAC.SetPrefixLen(0).ToPrefixBlock()
var macAllExtended = createMACZero(true).SetPrefixLen(0).ToPrefixBlock()

func createMACZero(extended bool) *MACAddress {
	segs := []*MACAddressSegment{zeroMACSeg, zeroMACSeg, zeroMACSeg, zeroMACSeg, zeroMACSeg, zeroMACSeg}
	if extended {
		segs = append(segs, zeroMACSeg, zeroMACSeg)
	}
	section := NewMACSection(segs)
	return newMACAddress(section)
}

// MACAddress represents a MAC address, or a collection of multiple individual MAC addresses.
// Each segment can represent a single byte value or a range of byte values.
//
// You can construct a MAC address from a byte slice, from a uint64, from a SegmentValueProvider,
// from a MACAddressSection of 6 or 8 segments, or from an array of 6 or 8 MACAddressSegment instances.
//
// To construct one from a string, use NewMACAddressString, then use the ToAddress or GetAddress method of [MACAddressString].
type MACAddress struct {
	addressInternal
}

// GetCount returns the count of addresses that this address or address collection represents.
//
// If just a single address, not a collection of multiple addresses, returns 1.
func (addr *MACAddress) init() *MACAddress {
	if addr.section == nil {
		return zeroMAC
	}
	return addr
}

// GetCount returns the count of addresses that this address or address collection represents.
//
// If just a single address, not a collection of multiple addresses, returns 1.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (addr *MACAddress) GetCount() *big.Int {
	if addr == nil {
		return bigZero()
	}
	return addr.getCount()
}

// IsMultiple returns true if this represents more than a single individual address, whether it is a collection of multiple addresses.
func (addr *MACAddress) IsMultiple() bool {
	return addr != nil && addr.isMultiple()
}

// IsPrefixed returns whether this address has an associated prefix length.
func (addr *MACAddress) IsPrefixed() bool {
	return addr != nil && addr.isPrefixed()
}

// IsFullRange returns whether this address covers the entire MAC address space for its MAC bit length.
//
// This is true if and only if both IncludesZero and IncludesMax return true.
func (addr *MACAddress) IsFullRange() bool {
	return addr.GetSection().IsFullRange()
}

// GetBitCount returns the number of bits comprising this address,
// or each address in the range.
func (addr *MACAddress) GetBitCount() BitCount {
	return addr.init().addressInternal.GetBitCount()
}

// GetByteCount returns the number of bytes required for this address,
// or each address in the range.
func (addr *MACAddress) GetByteCount() int {
	return addr.init().addressInternal.GetByteCount()
}

// GetBitsPerSegment returns the number of bits comprising each segment in this address.  Segments in the same address are equal length.
func (addr *MACAddress) GetBitsPerSegment() BitCount {
	return MACBitsPerSegment
}

// GetBytesPerSegment returns the number of bytes comprising each segment in this address.  Segments in the same address are equal length.
func (addr *MACAddress) GetBytesPerSegment() int {
	return MACBytesPerSegment
}

func (addr *MACAddress) checkIdentity(section *MACAddressSection) *MACAddress {
	if section == nil {
		return nil
	}
	sec := section.ToSectionBase()
	if sec == addr.section {
		return addr
	}
	return newMACAddress(section)
}

// GetValue returns the lowest address in this subnet or address as an integer value.
func (addr *MACAddress) GetValue() *big.Int {
	return addr.init().section.GetValue()
}

// GetUpperValue returns the highest address in this subnet or address as an integer value.
func (addr *MACAddress) GetUpperValue() *big.Int {
	return addr.init().section.GetUpperValue()
}

// GetLower returns the address in the collection with the lowest numeric value,
// which will be the receiver if it represents a single address.
// For example, for "1:1:1:2-3:4:5-6", the series "1:1:1:2:4:5" is returned.
func (addr *MACAddress) GetLower() *MACAddress {
	return addr.init().getLower().ToMAC()
}

// GetUpper returns the address in the collection with the highest numeric value,
// which will be the receiver if it represents a single address.
// For example, for "1:1:1:2-3:4:5-6", the series "1:1:1:3:4:6" is returned.
func (addr *MACAddress) GetUpper() *MACAddress {
	return addr.init().getUpper().ToMAC()
}

// Uint64Value returns the lowest address in the address collection as a uint64.
func (addr *MACAddress) Uint64Value() uint64 {
	return addr.GetSection().Uint64Value()
}

// UpperUint64Value returns the highest address in the address collection as a uint64.
func (addr *MACAddress) UpperUint64Value() uint64 {
	return addr.GetSection().UpperUint64Value()
}

// GetHardwareAddr returns the lowest address in this address or address collection as a net.HardwareAddr.
func (addr *MACAddress) GetHardwareAddr() net.HardwareAddr {
	return addr.Bytes()
}

// CopyHardwareAddr copies the value of the lowest individual address in the address collection into a net.HardwareAddr.
//
// If the value can fit in the given net.HardwareAddr, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new net.HardwareAddr is created and returned with the value.
func (addr *MACAddress) CopyHardwareAddr(bytes net.HardwareAddr) net.HardwareAddr {
	return addr.CopyBytes(bytes)
}

// GetUpperHardwareAddr returns the highest address in this address or address collection as a net.HardwareAddr.
func (addr *MACAddress) GetUpperHardwareAddr() net.HardwareAddr {
	return addr.UpperBytes()
}

// CopyUpperHardwareAddr copies the value of the highest individual address in the address collection into a net.HardwareAddr.
//
// If the value can fit in the given net.HardwareAddr, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new net.HardwareAddr is created and returned with the value.
func (addr *MACAddress) CopyUpperHardwareAddr(bytes net.HardwareAddr) net.HardwareAddr {
	return addr.CopyUpperBytes(bytes)
}

// Bytes returns the lowest address in this address or address collection as a byte slice.
func (addr *MACAddress) Bytes() []byte {
	return addr.init().section.Bytes()
}

// UpperBytes returns the highest address in this address or address collection as a byte slice.
func (addr *MACAddress) UpperBytes() []byte {
	return addr.init().section.UpperBytes()
}

// CopyBytes copies the value of the lowest individual address in the address collection into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (addr *MACAddress) CopyBytes(bytes []byte) []byte {
	return addr.init().section.CopyBytes(bytes)
}

// CopyUpperBytes copies the value of the highest individual address in the address collection into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (addr *MACAddress) CopyUpperBytes(bytes []byte) []byte {
	return addr.init().section.CopyUpperBytes(bytes)
}

// GetSection returns the backing section for this address or address collection, comprising all segments.
func (addr *MACAddress) GetSection() *MACAddressSection {
	return addr.init().section.ToMAC()
}

// GetTrailingSection gets the subsection from the series starting from the given index.
// The first segment is at index 0.
func (addr *MACAddress) GetTrailingSection(index int) *MACAddressSection {
	return addr.GetSection().GetTrailingSection(index)
}

// GetSubSection gets the subsection from the series starting from the given index and ending just before the give endIndex.
// The first segment is at index 0.
func (addr *MACAddress) GetSubSection(index, endIndex int) *MACAddressSection {
	return addr.GetSection().GetSubSection(index, endIndex)
}

// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
func (addr *MACAddress) CopySubSegments(start, end int, segs []*MACAddressSegment) (count int) {
	return addr.GetSection().CopySubSegments(start, end, segs)
}

// CopySegments copies the existing segments into the given slice,
// as much as can be fit into the slice, returning the number of segments copied.
func (addr *MACAddress) CopySegments(segs []*MACAddressSegment) (count int) {
	return addr.GetSection().CopySegments(segs)
}

// GetSegments returns a slice with the address segments.  The returned slice is not backed by the same array as this address.
func (addr *MACAddress) GetSegments() []*MACAddressSegment {
	return addr.GetSection().GetSegments()
}

// GetSegment returns the segment at the given index.
// The first segment is at index 0.
// GetSegment will panic given a negative index or an index matching or larger than the segment count.
func (addr *MACAddress) GetSegment(index int) *MACAddressSegment {
	return addr.init().getSegment(index).ToMAC()
}

// GetSegmentCount returns the segment/division count
func (addr *MACAddress) GetSegmentCount() int {
	return addr.GetDivisionCount()
}

// ForEachSegment visits each segment in order from most-significant to least, the most significant with index 0, calling the given function for each, terminating early if the function returns true.
// Returns the number of visited segments.
func (addr *MACAddress) ForEachSegment(consumer func(segmentIndex int, segment *MACAddressSegment) (stop bool)) int {
	return addr.GetSection().ForEachSegment(consumer)
}

// GetGenericDivision returns the segment at the given index as a DivisionType.
func (addr *MACAddress) GetGenericDivision(index int) DivisionType {
	return addr.init().getDivision(index)
}

// GetGenericSegment returns the segment at the given index as an AddressSegmentType.
// The first segment is at index 0.
// GetGenericSegment will panic given a negative index or an index matching or larger than the segment count.
func (addr *MACAddress) GetGenericSegment(index int) AddressSegmentType {
	return addr.init().getSegment(index)
}

// TestBit returns true if the bit in the lower value of this address at the given index is 1, where index 0 refers to the least significant bit.
// In other words, it computes (bits & (1 << n)) != 0), using the lower value of this address.
// TestBit will panic if n < 0, or if it matches or exceeds the bit count of this item.
func (addr *MACAddress) TestBit(n BitCount) bool {
	return addr.init().testBit(n)
}

// IsOneBit returns true if the bit in the lower value of this address at the given index is 1, where index 0 refers to the most significant bit.
// IsOneBit will panic if bitIndex is less than zero, or if it is larger than the bit count of this item.
func (addr *MACAddress) IsOneBit(bitIndex BitCount) bool {
	return addr.init().isOneBit(bitIndex)
}

// IsMax returns whether this address matches exactly the maximum possible value, the address whose bits are all ones.
func (addr *MACAddress) IsMax() bool {
	return addr.init().section.IsMax()
}

// IncludesMax returns whether this address includes the max address, the address whose bits are all ones, within its range.
func (addr *MACAddress) IncludesMax() bool {
	return addr.init().section.IncludesMax()
}

// GetDivisionCount returns the segment count, implementing the interface AddressDivisionSeries.
func (addr *MACAddress) GetDivisionCount() int {
	return addr.init().getDivisionCount()
}

// ToPrefixBlock returns the address associated with the prefix of this address or address collection,
// the address whose prefix matches the prefix of this address, and the remaining bits span all values.
// If this address has no prefix length, this address is returned.
//
// The returned address collection will include all addresses with the same prefix as this one, the prefix "block".
func (addr *MACAddress) ToPrefixBlock() *MACAddress {
	return addr.init().toPrefixBlock().ToMAC()
}

// ToPrefixBlockLen returns the address associated with the prefix length provided,
// the address collection whose prefix of that length matches the prefix of this address, and the remaining bits span all values.
//
// The returned address will include all addresses with the same prefix as this one, the prefix "block".
func (addr *MACAddress) ToPrefixBlockLen(prefLen BitCount) *MACAddress {
	return addr.init().toPrefixBlockLen(prefLen).ToMAC()
}

// ToBlock creates a new block of addresses by changing the segment at the given index to have the given lower and upper value,
// and changing the following segments to be full-range.
func (addr *MACAddress) ToBlock(segmentIndex int, lower, upper SegInt) *MACAddress {
	return addr.init().toBlock(segmentIndex, lower, upper).ToMAC()
}

// WithoutPrefixLen provides the same address but with no prefix length.  The values remain unchanged.
func (addr *MACAddress) WithoutPrefixLen() *MACAddress {
	if !addr.IsPrefixed() {
		return addr
	}
	return addr.init().withoutPrefixLen().ToMAC()
}

// SetPrefixLen sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the address.
// The provided prefix length will be adjusted to these boundaries if necessary.
func (addr *MACAddress) SetPrefixLen(prefixLen BitCount) *MACAddress {
	return addr.init().setPrefixLen(prefixLen).ToMAC()
}

// SetPrefixLenZeroed sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the address.
// The provided prefix length will be adjusted to these boundaries if necessary.
//
// If this address has a prefix length, and the prefix length is increased when setting the new prefix length, the bits moved within the prefix become zero.
// If this address has a prefix length, and the prefix length is decreased when setting the new prefix length, the bits moved outside the prefix become zero.
//
// In other words, bits that move from one side of the prefix length to the other (bits moved into the prefix or outside the prefix) are zeroed.
//
// If the result cannot be zeroed because zeroing out bits results in a non-contiguous segment, an error is returned.

func (addr *MACAddress) SetPrefixLenZeroed(prefixLen BitCount) (*MACAddress, addrerr.IncompatibleAddressError) {
	res, err := addr.init().setPrefixLenZeroed(prefixLen)
	return res.ToMAC(), err
}

// AdjustPrefixLen increases or decreases the prefix length by the given increment.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the address.
//
// If this address has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
func (addr *MACAddress) AdjustPrefixLen(prefixLen BitCount) *MACAddress {
	return addr.init().adjustPrefixLen(prefixLen).ToMAC()
}

// AdjustPrefixLenZeroed increases or decreases the prefix length by the given increment while zeroing out the bits that have moved into or outside the prefix.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the address.
//
// If this address has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
//
// When prefix length is increased, the bits moved within the prefix become zero.
// When a prefix length is decreased, the bits moved outside the prefix become zero.
//
// If the result cannot be zeroed because zeroing out bits results in a non-contiguous segment, an error is returned.
func (addr *MACAddress) AdjustPrefixLenZeroed(prefixLen BitCount) (*MACAddress, addrerr.IncompatibleAddressError) {
	res, err := addr.init().adjustPrefixLenZeroed(prefixLen)
	return res.ToMAC(), err
}

// AssignPrefixForSingleBlock returns the equivalent prefix block that matches exactly the range of values in this address.
// The returned block will have an assigned prefix length indicating the prefix length for the block.
//
// There may be no such address - it is required that the range of values match the range of a prefix block.
// If there is no such address, then nil is returned.
func (addr *MACAddress) AssignPrefixForSingleBlock() *MACAddress {
	return addr.init().assignPrefixForSingleBlock().ToMAC()
}

// AssignMinPrefixForBlock returns an equivalent subnet, assigned the smallest prefix length possible,
// such that the prefix block for that prefix length is in this subnet.
//
// In other words, this method assigns a prefix length to this subnet matching the largest prefix block in this subnet.
func (addr *MACAddress) AssignMinPrefixForBlock() *MACAddress {
	return addr.init().assignMinPrefixForBlock().ToMAC()
}

// ToSinglePrefixBlockOrAddress converts to a single prefix block or address.
// If the given address is a single prefix block, it is returned.
// If it can be converted to a single prefix block by assigning a prefix length, the converted block is returned.
// If it is a single address, any prefix length is removed and the address is returned.
// Otherwise, nil is returned.
// This method provides the address formats used by tries.
// ToSinglePrefixBlockOrAddress is quite similar to AssignPrefixForSingleBlock, which always returns prefixed addresses, while this does not.
func (addr *MACAddress) ToSinglePrefixBlockOrAddress() *MACAddress {
	return addr.init().toSinglePrefixBlockOrAddr().ToMAC()
}

func (addr *MACAddress) toSinglePrefixBlockOrAddress() (*MACAddress, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nil, &incompatibleAddressError{addressError{key: "ipaddress.error.address.not.block", str: addr.String()}}
	}
	res := addr.ToSinglePrefixBlockOrAddress()
	if res == nil {
		return nil, &incompatibleAddressError{addressError{key: "ipaddress.error.address.not.block", str: addr.String()}}
	}
	return res, nil
}

// ContainsPrefixBlock returns whether the range of this address or address collection contains the block of addresses for the given prefix length.
//
// Unlike ContainsSinglePrefixBlock, whether there are multiple prefix values in this item for the given prefix length makes no difference.
//
// Use GetMinPrefixLenForBlock to determine the smallest prefix length for which this method returns true.
func (addr *MACAddress) ContainsPrefixBlock(prefixLen BitCount) bool {
	return addr.init().addressInternal.ContainsPrefixBlock(prefixLen)
}

// ContainsSinglePrefixBlock returns whether this address contains a single prefix block for the given prefix length.
//
// This means there is only one prefix value for the given prefix length, and it also contains the full prefix block for that prefix, all addresses with that prefix.
//
// Use GetPrefixLenForSingleBlock to determine whether there is a prefix length for which this method returns true.
func (addr *MACAddress) ContainsSinglePrefixBlock(prefixLen BitCount) bool {
	return addr.init().addressInternal.ContainsSinglePrefixBlock(prefixLen)
}

// GetMinPrefixLenForBlock returns the smallest prefix length such that this includes the block of addresses for that prefix length.
//
// If the entire range can be described this way, then this method returns the same value as GetPrefixLenForSingleBlock.
//
// There may be a single prefix, or multiple possible prefix values in this item for the returned prefix length.
// Use GetPrefixLenForSingleBlock to avoid the case of multiple prefix values.
//
// If this represents just a single address, returns the bit length of this address.
func (addr *MACAddress) GetMinPrefixLenForBlock() BitCount {
	return addr.init().addressInternal.GetMinPrefixLenForBlock()
}

// GetPrefixLenForSingleBlock returns a prefix length for which the range of this address collection matches the block of addresses for that prefix.
//
// If the range can be described this way, then this method returns the same value as GetMinPrefixLenForBlock.
//
// If no such prefix exists, returns nil.
//
// If this segment grouping represents a single value, this returns the bit length of this address.
func (addr *MACAddress) GetPrefixLenForSingleBlock() PrefixLen {
	return addr.init().addressInternal.GetPrefixLenForSingleBlock()
}

// Compare returns a negative integer, zero, or a positive integer if this address or address collection is less than, equal, or greater than the given item.
// Any address item is comparable to any other.  All address items use CountComparator to compare.
func (addr *MACAddress) Compare(item AddressItem) int {
	return CountComparator.Compare(addr, item)
}

// PrefixEqual determines if the given address matches this address up to the prefix length of this address.
// It returns whether the two addresses share the same range of prefix values.
func (addr *MACAddress) PrefixEqual(other AddressType) bool {
	return addr.init().prefixEquals(other)
}

// PrefixContains returns whether the prefix values in the given address
// are prefix values in this address, using the prefix length of this address.
// If this address has no prefix length, the entire address is compared.
//
// It returns whether the prefix of this address contains all values of the same prefix length in the given address.
func (addr *MACAddress) PrefixContains(other AddressType) bool {
	return addr.init().prefixContains(other)
}

// containsSame returns whether this address contains all addresses in the given address or subnet of the same type.
func (addr *MACAddress) containsSame(other *MACAddress) bool {
	return addr.Contains(other)
}

// Contains returns whether this is the same type and version as the given address or subnet and whether it contains all addresses in the given address or subnet.
func (addr *MACAddress) Contains(other AddressType) bool {
	if addr == nil {
		return other == nil || other.ToAddressBase() == nil
	}
	// note: we don't use the same optimization as in IPv4/6 because we do need to check segment count with MAC
	return addr.init().contains(other)
}

// Overlaps returns true if this address overlaps the given address or address collection
func (addr *MACAddress) Overlaps(other AddressType) bool {
	if addr == nil {
		return true
	}
	return addr.init().overlaps(other)
}

// Equal returns whether the given address or address collection is equal to this address or address collection.
// Two address instances are equal if they represent the same set of addresses.
func (addr *MACAddress) Equal(other AddressType) bool {
	if addr == nil {
		return other == nil || other.ToAddressBase() == nil
	}
	// note: we don't use the same optimization as in IPv4/6 because we do need to check segment count with MAC
	return addr.init().equals(other)
}

// CompareSize compares the counts of two addresses or address collections or address items, the number of individual addresses or items within.
//
// Rather than calculating counts with GetCount, there can be more efficient ways of determining whether one address collection represents more individual addresses than another.
//
// CompareSize returns a positive integer if this address or address collection has a larger count than the one given, zero if they are the same, or a negative integer if the other has a larger count.
func (addr *MACAddress) CompareSize(other AddressItem) int { // this is here to take advantage of the CompareSize in IPAddressSection
	if addr == nil {
		if isNilItem(other) {
			return 0
		}
		// we have size 0, other has size >= 1
		return -1
	}
	return addr.init().compareSize(other)
}

// TrieCompare compares two addresses according to address trie ordering.
// It returns a number less than zero, zero, or a number greater than zero if the first address argument is less than, equal to, or greater than the second.
//
// The comparison is intended for individual addresses and CIDR prefix blocks.
// If an address is neither an individual address nor a prefix block, it is treated like one:
//
//   - ranges that occur inside the prefix length are ignored, only the lower value is used.
//   - ranges beyond the prefix length are assumed to be the full range across all hosts for that prefix length.
func (addr *MACAddress) TrieCompare(other *MACAddress) (int, addrerr.IncompatibleAddressError) {
	if addr.GetSegmentCount() != other.GetSegmentCount() {
		return 0, &incompatibleAddressError{addressError{key: "ipaddress.error.mismatched.bit.size"}}
	}
	return addr.init().trieCompare(other.ToAddressBase()), nil
}

// TrieIncrement returns the next address or block according to address trie ordering
//
// If an address is neither an individual address nor a prefix block, it is treated like one:
//
//   - ranges that occur inside the prefix length are ignored, only the lower value is used.
//   - ranges beyond the prefix length are assumed to be the full range across all hosts for that prefix length.
func (addr *MACAddress) TrieIncrement() *MACAddress {
	if res, ok := trieIncrement(addr); ok {
		return res
	}
	return nil
}

// TrieDecrement returns the previous address or block according to address trie ordering
//
// If an address is neither an individual address nor a prefix block, it is treated like one:
//
//   - ranges that occur inside the prefix length are ignored, only the lower value is used.
//   - ranges beyond the prefix length are assumed to be the full range across all hosts for that prefix length.
func (addr *MACAddress) TrieDecrement() *MACAddress {
	if res, ok := trieDecrement(addr); ok {
		return res
	}
	return nil
}

// GetMaxSegmentValue returns the maximum possible segment value for this type of address.
//
// Note this is not the maximum of the range of segment values in this specific address,
// this is the maximum value of any segment for this address type and version, determined by the number of bits per segment.
func (addr *MACAddress) GetMaxSegmentValue() SegInt {
	return addr.init().getMaxSegmentValue()
}

// IsMulticast returns whether this address or collection of addresses is entirely multicast.
// Multicast MAC addresses have the least significant bit of the first octet set to 1.
func (addr *MACAddress) IsMulticast() bool {
	return addr.GetSegment(0).MatchesWithMask(1, 0x1)
}

// IsUnicast returns whether this address or collection of addresses is entirely unicast.
// Unicast MAC addresses have the least significant bit of the first octet set to 0.
func (addr *MACAddress) IsUnicast() bool {
	return !addr.IsMulticast()
}

// IsUniversal returns whether this is a universal address.
// Universal MAC addresses have second the least significant bit of the first octet set to 0.
func (addr *MACAddress) IsUniversal() bool {
	return !addr.IsLocal()
}

// IsLocal returns whether this is a local address.
// Local MAC addresses have the second least significant bit of the first octet set to 1.
func (addr *MACAddress) IsLocal() bool {
	return addr.GetSegment(0).MatchesWithMask(2, 0x2)
}

// Iterator provides an iterator to iterate through the individual addresses of this address or subnet.
//
// When iterating, the prefix length is preserved.  Remove it using WithoutPrefixLen prior to iterating if you wish to drop it from all individual addresses.
//
// Call IsMultiple to determine if this instance represents multiple addresses, or GetCount for the count.
func (addr *MACAddress) Iterator() Iterator[*MACAddress] {
	if addr == nil {
		return macAddressIterator{nilAddrIterator()}
	}
	return macAddressIterator{addr.init().addrIterator(nil)}
}

// PrefixIterator provides an iterator to iterate through the individual prefixes of this subnet,
// each iterated element spanning the range of values for its prefix.
//
// It is similar to the prefix block iterator, except for possibly the first and last iterated elements, which might not be prefix blocks,
// instead constraining themselves to values from this subnet.
//
// If the subnet has no prefix length, then this is equivalent to Iterator.
func (addr *MACAddress) PrefixIterator() Iterator[*MACAddress] {
	return macAddressIterator{addr.init().prefixIterator(false)}
}

// PrefixBlockIterator provides an iterator to iterate through the individual prefix blocks, one for each prefix of this address or subnet.
// Each iterated address or subnet will be a prefix block with the same prefix length as this address or subnet.
//
// If this address has no prefix length, then this is equivalent to Iterator.
func (addr *MACAddress) PrefixBlockIterator() Iterator[*MACAddress] {
	return macAddressIterator{addr.init().prefixIterator(true)}
}

// BlockIterator iterates through the addresses that can be obtained by iterating through all the upper segments up to the given segment count.
// The segments following remain the same in all iterated addresses.
func (addr *MACAddress) BlockIterator(segmentCount int) Iterator[*MACAddress] {
	return macAddressIterator{addr.init().blockIterator(segmentCount)}
}

// SequentialBlockIterator iterates through the sequential subnets or addresses that make up this address or subnet.
//
// Practically, this means finding the count of segments for which the segments that follow are not full range, and then using BlockIterator with that segment count.
//
// For instance, given the IPv4 subnet "1-2.3-4.5-6.7-8", it will iterate through "1.3.5.7-8", "1.3.6.7-8", "1.4.5.7-8", "1.4.6.7-8", "2.3.5.7-8", "2.3.6.7-8", "2.4.6.7-8" and "2.4.6.7-8".
//
// Use GetSequentialBlockCount to get the number of iterated elements.
func (addr *MACAddress) SequentialBlockIterator() Iterator[*MACAddress] {
	return macAddressIterator{addr.init().sequentialBlockIterator()}
}

// GetSequentialBlockIndex gets the minimal segment index for which all following segments are full-range blocks.
//
// The segment at this index is not a full-range block itself, unless all segments are full-range.
// The segment at this index and all following segments form a sequential range.
// For the full address collection to be sequential, the preceding segments must be single-valued.
func (addr *MACAddress) GetSequentialBlockIndex() int {
	return addr.init().getSequentialBlockIndex()
}

// GetSequentialBlockCount provides the count of elements from the sequential block iterator, the minimal number of sequential address ranges that comprise this address collection.
func (addr *MACAddress) GetSequentialBlockCount() *big.Int {
	return addr.init().getSequentialBlockCount()
}

// IncrementBoundary returns the address that is the given increment from the range boundaries of this address collection.
//
// If the given increment is positive, adds the value to the upper address (GetUpper) in the range to produce a new address.
// If the given increment is negative, adds the value to the lower address (GetLower) in the range to produce a new address.
// If the increment is zero, returns this address.
//
// If this is a single address value, that address is simply incremented by the given increment value, positive or negative.
//
// On address overflow or underflow, IncrementBoundary returns nil.
func (addr *MACAddress) IncrementBoundary(increment int64) *MACAddress {
	return addr.init().incrementBoundary(increment).ToMAC()
}

// Increment returns the address from the address collection that is the given increment upwards into the address range,
// with the increment of 0 returning the first address in the range.
//
// If the increment i matches or exceeds the size count c, then i - c + 1
// is added to the upper address of the range.
// An increment matching the range count gives you the address just above the highest address in the range.
//
// If the increment is negative, it is added to the lower address of the range.
// To get the address just below the lowest address of the address range, use the increment -1.
//
// If this is just a single address value, the address is simply incremented by the given increment, positive or negative.
//
// If this is an address range with multiple values, a positive increment i is equivalent i + 1 values from the iterator and beyond.
// For instance, a increment of 0 is the first value from the iterator, an increment of 1 is the second value from the iterator, and so on.
// An increment of a negative value added to the range count is equivalent to the same number of iterator values preceding the upper bound of the iterator.
// For instance, an increment of count - 1 is the last value from the iterator, an increment of count - 2 is the second last value, and so on.
//
// On address overflow or underflow, Increment returns nil.
func (addr *MACAddress) Increment(increment int64) *MACAddress {
	return addr.init().increment(increment).ToMAC()
}

// Enumerate indicates where an address sits relative to the address collection ordering.
//
// Determines how many address elements of the address collection precede the given address element, if the address is in the address collection.
// If above the address collection range, it is the distance to the upper boundary added to the count less one, and if below the address collection range, the distance to the lower boundary.
//
// In other words, if the given address is not in the address collection but above it, returns the number of addresses preceding the address from the upper range boundary,
// added to one less than the total number of address collection addresses.  If the given address is not in the address collection but below it, returns the number of addresses following the address to the lower address collection boundary.
//
// If the argument is not in the address collection, but neither above nor below the range, then nil is returned.
//
// Enumerate returns nil when the argument is multi-valued. The argument must be an individual address.
//
// When this is also an individual address, the returned value is the distance (difference) between the two addresses.
//
// Enumerate is the inverse of the increment method:
//   - subnet.Enumerate(subnet.Increment(inc)) = inc
//   - subnet.Increment(subnet.Enumerate(newAddr)) = newAddr
//
// If the given address does not have the same MAC address type and size, then nil is returned.
func (addr *MACAddress) Enumerate(other AddressType) *big.Int {
	if other != nil {
		if otherAddr := other.ToAddressBase(); otherAddr != nil {
			return addr.GetSection().Enumerate(otherAddr.GetSection())
		}
	}
	return nil
}

// ReverseBytes returns a new address with the bytes reversed.  Any prefix length is dropped.
func (addr *MACAddress) ReverseBytes() *MACAddress {
	return addr.checkIdentity(addr.GetSection().ReverseBytes())
}

// ReverseBits returns a new address with the bits reversed.  Any prefix length is dropped.
//
// If the bits within a single segment cannot be reversed because the segment represents a range,
// and reversing the segment values results in a range that is not contiguous, this returns an error.
//
// In practice this means that to be reversible, a segment range must include all values except possibly the largest and/or smallest, which reverse to themselves.
//
// If perByte is true, the bits are reversed within each byte, otherwise all the bits are reversed.
func (addr *MACAddress) ReverseBits(perByte bool) (*MACAddress, addrerr.IncompatibleAddressError) {
	res, err := addr.GetSection().ReverseBits(perByte)
	if err != nil {
		return nil, err
	}
	return addr.checkIdentity(res), nil
}

// ReverseSegments returns a new address with the segments reversed.
func (addr *MACAddress) ReverseSegments() *MACAddress {
	return addr.checkIdentity(addr.GetSection().ReverseSegments())
}

// ReplaceLen replaces segments starting from startIndex and ending before endIndex with the same number of segments starting at replacementStartIndex from the replacement section.
// Mappings to or from indices outside the range of this or the replacement address are skipped.
func (addr *MACAddress) ReplaceLen(startIndex, endIndex int, replacement *MACAddress, replacementIndex int) *MACAddress {
	replacementSegCount := replacement.GetSegmentCount()
	if replacementIndex <= 0 {
		startIndex -= replacementIndex
		replacementIndex = 0
	} else if replacementIndex >= replacementSegCount {
		return addr
	}
	// We must do a 1 to 1 adjustment of indices before calling the section replace which would do an adjustment of indices not 1 to 1.
	// Here we assume replacementIndex is 0 and working on the subsection starting at that index.
	// In other words, a replacementIndex of x on the whole section is equivalent to replacementIndex of 0 on the shorter subsection starting at x.
	// Then afterwards we use the original replacement index to work on the whole section again, adjusting as needed.
	startIndex, endIndex, replacementIndexAdjustment := adjust1To1Indices(startIndex, endIndex, addr.GetSegmentCount(), replacementSegCount-replacementIndex)
	if startIndex == endIndex {
		return addr
	}
	replacementIndex += replacementIndexAdjustment
	count := endIndex - startIndex
	return addr.init().checkIdentity(addr.GetSection().ReplaceLen(startIndex, endIndex, replacement.GetSection(), replacementIndex, replacementIndex+count))
}

// Replace replaces segments starting from startIndex with segments from the replacement section.
func (addr *MACAddress) Replace(startIndex int, replacement *MACAddressSection) *MACAddress {
	// We must do a 1 to 1 adjustment of indices before calling the section replace which would do an adjustment of indices not 1 to 1.
	startIndex, endIndex, replacementIndex :=
		adjust1To1Indices(startIndex, startIndex+replacement.GetSegmentCount(), addr.GetSegmentCount(), replacement.GetSegmentCount())
	count := endIndex - startIndex
	return addr.init().checkIdentity(addr.GetSection().ReplaceLen(startIndex, endIndex, replacement, replacementIndex, replacementIndex+count))
}

// GetOUISection returns a section with the first 3 segments, the organizational unique identifier
func (addr *MACAddress) GetOUISection() *MACAddressSection {
	return addr.GetSubSection(0, MACOrganizationalUniqueIdentifierSegmentCount)
}

// GetODISection returns a section with the segments following the first 3 segments, the organizational distinct identifier
func (addr *MACAddress) GetODISection() *MACAddressSection {
	return addr.GetTrailingSection(MACOrganizationalUniqueIdentifierSegmentCount)
}

// ToOUIPrefixBlock returns a section in which the range of values match the full block for the OUI (organizationally unique identifier) bytes
func (addr *MACAddress) ToOUIPrefixBlock() *MACAddress {
	addr = addr.init()
	segmentCount := addr.getDivisionCount()
	currentPref := addr.getPrefixLen()
	newPref := BitCount(MACOrganizationalUniqueIdentifierSegmentCount) << 3 //ouiSegmentCount * MACAddress.BITS_PER_SEGMENT
	createNew := currentPref == nil || currentPref.bitCount() > newPref
	if !createNew {
		newPref = currentPref.bitCount()
		for i := MACOrganizationalUniqueIdentifierSegmentCount; i < segmentCount; i++ {
			segment := addr.GetSegment(i)
			if !segment.IsFullRange() {
				createNew = true
				break
			}
		}
	}
	if !createNew {
		return addr
	}
	segmentIndex := MACOrganizationalUniqueIdentifierSegmentCount
	newSegs := createSegmentArray(segmentCount)
	addr.GetSection().copySubDivisions(0, segmentIndex, newSegs)
	allRangeSegment := allRangeMACSeg.ToDiv()
	for i := segmentIndex; i < segmentCount; i++ {
		newSegs[i] = allRangeSegment
	}
	newSect := createSectionMultiple(newSegs, cacheBitCount(newPref), addr.getAddrType(), true).ToMAC()
	return newMACAddress(newSect)
}

var IPv6LinkLocalPrefix = createLinkLocalPrefix()

func createLinkLocalPrefix() *IPv6AddressSection {
	zeroSeg := zeroIPv6Seg.ToDiv()
	segs := []*AddressDivision{
		NewIPv6Segment(0xfe80).ToDiv(),
		zeroSeg,
		zeroSeg,
		zeroSeg,
	}
	return newIPv6Section(segs)
}

// ToLinkLocalIPv6 converts to a link-local Ipv6 address.  Any MAC prefix length is ignored.  Other elements of this address section are incorporated into the conversion.
// This will provide the latter 4 segments of an IPv6 address, to be paired with the link-local IPv6 prefix of 4 segments.
func (addr *MACAddress) ToLinkLocalIPv6() (*IPv6Address, addrerr.IncompatibleAddressError) {
	sect, err := addr.ToEUI64IPv6()
	if err != nil {
		return nil, err
	}
	return newIPv6Address(IPv6LinkLocalPrefix.Append(sect)), nil
}

// ToEUI64IPv6 converts to an Ipv6 address section.  Any MAC prefix length is ignored.  Other elements of this address section are incorporated into the conversion.
// This will provide the latter 4 segments of an IPv6 address, to be paired with an IPv6 prefix of 4 segments.
func (addr *MACAddress) ToEUI64IPv6() (*IPv6AddressSection, addrerr.IncompatibleAddressError) {
	return NewIPv6SectionFromMAC(addr.init())
}

// IsEUI64 returns whether this section is consistent with an IPv6 EUI64Size section,
// which means it came from an extended 8 byte address,
// and the corresponding segments in the middle match 0xff and 0xff/fe for MAC/not-MAC
func (addr *MACAddress) IsEUI64(asMAC bool) bool {
	if addr.GetSegmentCount() == ExtendedUniqueIdentifier64SegmentCount { //getSegmentCount() == EXTENDED_UNIQUE_IDENTIFIER_64_SEGMENT_COUNT
		section := addr.GetSection()
		seg3 := section.GetSegment(3)
		seg4 := section.GetSegment(4)
		if seg3.matches(0xff) {
			if asMAC {
				return seg4.matches(0xff)
			}
			return seg4.matches(0xfe)
		}
	}
	return false
}

// ToEUI64 converts to IPv6 EUI-64 section
//
// http://standards.ieee.org/develop/regauth/tut/eui64.pdf
//
// If asMAC if true, this address is considered MAC and the EUI-64 is extended using ff-ff, otherwise this address is considered EUI-48 and extended using ff-fe
// Note that IPv6 treats MAC as EUI-48 and extends MAC to IPv6 addresses using ff-fe
func (addr *MACAddress) ToEUI64(asMAC bool) (*MACAddress, addrerr.IncompatibleAddressError) {
	addr = addr.init()
	section := addr.GetSection()
	if addr.GetSegmentCount() == ExtendedUniqueIdentifier48SegmentCount {
		segs := createSegmentArray(ExtendedUniqueIdentifier64SegmentCount)
		section.copySubDivisions(0, 3, segs)
		segs[3] = ffMACSeg.ToDiv()
		if asMAC {
			segs[4] = ffMACSeg.ToDiv()
		} else {
			segs[4] = feMACSeg.ToDiv()
		}
		section.copySubDivisions(3, 6, segs[5:])
		prefixLen := addr.getPrefixLen()
		if prefixLen != nil {
			if prefixLen.bitCount() >= 24 {
				prefixLen = cacheBitCount(prefixLen.bitCount() + (MACBitsPerSegment << 1)) //two segments
			}
		}
		newSect := createInitializedSection(segs, prefixLen, addr.getAddrType()).ToMAC()
		return newMACAddress(newSect), nil
	}
	seg3 := section.GetSegment(3)
	seg4 := section.GetSegment(4)
	if seg3.matches(0xff) {
		if asMAC {
			if seg4.matches(0xff) {
				return addr, nil
			}
		} else {
			if seg4.matches(0xfe) {
				return addr, nil
			}
		}
	}
	return nil, &incompatibleAddressError{addressError{key: "ipaddress.mac.error.not.eui.convertible"}}
}

// String implements the [fmt.Stringer] interface, returning the canonical string provided by ToCanonicalString, or "<nil>" if the receiver is a nil pointer.
func (addr *MACAddress) String() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().addressInternal.toString()
}

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
func (addr MACAddress) Format(state fmt.State, verb rune) {
	addr.init().format(state, verb)
}

// GetSegmentStrings returns a slice with the string for each segment being the string that is normalized with wildcards.
func (addr *MACAddress) GetSegmentStrings() []string {
	if addr == nil {
		return nil
	}
	return addr.init().getSegmentStrings()
}

// ToCanonicalString produces a canonical string for the address.
//
// For MAC, it uses the canonical standardized IEEE 802 MAC address representation of xx-xx-xx-xx-xx-xx.  An example is "01-23-45-67-89-ab".
// For range segments, '|' is used: "11-22-33|44-55-66".
//
// Each MAC address has a unique canonical string.
func (addr *MACAddress) ToCanonicalString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toCanonicalString()
}

// ToNormalizedString produces a normalized string for the address.
//
// For MAC, it differs from the canonical string.  It uses the most common representation of MAC addresses: "xx:xx:xx:xx:xx:xx".  An example is "01:23:45:67:89:ab".
// For range segments, '-' is used: "11:22:33-44:55:66".
//
// Each address has a unique normalized string.
func (addr *MACAddress) ToNormalizedString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toNormalizedString()
}

// ToNormalizedWildcardString produces the normalized string.
func (addr *MACAddress) ToNormalizedWildcardString() string {
	return addr.toNormalizedWildcardString()
}

// ToCompressedString produces a short representation of this address while remaining within the confines of standard representation(s) of the address.
//
// For MAC, it differs from the canonical string.  It produces a shorter string for the address that has no leading zeros.
func (addr *MACAddress) ToCompressedString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toCompressedString()
}

// ToHexString writes this address as a single hexadecimal value (possibly two values if a range),
// the number of digits according to the bit count, with or without a preceding "0x" prefix.
//
// If an address collection cannot be written as a range of two values, an error is returned.
func (addr *MACAddress) ToHexString(with0xPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	return addr.init().toHexString(with0xPrefix)
}

// ToOctalString writes this address as a single octal value (possibly two values if a range),
// the number of digits according to the bit count, with or without a preceding "0" prefix.
//
// If a multiple-valued address collection cannot be written as a single prefix block or a range of two values, an error is returned.
func (addr *MACAddress) ToOctalString(with0Prefix bool) (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	return addr.init().toOctalString(with0Prefix)
}

// ToBinaryString writes this address as a single binary value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0b" prefix.
//
// If an address collection cannot be written as a range of two values, an error is returned.
func (addr *MACAddress) ToBinaryString(with0bPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	return addr.init().toBinaryString(with0bPrefix)
}

// GetDottedAddress returns an AddressDivisionGrouping which organizes the address into segments of bit-length 16, rather than the more typical 8 bits per segment.
//
// If this represents a collection of MAC addresses, this returns an error when unable to join two address segments,
// the first with a range of values, into a division of the larger bit-length that represents the same set of values.
func (addr *MACAddress) GetDottedAddress() (*AddressDivisionGrouping, addrerr.IncompatibleAddressError) {
	return addr.init().GetSection().GetDottedGrouping()
}

// ToDottedString produces the dotted hexadecimal format aaaa.bbbb.cccc
func (addr *MACAddress) ToDottedString() (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	return addr.init().GetSection().ToDottedString()
}

// ToSpaceDelimitedString produces a string delimited by spaces: aa bb cc dd ee ff
func (addr *MACAddress) ToSpaceDelimitedString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().GetSection().ToSpaceDelimitedString()
}

// ToDashedString produces a string delimited by dashes: "aa-bb-cc-dd-ee-ff".
// For range segments, '|' is used: "11-22-33|44-55-66".
// It returns the same string as ToCanonicalString.
func (addr *MACAddress) ToDashedString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().GetSection().ToDashedString()
}

// ToColonDelimitedString produces a string delimited by colons: "aa:bb:cc:dd:ee:ff".
// For range segments, '-' is used: "11:22:33-44:55:66".
// It returns the same string as ToNormalizedString.
func (addr *MACAddress) ToColonDelimitedString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().GetSection().ToColonDelimitedString()
}

// ToCustomString creates a customized string from this address or address collection according to the given string option parameters.
func (addr *MACAddress) ToCustomString(stringOptions addrstr.StringOptions) string {
	if addr == nil {
		return nilString()
	}
	return addr.init().GetSection().toCustomString(stringOptions)
}

// ToAddressString retrieves or generates a MACAddressString instance for this MACAddress instance.
// This may be the MACAddressString this instance was generated from, if it was generated from a MACAddressString.
//
// In general, users are intended to create MACAddress instances from MACAddressString instances,
// while the reverse direction is generally not common and not useful, except under specific circumstances.
//
// However, the reverse direction can be useful under certain circumstances,
// such as when maintaining a collection of MACAddressString instances.
func (addr *MACAddress) ToAddressString() *MACAddressString {
	addr = addr.init()
	cache := addr.cache
	if cache != nil {
		res := addr.cache.identifierStr
		if res != nil {
			hostIdStr := res.idStr
			return hostIdStr.(*MACAddressString)
		}
	}
	return newMACAddressStringFromAddr(addr.toCanonicalString(), addr)
}

func (addr *MACAddress) toMaxLower() *MACAddress {
	return addr.init().addressInternal.toMaxLower().ToMAC()
}

func (addr *MACAddress) toMinUpper() *MACAddress {
	return addr.init().addressInternal.toMinUpper().ToMAC()
}

// ToAddressBase converts to an Address, a polymorphic type usable with all addresses and subnets.
// Afterwards, you can convert back with ToMAC.
//
// ToAddressBase can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (addr *MACAddress) ToAddressBase() *Address {
	if addr != nil {
		addr = addr.init()
	}
	return (*Address)(addr)
}

// Wrap wraps this address, returning a WrappedAddress, an implementation of ExtendedSegmentSeries,
// which can be used to write code that works with both addresses and address sections.
func (addr *MACAddress) Wrap() WrappedAddress {
	return wrapAddress(addr.ToAddressBase())
}

// ToKey creates the associated address key.
// While addresses can be compared with the Compare, TrieCompare or Equal methods as well as various provided instances of AddressComparator,
// they are not comparable with Go operators.
// However, AddressKey instances are comparable with Go operators, and thus can be used as map keys.
func (addr *MACAddress) ToKey() MACAddressKey {
	key := MACAddressKey{
		additionalByteCount: uint8(addr.GetSegmentCount()) - MediaAccessControlSegmentCount,
	}
	section := addr.GetSection()
	divs := section.getDivArray()
	var lowerVal, upperVal uint64
	if addr.IsMultiple() {
		for _, div := range divs {
			seg := div.ToMAC()
			lowerVal = (lowerVal << MACBitsPerSegment) | uint64(seg.GetMACSegmentValue())
			upperVal = (upperVal << MACBitsPerSegment) | uint64(seg.GetMACUpperSegmentValue())
		}
	} else {
		for _, div := range divs {
			seg := div.ToMAC()
			lowerVal = (lowerVal << MACBitsPerSegment) | uint64(seg.GetMACSegmentValue())
		}
		upperVal = lowerVal
	}
	key.vals.lower = lowerVal
	key.vals.upper = upperVal
	return key
}

func fromMACKey(key MACAddressKey) *MACAddress {
	additionalByteCount := key.additionalByteCount
	segCount := int(additionalByteCount) + MediaAccessControlSegmentCount
	return NewMACAddressFromRangeExt(
		func(segmentIndex int) MACSegInt {
			segIndex := (segCount - 1) - segmentIndex
			return MACSegInt(key.vals.lower >> (segIndex << macBitsToSegmentBitshift))
		}, func(segmentIndex int) MACSegInt {
			segIndex := (segCount - 1) - segmentIndex
			return MACSegInt(key.vals.upper >> (segIndex << macBitsToSegmentBitshift))
		},
		additionalByteCount != 0,
	)
}

// ToGenericKey produces a generic Key[*MACAddress] that can be used with generic code working with [Address], [IPAddress], [IPv4Address], [IPv6Address] and [MACAddress].
// ToKey produces a more compact key for code that is MAC-specific.
func (addr *MACAddress) ToGenericKey() Key[*MACAddress] {
	// Note: We intentionally do not populate the "scheme" field for MAC-48.
	// With Key[*IPv4Address], by leaving the scheme zero for MAC-48, the zero Key[*MACAddress] matches up with the key produced here by the zero address.
	// We do not need the scheme field for Key[*MACAddress] since the generic type indicates MAC, but we do need a flag to distinguish 64-bit EUI-64.
	key := Key[*MACAddress]{}
	if isExtended := addr.GetSegmentCount() == ExtendedUniqueIdentifier64SegmentCount; isExtended {
		key.scheme = eui64Scheme
	}
	addr.init().toMACKey(&key.keyContents)
	return key
}

func (addr *MACAddress) fromKey(scheme addressScheme, key *keyContents) *MACAddress {
	// See ToGenericKey for details such as the fact that the scheme is populated only for eui64Scheme
	return fromMACAddrKey(scheme, key)
}

func (addr *MACAddress) toMACKey(contents *keyContents) {
	section := addr.GetSection()
	divs := section.getDivArray()
	if addr.IsMultiple() {
		for i, div := range divs {
			seg := div.ToMAC()
			val := &contents.vals[i>>3]
			val.lower = (val.lower << MACBitsPerSegment) | uint64(seg.GetMACSegmentValue())
			val.upper = (val.upper << MACBitsPerSegment) | uint64(seg.GetMACUpperSegmentValue())
		}
	} else {
		for i, div := range divs {
			seg := div.ToMAC()
			val := &contents.vals[i>>3]
			newLower := (val.lower << MACBitsPerSegment) | uint64(seg.GetMACSegmentValue())
			val.lower = newLower
			val.upper = newLower
		}
	}
}

func fromMACAddrKey(scheme addressScheme, key *keyContents) *MACAddress {
	segCount := MediaAccessControlSegmentCount
	isExtended := false
	// Note: the check here must be for eui64Scheme and not mac48Scheme
	// ToGenericKey will only populate the scheme to eui64Scheme, it will be left as 0 otherwise
	if isExtended = scheme == eui64Scheme; isExtended {
		segCount = ExtendedUniqueIdentifier64SegmentCount
	}
	return NewMACAddressFromRangeExt(
		func(segmentIndex int) MACSegInt {
			valsIndex := segmentIndex >> 3
			segIndex := ((segCount - 1) - segmentIndex) & 0x7
			return MACSegInt(key.vals[valsIndex].lower >> (segIndex << macBitsToSegmentBitshift))
		}, func(segmentIndex int) MACSegInt {
			valsIndex := segmentIndex >> 3
			segIndex := ((segCount - 1) - segmentIndex) & 0x7
			return MACSegInt(key.vals[valsIndex].upper >> (segIndex << macBitsToSegmentBitshift))
		},
		isExtended,
	)
}
