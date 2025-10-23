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
	"net/netip"
	"unsafe"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstr"
)

const (
	IPv4SegmentSeparator      = '.'
	IPv4SegmentSeparatorStr   = "."
	IPv4BitsPerSegment        = 8
	IPv4BytesPerSegment       = 1
	IPv4SegmentCount          = 4
	IPv4ByteCount             = 4
	IPv4BitCount              = 32
	IPv4DefaultTextualRadix   = 10
	IPv4MaxValuePerSegment    = 0xff
	IPv4MaxValue              = 0xffffffff
	IPv4ReverseDnsSuffix      = ".in-addr.arpa"
	IPv4SegmentMaxChars       = 3
	ipv4BitsToSegmentBitshift = 3
)

func newIPv4Address(section *IPv4AddressSection) *IPv4Address {
	return createAddress(section.ToSectionBase(), NoZone).ToIPv4()
}

// NewIPv4Address constructs an IPv4 address or subnet from the given address section.
// If the section does not have 4 segments, an error is returned.
func NewIPv4Address(section *IPv4AddressSection) (*IPv4Address, addrerr.AddressValueError) {
	if section == nil {
		return zeroIPv4, nil
	}
	segCount := section.GetSegmentCount()
	if segCount != IPv4SegmentCount {
		return nil, &addressValueError{
			addressError: addressError{key: "ipaddress.error.invalid.size"},
			val:          segCount,
		}
	}
	return createAddress(section.ToSectionBase(), NoZone).ToIPv4(), nil
}

// NewIPv4AddressFromSegs constructs an IPv4 address or subnet from the given segments.
// If the given slice does not have 4 segments, an error is returned.
func NewIPv4AddressFromSegs(segments []*IPv4AddressSegment) (*IPv4Address, addrerr.AddressValueError) {
	segCount := len(segments)
	if segCount != IPv4SegmentCount {
		return nil, &addressValueError{
			addressError: addressError{key: "ipaddress.error.invalid.size"},
			val:          segCount,
		}
	}
	section := NewIPv4Section(segments)
	return createAddress(section.ToSectionBase(), NoZone).ToIPv4(), nil
}

// NewIPv4AddressFromPrefixedSegs constructs an IPv4 address or subnet from the given segments and prefix length.
// If the given slice does not have 4 segments, an error is returned.
// If the address has a zero host for its prefix length, the returned address will be the prefix block.
func NewIPv4AddressFromPrefixedSegs(segments []*IPv4AddressSegment, prefixLength PrefixLen) (*IPv4Address, addrerr.AddressValueError) {
	segCount := len(segments)
	if segCount != IPv4SegmentCount {
		return nil, &addressValueError{
			addressError: addressError{key: "ipaddress.error.invalid.size"},
			val:          segCount,
		}
	}
	section := NewIPv4PrefixedSection(segments, prefixLength)
	return createAddress(section.ToSectionBase(), NoZone).ToIPv4(), nil
}

// NewIPv4AddressFromBytes constructs an IPv4 address from the given byte slice.
// An error is returned when the byte slice has too many bytes to match the IPv4 segment count of 4.
// There should be 4 bytes or less, although extra leading zeros are tolerated.
func NewIPv4AddressFromBytes(bytes []byte) (addr *IPv4Address, err addrerr.AddressValueError) {
	if ipv4 := net.IP(bytes).To4(); ipv4 != nil {
		bytes = ipv4
	}
	section, err := NewIPv4SectionFromSegmentedBytes(bytes, IPv4SegmentCount)
	if err == nil {
		addr = newIPv4Address(section)
	}
	return
}

// NewIPv4AddressFromPrefixedBytes constructs an IPv4 address or prefix block from the given byte slice and prefix length.
// An error is returned when the byte slice has too many bytes to match the IPv4 segment count of 4.
// There should be 4 bytes or less, although extra leading zeros are tolerated.
// If the address has a zero host for the given prefix length, the returned address will be the prefix block.
func NewIPv4AddressFromPrefixedBytes(bytes []byte, prefixLength PrefixLen) (addr *IPv4Address, err addrerr.AddressValueError) {
	if ipv4 := net.IP(bytes).To4(); ipv4 != nil {
		bytes = ipv4
	}
	section, err := NewIPv4SectionFromPrefixedBytes(bytes, IPv4SegmentCount, prefixLength)
	if err == nil {
		addr = newIPv4Address(section)
	}
	return
}

// NewIPv4AddressFromUint32 constructs an IPv4 address from the given value.
func NewIPv4AddressFromUint32(val uint32) *IPv4Address {
	section := NewIPv4SectionFromUint32(val, IPv4SegmentCount)
	return createAddress(section.ToSectionBase(), NoZone).ToIPv4()
}

// NewIPv4AddressFromPrefixedUint32 constructs an IPv4 address or prefix block from the given value and prefix length.
// If the address has a zero host for the given prefix length, the returned address will be the prefix block.
func NewIPv4AddressFromPrefixedUint32(val uint32, prefixLength PrefixLen) *IPv4Address {
	section := NewIPv4SectionFromPrefixedUint32(val, IPv4SegmentCount, prefixLength)
	return createAddress(section.ToSectionBase(), NoZone).ToIPv4()
}

// NewIPv4AddressFromVals constructs an IPv4 address from the given values.
func NewIPv4AddressFromVals(vals IPv4SegmentValueProvider) *IPv4Address {
	section := NewIPv4SectionFromVals(vals, IPv4SegmentCount)
	return newIPv4Address(section)
}

// NewIPv4AddressFromPrefixedVals constructs an IPv4 address or prefix block from the given values and prefix length.
// If the address has a zero host for the given prefix length, the returned address will be the prefix block.
func NewIPv4AddressFromPrefixedVals(vals IPv4SegmentValueProvider, prefixLength PrefixLen) *IPv4Address {
	section := NewIPv4SectionFromPrefixedVals(vals, IPv4SegmentCount, prefixLength)
	return newIPv4Address(section)
}

// NewIPv4AddressFromRange constructs an IPv4 subnet from the given values.
func NewIPv4AddressFromRange(vals, upperVals IPv4SegmentValueProvider) *IPv4Address {
	section := NewIPv4SectionFromRange(vals, upperVals, IPv4SegmentCount)
	return newIPv4Address(section)
}

// NewIPv4AddressFromPrefixedRange constructs an IPv4 subnet from the given values and prefix length.
// If the address has a zero host for the given prefix length, the returned address will be the prefix block.
func NewIPv4AddressFromPrefixedRange(vals, upperVals IPv4SegmentValueProvider, prefixLength PrefixLen) *IPv4Address {
	section := NewIPv4SectionFromPrefixedRange(vals, upperVals, IPv4SegmentCount, prefixLength)
	return newIPv4Address(section)
}

func newIPv4AddressFromPrefixedSingle(vals, upperVals IPv4SegmentValueProvider, prefixLength PrefixLen) *IPv4Address {
	section := newIPv4SectionFromPrefixedSingle(vals, upperVals, IPv4SegmentCount, prefixLength, true)
	return newIPv4Address(section)
}

var zeroIPv4 = initZeroIPv4()
var ipv4All = zeroIPv4.ToPrefixBlockLen(0)

func initZeroIPv4() *IPv4Address {
	div := zeroIPv4Seg
	segs := []*IPv4AddressSegment{div, div, div, div}
	section := NewIPv4Section(segs)
	return newIPv4Address(section)
}

// IPv4Address is an IPv4 address, or a subnet of multiple IPv4 addresses.
// An IPv4 address is composed of 4 1-byte segments and can optionally have an associated prefix length.
// Each segment can represent a single value or a range of values.
// The zero value is "0.0.0.0".
//
// To construct one from a string, use NewIPAddressString, then use the ToAddress or GetAddress method of [IPAddressString],
// and then use ToIPv4 to get an IPv4Address, assuming the string had an IPv4 format.
//
// For other inputs, use one of the multiple constructor functions like NewIPv4Address.
// You can also use one of the multiple constructors for [IPAddress] like NewIPAddress and then convert using ToIPv4.
type IPv4Address struct {
	ipAddressInternal
}

func (addr *IPv4Address) init() *IPv4Address {
	if addr.section == nil {
		return zeroIPv4
	}
	return addr
}

// GetCount returns the count of addresses that this address or subnet represents.
//
// If just a single address, not a subnet of multiple addresses, returns 1.
//
// For instance, the IP address subnet "1.2.0.0/15" has the count of 2 to the power of 17.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (addr *IPv4Address) GetCount() *big.Int {
	if addr == nil {
		return bigZero()
	}
	return addr.getCount()
}

// GetIPv4Count returns the count of possible distinct values for this section.
// It is the same as GetCount but returns the value as a uint64 instead of a big integer.
// If not representing multiple values, the count is 1.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (addr *IPv4Address) GetIPv4Count() uint64 {
	if addr == nil {
		return 0
	}
	return addr.GetSection().GetIPv4Count()
}

// GetIPv4PrefixCount returns the number of distinct prefix values in this section.
// It is the same as GetPrefixCount but returns the value as a uint64 instead of a big integer.
//
// The prefix length is given by GetPrefixLen.
//
// If this has a non-nil prefix length, returns the number of distinct prefix values.
//
// If this has a nil prefix length, returns the same value as GetIPv4Count.
func (addr *IPv4Address) GetIPv4PrefixCount() uint64 {
	return addr.GetSection().GetIPv4PrefixCount()
}

// GetIPv4PrefixCountLen gives count available as a uint64 instead of big.Int.
//
// It is the similar to GetPrefixCountLen but returns a uint64, not a *big.Int
func (addr *IPv4Address) GetIPv4PrefixCountLen(prefixLength BitCount) uint64 {
	return addr.GetSection().GetIPv4PrefixCountLen(prefixLength)
}

// GetIPv4BlockCount returns the count of distinct values in the given number of initial (more significant) segments.
//
// It is similar to GetBlockCount but returns a uint64 instead of a big integer.
func (addr *IPv4Address) GetIPv4BlockCount(segmentCount int) uint64 {
	return addr.GetSection().GetIPv4BlockCount(segmentCount)
}

// IsMultiple returns true if this represents more than a single individual address, whether it is a subnet of multiple addresses.
func (addr *IPv4Address) IsMultiple() bool {
	return addr != nil && addr.isMultiple()
}

