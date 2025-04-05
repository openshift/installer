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

// Package cidr provides utilities for working with CIDR blocks.
package cidr

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"

	"github.com/pkg/errors"
)

// SplitIntoSubnetsIPv4 splits a IPv4 CIDR into a specified number of subnets.
// If the number of required subnets isn't a power of 2 then CIDR will be split
// into the next highest power of 2, and you will end up with unused ranges.
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

	subnets := make([]*net.IPNet, 0, numSubnets)
	for i := range numSubnets {
		ip4 := parent.IP.To4()
		if ip4 == nil {
			return nil, errors.Errorf("unexpected IP address type: %s", parent)
		}

		n := binary.BigEndian.Uint32(ip4)
		n += uint32(i) << uint(32-modifiedNetworkLen) //#nosec G115
		subnetIP := make(net.IP, len(ip4))
		binary.BigEndian.PutUint32(subnetIP, n)

		subnets = append(subnets, &net.IPNet{
			IP:   subnetIP,
			Mask: net.CIDRMask(modifiedNetworkLen, 32),
		})
	}

	return subnets, nil
}

const subnetIDLocation = 7

// SplitIntoSubnetsIPv6 splits a IPv6 address into a specified number of subnets.
// AWS IPv6 based subnets **must always have a /64 prefix**. AWS provides an IPv6
// CIDR with /56 prefix. That's the initial CIDR. We must convert that to /64 and
// slice the subnets by increasing the subnet ID by 1.
// so given: 2600:1f14:e08:7400::/56
// sub1: 2600:1f14:e08:7400::/64
// sub2: 2600:1f14:e08:7401::/64
// sub3: 2600:1f14:e08:7402::/64
// sub4: 2600:1f14:e08:7403::/64
// This function can also be called with /64 prefix to further slice existing subnet
// addresses.
// When splitting further, we always have to take the LAST one to avoid collisions
// since the prefix stays the same, but the subnet ID increases.
// To see this restriction read https://docs.aws.amazon.com/vpc/latest/userguide/how-it-works.html#ipv4-ipv6-comparison
func SplitIntoSubnetsIPv6(cidrBlock string, numSubnets int) ([]*net.IPNet, error) {
	_, ipv6CidrBlock, err := net.ParseCIDR(cidrBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cidr block %s with error: %w", cidrBlock, err)
	}
	// update the prefix to 64.
	ipv6CidrBlock.Mask = net.CIDRMask(64, 128)
	subnets := make([]*net.IPNet, 0, numSubnets)
	for range numSubnets {
		ipv6CidrBlock.IP[subnetIDLocation]++
		newIP := net.ParseIP(ipv6CidrBlock.IP.String())
		v := &net.IPNet{
			IP:   newIP,
			Mask: net.CIDRMask(64, 128),
		}
		subnets = append(subnets, v)
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
