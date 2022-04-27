package manifests

import (
	"testing"

	aiv1beta1 "github.com/openshift/assisted-service/api/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func TestInfraEnv_Generate(t *testing.T) {

	installConfig := &installconfig.InstallConfig{
		Config: &types.InstallConfig{
			ObjectMeta: v1.ObjectMeta{
				Name:      "ocp-edge-cluster-0",
				Namespace: "cluster0",
			},
			PullSecret: "secret-agent",
			SSHKey:     "ssh-key",
		},
	}
	agentPullSecret := &AgentPullSecret{}

	parents := asset.Parents{}
	parents.Add(installConfig, agentPullSecret)

	asset := &InfraEnv{}
	err := asset.Generate(parents)
	assert.NoError(t, err)

	assert.NotEmpty(t, asset.Files())
	infraEnvFile := asset.Files()[0]
	assert.Equal(t, "manifests-ztp/infraenv.yml", infraEnvFile.Filename)

	infraEnv := &aiv1beta1.InfraEnv{}
	err = yaml.Unmarshal(infraEnvFile.Data, &infraEnv)
	assert.NoError(t, err)

	assert.Equal(t, "infraEnv", infraEnv.Name)
	assert.Equal(t, "cluster0", infraEnv.Namespace)
	assert.Equal(t, "ocp-edge-cluster-0", infraEnv.Spec.ClusterRef.Name)
	assert.Equal(t, "cluster0", infraEnv.Spec.ClusterRef.Namespace)
	assert.Equal(t, "ssh-key", infraEnv.Spec.SSHAuthorizedKey)
	assert.Equal(t, "pull-secret", infraEnv.Spec.PullSecretRef.Name)
}
