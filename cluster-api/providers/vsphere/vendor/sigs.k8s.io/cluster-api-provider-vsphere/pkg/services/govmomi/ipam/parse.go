/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ipam

import (
	"fmt"
	"net/netip"

	ipamv1 "sigs.k8s.io/cluster-api/exp/ipam/api/v1alpha1"
)

// prefixesAsStrings converts []netip.Prefix to []string.
func prefixesAsStrings(prefixes []netip.Prefix) []string {
	prefixSrings := make([]string, 0, len(prefixes))
	for _, prefix := range prefixes {
		prefixSrings = append(prefixSrings, prefix.String())
	}
	return prefixSrings
}

// parseAddressWithPrefix converts a *ipamv1.IPAddress to a string, e.g. '10.0.0.1/24'.
func parseAddressWithPrefix(ipamAddress *ipamv1.IPAddress) (netip.Prefix, error) {
	addressWithPrefix := fmt.Sprintf("%s/%d", ipamAddress.Spec.Address, ipamAddress.Spec.Prefix)
	parsedPrefix, err := netip.ParsePrefix(addressWithPrefix)
	if err != nil {
		return netip.Prefix{}, fmt.Errorf("IPAddress %s/%s has invalid ip address: %q",
			ipamAddress.Namespace,
			ipamAddress.Name,
			addressWithPrefix,
		)
	}

	return parsedPrefix, nil
}

// parseGateway parses the gateway address on a ipamv1.IPAddress and ensures it
// does not conflict with the gateway addresses parsed from other
// ipamv1.IPAddresses on the current device. Gateway addresses must be the same
// family as the address on the ipamv1.IPAddress. Gateway addresses of one
// family must match the other addresses of the same family. A gateway address
// is optional. If it is not set this function returns `nil, nil`.
func parseGateway(ipamAddress *ipamv1.IPAddress, addressWithPrefix netip.Prefix, ipamDeviceConfig ipamDeviceConfig) (*netip.Addr, error) {
	if ipamAddress.Spec.Gateway == "" {
		return nil, nil
	}

	gatewayAddr, err := netip.ParseAddr(ipamAddress.Spec.Gateway)
	if err != nil {
		return nil, fmt.Errorf("IPAddress %s/%s has invalid gateway: %q",
			ipamAddress.Namespace,
			ipamAddress.Name,
			ipamAddress.Spec.Gateway,
		)
	}

	if addressWithPrefix.Addr().Is4() != gatewayAddr.Is4() {
		return nil, fmt.Errorf("IPAddress %s/%s has mismatched gateway and address IP families",
			ipamAddress.Namespace,
			ipamAddress.Name,
		)
	}

	if gatewayAddr.Is4() {
		if areGatewaysMismatched(ipamDeviceConfig.NetworkSpecGateway4, ipamAddress.Spec.Gateway) {
			return nil, fmt.Errorf("the IPv4 Gateway for IPAddress %s does not match the Gateway4 already configured on device (index %d)",
				ipamAddress.Name,
				ipamDeviceConfig.DeviceIndex,
			)
		}
		if areGatewaysMismatched(ipamDeviceConfig.IPAMConfigGateway4, ipamAddress.Spec.Gateway) {
			return nil, fmt.Errorf("the IPv4 IPAddresses assigned to the same device (index %d) do not have the same gateway",
				ipamDeviceConfig.DeviceIndex,
			)
		}
	} else {
		if areGatewaysMismatched(ipamDeviceConfig.NetworkSpecGateway6, ipamAddress.Spec.Gateway) {
			return nil, fmt.Errorf("the IPv6 Gateway for IPAddress %s does not match the Gateway6 already configured on device (index %d)",
				ipamAddress.Name,
				ipamDeviceConfig.DeviceIndex,
			)
		}
		if areGatewaysMismatched(ipamDeviceConfig.IPAMConfigGateway6, ipamAddress.Spec.Gateway) {
			return nil, fmt.Errorf("the IPv6 IPAddresses assigned to the same device (index %d) do not have the same gateway",
				ipamDeviceConfig.DeviceIndex,
			)
		}
	}

	return &gatewayAddr, nil
}

// areGatewaysMismatched checks that a gateway for a device is equal to an
// IPAddresses gateway. We can assume that IPAddresses will always have
// gateways so we do not need to check for empty string. It is possible to
// configure a device and not a gateway, we don't want to fail in that case.
func areGatewaysMismatched(deviceGateway, ipAddressGateway string) bool {
	return deviceGateway != "" && deviceGateway != ipAddressGateway
}
