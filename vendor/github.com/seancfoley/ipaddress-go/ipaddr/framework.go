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

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
)

type BitItem interface {
	// GetByteCount returns the number of bytes required for each value comprising this address item,
	// rounding up if the bit count is not a multiple of 8.
	GetByteCount() int

	// GetBitCount returns the number of bits in each value comprising this address item.
	GetBitCount() BitCount
}

// AddressItem represents all addresses, division groupings, divisions, and sequential ranges.
// Any address item can be compared to any other.
type AddressItem interface {
	BitItem

	// GetValue returns the lowest individual address item in the address item range as an integer value.
	GetValue() *big.Int

	// GetUpperValue returns the highest individual address item in the address item range as an integer value.
	GetUpperValue() *big.Int

	// CopyBytes copies the value of the lowest individual address item in this address item range into a byte slice.
	//
	// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
	// Otherwise, a new slice is created and returned with the value.
	CopyBytes(bytes []byte) []byte

	// CopyUpperBytes copies the value of the highest individual address item in this address item range into a byte slice.
	//
	// If the value can fit in the given slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
	// Otherwise, a new slice is created and returned with the value.
	CopyUpperBytes(bytes []byte) []byte

	// Bytes returns the lowest individual address item in the address item range as a byte slice.
	Bytes() []byte

	// UpperBytes returns the highest individual address item in the address item range as a byte slice.
	UpperBytes() []byte

	// GetCount provides the number of address items represented by this AddressItem, for example the subnet size for IP addresses
	GetCount() *big.Int

	// IsMultiple returns whether this item represents multiple values (the count is larger than 1)
	IsMultiple() bool

	// IsFullRange returns whether this address item represents all possible values attainable by an address item of this type.
	//
	// This is true if and only if both IncludesZero and IncludesMax return true.
	IsFullRange() bool

	// IncludesZero returns whether this item includes the value of zero within its range.
	IncludesZero() bool

	// IncludesMax returns whether this item includes the max value, the value whose bits are all ones, within its range.
	IncludesMax() bool

	// IsZero returns whether this address item matches exactly the value of zero.
	IsZero() bool

	// IsMax returns whether this address item matches exactly the maximum possible value, the value whose bits are all ones.
	IsMax() bool

	// ContainsPrefixBlock returns whether the values of this item contains the prefix block for the given prefix length.
	// Unlike ContainsSinglePrefixBlock, whether there are multiple prefix values for the given prefix length makes no difference.
	ContainsPrefixBlock(BitCount) bool

	// ContainsSinglePrefixBlock returns whether the values of this series contains a single prefix block for the given prefix length.
	// This means there is only one prefix of the given length in this item, and this item contains the prefix block for that given prefix, all items with that same prefix.
	ContainsSinglePrefixBlock(BitCount) bool

	// GetPrefixLenForSingleBlock returns a prefix length for which there is only one prefix of that length in this item,
	// and the range of this item matches the block of all values for that prefix.
	//
	// If the entire range can be described this way, then this method returns the same value as GetMinPrefixLenForBlock.
	//
	// If no such prefix length exists, returns nil.
	//
	// If this item represents a single value, this returns the bit count.
	GetPrefixLenForSingleBlock() PrefixLen

	// GetMinPrefixLenForBlock returns the smallest prefix length possible such that this item includes the block of all values for that prefix length.
	//
	// If the entire range can be dictated this way, then this method returns the same value as GetPrefixLenForSingleBlock.
	//
	// There may be a single prefix, or multiple possible prefix values in this item for the returned prefix length.
	// Use GetPrefixLenForSingleBlock to avoid the case of multiple prefix values.
	//
	// If this item represents a single value, this returns the bit count.
	GetMinPrefixLenForBlock() BitCount

	// GetPrefixCountLen returns the count of the number of distinct values within the prefix part of the range of values for this item
	GetPrefixCountLen(BitCount) *big.Int

	// Compare returns a negative integer, zero, or a positive integer if this address item is less than, equal, or greater than the given item.
	// Any address item is comparable to any other.  All address items use CountComparator to compare.
	Compare(AddressItem) int

	// CompareSize compares the counts of two address items,
	// whether addresses in the subnet or address range, whether individual sections in the collection of sections, whether individual segments in the segment's range.
	// It compares the number of individual elements within each.
	//
	// Rather than calculating counts with GetCount, there can be more efficient ways of determining whether one item represents more individual addresses than another.
	//
	// CompareSize returns a positive integer if this item has a larger count than the one given, zero if they are the same, or a negative integer if the other has a larger count.
	CompareSize(AddressItem) int

	fmt.Stringer
	fmt.Formatter
}

