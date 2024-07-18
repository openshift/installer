package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

func validPlatform() *openstack.Platform {
	return &openstack.Platform{
		Cloud:           "test-cloud",
		ExternalNetwork: "test-network",
		DefaultMachinePlatform: &openstack.MachinePool{
			FlavorName: "test-flavor",
		},
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

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name                  string
		config                *types.InstallConfig
		platform              *openstack.Platform
		networking            *types.Networking
		noClouds              bool
		noNetworks            bool
		noFlavors             bool
		validMachinesSubnet   bool
		invalidMachinesSubnet bool
		valid                 bool
		expectedError         string
	}{
		{
			name:       "minimal",
			platform:   validPlatform(),
			networking: validNetworking(),
			valid:      true,
		},
		{
			name:     "allowed load balancer field with OpenShift managed default",
			platform: validPlatform(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					OpenStack: func() *openstack.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.OpenStackPlatformLoadBalancer{
							Type: configv1.LoadBalancerTypeOpenShiftManagedDefault,
						}
						return p
					}(),
				},
			},
			valid: true,
		},
		{
			name:     "allowed load balancer field with user-managed",
			platform: validPlatform(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					OpenStack: func() *openstack.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.OpenStackPlatformLoadBalancer{
							Type: configv1.LoadBalancerTypeUserManaged,
						}
						return p
					}(),
				},
			},
			valid: true,
		},
		{
			name:     "allowed load balancer field invalid type",
			platform: validPlatform(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					OpenStack: func() *openstack.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.OpenStackPlatformLoadBalancer{
							Type: "FooBar",
						}
						return p
					}(),
				},
			},
			expectedError: `^test-path\.loadBalancer.type: Invalid value: "FooBar": invalid load balancer type`,
			valid:         false,
		},
		{
			name: "non IP external dns",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ExternalDNS = []string{
					"invalid",
				}
				return p
			}(),
			networking:    validNetworking(),
			valid:         false,
			expectedError: `\"invalid\" is not a valid IP`,
		},
		{
			name: "valid external dns",
			platform: func() *openstack.Platform {
				p := validPlatform()
				p.ExternalDNS = []string{
					"192.168.1.1",
				}
				return p
			}(),
			networking: validNetworking(),
			valid:      true,
		},
		{
			name: "invalid subnet ID",
			platform: func() *openstack.Platform {
				p := validPlatform()
				fixedIP := openstack.FixedIP{
					Subnet: openstack.SubnetFilter{ID: "fake"},
				}
				p.ControlPlanePort = &openstack.PortTarget{
					FixedIPs: []openstack.FixedIP{fixedIP},
				}
				return p
			}(),
			networking:    validNetworking(),
			expectedError: `^test-path\.controlPlanePort.fixedIPs\[0\]\.subnet.id: Invalid value: "fake": invalid subnet ID`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Build default wrapping installConfig
			if tc.config == nil {
				tc.config = installConfig().build()
				tc.config.OpenStack = tc.platform
			}

			err := ValidatePlatform(tc.platform, tc.networking, field.NewPath("test-path"), tc.config).ToAggregate()
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expectedError, err)
			}
		})
	}
}

type installConfigBuilder struct {
	types.InstallConfig
}

func installConfig() *installConfigBuilder {
	return &installConfigBuilder{
		InstallConfig: types.InstallConfig{},
	}
}

func (icb *installConfigBuilder) build() *types.InstallConfig {
	return &icb.InstallConfig
}
