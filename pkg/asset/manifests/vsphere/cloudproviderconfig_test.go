package vsphere

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	vsphere "github.com/openshift/installer/pkg/types/vsphere"
)

var (
	expectedIniConfig = `[Global]
secret-name = "vsphere-creds"
secret-namespace = "kube-system"
insecure-flag = "1"

[VirtualCenter "test-vcenter"]
port = "443"

datacenters = "test-datacenter,test-datacenter2"

[Workspace]
server = "test-vcenter"
datacenter = "test-datacenter"
default-datastore = "test-datastore"
folder = "/test-datacenter/vm/test-folder"
resourcepool-path = "/test-datacenter/host/cluster/Resources/test-resourcepool"

`

	expectIniLabelsSection = `[Labels]
region = "openshift-region"
zone = "openshift-zone"
`

	expectedYamlConfig = `global:
  insecureFlag: true
  secretName: vsphere-creds
  secretNamespace: kube-system
labels:
  region: openshift-region
  zone: openshift-zone
vcenter:
  test-vcenter:
    datacenters:
    - test-datacenter
    - test-datacenter2
    port: 443
    server: test-vcenter
`
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
					"test-datacenter2",
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
	cases := []struct {
		name                string
		platform            *vsphere.Platform
		cloudProviderFunc   func(string, *vsphere.Platform) (string, error)
		expectedCloudConfig string
	}{
		{
			name:                "valid intree cloud provider config",
			platform:            validPlatform(),
			cloudProviderFunc:   CloudProviderConfigIni,
			expectedCloudConfig: expectedIniConfig + expectIniLabelsSection,
		},
		{
			name: "valid single failure domain intree cloud provider config",
			platform: func() *vsphere.Platform {
				p := validPlatform()

				p.FailureDomains = p.FailureDomains[0:1]
				p.VCenters[0].Datacenters = p.VCenters[0].Datacenters[0:1]

				return p
			}(),
			cloudProviderFunc: CloudProviderConfigIni,
			expectedCloudConfig: func() string {
				// only a single datacenter would be provided to the datacenters
				ini := strings.ReplaceAll(expectedIniConfig, ",test-datacenter2", "")
				return ini
			}(),
		},
		{
			name:                "valid out of tree yaml cloud provider config",
			platform:            validPlatform(),
			cloudProviderFunc:   CloudProviderConfigYaml,
			expectedCloudConfig: expectedYamlConfig,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var cloudConfig string
			var err error
			cloudConfig, err = tc.cloudProviderFunc("infraID", tc.platform)
			assert.NoError(t, err, "failed to create cloud provider config")
			assert.Equal(t, tc.expectedCloudConfig, cloudConfig, "unexpected cloud provider config")
		})
	}
}
