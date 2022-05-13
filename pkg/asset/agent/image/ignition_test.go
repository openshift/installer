package image

import (
	"testing"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/stretchr/testify/assert"
)

func TestIgnition_Generate(t *testing.T) {

	// agentClusterInstall := &hiveext.AgentClusterInstall{
	// 	ObjectMeta: v1.ObjectMeta{
	// 		Name:      "test-agent-cluster-install",
	// 		Namespace: "cluster0",
	// 	},
	// 	Spec: hiveext.AgentClusterInstallSpec{
	// 		SSHPublicKey: "ssh-rsa AAAAmyKey",
	// 	},
	// }

	// agentClusterInstallData, err := yaml.Marshal(agentClusterInstall)
	// assert.NoError(t, err)

	// aciConfig := &manifests.AgentClusterInstall{
	// 	File: &asset.File{
	// 		Filename: "/etc/assisted/manifests/agent-cluster-install",
	// 		Data:     agentClusterInstallData,
	// 	},
	// 	Config: agentClusterInstall,
	// }

	agentManifests := &manifests.AgentManifests{
		FileList: []*asset.File{
			{
				Filename: "file1",
				Data:     []byte("file1-content"),
			},
			{
				Filename: "file2",
				Data:     []byte("file2-content"),
			},
			// {
			// 	Filename: "/etc/assisted/manifests/agent-cluster-install",
			// 	Data:     agentClusterInstallData,
			// },
		},
	}

	parents := asset.Parents{}
	parents.Add(agentManifests)

	asset := &Ignition{}
	err := asset.Generate(parents)
	assert.NoError(t, err)

	ignition := asset.Config

	// assert.Equal(t, 1, len(ignition.Passwd.Users))
	// assert.Equal(t, "core", ignition.Passwd.Users[0].Name)
	// assert.Equal(t, agentClusterInstall.Spec.SSHPublicKey, string(ignition.Passwd.Users[0].SSHAuthorizedKeys[0]))

	assert.Equal(t, 2, len(ignition.Storage.Files))
}
