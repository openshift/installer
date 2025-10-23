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
	IPv6SegmentSeparator            = ':'
	IPv6SegmentSeparatorStr         = ":"
	IPv6ZoneSeparator               = '%'
	IPv6ZoneSeparatorStr            = "%"
	IPv6AlternativeZoneSeparator    = '\u00a7'
	IPv6AlternativeZoneSeparatorStr = "\u00a7" //'§'
	IPv6BitsPerSegment              = 16
	IPv6BytesPerSegment             = 2
	IPv6SegmentCount                = 8
	IPv6MixedReplacedSegmentCount   = 2
	IPv6MixedOriginalSegmentCount   = 6
	IPv6MixedOriginalByteCount      = 12
	IPv6ByteCount                   = 16
	IPv6BitCount                    = 128
	IPv6DefaultTextualRadix         = 16
	IPv6MaxValuePerSegment          = 0xffff
	IPv6ReverseDnsSuffix            = ".ip6.arpa"
	IPv6ReverseDnsSuffixDeprecated  = ".ip6.int"

	IPv6UncSegmentSeparator    = '-'
	IPv6UncSegmentSeparatorStr = "-"
	IPv6UncZoneSeparator       = 's'
	IPv6UncZoneSeparatorStr    = "s"
	IPv6UncRangeSeparator      = AlternativeRangeSeparator
	IPv6UncRangeSeparatorStr   = AlternativeRangeSeparatorStr
	IPv6UncSuffix              = ".ipv6-literal.net"

	IPv6SegmentMaxChars = 4

	ipv6BitsToSegmentBitshift = 4

	IPv6AlternativeRangeSeparatorStr = AlternativeRangeSeparatorStr
)

// Zone represents an IPv6 address zone or scope.
type Zone string

// IsEmpty returns whether the zone is the zero-zone, which is the lack of a zone, or the empty string zone.
func (zone Zone) IsEmpty() bool {
	return zone == ""
}

// String implements the [fmt.Stringer] interface, returning the zone characters as a string
func (zone Zone) String() string {
	return string(zone)
}

const NoZone = ""

func newIPv6Address(section *IPv6AddressSection) *IPv6Address {
	return createAddress(section.ToSectionBase(), NoZone).ToIPv6()
}

// NewIPv6Address constructs an IPv6 address or subnet from the given address section.
// If the section does not have 8 segments, an error is returned.
func NewIPv6Address(section *IPv6AddressSection) (*IPv6Address, addrerr.AddressValueError) {
	if section == nil {
		return zeroIPv6, nil
	}
	segCount := section.GetSegmentCount()
	if segCount != IPv6SegmentCount {
		return nil, &addressValueError{
			addressError: addressError{key: "ipaddress.error.invalid.size"},
			val:          segCount,
		}
	}
	return createAddress(section.ToSectionBase(), NoZone).ToIPv6(), nil
}

func newIPv6AddressZoned(section *IPv6AddressSection, zone string) *IPv6Address {
	zoneVal := Zone(zone)
	result := createAddress(section.ToSectionBase(), zoneVal).ToIPv6()
	assignIPv6Cache(zoneVal, result.cache)
	return result
}

func assignIPv6Cache(zoneVal Zone, cache *addressCache) {
	if zoneVal != NoZone { // will need to cache its own strings
		cache.stringCache = &stringCache{ipv6StringCache: &ipv6StringCache{}, ipStringCache: &ipStringCache{}}
	}
}

// NewIPv6AddressZoned constructs an IPv6 address or subnet from the given address section and zone.
// If the section does not have 8 segments, an error is returned.
func NewIPv6AddressZoned(section *IPv6AddressSection, zone string) (*IPv6Address, addrerr.AddressValueError) {
	if section == nil {
		return zeroIPv6.SetZone(zone), nil
	}
	segCount := section.GetSegmentCount()
	if segCount != IPv6SegmentCount {
		return nil, &addressValueError{
			addressError: addressError{key: "ipaddress.error.invalid.size"},
			val:          segCount,
		}
	}
	return newIPv6AddressZoned(section, zone), nil
}

// NewIPv6AddressFromSegs constructs an IPv6 address or subnet from the given segments.
// If the given slice does not have 8 segments, an error is returned.
func NewIPv6AddressFromSegs(segments []*IPv6AddressSegment) (addr *IPv6Address, err addrerr.AddressValueError) {
	segCount := len(segments)
	if segCount != IPv6SegmentCount {
		return nil, &addressValueError{
			addressError: addressError{key: "ipaddress.error.invalid.size"},
			val:          segCount,
		}
	}
	section := NewIPv6Section(segments)
	return NewIPv6Address(section)
}

// NewIPv6AddressFromPrefixedSegs constructs an IPv6 address or subnet from the given segments and prefix length.
// If the given slice does not have 8 segments, an error is returned.
// If the address has a zero host for its prefix length, the returned address will be the prefix block.
func NewIPv6AddressFromPrefixedSegs(segments []*IPv6AddressSegment, prefixLength PrefixLen) (addr *IPv6Address, err addrerr.AddressValueError) {
	segCount := len(segments)
	if segCount != IPv6SegmentCount {
		return nil, &addressValueError{
			addressError: addressError{key: "ipaddress.error.invalid.size"},
			val:          segCount,
		}
	}
	section := NewIPv6PrefixedSection(segments, prefixLength)
	return NewIPv6Address(section)
}

// NewIPv6AddressFromZonedSegs constructs an IPv6 address or subnet from the given segments and zone.
// If the given slice does not have 8 segments, an error is returned.
func NewIPv6AddressFromZonedSegs(segments []*IPv6AddressSegment, zone string) (addr *IPv6Address, err addrerr.AddressValueError) {
	segCount := len(segments)
	if segCount != IPv6SegmentCount {
		return nil, &addressValueError{
			addressError: addressError{key: "ipaddress.error.invalid.size"},
			val:          segCount,
		}
	}
	section := NewIPv6Section(segments)
	return NewIPv6AddressZoned(section, zone)
}

// NewIPv6AddressFromPrefixedZonedSegs constructs an IPv6 address or subnet from the given segments, prefix length, and zone.
// If the given slice does not have 8 segments, an error is returned.
// If the address has a zero host for its prefix length, the returned address will be the prefix block.
func NewIPv6AddressFromPrefixedZonedSegs(segments []*IPv6AddressSegment, prefixLength PrefixLen, zone string) (addr *IPv6Address, err addrerr.AddressValueError) {
	segCount := len(segments)
	if segCount != IPv6SegmentCount {
		return nil, &addressValueError{
			addressError: addressError{key: "ipaddress.error.invalid.size"},
			val:          segCount,
		}
	}
	section := NewIPv6PrefixedSection(segments, prefixLength)
	return NewIPv6AddressZoned(section, zone)
}

// NewIPv6AddressFromBytes constructs an IPv6 address from the given byte slice.
// An error is returned when the byte slice has too many bytes to match the IPv6 segment count of 8.
// There should be 16 bytes or less, although extra leading zeros are tolerated.
func NewIPv6AddressFromBytes(bytes []byte) (addr *IPv6Address, err addrerr.AddressValueError) {
	section, err := NewIPv6SectionFromSegmentedBytes(bytes, IPv6SegmentCount)
	if err == nil {
		addr = newIPv6Address(section)
	}
	return
}

// NewIPv6AddressFromPrefixedBytes constructs an IPv6 address from the given byte slice and prefix length.
// An error is returned when the byte slice has too many bytes to match the IPv6 segment count of 8.
// There should be 16 bytes or less, although extra leading zeros are tolerated.
// If the address has a zero host for the given prefix length, the returned address will be the prefix block.
func NewIPv6AddressFromPrefixedBytes(bytes []byte, prefixLength PrefixLen) (addr *IPv6Address, err addrerr.AddressValueError) {
	section, err := NewIPv6SectionFromPrefixedBytes(bytes, IPv6SegmentCount, prefixLength)
	if err == nil {
		addr = newIPv6Address(section)
	}
	return
}

// NewIPv6AddressFromZonedBytes constructs an IPv6 address from the given byte slice and zone.
// An error is returned when the byte slice has too many bytes to match the IPv6 segment count of 8.
// There should be 16 bytes or less, although extra leading zeros are tolerated.
func NewIPv6AddressFromZonedBytes(bytes []byte, zone string) (addr *IPv6Address, err addrerr.AddressValueError) {
	addr, err = NewIPv6AddressFromBytes(bytes)
	if err == nil {
		addr.zone = Zone(zone)
		assignIPv6Cache(addr.zone, addr.cache)
	}
	return
}

// NewIPv6AddressFromPrefixedZonedBytes constructs an IPv6 address from the given byte slice, prefix length, and zone.
// An error is returned when the byte slice has too many bytes to match the IPv6 segment count of 8.
// There should be 16 bytes or less, although extra leading zeros are tolerated.
// If the address has a zero host for the given prefix length, the returned address will be the prefix block.
func NewIPv6AddressFromPrefixedZonedBytes(bytes []byte, prefixLength PrefixLen, zone string) (addr *IPv6Address, err addrerr.AddressValueError) {
	addr, err = NewIPv6AddressFromPrefixedBytes(bytes, prefixLength)
	if err == nil {
		addr.zone = Zone(zone)
		assignIPv6Cache(addr.zone, addr.cache)
	}
	return
}

//   TODO LATER maybe integrate with net.Interface "Name"

// NewIPv6AddressFromInt constructs an IPv6 address from the given value.
// An error is returned when the values is negative or too large.
func NewIPv6AddressFromInt(val *big.Int) (addr *IPv6Address, err addrerr.AddressValueError) {
	section, err := NewIPv6SectionFromBigInt(val, IPv6SegmentCount)
	if err == nil {
		addr = newIPv6Address(section)
	}
	return
}

// NewIPv6AddressFromPrefixedInt constructs an IPv6 address from the given value and prefix length.
// An error is returned when the values is negative or too large.
// If the address has a zero host for the given prefix length, the returned address will be the prefix block.
func NewIPv6AddressFromPrefixedInt(val *big.Int, prefixLength PrefixLen) (addr *IPv6Address, err addrerr.AddressValueError) {
	section, err := NewIPv6SectionFromPrefixedBigInt(val, IPv6SegmentCount, prefixLength)
	if err == nil {
		addr = newIPv6Address(section)
	}
	return
}

// NewIPv6AddressFromZonedInt constructs an IPv6 address from the given value and zone.
// An error is returned when the values is negative or too large.
func NewIPv6AddressFromZonedInt(val *big.Int, zone string) (addr *IPv6Address, err addrerr.AddressValueError) {
	addr, err = NewIPv6AddressFromInt(val)
	if err == nil {
		addr.zone = Zone(zone)
		assignIPv6Cache(addr.zone, addr.cache)
	}
	return
}

// NewIPv6AddressFromPrefixedZonedInt constructs an IPv6 address from the given value, prefix length, and zone.
// An error is returned when the values is negative or too large.
// If the address has a zero host for the given prefix length, the returned address will be the prefix block.
func NewIPv6AddressFromPrefixedZonedInt(val *big.Int, prefixLength PrefixLen, zone string) (addr *IPv6Address, err addrerr.AddressValueError) {
	addr, err = NewIPv6AddressFromPrefixedInt(val, prefixLength)
	if err == nil {
		addr.zone = Zone(zone)
		assignIPv6Cache(addr.zone, addr.cache)
	}
	return
}

// NewIPv6AddressFromUint64 constructs an IPv6 address from the given values.
func NewIPv6AddressFromUint64(highBytes, lowBytes uint64) *IPv6Address {
	section := NewIPv6SectionFromUint64(highBytes, lowBytes, IPv6SegmentCount)
	return newIPv6Address(section)
}

