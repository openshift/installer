package machine

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// TestComputeGenerate tests generating the compute asset.
func TestComputeGenerate(t *testing.T) {
	installConfig := &installconfig.InstallConfig{
		Config: &types.InstallConfig{
			Networking: types.Networking{
				ServiceCIDR: *ipnet.MustParseCIDR("10.0.1.0/24"),
			},
			Platform: types.Platform{
				AWS: &aws.Platform{
					Region: "us-east",
				},
			},
		},
	}

	rootCA := &tls.RootCA{}
	err := rootCA.Generate(nil)
	assert.NoError(t, err, "unexpected error generating root CA")

	parents := asset.Parents{}
	parents.Add(installConfig, rootCA)

	compute := &Compute{}
	err = compute.Generate(parents)
	assert.NoError(t, err, "unexpected error generating compute asset")

	actualFiles := compute.Files()
	assert.Equal(t, 1, len(actualFiles), "unexpected number of files in compute state")
	assert.Equal(t, "compute.ign", actualFiles[0].Filename, "unexpected name for compute ignition config")
}
