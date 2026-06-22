package openstack

import (
	"net"
	"testing"

	machinev1 "github.com/openshift/api/machine/v1"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

func TestIsSingleStackIPv6(t *testing.T) {
	tests := []struct {
		name           string
		machineNetwork []types.MachineNetworkEntry
		expected       bool
	}{
		{
			name: "single IPv6 CIDR",
			machineNetwork: []types.MachineNetworkEntry{
				{
					CIDR: ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   net.ParseIP("2001:db8::"),
							Mask: net.CIDRMask(32, 128),
						},
					},
				},
			},
			expected: true,
		},
		{
			name: "single IPv4 CIDR",
			machineNetwork: []types.MachineNetworkEntry{
				{
					CIDR: ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   net.ParseIP("192.168.1.0"),
							Mask: net.CIDRMask(24, 32),
						},
					},
				},
			},
			expected: false,
		},
		{
			name: "multiple CIDRs",
			machineNetwork: []types.MachineNetworkEntry{
				{
					CIDR: ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   net.ParseIP("2001:db8::"),
							Mask: net.CIDRMask(32, 128),
						},
					},
				},
				{
					CIDR: ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   net.ParseIP("192.168.1.0"),
							Mask: net.CIDRMask(24, 32),
						},
					},
				},
			},
			expected: false,
		},
		{
			name:           "empty machine network",
			machineNetwork: []types.MachineNetworkEntry{},
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSingleStackIPv6(tt.machineNetwork)
			if result != tt.expected {
				t.Errorf("isSingleStackIPv6() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// newMinimalInstallConfig builds a minimal *types.InstallConfig for use in unit tests.
// It sets only the fields required by generateMachineSpec so that tests stay focused.
func newMinimalInstallConfig(platform *openstack.Platform) *types.InstallConfig {
	return &types.InstallConfig{
		Platform: types.Platform{
			OpenStack: platform,
		},
	}
}

// newMinimalMachinePool builds a minimal *openstack.MachinePool for use in unit tests.
func newMinimalMachinePool(flavorName string) *openstack.MachinePool {
	return &openstack.MachinePool{
		FlavorName: flavorName,
	}
}

// TestBootstrapFlavorSelection verifies that generateMachineSpec selects the correct
// flavor for bootstrap and master roles:
//   - Bootstrap with BootstrapFlavor set  → BootstrapFlavor is used
//   - Bootstrap with BootstrapFlavor empty → control plane FlavorName is used (fallback)
//   - Master with BootstrapFlavor set     → FlavorName is used (master is unaffected)
func TestBootstrapFlavorSelection(t *testing.T) {
	const (
		controlPlaneFlavor = "m1.xlarge"
		bootstrapFlavor    = "m1.medium"
	)

	// A minimal failure domain with no zones — generateMachineSpec needs this to
	// determine the root volume AZ (none here, since RootVolume is nil).
	emptyFD := machinev1.OpenStackFailureDomain{}
	configDrive := false

	tests := []struct {
		name            string
		role            string
		bootstrapFlavor string // value placed in platform.BootstrapFlavor
		wantFlavor      string
	}{
		{
			name:            "bootstrap uses BootstrapFlavor when specified",
			role:            bootstrapRole,
			bootstrapFlavor: bootstrapFlavor,
			wantFlavor:      bootstrapFlavor,
		},
		{
			name:            "bootstrap falls back to control plane flavor when BootstrapFlavor is empty",
			role:            bootstrapRole,
			bootstrapFlavor: "",
			wantFlavor:      controlPlaneFlavor,
		},
		{
			name:            "master always uses FlavorName regardless of BootstrapFlavor",
			role:            masterRole,
			bootstrapFlavor: bootstrapFlavor,
			wantFlavor:      controlPlaneFlavor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			platform := &openstack.Platform{
				BootstrapFlavor: tt.bootstrapFlavor,
			}
			config := newMinimalInstallConfig(platform)
			mpool := newMinimalMachinePool(controlPlaneFlavor)

			spec, err := generateMachineSpec(
				"test-cluster",
				config,
				mpool,
				"rhcos",
				tt.role,
				emptyFD,
				&configDrive,
			)
			if err != nil {
				t.Fatalf("generateMachineSpec() unexpected error: %v", err)
			}
			if spec.Flavor == nil {
				t.Fatal("generateMachineSpec() returned spec with nil Flavor")
			}
			if got := *spec.Flavor; got != tt.wantFlavor {
				t.Errorf("Flavor = %q, want %q", got, tt.wantFlavor)
			}
		})
	}
}
