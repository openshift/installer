package azure

import (
	"net"
	"testing"

	aznetwork "github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/network/mgmt/network"
)

func FuzzGetIPWithinCIDR(f *testing.F) {
	f.Add("10.0.0.0/16", "10.0.1.0/24", "10.0.0.100", false)
	f.Add("10.0.0.0/16", "", "192.168.1.1", false)
	f.Add("10.0.0.0/16", "10.0.1.0/24", "", false)
	f.Add("not-a-cidr", "", "10.0.0.5", false)
	// A nil first subnet used to be dereferenced after the scan skipped it.
	f.Add("10.0.0.0/16", "10.0.1.0/24", "192.168.1.1", true)

	f.Fuzz(func(t *testing.T, cidr1, cidr2, ip string, firstNil bool) {
		subnets := subnetNets(cidr1, cidr2, firstNil)
		got := getIPWithinCIDR(subnets, ip)
		want := expectedIPWithinCIDR(subnets, ip)
		if got != want {
			t.Fatalf("getIPWithinCIDR(%q, %q, firstNil=%t) ip %q = %q, want %q", cidr1, cidr2, firstNil, ip, got, want)
		}
	})
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
		subnets = append(subnets, nil)
	}
	for _, cidr := range []string{cidr1, cidr2} {
		if _, network, err := net.ParseCIDR(cidr); err == nil {
			subnets = append(subnets, network)
		}
	}
	return subnets
}

func expectedIPWithinCIDR(subnets []*net.IPNet, ip string) string {
	if len(subnets) == 0 || ip == "" {
		return ""
	}
	parsed := net.ParseIP(ip)
	var network net.IP
	for _, subnet := range subnets {
		if subnet == nil || len(subnet.IP) == 0 {
			continue
		}
		if network == nil {
			network = subnet.IP
		}
		if subnet.Contains(parsed) {
			return ip
		}
	}
	if network == nil {
		return ""
	}
	next := make(net.IP, len(network))
	copy(next, network)
	next[len(next)-1] += 4
	return next.String()
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
