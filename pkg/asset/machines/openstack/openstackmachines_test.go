package openstack

import (
	"net"
	"testing"

	"k8s.io/utils/ptr"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"

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

// newTestInstallConfig returns a minimal *types.InstallConfig configured for
// OpenStack that is sufficient to exercise GenerateMachines.
func newTestInstallConfig(bootstrapFlavor string) *types.InstallConfig {
	return &types.InstallConfig{
		Networking: &types.Networking{
			MachineNetwork: []types.MachineNetworkEntry{
				{
					CIDR: ipnet.IPNet{
						IPNet: net.IPNet{
							IP:   net.ParseIP("10.0.0.0"),
							Mask: net.CIDRMask(16, 32),
						},
					},
				},
			},
		},
		Platform: types.Platform{
			OpenStack: &openstack.Platform{
				Cloud:           "mycloud",
				BootstrapFlavor: bootstrapFlavor,
				APIVIPs:         []string{"10.0.0.5"},
				IngressVIPs:     []string{"10.0.0.7"},
			},
		},
	}
}

// newTestMachinePool returns a *types.MachinePool configured for OpenStack
// with the given control plane flavor name and replica count.
func newTestMachinePool(flavorName string, replicas int64) *types.MachinePool {
	return &types.MachinePool{
		Name:     "master",
		Replicas: ptr.To(replicas),
		Platform: types.MachinePoolPlatform{
			OpenStack: &openstack.MachinePool{
				FlavorName: flavorName,
			},
		},
	}
}

// TestGenerateMachinesBootstrapFlavor tests that GenerateMachines correctly
// assigns the bootstrap flavor to the bootstrap machine and leaves control
// plane machines using the pool's flavor.
func TestGenerateMachinesBootstrapFlavor(t *testing.T) {
	const (
		clusterID          = "test-cluster"
		osImage            = "rhcos-latest"
		controlPlaneFlavor = "m1.large"
		bootstrapFlavor    = "m1.xlarge"
	)

	tests := []struct {
		name                   string
		bootstrapFlavor        string
		controlPlaneFlavor     string
		role                   string
		wantFlavor             string
		wantErr                bool
	}{
		{
			name:               "bootstrap machine uses explicit bootstrap flavor",
			bootstrapFlavor:    bootstrapFlavor,
			controlPlaneFlavor: controlPlaneFlavor,
			role:               bootstrapRole,
			wantFlavor:         bootstrapFlavor,
		},
		{
			name:               "bootstrap machine inherits control plane flavor when bootstrap flavor not set",
			bootstrapFlavor:    "",
			controlPlaneFlavor: controlPlaneFlavor,
			role:               bootstrapRole,
			wantFlavor:         controlPlaneFlavor,
		},
		{
			name:               "control plane machines use pool flavor regardless of bootstrap flavor",
			bootstrapFlavor:    bootstrapFlavor,
			controlPlaneFlavor: controlPlaneFlavor,
			role:               masterRole,
			wantFlavor:         controlPlaneFlavor,
		},
		{
			name:               "control plane machines unaffected when bootstrap flavor is empty",
			bootstrapFlavor:    "",
			controlPlaneFlavor: controlPlaneFlavor,
			role:               masterRole,
			wantFlavor:         controlPlaneFlavor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := newTestInstallConfig(tt.bootstrapFlavor)
			pool := newTestMachinePool(tt.controlPlaneFlavor, 1)

			files, err := GenerateMachines(clusterID, config, pool, osImage, tt.role)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(files) == 0 {
				t.Fatalf("expected at least one file, got none")
			}

			// The first file in the result is always the OpenStackMachine infra object.
			// Extract and check its flavor.
			osMachine, ok := files[0].Object.(*capo.OpenStackMachine)
			if !ok {
				t.Fatalf("expected first object to be *capo.OpenStackMachine, got %T", files[0].Object)
			}

			if osMachine.Spec.Flavor == nil {
				t.Fatalf("OpenStackMachine.Spec.Flavor is nil")
			}
			if got := *osMachine.Spec.Flavor; got != tt.wantFlavor {
				t.Errorf("OpenStackMachine.Spec.Flavor = %q, want %q", got, tt.wantFlavor)
			}
		})
	}
}

// TestGenerateMachinesControlPlaneNotAffectedByBootstrapFlavor tests that
// generating multiple control plane machines does not apply the bootstrap
// flavor to any of them, even when bootstrapFlavor is set.
func TestGenerateMachinesControlPlaneNotAffectedByBootstrapFlavor(t *testing.T) {
	const (
		clusterID          = "test-cluster"
		osImage            = "rhcos-latest"
		controlPlaneFlavor = "m1.large"
		replicas           = int64(3)
	)

	config := newTestInstallConfig("m1.xlarge") // bootstrap flavor set
	pool := newTestMachinePool(controlPlaneFlavor, replicas)

	files, err := GenerateMachines(clusterID, config, pool, osImage, masterRole)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Each replica produces 2 files: an OpenStackMachine and a CAPI Machine.
	// Total = replicas * 2.
	expectedFiles := int(replicas) * 2
	if len(files) != expectedFiles {
		t.Fatalf("expected %d files, got %d", expectedFiles, len(files))
	}

	// Every OpenStackMachine (even-indexed files) must use the control plane flavor.
	for i := int64(0); i < replicas; i++ {
		fileIdx := i * 2 // OpenStackMachine is at even index
		osMachine, ok := files[fileIdx].Object.(*capo.OpenStackMachine)
		if !ok {
			t.Errorf("file[%d]: expected *capo.OpenStackMachine, got %T", fileIdx, files[fileIdx].Object)
			continue
		}
		if osMachine.Spec.Flavor == nil {
			t.Errorf("file[%d]: OpenStackMachine.Spec.Flavor is nil", fileIdx)
			continue
		}
		if got := *osMachine.Spec.Flavor; got != controlPlaneFlavor {
			t.Errorf("file[%d]: control plane machine flavor = %q, want %q", fileIdx, got, controlPlaneFlavor)
		}
	}
}
