//
// Copyright 2020-2022 Sean C Foley
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

// IPv6AddressConverter converts IP addresses to IPv6.
type IPv6AddressConverter interface {
	// ToIPv6 converts to IPv6.  If the given address is IPv6, or can be converted to IPv6, returns that IPv6Address.  Otherwise, returns nil.
	ToIPv6(address *IPAddress) *IPv6Address
}

// IPv4AddressConverter converts IP addresses to IPv4.
type IPv4AddressConverter interface {
	// ToIPv4 converts to IPv4.  If the given address is IPv4, or can be converted to IPv4, returns that IPv4Address.  Otherwise, returns nil.
	ToIPv4(address *IPAddress) *IPv4Address
}

// IPAddressConverter converts IP addresses to either IPv4 or IPv6.
type IPAddressConverter interface {
	IPv4AddressConverter

	IPv6AddressConverter

	// IsIPv4Convertible returns whether the address is IPv4 or can be converted to IPv4.  If true, ToIPv4 returns non-nil.
	IsIPv4Convertible(address *IPAddress) bool

	// IsIPv6Convertible returns whether the address is IPv6 or can be converted to IPv6.  If true, ToIPv6 returns non-nil.
	IsIPv6Convertible(address *IPAddress) bool
}

// DefaultAddressConverter converts to/from IPv4-mapped addresses, which maps IPv4 "a.b.c.d" to/from the IPv4-mapped IPv6 "::ffff:a.b.c.d".
// Converting from IPv6 to IPv4 requires that the IPV6 address have the prefix "0:0:0:0:0:ffff".
//
// Note that with some subnets, the mapping is not possible due to the range of values in segments.
// For example, "::ffff:0-100:0" cannot be mapped to an IPv4 address because the range 0-0x100 cannot be split into two smaller ranges.
// Similarly, "1-2.0.0.0" cannot be converted to an IPv4-mapped IPv6 address,
// because the two segments "1-2.0" cannot be joined into a single IPv6 segment with the same range of values, namely the two values 0x100 and 0x200.
type DefaultAddressConverter struct{}

var _ IPAddressConverter = DefaultAddressConverter{}

// ToIPv4 converts IPv4-mapped IPv6 addresses to IPv4, or returns the original address if IPv4 already, or returns nil if the address cannot be converted.
func (converter DefaultAddressConverter) ToIPv4(address *IPAddress) *IPv4Address {
	if addr := address.ToIPv4(); addr != nil {
		return addr
	} else if addr := address.ToIPv6(); addr != nil {
		if ipv4Addr, err := addr.GetEmbeddedIPv4Address(); err == nil {
			return ipv4Addr
		}
	}
	return nil
}

// ToIPv6 converts to an IPv4-mapped IPv6 address or returns the original address if IPv6 already.
func (converter DefaultAddressConverter) ToIPv6(address *IPAddress) *IPv6Address {
	if addr := address.ToIPv6(); addr != nil {
		return addr
	} else if addr := address.ToIPv4(); addr != nil {
		if ipv6Addr, err := addr.GetIPv4MappedAddress(); err == nil {
			return ipv6Addr
		}
	}
	return nil
}

// IsIPv4Convertible returns true if ToIPv4 returns non-nil.
func (converter DefaultAddressConverter) IsIPv4Convertible(address *IPAddress) bool {
	if addr := address.ToIPv6(); addr != nil {
		if addr.IsIPv4Mapped() {
			if _, _, _, _, err := addr.GetSegment(IPv6SegmentCount - 1).splitSegValues(); err != nil {
				return false
			} else if _, _, _, _, err := addr.GetSegment(IPv6SegmentCount - 2).splitSegValues(); err != nil {
				return false
			}
			return true
		}
	}
	return address.IsIPv4()
}

// IsIPv6Convertible returns true if ToIPv6 returns non-nil.
func (converter DefaultAddressConverter) IsIPv6Convertible(address *IPAddress) bool {
	if addr := address.ToIPv4(); addr != nil {
		return addr.GetSegment(0).isJoinableTo(addr.GetSegment(1)) && addr.GetSegment(2).isJoinableTo(addr.GetSegment(3))
	}
	return address.IsIPv6()
}
