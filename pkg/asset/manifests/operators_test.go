package manifests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
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
				MachineCIDR: ipnet.MustParseCIDR("1.2.3.4/5"),
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
				Name:     "control-plane",
				Replicas: pointer.Int64Ptr(3),
			},
			Compute: []types.MachinePool{
				{
					Name:     "compute",
					Replicas: pointer.Int64Ptr(3),
				},
			},
			Platform: types.Platform{
				VSphere: &vspheretypes.Platform{
					VirtualCenters: []vspheretypes.VirtualCenter{
						{
							Name:        "test-server-1",
							Username:    "test-user-1",
							Password:    "test-pass-1",
							Datacenters: []string{"test-datacenter-1"},
						},
						{
							Name:        "test-server-2",
							Username:    "test-user-2",
							Password:    "test-pass-2",
							Datacenters: []string{"test-datacenter-2"},
						},
					},
					Workspace: vspheretypes.Workspace{
						Server:           "test-server-1",
						Datacenter:       "test-datacenter-1",
						DefaultDatastore: "test-datastore",
						ResourcePoolPath: "test-resource-pool",
						Folder:           "test-folder",
					},
					SCSIControllerType: "test-controller-type",
					PublicNetwork:      "test-network",
				},
			},
			PullSecret: "test-pull-secret",
		}
	}
	expectedConfig := createInstallConfig()
	expectedYaml := `baseDomain: test-domain
compute:
- name: compute
  platform: {}
  replicas: 3
controlPlane:
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
  machineCIDR: 1.2.3.4/5
  networkType: test-network-type
  serviceNetwork:
  - 1.2.3.4/5
platform:
  vsphere:
    publicNetwork: test-network
    scsiControllerType: test-controller-type
    virtualCenters:
    - datacenters:
      - test-datacenter-1
      name: test-server-1
      password: ""
      username: ""
    - datacenters:
      - test-datacenter-2
      name: test-server-2
      password: ""
      username: ""
    workspace:
      datacenter: test-datacenter-1
      defaultDatastore: test-datastore
      folder: test-folder
      resourcePoolPath: test-resource-pool
      server: test-server-1
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