// NewIPv6AddressFromPrefixedUint64 constructs an IPv6 address or prefix block from the given values and prefix length.
// If the address has a zero host for the given prefix length, the returned address will be the prefix block.
func NewIPv6AddressFromPrefixedUint64(highBytes, lowBytes uint64, prefixLength PrefixLen) *IPv6Address {
	section := NewIPv6SectionFromPrefixedUint64(highBytes, lowBytes, IPv6SegmentCount, prefixLength)
	return newIPv6Address(section)
}

// NewIPv6AddressFromZonedUint64 constructs an IPv6 address from the given values and zone.
func NewIPv6AddressFromZonedUint64(highBytes, lowBytes uint64, zone string) *IPv6Address {
	section := NewIPv6SectionFromUint64(highBytes, lowBytes, IPv6SegmentCount)
	return newIPv6AddressZoned(section, zone)
}

// NewIPv6AddressFromPrefixedZonedUint64 constructs an IPv6 address or prefix block from the given values, prefix length, and zone
// If the address has a zero host for the given prefix length, the returned address will be the prefix block.
func NewIPv6AddressFromPrefixedZonedUint64(highBytes, lowBytes uint64, prefixLength PrefixLen, zone string) *IPv6Address {
	section := NewIPv6SectionFromPrefixedUint64(highBytes, lowBytes, IPv6SegmentCount, prefixLength)
	return newIPv6AddressZoned(section, zone)
}

// NewIPv6AddressFromVals constructs an IPv6 address from the given values.
func NewIPv6AddressFromVals(vals IPv6SegmentValueProvider) *IPv6Address {
	section := NewIPv6SectionFromVals(vals, IPv6SegmentCount)
	return newIPv6Address(section)
}

// NewIPv6AddressFromPrefixedVals constructs an IPv6 address or prefix block from the given values and prefix length.
// If the address has a zero host for the given prefix length, the returned address will be the prefix block.
func NewIPv6AddressFromPrefixedVals(vals IPv6SegmentValueProvider, prefixLength PrefixLen) *IPv6Address {
	section := NewIPv6SectionFromPrefixedVals(vals, IPv6SegmentCount, prefixLength)
	return newIPv6Address(section)
}

// NewIPv6AddressFromRange constructs an IPv6 subnet from the given values.
func NewIPv6AddressFromRange(vals, upperVals IPv6SegmentValueProvider) *IPv6Address {
	section := NewIPv6SectionFromRange(vals, upperVals, IPv6SegmentCount)
	return newIPv6Address(section)
}

// NewIPv6AddressFromPrefixedRange constructs an IPv6 subnet from the given values and prefix length.
// If the address has a zero host for the given prefix length, the returned address will be the prefix block.
func NewIPv6AddressFromPrefixedRange(vals, upperVals IPv6SegmentValueProvider, prefixLength PrefixLen) *IPv6Address {
	section := NewIPv6SectionFromPrefixedRange(vals, upperVals, IPv6SegmentCount, prefixLength)
	return newIPv6Address(section)
}

// NewIPv6AddressFromZonedRange constructs an IPv6 subnet from the given values and zone.
func NewIPv6AddressFromZonedRange(vals, upperVals IPv6SegmentValueProvider, zone string) *IPv6Address {
	section := NewIPv6SectionFromRange(vals, upperVals, IPv6SegmentCount)
	return newIPv6AddressZoned(section, zone)
}

// NewIPv6AddressFromPrefixedZonedRange constructs an IPv6 subnet from the given values, prefix length, and zone.
// If the address has a zero host for the given prefix length, the returned address will be the prefix block.
func NewIPv6AddressFromPrefixedZonedRange(vals, upperVals IPv6SegmentValueProvider, prefixLength PrefixLen, zone string) *IPv6Address {
	section := NewIPv6SectionFromPrefixedRange(vals, upperVals, IPv6SegmentCount, prefixLength)
	return newIPv6AddressZoned(section, zone)
}

func newIPv6AddressFromPrefixedSingle(vals, upperVals IPv6SegmentValueProvider, prefixLength PrefixLen, zone string) *IPv6Address {
	section := newIPv6SectionFromPrefixedSingle(vals, upperVals, IPv6SegmentCount, prefixLength, true)
	return newIPv6AddressZoned(section, zone)
}

// NewIPv6AddressFromMAC constructs an IPv6 address from a modified EUI-64 (Extended Unique Identifier) MAC address and an IPv6 address 64-bit prefix.
//
// If the supplied MAC address section is an 8-byte EUI-64, then it must match the required EUI-64 format of "xx-xx-ff-fe-xx-xx"
// with the "ff-fe" section in the middle.
//
// If the supplied MAC address section is a 6-byte MAC-48 or EUI-48, then the "ff-fe" pattern will be inserted when converting to IPv6.
//
// The constructor will toggle the MAC U/L (universal/local) bit as required with EUI-64.
//
// The IPv6 address section must be at least 8 bytes.  If it has a zone, then the resulting address will have the same zone.
//
// Any prefix length in the MAC address is ignored, while a prefix length in the IPv6 address is preserved but only up to the first 4 segments.
//
// The error is either an AddressValueError for sections that are of insufficient segment count,
// or IncompatibleAddressError when attempting to join two MAC segments, at least one with ranged values, into an equivalent IPV6 segment range.
func NewIPv6AddressFromMAC(prefix *IPv6Address, suffix *MACAddress) (*IPv6Address, addrerr.IncompatibleAddressError) {
	zone := prefix.GetZone()
	zoneStr := NoZone
	if len(zone) > 0 {
		zoneStr = string(zone)
	}
	prefixSection := prefix.GetSection()
	return newIPv6AddressFromMAC(prefixSection, suffix.GetSection(), zoneStr)
}

// when this is called, we know the sections are sufficient length
func newIPv6AddressFromMAC(prefixSection *IPv6AddressSection, suffix *MACAddressSection, zone string) (*IPv6Address, addrerr.IncompatibleAddressError) {
	prefixLen := prefixSection.getPrefixLen()
	if prefixLen != nil && prefixLen.bitCount() > getNetworkPrefixLen(IPv6BitsPerSegment, 0, 4).bitCount() {
		prefixLen = nil
	}
	segments := createSegmentArray(8)
	if err := toIPv6SegmentsFromEUI(segments, 4, suffix, prefixLen); err != nil {
		return nil, err
	}
	prefixSection.copySubDivisions(0, 4, segments)
	res := createIPv6Section(segments)
	res.prefixLength = prefixLen
	res.isMult = suffix.isMultiple() || prefixSection.isMultipleTo(4)
	return newIPv6AddressZoned(res, zone), nil
}

// NewIPv6AddressFromMACSection constructs an IPv6 address from a modified EUI-64 (Extended Unique Identifier) MAC address section and an IPv6 address section network prefix.
//
// If the supplied MAC address section is an 8-byte EUI-64, then it must match the required EUI-64 format of "xx-xx-ff-fe-xx-xx"
// with the "ff-fe" section in the middle.
//
// If the supplied MAC address section is a 6-byte MAC-48 or EUI-48, then the "ff-fe" pattern will be inserted when converting to IPv6.
//
// The constructor will toggle the MAC U/L (universal/local) bit as required with EUI-64.
//
// The IPv6 address section must be at least 8 bytes (4 segments) in length.
//
// Any prefix length in the MAC address is ignored, while a prefix length in the IPv6 address is preserved but only up to the first 4 segments.
//
// The error is either an AddressValueError for sections that are of insufficient segment count,
// or IncompatibleAddressError when unable to join two MAC segments, at least one with ranged values, into an equivalent IPV6 segment range.
func NewIPv6AddressFromMACSection(prefix *IPv6AddressSection, suffix *MACAddressSection) (*IPv6Address, addrerr.AddressError) {
	return newIPv6AddressFromZonedMAC(prefix, suffix, NoZone)
}

// NewIPv6AddressFromZonedMACSection constructs an IPv6 address from a modified EUI-64 (Extended Unique Identifier) MAC address section, an IPv6 address section network prefix, and a zone.
//
// It is similar to NewIPv6AddressFromMACSection but also allows you to specify a zone.
//
// It is similar to NewIPv6AddressFromMAC, which can supply a zone with the IPv6Address argument.
func NewIPv6AddressFromZonedMACSection(prefix *IPv6AddressSection, suffix *MACAddressSection, zone string) (*IPv6Address, addrerr.AddressError) {
	return newIPv6AddressFromZonedMAC(prefix, suffix, zone)
}

func newIPv6AddressFromZonedMAC(prefix *IPv6AddressSection, suffix *MACAddressSection, zone string) (*IPv6Address, addrerr.AddressError) {
	suffixSegCount := suffix.GetSegmentCount()
	if prefix.GetSegmentCount() < 4 || (suffixSegCount != ExtendedUniqueIdentifier48SegmentCount && suffixSegCount != ExtendedUniqueIdentifier64SegmentCount) {
		return nil, &addressValueError{addressError: addressError{key: "ipaddress.mac.error.not.eui.convertible"}}
	}
	return newIPv6AddressFromMAC(prefix, suffix, zone)
}

var zeroIPv6 = initZeroIPv6()
var ipv6All = zeroIPv6.ToPrefixBlockLen(0)

func initZeroIPv6() *IPv6Address {
	div := zeroIPv6Seg
	segs := []*IPv6AddressSegment{div, div, div, div, div, div, div, div}
	section := NewIPv6Section(segs)
	return newIPv6Address(section)
}

// IPv6Address is an IPv6 address, or a subnet of multiple IPv6 addresses.
// An IPv6 address is composed of 8 2-byte segments and can optionally have an associated prefix length.
// Each segment can represent a single value or a range of values.
// The zero value is "::".
//
// To construct one from a string, use NewIPAddressString, then use the ToAddress or GetAddress method of [IPAddressString],
// and then use ToIPv6 to get an IPv6Address, assuming the string had an IPv6 format.
//
// For other inputs, use one of the multiple constructor functions like NewIPv6Address.
// You can also use one of the multiple constructors for [IPAddress] like NewIPAddress and then convert using ToIPv6.
type IPv6Address struct {
	ipAddressInternal
}

func (addr *IPv6Address) init() *IPv6Address {
	if addr.section == nil {
		return zeroIPv6
	}
	return addr
}

// GetCount returns the count of addresses that this address or subnet represents.
//
// If just a single address, not a subnet of multiple addresses, returns 1.
//
// For instance, the IP address subnet "2001:db8::/64" has the count of 2 to the power of 64.
//
// Use IsMultiple if you simply want to know if the count is greater than 1.
func (addr *IPv6Address) GetCount() *big.Int {
	if addr == nil {
		return bigZero()
	}
	return addr.getCount()
}

// IsMultiple returns true if this represents more than a single individual address, whether it is a subnet of multiple addresses.
func (addr *IPv6Address) IsMultiple() bool {
	return addr != nil && addr.isMultiple()
}

// IsPrefixed returns whether this address has an associated prefix length.
func (addr *IPv6Address) IsPrefixed() bool {
	return addr != nil && addr.isPrefixed()
}

// IsFullRange returns whether this address covers the entire IPv6 address space.
//
// This is true if and only if both IncludesZero and IncludesMax return true.
func (addr *IPv6Address) IsFullRange() bool {
	return addr.GetSection().IsFullRange()
}

// GetBitCount returns the number of bits comprising this address,
// or each address in the range if a subnet, which is 128.
func (addr *IPv6Address) GetBitCount() BitCount {
	return IPv6BitCount
}

// GetByteCount returns the number of bytes required for this address,
// or each address in the range if a subnet, which is 16.
func (addr *IPv6Address) GetByteCount() int {
	return IPv6ByteCount
}

// GetBitsPerSegment returns the number of bits comprising each segment in this address.  Segments in the same address are equal length.
func (addr *IPv6Address) GetBitsPerSegment() BitCount {
	return IPv6BitsPerSegment
}

// GetBytesPerSegment returns the number of bytes comprising each segment in this address or subnet.  Segments in the same address are equal length.
func (addr *IPv6Address) GetBytesPerSegment() int {
	return IPv6BytesPerSegment
}

