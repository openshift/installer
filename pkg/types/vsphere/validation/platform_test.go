package validation

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/vsphere"
)

func validPlatform() *vsphere.Platform {
	return &vsphere.Platform{
		VCenter:          "test-vcenter",
		Username:         "test-username",
		Password:         "test-password",
		Datacenter:       "test-datacenter",
		DefaultDatastore: "test-datastore",
	}
}

func validMultiVCenterPlatform() *vsphere.Platform {
	return &vsphere.Platform{
		VCenter:          "test-vcenter",
		Username:         "test-username",
		Password:         "test-password",
		Datacenter:       "test-datacenter",
		DefaultDatastore: "test-datastore",
		Folder:           "/test-datacenter/vm/test-folder",
		Cluster:          "test-cluster",
		ResourcePool:     "/test-datacenter/host/test-cluster/Resources/test-resource-pool",
		Network:          "test-network",
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
					Datastore:      "test-datastore",
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
					Datastore:      "test-datastore",
					Networks:       []string{"test-portgroup"},
					Folder:         "/test-datacenter/vm/test-folder",
				},
			},
		},
	}
}

func checkIfFailureDomainDefaultsApplied(platform *vsphere.Platform) error {
	for _, failureDomain := range platform.FailureDomains {
		defaultsApplied := failureDomain.Server == platform.VCenter &&
			failureDomain.Topology.ResourcePool == platform.ResourcePool &&
			failureDomain.Topology.ComputeCluster == fmt.Sprintf("/%s/host/%s", platform.Datacenter, platform.Cluster) &&
			failureDomain.Topology.Networks[0] == platform.Network &&
			failureDomain.Topology.Datastore == platform.DefaultDatastore &&
			failureDomain.Topology.Folder == platform.Folder &&
			failureDomain.Topology.Datacenter == platform.Datacenter
		if defaultsApplied == false {
			return errors.New("defaults not applied")
		}
	}
	return nil
}