type Prefixed interface {
	// IsPrefixed returns whether this item has an associated prefix length.
	IsPrefixed() bool

	// GetPrefixLen returns the prefix length, or nil if there is no prefix length.
	//
	// A prefix length indicates the number of bits in the initial part (most significant bits) of the series that comprise the prefix.
	//
	// A prefix is a part of the series that is not specific to that series but common amongst a group, such as a CIDR prefix block subnet.
	GetPrefixLen() PrefixLen

	// IsPrefixBlock returns whether this item has a prefix length and includes the block associated with that prefix length.
	// If the prefix length matches the bit count, this returns true.
	//
	// This is different from ContainsPrefixBlock in that this method returns
	// false if this item has no prefix length, or it has a prefix length that differs from a prefix length for which ContainsPrefixBlock returns true.
	IsPrefixBlock() bool

	// IsSinglePrefixBlock returns whether the range of values matches a single subnet block for the prefix length.
	//
	// This is different from ContainsSinglePrefixBlock in that this method returns
	// false if this series has no prefix length or a prefix length that differs from a prefix length for which ContainsSinglePrefixBlock returns true.
	IsSinglePrefixBlock() bool
}

type PrefixedConstraint[T any] interface {
	Prefixed

	// WithoutPrefixLen provides the same item but with no prefix length.  The values remain unchanged.
	WithoutPrefixLen() T

	// ToPrefixBlock returns the item whose prefix matches the prefix of this item, while the remaining bits span all values.
	// If this item has no prefix length, then this item is returned.
	//
	// The returned item will include all items with the same prefix as this item, known as the prefix "block".
	ToPrefixBlock() T

	// ToPrefixBlockLen returns the item associated with the prefix length provided,
	// the item whose prefix of that length matches the prefix of that length in this item, and the remaining bits span all values.
	//
	// The returned address will include all items with the same prefix as this one, known as the prefix "block".
	ToPrefixBlockLen(BitCount) T

	// SetPrefixLen sets the prefix length, returning a new item with the same values but with the new prefix length.
	//
	// A prefix length will not be set to a value lower than zero or beyond the bit length of the item.
	// The provided prefix length will be adjusted to these boundaries if necessary.
	SetPrefixLen(BitCount) T
}

// AddressDivisionSeries serves as a common interface to all division groupings, address sections, and addresses.
type AddressDivisionSeries interface {
	AddressItem

	// GetDivisionCount returns the number of divisions.
	GetDivisionCount() int

	// GetPrefixCount returns the count of prefixes in this series for its prefix length, or the total count if it has no prefix length
	GetPrefixCount() *big.Int

	// GetBlockCount returns the count of distinct values in the given number of initial (more significant) segments.
	GetBlockCount(divisionCount int) *big.Int

	// GetSequentialBlockIndex gets the minimal division index for which all following divisions are full-range blocks.
	//
	// The division at this index is not a full-range block unless all divisions are full-range.
	// The division at this index and all following divisions form a sequential range.
	// For the full series to be sequential, the preceding divisions must be single-valued.
	GetSequentialBlockIndex() int

	// GetSequentialBlockCount provides the count of elements from the sequential block iterator, the minimal number of sequential address division series that comprise this address division series.
	GetSequentialBlockCount() *big.Int

	// IsSequential returns  whether the series represents a range of values that are sequential.
	//
	// Generally, this means that any division covering a range of values must be followed by divisions that are full range, covering all values.
	IsSequential() bool

	Prefixed

	// GetGenericDivision returns the division at the given index as a DivisionType.
	// The first division is at index 0.
	// GetGenericDivision will panic given a negative index or index larger than the division count.
	GetGenericDivision(index int) DivisionType // useful for comparisons
}

var _ AddressDivisionSeries = &IPAddressLargeDivisionGrouping{}

// StandardDivGroupingType represents any standard division grouping (division groupings or address sections where all divisions are 64 bits or less)
// including [AddressSection], [IPAddressSection], [IPv4AddressSection], [IPv6AddressSection], [MACAddressSection], and [AddressDivisionGrouping]
type StandardDivGroupingType interface {
	AddressDivisionSeries

	// IsAdaptiveZero returns true if the division grouping was originally created as an implicitly zero-valued section or grouping (e.g. IPv4AddressSection{}),
	// meaning it was not constructed using a constructor function.
	// Such a grouping, which has no divisions or segments, is convertible to an implicitly zero-valued grouping of any type or version, whether IPv6, IPv4, MAC, or other.
	// In other words, when a section or grouping is the zero-value, then it is equivalent and convertible to the zero value of any other section or grouping type.
	IsAdaptiveZero() bool

	// ToDivGrouping converts to an AddressDivisionGrouping, a polymorphic type usable with all address sections and division groupings.
	//
	// ToDivGrouping implementations can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
	ToDivGrouping() *AddressDivisionGrouping
}

var _, _ StandardDivGroupingType = &AddressDivisionGrouping{},
	&IPv6v4MixedAddressGrouping{}

