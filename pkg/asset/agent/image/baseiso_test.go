package image

import (
	"fmt"
	"testing"

	"github.com/openshift/installer/pkg/asset"
	"github.com/stretchr/testify/assert"
)

func TestInfraBaseIso_Generate(t *testing.T) {

	GetIsoPluggable = func() (string, error) {
		return "some-openshift-release.iso", nil
	}

	parents := asset.Parents{}

	asset := &BaseIso{}
	err := asset.Generate(parents)
	assert.NoError(t, err)

	assert.NotEmpty(t, asset.Files())
	baseIso := asset.Files()[0]
	assert.Equal(t, baseIso.Filename, "some-openshift-release.iso")

	GetIsoPluggable = func() (string, error) {
		return "", fmt.Errorf("no iso found")
	}
	asset = &BaseIso{}
	err = asset.Generate(parents)
	assert.Error(t, err, "no iso found")
}
