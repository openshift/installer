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
	"unsafe"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstrparam"
)

// All IP address strings corresponds to exactly one of these types.
// In cases where there is no corresponding default IPAddress value (invalidType, allType, and possibly emptyType), these types can be used for comparison.
// emptyType means a zero-length string (useful for validation, we can set validation to allow empty strings) that has no corresponding IPAddress value (validation options allow you to map empty to the loopback)
// invalidType means it is known that it is not any of the other allowed types (validation options can restrict the allowed types)
// allType means it is wildcard(s) with no separators, like "*", which represents all addresses, whether IPv4, IPv6 or other, and thus has no corresponding IPAddress value
// These constants are ordered by address space size, from smallest to largest, and the ordering affects comparisons
type ipType int

func fromVersion(version IPVersion) ipType {
	switch version {
	case IPv4:
		return ipv4AddrType
	case IPv6:
		return ipv6AddrType
	default:
	}
	return uninitializedType
}

func (t ipType) isUnknown() bool {
	return t == uninitializedType
}

const (
	uninitializedType ipType = iota
	invalidType
	emptyType
	ipv4AddrType
	ipv6AddrType
	allType
)

type ipAddressProvider interface {
	getType() ipType

	getProviderHostAddress() (*IPAddress, addrerr.IncompatibleAddressError)

	getProviderAddress() (*IPAddress, addrerr.IncompatibleAddressError)

	getVersionedAddress(version IPVersion) (*IPAddress, addrerr.IncompatibleAddressError)

	isSequential() bool

	getProviderSeqRange() *SequentialRange[*IPAddress]

	getProviderMask() *IPAddress

	// TODO LATER getDivisionGrouping
	//default IPAddressDivisionSeries getDivisionGrouping() throwsaddrerr.IncompatibleAddressError {
	//	return getProviderAddress();
	//}

	providerCompare(ipAddressProvider) (int, addrerr.IncompatibleAddressError)

	providerEquals(ipAddressProvider) (bool, addrerr.IncompatibleAddressError)

	getProviderIPVersion() IPVersion

	isProvidingIPAddress() bool

	isProvidingIPv4() bool

	isProvidingIPv6() bool

	// providing **all** addresses of any IP version, ie "*", not "*.*" or "*:*"
	isProvidingAllAddresses() bool

	isProvidingEmpty() bool

	isProvidingMixedIPv6() bool

	isProvidingBase85IPv6() bool

	getProviderNetworkPrefixLen() PrefixLen

	isInvalid() bool

	// If the address was created by parsing, this provides the parameters used when creating the address,
	// otherwise nil
	getParameters() addrstrparam.IPAddressStringParams

	// containsProvider is an optimized contains that does not need to create address objects to return an answer.
	// Unconventional addresses may require that the address objects are created, in such cases nil is returned.
	//
	// Addresses constructed from canonical or normalized representations with no wildcards will not return null.
	containsProvider(ipAddressProvider) boolSetting

	// prefixEqualsProvider is an optimized prefix comparison that does not need to create addresses to return an answer.
	//
	// Unconventional addresses may require the address objects, in such cases null is returned.
	prefixEqualsProvider(ipAddressProvider) boolSetting

	// prefixContainsProvider is an optimized prefix comparison that does not need to create addresses to return an answer.
	//
	// Unconventional addresses may require the address objects, in such cases null is returned.
	prefixContainsProvider(ipAddressProvider) boolSetting

	// parsedEquals is an optimized equality comparison that does not need to create addresses to return an answer.
	//
	// Unconventional addresses may require the address objects, in such cases null is returned.
	parsedEquals(ipAddressProvider) boolSetting
}

type ipAddrProvider struct{}

func (p *ipAddrProvider) getType() ipType {
	return uninitializedType
}

func (p *ipAddrProvider) isSequential() bool {
	return false
}

func (p *ipAddrProvider) getProviderHostAddress() (*IPAddress, addrerr.IncompatibleAddressError) {
	return nil, nil
}

func (p *ipAddrProvider) getProviderAddress() (*IPAddress, addrerr.IncompatibleAddressError) {
	return nil, nil
}

func (p *ipAddrProvider) getProviderSeqRange() *SequentialRange[*IPAddress] {
	return nil
}

