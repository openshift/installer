package machine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// TestWorkerGenerate tests generating the worker asset.
func TestWorkerGenerate(t *testing.T) {
	installConfig := installconfig.MakeAsset(
		&types.InstallConfig{
			Networking: &types.Networking{
				ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
			},
			Platform: types.Platform{
				AWS: &aws.Platform{
					Region: "us-east",
				},
			},
		})

	rootCA := &tls.RootCA{}
	err := rootCA.Generate(context.Background(), nil)
	assert.NoError(t, err, "unexpected error generating root CA")

	parents := asset.Parents{}
	parents.Add(installConfig, rootCA)

	worker := &Worker{}
	err = worker.Generate(context.Background(), parents)
	assert.NoError(t, err, "unexpected error generating worker asset")

	actualFiles := worker.Files()
	assert.Equal(t, 1, len(actualFiles), "unexpected number of files in worker state")
	assert.Equal(t, "worker.ign", actualFiles[0].Filename, "unexpected name for worker ignition config")
}
