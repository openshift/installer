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
	"strings"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstrparam"
)

// NewIPAddressStringParams constructs an IPAddressString that will parse the given string according to the given parameters.
func NewIPAddressStringParams(str string, params addrstrparam.IPAddressStringParams) *IPAddressString {
	var p addrstrparam.IPAddressStringParams
	if params == nil {
		p = defaultIPAddrParameters
	} else {
		p = addrstrparam.CopyIPAddressStringParams(params)
	}
	return parseIPAddressString(str, p)
}

// NewIPAddressString constructs an IPAddressString.
func NewIPAddressString(str string) *IPAddressString {
	return parseIPAddressString(str, defaultIPAddrParameters)
}

func newIPAddressStringFromAddr(str string, addr *IPAddress) *IPAddressString {
	return &IPAddressString{
		str:             str,
		addressProvider: addr.getProvider(),
	}
}

func parseIPAddressString(str string, params addrstrparam.IPAddressStringParams) *IPAddressString {
	str = strings.TrimSpace(str)
	res := &IPAddressString{
		str: str,
	}
	res.validate(params)
	return res
}

var validator hostIdentifierStringValidator = strValidator{}

var defaultIPAddrParameters = new(addrstrparam.IPAddressStringParamsBuilder).ToParams()

// IPAddressString parses the string representation of an IP address.  Such a string can represent just a single address like "1.2.3.4" or "1:2:3:4:6:7:8", or a subnet like "1.2.0.0/16" or "1.*.1-3.1-4" or "1111:222::/64".
//
// This supports a wide range of address string formats.  It supports subnet formats, provides specific error messages, and allows more specific configuration.
//
// You can control all the supported formats using an IPAddressStringParamsBuilder to build a parameters instance of IPAddressStringParams.
// When no IPAddressStringParams is supplied, a default instance of IPAddressStringParams is used that is generally permissive.
//
// # Supported Formats
//
// Both IPv4 and IPv6 are supported.
//
// Subnets are supported:
//   - wildcards '*' and ranges '-' (for example "1.*.2-3.4"), useful for working with subnets
//   - the wildcard '*' can span multiple segments, so you can represent all addresses with '*', all IPv4 with '*.*', or all IPv6 with '*:*'
//   - SQL wildcards '%' and '_', although '%' is considered an SQL wildcard only when it is not considered an IPv6 zone indicator
//   - CIDR network prefix length addresses, like "1.2.0.0/16", which is equivalent to "1.2.*.*" (all-zero hosts are the full subnet, non-zero hosts are single addresses)
//   - address/mask pairs, in which the mask is applied to the address, like "1.2.3.4/255.255.0.0", which is also equivalent to "1.2.*.*"
//
// You can combine these variations, such as "1.*.2-3.4/255.255.255.0".
//
// IPv6 is fully supported:
//   - IPv6 addresses like "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"
//   - IPv6 zones or scope identifiers, like "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff%zone"
//   - IPv6 mixed addresses are supported, which are addresses for which the last two IPv6 segments are represented as IPv4, like "ffff:ffff:ffff:ffff:ffff:ffff:255.255.255.255"
//   - IPv6 compressed addresses like "::1"
//   - A single value of 32 hex digits like "00aa00bb00cc00dd00ee00ff00aa00bb" with or without a preceding hex delimiter "0x"
//   - A base 85 address comprising 20 base 85 digits like "4)+k&amp;C#VzJ4br&gt;0wv%Yp" as in RFC 1924 https://tools.ietf.org/html/rfc1924
//   - Binary, preceded by "0b", either with binary segments that comprise all 16 bits like "::0b0000111100001111" or a single segment address of "0b" followed by 128 binary bits.
//
// All of the above subnet variations work for IPv6, whether network prefix lengths, masks, ranges or wildcards.
// Similarly, all the above subnet variations work for any supported IPv4 format, such as the standard dotted-decimal IPv4 format as well as the inet_aton formats listed below.
//
// This type support all address formats of the C routine inet_pton and the Java method java.net.InetAddress.getByName.
// This type supports all IPv4 address formats of the C routine inet_aton as follows:
//   - IPv4 hex: "0x1.0x2.0x3.0x4" ("0x" prefix)
//   - IPv4 octal: "01.02.03.0234".  Note this clashes with the same address interpreted as dotted decimal
//   - 3-part IPv4: "1.2.3" (which is interpreted as "1.2.0.3" (ie the third part covers the last two)
//   - 2-part IPv4: "1.2" (which is interpreted as "1.0.0.2" (ie the 2nd part covers the last 3)
//   - 1-part IPv4: "1" (which is interpreted as "0.0.0.1" (ie the number represents all 4 segments, and can be any number of digits less than the 32 digits which would be interpreted as IPv6)
//   - hex or octal variants of 1, 2, and 3 part, such as "0xffffffff" (which is interpreted as "255.255.255.255")
//
// Also supported are binary segments of a "0b" followed by binary digits like "0b1.0b1010.2.3", or a single segment address of "0b" followed by all 32 bits.
//
// inet_aton (and this type) allows mixing octal, hex and decimal (e.g. "0xa.11.013.11" which is equivalent to "11.11.11.11").
// String variations using prefixes, masks, ranges, and wildcards also work for inet_aton style.
// The same can be said of binary segments, they can be mixed with all other formats.
//
// Note that there is ambiguity when supporting both inet_aton octal and dotted-decimal leading zeros, like "010.010.010.010" which can
// be interpreted as octal or decimal, thus it can be either "8.8.8.8" or "10.10.10.10", with the default behaviour using the former interpretation.
// This behaviour can be controlled by IPAddressStringParamsBuilder.GetIPv4AddressParamsBuilder and
// IPv4AddressStringParametersBuilder.allowLeadingZeros(boolean)
//
// Some Additional Formats:
//   - empty strings are interpreted as the zero-address or the loopback
//   - as noted previously, the single wildcard address "*" represents all addresses both ipv4 and ipv6,
//
// although you need to give it some help when converting to [IPAddress] by specifying the IP version in GetVersionedAddress(IPVersion) or ToVersionedAddress(IPVersion).
//
// If you have an address in which segments have been delimited with commas, such as "1,2.3.4,5.6", you can parse this with ParseDelimitedSegments(string)
// which gives an iterator of strings.  For "1,2.3.4,5.6" you will iterate through "1.3.4.6", "1.3.5.6", "2.3.4.6" and "2.3.5.6".
// You can count the number of elements in such an iterator with CountDelimitedAddresses(String).
// Each string can then be used to construct an IPAddressString.
//
// # Usage
//
// Once you have constructed an IPAddressString object, you can convert it to an [IPAddress] object with various methods.
//
// Most address strings can be converted to an [IPAddress] object using GetAddress or ToAddress.  In most cases the IP version is determined by the string itself.
//
// There are a few exceptions, cases in which the version is unknown or ambiguous, for which GetAddress returns nil:
//   - strings which do not represent valid addresses (eg "bla")
//   - the "all" address "*" which represents all IPv4 and IPv6 addresses.  For this string you can provide the IPv4/IPv6 version to GetVersionedAddress to get an address representing either all IPv4 or all IPv6 addresses.
//   - empty string "" is interpreted as the zero-address, or optionally the default loopback address.  You can provide the IPv4/IPv6 version to GetVersionedAddress to get the version of your choice.
//
// The other exception is a subnet in which the range of values in a segment of the subnet are not sequential, for which ToAddress returns IncompatibleAddressError because there is no single [IPAddress] value, there would be many.
// An [IPAddress] instance requires that all segments can be represented as a range of values.
//
// There are only two unusual circumstances when this can occur:
//   - using masks on subnets specified with wildcard or range characters causing non-sequential segments such as the final IPv4 segment of "0.0.0.*" with mask "0.0.0.128",
//     this example translating to the two addresses "0.0.0.0" and "0.0.0.128", so the last IPv4 segment cannot be represented as a sequential range of values.
//   - using wildcards or range characters in the IPv4 section of an IPv6 mixed address causing non-sequential segments such as the last IPv6 segment of "::ffff:0.0.*.0",
//     this example translating to the addresses "::ffff:0:100", "::ffff:0:200", "::ffff:0:300", ..., so the last IPv6 segment cannot be represented as a sequential range of values.
//
// These exceptions do not occur with non-subnets (ie individual addresses), nor can they occur with standard CIDR prefix-based subnets.
//
// This type is concurrency-safe.  In fact, IPAddressString objects are immutable.
// An IPAddressString object represents a single IP address representation that cannot be changed after construction.
// Some derived state is created upon demand and cached, such as the derived [IPAddress] instances.
//
// This type has a few methods with analogs in [IPAddress], such as Contains, GetSequentialRange,
// PrefixEqual, IsIPv4, and IsIPv6.
// Such methods are provided to make creating the [IPAddress] instance unnecessary when no such [IPAddress] instance is needed for other reasons.
type IPAddressString struct {
	str             string
	addressProvider ipAddressProvider
	validateError   addrerr.AddressStringError
}

