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
	"net"
	"unsafe"
)

type addressNetwork interface {
	getAddressCreator() parsedAddressCreator
}

// IPAddressNetwork represents a network of addresses of a single IP version providing a collection of standard address components for that version, such as masks and loopbacks.
type IPAddressNetwork interface {
	GetLoopback() *IPAddress

	GetNetworkMask(prefixLength BitCount) *IPAddress

	GetPrefixedNetworkMask(prefixLength BitCount) *IPAddress

	GetHostMask(prefixLength BitCount) *IPAddress

	GetPrefixedHostMask(prefixLength BitCount) *IPAddress

	getIPAddressCreator() ipAddressCreator

	addressNetwork
}

type ipAddressNetwork struct {
	subnetsMasksWithPrefix, subnetMasks, hostMasksWithPrefix, hostMasks []*IPAddress
}

type ipv6AddressNetwork struct {
	ipAddressNetwork
	creator ipv6AddressCreator
}

func (network *ipv6AddressNetwork) getIPAddressCreator() ipAddressCreator {
	return &network.creator
}

func (network *ipv6AddressNetwork) getAddressCreator() parsedAddressCreator {
	return &network.creator
}

func (network *ipv6AddressNetwork) GetLoopback() *IPAddress {
	return ipv6loopback.ToIP()
}

func (network *ipv6AddressNetwork) GetNetworkMask(prefLen BitCount) *IPAddress {
	return network.subnetMasks[adjustBits(IPv6, prefLen)]
}

func (network *ipv6AddressNetwork) GetPrefixedNetworkMask(prefLen BitCount) *IPAddress {
	return network.subnetsMasksWithPrefix[adjustBits(IPv6, prefLen)]
}

func (network *ipv6AddressNetwork) GetHostMask(prefLen BitCount) *IPAddress {
	return network.hostMasks[adjustBits(IPv6, prefLen)]
}

func (network *ipv6AddressNetwork) GetPrefixedHostMask(prefLen BitCount) *IPAddress {
	return network.hostMasksWithPrefix[adjustBits(IPv6, prefLen)]
}

var _ IPAddressNetwork = &ipv6AddressNetwork{}

func adjustBits(version IPVersion, bits BitCount) BitCount {
	if bits < 0 {
		bits = 0
	} else {
		addressBitLength := version.GetBitCount()
		if bits > addressBitLength {
			bits = addressBitLength
		}
	}
	return bits
}

// IPv6AddressNetwork is the implementation of IPAddressNetwork for IPv6
type IPv6AddressNetwork struct {
	*ipv6AddressNetwork
}

func (network IPv6AddressNetwork) GetLoopback() *IPv6Address {
	return ipv6loopback
}

func (network IPv6AddressNetwork) GetNetworkMask(prefLen BitCount) *IPv6Address {
	return network.ipv6AddressNetwork.GetNetworkMask(prefLen).ToIPv6()
}

func (network IPv6AddressNetwork) GetPrefixedNetworkMask(prefLen BitCount) *IPv6Address {
	return network.ipv6AddressNetwork.GetPrefixedNetworkMask(prefLen).ToIPv6()
}

func (network IPv6AddressNetwork) GetHostMask(prefLen BitCount) *IPv6Address {
	return network.ipv6AddressNetwork.GetHostMask(prefLen).ToIPv6()
}

func (network IPv6AddressNetwork) GetPrefixedHostMask(prefLen BitCount) *IPv6Address {
	return network.ipv6AddressNetwork.GetPrefixedHostMask(prefLen).ToIPv6()
}

func createIPv6AddressNetwork() *ipv6AddressNetwork {
	network := &ipv6AddressNetwork{
		ipAddressNetwork: ipAddressNetwork{
			make([]*IPAddress, IPv6BitCount+1),
			make([]*IPAddress, IPv6BitCount+1),
			make([]*IPAddress, IPv6BitCount+1),
			make([]*IPAddress, IPv6BitCount+1),
		},
	}
	populateNetwork(IPv6, &network.ipAddressNetwork, zeroIPv6Seg.ToDiv())
	for i := 0; i <= IPv6BitCount; i++ {
		addr := network.subnetMasks[i].ToIPv6()
		high, low := addr.Uint64Values()
		ipv6NetworkMasks[i] = [2]uint64{high, low}
	}
	return network
}

func populateNetwork(version IPVersion, network *ipAddressNetwork, zeroDiv *AddressDivision) {
	addressBitLength := version.GetBitCount()
	for i := 0; i <= addressBitLength; i++ {
		_ = createMask(version, zeroDiv, i, network.subnetMasks, true, false)
		_ = createMask(version, zeroDiv, i, network.subnetsMasksWithPrefix, true, true)
		_ = createMask(version, zeroDiv, i, network.hostMasks, false, false)
		_ = createMask(version, zeroDiv, i, network.hostMasksWithPrefix, false, true)
	}
}

var ipv6Network = createIPv6AddressNetwork()

var IPv6Network = &IPv6AddressNetwork{ipv6Network}

var ipv6NetworkMasks [IPv6BitCount + 1][2]uint64