// IsPrefixed returns whether this address has an associated prefix length.
func (addr *IPv4Address) IsPrefixed() bool {
	return addr != nil && addr.isPrefixed()
}

// IsFullRange returns whether this address covers the entire IPv4 address space.
//
// This is true if and only if both IncludesZero and IncludesMax return true.
func (addr *IPv4Address) IsFullRange() bool {
	return addr.GetSection().IsFullRange()
}

// GetBitCount returns the number of bits comprising this address,
// or each address in the range if a subnet, which is 32.
func (addr *IPv4Address) GetBitCount() BitCount {
	return IPv4BitCount
}

// GetByteCount returns the number of bytes required for this address,
// or each address in the range if a subnet, which is 4.
func (addr *IPv4Address) GetByteCount() int {
	return IPv4ByteCount
}

// GetBitsPerSegment returns the number of bits comprising each segment in this address.  Segments in the same address are equal length.
func (addr *IPv4Address) GetBitsPerSegment() BitCount {
	return IPv4BitsPerSegment
}

// GetBytesPerSegment returns the number of bytes comprising each segment in this address or subnet.  Segments in the same address are equal length.
func (addr *IPv4Address) GetBytesPerSegment() int {
	return IPv4BytesPerSegment
}

// GetSection returns the backing section for this address or subnet, comprising all segments.
func (addr *IPv4Address) GetSection() *IPv4AddressSection {
	return addr.init().section.ToIPv4()
}

// GetTrailingSection gets the subsection from the series starting from the given index.
// The first segment is at index 0.
func (addr *IPv4Address) GetTrailingSection(index int) *IPv4AddressSection {
	return addr.GetSection().GetTrailingSection(index)
}

// GetSubSection gets the subsection from the series starting from the given index and ending just before the give endIndex.
// The first segment is at index 0.
func (addr *IPv4Address) GetSubSection(index, endIndex int) *IPv4AddressSection {
	return addr.GetSection().GetSubSection(index, endIndex)
}

// GetNetworkSection returns an address section containing the segments with the network of the address or subnet, the prefix bits.
// The returned section will have only as many segments as needed as determined by the existing CIDR network prefix length.
//
// If this series has no CIDR prefix length, the returned network section will
// be the entire series as a prefixed section with prefix length matching the address bit length.
func (addr *IPv4Address) GetNetworkSection() *IPv4AddressSection {
	return addr.GetSection().GetNetworkSection()
}

// GetNetworkSectionLen returns a section containing the segments with the network of the address or subnet, the prefix bits according to the given prefix length.
// The returned section will have only as many segments as needed to contain the network.
//
// The new section will be assigned the given prefix length,
// unless the existing prefix length is smaller, in which case the existing prefix length will be retained.
func (addr *IPv4Address) GetNetworkSectionLen(prefLen BitCount) *IPv4AddressSection {
	return addr.GetSection().GetNetworkSectionLen(prefLen)
}

// GetHostSection returns a section containing the segments with the host of the address or subnet, the bits beyond the CIDR network prefix length.
// The returned section will have only as many segments as needed to contain the host.
//
// If this series has no prefix length, the returned host section will be the full section.
func (addr *IPv4Address) GetHostSection() *IPv4AddressSection {
	return addr.GetSection().GetHostSection()
}

// GetHostSectionLen returns a section containing the segments with the host of the address or subnet, the bits beyond the given CIDR network prefix length.
// The returned section will have only as many segments as needed to contain the host.
func (addr *IPv4Address) GetHostSectionLen(prefLen BitCount) *IPv4AddressSection {
	return addr.GetSection().GetHostSectionLen(prefLen)
}

// GetNetworkMask returns the network mask associated with the CIDR network prefix length of this address or subnet.
// If this address or subnet has no prefix length, then the all-ones mask is returned.
func (addr *IPv4Address) GetNetworkMask() *IPv4Address {
	var prefLen BitCount
	if pref := addr.getPrefixLen(); pref != nil {
		prefLen = pref.bitCount()
	} else {
		prefLen = IPv4BitCount
	}
	return ipv4Network.GetNetworkMask(prefLen).ToIPv4()
}

// GetHostMask returns the host mask associated with the CIDR network prefix length of this address or subnet.
// If this address or subnet has no prefix length, then the all-ones mask is returned.
func (addr *IPv4Address) GetHostMask() *IPv4Address {
	return addr.getHostMask(ipv4Network).ToIPv4()
}

// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
func (addr *IPv4Address) CopySubSegments(start, end int, segs []*IPv4AddressSegment) (count int) {
	return addr.GetSection().CopySubSegments(start, end, segs)
}

// CopySegments copies the existing segments into the given slice,
// as much as can be fit into the slice, returning the number of segments copied.
func (addr *IPv4Address) CopySegments(segs []*IPv4AddressSegment) (count int) {
	return addr.GetSection().CopySegments(segs)
}

// GetSegments returns a slice with the address segments.  The returned slice is not backed by the same array as this address.
func (addr *IPv4Address) GetSegments() []*IPv4AddressSegment {
	return addr.GetSection().GetSegments()
}

// GetSegment returns the segment at the given index.
// The first segment is at index 0.
// GetSegment will panic given a negative index or an index matching or larger than the segment count.
func (addr *IPv4Address) GetSegment(index int) *IPv4AddressSegment {
	return addr.init().getSegment(index).ToIPv4()
}

// GetSegmentCount returns the segment count, the number of segments in this address, which is 4.
func (addr *IPv4Address) GetSegmentCount() int {
	return addr.GetDivisionCount()
}

// ForEachSegment visits each segment in order from most-significant to least, the most significant with index 0, calling the given function for each, terminating early if the function returns true.
// Returns the number of visited segments.
func (addr *IPv4Address) ForEachSegment(consumer func(segmentIndex int, segment *IPv4AddressSegment) (stop bool)) int {
	return addr.GetSection().ForEachSegment(consumer)
}

// GetGenericDivision returns the segment at the given index as a DivisionType.
func (addr *IPv4Address) GetGenericDivision(index int) DivisionType {
	return addr.init().getDivision(index)
}

// GetGenericSegment returns the segment at the given index as an AddressSegmentType.
// The first segment is at index 0.
// GetGenericSegment will panic given a negative index or an index matching or larger than the segment count.
func (addr *IPv4Address) GetGenericSegment(index int) AddressSegmentType {
	return addr.init().getSegment(index)
}

// GetDivisionCount returns the segment count.
func (addr *IPv4Address) GetDivisionCount() int {
	return addr.init().getDivisionCount()
}

// GetIPVersion returns IPv4, the IP version of this address.
func (addr *IPv4Address) GetIPVersion() IPVersion {
	return IPv4
}

func (addr *IPv4Address) checkIdentity(section *IPv4AddressSection) *IPv4Address {
	if section == nil {
		return nil
	}
	sec := section.ToSectionBase()
	if sec == addr.section {
		return addr
	}
	return newIPv4Address(section)
}

// Mask applies the given mask to all addresses represented by this IPv4Address.
// The mask is applied to all individual addresses.
//
// If this represents multiple addresses, and applying the mask to all addresses creates a set of addresses
// that cannot be represented as a sequential range within each segment, then an error is returned.
func (addr *IPv4Address) Mask(other *IPv4Address) (masked *IPv4Address, err addrerr.IncompatibleAddressError) {
	return addr.maskPrefixed(other, true)
}

func (addr *IPv4Address) maskPrefixed(other *IPv4Address, retainPrefix bool) (masked *IPv4Address, err addrerr.IncompatibleAddressError) {
	addr = addr.init()
	sect, err := addr.GetSection().maskPrefixed(other.GetSection(), retainPrefix)
	if err == nil {
		masked = addr.checkIdentity(sect)
	}
	return
}

// BitwiseOr does the bitwise disjunction with this address or subnet, useful when subnetting.
// It is similar to Mask which does the bitwise conjunction.
//
// The operation is applied to all individual addresses and the result is returned.
//
// If this is a subnet representing multiple addresses, and applying the operation to all addresses creates a set of addresses
// that cannot be represented as a sequential range within each segment, then an error is returned.
func (addr *IPv4Address) BitwiseOr(other *IPv4Address) (masked *IPv4Address, err addrerr.IncompatibleAddressError) {
	return addr.bitwiseOrPrefixed(other, true)
}

func (addr *IPv4Address) bitwiseOrPrefixed(other *IPv4Address, retainPrefix bool) (masked *IPv4Address, err addrerr.IncompatibleAddressError) {
	addr = addr.init()
	sect, err := addr.GetSection().bitwiseOrPrefixed(other.GetSection(), retainPrefix)
	if err == nil {
		masked = addr.checkIdentity(sect)
	}
	return
}

// Subtract subtracts the given subnet from this subnet, returning an array of subnets for the result (the subnets will not be contiguous so an array is required).
// Subtract computes the subnet difference, the set of addresses in this address subnet but not in the provided subnet.
// This is also known as the relative complement of the given argument in this subnet.
// This is set subtraction, not subtraction of address values (use Increment for the latter).  We have a subnet of addresses and we are removing those addresses found in the argument subnet.
// If there are no remaining addresses, nil is returned.
func (addr *IPv4Address) Subtract(other *IPv4Address) []*IPv4Address {
	addr = addr.init()
	sects, _ := addr.GetSection().Subtract(other.GetSection())
	sectLen := len(sects)
	if sectLen == 0 {
		return nil
	} else if sectLen == 1 {
		sec := sects[0]
		if sec.ToSectionBase() == addr.section {
			return []*IPv4Address{addr}
		}
	}
	res := make([]*IPv4Address, sectLen)
	for i, sect := range sects {
		res[i] = newIPv4Address(sect)
	}
	return res
}

// Intersect returns the subnet whose addresses are found in both this and the given subnet argument, or nil if no such addresses exist.
//
// This is also known as the conjunction of the two sets of addresses.
func (addr *IPv4Address) Intersect(other *IPv4Address) *IPv4Address {
	addr = addr.init()
	section, _ := addr.GetSection().Intersect(other.GetSection())
	if section == nil {
		return nil
	}
	return addr.checkIdentity(section)
}

// SpanWithRange returns an IPv4AddressSeqRange instance that spans this subnet to the given subnet.
// If the other address is a different version than this, then the other is ignored, and the result is equivalent to calling ToSequentialRange.
func (addr *IPv4Address) SpanWithRange(other *IPv4Address) *SequentialRange[*IPv4Address] {
	return NewSequentialRange(addr.init(), other)
}

