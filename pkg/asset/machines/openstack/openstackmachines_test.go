package openstack

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"

	machinev1 "github.com/openshift/api/machine/v1"
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

func TestGenerateMachineSpecTags(t *testing.T) {
	const clusterID = "test-infra-id"
	tests := []struct {
		name     string
		powervc  bool
		expected []string
	}{
		{
			name:     "openstack sets cluster ID tag",
			expected: []string{"openshiftClusterID=" + clusterID},
		},
		{
			name:     "powervc sets no tags",
			powervc:  true,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &types.InstallConfig{
				Platform: types.Platform{
					OpenStack: &openstack.Platform{},
				},
			}
			if tt.powervc {
				config.Platform.PowerVC = &powervc.Platform{}
			}
			spec, err := generateMachineSpec(clusterID, config, &openstack.MachinePool{}, "image", masterRole, machinev1.OpenStackFailureDomain{}, ptr.To(false))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assert.Equal(t, tt.expected, spec.Tags)
		})
	}
}

func TestGenerateMachineSpecServerGroup(t *testing.T) {
	const clusterID = "test-infra-id"
	tests := []struct {
		name     string
		powervc  bool
		role     string
		expected string
	}{
		{
			name:     "openstack master references the master server group",
			role:     masterRole,
			expected: clusterID + "-" + masterRole,
		},
		{
			name: "openstack bootstrap has no server group",
			role: bootstrapRole,
		},
		{
			name:    "powervc master has no server group",
			powervc: true,
			role:    masterRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &types.InstallConfig{
				Platform: types.Platform{
					OpenStack: &openstack.Platform{},
				},
			}
			if tt.powervc {
				config.Platform.PowerVC = &powervc.Platform{}
			}
			spec, err := generateMachineSpec(clusterID, config, &openstack.MachinePool{}, "image", tt.role, machinev1.OpenStackFailureDomain{}, ptr.To(false))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expected == "" {
				assert.Nil(t, spec.ServerGroup)
				return
			}
			if assert.NotNil(t, spec.ServerGroup) && assert.NotNil(t, spec.ServerGroup.Filter) && assert.NotNil(t, spec.ServerGroup.Filter.Name) {
				assert.Equal(t, tt.expected, *spec.ServerGroup.Filter.Name)
			}
		})
	}
}