func (p *ipAddrProvider) getVersionedAddress(_ IPVersion) (*IPAddress, addrerr.IncompatibleAddressError) {
	return nil, nil
}

func (p *ipAddrProvider) getProviderMask() *IPAddress {
	return nil
}

func (p *ipAddrProvider) getProviderIPVersion() IPVersion {
	return IndeterminateIPVersion
}

func (p *ipAddrProvider) isProvidingIPAddress() bool {
	return false
}

func (p *ipAddrProvider) isProvidingIPv4() bool {
	return false
}

func (p *ipAddrProvider) isProvidingIPv6() bool {
	return false
}

func (p *ipAddrProvider) isProvidingAllAddresses() bool {
	return false
}

func (p *ipAddrProvider) isProvidingEmpty() bool {
	return false
}

func (p *ipAddrProvider) isInvalid() bool {
	return false
}

func (p *ipAddrProvider) isProvidingMixedIPv6() bool {
	return false
}

func (p *ipAddrProvider) isProvidingBase85IPv6() bool {
	return false
}

func (p *ipAddrProvider) getProviderNetworkPrefixLen() PrefixLen {
	return nil
}

func (p *ipAddrProvider) getParameters() addrstrparam.IPAddressStringParams {
	return nil
}

func (p *ipAddrProvider) containsProvider(ipAddressProvider) boolSetting {
	return boolSetting{}
}

func (p *ipAddrProvider) prefixEqualsProvider(ipAddressProvider) boolSetting {
	return boolSetting{}
}

func (p *ipAddrProvider) prefixContainsProvider(ipAddressProvider) boolSetting {
	return boolSetting{}
}

func (p *ipAddrProvider) parsedEquals(ipAddressProvider) boolSetting {
	return boolSetting{}
}

func providerCompare(p, other ipAddressProvider) (res int, err addrerr.IncompatibleAddressError) {
	if p == other {
		return
	}
	value, err := p.getProviderAddress()
	if err != nil {
		return
	}
	if value != nil {
		var otherValue *IPAddress
		otherValue, err = other.getProviderAddress()
		if err != nil {
			return
		}
		if otherValue != nil {
			res = value.Compare(otherValue)
			return
		}
	}
	var thisType, otherType = p.getType(), other.getType()
	res = int(thisType - otherType)
	return
}

/**
* When a value provider produces no value, equality and comparison are based on the enum ipType,
* which can by null.
* @param o
* @return
 */
func providerEquals(p, other ipAddressProvider) (res bool, err addrerr.IncompatibleAddressError) {
	if p == other {
		res = true
		return
	}
	value, err := p.getProviderAddress()
	if err != nil {
		return
	}
	if value != nil {
		var otherValue *IPAddress
		otherValue, err = other.getProviderAddress()
		if err != nil {
			return
		}
		if otherValue != nil {
			res = value.Equal(otherValue)
			return
		} else {
			return // returns false
		}
	}
	res = p.getType() == other.getType()
	return
}

type nullProvider struct {
	ipAddrProvider

	ipType                ipType
	isInvalidVal, isEmpty bool
	params                addrstrparam.IPAddressStringParams
}

func (p *nullProvider) isInvalid() bool {
	return p.isInvalidVal
}

func (p *nullProvider) isProvidingEmpty() bool {
	return p.isEmpty
}

func (p *nullProvider) getType() ipType {
	return p.ipType
}

func (p *nullProvider) providerCompare(other ipAddressProvider) (int, addrerr.IncompatibleAddressError) {
	return providerCompare(p, other)
}

func (p *nullProvider) providerEquals(other ipAddressProvider) (bool, addrerr.IncompatibleAddressError) {
	return providerEquals(p, other)
}

var (
	invalidProvider = &nullProvider{isInvalidVal: true, ipType: invalidType}
	emptyProvider   = &nullProvider{isEmpty: true, ipType: emptyType}
)

// Wraps an IPAddress for IPAddressString in the cases where no parsing is provided, the address exists already
func getProviderFor(address, hostAddress *IPAddress) ipAddressProvider {
	return &cachedAddressProvider{addresses: &addressResult{address: address, hostAddress: hostAddress}}
}

