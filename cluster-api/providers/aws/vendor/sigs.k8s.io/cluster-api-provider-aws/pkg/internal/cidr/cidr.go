/*
Copyright 2020 The Kubernetes Authors.

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

package cidr

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"

	"github.com/pkg/errors"
)

// SplitIntoSubnetsIPv4 splits a IPv4 CIDR into a specified number of subnets.
// If the number of required subnets isn't a power of 2 then then CIDR will be split
// into the the next highest power of 2 and you will end up with unused ranges.
// NOTE: this code is adapted from kops https://github.com/kubernetes/kops/blob/c323819e6480d71bad8d21184516e3162eaeca8f/pkg/util/subnet/subnet.go#L46
func SplitIntoSubnetsIPv4(cidrBlock string, numSubnets int) ([]*net.IPNet, error) {
	_, parent, err := net.ParseCIDR(cidrBlock)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse CIDR")
	}

	subnetBits := math.Ceil(math.Log2(float64(numSubnets)))

	networkLen, addrLen := parent.Mask.Size()
	modifiedNetworkLen := networkLen + int(subnetBits)

	if modifiedNetworkLen > addrLen {
		return nil, errors.Errorf("cidr %s cannot accommodate %d subnets", cidrBlock, numSubnets)
	}

	var subnets []*net.IPNet
	for i := 0; i < numSubnets; i++ {
		ip4 := parent.IP.To4()
		if ip4 == nil {
			return nil, errors.Errorf("unexpected IP address type: %s", parent)
		}

		n := binary.BigEndian.Uint32(ip4)
		n += uint32(i) << uint(32-modifiedNetworkLen)
		subnetIP := make(net.IP, len(ip4))
		binary.BigEndian.PutUint32(subnetIP, n)

		subnets = append(subnets, &net.IPNet{
			IP:   subnetIP,
			Mask: net.CIDRMask(modifiedNetworkLen, 32),
		})
	}

	return subnets, nil
}

// GetIPv4Cidrs gets the IPv4 CIDRs from a string slice.
func GetIPv4Cidrs(cidrs []string) ([]string, error) {
	found := []string{}

	for i := range cidrs {
		cidr := cidrs[i]

		ip, _, err := net.ParseCIDR(cidr)
		if err != nil {
			return found, fmt.Errorf("parsing %s as cidr: %w", cidr, err)
		}

		ipv4 := ip.To4()
		if ipv4 != nil {
			found = append(found, cidr)
		}
	}

	return found, nil
}

// GetIPv6Cidrs gets the IPv6 CIDRs from a string slice.
func GetIPv6Cidrs(cidrs []string) ([]string, error) {
	found := []string{}

	for i := range cidrs {
		cidr := cidrs[i]

		ip, _, err := net.ParseCIDR(cidr)
		if err != nil {
			return found, fmt.Errorf("parsing %s as cidr: %w", cidr, err)
		}

		ipv4 := ip.To4()
		if ipv4 == nil {
			found = append(found, cidr)
		}
	}

	return found, nil
}