func (addrStr *IPAddressString) init() *IPAddressString {
	if addrStr.addressProvider == nil && addrStr.validateError == nil {
		return zeroIPAddressString
	}
	return addrStr
}

// GetValidationOptions returns the validation options supplied when constructing this address string,
// or the default options if no options were supplied.  It returns nil if no options were used to construct.
func (addrStr *IPAddressString) GetValidationOptions() addrstrparam.IPAddressStringParams {
	provider, _ := addrStr.getAddressProvider()
	if provider != nil {
		return provider.getParameters()
	}
	return nil
}

// IsPrefixed returns whether this address string has an associated prefix length.
// If so, the prefix length is given by GetNetworkPrefixLen.
func (addrStr *IPAddressString) IsPrefixed() bool {
	return addrStr.getNetworkPrefixLen() != nil
}

// GetNetworkPrefixLen returns the associated network prefix length.
//
// If this address is a valid address with an associated network prefix length then this returns that prefix length, otherwise returns nil.
// The prefix length may be expressed explicitly with the notation "/xx" where xx is a decimal value, or it may be expressed implicitly as a network mask such as "/255.255.0.0".
func (addrStr *IPAddressString) GetNetworkPrefixLen() PrefixLen {
	return addrStr.getNetworkPrefixLen().copy()
}