// AddressComponent represents all addresses, address sections, and address segments.
type AddressComponent interface { //AddressSegment and above, AddressSegmentSeries and above
	// TestBit returns true if the bit in the lower value of the address component at the given index is 1, where index 0 refers to the least significant bit.
	// In other words, it computes (bits & (1 << n)) != 0), using the lower value of this address component.
	// TestBit will panic if n < 0, or if it matches or exceeds the bit count of this address component.
	TestBit(index BitCount) bool

	// IsOneBit returns true if the bit in the lower value of this address component at the given index is 1, where index 0 refers to the most significant bit.
	// IsOneBit will panic if bitIndex is less than zero, or if it is larger than the bit count of this address component.
	IsOneBit(index BitCount) bool

	// ToHexString writes this address component as a single hexadecimal value (possibly two values if a range that is not a prefixed block),
	// the number of digits according to the bit count, with or without a preceding "0x" prefix.
	//
	// If a multiple-valued component cannot be written as a single prefix block or a range of two values, an error is returned.
	ToHexString(with0xPrefix bool) (string, addrerr.IncompatibleAddressError)

	// ToNormalizedString produces a string that is consistent for all address components of the same type and version.
	ToNormalizedString() string
}

// AddressSegmentSeries serves as a common interface to all address sections and addresses.
type AddressSegmentSeries interface { // Address and above, AddressSection and above, IPAddressSegmentSeries, ExtendedIPSegmentSeries
	AddressComponent

	AddressDivisionSeries

	// GetMaxSegmentValue returns the maximum possible segment value for this type of series.
	//
	// Note this is not the maximum of the range of segment values in this specific series,
	// this is the maximum value of any segment for this series type and version, determined by the number of bits per segment.
	GetMaxSegmentValue() SegInt

	// GetSegmentCount returns the number of segments, which is the same as the division count since the segments are also the divisions
	GetSegmentCount() int

	// GetBitsPerSegment returns the number of bits comprising each segment in this series.  Segments in the same series are equal length.
	GetBitsPerSegment() BitCount

	// GetBytesPerSegment returns the number of bytes comprising each segment in this series.  Segments in the same series are equal length.
	GetBytesPerSegment() int

	// ToCanonicalString produces a canonical string for the address series.
	//
	// For IPv4, dotted octet format, also known as dotted decimal format, is used.
	// https://datatracker.ietf.org/doc/html/draft-main-ipaddr-text-rep-00#section-2.1
	//
	// For IPv6, RFC 5952 describes the canonical string representation.
	// https://en.wikipedia.org/wiki/IPv6_address#Representation
	// http://tools.ietf.org/html/rfc5952
	//
	// For MAC, it uses the canonical standardized IEEE 802 MAC address representation of xx-xx-xx-xx-xx-xx.  An example is "01-23-45-67-89-ab".
	// For range segments, '|' is used: "11-22-33|44-55-66".
	//
	// Each address has a unique canonical string, not counting the prefix length.
	// With IP addresses and sections, the prefix length is included in the string, and the prefix length can cause two equal addresses to have different strings, for example "1.2.3.4/16" and "1.2.3.4".
	// It can also cause two different addresses to have the same string, such as "1.2.0.0/16" for the individual address "1.2.0.0" and also the prefix block "1.2.*.*".
	ToCanonicalString() string

	// ToNormalizedWildcardString produces a string similar to the normalized string but avoids the CIDR prefix length in the case of IP addresses.
	// Multiple-valued segments will be shown with wildcards and ranges (denoted by '*' and '-').
	ToNormalizedWildcardString() string

	// ToCompressedString produces a short representation of this series while remaining within the confines of standard representation(s) of the series.
	//
	// For IPv4, it is the same as the canonical string.
	//
	// For IPv6, it differs from the canonical string.  It compresses the maximum number of zeros and/or host segments with the IPv6 compression notation '::'.
	//
	// For MAC, it differs from the canonical string.  It produces a shorter string for the address that has no leading zeros.
	ToCompressedString() string

	// ToBinaryString writes this address series as a single binary value (possibly two values if a range that is not a prefixed block),
	// the number of digits according to the bit count, with or without a preceding "0b" prefix.
	//
	// If a multiple-valued series cannot be written as a single prefix block or a range of two values, an error is returned.
	ToBinaryString(with0bPrefix bool) (string, addrerr.IncompatibleAddressError)

	// ToOctalString writes this address series as a single octal value (possibly two values if a range that is not a prefixed block),
	// the number of digits according to the bit count, with or without a preceding "0" prefix.
	//
	// If a multiple-valued series cannot be written as a single prefix block or a range of two values, an error is returned.
	ToOctalString(withPrefix bool) (string, addrerr.IncompatibleAddressError)

	// GetSegmentStrings returns a slice with the string for each segment being the string that is normalized with wildcards.
	GetSegmentStrings() []string

	// GetGenericSegment returns the segment at the given index as an AddressSegmentType.
	// The first segment is at index 0.
	// GetGenericSegment will panic given a negative index or an index matching or larger than the segment count.
	GetGenericSegment(index int) AddressSegmentType
}