func checkIfFailureDomainOverridesApplied(platform *vsphere.Platform) error {
	for _, failureDomain := range platform.FailureDomains {
		defaultsNotApplied := failureDomain.Topology.ResourcePool != platform.ResourcePool &&
			failureDomain.Topology.ComputeCluster != fmt.Sprintf("/%s/host/%s", platform.Datacenter, platform.Cluster) &&
			failureDomain.Topology.Networks[0] != platform.Network &&
			failureDomain.Topology.Datastore != platform.DefaultDatastore &&
			failureDomain.Topology.Folder != platform.Folder &&
			failureDomain.Topology.Datacenter != platform.Datacenter
		if defaultsNotApplied == false {
			return errors.New("overrides not applied")
		}
	}
	return nil
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name                   string
		platform               *vsphere.Platform
		postValidationFunction func(*vsphere.Platform) error
		expectedError          string
	}{
		{
			name:     "minimal",
			platform: validPlatform(),
		},
		{
			name: "missing vCenter name",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenter = ""
				return p
			}(),
			expectedError: `^test-path\.vCenter: Required value: must specify the name of the vCenter$`,
		},
		{
			name: "missing username",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Username = ""
				return p
			}(),
			expectedError: `^test-path\.username: Required value: must specify the username$`,
		},
		{
			name: "missing password",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Password = ""
				return p
			}(),
			expectedError: `^test-path\.password: Required value: must specify the password$`,
		},
		{
			name: "missing datacenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Datacenter = ""
				return p
			}(),
			expectedError: `^test-path\.datacenter: Required value: must specify the datacenter$`,
		},
		{
			name: "missing default datastore",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.DefaultDatastore = ""
				return p
			}(),
			expectedError: `^test-path\.defaultDatastore: Required value: must specify the default datastore$`,
		},
		{
			name: "Invalid folder path",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Folder = "test-folder"
				return p
			}(),
			expectedError: `^test-path\.folder: Invalid value: "test-folder": folder must be absolute path: expected prefix /test-datacenter/vm/`,
		},
		{
			name: "Capital letters in vCenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenter = "tEsT-vCenter"
				return p
			}(),
			expectedError: `^test-path\.vCenter: Invalid value: "tEsT-vCenter": must be the domain name or IP address of the vCenter`,
		},
		{
			name: "URL as vCenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VCenter = "https://test-center"
				return p
			}(),
			expectedError: `^test-path\.vCenter: Invalid value: "https://test-center": must be the domain name or IP address of the vCenter$`,
		},
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
			name: "Valid Multi-zone platform",
			platform: func() *vsphere.Platform {
				return validMultiVCenterPlatform()
			}(),
		},
		{
			name: "Multi-zone platform missing failureDomains",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains = make([]vsphere.FailureDomain, 0)
				return p
			}(),
			expectedError: `^test-path.failureDomains: Required value: must be defined if vcenters is defined`,
		},
		{
			name: "Multi-zone platform vCenter missing server",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.VCenters[0].Server = ""
				return p
			}(),
			expectedError: `^test-path\.vcenters\.server: Required value: must be the domain name or IP address of the vCenter`,
		},
		{
			name: "Multi-zone platform more than one vCenter",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
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
				p := validMultiVCenterPlatform()
				p.VCenters[0].Server = "tEsT-vCenter"
				return p
			}(),
			expectedError: `^test-path\.vcenters.server: Invalid value: "tEsT-vCenter": must be the domain name or IP address of the vCenter`,
		},
		{
			name: "Multi-zone missing username",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.VCenters[0].Username = ""
				return p
			}(),
			expectedError: `^test-path\.vcenters.username: Required value: must specify the username$`,
		},
		{
			name: "Multi-zone missing password",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.VCenters[0].Password = ""
				return p
			}(),
			expectedError: `^test-path\.vcenters.password: Required value: must specify the password$`,
		},
		{
			name: "Multi-zone missing datacenter",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.VCenters[0].Datacenters = []string{}
				return p
			}(),
			expectedError: `^test-path\.vcenters.datacenters: Required value: must specify at least one datacenter$`,
		},
		{
			name: "Multi-zone platform wrong vCenter name in failureDomain zone",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Server = "bad-vcenter"
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.server: Invalid value: "bad-vcenter": server does not exist in vcenters`,
		},
		{
			name: "Multi-zone platform failure domain topology cluster relative path",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Topology.ComputeCluster = "incomplete-path"
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.topology\.computeCluster: Invalid value: "incomplete-path": full path of compute cluster must be provided in format /<datacenter>/host/<cluster>`,
		},
		{
			name: "Multi-zone platform failure domain topology compute cluster required",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Topology.ComputeCluster = ""
				return p
			}(),
			postValidationFunction: func(platform *vsphere.Platform) error {
				if platform.FailureDomains[0].Topology.ComputeCluster != fmt.Sprintf("/%s/host/%s", platform.Datacenter, platform.Cluster) {
					return errors.New("topology cluster not set to default")
				}
				return nil
			},
		},
		{
			name: "Multi-zone platform failure domain topology datastore required",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Topology.Datastore = ""
				return p
			}(),
			postValidationFunction: func(platform *vsphere.Platform) error {
				if platform.FailureDomains[0].Topology.Datastore != platform.DefaultDatastore {
					return errors.New("topology datastore not set to default")
				}
				return nil
			},
		},
		{
			name: "Multi-zone platform datacenter in failure domain topology doesn't match cluster datacenter",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Topology.ComputeCluster = "/other-datacenter/host/cluster"
				return p
			}(),
			expectedError: `^test-path.failureDomains.topology.computeCluster: Invalid value: "/other-datacenter/host/cluster": compute cluster must be in datacenter test-datacenter`,
		},
		{
			name: "Multi-zone platform failureDomain missing name",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Name = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.name: Required value: must specify the name`,
		},
		{
			name: "Multi-zone platform failureDomain region missing name",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Region = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.region: Required value: must specify region tag value`,
		},
		{
			name: "Multi-zone platform failureDomain zone missing name",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Name = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.name: Required value: must specify the name`,
		},
		{
			name: "Multi-zone platform failureDomain zone missing tag category",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Zone = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.zone: Required value: must specify zone tag value`,
		},
		{
			name: "Multi-zone platform failureDomain topology missing datacenter",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Topology.Datacenter = ""
				return p
			}(),
			postValidationFunction: func(platform *vsphere.Platform) error {
				if platform.FailureDomains[0].Topology.Datacenter != platform.Datacenter {
					return errors.New("topology datacenter not set to default")
				}
				return nil
			},
		},
		{
			name: "Multi-zone platform failureDomain topology empty networks",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Topology.Networks = []string{}
				return p
			}(),
			postValidationFunction: func(platform *vsphere.Platform) error {
				if len(platform.FailureDomains[0].Topology.Networks) != 1 ||
					platform.FailureDomains[0].Topology.Networks[0] != platform.Network {
					return errors.New("topology networks not set to default")
				}
				return nil
			},
		},
		{
			name: "Multi-zone platform failureDomain defaults applied",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.VCenters = []vsphere.VCenter{}
				p.FailureDomains[0].Server = ""
				p.FailureDomains[0].Topology = vsphere.Topology{}
				p.FailureDomains[1].Server = ""
				p.FailureDomains[1].Topology = vsphere.Topology{}
				return p
			}(),
			postValidationFunction: checkIfFailureDomainDefaultsApplied,
		},
		{
			name: "Multi-zone platform failureDomain overrides applied",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.VCenters = []vsphere.VCenter{}
				p.FailureDomains[0].Server = "test-vcenter"
				p.FailureDomains[0].Topology = vsphere.Topology{
					Datacenter:     "dc1",
					ComputeCluster: "/dc1/host/c1",
					Networks: []string{
						"network1",
					},
					Datastore:    "ds1",
					ResourcePool: "/dc1/host/c1/Resources/rp1",
					Folder:       "/dc1/vm/folder1",
				}
				p.FailureDomains[1].Server = "test-vcenter"
				p.FailureDomains[1].Topology = vsphere.Topology{
					Datacenter:     "dc2",
					ComputeCluster: "/dc2/host/c2",
					Networks: []string{
						"network2",
					},
					Datastore:    "ds2",
					ResourcePool: "/dc2/host/c2/Resources/rp2",
					Folder:       "/dc2/vm/folder2",
				}
				return p
			}(),
			postValidationFunction: checkIfFailureDomainOverridesApplied,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePlatform(tc.platform, field.NewPath("test-path")).ToAggregate()
			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expectedError, err)
			}
			if tc.postValidationFunction != nil {
				err := tc.postValidationFunction(tc.platform)
				if err != nil {
					assert.NoError(t, err)
				}
			}
		})
	}
}
