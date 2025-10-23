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

type addrPtrType[S any, T any] interface {
	*S
	getAddrType() addrType
	init() T
}

// returns a slice of addresses that match the same IP version as the given
func filterSeries[S any, T addrPtrType[S, T]](addr T, orig []T) []T {
	addr = addr.init()
	addrType := addr.getAddrType()
	result := append(make([]T, 0, len(orig)+1), addr)
	for _, a := range orig {
		if addrType == a.getAddrType() {
			result = append(result, a.init())
		}
	}
	return result
}

func cloneIPAddrs(sect *IPAddress, orig []*IPAddress) []ExtendedIPSegmentSeries {
	converter := func(a *IPAddress) ExtendedIPSegmentSeries { return wrapIPAddress(a) }
	if sect == nil {
		return cloneTo(orig, converter)
	}
	return cloneToExtra(sect, orig, converter)
}

func cloneIPSections(sect *IPAddressSection, orig []*IPAddressSection) []ExtendedIPSegmentSeries {
	converter := func(a *IPAddressSection) ExtendedIPSegmentSeries { return wrapIPSection(a) }
	if sect == nil {
		return cloneTo(orig, converter)
	}
	return cloneToExtra(sect, orig, converter)
}

func cloneTo[T any, U any](orig []T, conv func(T) U) (result []U) {
	result = make([]U, len(orig))
	for i, v := range orig {
		result[i] = conv(v)
	}
	return
}

func cloneToExtra[T any, U any](sect T, orig []T, conv func(T) U) (result []U) {
	origCount := len(orig)
	result = make([]U, origCount+1)
	result[origCount] = conv(sect)
	for i, v := range orig {
		result[i] = conv(v)
	}
	return
}

// ToIPv4Slice converts a slice of subnets, addresses, or components thereof into IPv4-specific components.
// The original slice element type can be one of *Address, *IPAddress, *AddressDivisionGrouping,
// *AddressSection, *IPAddressSection, *AddressDivision, *AddressSegment, *IPAddressSegment, *SequentialRange,
// ExtendedIPSegmentSeries, WrappedIPAddress, WrappedIPAddressSection, ExtendedSegmentSeries, WrappedAddress, or WrappedAddressSection.
// Each slice element will be converted if the element originated as an IPv4 component, otherwise the element will be converted to nil in the returned slice.
func ToIPv4Slice[T interface {
	ToIPv4() IPv4Type
}, IPv4Type any](orig []T) []IPv4Type {
	return cloneTo(orig, func(a T) IPv4Type { return a.ToIPv4() })
}

// ToIPv6Slice converts a slice of subnets, addresses, or components thereof into IPv6-specific components.
// The original slice element type can be one of *Address, *IPAddress, *AddressDivisionGrouping,
// *AddressSection, *IPAddressSection, *AddressDivision, *AddressSegment, *IPAddressSegment, *SequentialRange,
// ExtendedIPSegmentSeries, WrappedIPAddress, WrappedIPAddressSection, ExtendedSegmentSeries, WrappedAddress, or WrappedAddressSection.
// Each slice element will be converted if the element originated as an IPv6 component, otherwise the element will be converted to nil in the returned slice.
func ToIPv6Slice[T interface {
	ToIPv6() IPv6Type
}, IPv6Type any](orig []T) []IPv6Type {
	return cloneTo(orig, func(a T) IPv6Type { return a.ToIPv6() })
}

// ToIPSlice converts a slice of subnets, addresses, or components thereof into IP-specific components.
// The original slice element type can be one of *Address, *IPv4Address, *IPv6Address, *AddressDivisionGrouping,
// *AddressSection, *IPv4AddressSection, *IPv6AddressSection, *AddressDivision, *AddressSegment, *IPv4AddressSegment, *IPv6AddressSegment, *IPAddressSegment,
// *SequentialRange, ExtendedSegmentSeries, WrappedAddress, WrappedAddressSection, WrappedIPAddress, IPAddressType, *SequentialRange, or IPAddressSeqRangeType.
// Each slice element will be converted if the element originated as an IP component, otherwise the element will be converted to nil in the returned slice.
func ToIPSlice[T interface {
	ToIP() IPType
}, IPType any](orig []T) []IPType {
	return cloneTo(orig, func(a T) IPType { return a.ToIP() })
}

// ToAddressBaseSlice converts a slice of subnets or addresses into general addresses or subnets not specific to a version or address type.
// The original slice element type can be one of *Address, *IPv4Address, *IPv6Address, *MACAddress, *IPAddress, or AddressType.
// Each slice element will be converted if the element originated as an address component, otherwise the element will be converted to nil in the returned slice.
func ToAddressBaseSlice[T interface {
	ToAddressBase() AddrType
}, AddrType any](orig []T) []AddrType {
	return cloneTo(orig, func(a T) AddrType { return a.ToAddressBase() })
}

// ToMACSlice converts a slice of subnets, addresses, or components thereof into MAC-specific components.
// The original slice element type can be one of *Address, *AddressDivisionGrouping,
// *AddressSection, *AddressDivision, *AddressSegment,
// ExtendedSegmentSeries, WrappedAddress, or WrappedAddressSection.
// Each slice element will be converted if the element originated as a MAC component, otherwise the element will be converted to nil in the returned slice.
func ToMACSlice[T interface {
	ToMAC() MACType
}, MACType any](orig []T) []MACType {
	return cloneTo(orig, func(a T) MACType { return a.ToMAC() })
}
