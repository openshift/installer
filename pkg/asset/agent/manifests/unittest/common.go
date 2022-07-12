package unittest

import (
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// GetEmptyTestAssets returns an empty asset without an install config
func GetEmptyTestAssets() []asset.Asset {
	return []asset.Asset{
		&agent.OptionalInstallConfig{},
	}
}

// GetTestAssetsWithValidInstallConfig returns a valid install config
func GetTestAssetsWithValidInstallConfig() []asset.Asset {
	return []asset.Asset{
		&agent.OptionalInstallConfig{
			InstallConfig: installconfig.InstallConfig{
				Config: &types.InstallConfig{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "ocp-edge-cluster-0",
						Namespace: "cluster-0",
					},
					BaseDomain: "testing.com",
					SSHKey:     "ssh-key",
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
				},
			},
			Supplied: true,
		},
	}
}
