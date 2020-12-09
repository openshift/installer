package machine

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/tls"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// TestMasterIgnitionCustomizationsGenerate tests generating the master ignition check asset.
func TestMasterIgnitionCustomizationsGenerate(t *testing.T) {
	cases := []struct {
		name          string
		customize     bool
		assetExpected bool
	}{
		{
			name:          "not customized",
			customize:     false,
			assetExpected: false,
		},
		{
			name:          "pointer customized",
			customize:     true,
			assetExpected: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			installConfig := &installconfig.InstallConfig{
				Config: &types.InstallConfig{
					Networking: &types.Networking{
						ServiceNetwork: []ipnet.IPNet{*ipnet.MustParseCIDR("10.0.1.0/24")},
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

			master := &Master{}
			err = master.Generate(parents)
			assert.NoError(t, err, "unexpected error generating master asset")

			if tc.customize == true {
				// Modify the master config, emulating a customization to the pointer
				master.Config.Storage.Files = append(master.Config.Storage.Files,
					ignition.FileFromString("/etc/foo", "root", 0644, "foo"))
			}

			parents.Add(master)
			masterIgnCheck := &MasterIgnitionCustomizations{}
			err = masterIgnCheck.Generate(parents)
			assert.NoError(t, err, "unexpected error generating master ignition check asset")

			actualFiles := masterIgnCheck.Files()
			if tc.assetExpected == true {
				assert.Equal(t, 1, len(actualFiles), "unexpected number of files in master state")
				assert.Equal(t, masterMachineConfigFileName, actualFiles[0].Filename, "unexpected name for master ignition config")
			} else {
				assert.Equal(t, 0, len(actualFiles), "unexpected number of files in master state")
			}
		})
	}
}
