package openstack

import (
	"net"
	"testing"

	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/powervc"
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

// newMinimalInstallConfig returns a minimal InstallConfig for OpenStack suitable
// for use in unit tests of GenerateMachines.
func newMinimalInstallConfig(bootstrapFlavor string) *types.InstallConfig {
	return &types.InstallConfig{
		Platform: types.Platform{
			OpenStack: &openstack.Platform{
				Cloud:           "test-cloud",
				BootstrapFlavor: bootstrapFlavor,
			},
		},
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{
					CIDR: ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   net.ParseIP("10.0.0.0"),
							Mask: net.CIDRMask(24, 32),
						},
					},
				},
			},
		},
	}
}

// newMinimalPool returns a minimal MachinePool for OpenStack with the given flavor.
func newMinimalPool(flavorName string) *types.MachinePool {
	return &types.MachinePool{
		Name: "master",
		Platform: types.MachinePoolPlatform{
			OpenStack: &openstack.MachinePool{
				FlavorName: flavorName,
			},
		},
	}
}

// openStackMachineFlavorFromFiles finds the first OpenStackMachine in the given
// runtime files and returns its Flavor field.
func openStackMachineFlavorFromFiles(t *testing.T, files []*asset.RuntimeFile) string {
	t.Helper()
	for _, f := range files {
		if osm, ok := f.Object.(*capo.OpenStackMachine); ok {
			if osm.Spec.Flavor == nil {
				t.Fatal("OpenStackMachine.Spec.Flavor is nil")
			}
			return *osm.Spec.Flavor
		}
	}
	t.Fatal("no OpenStackMachine found in generated files")
	return ""
}

// newMinimalPowerVCInstallConfig returns a minimal InstallConfig for PowerVC
// suitable for use in unit tests of GenerateMachines. PowerVC is a thin
// platform that requires both Platform.PowerVC and Platform.OpenStack to be
// set; after ConvertPowerVCInstallConfig is called, Platform.OpenStack
// carries the resolved values. This helper pre-populates both fields so that
// tests can bypass the conversion step and call GenerateMachines directly.
func newMinimalPowerVCInstallConfig(bootstrapFlavor string) *types.InstallConfig {
	return &types.InstallConfig{
		Platform: types.Platform{
			// PowerVC must be non-nil so that Platform.Name() returns "powervc".
			PowerVC: &powervc.Platform{
				Cloud: "test-powervc-cloud",
			},
			// OpenStack must also be non-nil because openstackmachines.go
			// accesses config.Platform.OpenStack directly.  In a real install
			// ConvertPowerVCInstallConfig copies PowerVC data here; tests set
			// it explicitly.
			OpenStack: &openstack.Platform{
				Cloud:           "test-powervc-cloud",
				BootstrapFlavor: bootstrapFlavor,
			},
		},
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{
					CIDR: ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   net.ParseIP("10.0.0.0"),
							Mask: net.CIDRMask(24, 32),
						},
					},
				},
			},
		},
	}
}

// TestGenerateMachinesBootstrapFlavor verifies that GenerateMachines uses the
// correct flavor for both master and bootstrap roles.
//
// When called for the "bootstrap" role and the pool already has the bootstrap
// flavor resolved into FlavorName (as done by the caller in clusterapi.go),
// the resulting machine spec must reflect that flavor. When no special bootstrap
// flavor is requested the control-plane flavor is preserved.
func TestGenerateMachinesBootstrapFlavor(t *testing.T) {
	const (
		clusterID          = "test-cluster"
		controlPlaneFlavor = "m1.xlarge"
		bootstrapFlavor    = "m1.medium"
		// specialFlavor contains dots, dashes and digits — all characters that
		// are common in real OpenStack flavour names but could in theory be
		// mishandled by string manipulation.
		specialFlavor = "ocp.bootstrap-4cpu"
		osImage       = "rhcos-4.18"
	)

	tests := []struct {
		name            string
		bootstrapFlavor string // set on Platform.BootstrapFlavor; empty means not set
		role            string
		wantFlavor      string
	}{
		{
			name:            "bootstrap uses bootstrapFlavor when set",
			bootstrapFlavor: bootstrapFlavor,
			role:            bootstrapRole,
			wantFlavor:      bootstrapFlavor,
		},
		{
			name:            "bootstrap uses control plane flavor when bootstrapFlavor not set",
			bootstrapFlavor: "",
			role:            bootstrapRole,
			wantFlavor:      controlPlaneFlavor,
		},
		{
			name:            "master always uses control plane flavor regardless of bootstrapFlavor",
			bootstrapFlavor: bootstrapFlavor,
			role:            masterRole,
			wantFlavor:      controlPlaneFlavor,
		},
		{
			name:            "master uses control plane flavor when bootstrapFlavor not set",
			bootstrapFlavor: "",
			role:            masterRole,
			wantFlavor:      controlPlaneFlavor,
		},
		{
			// Flavor names that contain dots, dashes, and digits must be
			// passed through as-is without any modification.
			name:            "bootstrap flavor with special characters is used verbatim",
			bootstrapFlavor: specialFlavor,
			role:            bootstrapRole,
			wantFlavor:      specialFlavor,
		},
		{
			// Master machines must ignore the special-character bootstrapFlavor
			// and always use the control-plane flavor.
			name:            "master ignores special-character bootstrapFlavor",
			bootstrapFlavor: specialFlavor,
			role:            masterRole,
			wantFlavor:      controlPlaneFlavor,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ic := newMinimalInstallConfig(tc.bootstrapFlavor)

			// Resolve the effective flavor for this role, mirroring what clusterapi.go does.
			pool := newMinimalPool(controlPlaneFlavor)
			if tc.role == bootstrapRole {
				// Apply the same logic as clusterapi.go: resolve bootstrapFlavor from platform.
				resolvedFlavor := ic.Platform.OpenStack.ResolveBootstrapFlavor(pool.Platform.OpenStack.FlavorName)
				if resolvedFlavor != pool.Platform.OpenStack.FlavorName {
					// Clone the pool so master pool is unaffected.
					bootstrapMpool := *pool.Platform.OpenStack
					bootstrapMpool.FlavorName = resolvedFlavor
					pool = newMinimalPool(resolvedFlavor)
				}
			}

			files, err := GenerateMachines(clusterID, ic, pool, osImage, tc.role)
			if err != nil {
				t.Fatalf("GenerateMachines() returned error: %v", err)
			}
			if len(files) == 0 {
				t.Fatal("GenerateMachines() returned no files")
			}

			gotFlavor := openStackMachineFlavorFromFiles(t, files)
			if gotFlavor != tc.wantFlavor {
				t.Errorf("OpenStackMachine.Spec.Flavor = %q, want %q", gotFlavor, tc.wantFlavor)
			}
		})
	}
}