//
//
//
//
//

type ipv4AddressNetwork struct {
	ipAddressNetwork
	creator ipv4AddressCreator
}

func (network *ipv4AddressNetwork) getIPAddressCreator() ipAddressCreator {
	return &network.creator
}

func (network *ipv4AddressNetwork) getAddressCreator() parsedAddressCreator {
	return &network.creator
}

func (network *ipv4AddressNetwork) GetLoopback() *IPAddress {
	return ipv4loopback.ToIP()
}

func (network *ipv4AddressNetwork) GetNetworkMask(prefLen BitCount) *IPAddress {
	return network.subnetMasks[adjustBits(IPv4, prefLen)]
}

func (network *ipv4AddressNetwork) GetPrefixedNetworkMask(prefLen BitCount) *IPAddress {
	return network.subnetsMasksWithPrefix[adjustBits(IPv4, prefLen)]
}

func (network *ipv4AddressNetwork) GetHostMask(prefLen BitCount) *IPAddress {
	return network.hostMasks[adjustBits(IPv4, prefLen)]
}

func (network *ipv4AddressNetwork) GetPrefixedHostMask(prefLen BitCount) *IPAddress {
	return network.hostMasksWithPrefix[adjustBits(IPv4, prefLen)]
}

var _ IPAddressNetwork = &ipv4AddressNetwork{}

// IPv4AddressNetwork is the implementation of IPAddressNetwork for IPv4
type IPv4AddressNetwork struct {
	*ipv4AddressNetwork
}

func (network IPv4AddressNetwork) GetLoopback() *IPv4Address {
	return ipv4loopback
}

func (network IPv4AddressNetwork) GetNetworkMask(prefLen BitCount) *IPv4Address {
	return network.ipv4AddressNetwork.GetNetworkMask(prefLen).ToIPv4()
}

func (network IPv4AddressNetwork) GetPrefixedNetworkMask(prefLen BitCount) *IPv4Address {
	return network.ipv4AddressNetwork.GetPrefixedNetworkMask(prefLen).ToIPv4()
}

func (network IPv4AddressNetwork) GetHostMask(prefLen BitCount) *IPv4Address {
	return network.ipv4AddressNetwork.GetHostMask(prefLen).ToIPv4()
}

func (network IPv4AddressNetwork) GetPrefixedHostMask(prefLen BitCount) *IPv4Address {
	return network.ipv4AddressNetwork.GetPrefixedHostMask(prefLen).ToIPv4()
}

func createIPv4AddressNetwork() *ipv4AddressNetwork {
	network := &ipv4AddressNetwork{
		ipAddressNetwork: ipAddressNetwork{
			make([]*IPAddress, IPv4BitCount+1),
			make([]*IPAddress, IPv4BitCount+1),
			make([]*IPAddress, IPv4BitCount+1),
			make([]*IPAddress, IPv4BitCount+1),
		},
	}
	populateNetwork(IPv4, &network.ipAddressNetwork, zeroIPv4Seg.ToDiv())
	for i := 0; i <= IPv4BitCount; i++ {
		addr := network.subnetMasks[i].ToIPv4()
		ipv4NetworkMasks[i] = addr.Uint32Value()
	}
	return network
}

var ipv4Network = createIPv4AddressNetwork()

var IPv4Network = &IPv4AddressNetwork{ipv4Network}

var ipv4NetworkMasks [IPv4BitCount + 1]uint32