var _, _ AddressSegmentSeries = &Address{}, &AddressSection{}

// IPAddressSegmentSeries serves as a common interface to all IP address sections and IP addresses.
type IPAddressSegmentSeries interface { // IPAddress and above, IPAddressSection and above, ExtendedIPSegmentSeries
	AddressSegmentSeries

	// IncludesZeroHost returns whether the series contains an individual series with a host of zero.  If the series has no prefix length it returns false.
	// If the prefix length matches the bit count, then it returns true.
	//
	// Otherwise, it checks whether it contains an individual series for which all bits past the prefix are zero.
	IncludesZeroHost() bool

	// IncludesZeroHostLen returns whether the series contains an individual series with a host of zero, a series for which all bits past the given prefix length are zero.
	IncludesZeroHostLen(prefLen BitCount) bool

	// IncludesMaxHost returns whether the series contains an individual series with a host of all one-bits.  If the series has no prefix length it returns false.
	// If the prefix length matches the bit count, then it returns true.
	//
	// Otherwise, it checks whether it contains an individual series for which all bits past the prefix are one.
	IncludesMaxHost() bool

	// IncludesMaxHostLen returns whether the series contains an individual series with a host of all one-bits, a series for which all bits past the given prefix length are all ones.
	IncludesMaxHostLen(prefLen BitCount) bool

	// IsZeroHost returns whether this series has a prefix length and if so,
	// whether the host section is always zero for all individual series in this subnet or address section.
	//
	// If the host section is zero length (there are zero host bits), IsZeroHost returns true.
	IsZeroHost() bool

	// IsZeroHostLen returns whether the host section is always zero for all individual series in this address or address section,
	// for the given prefix length.
	//
	// If the host section is zero length (there are zero host bits), IsZeroHostLen returns true.
	IsZeroHostLen(BitCount) bool

	// IsMaxHost returns whether this address or address section has a prefix length and if so,
	// whether the host section is always all one-bits, the max value, for all individual series in this address or address section,
	//the host being the bits following the prefix.
	//
	// If the host section is zero length (there are zero host bits), IsMaxHost returns true.
	IsMaxHost() bool

	// IsMaxHostLen returns whether the host is all one-bits, the max value, for all individual series in this address or address section,
	// for the given prefix length, the host being the bits following the prefix.
	//
	// If the host is zero length (there are zero host bits), IsMaxHostLen returns true.
	IsMaxHostLen(BitCount) bool

	// IsSingleNetwork returns whether the network section of the IP address series, the prefix, consists of a single value.
	//
	// If it has no prefix length, it returns true if not multiple, if it contains only a single individual series.
	IsSingleNetwork() bool

	// GetIPVersion returns the IP version of this IP address or IP address section.
	GetIPVersion() IPVersion

	// GetBlockMaskPrefixLen returns the prefix length if this IP address or IP address section is equivalent to the mask for a CIDR prefix block.
	// Otherwise, it returns nil.
	// A CIDR network mask is a series with all ones in the network section and then all zeros in the host section.
	// A CIDR host mask is a series with all zeros in the network section and then all ones in the host section.
	// The prefix length is the bit-length of the network section.
	//
	// Also, keep in mind that the prefix length returned by this method is not equivalent to the prefix length of this instance,
	// indicating the network and host section of this series.
	// The prefix length returned here indicates the whether the value of this series can be used as a mask for the network and host
	// section of any other series.  Therefore, the two values can be different values, or one can be nil while the other is not.
	//
	// This method applies only to the lower value of the range if this series represents multiple values.
	GetBlockMaskPrefixLen(network bool) PrefixLen

	// GetLeadingBitCount returns the number of consecutive leading one or zero-bits.
	// If ones is true, returns the number of consecutive leading one-bits.
	// Otherwise, returns the number of consecutive leading zero bits.
	//
	// This method applies to the lower value of the range if this series represents multiple values.
	GetLeadingBitCount(ones bool) BitCount

	// GetTrailingBitCount returns the number of consecutive trailing one or zero-bits.
	// If ones is true, returns the number of consecutive trailing zero bits.
	// Otherwise, returns the number of consecutive trailing one-bits.
	//
	// This method applies to the lower value of the range if this series represents multiple values.
	GetTrailingBitCount(ones bool) BitCount

	// ToFullString produces a string with no compressed segments and all segments of full length with leading zeros.
	ToFullString() string

	// ToPrefixLenString returns a string with a CIDR network prefix length if this address has a network prefix length.
	// For IPv6, a zero host section will be compressed with "::". For IPv4 the string is equivalent to the canonical string.
	ToPrefixLenString() string

	// ToSubnetString produces a string with specific formats for subnets.
	// The subnet string looks like "1.2.*.*" or "1:2::/16".
	//
	// In the case of IPv4, this means that wildcards are used instead of a network prefix when a network prefix has been supplied.
	// In the case of IPv6, when a network prefix has been supplied, the prefix will be shown and the host section will be compressed with "::".
	ToSubnetString() string

	// ToCanonicalWildcardString produces a string similar to the canonical string but avoids the CIDR prefix length.
	// Series with a network prefix length will be shown with wildcards and ranges (denoted by '*' and '-') instead of using the CIDR prefix length notation.
	// IPv6 series will be compressed according to the canonical representation.
	ToCanonicalWildcardString() string

	// ToCompressedWildcardString produces a string similar to ToNormalizedWildcardString, avoiding the CIDR prefix, but with full IPv6 segment compression as well, including single zero-segments.
	// For IPv4 it is the same as ToNormalizedWildcardString.
	ToCompressedWildcardString() string

	// ToSegmentedBinaryString writes this IP address segment series as segments of binary values preceded by the "0b" prefix.
	ToSegmentedBinaryString() string

	// ToSQLWildcardString create a string similar to that from toNormalizedWildcardString except that
	// it uses SQL wildcards.  It uses '%' instead of '*' and also uses the wildcard '_'.
	ToSQLWildcardString() string

	// ToReverseDNSString generates the reverse-DNS lookup string,
	// returning an error if this address series is an IPv6 multiple-valued section for which the range cannot be represented.
	// For "8.255.4.4" it is "4.4.255.8.in-addr.arpa".
	// For "2001:db8::567:89ab" it is "b.a.9.8.7.6.5.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa".
	ToReverseDNSString() (string, addrerr.IncompatibleAddressError)
}

