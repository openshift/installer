package machine

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
)

// TestWorkerGenerate tests generating the worker asset.
func TestWorkerGenerate(t *testing.T) {
	installConfig := &installconfig.InstallConfig{
		Config: &types.InstallConfig{
			Networking: types.Networking{
				ServiceCIDR: ipnet.IPNet{
					IPNet: func(s string) net.IPNet {
						_, cidr, _ := net.ParseCIDR(s)
						return *cidr
					}("10.0.1.0/24"),
				},
			},
			Platform: types.Platform{
				AWS: &types.AWSPlatform{
					Region:       "us-east",
					InstanceType: "m4.large",
				},
			},
		},
	}

	rootCA := &tls.RootCA{}
	err := rootCA.Generate(nil)
	assert.NoError(t, err, "unexpected error generating root CA")

	parents := asset.Parents{}
	parents.Add(installConfig, rootCA)

	worker := &Worker{}
	err = worker.Generate(parents)
	assert.NoError(t, err, "unexpected error generating worker asset")

	actualFiles := worker.Files()
	assert.Equal(t, 1, len(actualFiles), "unexpected number of files in worker state")
	assert.Equal(t, "worker.ign", actualFiles[0].Filename, "unexpected name for worker ignition config")
}
