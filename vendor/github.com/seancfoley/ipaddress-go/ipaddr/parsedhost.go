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
	"github.com/seancfoley/ipaddress-go/ipaddr/addrstrparam"
	"strings"

	"github.com/seancfoley/ipaddress-go/ipaddr/addrerr"
)

type embeddedAddress struct {
	isUNCIPv6Literal, isReverseDNS bool

	addressStringError addrerr.AddressStringError

	addressProvider ipAddressProvider
}

var (
	noQualifier = &parsedHostIdentifierStringQualifier{}
)

type parsedHostCache struct {
	normalizedLabels []string
	host             string
}

type parsedHost struct {
	separatorIndices []int // can be nil
	normalizedFlags  []bool

	labelsQualifier parsedHostIdentifierStringQualifier

	embeddedAddress embeddedAddress

	originalStr string

	*parsedHostCache

	params addrstrparam.HostNameParams
}

func (host *parsedHost) getQualifier() *parsedHostIdentifierStringQualifier {
	return &host.labelsQualifier
}

func (host *parsedHost) isIPv6Address() bool {
	return host.hasEmbeddedAddress() && host.getAddressProvider().isProvidingIPv6()
}

func (host *parsedHost) getPort() Port {
	return host.labelsQualifier.getPort()
}

func (host *parsedHost) getService() string {
	return host.labelsQualifier.getService()
}

func (host *parsedHost) getNetworkPrefixLen() PrefixLen {
	return host.labelsQualifier.getNetworkPrefixLen()
}

func (host *parsedHost) getEquivalentPrefixLen() PrefixLen {
	return host.labelsQualifier.getEquivalentPrefixLen()
}

func (host *parsedHost) getMask() *IPAddress {
	return host.labelsQualifier.getMaskLower()
}

func (host *parsedHost) getAddressProvider() ipAddressProvider {
	return host.embeddedAddress.addressProvider
}

func (host *parsedHost) hasEmbeddedAddress() bool {
	return host.embeddedAddress.addressProvider != nil
}

func (host *parsedHost) isAddressString() bool {
	return host.getAddressProvider() != nil
}

func (host *parsedHost) asAddress() (*IPAddress, addrerr.IncompatibleAddressError) {
	if host.hasEmbeddedAddress() {
		return host.getAddressProvider().getProviderAddress()
	}
	return nil, nil
}

func (host *parsedHost) mapString(addressProvider ipAddressProvider) string {
	if addressProvider.isProvidingAllAddresses() {
		return SegmentWildcardStr
		//} else if addressProvider.isProvidingPrefixOnly() {
		//return IPAddressNetwork.getPrefixString(addressProvider.getProviderNetworkPrefixLength())
	} else if addressProvider.isProvidingEmpty() {
		return ""
	}
	return host.originalStr
}

func (host *parsedHost) asGenericAddressString() *IPAddressString {
	if host.hasEmbeddedAddress() {
		addressProvider := host.getAddressProvider()
		if addressProvider.isProvidingAllAddresses() {
			return NewIPAddressStringParams(SegmentWildcardStr, addressProvider.getParameters())
		} else if addressProvider.isProvidingEmpty() {
			return NewIPAddressStringParams("", addressProvider.getParameters())
		} else {
			addr, err := addressProvider.getProviderAddress()
			if err != nil {
				return NewIPAddressStringParams(host.originalStr, addressProvider.getParameters())
			}
			return addr.ToAddressString()
		}
	}
	return nil
}

func (host *parsedHost) getHost() string {
	return host.buildHostString()
}

func (host *parsedHost) buildHostString() string {
	if host.parsedHostCache == nil {
		var hostStr string
		if host.hasEmbeddedAddress() {
			addressProvider := host.getAddressProvider()
			addr, err := addressProvider.getProviderAddress()
			if err == nil && addr != nil {
				section := addr.GetSection()
				//port was stripped out
				//mask and prefix removed by toNormalizedWildcardString
				//getSection() removes zone
				hostStr = section.ToCanonicalWildcardString()
			} else {
				hostStr = host.mapString(addressProvider)
			}
		} else {
			var label string
			normalizedFlags := host.normalizedFlags
			var hostStrBuilder strings.Builder
			for i, lastSep := 0, -1; i < len(host.separatorIndices); i++ {
				index := host.separatorIndices[i]
				if len(normalizedFlags) > 0 && !normalizedFlags[i] {
					var normalizedLabelBuilder strings.Builder
					normalizedLabelBuilder.Grow((index - lastSep) - 1)
					for j := lastSep + 1; j < index; j++ {
						c := host.originalStr[j]
						if c >= 'A' && c <= 'Z' {
							c = c + ('a' - 'A')
						}
						normalizedLabelBuilder.WriteByte(c)
					}
					label = normalizedLabelBuilder.String()
				} else {
					label = host.originalStr[lastSep+1 : index]
				}
				if i > 0 {
					hostStrBuilder.WriteByte(LabelSeparator)
				}
				hostStrBuilder.WriteString(label)
				lastSep = index
			}
			hostStr = hostStrBuilder.String()
		}
		return hostStr
	}
	return host.parsedHostCache.host
}

func (host *parsedHost) buildNormalizedLabels() []string {
	if host.parsedHostCache == nil {
		var normalizedLabels []string
		if host.hasEmbeddedAddress() {
			addressProvider := host.getAddressProvider()
			addr, err := addressProvider.getProviderAddress()
			if err == nil && addr != nil {
				section := addr.GetSection()
				normalizedLabels = section.GetSegmentStrings()
			} else {
				hostStr := host.mapString(addressProvider)
				if addressProvider.isProvidingEmpty() {
					normalizedLabels = []string{}
				} else {
					normalizedLabels = []string{hostStr}
				}
			}
		} else {
			normalizedLabels = make([]string, len(host.separatorIndices))
			normalizedFlags := host.normalizedFlags

			for i, lastSep := 0, -1; i < len(normalizedLabels); i++ {
				index := host.separatorIndices[i]
				if len(normalizedFlags) > 0 && !normalizedFlags[i] {
					var normalizedLabelBuilder strings.Builder
					normalizedLabelBuilder.Grow((index - lastSep) - 1)
					for j := lastSep + 1; j < index; j++ {
						c := host.originalStr[j]
						if c >= 'A' && c <= 'Z' {
							c = c + ('a' - 'A')
						}
						normalizedLabelBuilder.WriteByte(c)
					}
					normalizedLabels[i] = normalizedLabelBuilder.String()
				} else {
					normalizedLabels[i] = host.originalStr[lastSep+1 : index]
				}
				lastSep = index
			}
		}
		return normalizedLabels
	}
	return host.parsedHostCache.normalizedLabels
}

func (host *parsedHost) getNormalizedLabels() []string {
	return host.buildNormalizedLabels()
}

func (host *parsedHost) isUNCIPv6Literal() bool {
	return host.embeddedAddress.isUNCIPv6Literal
}

func (host *parsedHost) isReverseDNS() bool {
	return host.embeddedAddress.isReverseDNS
}
