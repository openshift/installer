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

package addrerr

/*
Error hierarchy:

AddressError
	-IncompatibleAddressError
		- SizeMismatchError
	- HostIdentifierError
		- HostNameError
		- AddressStringError
	- AddressValueError

unused, but present in Java:
NetworkMismatchException
InconsistentPrefixException
AddressPositionException
AddressConversionException
PrefixLenException
PositionMismatchException
*/

// AddressError is a type used by all library errors in order to be able to provide internationalized error messages.
type AddressError interface {
	error

	// GetKey allows users to implement their own i18n error messages.
	// The keys and mappings are listed in IPAddressResources.properties,
	// so users of this library need only provide translations and implement
	// their own method of i18n to incorporate those translations,
	// such as the method provided by golang.org/x/text
	GetKey() string
}

// HostIdentifierError represents errors in string formats used to identify hosts.
type HostIdentifierError interface {
	AddressError
}

// AddressStringError represents errors in address string formats used to identify addresses.
type AddressStringError interface {
	HostIdentifierError
}

// HostNameError represents errors in host name string formats used to identify hosts.
type HostNameError interface {
	HostIdentifierError

	// GetAddrError returns the underlying address error, or nil if none.
	GetAddrError() AddressError
}

// IncompatibleAddressError represents situations when an address, address section, address segment, or address string represents a valid type or format but
// that type does not match the required type or format for a given operation.
//
// All such occurrences occur only from subnet addresses and sections.  These occurrences cannot happen with single-valued address objects.
// These occurrences cannot happen when using a standard prefix block subnet with standard masks.
//
// Examples include:
//
//   - producing non-segmented hex, octal or base 85 strings from a subnet with a range that cannot be represented as a single range of values,
//   - masking subnets in a way that produces a non-contiguous range of values in a segment,
//   - reversing values that are not reversible,
//   - producing strings that are single-segment ranges from subnets which cannot be represented that way,
//   - producing new formats for which the range of values are incompatible with the new segments (EUI-64, IPv4 inet_aton formats, IPv4 embedded within IPv6, dotted MAC addresses from standard mac addresses, reverse-DNS strings), or
//   - using a subnet for an operation that requires a single address, such as with ToCanonicalHostName in IPAddress
type IncompatibleAddressError interface {
	AddressError
}

// SizeMismatchError is an error that results from attempting an operation that requires address items of equal size,
// but the supplied arguments were not equal in size.
type SizeMismatchError interface {
	IncompatibleAddressError
}

// AddressValueError results from supplying an invalid value to an address operation.
// Used when an address or address component would be too large or small,
// when a prefix length is too large or small, or when prefixes across segments are inconsistent.
// Not used when constructing new address components.
// Not used when parsing strings to construct new address components, in which case AddressStringError is used instead.
type AddressValueError interface {
	AddressError
}
