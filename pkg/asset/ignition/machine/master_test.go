package machine

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
)

// TestMasterGenerate tests generating the master asset.
func TestMasterGenerate(t *testing.T) {
	installConfig := &installconfig.InstallConfig{
		Config: &types.InstallConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster",
			},
			BaseDomain: "test-domain",
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
					Region: "us-east",
				},
			},
			Machines: []types.MachinePool{
				{
					Name:     "master",
					Replicas: func(x int64) *int64 { return &x }(3),
				},
			},
		},
	}

	rootCA := &tls.RootCA{}
	err := rootCA.Generate(nil)
	assert.NoError(t, err, "unexpected error generating root CA")

	parents := asset.Parents{}
	parents.Add(installConfig, rootCA)

	master := &Master{}
	err = master.Generate(parents)
	assert.NoError(t, err, "unexpected error generating master asset")
	expectedIgnitionConfigNames := []string{
		"master-0.ign",
		"master-1.ign",
		"master-2.ign",
	}
	actualFiles := master.Files()
	actualIgnitionConfigNames := make([]string, len(actualFiles))
	for i, f := range actualFiles {
		actualIgnitionConfigNames[i] = f.Filename
	}
	assert.Equal(t, expectedIgnitionConfigNames, actualIgnitionConfigNames, "unexpected names for master ignition configs")
}
