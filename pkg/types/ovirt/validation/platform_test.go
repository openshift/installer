package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ovirt"
)

func validPlatform() *ovirt.Platform {
	return &ovirt.Platform{
		ClusterID:              "0953a3fb-6682-49a6-8767-43f561045a59",
		StorageDomainID:        "57e42205-02ac-46e1-a6d1-07459d94bc51",
		VNICProfileID:          "57e42205-02ac-46e1-a6d1-07459d94bc52",
		NetworkName:            "ocp-blue",
		APIVIPs:                []string{"10.0.0.1"},
		IngressVIPs:            []string{"10.0.0.3"},
		DefaultMachinePlatform: nil,
	}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name          string
		config        *types.InstallConfig
		platform      *ovirt.Platform
		valid         bool
		expectedError string
	}{
		{
			name:     "minimal",
			platform: validPlatform(),
			valid:    true,
		},
		{
			name:     "forbidden load balancer field",
			platform: validPlatform(),
			config: &types.InstallConfig{
				Platform: types.Platform{
					Ovirt: func() *ovirt.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.OvirtPlatformLoadBalancer{
							Type: configv1.LoadBalancerTypeOpenShiftManagedDefault,
						}
						return p
					}(),
				},
			},
			valid:         false,
			expectedError: `^test-path\.loadBalancer: Forbidden: load balancer is not supported in this feature set`,
		},
		{
			name:     "allowed load balancer field with OpenShift managed default",
			platform: validPlatform(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					Ovirt: func() *ovirt.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.OvirtPlatformLoadBalancer{
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
					Ovirt: func() *ovirt.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.OvirtPlatformLoadBalancer{
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
					Ovirt: func() *ovirt.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.OvirtPlatformLoadBalancer{
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
			name: "invalid when empty cluster ID",
			platform: func() *ovirt.Platform {
				p := validPlatform()
				p.ClusterID = ""
				return p
			}(),
			expectedError: `^test-path.ovirt_cluster_id: Invalid value: "": invalid UUID length: 0`,
			valid:         false,
		},
		{
			name: "invalid when empty storage ID",
			platform: func() *ovirt.Platform {
				p := validPlatform()
				p.StorageDomainID = ""
				return p
			}(),
			expectedError: `^test-path.ovirt_storage_domain_id: Invalid value: "": invalid UUID length: 0`,
			valid:         false,
		},
		{
			name: "malformed vnic profile id",
			platform: func() *ovirt.Platform {
				p := validPlatform()
				p.VNICProfileID = "abcd-sdf"
				return p
			}(),
			expectedError: `^test-path.vnicProfileID: Invalid value: "abcd-sdf": invalid UUID length: 8`,
			valid:         false,
		},
		{
			name: "valid empty vnic profile id",
			platform: func() *ovirt.Platform {
				p := validPlatform()
				p.VNICProfileID = ""
				return p
			}(),
			valid: true,
		},
		{
			name: "valid machine pool",
			platform: func() *ovirt.Platform {
				p := validPlatform()
				p.DefaultMachinePlatform = &ovirt.MachinePool{}
				return p
			}(),
			valid: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Build default wrapping installConfig
			if tc.config == nil {
				tc.config = installConfig().build()
				tc.config.Ovirt = tc.platform
			}

			err := ValidatePlatform(tc.platform, field.NewPath("test-path"), tc.config).ToAggregate()
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