// If this address is a valid address with an associated network prefix length then this returns that prefix length, otherwise returns nil.
// The prefix length may be expressed explicitly with the notation "/xx" where xx is a decimal value, or it may be expressed implicitly as a network mask such as "/255.255.0.0".
func (addrStr *IPAddressString) getNetworkPrefixLen() PrefixLen {
	addrStr = addrStr.init()
	if addrStr.IsValid() {
		return addrStr.addressProvider.getProviderNetworkPrefixLen()
	}
	return nil
}

// GetMask returns the mask, if any, that was provided with this address string.
func (addrStr *IPAddressString) GetMask() *IPAddress {
	addrStr = addrStr.init()
	if addrStr.IsValid() {
		return addrStr.addressProvider.getProviderMask()
	}
	return nil
}

// IsAllAddresses returns true if the string represents all IP addresses, such as the string "*".
// You can denote all IPv4 addresses with *.*, or all IPv6 addresses with *:*.
func (addrStr *IPAddressString) IsAllAddresses() bool {
	addrStr = addrStr.init()
	return addrStr.IsValid() && addrStr.addressProvider.isProvidingAllAddresses()
}

// IsEmpty returns true if the address string is empty (zero-length).
func (addrStr *IPAddressString) IsEmpty() bool {
	addrStr = addrStr.init()
	return addrStr.IsValid() && addrStr.addressProvider.isProvidingEmpty()
}

// IsIPv4 returns true if the address is IPv4.
func (addrStr *IPAddressString) IsIPv4() bool {
	addrStr = addrStr.init()
	return addrStr.IsValid() && addrStr.addressProvider.isProvidingIPv4()
}

// IsIPv6 returns true if the address is IPv6.
func (addrStr *IPAddressString) IsIPv6() bool {
	addrStr = addrStr.init()
	return addrStr.IsValid() && addrStr.addressProvider.isProvidingIPv6()
}

