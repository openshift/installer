package vsphere

import (
	"testing"

	"github.com/openshift/installer/pkg/types"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

func TestMetadataReturnsCredentialLookupError(t *testing.T) {
	config := &types.InstallConfig{
		Platform: types.Platform{VSphere: &vspheretypes.Platform{
			CredentialType: vspheretypes.CredentialTypeComponentScoped,
			VCenters: []vspheretypes.VCenter{{
				Server: "vcenter.example.com",
			}},
		}},
	}

	_, err := Metadata(config)
	if err == nil {
		t.Fatal("expected missing component credentials error")
	}
}
