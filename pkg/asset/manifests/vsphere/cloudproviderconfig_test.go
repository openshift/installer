package vsphere

import (
	"fmt"
	"testing"

	vsphere "github.com/openshift/installer/pkg/types/vsphere"
	"github.com/stretchr/testify/assert"
)

var (
	expectedInTreeConfig = `[Global]
secret-name = "vsphere-creds"
secret-namespace = "kube-system"
insecure-flag = "1"

[Workspace]
server = "test-vcenter"
datacenter = "test-datacenter"
default-datastore = "test-datastore"
folder = "/test-datacenter/vm/clusterID"

[VirtualCenter "test-vcenter"]
datacenters = "test-datacenter"
`

	expectedIniMultiZoneConfig = `[Global]
secret-name = "vsphere-creds"
secret-namespace = "kube-system"
insecure-flag = "1"

[VirtualCenter "test-vcenter"]
port = "443"

datacenters = "test-datacenter,test-datacenter2"

[Workspace]
folder = "/test-datacenter/vm/clusterID"

[Labels]
region = "openshift-region"
zone = "openshift-zone"
`
	expectedYamlMultiZoneConfig = `global:
  user: ""
  password: ""
  server: ""
  port: 0
  insecureFlag: false
  datacenters: []
  soapRoundtripCount: 0
  caFile: ""
  thumbprint: ""
  secretName: vsphere-creds
  secretNamespace: kube-system
  secretsDirectory: ""
  apiDisable: false
  apiBinding: ""
  ipFamily: []
vcenter:
  test-vcenter:
    user: ""
    password: ""
    tenantref: ""
    server: test-vcenter
    port: 443
    insecureFlag: false
    datacenters:
    - test-datacenter
    soapRoundtripCount: 0
    caFile: ""
    thumbprint: ""
    secretref: ""
    secretName: ""
    secretNamespace: ""
    ipFamily: []
labels:
  zone: openshift-zone
  region: openshift-region
`
)

func validInTreePlatform() *vsphere.Platform {
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
				Name:   "test-dz-east-1a",
				Server: "test-vcenter",
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter",
					ComputeCluster: "/test-datacenter/host/cluster",
					Networks: []string{
						"test-network-1",
					},
					Datastore:    "test-datastore",
					ResourcePool: "/test-datacenter/host/cluster/Resources/test-resourcepool",
					Folder:       "/test-datacenter/vm/test-folder",
				},
			},
			{
				Name:   "test-dz-east-2a",
				Server: "test-vcenter",
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter2",
					ComputeCluster: "/test-datacenter2/host/cluster",
					Networks: []string{
						"test-network-1",
					},
					Datastore:    "test-datastore2",
					ResourcePool: "/test-datacenter2/host/cluster/Resources/test-resourcepool",
					Folder:       "/test-datacenter2/vm/test-folder",
				},
			},
			{
				Name:   "test-dz-east-3a",
				Server: "test-vcenter",
				Topology: vsphere.Topology{
					Datacenter:     "test-datacenter2",
					ComputeCluster: "/test-datacenter2/host/cluster",
					Networks: []string{
						"test-network-1",
					},
					Datastore:    "test-datastore2",
					ResourcePool: "/test-datacenter2/host/cluster/Resources/test-resourcepool",
					Folder:       "/test-datacenter2/vm/test-folder",
				},
			},
		},
	}
}

func TestCloudProviderConfig(t *testing.T) {
	folderPath := fmt.Sprintf("/%s/vm/%s", "test-datacenter", "clusterID")
	cases := []struct {
		name                string
		platform            *vsphere.Platform
		expectedCloudConfig string
		outOfTree           bool
	}{
		{
			name:                "valid in-tree cloud provider config",
			platform:            validInTreePlatform(),
			expectedCloudConfig: expectedInTreeConfig,
			outOfTree:           false,
		},
		{
			name:                "valid out of tree cloud provider config",
			platform:            validMultiVCenterPlatform(),
			expectedCloudConfig: expectedIniMultiZoneConfig,
			outOfTree:           true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var cloudConfig string
			var err error
			if tc.outOfTree {
				cloudConfig, err = MultiZoneIniCloudProviderConfig(folderPath, tc.platform)
				fmt.Println(cloudConfig)
			} else {
				cloudConfig, err = InTreeCloudProviderConfig(folderPath, tc.platform)
			}
			assert.NoError(t, err, "failed to create cloud provider config")
			assert.Equal(t, tc.expectedCloudConfig, cloudConfig, "unexpected cloud provider config")
		})
	}
}
