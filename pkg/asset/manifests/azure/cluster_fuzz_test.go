package azure

import (
	"math/big"
	"net"
	"testing"

	aznetwork "github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
)

func FuzzGetIPWithinCIDR(f *testing.F) {
	f.Add("10.0.0.0/16", "10.0.1.0/24", "10.0.0.100", false)
	f.Add("10.0.0.252/30", "", "192.168.1.1", false)
	f.Add("10.0.0.0/16", "", "192.168.1.1", false)
	f.Add("10.0.0.0/16", "10.0.1.0/24", "", false)
	f.Add("not-a-cidr", "", "10.0.0.5", false)
	// If every subnet is nil or has an empty IP, fallback selection can index an empty network.
	f.Add("10.0.0.0/16", "10.0.1.0/24", "192.168.1.1", true)

	f.Fuzz(func(t *testing.T, cidr1, cidr2, ip string, firstNil bool) {
		subnets := subnetNets(cidr1, cidr2, firstNil)
		got := getIPWithinCIDR(subnets, ip)

		parsedInput := net.ParseIP(ip)
		for _, subnet := range subnets {
			if subnet != nil && len(subnet.IP) > 0 && parsedInput != nil && subnet.Contains(parsedInput) {
				if got != ip {
					t.Fatalf("getIPWithinCIDR(%q, %q, firstNil=%t) = %q, want the in-subnet IP %q", cidr1, cidr2, firstNil, got, ip)
				}
				return
			}
		}

		if got == "" {
			if ip != "" && hasFallbackIP(subnets) {
				t.Fatalf("getIPWithinCIDR(%q, %q, firstNil=%t) returned empty despite an available fallback IP", cidr1, cidr2, firstNil)
			}
			return
		}
		parsed := net.ParseIP(got)
		if parsed == nil {
			t.Fatalf("getIPWithinCIDR(%q, %q, firstNil=%t) returned invalid IP %q", cidr1, cidr2, firstNil, got)
		}
		for _, subnet := range subnets {
			if subnet != nil && len(subnet.IP) > 0 && subnet.Contains(parsed) {
				return
			}
		}
		t.Fatalf("getIPWithinCIDR(%q, %q, firstNil=%t) returned %q outside every subnet", cidr1, cidr2, firstNil, got)
	})
}

func hasFallbackIP(subnets []*net.IPNet) bool {
	for _, subnet := range subnets {
		if subnet == nil || len(subnet.IP) == 0 {
			continue
		}

		// Use big.Int to independently check whether the first usable subnet
		// contains the Azure-reserved-address offset (+4).
		candidate := new(big.Int).SetBytes(subnet.IP)
		candidate.Add(candidate, big.NewInt(4))
		if candidate.BitLen() > len(subnet.IP)*8 {
			return false
		}
		candidateIP := make(net.IP, len(subnet.IP))
		candidate.FillBytes(candidateIP)
		return subnet.Contains(candidateIP)
	}
	return false
}

func FuzzSubnetAddressPrefixes(f *testing.F) {
	f.Add("10.0.0.0/16", "10.0.1.0/24", true, true)
	f.Add("10.0.0.0/16", "not-a-cidr", true, true)
	f.Add("not-a-cidr", "10.0.0.0/16", true, false)
	f.Add("", "", false, false)
	f.Add("2001:db8::/64", "10.0.0.0/16", false, true)

	f.Fuzz(func(t *testing.T, listPrefix, singlePrefix string, hasList, hasSingle bool) {
		subnet := &aznetwork.Subnet{SubnetPropertiesFormat: &aznetwork.SubnetPropertiesFormat{}}
		if hasList {
			prefixes := []string{listPrefix}
			subnet.AddressPrefixes = &prefixes
		}
		if hasSingle {
			subnet.AddressPrefix = &singlePrefix
		}
		got, err := getSubnetAddressPrefixes(subnet)
		if err != nil && len(got) > 2 {
			t.Fatalf("parsed %d prefixes after an error", len(got))
		}
		for _, prefix := range got {
			if prefix == nil || prefix.IP == nil || prefix.Mask == nil {
				t.Fatalf("parsed prefix is incomplete: %#v", prefix)
			}
		}
		wantN, wantErr := expectedPrefixParse(listPrefix, singlePrefix, hasList, hasSingle)
		if (err != nil) != wantErr || len(got) != wantN {
			t.Fatalf("prefixes list=%q single=%q err=%v len=%d, want err=%t len=%d", listPrefix, singlePrefix, err, len(got), wantErr, wantN)
		}
	})
}

func subnetNets(cidr1, cidr2 string, firstNil bool) []*net.IPNet {
	var subnets []*net.IPNet
	if firstNil {
		subnets = append(subnets, nil, &net.IPNet{})
	}
	for _, cidr := range []string{cidr1, cidr2} {
		if _, network, err := net.ParseCIDR(cidr); err == nil {
			subnets = append(subnets, network)
		}
	}
	return subnets
}

func expectedPrefixParse(listPrefix, singlePrefix string, hasList, hasSingle bool) (int, bool) {
	n := 0
	if hasList {
		if _, _, err := net.ParseCIDR(listPrefix); err != nil {
			return n, true
		}
		n++
	}
	if hasSingle {
		if _, _, err := net.ParseCIDR(singlePrefix); err != nil {
			return n, true
		}
		n++
	}
	return n, false
}
