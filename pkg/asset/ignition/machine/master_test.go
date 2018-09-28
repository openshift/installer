package machine

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
)

// TestMasterGenerate tests generating the master asset.
func TestMasterGenerate(t *testing.T) {
	installConfig := `
metadata:
  name: test-cluster
baseDomain: test-domain
networking:
  ServiceCIDR: 10.0.1.0/24
platform:
  aws:
    region: us-east
machines:
- name: master
  replicas: 3
`
	installConfigAsset := &testAsset{"install-config"}
	rootCAAsset := &testAsset{"rootCA"}
	master := &master{
		installConfig: installConfigAsset,
		rootCA:        rootCAAsset,
	}
	dependencies := map[asset.Asset]*asset.State{
		installConfigAsset: stateWithContentsData(installConfig),
		rootCAAsset:        stateWithContentsData("test-rootCA-priv", "test-rootCA-pub"),
	}
	masterState, err := master.Generate(dependencies)
	assert.NoError(t, err, "unexpected error generating master asset")
	expectedIgnitionConfigNames := []string{
		"master-0.ign",
		"master-1.ign",
		"master-2.ign",
	}
	actualIgnitionConfigNames := make([]string, len(masterState.Contents))
	for i, c := range masterState.Contents {
		actualIgnitionConfigNames[i] = c.Name
	}
	assert.Equal(t, expectedIgnitionConfigNames, actualIgnitionConfigNames, "unexpected names for master ignition configs")
}
