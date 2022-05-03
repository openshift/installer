package manifests

import (
	"encoding/base64"
	"testing"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func TestAgentPullSecret_Generate(t *testing.T) {

	installconfigAsset := &installconfig.InstallConfig{
		Config: &types.InstallConfig{
			ObjectMeta: v1.ObjectMeta{
				Namespace: "cluster0",
			},
			PullSecret: "secret-agent",
		},
	}

	parents := asset.Parents{}
	parents.Add(installconfigAsset)

	asset := &AgentPullSecret{}
	err := asset.Generate(parents)
	assert.NoError(t, err)

	assert.NotEmpty(t, asset.Files())
	pullSecretFile := asset.Files()[0]

	assert.Equal(t, "manifests-ztp/pull-secret.yaml", pullSecretFile.Filename)
	secret := corev1.Secret{}
	err = yaml.Unmarshal(pullSecretFile.Data, &secret)
	assert.NoError(t, err)

	data, err := base64.StdEncoding.DecodeString(secret.StringData[".dockerconfigjson"])
	assert.NoError(t, err)
	assert.Equal(t, installconfigAsset.Config.PullSecret, string(data))
}
