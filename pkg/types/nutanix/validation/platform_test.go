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
		{
			name: "invalid Prism Central port zero",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral.Endpoint.Port = 0
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.endpoint\.port: Invalid value: 0: The Prism Central endpoint port is invalid, must be in the range of 1 to 65535$`,
		},
		{
			name: "invalid Prism Central port too high",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismCentral.Endpoint.Port = 65536
				return p
			}(),
			expectedError: `^test-path\.prismCentral\.endpoint\.port: Invalid value: 65536: The Prism Central endpoint port is invalid, must be in the range of 1 to 65535$`,
		},
		{
			name: "invalid Prism Element endpoint address with capitals",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismElements[0].Endpoint.Address = "tEsT-PrismElement"
				return p
			}(),
			expectedError: `^test-path\.prismElements\.endpoint\.address: Invalid value: "tEsT-PrismElement": must be the domain name or IP address of the Prism Element \(cluster\)$`,
		},
		{
			name: "invalid Prism Element endpoint address with URL",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismElements[0].Endpoint.Address = "https://test-pe"
				return p
			}(),
			expectedError: `^test-path\.prismElements\.endpoint\.address: Invalid value: "https://test-pe": must be the domain name or IP address of the Prism Element \(cluster\)$`,
		},
		{
			name: "invalid Prism Element port zero",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismElements[0].Endpoint.Port = 0
				return p
			}(),
			expectedError: `^test-path\.prismElements\.endpoint\.port: Invalid value: 0: The Prism Element endpoint port is invalid, must be in the range of 1 to 65535$`,
		},
		{
			name: "invalid Prism Element port too high",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismElements[0].Endpoint.Port = 65536
				return p
			}(),
			expectedError: `^test-path\.prismElements\.endpoint\.port: Invalid value: 65536: The Prism Element endpoint port is invalid, must be in the range of 1 to 65535$`,
		},
		{
			name: "too many Prism Elements",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.PrismElements = append(p.PrismElements, nutanix.PrismElement{
					UUID:     "second-pe-uuid",
					Endpoint: nutanix.PrismEndpoint{Address: "test-pe-2", Port: 8082},
				})
				return p
			}(),
			expectedError: `^test-path\.prismElements: Required value: must specify one Prism Element$`,
		},
		{
			name: "valid platform with clusterOSImage",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.ClusterOSImage = "http://example.com/rhcos-nutanix.qcow2"
				return p
			}(),
		},
		{
			name: "failureDomain with invalid name special characters",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.FailureDomains = []nutanix.FailureDomain{
					{
						Name:         "!!!",
						PrismElement: nutanix.PrismElement{UUID: "fd-pe-uuid", Endpoint: nutanix.PrismEndpoint{Address: "fd-pe", Port: 9440}},
						SubnetUUIDs:  []string{"fd-subnet-uuid"},
					},
				}
				return p
			}(),
			expectedError: `^test-path\.failureDomain\.name: Invalid value: "!!!": failureDomain name should match the pattern "\[a-z0-9\]\(\[-a-z0-9\]\*\[a-z0-9\]\)\?"\.$`,
		},
		{
			name: "failureDomain with empty prismElement UUID",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.FailureDomains = []nutanix.FailureDomain{
					{
						Name:         "fd-1",
						PrismElement: nutanix.PrismElement{UUID: "", Endpoint: nutanix.PrismEndpoint{Address: "fd-pe", Port: 9440}},
						SubnetUUIDs:  []string{"fd-subnet-uuid"},
					},
				}
				return p
			}(),
			expectedError: `^test-path\.failureDomain\.prismElement\.uuid: Required value: failureDomain prismElement uuid cannot be empty$`,
		},
		{
			name: "failureDomain with empty subnetUUIDs",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.FailureDomains = []nutanix.FailureDomain{
					{
						Name:         "fd-1",
						PrismElement: nutanix.PrismElement{UUID: "fd-pe-uuid", Endpoint: nutanix.PrismEndpoint{Address: "fd-pe", Port: 9440}},
						SubnetUUIDs:  []string{},
					},
				}
				return p
			}(),
			expectedError: `^test-path\.failureDomain\.subnetUUIDs: Required value: must specify at least one subnet$`,
		},
		{
			name: "failureDomain with empty storageContainer referenceName",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.FailureDomains = []nutanix.FailureDomain{
					{
						Name:         "fd-1",
						PrismElement: nutanix.PrismElement{UUID: "fd-pe-uuid", Endpoint: nutanix.PrismEndpoint{Address: "fd-pe", Port: 9440}},
						SubnetUUIDs:  []string{"fd-subnet-uuid"},
						StorageContainers: []nutanix.StorageResourceReference{
							{ReferenceName: "", UUID: "sc-uuid"},
						},
					},
				}
				return p
			}(),
			expectedError: `^test-path\.failureDomain\.storageContainers\.referenceName: Required value: failureDomain "fd-1": missing storageContainer referenceName$`,
		},
		{
			name: "failureDomain with empty storageContainer UUID",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.FailureDomains = []nutanix.FailureDomain{
					{
						Name:         "fd-1",
						PrismElement: nutanix.PrismElement{UUID: "fd-pe-uuid", Endpoint: nutanix.PrismEndpoint{Address: "fd-pe", Port: 9440}},
						SubnetUUIDs:  []string{"fd-subnet-uuid"},
						StorageContainers: []nutanix.StorageResourceReference{
							{ReferenceName: "sc-ref", UUID: ""},
						},
					},
				}
				return p
			}(),
			expectedError: `^test-path\.failureDomain\.storageContainers\.uuid: Required value: failureDomain "fd-1": missing storageContainer uuid$`,
		},
		{
			name: "failureDomain with empty dataSourceImage referenceName",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.FailureDomains = []nutanix.FailureDomain{
					{
						Name:         "fd-1",
						PrismElement: nutanix.PrismElement{UUID: "fd-pe-uuid", Endpoint: nutanix.PrismEndpoint{Address: "fd-pe", Port: 9440}},
						SubnetUUIDs:  []string{"fd-subnet-uuid"},
						DataSourceImages: []nutanix.StorageResourceReference{
							{ReferenceName: "", UUID: "dsi-uuid"},
						},
					},
				}
				return p
			}(),
			expectedError: `^test-path\.failureDomain\.dataSourceImages\.referenceName: Required value: failureDomain "fd-1": missing dataSourceImage referenceName$`,
		},
		{
			name: "failureDomain with both dataSourceImage UUID and name empty",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.FailureDomains = []nutanix.FailureDomain{
					{
						Name:         "fd-1",
						PrismElement: nutanix.PrismElement{UUID: "fd-pe-uuid", Endpoint: nutanix.PrismEndpoint{Address: "fd-pe", Port: 9440}},
						SubnetUUIDs:  []string{"fd-subnet-uuid"},
						DataSourceImages: []nutanix.StorageResourceReference{
							{ReferenceName: "dsi-ref", UUID: "", Name: ""},
						},
					},
				}
				return p
			}(),
			expectedError: `^test-path\.failureDomain\.dataSourceImages: Required value: failureDomain "fd-1": both the dataSourceImage's uuid and name are empty, you need to configure one\.$`,
		},
		{
			name: "valid failureDomain",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.FailureDomains = []nutanix.FailureDomain{
					{
						Name:         "fd-1",
						PrismElement: nutanix.PrismElement{UUID: "fd-pe-uuid", Endpoint: nutanix.PrismEndpoint{Address: "fd-pe", Port: 9440}},
						SubnetUUIDs:  []string{"fd-subnet-uuid"},
						StorageContainers: []nutanix.StorageResourceReference{
							{ReferenceName: "sc-ref", UUID: "sc-uuid"},
						},
						DataSourceImages: []nutanix.StorageResourceReference{
							{ReferenceName: "dsi-ref", UUID: "dsi-uuid"},
						},
					},
				}
				return p
			}(),
		},
		{
			name: "valid failureDomain with multiple subnets for multi-NIC",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.FailureDomains = []nutanix.FailureDomain{
					{
						Name:         "fd-1",
						PrismElement: nutanix.PrismElement{UUID: "fd-pe-uuid", Endpoint: nutanix.PrismEndpoint{Address: "fd-pe", Port: 9440}},
						SubnetUUIDs:  []string{"fd-subnet-uuid-1", "fd-subnet-uuid-2"},
					},
				}
				return p
			}(),
		},
		{
			name: "valid platform with multiple subnets for multi-NIC",
			platform: func() *nutanix.Platform {
				p := validPlatform()
				p.SubnetUUIDs = []string{
					"b06179c8-dea3-4f8e-818a-b2e88fbc2201",
					"c17280d9-efb4-5f9f-929b-c3f99gcd3312",
				}
				return p
			}(),
		},
		{
			name:     "external DNS records without user-managed load balancer",
			platform: validPlatform(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					Nutanix: func() *nutanix.Platform {
						p := validPlatform()
						p.DNSRecordsType = configv1.DNSRecordsTypeExternal
						p.LoadBalancer = &configv1.NutanixPlatformLoadBalancer{
							Type: configv1.LoadBalancerTypeOpenShiftManagedDefault,
						}
						return p
					}(),
				},
			},
			expectedError: `^test-path\.dnsRecordsType: Invalid value: "External": external DNS records can only be configured with user-managed loadbalancers$`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Build default wrapping installConfig
			if tc.config == nil {
				tc.config = installConfig().build()
				tc.config.Nutanix = tc.platform
			}
			err := ValidatePlatform(tc.platform, field.NewPath("test-path"), tc.config, false).ToAggregate()
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