var _, _ IPAddressSegmentSeries = &IPAddress{}, &IPAddressSection{}

// IPv6AddressSegmentSeries serves as a common interface to all IPv6 address sections and IPv6 addresses.
type IPv6AddressSegmentSeries interface {
	IPAddressSegmentSeries

	// GetTrailingSection returns an ending subsection of the full address or address section
	GetTrailingSection(index int) *IPv6AddressSection

	// GetSubSection returns a subsection of the full address or address section
	GetSubSection(index, endIndex int) *IPv6AddressSection

	// GetNetworkSection returns an address section containing the segments with the network of the series, the prefix bits.
	// The returned section will have only as many segments as needed as determined by the existing CIDR network prefix length.
	//
	// If this series has no CIDR prefix length, the returned network section will
	// be the entire series as a prefixed section with prefix length matching the address bit length.
	GetNetworkSection() *IPv6AddressSection

	// GetHostSection returns a section containing the segments with the host of the series, the bits beyond the CIDR network prefix length.
	// The returned section will have only as many segments as needed to contain the host.
	//
	// If this series has no prefix length, the returned host section will be the full section.
	GetHostSection() *IPv6AddressSection

	// GetNetworkSectionLen returns a section containing the segments with the network of the series, the prefix bits according to the given prefix length.
	// The returned section will have only as many segments as needed to contain the network.
	//
	// The new section will be assigned the given prefix length,
	// unless the existing prefix length is smaller, in which case the existing prefix length will be retained.
	GetNetworkSectionLen(BitCount) *IPv6AddressSection

	// GetHostSectionLen returns a section containing the segments with the host of the series, the bits beyond the given CIDR network prefix length.
	// The returned section will have only as many segments as needed to contain the host.
	GetHostSectionLen(BitCount) *IPv6AddressSection

	// GetSegments returns a slice with the address segments.  The returned slice is not backed by the same array as the receiver.
	GetSegments() []*IPv6AddressSegment

	// CopySegments copies the existing segments into the given slice,
	// as much as can be fit into the slice, returning the number of segments copied.
	CopySegments(segs []*IPv6AddressSegment) (count int)

	// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
	// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
	CopySubSegments(start, end int, segs []*IPv6AddressSegment) (count int)

	// GetSegment returns the segment at the given index.
	// The first segment is at index 0.
	// GetSegment will panic given a negative index or an index matching or larger than the segment count.
	GetSegment(index int) *IPv6AddressSegment
}

var _, _, _ IPv6AddressSegmentSeries = &IPv6Address{},
	&IPv6AddressSection{},
	&EmbeddedIPv6AddressSection{}