func createMask(version IPVersion, zeroSeg *AddressDivision, networkPrefixLength BitCount, cache []*IPAddress, network, withPrefixLength bool) *IPAddress {
	bits := networkPrefixLength
	addressBitLength := version.GetBitCount()

	bits = adjustBits(version, bits)

	cacheIndex := bits

	subnet := cache[cacheIndex]
	if subnet != nil {
		return subnet
	}

	var onesSubnetIndex, zerosSubnetIndex int
	if network {
		onesSubnetIndex = int(addressBitLength)
		zerosSubnetIndex = 0
	} else {
		onesSubnetIndex = 0
		zerosSubnetIndex = int(addressBitLength)
	}
	segmentCount := version.GetSegmentCount()
	bitsPerSegment := version.GetBitsPerSegment()
	maxSegmentValue := version.GetMaxSegmentValue()

	onesSubnet := (*IPAddress)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache[onesSubnetIndex]))))
	if onesSubnet == nil {
		newSegments := createSegmentArray(segmentCount)
		if withPrefixLength {
			if network {
				segment := createAddressDivision(zeroSeg.deriveNewSeg(maxSegmentValue, nil))
				lastSegment := createAddressDivision(zeroSeg.deriveNewSeg(maxSegmentValue, cacheBitCount(bitsPerSegment) /* bitsPerSegment */))
				lastIndex := len(newSegments) - 1
				fillDivs(newSegments[:lastIndex], segment)
				newSegments[lastIndex] = lastSegment
				onesSubnet = createIPAddress(createSection(newSegments, cacheBitCount(addressBitLength), version.toType()), NoZone)
			} else {
				segment := createAddressDivision(zeroSeg.deriveNewSeg(maxSegmentValue, cacheBitCount(0)))
				fillDivs(newSegments, segment)
				onesSubnet = createIPAddress(createSection(newSegments, cacheBitCount(0), version.toType()), NoZone)
			}
		} else {
			segment := createAddressDivision(zeroSeg.deriveNewSeg(maxSegmentValue, nil))
			fillDivs(newSegments, segment)
			onesSubnet = createIPAddress(createSection(newSegments, nil, version.toType()), NoZone) /* address creation */
		}
		cache[onesSubnetIndex] = onesSubnet
	}
	zerosSubnet := (*IPAddress)(atomicLoadPointer((*unsafe.Pointer)(unsafe.Pointer(&cache[zerosSubnetIndex]))))
	if zerosSubnet == nil {
		newSegments := createSegmentArray(segmentCount)
		if withPrefixLength {
			prefLen := cacheBitCount(0)
			if network {
				segment := createAddressDivision(zeroSeg.deriveNewSeg(0, prefLen))
				fillDivs(newSegments, segment)
				zerosSubnet = createIPAddress(createSection(newSegments, prefLen, version.toType()), NoZone)
			} else {
				lastSegment := createAddressDivision(zeroSeg.deriveNewSeg(0, cacheBitCount(bitsPerSegment) /* bitsPerSegment */))
				lastIndex := len(newSegments) - 1
				fillDivs(newSegments[:lastIndex], zeroSeg)
				newSegments[lastIndex] = lastSegment
				zerosSubnet = createIPAddress(createSection(newSegments, cacheBitCount(addressBitLength), version.toType()), NoZone)
			}
		} else {
			segment := createAddressDivision(zeroSeg.deriveNewSeg(0, nil))
			fillDivs(newSegments, segment)
			zerosSubnet = createIPAddress(createSection(newSegments, nil, version.toType()), NoZone)
		}
		cache[zerosSubnetIndex] = zerosSubnet
	}
	prefix := bits
	onesSegment := onesSubnet.getDivision(0)
	zerosSegment := zerosSubnet.getDivision(0)
	newSegments := createSegmentArray(segmentCount)[:0]
	i := 0
	for ; bits > 0; i, bits = i+1, bits-bitsPerSegment {
		if bits <= bitsPerSegment {
			var segment *AddressDivision

			//first do a check whether we have already created a segment like the one we need
			offset := ((bits - 1) % bitsPerSegment) + 1
			for j, entry := 0, offset; j < segmentCount; j, entry = j+1, entry+bitsPerSegment {
				//for j := 0, entry = offset; j < segmentCount; j++, entry += bitsPerSegment {
				if entry != cacheIndex { //we already know that the entry at cacheIndex is nil
					prev := cache[entry]
					if prev != nil {
						segment = prev.getDivision(j)
						break
					}
				}
			}

			//if none of the other addresses with a similar segment are created yet, we need a new segment.
			if segment == nil {
				if network {
					mask := maxSegmentValue & (maxSegmentValue << uint(bitsPerSegment-bits))
					if withPrefixLength {
						segment = createAddressDivision(zeroSeg.deriveNewSeg(mask, getDivisionPrefixLength(bitsPerSegment, bits)))
					} else {
						segment = createAddressDivision(zeroSeg.deriveNewSeg(mask, nil))
					}
				} else {
					mask := maxSegmentValue & ^(maxSegmentValue << uint(bitsPerSegment-bits))
					if withPrefixLength {
						segment = createAddressDivision(zeroSeg.deriveNewSeg(mask, getDivisionPrefixLength(bitsPerSegment, bits)))
					} else {
						segment = createAddressDivision(zeroSeg.deriveNewSeg(mask, nil))
					}
				}
			}
			newSegments = append(newSegments, segment)
		} else {
			if network {
				newSegments = append(newSegments, onesSegment)
			} else {
				newSegments = append(newSegments, zerosSegment)
			}
		}
	}
	for ; i < segmentCount; i++ {
		if network {
			newSegments = append(newSegments, zerosSegment)
		} else {
			newSegments = append(newSegments, onesSegment)
		}
	}
	var prefLen PrefixLen
	if withPrefixLength {
		prefLen = cacheBitCount(prefix)
	}
	subnet = createIPAddress(createSection(newSegments, prefLen, version.toType()), NoZone)
	cache[cacheIndex] = subnet
	return subnet
}

type macAddressNetwork struct {
	creator macAddressCreator
}

func (network *macAddressNetwork) getAddressCreator() parsedAddressCreator {
	return &network.creator
}

var macNetwork = &macAddressNetwork{}

var _ addressNetwork = &macAddressNetwork{}

var ipv4loopback = createIPv4Loopback()
var ipv6loopback = createIPv6Loopback()

func createIPv6Loopback() *IPv6Address {
	ipv6loopback, _ := NewIPv6AddressFromBytes(net.IPv6loopback)
	return ipv6loopback
}

func createIPv4Loopback() *IPv4Address {
	ipv4loopback, _ := NewIPv4AddressFromBytes([]byte{127, 0, 0, 1})
	return ipv4loopback
}