type addressResult struct {
	address, hostAddress *IPAddress

	// addrErr applies to address, hostErr to hostAddress
	addrErr, hostErr addrerr.IncompatibleAddressError

	// only used when no address can be obtained
	rng *SequentialRange[*IPAddress]
}

type cachedAddressProvider struct {
	ipAddrProvider

	// addressCreator creates two addresses, the host address and address with prefix/mask, at the same time
	addressCreator func() (address, hostAddress *IPAddress, addrErr, hosterr addrerr.IncompatibleAddressError)

	addresses *addressResult
}

func (cached *cachedAddressProvider) providerCompare(other ipAddressProvider) (int, addrerr.IncompatibleAddressError) {
	return providerCompare(cached, other)
}

func (cached *cachedAddressProvider) providerEquals(other ipAddressProvider) (bool, addrerr.IncompatibleAddressError) {
	return providerEquals(cached, other)
}

func (cached *cachedAddressProvider) isProvidingIPAddress() bool {
	return true
}

func (cached *cachedAddressProvider) getVersionedAddress(version IPVersion) (*IPAddress, addrerr.IncompatibleAddressError) {
	thisVersion := cached.getProviderIPVersion()
	if version != thisVersion {
		return nil, nil
	}
	return cached.getProviderAddress()
}

func (cached *cachedAddressProvider) getProviderSeqRange() *SequentialRange[*IPAddress] {
	addr, _ := cached.getProviderAddress()
	if addr != nil {
		return addr.ToSequentialRange()
	}
	return nil
}

func (cached *cachedAddressProvider) isSequential() bool {
	addr, _ := cached.getProviderAddress()
	if addr != nil {
		return addr.IsSequential()
	}
	return false
}

func (cached *cachedAddressProvider) getProviderHostAddress() (res *IPAddress, err addrerr.IncompatibleAddressError) {
	_, res, _, err = cached.getCachedAddresses()
	return
}

func (cached *cachedAddressProvider) getProviderAddress() (res *IPAddress, err addrerr.IncompatibleAddressError) {
	res, _, err, _ = cached.getCachedAddresses()
	return
}

func (cached *cachedAddressProvider) getCachedAddresses() (address, hostAddress *IPAddress, addrErr, hostErr addrerr.IncompatibleAddressError) {
	addrs := (*addressResult)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cached.addresses))))
	if addrs == nil {
		if cached.addressCreator != nil {
			address, hostAddress, addrErr, hostErr = cached.addressCreator()
			addresses := &addressResult{
				address:     address,
				hostAddress: hostAddress,
				addrErr:     addrErr,
				hostErr:     hostErr,
			}
			dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&cached.addresses))
			atomicStorePointer(dataLoc, unsafe.Pointer(addresses))
		}
	} else {
		address, hostAddress, addrErr, hostErr = addrs.address, addrs.hostAddress, addrs.addrErr, addrs.hostErr
	}
	return
}

func (cached *cachedAddressProvider) getProviderNetworkPrefixLen() (p PrefixLen) {
	if addr, _ := cached.getProviderAddress(); addr != nil {
		p = addr.getNetworkPrefixLen()
	}
	return
}

func (cached *cachedAddressProvider) getProviderIPVersion() IPVersion {
	if addr, _ := cached.getProviderAddress(); addr != nil {
		return addr.getIPVersion()
	}
	return IndeterminateIPVersion
}

func (cached *cachedAddressProvider) getType() ipType {
	return fromVersion(cached.getProviderIPVersion())
}

func (cached *cachedAddressProvider) isProvidingIPv4() bool {
	addr, _ := cached.getProviderAddress()
	return addr.IsIPv4()
}

func (cached *cachedAddressProvider) isProvidingIPv6() bool {
	addr, _ := cached.getProviderAddress()
	return addr.IsIPv6()
}

type versionedAddressCreator struct {
	cachedAddressProvider

	adjustedVersion IPVersion

	versionedAddressCreatorFunc func(IPVersion) (*IPAddress, addrerr.IncompatibleAddressError)

	versionedValues [2]*IPAddress

	parameters addrstrparam.IPAddressStringParams
}

