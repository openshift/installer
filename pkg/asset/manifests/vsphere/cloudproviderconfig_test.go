package vsphere

import (
	"testing"

	"github.com/stretchr/testify/assert"

	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

func TestCloudProviderConfig(t *testing.T) {
	cases := []struct {
		name           string
		platform       *vspheretypes.Platform
		expectedConfig string
	}{
		{
			name: "single virtualcenter",
			platform: &vspheretypes.Platform{
				VirtualCenters: []vspheretypes.VirtualCenter{
					{
						Name:        "test-name",
						Username:    "test-username",
						Password:    "test-password",
						Datacenters: []string{"test-datacenter"},
					},
				},
				Workspace: vspheretypes.Workspace{
					Server:           "test-name",
					Datacenter:       "test-datacenter",
					DefaultDatastore: "test-datastore",
					ResourcePoolPath: "test-resource-pool",
					Folder:           "test-folder",
				},
				SCSIControllerType: "test-scsi",
				PublicNetwork:      "test-network",
			},
			expectedConfig: `[Global]
secret-name      = vsphere-creds
secret-namespace = kube-system

[Workspace]
server            = test-name
datacenter        = test-datacenter
default-datastore = test-datastore
resourcepool-path = test-resource-pool
folder            = test-folder

[Disk]
scsicontrollertype = test-scsi

[Network]
public-network = test-network

[VirtualCenter "test-name"]
datacenters = test-datacenter

`,
		},
		{
			name: "multiple datacenters",
			platform: &vspheretypes.Platform{
				VirtualCenters: []vspheretypes.VirtualCenter{
					{
						Name:        "test-name",
						Username:    "test-username",
						Password:    "test-password",
						Datacenters: []string{"test-dc1", "test-dc2"},
					},
				},
				Workspace: vspheretypes.Workspace{
					Server:           "test-name",
					Datacenter:       "test-dc1",
					DefaultDatastore: "test-datastore",
					ResourcePoolPath: "test-resource-pool",
					Folder:           "test-folder",
				},
				SCSIControllerType: "test-scsi",
				PublicNetwork:      "test-network",
			},
			expectedConfig: `[Global]
secret-name      = vsphere-creds
secret-namespace = kube-system

[Workspace]
server            = test-name
datacenter        = test-dc1
default-datastore = test-datastore
resourcepool-path = test-resource-pool
folder            = test-folder

[Disk]
scsicontrollertype = test-scsi

[Network]
public-network = test-network

[VirtualCenter "test-name"]
datacenters = test-dc1,test-dc2

`,
		},
		{
			name: "multiple virtualcenters",
			platform: &vspheretypes.Platform{
				VirtualCenters: []vspheretypes.VirtualCenter{
					{
						Name:        "test-name1",
						Username:    "test-username1",
						Password:    "test-password1",
						Datacenters: []string{"test-datacenter1"},
					},
					{
						Name:        "test-name2",
						Username:    "test-username2",
						Password:    "test-password2",
						Datacenters: []string{"test-datacenter2"},
					},
				},
				Workspace: vspheretypes.Workspace{
					Server:           "test-name",
					Datacenter:       "test-datacenter",
					DefaultDatastore: "test-datastore",
					ResourcePoolPath: "test-resource-pool",
					Folder:           "test-folder",
				},
				SCSIControllerType: "test-scsi",
				PublicNetwork:      "test-network",
			},
			expectedConfig: `[Global]
secret-name      = vsphere-creds
secret-namespace = kube-system

[Workspace]
server            = test-name
datacenter        = test-datacenter
default-datastore = test-datastore
resourcepool-path = test-resource-pool
folder            = test-folder

[Disk]
scsicontrollertype = test-scsi

[Network]
public-network = test-network

[VirtualCenter "test-name1"]
datacenters = test-datacenter1

[VirtualCenter "test-name2"]
datacenters = test-datacenter2

`,
		},
		{
			name: "empty resource pool path",
			platform: &vspheretypes.Platform{
				VirtualCenters: []vspheretypes.VirtualCenter{
					{
						Name:        "test-name",
						Username:    "test-username",
						Password:    "test-password",
						Datacenters: []string{"test-datacenter"},
					},
				},
				Workspace: vspheretypes.Workspace{
					Server:           "test-name",
					Datacenter:       "test-datacenter",
					DefaultDatastore: "test-datastore",
					Folder:           "test-folder",
				},
				SCSIControllerType: "test-scsi",
				PublicNetwork:      "test-network",
			},
			expectedConfig: `[Global]
secret-name      = vsphere-creds
secret-namespace = kube-system

[Workspace]
server            = test-name
datacenter        = test-datacenter
default-datastore = test-datastore
folder            = test-folder

[Disk]
scsicontrollertype = test-scsi

[Network]
public-network = test-network

[VirtualCenter "test-name"]
datacenters = test-datacenter

`,
		},
		{
			name: "insecure flag",
			platform: &vspheretypes.Platform{
				VirtualCenters: []vspheretypes.VirtualCenter{
					{
						Name:        "test-name",
						Username:    "test-username",
						Password:    "test-password",
						Datacenters: []string{"test-datacenter"},
						Insecure:    true,
					},
				},
				Workspace: vspheretypes.Workspace{
					Server:           "test-name",
					Datacenter:       "test-datacenter",
					DefaultDatastore: "test-datastore",
					ResourcePoolPath: "test-resource-pool",
					Folder:           "test-folder",
				},
				SCSIControllerType: "test-scsi",
				PublicNetwork:      "test-network",
			},
			expectedConfig: `[Global]
secret-name      = vsphere-creds
secret-namespace = kube-system

[Workspace]
server            = test-name
datacenter        = test-datacenter
default-datastore = test-datastore
resourcepool-path = test-resource-pool
folder            = test-folder

[Disk]
scsicontrollertype = test-scsi

[Network]
public-network = test-network

[VirtualCenter "test-name"]
datacenters   = test-datacenter
insecure-flag = 1

`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config, err := CloudProviderConfig(tc.platform)
			assert.NoError(t, err, "failed to create cloud provider config")
			assert.Equal(t, tc.expectedConfig, config, "unexpected cloud provider config")
		})
	}
}