// IsMixedIPv6 returns whether the lower 4 bytes of the address string are represented as IPv4, if this address string represents an IPv6 address.
func (addrStr *IPAddressString) IsMixedIPv6() bool {
	addrStr = addrStr.init()
	return addrStr.IsIPv6() && addrStr.addressProvider.isProvidingMixedIPv6()
}

// IsBase85IPv6 returns whether this address string represents an IPv6 address, returns whether the string was base 85.
func (addrStr *IPAddressString) IsBase85IPv6() bool {
	return addrStr.IsIPv6() && addrStr.addressProvider.isProvidingBase85IPv6()
}

// IsIPv4Mapped returns true if the address is an IPv4-mapped IPv6 address.
func (addrStr *IPAddressString) IsIPv4Mapped() bool {
	return addrStr.IsIPv6() && ipv4MappedPrefix.PrefixEqual(addrStr)
}

// GetIPVersion returns the IP address version if this represents a valid IP address, otherwise it returns nil.
func (addrStr *IPAddressString) GetIPVersion() IPVersion {
	if addrStr.IsValid() {
		return addrStr.addressProvider.getProviderIPVersion()
	}
	return IndeterminateIPVersion
}

// IsLoopback returns whether this address string represents a loopback address,  such as "::1" or "127.0.0.1".
func (addrStr *IPAddressString) IsLoopback() bool {
	val := addrStr.GetAddress()
	return val != nil && val.IsLoopback()
}

// IsZero returns whether this string represents an IP address whose value is exactly zero.
func (addrStr *IPAddressString) IsZero() bool {
	value := addrStr.GetAddress()
	return value != nil && value.IsZero()
}

// String implements the [fmt.Stringer] interface,
// returning the original string used to create this IPAddressString (altered by strings.TrimSpace),
// or "<nil>" if the receiver is a nil pointer.
func (addrStr *IPAddressString) String() string {
	if addrStr == nil {
		return nilString()
	}
	return addrStr.str
}

// Format implements the [fmt.Formatter] interface.
// It accepts the verbs hat are applicable to strings,
// namely the verbs %s, %q, %x and %X.
func (addrStr IPAddressString) Format(state fmt.State, verb rune) {
	s := flagsFromState(state, verb)
	_, _ = state.Write([]byte(fmt.Sprintf(s, addrStr.str)))
}

// ToNormalizedString produces a normalized string for the address.
//
// For IPv4, it is the same as the canonical string.
//
// For IPv6, it differs from the canonical string.  Zero-segments are not compressed.
//
// If the address has a prefix length, it will be included in the string.
//
// If the original string is not a valid address string, the original string is used.
func (addrStr *IPAddressString) ToNormalizedString() string {
	addrStr = addrStr.init()
	if addrStr.IsValid() {
		if str, err := addrStr.toNormalizedString(addrStr.addressProvider); err == nil {
			return str
		}
	}
	return addrStr.String()
}

func (addrStr *IPAddressString) toNormalizedString(addressProvider ipAddressProvider) (result string, err addrerr.IncompatibleAddressError) {
	if addressProvider.isProvidingAllAddresses() {
		result = SegmentWildcardStr
	} else if addressProvider.isProvidingEmpty() {
		result = ""
	} else if addressProvider.isProvidingIPAddress() {
		var addr *IPAddress
		if addr, err = addressProvider.getProviderAddress(); err == nil {
			result = addr.ToNormalizedString()
		}
	}
	return
}

// IsValid returns whether this is a valid IP address string format.
// The accepted IP address formats are:
// an IPv4 address or subnet, an IPv6 address or subnet, the address representing all addresses of both versions, or an empty string.
// If this method returns false, and you want more details, call Validate and examine the error.
func (addrStr *IPAddressString) IsValid() bool {
	return addrStr.Validate() == nil
}

