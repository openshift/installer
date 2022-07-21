package manifests

import (
	"net"

	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/baremetal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

var (
	// TestSSHKey provides a ssh key for unit tests
	TestSSHKey = `|
	ssh-rsa AAAAB3NzaC1y1LJe3zew1ghc= root@localhost.localdomain`
	// TestSecret provides a ssh key for unit tests
	TestSecret = `'{"auths":{"cloud.openshift.com":{"auth":"b3BlUTA=","email":"test@redhat.com"}}}`
	// TestReleaseImage provides a release image url for unit tests
	TestReleaseImage = "registry.ci.openshift.org/origin/release:4.11"
)

// GetValidOptionalInstallConfig returns a valid optional install config
func getValidOptionalInstallConfig() *agent.OptionalInstallConfig {
	_, newCidr, _ := net.ParseCIDR("192.168.111.0/24")

	return &agent.OptionalInstallConfig{
		InstallConfig: installconfig.InstallConfig{
			Config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ocp-edge-cluster-0",
					Namespace: "cluster-0",
				},
				BaseDomain: "testing.com",
				PullSecret: TestSecret,
				SSHKey:     TestSSHKey,
				ControlPlane: &types.MachinePool{
					Name:     "master",
					Replicas: pointer.Int64Ptr(3),
					Platform: types.MachinePoolPlatform{},
				},
				Compute: []types.MachinePool{
					{
						Name:     "worker-machine-pool-1",
						Replicas: pointer.Int64Ptr(2),
					},
					{
						Name:     "worker-machine-pool-2",
						Replicas: pointer.Int64Ptr(3),
					},
				},
				Networking: &types.Networking{
					ClusterNetwork: []types.ClusterNetworkEntry{
						{
							CIDR:       ipnet.IPNet{IPNet: *newCidr},
							HostPrefix: 23,
						},
					},
					ServiceNetwork: []ipnet.IPNet{
						*ipnet.MustParseCIDR("172.30.0.0/16"),
					},
				},
				Platform: types.Platform{
					BareMetal: &baremetal.Platform{
						APIVIP:     "192.168.122.10",
						IngressVIP: "192.168.122.11",
					},
				},
			},
		},
		Supplied: true,
	}
}