// GetLower returns the lowest address in the subnet range,
// which will be the receiver if it represents a single address.
// For example, for "1.2-3.4.5-6", the series "1.2.4.5" is returned.
func (addr *IPv4Address) GetLower() *IPv4Address {
	return addr.init().getLower().ToIPv4()
}

// GetUpper returns the highest address in the subnet range,
// which will be the receiver if it represents a single address.
// For example, for "1.2-3.4.5-6", the address "1.3.4.6" is returned.
func (addr *IPv4Address) GetUpper() *IPv4Address {
	return addr.init().getUpper().ToIPv4()
}

// GetLowerIPAddress returns the address in the subnet or address collection with the lowest numeric value,
// which will be the receiver if it represents a single address.
// For example, for "1.2-3.4.5-6", the series "1.2.4.5" is returned.
// GetLowerIPAddress implements the IPAddressRange interface
func (addr *IPv4Address) GetLowerIPAddress() *IPAddress {
	return addr.GetLower().ToIP()
}

// GetUpperIPAddress returns the address in the subnet or address collection with the highest numeric value,
// which will be the receiver if it represents a single address.
// For example, for the subnet "1.2-3.4.5-6", the address "1.3.4.6" is returned.
// GetUpperIPAddress implements the IPAddressRange interface
func (addr *IPv4Address) GetUpperIPAddress() *IPAddress {
	return addr.GetUpper().ToIP()
}

// IsZeroHostLen returns whether the host section is always zero for all individual addresses in this subnet,
// for the given prefix length.
//
// If the host section is zero length (there are zero host bits), IsZeroHostLen returns true.
func (addr *IPv4Address) IsZeroHostLen(prefLen BitCount) bool {
	return addr.init().isZeroHostLen(prefLen)
}

// ToZeroHost converts the address or subnet to one in which all individual addresses have a host of zero,
// the host being the bits following the prefix length.
// If the address or subnet has no prefix length, then it returns an all-zero address.
//
// The returned address or subnet will have the same prefix and prefix length.
//
// For instance, the zero host of "1.2.3.4/16" is the individual address "1.2.0.0/16".
//
// This returns an error if the subnet is a range of addresses which cannot be converted to a range in which all addresses have zero hosts,
// because the conversion results in a subnet segment that is not a sequential range of values.
func (addr *IPv4Address) ToZeroHost() (*IPv4Address, addrerr.IncompatibleAddressError) {
	res, err := addr.init().toZeroHost(false)
	return res.ToIPv4(), err
}

// ToZeroHostLen converts the address or subnet to one in which all individual addresses have a host of zero,
// the host being the bits following the given prefix length.
// If this address or subnet has the same prefix length, then the returned one will too, otherwise the returned series will have no prefix length.
//
// For instance, the zero host of "1.2.3.4" for the prefix length of 16 is the address "1.2.0.0".
//
// This returns an error if the subnet is a range of addresses which cannot be converted to a range in which all addresses have zero hosts,
// because the conversion results in a subnet segment that is not a sequential range of values.
func (addr *IPv4Address) ToZeroHostLen(prefixLength BitCount) (*IPv4Address, addrerr.IncompatibleAddressError) {
	res, err := addr.init().toZeroHostLen(prefixLength)
	return res.ToIPv4(), err
}

// ToZeroNetwork converts the address or subnet to one in which all individual addresses have a network of zero,
// the network being the bits within the prefix length.
// If the address or subnet has no prefix length, then it returns an all-zero address.
//
// The returned address or subnet will have the same prefix length.
func (addr *IPv4Address) ToZeroNetwork() *IPv4Address {
	return addr.init().toZeroNetwork().ToIPv4()
}

// IsMaxHostLen returns whether the host is all one-bits, the max value, for all individual addresses in this subnet,
// for the given prefix length, the host being the bits following the prefix.
//
// If the host section is zero length (there are zero host bits), IsMaxHostLen returns true.
func (addr *IPv4Address) IsMaxHostLen(prefLen BitCount) bool {
	return addr.init().isMaxHostLen(prefLen)
}

// ToMaxHost converts the address or subnet to one in which all individual addresses have a host of all one-bits, the max value,
// the host being the bits following the prefix length.
// If the address or subnet has no prefix length, then it returns an all-ones address, the max address.
//
// The returned address or subnet will have the same prefix and prefix length.
//
// For instance, the max host of "1.2.3.4/16" gives the broadcast address "1.2.255.255/16".
//
// This returns an error if the subnet is a range of addresses which cannot be converted to a range in which all addresses have max hosts,
// because the conversion results in a subnet segment that is not a sequential range of values.
func (addr *IPv4Address) ToMaxHost() (*IPv4Address, addrerr.IncompatibleAddressError) {
	res, err := addr.init().toMaxHost()
	return res.ToIPv4(), err
}

// ToMaxHostLen converts the address or subnet to one in which all individual addresses have a host of all one-bits, the max host,
// the host being the bits following the given prefix length.
// If this address or subnet has the same prefix length, then the resulting one will too, otherwise the resulting address or subnet will have no prefix length.
//
// For instance, the zero host of "1.2.3.4" for the prefix length of 16 is the address "1.2.255.255".
//
// This returns an error if the subnet is a range of addresses which cannot be converted to a range in which all addresses have max hosts,
// because the conversion results in a subnet segment that is not a sequential range of values.
func (addr *IPv4Address) ToMaxHostLen(prefixLength BitCount) (*IPv4Address, addrerr.IncompatibleAddressError) {
	res, err := addr.init().toMaxHostLen(prefixLength)
	return res.ToIPv4(), err
}

// ToPrefixBlock returns the subnet associated with the prefix length of this address.
// If this address has no prefix length, this address is returned.
//
// The subnet will include all addresses with the same prefix as this one, the prefix "block".
// The network prefix will match the prefix of this address or subnet, and the host values will span all values.
//
// For example, if the address is "1.2.3.4/16" it returns the subnet "1.2.0.0/16", which can also be written as "1.2.*.*/16".
func (addr *IPv4Address) ToPrefixBlock() *IPv4Address {
	return addr.init().toPrefixBlock().ToIPv4()
}

// ToPrefixBlockLen returns the subnet associated with the given prefix length.
//
// The subnet will include all addresses with the same prefix as this one, the prefix "block" for that prefix length.
// The network prefix will match the prefix of this address or subnet, and the host values will span all values.
//
// For example, if the address is "1.2.3.4" and the prefix length provided is 16, it returns the subnet "1.2.0.0/16", which can also be written as "1.2.*.*/16".
func (addr *IPv4Address) ToPrefixBlockLen(prefLen BitCount) *IPv4Address {
	return addr.init().toPrefixBlockLen(prefLen).ToIPv4()
}

// ToBlock creates a new block of addresses by changing the segment at the given index to have the given lower and upper value,
// and changing the following segments to be full-range.
func (addr *IPv4Address) ToBlock(segmentIndex int, lower, upper SegInt) *IPv4Address {
	return addr.init().toBlock(segmentIndex, lower, upper).ToIPv4()
}

// WithoutPrefixLen provides the same address but with no prefix length.  The values remain unchanged.
func (addr *IPv4Address) WithoutPrefixLen() *IPv4Address {
	if !addr.IsPrefixed() {
		return addr
	}
	return addr.init().withoutPrefixLen().ToIPv4()
}