// IPv4AddressSegmentSeries serves as a common interface to all IPv4 address sections and IPv4 addresses.
type IPv4AddressSegmentSeries interface {
	IPAddressSegmentSeries

	// GetTrailingSection returns an ending subsection of the full address section.
	GetTrailingSection(index int) *IPv4AddressSection

	// GetSubSection returns a subsection of the full address section.
	GetSubSection(index, endIndex int) *IPv4AddressSection

	// GetNetworkSection returns an address section containing the segments with the network of the series, the prefix bits.
	// The returned section will have only as many segments as needed as determined by the existing CIDR network prefix length.
	//
	// If this series has no CIDR prefix length, the returned network section will
	// be the entire series as a prefixed section with prefix length matching the address bit length.
	GetNetworkSection() *IPv4AddressSection

	// GetHostSection returns a section containing the segments with the host of the series, the bits beyond the CIDR network prefix length.
	// The returned section will have only as many segments as needed to contain the host.
	//
	// If this series has no prefix length, the returned host section will be the full section.
	GetHostSection() *IPv4AddressSection

	// GetNetworkSectionLen returns a section containing the segments with the network of the series, the prefix bits according to the given prefix length.
	// The returned section will have only as many segments as needed to contain the network.
	//
	// The new section will be assigned the given prefix length,
	// unless the existing prefix length is smaller, in which case the existing prefix length will be retained.
	GetNetworkSectionLen(BitCount) *IPv4AddressSection

	// GetHostSectionLen returns a section containing the segments with the host of the series, the bits beyond the given CIDR network prefix length.
	// The returned section will have only as many segments as needed to contain the host.
	GetHostSectionLen(BitCount) *IPv4AddressSection

	// GetSegments returns a slice with the address segments.  The returned slice is not backed by the same array as the receiver.
	GetSegments() []*IPv4AddressSegment

	// CopySegments copies the existing segments into the given slice,
	// as much as can be fit into the slice, returning the number of segments copied.
	CopySegments(segs []*IPv4AddressSegment) (count int)

	// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
	// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
	CopySubSegments(start, end int, segs []*IPv4AddressSegment) (count int)

	// GetSegment returns the segment at the given index.
	// The first segment is at index 0.
	// GetSegment will panic given a negative index or an index matching or larger than the segment count.
	GetSegment(index int) *IPv4AddressSegment
}

var _, _ IPv4AddressSegmentSeries = &IPv4Address{}, &IPv4AddressSection{}

// MACAddressSegmentSeries serves as a common interface to all MAC address sections and MAC addresses.
type MACAddressSegmentSeries interface {
	AddressSegmentSeries

	// GetTrailingSection returns an ending subsection of the full address section.
	GetTrailingSection(index int) *MACAddressSection

	// GetSubSection returns a subsection of the full address section.
	GetSubSection(index, endIndex int) *MACAddressSection

	// GetSegments returns a slice with the address segments.  The returned slice is not backed by the same array as the receiver.
	GetSegments() []*MACAddressSegment

	// CopySegments copies the existing segments into the given slice,
	// as much as can be fit into the slice, returning the number of segments copied.
	CopySegments(segs []*MACAddressSegment) (count int)

	// CopySubSegments copies the existing segments from the given start index until but not including the segment at the given end index,
	// into the given slice, as much as can be fit into the slice, returning the number of segments copied.
	CopySubSegments(start, end int, segs []*MACAddressSegment) (count int)

	// GetSegment returns the segment at the given index.
	// The first segment is at index 0.
	// GetSegment will panic given a negative index or an index matching or larger than the segment count.
	GetSegment(index int) *MACAddressSegment
}

var _, _ MACAddressSegmentSeries = &MACAddress{}, &MACAddressSection{}

// AddressSectionType represents any address section
// that can be converted to/from the base type AddressSection,
// including [AddressSection], [IPAddressSection], [IPv4AddressSection], [IPv6AddressSection], and [MACAddressSection].
type AddressSectionType interface {
	StandardDivGroupingType

	// Equal returns whether the given address section is equal to this address section.
	// Two address sections are equal if they represent the same set of sections.
	// They must match:
	//  - type/version (IPv4, IPv6, MAC, etc.)
	//  - segment counts
	//  - bits per segment
	//  - segment value ranges
	// Prefix lengths are ignored.
	Equal(AddressSectionType) bool

	// Contains returns whether this is same type and version as the given address section and whether it contains all values in the given section.
	//
	// Sections must also have the same number of segments to be comparable, otherwise false is returned.
	Contains(AddressSectionType) bool

	// Overlaps returns whether this section overlaps with another.
	//
	// Sections must have the same number of segments to be comparable.
	//
	//For sections which are aware of their position in an address (IPv6 and MAC), their respective positions must match to be comparable.
	Overlaps(AddressSectionType) bool

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
	Enumerate(AddressSectionType) *big.Int

	// PrefixEqual determines if the given section matches this section up to the prefix length of this section.
	// It returns whether the argument section has the same address section prefix values as this.
	//
	// The entire prefix of this section must be present in the other section to be comparable.
	PrefixEqual(AddressSectionType) bool

	// PrefixContains returns whether the prefix values in the given address section
	// are prefix values in this address section, using the prefix length of this section.
	// If this address section has no prefix length, the entire address is compared.
	//
	// It returns whether the prefix of this address contains all values of the same prefix length in the given address.
	//
	// All prefix bits of this section must be present in the other section to be comparable.
	PrefixContains(AddressSectionType) bool

	// ToSectionBase converts to an AddressSection, a polymorphic type usable with all address sections.
	//
	// ToSectionBase implementations can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
	ToSectionBase() *AddressSection
}

