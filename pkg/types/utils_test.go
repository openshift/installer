package types

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/ipnet"
)

// TestStringsToIPs tests the StringsToIPs function.
func TestStringsToIPs(t *testing.T) {
	testcases := []struct {
		ips      []string
		expected []configv1.IP
	}{
		{
			[]string{"10.0.0.1", "10.0.0.2"},
			[]configv1.IP{"10.0.0.1", "10.0.0.2"},
		},
		{
			[]string{},
			[]configv1.IP{},
		},
		{
			[]string{"fe80:1:2:3::"},
			[]configv1.IP{"fe80:1:2:3::"},
		},
	}

	for _, tc := range testcases {
		res := StringsToIPs(tc.ips)
		assert.Equal(t, tc.expected, res, "conversion failed")
	}
}

// TestMachineNetworksToCIDRs tests the MachineNetworksToCIDRs function.
func TestMachineNetworksToCIDRs(t *testing.T) {
	testcases := []struct {
		networks []MachineNetworkEntry
		expected []configv1.CIDR
	}{
		{
			[]MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("10.0.0.1/32")},
				{CIDR: *ipnet.MustParseCIDR("10.0.0.2/32")},
			},
			[]configv1.CIDR{"10.0.0.1/32", "10.0.0.2/32"},
		},
		{
			[]MachineNetworkEntry{},
			[]configv1.CIDR{},
		},
		{
			[]MachineNetworkEntry{
				{CIDR: *ipnet.MustParseCIDR("fe80:1:2:3::/128")},
			},
			[]configv1.CIDR{"fe80:1:2:3::/128"},
		},
	}

	for _, tc := range testcases {
		res := MachineNetworksToCIDRs(tc.networks)
		assert.Equal(t, tc.expected, res, "conversion failed")
	}
}
