package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

func validPlatform() *openstack.Platform {
	return &openstack.Platform{
		Cloud:           "test-cloud",
		ExternalNetwork: "test-network",
	}
}

func validNetworking() *types.Networking {
	return &types.Networking{
		NetworkType: "OVNKubernetes",
		MachineNetwork: []types.MachineNetworkEntry{{
			CIDR: *ipnet.MustParseCIDR("10.0.0.0/16"),
		}},
	}
}

func TestSetPlatformDefaults(t *testing.T) {
	cases := []struct {
		name                string
		platform            *openstack.Platform
		config              *types.InstallConfig
		networking          *types.Networking
		expectedLB          configv1.PlatformLoadBalancerType
		expectedAPIVIPs     []string
		expectedIngressVIPs []string
	}{
		{
			name: "No load balancer provided",
			platform: func() *openstack.Platform {
				p := validPlatform()
				return p
			}(),
			networking:          validNetworking(),
			expectedAPIVIPs:     []string{"10.0.0.5"},
			expectedIngressVIPs: []string{"10.0.0.7"},
			expectedLB:          "OpenShiftManagedDefault",
		},
		{
			name: "Default Openshift Managed load balancer VIPs provided",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.LoadBalancer = &configv1.OpenStackPlatformLoadBalancer{
					Type: configv1.LoadBalancerTypeOpenShiftManagedDefault,
				}
				p.APIVIPs = []string{"10.0.0.2"}
				p.IngressVIPs = []string{"10.0.0.3"}
				return p
			}(),
			networking:          validNetworking(),
			expectedAPIVIPs:     []string{"10.0.0.2"},
			expectedIngressVIPs: []string{"10.0.0.3"},
			expectedLB:          "OpenShiftManagedDefault",
		},
		{
			name: "Default Openshift Managed load balancer no VIPs provided",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.LoadBalancer = &configv1.OpenStackPlatformLoadBalancer{
					Type: configv1.LoadBalancerTypeOpenShiftManagedDefault,
				}
				return p
			}(),
			networking:          validNetworking(),
			expectedAPIVIPs:     []string{"10.0.0.5"},
			expectedIngressVIPs: []string{"10.0.0.7"},
			expectedLB:          "OpenShiftManagedDefault",
		},
		{
			name: "User managed load balancer VIPs provided",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.LoadBalancer = &configv1.OpenStackPlatformLoadBalancer{
					Type: configv1.LoadBalancerTypeUserManaged,
				}
				p.APIVIPs = []string{"192.168.100.10"}
				p.IngressVIPs = []string{"192.168.100.11"}
				return p
			}(),
			networking:          validNetworking(),
			expectedAPIVIPs:     []string{"192.168.100.10"},
			expectedIngressVIPs: []string{"192.168.100.11"},
			expectedLB:          "UserManaged",
		},
		{
			name: "User managed load balancer no VIPs provided",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.LoadBalancer = &configv1.OpenStackPlatformLoadBalancer{
					Type: configv1.LoadBalancerTypeUserManaged,
				}
				return p
			}(),
			networking:          validNetworking(),
			expectedAPIVIPs:     []string(nil),
			expectedIngressVIPs: []string(nil),
			expectedLB:          "UserManaged",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			SetPlatformDefaults(tc.platform, tc.networking)
			assert.Equal(t, tc.expectedLB, tc.platform.LoadBalancer.Type, "unexpected loadbalancer")
			assert.Equal(t, tc.expectedAPIVIPs, tc.platform.APIVIPs, "unexpected APIVIPs")
			assert.Equal(t, tc.expectedIngressVIPs, tc.platform.IngressVIPs, "unexpected IngressVIPs")
		})
	}
}

func TestBootstrapFlavorName(t *testing.T) {
	cases := []struct {
		name           string
		platform       *openstack.Platform
		controlPlane   *openstack.MachinePool
		expectedFlavor string
	}{
		{
			name: "BootstrapFlavor is set",
			platform: &openstack.Platform{
				BootstrapFlavor: "bootstrap-flavor",
				DefaultMachinePlatform: &openstack.MachinePool{
					FlavorName: "default-flavor",
				},
			},
			controlPlane: &openstack.MachinePool{
				FlavorName: "control-plane-flavor",
			},
			expectedFlavor: "bootstrap-flavor",
		},
		{
			name: "BootstrapFlavor empty, use controlPlane flavor",
			platform: &openstack.Platform{
				DefaultMachinePlatform: &openstack.MachinePool{
					FlavorName: "default-flavor",
				},
			},
			controlPlane: &openstack.MachinePool{
				FlavorName: "control-plane-flavor",
			},
			expectedFlavor: "control-plane-flavor",
		},
		{
			name: "BootstrapFlavor and controlPlane empty, use DefaultMachinePlatform flavor",
			platform: &openstack.Platform{
				DefaultMachinePlatform: &openstack.MachinePool{
					FlavorName: "default-flavor",
				},
			},
			controlPlane:   nil,
			expectedFlavor: "default-flavor",
		},
		{
			name: "controlPlane with empty FlavorName, use DefaultMachinePlatform flavor",
			platform: &openstack.Platform{
				DefaultMachinePlatform: &openstack.MachinePool{
					FlavorName: "default-flavor",
				},
			},
			controlPlane: &openstack.MachinePool{
				FlavorName: "",
			},
			expectedFlavor: "default-flavor",
		},
		{
			name:           "No flavor configured anywhere",
			platform:       &openstack.Platform{},
			controlPlane:   nil,
			expectedFlavor: "",
		},
		{
			name: "DefaultMachinePlatform is nil, controlPlane is nil",
			platform: &openstack.Platform{
				DefaultMachinePlatform: nil,
			},
			controlPlane:   nil,
			expectedFlavor: "",
		},
		{
			name: "DefaultMachinePlatform with empty FlavorName",
			platform: &openstack.Platform{
				DefaultMachinePlatform: &openstack.MachinePool{
					FlavorName: "",
				},
			},
			controlPlane:   nil,
			expectedFlavor: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := BootstrapFlavorName(tc.platform, tc.controlPlane)
			assert.Equal(t, tc.expectedFlavor, result, "unexpected bootstrap flavor")
		})
	}
}
