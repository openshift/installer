package machine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// TestControlPlaneGenerate tests generating the control plane asset.
func TestControlPlaneGenerate(t *testing.T) {
	installConfig := &installconfig.InstallConfig{
		Config: &types.InstallConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-cluster",
			},
			BaseDomain: "test-domain",
			Networking: types.Networking{
				ServiceCIDR: *ipnet.MustParseCIDR("10.0.1.0/24"),
			},
			Platform: types.Platform{
				AWS: &aws.Platform{
					Region: "us-east",
				},
			},
			Machines: []types.MachinePool{
				{
					Name:     "controlplane",
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

	controlPlane := &ControlPlane{}
	err = controlPlane.Generate(parents)
	assert.NoError(t, err, "unexpected error generating control plane asset")
	expectedIgnitionConfigNames := []string{
		"controlplane.ign",
	}
	actualFiles := controlPlane.Files()
	actualIgnitionConfigNames := make([]string, len(actualFiles))
	for i, f := range actualFiles {
		actualIgnitionConfigNames[i] = f.Filename
	}
	assert.Equal(t, expectedIgnitionConfigNames, actualIgnitionConfigNames, "unexpected names for control plane ignition configs")
}
