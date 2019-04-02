package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/vsphere"
)

func validPlatform() *vsphere.Platform {
	return &vsphere.Platform{
		VirtualCenters: []vsphere.VirtualCenter{
			{
				Name:        "test-server-name",
				Username:    "test-username",
				Password:    "test-password",
				Datacenters: []string{"test-datacenter"},
			},
		},
		Workspace: vsphere.Workspace{
			Server:           "test-server-name",
			Datacenter:       "test-datacenter",
			DefaultDatastore: "test-datastore",
			ResourcePoolPath: "test-resource-pool",
			Folder:           "test-folder",
		},
		SCSIControllerType: "test-controller-type",
		PublicNetwork:      "test-network",
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
			name: "empty vCenters",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VirtualCenters = nil
				return p
			}(),
			expectedError: `^\[test-path\.virtualCenters: Required value: must include at least one vCenter, test-path\.workspace\.server: Invalid value: "test-server-name": workspace server must be a specified vCenter]$`,
		},
		{
			name: "missing vCenter name",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VirtualCenters[0].Name = ""
				return p
			}(),
			expectedError: `^\[test-path\.virtualCenters\[0]\.name: Required value: vCenter must have a name, test-path\.workspace\.server: Invalid value: "test-server-name": workspace server must be a specified vCenter]$`,
		},
		{
			name: "missing vCenter username",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VirtualCenters[0].Username = ""
				return p
			}(),
			expectedError: `^test-path\.virtualCenters\[0]\.username: Required value: username required for each vCenter$`,
		},
		{
			name: "missing vCenter password",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VirtualCenters[0].Password = ""
				return p
			}(),
			expectedError: `^test-path\.virtualCenters\[0]\.password: Required value: password required for each vCenter$`,
		},
		{
			name: "empty vCenter datacenters",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VirtualCenters[0].Datacenters = nil
				return p
			}(),
			expectedError: `^\[test-path\.virtualCenters\[0]\.datacenters: Required value: must include at least one datacenter, test-path\.workspace\.datacenter: Invalid value: "test-datacenter": workspace datacenter must be a datacenter in the workspace server]$`,
		},
		{
			name: "multiple vCenter datacenters",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VirtualCenters[0].Datacenters = []string{"test-datacenter", "other-datacenter"}
				return p
			}(),
		},
		{
			name: "duplicate vCenter datacenters",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.VirtualCenters[0].Datacenters = []string{"test-datacenter", "test-datacenter"}
				return p
			}(),
			expectedError: `^test-path\.virtualCenters\[0]\.datacenters\[1]: Duplicate value: "test-datacenter"$`,
		},
		{
			name: "missing workspace server",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Workspace.Server = ""
				return p
			}(),
			expectedError: `^test-path\.workspace\.server: Required value: must specify the workspace server$`,
		},
		{
			name: "no vCenter for workspace server",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Workspace.Server = "other-server-name"
				return p
			}(),
			expectedError: `^test-path\.workspace\.server: Invalid value: "other-server-name": workspace server must be a specified vCenter$`,
		},
		{
			name: "missing workspace datacenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Workspace.Datacenter = ""
				return p
			}(),
			expectedError: `^test-path\.workspace\.datacenter: Required value: must specify the workspace datacenter$`,
		},
		{
			name: "missing workspace datacenter",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Workspace.Datacenter = ""
				return p
			}(),
			expectedError: `^test-path\.workspace\.datacenter: Required value: must specify the workspace datacenter$`,
		},
		{
			name: "missing default datastore",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Workspace.DefaultDatastore = ""
				return p
			}(),
			expectedError: `^test-path\.workspace\.defaultDatastore: Required value: must specify the default datastore$`,
		},
		{
			name: "missing folder",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.Workspace.Folder = ""
				return p
			}(),
			expectedError: `^test-path\.workspace\.folder: Required value: must specify the VM folder$`,
		},
		{
			name: "missing SCSI controller type",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.SCSIControllerType = ""
				return p
			}(),
			expectedError: `^test-path\.scsiControllerType: Required value: must specify the SCSI controller type$`,
		},
		{
			name: "missing public network",
			platform: func() *vsphere.Platform {
				p := validPlatform()
				p.PublicNetwork = ""
				return p
			}(),
			expectedError: `^test-path\.publicNetwork: Required value: must specify the public VM network$`,
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
