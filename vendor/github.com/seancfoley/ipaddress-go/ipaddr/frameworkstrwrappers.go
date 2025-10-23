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

// ExtendedIdentifierString is a common interface for strings that identify hosts, namely [IPAddressString], [MACAddressString], and [HostName].
type ExtendedIdentifierString interface {
	HostIdentifierString

	// GetAddress returns the identified address or nil if none.
	GetAddress() AddressType

	// ToAddress returns the identified address or an error.
	ToAddress() (AddressType, error)

	// Unwrap returns the wrapped IPAddressString, MACAddressString or HostName as an interface, HostIdentifierString.
	Unwrap() HostIdentifierString
}

// WrappedIPAddressString wraps an IPAddressString to get an ExtendedIdentifierString, an extended polymorphic type.
type WrappedIPAddressString struct {
	*IPAddressString
}

// Unwrap returns the wrapped IPAddressString as an interface, HostIdentifierString.
func (str WrappedIPAddressString) Unwrap() HostIdentifierString {
	res := str.IPAddressString
	if res == nil {
		return nil
	}
	return res
}

// ToAddress returns the identified address or an error.
func (str WrappedIPAddressString) ToAddress() (AddressType, error) {
	addr, err := str.IPAddressString.ToAddress()
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// GetAddress returns the identified address or nil if none.
func (str WrappedIPAddressString) GetAddress() AddressType {
	if addr := str.IPAddressString.GetAddress(); addr != nil {
		return addr
	}
	return nil
}

// WrappedMACAddressString wraps a MACAddressString to get an ExtendedIdentifierString.
type WrappedMACAddressString struct {
	*MACAddressString
}

// Unwrap returns the wrapped MACAddressString as an interface, HostIdentifierString.
func (str WrappedMACAddressString) Unwrap() HostIdentifierString {
	res := str.MACAddressString
	if res == nil {
		return nil
	}
	return res
}

// ToAddress returns the identified address or an error.
func (str WrappedMACAddressString) ToAddress() (AddressType, error) {
	addr, err := str.MACAddressString.ToAddress()
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// GetAddress returns the identified address or nil if none.
func (str WrappedMACAddressString) GetAddress() AddressType {
	if addr := str.MACAddressString.GetAddress(); addr != nil {
		return addr
	}
	return nil
}

// WrappedHostName wraps a HostName to get an ExtendedIdentifierString.
type WrappedHostName struct {
	*HostName
}

// Unwrap returns the wrapped HostName as an interface, HostIdentifierString.
func (host WrappedHostName) Unwrap() HostIdentifierString {
	res := host.HostName
	if res == nil {
		return nil
	}
	return res
}

// ToAddress returns the identified address or an error.
func (host WrappedHostName) ToAddress() (AddressType, error) {
	addr, err := host.HostName.ToAddress()
	if err != nil {
		return nil, err
	}
	return addr, nil
}

// GetAddress returns the identified address or nil if none.
func (host WrappedHostName) GetAddress() AddressType {
	if addr := host.HostName.GetAddress(); addr != nil {
		return addr
	}
	return nil
}

var (
	_, _, _ ExtendedIdentifierString = WrappedIPAddressString{}, WrappedMACAddressString{}, WrappedHostName{}
)