// SetPrefixLen sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the address.
// The provided prefix length will be adjusted to these boundaries if necessary.
func (addr *IPv4Address) SetPrefixLen(prefixLen BitCount) *IPv4Address {
	return addr.init().setPrefixLen(prefixLen).ToIPv4()
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

func (addr *IPv4Address) SetPrefixLenZeroed(prefixLen BitCount) (*IPv4Address, addrerr.IncompatibleAddressError) {
	res, err := addr.init().setPrefixLenZeroed(prefixLen)
	return res.ToIPv4(), err
}

// AdjustPrefixLen increases or decreases the prefix length by the given increment.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the address.
//
// If this address has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
func (addr *IPv4Address) AdjustPrefixLen(prefixLen BitCount) *IPv4Address {
	return addr.init().adjustPrefixLen(prefixLen).ToIPv4()
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
// For example, "1.2.0.0/16" adjusted by -8 becomes "1.0.0.0/8".
// "1.2.0.0/16" adjusted by 8 becomes "1.2.0.0/24".
//
// If the result cannot be zeroed because zeroing out bits results in a non-contiguous segment, an error is returned.
func (addr *IPv4Address) AdjustPrefixLenZeroed(prefixLen BitCount) (*IPv4Address, addrerr.IncompatibleAddressError) {
	res, err := addr.init().adjustPrefixLenZeroed(prefixLen)
	return res.ToIPv4(), err
}

// AssignPrefixForSingleBlock returns the equivalent prefix block that matches exactly the range of values in this address.
// The returned block will have an assigned prefix length indicating the prefix length for the block.
//
// There may be no such address - it is required that the range of values match the range of a prefix block.
// If there is no such address, then nil is returned.
//
// Examples:
//   - 1.2.3.4 returns 1.2.3.4/32
//   - 1.2.*.* returns 1.2.0.0/16
//   - 1.2.*.0/24 returns 1.2.0.0/16
//   - 1.2.*.4 returns nil
//   - 1.2.0-1.* returns 1.2.0.0/23
//   - 1.2.1-2.* returns nil
//   - 1.2.252-255.* returns 1.2.252.0/22
//   - 1.2.3.4/16 returns 1.2.3.4/32
func (addr *IPv4Address) AssignPrefixForSingleBlock() *IPv4Address {
	return addr.init().assignPrefixForSingleBlock().ToIPv4()
}

// AssignMinPrefixForBlock returns an equivalent subnet, assigned the smallest prefix length possible,
// such that the prefix block for that prefix length is in this subnet.
//
// In other words, this method assigns a prefix length to this subnet matching the largest prefix block in this subnet.
//
// Examples:
//   - 1.2.3.4 returns 1.2.3.4/32
//   - 1.2.*.* returns 1.2.0.0/16
//   - 1.2.*.0/24 returns 1.2.0.0/16
//   - 1.2.*.4 returns 1.2.*.4/32
//   - 1.2.0-1.* returns 1.2.0.0/23
//   - 1.2.1-2.* returns 1.2.1-2.0/24
//   - 1.2.252-255.* returns 1.2.252.0/22
//   - 1.2.3.4/16 returns 1.2.3.4/32
func (addr *IPv4Address) AssignMinPrefixForBlock() *IPv4Address {
	return addr.init().assignMinPrefixForBlock().ToIPv4()
}

// ToSinglePrefixBlockOrAddress converts to a single prefix block or address.
// If the given address is a single prefix block, it is returned.
// If it can be converted to a single prefix block by assigning a prefix length, the converted block is returned.
// If it is a single address, any prefix length is removed and the address is returned.
// Otherwise, nil is returned.
// This method provides the address formats used by tries.
// ToSinglePrefixBlockOrAddress is quite similar to AssignPrefixForSingleBlock, which always returns prefixed addresses, while this does not.
func (addr *IPv4Address) ToSinglePrefixBlockOrAddress() *IPv4Address {
	return addr.init().toSinglePrefixBlockOrAddr().ToIPv4()
}

func (addr *IPv4Address) toSinglePrefixBlockOrAddress() (*IPv4Address, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nil, &incompatibleAddressError{addressError{key: "ipaddress.error.address.not.block", str: addr.String()}}
	}
	res := addr.ToSinglePrefixBlockOrAddress()
	if res == nil {
		return nil, &incompatibleAddressError{addressError{key: "ipaddress.error.address.not.block", str: addr.String()}}
	}
	return res, nil
}

// ContainsPrefixBlock returns whether the range of this address or subnet contains the block of addresses for the given prefix length.
//
// Unlike ContainsSinglePrefixBlock, whether there are multiple prefix values in this item for the given prefix length makes no difference.
//
// Use GetMinPrefixLenForBlock to determine the smallest prefix length for which this method returns true.
func (addr *IPv4Address) ContainsPrefixBlock(prefixLen BitCount) bool {
	return addr.init().ipAddressInternal.ContainsPrefixBlock(prefixLen)
}

// ContainsSinglePrefixBlock returns whether this address contains a single prefix block for the given prefix length.
//
// This means there is only one prefix value for the given prefix length, and it also contains the full prefix block for that prefix, all addresses with that prefix.
//
// Use GetPrefixLenForSingleBlock to determine whether there is a prefix length for which this method returns true.
func (addr *IPv4Address) ContainsSinglePrefixBlock(prefixLen BitCount) bool {
	return addr.init().ipAddressInternal.ContainsSinglePrefixBlock(prefixLen)
}

// GetMinPrefixLenForBlock returns the smallest prefix length such that this includes the block of addresses for that prefix length.
//
// If the entire range can be described this way, then this method returns the same value as GetPrefixLenForSingleBlock.
//
// There may be a single prefix, or multiple possible prefix values in this item for the returned prefix length.
// Use GetPrefixLenForSingleBlock to avoid the case of multiple prefix values.
//
// If this represents just a single address, returns the bit length of this address.
func (addr *IPv4Address) GetMinPrefixLenForBlock() BitCount {
	return addr.init().ipAddressInternal.GetMinPrefixLenForBlock()
}

// GetPrefixLenForSingleBlock returns a prefix length for which the range of this address subnet matches exactly the block of addresses for that prefix.
//
// If the range can be described this way, then this method returns the same value as GetMinPrefixLenForBlock.
//
// If no such prefix exists, returns nil.
//
// If this segment grouping represents a single value, returns the bit length of this address division series.
//
// Examples:
//   - 1.2.3.4 returns 32
//   - 1.2.3.4/16 returns 32
//   - 1.2.*.* returns 16
//   - 1.2.*.0/24 returns 16
//   - 1.2.0.0/16 returns 16
//   - 1.2.*.4 returns nil
//   - 1.2.252-255.* returns 22
func (addr *IPv4Address) GetPrefixLenForSingleBlock() PrefixLen {
	return addr.init().ipAddressInternal.GetPrefixLenForSingleBlock()
}

// GetValue returns the lowest address in this subnet or address as an integer value.
func (addr *IPv4Address) GetValue() *big.Int {
	return addr.init().section.GetValue()
}

// GetUpperValue returns the highest address in this subnet or address as an integer value.
func (addr *IPv4Address) GetUpperValue() *big.Int {
	return addr.init().section.GetUpperValue()
}

// Uint32Value returns the lowest address in the subnet range as a uint32.
func (addr *IPv4Address) Uint32Value() uint32 {
	return addr.GetSection().Uint32Value()
}

// UpperUint32Value returns the highest address in the subnet range as a uint32.
func (addr *IPv4Address) UpperUint32Value() uint32 {
	return addr.GetSection().UpperUint32Value()
}

// GetNetIPAddr returns the lowest address in this subnet or address as a net.IPAddr.
func (addr *IPv4Address) GetNetIPAddr() *net.IPAddr {
	return &net.IPAddr{
		IP: addr.GetNetIP(),
	}
}

// GetUpperNetIPAddr returns the highest address in this subnet or address as a net.IPAddr.
func (addr *IPv4Address) GetUpperNetIPAddr() *net.IPAddr {
	return &net.IPAddr{
		IP: addr.GetUpperNetIP(),
	}
}

// GetNetIP returns the lowest address in this subnet or address as a net.IP.
func (addr *IPv4Address) GetNetIP() net.IP {
	return addr.Bytes()
}

// GetUpperNetIP returns the highest address in this subnet or address as a net.IP.
func (addr *IPv4Address) GetUpperNetIP() net.IP {
	return addr.UpperBytes()
}

// GetNetNetIPAddr returns the lowest address in this subnet or address range as a netip.Addr.
func (addr *IPv4Address) GetNetNetIPAddr() netip.Addr {
	return addr.init().getNetNetIPAddr()
}

// GetUpperNetNetIPAddr returns the highest address in this subnet or address range as a netip.Addr.
func (addr *IPv4Address) GetUpperNetNetIPAddr() netip.Addr {
	return addr.init().getUpperNetNetIPAddr()
}

// CopyNetIP copies the value of the lowest individual address in the subnet into a net.IP.
//
// If the value can fit in the given net.IP slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (addr *IPv4Address) CopyNetIP(ip net.IP) net.IP {
	if ipv4 := ip.To4(); ipv4 != nil { // this shrinks the arg to 4 bytes if it was 16
		ip = ipv4
	}
	return addr.CopyBytes(ip)
}

// CopyUpperNetIP copies the value of the highest individual address in the subnet into a net.IP.
//
// If the value can fit in the given net.IP slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (addr *IPv4Address) CopyUpperNetIP(ip net.IP) net.IP {
	if ipv4 := ip.To4(); ipv4 != nil { // this shrinks the arg to 4 bytes if it was 16
		ip = ipv4
	}
	return addr.CopyUpperBytes(ip)
}

// Bytes returns the lowest address in this subnet or address as a byte slice.
func (addr *IPv4Address) Bytes() []byte {
	return addr.init().section.Bytes()
}

// UpperBytes returns the highest address in this subnet or address as a byte slice.
func (addr *IPv4Address) UpperBytes() []byte {
	return addr.init().section.UpperBytes()
}

// CopyBytes copies the value of the lowest individual address in the subnet into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (addr *IPv4Address) CopyBytes(bytes []byte) []byte {
	return addr.init().section.CopyBytes(bytes)
}

// CopyUpperBytes copies the value of the highest individual address in the subnet into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (addr *IPv4Address) CopyUpperBytes(bytes []byte) []byte {
	return addr.init().section.CopyUpperBytes(bytes)
}

// IsMax returns whether this address matches exactly the maximum possible value, the address whose bits are all ones.
func (addr *IPv4Address) IsMax() bool {
	return addr.init().section.IsMax()
}

// IncludesMax returns whether this address includes the max address, the address whose bits are all ones, within its range.
func (addr *IPv4Address) IncludesMax() bool {
	return addr.init().section.IncludesMax()
}

// TestBit returns true if the bit in the lower value of this address at the given index is 1, where index 0 refers to the least significant bit.
// In other words, it computes (bits & (1 << n)) != 0), using the lower value of this address.
// TestBit will panic if n < 0, or if it matches or exceeds the bit count of this item.
func (addr *IPv4Address) TestBit(n BitCount) bool {
	return addr.init().testBit(n)
}

// IsOneBit returns true if the bit in the lower value of this address at the given index is 1, where index 0 refers to the most significant bit.
// IsOneBit will panic if bitIndex is less than zero, or if it is larger than the bit count of this item.
func (addr *IPv4Address) IsOneBit(bitIndex BitCount) bool {
	return addr.init().isOneBit(bitIndex)
}

// PrefixEqual determines if the given address matches this address up to the prefix length of this address.
// It returns whether the two addresses share the same range of prefix values.
func (addr *IPv4Address) PrefixEqual(other AddressType) bool {
	return addr.init().prefixEquals(other)
}

// PrefixContains returns whether the prefix values in the given address or subnet
// are prefix values in this address or subnet, using the prefix length of this address or subnet.
// If this address has no prefix length, the entire address is compared.
//
// It returns whether the prefix of this address contains all values of the same prefix length in the given address.
func (addr *IPv4Address) PrefixContains(other AddressType) bool {
	return addr.init().prefixContains(other)
}

// containsSame returns whether this address contains all addresses in the given address or subnet of the same type.
func (addr *IPv4Address) containsSame(other *IPv4Address) bool {
	return addr.Contains(other)
}

// Contains returns whether this is the same type and version as the given address or subnet and whether it contains all addresses in the given address or subnet.
func (addr *IPv4Address) Contains(other AddressType) bool {
	if other == nil || other.ToAddressBase() == nil {
		return true
	} else if addr == nil {
		return false
	}
	addr = addr.init()
	otherAddr := other.ToAddressBase() // runs init before calling getAddrType below
	if addr.ToAddressBase() == otherAddr {
		return true
	}
	return otherAddr.getAddrType() == ipv4Type && addr.section.sameCountTypeContains(otherAddr.GetSection())
}

// ContainsRange returns true if this address contains the given sequential range
func (addr *IPv4Address) ContainsRange(other IPAddressSeqRangeType) bool {
	return isContainedBy(other, addr.ToIP())
}