// HasZone returns whether this IPv6 address includes a zone or scope.
func (addr *IPv6Address) HasZone() bool {
	return addr != nil && addr.zone != NoZone
}

// GetZone returns the zone it it has one, otherwise it returns NoZone, which is an empty string.
func (addr *IPv6Address) GetZone() Zone {
	if addr == nil {
		return NoZone
	}
	return addr.zone
}

// GetSection returns the backing section for this address or subnet, comprising all segments.
func (addr *IPv6Address) GetSection() *IPv6AddressSection {
	return addr.init().section.ToIPv6()
}

// GetTrailingSection gets the subsection from the series starting from the given index.
// The first segment is at index 0.
func (addr *IPv6Address) GetTrailingSection(index int) *IPv6AddressSection {
	return addr.GetSection().GetTrailingSection(index)
}

// GetSubSection gets the subsection from the series starting from the given index and ending just before the give endIndex.
// The first segment is at index 0.
func (addr *IPv6Address) GetSubSection(index, endIndex int) *IPv6AddressSection {
	return addr.GetSection().GetSubSection(index, endIndex)
}

// GetNetworkSection returns an address section containing the segments with the network of the address or subnet, the prefix bits.
// The returned section will have only as many segments as needed as determined by the existing CIDR network prefix length.
//
// If this series has no CIDR prefix length, the returned network section will
// be the entire series as a prefixed section with prefix length matching the address bit length.
func (addr *IPv6Address) GetNetworkSection() *IPv6AddressSection {
	return addr.GetSection().GetNetworkSection()
}

// GetNetworkSectionLen returns a section containing the segments with the network of the address or subnet, the prefix bits according to the given prefix length.
// The returned section will have only as many segments as needed to contain the network.
//
// The new section will be assigned the given prefix length,
// unless the existing prefix length is smaller, in which case the existing prefix length will be retained.
func (addr *IPv6Address) GetNetworkSectionLen(prefLen BitCount) *IPv6AddressSection {
	return addr.GetSection().GetNetworkSectionLen(prefLen)
}

// GetHostSection returns a section containing the segments with the host of the address or subnet, the bits beyond the CIDR network prefix length.
// The returned section will have only as many segments as needed to contain the host.
//
// If this series has no prefix length, the returned host section will be the full section.
func (addr *IPv6Address) GetHostSection() *IPv6AddressSection {
	return addr.GetSection().GetHostSection()
}

// GetHostSectionLen returns a section containing the segments with the host of the address or subnet, the bits beyond the given CIDR network prefix length.
// The returned section will have only as many segments as needed to contain the host.
func (addr *IPv6Address) GetHostSectionLen(prefLen BitCount) *IPv6AddressSection {
	return addr.GetSection().GetHostSectionLen(prefLen)
}

// GetNetworkMask returns the network mask associated with the CIDR network prefix length of this address or subnet.
// If this address or subnet has no prefix length, then the all-ones mask is returned.
func (addr *IPv6Address) GetNetworkMask() *IPv6Address {
	var prefLen BitCount
	if pref := addr.getPrefixLen(); pref != nil {
		prefLen = pref.bitCount()
	} else {
		prefLen = IPv6BitCount
	}
	return ipv6Network.GetNetworkMask(prefLen).ToIPv6()
}

// GetHostMask returns the host mask associated with the CIDR network prefix length of this address or subnet.
// If this address or subnet has no prefix length, then the all-ones mask is returned.
func (addr *IPv6Address) GetHostMask() *IPv6Address {
	return addr.getHostMask(ipv6Network).ToIPv6()
}

// GetMixedAddressGrouping creates a grouping by combining an IPv6 address section comprising the first six segments (most significant) in this address
// with the IPv4 section corresponding to the lowest (least-significant) two segments in this address, as produced by GetEmbeddedIPv4Address.
func (addr *IPv6Address) GetMixedAddressGrouping() (*IPv6v4MixedAddressGrouping, addrerr.IncompatibleAddressError) {
	return addr.init().GetSection().getMixedAddressGrouping()
}

// GetEmbeddedIPv4AddressSection gets the IPv4 section corresponding to the lowest (least-significant) 2 segments (4 bytes) in this address.
// Many IPv4 to IPv6 mapping schemes (but not all) use these 4 bytes for a mapped IPv4 address.
// An error can result when one of the associated IPv6 segments has a range of values that cannot be split into two ranges.
func (addr *IPv6Address) GetEmbeddedIPv4AddressSection() (*IPv4AddressSection, addrerr.IncompatibleAddressError) {
	return addr.init().GetSection().getEmbeddedIPv4AddressSection()
}

// GetEmbeddedIPv4Address gets the IPv4 address corresponding to the lowest (least-significant) 2 segments (4 bytes) in this address.
// Many IPv4 to IPv6 mapping schemes (but not all) use these 4 bytes for a mapped IPv4 address.
// An error can result when one of the associated IPv6 segments has a range of values that cannot be split into two ranges.
func (addr *IPv6Address) GetEmbeddedIPv4Address() (*IPv4Address, addrerr.IncompatibleAddressError) {
	section, err := addr.GetEmbeddedIPv4AddressSection()
	if err != nil {
		return nil, err
	}
	return newIPv4Address(section), nil
}

// GetIPv4AddressSection produces an IPv4 address section corresponding to any sequence of bytes in this IPv6 address section
func (addr *IPv6Address) GetIPv4AddressSection(startIndex, endIndex int) (*IPv4AddressSection, addrerr.IncompatibleAddressError) {
	return addr.init().GetSection().GetIPv4AddressSection(startIndex, endIndex)
}

// Get6To4IPv4Address Returns the second and third segments as an IPv4Address.
func (addr *IPv6Address) Get6To4IPv4Address() (*IPv4Address, addrerr.IncompatibleAddressError) {
	return addr.GetEmbeddedIPv4AddressAt(2)
}

// GetEmbeddedIPv4AddressAt produces an IPv4 address corresponding to any sequence of 4 bytes in this IPv6 address, starting at the given index.
func (addr *IPv6Address) GetEmbeddedIPv4AddressAt(byteIndex int) (*IPv4Address, addrerr.IncompatibleAddressError) {
	if byteIndex == IPv6MixedOriginalSegmentCount*IPv6BytesPerSegment {
		return addr.GetEmbeddedIPv4Address()
	}
	if byteIndex > IPv6ByteCount-IPv4ByteCount {
		return nil, &addressValueError{
			addressError: addressError{key: "ipaddress.error.invalid.size"},
			val:          byteIndex,
		}
	}
	section, err := addr.init().GetSection().GetIPv4AddressSection(byteIndex, byteIndex+IPv4ByteCount)
	if err != nil {
		return nil, err
	}
	return newIPv4Address(section), nil
}

// GetIPv6Address creates an IPv6 mixed address using the given address for the trailing embedded IPv4 segments
func (addr *IPv6Address) GetIPv6Address(embedded IPv4Address) (*IPv6Address, addrerr.IncompatibleAddressError) {
	return embedded.getIPv6Address(addr.WithoutPrefixLen().getDivisionsInternal())
}

// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
func (addr *IPv6Address) CopySubSegments(start, end int, segs []*IPv6AddressSegment) (count int) {
	return addr.GetSection().CopySubSegments(start, end, segs)
}

// CopySegments copies the existing segments into the given slice,
// as much as can be fit into the slice, returning the number of segments copied.
func (addr *IPv6Address) CopySegments(segs []*IPv6AddressSegment) (count int) {
	return addr.GetSection().CopySegments(segs)
}

// GetSegments returns a slice with the address segments.  The returned slice is not backed by the same array as this address.
func (addr *IPv6Address) GetSegments() []*IPv6AddressSegment {
	return addr.GetSection().GetSegments()
}

// GetSegment returns the segment at the given index.
// The first segment is at index 0.
// GetSegment will panic given a negative index or an index matching or larger than the segment count.
func (addr *IPv6Address) GetSegment(index int) *IPv6AddressSegment {
	return addr.init().getSegment(index).ToIPv6()
}

// GetSegmentCount returns the segment count, the number of segments in this address, which is 8
func (addr *IPv6Address) GetSegmentCount() int {
	return addr.GetDivisionCount()
}

// ForEachSegment visits each segment in order from most-significant to least, the most significant with index 0, calling the given function for each, terminating early if the function returns true.
// Returns the number of visited segments.
func (addr *IPv6Address) ForEachSegment(consumer func(segmentIndex int, segment *IPv6AddressSegment) (stop bool)) int {
	return addr.GetSection().ForEachSegment(consumer)
}

// GetGenericDivision returns the segment at the given index as a DivisionType.
func (addr *IPv6Address) GetGenericDivision(index int) DivisionType {
	return addr.init().getDivision(index)
}

// GetGenericSegment returns the segment at the given index as an AddressSegmentType.
// The first segment is at index 0.
// GetGenericSegment will panic given a negative index or an index matching or larger than the segment count.
func (addr *IPv6Address) GetGenericSegment(index int) AddressSegmentType {
	return addr.init().getSegment(index)
}

// GetDivisionCount returns the segment count.
func (addr *IPv6Address) GetDivisionCount() int {
	return addr.init().getDivisionCount()
}

// GetIPVersion returns IPv6, the IP version of this address.
func (addr *IPv6Address) GetIPVersion() IPVersion {
	return IPv6
}

func (addr *IPv6Address) checkIdentity(section *IPv6AddressSection) *IPv6Address {
	if section == nil {
		return nil
	}
	sec := section.ToSectionBase()
	if sec == addr.section {
		return addr
	}
	return newIPv6AddressZoned(section, string(addr.zone))
}

// Mask applies the given mask to all addresses represented by this IPv6Address.
// The mask is applied to all individual addresses.
//
// If this represents multiple addresses, and applying the mask to all addresses creates a set of addresses
// that cannot be represented as a sequential range within each segment, then an error is returned.
func (addr *IPv6Address) Mask(other *IPv6Address) (masked *IPv6Address, err addrerr.IncompatibleAddressError) {
	return addr.maskPrefixed(other, true)
}