//Note: if we had an IPAddressSectionType we could add Wrap() WrappedIPAddressSection to it, but I guess not much else.

var _, _, _, _, _ AddressSectionType = &AddressSection{},
	&IPAddressSection{},
	&IPv4AddressSection{},
	&IPv6AddressSection{},
	&MACAddressSection{}

// AddressType represents any address, all of which can be represented by the base type [Address].
// This includes [IPAddress], [IPv4Address], [IPv6Address], and [MACAddress].
// You must use the pointer types *Address, *IPAddress, *IPv4Address, *IPv6Address, and *MACAddress when implementing AddressType.
// It can be useful as a parameter for functions to take any address type, while inside the function you can convert to [Address] using ToAddressBase.
type AddressType interface {
	AddressSegmentSeries

	// Equal returns whether the given address or subnet is equal to this address or subnet.
	// Two address instances are equal if they represent the same set of addresses.
	Equal(AddressType) bool

	// Contains returns whether this is the same type and version as the given address or subnet and whether it contains all addresses in the given address or subnet.
	Contains(AddressType) bool

	// Overlaps returns whether this is the same type and version as the given address and whether it overlaps with the other, containing at least one individual address common to both.
	Overlaps(AddressType) bool

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
	// If the given address does not have the same version or type, then nil is returned.
	Enumerate(AddressType) *big.Int

	// PrefixEqual determines if the given address matches this address up to the prefix length of this address.
	// If this address has no prefix length, the entire address is compared.
	//
	// It returns whether the two addresses share the same range of prefix values.
	PrefixEqual(AddressType) bool

	// PrefixContains returns whether the prefix values in the given address or subnet
	// are prefix values in this address or subnet, using the prefix length of this address or subnet.
	// If this address has no prefix length, the entire address is compared.
	//
	// It returns whether the prefix of this address contains all values of the same prefix length in the given address.
	PrefixContains(AddressType) bool

	// ToAddressBase converts to an Address, a polymorphic type usable with all addresses and subnets.
	//
	// ToAddressBase implementations can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
	ToAddressBase() *Address
}

var _, _ AddressType = &Address{}, &MACAddress{}

// IPAddressRange represents all IPAddress instances and all IPAddress sequential range instances.
type IPAddressRange interface {
	// GetIPVersion returns the IP version of this IP address range
	GetIPVersion() IPVersion

	// GetLowerIPAddress returns the address in the subnet or address range with the lowest numeric value,
	// which will be the receiver if it represents a single address.
	// For example, for "1.2-3.4.5-6", the series "1.2.4.5" is returned.
	GetLowerIPAddress() *IPAddress

	// GetUpperIPAddress returns the address in the subnet or address range with the highest numeric value,
	// which will be the receiver if it represents a single address.
	// For example, for the subnet "1.2-3.4.5-6", the address "1.3.4.6" is returned.
	GetUpperIPAddress() *IPAddress

	// CopyNetIP copies the value of the lowest individual address in the subnet or address range into a net.IP.
	//
	// If the value can fit in the given net.IP slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
	// Otherwise, a new slice is created and returned with the value.
	CopyNetIP(bytes net.IP) net.IP

	// CopyUpperNetIP copies the value of the highest individual address in the subnet or address range into a net.IP.
	//
	// If the value can fit in the given net.IP slice, the value is copied into that slice and a length-adjusted sub-slice is returned.
	// Otherwise, a new slice is created and returned with the value.
	CopyUpperNetIP(bytes net.IP) net.IP

	// GetNetIP returns the lowest address in this subnet or address range as a net.IP.
	GetNetIP() net.IP

	// GetUpperNetIP returns the highest address in this subnet or address range as a net.IP.
	GetUpperNetIP() net.IP

	// GetNetNetIPAddr returns the lowest address in this subnet or address range as a netip.Addr.
	GetNetNetIPAddr() netip.Addr

	// GetUpperNetNetIPAddr returns the highest address in this subnet or address range as a netip.Addr.
	GetUpperNetNetIPAddr() netip.Addr

	// IsSequential returns whether the address item represents a range of addresses that are sequential.
	//
	// IP Address sequential ranges are sequential by definition.
	//
	// Generally, for a subnet this means that any segment covering a range of values must be followed by segments that are full range, covering all values.
	//
	// Individual addresses are sequential and CIDR prefix blocks are sequential.
	// The subnet "1.2.3-4.5" is not sequential, since the two addresses it represents, "1.2.3.5" and "1.2.4.5", are not ("1.2.3.6" is in-between the two but not in the subnet).
	IsSequential() bool
}

var _, _, _, _, _, _ IPAddressRange = &IPAddress{},
	&IPv4Address{},
	&IPv6Address{},
	&SequentialRange[*IPAddress]{},
	&SequentialRange[*IPv4Address]{},
	&SequentialRange[*IPv6Address]{}

