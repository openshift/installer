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
	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstrparam"
)

type parsedHostIdentifierStringQualifier struct {

	// if there is a port for the host, this will be its numeric value
	port    Port   // non-nil for a host with port
	service string // non-empty for host with a service instead of a port

	// if there is a prefix length for the address, this will be its numeric value
	networkPrefixLength PrefixLen //non-nil for a prefix-only address, sometimes non-nil for IPv4, IPv6

	// If instead of a prefix length a mask was provided, this is the mask.
	// We can also have both a prefix length and mask if one is added when merging qualifiers  */'
	mask *parsedIPAddress

	// overrides the parsed mask if present
	mergedMask *IPAddress

	// this is the IPv6 scope id or network interface name
	zone    Zone
	isZoned bool
}

func (parsedQual *parsedHostIdentifierStringQualifier) clearPortOrService() {
	parsedQual.port = nil
	parsedQual.service = ""
}

func (parsedQual *parsedHostIdentifierStringQualifier) clearPrefixOrMask() {
	parsedQual.networkPrefixLength = nil
	parsedQual.mask = nil
}

func (parsedQual *parsedHostIdentifierStringQualifier) merge(other *parsedHostIdentifierStringQualifier) (err addrerr.IncompatibleAddressError) {
	if parsedQual.networkPrefixLength == nil ||
		(other.networkPrefixLength != nil && other.networkPrefixLength.bitCount() < parsedQual.networkPrefixLength.bitCount()) {
		parsedQual.networkPrefixLength = other.networkPrefixLength
	}
	if parsedQual.mask == nil {
		parsedQual.mask = other.mask
	} else {
		otherMask := other.getMaskLower()
		if otherMask != nil {
			parsedQual.mergedMask, err = parsedQual.getMaskLower().Mask(otherMask)
		}
	}
	return
}

func (parsedQual *parsedHostIdentifierStringQualifier) getMaskLower() *IPAddress {
	if mask := parsedQual.mergedMask; mask != nil {
		return mask
	}
	if mask := parsedQual.mask; mask != nil {
		return mask.getValForMask()
	}
	return nil
}

func (parsedQual *parsedHostIdentifierStringQualifier) getNetworkPrefixLen() PrefixLen {
	return parsedQual.networkPrefixLength
}

func (parsedQual *parsedHostIdentifierStringQualifier) getEquivalentPrefixLen() PrefixLen {
	pref := parsedQual.getNetworkPrefixLen()
	if pref == nil {
		mask := parsedQual.getMaskLower()
		if mask != nil {
			pref = mask.GetBlockMaskPrefixLen(true)
		}
	}
	return pref
}

// we distinguish callers with empty zones vs callers in which there was no zone indicator
func (parsedQual *parsedHostIdentifierStringQualifier) setZone(z *Zone) {
	if z != nil {
		parsedQual.zone = *z
		parsedQual.isZoned = true
	}
}

func (parsedQual *parsedHostIdentifierStringQualifier) getZone() Zone {
	return parsedQual.zone
}

func (parsedQual *parsedHostIdentifierStringQualifier) getPort() Port {
	return parsedQual.port
}

func (parsedQual *parsedHostIdentifierStringQualifier) getService() string {
	return parsedQual.service
}

func (parsedQual *parsedHostIdentifierStringQualifier) inferVersion(validationOptions addrstrparam.IPAddressStringParams) IPVersion {
	if parsedQual.networkPrefixLength != nil {
		if parsedQual.networkPrefixLength.bitCount() > IPv4BitCount &&
			!validationOptions.GetIPv4Params().AllowsPrefixesBeyondAddressSize() {
			return IPv6
		}
	} else if mask := parsedQual.mask; mask != nil {
		if mask.isProvidingIPv6() {
			return IPv6
		} else if mask.isProvidingIPv4() {
			return IPv4
		}
	}
	if parsedQual.isZoned {
		//if parsedQual.zone != "" {
		return IPv6
	}
	return IndeterminateIPVersion
}