func (addr *IPv6Address) maskPrefixed(other *IPv6Address, retainPrefix bool) (masked *IPv6Address, err addrerr.IncompatibleAddressError) {
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
func (addr *IPv6Address) BitwiseOr(other *IPv6Address) (masked *IPv6Address, err addrerr.IncompatibleAddressError) {
	return addr.bitwiseOrPrefixed(other, true)
}

func (addr *IPv6Address) bitwiseOrPrefixed(other *IPv6Address, retainPrefix bool) (masked *IPv6Address, err addrerr.IncompatibleAddressError) {
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
func (addr *IPv6Address) Subtract(other *IPv6Address) []*IPv6Address {
	addr = addr.init()
	sects, _ := addr.GetSection().Subtract(other.GetSection())
	sectLen := len(sects)
	if sectLen == 0 {
		return nil
	} else if sectLen == 1 {
		sec := sects[0]
		if sec.ToSectionBase() == addr.section {
			return []*IPv6Address{addr}
		}
	}
	res := make([]*IPv6Address, sectLen)
	for i, sect := range sects {
		res[i] = newIPv6AddressZoned(sect, string(addr.zone))
	}
	return res
}

// Intersect returns the subnet whose addresses are found in both this and the given subnet argument, or nil if no such addresses exist.
//
// This is also known as the conjunction of the two sets of addresses.
func (addr *IPv6Address) Intersect(other *IPv6Address) *IPv6Address {
	addr = addr.init()
	section, _ := addr.GetSection().Intersect(other.GetSection())
	if section == nil {
		return nil
	}
	return addr.checkIdentity(section)
}

// SpanWithRange returns an IPv6AddressSeqRange instance that spans this subnet to the given subnet.
// If the other address is a different version than this, then the other is ignored, and the result is equivalent to calling ToSequentialRange.
func (addr *IPv6Address) SpanWithRange(other *IPv6Address) *SequentialRange[*IPv6Address] {
	return NewSequentialRange(addr.init(), other)
}

// GetLower returns the lowest address in the subnet range,
// which will be the receiver if it represents a single address.
// For example, for "1::1:2-3:4:5-6", the series "1::1:2:4:5" is returned.
func (addr *IPv6Address) GetLower() *IPv6Address {
	return addr.init().getLower().ToIPv6()
}

// GetUpper returns the highest address in the subnet range,
// which will be the receiver if it represents a single address.
// For example, for "1::1:2-3:4:5-6", the series "1::1:3:4:6" is returned.
func (addr *IPv6Address) GetUpper() *IPv6Address {
	return addr.init().getUpper().ToIPv6()
}

// GetLowerIPAddress returns the address in the subnet or address collection with the lowest numeric value,
// which will be the receiver if it represents a single address.
// GetLowerIPAddress implements the IPAddressRange interface
func (addr *IPv6Address) GetLowerIPAddress() *IPAddress {
	return addr.GetLower().ToIP()
}

// GetUpperIPAddress returns the address in the subnet or address collection with the highest numeric value,
// which will be the receiver if it represents a single address.
// GetUpperIPAddress implements the IPAddressRange interface
func (addr *IPv6Address) GetUpperIPAddress() *IPAddress {
	return addr.GetUpper().ToIP()
}

// IsZeroHostLen returns whether the host section is always zero for all individual addresses in this subnet,
// for the given prefix length.
//
// If the host section is zero length (there are zero host bits), IsZeroHostLen returns true.
func (addr *IPv6Address) IsZeroHostLen(prefLen BitCount) bool {
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
func (addr *IPv6Address) ToZeroHost() (*IPv6Address, addrerr.IncompatibleAddressError) {
	res, err := addr.init().toZeroHost(false)
	return res.ToIPv6(), err
}

// ToZeroHostLen converts the address or subnet to one in which all individual addresses have a host of zero,
// the host being the bits following the given prefix length.
// If this address or subnet has the same prefix length, then the returned one will too, otherwise the returned series will have no prefix length.
//
// For instance, the zero host of "1.2.3.4" for the prefix length of 16 is the address "1.2.0.0".
//
// This returns an error if the subnet is a range of addresses which cannot be converted to a range in which all addresses have zero hosts,
// because the conversion results in a subnet segment that is not a sequential range of values.
func (addr *IPv6Address) ToZeroHostLen(prefixLength BitCount) (*IPv6Address, addrerr.IncompatibleAddressError) {
	res, err := addr.init().toZeroHostLen(prefixLength)
	return res.ToIPv6(), err
}

// ToZeroNetwork converts the address or subnet to one in which all individual addresses have a network of zero,
// the network being the bits within the prefix length.
// If the address or subnet has no prefix length, then it returns an all-zero address.
//
// The returned address or subnet will have the same prefix length.
func (addr *IPv6Address) ToZeroNetwork() *IPv6Address {
	return addr.init().toZeroNetwork().ToIPv6()
}

// IsMaxHostLen returns whether the host is all one-bits, the max value, for all individual addresses in this subnet,
// for the given prefix length, the host being the bits following the prefix.
//
// If the host section is zero length (there are zero host bits), IsMaxHostLen returns true.
func (addr *IPv6Address) IsMaxHostLen(prefLen BitCount) bool {
	return addr.init().isMaxHostLen(prefLen)
}

// ToMaxHost converts the address or subnet to one in which all individual addresses have a host of all one-bits, the max value,
// the host being the bits following the prefix length.
// If the address or subnet has no prefix length, then it returns an all-ones address, the max address.
//
// The returned address or subnet will have the same prefix and prefix length.
//
// This returns an error if the subnet is a range of addresses which cannot be converted to a range in which all addresses have max hosts,
// because the conversion results in a subnet segment that is not a sequential range of values.
func (addr *IPv6Address) ToMaxHost() (*IPv6Address, addrerr.IncompatibleAddressError) {
	res, err := addr.init().toMaxHost()
	return res.ToIPv6(), err
}

// ToMaxHostLen converts the address or subnet to one in which all individual addresses have a host of all one-bits, the max host,
// the host being the bits following the given prefix length.
// If this address or subnet has the same prefix length, then the resulting one will too, otherwise the resulting address or subnet will have no prefix length.
//
// For instance, the zero host of "1.2.3.4" for the prefix length of 16 is the address "1.2.255.255".
//
// This returns an error if the subnet is a range of addresses which cannot be converted to a range in which all addresses have max hosts,
// because the conversion results in a subnet segment that is not a sequential range of values.
func (addr *IPv6Address) ToMaxHostLen(prefixLength BitCount) (*IPv6Address, addrerr.IncompatibleAddressError) {
	res, err := addr.init().toMaxHostLen(prefixLength)
	return res.ToIPv6(), err
}

// ToPrefixBlock returns the subnet associated with the prefix length of this address.
// If this address has no prefix length, this address is returned.
//
// The subnet will include all addresses with the same prefix as this one, the prefix "block".
// The network prefix will match the prefix of this address or subnet, and the host values will span all values.
//
// For example, if the address is "1:2:3:4:5:6:7:8/64" it returns the subnet "1:2:3:4::/64" which can also be written as "1:2:3:4:*:*:*:*/64".
func (addr *IPv6Address) ToPrefixBlock() *IPv6Address {
	return addr.init().toPrefixBlock().ToIPv6()
}

// ToPrefixBlockLen returns the subnet associated with the given prefix length.
//
// The subnet will include all addresses with the same prefix as this one, the prefix "block" for that prefix length.
// The network prefix will match the prefix of this address or subnet, and the host values will span all values.
//
// For example, if the address is "1:2:3:4:5:6:7:8" and the prefix length provided is 64, it returns the subnet "1:2:3:4::/64" which can also be written as "1:2:3:4:*:*:*:*/64".
func (addr *IPv6Address) ToPrefixBlockLen(prefLen BitCount) *IPv6Address {
	return addr.init().toPrefixBlockLen(prefLen).ToIPv6()
}

// ToBlock creates a new block of addresses by changing the segment at the given index to have the given lower and upper value,
// and changing the following segments to be full-range.
func (addr *IPv6Address) ToBlock(segmentIndex int, lower, upper SegInt) *IPv6Address {
	return addr.init().toBlock(segmentIndex, lower, upper).ToIPv6()
}

// WithoutPrefixLen provides the same address but with no prefix length.  The values remain unchanged.
func (addr *IPv6Address) WithoutPrefixLen() *IPv6Address {
	if !addr.IsPrefixed() {
		return addr
	}
	return addr.init().withoutPrefixLen().ToIPv6()
}

// SetPrefixLen sets the prefix length.
//
// A prefix length will not be set to a value lower than zero or beyond the bit length of the address.
// The provided prefix length will be adjusted to these boundaries if necessary.
func (addr *IPv6Address) SetPrefixLen(prefixLen BitCount) *IPv6Address {
	return addr.init().setPrefixLen(prefixLen).ToIPv6()
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

func (addr *IPv6Address) SetPrefixLenZeroed(prefixLen BitCount) (*IPv6Address, addrerr.IncompatibleAddressError) {
	res, err := addr.init().setPrefixLenZeroed(prefixLen)
	return res.ToIPv6(), err
}

// AdjustPrefixLen increases or decreases the prefix length by the given increment.
//
// A prefix length will not be adjusted lower than zero or beyond the bit length of the address.
//
// If this address has no prefix length, then the prefix length will be set to the adjustment if positive,
// or it will be set to the adjustment added to the bit count if negative.
func (addr *IPv6Address) AdjustPrefixLen(prefixLen BitCount) *IPv6Address {
	return addr.init().adjustPrefixLen(prefixLen).ToIPv6()
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
func (addr *IPv6Address) AdjustPrefixLenZeroed(prefixLen BitCount) (*IPv6Address, addrerr.IncompatibleAddressError) {
	res, err := addr.init().adjustPrefixLenZeroed(prefixLen)
	return res.ToIPv6(), err
}

// AssignPrefixForSingleBlock returns the equivalent prefix block that matches exactly the range of values in this address.
// The returned block will have an assigned prefix length indicating the prefix length for the block.
//
// There may be no such address - it is required that the range of values match the range of a prefix block.
// If there is no such address, then nil is returned.
func (addr *IPv6Address) AssignPrefixForSingleBlock() *IPv6Address {
	return addr.init().assignPrefixForSingleBlock().ToIPv6()
}

// AssignMinPrefixForBlock returns an equivalent subnet, assigned the smallest prefix length possible,
// such that the prefix block for that prefix length is in this subnet.
//
// In other words, this method assigns a prefix length to this subnet matching the largest prefix block in this subnet.
func (addr *IPv6Address) AssignMinPrefixForBlock() *IPv6Address {
	return addr.init().assignMinPrefixForBlock().ToIPv6()
}

// ToSinglePrefixBlockOrAddress converts to a single prefix block or address.
// If the given address is a single prefix block, it is returned.
// If it can be converted to a single prefix block by assigning a prefix length, the converted block is returned.
// If it is a single address, any prefix length is removed and the address is returned.
// Otherwise, nil is returned.
// This method provides the address formats used by tries.
// ToSinglePrefixBlockOrAddress is quite similar to AssignPrefixForSingleBlock, which always returns prefixed addresses, while this does not.
func (addr *IPv6Address) ToSinglePrefixBlockOrAddress() *IPv6Address {
	return addr.init().toSinglePrefixBlockOrAddr().ToIPv6()
}

func (addr *IPv6Address) toSinglePrefixBlockOrAddress() (*IPv6Address, addrerr.IncompatibleAddressError) {
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
func (addr *IPv6Address) ContainsPrefixBlock(prefixLen BitCount) bool {
	return addr.init().ipAddressInternal.ContainsPrefixBlock(prefixLen)
}

// ContainsSinglePrefixBlock returns whether this address contains a single prefix block for the given prefix length.
//
// This means there is only one prefix value for the given prefix length, and it also contains the full prefix block for that prefix, all addresses with that prefix.
//
// Use GetPrefixLenForSingleBlock to determine whether there is a prefix length for which this method returns true.
func (addr *IPv6Address) ContainsSinglePrefixBlock(prefixLen BitCount) bool {
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
func (addr *IPv6Address) GetMinPrefixLenForBlock() BitCount {
	return addr.init().ipAddressInternal.GetMinPrefixLenForBlock()
}

// GetPrefixLenForSingleBlock returns a prefix length for which the range of this address subnet matches exactly the block of addresses for that prefix.
//
// If the range can be described this way, then this method returns the same value as GetMinPrefixLenForBlock.
//
// If no such prefix exists, returns nil.
//
// If this segment grouping represents a single value, returns the bit length of this address division series.
func (addr *IPv6Address) GetPrefixLenForSingleBlock() PrefixLen {
	return addr.init().ipAddressInternal.GetPrefixLenForSingleBlock()
}

// GetValue returns the lowest address in this subnet or address as an integer value.
func (addr *IPv6Address) GetValue() *big.Int {
	return addr.init().section.GetValue()
}

// GetUpperValue returns the highest address in this subnet or address as an integer value.
func (addr *IPv6Address) GetUpperValue() *big.Int {
	return addr.init().section.GetUpperValue()
}

// Uint64Values returns the lowest address in the address range as a pair of uint64 values.
func (addr *IPv6Address) Uint64Values() (high, low uint64) {
	return addr.GetSection().Uint64Values()
}

// UpperUint64Values returns the highest address in the address section range as a pair of uint64 values.
func (addr *IPv6Address) UpperUint64Values() (high, low uint64) {
	return addr.GetSection().UpperUint64Values()
}

// GetNetIPAddr returns the lowest address in this subnet or address as a net.IPAddr.
func (addr *IPv6Address) GetNetIPAddr() *net.IPAddr {
	return addr.ToIP().GetNetIPAddr()
}

// GetUpperNetIPAddr returns the highest address in this subnet or address as a net.IPAddr.
func (addr *IPv6Address) GetUpperNetIPAddr() *net.IPAddr {
	return addr.ToIP().GetUpperNetIPAddr()
}

// GetNetIP returns the lowest address in this subnet or address as a net.IP.
func (addr *IPv6Address) GetNetIP() net.IP {
	return addr.Bytes()
}

// GetUpperNetIP returns the highest address in this subnet or address as a net.IP.
func (addr *IPv6Address) GetUpperNetIP() net.IP {
	return addr.UpperBytes()
}

// GetNetNetIPAddr returns the lowest address in this subnet or address range as a netip.Addr.
func (addr *IPv6Address) GetNetNetIPAddr() netip.Addr {
	res := addr.init().getNetNetIPAddr()
	if addr.hasZone() {
		res = res.WithZone(string(addr.zone))
	}
	return res
}

// GetUpperNetNetIPAddr returns the highest address in this subnet or address range as a netip.Addr.
func (addr *IPv6Address) GetUpperNetNetIPAddr() netip.Addr {
	return addr.init().getUpperNetNetIPAddr()
}

// CopyNetIP copies the value of the lowest individual address in the subnet into a net.IP.
//
// If the value can fit in the given net.IP slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (addr *IPv6Address) CopyNetIP(bytes net.IP) net.IP {
	return addr.CopyBytes(bytes)
}

// CopyUpperNetIP copies the value of the highest individual address in the subnet into a net.IP.
//
// If the value can fit in the given net.IP slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (addr *IPv6Address) CopyUpperNetIP(bytes net.IP) net.IP {
	return addr.CopyUpperBytes(bytes)
}

// Bytes returns the lowest address in this subnet or address as a byte slice.
func (addr *IPv6Address) Bytes() []byte {
	return addr.init().section.Bytes()
}

// UpperBytes returns the highest address in this subnet or address as a byte slice.
func (addr *IPv6Address) UpperBytes() []byte {
	return addr.init().section.UpperBytes()
}

// CopyBytes copies the value of the lowest individual address in the subnet into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (addr *IPv6Address) CopyBytes(bytes []byte) []byte {
	return addr.init().section.CopyBytes(bytes)
}

// CopyUpperBytes copies the value of the highest individual address in the subnet into a byte slice.
//
// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
// Otherwise, a new slice is created and returned with the value.
func (addr *IPv6Address) CopyUpperBytes(bytes []byte) []byte {
	return addr.init().section.CopyUpperBytes(bytes)
}

// IsMax returns whether this address matches exactly the maximum possible value, the address whose bits are all ones.
func (addr *IPv6Address) IsMax() bool {
	return addr.init().section.IsMax()
}

// IncludesMax returns whether this address includes the max address, the address whose bits are all ones, within its range.
func (addr *IPv6Address) IncludesMax() bool {
	return addr.init().section.IncludesMax()
}

// TestBit returns true if the bit in the lower value of this address at the given index is 1, where index 0 refers to the least significant bit.
// In other words, it computes (bits & (1 << n)) != 0), using the lower value of this address.
// TestBit will panic if n < 0, or if it matches or exceeds the bit count of this item.
func (addr *IPv6Address) TestBit(n BitCount) bool {
	return addr.init().testBit(n)
}

// IsOneBit returns true if the bit in the lower value of this address at the given index is 1, where index 0 refers to the most significant bit.
// IsOneBit will panic if bitIndex is less than zero, or if it is larger than the bit count of this item.
func (addr *IPv6Address) IsOneBit(bitIndex BitCount) bool {
	return addr.init().isOneBit(bitIndex)
}

// PrefixEqual determines if the given address matches this address up to the prefix length of this address.
// It returns whether the two addresses share the same range of prefix values.
func (addr *IPv6Address) PrefixEqual(other AddressType) bool {
	return addr.init().prefixEquals(other)
}

// PrefixContains returns whether the prefix values in the given address or subnet
// are prefix values in this address or subnet, using the prefix length of this address or subnet.
// If this address has no prefix length, the entire address is compared.
//
// It returns whether the prefix of this address contains all values of the same prefix length in the given address.
func (addr *IPv6Address) PrefixContains(other AddressType) bool {
	return addr.init().prefixContains(other)
}

// containsSame returns whether this address contains all addresses in the given address or subnet of the same type.
func (addr *IPv6Address) containsSame(other *IPv6Address) bool {
	return addr.Contains(other)
}

// Contains returns whether this is the same type and version as the given address or subnet and whether it contains all addresses in the given address or subnet.
func (addr *IPv6Address) Contains(other AddressType) bool {
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
	return otherAddr.getAddrType() == ipv6Type && addr.section.sameCountTypeContains(otherAddr.GetSection()) &&
		addr.isSameZone(other.ToAddressBase())
}

// ContainsRange returns true if this address contains the given sequential range
func (addr *IPv6Address) ContainsRange(other IPAddressSeqRangeType) bool {
	return isContainedBy(other, addr.ToIP())
}

// Overlaps returns true if this address overlaps the given address or subnet
func (addr *IPv6Address) Overlaps(other AddressType) bool {
	if addr == nil {
		return true
	}
	return addr.init().overlaps(other)
}

// Overlaps returns true if this address overlaps the given sequential range
func (addr *IPv6Address) OverlapsRange(other IPAddressSeqRangeType) bool {
	if other == nil {
		return true
	}
	return other.OverlapsAddress(addr)
}

// Compare returns a negative integer, zero, or a positive integer if this address or subnet is less than, equal, or greater than the given item.
// Any address item is comparable to any other.  All address items use CountComparator to compare.
func (addr *IPv6Address) Compare(item AddressItem) int {
	return CountComparator.Compare(addr, item)
}

// Equal returns whether the given address or subnet is equal to this address or subnet.
// Two address instances are equal if they represent the same set of addresses.
func (addr *IPv6Address) Equal(other AddressType) bool {
	if addr == nil {
		return other == nil || other.ToAddressBase() == nil
	} else if other.ToAddressBase() == nil {
		return false
	}
	return other.ToAddressBase().getAddrType() == ipv6Type && addr.init().section.sameCountTypeEquals(other.ToAddressBase().GetSection()) &&
		addr.isSameZone(other.ToAddressBase())
}

// CompareSize compares the counts of two subnets or addresses or items, the number of individual addresses or items within.
//
// Rather than calculating counts with GetCount, there can be more efficient ways of determining whether this subnet represents more individual addresses or items than another.
//
// CompareSize returns a positive integer if this address or subnet has a larger count than the one given, zero if they are the same, or a negative integer if the other has a larger count.
func (addr *IPv6Address) CompareSize(other AddressItem) int {
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
func (addr *IPv6Address) TrieCompare(other *IPv6Address) int {
	return addr.init().trieCompare(other.ToAddressBase())
}

// TrieIncrement returns the next address or block according to address trie ordering
//
// If an address is neither an individual address nor a prefix block, it is treated like one:
//
//   - ranges that occur inside the prefix length are ignored, only the lower value is used.
//   - ranges beyond the prefix length are assumed to be the full range across all hosts for that prefix length.
func (addr *IPv6Address) TrieIncrement() *IPv6Address {
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
func (addr *IPv6Address) TrieDecrement() *IPv6Address {
	if res, ok := trieDecrement(addr); ok {
		return res
	}
	return nil
}

// MatchesWithMask applies the mask to this address and then compares the result with the given address,
// returning true if they match, false otherwise.
func (addr *IPv6Address) MatchesWithMask(other *IPv6Address, mask *IPv6Address) bool {
	return addr.init().GetSection().MatchesWithMask(other.GetSection(), mask.GetSection())
}

// GetMaxSegmentValue returns the maximum possible segment value for this type of address.
//
// Note this is not the maximum of the range of segment values in this specific address,
// this is the maximum value of any segment for this address type and version, determined by the number of bits per segment.
func (addr *IPv6Address) GetMaxSegmentValue() SegInt {
	return addr.init().getMaxSegmentValue()
}

// WithoutZone returns the same address but with no zone.
func (addr *IPv6Address) WithoutZone() *IPv6Address {
	if addr.HasZone() {
		return newIPv6Address(addr.GetSection())
	}
	return addr
}

// SetZone returns the same address associated with the given zone,  The existing zone, if any, is replaced.
func (addr *IPv6Address) SetZone(zone string) *IPv6Address {
	if Zone(zone) == addr.GetZone() {
		return addr
	}
	return newIPv6AddressZoned(addr.GetSection(), zone)
}

// ToSequentialRange creates a sequential range instance from the lowest and highest addresses in this subnet.
//
// The two will represent the same set of individual addresses if and only if IsSequential is true.
// To get a series of ranges that represent the same set of individual addresses use the SequentialBlockIterator (or PrefixIterator),
// and apply this method to each iterated subnet.
//
// If this represents just a single address then the returned instance covers just that single address as well.
func (addr *IPv6Address) ToSequentialRange() *SequentialRange[*IPv6Address] {
	if addr == nil {
		return nil
	}
	addr = addr.init().WithoutPrefixLen().WithoutZone()
	return newSequRangeUnchecked(
		addr.GetLower(),
		addr.GetUpper(),
		addr.isMultiple())
}

func (addr *IPv6Address) getLowestHighestAddrs() (lower, upper *IPv6Address) {
	l, u := addr.ipAddressInternal.getLowestHighestAddrs()
	return l.ToIPv6(), u.ToIPv6()
}

// ToAddressString retrieves or generates an IPAddressString instance for this IPAddress instance.
// This may be the IPAddressString this instance was generated from, if it was generated from an IPAddressString.
//
// In general, users are intended to create IPAddress instances from IPAddressString instances,
// while the reverse direction is generally not common and not useful, except under specific circumstances.
//
// However, the reverse direction can be useful under certain circumstances,
// such as when maintaining a collection of HostIdentifierString instances.
func (addr *IPv6Address) ToAddressString() *IPAddressString {
	return addr.init().ToIP().ToAddressString()
}

// IncludesZeroHostLen returns whether the subnet contains an individual address with a host of zero, an individual address for which all bits past the given prefix length are zero.
func (addr *IPv6Address) IncludesZeroHostLen(networkPrefixLength BitCount) bool {
	return addr.init().includesZeroHostLen(networkPrefixLength)
}

// IncludesMaxHostLen returns whether the subnet contains an individual address with a host of all one-bits, an individual address for which all bits past the given prefix length are all ones.
func (addr *IPv6Address) IncludesMaxHostLen(networkPrefixLength BitCount) bool {
	return addr.init().includesMaxHostLen(networkPrefixLength)
}

// IsLinkLocal returns whether the address is link local, whether unicast or multicast.
func (addr *IPv6Address) IsLinkLocal() bool {
	firstSeg := addr.GetSegment(0)
	return (addr.IsMulticast() && firstSeg.matchesWithMask(2, 0xf)) || // ffx2::/16
		//1111 1110 10 .... fe8x currently only in use
		firstSeg.MatchesWithPrefixMask(0xfe80, 10)
}

// IsLocal returns true if the address is link local, site local, organization local, administered locally, or unspecified.
// This includes both unicast and multicast.
func (addr *IPv6Address) IsLocal() bool {
	if addr.IsMulticast() {
		/*
				 [RFC4291][RFC7346]
				 11111111|flgs|scop
					scope 4 bits
					 1  Interface-Local scope
			         2  Link-Local scope
			         3  Realm-Local scope
			         4  Admin-Local scope
			         5  Site-Local scope
			         8  Organization-Local scope
			         E  Global scope
		*/
		firstSeg := addr.GetSegment(0)
		if firstSeg.matchesWithMask(8, 0xf) {
			return true
		}
		if firstSeg.GetValueCount() <= 5 &&
			(firstSeg.getSegmentValue()&0xf) >= 1 && (firstSeg.getUpperSegmentValue()&0xf) <= 5 {
			//all values fall within the range from interface local to site local
			return true
		}
		//source specific multicast
		//rfc4607 and https://www.iana.org/assignments/multicast-addresses/multicast-addresses.xhtml
		//FF3X::8000:0 - FF3X::FFFF:FFFF	Reserved for local host allocation	[RFC4607]
		if firstSeg.MatchesWithPrefixMask(0xff30, 12) && addr.GetSegment(6).MatchesWithPrefixMask(0x8000, 1) {
			return true
		}
	}
	return addr.IsLinkLocal() || addr.IsSiteLocal() || addr.IsUniqueLocal() || addr.IsAnyLocal()
}

// IsUnspecified returns whether this is the unspecified address.  The unspecified address is the address that is all zeros.
func (addr *IPv6Address) IsUnspecified() bool {
	return addr.section == nil || addr.IsZero()
}

// IsAnyLocal returns whether this address is the address which binds to any address on the local host.
// This is the address that has the value of 0, aka the unspecified address.
func (addr *IPv6Address) IsAnyLocal() bool {
	return addr.section == nil || addr.IsZero()
}

// IsSiteLocal returns true if the address is site-local, or all addresses in the subnet are site-local, see rfc 3513, 3879, and 4291.
func (addr *IPv6Address) IsSiteLocal() bool {
	firstSeg := addr.GetSegment(0)
	return (addr.IsMulticast() && firstSeg.matchesWithMask(5, 0xf)) || // ffx5::/16
		//1111 1110 11 ...
		firstSeg.MatchesWithPrefixMask(0xfec0, 10) // deprecated RFC 3879
}

// IsUniqueLocal returns true if the address is unique-local, or all addresses in the subnet are unique-local, see RFC 4193.
func (addr *IPv6Address) IsUniqueLocal() bool {
	//RFC 4193
	return addr.GetSegment(0).MatchesWithPrefixMask(0xfc00, 7)
}

// IsIPv4Mapped returns whether the address or all addresses in the subnet are IPv4-mapped.
//
// "::ffff:x:x/96" indicates an IPv6 address mapped to IPv4.
func (addr *IPv6Address) IsIPv4Mapped() bool {
	//::ffff:x:x/96 indicates IPv6 address mapped to IPv4
	if addr.GetSegment(5).Matches(IPv6MaxValuePerSegment) {
		for i := 0; i < 5; i++ {
			if !addr.GetSegment(i).IsZero() {
				return false
			}
		}
		return true
	}
	return false
}

// IsIPv4Compatible returns whether the address or all addresses in the subnet are IPv4-compatible.
func (addr *IPv6Address) IsIPv4Compatible() bool {
	return addr.GetSegment(0).IsZero() && addr.GetSegment(1).IsZero() && addr.GetSegment(2).IsZero() &&
		addr.GetSegment(3).IsZero() && addr.GetSegment(4).IsZero() && addr.GetSegment(5).IsZero()
}

// Is6To4 returns whether the address or subnet is IPv6 to IPv4 relay.
func (addr *IPv6Address) Is6To4() bool {
	//2002::/16
	return addr.GetSegment(0).Matches(0x2002)
}

// Is6Over4 returns whether the address or all addresses in the subnet are 6over4.
func (addr *IPv6Address) Is6Over4() bool {
	return addr.GetSegment(0).Matches(0xfe80) &&
		addr.GetSegment(1).IsZero() && addr.GetSegment(2).IsZero() &&
		addr.GetSegment(3).IsZero() && addr.GetSegment(4).IsZero() &&
		addr.GetSegment(5).IsZero()
}

// IsTeredo returns whether the address or all addresses in the subnet are Teredo.
func (addr *IPv6Address) IsTeredo() bool {
	//2001::/32
	return addr.GetSegment(0).Matches(0x2001) && addr.GetSegment(1).IsZero()
}

// IsIsatap returns whether the address or all addresses in the subnet are ISATAP.
func (addr *IPv6Address) IsIsatap() bool {
	// 0,1,2,3 is fe80::
	// 4 can be 0200
	return addr.GetSegment(0).Matches(0xfe80) &&
		addr.GetSegment(1).IsZero() &&
		addr.GetSegment(2).IsZero() &&
		addr.GetSegment(3).IsZero() &&
		(addr.GetSegment(4).IsZero() || addr.GetSegment(4).Matches(0x200)) &&
		addr.GetSegment(5).Matches(0x5efe)
}

// IsIPv4Translatable returns whether the address or subnet is IPv4 translatable as in RFC 2765.
func (addr *IPv6Address) IsIPv4Translatable() bool { //rfc 2765
	//::ffff:0:x:x/96 indicates IPv6 addresses translated from IPv4
	return addr.GetSegment(4).Matches(0xffff) &&
		addr.GetSegment(5).IsZero() &&
		addr.GetSegment(0).IsZero() &&
		addr.GetSegment(1).IsZero() &&
		addr.GetSegment(2).IsZero() &&
		addr.GetSegment(3).IsZero()
}

// IsWellKnownIPv4Translatable returns whether the address has the well-known prefix for IPv4-translatable addresses as in RFC 6052 and RFC 6144.
func (addr *IPv6Address) IsWellKnownIPv4Translatable() bool { //rfc 6052 rfc 6144
	//64:ff9b::/96 prefix for auto ipv4/ipv6 translation
	if addr.GetSegment(0).Matches(0x64) && addr.GetSegment(1).Matches(0xff9b) {
		for i := 2; i <= 5; i++ {
			if !addr.GetSegment(i).IsZero() {
				return false
			}
		}
		return true
	}
	return false
}

// IsMulticast returns whether this address or subnet is entirely multicast.
func (addr *IPv6Address) IsMulticast() bool {
	// 11111111...
	return addr.GetSegment(0).MatchesWithPrefixMask(0xff00, 8)
}

// IsLoopback returns whether this address is a loopback address, namely "::1".
func (addr *IPv6Address) IsLoopback() bool {
	if addr.section == nil {
		return false
	}
	//::1
	i := 0
	lim := addr.GetSegmentCount() - 1
	for ; i < lim; i++ {
		if !addr.GetSegment(i).IsZero() {
			return false
		}
	}
	return addr.GetSegment(i).Matches(1)
}

// Iterator provides an iterator to iterate through the individual addresses of this address or subnet.
//
// When iterating, the prefix length is preserved.  Remove it using WithoutPrefixLen prior to iterating if you wish to drop it from all individual addresses.
//
// Call IsMultiple to determine if this instance represents multiple addresses, or GetCount for the count.
func (addr *IPv6Address) Iterator() Iterator[*IPv6Address] {
	if addr == nil {
		return ipv6AddressIterator{nilAddrIterator()}
	}
	return ipv6AddressIterator{addr.init().addrIterator(nil)}
}

// PrefixIterator provides an iterator to iterate through the individual prefixes of this subnet,
// each iterated element spanning the range of values for its prefix.
//
// It is similar to the prefix block iterator, except for possibly the first and last iterated elements, which might not be prefix blocks,
// instead constraining themselves to values from this subnet.
//
// If the subnet has no prefix length, then this is equivalent to Iterator.
func (addr *IPv6Address) PrefixIterator() Iterator[*IPv6Address] {
	return ipv6AddressIterator{addr.init().prefixIterator(false)}
}

// PrefixBlockIterator provides an iterator to iterate through the individual prefix blocks, one for each prefix of this address or subnet.
// Each iterated address or subnet will be a prefix block with the same prefix length as this address or subnet.
//
// If this address has no prefix length, then this is equivalent to Iterator.
func (addr *IPv6Address) PrefixBlockIterator() Iterator[*IPv6Address] {
	return ipv6AddressIterator{addr.init().prefixIterator(true)}
}

// BlockIterator iterates through the addresses that can be obtained by iterating through all the upper segments up to the given segment count.
// The segments following remain the same in all iterated addresses.
func (addr *IPv6Address) BlockIterator(segmentCount int) Iterator[*IPv6Address] {
	return ipv6AddressIterator{addr.init().blockIterator(segmentCount)}
}

// SequentialBlockIterator iterates through the sequential subnets or addresses that make up this address or subnet.
//
// Practically, this means finding the count of segments for which the segments that follow are not full range, and then using BlockIterator with that segment count.
//
// For instance, given the IPv4 subnet "1-2.3-4.5-6.7-8", it will iterate through "1.3.5.7-8", "1.3.6.7-8", "1.4.5.7-8", "1.4.6.7-8", "2.3.5.7-8", "2.3.6.7-8", "2.4.6.7-8" and "2.4.6.7-8".
//
// Use GetSequentialBlockCount to get the number of iterated elements.
func (addr *IPv6Address) SequentialBlockIterator() Iterator[*IPv6Address] {
	return ipv6AddressIterator{addr.init().sequentialBlockIterator()}
}

// GetSequentialBlockIndex gets the minimal segment index for which all following segments are full-range blocks.
//
// The segment at this index is not a full-range block itself, unless all segments are full-range.
// The segment at this index and all following segments form a sequential range.
// For the full subnet to be sequential, the preceding segments must be single-valued.
func (addr *IPv6Address) GetSequentialBlockIndex() int {
	return addr.init().getSequentialBlockIndex()
}

// GetSequentialBlockCount provides the count of elements from the sequential block iterator, the minimal number of sequential subnets that comprise this subnet.
func (addr *IPv6Address) GetSequentialBlockCount() *big.Int {
	return addr.getSequentialBlockCount()
}

func (addr *IPv6Address) rangeIterator(
	upper *IPv6Address,
	valsAreMultiple bool,
	prefixLen PrefixLen,
	segProducer func(addr *IPAddress, index int) *IPAddressSegment,
	segmentIteratorProducer func(seg *IPAddressSegment, index int) Iterator[*IPAddressSegment],
	segValueComparator func(seg1, seg2 *IPAddress, index int) bool,
	networkSegmentIndex,
	hostSegmentIndex int,
	prefixedSegIteratorProducer func(seg *IPAddressSegment, index int) Iterator[*IPAddressSegment],
) Iterator[*IPv6Address] {
	return ipv6AddressIterator{addr.ipAddressInternal.rangeIterator(upper.ToIP(), valsAreMultiple, prefixLen, segProducer, segmentIteratorProducer, segValueComparator, networkSegmentIndex, hostSegmentIndex, prefixedSegIteratorProducer)}
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
func (addr *IPv6Address) IncrementBoundary(increment int64) *IPv6Address {
	return addr.init().incrementBoundary(increment).ToIPv6()
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
func (addr *IPv6Address) Increment(increment int64) *IPv6Address {
	return addr.init().increment(increment).ToIPv6()
}

// IncrementBig increments the address or subnet.  It is the same as Increment but allows for a larger increment value.
// See Increment for more details.
func (addr *IPv6Address) IncrementBig(bigIncrement *big.Int) *IPv6Address {
	return addr.checkIdentity(addr.GetSection().IncrementBig(bigIncrement))
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
func (addr *IPv6Address) Enumerate(other AddressType) *big.Int {
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
func (addr *IPv6Address) SpanWithPrefixBlocks() []*IPv6Address {
	if addr.IsSequential() {
		if addr.IsSinglePrefixBlock() {
			return []*IPv6Address{addr}
		}
		return getSpanningPrefixBlocks(addr, addr)
	}
	return spanWithPrefixBlocks(addr)
}

// SpanWithPrefixBlocksTo returns the smallest slice of prefix block subnets that span from this subnet to the given subnet.
//
// The resulting slice is sorted from lowest address value to highest, regardless of the size of each prefix block.
//
// From the list of returned subnets you can recover the original range (from this to other) by converting each to [SequentialRange] with ToSequentialRange
// and them joining them into a single range with the Join method of [SequentialRange].
func (addr *IPv6Address) SpanWithPrefixBlocksTo(other *IPv6Address) []*IPv6Address {
	return getSpanningPrefixBlocks(addr, other)
}

// SpanWithSequentialBlocks produces the smallest slice of sequential blocks that cover the same set of addresses as this subnet.
//
// This slice can be shorter than that produced by SpanWithPrefixBlocks and is never longer.
//
// Unlike SpanWithSequentialBlocksTo, this method only includes addresses that are a part of this subnet.
func (addr *IPv6Address) SpanWithSequentialBlocks() []*IPv6Address {
	if addr.IsSequential() {
		return []*IPv6Address{addr}
	}
	return spanWithSequentialBlocks(addr)
}

// SpanWithSequentialBlocksTo produces the smallest slice of sequential block subnets that span all values from this subnet to the given subnet.
// The span will cover all addresses in both subnets and everything in between.
//
// The resulting slice is sorted from lowest address value to highest, regardless of the size of each prefix block.
func (addr *IPv6Address) SpanWithSequentialBlocksTo(other *IPv6Address) []*IPv6Address {
	return getSpanningSequentialBlocks(addr, other)
}

// CoverWithPrefixBlockTo returns the minimal-size prefix block that covers all the addresses spanning from this subnet to the given subnet.
func (addr *IPv6Address) CoverWithPrefixBlockTo(other *IPv6Address) *IPv6Address {
	return addr.init().coverWithPrefixBlockTo(other.ToIP()).ToIPv6()
}

// CoverWithPrefixBlock returns the minimal-size prefix block that covers all the addresses in this subnet.
// The resulting block will have a larger subnet size than this, unless this subnet is already a prefix block.
func (addr *IPv6Address) CoverWithPrefixBlock() *IPv6Address {
	return addr.init().coverWithPrefixBlock().ToIPv6()
}

// MergeToSequentialBlocks merges this with the list of addresses to produce the smallest array of sequential blocks.
//
// The resulting slice is sorted from lowest address value to highest, regardless of the size of each prefix block.
func (addr *IPv6Address) MergeToSequentialBlocks(addrs ...*IPv6Address) []*IPv6Address {
	return getMergedSequentialBlocks(cloneSeries(addr, addrs))
}

// MergeToPrefixBlocks merges this subnet with the list of subnets to produce the smallest array of prefix blocks.
//
// The resulting slice is sorted from lowest address value to highest, regardless of the size of each prefix block.
func (addr *IPv6Address) MergeToPrefixBlocks(addrs ...*IPv6Address) []*IPv6Address {
	return getMergedPrefixBlocks(cloneSeries(addr, addrs))
}

// ReverseBytes returns a new address with the bytes reversed.  Any prefix length is dropped.
//
// If the bytes within a single segment cannot be reversed because the segment represents a range,
// and reversing the segment values results in a range that is not contiguous, then this returns an error.
//
// In practice this means that to be reversible, a segment range must include all values except possibly the largest and/or smallest, which reverse to themselves.
func (addr *IPv6Address) ReverseBytes() (*IPv6Address, addrerr.IncompatibleAddressError) {
	res, err := addr.GetSection().ReverseBytes()
	if err != nil {
		return nil, err
	}
	return addr.checkIdentity(res), nil
}

// ReverseBits returns a new address with the bits reversed.  Any prefix length is dropped.
//
// If the bits within a single segment cannot be reversed because the segment represents a range,
// and reversing the segment values results in a range that is not contiguous, this returns an error.
//
// In practice this means that to be reversible, a segment range must include all values except possibly the largest and/or smallest, which reverse to themselves.
//
// If perByte is true, the bits are reversed within each byte, otherwise all the bits are reversed.
func (addr *IPv6Address) ReverseBits(perByte bool) (*IPv6Address, addrerr.IncompatibleAddressError) {
	res, err := addr.GetSection().ReverseBits(perByte)
	if err != nil {
		return nil, err
	}
	return addr.checkIdentity(res), nil
}

// ReverseSegments returns a new address with the segments reversed.
func (addr *IPv6Address) ReverseSegments() *IPv6Address {
	return addr.checkIdentity(addr.GetSection().ReverseSegments())
}

// ReplaceLen replaces segments starting from startIndex and ending before endIndex with the same number of segments starting at replacementStartIndex from the replacement section.
// Mappings to or from indices outside the range of this or the replacement address are skipped.
func (addr *IPv6Address) ReplaceLen(startIndex, endIndex int, replacement *IPv6Address, replacementIndex int) *IPv6Address {
	if replacementIndex <= 0 {
		startIndex -= replacementIndex
		replacementIndex = 0
	} else if replacementIndex >= IPv6SegmentCount {
		return addr
	}
	// We must do a 1 to 1 adjustment of indices before calling the section replace which would do an adjustment of indices not 1 to 1.
	// Here we assume replacementIndex is 0 and working on the subsection starting at that index.
	// In other words, a replacementIndex of x on the whole section is equivalent to replacementIndex of 0 on the shorter subsection starting at x.
	// Then afterwards we use the original replacement index to work on the whole section again, adjusting as needed.
	startIndex, endIndex, replacementIndexAdjustment := adjust1To1Indices(startIndex, endIndex, IPv6SegmentCount, IPv6SegmentCount-replacementIndex)
	if startIndex == endIndex {
		return addr
	}
	replacementIndex += replacementIndexAdjustment
	count := endIndex - startIndex
	return addr.init().checkIdentity(addr.GetSection().ReplaceLen(startIndex, endIndex, replacement.GetSection(), replacementIndex, replacementIndex+count))
}

// Replace replaces segments starting from startIndex with segments from the replacement section.
func (addr *IPv6Address) Replace(startIndex int, replacement *IPv6AddressSection) *IPv6Address {
	// We must do a 1 to 1 adjustment of indices before calling the section replace which would do an adjustment of indices not 1 to 1.
	startIndex, endIndex, replacementIndex :=
		adjust1To1Indices(startIndex, startIndex+replacement.GetSegmentCount(), IPv6SegmentCount, replacement.GetSegmentCount())
	count := endIndex - startIndex
	return addr.init().checkIdentity(addr.GetSection().ReplaceLen(startIndex, endIndex, replacement, replacementIndex, replacementIndex+count))
}

// GetLeadingBitCount returns the number of consecutive leading one or zero bits.
// If ones is true, returns the number of consecutive leading one bits.
// Otherwise, returns the number of consecutive leading zero bits.
//
// This method applies to the lower value of the range if this is a subnet representing multiple values.
func (addr *IPv6Address) GetLeadingBitCount(ones bool) BitCount {
	return addr.init().getLeadingBitCount(ones)
}

// GetTrailingBitCount returns the number of consecutive trailing one or zero bits.
// If ones is true, returns the number of consecutive trailing zero bits.
// Otherwise, returns the number of consecutive trailing one bits.
//
// This method applies to the lower value of the range if this is a subnet representing multiple values.
func (addr *IPv6Address) GetTrailingBitCount(ones bool) BitCount {
	return addr.init().getTrailingBitCount(ones)
}

// GetNetwork returns the singleton IPv6 network instance.
func (addr *IPv6Address) GetNetwork() IPAddressNetwork {
	return ipv6Network
}

// IsEUI64 returns whether this address is consistent with EUI64,
// which means the 12th and 13th bytes of the address match 0xff and 0xfe.
func (addr *IPv6Address) IsEUI64() bool {
	return addr.GetSegment(6).MatchesWithPrefixMask(0xfe00, 8) &&
		addr.GetSegment(5).MatchesWithMask(0xff, 0xff)
}

// ToEUI converts to the associated MACAddress.
// An error is returned if the 0xfffe pattern is missing in segments 5 and 6,
// or if an IPv6 segment's range of values cannot be split into two ranges of values.
func (addr *IPv6Address) ToEUI(extended bool) (*MACAddress, addrerr.IncompatibleAddressError) {
	segs, err := addr.toEUISegments(extended)
	if err != nil {
		return nil, err
	}
	sect := newMACSectionEUI(segs)
	return newMACAddress(sect), nil
}

// prefix length in this section is ignored when converting to MAC.
func (addr *IPv6Address) toEUISegments(extended bool) ([]*AddressDivision, addrerr.IncompatibleAddressError) {
	seg1 := addr.GetSegment(5)
	seg2 := addr.GetSegment(6)
	if !seg1.MatchesWithMask(0xff, 0xff) || !seg2.MatchesWithPrefixMask(0xfe00, 8) {
		return nil, &incompatibleAddressError{addressError{key: "ipaddress.mac.error.not.eui.convertible"}}
	}
	macStartIndex := 0
	var macSegCount int
	if extended {
		macSegCount = ExtendedUniqueIdentifier64SegmentCount
	} else {
		macSegCount = ExtendedUniqueIdentifier48SegmentCount
	}
	newSegs := createSegmentArray(macSegCount)
	seg0 := addr.GetSegment(4)
	if err := seg0.splitIntoMACSegments(newSegs, macStartIndex); err != nil {
		return nil, err
	}
	//toggle the u/l bit
	macSegment0 := newSegs[0].ToMAC()
	lower0 := macSegment0.GetSegmentValue()
	upper0 := macSegment0.GetUpperSegmentValue()
	mask2ndBit := SegInt(0x2)
	if !macSegment0.MatchesWithMask(mask2ndBit&lower0, mask2ndBit) { // ensures that bit remains constant
		return nil, &incompatibleAddressError{addressError{key: "ipaddress.mac.error.not.eui.convertible"}}
	}
	lower0 ^= mask2ndBit //flip the universal/local bit
	upper0 ^= mask2ndBit
	newSegs[0] = NewMACRangeSegment(MACSegInt(lower0), MACSegInt(upper0)).ToDiv()
	macStartIndex += 2
	if err := seg1.splitIntoMACSegments(newSegs, macStartIndex); err != nil { //a ff fe b
		return nil, err
	}
	if extended {
		macStartIndex += 2
		if err := seg2.splitIntoMACSegments(newSegs, macStartIndex); err != nil {
			return nil, err
		}
	} else {
		first := newSegs[macStartIndex]
		if err := seg2.splitIntoMACSegments(newSegs, macStartIndex); err != nil {
			return nil, err
		}
		newSegs[macStartIndex] = first
	}
	macStartIndex += 2
	seg3 := addr.GetSegment(7)
	if err := seg3.splitIntoMACSegments(newSegs, macStartIndex); err != nil {
		return nil, err
	}
	return newSegs, nil
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
func (addr IPv6Address) Format(state fmt.State, verb rune) {
	addr.init().format(state, verb)
}

// String implements the [fmt.Stringer] interface, returning the canonical string provided by ToCanonicalString, or "<nil>" if the receiver is a nil pointer.
func (addr *IPv6Address) String() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().addressInternal.toString()
}

// GetSegmentStrings returns a slice with the string for each segment being the string that is normalized with wildcards.
func (addr *IPv6Address) GetSegmentStrings() []string {
	if addr == nil {
		return nil
	}
	return addr.init().getSegmentStrings()
}

// ToCanonicalString produces a canonical string for the address.
//
// For IPv6, RFC 5952 describes canonical string representation.
// https://en.wikipedia.org/wiki/IPv6_address#Representation
// http://tools.ietf.org/html/rfc5952
//
// Each address has a unique canonical string, not counting the prefix length.
// With IP addresses, the prefix length can cause two equal addresses to have different strings, for example "1.2.3.4/16" and "1.2.3.4".
// It can also cause two different addresses to have the same string, such as "1.2.0.0/16" for the individual address "1.2.0.0" and also the prefix block "1.2.*.*".
// Use ToCanonicalWildcardString for a unique string for each IP address and subnet.
func (addr *IPv6Address) ToCanonicalString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toCanonicalString()
}

// ToNormalizedString produces a normalized string for the address.
//
// For IPv6, it differs from the canonical string.  Zero-segments are not compressed.
//
// Each address has a unique normalized string, not counting the prefix length.
// With IP addresses, the prefix length can cause two equal addresses to have different strings, for example "1.2.3.4/16" and "1.2.3.4".
// It can also cause two different addresses to have the same string, such as "1.2.0.0/16" for the individual address "1.2.0.0" and also the prefix block "1.2.*.*".
// Use the method ToNormalizedWildcardString for a unique string for each IP address and subnet.
func (addr *IPv6Address) ToNormalizedString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toNormalizedString()
}

// ToCompressedString produces a short representation of this address while remaining within the confines of standard representation(s) of the address.
//
// For IPv6, it differs from the canonical string.  It compresses the maximum number of zeros and/or host segments with the IPv6 compression notation '::'.
func (addr *IPv6Address) ToCompressedString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toCompressedString()
}

// ToCanonicalWildcardString produces a string similar to the canonical string and avoids the CIDR prefix length.
// Addresses and subnets with a network prefix length will be shown with wildcards and ranges (denoted by '*' and '-') instead of using the CIDR prefix length notation.
// IPv6 addresses will be compressed according to the canonical representation.
func (addr *IPv6Address) ToCanonicalWildcardString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toCanonicalWildcardString()
}

// ToNormalizedWildcardString produces a string similar to the normalized string but avoids the CIDR prefix length.
// CIDR addresses will be shown with wildcards and ranges (denoted by '*' and '-') instead of using the CIDR prefix notation.
func (addr *IPv6Address) ToNormalizedWildcardString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toNormalizedWildcardString()
}

// ToSegmentedBinaryString writes this address as segments of binary values preceded by the "0b" prefix.
func (addr *IPv6Address) ToSegmentedBinaryString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toSegmentedBinaryString()
}

// ToSQLWildcardString create a string similar to that from toNormalizedWildcardString except that
// it uses SQL wildcards.  It uses '%' instead of '*' and also uses the ending single-digit wildcard '_'.
func (addr *IPv6Address) ToSQLWildcardString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toSQLWildcardString()
}

// ToFullString produces a string with no compressed segments and all segments of full length with leading zeros,
// which is 4 characters for IPv6 segments.
func (addr *IPv6Address) ToFullString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toFullString()
}

// ToPrefixLenString returns a string with a CIDR network prefix length if this address has a network prefix length.
// For IPv6, a zero host section will be compressed with "::". For IPv4 the string is equivalent to the canonical string.
func (addr *IPv6Address) ToPrefixLenString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toPrefixLenString()
}

// ToSubnetString produces a string with specific formats for subnets.
// The subnet string looks like "1.2.*.*" or "1:2::/16".
//
// In the case of IPv6, when a network prefix has been supplied, the prefix will be shown and the host section will be compressed with "::".
func (addr *IPv6Address) ToSubnetString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toSubnetString()
}

