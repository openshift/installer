package machine

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
)

// TestWorkerGenerate tests generating the worker asset.
func TestWorkerGenerate(t *testing.T) {
	installConfig := `
networking:
  ServiceCIDR: 10.0.1.0/24
platform:
  aws:
region: us-east
`
	installConfigAsset := &testAsset{"install-config"}
	rootCAAsset := &testAsset{"rootCA"}
	worker := &worker{
		installConfig: installConfigAsset,
		rootCA:        rootCAAsset,
	}
	dependencies := map[asset.Asset]*asset.State{
		installConfigAsset: stateWithContentsData(installConfig),
		rootCAAsset:        stateWithContentsData("test-rootCA-priv", "test-rootCA-pub"),
	}
	workerState, err := worker.Generate(dependencies)
	assert.NoError(t, err, "unexpected error generating worker asset")
	assert.Equal(t, 1, len(workerState.Contents), "unexpected number of contents in worker state")
	assert.Equal(t, "worker.ign", workerState.Contents[0].Name, "unexpected name for worker ignition config")
}