// GetAddress returns the IP address if this IPAddressString is a valid string representing an IP address or subnet.  Otherwise, it returns nil.
//
// Use ToAddress for an equivalent method that returns an error when the format is invalid.
//
// If you have a prefixed address and you wish to get only the host without the prefix, use GetHostAddress.
func (addrStr *IPAddressString) GetAddress() *IPAddress {
	addr, _ := addrStr.ToAddress()
	return addr
}

// ToAddress produces the IPAddress corresponding to this IPAddressString.
//
// If this object does not represent a specific IPAddress or a subnet, nil is returned.
//
// If the string used to construct this object is not a known format (empty string, address, or range of addresses) then this method returns an error.
//
// An equivalent method that does not return the error is GetAddress.
//
// If you have a prefixed address and you wish to get only the host rather than the address with the prefix, use ToHostAddress.
//
// The error can be addrerr.AddressStringError or addrerr.IncompatibleAddressError
func (addrStr *IPAddressString) ToAddress() (*IPAddress, addrerr.AddressError) {
	provider, err := addrStr.getAddressProvider()
	if err != nil {
		return nil, err
	}
	return provider.getProviderAddress()
}

// GetVersionedAddress is similar to ToVersionedAddress, but returns nil rather than an error when the address is invalid or does not match the supplied version.
func (addrStr *IPAddressString) GetVersionedAddress(version IPVersion) *IPAddress {
	addr, _ := addrStr.ToVersionedAddress(version)
	return addr
}

// ToVersionedAddress Produces the IPAddress of the specified address version corresponding to this IPAddressString.
//
// In most cases the string indicates the address version and calling ToAddress() is sufficient, with a few exceptions.
//
// When this object represents only a network prefix length,
// specifying the address version allows the conversion to take place to the associated mask for that prefix length.
//
// When this object represents all addresses, specifying the address version allows the conversion to take place
// to the associated representation of all IPv4 or all IPv6 addresses.
//
// When this object represents the empty string and that string is interpreted as a loopback or zero address, then it returns
// the corresponding address for the given version.
//
// When this object represents an ipv4 or ipv6 address, it returns that address if and only if that address matches the provided version.
//
// If the string used to construct this object is an invalid format,
// or a format that does not match the provided version, then an error is returned.
func (addrStr *IPAddressString) ToVersionedAddress(version IPVersion) (*IPAddress, addrerr.AddressError) {
	provider, err := addrStr.getAddressProvider()
	if err != nil {
		return nil, err
	}
	return provider.getVersionedAddress(version)
}

// GetHostAddress parses the address while ignoring the prefix length or mask.
// GetHostAddress returns nil for an invalid string.  If you wish to receive an error instead, use ToHostAddress.
func (addrStr *IPAddressString) GetHostAddress() *IPAddress {
	addr, _ := addrStr.ToHostAddress()
	return addr
}

// ToHostAddress parses the address while ignoring the prefix length or mask.
// The error can be addrerr.AddressStringError for invalid strings or addrerr.IncompatibleAddressError.
// GetHostAddress is similar but does not return errors.
// Standard address formats do not result in errors.
func (addrStr *IPAddressString) ToHostAddress() (*IPAddress, addrerr.AddressError) {
	provider, err := addrStr.getAddressProvider()
	if err != nil {
		return nil, err
	}
	return provider.getProviderHostAddress()
}

// TODO LATER getDivisionGrouping: allows for isSequential
//// IsSequential returns whether the addresses returned by this IPAddressString are sequential,
//// meaning that if any address has a numerical value that lies in between the numerical values of two addresses represented by this IPAddressString,
//// then that address is also represented by this IPAddressString.  In other words, the range of addresses is sequential.
////
//// When the IPAddressString is sequential, it can be represented exactly by the IPAddressSeqRange returned from GetSequentialRange.
//// In some cases, no IPAddress instance can be obtained from GetAddress or ToAddress,
//// but if the IPAddressString is sequential, you can obtain a IPAddressSeqRange to represent the IPAddressString instead.
//func (addrStr *IPAddressString) IsSequential() bool {
//	addrStr = addrStr.init()
//	return addrStr.IsValid() && addrStr.addressProvider.isSequential()
//}
//
//// Also restore this part of the godoc for GetSequentialRange below once IsSequential  is back:
////
//// The sequential range matches the same set of addresses as the address string or the address when IsSequential is true.
//// Otherwise, the range includes addresses not specified by the address string.
////

