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

func TestMetadataPreservesPort(t *testing.T) {
	config := &types.InstallConfig{
		Platform: types.Platform{VSphere: &vspheretypes.Platform{
			VCenters: []vspheretypes.VCenter{{
				Server:   "vcenter.example.com",
				Port:     8443,
				Username: "user",
				Password: "pass",
			}},
		}},
	}

	meta, err := Metadata(config)
	if err != nil {
		t.Fatalf("Metadata() returned unexpected error: %v", err)
	}
	if len(meta.VCenters) != 1 {
		t.Fatalf("len(meta.VCenters) = %d, want 1", len(meta.VCenters))
	}
	if meta.VCenters[0].VCenter != "vcenter.example.com" {
		t.Errorf("VCenter = %q, want %q", meta.VCenters[0].VCenter, "vcenter.example.com")
	}
	if meta.VCenters[0].Port != 8443 {
		t.Errorf("Port = %d, want %d", meta.VCenters[0].Port, 8443)
	}
}
