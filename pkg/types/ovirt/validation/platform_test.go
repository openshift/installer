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
			expectedError: `^test-path: Forbidden: Platform oVirt is no longer supported`,
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
