package vsphere

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

var cloudProviderLegacyPlatform = &vspheretypes.Platform{
	VCenter:          "test-name",
	Username:         "test-username",
	Password:         "test-password",
	Datacenter:       "test-datacenter",
	DefaultDatastore: "test-datastore",
}

var cloudProviderLegacyExpectedConfig = `[Global]
secret-name = "vsphere-creds"
secret-namespace = "kube-system"
insecure-flag = "1"

[Workspace]
server = "test-name"
datacenter = "test-datacenter"
default-datastore = "test-datastore"
folder = "/test-datacenter/vm/clusterID"

[VirtualCenter "test-name"]
datacenters = "test-datacenter"
`

var cloudProviderMultiClusterPlatform = &vspheretypes.Platform{
	VCenters: []vspheretypes.VCenter{
		{
			Server:   "test-vcenter-1",
			User:     "test-user",
			Password: "test-password",
			Port:     443,
			Regions: []vspheretypes.Region{
				{
					Name:       "region-a",
					Datacenter: "dc-region-a",
					Zones: []vspheretypes.Zone{
						{
							Name:         "zone-a-1",
							ResourcePool: "zone-a-pool-1",
							Cluster:      "zone-a-cluster-1",
							Network:      "VM Network",
							Datastore:    "zone-a-datastore-1",
						},
						{
							Name:         "zone-a-2",
							ResourcePool: "zone-a-pool-2",
							Cluster:      "zone-a-cluster-2",
							Network:      "VM Network",
							Datastore:    "zone-a-datastore-2",
						},
					},
				},
				{
					Name:       "region-b",
					Datacenter: "dc-region-b",
					Zones: []vspheretypes.Zone{
						{
							Name:         "zone-b-1",
							ResourcePool: "zone-b-pool-1",
							Cluster:      "zone-b-cluster-1",
							Network:      "VM Network",
							Datastore:    "zone-b-datastore-1",
						},
						{
							Name:         "zone-b-2",
							ResourcePool: "zone-b-pool-2",
							Cluster:      "zone-b-cluster-2",
							Network:      "VM Network",
							Datastore:    "zone-b-datastore-2",
						},
					},
				},
			},
		},
		{
			Server:   "test-vcenter-2",
			User:     "test-user",
			Password: "test-password",
			Port:     443,
			Regions: []vspheretypes.Region{
				{
					Name:       "region-c",
					Datacenter: "dc-region-c",
					Zones: []vspheretypes.Zone{
						{
							Name:         "zone-c-1",
							ResourcePool: "zone-c-pool-1",
							Cluster:      "zone-c-cluster-1",
							Network:      "VM Network",
							Datastore:    "zone-c-datastore-1",
						},
						{
							Name:         "zone-c-2",
							ResourcePool: "zone-c-pool-2",
							Cluster:      "zone-c-cluster-2",
							Network:      "VM Network",
							Datastore:    "zone-c-datastore-2",
						},
					},
				},
				{
					Name:       "region-d",
					Datacenter: "dc-region-d",
					Zones: []vspheretypes.Zone{
						{
							Name:         "zone-d-1",
							ResourcePool: "zone-d-pool-1",
							Cluster:      "zone-d-cluster-1",
							Network:      "VM Network",
							Datastore:    "zone-d-datastore-1",
						},
						{
							Name:         "zone-d-2",
							ResourcePool: "zone-d-pool-2",
							Cluster:      "zone-d-cluster-2",
							Network:      "VM Network",
							Datastore:    "zone-d-datastore-2",
						},
					},
				},
			},
		},
	},
}

var cloudProviderMultiClusteryExpectedConfig = `[Global]
secret-name = "vsphere-creds"
secret-namespace = "kube-system"
insecure-flag = "1"

[Workspace]
folder = "/test-datacenter/vm/clusterID"

[VirtualCenter "test-vcenter-1"]
datacenters = "dc-region-a,dc-region-b"

[VirtualCenter "test-vcenter-2"]
datacenters = "dc-region-c,dc-region-d"

[Labels]
region = "region-a,region-b,region-c,region-d"
zone = "zone-a-1,zone-a-2,zone-b-1,zone-b-2,zone-c-1,zone-c-2,zone-d-1,zone-d-2"
`

func TestCloudProviderConfig(t *testing.T) {
	folderPath := fmt.Sprintf("/%s/vm/%s", "test-datacenter", "clusterID")

	testCases := []struct {
		testCase       string
		platform       *vspheretypes.Platform
		expectedConfig string
	}{
		{
			testCase:       "cloud provider configuration matches legacy platform specification",
			platform:       cloudProviderLegacyPlatform,
			expectedConfig: cloudProviderLegacyExpectedConfig,
		},
		{
			testCase:       "cloud provider configuration matches multi-vcenter platform specification",
			platform:       cloudProviderMultiClusterPlatform,
			expectedConfig: cloudProviderMultiClusteryExpectedConfig,
		},
	}

	for _, tc := range testCases {
		actualConfig, err := CloudProviderConfig(folderPath, tc.platform)
		assert.NoError(t, err, "failed to create cloud provider config")
		assert.Equal(t, tc.expectedConfig, actualConfig, "unexpected cloud provider config")
	}

}