func (versioned *versionedAddressCreator) getParameters() addrstrparam.IPAddressStringParams {
	return versioned.parameters
}

func (versioned *versionedAddressCreator) isProvidingIPAddress() bool {
	return versioned.adjustedVersion != IndeterminateIPVersion
}

func (versioned *versionedAddressCreator) isProvidingIPv4() bool {
	return versioned.adjustedVersion == IPv4
}

func (versioned *versionedAddressCreator) isProvidingIPv6() bool {
	return versioned.adjustedVersion == IPv6
}

func (versioned *versionedAddressCreator) getProviderIPVersion() IPVersion {
	return versioned.adjustedVersion
}

func (versioned *versionedAddressCreator) getType() ipType {
	return fromVersion(versioned.adjustedVersion)
}

func (versioned *versionedAddressCreator) getVersionedAddress(version IPVersion) (addr *IPAddress, err addrerr.IncompatibleAddressError) {
	index := version.index()
	if index >= IndeterminateIPVersion.index() {
		return
	}
	if versioned.versionedAddressCreatorFunc != nil {
		addr = (*IPAddress)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&versioned.versionedValues[index]))))
		if addr == nil {
			addr, err = versioned.versionedAddressCreatorFunc(version)
			if err == nil {
				dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&versioned.versionedValues[index]))
				atomicStorePointer(dataLoc, unsafe.Pointer(addr))
			}
		}
		return
	}
	addr = versioned.versionedValues[index]
	return
}

func emptyAddressCreator(emptyStrOption addrstrparam.EmptyStrOption, version IPVersion, zone Zone) (addrCreator func() (address, hostAddress *IPAddress), versionedCreator func() *IPAddress) {
	preferIPv6 := version.IsIPv6()
	double := func(one *IPAddress) (address, hostAddress *IPAddress) {
		return one, one
	}
	if emptyStrOption == addrstrparam.NoAddressOption {
		addrCreator = func() (*IPAddress, *IPAddress) { return double(nil) }
		versionedCreator = func() *IPAddress { return nil }
	} else if emptyStrOption == addrstrparam.LoopbackOption {
		if preferIPv6 {
			if len(zone) > 0 {
				ipv6WithZoneLoop := func() *IPAddress {
					network := ipv6Network
					creator := network.getIPAddressCreator()
					return creator.createAddressInternalFromBytes(network.GetLoopback().Bytes(), zone)
				}
				versionedCreator = ipv6WithZoneLoop
				addrCreator = func() (*IPAddress, *IPAddress) { return double(ipv6WithZoneLoop()) }
			} else {
				ipv6Loop := func() *IPAddress {
					return ipv6Network.GetLoopback()
				}
				versionedCreator = ipv6Loop
				addrCreator = func() (*IPAddress, *IPAddress) { return double(ipv6Loop()) }
			}
		} else {
			ipv4Loop := func() *IPAddress {
				return ipv4Network.GetLoopback()
			}
			addrCreator = func() (*IPAddress, *IPAddress) { return double(ipv4Loop()) }
			versionedCreator = ipv4Loop
		}
	} else { // EmptyStrParsedAs() == ZeroAddressOption
		if preferIPv6 {
			if len(zone) > 0 {
				ipv6WithZoneZero := func() *IPAddress {
					network := ipv6Network
					creator := network.getIPAddressCreator()
					return creator.createAddressInternalFromBytes(zeroIPv6.Bytes(), zone)
				}
				versionedCreator = ipv6WithZoneZero
				addrCreator = func() (*IPAddress, *IPAddress) { return double(ipv6WithZoneZero()) }
			} else {
				ipv6Zero := func() *IPAddress {
					return zeroIPv6.ToIP()
				}
				versionedCreator = ipv6Zero
				addrCreator = func() (*IPAddress, *IPAddress) { return double(ipv6Zero()) }
			}
		} else {
			ipv4Zero := func() *IPAddress {
				return zeroIPv4.ToIP()
			}
			addrCreator = func() (*IPAddress, *IPAddress) { return double(ipv4Zero()) }
			versionedCreator = ipv4Zero
		}
	}
	return
}