// ToCompressedWildcardString produces a string similar to ToNormalizedWildcardString, avoiding the CIDR prefix, but with full IPv6 segment compression as well, including single zero-segments.
func (addr *IPv6Address) ToCompressedWildcardString() string {
	if addr == nil {
		return nilString()
	}
	return addr.init().toCompressedWildcardString()
}

// ToReverseDNSString generates the reverse-DNS lookup string,
// returning an error if this address is a multiple-valued subnet for which the range cannot be represented.
// For "2001:db8::567:89ab" it is "b.a.9.8.7.6.5.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa".
func (addr *IPv6Address) ToReverseDNSString() (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	return addr.init().toReverseDNSString()
}

// ToHexString writes this address as a single hexadecimal value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0x" prefix.
//
// If a subnet cannot be written as a single prefix block or a range of two values, an error is returned.
func (addr *IPv6Address) ToHexString(with0xPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	return addr.init().toHexString(with0xPrefix)
}

// ToOctalString writes this address as a single octal value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0" prefix.
//
// If a subnet cannot be written as a single prefix block or a range of two values, an error is returned.
func (addr *IPv6Address) ToOctalString(with0Prefix bool) (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	return addr.init().toOctalString(with0Prefix)
}

// ToBinaryString writes this address as a single binary value (possibly two values if a range that is not a prefixed block),
// the number of digits according to the bit count, with or without a preceding "0b" prefix.
//
// If a subnet cannot be written as a single prefix block or a range of two values, an error is returned.
func (addr *IPv6Address) ToBinaryString(with0bPrefix bool) (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	return addr.init().toBinaryString(with0bPrefix)
}