// Overlaps returns true if this address overlaps the given address or subnet
func (addr *IPv4Address) Overlaps(other AddressType) bool {
	if addr == nil {
		return true
	}
	return addr.init().overlaps(other)
}

// OverlapsRange returns true if this address overlaps the given sequential range
func (addr *IPv4Address) OverlapsRange(other IPAddressSeqRangeType) bool {
	if other == nil {
		return true
	}
	return other.OverlapsAddress(addr)
}

// Compare returns a negative integer, zero, or a positive integer if this address or subnet is less than, equal, or greater than the given item.
// Any address item is comparable to any other.
func (addr *IPv4Address) Compare(item AddressItem) int {
	return CountComparator.Compare(addr, item)
}

// Equal returns whether the given address or subnet is equal to this address or subnet.
// Two address instances are equal if they represent the same set of addresses.
func (addr *IPv4Address) Equal(other AddressType) bool {
	if addr == nil {
		return other == nil || other.ToAddressBase() == nil
	} else if other.ToAddressBase() == nil {
		return false
	}
	return other.ToAddressBase().getAddrType() == ipv4Type && addr.init().section.sameCountTypeEquals(other.ToAddressBase().GetSection())
}

// CompareSize compares the counts of two subnets or addresses or other items, the number of individual addresses or items within.
//
// Rather than calculating counts with GetCount, there can be more efficient ways of determining whether this subnet represents more individual addresses than another item.
//
// CompareSize returns a positive integer if this address or subnet has a larger count than the one given, zero if they are the same, or a negative integer if the other has a larger count.
func (addr *IPv4Address) CompareSize(other AddressItem) int {
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
func (addr *IPv4Address) TrieCompare(other *IPv4Address) int {
	return addr.init().trieCompare(other.ToAddressBase())
}

// TrieIncrement returns the next address or block according to address trie ordering
//
// If an address is neither an individual address nor a prefix block, it is treated like one:
//
//   - ranges that occur inside the prefix length are ignored, only the lower value is used.
//   - ranges beyond the prefix length are assumed to be the full range across all hosts for that prefix length.
func (addr *IPv4Address) TrieIncrement() *IPv4Address {
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
func (addr *IPv4Address) TrieDecrement() *IPv4Address {
	if res, ok := trieDecrement(addr); ok {
		return res
	}
	return nil
}

// MatchesWithMask applies the mask to this address and then compares the result with the given address,
// returning true if they match, false otherwise.
func (addr *IPv4Address) MatchesWithMask(other *IPv4Address, mask *IPv4Address) bool {
	return addr.init().GetSection().MatchesWithMask(other.GetSection(), mask.GetSection())
}

// GetMaxSegmentValue returns the maximum possible segment value for this type of address.
//
// Note this is not the maximum of the range of segment values in this specific address,
// this is the maximum value of any segment for this address type and version, determined by the number of bits per segment.
func (addr *IPv4Address) GetMaxSegmentValue() SegInt {
	return addr.init().getMaxSegmentValue()
}

// ToSequentialRange creates a sequential range instance from the lowest and highest addresses in this subnet.
//
// The two will represent the same set of individual addresses if and only if IsSequential is true.
// To get a series of ranges that represent the same set of individual addresses use the SequentialBlockIterator (or PrefixIterator),
// and apply this method to each iterated subnet.
//
// If this represents just a single address then the returned instance covers just that single address as well.
func (addr *IPv4Address) ToSequentialRange() *SequentialRange[*IPv4Address] {
	if addr == nil {
		return nil
	}
	addr = addr.init().WithoutPrefixLen()
	return newSequRangeUnchecked(
		addr.GetLower(),
		addr.GetUpper(),
		addr.isMultiple())
}

func (addr *IPv4Address) getLowestHighestAddrs() (lower, upper *IPv4Address) {
	l, u := addr.ipAddressInternal.getLowestHighestAddrs()
	return l.ToIPv4(), u.ToIPv4()
}

// ToBroadcastAddress returns the IPv4 broadcast address.
// The broadcast address has the same prefix but a host that is all 1 bits.
// If this address or subnet is not prefixed, this returns the address of all 1 bits, the "max" address.
// This returns an error if a prefixed and ranged-valued segment cannot be converted to a host of all ones and remain a range of consecutive values.
func (addr *IPv4Address) ToBroadcastAddress() (*IPv4Address, addrerr.IncompatibleAddressError) {
	return addr.ToMaxHost()
}

// ToNetworkAddress returns the IPv4 network address.
// The network address has the same prefix but a zero host.
// If this address or subnet is not prefixed, this returns the zero "any" address.
// This returns an error if a prefixed and ranged-valued segment cannot be converted to a host of all zeros and remain a range of consecutive values.
func (addr *IPv4Address) ToNetworkAddress() (*IPv4Address, addrerr.IncompatibleAddressError) {
	return addr.ToZeroHost()
}

// ToAddressString retrieves or generates an IPAddressString instance for this IPAddress instance.
// This may be the IPAddressString this instance was generated from, if it was generated from an IPAddressString.
//
// In general, users are intended to create IPAddress instances from IPAddressString instances,
// while the reverse direction is generally not common and not useful, except under specific circumstances.
//
// However, the reverse direction can be useful under certain circumstances,
// such as when maintaining a collection of HostIdentifierString instances.
func (addr *IPv4Address) ToAddressString() *IPAddressString {
	return addr.init().ToIP().ToAddressString()
}

// IncludesZeroHostLen returns whether the subnet contains an individual address with a host of zero, an individual address for which all bits past the given prefix length are zero.
func (addr *IPv4Address) IncludesZeroHostLen(networkPrefixLength BitCount) bool {
	return addr.init().includesZeroHostLen(networkPrefixLength)
}

// IncludesMaxHostLen returns whether the subnet contains an individual address with a host of all one-bits, an individual address for which all bits past the given prefix length are all ones.
func (addr *IPv4Address) IncludesMaxHostLen(networkPrefixLength BitCount) bool {
	return addr.init().includesMaxHostLen(networkPrefixLength)
}

// IsLinkLocal returns whether the address is link local, whether unicast or multicast.
func (addr *IPv4Address) IsLinkLocal() bool {
	if addr.IsMulticast() {
		//224.0.0.252	Link-local Multicast Name Resolution	[RFC4795]
		return addr.GetSegment(0).Matches(224) && addr.GetSegment(1).IsZero() && addr.GetSegment(2).IsZero() && addr.GetSegment(3).Matches(252)
	}
	return addr.GetSegment(0).Matches(169) && addr.GetSegment(1).Matches(254)
}

// IsPrivate returns whether this is a unicast addresses allocated for private use,
// as defined by RFC 1918.
func (addr *IPv4Address) IsPrivate() bool {
	// refer to RFC 1918
	// 10/8 prefix
	// 172.16/12 prefix (172.16.0.0 – 172.31.255.255)
	// 192.168/16 prefix
	seg0, seg1 := addr.GetSegment(0), addr.GetSegment(1)
	return seg0.Matches(10) ||
		(seg0.Matches(172) && seg1.MatchesWithPrefixMask(16, 4)) ||
		(seg0.Matches(192) && seg1.Matches(168))
}

// IsMulticast returns whether this address or subnet is entirely multicast.
func (addr *IPv4Address) IsMulticast() bool {
	// 1110...
	//224.0.0.0/4
	return addr.GetSegment(0).MatchesWithPrefixMask(0xe0, 4)
}

// IsLocal returns true if the address is link local, site local, organization local, administered locally, or unspecified.
// This includes both unicast and multicast.
func (addr *IPv4Address) IsLocal() bool {
	if addr.IsMulticast() {
		//1110...
		seg0 := addr.GetSegment(0)
		//http://www.tcpipguide.com/free/t_IPMulticastAddressing.htm
		//RFC 4607 and https://www.iana.org/assignments/multicast-addresses/multicast-addresses.xhtml

		//239.0.0.0-239.255.255.255 organization local
		if seg0.matches(239) {
			return true
		}
		seg1, seg2 := addr.GetSegment(1), addr.GetSegment(2)

		// 224.0.0.0 to 224.0.0.255 local
		// includes link local multicast name resolution https://tools.ietf.org/html/rfc4795 224.0.0.252
		return (seg0.matches(224) && seg1.IsZero() && seg2.IsZero()) ||
			//232.0.0.1 - 232.0.0.255	Reserved for IANA allocation	[RFC4607]
			//232.0.1.0 - 232.255.255.255	Reserved for local host allocation	[RFC4607]
			(seg0.matches(232) && !(seg1.IsZero() && seg2.IsZero()))
	}
	return addr.IsLinkLocal() || addr.IsPrivate() || addr.IsAnyLocal()
}

// IsUnspecified returns whether this is the unspecified address.  The unspecified address is the address that is all zeros.
func (addr *IPv4Address) IsUnspecified() bool {
	return addr.section == nil || addr.IsZero()
}

// IsAnyLocal returns whether this address is the address which binds to any address on the local host.
// This is the address that has the value of 0, aka the unspecified address.
func (addr *IPv4Address) IsAnyLocal() bool {
	return addr.section == nil || addr.IsZero()
}

// IsLoopback returns whether this address is a loopback address, such as "127.0.0.1".
func (addr *IPv4Address) IsLoopback() bool {
	return addr.section != nil && addr.GetSegment(0).Matches(127)
}

// Iterator provides an iterator to iterate through the individual addresses of this address or subnet.
//
// When iterating, the prefix length is preserved.  Remove it using WithoutPrefixLen prior to iterating if you wish to drop it from all individual addresses.
//
// Call IsMultiple to determine if this instance represents multiple addresses, or GetCount for the count.
func (addr *IPv4Address) Iterator() Iterator[*IPv4Address] {
	if addr == nil {
		return ipv4AddressIterator{nilAddrIterator()}
	}
	return ipv4AddressIterator{addr.init().addrIterator(nil)}
}

// PrefixIterator provides an iterator to iterate through the individual prefixes of this subnet,
// each iterated element spanning the range of values for its prefix.
//
// It is similar to the prefix block iterator, except for possibly the first and last iterated elements, which might not be prefix blocks,
// instead constraining themselves to values from this subnet.
//
// If the subnet has no prefix length, then this is equivalent to Iterator.
func (addr *IPv4Address) PrefixIterator() Iterator[*IPv4Address] {
	return ipv4AddressIterator{addr.init().prefixIterator(false)}
}

// PrefixBlockIterator provides an iterator to iterate through the individual prefix blocks, one for each prefix of this address or subnet.
// Each iterated address or subnet will be a prefix block with the same prefix length as this address or subnet.
//
// If this address has no prefix length, then this is equivalent to Iterator.
func (addr *IPv4Address) PrefixBlockIterator() Iterator[*IPv4Address] {
	return ipv4AddressIterator{addr.init().prefixIterator(true)}
}

// BlockIterator iterates through the addresses that can be obtained by iterating through all the upper segments up to the given segment count.
// The segments following remain the same in all iterated addresses.
//
// For instance, given the IPv4 subnet "1-2.3-4.5-6.7" and the count argument 2,
// BlockIterator will iterate through "1.3.5-6.7", "1.4.5-6.7", "2.3.5-6.7" and "2.4.5-6.7".
func (addr *IPv4Address) BlockIterator(segmentCount int) Iterator[*IPv4Address] {
	return ipv4AddressIterator{addr.init().blockIterator(segmentCount)}
}

// SequentialBlockIterator iterates through the sequential subnets or addresses that make up this address or subnet.
//
// Practically, this means finding the count of segments for which the segments that follow are not full range, and then using BlockIterator with that segment count.
//
// For instance, given the IPv4 subnet "1-2.3-4.5-6.7-8", it will iterate through "1.3.5.7-8", "1.3.6.7-8", "1.4.5.7-8", "1.4.6.7-8", "2.3.5.7-8", "2.3.6.7-8", "2.4.6.7-8" and "2.4.6.7-8".
//
// Use GetSequentialBlockCount to get the number of iterated elements.
func (addr *IPv4Address) SequentialBlockIterator() Iterator[*IPv4Address] {
	return ipv4AddressIterator{addr.init().sequentialBlockIterator()}
}

// GetSequentialBlockIndex gets the minimal segment index for which all following segments are full-range blocks.
//
// The segment at this index is not a full-range block itself, unless all segments are full-range.
// The segment at this index and all following segments form a sequential range.
// For the full subnet to be sequential, the preceding segments must be single-valued.
func (addr *IPv4Address) GetSequentialBlockIndex() int {
	return addr.init().getSequentialBlockIndex()
}

// GetSequentialBlockCount provides the count of elements from the sequential block iterator, the minimal number of sequential subnets that comprise this subnet.
func (addr *IPv4Address) GetSequentialBlockCount() *big.Int {
	return addr.getSequentialBlockCount()
}

func (addr *IPv4Address) rangeIterator(
	upper *IPv4Address,
	valsAreMultiple bool,
	prefixLen PrefixLen,
	segProducer func(addr *IPAddress, index int) *IPAddressSegment,
	segmentIteratorProducer func(seg *IPAddressSegment, index int) Iterator[*IPAddressSegment],
	segValueComparator func(seg1, seg2 *IPAddress, index int) bool,
	networkSegmentIndex,
	hostSegmentIndex int,
	prefixedSegIteratorProducer func(seg *IPAddressSegment, index int) Iterator[*IPAddressSegment],
) Iterator[*IPv4Address] {
	return ipv4AddressIterator{addr.ipAddressInternal.rangeIterator(upper.ToIP(), valsAreMultiple, prefixLen, segProducer, segmentIteratorProducer, segValueComparator, networkSegmentIndex, hostSegmentIndex, prefixedSegIteratorProducer)}
}

// IncrementBoundary returns the address that is the given increment from the range boundaries of this subnet.
//
// If the given increment is positive, adds the value to the upper address (GetUpper) in the subnet range to produce a new address.
// If the given increment is negative, adds the value to the lower address (GetLower) in the subnet range to produce a new address.
// If the increment is zero, returns this address.
//
// If this is a single address value, that address is simply incremented by the given increment value, positive or negative.
//
// On address overflow or underflow, IncrementBoundary returns nil.
func (addr *IPv4Address) IncrementBoundary(increment int64) *IPv4Address {
	return addr.init().incrementBoundary(increment).ToIPv4()
}

// Increment returns the address from the subnet that is the given increment upwards into the subnet range,
// with the increment of 0 returning the first address in the range.
//
// If the increment i matches or exceeds the subnet size count c, then i - c + 1
// is added to the upper address of the range.
// An increment matching the subnet count gives you the address just above the highest address in the subnet.
//
// If the increment is negative, it is added to the lower address of the range.
// To get the address just below the lowest address of the subnet, use the increment -1.
//
// If this is just a single address value, the address is simply incremented by the given increment, positive or negative.
//
// If this is a subnet with multiple values, a positive increment i is equivalent i + 1 values from the subnet iterator and beyond.
// For instance, a increment of 0 is the first value from the iterator, an increment of 1 is the second value from the iterator, and so on.
// An increment of a negative value added to the subnet count is equivalent to the same number of iterator values preceding the upper bound of the iterator.
// For instance, an increment of count - 1 is the last value from the iterator, an increment of count - 2 is the second last value, and so on.
//
// On address overflow or underflow, Increment returns nil.
func (addr *IPv4Address) Increment(increment int64) *IPv4Address {
	return addr.init().increment(increment).ToIPv4()
}

// Enumerate indicates where an address sits relative to the subnet ordering.
//
// Determines how many address elements of the subnet precede the given address element, if the address is in the subnet.
// If above the subnet range, it is the distance to the upper boundary added to the subnet count less one, and if below the subnet range, the distance to the lower boundary.
// If not in the subnet, but neither above nor below the range, then 0 is returned and ok is false.
//
// In other words, if the given address is not in the subnet but above it, returns the number of addresses preceding the address from the upper range boundary,
// added to one less than the total number of subnet addresses.  If the given address is not in the subnet but below it, returns the number of addresses following the address to the lower subnet boundary.
//
// The ok value is false when the argument is multi-valued. The argument must be an individual address.
//
// When this is also an individual address, the returned value is the distance (difference) between the two addresses.
//
// If the given address does not have the same version or type, then ok is false.
func (addr *IPv4Address) EnumerateIPv4(other AddressType) (val int64, ok bool) {
	if other != nil {
		if otherAddr := other.ToAddressBase(); otherAddr != nil {
			return addr.GetSection().enumerateAddrIPv4(otherAddr.GetSection())
		}
	}
	return
}

// Enumerate indicates where an address sits relative to the subnet ordering.
//
// Determines how many address elements of the subnet precede the given address element, if the address is in the subnet.
// If above the subnet range, it is the distance to the upper boundary added to the subnet count less one, and if below the subnet range, the distance to the lower boundary.
//
// In other words, if the given address is not in the subnet but above it, returns the number of addresses preceding the address from the upper range boundary,
// added to one less than the total number of subnet addresses.  If the given address is not in the subnet but below it, returns the number of addresses following the address to the lower subnet boundary.
//
// If the argument is not in the subnet, but neither above nor below the range, then nil is returned.
//
// Enumerate returns nil when the argument is multi-valued. The argument must be an individual address.
//
// When this is also an individual address, the returned value is the distance (difference) between the two addresses.
//
// Enumerate is the inverse of the increment method:
//   - subnet.Enumerate(subnet.Increment(inc)) = inc
//   - subnet.Increment(subnet.Enumerate(newAddr)) = newAddr
//
// If the given address does not have the same version or type, then nil is returned.
func (addr *IPv4Address) Enumerate(other AddressType) *big.Int {
	if other != nil {
		if otherAddr := other.ToAddressBase(); otherAddr != nil {
			return addr.GetSection().enumerateAddr(otherAddr.GetSection())
		}
	}
	return nil
}

// SpanWithPrefixBlocks returns an array of prefix blocks that cover the same set of addresses as this subnet.
//
// Unlike SpanWithPrefixBlocksTo, the result only includes addresses that are a part of this subnet.
func (addr *IPv4Address) SpanWithPrefixBlocks() []*IPv4Address {
	if addr.IsSequential() {
		if addr.IsSinglePrefixBlock() {
			return []*IPv4Address{addr}
		}
		return getSpanningPrefixBlocks(addr, addr)
	}
	return spanWithPrefixBlocks(addr)
}

// SpanWithPrefixBlocksTo returns the smallest slice of prefix block subnets that span from this subnet to the given subnet.
//
// The resulting slice is sorted from lowest address value to highest, regardless of the size of each prefix block.
//
// From the list of returned subnets you can recover the original range (this to other) by converting each to IPAddressRange with ToSequentialRange
// and them joining them into a single range with the Join method of IPAddressSeqRange.
func (addr *IPv4Address) SpanWithPrefixBlocksTo(other *IPv4Address) []*IPv4Address {
	return getSpanningPrefixBlocks(addr, other)
}

// SpanWithSequentialBlocks produces the smallest slice of sequential blocks that cover the same set of addresses as this subnet.
//
// This slice can be shorter than that produced by SpanWithPrefixBlocks and is never longer.
//
// Unlike SpanWithSequentialBlocksTo, this method only includes addresses that are a part of this subnet.
func (addr *IPv4Address) SpanWithSequentialBlocks() []*IPv4Address {
	if addr.IsSequential() {
		return []*IPv4Address{addr}
	}
	return spanWithSequentialBlocks(addr)
}

// SpanWithSequentialBlocksTo produces the smallest slice of sequential block subnets that span all values from this subnet to the given subnet.
// The span will cover all addresses in both subnets and everything in between.
//
// Individual block subnets come in the form "1-3.1-4.5.6-8", however that particular subnet is not sequential since address "1.1.5.8" is in the subnet,
// the next sequential address "1.1.5.9" is not in the subnet, and a higher address "1.2.5.6" is in the subnet.
// Blocks are sequential when the first segment with a range of values is followed by segments that span all values.
//
// The resulting slice is sorted from lowest address value to highest, regardless of the size of each prefix block.
func (addr *IPv4Address) SpanWithSequentialBlocksTo(other *IPv4Address) []*IPv4Address {
	return getSpanningSequentialBlocks(addr, other)
}

// CoverWithPrefixBlockTo returns the minimal-size prefix block that covers all the addresses spanning from this subnet to the given subnet.
func (addr *IPv4Address) CoverWithPrefixBlockTo(other *IPv4Address) *IPv4Address {
	return addr.init().coverWithPrefixBlockTo(other.ToIP()).ToIPv4()
}

// CoverWithPrefixBlock returns the minimal-size prefix block that covers all the addresses in this subnet.
// The resulting block will have a larger subnet size than this, unless this subnet is already a prefix block.
func (addr *IPv4Address) CoverWithPrefixBlock() *IPv4Address {
	return addr.init().coverWithPrefixBlock().ToIPv4()
}

// MergeToSequentialBlocks merges this with the list of addresses to produce the smallest array of sequential blocks.
//
// The resulting slice is sorted from lowest address value to highest, regardless of the size of each prefix block.
func (addr *IPv4Address) MergeToSequentialBlocks(addrs ...*IPv4Address) []*IPv4Address {
	return getMergedSequentialBlocks(cloneSeries(addr, addrs))
}

// MergeToPrefixBlocks merges this subnet with the list of subnets to produce the smallest array of CIDR prefix blocks.
//
// The resulting slice is sorted from lowest address value to highest, regardless of the size of each prefix block.
func (addr *IPv4Address) MergeToPrefixBlocks(addrs ...*IPv4Address) []*IPv4Address {
	return getMergedPrefixBlocks(cloneSeries(addr, addrs))
}

// ReverseBytes returns a new address with the bytes reversed.  Any prefix length is dropped.
func (addr *IPv4Address) ReverseBytes() *IPv4Address {
	addr = addr.init()
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
func (addr *IPv4Address) ReverseBits(perByte bool) (*IPv4Address, addrerr.IncompatibleAddressError) {
	addr = addr.init()
	res, err := addr.GetSection().ReverseBits(perByte)
	if err != nil {
		return nil, err
	}
	return addr.checkIdentity(res), nil
}

// ReverseSegments returns a new address with the segments reversed.
func (addr *IPv4Address) ReverseSegments() *IPv4Address {
	addr = addr.init()
	return addr.checkIdentity(addr.GetSection().ReverseSegments())
}

// ReplaceLen replaces segments starting from startIndex and ending before endIndex with the same number of segments starting at replacementStartIndex from the replacement section.
// Mappings to or from indices outside the range of this or the replacement address are skipped.
func (addr *IPv4Address) ReplaceLen(startIndex, endIndex int, replacement *IPv4Address, replacementIndex int) *IPv4Address {
	if replacementIndex <= 0 {
		startIndex -= replacementIndex
		replacementIndex = 0
	} else if replacementIndex >= IPv4SegmentCount {
		return addr
	}
	// We must do a 1 to 1 adjustment of indices before calling the section replace which would do an adjustment of indices not 1 to 1.
	// Here we assume replacementIndex is 0 and working on the subsection starting at that index.
	// In other words, a replacementIndex of x on the whole section is equivalent to replacementIndex of 0 on the shorter subsection starting at x.
	// Then afterwards we use the original replacement index to work on the whole section again, adjusting as needed.
	startIndex, endIndex, replacementIndexAdjustment := adjust1To1Indices(startIndex, endIndex, IPv4SegmentCount, IPv4SegmentCount-replacementIndex)
	if startIndex == endIndex {
		return addr
	}
	replacementIndex += replacementIndexAdjustment
	count := endIndex - startIndex
	return addr.init().checkIdentity(addr.GetSection().ReplaceLen(startIndex, endIndex, replacement.GetSection(), replacementIndex, replacementIndex+count))
}

// Replace replaces segments starting from startIndex with segments from the replacement section.
// Mappings to or from indices outside the range of this address or the replacement section are skipped.
func (addr *IPv4Address) Replace(startIndex int, replacement *IPv4AddressSection) *IPv4Address {
	// We must do a 1 to 1 adjustment of indices before calling the section replace which would do an adjustment of indices not 1 to 1.
	startIndex, endIndex, replacementIndex :=
		adjust1To1Indices(startIndex, startIndex+replacement.GetSegmentCount(), IPv4SegmentCount, replacement.GetSegmentCount())
	count := endIndex - startIndex
	return addr.init().checkIdentity(addr.GetSection().ReplaceLen(startIndex, endIndex, replacement, replacementIndex, replacementIndex+count))
}

// GetLeadingBitCount returns the number of consecutive leading one or zero bits.
// If ones is true, returns the number of consecutive leading one bits.
// Otherwise, returns the number of consecutive leading zero bits.
//
// This method applies to the lower value of the range if this is a subnet representing multiple values.
func (addr *IPv4Address) GetLeadingBitCount(ones bool) BitCount {
	return addr.init().getLeadingBitCount(ones)
}

// GetTrailingBitCount returns the number of consecutive trailing one or zero bits.
// If ones is true, returns the number of consecutive trailing zero bits.
// Otherwise, returns the number of consecutive trailing one bits.
//
// This method applies to the lower value of the range if this is a subnet representing multiple values.
func (addr *IPv4Address) GetTrailingBitCount(ones bool) BitCount {
	return addr.init().getTrailingBitCount(ones)
}

// GetNetwork returns the singleton IPv4 network instance.
func (addr *IPv4Address) GetNetwork() IPAddressNetwork {
	return ipv4Network
}

// GetIPv6Address creates an IPv6 mixed address using the given ipv6 segments and using this address for the embedded IPv4 segments
func (addr *IPv4Address) GetIPv6Address(section *IPv6AddressSection) (*IPv6Address, addrerr.AddressError) {
	if section.GetSegmentCount() < IPv6MixedOriginalSegmentCount {
		return nil, &addressValueError{addressError: addressError{key: "ipaddress.mac.error.not.eui.convertible"}}
	}
	newSegs := createSegmentArray(IPv6SegmentCount)
	section = section.WithoutPrefixLen()
	section.copyDivisions(newSegs)
	sect, err := createMixedSection(newSegs, addr)
	if err != nil {
		return nil, err
	}
	return newIPv6Address(sect), nil
}

// GetIPv4MappedAddress returns the IPv4-mapped IPv6 address corresponding to this IPv4 address.
// The IPv4-mapped IPv6 address is all zeros in the first 12 bytes, with the last 4 bytes matching the bytes of this IPv4 address.
// See rfc 5156 for details.
// If this is a subnet with segment ranges which cannot be converted to two IPv6 segment ranges, than an error is returned.
func (addr *IPv4Address) GetIPv4MappedAddress() (*IPv6Address, addrerr.IncompatibleAddressError) {
	zero := zeroIPv6Seg.ToDiv()
	segs := createSegmentArray(IPv6SegmentCount)
	segs[0], segs[1], segs[2], segs[3], segs[4] = zero, zero, zero, zero, zero
	segs[5] = NewIPv6Segment(IPv6MaxValuePerSegment).ToDiv()
	var sect *IPv6AddressSection
	sect, err := createMixedSection(segs, addr.WithoutPrefixLen())
	if err != nil {
		return nil, err
	}
	return newIPv6Address(sect), nil
}

// returns an error if the first or 3rd segments have a range of values that cannot be combined with their neighbouting segments into IPv6 segments
func (addr *IPv4Address) getIPv6Address(ipv6Segs []*AddressDivision) (*IPv6Address, addrerr.IncompatibleAddressError) {
	newSegs := createSegmentArray(IPv6SegmentCount)
	copy(newSegs, ipv6Segs)
	sect, err := createMixedSection(newSegs, addr)
	if err != nil {
		return nil, err
	}
	return newIPv6Address(sect), nil
}

func createMixedSection(newIPv6Divisions []*AddressDivision, mixedSection *IPv4Address) (res *IPv6AddressSection, err addrerr.IncompatibleAddressError) {
	ipv4Section := mixedSection.GetSection().WithoutPrefixLen()
	var seg *IPv6AddressSegment
	if seg, err = ipv4Section.GetSegment(0).Join(ipv4Section.GetSegment(1)); err == nil {
		newIPv6Divisions[6] = seg.ToDiv()
		if seg, err = ipv4Section.GetSegment(2).Join(ipv4Section.GetSegment(3)); err == nil {
			newIPv6Divisions[7] = seg.ToDiv()
			res = newIPv6SectionFromMixed(newIPv6Divisions)
			if res.cache != nil {
				nonMixedSection := res.createNonMixedSection()
				mixedGrouping := newIPv6v4MixedGrouping(
					nonMixedSection,
					ipv4Section,
				)
				mixed := &mixedCache{
					defaultMixedAddressSection: mixedGrouping,
					embeddedIPv6Section:        nonMixedSection,
					embeddedIPv4Section:        ipv4Section,
				}
				res.cache.mixed = mixed
			}
		}
	}
	return
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
func (addr IPv4Address) Format(state fmt.State, verb rune) {
	addr.init().format(state, verb)
}

// String implements the [fmt.Stringer] interface, returning the canonical string provided by ToCanonicalString, or "<nil>" if the receiver is a nil pointer.
func (addr *IPv4Address) String() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toString()
}

// GetSegmentStrings returns a slice with the string for each segment being the string that is normalized with wildcards.
func (addr *IPv4Address) GetSegmentStrings() []string {
	if addr == nil {
		return nil
	}
	return addr.init().getSegmentStrings()
}

// ToCanonicalString produces a canonical string for the address.
//
// For IPv4, dotted octet format, also known as dotted decimal format, is used.
// https://datatracker.ietf.org/doc/html/draft-main-ipaddr-text-rep-00#section-2.1
//
// Each address has a unique canonical string, not counting the prefix length.
// With IP addresses, the prefix length can cause two equal addresses to have different strings, for example "1.2.3.4/16" and "1.2.3.4".
// It can also cause two different addresses to have the same string, such as "1.2.0.0/16" for the individual address "1.2.0.0" and also the prefix block "1.2.*.*".
// Use ToCanonicalWildcardString for a unique string for each IP address and subnet.
func (addr *IPv4Address) ToCanonicalString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toCanonicalString()
}

// ToNormalizedString produces a normalized string for the address.
//
// For IPv4, it is the same as the canonical string.
//
// Each address has a unique normalized string, not counting the prefix length.
// With IP addresses, the prefix length can cause two equal addresses to have different strings, for example "1.2.3.4/16" and "1.2.3.4".
// It can also cause two different addresses to have the same string, such as "1.2.0.0/16" for the individual address "1.2.0.0" and also the prefix block "1.2.*.*".
// Use the method ToNormalizedWildcardString for a unique string for each IP address and subnet.
func (addr *IPv4Address) ToNormalizedString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toNormalizedString()
}

// ToCompressedString produces a short representation of this address while remaining within the confines of standard representation(s) of the address.
//
// For IPv4, it is the same as the canonical string.
func (addr *IPv4Address) ToCompressedString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toCompressedString()
}

