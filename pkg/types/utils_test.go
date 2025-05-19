package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types/azure"
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

func TestCreateAzureIdentity(t *testing.T) {
	baseInstallConfig := func() *InstallConfig {
		return &InstallConfig{
			Compute: []MachinePool{
				{
					Platform: MachinePoolPlatform{
						Azure: &azure.MachinePool{},
					},
				},
			},
			ControlPlane: &MachinePool{
				Platform: MachinePoolPlatform{
					Azure: &azure.MachinePool{},
				},
			},
			Platform: Platform{
				Azure: &azure.Platform{},
			},
		}
	}

	cases := []struct {
		name           string
		installConfig  *InstallConfig
		expectedResult bool
	}{
		{
			name: "Create Identities with minimal install config (default)",
			installConfig: &InstallConfig{
				Platform: Platform{
					Azure: &azure.Platform{},
				},
			},
			expectedResult: true,
		},
		{
			name: "Create Identities by Default",
			installConfig: func() *InstallConfig {
				return baseInstallConfig()
			}(),
			expectedResult: true,
		},
		{
			name: "Don't create identities when default machine pool identity is none",
			installConfig: func() *InstallConfig {
				ic := baseInstallConfig()
				ic.Platform.Azure.DefaultMachinePlatform = &azure.MachinePool{
					Identity: &azure.VMIdentity{
						Type: capz.VMIdentityNone,
					},
				}
				return ic
			}(),
			expectedResult: false,
		},
		{
			name: "create identities when identity type is user assigned but none are supplied",
			installConfig: func() *InstallConfig {
				ic := baseInstallConfig()
				ic.Platform.Azure.DefaultMachinePlatform = &azure.MachinePool{
					Identity: &azure.VMIdentity{
						Type: capz.VMIdentityUserAssigned,
					},
				}
				return ic
			}(),
			expectedResult: true,
		},
		{
			name: "create identities when byo control plane identity",
			installConfig: func() *InstallConfig {
				ic := baseInstallConfig()
				ic.Platform.Azure.DefaultMachinePlatform = &azure.MachinePool{
					Identity: &azure.VMIdentity{
						Type: capz.VMIdentityUserAssigned,
					},
				}
				ic.ControlPlane.Platform.Azure.Identity = &azure.VMIdentity{
					Type: capz.VMIdentityUserAssigned,
					UserAssignedIdentities: []azure.UserAssignedIdentity{
						{
							Name:          "test",
							ResourceGroup: "test",
							Subscription:  "test",
						},
					},
				}
				return ic
			}(),
			expectedResult: true,
		},
		{
			name: "Do not create identities when byo identities",
			installConfig: func() *InstallConfig {
				ic := baseInstallConfig()
				ic.Platform.Azure.DefaultMachinePlatform = &azure.MachinePool{
					Identity: &azure.VMIdentity{
						Type: capz.VMIdentityUserAssigned,
						UserAssignedIdentities: []azure.UserAssignedIdentity{
							{
								Name:          "test",
								ResourceGroup: "test",
								Subscription:  "test",
							},
						},
					},
				}
				return ic
			}(),
			expectedResult: false,
		},
	}
	for _, tc := range cases {
		res := tc.installConfig.CreateAzureIdentity()
		assert.Equal(t, tc.expectedResult, res, tc.name)
	}
}
