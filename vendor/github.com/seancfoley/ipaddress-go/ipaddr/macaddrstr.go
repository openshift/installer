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

package ipaddr

import (
	"fmt"
	"strings"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstrparam"
)

var defaultMACAddrParameters = new(addrstrparam.MACAddressStringParamsBuilder).ToParams()

// NewMACAddressStringParams constructs a MACAddressString that will parse the given string according to the given parameters.
func NewMACAddressStringParams(str string, params addrstrparam.MACAddressStringParams) *MACAddressString {
	var p addrstrparam.MACAddressStringParams
	if params == nil {
		p = defaultMACAddrParameters
	} else {
		p = addrstrparam.CopyMACAddressStringParams(params)
	}
	return parseMACAddressString(str, p)
}

// NewMACAddressString constructs a MACAddressString that will parse the given string according to the default parameters.
func NewMACAddressString(str string) *MACAddressString {
	return parseMACAddressString(str, defaultMACAddrParameters)
}

func newMACAddressStringFromAddr(str string, addr *MACAddress) *MACAddressString {
	return &MACAddressString{
		str:             str,
		addressProvider: wrappedMACAddressProvider{addr},
	}
}

func parseMACAddressString(str string, params addrstrparam.MACAddressStringParams) *MACAddressString {
	str = strings.TrimSpace(str)
	res := &MACAddressString{str: str}
	res.validate(params)
	return res
}

var zeroMACAddressString = NewMACAddressString("")

// MACAddressString parses the string representation of a MAC address.  Such a string can represent just a single address or a collection of addresses like "1:*:1-3:1-4:5:6".
//
// This supports a wide range of address formats and provides specific error messages, and allows specific configuration.
//
// You can control all the supported formats using MACAddressStringParamsBuilder to build a parameters instance of  MACAddressStringParams.
// When not using the constructor that takes a MACAddressStringParams, a default instance of MACAddressStringParams is used that is generally permissive.
//
// # Supported Formats
//
// Ranges are supported:
//
//   - wildcards '*' and ranges '-' (for example "1:*:1-3:1-4:5:6"), useful for working with MAC address collections
//   - SQL wildcards '%" and "_", although '%' is considered an SQL wildcard only when it is not considered an IPv6 zone indicator
//
// The different methods of representing MAC addresses are supported:
//
//   - 6 or 8 bytes in hex representation like "aa:bb:cc:dd:ee:ff"
//   - The same but with a hyphen separator like "aa-bb-cc-dd-ee-ff" (the range separator in this case becomes '/')
//   - The same but with space separator like "aa bb cc dd ee ff"
//   - The dotted representation, 4 sets of 12 bits in hex representation like "aaa.bbb.ccc.ddd"
//   - The 12 or 16 hex representation with no separators like "aabbccddeeff"
//
// All of the above range variations also work for each of these ways of representing MAC addresses.
//
// Some additional formats:
//
//   - null or empty strings representing an unspecified address
//   - the single wildcard address "*" which represents all MAC addresses
//
// Usage
// Once you have constructed a MACAddressString object, you can convert it to a [MACAddress] object with GetAddress or ToAddress.
//
// For empty addresses, both ToAddress and GetAddress return nil.  For invalid addresses, GetAddress and ToAddress return nil, with ToAddress also returning an error.
//
// This type is concurrency-safe.  In fact, MACAddressString objects are immutable.
// A MACAddressString object represents a single MAC address representation that cannot be changed after construction.
// Some derived state is created upon demand and cached, such as the derived [MACAddress] instances.
type MACAddressString struct {
	str             string
	addressProvider macAddressProvider
	validateError   addrerr.AddressStringError
}

func (addrStr *MACAddressString) init() *MACAddressString {
	if addrStr.addressProvider == nil && addrStr.validateError == nil {
		return zeroMACAddressString
	}
	return addrStr
}

// GetValidationOptions returns the validation options supplied when constructing this address string,
// or the default options if no options were supplied.  It returns nil if no parameters were used to construct the address.
func (addrStr *MACAddressString) GetValidationOptions() addrstrparam.MACAddressStringParams {
	provider, _ := addrStr.getAddressProvider()
	if provider != nil {
		return provider.getParameters()
	}
	return nil
}

// String implements the [fmt.Stringer] interface,
// returning the original string used to create this MACAddressString (altered by strings.TrimSpace),
// or "<nil>" if the receiver is a nil pointer.
func (addrStr *MACAddressString) String() string {
	if addrStr == nil {
		return nilString()
	}
	return addrStr.str
}

// Format implements the [fmt.Formatter] interface.
// It accepts the verbs hat are applicable to strings,
// namely the verbs %s, %q, %x and %X.
func (addrStr MACAddressString) Format(state fmt.State, verb rune) {
	s := flagsFromState(state, verb)
	_, _ = state.Write([]byte(fmt.Sprintf(s, addrStr.str)))
}

// ToNormalizedString produces a normalized string for the address.
//
// For MAC, it differs from the canonical string.  It uses the most common representation of MAC addresses: "xx:xx:xx:xx:xx:xx".  An example is "01:23:45:67:89:ab".
// For range segments, '-' is used: "11:22:33-44:55:66".
//
// If the original string is not a valid address string, the original string is used.
func (addrStr *MACAddressString) ToNormalizedString() string {
	addr := addrStr.GetAddress()
	if addr != nil {
		return addr.toNormalizedString()
	}
	return addrStr.String()
}

// GetAddress returns the MAC address if this MACAddressString is a valid string representing a MAC address or address collection.  Otherwise, it returns nil.
//
// Use ToAddress for an equivalent method that returns an error when the format is invalid.
func (addrStr *MACAddressString) GetAddress() *MACAddress {
	addr, _ := addrStr.ToAddress()
	return addr
}