// GetSequentialRange returns the range of sequential addresses from the lowest address specified in this address string to the highest.
//
// Since not all IPAddressString instances describe a sequential series of addresses,
// this does not necessarily match the exact set of addresses specified by the string.
// For example, "1-2.3.4.1-2" produces the sequential range "1.3.4.1" to "2.3.4.2" that includes the address "1.255.255.2" not specified by the string.
//
// This method can also produce a range for a string for which no IPAddress instance can be created,
// those cases where IsValid returns true but ToAddress returns addrerr.IncompatibleAddressError and GetAddress returns nil.
// The range cannot be produced for the other cases where GetAddress returns nil
//
// This is similar to ToSequentialRange except that ToSequentialRange provides a descriptive error when nil is returned.
func (addrStr *IPAddressString) GetSequentialRange() (res *IPAddressSeqRange) {
	res, _ = addrStr.ToSequentialRange()
	return
}

// ToSequentialRange returns the range of sequential addresses from the lowest address specified in this address string to the highest.
//
// This is similar to GetSequentialRange except that this method provides a descriptive error when nil is returned. See GetSequentialRange for more details.
func (addrStr *IPAddressString) ToSequentialRange() (*IPAddressSeqRange, addrerr.AddressStringError) {
	provider, err := addrStr.getAddressProvider()
	if err != nil {
		return nil, err
	}
	return provider.getProviderSeqRange(), nil
}

// ValidateIPv4 validates that this string is a valid IPv4 address, returning nil, and if not, returns an error with a descriptive message indicating why it is not.
func (addrStr *IPAddressString) ValidateIPv4() addrerr.AddressStringError {
	return addrStr.ValidateVersion(IPv4)
}

// ValidateIPv6 validates that this string is a valid IPv6 address, returning nil, and if not, returns an error with a descriptive message indicating why it is not.
func (addrStr *IPAddressString) ValidateIPv6() addrerr.AddressStringError {
	return addrStr.ValidateVersion(IPv6)
}

func (addrStr *IPAddressString) getAddressProvider() (ipAddressProvider, addrerr.AddressStringError) {
	addrStr = addrStr.init()
	err := addrStr.Validate()
	return addrStr.addressProvider, err
}

func (addrStr *IPAddressString) validate(validationOptions addrstrparam.IPAddressStringParams) {
	addrStr.addressProvider, addrStr.validateError = validator.validateIPAddressStr(addrStr, validationOptions)
}

// Validate validates that this string is a valid IP address, returning nil, and if not, returns an error with a descriptive message indicating why it is not.
func (addrStr *IPAddressString) Validate() addrerr.AddressStringError {
	return addrStr.init().validateError
}

// ValidateVersion validates that this string is a valid IP address of the given version.
// If it is, it returns nil, otherwise it returns an error with a descriptive message indicating why it is not.
func (addrStr *IPAddressString) ValidateVersion(version IPVersion) addrerr.AddressStringError {
	addrStr = addrStr.init()
	err := addrStr.Validate()
	if err != nil {
		return err
	} else if version.IsIndeterminate() {
		return &addressStringError{addressError{str: addrStr.str, key: "ipaddress.error.ipVersionIndeterminate"}}
	} else {
		addrVersion := addrStr.addressProvider.getProviderIPVersion()
		if version.IsIPv4() {
			if addrVersion.IsIPv6() {
				return &addressStringError{addressError{str: addrStr.str, key: "ipaddress.error.address.is.ipv6"}}
			}
		} else if version.IsIPv6() {
			if addrVersion.IsIPv4() {
				return &addressStringError{addressError{str: addrStr.str, key: "ipaddress.error.address.is.ipv4"}}
			}
		}
	}
	return nil
}

