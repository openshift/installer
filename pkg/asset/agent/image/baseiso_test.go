package image

import (
	"testing"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/agent"
	"github.com/openshift/installer/pkg/asset/agent/manifests"
	"github.com/stretchr/testify/assert"
)

func TestInfraBaseIso_Generate(t *testing.T) {

	GetIsoPluggable = func() (string, error) {
		return "some-openshift-release.iso", nil
	}

	parents := asset.Parents{}
	manifests := &manifests.AgentManifests{}
	installConfig := &agent.OptionalInstallConfig{}
	parents.Add(manifests, installConfig)

	asset := &BaseIso{}
	err := asset.Generate(parents)
	assert.NoError(t, err)

	assert.NotEmpty(t, asset.Files())
	baseIso := asset.Files()[0]
	assert.Equal(t, baseIso.Filename, "some-openshift-release.iso")
}