// ToAddress produces the MACAddress corresponding to this MACAddressString.
//
// If this object does not represent a specific MACAddress or address collection, nil is returned.
//
// If the string used to construct this object is not a known format (empty string, address, or range of addresses) then this method returns an error.
//
// An equivalent method that does not return the error is GetAddress.
//
// The error can be addrerr.AddressStringError for an invalid string, or addrerr.IncompatibleAddressError for non-standard strings that cannot be converted to MACAddress.
func (addrStr *MACAddressString) ToAddress() (*MACAddress, addrerr.AddressError) {
	provider, err := addrStr.getAddressProvider()
	if err != nil {
		return nil, err
	}
	return provider.getAddress()
}

// IsPrefixed returns whether this address has an associated prefix length,
// which for MAC means that the string represents the set of all addresses with the same prefix.
func (addrStr *MACAddressString) IsPrefixed() bool {
	return addrStr.getPrefixLen() != nil
}

// GetPrefixLen returns the prefix length if this address is a prefixed address, otherwise it returns nil.
//
// For MAC addresses, the prefix is initially inferred from the range, so "1:2:3:*:*:*" has a prefix length of 24.
// Addresses derived from the original may retain the original prefix length regardless of their range.
func (addrStr *MACAddressString) GetPrefixLen() PrefixLen {
	return addrStr.getPrefixLen().copy()
}

func (addrStr *MACAddressString) getPrefixLen() PrefixLen {
	addr := addrStr.GetAddress()
	if addr != nil {
		return addr.getPrefixLen()
	}
	return nil
}

// IsFullRange returns whether the address represents the set of all valid MAC addresses for its address length
func (addrStr *MACAddressString) IsFullRange() bool {
	addr := addrStr.GetAddress()
	return addr != nil && addr.IsFullRange()
}

// IsEmpty returns true if the address is empty (zero-length).
func (addrStr *MACAddressString) IsEmpty() bool {
	addr, err := addrStr.ToAddress()
	return err == nil && addr == nil
}

// IsZero returns whether this string represents a MAC address whose value is exactly zero.
func (addrStr *MACAddressString) IsZero() bool {
	addr := addrStr.GetAddress()
	return addr != nil && addr.IsZero()
}

// IsValid returns whether this is a valid MAC address string format.
// The accepted MAC address formats are:
// a MAC address or address collection, the address representing all MAC addresses, or an empty string.
// If this method returns false, and you want more details, call Validate and examine the error.
func (addrStr *MACAddressString) IsValid() bool {
	return addrStr.Validate() == nil
}

func (addrStr *MACAddressString) getAddressProvider() (macAddressProvider, addrerr.AddressStringError) {
	addrStr = addrStr.init()
	err := addrStr.Validate()
	return addrStr.addressProvider, err
}

func (addrStr *MACAddressString) validate(validationOptions addrstrparam.MACAddressStringParams) {
	addrStr.addressProvider, addrStr.validateError = validator.validateMACAddressStr(addrStr, validationOptions)
}

// Validate validates that this string is a valid address, and if not, throws an exception with a descriptive message indicating why it is not.
func (addrStr *MACAddressString) Validate() addrerr.AddressStringError {
	return addrStr.init().validateError
}

// Compare compares this address string with another,
// returning a negative number, zero, or a positive number if this address string is less than, equal to, or greater than the other.
//
// All address strings are comparable.  If two address strings are invalid, their strings are compared.
// Two valid address trings are compared using the comparison rules for their respective addresses.
func (addrStr *MACAddressString) Compare(other *MACAddressString) int {
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
			addr := addrStr.GetAddress()
			if addr != nil {
				otherAddr := other.GetAddress()
				if otherAddr != nil {
					return addr.Compare(otherAddr)
				}
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

// Equal returns whether this MACAddressString is equal to the given one.
// Two MACAddressString objects are equal if they represent the same set of addresses.
//
// If a MACAddressString is invalid, it is equal to another address only if the other address was constructed from the same string.
func (addrStr *MACAddressString) Equal(other *MACAddressString) bool {
	if addrStr == nil {
		return other == nil
	} else if other == nil {
		return false
	}
	addrStr = addrStr.init()
	other = other.init()
	if addrStr == other {
		return true
	}

	//if they have the same string, they must be the same,
	//but the converse is not true, if they have different strings, they can still be the same

	// Also note that we do not call equals() on the validation options, this is intended as an optimization,
	// and probably better to avoid going through all the validation objects here
	stringsMatch := addrStr.String() == other.String()
	if stringsMatch && addrStr.GetValidationOptions() == other.GetValidationOptions() {
		return true
	}
	if addrStr.IsValid() {
		if other.IsValid() {
			value := addrStr.GetAddress()
			if value != nil {
				otherValue := other.GetAddress()
				if otherValue != nil {
					return value.equals(otherValue)
				} else {
					return false
				}
			} else if other.GetAddress() != nil {
				return false
			}
			// both are nil, either empty or addrerr.IncompatibleAddressError
			return stringsMatch
		}
	} else if !other.IsValid() { // both are invalid
		return stringsMatch // Two invalid addresses are not equal unless strings match, regardless of validation options
	}
	return false
}

// Wrap wraps this address string, returning a WrappedMACAddressString as an implementation of ExtendedIdentifierString,
// which can be used to write code that works with different host identifier types polymorphically,  including IPAddressString, MACAddressString, and HostName.
func (addrStr *MACAddressString) Wrap() ExtendedIdentifierString {
	return WrappedMACAddressString{addrStr}
}
