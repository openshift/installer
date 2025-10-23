//
// Copyright 2020-2023 Sean C Foley
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

/*
IPAddress is a library for handling IP addresses and subnets, both IPv4 and IPv6.

# Benefits of this Library

The primary goals are:
  - Parse all IPv4 and IPv6 address formats and host name formats in common usage, plus some additional formats.
  - Parse and represent subnets, either those specified by network prefix length or those specified with ranges of segment values.
  - Allow the separation of address parsing from host parsing.
  - Allow control over which formats are allowed when parsing, whether IPv4, IPv6, subnet formats, inet_aton formats, or other.
  - Produce all common address strings of different formats for a given IPv4 or IPv6 address and produce collections of such strings.
  - Parse all common MAC Address formats in usage and produce all common MAC address strings of different formats.
  - Integrate MAC addresses with IPv6 with standard conversions.
  - Integrate IPv4 Addresses with IPv6 through common address conversions.
  - Polymorphism is a key goal. The library maintains an address framework of interfaces that allow most library functionality to be independent of address type or version, whether IPv4, IPv6 or MAC. This allows for code which supports both IPv4 and IPv6 transparently.
  - Provide concurrency-safety, primarily through immutability. The core types (host names, address strings, addresses, address sections, address segments, address ranges) are all immutable. They do not change their underlying value. For sharing amongst goroutines, this is valuable.
  - Modify addresses, such as altering prefix lengths, masking, splitting into sections and segments, splitting into network and host sections, reconstituting from sections and segments.
  - Provide address operations and subnetting, such as obtaining the prefix block subnet for a prefixed address, iterating through subnets, iterating through prefixes, prefix blocks, or segment blocks of subnets, incrementing and decrementing addresses by integer values, reversing address bits for endianness or DNS lookup, set-subtracting subnets from other subnets, subnetting, intersections of subnets, merging subnets, checking containment of addresses in subnets, listing subnets covering a span of addresses.
  - Sorting and comparison of host names, addresses, address strings and subnets.  All the address component types are compararable.
  - Integrate with the Go language primitive types and the standard library types [net.IP], [net.IPAddr], [net.IPMask], [net.IPNet], [net.TCPAddr], [net.UDPAddr], [net/netip.Addr], [net/netip.Prefix], [net/netip.AddrPort], and [math/big.Int].
  - Make address manipulations easy, so you do not have to worry about longs/ints/shorts/bytes/bits, signed/unsigned, sign extension, ipv4/v6, masking, iterating, and other implementation details.

# Design Overview

The core types are [IPAddressString], [HostName], and [MACAddressString] along with the [Address] base type and its associated types [IPAddress], [IPv4Address], [IPv6Address], and [MACAddress], as well as the sequential address type [SequentialRange].
If you have a textual representation of an IP address, then start with [IPAddressString] or [HostName].  If you have a textual representation of a MAC address, then start with [MACAddressString].
Note that address instances can represent either a single address or a subnet. If you have either an address or host name, or you have something with a port or service name, then use [HostName].
If you have numeric bytes or integers, then start with [IPv4Address], [IPv6Address], [IPAddress], or [MACAddress].

This library allows you to scale down from more specific address types to more generic address types, and then to scale back up again.
The polymorphism is useful for IP-version ambiguous code.  The most-specific types allow for method sets tailored to the address version or type.
You can only scale up to a specific version or address type if the more generic instance was originally derived from an instance of the specific type.
So, for instance, an [IPv6Address] can be converted to an [IPAddress] using [IPv6Address.ToIP], or to an [Address] using [IPv6Address.ToAddressBase], which can then be converted back to [IPAddress] or an [IPv6Address] using [Address.ToIP] or [Address.ToIPv6].
But that [IPv6Address] cannot be scaled back to IPv4.  If you wish to convert that [IPv6Address] to IPv4, you would need to use an implementation of [IPv4AddressConverter].

This library has some similarities in design to the [Java IPAddress library].
Notable divergences derive from the differences between the Java and Go languages,
such as the differences in error handling and the lack of inheritance in Go, amongst many other differences.
Other divergences derive from common Go language idioms and practices which differ from standard Java idioms and practices.
Some similarities include the inclusion of an address framework of interfaces, and the data structures in use such as tries, and the segment/section/address architecture of addresses.
Both share many of the same operations, such as spanning with prefix blocks, merging into prefix blocks, iterating subnets, and so on.

# Code Examples

For common use-cases, you may wish to go straight to the [wiki code examples] which cover a wide breadth of common use-cases.

# Further Documentation

You can read [further documentation] with more depth.

[Java IPAddress library]: https://github.com/seancfoley/IPAddress
[wiki code examples]: https://github.com/seancfoley/ipaddress-go/wiki/Code-Examples
[further documentation]: https://seancfoley.github.io/IPAddress/
*/
package ipaddr
