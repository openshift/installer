package manifests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

func TestRedactedInstallConfigRemovesComponentCredentials(t *testing.T) {
	config := types.InstallConfig{
		Platform: types.Platform{VSphere: &vspheretypes.Platform{
			CredentialType: vspheretypes.CredentialTypeComponentScoped,
			VCenters: []vspheretypes.VCenter{{
				Server:   "vcenter.example.com",
				Username: "must-not-leak",
				Password: "must-not-leak-either",
				ComponentCredentials: &vspheretypes.ComponentCredentials{
					MachineManagement:      vspheretypes.Credential{User: "machine-user", Password: "machine-pass"},
					Storage:                vspheretypes.Credential{User: "csi-user", Password: "csi-pass"},
					CloudControllerManager: vspheretypes.Credential{User: "ccm-user", Password: "ccm-pass"},
					VSphereProblemDetector: vspheretypes.Credential{User: "diag-user", Password: "diag-pass"},
				},
			}},
		}},
	}

	redacted, err := redactedInstallConfig(config)
	if !assert.NoError(t, err) {
		return
	}
	for _, secret := range []string{
		"must-not-leak", "must-not-leak-either", "machine-user", "machine-pass",
		"csi-user", "csi-pass", "ccm-user", "ccm-pass", "diag-user", "diag-pass",
	} {
		assert.NotContains(t, string(redacted), secret)
	}
	assert.Equal(t, "machine-pass", config.Platform.VSphere.VCenters[0].ComponentCredentials.MachineManagement.Password)
}