// ToCanonicalWildcardString produces a string similar to the canonical string and avoids the CIDR prefix length.
// Addresses and subnets with a network prefix length will be shown with wildcards and ranges (denoted by '*' and '-') instead of using the CIDR prefix length notation.
// For IPv4 it is the same as ToNormalizedWildcardString.
func (addr *IPv4Address) ToCanonicalWildcardString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toCanonicalWildcardString()
}

// ToNormalizedWildcardString produces a string similar to the normalized string but avoids the CIDR prefix length.
// CIDR addresses will be shown with wildcards and ranges (denoted by '*' and '-') instead of using the CIDR prefix notation.
func (addr *IPv4Address) ToNormalizedWildcardString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toNormalizedWildcardString()
}

// ToSegmentedBinaryString writes this address as segments of binary values preceded by the "0b" prefix.
func (addr *IPv4Address) ToSegmentedBinaryString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toSegmentedBinaryString()
}

// ToSQLWildcardString create a string similar to that from toNormalizedWildcardString except that
// it uses SQL wildcards.  It uses '%' instead of '*' and also uses the wildcard '_'.
func (addr *IPv4Address) ToSQLWildcardString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toSQLWildcardString()
}

// ToFullString produces a string with no compressed segments and all segments of full length with leading zeros,
// which is 3 characters for IPv4 segments.
func (addr *IPv4Address) ToFullString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toFullString()
}