// TestGenerateMachinesBootstrapFlavorPowerVC verifies that the bootstrap flavor
// feature works for the PowerVC platform. PowerVC shares the OpenStack CAPI
// code path: after ConvertPowerVCInstallConfig runs, Platform.OpenStack is
// populated and openstackmachines.GenerateMachines is called for both master
// and bootstrap roles. BootstrapFlavor lives on Platform.OpenStack (not on
// Platform.PowerVC) so the same resolution logic applies.
func TestGenerateMachinesBootstrapFlavorPowerVC(t *testing.T) {
	const (
		clusterID          = "powervc-cluster"
		controlPlaneFlavor = "powervc-large"
		bootstrapFlavor    = "powervc-medium"
		specialFlavor      = "ocp.bootstrap-4cpu"
		osImage            = "rhcos-4.18"
	)

	tests := []struct {
		name            string
		bootstrapFlavor string // set on Platform.OpenStack.BootstrapFlavor
		role            string
		wantFlavor      string
	}{
		{
			name:            "PowerVC bootstrap uses bootstrapFlavor when set",
			bootstrapFlavor: bootstrapFlavor,
			role:            bootstrapRole,
			wantFlavor:      bootstrapFlavor,
		},
		{
			name:            "PowerVC bootstrap falls back to control plane flavor when not set",
			bootstrapFlavor: "",
			role:            bootstrapRole,
			wantFlavor:      controlPlaneFlavor,
		},
		{
			name:            "PowerVC master always uses control plane flavor",
			bootstrapFlavor: bootstrapFlavor,
			role:            masterRole,
			wantFlavor:      controlPlaneFlavor,
		},
		{
			name:            "PowerVC bootstrap flavor with special characters is used verbatim",
			bootstrapFlavor: specialFlavor,
			role:            bootstrapRole,
			wantFlavor:      specialFlavor,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Use a PowerVC install config (Platform.Name() == "powervc").
			ic := newMinimalPowerVCInstallConfig(tc.bootstrapFlavor)

			// Resolve the effective flavor for this role, mirroring clusterapi.go.
			pool := newMinimalPool(controlPlaneFlavor)
			if tc.role == bootstrapRole {
				resolvedFlavor := ic.Platform.OpenStack.ResolveBootstrapFlavor(pool.Platform.OpenStack.FlavorName)
				if resolvedFlavor != pool.Platform.OpenStack.FlavorName {
					pool = newMinimalPool(resolvedFlavor)
				}
			}

			files, err := GenerateMachines(clusterID, ic, pool, osImage, tc.role)
			if err != nil {
				t.Fatalf("GenerateMachines() returned error: %v", err)
			}
			if len(files) == 0 {
				t.Fatal("GenerateMachines() returned no files")
			}

			gotFlavor := openStackMachineFlavorFromFiles(t, files)
			if gotFlavor != tc.wantFlavor {
				t.Errorf("OpenStackMachine.Spec.Flavor = %q, want %q", gotFlavor, tc.wantFlavor)
			}
		})
	}
}