// IPAddressType represents any IP address, all of which can be represented by the base type [IPAddress].
// This includes [IPv4Address] and [IPv6Address].
// You must use the pointer types *IPAddress, *IPv4Address, and *IPv6Address when implementing IPAddressType.
type IPAddressType interface {
	AddressType

	IPAddressRange

	// ContainsRange returns whether all the addresses in the given sequential range are also contained in this sequential range.
	ContainsRange(IPAddressSeqRangeType) bool

	// Overlaps returns whether this IP address is the same version as the given range and whether it overlaps with the given range, containing at least one individual address common to both.
	OverlapsRange(IPAddressSeqRangeType) bool

	// Wrap wraps this IP address, returning a WrappedIPAddress, an implementation of ExtendedIPSegmentSeries,
	// which can be used to write code that works with both IP addresses and IP address sections.
	Wrap() WrappedIPAddress

	// ToIP converts to an IPAddress, a polymorphic type usable with all IP addresses and subnets.
	//
	// ToIP can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
	ToIP() *IPAddress

	// ToAddressString retrieves or generates an IPAddressString instance for this IP address.
	// This may be the IPAddressString this instance was generated from, if it was generated from an IPAddressString.
	//
	// In general, users are intended to create IP address instances from IPAddressString instances,
	// while the reverse direction, calling this method, is generally not encouraged and not useful, except under specific circumstances.
	//
	// Those specific circumstances may include when maintaining a collection of HostIdentifierString or IPAddressString instances.
	ToAddressString() *IPAddressString
}

var _, _, _ IPAddressType = &IPAddress{},
	&IPv4Address{},
	&IPv6Address{}

// IPAddressSeqRangeType represents any IP address sequential range, all of which can be represented by the base type IPAddressSeqRange.
// This includes IPv4AddressSeqRange and IPv6AddressSeqRange.
type IPAddressSeqRangeType interface {
	AddressItem

	IPAddressRange

	// ContainsRange returns whether all the addresses in the given sequential range are also contained in this sequential range.
	ContainsRange(IPAddressSeqRangeType) bool

	// Contains returns whether this range contains all IP addresses in the given address or subnet.
	Contains(IPAddressType) bool // this is not in IPAddressRange because addresses use Contains(AddressType)

	// OverlapsAddress indicates whether this range is the same type and version as the given address and whether it overlaps with the given address, containing at least one individual address common to both.
	OverlapsAddress(IPAddressType) bool

	// Enumerate indicates where an address sits relative to the range ordering.
	//
	// Determines how many address elements of a range precede the given address element, if the address is in the range.
	// If above the range, it is the distance to the upper boundary added to the range count less one, and if below the range, the distance to the lower boundary.
	//
	// In other words, if the given address is not in the range but above it, returns the number of addresses preceding the address from the upper range boundary,
	// added to one less than the total number of range addresses.  If the given address is not in the subnet but below it, returns the number of addresses following the address to the lower subnet boundary.
	//
	// Returns nil when the argument is multi-valued. The argument must be an individual address.
	//
	// If the given address is not the same version, then nil is returned.
	Enumerate(IPAddressType) *big.Int // this is not in IPAddressRange because addresses use Enumerate(AddressType)

	// Equal returns whether the given sequential address range is equal to this sequential address range.
	// Two sequential address ranges are equal if their lower and upper range boundaries are equal.
	Equal(IPAddressSeqRangeType) bool

	// ToCanonicalString produces a canonical string for the address range.
	// It has the format "lower -> upper" where lower and upper are the canonical strings for the lowest and highest addresses in the range, given by GetLower and GetUpper.
	ToCanonicalString() string

	// ToNormalizedString produces a normalized string for the address range.
	// It has the format "lower -> upper" where lower and upper are the normalized strings for the lowest and highest addresses in the range, given by GetLower and GetUpper.
	ToNormalizedString() string

	// ToIP converts to an IPAddressSeqRange, a polymorphic type usable with all IP address sequential ranges.
	//
	// ToIP can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
	ToIP() *SequentialRange[*IPAddress]
}

var _, _, _ IPAddressSeqRangeType = &SequentialRange[*IPAddress]{},
	&SequentialRange[*IPv4Address]{},
	&SequentialRange[*IPv6Address]{}

// HostIdentifierString represents a string that is used to identify a host.
type HostIdentifierString interface {

	// ToNormalizedString provides a normalized String representation for the host identified by this HostIdentifierString instance.
	ToNormalizedString() string

	// IsValid returns whether the wrapped string is a valid identifier for a host.
	IsValid() bool

	// Wrap wraps the identifier string into the extended type that is polymorphic with other identifier strings.
	Wrap() ExtendedIdentifierString

	fmt.Stringer
	fmt.Formatter
}

var _, _, _ HostIdentifierString = &IPAddressString{}, &MACAddressString{}, &HostName{}