// ToUNCHostName Generates the Microsoft UNC path component for this address.  For examples see https://ipv6-literal.com/
//
// For IPv6, it is the canonical string but with colons replaced by dashes, percent signs with the letter “s”, and then appended with the root domain ".ipv6-literal.net".
func (addr *IPv6Address) ToUNCHostName() string {
	if addr == nil {
		return nilString()
	}
	cache := addr.getStringCache()
	if cache == nil {
		res, _ := addr.GetSection().toCustomString(uncParams, addr.zone)
		return res
	}
	var cacheField **string
	cacheField = &cache.uncString
	return cacheStr(cacheField,
		func() string {
			res, _ := addr.GetSection().toCustomString(uncParams, addr.zone)
			return res
		})
}

// ToBase85String creates the base 85 string, which is described by RFC 1924, "A Compact Representation of IPv6 Addresses".
// See https://www.rfc-editor.org/rfc/rfc1924.html
//
// It may be written as a range of two values if a range that is not a prefixed block.
//
// If a subnet cannot be written as a single prefix block or a range of two values, an error is returned.
func (addr *IPv6Address) ToBase85String() (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	if addr.hasZone() {
		cache := addr.getStringCache()
		if cache == nil {
			return addr.GetSection().toBase85String(addr.zone)
		}
		var cacheField **string
		cacheField = &cache.base85String
		return cacheStrErr(cacheField,
			func() (string, addrerr.IncompatibleAddressError) {
				return addr.GetSection().toBase85String(addr.zone)
			})
	}
	return addr.GetSection().ToBase85String()
}

