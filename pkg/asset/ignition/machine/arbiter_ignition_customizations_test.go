package machine

import (
	"context"
	"testing"

	igntypes "github.com/coreos/ignition/v2/config/v3_2/types"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// TestArbiterIgnitionCustomizationsGenerate tests generating the arbiter ignition check asset.
func TestArbiterIgnitionCustomizationsGenerate(t *testing.T) {
	cases := []struct {
		name          string
		customize     bool
		assetExpected bool
		installConfig *installconfig.InstallConfig
	}{
		{
			name:          "not customized",
			customize:     false,
			assetExpected: false,
			installConfig: installconfig.MakeAsset(
				&types.InstallConfig{
					Networking: &types.Networking{
						ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
					},
					Platform: types.Platform{
						AWS: &aws.Platform{
							Region: "us-east",
						},
					},
					Arbiter: &types.MachinePool{
						Name:     "arbiter",
						Replicas: ptr.To(int64(1)),
					},
				}),
		},
		{
			name:          "pointer customized",
			customize:     true,
			assetExpected: true,
			installConfig: installconfig.MakeAsset(
				&types.InstallConfig{
					Networking: &types.Networking{
						ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
					},
					Platform: types.Platform{
						AWS: &aws.Platform{
							Region: "us-east",
						},
					},
					Arbiter: &types.MachinePool{
						Name:     "arbiter",
						Replicas: ptr.To(int64(1)),
					},
				}),
		},
		{
			name:          "pointer customized but arbiter not set",
			customize:     true,
			assetExpected: false,
			installConfig: installconfig.MakeAsset(
				&types.InstallConfig{
					Networking: &types.Networking{
						ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
					},
					Platform: types.Platform{
						AWS: &aws.Platform{
							Region: "us-east",
						},
					},
				}),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rootCA := &tls.RootCA{}
			err := rootCA.Generate(context.Background(), nil)
			assert.NoError(t, err, "unexpected error generating root CA")

			parents := asset.Parents{}
			parents.Add(tc.installConfig, rootCA)

			arbiter := &Arbiter{}
			err = arbiter.Generate(context.Background(), parents)
			assert.NoError(t, err, "unexpected error generating arbiter asset")

			if tc.customize {
				// Create empty config so that we force the pointer check to validate
				// arbiter skip is also happening when customizations are triggered
				// on non arbiter cluster.
				if arbiter.Config == nil {
					arbiter.Config = &igntypes.Config{}
				}
				// Modify the arbiter config, emulating a customization to the pointer.
				arbiter.Config.Storage.Files = append(arbiter.Config.Storage.Files,
					ignition.FileFromString("/etc/foo", "root", 0644, "foo"))
			}

			parents.Add(arbiter)
			arbiterIgnCheck := &ArbiterIgnitionCustomizations{}
			err = arbiterIgnCheck.Generate(context.Background(), parents)
			assert.NoError(t, err, "unexpected error generating arbiter ignition check asset")

			actualFiles := arbiterIgnCheck.Files()
			if tc.assetExpected {
				assert.Equal(t, 1, len(actualFiles), "unexpected number of files in arbiter state")
				assert.Equal(t, arbiterMachineConfigFileName, actualFiles[0].Filename, "unexpected name for arbiter ignition config")
			} else {
				assert.Equal(t, 0, len(actualFiles), "unexpected number of files in arbiter state")
			}
		})
	}
}