// ToReverseDNSString generates the reverse-DNS lookup string.
// For IPV4, the error is always nil.
// For "8.255.4.4" it is "4.4.255.8.in-addr.arpa".
func (addr *IPv4Address) ToReverseDNSString() (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	str, _ := addr.init().toReverseDNSString()
	return str, nil
}

// ToPrefixLenString returns a string with a CIDR network prefix length if this address has a network prefix length.
// For IPv6, a zero host section will be compressed with "::". For IPv4 the string is equivalent to the canonical string.
func (addr *IPv4Address) ToPrefixLenString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toPrefixLenString()
}

// ToSubnetString produces a string with specific formats for subnets.
// The subnet string looks like "1.2.*.*" or "1:2::/16".
//
// In the case of IPv4, this means that wildcards are used instead of a network prefix when a network prefix has been supplied.
func (addr *IPv4Address) ToSubnetString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toSubnetString()
}

// ToCompressedWildcardString produces a string similar to ToNormalizedWildcardString, and in fact
// for IPv4 it is the same as ToNormalizedWildcardString.
func (addr *IPv4Address) ToCompressedWildcardString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toCompressedWildcardString()
}

// ToHexString writes this address as a single hexadecimal value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0x" prefix.
//
// If a subnet cannot be written as a single prefix block or a range of two values, an error is returned.
func (addr *IPv4Address) ToHexString(with0xPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	return addr.init().toHexString(with0xPrefix)
}

