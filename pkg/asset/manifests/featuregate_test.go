package manifests

import (
	"os"
	"testing"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"

	"github.com/stretchr/testify/assert"
)

var fixtureCustomNoUpgrade = `apiVersion: config.openshift.io/v1
kind: FeatureGate
metadata:
  creationTimestamp: null
  name: cluster
spec:
  customNoUpgrade:
    disabled:
    - ""
    enabled:
    - DisableCloudProviders
    - CSIMigrationAWS
  featureSet: CustomNoUpgrade
status: {}
`

var fixtureIPv6DualStackNoUpgrade = `apiVersion: config.openshift.io/v1
kind: FeatureGate
metadata:
  creationTimestamp: null
  name: cluster
spec:
  featureSet: IPv6DualStackNoUpgrade
status: {}
`

func TestCustomNoUpgradeFeatureGate(t *testing.T) {
	os.Setenv("OPENSHIFT_INSTALL_CUSTOM_FEATURESET_NAME", "CustomNoUpgrade")
	os.Setenv("OPENSHIFT_INSTALL_CUSTOM_FEATUREGATES_ENABLED", "DisableCloudProviders,CSIMigrationAWS")
	defer os.Unsetenv("OPENSHIFT_INSTALL_CUSTOM_FEATURESET_NAME")
	defer os.Unsetenv("OPENSHIFT_INSTALL_CUSTOM_FEATUREGATES_ENABLED")

	fg := FeatureGate{}

	parents := asset.Parents{}
	parents.Add(&installconfig.InstallConfig{})

	err := fg.Generate(parents)
	assert.NoError(t, err)

	files := fg.Files()
	assert.Equal(t, 1, len(files))

	assert.Equal(t, featureGateFilename, files[0].Filename)
	assert.Equal(t, fixtureCustomNoUpgrade, string(files[0].Data))
}

func TestIPv6DualStackNoUpgradeFeatureGate(t *testing.T) {
	os.Setenv("OPENSHIFT_INSTALL_CUSTOM_FEATURESET_NAME", "IPv6DualStackNoUpgrade")
	os.Setenv("OPENSHIFT_INSTALL_CUSTOM_FEATUREGATES_ENABLED", "DisableCloudProviders,CSIMigrationAWS")
	defer os.Unsetenv("OPENSHIFT_INSTALL_CUSTOM_FEATURESET_NAME")
	defer os.Unsetenv("OPENSHIFT_INSTALL_CUSTOM_FEATUREGATES_ENABLED")

	fg := FeatureGate{}

	parents := asset.Parents{}
	parents.Add(&installconfig.InstallConfig{})

	err := fg.Generate(parents)
	assert.NoError(t, err)

	files := fg.Files()
	assert.Equal(t, 1, len(files))

	assert.Equal(t, featureGateFilename, files[0].Filename)
	assert.Equal(t, fixtureIPv6DualStackNoUpgrade, string(files[0].Data))
}

func TestNoFeatureGate(t *testing.T) {
	fg := FeatureGate{}

	parents := asset.Parents{}
	parents.Add(&installconfig.InstallConfig{})

	err := fg.Generate(parents)
	assert.NoError(t, err)

	files := fg.Files()
	assert.Equal(t, 0, len(files))
}