func newEmptyAddrCreator(options addrstrparam.IPAddressStringParams, zone Zone) *emptyAddrCreator {
	var version = IPVersion(options.GetPreferredVersion())
	// EmptyStrParsedAs chooses whether to produce loopbacks, zero addresses, or nothing for the empty string ""
	addrCreator, versionedCreator := emptyAddressCreator(options.EmptyStrParsedAs(), version, zone)
	cached := cachedAddressProvider{
		addressCreator: func() (address, hostAddress *IPAddress, addrErr, hosterr addrerr.IncompatibleAddressError) {
			address, hostAddress = addrCreator()
			return
		},
	}
	versionedCreatorFunc := func(v IPVersion) *IPAddress {
		addresses := cached.addresses
		if addresses != nil {
			addr := addresses.address
			if v == addr.GetIPVersion() {
				return addr
			}
		}
		if v.IsIndeterminate() {
			return versionedCreator()
		}
		_, vCreator := emptyAddressCreator(options.EmptyStrParsedAs(), v, zone)
		return vCreator()
	}
	versionedAddressCreatorFunc := func(version IPVersion) (*IPAddress, addrerr.IncompatibleAddressError) {
		return versionedCreatorFunc(version), nil
	}
	return &emptyAddrCreator{
		versionedAddressCreator: versionedAddressCreator{
			adjustedVersion:             version,
			parameters:                  options,
			cachedAddressProvider:       cached,
			versionedAddressCreatorFunc: versionedAddressCreatorFunc,
		},
		zone: zone,
	}
}

type emptyAddrCreator struct {
	versionedAddressCreator

	zone Zone
}

func (loop *emptyAddrCreator) providerCompare(other ipAddressProvider) (int, addrerr.IncompatibleAddressError) {
	return providerCompare(loop, other)
}

func (loop *emptyAddrCreator) providerEquals(other ipAddressProvider) (bool, addrerr.IncompatibleAddressError) {
	return providerEquals(loop, other)
}

func (loop *emptyAddrCreator) getProviderNetworkPrefixLen() PrefixLen {
	return nil
}

type adjustedAddressCreator struct {
	versionedAddressCreator

	networkPrefixLength PrefixLen
}

func (adjusted *adjustedAddressCreator) getProviderNetworkPrefixLen() PrefixLen {
	return adjusted.networkPrefixLength
}

func (adjusted *adjustedAddressCreator) getProviderAddress() (*IPAddress, addrerr.IncompatibleAddressError) {
	if !adjusted.isProvidingIPAddress() {
		return nil, nil
	}
	return adjusted.versionedAddressCreator.getProviderAddress()
}

func (adjusted *adjustedAddressCreator) getProviderHostAddress() (*IPAddress, addrerr.IncompatibleAddressError) {
	if !adjusted.isProvidingIPAddress() {
		return nil, nil
	}
	return adjusted.versionedAddressCreator.getProviderHostAddress()
}

func newMaskCreator(options addrstrparam.IPAddressStringParams, adjustedVersion IPVersion, networkPrefixLength PrefixLen) *maskCreator {
	if adjustedVersion == IndeterminateIPVersion {
		adjustedVersion = IPVersion(options.GetPreferredVersion())
	}
	createVersionedMask := func(version IPVersion, prefLen PrefixLen, withPrefixLength bool) *IPAddress {
		_ = withPrefixLength
		if version == IPv4 {
			network := ipv4Network
			return network.GetNetworkMask(prefLen.bitCount())
		} else if version == IPv6 {
			network := ipv6Network
			return network.GetNetworkMask(prefLen.bitCount())
		}
		return nil
	}
	versionedAddressCreatorFunc := func(version IPVersion) (*IPAddress, addrerr.IncompatibleAddressError) {
		return createVersionedMask(version, networkPrefixLength, true), nil
	}
	maskCreatorFunc := func() (address, hostAddress *IPAddress) {
		prefLen := networkPrefixLength
		return createVersionedMask(adjustedVersion, prefLen, true),
			createVersionedMask(adjustedVersion, prefLen, false)
	}
	addrCreator := func() (address, hostAddress *IPAddress, addrErr, hosterr addrerr.IncompatibleAddressError) {
		address, hostAddress = maskCreatorFunc()
		return
	}
	cached := cachedAddressProvider{addressCreator: addrCreator}
	return &maskCreator{
		adjustedAddressCreator{
			networkPrefixLength: networkPrefixLength,
			versionedAddressCreator: versionedAddressCreator{
				adjustedVersion:             adjustedVersion,
				parameters:                  options,
				cachedAddressProvider:       cached,
				versionedAddressCreatorFunc: versionedAddressCreatorFunc,
			},
		},
	}
}