// ToOctalString writes this address as a single octal value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0" prefix.
//
// If a subnet cannot be written as a single prefix block or a range of two values, an error is returned.
func (addr *IPv4Address) ToOctalString(with0Prefix bool) (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	return addr.init().toOctalString(with0Prefix)
}

// ToBinaryString writes this address as a single binary value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0b" prefix.
//
// If a subnet cannot be written as a single prefix block or a range of two values, an error is returned.
func (addr *IPv4Address) ToBinaryString(with0bPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	return addr.init().toBinaryString(with0bPrefix)
}

// ToUNCHostName Generates the Microsoft UNC path component for this address.
//
// For IPv4 it is the canonical string.
func (addr *IPv4Address) ToUNCHostName() string {
	return addr.ToCanonicalString()
}

// ToInetAtonString returns a string with a format that is styled from the inet_aton routine.
// The string can have an octal or hexadecimal radix rather than decimal.
// When using octal, the octal segments each have a leading zero prefix of "0", and when using hex, a prefix of "0x".
func (addr *IPv4Address) ToInetAtonString(radix Inet_aton_radix) string {
	if addr == nil {
		return nilString()
	}
	return addr.GetSection().ToInetAtonString(radix)
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
func (addr *IPv4Address) ToInetAtonJoinedString(radix Inet_aton_radix, joinedCount int) (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	return addr.GetSection().ToInetAtonJoinedString(radix, joinedCount)
}

// ToCustomString creates a customized string from this address or subnet according to the given string option parameters.
func (addr *IPv4Address) ToCustomString(stringOptions addrstr.IPStringOptions) string {
	if addr == nil {
		return nilString()
	}
	return addr.GetSection().toCustomString(stringOptions)
}

// ToAddressBase converts to an Address, a polymorphic type usable with all addresses and subnets.
// Afterwards, you can convert back with ToIPv4.
//
// ToAddressBase can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (addr *IPv4Address) ToAddressBase() *Address {
	if addr != nil {
		addr = addr.init()
	}
	return (*Address)(unsafe.Pointer(addr))
}

// ToIP converts to an IPAddress, a polymorphic type usable with all IP addresses and subnets.
// Afterwards, you can convert back with ToIPv4.
//
// ToIP can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (addr *IPv4Address) ToIP() *IPAddress {
	if addr != nil {
		addr = addr.init()
	}
	return (*IPAddress)(addr)
}

// Wrap wraps this IP address, returning a WrappedIPAddress, an implementation of ExtendedIPSegmentSeries,
// which can be used to write code that works with both IP addresses and IP address sections.
// Wrap can be called with a nil receiver, wrapping a nil address.
func (addr *IPv4Address) Wrap() WrappedIPAddress {
	return wrapIPAddress(addr.ToIP())
}

// WrapAddress wraps this IP address, returning a WrappedAddress, an implementation of ExtendedSegmentSeries,
// which can be used to write code that works with both addresses and address sections.
// WrapAddress can be called with a nil receiver, wrapping a nil address.
func (addr *IPv4Address) WrapAddress() WrappedAddress {
	return wrapAddress(addr.ToAddressBase())
}

func (addr *IPv4Address) toMaxLower() *IPv4Address {
	return addr.init().addressInternal.toMaxLower().ToIPv4()
}

func (addr *IPv4Address) toMinUpper() *IPv4Address {
	return addr.init().addressInternal.toMinUpper().ToIPv4()
}

// ToKey creates the associated address key.
// While addresses can be compared with the Compare, TrieCompare or Equal methods as well as various provided instances of AddressComparator,
// they are not comparable with Go operators.
// However, AddressKey instances are comparable with Go operators, and thus can be used as map keys.
func (addr *IPv4Address) ToKey() IPv4AddressKey {
	addr = addr.init()
	key := IPv4AddressKey{}
	section := addr.GetSection()
	divs := section.getDivArray()
	var newVal uint64
	if addr.IsMultiple() {
		for _, div := range divs {
			seg := div.ToIPv4()
			newVal = (newVal << IPv4BitsPerSegment) | uint64(seg.GetIPv4SegmentValue()) | (uint64(seg.GetIPv4UpperSegmentValue()) << IPv4BitCount)
		}
	} else {
		for _, div := range divs {
			seg := div.ToIPv4()
			newVal = (newVal << IPv4BitsPerSegment) | uint64(seg.GetIPv4SegmentValue())
		}
		newVal |= newVal << IPv4BitCount
	}
	key.vals = newVal
	return key
}

func fromIPv4Key(key IPv4AddressKey) *IPv4Address {
	keyVal := key.vals
	return NewIPv4AddressFromRange(
		func(segmentIndex int) IPv4SegInt {
			segIndex := (IPv4SegmentCount - 1) - segmentIndex
			return IPv4SegInt(keyVal >> (segIndex << ipv4BitsToSegmentBitshift))
		},
		func(segmentIndex int) IPv4SegInt {
			segIndex := ((IPv4SegmentCount << 1) - 1) - segmentIndex
			return IPv4SegInt(keyVal >> (segIndex << ipv4BitsToSegmentBitshift))
		},
	)
}

// ToGenericKey produces a generic Key[*IPv4Address] that can be used with generic code working with [Address], [IPAddress], [IPv4Address], [IPv6Address] and [MACAddress].
// ToKey produces a more compact key for code that is IPv4-specific.
func (addr *IPv4Address) ToGenericKey() Key[*IPv4Address] {
	// Note: We intentionally do not populate the "scheme" field, which is used with Key[*Address] and Key[*IPAddress] and 64-bit Key[*MACAddress].
	// With Key[*IPv4Address], by leaving the scheme zero, the zero Key[*IPv4Address] matches up with the key produced here by the zero address.
	// We do not need the scheme field for Key[*IPv4Address] since the generic type indicates IPv4.
	key := Key[*IPv4Address]{}
	addr.init().toIPv4Key(&key.keyContents)
	return key
}

func (addr *IPv4Address) fromKey(_ addressScheme, key *keyContents) *IPv4Address {
	return fromIPv4IPKey(key)
}

func (addr *IPv4Address) toIPv4Key(contents *keyContents) {
	section := addr.GetSection()
	divs := section.getDivArray()
	val := &contents.vals[0]
	if addr.IsMultiple() {
		for _, div := range divs {
			seg := div.ToIPv4()
			val.lower = (val.lower << IPv4BitsPerSegment) | uint64(seg.GetIPv4SegmentValue())
			val.upper = (val.upper << IPv4BitsPerSegment) | uint64(seg.GetIPv4UpperSegmentValue())
		}
	} else {
		for _, div := range divs {
			seg := div.ToIPv4()
			val.lower = (val.lower << IPv4BitsPerSegment) | uint64(seg.GetIPv4SegmentValue())
			val.upper = val.lower
		}
	}
}

func fromIPv4IPKey(key *keyContents) *IPv4Address {
	return NewIPv4AddressFromRange(
		func(segmentIndex int) IPv4SegInt {
			segIndex := (IPv4SegmentCount - 1) - segmentIndex
			return IPv4SegInt(key.vals[0].lower >> (segIndex << ipv4BitsToSegmentBitshift))
		}, func(segmentIndex int) IPv4SegInt {
			segIndex := (IPv4SegmentCount - 1) - segmentIndex
			return IPv4SegInt(key.vals[0].upper >> (segIndex << ipv4BitsToSegmentBitshift))
		},
	)
}