// Compare compares this address string with another,
// returning a negative number, zero, or a positive number if this address string is less than, equal to, or greater than the other.
//
// All address strings are comparable.  If two address strings are invalid, their strings are compared.
// Otherwise, address strings are compared according to which type or version of string, and then within each type or version
// they are compared using the comparison rules for addresses.
func (addrStr *IPAddressString) Compare(other *IPAddressString) int {
	if addrStr == other {
		return 0
	} else if addrStr == nil {
		return -1
	} else if other == nil {
		return 1
	}
	addrStr = addrStr.init()
	other = other.init()
	if addrStr == other {
		return 0
	}
	if addrStr.IsValid() {
		if other.IsValid() {
			if res, err := addrStr.addressProvider.providerCompare(other.addressProvider); err == nil {
				return res
			}
			// one or the other is nil, either empty or IncompatibleAddressException
			return strings.Compare(addrStr.String(), other.String())
		}
		return 1
	} else if other.IsValid() {
		return -1
	}
	return strings.Compare(addrStr.String(), other.String())
}

// PrefixEqual is similar to Equal, but instead returns whether the prefix of this address string matches the same of the given address string,
// using the prefix length of this address string.  It returns whether the argument address string has the same address prefix values as this.
//
// In other words, it determines if the other address has the same prefix subnet using the prefix length of this address.
//
// If an address has no prefix length, the whole address is compared.
//
// If this address string or the given address string is invalid, it returns false.
func (addrStr *IPAddressString) PrefixEqual(other *IPAddressString) bool {
	// getting the prefix
	addrStr = addrStr.init()
	other = other.init()
	if other == addrStr {
		return true
	} else if !addrStr.IsValid() {
		return false
	} else if other.IsValid() {
		directResult := addrStr.addressProvider.prefixEqualsProvider(other.addressProvider)
		if directResult.isSet {
			return directResult.val
		}
		thisAddress := addrStr.GetAddress()
		if thisAddress != nil {
			otherAddress := other.GetAddress()
			if otherAddress != nil {
				return thisAddress.prefixEquals(otherAddress)
			}
		}
		// one or both addresses are nil, so there is no prefix to speak of
	}
	return false
}

// PrefixContains is similar to PrefixEqual, but instead returns whether the prefix of this address contains the same of the given address,
// using the prefix length of this address.  It returns whether the argument address string prefix values of that length are also prefix values in this address string.
//
// In other words, determines if the other address is in one of the same prefix subnets using the prefix length of this address.
//
// If an address has no prefix length, the whole address is used as the prefix.
//
// If this address string or the given address string is invalid, it returns false.
func (addrStr *IPAddressString) PrefixContains(other *IPAddressString) bool {
	addrStr = addrStr.init()
	other = other.init()
	if other == addrStr {
		return true
	} else if !addrStr.IsValid() {
		return false
	} else if other.IsValid() {
		directResult := addrStr.addressProvider.prefixContainsProvider(other.addressProvider)
		if directResult.isSet {
			return directResult.val
		}
		thisAddress := addrStr.GetAddress()
		if thisAddress != nil {
			otherAddress := other.GetAddress()
			if otherAddress != nil {
				return thisAddress.prefixContains(otherAddress)
			}
		}
		// one or both addresses are nil, so there is no prefix to speak of
	}
	return false
}

// Contains returns whether the address or subnet identified by this address string contains the address or subnet identified by the given string.
// If this address string or the given address string is invalid then Contains returns false.
//
// With standard subnet and address formats, using the strings alone are enough for this method to return the result.
// Only with non-standard strings will address objects be constructed to determine the result.
func (addrStr *IPAddressString) Contains(other *IPAddressString) bool {
	addrStr = addrStr.init()
	other = other.init()
	if addrStr.IsValid() {
		if other == addrStr {
			return true
		}
		if other.IsValid() {
			// note the quick result also handles the case of "all addresses"
			directResult := addrStr.addressProvider.containsProvider(other.addressProvider)
			if directResult.isSet {
				return directResult.val
			}
			// defer to the constructed addresses
			addr := addrStr.GetAddress()
			if addr != nil {
				otherAddress := other.GetAddress()
				if otherAddress != nil {
					return addr.Contains(otherAddress)
				}
			}
		}
	}
	return false
}

