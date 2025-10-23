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

import (
	"sync"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstrparam"
)

type macAddressProvider interface {
	getAddress() (*MACAddress, addrerr.IncompatibleAddressError)

	// If the address was created by parsing, this provides the parameters used when creating the address,
	// otherwise nil
	getParameters() addrstrparam.MACAddressStringParams
}

type macAddressNullProvider struct {
	validationOptions addrstrparam.MACAddressStringParams
}

var invalidMACProvider = macAddressEmptyProvider{macAddressNullProvider{defaultMACAddrParameters}}

func (provider macAddressNullProvider) getParameters() addrstrparam.MACAddressStringParams {
	return provider.validationOptions
}

func (provider macAddressNullProvider) getAddress() (*MACAddress, addrerr.IncompatibleAddressError) {
	return nil, nil
}

type macAddressEmptyProvider struct {
	macAddressNullProvider
}

var defaultMACAddressEmptyProvider = macAddressEmptyProvider{macAddressNullProvider{defaultMACAddrParameters}}

type macAddressAllProvider struct {
	validationOptions addrstrparam.MACAddressStringParams
	address           *MACAddress
	creationLock      *sync.Mutex
}

func (provider *macAddressAllProvider) getParameters() addrstrparam.MACAddressStringParams {
	return provider.validationOptions
}

func (provider *macAddressAllProvider) getAddress() (*MACAddress, addrerr.IncompatibleAddressError) {
	addr := provider.address
	if addr == nil {
		provider.creationLock.Lock()
		addr = provider.address
		if addr == nil {
			validationOptions := provider.validationOptions
			size := validationOptions.GetPreferredLen()
			creator := macType.getNetwork().getAddressCreator()
			var segCount int
			if size == addrstrparam.EUI64Len {
				segCount = ExtendedUniqueIdentifier64SegmentCount
			} else {
				segCount = MediaAccessControlSegmentCount
			}
			allRangeSegment := creator.createRangeSegment(0, MACMaxValuePerSegment)
			segments := make([]*AddressDivision, segCount)
			for i := range segments {
				segments[i] = allRangeSegment
			}
			section := creator.createSectionInternal(segments, true)
			addr = creator.createAddressInternal(section.ToSectionBase(), nil).ToMAC()
		}
		provider.creationLock.Unlock()
	}
	return addr, nil
}

var macAddressDefaultAllProvider = &macAddressAllProvider{validationOptions: defaultMACAddrParameters, creationLock: &sync.Mutex{}}

type wrappedMACAddressProvider struct {
	address *MACAddress
}

func (provider wrappedMACAddressProvider) getParameters() addrstrparam.MACAddressStringParams {
	return nil
}

func (provider wrappedMACAddressProvider) getAddress() (*MACAddress, addrerr.IncompatibleAddressError) {
	return provider.address, nil
}

var (
	_, _, _ macAddressProvider = macAddressEmptyProvider{},
		&macAddressAllProvider{},
		&wrappedMACAddressProvider{}
)
