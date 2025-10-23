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

// DivisionType serves as a common interface to all divisions
type DivisionType interface {
	AddressItem

	getAddrType() addrType

	// getStringAsLower caches the string from getDefaultLowerString.
	getStringAsLower() string

	// GetString produces a string that avoids wildcards when a prefix length is part of the string.  Equivalent to GetWildcardString when the prefix length is not part of the string.
	GetString() string

	// GetWildcardString produces a string that uses wildcards and avoids prefix length.
	GetWildcardString() string

	// IsSinglePrefix determines if the division has a single prefix for the given prefix length.  You can call GetPrefixCountLen to get the count of prefixes.
	IsSinglePrefix(BitCount) bool

	// methods for string generation used by the string params and string writer.
	divStringProvider
}

var _ DivisionType = &IPAddressLargeDivision{}

// StandardDivisionType represents any standard address division, which is a division of size 64 bits or less.
// All can be converted to/from [AddressDivision].
type StandardDivisionType interface {
	DivisionType

	// ToDiv converts to an AddressDivision, a polymorphic type usable with all address segments and divisions.
	//
	// ToDiv implementations can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
	ToDiv() *AddressDivision
}

var _ StandardDivisionType = &AddressDivision{}

// AddressSegmentType serves as a common interface to all segments, including [AddressSegment], [IPAddressSegment], [IPv6AddressSegment], [IPv4AddressSegment] and [MACAddressSegment].
type AddressSegmentType interface {
	AddressComponent

	StandardDivisionType

	// Equal returns whether the given segment is equal to this segment.
	// Two segments are equal if they match:
	//  - type/version (IPv4, IPv6, MAC)
	//  - value range
	// Prefix lengths are ignored.
	Equal(AddressSegmentType) bool

	// Contains returns whether this segment is same type and version as the given segment and whether it contains all values in the given segment.
	Contains(AddressSegmentType) bool

	// Overlaps returns whether this segment is same type and version as the given segment and whether it overlaps with the values in the given segment.
	Overlaps(AddressSegmentType) bool

	// GetSegmentValue returns the lower value of the segment value range as a SegInt.
	GetSegmentValue() SegInt

	// GetUpperSegmentValue returns the upper value of the segment value range as a SegInt.
	GetUpperSegmentValue() SegInt

	// ToSegmentBase converts to an AddressSegment, a polymorphic type usable with all address segments.
	//
	// ToSegmentBase implementations can be called with a nil receiver, enabling you to chain this method with methods that might return a nil pointer.
	ToSegmentBase() *AddressSegment
}

var _, _, _, _, _ AddressSegmentType = &AddressSegment{},
	&IPAddressSegment{},
	&IPv6AddressSegment{},
	&IPv4AddressSegment{},
	&MACAddressSegment{}