type maskCreator struct {
	adjustedAddressCreator
}

func newAllCreator(qualifier *parsedHostIdentifierStringQualifier, adjustedVersion IPVersion, originator HostIdentifierString, options addrstrparam.IPAddressStringParams) ipAddressProvider {
	result := &allCreator{
		adjustedAddressCreator: adjustedAddressCreator{
			networkPrefixLength: qualifier.getEquivalentPrefixLen(),
			versionedAddressCreator: versionedAddressCreator{
				adjustedVersion: adjustedVersion,
				parameters:      options,
			},
		},
		originator: originator,
		qualifier:  *qualifier,
	}
	result.addressCreator = result.createAddrs
	result.versionedAddressCreatorFunc = result.versionedCreate
	return result
}

type allCreator struct {
	adjustedAddressCreator

	originator HostIdentifierString
	qualifier  parsedHostIdentifierStringQualifier
}

func (all *allCreator) getType() ipType {
	if !all.adjustedVersion.IsIndeterminate() {
		return fromVersion(all.adjustedVersion)
	}
	return allType
}

func (all *allCreator) providerCompare(other ipAddressProvider) (int, addrerr.IncompatibleAddressError) {
	return providerCompare(all, other)
}

func (all *allCreator) providerEquals(other ipAddressProvider) (bool, addrerr.IncompatibleAddressError) {
	return providerEquals(all, other)
}

// providing **all** addresses of any IP version, ie "*", not "*.*" or "*:*"
func (all *allCreator) isProvidingAllAddresses() bool {
	return all.adjustedVersion == IndeterminateIPVersion
}

func (all *allCreator) getProviderNetworkPrefixLen() PrefixLen {
	return all.qualifier.getEquivalentPrefixLen()
}

func (all *allCreator) getProviderMask() *IPAddress {
	return all.qualifier.getMaskLower()
}

func (all *allCreator) isSequential() bool {
	addr, _ := all.getProviderAddress()
	if addr != nil {
		return addr.IsSequential()
	}
	return false
}

func (all *allCreator) getProviderHostAddress() (res *IPAddress, err addrerr.IncompatibleAddressError) {
	if !all.isProvidingIPAddress() {
		return nil, nil
	}
	_, res, _, err = all.createAddrs()
	return
}

func (all *allCreator) getProviderAddress() (res *IPAddress, err addrerr.IncompatibleAddressError) {
	if !all.isProvidingIPAddress() {
		return
	}
	res, _, err, _ = all.createAddrs()
	return
}

func (all *allCreator) getProviderSeqRange() *SequentialRange[*IPAddress] {
	if all.isProvidingAllAddresses() {
		return nil
	}
	return all.createRange()
}

func (all *allCreator) createRange() (rng *SequentialRange[*IPAddress]) {
	rng, _, _, _, _ = all.createAll()
	return
}

func (all *allCreator) createAddrs() (addr *IPAddress, hostAddr *IPAddress, addrErr, hostErr addrerr.IncompatibleAddressError) {
	_, addr, hostAddr, addrErr, hostErr = all.createAll()
	return
}

func (all *allCreator) createAll() (rng *SequentialRange[*IPAddress], addr *IPAddress, hostAddr *IPAddress, addrErr, hostErr addrerr.IncompatibleAddressError) {
	addrs := (*addressResult)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&all.addresses))))
	if addrs == nil {
		var lower, upper *IPAddress
		addr, hostAddr, lower, upper, addrErr = createAllAddress(
			all.adjustedVersion,
			&all.qualifier,
			all.originator)
		rng = lower.SpanWithRange(upper)
		addresses := &addressResult{
			address:     addr,
			hostAddress: hostAddr,
			addrErr:     addrErr,
			rng:         rng,
		}
		dataLoc := (*unsafe.Pointer)(unsafe.Pointer(&all.addresses))
		atomicStorePointer(dataLoc, unsafe.Pointer(addresses))
	} else {
		rng, addr, hostAddr, addrErr, hostErr = addrs.rng, addrs.address, addrs.hostAddress, addrs.addrErr, addrs.hostErr
	}
	return
}

