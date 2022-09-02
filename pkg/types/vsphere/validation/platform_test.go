package validation

import (
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
		DeploymentZones: []vsphere.DeploymentZone{
			{
				Name:          "test-dz-east-1a",
				Server:        "test-vcenter",
				FailureDomain: "test-east-1a",
				ControlPlane:  vsphere.Allowed,
				PlacementConstraint: vsphere.PlacementConstraint{
					ResourcePool: "/test-datacenter/host/cluster/Resources/test-resourcepool",
					Folder:       "/test-datacenter/vm/test-folder",
				},
			},
			{
				Name:          "test-dz-east-2a",
				Server:        "test-vcenter",
				FailureDomain: "test-east-1a",
				ControlPlane:  vsphere.NotAllowed,
				PlacementConstraint: vsphere.PlacementConstraint{
					ResourcePool: "/test-datacenter/host/cluster/Resources/test-resourcepool",
					Folder:       "/test-datacenter/vm/test-folder",
				},
			},
		},
		FailureDomains: []vsphere.FailureDomain{
			{
				Name: "test-east-1a",
				Region: vsphere.FailureDomainCoordinate{
					Name:        "test-region-east",
					Type:        "Datacenter",
					TagCategory: "openshift-region",
				},
				Zone: vsphere.FailureDomainCoordinate{
					Name:        "test-zone-1a",
					Type:        "ComputeCluster",
					TagCategory: "openshift-zone",
				},
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter",
					ComputeCluster: "/test-datacenter/host/cluster",
					Hosts:          nil,
					Networks: []string{
						"test-network-1",
					},
					Datastore: "test-datastore",
				}},
		},
	}
}

func TestValidatePlatform(t *testing.T) {
	cases := []struct {
		name          string
		platform      *vsphere.Platform
		expectedError string
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
			name: "Multi-zone platform missing deploymentZones",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.DeploymentZones = make([]vsphere.DeploymentZone, 0)
				return p
			}(),
			expectedError: `^test-path.deploymentZones: Required value: must be defined if vcenters is defined`,
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
			name: "Multi-zone platform wrong vCenter name in deployment zone",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.DeploymentZones[0].Server = "bad-vcenter"
				return p
			}(),
			expectedError: `^test-path\.deploymentZones\.server: Invalid value: "bad-vcenter": server does not exist in vcenters`,
		},
		{
			name: "Multi-zone platform missing deployment zone name",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.DeploymentZones[0].Name = ""
				return p
			}(),
			expectedError: `^test-path\.deploymentZones\.name: Required value: must specify the name`,
		},
		{
			name: "Multi-zone platform missing failureDomain name in deployment zone",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.DeploymentZones[0].FailureDomain = ""
				return p
			}(),
			expectedError: `^test-path\.deploymentZones\.failureDomain: Required value: must specify the failureDomain name`,
		},
		{
			name: "Multi-zone platform failureDomain name does not exist in failureDomains",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.DeploymentZones[0].FailureDomain = "bad-domain"
				return p
			}(),
			expectedError: `^test-path\.deploymentZones\.failureDomain: Invalid value: "bad-domain": does not exist in failureDomains`,
		},
		{
			name: "Multi-zone platform controlPlane value invalid",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.DeploymentZones[0].ControlPlane = "maybe"
				return p
			}(),
			expectedError: `^test-path\.deploymentZones\.controlPlane: Invalid value: "maybe": valid values are Allowed and NotAllowed`,
		},
		{
			name: "Multi-zone platform folder placement constraint relative path",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.DeploymentZones[0].PlacementConstraint.Folder = "incomplete-path"
				return p
			}(),
			expectedError: `^test-path\.deploymentZones\.placementConstraint\.folder: Invalid value: "incomplete-path": full path of folder must be provided in format /<datacenter>/vm/<folder>`,
		},
		{
			name: "Multi-zone platform resource pool placement constraint relative path",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.DeploymentZones[0].PlacementConstraint.ResourcePool = "incomplete-path"
				return p
			}(),
			expectedError: `^test-path\.deploymentZones\.placementConstraint\.resourcePool: Invalid value: "incomplete-path": full path of resource pool must be provided in format /<datacenter>/host/<cluster>/Resources/<resource-pool>`,
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
			expectedError: `^test-path.failureDomains.topology.computeCluster: Required value: must specify a computeCluster`,
		},
		{
			name: "Multi-zone platform failure domain topology datastore required",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Topology.Datastore = ""
				return p
			}(),
			expectedError: `^test-path.failureDomains.topology.datastore: Required value: must specify a datastore`,
		},
		{
			name: "Multi-zone platform datacenter in failure domain topology is not defined",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Topology.Datacenter = "unknown-datacenter"
				p.FailureDomains[0].Topology.ComputeCluster = "/unknown-datacenter/host/cluster"
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.topology\.datacenter: Invalid value: "unknown-datacenter": datacenter unknown-datacenter in failure domain topology does not exist in associated vCenter`,
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
				p.FailureDomains[0].Region.Name = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.region\.name: Required value: must specify the name`,
		},
		{
			name: "Multi-zone platform failureDomain region invalid type",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Region.Type = "invalid"
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.region\.type: Invalid value: "invalid": must be ComputeCluster or Datacenter`,
		},
		{
			name: "Multi-zone platform failureDomain region missing tag category",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Region.TagCategory = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.region\.tagCategory: Required value: must specify a tag category`,
		},
		{
			name: "Multi-zone platform failureDomain zone missing name",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Zone.Name = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.zone\.name: Required value: must specify the name`,
		},
		{
			name: "Multi-zone platform failureDomain zone invalid type",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Zone.Type = "invalid"
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.zone\.type: Invalid value: "invalid": must be ComputeCluster or Datacenter`,
		},
		{
			name: "Multi-zone platform failureDomain HostGroup is not supported",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Zone.Type = "HostGroup"
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.zone\.type: Invalid value: "HostGroup": is not supported`,
		},
		{
			name: "Multi-zone platform failureDomain zone missing tag category",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Zone.TagCategory = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.zone\.tagCategory: Required value: must specify a tag category`,
		},
		{
			name: "Multi-zone platform failureDomain topology missing datacenter",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Topology.Datacenter = ""
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.topology\.datacenter: Required value: must specify a datacenter`,
		},
		{
			name: "Multi-zone platform failureDomain topology hosts missing vmGroupName",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Topology.Hosts = &vsphere.FailureDomainHosts{
					VMGroupName:   "",
					HostGroupName: "test-host-group-name",
				}
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.topology\.hosts\.vmGroupName: Required value: must specify the vmGroupName`,
		},
		{
			name: "Multi-zone platform failureDomain topology hosts missing vmGroupName",
			platform: func() *vsphere.Platform {
				p := validMultiVCenterPlatform()
				p.FailureDomains[0].Topology.Hosts = &vsphere.FailureDomainHosts{
					VMGroupName:   "test-vm-group-name",
					HostGroupName: "",
				}
				return p
			}(),
			expectedError: `^test-path\.failureDomains\.topology\.hosts.hostGroupName: Required value: must specify the hostGroupName`,
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
		})
	}
}
