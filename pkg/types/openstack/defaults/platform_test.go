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
		name                    string
		platform                *openstack.Platform
		config                  *types.InstallConfig
		networking              *types.Networking
		expectedLB              configv1.PlatformLoadBalancerType
		expectedAPIVIPs         []string
		expectedIngressVIPs     []string
		expectedBootstrapFlavor string
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
		{
			name: "Bootstrap flavor is defaulted from DefaultMachinePlatform",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.DefaultMachinePlatform = &openstack.MachinePool{
					FlavorName: "m1.large",
				}
				return p
			}(),
			networking:              validNetworking(),
			expectedBootstrapFlavor: "m1.large",
			expectedAPIVIPs:         []string{"10.0.0.5"},
			expectedIngressVIPs:     []string{"10.0.0.7"},
			expectedLB:              "OpenShiftManagedDefault",
		},
		{
			name: "Bootstrap flavor is not defaulted if already set",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.BootstrapFlavor = "m1.xlarge"
				p.DefaultMachinePlatform = &openstack.MachinePool{
					FlavorName: "m1.large",
				}
				return p
			}(),
			networking:              validNetworking(),
			expectedBootstrapFlavor: "m1.xlarge",
			expectedAPIVIPs:         []string{"10.0.0.5"},
			expectedIngressVIPs:     []string{"10.0.0.7"},
			expectedLB:              "OpenShiftManagedDefault",
		},
		{
			name: "Bootstrap flavor is defaulted from DefaultMachinePlatform when it contains only whitespaces",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.BootstrapFlavor = "    "
				p.DefaultMachinePlatform = &openstack.MachinePool{
					FlavorName: "m1.large",
				}
				return p
			}(),
			networking:              validNetworking(),
			expectedBootstrapFlavor: "m1.large",
			expectedAPIVIPs:         []string{"10.0.0.5"},
			expectedIngressVIPs:     []string{"10.0.0.7"},
			expectedLB:              "OpenShiftManagedDefault",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			SetPlatformDefaults(tc.platform, tc.networking)
			assert.Equal(t, tc.expectedLB, tc.platform.LoadBalancer.Type, "unexpected loadbalancer")
			assert.Equal(t, tc.expectedAPIVIPs, tc.platform.APIVIPs, "unexpected APIVIPs")
			assert.Equal(t, tc.expectedIngressVIPs, tc.platform.IngressVIPs, "unexpected IngressVIPs")
			if tc.expectedBootstrapFlavor != "" {
				assert.Equal(t, tc.expectedBootstrapFlavor, tc.platform.BootstrapFlavor, "unexpected bootstrap flavor")
			}
		})
	}
}