// ToMixedString produces the mixed IPv6/IPv4 string.  It is the shortest such string (ie fully compressed).
// For some address sections with ranges of values in the IPv4 part of the address, there is not mixed string, and an error is returned.
func (addr *IPv6Address) ToMixedString() (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	if addr.hasZone() {
		cache := addr.getStringCache()
		if cache == nil {
			return addr.GetSection().toMixedStringZoned(addr.zone)
		}
		var cacheField **string
		cacheField = &cache.mixedString
		return cacheStrErr(cacheField,
			func() (string, addrerr.IncompatibleAddressError) {
				return addr.GetSection().toMixedStringZoned(addr.zone)
			})
	}
	return addr.GetSection().toMixedString()

}

// ToCustomString creates a customized string from this address or subnet according to the given string option parameters.
//
// Errors can result from split digits with ranged values, or mixed IPv4/v6 with ranged values, when a range cannot be split up.
// Options without split digits or mixed addresses do not produce errors.
// Single addresses do not produce errors.
func (addr *IPv6Address) ToCustomString(stringOptions addrstr.IPv6StringOptions) (string, addrerr.IncompatibleAddressError) {
	if addr == nil {
		return nilString(), nil
	}
	return addr.GetSection().toCustomString(stringOptions, addr.zone)
}

