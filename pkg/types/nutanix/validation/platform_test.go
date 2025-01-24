package validation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/nutanix"
)

func validPlatform() *nutanix.Platform {
	return &nutanix.Platform{
		PrismCentral: nutanix.PrismCentral{
			Endpoint: nutanix.PrismEndpoint{Address: "test-pc", Port: 8080},
			Username: "test-username-pc",
			Password: "test-password-pc",
		},
		PrismElements: []nutanix.PrismElement{{
			UUID:     "test-pe-uuid",
			Endpoint: nutanix.PrismEndpoint{Address: "test-pe", Port: 8081},
		}},
		SubnetUUIDs: []string{"b06179c8-dea3-4f8e-818a-b2e88fbc2201"},
	}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name          string
		config        *types.InstallConfig
		platform      *nutanix.Platform
		expectedError string
	}{
		{
			name:     "minimal",
			platform: validPlatform(),
		},
		{
			name: "missing Prism Central address",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral.Endpoint.Address = ""
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.endpoint\.address: Required value: must specify the Prism Central endpoint address$`,
		},
		{
			name:     "allowed load balancer field with OpenShift managed default",
			platform: validPlatform(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					Nutanix: func() *nutanix.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.NutanixPlatformLoadBalancer{
							Type: configv1.LoadBalancerTypeOpenShiftManagedDefault,
						}
						return p
					}(),
				},
			},
		},
		{
			name:     "allowed load balancer field with user-managed",
			platform: validPlatform(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					Nutanix: func() *nutanix.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.NutanixPlatformLoadBalancer{
							Type: configv1.LoadBalancerTypeUserManaged,
						}
						return p
					}(),
				},
			},
		},
		{
			name:     "allowed load balancer field invalid type",
			platform: validPlatform(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					Nutanix: func() *nutanix.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.NutanixPlatformLoadBalancer{
							Type: "FooBar",
						}
						return p
					}(),
				},
			},
			expectedError: `^test-path\.loadBalancer.type: Invalid value: "FooBar": invalid load balancer type`,
		},
		{
			name: "missing Prism Central username",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral.Username = ""
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.username: Required value: must specify the Prism Central username$`,
		},
		{
			name: "missing Prism Central password",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral.Password = ""
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.password: Required value: must specify the Prism Central password$`,
		},
		{
			name: "missing prism elements",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismElements = []nutanix.PrismElement{}
				return p
			}(),
			expectedError: `^test-path\.prismElements: Required value: must specify one Prism Element$`,
		},
		{
			name: "missing prism element uuid",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismElements[0].UUID = ""
				return p
			}(),
			expectedError: `^test-path\.prismElements\.uuid: Required value: must specify the Prism Element UUID$`,
		},
		{
			name: "missing prism element address",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismElements[0].Endpoint.Address = ""
				return p
			}(),
			expectedError: `^test-path\.prismElements\.endpoint\.address: Required value: must specify the Prism Element endpoint address$`,
		},
		{
			name: "missing subnet",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.SubnetUUIDs[0] = ""
				return p
			}(),
			expectedError: `^test-path\.subnetUUIDs: Required value: must specify at least one subnet$`,
		},
		{
			name: "too many subnet items",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				for i := 1; i <= 32; i++ {
					p.SubnetUUIDs = append(p.SubnetUUIDs, fmt.Sprintf("subnet-%v", i))
				}
				return p
			}(),
			expectedError: `^test-path\.subnetUUIDs: Too many: 33: must have at most 32 items$`,
		},
		{
			name: "duplicate subnet item",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.SubnetUUIDs = append(p.SubnetUUIDs, p.SubnetUUIDs[0])
				return p
			}(),
			expectedError: `^test-path\.subnetUUIDs: Invalid value: "b06179c8-dea3-4f8e-818a-b2e88fbc2201": should not configure duplicate value$`,
		},
		{
			name: "Capital letters in Prism Central",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral.Endpoint.Address = "tEsT-PrismCentral"
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.endpoint\.address: Invalid value: "tEsT-PrismCentral": must be the domain name or IP address of the Prism Central$`,
		},
		{
			name: "URL as Prism Central",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral.Endpoint.Address = "https://test-pc"
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.endpoint\.address: Invalid value: "https://test-pc": must be the domain name or IP address of the Prism Central$`,
		},
		{
			name: "prismAPICallTimeout must be a positive integer",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				timeout := -1
				p.PrismAPICallTimeout = &timeout
				return p
			}(),
			expectedError: `^test-path\.prismAPICallTimeout: Invalid value: -1: must be a positive integer value$`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Build default wrapping installConfig
			if tc.config == nil {
				tc.config = installConfig().build()
				tc.config.Nutanix = tc.platform
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
