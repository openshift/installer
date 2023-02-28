package validation

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

func validPlatform() *vsphere.Platform {
	return &vsphere.Platform{
		VCenters: []vsphere.VCenter{
			{
				Server:   "test-vcenter",
				Port:     443,
				Username: "test-username",
				Password: "test-password",
				Datacenters: []string{
					"test-datacenter",
				},
			},
		},
		FailureDomains: []vsphere.FailureDomain{
			{
				Name:   "test-east-1a",
				Region: "test-east",
				Zone:   "test-east-1a",
				Server: "test-vcenter",
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter",
					ComputeCluster: "/test-datacenter/host/test-cluster",
					Datastore:      "/test-datacenter/datastore/test-datastore",
					Networks:       []string{"test-portgroup"},
					ResourcePool:   "/test-datacenter/host/test-cluster/Resources/test-resourcepool",
					Folder:         "/test-datacenter/vm/test-folder",
				},
			},
			{
				Name:   "test-east-2a",
				Region: "test-east",
				Zone:   "test-east-2a",
				Server: "test-vcenter",
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter",
					ComputeCluster: "/test-datacenter/host/test-cluster",
					Datastore:      "/test-datacenter/datastore/test-datastore",
					Networks:       []string{"test-portgroup"},
					Folder:         "/test-datacenter/vm/test-folder",
				},
			},
		},
	}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name          string
		config        *types.InstallConfig
		platform      *vsphere.Platform
		expectedError string
	}{
		{
			name: "Valid diskType",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.DiskType = "eagerZeroedThick"
				return p
			}(),
		},
		{
			name: "Invalid diskType",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.DiskType = "invalidDiskType"
				return p
			}(),
			expectedError: `^test-path\.diskType: Invalid value: "invalidDiskType": diskType must be one of \[eagerZeroedThick thick thin\]$`,
		},

		{
			name:     "Valid Multi-zone platform",
			platform: validPlatform(),
		},
		{
			name: "Multi-zone platform missing failureDomains",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.FailureDomains = make([]vsphere.FailureDomain, 0)
				return p
			}(),
			expectedError: `^test-path.failureDomains: Required value: must be defined`,
		},
		{
			name: "Multi-zone platform vCenter missing server",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenters[0].Server = ""
				return p
			}(),
			expectedError: `test-path\.vcenters\.server: Required value: must be the domain name or IP address of the vCenter(.*)`,
		},
		{
			name: "Multi-zone platform more than one vCenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenters = append(p.VCenters, vsphere.VCenter{
					Server: "additional-vcenter",
				})
				return p
			}(),
			expectedError: `^test-path\.vcenters: Too many: 2: must have at most 1 items`,
		},
		{
			name: "Multi-zone platform Capital letters in vCenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenters[0].Server = "tEsT-vCenter"
				return p
			}(),
			expectedError: `(.*)test-path\.vcenters.server: Invalid value: "tEsT-vCenter": must be the domain name or IP address of the vCenter`,
		},
		{
			name: "Multi-zone missing username",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenters[0].Username = ""
				return p
			}(),
			expectedError: `^test-path\.vcenters.username: Required value: must specify the username$`,
		},
		{
			name: "Multi-zone missing password",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenters[0].Password = ""
				return p
			}(),
			expectedError: `^test-path\.vcenters.password: Required value: must specify the password$`,
		},
		{
			name: "Multi-zone missing datacenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenters[0].Datacenters = []string{}
				return p
			}(),
			expectedError: `^test-path\.vcenters.datacenters: Required value: must specify at least one datacenter$`,
		},
		{
			name: "Multi-zone platform wrong vCenter name in failureDomain zone",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.FailureDomains[0].Server = "bad-vcenter"
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.server: Invalid value: "bad-vcenter": server does not exist in vcenters`,
		},
		{
			name: "Multi-zone platform failure domain topology cluster relative path",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.FailureDomains[0].Topology.ComputeCluster = "incomplete-path"
				p.FailureDomains[0].Topology.ResourcePool = "/test-datacenter/host/incomplete-path/Resources/test-resourcepool"
				return p
			}(),
			expectedError: `(.*)test-path\.failureDomains\.topology\.computeCluster: Invalid value: "incomplete-path": full path of compute cluster must be provided in format /<datacenter>/host/<cluster>`,
		},
		{
			name: "Multi-zone platform datacenter in failure domain topology doesn't match cluster datacenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.FailureDomains[0].Topology.ComputeCluster = "/other-datacenter/host/cluster"
				return p
			}(),
			expectedError: `^test-path.failureDomains.topology.computeCluster: Invalid value: "/other-datacenter/host/cluster": compute cluster must be in datacenter test-datacenter`,
		},
		{
			name: "Multi-zone platform failureDomain missing name",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.FailureDomains[0].Name = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.name: Required value: must specify the name`,
		},
		{
			name: "Multi-zone platform failureDomain region missing name",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.FailureDomains[0].Region = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.region: Required value: must specify region tag value`,
		},
		{
			name: "Multi-zone platform failureDomain zone missing name",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.FailureDomains[0].Name = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.name: Required value: must specify the name`,
		},
		{
			name: "Multi-zone platform failureDomain zone missing tag category",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.FailureDomains[0].Zone = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.zone: Required value: must specify zone tag value`,
		},
		{
			name:     "forbidden load balancer field",
			platform: validPlatform(),
			config: &types.InstallConfig{
				Platform: types.Platform{
					VSphere: func() *vsphere.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.VSpherePlatformLoadBalancer{
							Type: configv1.LoadBalancerTypeOpenShiftManagedDefault,
						}
						return p
					}(),
				},
			},
			expectedError: `^test-path\.loadBalancer: Forbidden: load balancer is not supported in this feature set`,
		},
		{
			name:     "allowed load balancer field with OpenShift managed default",
			platform: validPlatform(),
			config: &types.InstallConfig{
				FeatureSet: configv1.TechPreviewNoUpgrade,
				Platform: types.Platform{
					VSphere: func() *vsphere.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.VSpherePlatformLoadBalancer{
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
					VSphere: func() *vsphere.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.VSpherePlatformLoadBalancer{
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
					VSphere: func() *vsphere.Platform {
						p := validPlatform()
						p.LoadBalancer = &configv1.VSpherePlatformLoadBalancer{
							Type: "FooBar",
						}
						return p
					}(),
				},
			},
			expectedError: `^test-path\.loadBalancer.type: Invalid value: "FooBar": invalid load balancer type`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Build default wrapping installConfig
			if tc.config == nil {
				tc.config = installConfig().build()
				tc.config.VSphere = tc.platform
			}
			err := ValidatePlatform(tc.platform, false, field.NewPath("test-path"), tc.config).ToAggregate()
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, regexp.MustCompile(tc.expectedError), err)
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
