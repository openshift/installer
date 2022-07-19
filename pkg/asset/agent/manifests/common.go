package manifests

import (
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/releaseimage"
	"github.com/openshift/installer/pkg/types"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

var TestSSHKey = `|
	ssh-rsa AAAAB3NzaC1y1LJe3zew1ghc= root@localhost.localdomain`

// GetValidAgentPullSecret returns a valid agent pull secret
func GetValidAgentPullSecret() *AgentPullSecret {
	return &AgentPullSecret{
		Config: &corev1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pull-secret",
				Namespace: "cluster-0",
			},
			StringData: map[string]string{
				".dockerconfigjson": "c2VjcmV0LWFnZW50",
			},
		},
	}
}

// GetValidOptionalInstallConfig returns a valid optional install config
func GetValidOptionalInstallConfig() *agent.OptionalInstallConfig {
	return &agent.OptionalInstallConfig{
		InstallConfig: installconfig.InstallConfig{
			Config: &types.InstallConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ocp-edge-cluster-0",
					Namespace: "cluster-0",
				},
				BaseDomain: "testing.com",
				PullSecret: "secret-agent",
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
			},
		},
		Supplied: true,
	}
}

// GetValidReleaseimage returns a valid release image
func GetValidReleaseimage() *releaseimage.Image {
	return &releaseimage.Image{
		PullSpec: "registry.ci.openshift.org/origin/release:4.11",
	}
}
