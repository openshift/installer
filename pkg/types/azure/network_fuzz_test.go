package azure

import (
	"net"
	"testing"

	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/installer/pkg/ipnet"
)

func FuzzAddressFamilySubnets(f *testing.F) {
	f.Add("10.0.0.0/16", "", "", "", 3)
	f.Add("10.0.0.0/32", "", "", "", 2)
	f.Add("10.0.0.0/16", "10.1.0.0/16", "2001:db8::/56", "", 4)
	f.Add("10.0.0.0/16", "2001:db8::/56", "2001:db8:1::/56", "2001:db8:2::/56", 2)
	f.Add("2001:db8::/56", "", "", "", 2)
	f.Add("10.0.0.0/16", "10.1.0.0/16", "10.2.0.0/16", "", 0)
	f.Add("10.0.0.0/16", "", "", "", -1)
	// No IPv4 CIDR leaves an empty subnet list. Splitting that used to panic.
	f.Add("", "", "", "", 3)

	f.Fuzz(func(t *testing.T, cidr1, cidr2, cidr3, cidr4 string, numSubnets int) {
		cidrs := parseCIDRs(cidr1, cidr2, cidr3, cidr4)
		subnets, err := GetAddressFamilySubnets(cidrs)
		if err != nil {
			if subnets.length != 0 || len(subnets.addressFamilySubnets) != 0 {
				t.Fatalf("GetAddressFamilySubnets error %v left %d subnets", err, subnets.length)
			}
		} else {
			assertAddressFamilySubnets(t, subnets)
		}

		splitCount := numSubnets
		if splitCount > 0 {
			splitCount = (splitCount % 8) + 1
		}
		lengthBefore := subnets.length
		var controlPlane string
		if lengthBefore > 0 && subnets.addressFamilySubnets[0].IPv4Subnet != nil {
			controlPlane = subnets.addressFamilySubnets[0].IPv4Subnet.String()
		}
		splitErr := subnets.SplitComputeSubnet(splitCount)
		if numSubnets <= 0 || lengthBefore == 0 {
			if splitErr != nil || subnets.length != lengthBefore {
				t.Fatalf("split %d on length %d = %v, length %d", splitCount, lengthBefore, splitErr, subnets.length)
			}
			return
		}
		if splitErr != nil {
			if subnets.length != lengthBefore {
				t.Fatalf("failed split changed length from %d to %d: %v", lengthBefore, subnets.length, splitErr)
			}
			return
		}
		if subnets.length < 1 || subnets.addressFamilySubnets[0].IPv4Subnet == nil || subnets.addressFamilySubnets[0].IPv4Subnet.String() != controlPlane {
			t.Fatalf("split dropped control plane %q, length %d", controlPlane, subnets.length)
		}
		if subnets.length > lengthBefore+splitCount {
			t.Fatalf("split length %d exceeds bound %d", subnets.length, lengthBefore+splitCount)
		}
	})
}

func parseCIDRs(values ...string) []ipnet.IPNet {
	var cidrs []ipnet.IPNet
	for _, value := range values {
		if value == "" {
			continue
		}
		_, network, err := net.ParseCIDR(value)
		if err != nil {
			continue
		}
		cidrs = append(cidrs, ipnet.IPNet{IPNet: *network})
	}
	return cidrs
}

func assertAddressFamilySubnets(t *testing.T, subnets AddressFamilySubnets) {
	t.Helper()
	if subnets.length < 2 || subnets.length != len(subnets.addressFamilySubnets) || subnets.ipv4Count != subnets.length {
		t.Fatalf("length=%d slice=%d ipv4=%d", subnets.length, len(subnets.addressFamilySubnets), subnets.ipv4Count)
	}
	if subnets.ipv6Count > subnets.ipv4Count {
		t.Fatalf("ipv6 count %d exceeds ipv4 count %d", subnets.ipv6Count, subnets.ipv4Count)
	}
	for i, subnet := range subnets.addressFamilySubnets {
		if subnet.IPv4Subnet == nil || subnet.IPv4Subnet.IP == nil || subnet.IPv4Subnet.Mask == nil {
			t.Fatalf("subnet %d is missing an IPv4 network", i)
		}
	}
	if subnets.addressFamilySubnets[0].SubnetRole == nil || *subnets.addressFamilySubnets[0].SubnetRole != capz.SubnetControlPlane {
		t.Fatalf("first subnet role = %v, want control plane", subnets.addressFamilySubnets[0].SubnetRole)
	}
	if subnets.addressFamilySubnets[1].SubnetRole == nil || *subnets.addressFamilySubnets[1].SubnetRole != capz.SubnetNode {
		t.Fatalf("second subnet role = %v, want compute", subnets.addressFamilySubnets[1].SubnetRole)
	}
}