func (addr *IPv6Address) toMaxLower() *IPv6Address {
	return addr.init().addressInternal.toMaxLower().ToIPv6()
}

func (addr *IPv6Address) toMinUpper() *IPv6Address {
	return addr.init().addressInternal.toMinUpper().ToIPv6()
}

// ToAddressBase converts to an Address, a polymorphic type usable with all addresses and subnets.
// Afterwards, you can convert back with ToIPv6.
//
// ToAddressBase can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (addr *IPv6Address) ToAddressBase() *Address {
	if addr != nil {
		addr = addr.init()
	}
	return (*Address)(unsafe.Pointer(addr))
}

// ToIP converts to an IPAddress, a polymorphic type usable with all IP addresses and subnets.
//
// ToIP can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
func (addr *IPv6Address) ToIP() *IPAddress {
	if addr != nil {
		addr = addr.init()
	}
	return (*IPAddress)(addr)
}

// Wrap wraps this IP address, returning a WrappedIPAddress, an implementation of ExtendedIPSegmentSeries,
// which can be used to write code that works with both IP addresses and IP address sections.
// Wrap can be called with a nil receiver, wrapping a nil address.
func (addr *IPv6Address) Wrap() WrappedIPAddress {
	return wrapIPAddress(addr.ToIP())
}

// WrapAddress wraps this IP address, returning a WrappedAddress, an implementation of ExtendedSegmentSeries,
// which can be used to write code that works with both addresses and address sections.
// WrapAddress can be called with a nil receiver, wrapping a nil address.
func (addr *IPv6Address) WrapAddress() WrappedAddress {
	return wrapAddress(addr.ToAddressBase())
}

// ToKey creates the associated address key.
// While addresses can be compared with the Compare, TrieCompare or Equal methods as well as various provided instances of AddressComparator,
// they are not comparable with Go operators.
// However, AddressKey instances are comparable with Go operators, and thus can be used as map keys.
func (addr *IPv6Address) ToKey() IPv6AddressKey {
	addr = addr.init()
	key := IPv6AddressKey{}
	addr.toIPv6Key(&key.keyContents)
	return key
}

// ToGenericKey produces a generic Key[*IPv6Address] that can be used with generic code working with [Address], [IPAddress], [IPv4Address], [IPv6Address] and [MACAddress].
// ToKey produces a more compact key for code that is IPv6-specific.
func (addr *IPv6Address) ToGenericKey() Key[*IPv6Address] {
	// Note: We intentionally do not populate the "scheme" field, which is used with Key[*Address] and Key[*IPAddress] and 64-bit Key[*MACAddress].
	// With Key[*IPv6Address], by leaving the scheme zero, the zero Key[*IPv6Address] matches up with the key produced here by the zero address.
	// We do not need the scheme field for Key[*IPv6Address] since the generic type indicates IPv6.
	key := Key[*IPv6Address]{}
	addr.init().toIPv6Key(&key.keyContents)
	return key
}

func (addr *IPv6Address) fromKey(_ addressScheme, key *keyContents) *IPv6Address {
	// See ToGenericKey for details such as the fact that the scheme is ignored
	return fromIPv6IPKey(key)
}

func (addr *IPv6Address) toIPv6Key(contents *keyContents) {
	contents.zone = addr.GetZone()
	section := addr.GetSection()
	divs := section.getDivArray()
	if addr.IsMultiple() {
		for i, div := range divs {
			seg := div.ToIPv6()
			val := &contents.vals[i>>2]
			val.lower = (val.lower << IPv6BitsPerSegment) | uint64(seg.GetIPv6SegmentValue())
			val.upper = (val.upper << IPv6BitsPerSegment) | uint64(seg.GetIPv6UpperSegmentValue())
		}
	} else {
		for i, div := range divs {
			seg := div.ToIPv6()
			val := &contents.vals[i>>2]
			newLower := (val.lower << IPv6BitsPerSegment) | uint64(seg.GetIPv6SegmentValue())
			val.lower = newLower
			val.upper = newLower
		}
	}
}

func fromIPv6Key(key IPv6AddressKey) *IPv6Address {
	return fromIPv6IPKey(&key.keyContents)
}

func fromIPv6IPKey(contents *keyContents) *IPv6Address {
	return NewIPv6AddressFromZonedRange(
		func(segmentIndex int) IPv6SegInt {
			valsIndex := segmentIndex >> 2
			segIndex := ((IPv6SegmentCount - 1) - segmentIndex) & 0x3
			return IPv6SegInt(contents.vals[valsIndex].lower >> (segIndex << ipv6BitsToSegmentBitshift))
		},
		func(segmentIndex int) IPv6SegInt {
			valsIndex := segmentIndex >> 2
			segIndex := ((IPv6SegmentCount - 1) - segmentIndex) & 0x3
			return IPv6SegInt(contents.vals[valsIndex].upper >> (segIndex << ipv6BitsToSegmentBitshift))
		},
		string(contents.zone))
}
