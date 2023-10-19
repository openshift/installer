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
	"net"
	"testing"

	. "github.com/onsi/gomega"
)

func TestSplitIntoSubnetsIPv4(t *testing.T) {
	RegisterTestingT(t)
	tests := []struct {
		name        string
		cidrblock   string
		subnetcount int
		expected    []*net.IPNet
	}{
		{
			// https://aws.amazon.com/about-aws/whats-new/2018/10/amazon-eks-now-supports-additional-vpc-cidr-blocks/
			name:        "default secondary cidr block configuration with primary cidr",
			cidrblock:   "100.64.0.0/10",
			subnetcount: 3,
			expected: []*net.IPNet{
				{
					IP:   net.IPv4(100, 64, 0, 0).To4(),
					Mask: net.IPv4Mask(255, 240, 0, 0),
				},
				{
					IP:   net.IPv4(100, 80, 0, 0).To4(),
					Mask: net.IPv4Mask(255, 240, 0, 0),
				},
				{
					IP:   net.IPv4(100, 96, 0, 0).To4(),
					Mask: net.IPv4Mask(255, 240, 0, 0),
				},
			},
		},
		{
			// https://aws.amazon.com/about-aws/whats-new/2018/10/amazon-eks-now-supports-additional-vpc-cidr-blocks/
			name:        "default secondary cidr block configuration with alternative cidr",
			cidrblock:   "198.19.0.0/16",
			subnetcount: 3,
			expected: []*net.IPNet{
				{
					IP:   net.IPv4(198, 19, 0, 0).To4(),
					Mask: net.IPv4Mask(255, 255, 192, 0),
				},
				{
					IP:   net.IPv4(198, 19, 64, 0).To4(),
					Mask: net.IPv4Mask(255, 255, 192, 0),
				},
				{
					IP:   net.IPv4(198, 19, 128, 0).To4(),
					Mask: net.IPv4Mask(255, 255, 192, 0),
				},
			},
		},
		{
			name:        "slash 16 cidr with one subnet",
			cidrblock:   "1.1.0.0/16",
			subnetcount: 1,
			expected: []*net.IPNet{
				{
					IP:   net.IPv4(1, 1, 0, 0).To4(),
					Mask: net.IPv4Mask(255, 255, 0, 0),
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := SplitIntoSubnetsIPv4(tc.cidrblock, tc.subnetcount)
			Expect(err).NotTo(HaveOccurred())
			Expect(output).To(ConsistOf(tc.expected))
		})
	}
}

var (
	block = "2001:db8:1234:1a00::/56"
)

func TestParseIPv4CIDR(t *testing.T) {
	RegisterTestingT(t)

	input := []string{
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334/64",
		"2001:db8::/32",
		"193.168.3.20/7",
	}

	output, err := GetIPv4Cidrs(input)
	Expect(err).NotTo(HaveOccurred())
	Expect(output).To(HaveLen(1))
}

func TestParseIPv6CIDR(t *testing.T) {
	RegisterTestingT(t)

	input := []string{
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334/64",
		"2001:db8::/32",
		"193.168.3.20/7",
	}

	output, err := GetIPv6Cidrs(input)
	Expect(err).NotTo(HaveOccurred())
	Expect(output).To(HaveLen(2))
}

func TestSplitIntoSubnetsIPv6(t *testing.T) {
	RegisterTestingT(t)
	ip1, _, _ := net.ParseCIDR("2001:db8:1234:1a01::/64")
	ip2, _, _ := net.ParseCIDR("2001:db8:1234:1a02::/64")
	ip3, _, _ := net.ParseCIDR("2001:db8:1234:1a03::/64")
	ip4, _, _ := net.ParseCIDR("2001:db8:1234:1a04::/64")
	output, err := SplitIntoSubnetsIPv6(block, 4)
	Expect(err).NotTo(HaveOccurred())
	Expect(output).To(ConsistOf(
		&net.IPNet{
			IP:   ip1,
			Mask: net.CIDRMask(64, 128),
		},
		&net.IPNet{
			IP:   ip2,
			Mask: net.CIDRMask(64, 128),
		},
		&net.IPNet{
			IP:   ip3,
			Mask: net.CIDRMask(64, 128),
		},
		&net.IPNet{
			IP:   ip4,
			Mask: net.CIDRMask(64, 128),
		},
	))
}

func TestSplitIntoSubnetsIPv6WithFurtherSplitting(t *testing.T) {
	RegisterTestingT(t)
	ip1, _, _ := net.ParseCIDR("2001:db8:1234:1a01::/64")
	ip2, _, _ := net.ParseCIDR("2001:db8:1234:1a02::/64")
	ip3, _, _ := net.ParseCIDR("2001:db8:1234:1a03::/64")
	ip4, _, _ := net.ParseCIDR("2001:db8:1234:1a04::/64")
	output, err := SplitIntoSubnetsIPv6(block, 4)
	Expect(err).NotTo(HaveOccurred())
	Expect(output).To(ConsistOf(
		&net.IPNet{
			IP:   ip1,
			Mask: net.CIDRMask(64, 128),
		},
		&net.IPNet{
			IP:   ip2,
			Mask: net.CIDRMask(64, 128),
		},
		&net.IPNet{
			IP:   ip3,
			Mask: net.CIDRMask(64, 128),
		},
		&net.IPNet{
			IP:   ip4,
			Mask: net.CIDRMask(64, 128),
		},
	))
	output, err = SplitIntoSubnetsIPv6(output[len(output)-1].String(), 3)
	Expect(err).NotTo(HaveOccurred())
	ip1, _, _ = net.ParseCIDR("2001:db8:1234:1a05::/64")
	ip2, _, _ = net.ParseCIDR("2001:db8:1234:1a06::/64")
	ip3, _, _ = net.ParseCIDR("2001:db8:1234:1a07::/64")
	Expect(output).To(ContainElements(
		&net.IPNet{
			IP:   ip1,
			Mask: net.CIDRMask(64, 128),
		},
		&net.IPNet{
			IP:   ip2,
			Mask: net.CIDRMask(64, 128),
		},
		&net.IPNet{
			IP:   ip3,
			Mask: net.CIDRMask(64, 128),
		},
	))
}

func TestSplitIntoSubnetsIPv6HigherSubnetSplitting(t *testing.T) {
	RegisterTestingT(t)
	output, err := SplitIntoSubnetsIPv6("2001:db8:cad:ffff::/56", 6)
	Expect(err).NotTo(HaveOccurred())
	ip1, _, _ := net.ParseCIDR("2001:db8:cad:ff01::/64")
	ip2, _, _ := net.ParseCIDR("2001:db8:cad:ff02::/64")
	ip3, _, _ := net.ParseCIDR("2001:db8:cad:ff03::/64")
	ip4, _, _ := net.ParseCIDR("2001:db8:cad:ff04::/64")
	Expect(output).To(ContainElements(
		&net.IPNet{
			IP:   ip1,
			Mask: net.CIDRMask(64, 128),
		},
		&net.IPNet{
			IP:   ip2,
			Mask: net.CIDRMask(64, 128),
		},
		&net.IPNet{
			IP:   ip3,
			Mask: net.CIDRMask(64, 128),
		},
		&net.IPNet{
			IP:   ip4,
			Mask: net.CIDRMask(64, 128),
		},
	))
}

func TestSplitIntoSubnetsIPv6NoCompression(t *testing.T) {
	RegisterTestingT(t)
	output, err := SplitIntoSubnetsIPv6("2001:0db8:85a3:0010:1111:8a2e:0370:7334/56", 5)
	Expect(err).NotTo(HaveOccurred())
	ip1, _, _ := net.ParseCIDR("2001:db8:85a3:1::/64")
	ip2, _, _ := net.ParseCIDR("2001:db8:85a3:2::/64")
	ip3, _, _ := net.ParseCIDR("2001:db8:85a3:3::/64")
	ip4, _, _ := net.ParseCIDR("2001:db8:85a3:4::/64")
	ip5, _, _ := net.ParseCIDR("2001:db8:85a3:5::/64")
	Expect(output).To(ContainElements(
		&net.IPNet{
			IP:   ip1,
			Mask: net.CIDRMask(64, 128),
		},
		&net.IPNet{
			IP:   ip2,
			Mask: net.CIDRMask(64, 128),
		},
		&net.IPNet{
			IP:   ip3,
			Mask: net.CIDRMask(64, 128),
		},
		&net.IPNet{
			IP:   ip4,
			Mask: net.CIDRMask(64, 128),
		},
		&net.IPNet{
			IP:   ip5,
			Mask: net.CIDRMask(64, 128),
		},
	))
}

func TestSplitIntoSubnetsIPv6InvalidCIDR(t *testing.T) {
	RegisterTestingT(t)
	_, err := SplitIntoSubnetsIPv6("2001:db8:cad::", 60)
	Expect(err).To(MatchError(ContainSubstring("failed to parse cidr block 2001:db8:cad:: with error: invalid CIDR address: 2001:db8:cad::")))
}
