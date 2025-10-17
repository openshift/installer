package manifests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	nutanixtypes "github.com/openshift/installer/pkg/types/nutanix"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

// TestRedactedInstallConfig tests the redactedInstallConfig function.
func TestRedactedInstallConfig(t *testing.T) {
	createInstallConfig := func() *types.InstallConfig {
		return &types.InstallConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster",
			},
			SSHKey:     "test-ssh-key",
			BaseDomain: "test-domain",
			Networking: &types.Networking{
				MachineNetwork: []types.MachineNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("1.2.3.4/5")},
				},
				NetworkType: "test-network-type",
				ClusterNetwork: []types.ClusterNetworkEntry{
					{
						CIDR:       *ipnet.MustParseCIDR("1.2.3.4/5"),
						HostPrefix: 6,
					},
				},
				ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("1.2.3.4/5")},
			},
			ControlPlane: &types.MachinePool{
				Name:         "control-plane",
				Replicas:     ptr.To(int64(3)),
				Architecture: types.ArchitectureAMD64,
			},
			Compute: []types.MachinePool{
				{
					Name:         "compute",
					Replicas:     ptr.To(int64(3)),
					Architecture: types.ArchitectureAMD64,
				},
			},
			Platform: types.Platform{
				VSphere: &vspheretypes.Platform{
					DeprecatedUsername: "test-username",
					DeprecatedPassword: "test-password",
					VCenters: []vspheretypes.VCenter{{
						Server:      "test-server-1",
						Port:        443,
						Username:    "",
						Password:    "",
						Datacenters: []string{"test-datacenter"},
					}},
					FailureDomains: []vspheretypes.FailureDomain{{
						Name:   "test-failuredomain",
						Region: "test-region",
						Zone:   "test-zone",
						Server: "test-server-1",
						Topology: vspheretypes.Topology{
							Datacenter:     "test-datacenter",
							ComputeCluster: "test-computecluster",
							Networks:       []string{"test-network"},
							Datastore:      "test-datastore",
							ResourcePool:   "test-resourcepool",
							Folder:         "test-folder",
						},
					}},
				},
			},
			PullSecret: "test-pull-secret",
		}
	}
	expectedConfig := createInstallConfig()
	expectedYaml := `baseDomain: test-domain
compute:
- architecture: amd64
  name: compute
  platform: {}
  replicas: 3
controlPlane:
  architecture: amd64
  name: control-plane
  platform: {}
  replicas: 3
metadata:
  creationTimestamp: null
  name: test-cluster
networking:
  clusterNetwork:
  - cidr: 1.2.3.4/5
    hostPrefix: 6
  machineNetwork:
  - cidr: 1.2.3.4/5
  networkType: test-network-type
  serviceNetwork:
  - 1.2.3.4/5
platform:
  vsphere:
    failureDomains:
    - name: test-failuredomain
      region: test-region
      server: test-server-1
      topology:
        computeCluster: test-computecluster
        datacenter: test-datacenter
        datastore: test-datastore
        folder: test-folder
        networks:
        - test-network
        resourcePool: test-resourcepool
      zone: test-zone
    vcenters:
    - datacenters:
      - test-datacenter
      password: ""
      server: test-server-1
      user: ""
pullSecret: ""
sshKey: test-ssh-key
`
	ic := createInstallConfig()
	actualYaml, err := redactedInstallConfig(*ic)
	if assert.NoError(t, err, "unexpected error") {
		assert.Equal(t, expectedYaml, string(actualYaml), "unexpected yaml")
	}
	assert.Equal(t, expectedConfig, ic, "install config was unexpectedly modified")
}

// TestRedactedInstallConfigNutanix tests the redactedInstallConfig function for Nutanix platform.
func TestRedactedInstallConfigNutanix(t *testing.T) {
	createInstallConfigWithNutanix := func() *types.InstallConfig {
		return &types.InstallConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-nutanix-cluster",
			},
			SSHKey:     "test-ssh-key",
			BaseDomain: "test-domain",
			Networking: &types.Networking{
				MachineNetwork: []types.MachineNetworkEntry{
					{CIDR: *ipnet.MustParseCIDR("1.2.3.4/5")},
				},
				NetworkType: "test-network-type",
				ClusterNetwork: []types.ClusterNetworkEntry{
					{
						CIDR:       *ipnet.MustParseCIDR("1.2.3.4/5"),
						HostPrefix: 6,
					},
				},
				ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("1.2.3.4/5")},
			},
			ControlPlane: &types.MachinePool{
				Name:         "control-plane",
				Replicas:     ptr.To(int64(3)),
				Architecture: types.ArchitectureAMD64,
			},
			Compute: []types.MachinePool{
				{
					Name:         "compute",
					Replicas:     ptr.To(int64(3)),
					Architecture: types.ArchitectureAMD64,
				},
			},
			Platform: types.Platform{
				Nutanix: &nutanixtypes.Platform{
					PrismCentral: nutanixtypes.PrismCentral{
						Endpoint: nutanixtypes.PrismEndpoint{
							Address: "test-prism-central.test.com",
							Port:    9440,
						},
						Username: "test-username",
						Password: "test-password",
					},
					PrismElements: []nutanixtypes.PrismElement{{
						UUID: "test-uuid",
						Endpoint: nutanixtypes.PrismEndpoint{
							Address: "test-prism-element.test.com",
							Port:    9440,
						},
						Name: "test-element",
					}},
					SubnetUUIDs: []string{"test-subnet-uuid"},
				},
			},
			PullSecret: "test-pull-secret",
		}
	}

	expectedYaml := `baseDomain: test-domain
compute:
- architecture: amd64
  name: compute
  platform: {}
  replicas: 3
controlPlane:
  architecture: amd64
  name: control-plane
  platform: {}
  replicas: 3
metadata:
  creationTimestamp: null
  name: test-nutanix-cluster
networking:
  clusterNetwork:
  - cidr: 1.2.3.4/5
    hostPrefix: 6
  machineNetwork:
  - cidr: 1.2.3.4/5
  networkType: test-network-type
  serviceNetwork:
  - 1.2.3.4/5
platform:
  nutanix:
    prismCentral:
      endpoint:
        address: test-prism-central.test.com
        port: 9440
      password: ""
      username: ""
    prismElements:
    - endpoint:
        address: test-prism-element.test.com
        port: 9440
      name: test-element
      uuid: test-uuid
    subnetUUIDs:
    - test-subnet-uuid
pullSecret: ""
sshKey: test-ssh-key
`

	expectedConfig := createInstallConfigWithNutanix()
	ic := createInstallConfigWithNutanix()
	actualYaml, err := redactedInstallConfig(*ic)
	if assert.NoError(t, err, "unexpected error") {
		assert.Equal(t, expectedYaml, string(actualYaml), "unexpected yaml")
	}
	assert.Equal(t, expectedConfig, ic, "install config was unexpectedly modified")
}