func (all *allCreator) versionedCreate(version IPVersion) (addr *IPAddress, addrErr addrerr.IncompatibleAddressError) {
	if version == all.adjustedVersion {
		return all.getProviderAddress()
	} else if all.adjustedVersion != IndeterminateIPVersion {
		return nil, nil
	}
	addr, _, _, _, addrErr = createAllAddress(
		version,
		&all.qualifier,
		all.originator)
	return
}

func (all *allCreator) prefixContainsProvider(otherProvider ipAddressProvider) boolSetting {
	return all.containsProviderFunc(otherProvider, (*IPAddress).prefixContains)
}

func (all *allCreator) containsProvider(otherProvider ipAddressProvider) (res boolSetting) {
	return all.containsProviderFunc(otherProvider, (*IPAddress).contains)
}

func (all *allCreator) containsProviderFunc(otherProvider ipAddressProvider, functor func(*IPAddress, AddressType) bool) (res boolSetting) {
	if otherProvider.isInvalid() {
		return boolSetting{true, false}
	} else if all.adjustedVersion == IndeterminateIPVersion {
		return boolSetting{true, true}
	} else if all.adjustedVersion != otherProvider.getProviderIPVersion() {
		return boolSetting{true, false}
	} else if all.qualifier.getMaskLower() == nil && all.qualifier.getZone() == NoZone {
		return boolSetting{true, true}
	} else if addr, err := all.getProviderAddress(); err != nil {
		return boolSetting{true, false}
	} else if otherAddr, err := all.getProviderAddress(); err != nil {
		return boolSetting{true, false}
	} else {
		return boolSetting{true, functor(addr, otherAddr)}
		//return boolSetting{true, addr.Contains(otherAddr)}
	}
}

// TODO LATER getDivisionGrouping()
//
//		@Override
//		public IPAddressDivisionSeries getDivisionGrouping() throwsaddrerr.IncompatibleAddressError {
//			if(isProvidingAllAddresses()) {
//				return null;
//			}
//			IPAddressNetwork<?, ?, ?, ?, ?> network = adjustedVersion.IsIPv4() ?
//					options.getIPv4Parameters().getNetwork() : options.getIPv6Parameters().getNetwork();
//			IPAddress mask = getProviderMask();
//			if(mask != null && mask.getBlockMaskPrefixLen(true) == null) {
//				// there is a mask
//				Integer hostMaskPrefixLen = mask.getBlockMaskPrefixLen(false);
//				if(hostMaskPrefixLen == null) { // not a host mask
//					throw newaddrerr.IncompatibleAddressError(getProviderAddress(), mask, "ipaddress.error.maskMismatch");
//				}
//				IPAddress hostMask = network.getHostMask(hostMaskPrefixLen);
//				return hostMask.toPrefixBlock();
//			}
//			IPAddressDivisionSeries grouping;
//			if(adjustedVersion.IsIPv4()) {
//				grouping = new IPAddressDivisionGrouping(new IPAddressBitsDivision[] {
//							new IPAddressBitsDivision(0, IPv4Address.MAX_VALUE, IPv4Address.BIT_COUNT, IPv4Address.DEFAULT_TEXTUAL_RADIX, network, qualifier.getEquivalentPrefixLength())
//						}, network);
//			} else if(adjustedVersion.IsIPv6()) {
//				byte upperBytes[] = new byte[16];
//				Arrays.fill(upperBytes, (byte) 0xff);
//				grouping = new IPAddressLargeDivisionGrouping(new IPAddressLargeDivision[] {new IPAddressLargeDivision(new byte[IPv6Address.BYTE_COUNT], upperBytes, IPv6Address.BIT_COUNT, IPv6Address.DEFAULT_TEXTUAL_RADIX, network, qualifier.getEquivalentPrefixLength())}, network);
//			} else {
//				grouping = null;
//			}
//			return grouping;
//		}
//	}