// Equal compares two IP address strings for equality.
// Two IPAddressString objects are equal if they represent the same set of addresses.
// Whether one or the other has an associated network prefix length is not considered.
//
// If an IPAddressString is invalid, it is equal to another address only if the other address was constructed from the same string.
func (addrStr *IPAddressString) Equal(other *IPAddressString) bool {
	if addrStr == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	addrStr = addrStr.init()
	other = other.init()
	if other == addrStr {
		return true
	}
	// if they have the same string, they must be the same,
	// but the converse is not true, if they have different strings, they can
	// still be the same because IPv6 addresses have many representations
	// and additional things like leading zeros can have an effect for IPv4

	// Also note that we do not call equals() on the validation options, this is intended as an optimization,
	// and probably better to avoid going through all the validation objects here
	stringsMatch := addrStr.String() == other.String()
	if stringsMatch && addrStr.GetValidationOptions() == other.GetValidationOptions() {
		return true
	}
	if addrStr.IsValid() {
		if other.IsValid() {
			directResult := addrStr.addressProvider.parsedEquals(other.addressProvider)
			if directResult.isSet {
				return directResult.val
			}
			// When a value provider produces no value, equality and comparison are based on the enum ipType
			var err addrerr.AddressError
			addrProvider, err := addrStr.getAddressProvider()
			if err != nil {
				return stringsMatch
			}
			equals, err := addrProvider.providerEquals(other.addressProvider)
			if err != nil {
				return stringsMatch
			}
			return equals
		}
	} else if !other.IsValid() {
		return stringsMatch // Two invalid addresses are not equal unless strings match, regardless of validation options
	}
	return false
}

// AdjustPrefixLen increases or decreases the prefix length by the given increment.
//
// If the address string has prefix length 0 and represents all addresses of the same version,
// and the prefix length is being decreased, then the address representing all addresses of any version is returned.
//
// When there is an associated address value and the prefix length is increased, the bits moved within the prefix become zero,
// and if prefix length is extended beyond the segment series boundary, it is removed.
// When there is an associated address value and the prefix length is decreased, the bits moved outside the prefix become zero.
//
// If the address string represents a prefix block, then the result will also represent a prefix block.
func (addrStr *IPAddressString) AdjustPrefixLen(adjustment BitCount) (*IPAddressString, addrerr.IncompatibleAddressError) {
	address := addrStr.GetAddress()
	if address == nil {
		return nil, nil
	}
	if adjustment == 0 && addrStr.IsPrefixed() {
		return addrStr, nil
	}
	prefix := address.getNetworkPrefixLen()
	isPrefBlock := address.IsPrefixBlock()
	var addr *IPAddress
	var err addrerr.IncompatibleAddressError
	if adjustment < 0 && isPrefBlock {
		if prefix != nil && prefix.bitCount()+adjustment < 0 {
			return NewIPAddressStringParams(SegmentWildcardStr, addrStr.GetValidationOptions()), nil
		}
		addr, err = address.AdjustPrefixLenZeroed(adjustment)
		if err != nil {
			return nil, err
		}
		addr = addr.ToPrefixBlock()
	} else {
		addr, err = address.AdjustPrefixLenZeroed(adjustment)
		if err != nil {
			return nil, err
		}
	}
	return addr.ToAddressString(), nil
}

// Wrap wraps this address string, returning a WrappedIPAddressString as an implementation of ExtendedIdentifierString,
// which can be used to write code that works with different host identifier types polymorphically,  including IPAddressString, MACAddressString, and HostName.
func (addrStr *IPAddressString) Wrap() ExtendedIdentifierString {
	return WrappedIPAddressString{addrStr}
}

// ValidatePrefixLenStr validates that the string represents a valid prefix length, such as "24".
// The string should not include a beginning '/' character.
// If invalid, it returns an error with an appropriate message.
// You can specify the IP version or IndeterminateIPVersion if unknown.
// An error is returned if the format is invalid.
func ValidatePrefixLenStr(str string, version IPVersion) (prefixLen PrefixLen, err addrerr.AddressStringError) {
	return validator.validatePrefixLenStr(str, version)
}
